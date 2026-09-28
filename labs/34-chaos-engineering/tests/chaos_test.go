package tests

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"labs/34-chaos-engineering/internal/circuitbreaker"
	"labs/34-chaos-engineering/internal/experiment"
	"labs/34-chaos-engineering/internal/fault"
	"labs/34-chaos-engineering/internal/monitor"
)

func TestFaultInjector(t *testing.T) {
	fi := fault.NewInjector()
	ctx := context.Background()

	// Initially disabled
	if err := fi.Execute(ctx); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// Set fault
	fi.SetFault(10*time.Millisecond, true)
	if !fi.IsEnabled() {
		t.Fatalf("expected injector to be enabled")
	}

	start := time.Now()
	err := fi.Execute(ctx)
	elapsed := time.Since(start)

	if !errors.Is(err, fault.ErrInjectedFault) {
		t.Fatalf("expected ErrInjectedFault, got %v", err)
	}
	if elapsed < 10*time.Millisecond {
		t.Fatalf("expected latency injection >= 10ms, got %v", elapsed)
	}

	// Clear fault
	fi.Clear()
	if fi.IsEnabled() {
		t.Fatalf("expected injector to be disabled")
	}
	if err := fi.Execute(ctx); err != nil {
		t.Fatalf("expected nil error after clear, got %v", err)
	}
}

func TestCircuitBreakerStateTransitions(t *testing.T) {
	cb := circuitbreaker.NewCircuitBreaker(2, 50*time.Millisecond)

	failFn := func() error { return errors.New("remote error") }
	okFn := func() error { return nil }

	// First failure - still closed
	_ = cb.Execute(failFn, nil)
	if cb.State() != circuitbreaker.StateClosed {
		t.Fatalf("expected StateClosed, got %v", cb.State())
	}

	// Second failure - trips to open
	_ = cb.Execute(failFn, nil)
	if cb.State() != circuitbreaker.StateOpen {
		t.Fatalf("expected StateOpen, got %v", cb.State())
	}

	// Fast-fail when open
	err := cb.Execute(okFn, nil)
	if !errors.Is(err, circuitbreaker.ErrCircuitOpen) {
		t.Fatalf("expected ErrCircuitOpen, got %v", err)
	}

	// Wait for cooldown
	time.Sleep(60 * time.Millisecond)
	if cb.State() != circuitbreaker.StateHalfOpen {
		t.Fatalf("expected StateHalfOpen after cooldown, got %v", cb.State())
	}

	// Successful execution resets to closed
	err = cb.Execute(okFn, nil)
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if cb.State() != circuitbreaker.StateClosed {
		t.Fatalf("expected StateClosed after recovery, got %v", cb.State())
	}
}

func TestCircuitBreakerGracefulDegradation(t *testing.T) {
	cb := circuitbreaker.NewCircuitBreaker(1, 50*time.Millisecond)
	failFn := func() error { return errors.New("upstream failed") }

	fallbackInvoked := false
	fallback := func() error {
		fallbackInvoked = true
		return nil
	}

	err := cb.Execute(failFn, fallback)
	if err != nil {
		t.Fatalf("expected fallback to swallow error, got %v", err)
	}
	if !fallbackInvoked {
		t.Fatalf("expected fallback to be invoked")
	}
}

func TestExperimentAutoAbortOnSteadyStateViolation(t *testing.T) {
	inj := fault.NewInjector()
	mon := monitor.NewMonitor(0.20) // 20% max error rate

	// Seed with some baseline
	for i := 0; i < 4; i++ {
		mon.RecordSuccess()
	}

	cfg := experiment.Config{
		Name:            "abort-test",
		Duration:        1 * time.Second,
		ForceError:      true,
		MonitorInterval: 10 * time.Millisecond,
	}
	exp := experiment.NewExperiment(cfg, inj, mon)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- exp.Run(ctx)
	}()

	// Inject failures into monitor to breach threshold
	mon.RecordFailure()
	mon.RecordFailure()

	select {
	case err := <-errCh:
		if err == nil {
			t.Fatalf("expected error from auto-abort, got nil")
		}
		if exp.State() != experiment.StateAborted {
			t.Fatalf("expected StateAborted, got %v", exp.State())
		}
		if inj.IsEnabled() {
			t.Fatalf("expected injector to be neutralized after abort")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatalf("timeout waiting for experiment auto-abort")
	}
}

func TestConcurrencyAndRace(t *testing.T) {
	inj := fault.NewInjector()
	mon := monitor.NewMonitor(0.50)
	cb := circuitbreaker.NewCircuitBreaker(5, 20*time.Millisecond)

	var wg sync.WaitGroup
	workers := 20
	iterations := 50

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				if id%2 == 0 {
					inj.SetFault(time.Microsecond, true)
				} else {
					inj.Clear()
				}

				_ = cb.Execute(func() error {
					return inj.Execute(context.Background())
				}, func() error {
					return nil
				})

				mon.RecordSuccess()
				_ = mon.IsHealthy()
			}
		}(i)
	}

	wg.Wait()
}
