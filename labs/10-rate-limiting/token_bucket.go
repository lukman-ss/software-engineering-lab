// Package ratelimit implements rate limiting patterns.
// This lab covers token bucket rate limiting, per-client limits,
// and distributed rate limiting with Redis.
package ratelimit

import (
	"errors"
	"sync"
	"time"
)

// ErrRateLimitExceeded is returned when a request would exceed the rate limit.
var ErrRateLimitExceeded = errors.New("rate limit exceeded")

// Clock provides a time source for testability.
type Clock interface {
	Now() time.Time
}

// SystemClock uses the real system clock.
type SystemClock struct{}

func (SystemClock) Now() time.Time { return time.Now() }

// TokenBucket implements a token bucket rate limiter.
// Tokens are added at a fixed rate up to a maximum capacity.
type TokenBucket struct {
	mu sync.Mutex

	capacity    float64
	tokens      float64
	ratePerSec  float64
	lastRefill  time.Time
	closer      chan struct{}

	// Clock for testability (Prompt 052)
	clock Clock
}

// TokenBucketConfig holds configuration for a token bucket.
type TokenBucketConfig struct {
	Capacity   float64       // Maximum tokens (burst size)
	Rate       float64       // Tokens added per second
	Now        func() time.Time // Time source for testing
}

// NewTokenBucket creates a new token bucket limiters.
func NewTokenBucket(config TokenBucketConfig) *TokenBucket {
	now := config.Now
	if now == nil {
		now = time.Now
	}

	return &TokenBucket{
		capacity:   config.Capacity,
		tokens:     config.Capacity, // Start full
		ratePerSec: config.Rate,
		lastRefill: now(),
		closer:     make(chan struct{}),
		clock:      SystemClock{},
	}
}

// Use attempts to consume n tokens, returns true if allowed.
func (tb *TokenBucket) Use(n float64) (bool, error) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()

	// Refill tokens based on elapsed time
	elapsed := now.Sub(tb.lastRefill).Seconds()
	tb.tokens += elapsed * tb.ratePerSec

	// Cap at maximum capacity
	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}
	tb.lastRefill = now()

	// Check if we have enough tokens
	if tb.tokens >= n {
		tb.tokens -= n
		return true, nil
	}

	return false, ErrRateLimitExceeded
}

// Tokens returns current token count (for testing/debugging).
func (tb *TokenBucket) Tokens() float64 {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	return tb.tokens
}

// MockClock provides controllable time for testing without flaky time.Sleep.
type MockClock struct {
	mu   sync.Mutex
	now  time.Time
}

// NewMockClock creates a mock clock starting at t.
func NewMockClock(t time.Time) *MockClock {
	return &MockClock{now: t}
}

// Now implements Clock interface.
func (mc *MockClock) Now() time.Time {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	return mc.now
}

// Advance moves the clock forward by duration.
func (mc *MockClock) Advance(d time.Duration) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.now = mc.now.Add(d)
}

// Refill advances time and returns the number of new tokens.
// Returns the actual number of tokens added.
func (mc *MockClock) Refill(bucket *TokenBucket, duration time.Duration) float64 {
	mc.mu.Lock()
	mc.now = mc.now.Add(duration)
	current := mc.now
	mc.mu.Unlock()

	bucket.mu.Lock()
	defer bucket.mu.Unlock()
	elapsed := current.Sub(bucket.lastRefill).Seconds()
	newTokens := elapsed * bucket.ratePerSec
	bucket.tokens += newTokens
	if bucket.tokens > bucket.capacity {
		bucket.tokens = bucket.capacity
	}
	bucket.lastRefill = current
	return newTokens
}

// SetTokens sets the token count (for testing).
func (tb *TokenBucket) SetTokens(n float64) {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.tokens = n
}

// TestableTokenBucket wraps TokenBucket with mock clock support.
type TestableTokenBucket struct {
	*TokenBucket
	clock *MockClock
}

// NewTestableTokenBucket creates a rate limiter with mock clock for testing.
// This avoids flaky time.Sleep in tests.
func NewTestableTokenBucket(capacity, rate float64, initialTime time.Time) *TestableTokenBucket {
	mockClock := NewMockClock(initialTime)
	tb := &TokenBucket{
		capacity:   capacity,
		tokens:     capacity, // Start full
		ratePerSec: rate,
		lastRefill: mockClock.now,
		closer:     make(chan struct{}),
		clock:      mockClock,
	}
	return &TestableTokenBucket{
		TokenBucket: tb,
		clock:       mockClock,
	}
}

// Clock returns the mock clock for precise time control in tests.
func (t *TestableTokenBucket) Clock() *MockClock {
	return t.clock
}

// Refill advances time by the given duration.
func (t *TestableTokenBucket) Refill(d time.Duration) {
	t.clock.Refill(t.TokenBucket, d)
}