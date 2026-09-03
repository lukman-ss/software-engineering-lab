package caching_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	caching "github.com/lukman-ss/software-engineering-lab/labs/04-caching"
)

// Test 1: Cache Miss → Cache Hit
func TestDashboardCacheMissThenHit(t *testing.T) {
	cache := caching.NewMockCache()
	ctx := context.Background()

	branchID := int64(1)
	fixedDate := time.Date(2026, 9, 3, 0, 0, 0, 0, time.UTC)
	key := caching.NewTenantDashboardKey(1, branchID, fixedDate).Build()

	// Step 1: Cache empty
	_, err := cache.Get(ctx, key)
	if err == nil {
		t.Fatal("expected cache miss initially")
	}

	// Step 2: Populate cache
	dashboard := caching.Dashboard{
		BranchID:          branchID,
		InvoiceCountToday: 50,
		TotalRevenueToday: 250000.0,
		Date:              "2026-09-03",
	}
	data, _ := json.Marshal(dashboard)
	cache.Set(ctx, key, string(data), 30*time.Second)

	// Step 3: Cache hit
	cached, err := cache.Get(ctx, key)
	if err != nil {
		t.Fatalf("expected cache hit, got error: %v", err)
	}

	var result caching.Dashboard
	json.Unmarshal([]byte(cached), &result)
	if result.InvoiceCountToday != 50 {
		t.Errorf("expected 50 invoices, got %d", result.InvoiceCountToday)
	}

	t.Log("Cache miss then hit behavior validated")
}

// Test 2: Cache Invalidation After Mutation
func TestDashboardCacheInvalidation(t *testing.T) {
	cache := caching.NewMockCache()
	ctx := context.Background()
	svc := caching.NewDashboardCacheService(nil, cache)

	branchID := int64(1)
	fixedDate := time.Date(2026, 9, 3, 0, 0, 0, 0, time.UTC)
	key := caching.NewTenantDashboardKey(1, branchID, fixedDate).Build()

	// Step 1: Request dashboard - cache miss
	initialData := caching.Dashboard{
		BranchID:          branchID,
		InvoiceCountToday: 10,
		TotalRevenueToday: 100000.0,
		Date:              "2026-09-03",
	}
	jsonData, _ := json.Marshal(initialData)
	cache.Set(ctx, key, string(jsonData), 30*time.Second)

	// Verify cache has data
	cachedData, err := cache.Get(ctx, key)
	if err != nil {
		t.Fatalf("cache should have initial data: %v", err)
	}

	var beforeMutation caching.Dashboard
	json.Unmarshal([]byte(cachedData), &beforeMutation)

	// Step 2: Invalidate cache (happens after commit)
	err = svc.InvalidateDashboard(ctx, 1, branchID, fixedDate)
	if err != nil {
		t.Fatalf("invalidation failed: %v", err)
	}

	// Step 3: Cache should be invalidated
	_, err = cache.Get(ctx, key)
	if err == nil {
		t.Fatal("cache should be invalidated (miss) after mutation")
	}

	t.Log("Cache invalidation after mutation validated")
}

// Test 3: Cache Invalidation Required After Data Change
func TestCacheInvalidationRequiredAfterDataChange(t *testing.T) {
	cache := caching.NewMockCache()
	metrics := caching.NewCacheMetrics()
	repo := caching.NewFakeDashboardRepository()
	svc := caching.NewRobustDashboardService(repo, cache, metrics)
	ctx := context.Background()

	branchID := int64(42)
	fixedDate := time.Date(2026, 9, 3, 0, 0, 0, 0, time.UTC)
	key := caching.NewTenantDashboardKey(1, branchID, fixedDate).Build()

	// Step 1: Populate cache via service read
	initial, err := svc.GetDashboardWithTenant(ctx, 1, branchID, fixedDate)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if initial.InvoiceCountToday != 43 {
		t.Fatalf("expected initial repo invoice count 43, got %d", initial.InvoiceCountToday)
	}

	// Step 2: Authoritative source changes
	repo.SetNextValue(func() caching.Dashboard {
		return caching.Dashboard{BranchID: branchID, InvoiceCountToday: 150, Date: "2026-09-03"}
	})

	// Step 3: Assert stale value is still returned from cache before invalidation
	staleResult, err := svc.GetDashboardWithTenant(ctx, 1, branchID, fixedDate)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if staleResult.InvoiceCountToday != 43 {
		t.Errorf("expected stale value 43 before invalidation, got %d", staleResult.InvoiceCountToday)
	}

	// Step 4: Invalidate cache
	if err := svc.InvalidateDashboard(ctx, 1, branchID, fixedDate); err != nil {
		t.Fatalf("invalidation failed: %v", err)
	}

	// Step 5: Assert cache miss on key
	_, err = cache.Get(ctx, key)
	if err == nil {
		t.Fatal("expected cache miss after invalidation")
	}

	// Step 6: Service fetches fresh value from repository
	freshResult, err := svc.GetDashboardWithTenant(ctx, 1, branchID, fixedDate)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if freshResult.InvoiceCountToday != 150 {
		t.Errorf("expected fresh value 150 after invalidation, got %d", freshResult.InvoiceCountToday)
	}

	t.Log("Cache invalidation flow validated: stale cache -> mutation -> invalidation -> miss -> fresh repo value")
}

// Test 4: Repeated Reads Hit Cache
func TestRepeatedReadsHitCache(t *testing.T) {
	cache := caching.NewMockCache()
	ctx := context.Background()

	branchID := int64(99)
	fixedDate := time.Date(2026, 9, 3, 0, 0, 0, 0, time.UTC)
	key := caching.NewTenantDashboardKey(1, branchID, fixedDate).Build()
	hits, misses := 0, 0

	// Populate cache
	dashboard := caching.Dashboard{InvoiceCountToday: 50, Date: "2026-09-03"}
	data, _ := json.Marshal(dashboard)
	cache.Set(ctx, key, string(data), 30*time.Second)

	// Simulate 100 read requests
	for i := 0; i < 100; i++ {
		_, err := cache.Get(ctx, key)
		if err == nil {
			hits++
		} else {
			misses++
		}
	}

	if hits != 100 {
		t.Errorf("expected 100 cache hits, got %d hits, %d misses", hits, misses)
	}

	ratio := float64(hits) / float64(hits+misses) * 100.0
	if ratio != 100.0 {
		t.Errorf("expected hit ratio 100.0%%, got %.2f%%", ratio)
	}

	t.Logf("Cache hit ratio: %.2f%%", ratio)
	t.Log("Cache reduces DB load significantly")
}

// Test 5: Dashboard Serves Stale Value Before Invalidation
func TestDashboardServesStaleValueBeforeInvalidation(t *testing.T) {
	cache := caching.NewMockCache()
	metrics := caching.NewCacheMetrics()
	repo := caching.NewFakeDashboardRepository()
	svc := caching.NewRobustDashboardService(repo, cache, metrics)
	ctx := context.Background()

	branchID := int64(7)
	fixedDate := time.Date(2026, 9, 3, 0, 0, 0, 0, time.UTC)

	// Populate cache with stale value (100)
	staleDashboard := caching.Dashboard{BranchID: branchID, InvoiceCountToday: 100, Date: "2026-09-03"}
	data, _ := json.Marshal(staleDashboard)
	cache.Set(ctx, caching.NewTenantDashboardKey(1, branchID, fixedDate).Build(), string(data), 30*time.Second)

	// Authoritative source has newer value (120)
	repo.SetNextValue(func() caching.Dashboard {
		return caching.Dashboard{BranchID: branchID, InvoiceCountToday: 120, Date: "2026-09-03"}
	})

	// Read before invalidation returns stale value (100)
	read1, err := svc.GetDashboardWithTenant(ctx, 1, branchID, fixedDate)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if read1.InvoiceCountToday != 100 {
		t.Errorf("expected stale value 100 before invalidation, got %d", read1.InvoiceCountToday)
	}

	// Invalidate cache
	if err := svc.InvalidateDashboard(ctx, 1, branchID, fixedDate); err != nil {
		t.Fatalf("invalidation failed: %v", err)
	}

	// Read after invalidation returns fresh value (120)
	read2, err := svc.GetDashboardWithTenant(ctx, 1, branchID, fixedDate)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if read2.InvoiceCountToday != 120 {
		t.Errorf("expected authoritative value 120 after invalidation, got %d", read2.InvoiceCountToday)
	}

	t.Log("Verified stale value served before invalidation and authoritative value served after invalidation")
}
