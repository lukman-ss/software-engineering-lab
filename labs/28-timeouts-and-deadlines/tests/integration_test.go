package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"timeouts-and-deadlines/internal/circuit"
	"timeouts-and-deadlines/internal/idempotency"
	"timeouts-and-deadlines/internal/retry"
)

func TestIntegration_RetryWithCircuitBreaker(t *testing.T) {
	cb := circuit.NewBreaker(circuit.Config{
		FailureThreshold: 2,
		SuccessThreshold: 1,
		Cooldown:         50 * time.Millisecond,
	})

	r := retry.NewRetrier(retry.Config{
		MaxAttempts: 4,
		BaseBackoff: 5 * time.Millisecond,
		MaxBackoff:  20 * time.Millisecond,
	})

	attempts := 0
	err := r.Do(context.Background(), func(ctx context.Context) error {
		return cb.Execute(func() error {
			attempts++
			return errors.New("backend down")
		})
	})

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if cb.State() != circuit.StateOpen {
		t.Fatalf("expected circuit open, got %s", cb.State())
	}
}

func TestIntegration_IdempotentRetry(t *testing.T) {
	store := idempotency.NewStore(time.Minute)
	r := retry.NewRetrier(retry.Config{
		MaxAttempts: 3,
		BaseBackoff: 2 * time.Millisecond,
		MaxBackoff:  10 * time.Millisecond,
	})

	const idKey = "unique-order-key"
	actualExecutions := 0

	process := func() (string, error) {
		if val, ok := store.Get(idKey); ok {
			return val, nil
		}
		actualExecutions++
		res := "order-created"
		store.Set(idKey, res)
		return res, nil
	}

	// 1st run
	res, err := process()
	if err != nil || res != "order-created" {
		t.Fatalf("unexpected result: %v, %s", err, res)
	}

	// Retry loop invoking process
	err = r.Do(context.Background(), func(ctx context.Context) error {
		_, innerErr := process()
		return innerErr
	})

	if err != nil {
		t.Fatalf("expected retries to succeed, got %v", err)
	}

	if actualExecutions != 1 {
		t.Fatalf("expected exactly 1 actual execution due to deduplication, got %d", actualExecutions)
	}
}
