//go:build integration
// +build integration

package race

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/lukman-ss/software-engineering-lab/pkg/database"
)

// TestPostgresRowLock_ConcurrentStock menguji row locking dengan PostgreSQL.
//
// Menggunakan barrier synchronization untuk memaksa interleaving:
// A memperoleh lock duluan → B mencoba lock dan ter-block → A commit → B unblocked.
//
// Scenario (stock = 1, 2 concurrent):
//   - Transaction A memperoleh lock via SELECT ... FOR UPDATE
//   - Transaction B ter-block hingga A COMMIT
//   - A mengurangi stock dari 1 → 0
//   - B membaca stock = 0, CHECK gagal → Reject
//
// Expected:
//   - success = 1
//   - out of stock = 1
//   - final_stock = 0
func TestPostgresRowLock_ConcurrentStock(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, err := database.Connect(ctx, database.FromEnv())
	if err != nil {
		t.Skipf("PostgreSQL not available: %v", err)
	}
	defer db.Close()

	productID := "unit-oli-mesin-rowlock"
	initialStock := 1

	if err := setupTestInventory(ctx, db, productID, initialStock); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	defer cleanupTestInventory(ctx, db, productID)

	// Channels for deterministic coordination:
	// lockAcquired: A has SELECT ... FOR UPDATE completed
	// releaseA: signal A to proceed with UPDATE/COMMIT
	// bStarted: B has started and is blocked
	// bFinished: B has completed (or rejected)
	lockAcquired := make(chan struct{})
	releaseA := make(chan struct{})
	bStarted := make(chan struct{})
	bFinished := make(chan struct{})

	var wg sync.WaitGroup
	var successCount int
	var outOfStockCount int
	var unexpectedErrorCount int

	var aErr, bErr error
	wg.Add(2)

	// Transaction A — starts first, acquires lock, then waits for release signal
	go func() {
		defer wg.Done()
		aErr = trySellWithLockA(ctx, db, productID, lockAcquired, releaseA)
	}()

	// Transaction B — waits for lock to be acquired, then starts (gets blocked)
	go func() {
		defer wg.Done()
		<-lockAcquired // wait until A has the lock
		close(bStarted) // signal that B has started
		bErr = trySellWithLockB(ctx, db, productID)
		close(bFinished) // signal that B has completed
	}()

	// Wait for A to acquire lock, then B starts
	<-lockAcquired
	// B should now be blocked waiting for the lock
	// Verify B has started but not finished yet
	select {
	case <-bStarted:
		// B has started correctly - it's now blocked on SELECT FOR UPDATE
	default:
		t.Fatal("B should have started and be blocked")
	}

	select {
	case <-bFinished:
		t.Fatal("B should NOT have finished yet - it should be blocked on A's lock")
	default:
		// B is still blocked - this is correct
	}

	// Now release A to continue (UPDATE + COMMIT)
	close(releaseA)

	wg.Wait()

	// Classify results
	if aErr == nil {
		successCount++
	} else if errors.Is(aErr, ErrOutOfStock) {
		outOfStockCount++
	} else {
		unexpectedErrorCount++
	}

	if bErr == nil {
		successCount++
	} else if errors.Is(bErr, ErrOutOfStock) {
		outOfStockCount++
	} else {
		unexpectedErrorCount++
	}

	repo := NewPostgresRowLockRepository(db)
	finalStock, err := repo.GetStock(ctx, productID)
	if err != nil {
		t.Fatalf("get final stock: %v", err)
	}

	t.Logf("Request A: err=%v", aErr)
	t.Logf("Request B: err=%v", bErr)
	t.Logf("=== ROW LOCK RESULTS ===")
	t.Logf("Initial stock: %d", initialStock)
	t.Logf("Success: %d", successCount)
	t.Logf("Out of stock: %d", outOfStockCount)
	t.Logf("Unexpected errors: %d", unexpectedErrorCount)
	t.Logf("Final stock: %d", finalStock)

	// Assertions
	if unexpectedErrorCount > 0 {
		t.Errorf("expected 0 unexpected errors, got %d", unexpectedErrorCount)
	}
	if successCount != 1 {
		t.Errorf("expected success = 1, got %d", successCount)
	}
	if outOfStockCount != 1 {
		t.Errorf("expected out of stock = 1, got %d", outOfStockCount)
	}
	if finalStock != 0 {
		t.Errorf("expected final_stock = 0, got %d", finalStock)
	}

	invariant := successCount + finalStock
	if invariant != initialStock {
		t.Errorf("INVARIANT BROKEN: %d + %d != %d", successCount, finalStock, initialStock)
	}

	// Ensure B was actually blocked until A committed (not just racing past)
	if bErr == nil {
		t.Error("B should have been rejected (stock=0 after A), but succeeded")
	}

	t.Logf("✅ INVARIANT HOLDS: initial_stock = successful + final_stock")
}

// TestPostgresRowLock_HighContention menguji row locking dengan 500 concurrent requests.
func TestPostgresRowLock_HighContention(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	db, err := database.Connect(ctx, database.FromEnv())
	if err != nil {
		t.Skipf("PostgreSQL not available: %v", err)
	}
	defer db.Close()

	productID := "unit-oli-mesin-high-contention"
	const initialStock = 100
	const attempts = 500

	if err := setupTestInventory(ctx, db, productID, initialStock); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	defer cleanupTestInventory(ctx, db, productID)

	repo := NewPostgresRowLockRepository(db)

	var successCount int
	var outOfStockCount int
	var unexpectedErrorCount int
	var mu sync.Mutex

	ready, release := startGate(attempts)

	var wg sync.WaitGroup
	wg.Add(attempts)

	for i := 0; i < attempts; i++ {
		go func(i int) {
			defer wg.Done()
			ready <- struct{}{}
			<-release

			err := repo.TrySell(ctx, productID)
			mu.Lock()
			if err == nil {
				successCount++
			} else if errors.Is(err, ErrOutOfStock) {
				outOfStockCount++
			} else {
				unexpectedErrorCount++
			}
			mu.Unlock()
		}(i)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-ctx.Done():
		t.Fatal("test timeout: potential deadlock or DB connection exhaustion")
	}

	finalStock, err := repo.GetStock(ctx, productID)
	if err != nil {
		t.Fatalf("get final stock: %v", err)
	}

	t.Logf("=== ROW LOCK HIGH CONTENTION RESULTS ===")
	t.Logf("Initial stock: %d", initialStock)
	t.Logf("Success: %d", successCount)
	t.Logf("Out of stock: %d", outOfStockCount)
	t.Logf("Unexpected errors: %d", unexpectedErrorCount)
	t.Logf("Final stock: %d", finalStock)

	if unexpectedErrorCount > 0 {
		t.Errorf("expected 0 unexpected errors, got %d", unexpectedErrorCount)
	}
	if successCount != initialStock {
		t.Errorf("expected success = %d, got %d", initialStock, successCount)
	}
	if outOfStockCount != 400 {
		t.Errorf("expected out of stock = 400, got %d", outOfStockCount)
	}
	if finalStock != 0 {
		t.Errorf("expected final_stock = 0, got %d", finalStock)
	}

	invariant := successCount + finalStock
	if invariant != initialStock {
		t.Errorf("INVARIANT BROKEN: %d + %d != %d", successCount, finalStock, initialStock)
	}

	t.Logf("✅ INVARIANT HOLDS: initial_stock = successful + final_stock")
}

// trySellWithLockA performs TrySell in Transaction A (lock holder).
// Signals lockAcquired after SELECT ... FOR UPDATE, waits for release, then updates and commits.
func trySellWithLockA(ctx context.Context, db *sql.DB, productID string, lockAcquired chan struct{}, release chan struct{}) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	var stock int
	err = tx.QueryRowContext(ctx,
		"SELECT stock FROM inventory_products WHERE id = $1 FOR UPDATE",
		productID).Scan(&stock)
	if err != nil {
		return fmt.Errorf("lock and read: %w", err)
	}

	// Signal lock acquired so B can start its blocked attempt
	if lockAcquired != nil {
		close(lockAcquired)
	}

	// Wait until B has started and is blocked on the lock
	<-release

	if stock <= 0 {
		_ = tx.Rollback()
		return ErrOutOfStock
	}

	_, err = tx.ExecContext(ctx,
		"UPDATE inventory_products SET stock = $1 WHERE id = $2",
		stock-1, productID)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("update stock: %w", err)
	}

	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}

// trySellWithLockB performs TrySell in Transaction B (blocked until A commits).
func trySellWithLockB(ctx context.Context, db *sql.DB, productID string) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	var stock int
	err = tx.QueryRowContext(ctx,
		"SELECT stock FROM inventory_products WHERE id = $1 FOR UPDATE",
		productID).Scan(&stock)
	if err != nil {
		return fmt.Errorf("lock and read: %w", err)
	}

	if stock <= 0 {
		_ = tx.Rollback()
		return ErrOutOfStock
	}

	_, err = tx.ExecContext(ctx,
		"UPDATE inventory_products SET stock = $1 WHERE id = $2",
		stock-1, productID)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("update stock: %w", err)
	}

	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}
