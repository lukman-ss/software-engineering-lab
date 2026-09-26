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
			store.Get("concurrent-key")
		}(i)
	}

	wg.Wait()
}
