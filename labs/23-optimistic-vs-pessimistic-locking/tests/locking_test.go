package tests

import (
	"sync"
	"testing"

	"github.com/lukman/software-engineering-lab/labs/23-optimistic-vs-pessimistic-locking/internal/inventory"
)

func TestNaiveLostUpdate(t *testing.T) {
	store := inventory.NewStore()
	store.Seed(1, "Item A", 100)
	svc := inventory.NewService(store)

	const goroutines = 50
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			_ = svc.DeductNaive(1, 1)
		}()
	}
	wg.Wait()

	p, err := store.Get(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Lost update manifests because multiple goroutines overwrite each other's updates
	// Stock will NOT be 50 despite 50 successful deduct calls recorded
	if p.Stock == 50 {
		t.Errorf("expected lost update anomaly to manifest (stock != 50), got stock = %d", p.Stock)
	}
	t.Logf("Lost Update Demonstrated: 50 deduct calls occurred, but final stock is %d (expected 50 under proper locking)", p.Stock)
}

func TestPessimisticLocking(t *testing.T) {
	store := inventory.NewStore()
	store.Seed(2, "Item B", 100)
	svc := inventory.NewService(store)

	const goroutines = 50
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			err := svc.DeductPessimistic(2, 1)
			if err != nil {
				t.Errorf("deduct failed: %v", err)
			}
		}()
	}
	wg.Wait()

	p, err := store.Get(2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Stock != 50 {
		t.Fatalf("expected stock 50, got %d", p.Stock)
	}
}

func TestPessimisticLockingInsufficientStock(t *testing.T) {
	store := inventory.NewStore()
	store.Seed(20, "Limited Item", 2)
	svc := inventory.NewService(store)

	err1 := svc.DeductPessimistic(20, 2)
	if err1 != nil {
		t.Fatalf("expected success, got %v", err1)
	}

	err2 := svc.DeductPessimistic(20, 1)
	if err2 != inventory.ErrInsufficientStock {
		t.Fatalf("expected ErrInsufficientStock, got %v", err2)
	}
}

func TestOptimisticLockingConflict(t *testing.T) {
	store := inventory.NewStore()
	store.Seed(3, "Item C", 100)
	svc := inventory.NewService(store)

	const goroutines = 20
	var wg sync.WaitGroup
	wg.Add(goroutines)

	var conflictCount int64
	var successCount int64
	var mu sync.Mutex

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			err := svc.DeductOptimisticDirect(3, 1)
			mu.Lock()
			if err == inventory.ErrOptimisticLock {
				conflictCount++
			} else if err == nil {
				successCount++
			}
			mu.Unlock()
		}()
	}
	wg.Wait()

	p, _ := store.Get(3)
	if successCount+int64(p.Stock) != 100 {
		t.Fatalf("stock invariant violated: success %d + stock %d != 100", successCount, p.Stock)
	}
	if conflictCount == 0 {
		t.Fatalf("expected at least one optimistic lock conflict, got 0")
	}
	t.Logf("Optimistic Locking Conflict Demonstrated: %d succeeded, %d conflicts rejected", successCount, conflictCount)
}

func TestOptimisticLockingWithRetry(t *testing.T) {
	store := inventory.NewStore()
	store.Seed(4, "Item D", 100)
	svc := inventory.NewService(store)

	const goroutines = 20
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			_ = svc.DeductOptimisticWithRetry(4, 1, 10)
		}()
	}
	wg.Wait()

	p, _ := store.Get(4)
	if p.Stock != 100-int(store.Optimistically) {
		t.Fatalf("inconsistent stock: expected %d, got %d", 100-int(store.Optimistically), p.Stock)
	}
	t.Logf("Optimistic Retry Successful: %d total updates applied, retried conflicts resolved", store.Optimistically)
}

func TestAtomicConditionalUpdate(t *testing.T) {
	store := inventory.NewStore()
	store.Seed(5, "Item E", 100)
	svc := inventory.NewService(store)

	const goroutines = 50
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			err := svc.DeductAtomic(5, 1)
			if err != nil {
				t.Errorf("unexpected atomic deduct failure: %v", err)
			}
		}()
	}
	wg.Wait()

	p, _ := store.Get(5)
	if p.Stock != 50 {
		t.Fatalf("expected stock 50, got %d", p.Stock)
	}
}
