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
// Menggunakan PostgreSQL lock introspection untuk membuktikan blocking secara deterministik:
// Transaction A memperoleh lock → Transaction B menunggu pada PG_STATE = "active" + wait_event_type = "Lock"
//
// Sequence yang dibuktikan:
//   1. A memperoleh SELECT ... FOR UPDATE (lock acquired)
//   2. A signal lockAcquired
//   3. B memulai transaksi, mengeksekusi SELECT ... FOR UPDATE (blocked di DB backend)
//   4. B verified: pg_stat_activity menunjukkan wait_event_type = "Lock"
//   5. Main test memverifikasi kondisi ini sebelum melanjutkan
//   6. close(releaseA) → A UPDATE + COMMIT
//   7. B unblocked → membaca stock=0 → ErrOutOfStock
//
// Expected:
//   - success_count = 1
//   - out_of_stock_count = 1
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
	const initialStock = 1

	if err := setupTestInventory(ctx, db, productID, initialStock); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	defer cleanupTestInventory(ctx, db, productID)

	// Channels for deterministic coordination:
	// lockAcquired: A has SELECT ... FOR UPDATE completed
	// releaseA: signal A to proceed with UPDATE/COMMIT
	// bRunning: B transaction is running (started)
	// bFinished: B transaction completed (or rejected)
	lockAcquired := make(chan struct{})
	releaseA := make(chan struct{})
	bRunning := make(chan struct{})
	bFinished := make(chan struct{})

	var wg sync.WaitGroup
	var successCount int
	var outOfStockCount int
	var unexpectedErrorCount int
	var mu sync.Mutex

	var aErr, bErr error
	wg.Add(2)

	// Transaction A — starts first, acquires lock, then waits for release signal
	go func() {
		defer wg.Done()
		aErr = trySellWithLockA(ctx, db, productID, lockAcquired, releaseA)
	}()

	// Transaction B goroutine — will block on SELECT FOR UPDATE
	go func() {
		defer wg.Done()

		// Wait for A to acquire lock first
		<-lockAcquired

		// Signal B goroutine is now active
		close(bRunning)

		// Attempt to sell - this will block if A holds the lock
		bErr = trySellWithLockB(ctx, db, productID)
		close(bFinished)
	}()

	// Wait for A to acquire lock
	select {
	case <-lockAcquired:
		t.Log("A has acquired row lock via SELECT ... FOR UPDATE")
	case <-ctx.Done():
		t.Fatal("transaction A failed to acquire row lock in time")
	}

	// Now poll to verify B is waiting for lock using pg_stat_activity
	// B's goroutine has started and is trying to acquire the lock
	// We need to detect that B is blocked on the lock
	verifiedBWaiting := waitForBWaitingOnLock(ctx, db)

	if !verifiedBWaiting {
		t.Fatal("Could not verify B is waiting on lock - test may be unreliable")
	}

	t.Log("✓ Verified: Transaction B is blocked waiting for row lock")

	// Now release A to continue (UPDATE + COMMIT)
	t.Log("Releasing A to continue with UPDATE + COMMIT")
	close(releaseA)

	wg.Wait()

	// Classify results
	mu.Lock()
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
	mu.Unlock()

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

	t.Logf("✅ ROW LOCK BLOCKING PROOF COMPLETE: A lock → B waiting → A release → B rejected")
}

// waitForBWaitingOnLock polls PostgreSQL to verify B is waiting for a row lock.
// Uses pg_stat_activity to check wait_event_type = 'Lock'.
func waitForBWaitingOnLock(ctx context.Context, db *sql.DB) bool {
	// Get count of backends waiting on locks (excluding this one)
	const pollInterval = 10 * time.Millisecond
	const maxPolls = 100 // 1 second total timeout

	for i := 0; i < maxPolls; i++ {
		var waitingCount int
		err := db.QueryRowContext(ctx,
			`SELECT COUNT(*) FROM pg_stat_activity WHERE state = 'active' AND wait_event_type = 'Lock' AND pid != pg_backend_pid()`).Scan(&waitingCount)
		if err == nil && waitingCount > 0 {
			t.Logf("Verified: %d backend(s) waiting on lock", waitingCount)
			return true
		}
		time.Sleep(pollInterval)
	}

	return false
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
// Always ends with COMMIT or ROLLBACK - never leaves transaction open.
func trySellWithLockA(ctx context.Context, db *sql.DB, productID string, lockAcquired chan struct{}, release chan struct{}) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if tx != nil {
			_ = tx.Rollback()
		}
	}()

	var stock int
	err = tx.QueryRowContext(ctx,
		"SELECT stock FROM inventory_products WHERE id = $1 FOR UPDATE",
		productID).Scan(&stock)
	if err != nil {
		return fmt.Errorf("lock and read: %w", err)
	}

	if lockAcquired != nil {
		close(lockAcquired)
	}

	select {
	case <-release:
	case <-ctx.Done():
		return fmt.Errorf("context cancelled before release: %w", ctx.Err())
	}

	if stock <= 0 {
		return ErrOutOfStock
	}

	_, err = tx.ExecContext(ctx,
		"UPDATE inventory_products SET stock = $1 WHERE id = $2",
		stock-1, productID)
	if err != nil {
		return fmt.Errorf("update stock: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	tx = nil
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