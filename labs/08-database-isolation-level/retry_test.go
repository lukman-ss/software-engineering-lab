package isolation_test

import (
	"context"
	"errors"
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

func TestWrapTxError_SerializationFailure(t *testing.T) {
	original := &pq.Error{Code: "40001", Message: "could not serialize access"}
	err := isolation.WrapTxError("deduct", original)

	if !errors.Is(err, isolation.ErrSerializationFailure) {
		t.Fatalf("expected ErrSerializationFailure, got %v", err)
	}

	var pqErr *pq.Error
	if !errors.As(err, &pqErr) {
		t.Fatalf("expected *pq.Error via errors.As, got %T", err)
	}
	if pqErr.Code != "40001" {
		t.Fatalf("expected SQLSTATE 40001, got %s", pqErr.Code)
	}
}

func TestWrapTxError_Deadlock(t *testing.T) {
	original := &pq.Error{Code: "40P01", Message: "deadlock detected"}
	err := isolation.WrapTxError("commit", original)

	if !errors.Is(err, isolation.ErrDeadlockDetected) {
		t.Fatalf("expected ErrDeadlockDetected, got %v", err)
	}

	var pqErr *pq.Error
	if !errors.As(err, &pqErr) {
		t.Fatalf("expected *pq.Error via errors.As, got %T", err)
	}
	if pqErr.Code != "40P01" {
		t.Fatalf("expected SQLSTATE 40P01, got %s", pqErr.Code)
	}
}

func TestRetryTransaction_SerializationFailureThenSuccess(t *testing.T) {
	ctx := context.Background()
	attempts := 0

	policy := isolation.RetryPolicy{
		BaseDelay: 1 * time.Millisecond,
		MaxDelay:  100 * time.Millisecond,
		RandSrc:   func() float64 { return 0.5 },
		Sleep:     func(ctx context.Context, d time.Duration) error { return nil },
	}

	err := isolation.RetryTransaction(ctx, 3, policy, func(ctx context.Context) error {
		attempts++
		if attempts < 2 {
			return &pq.Error{Code: "40001", Message: "serialization failure"}
		}
		return nil
	})

	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if attempts != 2 {
		t.Fatalf("expected 2 attempts, got %d", attempts)
	}
}

func TestRetryTransaction_DeadlockThenSuccess(t *testing.T) {
	ctx := context.Background()
	attempts := 0

	policy := isolation.RetryPolicy{
		BaseDelay: 1 * time.Millisecond,
		MaxDelay:  100 * time.Millisecond,
		RandSrc:   func() float64 { return 0.5 },
		Sleep:     func(ctx context.Context, d time.Duration) error { return nil },
	}

	err := isolation.RetryTransaction(ctx, 3, policy, func(ctx context.Context) error {
		attempts++
		if attempts < 3 {
			return &pq.Error{Code: "40P01", Message: "deadlock detected"}
		}
		return nil
	})

	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempts)
	}
}

func TestRetryTransaction_NonRetryableStopsImmediately(t *testing.T) {
	ctx := context.Background()
	attempts := 0

	policy := isolation.RetryPolicy{
		BaseDelay: 1 * time.Millisecond,
		MaxDelay:  100 * time.Millisecond,
		RandSrc:   func() float64 { return 0.5 },
		Sleep:     func(ctx context.Context, d time.Duration) error { return nil },
	}

	err := isolation.RetryTransaction(ctx, 5, policy, func(ctx context.Context) error {
		attempts++
		return &pq.Error{Code: "23505", Message: "duplicate key value violates unique constraint"}
	})

	var pqErr *pq.Error
	if !errors.As(err, &pqErr) {
		t.Fatalf("expected pq.Error via errors.As, got %T", err)
	}
	if pqErr.Code != "23505" {
		t.Fatalf("expected SQLSTATE 23505, got %s", pqErr.Code)
	}
	if attempts != 1 {
		t.Fatalf("expected 1 attempt (stopped immediately), got %d", attempts)
	}
}

func TestRetryTransaction_MaxAttemptsExceeded(t *testing.T) {
	ctx := context.Background()
	attempts := 0

	policy := isolation.RetryPolicy{
		BaseDelay: 1 * time.Millisecond,
		MaxDelay:  100 * time.Millisecond,
		RandSrc:   func() float64 { return 0.5 },
		Sleep:     func(ctx context.Context, d time.Duration) error { return nil },
	}

	err := isolation.RetryTransaction(ctx, 3, policy, func(ctx context.Context) error {
		attempts++
		return &pq.Error{Code: "40001", Message: "serialization failure"}
	})

	if !errors.Is(err, isolation.ErrMaxRetryExceeded) {
		t.Fatalf("expected ErrMaxRetryExceeded, got %v", err)
	}

	var pqErr *pq.Error
	if !errors.As(err, &pqErr) {
		t.Fatalf("expected pq.Error via errors.As, got %v", err)
	}
	if pqErr.Code != "40001" {
		t.Fatalf("expected underlying SQLSTATE 40001, got %s", pqErr.Code)
	}
	if attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempts)
	}
}

func TestRetryTransaction_ContextCancelledBeforeFirstAttempt(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	attempts := 0

	policy := isolation.RetryPolicy{
		BaseDelay: 1 * time.Millisecond,
		MaxDelay:  100 * time.Millisecond,
		RandSrc:   func() float64 { return 0.5 },
		Sleep:     func(ctx context.Context, d time.Duration) error { return nil },
	}

	err := isolation.RetryTransaction(ctx, 5, policy, func(ctx context.Context) error {
		attempts++
		return nil
	})

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if attempts != 0 {
		t.Fatalf("expected 0 attempts, got %d", attempts)
	}
}

func TestRetryTransaction_ContextCancelledDuringBackoff(t *testing.T) {
	attempts := 0
	sleepCancelled := make(chan struct{})

	policy := isolation.RetryPolicy{
		BaseDelay: 1 * time.Second, // Long delay so context times out during backoff
		MaxDelay:  2 * time.Second,
		RandSrc:   func() float64 { return 0.5 },
		Sleep: func(ctx context.Context, d time.Duration) error {
			select {
			case <-ctx.Done():
				close(sleepCancelled)
				return ctx.Err()
			case <-time.After(d):
				return nil
			}
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := isolation.RetryTransaction(ctx, 5, policy, func(ctx context.Context) error {
		attempts++
		return &pq.Error{Code: "40001", Message: "serialization failure"}
	})

	if !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected context.Canceled or DeadlineExceeded, got %v", err)
	}
	if attempts != 1 {
		t.Fatalf("expected 1 attempt, got %d", attempts)
	}
	select {
	case <-sleepCancelled:
		t.Logf("sleep was cancelled by context")
	default:
		t.Logf("sleep may not have been cancelled, attempts=%d", attempts)
	}
}

func TestRetryTransaction_SleepErrorPropagated(t *testing.T) {
	ctx := context.Background()
	expectedErr := errors.New("custom sleep error")

	policy := isolation.RetryPolicy{
		BaseDelay: 1 * time.Millisecond,
		MaxDelay:  10 * time.Millisecond,
		RandSrc:   func() float64 { return 0.5 },
		Sleep: func(ctx context.Context, d time.Duration) error {
			return expectedErr
		},
	}

	err := isolation.RetryTransaction(ctx, 3, policy, func(ctx context.Context) error {
		return &pq.Error{Code: "40001", Message: "serialization failure"}
	})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected sleep error %v, got %v", expectedErr, err)
	}
}

func TestRetryTransaction_SuccessFirstAttempt(t *testing.T) {
	ctx := context.Background()
	attempts := 0

	policy := isolation.RetryPolicy{
		BaseDelay: 1 * time.Millisecond,
		MaxDelay:  100 * time.Millisecond,
		RandSrc:   func() float64 { return 0.5 },
		Sleep:     func(ctx context.Context, d time.Duration) error { return nil },
	}

	err := isolation.RetryTransaction(ctx, 3, policy, func(ctx context.Context) error {
		attempts++
		return nil
	})

	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if attempts != 1 {
		t.Fatalf("expected 1 attempt, got %d", attempts)
	}
}

func TestRetryTransaction_InvalidMaxAttempts(t *testing.T) {
	ctx := context.Background()

	policy := isolation.RetryPolicy{
		BaseDelay: 1 * time.Millisecond,
		MaxDelay:  100 * time.Millisecond,
		RandSrc:   func() float64 { return 0.5 },
		Sleep:     func(ctx context.Context, d time.Duration) error { return nil },
	}

	err := isolation.RetryTransaction(ctx, 0, policy, func(ctx context.Context) error {
		return nil
	})
	if !errors.Is(err, isolation.ErrInvalidMaxAttempts) {
		t.Fatalf("expected ErrInvalidMaxAttempts for maxAttempts=0, got %v", err)
	}

	err = isolation.RetryTransaction(ctx, -1, policy, func(ctx context.Context) error {
		return nil
	})
	if !errors.Is(err, isolation.ErrInvalidMaxAttempts) {
		t.Fatalf("expected ErrInvalidMaxAttempts for maxAttempts=-1, got %v", err)
	}
}

func TestRetryTransaction_DelayBounds(t *testing.T) {
	ctx := context.Background()
	attempts := 0
	var delays []time.Duration

	policy := isolation.RetryPolicy{
		BaseDelay: 10 * time.Millisecond,
		MaxDelay:  1 * time.Second,
		RandSrc: func() float64 {
			v := 0.0
			return v
		},
		Sleep: func(ctx context.Context, d time.Duration) error {
			delays = append(delays, d)
			return nil
		},
	}

	err := isolation.RetryTransaction(ctx, 4, policy, func(ctx context.Context) error {
		attempts++
		return &pq.Error{Code: "40001", Message: "serialization failure"}
	})

	if !errors.Is(err, isolation.ErrMaxRetryExceeded) {
		t.Fatalf("expected ErrMaxRetryExceeded, got %v", err)
	}
	if attempts != 4 {
		t.Fatalf("expected 4 attempts, got %d", attempts)
	}
	if len(delays) != 3 {
		t.Fatalf("expected 3 backoff sleeps, got %d", len(delays))
	}

	for i, d := range delays {
		cap := time.Duration(float64(policy.BaseDelay) * float64(int(1) << uint(i)))
		if cap > policy.MaxDelay {
			cap = policy.MaxDelay
		}
		if d < 0 || d > cap {
			t.Fatalf("delay %d out of bounds [0, %v]: %v", i+1, cap, d)
		}
	}
}