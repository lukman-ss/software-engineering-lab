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
	// lockVerified: B's lock wait state verified by pg_stat_activity
	lockAcquired := make(chan struct{})
	releaseA := make(chan struct{})
	lockVerified := make(chan struct{})

	var wg sync.WaitGroup
	var successCount int
	var outOfStockCount int
	var unexpectedErrorCount int
	var mu sync.Mutex

	var aErr, bErr error
	wg.Add(2)

	// Get B's DB connection for introspection
	connB, err := db.Conn(ctx)
	if err != nil {
		t.Fatalf("failed to get B connection: %v", err)
	}
	defer connB.Close()

	var backendPIDB int
	err = connB.QueryRowContext(ctx, "SELECT pg_backend_pid()").Scan(&backendPIDB)
	if err != nil {
		t.Fatalf("failed to get backend PID: %v", err)
	}
	t.Logf("Transaction B backend PID: %d", backendPIDB)

	// Transaction A — starts first, acquires lock, then waits for release signal
	go func() {
		defer wg.Done()
		aErr = trySellWithLockA(ctx, db, productID, lockAcquired, releaseA)
	}()

	// Transaction B verification goroutine — polls pg_stat_activity until B is waiting on lock
	go func() {
		defer wg.Done()

		// Wait for A to hold lock before starting B
		<-lockAcquired

		// Start transaction B that will block on the same row
		err := trySellWithLockBAsync(ctx, connB, productID)
		if err != nil {
			t.Errorf("Transaction B failed to start: %v", err)
			return
		}

		// Poll pg_stat_activity to verify B is waiting for lock
		// Using bounded polling with short intervals (not arbitrary sleep)
		const pollInterval = 10 * time.Millisecond
		const maxPolls = 100 // 1 second total timeout for lock detection
		var isWaiting bool

		for i := 0; i < maxPolls; i++ {
			var waitEventType string
			err = connB.QueryRowContext(ctx,
				`SELECT wait_event_type FROM pg_stat_activity WHERE pid = $1`,
				backendPIDB).Scan(&waitEventType)
			if err == nil && waitEventType == "Lock" {
				isWaiting = true
				break
			}
			time.Sleep(pollInterval)
		}

		if !isWaiting {
			t.Error("B did not enter lock wait state - test may be flaky or DB race condition")
			return
		}

		t.Logf("Verified: Backend %d is waiting on Lock (wait_event_type='Lock')", backendPIDB)
		close(lockVerified)
	}()

	// Wait for A to acquire lock
	select {
	case <-lockAcquired:
		t.Log("A has acquired row lock via SELECT ... FOR UPDATE")
	case <-ctx.Done():
		t.Fatal("transaction A failed to acquire row lock in time")
	}

	// Wait for B to be verified as waiting on the lock
	select {
	case <-lockVerified:
		t.Log("B is confirmed waiting on row lock (PostgreSQL verified)")
	case <-ctx.Done():
		t.Fatal("transaction B did not enter lock wait state in time")
	}

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
		// Ensure transaction is always cleaned up
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

	// Signal lock acquired so B can start its blocked attempt
	if lockAcquired != nil {
		close(lockAcquired)
	}

	// Wait until B has started and is blocked on the lock (with context timeout)
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

	// Clear tx so defer doesn't double-rollback
	tx = nil
	return nil
}

// trySellWithLockBAsync starts Transaction B on a dedicated connection for lock introspection.
// The transaction B will block on SELECT ... FOR UPDATE until A commits.
func trySellWithLockBAsync(ctx context.Context, conn *sql.Conn, productID string) error {
	// Start transaction using the dedicated connection
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx B: %w", err)
	}

	// Execute SELECT ... FOR UPDATE which will block if A holds the lock
	var stock int
	err = tx.QueryRowContext(ctx,
		"SELECT stock FROM inventory_products WHERE id = $1 FOR UPDATE",
		productID).Scan(&stock)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("lock B: %w", err)
	}

	// At this point, either B got the lock (A already committed) or B is blocked
	// If blocked, pg_stat_activity will show wait_event_type = "Lock"

	// Check if we got the lock and stock > 0
	if stock <= 0 {
		_ = tx.Rollback()
		return ErrOutOfStock
	}

	// Update and commit
	_, err = tx.ExecContext(ctx,
		"UPDATE inventory_products SET stock = $1 WHERE id = $2",
		stock-1, productID)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("update stock B: %w", err)
	}

	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("commit B: %w", err)
	}

	return nil
}