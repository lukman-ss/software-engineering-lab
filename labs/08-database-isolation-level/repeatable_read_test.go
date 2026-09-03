package isolation_test

import (
	"context"
	"database/sql"
	"sync"
	"testing"
	"time"

	_ "github.com/lib/pq"
	isolation "github.com/lukman-ss/software-engineering-lab/labs/08-database-isolation-level"
)

// TestRepeatableRead_SnapshotIsolation proves that in REPEATABLE READ,
// a transaction sees a consistent snapshot taken at the start of transaction.
// Even when another transaction commits an update in between, the read remains consistent.
func TestRepeatableRead_SnapshotIsolation(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()
	repo := isolation.NewPostgresWalletRepo(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resetTestState(t, ctx, db, repo, 1000000, 1000000, 1000000)

	tx1Ready := make(chan struct{})
	tx2Committed := make(chan struct{})

	var tx1FirstRead, tx1SecondRead int64
	var tx1Err, tx2Err error

	var wg sync.WaitGroup
	wg.Add(2)

	// Transaction 1: REPEATABLE READ reader
	go func() {
		defer wg.Done()
		tx1, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
		if err != nil {
			tx1Err = err
			return
		}
		defer tx1.Rollback()

		// 1st Read (snapshot acquired)
		tx1FirstRead, tx1Err = repo.GetBalance(ctx, tx1, 1)
		if tx1Err != nil {
			return
		}

		close(tx1Ready) // Signal TX2 to modify data

		select {
		case <-tx2Committed: // Wait for TX2 commit
		case <-ctx.Done():
			tx1Err = ctx.Err()
			return
		}

		// 2nd Read (must see snapshot, ignoring TX2's commit)
		tx1SecondRead, tx1Err = repo.GetBalance(ctx, tx1, 1)
		if tx1Err != nil {
			return
		}

		if err := tx1.Commit(); err != nil {
			tx1Err = err
		}
	}()

	// Transaction 2: Writer
	go func() {
		defer wg.Done()
		defer close(tx2Committed)

		select {
		case <-tx1Ready:
		case <-ctx.Done():
			tx2Err = ctx.Err()
			return
		}

		tx2, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
		if err != nil {
			tx2Err = err
			return
		}
		defer tx2.Rollback()

		_, tx2Err = tx2.ExecContext(ctx, "UPDATE isolation_accounts SET balance = balance + 500000 WHERE id = 1")
		if tx2Err != nil {
			return
		}

		tx2Err = tx2.Commit()
	}()

	wg.Wait()

	if tx1Err != nil || tx2Err != nil {
		t.Fatalf("unexpected error: tx1=%v, tx2=%v", tx1Err, tx2Err)
	}

	t.Logf("REPEATABLE READ Snapshot Isolation Proof:")
	t.Logf("  TX1 First Read:  %d", tx1FirstRead)
	t.Logf("  TX1 Second Read: %d", tx1SecondRead)

	if tx1FirstRead != tx1SecondRead {
		t.Fatalf("REPEATABLE READ failed to provide snapshot isolation: 1st=%d != 2nd=%d", tx1FirstRead, tx1SecondRead)
	}
	if tx1FirstRead != 1000000 {
		t.Errorf("expected balance 1000000, got %d", tx1FirstRead)
	}
}

// TestRepeatableRead_ConcurrentUpdate_SerializationFailure proves that PostgreSQL's
// REPEATABLE READ prevents lost updates by raising error 40001 (serialization failure)
// when two transactions try to update the same row concurrently.
func TestRepeatableRead_ConcurrentUpdate_SerializationFailure(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()
	repo := isolation.NewPostgresWalletRepo(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resetTestState(t, ctx, db, repo, 1000000, 1000000, 1000000)

	var wg sync.WaitGroup
	var errA, errB error
	wg.Add(2)

	tx1Read := make(chan struct{})
	tx2Read := make(chan struct{})
	release := make(chan struct{})

	go func() {
		defer wg.Done()
		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
		if err != nil {
			errA = err
			return
		}
		defer tx.Rollback()

		if _, err := repo.GetBalance(ctx, tx, 1); err != nil {
			errA = err
			return
		}
		close(tx1Read)

		select {
		case <-release:
		case <-ctx.Done():
			errA = ctx.Err()
			return
		}

		_, errA = tx.ExecContext(ctx, "UPDATE isolation_accounts SET balance = balance - 800000 WHERE id = 1")
		if errA == nil {
			errA = tx.Commit()
		}
	}()

	go func() {
		defer wg.Done()
		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
		if err != nil {
			errB = err
			return
		}
		defer tx.Rollback()

		if _, err := repo.GetBalance(ctx, tx, 1); err != nil {
			errB = err
			return
		}
		close(tx2Read)

		select {
		case <-release:
		case <-ctx.Done():
			errB = ctx.Err()
			return
		}

		_, errB = tx.ExecContext(ctx, "UPDATE isolation_accounts SET balance = balance - 800000 WHERE id = 1")
		if errB == nil {
			errB = tx.Commit()
		}
	}()

	select {
	case <-tx1Read:
	case <-ctx.Done():
		t.Fatalf("timeout waiting for tx1Read: %v", ctx.Err())
	}

	select {
	case <-tx2Read:
	case <-ctx.Done():
		t.Fatalf("timeout waiting for tx2Read: %v", ctx.Err())
	}

	close(release)
	wg.Wait()

	t.Logf("Repeatable Read Concurrent Update Error: errA=%v, errB=%v", errA, errB)

	successCount := 0
	serializationFailureCount := 0

	for _, err := range []error{errA, errB} {
		if err == nil {
			successCount++
		} else if isolation.IsSerializationError(err) {
			serializationFailureCount++
		} else {
			t.Fatalf("unexpected error (neither success nor 40001): %v", err)
		}
	}

	if successCount != 1 || serializationFailureCount != 1 {
		t.Fatalf("expected exactly 1 success and 1 serialization failure (40001), got success=%d, failures=%d", successCount, serializationFailureCount)
	}
}
