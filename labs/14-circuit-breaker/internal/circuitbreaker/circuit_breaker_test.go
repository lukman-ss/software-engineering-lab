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
		var mockErrors atomic.Int32
		var openErrors atomic.Int32

		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := cb.Execute(func() error {
					return errMock
				})
				if errors.Is(err, ErrCircuitOpen) {
					openErrors.Add(1)
				} else if errors.Is(err, errMock) {
					mockErrors.Add(1)
				} else {
					t.Errorf("unexpected error: %v", err)
				}
			}()
		}
		wg.Wait()

		if mockErrors.Load() < 5 {
			t.Fatalf("expected at least 5 mock errors, got %d", mockErrors.Load())
		}
		if openErrors.Load() <= 0 {
			t.Fatalf("expected circuit open errors, got %d", openErrors.Load())
		}
		if total := mockErrors.Load() + openErrors.Load(); total != 50 {
			t.Fatalf("expected 50 total errors, got %d", total)
		}
		if got := cb.State(); got != StateOpen {
			t.Fatalf("expected StateOpen after concurrent failures, got %s", got)
		}
	})

	t.Run("12. HalfOpenMaxCalls > 1 limits probes and requires consecutive successes", func(t *testing.T) {
		currentTime := time.Now()
		cb := New(Config{FailureThreshold: 1, OpenTimeout: 50 * time.Millisecond, HalfOpenMaxCalls: 2})
		cb.now = func() time.Time { return currentTime }

		// Trip to OPEN
		_ = cb.Execute(func() error { return errMock })
		currentTime = currentTime.Add(60 * time.Millisecond)

		// First probe succeeds, state stays HALF_OPEN
		if err := cb.Execute(func() error { return nil }); err != nil {
			t.Fatalf("unexpected probe error: %v", err)
		}
		if got := cb.State(); got != StateHalfOpen {
			t.Fatalf("expected StateHalfOpen after 1 of 2 probes, got %s", got)
		}

		// Second probe succeeds, moves to CLOSED
		if err := cb.Execute(func() error { return nil }); err != nil {
			t.Fatalf("unexpected probe error: %v", err)
		}
		if got := cb.State(); got != StateClosed {
			t.Fatalf("expected StateClosed after 2 of 2 probes, got %s", got)
		}

		// Now test probe limit enforcement
		_ = cb.Execute(func() error { return errMock })
		currentTime = currentTime.Add(60 * time.Millisecond)

		// In HALF_OPEN, simulate slow probe that consumes allowed probes
		var blockProbe sync.WaitGroup
		var probeStarted sync.WaitGroup
		blockProbe.Add(1)
		probeStarted.Add(2)

		for i := 0; i < 2; i++ {
			go func() {
				_ = cb.Execute(func() error {
					probeStarted.Done()
					blockProbe.Wait()
					return nil
				})
			}()
		}

		probeStarted.Wait()
		// Third concurrent call should be rejected with ErrCircuitOpen because halfOpenCalls == 2
		err := cb.Execute(func() error { return nil })
		if !errors.Is(err, ErrCircuitOpen) {
			t.Fatalf("expected ErrCircuitOpen for probe exceeding HalfOpenMaxCalls, got %v", err)
		}
		blockProbe.Done()
	})

	t.Run("13. panic in fn during HALF_OPEN sets state back to OPEN", func(t *testing.T) {
		currentTime := time.Now()
		cb := New(Config{FailureThreshold: 1, OpenTimeout: 50 * time.Millisecond, HalfOpenMaxCalls: 1})
		cb.now = func() time.Time { return currentTime }

		_ = cb.Execute(func() error { return errMock })
		currentTime = currentTime.Add(60 * time.Millisecond)

		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic to propagate")
			}
			if got := cb.State(); got != StateOpen {
				t.Fatalf("expected StateOpen after panic, got %s", got)
			}
		}()

		_ = cb.Execute(func() error {
			panic("downstream panic")
		})
	})

	t.Run("14. default config values applied", func(t *testing.T) {
		cb := New(Config{})
		if cb.config.FailureThreshold != 3 {
			t.Fatalf("expected FailureThreshold=3, got %d", cb.config.FailureThreshold)
		}
		if cb.config.OpenTimeout != 5*time.Second {
			t.Fatalf("expected OpenTimeout=5s, got %v", cb.config.OpenTimeout)
		}
		if cb.config.HalfOpenMaxCalls != 1 {
			t.Fatalf("expected HalfOpenMaxCalls=1, got %d", cb.config.HalfOpenMaxCalls)
		}
	})

	t.Run("15. panic in fn during CLOSED records failure", func(t *testing.T) {
		cb := New(Config{FailureThreshold: 2})

		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic to propagate")
			}
			if cb.failureCount != 1 {
				t.Fatalf("expected failureCount=1 after panic, got %d", cb.failureCount)
			}
			if got := cb.State(); got != StateClosed {
				t.Fatalf("expected StateClosed, got %s", got)
			}
		}()

		_ = cb.Execute(func() error {
			panic("downstream panic")
		})
	})

	t.Run("16. success in CLOSED resets consecutive failure count", func(t *testing.T) {
		cb := New(Config{FailureThreshold: 3})

		// 2 failures (threshold - 1), state remains CLOSED
		_ = cb.Execute(func() error { return errMock })
		_ = cb.Execute(func() error { return errMock })
		if got := cb.State(); got != StateClosed {
			t.Fatalf("expected StateClosed after 2 failures, got %s", got)
		}

		// 1 success resets consecutive failure count
		if err := cb.Execute(func() error { return nil }); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// 2 more failures should NOT trip breaker since count was reset to 0
		_ = cb.Execute(func() error { return errMock })
		_ = cb.Execute(func() error { return errMock })
		if got := cb.State(); got != StateClosed {
			t.Fatalf("expected StateClosed after reset and 2 failures, got %s", got)
		}

		// 3rd failure reaches threshold and trips to OPEN
		_ = cb.Execute(func() error { return errMock })
		if got := cb.State(); got != StateOpen {
			t.Fatalf("expected StateOpen after reaching threshold, got %s", got)
		}
	})
}
