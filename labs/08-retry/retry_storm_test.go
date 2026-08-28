package retry_test

import (
	"context"
	"testing"
	"time"

	"github.com/lukman-ss/software-engineering-lab/labs/08-retry"
)

func TestRetryStormExperiment(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Provider that always fails (simulating down service)
	provider := retry.NewMockProvider().
		WithMaxFailures(100). // Fail every request
		WithStatusCode(retry.StatusInternalServerError)

	cfg := retry.StormConfig{
		WorkerCount:  10,
		MaxRetries:   3,
		BaseDelay:    10 * time.Millisecond,
		ShouldJitter: false,
	}

	result := retry.RunStormExperiment(ctx, provider, cfg)

	t.Logf("Retry Storm Experiment Results:")
	t.Logf("  Total Requests Generated: %d", result.TotalRequests)
	t.Logf("  Avg Requests Per Worker: %.2f", result.AvgRequestsPerWorker)

	if result.TotalRequests == 0 {
		t.Fatal("expected requests to be generated")
	}
}

func TestJitterStrategies(t *testing.T) {
	base := 100 * time.Millisecond

	// Test Full Jitter
	fj := retry.FullJitter(base, 2)
	if fj < 0 || fj > 400*time.Millisecond {
		t.Errorf("full jitter out of bounds: %v", fj)
	}

	// Test Equal Jitter
	eq := retry.EqualJitter(base, 2)
	if eq < base || eq > 500*time.Millisecond {
		t.Errorf("equal jitter out of bounds: %v", eq)
	}

	// Test Decorrelated Jitter
	dec := retry.DecorrelatedJitter(base, 200*time.Millisecond)
	if dec < base {
		t.Errorf("decorrelated jitter below base: %v", dec)
	}
}
