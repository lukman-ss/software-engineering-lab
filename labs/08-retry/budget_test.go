package retry_test

import (
	"context"
	"testing"
	"time"

	"github.com/lukman-ss/software-engineering-lab/labs/08-retry"
)

func TestRetryBudget(t *testing.T) {
	// Budget of 3 retries per second
	budget := retry.NewRetryBudget(3, time.Second)
	defer budget.Reset()

	// First 3 should succeed
	for i := 0; i < 3; i++ {
		if !budget.TryConsume() {
			t.Errorf("attempt %d should be allowed by budget", i+1)
		}
	}

	// 4th should be rejected
	if budget.TryConsume() {
		t.Error("4th attempt should be rejected by budget")
	}
}

func TestRetryBudgetReplenish(t *testing.T) {
	budget := retry.NewRetryBudget(2, 100*time.Millisecond)
	defer budget.Reset()

	// Use up budget
	budget.TryConsume()
	budget.TryConsume()

	if budget.TryConsume() {
		t.Error("3rd attempt should be rejected")
	}

	// Wait for window to pass
	time.Sleep(150 * time.Millisecond)

	// Should be replenished
	if !budget.TryConsume() {
		t.Error("budget should be replenished after window")
	}
}

func TestBudgetedClient(t *testing.T) {
	ctx := context.Background()

	// Provider fails 100 times (always down)
	provider := retry.NewMockProvider().
		WithMaxFailures(100).
		WithStatusCode(retry.StatusInternalServerError)

	baseClient := retry.NewRetryableClient(provider,
		retry.WithRetryAttempts(10),
		retry.WithBaseDelay(5*time.Millisecond),
	)

	// Retry budget: only allows 3 retries per second
	budget := retry.NewRetryBudget(3, time.Second)
	client := retry.NewBudgetedClient(baseClient, budget)

	_, err := client.Get(ctx, "http://example.com")
	if err == nil {
		t.Fatal("expected error after budget exhaustion")
	}

	t.Logf("Budgeted client correctly rejected retries after budget exhausted: %v", err)
}

func TestAmplificationFactor(t *testing.T) {
	tests := []struct {
		layers         int
		retriesPerLayer int
		expected       int
	}{
		{1, 3, 4},
		{2, 3, 16},
		{3, 3, 64},
		{3, 4, 125},
	}

	for _, tt := range tests {
		result := retry.AmplificationFactor(tt.layers, tt.retriesPerLayer)
		if result != tt.expected {
			t.Errorf("AmplificationFactor(%d, %d) = %d, want %d",
				tt.layers, tt.retriesPerLayer, result, tt.expected)
		}
	}

	// Document the amplification
	t.Logf("3 layers, 3 retries each = %dx amplification!", retry.AmplificationFactor(3, 3))
}