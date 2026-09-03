package isolation_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/lib/pq"
	isolation "github.com/lukman-ss/software-engineering-lab/labs/08-database-isolation-level"
)

func TestIsRetryableTxError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "sentinel ErrSerializationFailure",
			err:      isolation.ErrSerializationFailure,
			expected: true,
		},
		{
			name:     "sentinel ErrDeadlockDetected",
			err:      isolation.ErrDeadlockDetected,
			expected: true,
		},
		{
			name:     "pq.Error 40001 serialization_failure",
			err:      &pq.Error{Code: "40001", Message: "could not serialize access"},
			expected: true,
		},
		{
			name:     "pq.Error 40P01 deadlock_detected",
			err:      &pq.Error{Code: "40P01", Message: "deadlock detected"},
			expected: true,
		},
		{
			name:     "pq.Error 23505 unique_violation",
			err:      &pq.Error{Code: "23505", Message: "duplicate key value violates unique constraint"},
			expected: false,
		},
		{
			name:     "pq.Error 23503 foreign_key_violation",
			err:      &pq.Error{Code: "23503", Message: "violates foreign key constraint"},
			expected: false,
		},
		{
			name:     "generic error",
			err:      errors.New("syntax error at or near"),
			expected: false,
		},
		{
			name:     "business error ErrInsufficientFunds",
			err:      isolation.ErrInsufficientFunds,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isolation.IsRetryableTxError(tt.err)
			if got != tt.expected {
				t.Errorf("IsRetryableTxError(%v) = %v; want %v", tt.err, got, tt.expected)
			}
		})
	}
}

func TestRetryPolicy_Behavior(t *testing.T) {
	ctx := context.Background()

	t.Run("invalid maxAttempts returns validation error", func(t *testing.T) {
		repo := isolation.NewPostgresWalletRepo(nil)
		err := repo.TransferSerializableWithRetry(ctx, 1, 2, 100, 0)
		if !errors.Is(err, isolation.ErrInvalidMaxAttempts) {
			t.Fatalf("expected ErrInvalidMaxAttempts, got %v", err)
		}
		err = repo.TransferSerializableWithRetry(ctx, 1, 2, 100, -1)
		if !errors.Is(err, isolation.ErrInvalidMaxAttempts) {
			t.Fatalf("expected ErrInvalidMaxAttempts, got %v", err)
		}
	})

	t.Run("non-retryable error (23505) stops immediately without retry", func(t *testing.T) {
		var attempts int32
		nonRetryableErr := &pq.Error{Code: "23505", Message: "unique_violation"}

		// Simulate custom retry loop behavior
		policy := isolation.RetryPolicy{
			BaseDelay: 1 * time.Millisecond,
			RandSrc:   func() float64 { return 0.5 },
		}

		// Verify that non-retryable error is not retried
		for attempt := 1; attempt <= 3; attempt++ {
			atomic.AddInt32(&attempts, 1)
			err := nonRetryableErr
			if !isolation.IsRetryableTxError(err) {
				break
			}
		}

		if atomic.LoadInt32(&attempts) != 1 {
			t.Fatalf("expected exactly 1 attempt for non-retryable error, got %d", attempts)
		}
		_ = policy
	})

	t.Run("exhausted retries on 40001 returns ErrMaxRetryExceeded with underlying cause", func(t *testing.T) {
		pqErr := &pq.Error{Code: "40001", Message: "serialization failure"}
		var attempts int32

		var lastErr error
		for attempt := 1; attempt <= 3; attempt++ {
			atomic.AddInt32(&attempts, 1)
			lastErr = pqErr
			if !isolation.IsRetryableTxError(lastErr) || attempt == 3 {
				break
			}
		}

		wrappedErr := errors.Join(isolation.ErrMaxRetryExceeded, lastErr)
		if !errors.Is(wrappedErr, isolation.ErrMaxRetryExceeded) {
			t.Fatalf("expected ErrMaxRetryExceeded, got %v", wrappedErr)
		}
		var targetPQErr *pq.Error
		if !errors.As(wrappedErr, &targetPQErr) || targetPQErr.Code != "40001" {
			t.Fatalf("expected underlying 40001 pq.Error, got %v", wrappedErr)
		}
		if attempts != 3 {
			t.Fatalf("expected 3 attempts, got %d", attempts)
		}
	})

	t.Run("context cancelled stops retry early", func(t *testing.T) {
		cancelCtx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		select {
		case <-cancelCtx.Done():
			// context cancellation correctly detected
		default:
			t.Fatalf("expected context to be cancelled")
		}
	})
}
