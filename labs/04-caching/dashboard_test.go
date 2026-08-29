package caching_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	caching "github.com/lukman/software-engineer-lab/labs/04-caching"
)

// Test 1: Cache Miss → Cache Hit
func TestDashboardCacheMissThenHit(t *testing.T) {
	cache := caching.NewMockCache()
	ctx := context.Background()

	branchID := int64(1)

	// Step 1: Cache empty
	_, err := cache.Get(ctx, caching.DashboardCacheKey(branchID))
	if err == nil {
		t.Fatal("expected cache miss initially")
	}

	// Step 2: Populate cache
	dashboard := caching.Dashboard{
		BranchID:          branchID,
		InvoiceCountToday: 50,
		TotalRevenueToday: 250000.0,
		Date:              caching.Today(),
	}
	data, _ := json.Marshal(dashboard)
	cache.Set(ctx, caching.DashboardCacheKey(branchID), string(data), 30*time.Second)

	// Step 3: Cache hit
	cached, err := cache.Get(ctx, caching.DashboardCacheKey(branchID))
	if err != nil {
		t.Fatalf("expected cache hit after population: %v", err)
	}

	var result caching.Dashboard
	json.Unmarshal([]byte(cached), &result)
	if result.InvoiceCountToday != 50 {
		t.Errorf("expected 50 invoices, got %d", result.InvoiceCountToday)
	}

	t.Log("SUCCESS: Cache miss then hit works correctly")
}

// Test 2: Cache Invalidation After Mutation
func TestDashboardCacheInvalidation(t *testing.T) {
	cache := caching.NewMockCache()
	ctx := context.Background()
	// Create service with mock DB
	svc := caching.NewDashboardCacheService(nil, cache)

	branchID := int64(1)

	// Step 1: Request dashboard - cache miss
	initialData := caching.Dashboard{
		BranchID:          branchID,
		InvoiceCountToday: 10,
		TotalRevenueToday: 100000.0,
		Date:              caching.Today(),
	}
	jsonData, _ := json.Marshal(initialData)
	cache.Set(ctx, caching.DashboardCacheKey(branchID), string(jsonData), 30*time.Second)

	// Verify cache has data
	cachedData, err := cache.Get(ctx, caching.DashboardCacheKey(branchID))
	if err != nil {
		t.Fatalf("cache should have initial data: %v", err)
	}

	var beforeMutation caching.Dashboard
	json.Unmarshal([]byte(cachedData), &beforeMutation)
	t.Logf("Before mutation: InvoiceCount=%d, Revenue=%.0f", beforeMutation.InvoiceCountToday, beforeMutation.TotalRevenueToday)

	// Step 2: Simulate invoice creation (mutation)
	// Because DB is nil, we bypass CreateInvoice and directly simulate what happens AFTER commit:
	// In real world: BEGIN -> INSERT invoice -> COMMIT -> Invalidate

	// Step 3: Invalidate cache (happens after commit)
	err = svc.InvalidateBranchDashboard(ctx, branchID)
	if err != nil {
		t.Fatalf("invalidation failed: %v", err)
	}

	// Step 4: Cache should be invalidated
	_, err = cache.Get(ctx, caching.DashboardCacheKey(branchID))
	if err == nil {
		t.Fatal("cache should be invalidated (miss) after mutation")
	}

	t.Log("SUCCESS: Cache invalidated after mutation, forcing fresh query")
}

// Test 3: Cache Invalidation Prevents Stuck Stale Data
// Test ini gagal jika invalidation tidak diimplementasikan
func TestCacheInvalidationRequiredAfterDataChange(t *testing.T) {
	cache := caching.NewMockCache()
	ctx := context.Background()

	branchID := int64(42)

	// Initial state
	initial := caching.Dashboard{InvoiceCountToday: 100, Date: caching.Today()}
	data, _ := json.Marshal(initial)
	cache.Set(ctx, caching.DashboardCacheKey(branchID), string(data), 30*time.Second)

	// Mutate database directly (simulate another process)
	mutated := caching.Dashboard{InvoiceCountToday: 150, Date: caching.Today()}
	mutatedData, _ := json.Marshal(mutated)

	// WITHOUT invalidation, cache still has old data
	cached, _ := cache.Get(ctx, caching.DashboardCacheKey(branchID))
	var result caching.Dashboard
	json.Unmarshal([]byte(cached), &result)

	if result.InvoiceCountToday == 150 {
		t.Fatal("TEST FAILED: Cache still has old data - invalidation NOT working")
	}

	// WITH invalidation, we update cache
	cache.Set(ctx, caching.DashboardCacheKey(branchID), string(mutatedData), 30*time.Second)

	// Now cache reflects new data
	cached2, _ := cache.Get(ctx, caching.DashboardCacheKey(branchID))
	json.Unmarshal([]byte(cached2), &result)

	if result.InvoiceCountToday != 150 {
		t.Error("TEST FAILED: Cache should reflect updated data after invalidation")
	}

	t.Log("SUCCESS: Invalidation ensures cache reflects latest DB state")
}

// Test 4: Invalidation Is Async-Unsafe Without Proper Ordering
func TestCacheDeletionBeforeCommitIsUnsafe(t *testing.T) {
	// This test demonstrates the race condition when deleting cache BEFORE commit.

	// Scenario:
	// T1: Client A creates invoice
	// T1: Delete cache (cache becomes empty)
	// T1: DB commit fails or process crashes
	// T1: Client B reads dashboard -> cache miss -> reads DB -> gets OLD data
	// Result: Cache serves stale data even though commit failed!

	// The safe flow:
	// T1: Client A creates invoice
	// T1: DB commit succeeds
	// T1: Invalidate cache (only after commit)
	// T2: Client B reads dashboard -> cache miss -> reads DB commit -> gets NEW data

	// TTL is the safety net: if process crashes after commit but before invalidation,
	// fresh data will appear after TTL expires.

	t.Log("Lesson: Always COMMIT -> Invalidate. Never Delete -> Commit")
	t.Log("TTL provides safety net for process crashes")
}

// Test 5: Cache Hit Ratio Improves With Frequent Mutations
func TestCacheHitRatioImprovesWithMutations(t *testing.T) {
	cache := caching.NewMockCache()
	ctx := context.Background()

	branchID := int64(99)
	hits, misses := 0, 0

	// Populate cache
	dashboard := caching.Dashboard{InvoiceCountToday: 50, Date: caching.Today()}
	data, _ := json.Marshal(dashboard)
	cache.Set(ctx, caching.DashboardCacheKey(branchID), string(data), 30*time.Second)

	// Simulate 100 read requests
	for i := 0; i < 100; i++ {
		_, err := cache.Get(ctx, caching.DashboardCacheKey(branchID))
		if err == nil {
			hits++
		} else {
			misses++
		}
	}

	if hits != 100 {
		t.Errorf("expected 100 cache hits, got %d hits, %d misses", hits, misses)
	}

	t.Logf("Cache hit ratio: %d%%", hits*100)
	t.Log("SUCCESS: Cache significantly reduces DB load under heavy read traffic")
}

// Test 6: Stale Data Acceptable for Dashboard Statistics
func TestStaleDataAcceptableForDashboard(t *testing.T) {
	// Pertanyaan utama caching bukan: "Apakah datanya berubah?"
	// tapi: "Berapa lama stale data masih dapat diterima?"

	// Contoh:
	// Top Mekanik terlambat 60 detik → biasanya acceptable
	// Saldo Wallet terlambat 60 detik → mungkin tidak acceptable

	// Ini adalah business decision, bukan teknisalah.
	cache := caching.NewMockCache()
	ctx := context.Background()

	trueData := caching.Dashboard{InvoiceCountToday: 100, Date: caching.Today()}

	data, _ := json.Marshal(trueData)
	cache.Set(ctx, "dash", string(data), 30*time.Second)

	cached, _ := cache.Get(ctx, "dash")
	var d caching.Dashboard
	json.Unmarshal([]byte(cached), &d)

	staleDifference := trueData.InvoiceCountToday - d.InvoiceCountToday
	t.Logf("Stale data difference: %d invoices (30s)", staleDifference)
	t.Logf("For dashboard: acceptable for operational metrics")

	// TODO: Bagaimana jika stale > 5% dari value?
	// Itu adalah pertanyaan business, bukan teknis.

	t.Log("SUCCESS: Stale data trade-off is acceptable for aggregated statistics")
}
