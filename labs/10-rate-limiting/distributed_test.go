package ratelimit_test

import (
	"context"
	"testing"

	"github.com/lukman-ss/software-engineering-lab/labs/10-rate-limiting"
)

func TestRedisRateLimiterConcept(t *testing.T) {
	limiter := ratelimit.NewRedisRateLimiter(100, 60) // 100 req/min

	// The Lua script ensures atomicity across multiple app instances

	// Token Bucket Lua returns: 1 if allowed, 0 if rate limited
	// It atomically:
	// 1. Reads current tokens and last_update
	// 2. Calculates tokens generated since last_update
	// 3. Checks if tokens >= requested
	// 4. Deducts tokens if allowed
	// 5. Saves new state

	t.Log("TOKEN BUCKET LUA SCRIPT:")
	t.Log(ratelimit.TokenBucketLuaScript)

	// Sliding Window ZSET returns: 1 if allowed, 0 if rate limited
	t.Log("SLIDING WINDOW LUA SCRIPT:")
	t.Log(ratelimit.SlidingWindowLuaScript)
}

func TestTradeOffs(t *testing.T) {
	trades := ratelimit.TradeOffsDocument()

	t.Log("DISTRIBUTED RATE LIMITING TRADE-OFFS:")
	for k, v := range trades {
		t.Logf("  %s: %s", k, v)
	}

	t.Log("")
	t.Log("DECISION FRAMEWORK:")
	t.Log("  - Use Redis SLIDING WINDOW for critical:login endpoints (accuracy matters)")
	t.Log("  - Use Redis TOKEN BUCKET for api:limit (bursts are OK)")
	t.Log("  - Use LOCAL MEMORY + shared cache for read-heavy endpoints")
	t.Log("  - Always implement fallback to local memory during Redis outage")
}

func TestCircuitBreakerIntegrationConcept(t *testing.T) {
	ctx := context.Background()
	_ = ctx

	// When Redis is unavailable, fall back to local memory
	// This prevents cascading failures
	t.Log("FAIL-SAFE PATTERN:")
	t.Log("  Redis unavailable -> Use local memory bucket (less accurate but functional)")
	t.Log("  Redis recovers -> Gradually re-sync local state to Redis")
}
