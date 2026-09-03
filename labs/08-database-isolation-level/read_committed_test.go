package isolation_test

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	_ "github.com/lib/pq"
	isolation "github.com/lukman-ss/software-engineering-lab/labs/08-database-isolation-level"
	"github.com/lukman-ss/software-engineering-lab/pkg/database"
)

func getTestDB(t *testing.T) *sql.DB {
	t.Helper()

	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("POSTGRES_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("POSTGRES_USER")
	if user == "" {
		user = "postgres"
	}
	pass := os.Getenv("POSTGRES_PASSWORD")
	if pass == "" {
		pass = "postgres"
	}
	dbname := os.Getenv("POSTGRES_DB")
	if dbname == "" {
		dbname = "se_lab"
	}

	cfg := database.FromEnv()
	cfg.Host = host
	cfg.Port = port
	cfg.User = user
	cfg.Password = pass
	cfg.Database = dbname
	cfg.MaxOpenConns = 25
	cfg.MaxIdleConns = 5

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	db, err := database.Connect(ctx, cfg)
	if err != nil {
		if os.Getenv("REQUIRE_POSTGRES") == "1" {
			t.Fatalf("PostgreSQL required but unavailable: %v", err)
		}
		t.Skipf("PostgreSQL not available: %v. Skipping integration test.", err)
	}

	// Read and execute schema.sql reliably regardless of execution working directory
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("failed to get caller location")
	}
	schemaPath := filepath.Join(filepath.Dir(currentFile), "schema.sql")
	schemaSQL, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("failed to read schema.sql: %v", err)
	}

	if _, err := db.Exec(string(schemaSQL)); err != nil {
		t.Fatalf("failed to initialize schema: %v", err)
	}

	return db
}

func resetTestState(t *testing.T, ctx context.Context, db *sql.DB, repo *isolation.PostgresWalletRepo, alice, bob, charlie int64) {
	t.Helper()
	if err := repo.ResetAccounts(ctx, alice, bob, charlie); err != nil {
		t.Fatalf("reset accounts: %v", err)
	}
	if _, err := db.ExecContext(ctx, "DELETE FROM isolation_invoices"); err != nil {
		t.Fatalf("reset invoices: %v", err)
	}
}

// TestReadUncommitted_PostgresDoesNotAllowDirtyRead proves that in PostgreSQL,
// READ UNCOMMITTED behaves as READ COMMITTED and does NOT permit Dirty Reads.
func TestReadUncommitted_PostgresDoesNotAllowDirtyRead(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()
	repo := isolation.NewPostgresWalletRepo(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resetTestState(t, ctx, db, repo, 1000000, 1000000, 1000000)

	tx1UncommittedUpdateDone := make(chan struct{})
	tx2ReadDone := make(chan struct{})

	var tx2ObservedBalance int64
	var tx1Err, tx2Err error

	var wg sync.WaitGroup
	wg.Add(2)

	// TX1: Modifies balance from 1,000,000 -> 500,000 but does NOT commit yet
	go func() {
		defer wg.Done()
		tx1, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
		if err != nil {
			tx1Err = err
			return
		}
		defer tx1.Rollback()

		_, tx1Err = tx1.ExecContext(ctx, "UPDATE isolation_accounts SET balance = 500000 WHERE id = 1")
		if tx1Err != nil {
			return
		}

		close(tx1UncommittedUpdateDone)

		select {
		case <-tx2ReadDone:
		case <-ctx.Done():
			tx1Err = ctx.Err()
			return
		}

		// Rollback explicitly at end
	}()

	// TX2: READ UNCOMMITTED Reader
	go func() {
		defer wg.Done()
		defer close(tx2ReadDone)

		select {
		case <-tx1UncommittedUpdateDone:
		case <-ctx.Done():
			tx2Err = ctx.Err()
			return
		}

		tx2, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadUncommitted})
		if err != nil {
			tx2Err = err
			return
		}
		defer tx2.Rollback()

		tx2ObservedBalance, tx2Err = repo.GetBalance(ctx, tx2, 1)
		if tx2Err != nil {
			return
		}

		tx2Err = tx2.Commit()
	}()

	wg.Wait()

	if tx1Err != nil || tx2Err != nil {
		t.Fatalf("unexpected error: tx1=%v, tx2=%v", tx1Err, tx2Err)
	}

	t.Logf("READ UNCOMMITTED Observed Balance: %d", tx2ObservedBalance)

	// In PostgreSQL, READ UNCOMMITTED treats isolation as READ COMMITTED.
	// Therefore, TX2 must observe the committed balance (1,000,000), not the dirty uncommitted 500,000.
	if tx2ObservedBalance != 1000000 {
		t.Fatalf("expected committed balance 1000000 (no dirty read in Postgres), got %d", tx2ObservedBalance)
	}
}

// TestReadCommitted_NonRepeatableRead proves that in READ COMMITTED,
// two consecutive SELECTs in the same transaction observe different values
// if another transaction commits an UPDATE in between.
func TestReadCommitted_NonRepeatableRead(t *testing.T) {
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

	// Transaction 1: READ COMMITTED reader
	go func() {
		defer wg.Done()
		tx1, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
		if err != nil {
			tx1Err = err
			return
		}
		defer tx1.Rollback()

		// 1st Read
		tx1FirstRead, tx1Err = repo.GetBalance(ctx, tx1, 1)
		if tx1Err != nil {
			return
		}

		close(tx1Ready) // Signal TX2 to proceed

		select {
		case <-tx2Committed: // Wait until TX2 has updated & committed
		case <-ctx.Done():
			tx1Err = ctx.Err()
			return
		}

		// 2nd Read in same transaction
		tx1SecondRead, tx1Err = repo.GetBalance(ctx, tx1, 1)
		if tx1Err != nil {
			return
		}

		if err := tx1.Commit(); err != nil {
			tx1Err = err
		}
	}()

	// Transaction 2: Writer modifies Alice's balance
	go func() {
		defer wg.Done()
		defer close(tx2Committed) // Signal TX1 that TX2 goroutine is done

		select {
		case <-tx1Ready: // Wait until TX1 has performed its first read
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

	t.Logf("READ COMMITTED Non-Repeatable Read Proof:")
	t.Logf("  TX1 First Read:  %d", tx1FirstRead)
	t.Logf("  TX1 Second Read: %d", tx1SecondRead)

	if tx1FirstRead == tx1SecondRead {
		t.Fatalf("expected non-repeatable read under READ COMMITTED, but reads were identical (%d)", tx1FirstRead)
	}
	if tx1FirstRead != 1000000 || tx1SecondRead != 1500000 {
		t.Errorf("unexpected read values: 1st=%d (want 1000000), 2nd=%d (want 1500000)", tx1FirstRead, tx1SecondRead)
	}
}
