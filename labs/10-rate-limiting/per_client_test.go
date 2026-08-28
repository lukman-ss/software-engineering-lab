package ratelimit_test

import (
	"testing"
	"time"

	"github.com/lukman-ss/software-engineering-lab/labs/10-rate-limiting"
)

func TestPerClientLimiter(t *testing.T) {
	limiter := ratelimit.NewPerClientLimiter(ratelimit.PerClientConfig{
		GlobalRate: 0, // Unused in this test
	})

	cfg := ratelimit.PerClientConfig{
		PerClientRate:     10,
		PerClientCapacity: 5,
	}

	// Client 1 - Uses up its burst
	bucket1 := limiter.TokenBucket("client-1", cfg)
	allowed, _ := bucket1.Use(5)
	if !allowed {
		t.Error("client 1 should be allowed 5 tokens")
	}

	// Client 1 - Empty
	allowed, _ = bucket1.Use(1)
	if allowed {
		t.Error("client 1 should be rejected (empty)")
	}

	// Client 2 - Independent, full capacity
	bucket2 := limiter.TokenBucket("client-2", cfg)
	allowed, _ = bucket2.Use(5)
	if !allowed {
		t.Error("client 2 should be allowed 5 tokens (independent of client 1)")
	}

	// Verify isolation
	if limiter.ClientCount() != 2 {
		t.Errorf("expected 2 clients, got %d", limiter.ClientCount())
	}
}

func TestUnboundedMapProtection(t *testing.T) {
	// 5 max clients
	limiter := ratelimit.NewSafePerClientLimiterWithSts(5)
	cfg := ratelimit.TokenBucketConfig{
		Rate:     10,
		Capacity: 10,
	}

	// Add 5 clients
	for i := 0; i < 5; i++ {
		_, created := limiter.GetOrCreate(string(rune('A'+i)), cfg)
		if !created {
			t.Errorf("expected client %c to be created", 'A'+i)
		}
	}

	// Add 6th client - should trigger eviction
	_, created := limiter.GetOrCreate("F", cfg)
	if !created {
		t.Error("expected client F to be created (evicting one)")
	}
}
