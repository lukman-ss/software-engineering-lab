package caching_test

import (
	"context"
	"sync"
	"testing"
	"time"

	caching "github.com/lukman-ss/software-engineering-lab/labs/04-caching"
)

var testFixedDate = time.Date(2026, 9, 3, 0, 0, 0, 0, time.UTC)

// Integration tests for essential cache behaviors.

// TestCacheMissReadsDatabase verifies: cache miss reads from repository
func TestCacheMissReadsDatabase(t *testing.T) {
	cache := caching.NewMockCache()
	metrics := caching.NewCacheMetrics()
	repo := caching.NewFakeDashboardRepository()
	svc := caching.NewRobustDashboardService(repo, cache, metrics)
	ctx := context.Background()

	// First request = cache miss → database read
	_, err := svc.GetDashboardWithTenant(ctx, 1, 1, testFixedDate)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.CallCount() == 0 {
		t.Error("cache miss should trigger repository call")
	}
	if metrics.Misses() == 0 {
		t.Error("cache miss count should be incremented")
	}
	t.Log("✓ Cache miss reads from repository")
	t.Logf("✓ Repository CallCount = %d (expected 1)", repo.CallCount())
}

// TestCacheMissWritesCache verifies: cache miss writes to cache after read
func TestCacheMissWritesCache(t *testing.T) {
	cache := caching.NewMockCacheWithStats()
	metrics := caching.NewCacheMetrics()
	repo := caching.NewFakeDashboardRepository()
	svc := caching.NewRobustDashboardService(repo, cache, metrics)
	ctx := context.Background()

	// First request
	_, _ = svc.GetDashboardWithTenant(ctx, 1, 1, testFixedDate)

	// Second request - should be cache hit
	_, _ = svc.GetDashboardWithTenant(ctx, 1, 1, testFixedDate)

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
	_, _ = svc.GetDashboardWithTenant(ctx, 1, 1, testFixedDate)

	initialCalls := repo.CallCount()

	// Multiple cache hits
	for i := 0; i < 5; i++ {
		_, _ = svc.GetDashboardWithTenant(ctx, 1, 1, testFixedDate)
	}

	if repo.CallCount() > initialCalls {
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
	_, _ = svc.GetDashboardWithTenant(ctx, 1, 1, testFixedDate)
	initialRebuilds := metrics.Rebuilds()

	// Force miss by deleting
	key := caching.DashboardCacheKey(1, 1, testFixedDate)
	_ = cache.Delete(ctx, key)

	// Next request rebuilds
	_, _ = svc.GetDashboardWithTenant(ctx, 1, 1, testFixedDate)

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
	key := caching.DashboardCacheKey(1, branchID, testFixedDate)

	// Step 1: Populate cache with initial data
	repo.Reset()
	cache.Set(ctx, key, `{"branch_id":1,"invoice_count":10,"date":"2026-09-03"}`, time.Minute)

	// Verify cache has initial data
	cached, _ := cache.Get(ctx, key)
	if cached == "" {
		t.Error("cache should have initial data")
	}

	// Step 2: Invalidate
	_ = svc.InvalidateDashboard(ctx, 1, branchID, testFixedDate)

	// Step 3: Next request should query repo and update cache
	repo.Reset()
	_, _ = svc.GetDashboardWithTenant(ctx, 1, branchID, testFixedDate)

	if repo.CallCount() != 1 {
		t.Errorf("expected 1 repository call after invalidation, got %d", repo.CallCount())
	}
	t.Log("✓ Cache invalidation returns fresh data on next request")
}

// TestDifferentBranchDoesNotShareCache verifies: different branch has different cache
func TestDifferentBranchDoesNotShareCache(t *testing.T) {
	cache := caching.NewMockCache()
	ctx := context.Background()

	// Set data for branch 1
	key1 := caching.DashboardCacheKey(1, 1, testFixedDate)
	data1 := `{"branch_id":1,"invoice_count":50}`
	cache.Set(ctx, key1, data1, time.Minute)

	// Set data for branch 2
	key2 := caching.DashboardCacheKey(1, 2, testFixedDate)
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

// TestDifferentTenantDoesNotShareCache verifies tenant isolation
func TestDifferentTenantDoesNotShareCache(t *testing.T) {
	cache := caching.NewMockCache()
	ctx := context.Background()

	keyTenantA := caching.DashboardCacheKey(1, 1, testFixedDate)
	keyTenantB := caching.DashboardCacheKey(2, 1, testFixedDate)

	t.Logf("Tenant A key: %s", keyTenantA)
	t.Logf("Tenant B key: %s", keyTenantB)

	if keyTenantA == keyTenantB {
		t.Fatal("Tenant A and Tenant B must not share the same cache key")
	}

	cache.Set(ctx, keyTenantA, `{"branch_id":1,"invoice_count":10}`, time.Minute)
	cache.Set(ctx, keyTenantB, `{"branch_id":1,"invoice_count":99}`, time.Minute)

	cachedA, _ := cache.Get(ctx, keyTenantA)
	cachedB, _ := cache.Get(ctx, keyTenantB)

	if cachedA == cachedB {
		t.Fatal("Data leakage: Tenant A and Tenant B received identical cache data")
	}

	t.Log("✓ Different tenant does not share cache")
}

// TestTenantIsolationInService verifies multi-tenant isolation in RobustDashboardService
func TestTenantIsolationInService(t *testing.T) {
	cache := caching.NewMockCache()
	metrics := caching.NewCacheMetrics()
	repo := caching.NewFakeDashboardRepository()
	svc := caching.NewRobustDashboardService(repo, cache, metrics)
	ctx := context.Background()

	keyA1 := caching.DashboardCacheKey(1, 1, testFixedDate)
	keyB1 := caching.DashboardCacheKey(2, 1, testFixedDate)

	if keyA1 == keyB1 {
		t.Fatal("Tenant A and Tenant B keys must be different")
	}

	resultA1, err := svc.GetDashboardWithTenant(ctx, 1, 1, testFixedDate)
	if err != nil {
		t.Fatalf("Tenant A request failed: %v", err)
	}

	if resultA1.InvoiceCountToday != 43 {
		t.Errorf("expected Tenant A invoice count 43, got %d", resultA1.InvoiceCountToday)
	}

	resultB1, err := svc.GetDashboardWithTenant(ctx, 2, 1, testFixedDate)
	if err != nil {
		t.Fatalf("Tenant B request failed: %v", err)
	}

	if resultB1.InvoiceCountToday != 44 {
		t.Errorf("expected Tenant B invoice count 44, got %d", resultB1.InvoiceCountToday)
	}

	if resultA1.InvoiceCountToday == resultB1.InvoiceCountToday {
		t.Fatalf("DATA LEAK: Tenant A data matches Tenant B data (%d == %d)",
			resultA1.InvoiceCountToday, resultB1.InvoiceCountToday)
	}

	t.Log("✓ Tenant isolation verified: Tenant A Branch 1 ≠ Tenant B Branch 1")
}

// TestCorruptCacheFallsBackAppropriately verifies fallback on corrupt cache
func TestCorruptCacheFallsBackAppropriately(t *testing.T) {
	cache := caching.NewMockCache()
	metrics := caching.NewCacheMetrics()
	repo := caching.NewFakeDashboardRepository()
	svc := caching.NewRobustDashboardService(repo, cache, metrics)
	ctx := context.Background()

	branchID := int64(1)
	key := caching.DashboardCacheKey(1, branchID, testFixedDate)

	cache.Set(ctx, key, `{"invalid json`, time.Minute)

	result, err := svc.GetDashboardWithTenant(ctx, 1, branchID, testFixedDate)
	if err != nil {
		t.Fatalf("expected fallback to succeed, got error: %v", err)
	}

	if result.BranchID != branchID {
		t.Errorf("expected BranchID %d, got %d", branchID, result.BranchID)
	}

	if metrics.Errors() == 0 {
		t.Error("corrupt cache should increment error counter")
	}

	t.Logf("Result from database fallback: %+v", result)
	t.Log("✓ Corrupt cache falls back appropriately")
}

// TestCacheFailureFallsBackToDatabase verifies fallback when cache is down
func TestCacheFailureFallsBackToDatabase(t *testing.T) {
	cache := caching.NewFailingMockCache()
	metrics := caching.NewCacheMetrics()
	repo := caching.NewFakeDashboardRepository()
	svc := caching.NewRobustDashboardService(repo, cache, metrics)
	ctx := context.Background()

	result, err := svc.GetDashboardWithTenant(ctx, 1, 1, testFixedDate)
	if err != nil {
		t.Fatalf("expected fallback to succeed, got error: %v", err)
	}

	if result.BranchID != 1 {
		t.Errorf("expected BranchID 1, got %d", result.BranchID)
	}

	if repo.CallCount() == 0 {
		t.Error("repository should have been called for fallback")
	}

	if metrics.Errors() == 0 {
		t.Error("cache failure should increment error counter")
	}

	t.Logf("Result from database fallback: %+v", result)
	t.Log("✓ Cache failure falls back to database")
}

// TestConcurrentCacheMissProtectedBySingleflight verifies singleflight deduplication deterministically
func TestConcurrentCacheMissProtectedBySingleflight(t *testing.T) {
	cache := caching.NewMockCache()
	repo := caching.NewCounterRepository()
	svc := caching.NewProtectedStampedeService(cache, repo)
	ctx := context.Background()

	numGoroutines := 10
	waiterEntered := make(chan struct{}, numGoroutines)
	svc.SetOnWaitEntry(func(key string) {
		waiterEntered <- struct{}{}
	})

	repo.Block()

	var wg sync.WaitGroup
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = svc.GetData(ctx, 1)
		}()
	}

	// 1. Leader enters repository
	repo.WaitUntilEntered()

	// 2. Wait until all goroutines have entered/joined singleflight
	for i := 0; i < numGoroutines; i++ {
		<-waiterEntered
	}

	// 3. Unblock repository
	repo.Unblock()
	wg.Wait()

	if repo.CallCount() != 1 {
		t.Errorf("singleflight should dedupe concurrent misses to 1 DB query, got %d", repo.CallCount())
	}
	t.Log("✓ Concurrent cache miss protected by singleflight (deterministic)")
}

// TestTTLJitterRange verifies: TTL jitter stays within configured range [base, base + maxJitter)
func TestTTLJitterRange(t *testing.T) {
	baseTTL := 60 * time.Second
	maxJitter := 15 * time.Second

	for i := 0; i < 1000; i++ {
		jitteredTTL := caching.TTLWithJitter(baseTTL, maxJitter)
		if jitteredTTL < baseTTL || jitteredTTL >= baseTTL+maxJitter {
			t.Errorf("TTL jitter out of range [base, base+maxJitter): %v", jitteredTTL)
		}
	}

	if ttl := caching.TTLWithJitter(baseTTL, 0); ttl != baseTTL {
		t.Errorf("expected ttl == baseTTL when maxJitter == 0, got %v", ttl)
	}
	if ttl := caching.TTLWithJitter(baseTTL, -5*time.Second); ttl != baseTTL {
		t.Errorf("expected ttl == baseTTL when maxJitter < 0, got %v", ttl)
	}

	t.Log("✓ TTL jitter stays within configured range [base, base+maxJitter)")
}

// TestCacheHitRatioCalculatesCorrectly verifies hit ratio calculation
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
func TestHistoricalDateInvariant(t *testing.T) {
	cache := caching.NewMockCache()
	metrics := caching.NewCacheMetrics()
	repo := caching.NewFakeDashboardRepository()
	svc := caching.NewRobustDashboardService(repo, cache, metrics)
	ctx := context.Background()

	historicalDate := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	tenantID, branchID := int64(1), int64(42)

	result, err := svc.GetDashboardWithTenant(ctx, tenantID, branchID, historicalDate)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedDateStr := "2026-08-01"
	if result.Date != expectedDateStr {
		t.Errorf("Result date is %s, expected %s", result.Date, expectedDateStr)
	}

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

	branchID := int64(10)

	res1, _ := svc.GetDashboardWithTenant(ctx, 1, branchID, testFixedDate)
	res2, _ := svc.GetDashboardWithTenant(ctx, 2, branchID, testFixedDate)

	expectedT1 := 42 + 1 // 43
	expectedT2 := 42 + 2 // 44

	if res1.InvoiceCountToday != expectedT1 {
		t.Errorf("Tenant 1 expected %d invoices, got %d", expectedT1, res1.InvoiceCountToday)
	}

	if res2.InvoiceCountToday != expectedT2 {
		t.Errorf("Tenant 2 expected %d invoices, got %d", expectedT2, res2.InvoiceCountToday)
	}

	if repo.CallCount() != 2 {
		t.Errorf("Expected 2 repo calls, got %d", repo.CallCount())
	}

	res1Cached, _ := svc.GetDashboardWithTenant(ctx, 1, branchID, testFixedDate)
	res2Cached, _ := svc.GetDashboardWithTenant(ctx, 2, branchID, testFixedDate)

	if res1Cached.InvoiceCountToday != expectedT1 {
		t.Errorf("Tenant 1 cached expected %d, got %d", expectedT1, res1Cached.InvoiceCountToday)
	}

	if res2Cached.InvoiceCountToday != expectedT2 {
		t.Errorf("Tenant 2 cached expected %d, got %d", expectedT2, res2Cached.InvoiceCountToday)
	}

	if repo.CallCount() != 2 {
		t.Errorf("Expected repo calls to remain 2, got %d", repo.CallCount())
	}

	t.Log("✓ Multi-tenant behavioral isolation validated")
}

// TestSingleflightContextCancellation verifies cancellation
func TestSingleflightContextCancellation(t *testing.T) {
	cache := caching.NewMockCache()
	repo := caching.NewCounterRepository()
	svc := caching.NewProtectedStampedeService(cache, repo)

	ctx1 := context.Background()
	ctx2, cancel := context.WithCancel(context.Background())

	errCh := make(chan error, 2)

	repo.Block()

	go func() {
		_, err := svc.GetData(ctx1, 99)
		errCh <- err
	}()

	repo.WaitUntilEntered()

	go func() {
		_, err := svc.GetData(ctx2, 99)
		errCh <- err
	}()

	cancel()

	err2 := <-errCh
	if err2 != context.Canceled {
		t.Errorf("Expected context.Canceled for second request, got: %v", err2)
	}

	repo.Unblock()

	err1 := <-errCh
	if err1 != nil {
		t.Errorf("First request should complete successfully, got error: %v", err1)
	}

	t.Log("✓ Context cancellation during singleflight wait validated")
}
