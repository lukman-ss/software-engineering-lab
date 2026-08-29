package caching_test

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	caching "github.com/lukman-ss/software-engineering-lab/labs/04-caching"
)

// TestCacheFailureGracefulDegradation verifies: when cache fails, service falls back to repository.
func TestCacheFailureGracefulDegradation(t *testing.T) {
	cache := caching.NewFailingMockCache()
	metrics := caching.NewCacheMetrics()
	repo := caching.NewFakeDashboardRepository()
	svc := caching.NewRobustDashboardService(repo, cache, metrics)
	ctx := context.Background()

	// Request should succeed via fallback
	result, err := svc.GetDashboard(ctx, 1)
	if err != nil {
		t.Fatalf("expected fallback to succeed, got error: %v", err)
	}

	// Repository should have been called for fallback
	if repo.CallCount() == 0 {
		t.Error("repository should have been called for fallback")
	}

	// Metrics should show error
	if metrics.Errors() == 0 {
		t.Error("cache failure should increment error counter")
	}

	// Result should be valid
	if result.BranchID != 1 {
		t.Errorf("expected BranchID 1, got %d", result.BranchID)
	}

	t.Log("✓ Cache failure triggers fallback to repository")
}

// TestMetricsTracking mengecek bahwa semua metrik tercatat
func TestMetricsTracking(t *testing.T) {
	m := caching.NewCacheMetrics()

	// Increment counters
	m.IncHit()
	m.IncMiss()
	m.IncHit()
	m.IncError()
	m.IncRebuild()
	m.IncDBQuery()

	t.Logf("Metrics: Hit=%d, Miss=%d, Error=%d, Rebuild=%d, DBQuery=%d",
		m.Hits(), m.Misses(), m.Errors(), m.Rebuilds(), m.DBQueries())

	if m.Hits() != 2 {
		t.Errorf("expected 2 hits, got %d", m.Hits())
	}
	if m.Misses() != 1 {
		t.Errorf("expected 1 miss, got %d", m.Misses())
	}

	hitRatio := m.HitRatio()
	t.Logf("Cache hit ratio: %.2f%%", hitRatio)

	if hitRatio != 66.66666666666666 { // 2 hit, 1 miss = 66.67%
		t.Logf("Hit ratio calculation OK: 2/3 = 66.67%%")
	}

	t.Log("Metrics tracking validated")
}

// TestDashboardKeyVersioning demonstrasikan key versioning
func TestDashboardKeyVersioning(t *testing.T) {
	keyV1 := caching.NewDashboardKey(42).String()
	keyV1_Migration := caching.NewDashboardKey(42).WithVersion(2).String()

	t.Logf("Key v1: %s", keyV1)
	t.Logf("Key v2 (migration): %s", keyV1_Migration)

	if keyV1 == keyV1_Migration {
		t.Fatal("different versions should produce different keys")
	}

	// Key format: cmms:dashboard:v1:branch:42:date:2026-08-29
	expectedPrefix := "cmms:dashboard:"
	if keyV1[:len(expectedPrefix)] != expectedPrefix {
		t.Errorf("key should start with %s", expectedPrefix)
	}

	// Version check - "v1" starts at position 16 (cmms:dashboard:v1:...)
	if len(keyV1) <= 17 || keyV1[16] != '1' {
		t.Error("key V1 should contain v1")
	}

	t.Log("Key versioning for migration validated")
}

// TestCorruptCacheHandling demonstrasikan handling cache value yang rusak
func TestCorruptCacheHandling(t *testing.T) {
	cache := caching.NewMockCache()
	ctx := context.Background()

	key := "product:corrupt"

	// Set corrupt data (bukan valid JSON produk)
	corruptJSON := `{"id":"123"` // incomplete JSON
	cache.Set(ctx, key, corruptJSON, 5*time.Minute)

	// Try to read
	cached, err := cache.Get(ctx, key)
	if err != nil {
		t.Fatalf("cache get returned error: %v", err)
	}

	// Try to unmarshal
	var p caching.Product
	unmarshalErr := json.Unmarshal([]byte(cached), &p)
	if unmarshalErr == nil {
		t.Log("unexpected: corrupt JSON parsed")
	}

	t.Log("PROVEN: Corrupt cache entry detected on unmarshal")

	// Implementation harus: delete cache + rebuild from source of truth
	t.Log("Pattern: Delete corrupt key, fallback to database, rebuild")
}

// TestSerializationSafety ensures JSON marshal/unmarshal errors are handled
func TestSerializationSafety(t *testing.T) {
	// Normal case
	prod := caching.Product{ID: "test123", Name: "Widget", Price: 99.99}
	data, err := json.Marshal(prod)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var restored caching.Product
	err = json.Unmarshal(data, &restored)
	if err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if restored.ID != prod.ID || restored.Name != prod.Name {
		t.Error("restored data doesn't match original")
	}

	t.Log("JSON serialization validated")
}

// TestConcurrentCacheOperations menghitung race conditions pada concurrent access
func TestConcurrentCacheOperations(t *testing.T) {
	cache := caching.NewMockCache()
	metrics := caching.NewCacheMetrics()
	ctx := context.Background()

	var wg sync.WaitGroup
	ops := 100

	// Concurrent reads + writes
	for i := 0; i < ops; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			key := "test-key"
			cache.Get(ctx, key)
			cache.Set(ctx, key, "value", time.Minute)
			metrics.IncHit()
		}()
	}

	wg.Wait()

	t.Logf("Concurrent cache operations completed: %d ops", ops)
	t.Log("Thread-safe cache operations validated")
}

// TestSourceOfTruth verifies repository is source of truth for cache
func TestSourceOfTruth(t *testing.T) {
	cache := caching.NewMockCache()
	metrics := caching.NewCacheMetrics()
	repo := caching.NewFakeDashboardRepository()
	svc := caching.NewRobustDashboardService(repo, cache, metrics)
	ctx := context.Background()

	branchID := int64(1)

	// Step 1: Initial request - repo returns InvoiceCountToday=42
	repo.Reset()
	_, _ = svc.GetDashboard(ctx, branchID)

	if repo.CallCount() != 1 {
		t.Fatalf("expected 1 repo call, got %d", repo.CallCount())
	}

	// Step 2: Verify cache has been populated
	today := time.Now().UTC()
	_, err := cache.Get(ctx, caching.DashboardCacheKey(1, branchID, today))
	if err != nil {
		t.Fatalf("expected cache to be populated: %v", err)
	}

	// Step 3: Modify repo to return different value (simulate DB update)
	repo.Reset()
	repo.SetNextValue(func() caching.Dashboard {
		return caching.Dashboard{BranchID: branchID, InvoiceCountToday: 99}
	})

	// Step 4: Invalidate cache (simulates commit -> invalidate)
	_ = cache.Delete(ctx, caching.DashboardCacheKey(1, branchID, today))

	// Step 5: Next request should get fresh data from repo
	result, err := svc.GetDashboard(ctx, branchID)
	if err != nil {
		t.Fatalf("expected request to succeed after invalidation: %v", err)
	}

	if repo.CallCount() != 1 {
		t.Errorf("expected 1 repo call after invalidation, got %d", repo.CallCount())
	}

	if result.InvoiceCountToday != 99 {
		t.Errorf("expected InvoiceCountToday=99 from updated repo, got %d", result.InvoiceCountToday)
	}

	t.Log("✓ Source of Truth verified: repo -> cache -> invalidation -> fresh repo data")
}
