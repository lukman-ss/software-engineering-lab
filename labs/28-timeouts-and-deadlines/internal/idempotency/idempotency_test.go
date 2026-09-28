package idempotency

import (
	"sync"
	"testing"
	"time"
)

func TestStore_GetSet(t *testing.T) {
	store := NewStore(100 * time.Millisecond)

	_, ok := store.Get("key-1")
	if ok {
		t.Fatalf("expected key-1 not found")
	}

	store.Set("key-1", "order-123")

	res, ok := store.Get("key-1")
	if !ok || res != "order-123" {
		t.Fatalf("expected order-123, got %s (ok=%v)", res, ok)
	}

	time.Sleep(120 * time.Millisecond)
	_, ok = store.Get("key-1")
	if ok {
		t.Fatalf("expected key-1 to expire")
	}
}

func TestStore_ConcurrentAccess(t *testing.T) {
	store := NewStore(1 * time.Second)
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func(id int) {
			defer wg.Done()
			store.Set("concurrent-key", "result")
		}(i)
		go func(id int) {
			defer wg.Done()
			res, ok := store.Get("concurrent-key")
			if ok && res != "result" {
				t.Errorf("unexpected value: %s", res)
			}
		}(i)
	}

	wg.Wait()
}

func TestStore_LazyEvictionOnGet(t *testing.T) {
	store := NewStore(20 * time.Millisecond)
	store.Set("expire-key", "val")

	time.Sleep(30 * time.Millisecond)

	val, ok := store.Get("expire-key")
	if ok || val != "" {
		t.Fatalf("expected key to be expired and empty, got %q, ok=%v", val, ok)
	}

	// Verify key was removed from map
	store.mu.RLock()
	_, exists := store.records["expire-key"]
	store.mu.RUnlock()
	if exists {
		t.Fatalf("expected key to be deleted from map on lazy eviction")
	}
}
