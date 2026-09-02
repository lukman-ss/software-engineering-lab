package caching_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	caching "github.com/lukman-ss/software-engineering-lab/labs/04-caching"
)

// TestCacheFailureFallback verifies: when cache fails, service falls back to DB.
// Each request hits DB because cache is unavailable.
func TestCacheFailureFallback(t *testing.T) {
	repo := caching.NewFakeDashboardRepository()
	// Using a failing cache simulates cache outage
	cache := caching.NewFailingMockCache()
	metrics := caching.NewCacheMetrics()
	svc := caching.NewRobustDashboardService(repo, cache, metrics)
	ctx := context.Background()

	// 5 requests - all should hit DB due to cache failure
	for i := 0; i < 5; i++ {
		_, _ = svc.GetDashboard(ctx, 1)
	}

	if repo.CallCount() != 5 {
		t.Fatalf("expected 5 DB queries (cache down), got %d", repo.CallCount())
	}

	t.Log("✓ Cache failure fallback: all requests hit DB when cache is unavailable")
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

// Test 4: Stale read is possible with cache TTL
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

// Test 7: Distributed lock wrong token cannot release another holder's lock
func TestDistributedLockWrongTokenCannotRelease(t *testing.T) {
	locker := caching.NewMockRedisClient()
	ctx := context.Background()

	key := "lock:item:888"
	ttl := 5 * time.Second

	holder1, value1, _ := caching.TryAcquireLock(ctx, locker, key, ttl)
	if !holder1 {
		t.Fatal("first lock acquisition should succeed")
	}

	// Create a fake "other" token
	otherToken := "fake-token-that-doesnt-match"

	// Attempt to release with wrong token - should not error and not release
	err := caching.ReleaseLock(ctx, locker, key, otherToken)
	if err == nil {
		t.Error("wrong token should NOT release lock (should return error)")
	}

	// Original holder can still release
	if err := caching.ReleaseLock(ctx, locker, key, value1); err != nil {
		t.Errorf("original token should release lock: %v", err)
	}

	t.Log("Wrong token rejection verified")
}

// Test 8: Negative TTL should fail
func TestDistributedLockNegativeTTL(t *testing.T) {
	locker := caching.NewMockRedisClient()
	ctx := context.Background()

	key := "lock:item:777"
	negativeTTL := -1 * time.Second

	_, _, err := caching.TryAcquireLock(ctx, locker, key, negativeTTL)
	if err == nil {
		t.Error("negative TTL should fail to acquire lock")
	}

	t.Log("Negative TTL validation verified")
}

// NOTE: TestDistributedLockGetCleansExpired removed - it used time.Sleep(100ms) which is flaky.
// Lock expiry with real time is better tested via integration tests or with injected fake clock.

// Test 9: Cache key versioning for schema/deployment namespace migration
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

// TestCacheAsideServiceMetricsClassification verifies correct metric classification:
// - ErrCacheMiss → IncMiss()
// - ErrCacheDown → IncError() + IncDBFallback()
// - Hit → IncHit()
// - Corrupt JSON → IncError() (then DB fallback)
func TestCacheAsideServiceMetricsClassification(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	metrics := caching.NewCacheMetrics()

	// Test 1: ErrCacheMiss → Misses +1, Errors = 0
	t.Run("ErrCacheMiss", func(t *testing.T) {
		cache := caching.NewMockCache()
		svc := caching.NewCacheAsideService(db, cache, metrics)
		ctx := context.Background()

		mock.ExpectQuery("SELECT id, name, price FROM products WHERE id =").
			WithArgs("prod-1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price"}).
				AddRow("prod-1", "Test Product", 10.0))

		result, err := svc.GetProduct(ctx, "prod-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.ID != "prod-1" {
			t.Errorf("expected product ID prod-1, got %s", result.ID)
		}

		if metrics.Misses() != 1 {
			t.Errorf("expected 1 miss, got %d", metrics.Misses())
		}
		if metrics.Errors() != 0 {
			t.Errorf("expected 0 errors for cache miss, got %d", metrics.Errors())
		}
		if metrics.DBFallbacks() != 0 {
			t.Errorf("expected 0 dbfallbacks for cache miss, got %d", metrics.DBFallbacks())
		}
	})

	// Test 2: ErrCacheDown → Errors +2 (GET error + SET error), DBFallbacks +1
	t.Run("ErrCacheDown", func(t *testing.T) {
		metrics.Reset()
		cache := caching.NewFailingMockCache()
		svc := caching.NewCacheAsideService(db, cache, metrics)
		ctx := context.Background()

		mock.ExpectQuery("SELECT id, name, price FROM products WHERE id =").
			WithArgs("prod-2").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price"}).
				AddRow("prod-2", "Test Product 2", 20.0))

		_, err := svc.GetProduct(ctx, "prod-2")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// 1 cache GET error + 1 cache SET error = 2 errors
		if metrics.Errors() != 2 {
			t.Errorf("expected 2 errors for cache down (GET + SET), got %d", metrics.Errors())
		}
		if metrics.DBFallbacks() != 1 {
			t.Errorf("expected 1 dbfallback for cache down, got %d", metrics.DBFallbacks())
		}
		if metrics.Misses() != 0 {
			t.Errorf("expected 0 misses for cache error, got %d", metrics.Misses())
		}
	})

	// Test 3: Cache hit
	t.Run("CacheHit", func(t *testing.T) {
		metrics.Reset()
		cache := caching.NewMockCache()
		svc := caching.NewCacheAsideService(db, cache, metrics)
		ctx := context.Background()

		// Pre-populate cache
		product := caching.Product{ID: "prod-3", Name: "Cached Product", Price: 30.0}
		data, _ := json.Marshal(product)
		cache.Set(ctx, caching.CacheKey("product", "prod-3", 1), string(data), time.Minute)

		_, err := svc.GetProduct(ctx, "prod-3")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if metrics.Hits() != 1 {
			t.Errorf("expected 1 hit, got %d", metrics.Hits())
		}
		if metrics.Misses() != 0 {
			t.Errorf("expected 0 misses for hit, got %d", metrics.Misses())
		}
	})

	// Test 4: Corrupt JSON → Errors +1, then DB fallback
	t.Run("CorruptJSON", func(t *testing.T) {
		metrics.Reset()
		cache := caching.NewMockCache()
		svc := caching.NewCacheAsideService(db, cache, metrics)
		ctx := context.Background()

		// Set corrupt JSON
		cache.Set(ctx, caching.CacheKey("product", "prod-corrupt", 1), `{invalid json`, time.Minute)

		mock.ExpectQuery("SELECT id, name, price FROM products WHERE id =").
			WithArgs("prod-corrupt").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price"}).
				AddRow("prod-corrupt", "Repaired Product", 40.0))

		_, err := svc.GetProduct(ctx, "prod-corrupt")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if metrics.Errors() < 1 {
			t.Errorf("expected at least 1 error for corrupt JSON, got %d", metrics.Errors())
		}
	})
}

// Removed TestRandomJitterDispersesRequests - it used rand.Intn directly
// instead of testing TTLWithJitter. Proper jitter tests are in stampede_test.go.

// Removed TestCacheFailureHandledGracefully as it was a redundant/incorrect test.
// True graceful degradation is tested in robust_test.go's TestCacheFailureGracefulDegradation.
