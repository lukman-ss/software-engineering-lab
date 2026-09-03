package isolation_test

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"testing"
	"time"

	_ "github.com/lib/pq"
	isolation "github.com/lukman-ss/software-engineering-lab/labs/08-database-isolation-level"
	"errors"
)

func TestTransfer_SelfTransferRejected(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()
	repo := isolation.NewPostgresWalletRepo(db)

	ctx := context.Background()
	resetTestState(t, ctx, db, repo, 1000, 1000, 1000)

	strategies := []struct {
		name string
		fn   func(ctx context.Context, fromID, toID int, amount int64) error
	}{
		{"Naive", repo.TransferNaive},
		{"WithLock", repo.TransferWithLock},
		{"RepeatableRead", repo.TransferRepeatableRead},
		{"Serializable", repo.TransferSerializable},
	}

	for _, s := range strategies {
		t.Run(s.name, func(t *testing.T) {
			err := s.fn(ctx, 1, 1, 100)
			if !errors.Is(err, isolation.ErrSameAccountTransfer) {
				t.Errorf("expected ErrSameAccountTransfer, got %v", err)
			}
			
			// Verify audit log empty
			var auditCount int
			err = db.QueryRow("SELECT COUNT(*) FROM isolation_transfer_audit").Scan(&auditCount)
			if err != nil {
				t.Fatalf("count audit: %v", err)
			}
			if auditCount != 0 {
				t.Errorf("expected 0 audit records, got %d", auditCount)
			}
			// Verify balance unchanged
			alice, err := repo.GetAccount(ctx, 1)
			if err != nil {
				t.Fatalf("get alice: %v", err)
			}
			if alice.Balance != 1000 {
				t.Errorf("expected balance 1000, got %d", alice.Balance)
			}
		})
	}
}

// TestNaiveTransfer_LostUpdate demonstrates the classic double-spend / lost-update bug:
// Alice has 1,000,000.
// Transfer A = 800,000 to Bob.
// Transfer B = 800,000 to Charlie.
// Result: Total money invariant breaks!
func TestNaiveTransfer_LostUpdate(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()
	repo := isolation.NewPostgresWalletRepo(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resetTestState(t, ctx, db, repo, 1000000, 1000000, 1000000)

	tx1Read := make(chan struct{})
	tx2Read := make(chan struct{})
	release := make(chan struct{})

	type txResult struct {
		err error
	}
	resultCh := make(chan txResult, 2)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		var res txResult
		defer func() { resultCh <- res }()

		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
		if err != nil {
			res.err = err
			return
		}
		defer tx.Rollback()

		var b1, b2 int64
		b1, res.err = repo.GetBalance(ctx, tx, 1)
		if res.err != nil {
			return
		}
		b2, res.err = repo.GetBalance(ctx, tx, 2)
		if res.err != nil {
			return
		}
		close(tx1Read)
		<-release

		if b1 >= 800000 {
			_, res.err = tx.ExecContext(ctx, "UPDATE isolation_accounts SET balance = $1 WHERE id = 1", b1-800000)
			if res.err == nil {
				_, res.err = tx.ExecContext(ctx, "UPDATE isolation_accounts SET balance = $1 WHERE id = 2", b2+800000)
			}
			if res.err == nil {
				res.err = tx.Commit()
			}
		}
	}()

	go func() {
		defer wg.Done()
		var res txResult
		defer func() { resultCh <- res }()

		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
		if err != nil {
			res.err = err
			return
		}
		defer tx.Rollback()

		var b1, b3 int64
		b1, res.err = repo.GetBalance(ctx, tx, 1)
		if res.err != nil {
			return
		}
		b3, res.err = repo.GetBalance(ctx, tx, 3)
		if res.err != nil {
			return
		}
		close(tx2Read)
		<-release

		if b1 >= 800000 {
			_, res.err = tx.ExecContext(ctx, "UPDATE isolation_accounts SET balance = $1 WHERE id = 1", b1-800000)
			if res.err == nil {
				_, res.err = tx.ExecContext(ctx, "UPDATE isolation_accounts SET balance = $1 WHERE id = 3", b3+800000)
			}
			if res.err == nil {
				res.err = tx.Commit()
			}
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

	res1 := <-resultCh
	res2 := <-resultCh
	if res1.err != nil {
		t.Fatalf("tx1 unexpected error: %v", res1.err)
	}
	if res2.err != nil {
		t.Fatalf("tx2 unexpected error: %v", res2.err)
	}

	alice, err := repo.GetAccount(ctx, 1)
	if err != nil {
		t.Fatalf("get alice: %v", err)
	}
	bob, err := repo.GetAccount(ctx, 2)
	if err != nil {
		t.Fatalf("get bob: %v", err)
	}
	charlie, err := repo.GetAccount(ctx, 3)
	if err != nil {
		t.Fatalf("get charlie: %v", err)
	}

	t.Logf("Naive Transfer Concurrency Result:")
	t.Logf("  Alice:   %d", alice.Balance)
	t.Logf("  Bob:     %d", bob.Balance)
	t.Logf("  Charlie: %d", charlie.Balance)

	// Both TXs read Alice=1,000,000, and each overwrote Alice balance with 1,000,000 - 800,000 = 200,000.
	// Bob received 800,000 (now 1,800,000).
	// Charlie received 800,000 (now 1,800,000).
	// Alice lost update manifests as Alice=200,000, Total=3,800,000 (800,000 created out of thin air).
	if alice.Balance != 200000 {
		t.Fatalf("expected Alice balance to be 200000 due to lost update, got %d", alice.Balance)
	}
	if bob.Balance != 1800000 {
		t.Fatalf("expected Bob balance to be 1800000, got %d", bob.Balance)
	}
	if charlie.Balance != 1800000 {
		t.Fatalf("expected Charlie balance to be 1800000, got %d", charlie.Balance)
	}

	totalMoney := alice.Balance + bob.Balance + charlie.Balance
	if totalMoney != 3800000 {
		t.Fatalf("expected total money to be 3800000 due to double-spend, got %d", totalMoney)
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resetTestState(t, ctx, db, repo, 1000000, 1000000, 1000000)

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

	alice, err := repo.GetAccount(ctx, 1)
	if err != nil {
		t.Fatalf("get alice: %v", err)
	}
	bob, err := repo.GetAccount(ctx, 2)
	if err != nil {
		t.Fatalf("get bob: %v", err)
	}
	charlie, err := repo.GetAccount(ctx, 3)
	if err != nil {
		t.Fatalf("get charlie: %v", err)
	}

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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resetTestState(t, ctx, db, repo, 1000000, 1000000, 1000000)

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

	alice, err := repo.GetAccount(ctx, 1)
	if err != nil {
		t.Fatalf("get alice: %v", err)
	}
	bob, err := repo.GetAccount(ctx, 2)
	if err != nil {
		t.Fatalf("get bob: %v", err)
	}
	charlie, err := repo.GetAccount(ctx, 3)
	if err != nil {
		t.Fatalf("get charlie: %v", err)
	}

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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Alice: 50,000, Bob: 50,000, Charlie: 50,000. Total = 150,000
	resetTestState(t, ctx, db, repo, 50000, 50000, 50000)

	const totalTransfers = 100
	var wg sync.WaitGroup
	wg.Add(totalTransfers)

	startGate := make(chan struct{})
	errChan := make(chan error, totalTransfers)

	for i := 0; i < totalTransfers; i++ {
		from := (i % 3) + 1
		to := ((i + 1) % 3) + 1

		go func(fromID, toID int) {
			defer wg.Done()
			<-startGate // wait for barrier release to maximize contention

			// Each worker attempts to transfer 1,000
			if err := repo.TransferWithLock(ctx, fromID, toID, 1000); err != nil {
				errChan <- fmt.Errorf("transfer %d->%d: %w", fromID, toID, err)
			}
		}(from, to)
	}

	close(startGate)
	wg.Wait()
	close(errChan)

	for err := range errChan {
		t.Fatalf("unexpected transfer error in high contention: %v", err)
	}

	alice, err := repo.GetAccount(ctx, 1)
	if err != nil {
		t.Fatalf("get alice: %v", err)
	}
	bob, err := repo.GetAccount(ctx, 2)
	if err != nil {
		t.Fatalf("get bob: %v", err)
	}
	charlie, err := repo.GetAccount(ctx, 3)
	if err != nil {
		t.Fatalf("get charlie: %v", err)
	}

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
