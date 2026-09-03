package caching_test

import (
	"context"
	"sync"
	"testing"
	"time"

	caching "github.com/lukman-ss/software-engineering-lab/labs/04-caching"
)

// TestStampedeBrokenVersion demonstrates the BROKEN stampede pattern deterministically.
func TestStampedeBrokenVersion(t *testing.T) {
	db := caching.NewCounterRepository()
	cache := caching.NewMockCache()
	svc := caching.NewBrokenStampedeService(cache, db)
	ctx := context.Background()

	key := int64(1)
	numRequests := 10

	db.Block()

	var wg sync.WaitGroup
	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = svc.GetData(ctx, key)
		}()
	}

	// Ensure multiple requests reach the blocked repository
	for i := 0; i < numRequests; i++ {
		db.WaitUntilEntered()
	}

	// Unblock repository after all requests are waiting inside repo
	db.Unblock()
	wg.Wait()

	dbCalls := db.CallCount()
	t.Logf("Broken version - DB calls: %d (expected > 1)", dbCalls)

	if dbCalls <= 1 {
		t.Errorf("Broken version should have > 1 DB calls for stampede, got %d", dbCalls)
	}
	t.Log("PROVEN: Broken version causes stampede - multiple concurrent DB queries without protection")
}

// TestStampedeProtectedVersion demonstrates singleflight protection deterministically.
func TestStampedeProtectedVersion(t *testing.T) {
	db := caching.NewCounterRepository()
	cache := caching.NewMockCache()
	svc := caching.NewProtectedStampedeService(cache, db)
	ctx := context.Background()

	key := int64(2)
	numRequests := 10

	db.Block()

	var wg sync.WaitGroup
	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = svc.GetData(ctx, key)
		}()
	}

	// Singleflight leader enters repository
	db.WaitUntilEntered()

	// Allow waiters to join singleflight queue
	time.Sleep(5 * time.Millisecond)

	// Unblock repository after leader and waiters are ready
	db.Unblock()
	wg.Wait()

	dbCalls := db.CallCount()
	t.Logf("Protected version - DB rebuild count: %d (expected 1)", dbCalls)

	if dbCalls != 1 {
		t.Errorf("expected 1 rebuild with singleflight, got %d", dbCalls)
	}
	t.Log("✓ Single-flight deduplication validated - only 1 DB query for concurrent requests")
}

// TestTTLWithJitter verifies TTL jitter boundary contract.
func TestTTLWithJitter(t *testing.T) {
	baseTTL := 60 * time.Second
	maxJitter := 15 * time.Second

	for i := 0; i < 1000; i++ {
		ttl := caching.TTLWithJitter(baseTTL, maxJitter)
		if ttl < baseTTL || ttl >= baseTTL+maxJitter {
			t.Fatalf("TTL out of range [base, base+maxJitter): %v", ttl)
		}
	}

	// Contract edge cases
	if ttl := caching.TTLWithJitter(baseTTL, 0); ttl != baseTTL {
		t.Fatalf("expected ttl == baseTTL when maxJitter == 0, got %v", ttl)
	}
	if ttl := caching.TTLWithJitter(baseTTL, -5*time.Second); ttl != baseTTL {
		t.Fatalf("expected ttl == baseTTL when maxJitter < 0, got %v", ttl)
	}

	t.Log("✓ TTL jitter boundary contract validated (base <= TTL < base + maxJitter, zero/negative maxJitter returns base)")
}

// TestSingleflightLeaderCancelDoesNotKillRebuild verifies the bounded context fix deterministically.
func TestSingleflightLeaderCancelDoesNotKillRebuild(t *testing.T) {
	repo := caching.NewCounterRepository()
	cache := caching.NewMockCache()
	svc := caching.NewProtectedStampedeService(cache, repo)

	ctxA, cancelA := context.WithCancel(context.Background())
	ctxB := context.Background()

	errA := make(chan error, 1)
	errB := make(chan error, 1)

	repo.Block()

	// A becomes leader (first caller)
	go func() {
		_, err := svc.GetData(ctxA, 200)
		errA <- err
	}()

	// Wait until A/shared rebuild is actually inside repository
	repo.WaitUntilEntered()

	// B joins as waiter (second caller)
	go func() {
		_, err := svc.GetData(ctxB, 200)
		errB <- err
	}()

	// Ensure B has registered in singleflight
	time.Sleep(5 * time.Millisecond)

	// A cancels its context while waiting
	cancelA()

	// A should return context.Canceled immediately
	if err := <-errA; err != context.Canceled {
		t.Errorf("expected Caller A to return context.Canceled, got: %v", err)
	}

	// Unblock repository execution
	repo.Unblock()

	// B should complete successfully
	if err := <-errB; err != nil {
		t.Errorf("expected Caller B to succeed, got: %v", err)
	}

	// DB must be called exactly once
	if repo.CallCount() != 1 {
		t.Errorf("expected 1 DB call, got %d", repo.CallCount())
	}

	// Verify cache is populated
	cacheKey := "dash:200"
	_, err := cache.Get(ctxB, cacheKey)
	if err != nil {
		t.Errorf("cache should be populated after leader cancel, got: %v", err)
	}

	// Third request hits cache -> DB CallCount stays 1
	_, err = svc.GetData(context.Background(), 200)
	if err != nil {
		t.Errorf("third request should hit cache, got: %v", err)
	}
	if repo.CallCount() != 1 {
		t.Errorf("third request should hit cache without DB call, got %d", repo.CallCount())
	}

	t.Log("✓ Deterministic singleflight leader cancellation test passed")
}

// TestBoundedContextIsolatedFromLeaderCancel specifically verifies bounded context isolation.
func TestBoundedContextIsolatedFromLeaderCancel(t *testing.T) {
	repo := caching.NewCounterRepository()
	cache := caching.NewMockCache()
	svc := caching.NewProtectedStampedeService(cache, repo)
	ctx := context.Background()

	ctxA, cancelA := context.WithCancel(ctx)

	doneA := make(chan error, 1)
	doneB := make(chan error, 1)

	repo.Block()

	go func() {
		_, err := svc.GetData(ctxA, 300)
		doneA <- err
	}()

	repo.WaitUntilEntered()

	go func() {
		_, err := svc.GetData(ctx, 300)
		doneB <- err
	}()

	time.Sleep(5 * time.Millisecond)

	cancelA()

	if err := <-doneA; err != context.Canceled {
		t.Errorf("expected A context.Canceled, got: %v", err)
	}

	repo.Unblock()

	if err := <-doneB; err != nil {
		t.Errorf("expected B success, got: %v", err)
	}
	if repo.CallCount() != 1 {
		t.Errorf("expected 1 DB call, got %d", repo.CallCount())
	}

	// Third request should hit cache
	_, err := svc.GetData(ctx, 300)
	if err != nil {
		t.Errorf("third request should succeed from cache, got: %v", err)
	}
	if repo.CallCount() != 1 {
		t.Errorf("cache hit: expected DB count stay 1, got %d", repo.CallCount())
	}

	t.Log("✓ Bounded context isolated from leader cancellation validated")
}
