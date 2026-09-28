package circuit

import (
	"errors"
	"testing"
	"time"
)

func TestCircuitBreaker_StateTransitions(t *testing.T) {
	cb := NewBreaker(Config{
		FailureThreshold: 2,
		SuccessThreshold: 2,
		Cooldown:         50 * time.Millisecond,
	})

	if cb.State() != StateClosed {
		t.Fatalf("expected initial state CLOSED, got %v", cb.State())
	}

	dummyErr := errors.New("service fail")

	// Trigger failures to trip Open
	cb.Execute(func() error { return dummyErr })
	cb.Execute(func() error { return dummyErr })

	if cb.State() != StateOpen {
		t.Fatalf("expected state OPEN after 2 failures, got %v", cb.State())
	}

	err := cb.Execute(func() error { return nil })
	if !errors.Is(err, ErrCircuitOpen) {
		t.Fatalf("expected ErrCircuitOpen when OPEN, got %v", err)
	}

	// Wait for cooldown
	time.Sleep(60 * time.Millisecond)

	if cb.State() != StateHalfOpen {
		t.Fatalf("expected state HALF_OPEN after cooldown, got %v", cb.State())
	}

	// 2 successes in HALF_OPEN to transition back to CLOSED
	cb.Execute(func() error { return nil })
	cb.Execute(func() error { return nil })

	if cb.State() != StateClosed {
		t.Fatalf("expected state CLOSED after 2 successes in HALF_OPEN, got %v", cb.State())
	}
}

func TestCircuitBreaker_HalfOpenFailureTripsOpen(t *testing.T) {
	cb := NewBreaker(Config{
		FailureThreshold: 2,
		SuccessThreshold: 2,
		Cooldown:         30 * time.Millisecond,
	})

	dummyErr := errors.New("service fail")
	cb.Execute(func() error { return dummyErr })
	cb.Execute(func() error { return dummyErr })

	if cb.State() != StateOpen {
		t.Fatalf("expected OPEN state, got %v", cb.State())
	}

	time.Sleep(40 * time.Millisecond)

	if cb.State() != StateHalfOpen {
		t.Fatalf("expected HALF_OPEN state, got %v", cb.State())
	}

	// Single failure in HALF_OPEN must trip back to OPEN immediately
	err := cb.Execute(func() error { return dummyErr })
	if !errors.Is(err, dummyErr) {
		t.Fatalf("expected dummyErr, got %v", err)
	}

	if cb.State() != StateOpen {
		t.Fatalf("expected state OPEN immediately after failure in HALF_OPEN, got %v", cb.State())
	}
}

func TestCircuitBreaker_DefaultZeroConfig(t *testing.T) {
	cb := NewBreaker(Config{})
	if cb.State() != StateClosed {
		t.Fatalf("expected initial CLOSED, got %v", cb.State())
	}
	if cb.cfg.FailureThreshold <= 0 || cb.cfg.SuccessThreshold <= 0 || cb.cfg.Cooldown <= 0 {
		t.Fatalf("expected defaults for zero config, got %+v", cb.cfg)
	}
}
