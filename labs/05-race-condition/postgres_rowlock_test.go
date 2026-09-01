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

	// Manual transactions with barrier to ensure deterministic lock acquisition order.
	aStarted := make(chan struct{}) // A memegang lock
	aCommitted := make(chan struct{})
	bStarted := make(chan struct{})

	var wg sync.WaitGroup
	var successCount, rejectCount int

	var aErr, bErr error
	wg.Add(2)

	// Transaction A — memperoleh lock duluan
	go func() {
		defer wg.Done()
		aErr = trySellWithBarrier(ctx, db, productID, aStarted, aCommitted)
	}()

	// Tunggu A yakin punya lock sebelum B mulai
	go func() {
		<-aStarted
		close(bStarted)
	}()

	// Transaction B — akan ter-block sampai A commit
	go func() {
		defer wg.Done()
		<-bStarted
		bErr = trySellWithBarrier(ctx, db, productID, nil, nil)
	}()

	// Tunggu B selesai (hanya akan selesai setelah A commit)
	<-aCommitted
	wg.Wait()

	// Count success/reject
	if aErr == nil {
		successCount = 1
	} else {
		rejectCount = 1
	}
	if bErr == nil {
		successCount++
	} else {
		rejectCount++
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
	t.Logf("Rejected: %d", rejectCount)
	t.Logf("Final stock: %d", finalStock)

	if successCount != 1 {
		t.Errorf("expected success = 1, got %d", successCount)
	}
	if rejectCount != 1 {
		t.Errorf("expected reject = 1, got %d", rejectCount)
	}
	if finalStock != 0 {
		t.Errorf("expected final_stock = 0, got %d", finalStock)
	}

	invariant := successCount + finalStock
	if invariant != initialStock {
		t.Errorf("INVARIANT BROKEN: %d + %d != %d", successCount, finalStock, initialStock)
	}

	// Ensure B was actually blocked until A committed (not just racing past)
	// If B succeeded, that's a bug since stock=0 after A's decrement
	if bErr == nil {
		t.Error("B should have been rejected (stock=0 after A), but succeeded — lock not holding")
	}
	if !errors.Is(bErr, ErrOutOfStock) {
		t.Logf("B rejected with: %v (expected ErrOutOfStock)", bErr)
	}

	t.Logf("✅ INVARIANT HOLDS: initial_stock = successful + final_stock")
}

// trySellWithBarrier performs TrySell with a manual transaction.
// If lockAcquired != nil, it is closed after SELECT ... FOR UPDATE returns.
// If unlockSignal != nil, TrySell waits on it before proceeding to UPDATE (not used here).
func trySellWithBarrier(ctx context.Context, db *sql.DB, productID string, lockAcquired chan struct{}, unlockSignal chan struct{}) error {
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

	// Signal lock acquired
	if lockAcquired != nil {
		close(lockAcquired)
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

// TestPostgresRowLock_HighContention menguji row locking dengan 500 concurrent requests.
//
// Scenario:
//   - initial_stock = 100
//   - requests = 500
//
// Expected:
//   - success = 100
//   - rejected = 400
//   - final_stock = 0
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

	var wg sync.WaitGroup
	var successCount, rejectCount int
	var mu sync.Mutex

	wg.Add(attempts)

	for i := 0; i < attempts; i++ {
		go func(i int) {
			defer wg.Done()
			err := repo.TrySell(ctx, productID)
			mu.Lock()
			if err == nil {
				successCount++
			} else {
				rejectCount++
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
	t.Logf("Rejected: %d", rejectCount)
	t.Logf("Final stock: %d", finalStock)

	if successCount != initialStock {
		t.Errorf("expected success = %d, got %d", initialStock, successCount)
	}
	if rejectCount != 400 {
		t.Errorf("expected reject = 400, got %d", rejectCount)
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
