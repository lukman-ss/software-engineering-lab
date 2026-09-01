package race

import (
	"sync"
	"sync/atomic"
)

// 1. Unsafe Counter (Data Race)
type UnsafeCounter struct {
	val int
}

func (c *UnsafeCounter) Inc() {
	c.val++ // Data race here
}

func (c *UnsafeCounter) Value() int {
	return c.val
}

// 2. Mutex Counter
type MutexCounter struct {
	mu  sync.Mutex
	val int
}

func (c *MutexCounter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.val++
}

func (c *MutexCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.val
}

// 3. Atomic Counter
type AtomicCounter struct {
	val int64
}

func (c *AtomicCounter) Inc() {
	atomic.AddInt64(&c.val, 1)
}

func (c *AtomicCounter) Value() int64 {
	return atomic.LoadInt64(&c.val)
}

// 4. Channel Ownership Counter
type ChannelCounter struct {
	incChan chan struct{}
	valChan chan int
}

func NewChannelCounter() *ChannelCounter {
	c := &ChannelCounter{
		incChan: make(chan struct{}),
		valChan: make(chan int),
	}
	go c.run()
	return c
}

func (c *ChannelCounter) run() {
	val := 0
	for {
		select {
		case <-c.incChan:
			val++
		case c.valChan <- val:
		}
	}
}

func (c *ChannelCounter) Inc() {
	c.incChan <- struct{}{}
}

func (c *ChannelCounter) Value() int {
	return <-c.valChan
}
