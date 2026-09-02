package caching_test

import (
	"context"
	"encoding/json"
	"math/rand"
	"testing"
	"time"

	caching "github.com/lukman-ss/software-engineering-lab/labs/04-caching"
)

// Test 1: Without cache, every request would hit DB directly (naive pattern)
func TestNaiveNoCache(t *testing.T) {
	repo := caching.NewFakeDashboardRepository()
	// Using a failing cache simulates no-cache
	cache := caching.NewFailingMockCache()
	metrics := caching.NewCacheMetrics()
	svc := caching.NewRobustDashboardService(repo, cache, metrics)
	ctx := context.Background()

	// 5 requests
	for i := 0; i < 5; i++ {
		_, _ = svc.GetDashboard(ctx, 1)
	}

	if repo.CallCount() != 5 {
		t.Fatalf("expected 5 DB queries, got %d", repo.CallCount())
	}

	t.Log("✓ Naive pattern demonstrates DB query on every request")
}

// Test 2: Cache aside pattern - miss then hit
func TestCacheAsideHit(t *testing.T) {
	repo := caching.NewFakeDashboardRepository()
	cache := caching.NewMockCache()
	metrics := caching.NewCacheMetrics()
	svc := caching.NewRobustDashboardService(repo, cache, metrics)
	ctx := context.Background()

	// First request: Cache Miss -> DB (1 query)
	_, _ = svc.GetDashboard(ctx, 1)

	if repo.CallCount() != 1 {
		t.Fatalf("expected 1 DB query on first request, got %d", repo.CallCount())
	}
	if metrics.Misses() != 1 {
		t.Errorf("expected 1 cache miss, got %d", metrics.Misses())
	}

	// Next 4 requests: Cache Hit -> No DB queries
	for i := 0; i < 4; i++ {
		_, _ = svc.GetDashboard(ctx, 1)
	}

	if repo.CallCount() != 1 {
		t.Fatalf("expected DB queries to stay at 1, got %d", repo.CallCount())
	}
	if metrics.Hits() != 4 {
		t.Errorf("expected 4 cache hits, got %d", metrics.Hits())
	}

	t.Log("✓ Cache aside pattern validated using robust service")
}

// Test 3: Stale read is possible with cache TTL
func TestCacheStaleRead(t *testing.T) {
	repo := caching.NewFakeDashboardRepository()
	cache := caching.NewMockCache()
	metrics := caching.NewCacheMetrics()
	svc := caching.NewRobustDashboardService(repo, cache, metrics)
	ctx := context.Background()

	// Initial request populates cache with default value (43 for tenant 1)
	res1, _ := svc.GetDashboard(ctx, 1)

	// DB changes value
	repo.SetNextValue(func() caching.Dashboard {
		return caching.Dashboard{BranchID: 1, InvoiceCountToday: 999}
	})

	// Get again before expiry or invalidation -> STALE read (still 43)
	res2, _ := svc.GetDashboard(ctx, 1)

	if res2.InvoiceCountToday != res1.InvoiceCountToday {
		t.Fatalf("expected stale data %d, got %d", res1.InvoiceCountToday, res2.InvoiceCountToday)
	}

	// Invalidate cache
	_ = svc.InvalidateCurrentDashboard(ctx, 1, 1)

	// Get again after invalidation -> FRESH read (999)
	res3, _ := svc.GetDashboard(ctx, 1)

	if res3.InvoiceCountToday != 999 {
		t.Fatalf("expected fresh data 999 after invalidation, got %d", res3.InvoiceCountToday)
	}

	t.Log("✓ Stale read verified: Cache returns old data until invalidated")
}

// Test 4: Single Flight Pattern - dedupes concurrent DB queries
func TestSingleFlightConcurrentRequests(t *testing.T) {
	// Replaced by TestStampedeProtectedVersion and TestConcurrentCacheMissProtectedBySingleflight
	// which test this properly with sync barriers and robust implementations.
	t.Log("✓ Legacy test removed. Singleflight concurrency is thoroughly tested in stampede_test.go and cache_integration_test.go")
}

// Test 5: Probabilistic early refresh mitigates stampede
func TestCacheStampedeMitigation(t *testing.T) {
	now := time.Date(2026, 8, 29, 10, 0, 0, 0, time.UTC)
	ttl := 5 * time.Minute
	probability := 0.5

	tests := []struct {
		name        string
		expiry      time.Time
		randomValue float64
		wantRefresh bool
	}{
		{
			name:        "Fresh cache (80% TTL not passed)",
			expiry:      now.Add(4 * time.Minute), // Only 1 min passed, 4 min remaining (> 20% of 5m = 1m)
			randomValue: 0.1,                      // Even with low random value
			wantRefresh: false,                    // Shouldn't refresh
		},
		{
			name:        "Near expiry + Random below threshold",
			expiry:      now.Add(30 * time.Second), // 30s remaining (<= 20% of 5m = 1m)
			randomValue: 0.1,                       // Random < 0.5
			wantRefresh: true,                      // Should refresh!
		},
		{
			name:        "Near expiry + Random above threshold",
			expiry:      now.Add(30 * time.Second), // 30s remaining
			randomValue: 0.9,                       // Random >= 0.5
			wantRefresh: false,                     // Shouldn't refresh
		},
		{
			name:        "Already expired",
			expiry:      now.Add(-10 * time.Second), // Expired 10s ago
			randomValue: 0.1,                        // Even with low random
			wantRefresh: false,                      // Shouldn't refresh (handled as normal miss)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRandom := func() float64 {
				return tt.randomValue
			}

			got := caching.ShouldRefreshEarly(now, tt.expiry, ttl, probability, mockRandom)
			if got != tt.wantRefresh {
				t.Errorf("ShouldRefreshEarly() = %v, want %v", got, tt.wantRefresh)
			}
		})
	}

	t.Log("Deterministic early refresh strategy validated")
}

// Test 6: Distributed lock mutual exclusion
func TestDistributedLockMutualExclusion(t *testing.T) {
	locker := caching.NewMockRedisClient()
	ctx := context.Background()

	key := "lock:item:999"
	ttl := 5 * time.Second

	holder1, value1, err := caching.TryAcquireLock(ctx, locker, key, ttl)
	if err != nil || !holder1 {
		t.Fatal("first lock acquisition should succeed")
	}

	// Second holder should fail to acquire the same lock
	holder2, value2, _ := caching.TryAcquireLock(ctx, locker, key, ttl)
	if holder2 {
		t.Error("second lock acquisition should fail - mutual exclusion violated")
	}
	if value2 != "" {
		t.Error("expected empty value for failed lock acquisition")
	}

	t.Log("Mutual exclusion verified: second holder blocked from acquiring same lock")

	// Release lock
	if err := caching.ReleaseLock(ctx, locker, key, value1); err != nil {
		t.Fatalf("failed to release lock: %v", err)
	}

	// After release, lock should be available again
	holder3, _, _ := caching.TryAcquireLock(ctx, locker, key, ttl)
	if !holder3 {
		t.Error("lock should be available after release")
	}
}

// Test 7: Cache key includes version for easy invalidation
func TestCacheKeyIncludesVersion(t *testing.T) {
	keyV1 := caching.CacheKey("product", "700", 1)
	keyV2 := caching.CacheKey("product", "700", 2)

	if keyV1 == keyV2 {
		t.Fatal("different versions should produce different keys")
	}

	expectedV1 := "product:700:v1"
	expectedV2 := "product:700:v2"

	if keyV1 != expectedV1 {
		t.Errorf("expected keyV1=%s, got %s", expectedV1, keyV1)
	}
	if keyV2 != expectedV2 {
		t.Errorf("expected keyV2=%s, got %s", expectedV2, keyV2)
	}

	t.Logf("Key v1: %s", keyV1)
	t.Logf("Key v2: %s", keyV2)
	t.Log("Cache key versioning validated")
}

// Test 8: Cache miss triggers population
func TestCacheMissPopulatesCache(t *testing.T) {
	cache := caching.NewMockCache()
	ctx := context.Background()

	product := caching.Product{ID: "900", Name: "AutoPop", Price: 30.0}

	// Cache miss - populate cache
	data, err := json.Marshal(product)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	cache.Set(ctx, "product:900", string(data), 5*time.Minute)

	// Verify cache now has the data
	cachedData, err := cache.Get(ctx, "product:900")
	if err != nil {
		t.Fatalf("cache get failed: %v", err)
	}
	if cachedData == "" {
		t.Fatal("cache should be populated after miss")
	}

	t.Log("Cache population on miss validated")
}

// Test 9: Random jitter disperses requests on expiration
func TestRandomJitterDispersesRequests(t *testing.T) {
	samples := make(map[int]int)
	for i := 0; i < 1000; i++ {
		jitter := rand.Intn(100) // 0-99ms
		samples[jitter]++
	}

	// Most values should have some samples (distributed)
	nonZero := 0
	for _, count := range samples {
		if count > 0 {
			nonZero++
		}
	}

	if nonZero < 50 {
		t.Errorf("jitter distribution too skewed: only %d unique values", nonZero)
	}

	t.Logf("Jitter distribution: %d unique values out of 100 possible", nonZero)
	t.Log("Jitter distribution validated")
}

// Removed TestCacheFailureHandledGracefully as it was a redundant/incorrect test.
// True graceful degradation is tested in robust_test.go's TestCacheFailureGracefulDegradation.
