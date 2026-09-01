//go:build integration
// +build integration

package race

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/lukman-ss/software-engineering-lab/pkg/database"
)

// TestPostgresRowLock_ConcurrentStock menguji row locking dengan PostgreSQL.
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

	repo := NewPostgresRowLockRepository(db)

	var wg sync.WaitGroup
	var successCount, rejectCount int
	var mu sync.Mutex

	wg.Add(2)

	for i := 0; i < 2; i++ {
		go func(id int) {
			defer wg.Done()
			err := repo.TrySell(ctx, productID)
			mu.Lock()
			if err == nil {
				successCount++
			} else {
				rejectCount++
			}
			mu.Unlock()
			t.Logf("Request %d: success=%v, err=%v", id, err == nil, err)
		}(i)
	}

	wg.Wait()

	finalStock, err := repo.GetStock(ctx, productID)
	if err != nil {
		t.Fatalf("get final stock: %v", err)
	}

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

	t.Logf("✅ INVARIANT HOLDS: initial_stock = successful + final_stock")
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
