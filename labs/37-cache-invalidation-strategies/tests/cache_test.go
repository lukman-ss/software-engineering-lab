package tests

import (
	"context"
	"math"
	"sync"
	"testing"
	"time"

	"github.com/lukman/labs/37-cache-invalidation-strategies/internal/cache"
)

func TestCachePatterns(t *testing.T) {
	ctx := context.Background()

	t.Run("Cache-Aside Read & Write", func(t *testing.T) {
		mem := cache.NewMemoryCache()
		db := cache.NewMockDB(0)
		db.SetData("k1", "v1")

		svc := cache.NewCacheAsideService(mem, db, 1*time.Minute)

		// 1. Initial read misses cache, queries DB
		val, err := svc.Get(ctx, "k1")
		if err != nil || val != "v1" {
			t.Fatalf("expected v1, got %v (err %v)", val, err)
		}
		if db.QueryCount() != 1 {
			t.Fatalf("expected 1 query, got %d", db.QueryCount())
		}

		// 2. Second read hits cache
		val, err = svc.Get(ctx, "k1")
		if err != nil || val != "v1" {
			t.Fatalf("expected v1, got %v", val)
		}
		if db.QueryCount() != 1 {
			t.Fatalf("expected still 1 query, got %d", db.QueryCount())
		}

		// 3. Update writes to DB and invalidates cache
		if err := svc.Update(ctx, "k1", "v2"); err != nil {
			t.Fatalf("update failed: %v", err)
		}

		// 4. Next read queries DB again with new value
		val, err = svc.Get(ctx, "k1")
		if err != nil || val != "v2" {
			t.Fatalf("expected v2, got %v", val)
		}
		if db.QueryCount() != 2 {
			t.Fatalf("expected 2 queries, got %d", db.QueryCount())
		}
	})

	t.Run("Write-Through Read & Write", func(t *testing.T) {
		mem := cache.NewMemoryCache()
		db := cache.NewMockDB(0)
		db.SetData("k1", "v1")

		svc := cache.NewWriteThroughService(mem, db, 1*time.Minute)

		// 1. Read populates cache
		_, _ = svc.Get(ctx, "k1")
		if db.QueryCount() != 1 {
			t.Fatalf("expected 1 query, got %d", db.QueryCount())
		}

		// 2. Update writes to DB and cache synchronously
		if err := svc.Update(ctx, "k1", "v-updated"); err != nil {
			t.Fatalf("update failed: %v", err)
		}

		// 3. Subsequent read directly hits updated cache (no DB query)
		val, err := svc.Get(ctx, "k1")
		if err != nil || val != "v-updated" {
			t.Fatalf("expected v-updated, got %v", val)
		}
		if db.QueryCount() != 1 {
			t.Fatalf("expected write-through to cache without query, got %d queries", db.QueryCount())
		}
	})

	t.Run("Write-Behind Asynchronous Flush", func(t *testing.T) {
		mem := cache.NewMemoryCache()
		db := cache.NewMockDB(0)
		svc := cache.NewWriteBehindService(mem, db, 1*time.Minute, 10)
		defer svc.Close()

		svc.Update("wb-key", "wb-val")

		// Cache is updated immediately
		val, err := svc.Get(ctx, "wb-key")
		if err != nil || val != "wb-val" {
			t.Fatalf("expected immediate cache read wb-val, got %v", val)
		}

		// DB is updated asynchronously
		time.Sleep(50 * time.Millisecond)
		if db.WriteCount() != 1 {
			t.Fatalf("expected DB write count 1, got %d", db.WriteCount())
		}
	})
}

func TestStampedeMitigation(t *testing.T) {
	ctx := context.Background()

	t.Run("Naive Stampede Queries DB Concurrently", func(t *testing.T) {
		mem := cache.NewMemoryCache()
		db := cache.NewMockDB(10 * time.Millisecond)
		db.SetData("hot-key", "hot-value")
		svc := cache.NewNaiveStampedeService(mem, db, 1*time.Minute)

		concurrency := 20
		var wg sync.WaitGroup
		wg.Add(concurrency)

		for i := 0; i < concurrency; i++ {
			go func() {
				defer wg.Done()
				_, _ = svc.Get(ctx, "hot-key")
			}()
		}
		wg.Wait()

		// Naive service causes multiple simultaneous queries on a miss
		if db.QueryCount() <= 1 {
			t.Fatalf("expected stampede with > 1 query, got %d", db.QueryCount())
		}
	})

	t.Run("SingleFlight Coalesces To Single Query", func(t *testing.T) {
		mem := cache.NewMemoryCache()
		db := cache.NewMockDB(20 * time.Millisecond)
		db.SetData("hot-key", "hot-value")
		svc := cache.NewSingleFlightService(mem, db, 1*time.Minute)

		concurrency := 20
		var wg sync.WaitGroup
		wg.Add(concurrency)

		for i := 0; i < concurrency; i++ {
			go func() {
				defer wg.Done()
				val, err := svc.Get(ctx, "hot-key")
				if err != nil || val != "hot-value" {
					t.Errorf("expected hot-value, got %v, err: %v", val, err)
				}
			}()
		}
		wg.Wait()

		// Singleflight coalesces all 20 calls to 1 query
		if db.QueryCount() != 1 {
			t.Fatalf("expected exactly 1 DB query with singleflight, got %d", db.QueryCount())
		}
	})
}

func TestXFetchLogic(t *testing.T) {
	// Formula: -delta * beta * ln(u) > ttl_remaining
	delta := 100 * time.Millisecond
	beta := 1.0
	remaining := 50 * time.Millisecond

	// Test 1: High remaining TTL -> Should not recompute
	// u = 0.5 -> -ln(0.5) = 0.6931 -> delta * 0.6931 = 69.31ms < 200ms
	if cache.ShouldRecompute(delta, beta, 200*time.Millisecond, 0.5) {
		t.Fatalf("should not recompute when remaining TTL is large")
	}

	// Test 2: Low remaining TTL -> Recompute triggered
	// u = 0.3 -> -ln(0.3) = 1.2039 -> delta * 1.2039 = 120.39ms > 50ms
	if !cache.ShouldRecompute(delta, beta, remaining, 0.3) {
		t.Fatalf("should recompute when -delta*beta*ln(u) > remaining")
	}

	// Test 3: Verify erroneous formula (without minus sign) produces false
	erroneousResult := (delta.Seconds() * beta * math.Log(0.5)) > remaining.Seconds()
	if erroneousResult {
		t.Fatalf("erroneous formula sign unexpectedly yielded true")
	}
}

func TestStaleWhileRevalidate(t *testing.T) {
	ctx := context.Background()
	mem := cache.NewMemoryCache()
	db := cache.NewMockDB(5 * time.Millisecond)
	db.SetData("swr-key", "initial-data")

	// TTL 10ms, Stale Delta 200ms
	svc := cache.NewSWRService(mem, db, 20*time.Millisecond, 300*time.Millisecond)

	// Initial fetch
	val, err := svc.Get(ctx, "swr-key")
	if err != nil || val != "initial-data" {
		t.Fatalf("expected initial-data, got %v", val)
	}

	// Wait for TTL to expire, but within stale delta
	time.Sleep(30 * time.Millisecond)

	// Update DB behind the scenes
	db.SetData("swr-key", "updated-data")

	// Read should return stale "initial-data" immediately, while triggering async revalidate
	val, err = svc.Get(ctx, "swr-key")
	if err != nil || val != "initial-data" {
		t.Fatalf("expected stale data 'initial-data', got %v", val)
	}

	// Wait for background revalidation
	time.Sleep(50 * time.Millisecond)

	// Now cache should have "updated-data"
	val, err = svc.Get(ctx, "swr-key")
	if err != nil || val != "updated-data" {
		t.Fatalf("expected updated data 'updated-data', got %v", val)
	}
}

func TestJitter(t *testing.T) {
	base := 10 * time.Second
	jitter := 2 * time.Second

	for i := 0; i < 100; i++ {
		res := cache.TTLWithJitter(base, jitter)
		if res < base || res >= base+jitter {
			t.Fatalf("jitter result %v out of range [%v, %v)", res, base, base+jitter)
		}
	}
}

func TestCachePatterns_FailurePaths(t *testing.T) {
	ctx := context.Background()

	t.Run("Cache-Aside DB Read Error", func(t *testing.T) {
		mem := cache.NewMemoryCache()
		db := cache.NewMockDB(0) // empty DB produces ErrNotFound
		svc := cache.NewCacheAsideService(mem, db, 1*time.Minute)

		val, err := svc.Get(ctx, "nonexistent")
		if err == nil || val != "" {
			t.Fatalf("expected error on nonexistent key, got val=%q, err=%v", val, err)
		}
	})

	t.Run("Write-Through DB Read Error", func(t *testing.T) {
		mem := cache.NewMemoryCache()
		db := cache.NewMockDB(0)
		svc := cache.NewWriteThroughService(mem, db, 1*time.Minute)

		val, err := svc.Get(ctx, "nonexistent")
		if err == nil || val != "" {
			t.Fatalf("expected error on nonexistent key, got val=%q, err=%v", val, err)
		}
	})
}

func TestXFetchService_Get(t *testing.T) {
	ctx := context.Background()
	mem := cache.NewMemoryCache()
	db := cache.NewMockDB(50 * time.Millisecond)
	db.SetData("xf-key", "initial-val")

	svc := cache.NewXFetchService(mem, db, 1*time.Second, 1.0)

	// Inject deterministic randFunc returning 0.99 (high u -> no early recompute)
	svc.SetRandFunc(func() float64 { return 0.99 })

	// 1. Initial read misses cache, queries DB
	val, err := svc.Get(ctx, "xf-key")
	if err != nil || val != "initial-val" {
		t.Fatalf("expected initial-val, got %v (err %v)", val, err)
	}
	if db.QueryCount() != 1 {
		t.Fatalf("expected 1 query on miss, got %d", db.QueryCount())
	}

	// 2. Immediate second read hits cache without recompute
	val, err = svc.Get(ctx, "xf-key")
	if err != nil || val != "initial-val" {
		t.Fatalf("expected initial-val from cache, got %v", val)
	}
	if db.QueryCount() != 1 {
		t.Fatalf("expected still 1 query, got %d", db.QueryCount())
	}

	// Update DB value
	db.SetData("xf-key", "refreshed-val")

	// Inject deterministic randFunc returning small u (0.00001) -> -50ms * 1.0 * ln(0.00001) = 575ms > ~950ms remaining? Wait, 575ms > remaining only when remaining < 575ms.
	// Let's set u extremely small (1e-12 -> -50ms * ln(1e-12) = 1381ms > 950ms remaining).
	svc.SetRandFunc(func() float64 { return 1e-12 })

	// 3. Read triggers early recompute due to low u
	val, err = svc.Get(ctx, "xf-key")
	if err != nil || val != "refreshed-val" {
		t.Fatalf("expected refreshed-val after early recompute, got %v", val)
	}
	if db.QueryCount() != 2 {
		t.Fatalf("expected 2 queries after proactive recompute, got %d", db.QueryCount())
	}
}

func TestWriteBehindService_QueueOverflow(t *testing.T) {
	mem := cache.NewMemoryCache()
	db := cache.NewMockDB(50 * time.Millisecond) // slow DB to keep worker busy
	bufSize := 2
	svc := cache.NewWriteBehindService(mem, db, 1*time.Minute, bufSize)
	defer svc.Close()

	// Fill queue past buffer capacity (1 in flight, 2 queued, remaining dropped)
	for i := 0; i < 10; i++ {
		svc.Update("overflow-key", "overflow-val")
	}

	// Immediate cache read reflects latest update
	ctx := context.Background()
	val, err := svc.Get(ctx, "overflow-key")
	if err != nil || val != "overflow-val" {
		t.Fatalf("expected immediate cache read overflow-val, got %v", val)
	}
}
