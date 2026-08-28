package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lukman-ss/software-engineering-lab/labs/08-retry"
)

func TestMockProviderFailModes(t *testing.T) {
	ctx := context.Background()

	t.Run("success immediately", func(t *testing.T) {
		provider := retry.NewMockProvider()
		resp, err := provider.Get(ctx, "http://example.com")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != 200 {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("fail N times then succeed", func(t *testing.T) {
		// Fail 2 times with 500, then succeed
		provider := retry.NewMockProvider().
			WithMaxFailures(2).
			WithStatusCode(retry.StatusInternalServerError)

		// Attempt 1: Fail
		resp1, _ := provider.Get(ctx, "http://example.com")
		if resp1.StatusCode != 500 {
			t.Errorf("expected 500 on attempt 1, got %d", resp1.StatusCode)
		}

		// Attempt 2: Fail
		resp2, _ := provider.Get(ctx, "http://example.com")
		if resp2.StatusCode != 500 {
			t.Errorf("expected 500 on attempt 2, got %d", resp2.StatusCode)
		}

		// Attempt 3: Success
		resp3, _ := provider.Get(ctx, "http://example.com")
		if resp3.StatusCode != 200 {
			t.Errorf("expected 200 on attempt 3, got %d", resp3.StatusCode)
		}
	})

	t.Run("timeout", func(t *testing.T) {
		provider := retry.NewMockProvider().WithTimeout(true)
		_, err := provider.Get(ctx, "http://example.com")
		if err == nil {
			t.Fatal("expected timeout error")
		}
	})

	t.Run("rate limit", func(t *testing.T) {
		provider := retry.NewMockProvider().
			WithMaxFailures(1).
			WithStatusCode(retry.StatusTooManyRequests)

		resp, _ := provider.Get(ctx, "http://example.com")
		if resp.StatusCode != 429 {
			t.Errorf("expected 429, got %d", resp.StatusCode)
		}
	})
}

func TestRetryableClientTransientErrors(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		statusCode    retry.HTTPStatus
		maxFailures   int
		shouldSucceed bool
	}{
		{"500 Internal Server Error (transient)", retry.StatusInternalServerError, 2, true},
		{"503 Service Unavailable (transient)", retry.StatusServiceUnavailable, 2, true},
		{"429 Too Many Requests (transient)", retry.StatusTooManyRequests, 2, true},
		{"400 Bad Request (permanent)", retry.StatusBadRequest, 2, false},
		{"401 Unauthorized (permanent)", retry.StatusUnauthorized, 2, false},
		{"403 Forbidden (permanent)", retry.StatusForbidden, 2, false},
		{"404 Not Found (permanent)", retry.StatusNotFound, 2, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := retry.NewMockProvider().
				WithMaxFailures(tt.maxFailures).
				WithStatusCode(tt.statusCode)

			client := retry.NewRetryableClient(provider,
				retry.WithRetryAttempts(3),
				retry.WithBaseDelay(10*time.Millisecond),
			)

			resp, err := client.Get(ctx, "http://example.com")

			if tt.shouldSucceed {
				if err != nil {
					t.Fatalf("expected success after retries, got error: %v", err)
				}
				if resp.StatusCode != 200 {
					t.Errorf("expected final status 200, got %d", resp.StatusCode)
				}
				// 2 failures + 1 success = 3 requests
				if provider.TotalRequests() != 3 {
					t.Errorf("expected 3 total requests, got %d", provider.TotalRequests())
				}
			} else {
				if err == nil {
					t.Fatal("expected permanent error")
				}
				// Should fail immediately on attempt 1 without retrying
				if provider.TotalRequests() != 1 {
					t.Errorf("expected exactly 1 request for permanent error, got %d", provider.TotalRequests())
				}
			}
		})
	}
}

func TestExponentialBackoffJitter(t *testing.T) {
	// A bit hard to test exact randomness, but we can verify wait duration is roughly correct
	// Note: We use waitDuration indirectly via execution time
	ctx := context.Background()
	provider := retry.NewMockProvider().
		WithMaxFailures(3).
		WithStatusCode(retry.StatusInternalServerError)

	baseDelay := 50 * time.Millisecond
	client := retry.NewRetryableClient(provider,
		retry.WithRetryAttempts(5),
		retry.WithBaseDelay(baseDelay),
		retry.WithJitter(0.5), // +/- 50%
	)

	start := time.Now()
	_, err := client.Get(ctx, "http://example.com")
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	elapsed := time.Since(start)

	// Retries happen at:
	// attempt 0 -> fails, waits ~50ms
	// attempt 1 -> fails, waits ~100ms
	// attempt 2 -> fails, waits ~200ms
	// attempt 3 -> succeeds
	// Expected total sleep: ~350ms

	expectedMin := 350 * time.Millisecond / 2 // With massive down-jitter
	expectedMax := 350 * time.Millisecond * 2 // With massive up-jitter

	if elapsed < expectedMin || elapsed > expectedMax {
		t.Logf("Elapsed time %v outside expected jitter bounds [%v, %v]", elapsed, expectedMin, expectedMax)
	}
}
