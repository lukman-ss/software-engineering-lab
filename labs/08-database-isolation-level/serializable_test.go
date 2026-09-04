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

// TestSerializable_ConcurrentUpdate_SerializationFailure proves that
// SERIALIZABLE isolation raises 40001 when conflicts are detected.
func TestSerializable_ConcurrentUpdate_SerializationFailure(t *testing.T) {
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
		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
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

		_, errA = tx.ExecContext(ctx, "UPDATE isolation_accounts SET balance = balance - 500000 WHERE id = 1")
		if errA == nil {
			errA = tx.Commit()
		}
	}()

	go func() {
		defer wg.Done()
		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
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

		_, errB = tx.ExecContext(ctx, "UPDATE isolation_accounts SET balance = balance - 500000 WHERE id = 1")
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

	t.Logf("Serializable Concurrent Update Error: errA=%v, errB=%v", errA, errB)

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

	t.Logf("Result: 1 TX Success, 1 TX Failure (SQLSTATE 40001).")
	t.Logf("This proves SERIALIZABLE does not 'just lock globally', but aborts conflicts.")
	t.Logf("Applications MUST catch 40001 and retry.")
}
