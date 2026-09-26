package ratelimit

import (
	"sync"
	"testing"
	"time"
)

func TestTokenBucket_BurstAndRefill(t *testing.T) {
	tb := NewTokenBucket(3, 10) // capacity 3, refill 10/sec

	// Initial burst of 3 allowed
	for i := 0; i < 3; i++ {
		if !tb.Allow() {
			t.Fatalf("expected token at attempt %d", i+1)
		}
	}

	// 4th fails
	if tb.Allow() {
		t.Fatalf("expected exhaustion on 4th attempt")
	}

	// Wait 200ms -> should refill ~2 tokens
	time.Sleep(200 * time.Millisecond)

	if !tb.Allow() {
		t.Fatalf("expected token after refill sleep")
	}
}

func TestLeakyBucket_LeakRate(t *testing.T) {
	lb := NewLeakyBucket(2, 5) // capacity 2, leaks 5/sec

	if !lb.Allow() {
		t.Fatal("expected 1st allow")
	}
	if !lb.Allow() {
		t.Fatal("expected 2nd allow")
	}
	if lb.Allow() {
		t.Fatal("expected burst exceeding capacity to be rejected")
	}

	time.Sleep(250 * time.Millisecond) // leaks > 1 item
	if !lb.Allow() {
		t.Fatal("expected allow after leak")
	}
}

func TestRegistry_TenantIsolation(t *testing.T) {
	reg := NewRegistry(1, 10)
	t1 := reg.Get("tenant-a")
	t2 := reg.Get("tenant-b")

	if !t1.Allow() {
		t.Fatal("tenant-a first request allowed")
	}
	if t1.Allow() {
		t.Fatal("tenant-a second request blocked")
	}

	// tenant-b should still have full quota
	if !t2.Allow() {
		t.Fatal("tenant-b should be isolated from tenant-a")
	}
}

func TestTokenBucket_RetryAfterSeconds(t *testing.T) {
	tb := NewTokenBucket(2, 2) // cap 2, refill 2/sec
	if !tb.Allow() || !tb.Allow() {
		t.Fatal("expected initial burst allowed")
	}
	if tb.RetryAfterSeconds(1.0) == 0 {
		t.Fatal("expected retry after > 0 when bucket empty")
	}
	time.Sleep(600 * time.Millisecond)
	if tb.RetryAfterSeconds(1.0) != 0 {
		t.Fatal("expected retry after == 0 after refilling enough tokens")
	}
}

func TestTokenBucket_ConcurrencyRace(t *testing.T) {
	tb := NewTokenBucket(100, 100)
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				_ = tb.Allow()
			}
		}()
	}

	wg.Wait()
}
