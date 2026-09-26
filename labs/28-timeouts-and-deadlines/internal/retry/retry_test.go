package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRetrier_SuccessOnFirstTry(t *testing.T) {
	r := NewRetrier(Config{MaxAttempts: 3, BaseBackoff: 5 * time.Millisecond, MaxBackoff: 50 * time.Millisecond})
	attempts := 0
	err := r.Do(context.Background(), func(ctx context.Context) error {
		attempts++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if attempts != 1 {
		t.Fatalf("expected 1 attempt, got %d", attempts)
	}
}

func TestRetrier_RetryUntilSuccess(t *testing.T) {
	r := NewRetrier(Config{MaxAttempts: 3, BaseBackoff: 2 * time.Millisecond, MaxBackoff: 20 * time.Millisecond})
	attempts := 0
	err := r.Do(context.Background(), func(ctx context.Context) error {
		attempts++
		if attempts < 3 {
			return errors.New("transient error")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempts)
	}
}

func TestRetrier_ExceedMaxAttempts(t *testing.T) {
	r := NewRetrier(Config{MaxAttempts: 2, BaseBackoff: 2 * time.Millisecond, MaxBackoff: 10 * time.Millisecond})
	attempts := 0
	err := r.Do(context.Background(), func(ctx context.Context) error {
		attempts++
		return errors.New("persistent error")
	})
	if !errors.Is(err, ErrMaxRetriesExceeded) {
		t.Fatalf("expected ErrMaxRetriesExceeded, got %v", err)
	}
	if attempts != 2 {
		t.Fatalf("expected 2 attempts, got %d", attempts)
	}
}

func TestRetrier_ContextCanceled(t *testing.T) {
	r := NewRetrier(Config{MaxAttempts: 5, BaseBackoff: 50 * time.Millisecond, MaxBackoff: 200 * time.Millisecond})
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	err := r.Do(ctx, func(c context.Context) error {
		return errors.New("fail")
	})
	if !errors.Is(err, context.DeadlineExceeded) && !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context deadline/canceled error, got %v", err)
	}
}
