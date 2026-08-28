package ratelimit_test

import (
	"testing"
	"time"

	"github.com/lukman-ss/software-engineering-lab/labs/10-rate-limiting"
)

func TestTokenBucketAllowsWhenFull(t *testing.T) {
	start := time.Now()
	tb := ratelimit.NewTestableTokenBucket(10, 100, start) // 10 tokens, 100/sec

	// Full bucket should allow requests up to capacity
	allowed, err := tb.Use(5)
	if !allowed {
		t.Error("expected request to be allowed when bucket is full")
	}
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if tokens := tb.Tokens(); tokens != 5 {
		t.Errorf("expected 5 tokens remaining, got %.2f", tokens)
	}
}

func TestTokenBucketRejectsWhenEmpty(t *testing.T) {
	start := time.Now()
	tb := ratelimit.NewTestableTokenBucket(5, 100, start) // 5 tokens capacity
	tb.SetTokens(0)                                       // Start empty

	allowed, err := tb.Use(1)
	if allowed {
		t.Error("expected request to be denied when bucket is empty")
	}
	if err != ratelimit.ErrRateLimitExceeded {
		t.Errorf("expected ErrRateLimitExceeded, got %v", err)
	}
}

func TestTokenBucketRefill(t *testing.T) {
	start := time.Now()
	tb := ratelimit.NewTestableTokenBucket(10, 100, start) // 10 tokens, 100/sec = 100 tokens/sec

	// Use 10 tokens
	tb.SetTokens(10)
	tb.Use(10)
	if tb.Tokens() != 0 {
		t.Errorf("expected 0 tokens after using all")
	}

	// Advance time by 1 second - should get 100 new tokens (capped at 10)
	tb.Refill(time.Second)

	// Should be at capacity again
	tokens := tb.Tokens()
	if tokens < 9 || tokens > 10 {
		t.Errorf("expected ~10 tokens after 1 second refill, got %.2f", tokens)
	}
}

func TestTokenBurst(t *testing.T) {
	start := time.Now()
	tb := ratelimit.NewTestableTokenBucket(100, 10, start) // 100 capacity, 10/sec refill rate

	// Burst: use 50 tokens immediately
	tb.SetTokens(100)
	allowed, _ := tb.Use(50)
	if !allowed {
		t.Error("burst request should be allowed")
	}

	// Remaining: 50 tokens
	if tb.Tokens() != 50 {
		t.Errorf("expected 50 tokens remaining after burst, got %.2f", tb.Tokens())
	}

	// Wait 2 seconds - should get 20 more tokens (10/sec * 2)
	tb.Refill(2 * time.Second)
	if tb.Tokens() < 69 || tb.Tokens() > 70 {
		t.Errorf("unexpected token count after refill: %.2f", tb.Tokens())
	}
}

func TestTokenBucketRateLimiting(t *testing.T) {
	// Rate: 10 tokens/second
	// Capacity: 10 (burst of 10 allowed)
	start := time.Now()
	tb := ratelimit.NewTestableTokenBucket(10, 10, start) // 10 tokens/sec

	tb.SetTokens(0)

	// First refill: 0.1 seconds = 1 token
	tb.Refill(100 * time.Millisecond)

	// Still empty after tiny refill
	if tb.Tokens() < 0.5 {
		t.Error("should have some tokens after 100ms")
	}

	// Consume 1 token
	allowed, _ := tb.Use(1)
	if !allowed {
		t.Error("expected to use 1 token")
	}

	// Now refill 1 second: should have ~10 tokens (capped)
	tb.Refill(time.Second)

	// Should be at capacity
	if tb.Tokens() < 9 {
		t.Errorf("expected ~10 tokens after 1s refill, got %.2f", tb.Tokens())
	}
}
