package isolation_test

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"testing"

	_ "github.com/lib/pq"
	isolation "github.com/lukman-ss/software-engineering-lab/labs/08-database-isolation-level"
)

// TestNaiveTransfer_LostUpdate demonstrates the classic double-spend / lost-update bug:
// Alice has 1,000,000.
// Transfer A = 800,000 to Bob.
// Transfer B = 800,000 to Charlie.
// Result: Total money invariant breaks!
func TestNaiveTransfer_LostUpdate(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()
	repo := isolation.NewPostgresWalletRepo(db)

	ctx := context.Background()
	_ = repo.ResetAccounts(ctx, 1000000, 1000000, 1000000)

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

	totalMoney := alice.Balance + bob.Balance + charlie.Balance
	if totalMoney > 3000000 {
		t.Logf("✅ Successfully reproduced Lost Update bug! Total money grew from 3,000,000 to %d (money created out of thin air)", totalMoney)
	} else {
		t.Fatalf("expected lost update to manifest under naive concurrent transfer, total was %d", totalMoney)
	}
}

// TestSafeTransferWithLock_DeterministicLocking verifies that SELECT ... FOR UPDATE
// with deterministic lock ordering strictly enforces invariant:
// 1. balance >= 0
// 2. total money before == total money after
func TestSafeTransferWithLock_DeterministicLocking(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()
	repo := isolation.NewPostgresWalletRepo(db)

	ctx := context.Background()
	_ = repo.ResetAccounts(ctx, 1000000, 1000000, 1000000)

	var wg sync.WaitGroup
	var errA, errB error
	wg.Add(2)

	// Transfer A: 1 -> 2 (Alice to Bob: 800k)
	go func() {
		defer wg.Done()
		errA = repo.TransferWithLock(ctx, 1, 2, 800000)
	}()

	// Transfer B: 1 -> 3 (Alice to Charlie: 800k)
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

	// Invariant Checks
	if alice.Balance < 0 || bob.Balance < 0 || charlie.Balance < 0 {
		t.Fatalf("Invariant violated: negative balance found! Alice=%d, Bob=%d, Charlie=%d", alice.Balance, bob.Balance, charlie.Balance)
	}
	totalMoney := alice.Balance + bob.Balance + charlie.Balance
	if totalMoney != 3000000 {
		t.Fatalf("Invariant violated: total money changed from 3000000 to %d", totalMoney)
	}
}

// TestBidirectionalTransfers_DeadlockFree tests concurrent A->B and B->A transfers.
// Without deterministic locking order (MIN id first), this would cause a Deadlock (40P01).
func TestBidirectionalTransfers_DeadlockFree(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()
	repo := isolation.NewPostgresWalletRepo(db)

	ctx := context.Background()
	_ = repo.ResetAccounts(ctx, 1000000, 1000000, 1000000)

	const iterations = 50
	var wg sync.WaitGroup
	wg.Add(2)

	errChan := make(chan error, iterations*2)

	// Worker 1: Alice -> Bob (1 -> 2)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			if err := repo.TransferWithLock(ctx, 1, 2, 1000); err != nil {
				errChan <- fmt.Errorf("1->2 failed: %w", err)
			}
		}
	}()

	// Worker 2: Bob -> Alice (2 -> 1) concurrently
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			if err := repo.TransferWithLock(ctx, 2, 1, 1000); err != nil {
				errChan <- fmt.Errorf("2->1 failed: %w", err)
			}
		}
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		t.Fatalf("bidirectional transfer encountered error (potential deadlock!): %v", err)
	}

	alice, _ := repo.GetAccount(ctx, 1)
	bob, _ := repo.GetAccount(ctx, 2)
	charlie, _ := repo.GetAccount(ctx, 3)

	totalMoney := alice.Balance + bob.Balance + charlie.Balance
	if totalMoney != 3000000 {
		t.Fatalf("total money invariant violated: %d != 3000000", totalMoney)
	}
	t.Logf("Bidirectional transfers completed cleanly with zero deadlocks: Alice=%d, Bob=%d", alice.Balance, bob.Balance)
}

// TestHighContention_100ConcurrentTransfers verifies the system holds invariants
// under extremely high contention.
func TestHighContention_100ConcurrentTransfers(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()
	repo := isolation.NewPostgresWalletRepo(db)

	ctx := context.Background()
	// Alice: 50,000, Bob: 50,000, Charlie: 50,000. Total = 150,000
	_ = repo.ResetAccounts(ctx, 50000, 50000, 50000)

	const totalTransfers = 100
	var wg sync.WaitGroup
	wg.Add(totalTransfers)

	startGate := make(chan struct{})

	for i := 0; i < totalTransfers; i++ {
		from := (i % 3) + 1
		to := ((i + 1) % 3) + 1

		go func(fromID, toID int) {
			defer wg.Done()
			<-startGate // wait for barrier release to maximize contention

			// Each worker attempts to transfer 1,000
			_ = repo.TransferWithLock(ctx, fromID, toID, 1000)
		}(from, to)
	}

	close(startGate)
	wg.Wait()

	alice, _ := repo.GetAccount(ctx, 1)
	bob, _ := repo.GetAccount(ctx, 2)
	charlie, _ := repo.GetAccount(ctx, 3)

	t.Logf("100 Concurrent Transfers Final Balances:")
	t.Logf("  Alice:   %d", alice.Balance)
	t.Logf("  Bob:     %d", bob.Balance)
	t.Logf("  Charlie: %d", charlie.Balance)

	// INVARIANT ASSERTIONS:
	// 1. Balance must NEVER be negative
	if alice.Balance < 0 || bob.Balance < 0 || charlie.Balance < 0 {
		t.Fatalf("Invariant violated: negative balance observed! Alice=%d, Bob=%d, Charlie=%d", alice.Balance, bob.Balance, charlie.Balance)
	}

	// 2. Total money in the system must remain exactly 150,000
	total := alice.Balance + bob.Balance + charlie.Balance
	if total != 150000 {
		t.Fatalf("Invariant violated: total money changed from 150000 to %d (leak or phantom money!)", total)
	}
}
