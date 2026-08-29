package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	caching "github.com/lukman-ss/software-engineering-lab/labs/04-caching"
)

func main() {
	scenario := flag.String("scenario", "", "Scenario to run: without-cache, cache-aside, stampede-unprotected, stampede-protected")
	flag.Parse()

	if *scenario == "" {
		fmt.Println("Usage: go run ./cmd/demo -scenario=[without-cache|cache-aside|stampede-unprotected|stampede-protected]")
		fmt.Println("\nAvailable Scenarios:")
		fmt.Println("  without-cache         : Baseline: 100 requests = 100 repository calls")
		fmt.Println("  cache-aside           : Standard cache: 100 requests = 1 repo call, 99 cache hits")
		fmt.Println("  stampede-unprotected  : Cache expire + concurrent requests = DB overloaded")
		fmt.Println("  stampede-protected    : Cache expire + singleflight = 1 DB rebuild")
		os.Exit(1)
	}

	ctx := context.Background()

	switch *scenario {
	case "without-cache":
		runWithoutCache(ctx)
	case "cache-aside":
		runCacheAside(ctx)
	case "stampede-unprotected":
		runStampedeUnprotected(ctx)
	case "stampede-protected":
		runStampedeProtected(ctx)
	default:
		log.Fatalf("Unknown scenario: %s", *scenario)
	}
}

// Scenario 1: without-cache (baseline demo)
// Simulates 100 direct repository calls with NO caching
func runWithoutCache(ctx context.Context) {
	fmt.Println("=== SCENARIO: without-cache (baseline demo) ===")
	fmt.Println("Explanation: This simulation shows the pattern WITHOUT any caching layer.")
	fmt.Println("Every request goes directly to the repository/database.")

	repo := caching.NewFakeDashboardRepository()

	const reqCount = 100
	fmt.Printf("Simulating %d sequential requests to dashboard...\n", reqCount)

	// Direct repository calls - no caching
	for i := 0; i < reqCount; i++ {
		_, _ = repo.GetDashboard(ctx, 1, time.Now())
	}

	fmt.Printf("\nResults:\n")
	fmt.Printf("Requests: %d\n", reqCount)
	fmt.Printf("Repository Calls: %d\n", repo.CallCount)
	fmt.Printf("Cache Hits: 0\n")
	fmt.Printf("Cache Misses: 0\n")
	fmt.Printf("\nConclusion: Without cache, every request puts full load on database.\n")
}

// Scenario 2: cache-aside
// First request misses and populates cache, subsequent requests hit
func runCacheAside(ctx context.Context) {
	fmt.Println("=== SCENARIO: cache-aside ===")
	cache := caching.NewMockCacheWithStats()
	metrics := caching.NewCacheMetrics()

	repo := caching.NewFakeDashboardRepository()
	svc := caching.NewRobustDashboardService(repo, cache, metrics)

	const reqCount = 100
	fmt.Printf("Simulating %d sequential requests to dashboard...\n", reqCount)

	for i := 0; i < reqCount; i++ {
		_, _ = svc.GetDashboard(ctx, 1)
	}

	fmt.Printf("\nResults:\n")
	fmt.Printf("Requests: %d\n", reqCount)
	fmt.Printf("Repository Calls: %d\n", repo.CallCount)
	fmt.Printf("Cache Hits: %d\n", metrics.Hits())
	fmt.Printf("Cache Misses: %d\n", metrics.Misses())
	fmt.Printf("\nConclusion: Cache Aside absorbs %d%% of the traffic, protecting the DB.\n",
		int(metrics.HitRatio()))
}

// Scenario 3: stampede-unprotected
// Multiple concurrent requests when cache is empty hit DB multiple times
func runStampedeUnprotected(ctx context.Context) {
	fmt.Println("=== SCENARIO: stampede-unprotected ===")
	fmt.Println("Context: Cache key expired/empty. 100 users request dashboard simultaneously.")

	cache := caching.NewMockCache()
	repo := caching.NewCounterRepository()
	svc := caching.NewBrokenStampedeService(cache, repo)

	const concurrentReqs = 100
	fmt.Printf("Simulating %d concurrent requests to dashboard...\n", concurrentReqs)

	var wg sync.WaitGroup
	for i := 0; i < concurrentReqs; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = svc.GetData(ctx, 1)
		}()
	}

	wg.Wait()

	fmt.Printf("\nResults:\n")
	fmt.Printf("Concurrent Requests: %d\n", concurrentReqs)
	fmt.Printf("Repository Calls: %d\n", repo.CallCount())
	fmt.Printf("\nConclusion: Missing protection during expiration causes DB overload (Stampede/Thundering Herd).\n")
}

// Scenario 4: stampede-protected
// Singleflight deduplicates concurrent requests into one DB call
func runStampedeProtected(ctx context.Context) {
	fmt.Println("=== SCENARIO: stampede-protected ===")
	fmt.Println("Context: Cache key expired/empty. 100 users request dashboard simultaneously.")
	fmt.Println("Protection: singleflight")

	cache := caching.NewMockCache()
	repo := caching.NewCounterRepository()
	svc := caching.NewProtectedStampedeService(cache, repo)

	const concurrentReqs = 100
	fmt.Printf("Simulating %d concurrent requests to dashboard...\n", concurrentReqs)

	var wg sync.WaitGroup
	for i := 0; i < concurrentReqs; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = svc.GetData(ctx, 2)
		}()
	}

	wg.Wait()

	fmt.Printf("\nResults:\n")
	fmt.Printf("Concurrent Requests: %d\n", concurrentReqs)
	fmt.Printf("Repository Calls: %d\n", repo.CallCount())
	fmt.Printf("\nConclusion: Singleflight deduplicates concurrent requests into a single DB query.\n")
}
