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

// TestPhantomRead_ReadCommitted demonstrates that in READ COMMITTED,
// a transaction querying COUNT(*) or range queries observes phantom rows inserted by another transaction.
func TestPhantomRead_ReadCommitted(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()
	repo := isolation.NewPostgresWalletRepo(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resetTestState(t, ctx, db, repo, 1000000, 1000000, 1000000)

	if _, err := db.ExecContext(ctx, "INSERT INTO isolation_invoices (amount, status) VALUES (10000, 'PAID'), (20000, 'PAID')"); err != nil {
		t.Fatalf("insert invoices: %v", err)
	}

	tx1Ready := make(chan struct{})
	tx2Committed := make(chan struct{})

	var count1, count2 int
	var tx1Err, tx2Err error

	var wg sync.WaitGroup
	wg.Add(2)

	// TX1: Read Committed Range Reader
	go func() {
		defer wg.Done()
		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
		if err != nil {
			tx1Err = err
			return
		}
		defer tx.Rollback()

		// 1st Count query
		tx1Err = tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM isolation_invoices WHERE status = 'PAID'").Scan(&count1)
		if tx1Err != nil {
			return
		}

		close(tx1Ready)

		select {
		case <-tx2Committed:
		case <-ctx.Done():
			tx1Err = ctx.Err()
			return
		}

		// 2nd Count query (in same transaction)
		tx1Err = tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM isolation_invoices WHERE status = 'PAID'").Scan(&count2)
		if tx1Err != nil {
			return
		}

		if err := tx.Commit(); err != nil {
			tx1Err = err
		}
	}()

	// TX2: Insert Phantom Row & Commit
	go func() {
		defer wg.Done()
		defer close(tx2Committed)

		select {
		case <-tx1Ready:
		case <-ctx.Done():
			tx2Err = ctx.Err()
			return
		}

		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
		if err != nil {
			tx2Err = err
			return
		}
		defer tx.Rollback()

		_, tx2Err = tx.ExecContext(ctx, "INSERT INTO isolation_invoices (amount, status) VALUES (30000, 'PAID')")
		if tx2Err != nil {
			return
		}

		tx2Err = tx.Commit()
	}()

	wg.Wait()

	if tx1Err != nil || tx2Err != nil {
		t.Fatalf("unexpected error: tx1=%v, tx2=%v", tx1Err, tx2Err)
	}

	t.Logf("READ COMMITTED Phantom Read Result: count1=%d, count2=%d", count1, count2)
	if count1 == count2 {
		t.Fatalf("expected phantom read in READ COMMITTED, but count remained %d", count1)
	}
	if count1 != 2 || count2 != 3 {
		t.Errorf("unexpected counts: count1=%d (want 2), count2=%d (want 3)", count1, count2)
	}
}

// TestPhantomRead_RepeatableRead proves that in PostgreSQL,
// REPEATABLE READ prevents Phantom Read due to MVCC snapshot isolation taken at tx start.
func TestPhantomRead_RepeatableRead(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()
	repo := isolation.NewPostgresWalletRepo(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resetTestState(t, ctx, db, repo, 1000000, 1000000, 1000000)

	if _, err := db.ExecContext(ctx, "INSERT INTO isolation_invoices (amount, status) VALUES (10000, 'PAID'), (20000, 'PAID')"); err != nil {
		t.Fatalf("insert invoices: %v", err)
	}

	tx1Ready := make(chan struct{})
	tx2Committed := make(chan struct{})

	var count1, count2 int
	var tx1Err, tx2Err error

	var wg sync.WaitGroup
	wg.Add(2)

	// TX1: Repeatable Read
	go func() {
		defer wg.Done()
		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
		if err != nil {
			tx1Err = err
			return
		}
		defer tx.Rollback()

		// 1st Count query
		tx1Err = tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM isolation_invoices WHERE status = 'PAID'").Scan(&count1)
		if tx1Err != nil {
			return
		}

		close(tx1Ready)

		select {
		case <-tx2Committed:
		case <-ctx.Done():
			tx1Err = ctx.Err()
			return
		}

		// 2nd Count query
		tx1Err = tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM isolation_invoices WHERE status = 'PAID'").Scan(&count2)
		if tx1Err != nil {
			return
		}

		if err := tx.Commit(); err != nil {
			tx1Err = err
		}
	}()

	// TX2: Insert Phantom Row & Commit
	go func() {
		defer wg.Done()
		defer close(tx2Committed)

		select {
		case <-tx1Ready:
		case <-ctx.Done():
			tx2Err = ctx.Err()
			return
		}

		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
		if err != nil {
			tx2Err = err
			return
		}
		defer tx.Rollback()

		_, tx2Err = tx.ExecContext(ctx, "INSERT INTO isolation_invoices (amount, status) VALUES (30000, 'PAID')")
		if tx2Err != nil {
			return
		}

		tx2Err = tx.Commit()
	}()

	wg.Wait()

	if tx1Err != nil || tx2Err != nil {
		t.Fatalf("unexpected error: tx1=%v, tx2=%v", tx1Err, tx2Err)
	}

	t.Logf("REPEATABLE READ Phantom Prevention Result: count1=%d, count2=%d", count1, count2)
	if count1 != count2 {
		t.Fatalf("PostgreSQL REPEATABLE READ snapshot isolation violated: count1=%d, count2=%d", count1, count2)
	}
	if count1 != 2 {
		t.Errorf("expected count 2, got %d", count1)
	}
}
