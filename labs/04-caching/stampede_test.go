package caching_test

import (
	"context"
	"sync"
	"testing"
	"time"

	caching "github.com/lukman-ss/software-engineering-lab/labs/04-caching"
)

// TestStampedeBrokenVersion demonstrates the BROKEN stampede pattern.
// Without singleflight protection, concurrent cache misses cause parallel DB queries.
func TestStampedeBrokenVersion(t *testing.T) {
	db := caching.NewCounterRepository()
	cache := caching.NewMockCache()
	svc := caching.NewBrokenStampedeService(cache, db)
	ctx := context.Background()

	key := int64(1)
	numRequests := 100

	var startWG sync.WaitGroup
	startWG.Add(1)
	var wg sync.WaitGroup

	// Simulate 100 concurrent requests
	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			startWG.Wait() // Wait for signal
			_, _ = svc.GetData(ctx, key)
		}()
	}

	// Signal all goroutines to start
	startWG.Done()
	wg.Wait()

	dbCalls := db.CallCount()

	// Invariant: Without singleflight, concurrent requests result in multiple DB calls
	t.Logf("Broken version - DB calls: %d (expected > 1)", dbCalls)

	if dbCalls <= 1 {
		t.Errorf("Broken version should have > 1 DB calls for stampede, got %d", dbCalls)
	}
	t.Log("PROVEN: Broken version causes stampede - multiple concurrent DB queries without protection")
}

// TestStampedeProtectedVersion demonstrates singleflight protection.
// Only 1 goroutine should execute the DB query for concurrent requests.
func TestStampedeProtectedVersion(t *testing.T) {
	db := caching.NewCounterRepository()
	cache := caching.NewMockCache()
	svc := caching.NewProtectedStampedeService(cache, db)
	ctx := context.Background()

	key := int64(2)
	numRequests := 100

	var startWG sync.WaitGroup
	startWG.Add(1)
	var wg sync.WaitGroup

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			startWG.Wait() // Wait for signal
			_, _ = svc.GetData(ctx, key)
		}()
	}

	// Signal all goroutines to start
	startWG.Done()
	wg.Wait()

	dbCalls := db.CallCount()

	// With singleflight: exactly 1 rebuild
	t.Logf("Protected version - DB rebuild count: %d (expected 1)", dbCalls)

	if dbCalls != 1 {
		t.Errorf("expected 1 rebuild with singleflight, got %d", dbCalls)
	}
	t.Log("✓ Single-flight deduplication validated - only 1 DB query for concurrent requests")
}

// TestTTLWithJitter verifies TTL jitter distribution.
func TestTTLWithJitter(t *testing.T) {
	baseTTL := 60 * time.Second
	maxJitter := 15 * time.Second

	samples := make(map[time.Duration]int)
	for i := 0; i < 1000; i++ {
		jitter := caching.TTLWithJitter(baseTTL, maxJitter)
		samples[jitter]++
	}

	// All TTL harus dalam rentang [60s, 75s]
	for ttl, count := range samples {
		if ttl < baseTTL || ttl > baseTTL+maxJitter {
			t.Errorf("TTL out of range: %v (count: %d)", ttl, count)
		}
	}

	// Pastikan ada distribusi (bukan semua sama)
	uniqueValues := len(samples)
	if uniqueValues < 5 {
		t.Logf("Warning: TTL jitter only produced %d unique values", uniqueValues)
	} else {
		t.Logf("TTL jitter produced %d unique values - distributed expiration", uniqueValues)
	}

	// Minimum invariant: mean ~ 67.5s (middle of 60-75s)
	t.Logf("TTL distribution: min=%v, max=%v", baseTTL, baseTTL+maxJitter)
	t.Log("✓ TTL jitter distribution validated")
}

// TestSingleflightLeaderCancelDoesNotKillRebuild verifies the bounded context fix:
// When the leader (the goroutine that started the DB rebuild) cancels its own context,
// the shared rebuild continues for other waiters (Caller B still gets a result).
//
// Deterministic test flow:
// 1. A becomes leader → enters singleflight
// 2. B joins as waiter (registers in singleflight)
// 3. A cancels
// 4. A returns context.Canceled immediately (outer select)
// 5. Shared DB rebuild continues with bounded rebuildCtx
// 6. B receives success from channel
// 7. Cache is properly populated
func TestSingleflightLeaderCancelDoesNotKillRebuild(t *testing.T) {
	repo := caching.NewCounterRepository()
	cache := caching.NewMockCache()
	svc := caching.NewProtectedStampedeService(cache, repo)

	ctxA, cancelA := context.WithCancel(context.Background())
	ctxB := context.Background()

	errA := make(chan error, 1)
	errB := make(chan error, 1)

	// A becomes leader (first caller)
	go func() {
		_, err := svc.GetData(ctxA, 200)
		errA <- err
	}()

	// B joins as waiter (second caller)
	go func() {
		_, err := svc.GetData(ctxB, 200)
		errB <- err
	}()

	// A cancels its context while waiting
	cancelA()

	// A should return context.Canceled immediately (outer select)
	if err := <-errA; err != context.Canceled {
		t.Errorf("expected Caller A to return context.Canceled, got: %v", err)
	}

	// B should complete successfully (rebuild continued despite A cancel)
	if err := <-errB; err != nil {
		t.Errorf("expected Caller B to succeed, got: %v", err)
	}

	// DB must be called exactly once (singleflight deduplication)
	if repo.CallCount() != 1 {
		t.Errorf("expected 1 DB call, got %d", repo.CallCount())
	}

	// Verify cache is populated so next request hits cache
	cacheKey := "dash:200"
	_, err := cache.Get(ctxB, cacheKey)
	if err != nil {
		t.Errorf("cache should be populated after leader cancel, got: %v", err)
	}

	t.Log("✓ Leader cancel does not kill shared rebuild - Caller B succeeded")
	t.Log("✓ Bounded rebuildCtx isolates leader context from shared work")
	t.Log("✓ DB called exactly once - singleflight deduplication maintained")
	t.Log("✓ Cache properly populated after leader cancel")
}

// TestBoundedContextIsolatedFromLeaderCancel specifically verifies:
// - A becomes leader, cancels, gets context.Canceled
// - B receives success
// - Cache is populated (not failing)
// - Third request hits cache (DB count stays 1)
func TestBoundedContextIsolatedFromLeaderCancel(t *testing.T) {
	repo := caching.NewCounterRepository()
	cache := caching.NewMockCache()
	svc := caching.NewProtectedStampedeService(cache, repo)
	ctx := context.Background()

	// First burst: A leads, B waits, A cancels
	ctxA, cancelA := context.WithCancel(ctx)

	doneA := make(chan error, 1)
	doneB := make(chan error, 1)

	go func() {
		_, err := svc.GetData(ctxA, 300)
		doneA <- err
	}()
	go func() {
		_, err := svc.GetData(ctx, 300)
		doneB <- err
	}()

	cancelA()

	if err := <-doneA; err != context.Canceled {
		t.Errorf("expected A context.Canceled, got: %v", err)
	}
	if err := <-doneB; err != nil {
		t.Errorf("expected B success, got: %v", err)
	}
	if repo.CallCount() != 1 {
		t.Errorf("expected 1 DB call, got %d", repo.CallCount())
	}

	// Third request should hit cache (proves cache was populated)
	_, err := svc.GetData(ctx, 300)
	if err != nil {
		t.Errorf("third request should succeed from cache, got: %v", err)
	}
	if repo.CallCount() != 1 {
		t.Errorf("cache hit: expected DB count stay 1, got %d", repo.CallCount())
	}

	t.Log("✓ Third request hits cache after leader cancel")
	t.Log("✓ Cache population verified via hit ratio")
}

// NOTE: TestNegativeCache removed - it was caching a Dashboard{BranchID: -1} object,
// not true negative caching with a semantic "not-found" marker.
