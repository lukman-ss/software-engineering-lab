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

func TestRetrier_JitterBoundsAndZeroConfig(t *testing.T) {
	r := NewRetrier(Config{})
	if r.cfg.MaxAttempts <= 0 || r.cfg.BaseBackoff <= 0 || r.cfg.MaxBackoff <= 0 {
		t.Fatalf("expected positive default configs, got %+v", r.cfg)
	}

	rCustom := NewRetrier(Config{
		MaxAttempts: 3,
		BaseBackoff: 10 * time.Millisecond,
		MaxBackoff:  40 * time.Millisecond,
	})

	for attempt := 1; attempt <= 5; attempt++ {
		backoff := rCustom.CalculateBackoff(attempt)
		if backoff < 0 {
			t.Fatalf("backoff should never be negative, got %v", backoff)
		}
		if backoff > 40*time.Millisecond {
			t.Fatalf("backoff should not exceed MaxBackoff (40ms), got %v", backoff)
		}
	}
}
