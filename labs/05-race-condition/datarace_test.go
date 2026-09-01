package race

import (
	"sync"
	"testing"
)

// TestUnsafeCounter_Race deliberately demonstrates a Go memory data race.
// Run with `go test -race` — this test is EXPECTED to fail under the race detector.
// Its purpose is to show that `counter++` from multiple goroutines without
// synchronization triggers the race detector. This is a demo, not a bug.
func TestUnsafeCounter_Race(t *testing.T) {
	c := &UnsafeCounter{}
	var wg sync.WaitGroup
	workers := 10
	iterations := 1000

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				c.Inc()
			}
		}()
	}
	wg.Wait()

	// With -race, this test will fail on race detector.
	// Final value will likely be less than workers * iterations due to lost updates.
	t.Logf("Unsafe counter final value: %d (expected %d)", c.Value(), workers*iterations)
}

func TestMutexCounter(t *testing.T) {
	c := &MutexCounter{}
	var wg sync.WaitGroup
	workers := 10
	iterations := 1000

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				c.Inc()
			}
		}()
	}
	wg.Wait()

	if c.Value() != workers*iterations {
		t.Errorf("got %d, want %d", c.Value(), workers*iterations)
	}
}

func TestAtomicCounter(t *testing.T) {
	c := &AtomicCounter{}
	var wg sync.WaitGroup
	workers := 10
	iterations := 1000

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				c.Inc()
			}
		}()
	}
	wg.Wait()

	if c.Value() != int64(workers*iterations) {
		t.Errorf("got %d, want %d", c.Value(), workers*iterations)
	}
}

func TestChannelCounter(t *testing.T) {
	c := NewChannelCounter()
	var wg sync.WaitGroup
	workers := 10
	iterations := 1000

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				c.Inc()
			}
		}()
	}
	wg.Wait()

	if c.Value() != workers*iterations {
		t.Errorf("got %d, want %d", c.Value(), workers*iterations)
	}
}
