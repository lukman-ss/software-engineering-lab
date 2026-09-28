package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/lukman/labs/37-cache-invalidation-strategies/internal/cache"
)

func main() {
	ctx := context.Background()

	fmt.Println("========================================")
	fmt.Println("Lab 37: Cache Invalidation Strategies")
	fmt.Println("========================================")

	demoCachePatterns(ctx)
	demoStampede(ctx)
	demoXFetch(ctx)
	demoSWR(ctx)
	demoJitter()

	fmt.Println("\n[Demo Complete]")
}

func demoCachePatterns(ctx context.Context) {
	fmt.Println("\n--- CACHE PATTERNS ---")

	// Cache-Aside
	{
		mem := cache.NewMemoryCache()
		db := cache.NewMockDB(5 * time.Millisecond)
		db.SetData("product:1", "Laptop")
		svc := cache.NewCacheAsideService(mem, db, 30*time.Second)

		fmt.Println("\n[Cache-Aside]")
		val, _ := svc.Get(ctx, "product:1")
		fmt.Printf("  1st GET (cache miss)  -> DB queries: %d, val: %s\n", db.QueryCount(), val)
		val, _ = svc.Get(ctx, "product:1")
		fmt.Printf("  2nd GET (cache hit)   -> DB queries: %d, val: %s\n", db.QueryCount(), val)
		_ = svc.Update(ctx, "product:1", "Laptop Pro")
		val, _ = svc.Get(ctx, "product:1")
		fmt.Printf("  After update & GET    -> DB queries: %d, val: %s\n", db.QueryCount(), val)
	}

	// Write-Through
	{
		mem := cache.NewMemoryCache()
		db := cache.NewMockDB(5 * time.Millisecond)
		db.SetData("product:2", "Phone")
		svc := cache.NewWriteThroughService(mem, db, 30*time.Second)

		fmt.Println("\n[Write-Through]")
		_, _ = svc.Get(ctx, "product:2")
		fmt.Printf("  1st GET (cache miss)  -> DB queries: %d\n", db.QueryCount())
		_ = svc.Update(ctx, "product:2", "Phone Pro")
		val, _ := svc.Get(ctx, "product:2")
		fmt.Printf("  After update & GET    -> DB queries: %d, val: %s\n", db.QueryCount(), val)
	}

	// Write-Behind
	{
		mem := cache.NewMemoryCache()
		db := cache.NewMockDB(5 * time.Millisecond)
		svc := cache.NewWriteBehindService(mem, db, 30*time.Second, 10)

		fmt.Println("\n[Write-Behind]")
		svc.Update("product:3", "Tablet")
		val, _ := svc.Get(ctx, "product:3")
		fmt.Printf("  After Update, GET from cache -> DB writes: %d, val: %s\n", db.WriteCount(), val)
		time.Sleep(50 * time.Millisecond)
		fmt.Printf("  After 50ms flush delay       -> DB writes: %d\n", db.WriteCount())
		svc.Close()
	}
}

func demoStampede(ctx context.Context) {
	fmt.Println("\n--- STAMPEDE DEMONSTRATION ---")

	// Naive
	{
		mem := cache.NewMemoryCache()
		db := cache.NewMockDB(20 * time.Millisecond)
		db.SetData("hot", "data")
		svc := cache.NewNaiveStampedeService(mem, db, 100*time.Millisecond)

		concurrency := 20
		var wg sync.WaitGroup
		wg.Add(concurrency)
		for i := 0; i < concurrency; i++ {
			go func() { defer wg.Done(); _, _ = svc.Get(ctx, "hot") }()
		}
		wg.Wait()
		fmt.Printf("\n[Naive] %d concurrent goroutines -> %d DB queries (stampede!)\n", concurrency, db.QueryCount())
	}

	// SingleFlight
	{
		mem := cache.NewMemoryCache()
		db := cache.NewMockDB(20 * time.Millisecond)
		db.SetData("hot", "data")
		svc := cache.NewSingleFlightService(mem, db, 100*time.Millisecond)

		concurrency := 20
		var (
			wg      sync.WaitGroup
			results atomic.Int64
		)
		wg.Add(concurrency)
		for i := 0; i < concurrency; i++ {
			go func() {
				defer wg.Done()
				val, _ := svc.Get(ctx, "hot")
				if val == "data" {
					results.Add(1)
				}
			}()
		}
		wg.Wait()
		fmt.Printf("[SF]    %d concurrent goroutines -> %d DB queries (singleflight), %d correct results\n",
			concurrency, db.QueryCount(), results.Load())
	}
}

func demoXFetch(ctx context.Context) {
	fmt.Println("\n--- XFETCH (Probabilistic Early Expiration) ---")

	mem := cache.NewMemoryCache()
	db := cache.NewMockDB(50 * time.Millisecond)
	db.SetData("popular", "hot-data")
	svc := cache.NewXFetchService(mem, db, 500*time.Millisecond, 1.0)

	// Initial load
	val, _ := svc.Get(ctx, "popular")
	fmt.Printf("  Initial GET       -> DB queries: %d, val: %s\n", db.QueryCount(), val)

	// Force-simulate: use rand that always triggers early refresh
	svc.SetRandFunc(func() float64 { return 1e-9 }) // -ln(1e-9) ≈ 20.7 -> always early trigger
	val, _ = svc.Get(ctx, "popular")
	fmt.Printf("  Low rand draw     -> DB queries: %d, val: %s (early proactive refresh)\n", db.QueryCount(), val)

	// Force-simulate: rand never triggers early refresh
	db2 := cache.NewMockDB(5 * time.Millisecond)
	db2.SetData("popular", "fresh-data")
	svc2 := cache.NewXFetchService(cache.NewMemoryCache(), db2, 500*time.Millisecond, 1.0)
	_, _ = svc2.Get(ctx, "popular")
	svc2.SetRandFunc(func() float64 { return 0.9999 }) // -ln(0.9999) ≈ 0.0001 -> never triggers early
	val2, _ := svc2.Get(ctx, "popular")
	fmt.Printf("  High rand draw    -> DB queries: %d, val: %s (served from cache, no early refresh)\n",
		db2.QueryCount(), val2)
}

func demoSWR(ctx context.Context) {
	fmt.Println("\n--- STALE-WHILE-REVALIDATE ---")

	mem := cache.NewMemoryCache()
	db := cache.NewMockDB(5 * time.Millisecond)
	db.SetData("swr-key", "v1")
	svc := cache.NewSWRService(mem, db, 20*time.Millisecond, 300*time.Millisecond)

	val, _ := svc.Get(ctx, "swr-key")
	fmt.Printf("  Initial GET       -> val: %s\n", val)

	time.Sleep(30 * time.Millisecond) // TTL expired
	db.SetData("swr-key", "v2")

	val, _ = svc.Get(ctx, "swr-key")
	fmt.Printf("  GET (stale)       -> val: %s  (stale served immediately)\n", val)

	time.Sleep(50 * time.Millisecond) // wait for async revalidation
	val, _ = svc.Get(ctx, "swr-key")
	fmt.Printf("  GET (fresh)       -> val: %s  (async revalidation complete)\n", val)
}

func demoJitter() {
	fmt.Println("\n--- TTL JITTER ---")
	base := 5 * time.Minute
	maxJitter := 30 * time.Second
	results := make([]time.Duration, 5)
	for i := range results {
		results[i] = cache.TTLWithJitter(base, maxJitter)
	}
	fmt.Printf("  Base TTL: %v, Max Jitter: %v\n", base, maxJitter)
	fmt.Printf("  Sampled TTLs: %v, %v, %v, %v, %v\n",
		results[0], results[1], results[2], results[3], results[4])
}
