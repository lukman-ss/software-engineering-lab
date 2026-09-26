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
