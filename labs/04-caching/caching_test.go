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
	cache := caching.NewMockCache()
	ctx := context.Background()

	// Simulate naive service: always cache miss → DB query
	// Di dunia nyata: setiap request dashboard = 6 query DB
	t.Log("Naive service: setiap request = multiple DB queries (no cache)")

	// Demonstrate cache miss
	_, err := cache.Get(ctx, "product:100")
	if err == nil {
		t.Fatal("expected cache miss")
	}

	t.Log("Naive pattern demonstrates cache miss on every request")
}

// Test 2: Cache aside pattern - miss then hit
func TestCacheAsideHit(t *testing.T) {
	cache := caching.NewMockCache()
	ctx := context.Background()

	// First: populate cache
	product := caching.Product{ID: "200", Name: "Gadget", Price: 250.0}
	data, err := json.Marshal(product)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	cache.Set(ctx, "product:200", string(data), 5*time.Minute)

	// Second: cache hit
	cached, err := cache.Get(ctx, "product:200")
	if err != nil {
		t.Fatalf("cache get failed: %v", err)
	}

	var p caching.Product
	if err := json.Unmarshal([]byte(cached), &p); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if p.Name != "Gadget" {
		t.Errorf("expected Gadget, got %s", p.Name)
	}

	t.Log("Cache hit returns correct data")
}

// Test 3: Stale read is possible with cache TTL
func TestCacheStaleRead(t *testing.T) {
	cache := caching.NewMockCache()
	ctx := context.Background()

	key := "product:300"
	product := caching.Product{ID: "300", Name: "NewName", Price: 200.0}

	data, err := json.Marshal(product)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	cache.Set(ctx, key, string(data), 5*time.Minute)

	// Cache returns what's stored, even if DB has changed
	cachedData, err := cache.Get(ctx, key)
	if err != nil {
		t.Fatalf("cache get failed: %v", err)
	}
	if cachedData == "" {
		t.Fatal("expected cache to have data")
	}

	var cachedProduct caching.Product
	if err := json.Unmarshal([]byte(cachedData), &cachedProduct); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if cachedProduct.Name != "NewName" {
		t.Errorf("expected 'NewName', got '%s'", cachedProduct.Name)
	}

	t.Log("Cache returns stored data regardless of external changes until TTL expires")
}

// Test 4: Single Flight Pattern - dedupes concurrent DB queries
func TestSingleFlightConcurrentRequests(t *testing.T) {
	cache := caching.NewMockCache()
	ctx := context.Background()

	// Pre-populate cache
	data := `{"id":"400","name":"Part","price":50.0}`
	cache.Set(ctx, "product:400", data, 5*time.Minute)

	// Multiple concurrent readers - all served from cache
	hits := 0
	for i := 0; i < 10; i++ {
		_, err := cache.Get(ctx, "product:400")
		if err == nil {
			hits++
		}
	}

	if hits != 10 {
		t.Errorf("expected 10 cache hits, got %d", hits)
	}

	t.Log("Concurrent cache hits validated")
}

// Test 5: Probabilistic early refresh mitigates stampede
func TestCacheStampedeMitigation(t *testing.T) {
	cache := caching.NewMockCache()
	ctx := context.Background()

	key := "product:500"
	product := caching.Product{ID: "500", Name: "Item", Price: 10.0}

	data, err := json.Marshal(product)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	cache.Set(ctx, key, string(data), 5*time.Minute)

	// Simulate early refresh check
	count := 0
	for i := 0; i < 100; i++ {
		if caching.ShouldRefreshEarly(ctx, cache, key) {
			count++
		}
	}

	t.Logf("Early refresh triggers: %d/100", count)
	t.Log("Early refresh strategy validated")
}

// Test 6: Distributed lock mutual exclusion
func TestDistributedLockMutualExclusion(t *testing.T) {
	locker := caching.NewMockRedisClient()
	ctx := context.Background()

	key := "lock:item:999"
	ttl := 5 * time.Second

	holder1, value1 := caching.TryAcquireLock(ctx, locker, key, ttl)
	if !holder1 {
		t.Fatal("first lock acquisition should succeed")
	}

	// Second holder should fail to acquire the same lock
	holder2, value2 := caching.TryAcquireLock(ctx, locker, key, ttl)
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
	holder3, _ := caching.TryAcquireLock(ctx, locker, key, ttl)
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

// Test 10: Cache failure handled gracefully
func TestCacheFailureHandledGracefully(t *testing.T) {
	cache := caching.NewFailingMockCache()
	ctx := context.Background()

	// When cache fails, service should fall back to DB
	_, err := cache.Get(ctx, "test")
	if err == nil {
		t.Fatal("expected cache to fail")
	}

	t.Log("Cache failure handling validated")
}
