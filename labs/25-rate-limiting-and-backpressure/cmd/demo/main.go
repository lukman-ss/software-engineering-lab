package main

import (
	"context"
	"fmt"
	"time"

	"labs/25-rate-limiting-and-backpressure/internal/backpressure"
	"labs/25-rate-limiting-and-backpressure/internal/ratelimit"
	"labs/25-rate-limiting-and-backpressure/internal/retry"
)

func main() {
	fmt.Println("=== 1. Token Bucket Burst & Rate Limiting ===")
	tb := ratelimit.NewTokenBucket(3, 5) // Cap 3, refill 5/s
	for i := 1; i <= 5; i++ {
		allowed := tb.Allow()
		fmt.Printf("Request #%d: Allowed=%v (Remaining Tokens: %.1f)\n", i, allowed, tb.Tokens())
	}

	time.Sleep(300 * time.Millisecond)
	fmt.Printf("After 300ms pause: Allowed=%v (Remaining Tokens: %.1f)\n", tb.Allow(), tb.Tokens())

	fmt.Println("\n=== 2. Leaky Bucket Traffic Smoothing ===")
	lb := ratelimit.NewLeakyBucket(3, 10) // Cap 3, leak 10/s
	for i := 1; i <= 5; i++ {
		allowed := lb.Allow()
		fmt.Printf("Request #%d: Allowed=%v (Current Water Level: %.1f)\n", i, allowed, lb.Water())
	}

	fmt.Println("\n=== 3. Bounded Queue Backpressure (Load Shedding) ===")
	bq := backpressure.NewBoundedQueue(3, 1) // 3 queue cap, 1 worker
	defer bq.Stop()

	for i := 1; i <= 6; i++ {
		jobID := i
		err := bq.TrySubmit(func(ctx context.Context) error {
			time.Sleep(50 * time.Millisecond)
			return nil
		})
		if err != nil {
			fmt.Printf("Job #%d: REJECTED (Backpressure Shedding: %v)\n", jobID, err)
		} else {
			fmt.Printf("Job #%d: ACCEPTED into bounded buffer\n", jobID)
		}
	}

	time.Sleep(100 * time.Millisecond)
	accepted, rejected, processed, _ := bq.Stats()
	fmt.Printf("Stats: Accepted=%d, Rejected=%d, Processed=%d\n", accepted, rejected, processed)

	fmt.Println("\n=== 4. AWS Retry Backoff Strategies (Attempts 0..3) ===")
	cfg := retry.Config{
		Base: 100 * time.Millisecond,
		Cap:  1000 * time.Millisecond,
	}

	for attempt := 0; attempt < 4; attempt++ {
		noJitter := retry.ComputeBackoff(retry.NoJitter, attempt, cfg, 0)
		fullJitter := retry.ComputeBackoff(retry.FullJitter, attempt, cfg, 0)
		equalJitter := retry.ComputeBackoff(retry.EqualJitter, attempt, cfg, 0)
		fmt.Printf("Attempt %d -> NoJitter: %-6v | FullJitter: %-6v | EqualJitter: %-6v\n",
			attempt, noJitter.Round(time.Millisecond), fullJitter.Round(time.Millisecond), equalJitter.Round(time.Millisecond))
	}
}
