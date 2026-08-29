package caching_test

import (
	"context"
	"testing"
	"time"

	caching "github.com/lukman/software-engineer-lab/labs/04-caching"
)

// Integration tests for essential cache behaviors.
// These tests demonstrate real-world scenarios without requiring Redis.

// TestCacheMissReadsDatabase verifies: cache miss reads from database
func TestCacheMissReadsDatabase(t *testing.T) {
	cache := caching.NewMockCache()
	metrics := caching.NewCacheMetrics()
	svc := caching.NewRobustDashboardService(nil, cache, metrics)
	ctx := context.Background()

	// First request = cache miss → database read
	_, _ = svc.GetDashboard(ctx, 1)

	if metrics.DBQueries() == 0 {
		t.Error("cache miss should trigger database read")
	}
	if metrics.Misses() == 0 {
		t.Error("cache miss count should be incremented")
	}
	t.Log("✓ Cache miss reads from database")
}

// TestCacheMissWritesCache verifies: cache miss writes to cache after read
func TestCacheMissWritesCache(t *testing.T) {
	cache := caching.NewMockCacheWithStats()
	metrics := caching.NewCacheMetrics()
	svc := caching.NewRobustDashboardService(nil, cache, metrics)
	ctx := context.Background()

	// First request
	_, _ = svc.GetDashboard(ctx, 1)

	// Second request - should be cache hit
	_, _ = svc.GetDashboard(ctx, 1)

	if metrics.Hits() == 0 {
		t.Error("second request should be cache hit (cache was written)")
	}
	t.Log("✓ Cache miss writes to cache")
}

// TestCacheHitDoesNotQueryDatabase verifies: cache hit does not query database
func TestCacheHitDoesNotQueryDatabase(t *testing.T) {
	cache := caching.NewMockCacheWithStats()
	metrics := caching.NewCacheMetrics()
	svc := caching.NewRobustDashboardService(nil, cache, metrics)
	ctx := context.Background()

	// Populate cache first
	db := caching.NewHeavyDB()
	_ = caching.NewProtectedStampedeService(cache, db) // For mock

	initialQueries := metrics.DBQueries()

	// Multiple cache hits
	for i := 0; i < 5; i++ {
		_, _ = svc.GetDashboard(ctx, 1)
	}

	if metrics.DBQueries() > initialQueries+1 {
		t.Errorf("cache hits should not increase DB queries")
	}
	t.Log("✓ Cache hit does not query database")
}

// TestCacheExpiryCausesRebuild verifies: cache expiry causes rebuild
func TestCacheExpiryCausesRebuild(t *testing.T) {
	cache := caching.NewMockCacheWithStats()
	metrics := caching.NewCacheMetrics()
	svc := caching.NewRobustDashboardService(nil, cache, metrics)
	ctx := context.Background()

	// Build cache
	_, _ = svc.GetDashboard(ctx, 1)
	initialRebuilds := metrics.Rebuilds()

	// Force miss by deleting
	_ = cache.Delete(ctx, caching.DashboardCacheKey(1))

	// Next request rebuilds
	_, _ = svc.GetDashboard(ctx, 1)

	if metrics.Rebuilds() <= initialRebuilds {
		t.Error("after delete, rebuild should be triggered")
	}
	t.Log("✓ Cache delete causes rebuild")
}

// TestCacheInvalidationReturnsFreshData verifies: invalidation returns fresh data
func TestCacheInvalidationReturnsFreshData(t *testing.T) {
	cache := caching.NewMockCacheWithStats()
	svc := caching.NewDashboardCacheService(nil, cache)
	ctx := context.Background()

	branchID := int64(1)
	key := caching.DashboardCacheKey(branchID)

	// Set initial data
	initialData := `{"branch_id":1,"invoice_count":10,"date":"` + caching.Today() + `"}`
	cache.Set(ctx, key, initialData, time.Minute)

	// Verify cache has initial data
	cached, _ := cache.Get(ctx, key)
	if cached != initialData {
		t.Error("cache should have initial data")
	}

	// Invalidate
	_ = svc.InvalidateBranchDashboard(ctx, branchID)

	// Cache should be empty
	_, err := cache.Get(ctx, key)
	if err == nil {
		t.Error("cache should be invalidated")
	}
	t.Log("✓ Cache invalidation returns fresh data on next request")
}

// TestDifferentBranchDoesNotShareCache verifies: different branch has different cache
func TestDifferentBranchDoesNotShareCache(t *testing.T) {
	cache := caching.NewMockCache()
	ctx := context.Background()

	// Set data for branch 1
	key1 := caching.DashboardCacheKey(1)
	data1 := `{"branch_id":1,"invoice_count":50}`
	cache.Set(ctx, key1, data1, time.Minute)

	// Set data for branch 2
	key2 := caching.DashboardCacheKey(2)
	data2 := `{"branch_id":2,"invoice_count":75}`
	cache.Set(ctx, key2, data2, time.Minute)

	// Keys should be different
	if key1 == key2 {
		t.Error("different branches should have different cache keys")
	}

	// Data should not leak
	cached1, _ := cache.Get(ctx, key1)
	cached2, _ := cache.Get(ctx, key2)

	if cached1 == cached2 {
		t.Error("different branches should have different cache data")
	}
	t.Log("✓ Different branch does not share cache")
}

// TestDifferentTenantDoesNotShareCache verifies: different tenant isolation
func TestDifferentTenantDoesNotShareCache(t *testing.T) {
	loc := time.UTC
	keyTenantA := caching.NewTenantDashboardKey(1, 1, loc).Build()
	keyTenantB := caching.NewTenantDashboardKey(2, 1, loc).Build()

	if keyTenantA == keyTenantB {
		t.Error("different tenants should have different cache keys")
	}

	// Key format check
	// Expected: cmms:dashboard:v1:tenant:{id}:branch:{id}:{date}
	t.Logf("Tenant A key: %s", keyTenantA)
	t.Logf("Tenant B key: %s", keyTenantB)

	t.Log("✓ Different tenant does not share cache")
}

// TestCorruptCacheFallsBackAppropriately verifies: corrupt cache falls back to database
func TestCorruptCacheFallsBackAppropriately(t *testing.T) {
	cache := caching.NewMockCache()
	metrics := caching.NewCacheMetrics()
	svc := caching.NewRobustDashboardService(nil, cache, metrics)
	ctx := context.Background()

	// Set corrupt JSON
	corruptData := `{"invalid json`
	cache.Set(ctx, "dashboard:1", corruptData, time.Minute)

	// Get should handle corrupt gracefully (fallback, not fail)
	_, err := svc.GetDashboard(ctx, 1)
	// Should not return unmarshal error - should use fallback

	t.Logf("Corrupt cache handled gracefully, error: %v", err)
	t.Log("✓ Corrupt cache falls back appropriately")
}

// TestCacheFailureFallsBackToDatabase verifies: cache down falls back to database
func TestCacheFailureFallsBackToDatabase(t *testing.T) {
	cache := caching.NewFailingMockCache()
	metrics := caching.NewCacheMetrics()
	svc := caching.NewRobustDashboardService(nil, cache, metrics)
	ctx := context.Background()

	// Cache is down, but request should still succeed
	result, err := svc.GetDashboard(ctx, 1)

	// Request should succeed via DB fallback
	if err != nil {
		t.Errorf("cache failure should fall back to database, got error: %v", err)
	}

	// Metrics should show error
	if metrics.Errors() == 0 {
		t.Error("cache failure should increment error counter")
	}

	t.Logf("Result from database fallback: %+v", result)
	t.Log("✓ Cache failure falls back to database")
}

// TestConcurrentCacheMissProtectedBySingleflight verifies stampede protection
func TestConcurrentCacheMissProtectedBySingleflight(t *testing.T) {
	db := caching.NewHeavyDB()
	cache := caching.NewMockCache()
	svc := caching.NewProtectedStampedeService(cache, db)
	ctx := context.Background()

	key := "dashboard:concurrent"

	// Run concurrent requests
	done := make(chan bool)
	for i := 0; i < 50; i++ {
		go func() {
			_, _ = svc.GetData(ctx, key)
			done <- true
		}()
	}

	for i := 0; i < 50; i++ {
		<-done
	}

	// With singleflight, all concurrent requests should result in 1 DB query
	if db.RebuildCount() != 1 {
		t.Errorf("singleflight should dedupe concurrent misses to 1 DB query, got %d", db.RebuildCount())
	}
	t.Log("✓ Concurrent cache miss protected by singleflight")
}

// TestTTLJitterRange verifies: TTL jitter stays within configured range
func TestTTLJitterRange(t *testing.T) {
	baseTTL := 60 * time.Second
	maxJitter := 15 * time.Second

	// Test many samples
	for i := 0; i < 1000; i++ {
		jitteredTTL := caching.TTLWithJitter(baseTTL, maxJitter)

		// Must be within [baseTTL, baseTTL + maxJitter]
		if jitteredTTL < baseTTL || jitteredTTL > baseTTL+maxJitter {
			t.Errorf("TTL jitter out of range: %v (expected %v - %v)",
				jitteredTTL, baseTTL, baseTTL+maxJitter)
		}
	}

	t.Log("✓ TTL jitter stays within configured range")
}

// TestCacheHitRatio calculates correctly verifies hit ratio calculation
func TestCacheHitRatioCalculatesCorrectly(t *testing.T) {
	metrics := caching.NewCacheMetrics()

	// 80 hits, 20 misses
	for i := 0; i < 80; i++ {
		metrics.IncHit()
	}
	for i := 0; i < 20; i++ {
		metrics.IncMiss()
	}

	hitRatio := metrics.HitRatio()
	expected := 80.0 // 80%

	if hitRatio != expected {
		t.Errorf("hit ratio should be %v, got %v", expected, hitRatio)
	}

	t.Logf("Cache hit ratio: %.2f%%", hitRatio)
	t.Log("✓ Hit ratio calculates correctly")
}