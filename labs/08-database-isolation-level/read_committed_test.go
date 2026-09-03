package isolation_test

import (
	"context"
	"database/sql"
	"os"
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
		t.Skipf("PostgreSQL not available: %v. Skipping integration test.", err)
	}

	schemaSQL := `
		CREATE TABLE IF NOT EXISTS isolation_accounts (
			id SERIAL PRIMARY KEY,
			owner VARCHAR(100) NOT NULL,
			balance BIGINT NOT NULL CHECK (balance >= 0)
		);
		CREATE TABLE IF NOT EXISTS isolation_transfer_audit (
			id SERIAL PRIMARY KEY,
			from_account_id INT NOT NULL,
			to_account_id INT NOT NULL,
			amount BIGINT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);
		CREATE TABLE IF NOT EXISTS isolation_invoices (
			id SERIAL PRIMARY KEY,
			amount BIGINT NOT NULL,
			status VARCHAR(20) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);
		INSERT INTO isolation_accounts (id, owner, balance) VALUES
			(1, 'Alice', 1000000),
			(2, 'Bob', 1000000),
			(3, 'Charlie', 1000000)
		ON CONFLICT (id) DO UPDATE SET balance = EXCLUDED.balance;
	`
	if _, err := db.Exec(schemaSQL); err != nil {
		t.Fatalf("failed to initialize schema: %v", err)
	}

	return db
}

// TestReadCommitted_NonRepeatableRead proves that in READ COMMITTED,
// two consecutive SELECTs in the same transaction observe different values
// if another transaction commits an UPDATE in between.
func TestReadCommitted_NonRepeatableRead(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()
	repo := isolation.NewPostgresWalletRepo(db)

	ctx := context.Background()
	_ = repo.ResetAccounts(ctx, 1000000, 1000000, 1000000)

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
		<-tx2Committed  // Wait until TX2 has updated & committed

		// 2nd Read in same transaction
		tx1SecondRead, tx1Err = repo.GetBalance(ctx, tx1, 1)
		if tx1Err != nil {
			return
		}

		_ = tx1.Commit()
	}()

	// Transaction 2: Writer modifies Alice's balance
	go func() {
		defer wg.Done()
		<-tx1Ready // Wait until TX1 has performed its first read

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
		close(tx2Committed) // Signal TX1 that commit is done
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
