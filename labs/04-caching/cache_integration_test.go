package caching_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	caching "github.com/lukman-ss/software-engineering-lab/labs/04-caching"
)

// Integration tests for essential cache behaviors.
// These tests demonstrate real-world scenarios without requiring Redis.

// TestCacheMissReadsDatabase verifies: cache miss reads from repository
func TestCacheMissReadsDatabase(t *testing.T) {
	cache := caching.NewMockCache()
	metrics := caching.NewCacheMetrics()
	repo := caching.NewFakeDashboardRepository()
	svc := caching.NewRobustDashboardService(repo, cache, metrics)
	ctx := context.Background()

	// First request = cache miss → database read
	_, _ = svc.GetDashboard(ctx, 1)

	if repo.CallCount == 0 {
		t.Error("cache miss should trigger repository call")
	}
	if metrics.Misses() == 0 {
		t.Error("cache miss count should be incremented")
	}
	t.Log("✓ Cache miss reads from repository")
	t.Logf("✓ Repository CallCount = %d (expected 1)", repo.CallCount)
}

// TestCacheMissWritesCache verifies: cache miss writes to cache after read
func TestCacheMissWritesCache(t *testing.T) {
	cache := caching.NewMockCacheWithStats()
	metrics := caching.NewCacheMetrics()
	repo := caching.NewFakeDashboardRepository()
	svc := caching.NewRobustDashboardService(repo, cache, metrics)
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
	ctx := context.Background()
	cache := caching.NewMockCacheWithStats()
	metrics := caching.NewCacheMetrics()
	repo := caching.NewFakeDashboardRepository()
	svc := caching.NewRobustDashboardService(repo, cache, metrics)

	// Populate cache first
	_, _ = svc.GetDashboard(ctx, 1)

	initialCalls := repo.CallCount

	// Multiple cache hits
	for i := 0; i < 5; i++ {
		_, _ = svc.GetDashboard(ctx, 1)
	}

	if repo.CallCount > initialCalls {
		t.Errorf("cache hits should not increase repository calls")
	}
	t.Log("✓ Cache hit does not query database")
}

// TestCacheExpiryCausesRebuild verifies: cache expiry causes rebuild
func TestCacheExpiryCausesRebuild(t *testing.T) {
	cache := caching.NewMockCache()
	metrics := caching.NewCacheMetrics()
	repo := caching.NewFakeDashboardRepository()
	svc := caching.NewRobustDashboardService(repo, cache, metrics)
	ctx := context.Background()

	// Build cache
	_, _ = svc.GetDashboard(ctx, 1)
	initialRebuilds := metrics.Rebuilds()

	// Force miss by deleting
	_ = cache.Delete(ctx, caching.DashboardCacheKey(1, 1, time.Now().UTC()))

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
	repo := caching.NewFakeDashboardRepository()
	svc := caching.NewDashboardCacheService(repo, cache)
	ctx := context.Background()

	branchID := int64(1)
	today := time.Now().UTC()
	key := caching.DashboardCacheKey(1, branchID, today)

	// Step 1: Populate cache with initial data
	repo.CallCount = 0
	cache.Set(ctx, key, `{"branch_id":1,"invoice_count":10,"date":"`+caching.Today()+`"}`, time.Minute)
	repo.CallCount = 1 // First request already counted

	// Verify cache has initial data
	cached, _ := cache.Get(ctx, key)
	if cached == "" {
		t.Error("cache should have initial data")
	}

	// Step 2: Invalidate
	_ = svc.InvalidateCurrentDashboard(ctx, 1, branchID)

	// Step 3: Next request should query repo and update cache
	repo.CallCount = 0
	_, _ = svc.GetDashboard(ctx, 1)

	if repo.CallCount != 1 {
		t.Errorf("expected 1 repository call after invalidation, got %d", repo.CallCount)
	}
	t.Log("✓ Cache invalidation returns fresh data on next request")
}

// TestDifferentBranchDoesNotShareCache verifies: different branch has different cache
func TestDifferentBranchDoesNotShareCache(t *testing.T) {
	cache := caching.NewMockCache()
	ctx := context.Background()

	// Set data for branch 1
	key1 := caching.DashboardCacheKey(1, 1, time.Now().UTC())
	data1 := `{"branch_id":1,"invoice_count":50}`
	cache.Set(ctx, key1, data1, time.Minute)

	// Set data for branch 2
	key2 := caching.DashboardCacheKey(1, 2, time.Now().UTC())
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
	now := time.Now()
	keyTenantA := caching.DashboardCacheKey(1, 1, now)
	keyTenantB := caching.DashboardCacheKey(2, 1, now)

	if keyTenantA == keyTenantB {
		t.Error("different tenants should have different cache keys")
	}

	// Key format check
	// Expected: cmms:dashboard:v1:tenant:{id}:branch:{id}:{date}
	t.Logf("Tenant A key: %s", keyTenantA)
	t.Logf("Tenant B key: %s", keyTenantB)

	t.Log("✓ Different tenant does not share cache")
}

// TestTenantIsolationInService verifies: Tenant A Branch 1 != Tenant B Branch 1
func TestTenantIsolationInService(t *testing.T) {
	ctx := context.Background()
	cache := caching.NewMockCache()
	metrics := caching.NewCacheMetrics()
	repo := caching.NewFakeDashboardRepository()
	svc := caching.NewRobustDashboardService(repo, cache, metrics)

	today := time.Now().UTC().Truncate(24 * time.Hour)

	// Tenant A, Branch 1 - set data
	keyA1 := caching.NewDashboardKey(1).WithTenant(1).WithDate(today).Build()
	dataA1 := caching.Dashboard{BranchID: 1, InvoiceCountToday: 100, Date: today.Format("2006-01-02")}
	dataA1Bytes, _ := json.Marshal(dataA1)
	cache.Set(ctx, keyA1, string(dataA1Bytes), time.Minute)

	// Tenant B, Branch 1 - set different data
	keyB1 := caching.NewDashboardKey(1).WithTenant(2).WithDate(today).Build()
	dataB1 := caching.Dashboard{BranchID: 1, InvoiceCountToday: 200, Date: today.Format("2006-01-02")}
	dataB1Bytes, _ := json.Marshal(dataB1)
	cache.Set(ctx, keyB1, string(dataB1Bytes), time.Minute)

	// Tenant A Branch 1 should get 100
	resultA1, err := svc.GetDashboardWithTenant(ctx, 1, 1, today)
	if err != nil {
		t.Fatalf("Tenant A get failed: %v", err)
	}
	if resultA1.InvoiceCountToday != 100 {
		t.Errorf("Tenant A Branch 1: expected 100, got %d", resultA1.InvoiceCountToday)
	}

	// Tenant B Branch 1 should get 200
	resultB1, err := svc.GetDashboardWithTenant(ctx, 2, 1, today)
	if err != nil {
		t.Fatalf("Tenant B get failed: %v", err)
	}
	if resultB1.InvoiceCountToday != 200 {
		t.Errorf("Tenant B Branch 1: expected 200, got %d", resultB1.InvoiceCountToday)
	}

	t.Log("✓ Tenant isolation verified: Tenant A Branch 1 ≠ Tenant B Branch 1")
}

// TestCorruptCacheFallsBackAppropriately verifies: corrupt cache falls back to database
// This test puts corrupt JSON in the exact key that will be read by the service.
func TestCorruptCacheFallsBackAppropriately(t *testing.T) {
	cache := caching.NewMockCache()
	repo := caching.NewFakeDashboardRepository()
	metrics := caching.NewCacheMetrics()
	svc := caching.NewRobustDashboardService(repo, cache, metrics)
	ctx := context.Background()

	branchID := int64(1)
	today := time.Now().UTC()
	key := caching.DashboardCacheKey(1, branchID, today)

	// Set corrupt JSON in the exact key that the service will read
	corruptData := `{"id":"123"` // incomplete JSON - corrupt
	cache.Set(ctx, key, corruptData, time.Minute)

	// GetDashboard will:
	// 1. Read from cache
	// 2. Fail to unmarshal
	// 3. Delete corrupt entry, fallback to repo, rebuild cache
	repo.CallCount = 0
	result, err := svc.GetDashboard(ctx, branchID)

	// Should succeed via fallback
	if err != nil {
		t.Fatalf("expected fallback to succeed, got error: %v", err)
	}

	// Repository should have been called once to rebuild
	if repo.CallCount != 1 {
		t.Errorf("expected 1 repository call after corrupt cache bypass, got %d", repo.CallCount)
	}

	// Cache hit should now work
	cached, err := cache.Get(ctx, key)
	if err != nil {
		t.Error("cache should now have valid data after rebuild")
	}

	// Verify the cached data is valid JSON
	var d caching.Dashboard
	if err := json.Unmarshal([]byte(cached), &d); err != nil {
		t.Error("cache should contain valid JSON after rebuild")
	}

	t.Logf("Result from database fallback: %+v", result)
	t.Log("✓ Corrupt cache falls back appropriately")
}

// TestCacheFailureFallsBackToDatabase verifies: cache down falls back to database
func TestCacheFailureFallsBackToDatabase(t *testing.T) {
	cache := caching.NewFailingMockCache()
	metrics := caching.NewCacheMetrics()
	repo := caching.NewFakeDashboardRepository()
	svc := caching.NewRobustDashboardService(repo, cache, metrics)
	ctx := context.Background()

	// Cache is down, but request should still succeed
	result, err := svc.GetDashboard(ctx, 1)

	// Request should succeed via DB fallback
	if err != nil {
		t.Errorf("cache failure should fall back to database, got error: %v", err)
	}

	// Repository should have been called
	if repo.CallCount == 0 {
		t.Error("repository should have been called for fallback")
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
	cache := caching.NewMockCache()
	repo := caching.NewCounterRepository()
	svc := caching.NewProtectedStampedeService(cache, repo)
	ctx := context.Background()

	// Run concurrent requests
	done := make(chan bool)
	for i := 0; i < 50; i++ {
		go func() {
			_, _ = svc.GetData(ctx, 1)
			done <- true
		}()
	}

	for i := 0; i < 50; i++ {
		<-done
	}

	// With singleflight, all concurrent requests should result in 1 DB query
	if repo.CallCount() != 1 {
		t.Errorf("singleflight should dedupe concurrent misses to 1 DB query, got %d", repo.CallCount())
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

// TestHistoricalDateInvariant verifies that historical date requests use exact dates
// throughout the whole pipeline (key, repo, result)
func TestHistoricalDateInvariant(t *testing.T) {
	cache := caching.NewMockCache()
	metrics := caching.NewCacheMetrics()
	repo := caching.NewFakeDashboardRepository()
	svc := caching.NewRobustDashboardService(repo, cache, metrics)
	ctx := context.Background()

	// Given a specific historical business date
	historicalDate := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	tenantID, branchID := int64(1), int64(42)

	// When dashboard is requested
	result, err := svc.GetDashboardWithTenant(ctx, tenantID, branchID, historicalDate)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Then the result date must match the requested historical date
	expectedDateStr := "2026-08-01"
	if result.Date != expectedDateStr {
		t.Errorf("Result date is %s, expected %s", result.Date, expectedDateStr)
	}

	// And cache key must have exactly that date
	expectedKey := caching.DashboardCacheKey(tenantID, branchID, historicalDate)
	cached, err := cache.Get(ctx, expectedKey)
	if err != nil || cached == "" {
		t.Errorf("Cache missing for historical date key: %s", expectedKey)
	}

	t.Log("✓ Historical date invariant validated")
}

// TestMultiTenantBehavioralIsolation verifies exact behavior across tenants
func TestMultiTenantBehavioralIsolation(t *testing.T) {
	cache := caching.NewMockCacheWithStats()
	metrics := caching.NewCacheMetrics()
	repo := caching.NewFakeDashboardRepository()
	svc := caching.NewRobustDashboardService(repo, cache, metrics)
	ctx := context.Background()

	date := time.Now()
	branchID := int64(10)

	// Tenant 1 request
	res1, _ := svc.GetDashboardWithTenant(ctx, 1, branchID, date)
	// Tenant 2 request
	res2, _ := svc.GetDashboardWithTenant(ctx, 2, branchID, date)

	// In FakeDashboardRepository, InvoiceCountToday = 42 + int(tenantID%10)
	expectedT1 := 42 + 1 // 43
	expectedT2 := 42 + 2 // 44

	if res1.InvoiceCountToday != expectedT1 {
		t.Errorf("Tenant 1 expected %d invoices, got %d", expectedT1, res1.InvoiceCountToday)
	}

	if res2.InvoiceCountToday != expectedT2 {
		t.Errorf("Tenant 2 expected %d invoices, got %d", expectedT2, res2.InvoiceCountToday)
	}

	// 2 distinct repository calls
	if repo.CallCount != 2 {
		t.Errorf("Expected 2 repo calls, got %d", repo.CallCount)
	}

	// Subsequent requests should hit cache, retaining correct tenant values
	res1Cached, _ := svc.GetDashboardWithTenant(ctx, 1, branchID, date)
	res2Cached, _ := svc.GetDashboardWithTenant(ctx, 2, branchID, date)

	if res1Cached.InvoiceCountToday != expectedT1 {
		t.Errorf("Tenant 1 cached expected %d, got %d", expectedT1, res1Cached.InvoiceCountToday)
	}

	if res2Cached.InvoiceCountToday != expectedT2 {
		t.Errorf("Tenant 2 cached expected %d, got %d", expectedT2, res2Cached.InvoiceCountToday)
	}

	// Still 2 repository calls (cache hit)
	if repo.CallCount != 2 {
		t.Errorf("Expected repo calls to remain 2, got %d", repo.CallCount)
	}

	t.Log("✓ Multi-tenant behavioral isolation validated")
}
