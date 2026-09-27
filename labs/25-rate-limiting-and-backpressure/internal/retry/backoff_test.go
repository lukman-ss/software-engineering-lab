package retry

import (
	"testing"
	"time"
)

func TestComputeBackoff_Bounds(t *testing.T) {
	cfg := Config{
		Base: 100 * time.Millisecond,
		Cap:  2 * time.Second,
	}

	for attempt := 0; attempt < 10; attempt++ {
		// Full Jitter: 0 <= sleep <= min(cap, base * 2^attempt)
		sleepFull := ComputeBackoff(FullJitter, attempt, cfg, 0)
		if sleepFull < 0 || sleepFull > cfg.Cap {
			t.Fatalf("FullJitter sleep out of bounds: %v", sleepFull)
		}

		// Equal Jitter: min/2 <= sleep <= min
		sleepEqual := ComputeBackoff(EqualJitter, attempt, cfg, 0)
		if sleepEqual < 0 || sleepEqual > cfg.Cap {
			t.Fatalf("EqualJitter sleep out of bounds: %v", sleepEqual)
		}

		// No Jitter: exactly min(cap, base * 2^attempt)
		sleepNo := ComputeBackoff(NoJitter, attempt, cfg, 0)
		if sleepNo < cfg.Base || sleepNo > cfg.Cap {
			t.Fatalf("NoJitter sleep out of bounds: %v", sleepNo)
		}
	}
}

func TestDecorrelatedJitter_Bounds(t *testing.T) {
	cfg := Config{
		Base: 50 * time.Millisecond,
		Cap:  500 * time.Millisecond,
	}

	prev := cfg.Base
	for i := 0; i < 5; i++ {
		sleep := ComputeBackoff(DecorrelatedJitter, i, cfg, prev)
		if sleep < cfg.Base || sleep > cfg.Cap {
			t.Fatalf("DecorrelatedJitter out of bounds: %v", sleep)
		}
		prev = sleep
	}
}

func TestComputeBackoff_UnknownStrategy(t *testing.T) {
	cfg := Config{
		Base: 100 * time.Millisecond,
		Cap:  2 * time.Second,
	}

	sleep := ComputeBackoff(BackoffStrategy("Unknown"), 1, cfg, 0)
	expected := 200 * time.Millisecond
	if sleep != expected {
		t.Fatalf("expected unknown strategy default fallback to %v, got %v", expected, sleep)
	}
}
