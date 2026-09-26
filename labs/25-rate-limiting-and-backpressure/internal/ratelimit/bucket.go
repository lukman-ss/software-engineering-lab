package ratelimit

import (
	"sync"
	"time"
)

type TokenBucket struct {
	mu         sync.Mutex
	capacity   float64
	tokens     float64
	refillRate float64 // tokens per second
	lastRefill time.Time
}

func NewTokenBucket(capacity float64, refillRate float64) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     capacity,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

func (tb *TokenBucket) Allow() bool {
	return tb.AllowN(1.0)
}

func (tb *TokenBucket) AllowN(n float64) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	tb.tokens = tb.tokens + elapsed*tb.refillRate
	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}
	tb.lastRefill = now

	if tb.tokens >= n {
		tb.tokens -= n
		return true
	}
	return false
}

func (tb *TokenBucket) Tokens() float64 {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	return tb.tokens
}

// RetryAfterSeconds calculates how long caller should wait for n tokens.
func (tb *TokenBucket) RetryAfterSeconds(n float64) int {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	tokens := tb.tokens + elapsed*tb.refillRate
	if tokens > tb.capacity {
		tokens = tb.capacity
	}

	if tokens >= n {
		return 0
	}
	needed := n - tokens
	secs := needed / tb.refillRate
	if secs <= 0 {
		return 1
	}
	rounded := int(secs)
	if float64(rounded) < secs {
		rounded++
	}
	return rounded
}

type LeakyBucket struct {
	mu         sync.Mutex
	capacity   float64
	water      float64
	leakRate   float64 // units leaked per second
	lastLeak   time.Time
}

func NewLeakyBucket(capacity float64, leakRate float64) *LeakyBucket {
	return &LeakyBucket{
		capacity: capacity,
		water:    0,
		leakRate: leakRate,
		lastLeak: time.Now(),
	}
}

func (lb *LeakyBucket) Allow() bool {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(lb.lastLeak).Seconds()
	lb.water = lb.water - elapsed*lb.leakRate
	if lb.water < 0 {
		lb.water = 0
	}
	lb.lastLeak = now

	if lb.water+1.0 <= lb.capacity {
		lb.water += 1.0
		return true
	}
	return false
}

func (lb *LeakyBucket) Water() float64 {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	return lb.water
}
