package retry

import (
	"context"
	"errors"
	"math/rand/v2"
	"time"
)

var ErrMaxRetriesExceeded = errors.New("max retries exceeded")

type Config struct {
	MaxAttempts int
	BaseBackoff time.Duration
	MaxBackoff  time.Duration
}

type Retrier struct {
	cfg Config
}

func NewRetrier(cfg Config) *Retrier {
	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = 3
	}
	if cfg.BaseBackoff <= 0 {
		cfg.BaseBackoff = 10 * time.Millisecond
	}
	if cfg.MaxBackoff <= 0 {
		cfg.MaxBackoff = 100 * time.Millisecond
	}
	return &Retrier{cfg: cfg}
}

func (r *Retrier) CalculateBackoff(attempt int) time.Duration {
	if attempt <= 0 {
		return 0
	}
	multiplier := 1 << uint(attempt-1)
	temp := float64(r.cfg.BaseBackoff) * float64(multiplier)
	maxVal := float64(r.cfg.MaxBackoff)
	if temp > maxVal {
		temp = maxVal
	}
	// Full Jitter: random duration in [0, temp]
	sleep := rand.Float64() * temp
	return time.Duration(sleep)
}

func (r *Retrier) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	var lastErr error
	for attempt := 1; attempt <= r.cfg.MaxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err := fn(ctx)
		if err == nil {
			return nil
		}
		lastErr = err

		if attempt == r.cfg.MaxAttempts {
			break
		}

		backoff := r.CalculateBackoff(attempt)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff):
		}
	}
	return errors.Join(ErrMaxRetriesExceeded, lastErr)
}
