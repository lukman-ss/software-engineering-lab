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

	// Ensure schema exists
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
// two consecutive SELECTs in the same transaction can observe different values
// when another transaction modifies and commits data in between.
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

// TestRepeatableRead_SnapshotIsolation proves that in REPEATABLE READ,
// a transaction sees a consistent snapshot of the database taken at transaction start.
// Even if TX2 commits an update, TX1's subsequent SELECT still returns the initial value.
func TestRepeatableRead_SnapshotIsolation(t *testing.T) {
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
		<-tx2Committed  // Wait for TX2 commit

		// 2nd Read (must see snapshot, ignoring TX2's commit)
		tx1SecondRead, tx1Err = repo.GetBalance(ctx, tx1, 1)
		if tx1Err != nil {
			return
		}

		_ = tx1.Commit()
	}()

	// Transaction 2: Writer
	go func() {
		defer wg.Done()
		<-tx1Ready

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
		close(tx2Committed)
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

// TestNaiveTransfer_LostUpdate demonstrates the classic double-spend / lost-update bug:
// Alice has 1,000,000.
// Transfer A = 800,000 to Bob.
// Transfer B = 800,000 to Charlie.
// Both read Alice balance as 1,000,000 concurrently. Both think Alice has enough funds!
// Result: Both succeed, Alice loses 800k only once, or invariant total money breaks!
func TestNaiveTransfer_LostUpdate(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()
	repo := isolation.NewPostgresWalletRepo(db)

	ctx := context.Background()
	_ = repo.ResetAccounts(ctx, 1000000, 1000000, 1000000)

	// Barrier to force concurrent read before write
	tx1Read := make(chan struct{})
	tx2Read := make(chan struct{})
	release := make(chan struct{})

	var tx1Err, tx2Err error
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
		if err != nil {
			tx1Err = err
			return
		}
		defer tx.Rollback()

		b1, _ := repo.GetBalance(ctx, tx, 1)
		b2, _ := repo.GetBalance(ctx, tx, 2)
		close(tx1Read)
		<-release

		if b1 >= 800000 {
			_, tx1Err = tx.ExecContext(ctx, "UPDATE isolation_accounts SET balance = $1 WHERE id = 1", b1-800000)
			_, tx1Err = tx.ExecContext(ctx, "UPDATE isolation_accounts SET balance = $1 WHERE id = 2", b2+800000)
			_ = tx.Commit()
		}
	}()

	go func() {
		defer wg.Done()
		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
		if err != nil {
			tx2Err = err
			return
		}
		defer tx.Rollback()

		b1, _ := repo.GetBalance(ctx, tx, 1)
		b3, _ := repo.GetBalance(ctx, tx, 3)
		close(tx2Read)
		<-release

		if b1 >= 800000 {
			_, tx2Err = tx.ExecContext(ctx, "UPDATE isolation_accounts SET balance = $1 WHERE id = 1", b1-800000)
			_, tx2Err = tx.ExecContext(ctx, "UPDATE isolation_accounts SET balance = $1 WHERE id = 3", b3+800000)
			_ = tx.Commit()
		}
	}()

	<-tx1Read
	<-tx2Read
	close(release)
	wg.Wait()
	
	_ = tx1Err
	_ = tx2Err

	alice, _ := repo.GetAccount(ctx, 1)
	bob, _ := repo.GetAccount(ctx, 2)
	charlie, _ := repo.GetAccount(ctx, 3)

	t.Logf("Naive Transfer Concurrency Result:")
	t.Logf("  Alice:   %d", alice.Balance)
	t.Logf("  Bob:     %d", bob.Balance)
	t.Logf("  Charlie: %d", charlie.Balance)

	// Total initial money was 3,000,000.
	// With naive transfer, Bob got 800k (1.8M), Charlie got 800k (1.8M), and Alice balance is 200k.
	// Total = 3,800,000! 800,000 created out of thin air due to lost update!
	totalMoney := alice.Balance + bob.Balance + charlie.Balance
	if totalMoney > 3000000 {
		t.Logf("✅ Successfully reproduced Lost Update bug! Total money grew from 3,000,000 to %d (money created out of thin air)", totalMoney)
	} else {
		t.Fatalf("expected lost update to manifest under naive concurrent transfer, total was %d", totalMoney)
	}
}

// TestSafeTransferWithLock proves that SELECT ... FOR UPDATE prevents double-spending / lost update
func TestSafeTransferWithLock(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()
	repo := isolation.NewPostgresWalletRepo(db)

	ctx := context.Background()
	_ = repo.ResetAccounts(ctx, 1000000, 1000000, 1000000)

	var wg sync.WaitGroup
	var errA, errB error
	wg.Add(2)

	// Transfer A: 800,000 from Alice to Bob
	go func() {
		defer wg.Done()
		errA = repo.TransferWithLock(ctx, 1, 2, 800000)
	}()

	// Transfer B: 800,000 from Alice to Charlie
	go func() {
		defer wg.Done()
		errB = repo.TransferWithLock(ctx, 1, 3, 800000)
	}()

	wg.Wait()

	alice, _ := repo.GetAccount(ctx, 1)
	bob, _ := repo.GetAccount(ctx, 2)
	charlie, _ := repo.GetAccount(ctx, 3)

	t.Logf("Safe Transfer with Row Lock Result:")
	t.Logf("  ErrA: %v, ErrB: %v", errA, errB)
	t.Logf("  Alice: %d, Bob: %d, Charlie: %d", alice.Balance, bob.Balance, charlie.Balance)

	// Exactly one transfer MUST succeed and the other MUST fail with ErrInsufficientFunds
	successCount := 0
	if errA == nil {
		successCount++
	}
	if errB == nil {
		successCount++
	}

	if successCount != 1 {
		t.Fatalf("expected exactly 1 transfer to succeed, got %d", successCount)
	}

	totalMoney := alice.Balance + bob.Balance + charlie.Balance
	if totalMoney != 3000000 {
		t.Fatalf("money conservation invariant broken: expected total 3000000, got %d", totalMoney)
	}
	if alice.Balance != 200000 {
		t.Fatalf("expected Alice balance 200000, got %d", alice.Balance)
	}
}

// TestRepeatableRead_ConcurrentUpdate_SerializationFailure proves that PostgreSQL's
// REPEATABLE READ prevents lost updates by raising error 40001 (serialization failure)
// when two transactions try to update the same row concurrently.
func TestRepeatableRead_ConcurrentUpdate_SerializationFailure(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()
	repo := isolation.NewPostgresWalletRepo(db)

	ctx := context.Background()
	_ = repo.ResetAccounts(ctx, 1000000, 1000000, 1000000)

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

		_, _ = repo.GetBalance(ctx, tx, 1)
		close(tx1Read)
		<-release

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

		_, _ = repo.GetBalance(ctx, tx, 1)
		close(tx2Read)
		<-release

		_, errB = tx.ExecContext(ctx, "UPDATE isolation_accounts SET balance = balance - 800000 WHERE id = 1")
		if errB == nil {
			errB = tx.Commit()
		}
	}()

	<-tx1Read
	<-tx2Read
	close(release)
	wg.Wait()

	t.Logf("Repeatable Read Concurrent Update Error: errA=%v, errB=%v", errA, errB)

	// One should succeed, the other MUST fail with 40001 (serialization failure)
	hasSerializationFailure := (errA != nil && isolation.IsSerializationError(errA)) ||
		(errB != nil && isolation.IsSerializationError(errB))

	if !hasSerializationFailure {
		t.Fatalf("expected serialization failure (40001) in one of the concurrent transactions")
	}
}

// TestSerializableWithRetry demonstrates how applications handle SERIALIZABLE isolation level
// by retrying on serialization failure (40001).
func TestSerializableWithRetry(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()
	repo := isolation.NewPostgresWalletRepo(db)

	ctx := context.Background()
	_ = repo.ResetAccounts(ctx, 2000000, 1000000, 1000000)

	// 10 concurrent transfers of 100,000 from Alice to Bob
	numTransfers := 10
	var wg sync.WaitGroup
	wg.Add(numTransfers)

	errChan := make(chan error, numTransfers)

	for i := 0; i < numTransfers; i++ {
		go func() {
			defer wg.Done()
			err := repo.TransferSerializableWithRetry(ctx, 1, 2, 100000, 10)
			if err != nil {
				errChan <- err
			}
		}()
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		t.Fatalf("transfer failed despite retries: %v", err)
	}

	alice, _ := repo.GetAccount(ctx, 1)
	bob, _ := repo.GetAccount(ctx, 2)

	t.Logf("Serializable with Retry Results:")
	t.Logf("  Alice: %d (expected 1000000)", alice.Balance)
	t.Logf("  Bob:   %d (expected 2000000)", bob.Balance)

	if alice.Balance != 1000000 || bob.Balance != 2000000 {
		t.Fatalf("incorrect balance after serializable transfers: Alice=%d, Bob=%d", alice.Balance, bob.Balance)
	}
}
