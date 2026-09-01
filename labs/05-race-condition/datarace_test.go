package race

import (
	"sync"
	"testing"
)

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
