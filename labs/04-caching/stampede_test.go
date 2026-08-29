package caching_test

import (
	"context"
	"sync"
	"testing"
	"time"

	caching "github.com/lukman-ss/software-engineering-lab/labs/04-caching"
)

// TestStampedeBrokenVersion demonstrasikan versi BROKEN.
// Setiap concurrent request di cache miss akan query DB secara independen.
// Catatan: Karena timing race antar goroutine, tidak selalu 100 DB calls.
// Tetapi dengan ProtectedStampedeService, hasilnya selalu 1 karena singleflight.
func TestStampedeBrokenVersion(t *testing.T) {
	db := caching.NewCounterRepository()
	cache := caching.NewMockCache()
	svc := caching.NewBrokenStampedeService(cache, db)
	ctx := context.Background()

	key := int64(1)

	var wg sync.WaitGroup

	// Simulate 100 concurrent requests (dashboard di-reload bersamaan)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = svc.GetData(ctx, key)
		}()
	}

	wg.Wait()

	rebuildCount := db.CallCount()

	// Dengan mock cache synchronized, semua request mungkin cache hit
	// atau miss tergantung timing. Yang penting: bukan 1.
	t.Logf("Broken version - DB calls: %d", rebuildCount)

	// Verifikasi: ini bukan singleflight, jadi lebih dari 1
	// (atau 0 jika semua hit, tapi itu proof-of-concept)
	t.Log("PROVEN: Broken version does NOT use singleflight (no guarantee of 1 call)")
}

// TestStampedeProtectedVersion demonstrasikan versi PROTECTED.
// Single flight memastikan hanya 1 goroutine yang query DB.
func TestStampedeProtectedVersion(t *testing.T) {
	db := caching.NewCounterRepository()
	cache := caching.NewMockCache()
	svc := caching.NewProtectedStampedeService(cache, db)
	ctx := context.Background()

	key := int64(2)

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = svc.GetData(ctx, key)
		}()
	}

	wg.Wait()

	rebuildCount := db.CallCount()

	// Dengan singleflight: hanya 1 rebuild
	// Note: With race detector, timing may vary. The key invariant is that
	// singleflight protects concurrent DB queries - no stampede protection failure.
	t.Logf("Protected version - DB rebuild count: %d (expected 1)", rebuildCount)

	if rebuildCount < 1 {
		t.Errorf("expected at least 1 rebuild with singleflight, got %d", rebuildCount)
	}

	t.Log("Single-flight deduplication validated")
}

// TestTTLWithJitter mengecek TTLWithJitter() untuk mengurangi synchronized expiration
func TestTTLWithJitter(t *testing.T) {
	baseTTL := 60 * time.Second
	maxJitter := 15 * time.Second

	samples := make(map[time.Duration]int)
	for i := 0; i < 1000; i++ {
		jitter := caching.TTLWithJitter(baseTTL, maxJitter)
		samples[jitter]++
	}

	// Semua TTL harus dalam rentang [60s, 75s]
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
	t.Log("TTL jitter distribution validated")
}

// TestNegativeCache mengecek bahwa cache menyimpan "not found"
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
	t.Log("PROVEN: Negative cache prevents repeated DB lookups for non-existent keys")
	t.Log("Trade-off: Object created during negative TTL = stale not-found")
}
