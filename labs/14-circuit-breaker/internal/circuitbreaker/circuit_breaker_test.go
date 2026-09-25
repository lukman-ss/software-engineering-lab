package circuitbreaker

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestCircuitBreaker(t *testing.T) {
	errMock := errors.New("mock downstream error")

	t.Run("1. initial state is CLOSED", func(t *testing.T) {
		cb := New(Config{FailureThreshold: 2})
		if got := cb.State(); got != StateClosed {
			t.Fatalf("expected StateClosed, got %s", got)
		}
	})

	t.Run("2. successful calls stay CLOSED", func(t *testing.T) {
		cb := New(Config{FailureThreshold: 2})
		err := cb.Execute(func() error { return nil })
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got := cb.State(); got != StateClosed {
			t.Fatalf("expected StateClosed, got %s", got)
		}
	})

	t.Run("3. failures below threshold stay CLOSED", func(t *testing.T) {
		cb := New(Config{FailureThreshold: 3})
		_ = cb.Execute(func() error { return errMock })
		_ = cb.Execute(func() error { return errMock })
		if got := cb.State(); got != StateClosed {
			t.Fatalf("expected StateClosed, got %s", got)
		}
	})

	t.Run("4. threshold reached changes state to OPEN", func(t *testing.T) {
		cb := New(Config{FailureThreshold: 2})
		_ = cb.Execute(func() error { return errMock })
		_ = cb.Execute(func() error { return errMock })
		if got := cb.State(); got != StateOpen {
			t.Fatalf("expected StateOpen, got %s", got)
		}
	})

	t.Run("5. OPEN calls fail fast", func(t *testing.T) {
		cb := New(Config{FailureThreshold: 1, OpenTimeout: time.Hour})
		_ = cb.Execute(func() error { return errMock })
		err := cb.Execute(func() error { return nil })
		if !errors.Is(err, ErrCircuitOpen) {
			t.Fatalf("expected ErrCircuitOpen, got %v", err)
		}
	})

	t.Run("6. OPEN calls do not execute downstream function", func(t *testing.T) {
		cb := New(Config{FailureThreshold: 1, OpenTimeout: time.Hour})
		_ = cb.Execute(func() error { return errMock })

		called := false
		_ = cb.Execute(func() error {
			called = true
			return nil
		})
		if called {
			t.Fatal("downstream function was unexpectedly called while circuit is OPEN")
		}
	})

	t.Run("7. cooldown moves breaker toward HALF_OPEN behavior", func(t *testing.T) {
		currentTime := time.Now()
		cb := New(Config{FailureThreshold: 1, OpenTimeout: 50 * time.Millisecond})
		cb.now = func() time.Time { return currentTime }

		_ = cb.Execute(func() error { return errMock })
		if got := cb.State(); got != StateOpen {
			t.Fatalf("expected StateOpen, got %s", got)
		}

		currentTime = currentTime.Add(60 * time.Millisecond)
		if got := cb.State(); got != StateHalfOpen {
			t.Fatalf("expected StateHalfOpen, got %s", got)
		}
	})

	t.Run("8. successful HALF_OPEN probe closes circuit", func(t *testing.T) {
		currentTime := time.Now()
		cb := New(Config{FailureThreshold: 1, OpenTimeout: 50 * time.Millisecond, HalfOpenMaxCalls: 1})
		cb.now = func() time.Time { return currentTime }

		_ = cb.Execute(func() error { return errMock })
		currentTime = currentTime.Add(60 * time.Millisecond)

		err := cb.Execute(func() error { return nil })
		if err != nil {
			t.Fatalf("unexpected probe error: %v", err)
		}
		if got := cb.State(); got != StateClosed {
			t.Fatalf("expected StateClosed after successful probe, got %s", got)
		}
	})

	t.Run("9. failed HALF_OPEN probe opens circuit again", func(t *testing.T) {
		currentTime := time.Now()
		cb := New(Config{FailureThreshold: 1, OpenTimeout: 50 * time.Millisecond, HalfOpenMaxCalls: 1})
		cb.now = func() time.Time { return currentTime }

		_ = cb.Execute(func() error { return errMock })
		currentTime = currentTime.Add(60 * time.Millisecond)

		err := cb.Execute(func() error { return errMock })
		if !errors.Is(err, errMock) {
			t.Fatalf("expected errMock, got %v", err)
		}
		if got := cb.State(); got != StateOpen {
			t.Fatalf("expected StateOpen after failed probe, got %s", got)
		}
	})

	t.Run("10. circuit recovers after dependency becomes healthy", func(t *testing.T) {
		currentTime := time.Now()
		cb := New(Config{FailureThreshold: 2, OpenTimeout: 50 * time.Millisecond, HalfOpenMaxCalls: 1})
		cb.now = func() time.Time { return currentTime }

		// Fail until open
		_ = cb.Execute(func() error { return errMock })
		_ = cb.Execute(func() error { return errMock })
		if got := cb.State(); got != StateOpen {
			t.Fatalf("expected StateOpen, got %s", got)
		}

		// Cooldown elapsed
		currentTime = currentTime.Add(60 * time.Millisecond)

		// Probe success
		_ = cb.Execute(func() error { return nil })
		if got := cb.State(); got != StateClosed {
			t.Fatalf("expected StateClosed, got %s", got)
		}

		// Subsequent calls should pass normally
		err := cb.Execute(func() error { return nil })
		if err != nil {
			t.Fatalf("unexpected error on subsequent call: %v", err)
		}
	})

	t.Run("11. concurrency and race safety", func(t *testing.T) {
		cb := New(Config{FailureThreshold: 5, OpenTimeout: 10 * time.Millisecond})
		var wg sync.WaitGroup
		var calls atomic.Int64

		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				_ = cb.Execute(func() error {
					calls.Add(1)
					if idx%2 == 0 {
						return errMock
					}
					return nil
				})
			}(i)
		}
		wg.Wait()
	})
}
