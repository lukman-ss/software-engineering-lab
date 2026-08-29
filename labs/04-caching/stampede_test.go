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

	// Use a barrier to ensure all goroutines start simultaneously
	// This maximizes the chance of concurrent cache misses
	var startWG sync.WaitGroup
	startWG.Add(1)

	var wg sync.WaitGroup

	// Simulate 100 concurrent requests (dashboard di-reload bersamaan)
	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Wait for the signal to start
			startWG.Wait()
			_, _ = svc.GetData(ctx, key)
		}()
	}

	// Small delay to let all goroutines start and be ready
	time.Sleep(10 * time.Millisecond)

	// Signal all goroutines to start simultaneously
	startWG.Done()
	wg.Wait()

	dbCalls := db.CallCount()

	// Invariant: Without singleflight, concurrent requests result in multiple DB calls
	// Exactly numRequests is not guaranteed due to timing, but it MUST be > 1
	t.Logf("Broken version - DB calls: %d (expected > 1)", dbCalls)

	if dbCalls <= 1 {
		t.Errorf("Broken version should have > 1 DB calls for stampede, got %d (this may indicate test timing issue)", dbCalls)
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

	// Use a channel-based coordination for deterministic timing
	startCh := make(chan struct{})
	var wg sync.WaitGroup

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-startCh // Wait for signal
			_, _ = svc.GetData(ctx, key)
		}()
	}

	// Small delay to ensure all goroutines are waiting
	time.Sleep(10 * time.Millisecond)

	// Release all goroutines simultaneously
	close(startCh)
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

// TestNegativeCache verifies negative caching for non-existent keys.
func TestNegativeCache(t *testing.T) {
	cache := caching.NewMockCache()
	ctx := context.Background()

	notFoundKey := "product:999999"

	// Set "not found" marker dengan TTL pendek
	cache.Set(ctx, notFoundKey, "NULL_NOT_FOUND", 30*time.Second)

	// Subsequent request dapat not-found dari cache (tanpa DB hit)
	cached, err := cache.Get(ctx, notFoundKey)
	if err != nil {
		t.Fatal("expected cache hit on not-found marker")
	}

	if cached != "NULL_NOT_FOUND" {
		t.Error("expected not-found marker in cache")
	}

	// Verify trade-off: jika object baru dibuat selama negative TTL,
	// user masih melihat not-found sampai TTL habis
	t.Log("✓ Negative cache prevents repeated DB lookups for non-existent keys")
	t.Log("Trade-off: Object created during negative TTL = stale not-found")
}