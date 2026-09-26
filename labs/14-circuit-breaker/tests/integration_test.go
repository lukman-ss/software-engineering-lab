package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"circuitbreaker/internal/checkout"
	"circuitbreaker/internal/circuitbreaker"
	"circuitbreaker/internal/payment"
)

func TestCircuitBreakerIntegration(t *testing.T) {
	fakeServer := payment.NewFakeServer(250 * time.Millisecond)
	defer fakeServer.Close()

	clientTimeout := 50 * time.Millisecond
	paymentClient := payment.NewClient(fakeServer.URL(), clientTimeout)

	cb := circuitbreaker.New(circuitbreaker.Config{
		FailureThreshold: 2,
		OpenTimeout:      100 * time.Millisecond,
		HalfOpenMaxCalls: 1,
	})

	svc := checkout.NewService(paymentClient, cb)

	t.Run("downstream fails, CB trips open", func(t *testing.T) {
		fakeServer.SetMode(payment.ModeDown)

		// 1st failure
		_ = svc.Checkout(context.Background())
		// 2nd failure trips CB
		_ = svc.Checkout(context.Background())

		if cb.State() != circuitbreaker.StateOpen {
			t.Fatalf("expected OPEN state, got %s", cb.State())
		}

		// 3rd request fails fast
		err := svc.Checkout(context.Background())
		if !errors.Is(err, circuitbreaker.ErrCircuitOpen) {
			t.Fatalf("expected ErrCircuitOpen, got %v", err)
		}

		if got := fakeServer.RequestCount(); got != 2 {
			t.Fatalf("expected 2 downstream requests, got %d", got)
		}
	})

	t.Run("cooldown and recovery", func(t *testing.T) {
		time.Sleep(150 * time.Millisecond)
		fakeServer.SetMode(payment.ModeHealthy)

		if cb.State() != circuitbreaker.StateHalfOpen {
			t.Fatalf("expected HALF_OPEN state, got %s", cb.State())
		}

		err := svc.Checkout(context.Background())
		if err != nil {
			t.Fatalf("probe failed unexpectedly: %v", err)
		}

		if cb.State() != circuitbreaker.StateClosed {
			t.Fatalf("expected CLOSED state after probe, got %s", cb.State())
		}
	})
}
