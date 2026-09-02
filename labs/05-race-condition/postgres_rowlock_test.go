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

// workerResult stores result from goroutine for deterministic test failure handling
type workerResult struct {
	name string
	err  error
}

// TestPostgresRowLock_ConcurrentStock menguji row locking dengan PostgreSQL.
//
// Menggunakan PostgreSQL lock introspection untuk membuktikan blocking secara deterministik:
// Transaction A memperoleh lock → Transaction B menunggu pada pg_stat_activity (wait_event_type='Lock')
//
// Sequence yang dibuktikan:
//  1. A BEGIN, SELECT ... FOR UPDATE, acquires row lock
//  2. A closes lockAcquired channel (signal)
//  3. B gets dedicated connection, starts transaction, tries SELECT ... FOR UPDATE (blocked)
//  4. Main verifies B is waiting via pg_stat_activity (wait_event_type = 'Lock')
//  5. close(releaseA) → A UPDATE + COMMIT
//  6. B unblocked → reads stock=0 → ErrOutOfStock
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
	// lockAcquired: A has completed SELECT ... FOR UPDATE
	// releaseA: signal A to proceed with UPDATE/COMMIT
	lockAcquired := make(chan struct{})
	releaseA := make(chan struct{})

	// worker results channel for deterministic error reporting
	resultsCh := make(chan workerResult, 2)

	// Get dedicated connection for B BEFORE starting goroutine
	// This ensures pg_backend_pid() is the same throughout B's transaction
	connB, err := db.Conn(ctx)
	if err != nil {
		t.Fatalf("failed to get dedicated connection for B: %v", err)
	}
	defer connB.Close()

	// Get B's backend PID
	var backendPIDB int
	err = connB.QueryRowContext(ctx, "SELECT pg_backend_pid()").Scan(&backendPIDB)
	if err != nil {
		t.Fatalf("failed to get B's backend PID: %v", err)
	}
	t.Logf("Transaction B backend PID: %d", backendPIDB)

	// Transaction A — starts first, acquires lock, then waits for release signal
	go func() {
		defer func() {
			if r := recover(); r != nil {
				resultsCh <- workerResult{name: "A", err: fmt.Errorf("panic in A: %v", r)}
			}
		}()
		aErr := trySellWithLockA(ctx, db, productID, lockAcquired, releaseA)
		resultsCh <- workerResult{name: "A", err: aErr}
	}()

	// Transaction B — gets wait of A's lock, then attempts to acquire (will block)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				resultsCh <- workerResult{name: "B", err: fmt.Errorf("panic in B: %v", r)}
			}
		}()
		bErr := trySellWithLockBOnConn(ctx, connB, productID)
		resultsCh <- workerResult{name: "B", err: bErr}
	}()

	// Wait for A to acquire lock (bounded)
	select {
	case <-lockAcquired:
		t.Log("A has acquired row lock via SELECT ... FOR UPDATE")
	case <-ctx.Done():
		t.Fatal("transaction A failed to acquire row lock in time")
	}

	// Now poll to verify B is waiting for lock using pg_stat_activity
	// B is blocked on SELECT FOR UPDATE (same row as A)
	verifiedBWaiting := waitForBWaitingOnLock(ctx, db, backendPIDB)

	if !verifiedBWaiting {
		t.Fatal("Could not verify B is waiting on lock - test may be unreliable")
	}
	t.Log("✓ Verified: Transaction B is blocked waiting for row lock via pg_stat_activity")

	// Critical: A still holds lock! B is still waiting.
	// Verify B has NOT finished yet (did not get lock and complete)
	select {
	case res := <-resultsCh:
		if res.err == nil {
			t.Fatal("B should NOT have completed yet - it should still be blocked on A's lock")
		}
		t.Logf("A goroutine or B goroutine completed early with error: %v", res.err)
		t.FailNow()
	default:
		// Neither has finished - this is the expected state
		t.Log("✓ Confirmed: B has NOT finished (still blocked)")
	}

	// Now release A to continue (UPDATE + COMMIT)
	t.Log("Releasing A to continue with UPDATE + COMMIT")
	close(releaseA)

	// Wait for both workers with bounded timeout
	var aErr, bErr error
	var gotA, gotB bool

	// Drain results
	for i := 0; i < 2; i++ {
		select {
		case res := <-resultsCh:
			if res.name == "A" {
				aErr = res.err
				gotA = true
			} else {
				bErr = res.err
				gotB = true
			}
		case <-ctx.Done():
			t.Fatal("test timeout: deadlock or goroutine leak waiting for results")
		}
	}

	if !gotA || !gotB {
		t.Fatal("failed to get results from both transactions")
	}

	repo := NewPostgresRowLockRepository(db)
	finalStock, err := repo.GetStock(ctx, productID)
	if err != nil {
		t.Fatalf("get final stock: %v", err)
	}

	// Classify results
	var successCount, outOfStockCount, unexpectedErrorCount int

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
// Uses pg_stat_activity to check wait_event_type = 'Lock' for the specific backend PID.
// Returns true when the backend with PID is in lock wait state.
func waitForBWaitingOnLock(ctx context.Context, db *sql.DB, backendPID int) bool {
	const pollInterval = 10 * time.Millisecond
	const maxPolls = 200 // 2 seconds total timeout

	for i := 0; i < maxPolls; i++ {
		var waitEventType sql.NullString
		err := db.QueryRowContext(ctx,
			`SELECT wait_event_type FROM pg_stat_activity WHERE pid = $1 AND state = 'active'`,
			backendPID).Scan(&waitEventType)
		if err == nil && waitEventType.Valid && waitEventType.String == "Lock" {
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

// trySellWithLockBOnConn performs TrySell in Transaction B on a dedicated connection.
// This allows reliable introspection via pg_stat_activity for PID tracking.
func trySellWithLockBOnConn(ctx context.Context, conn *sql.Conn, productID string) error {
	tx, err := conn.BeginTx(ctx, nil)
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

	return nil
}
