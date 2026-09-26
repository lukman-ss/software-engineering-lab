package circuitbreaker

import (
	"errors"
	"sync"
	"testing"
	"time"
)

func TestInitialStateIsClosed(t *testing.T) {
	b := New(DefaultConfig())
	if b.State() != Closed {
		t.Fatalf("expected CLOSED, got %s", b.State())
	}
}

func TestSuccessfulCallsStayClosed(t *testing.T) {
	b := New(Config{FailureThreshold: 3, OpenTimeout: 50 * time.Millisecond, HalfOpenMaxCalls: 1})
	for i := 0; i < 10; i++ {
		if err := b.Execute(func() error { return nil }); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	if b.State() != Closed {
		t.Fatalf("expected CLOSED, got %s", b.State())
	}
}

func TestFailuresBelowThresholdStayClosed(t *testing.T) {
	b := New(Config{FailureThreshold: 3, OpenTimeout: 50 * time.Millisecond, HalfOpenMaxCalls: 1})
	for i := 0; i < 2; i++ {
		if err := b.Execute(func() error { return errors.New("fail") }); err == nil {
			t.Fatal("expected error")
		}
	}
	if b.State() != Closed {
		t.Fatalf("expected CLOSED, got %s", b.State())
	}
}

func TestThresholdReachedOpens(t *testing.T) {
	b := New(Config{FailureThreshold: 3, OpenTimeout: 50 * time.Millisecond, HalfOpenMaxCalls: 1})
	for i := 0; i < 3; i++ {
		if err := b.Execute(func() error { return errors.New("fail") }); err == nil {
			t.Fatal("expected error")
		}
	}
	if b.State() != Open {
		t.Fatalf("expected OPEN, got %s", b.State())
	}
}

func TestOpenFailsFast(t *testing.T) {
	b := New(Config{FailureThreshold: 1, OpenTimeout: 50 * time.Millisecond, HalfOpenMaxCalls: 1})
	if err := b.Execute(func() error { return errors.New("fail") }); err == nil {
		t.Fatal("expected error")
	}
	start := time.Now()
	err := b.Execute(func() error { return nil })
	dur := time.Since(start)
	if err != ErrCircuitOpen {
		t.Fatalf("expected ErrCircuitOpen, got %v", err)
	}
	if dur > 10*time.Millisecond {
		t.Fatalf("fail-fast took too long: %v", dur)
	}
}

func TestOpenDoesNotCallDownstream(t *testing.T) {
	b := New(Config{FailureThreshold: 1, OpenTimeout: 50 * time.Millisecond, HalfOpenMaxCalls: 1})
	calls := 0
	fn := func() error {
		calls++
		return errors.New("fail")
	}
	if err := b.Execute(fn); err == nil {
		t.Fatal("expected error")
	}
	for i := 0; i < 5; i++ {
		if err := b.Execute(fn); err != ErrCircuitOpen {
			t.Fatalf("expected ErrCircuitOpen, got %v", err)
		}
	}
	if calls != 1 {
		t.Fatalf("expected 1 downstream call, got %d", calls)
	}
}

func TestCooldownMovesToHalfOpenBehavior(t *testing.T) {
	b := New(Config{FailureThreshold: 1, OpenTimeout: 30 * time.Millisecond, HalfOpenMaxCalls: 1})
	if err := b.Execute(func() error { return errors.New("fail") }); err == nil {
		t.Fatal("expected error")
	}
	if b.State() != Open {
		t.Fatalf("expected OPEN, got %s", b.State())
	}
	time.Sleep(40 * time.Millisecond)
	if b.State() != HalfOpen {
		t.Fatalf("expected HALF_OPEN after cooldown, got %s", b.State())
	}
}

func TestSuccessfulHalfOpenProbeCloses(t *testing.T) {
	b := New(Config{FailureThreshold: 1, OpenTimeout: 30 * time.Millisecond, HalfOpenMaxCalls: 1})
	if err := b.Execute(func() error { return errors.New("fail") }); err == nil {
		t.Fatal("expected error")
	}
	time.Sleep(40 * time.Millisecond)
	if err := b.Execute(func() error { return nil }); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.State() != Closed {
		t.Fatalf("expected CLOSED after successful probe, got %s", b.State())
	}
}

func TestFailedHalfOpenProbeReopens(t *testing.T) {
	b := New(Config{FailureThreshold: 1, OpenTimeout: 30 * time.Millisecond, HalfOpenMaxCalls: 1})
	if err := b.Execute(func() error { return errors.New("fail") }); err == nil {
		t.Fatal("expected error")
	}
	time.Sleep(40 * time.Millisecond)
	if err := b.Execute(func() error { return errors.New("fail") }); err == nil {
		t.Fatal("expected error")
	}
	if b.State() != Open {
		t.Fatalf("expected OPEN after failed probe, got %s", b.State())
	}
}

func TestRecoveryAfterDependencyHealthy(t *testing.T) {
	b := New(Config{FailureThreshold: 2, OpenTimeout: 30 * time.Millisecond, HalfOpenMaxCalls: 1})
	down := true
	fn := func() error {
		if down {
			return errors.New("down")
		}
		return nil
	}
	for i := 0; i < 2; i++ {
		if err := b.Execute(fn); err == nil {
			t.Fatal("expected error")
		}
	}
	if b.State() != Open {
		t.Fatalf("expected OPEN, got %s", b.State())
	}
	time.Sleep(40 * time.Millisecond)
	down = false
	if err := b.Execute(fn); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.State() != Closed {
		t.Fatalf("expected CLOSED after recovery, got %s", b.State())
	}
	for i := 0; i < 3; i++ {
		if err := b.Execute(fn); err != nil {
			t.Fatalf("unexpected error after recovery: %v", err)
		}
	}
}

func TestConcurrentAccess(t *testing.T) {
	b := New(Config{FailureThreshold: 5, OpenTimeout: 30 * time.Millisecond, HalfOpenMaxCalls: 2})
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_ = b.Execute(func() error {
				if i%3 == 0 {
					return errors.New("fail")
				}
				return nil
			})
			_ = b.State()
		}(i)
	}
	wg.Wait()
	if b.State() != Open && b.State() != Closed && b.State() != HalfOpen {
		t.Fatalf("invalid state: %s", b.State())
	}
}

func TestSuccessInClosedResetsFailures(t *testing.T) {
	b := New(Config{FailureThreshold: 3, OpenTimeout: 50 * time.Millisecond, HalfOpenMaxCalls: 1})
	for i := 0; i < 2; i++ {
		_ = b.Execute(func() error { return errors.New("fail") })
	}
	if b.State() != Closed {
		t.Fatalf("expected CLOSED, got %s", b.State())
	}
	if err := b.Execute(func() error { return nil }); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 2; i++ {
		_ = b.Execute(func() error { return errors.New("fail") })
	}
	if b.State() != Closed {
		t.Fatalf("expected CLOSED, got %s", b.State())
	}
	_ = b.Execute(func() error { return errors.New("fail") })
	if b.State() != Open {
		t.Fatalf("expected OPEN, got %s", b.State())
	}
}

func TestHalfOpenThrottlesExcessCalls(t *testing.T) {
	b := New(Config{FailureThreshold: 1, OpenTimeout: 30 * time.Millisecond, HalfOpenMaxCalls: 1})
	_ = b.Execute(func() error { return errors.New("fail") })
	time.Sleep(40 * time.Millisecond)

	probeStarted := make(chan struct{})
	probeRelease := make(chan struct{})

	go func() {
		_ = b.Execute(func() error {
			close(probeStarted)
			<-probeRelease
			return nil
		})
	}()

	<-probeStarted
	err := b.Execute(func() error { return nil })
	if err != ErrCircuitOpen {
		t.Fatalf("expected ErrCircuitOpen when HalfOpenMaxCalls exceeded, got %v", err)
	}
	close(probeRelease)
}

func TestPanicInHalfOpenCleansUpState(t *testing.T) {
	b := New(Config{FailureThreshold: 1, OpenTimeout: 30 * time.Millisecond, HalfOpenMaxCalls: 1})
	_ = b.Execute(func() error { return errors.New("fail") })
	time.Sleep(40 * time.Millisecond)

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic to propagate")
		}

		// Verify that breaker is not stuck in HalfOpen with exhausted calls
		if b.State() != Open {
			t.Fatalf("expected breaker to return to OPEN on panic, got %s", b.State())
		}

		time.Sleep(40 * time.Millisecond)
		// Should now be HalfOpen again and allow a new probe call
		err := b.Execute(func() error { return nil })
		if err != nil {
			t.Fatalf("expected probe to succeed after recovery, got %v", err)
		}
		if b.State() != Closed {
			t.Fatalf("expected breaker CLOSED, got %s", b.State())
		}
	}()

	_ = b.Execute(func() error {
		panic("boom downstream")
	})
}

func TestTrailingInFlightRequestDoesNotCorruptNewState(t *testing.T) {
	b := New(Config{FailureThreshold: 1, OpenTimeout: 50 * time.Millisecond, HalfOpenMaxCalls: 1})

	started := make(chan struct{})
	finish := make(chan struct{})

	// Goroutine 1: Initiated during Closed (generation 0), but takes a long time
	go func() {
		_ = b.Execute(func() error {
			close(started)
			<-finish
			return errors.New("slow failure from past")
		})
	}()

	<-started

	// Trigger trip to OPEN via another request
	_ = b.Execute(func() error { return errors.New("immediate fail") })
	if b.State() != Open {
		t.Fatalf("expected OPEN state, got %s", b.State())
	}

	// Advance time until HalfOpen
	time.Sleep(60 * time.Millisecond)
	if b.State() != HalfOpen {
		t.Fatalf("expected HALF_OPEN state, got %s", b.State())
	}

	// Release the trailing slow request from generation 0
	close(finish)
	time.Sleep(10 * time.Millisecond)

	// Trailing request should NOT have flipped HalfOpen to Open
	if b.State() != HalfOpen {
		t.Fatalf("expected state to remain HALF_OPEN despite trailing failure, got %s", b.State())
	}

	// New probe succeeds, closing the breaker
	err := b.Execute(func() error { return nil })
	if err != nil {
		t.Fatalf("expected probe to succeed, got %v", err)
	}
	if b.State() != Closed {
		t.Fatalf("expected CLOSED state, got %s", b.State())
	}
}

func TestInterleavedConcurrentTransitions(t *testing.T) {
	b := New(Config{FailureThreshold: 2, OpenTimeout: 20 * time.Millisecond, HalfOpenMaxCalls: 2})
	var wg sync.WaitGroup

	for i := 0; i < 30; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			_ = b.Execute(func() error {
				if id%2 == 0 {
					time.Sleep(15 * time.Millisecond)
					return errors.New("err")
				}
				time.Sleep(5 * time.Millisecond)
				return nil
			})
		}(i)
	}

	wg.Wait()
}
