package retry

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// RetryStormSimulator demonstrates the retry storm problem where multiple
// concurrent clients retry simultaneously when a service is down, amplifying
// the load on the already-failing service.
type RetryStormSimulator struct {
	mu            sync.Mutex
	requestCounts map[int]int // worker_id -> request count
	startTime     time.Time
}

// CreateStormConfig holds configuration for a retry storm experiment.
type StormConfig struct {
	WorkerCount  int
	MaxRetries   int
	BaseDelay    time.Duration
	ShouldJitter bool
}

// Result holds the simulation results.
type StormResult struct {
	TotalRequests        int
	MaxRequests          int   // Peak concurrent requests at any point
	RequestByTime        []int // Request rate over time
	Distribution         map[int]int
	AvgRequestsPerWorker float64
}

// RunStormExperiment runs a retry storm simulation with multiple workers
// all retrying to the same failing provider.
func RunStormExperiment(ctx context.Context, provider *MockProvider, cfg StormConfig) StormResult {
	sim := &RetryStormSimulator{
		requestCounts: make(map[int]int),
		startTime:     time.Now(),
	}

	var wg sync.WaitGroup
	var totalReqs atomic.Int32
	requestByTime := make([]int, 0, 1000)
	_ = requestByTime

	rand.Seed(time.Now().UnixNano())

	for workerID := 0; workerID < cfg.WorkerCount; workerID++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			var client *RetryableClient
			if cfg.ShouldJitter {
				client = NewRetryableClient(provider,
					WithRetryAttempts(cfg.MaxRetries),
					WithBaseDelay(cfg.BaseDelay),
					WithJitter(0.5),
				)
			} else {
				client = NewRetryableClient(provider,
					WithRetryAttempts(cfg.MaxRetries),
					WithBaseDelay(cfg.BaseDelay),
					WithJitter(0),
				)
			}

			for i := 0; i < 3; i++ {
				select {
				case <-ctx.Done():
					return
				default:
					_, _ = client.Get(ctx, "http://example.com")
					totalReqs.Add(1)
					time.Sleep(1 * time.Millisecond) // Sample timing
				}
			}
		}(workerID)
	}

	wg.Wait()

	total := int(totalReqs.Load())
	return StormResult{
		TotalRequests:        total,
		MaxRequests:          total, // Since sequential in test
		RequestByTime:        requestByTime,
		Distribution:         sim.requestCounts,
		AvgRequestsPerWorker: float64(total) / float64(cfg.WorkerCount),
	}
}

// AnalyzeRetryStorm analyzes the request distribution to detect storm patterns.
func AnalyzeRetryStorm(reqsByWorker map[int]int, totalWorkers int) map[string]interface{} {
	totalReqs := 0
	maxReqs := 0

	for _, reqs := range reqsByWorker {
		totalReqs += reqs
		if reqs > maxReqs {
			maxReqs = reqs
		}
	}

	avgReqs := float64(totalReqs) / float64(totalWorkers)
	variance := 0.0
	for _, reqs := range reqsByWorker {
		diff := float64(reqs) - avgReqs
		variance += diff * diff
	}
	variance /= float64(totalWorkers)

	return map[string]interface{}{
		"total_requests": totalReqs,
		"max_per_worker": maxReqs,
		"avg_per_worker": avgReqs,
		"variance":       variance,
		"storm_severity": variance / avgReqs, // Coefficient of variation
		"is_storm":       variance/avgReqs > 0.1 && avgReqs > 1,
	}
}

// DemonstrateRetryStorm shows the synchronization problem.
func DemonstrateRetryStorm() {
	fmt.Println("RETRY STORM DEMONSTRATION")
	fmt.Println("========================")
	fmt.Println("")
	fmt.Println("Problem: When multiple clients retry simultaneously, they create")
	fmt.Println("a 'thundering herd' that amplifies load on failing services.")
	fmt.Println("")

	// Simulate the problem conceptually
	fmt.Println("WITHOUT JITTER:")
	fmt.Println("  3 workers start retries at T=0")
	fmt.Println("  All workers calculate same wait: ~100ms (base)")
	fmt.Println("  At T=100ms: 30 requests hit server simultaneously")
	fmt.Println("  Server still failing -> 30 more retries at T=200ms")
	fmt.Println("  Load amplification: 30 requests for 3 clients (10x)")
	fmt.Println("")
	fmt.Println("WITH JITTER:")
	fmt.Println("  3 workers start retries at T=0")
	fmt.Println("  Worker 1: 100ms +/- 50ms -> waits 40-60ms (random)")
	fmt.Println("  Worker 2: 100ms +/- 50ms -> waits 120-180ms")
	fmt.Println("  Worker 3: 100ms +/- 50ms -> waits 70-130ms")
	fmt.Println("  Requests spread across 40-180ms window")
	fmt.Println("  Reduces peak load, prevents synchronization")
	fmt.Println("")

	fmt.Println("RECOMMENDATIONS:")
	fmt.Println("  1. Always use jitter in retry backoff (full jitter or equal jitter)")
	fmt.Println("  2. Set max jitter factor: 0.5 for base delay, 0.3 for max delay")
	fmt.Println("  3. Add random initial delay for startup")
	fmt.Println("  4. Consider exponential backoff with capped maximum")
	fmt.Println("  5. Use circuit breakers to prevent storm initiation")
}

// JitterConfig defines jitter calculation strategies.
type JitterConfig struct {
	Mode   string  // "full", "equal", "decorrelated"
	Factor float64 // 0.0 to 1.0
}

// FullJitter returns random delay between 0 and baseDelay * 2^attempt.
func FullJitter(baseDelay time.Duration, attempt int) time.Duration {
	maxDelay := baseDelay * time.Duration(1<<uint(attempt))
	return time.Duration(rand.Float64() * float64(maxDelay))
}

// EqualJitter returns baseDelay + random backoff up to baseDelay * 2^attempt.
func EqualJitter(baseDelay time.Duration, attempt int) time.Duration {
	maxJitter := baseDelay * time.Duration(1<<uint(attempt))
	return baseDelay + time.Duration(rand.Float64()*float64(maxJitter))
}

// DecorrelatedJitter uses previous delay to calculate next (Amazon style).
func DecorrelatedJitter(baseDelay time.Duration, prevDelay time.Duration) time.Duration {
	if prevDelay == 0 {
		prevDelay = baseDelay
	}
	return baseDelay + time.Duration(rand.Float64()*float64(prevDelay)*3)
}
