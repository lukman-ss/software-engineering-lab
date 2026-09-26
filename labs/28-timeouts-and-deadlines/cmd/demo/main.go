package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"timeouts-and-deadlines/internal/circuit"
	"timeouts-and-deadlines/internal/deadline"
	"timeouts-and-deadlines/internal/idempotency"
	"timeouts-and-deadlines/internal/retry"
)

func main() {
	fmt.Println("=== Timeouts & Deadlines Laboratory Demo ===")

	// 1. Deadline propagation demo
	fmt.Println("\n--- Demo 1: Context Deadline & Budget Propagation ---")
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := deadline.ExecuteWithBudget(ctx, 100*time.Millisecond, func(childCtx context.Context) error {
		select {
		case <-time.After(80 * time.Millisecond):
			return nil
		case <-childCtx.Done():
			return childCtx.Err()
		}
	})
	fmt.Printf("Deadline propagation result: %v\n", err)

	// 2. Retry with full jitter demo
	fmt.Println("\n--- Demo 2: Exponential Backoff with Full Jitter ---")
	retrier := retry.NewRetrier(retry.Config{
		MaxAttempts: 3,
		BaseBackoff: 10 * time.Millisecond,
		MaxBackoff:  50 * time.Millisecond,
	})
	attemptCount := 0
	err = retrier.Do(context.Background(), func(ctx context.Context) error {
		attemptCount++
		fmt.Printf("  Attempt #%d executed\n", attemptCount)
		if attemptCount < 3 {
			return errors.New("transient network hiccup")
		}
		return nil
	})
	fmt.Printf("Retry execution result: %v (total attempts: %d)\n", err, attemptCount)

	// 3. Circuit breaker integration demo
	fmt.Println("\n--- Demo 3: Circuit Breaker State Transitions ---")
	cb := circuit.NewBreaker(circuit.Config{
		FailureThreshold: 2,
		SuccessThreshold: 1,
		Cooldown:         50 * time.Millisecond,
	})

	dummyErr := errors.New("service backend error")
	fmt.Printf("Initial state: %s\n", cb.State())

	cb.Execute(func() error { return dummyErr })
	cb.Execute(func() error { return dummyErr })
	fmt.Printf("State after 2 failures: %s\n", cb.State())

	err = cb.Execute(func() error { return nil })
	fmt.Printf("Execution attempt while OPEN: %v\n", err)

	time.Sleep(60 * time.Millisecond)
	fmt.Printf("State after cooldown: %s\n", cb.State())

	err = cb.Execute(func() error { return nil })
	fmt.Printf("Execution attempt in HALF_OPEN (success): %v -> New State: %s\n", err, cb.State())

	// 4. Idempotency deduplication demo
	fmt.Println("\n--- Demo 4: Idempotent Request Retry Protection ---")
	store := idempotency.NewStore(5 * time.Minute)
	idKey := "req-tx-99231"

	processPayment := func(key string, amount int) (string, error) {
		if res, ok := store.Get(key); ok {
			return fmt.Sprintf("%s (DEDUPLICATED)", res), nil
		}
		// Process payment logic...
		result := fmt.Sprintf("Charged $%d successfully", amount)
		store.Set(key, result)
		return result, nil
	}

	res1, _ := processPayment(idKey, 100)
	fmt.Printf("First request execution: %s\n", res1)

	res2, _ := processPayment(idKey, 100)
	fmt.Printf("Retried request execution: %s\n", res2)

	fmt.Println("\n=== Demo Completed Successfully ===")
}
