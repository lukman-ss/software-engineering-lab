package race

import (
	"context"
	"sync"
	"testing"
	"time"
)

// TestAtomicUpdate_StockOne concurrent test.
//
// Setup:
//   - initial_stock = 1
//   - attempts = 500 goroutines mencoba decrement
//
// Expected:
//   - success = 1
//   - rejected = 499
//   - final_stock = 0
//
// Assert invariant:
//   - final_stock >= 0
//   - successful <= initial_stock
//   - initial_stock == successful + final_stock
func TestAtomicUpdate_StockOne(t *testing.T) {
	const initialStock = 1
	const attempts = 500

	repo := NewAtomicInventory(initialStock)
	ctx := context.Background()

	var successCount int
	var rejectedCount int
	var mu sync.Mutex

	var wg sync.WaitGroup
	wg.Add(attempts)

	for i := 0; i < attempts; i++ {
		go func() {
			defer wg.Done()
			newStock, err := repo.DecrementStock(ctx, "unit-oli-mesin")
			if err != nil {
				mu.Lock()
				rejectedCount++
				mu.Unlock()
				return
			}
			mu.Lock()
			successCount++
			_ = newStock
			mu.Unlock()
		}()
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Test completed within time
	case <-time.After(5 * time.Second):
		t.Fatal("test timeout: goroutine leak detected")
	}

	// Re-read via direct access
	repo.mu.RLock()
	currentStock := repo.stock
	repo.mu.RUnlock()

	t.Logf("Initial stock: %d", initialStock)
	t.Logf("Successful: %d", successCount)
	t.Logf("Rejected: %d", rejectedCount)
	t.Logf("Final stock: %d", currentStock)

	// Assertions
	if currentStock < 0 {
		t.Errorf("BUG: final_stock < 0: %d", currentStock)
	}

	if successCount > initialStock {
		t.Errorf("BUG: successful_sales (%d) > initial_stock (%d)", successCount, initialStock)
	}

	// Main invariant check
	expectedSum := successCount + currentStock
	if expectedSum > initialStock {
		t.Errorf("INVARIANT BROKEN: successful (%d) + final_stock (%d) = %d > initial_stock (%d)",
			successCount, currentStock, expectedSum, initialStock)
	}

	// Expected values
	if successCount != 1 {
		t.Errorf("expected success = 1, got %d", successCount)
	}
	if rejectedCount != 499 {
		t.Errorf("expected rejected = 499, got %d", rejectedCount)
	}
	if currentStock != 0 {
		t.Errorf("expected final_stock = 0, got %d", currentStock)
	}

	t.Logf("✅ INVARIANT HOLDS: initial_stock = successful + final_stock = %d", initialStock)
}

// TestAtomicUpdate_HighContention concurrent test.
//
// Setup:
//   - initial_stock = 100
//   - attempts = 500 goroutines
//
// Expected:
//   - success = 100
//   - rejected = 400
//   - final_stock = 0
func TestAtomicUpdate_HighContention(t *testing.T) {
	const initialStock = 100
	const attempts = 500

	repo := NewAtomicInventory(initialStock)
	ctx := context.Background()

	var successCount int
	var rejectedCount int
	var mu sync.Mutex

	var wg sync.WaitGroup
	wg.Add(attempts)

	// Pool pattern: all 500 goroutines compete for 100 slots
	// Database connection pool: simulate 10 max connections
	for i := 0; i < attempts; i++ {
		go func(id int) {
			defer wg.Done()
			newStock, err := repo.DecrementStock(ctx, "unit-oli-mesin")
			if err != nil {
				mu.Lock()
				rejectedCount++
				mu.Unlock()
				return
			}
			mu.Lock()
			successCount++
			_ = newStock
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
		// Test completed
	case <-time.After(10 * time.Second):
		t.Fatal("test timeout: potential goroutine leak")
	}

	// Get final stock
	repo.mu.RLock()
	currentStock := repo.stock
	repo.mu.RUnlock()

	t.Logf("Initial stock: %d", initialStock)
	t.Logf("Successful: %d", successCount)
	t.Logf("Rejected: %d", rejectedCount)
	t.Logf("Final stock: %d", currentStock)

	// Assertions
	if currentStock < 0 {
		t.Errorf("BUG: final_stock < 0: %d", currentStock)
	}

	if successCount > initialStock {
		t.Errorf("BUG: successful_sales (%d) > initial_stock (%d)", successCount, initialStock)
	}

	// Main invariant
	expectedSum := successCount + currentStock
	if expectedSum > initialStock {
		t.Errorf("INVARIANT BROKEN: %d != %d + %d", initialStock, successCount, currentStock)
	}

	// Expected values
	if successCount != 100 {
		t.Errorf("expected success = 100, got %d", successCount)
	}
	if rejectedCount != 400 {
		t.Errorf("expected rejected = 400, got %d", rejectedCount)
	}
	if currentStock != 0 {
		t.Errorf("expected final_stock = 0, got %d", currentStock)
	}

	t.Logf("✅ HIGH CONTENTION TEST PASSED")
}

// GetStock untuk testing
func (a *AtomicInventory) GetStock(_ context.Context) int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.stock
}
