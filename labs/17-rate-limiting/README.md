# Lab 17: Rate Limiting Patterns

Rate limiting is essential for protecting services from abuse and overload.

## Prompt 052 — Token Bucket

### Implementation

`token_bucket.go` implements a token bucket rate limiter with:
- Configurable capacity and refill rate
- Mock clock for deterministic, non-flaky tests

### Testing Without `time.Sleep`

```go
// Bad (flaky):
time.Sleep(1 * time.Second)
allowed = bucket.Use(10)

// Good (deterministic):
tb := NewTestableTokenBucket(100, 10, time.Now()) // 100 cap, 10/sec
tb.SetTokens(0)
tb.Refill(time.Second) // Advance mock clock + refill tokens
allowed, _ := tb.Use(10) // Deterministic result
```

### Burst Behavior

Token buckets allow bursts up to capacity:
```
Capacity: 100 tokens
Rate: 10 tokens/second

Request at T=0: 100 tokens available (burst of 100)
Request at T=0.5s: 105 tokens (refilled 5) → capped at 100
Request at T=2s: 120 tokens (refilled 20) → capped at 100
```

## Prompt 053 — Per-Client Rate Limiting

### Implementation

`per_client.go` provides rate limiting per client with safety measures.

### Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| **Unbounded Map Growth** | Memory exhaustion attack | Max client limit, TTL-based eviction |
| **IP Spoofing** | Bypass rate limit | Use session/token ID, validate X-Forwarded-For |
| **Proxy Headers** | Shared limits for wrong clients | Configure trusted proxies at gateway |
| **NAT Sharing** | One IP = one limit for many users | Combine IP with device fingerprint |

### Client Identification Strategy

```go
type ClientIdentifier struct {
    SessionID string  // Highest priority: session token
    UserID    string  // Auth token user ID
    RemoteIP  string  // Fallback: direct connection IP
    Headers   string  // User-Agent, etc. for fingerprinting
}
```

## Prompt 054 — Distributed Rate Limiting (Redis)

### Redis Lua Scripts

Lua scripts in Redis provide atomic operations, crucial for distributed rate limiting.

#### Token Bucket (Lua)

```lua
-- Atomic check-and-decrement
local tokens = redis.call('HGET', key, 'tokens') or capacity
local elapsed = now - last_update
tokens = min(capacity, tokens + (elapsed * rate))
if tokens >= requested then
    redis.call('HSET', key, 'tokens', tokens - requested)
    return 1  -- allowed
else
    return 0  -- rate limited
end
```

#### Sliding Window (ZSET)

```lua
-- Remove expired entries and count current window
redis.call('ZREMRANGEBYSCORE', key, 0, now-window)
local count = redis.call('ZCARD', key)
if count < limit then
    redis.call('ZADD', key, now, now)
    return 1  -- allowed
else
    return 0  -- rate limited
end
```

### Consistency vs Latency Trade-offs

| Approach | Latency | Consistency | Best For |
|----------|---------|-------------|----------|
| Redis Standalone | 1-3ms | Strong | Critical endpoints |
| Redis Cluster | 1-3ms | Eventual | High availability |
| Local Memory | ~0ms | Weak | Read-heavy, can tolerate abuse |
| Hybrid | 0ms (fallback) | Eventual | Resilience |

### Atomic Operations

Using Redis Lua ensures atomicity:
- No race conditions between check and decrement
- Multiple app instances see consistent state
- Lua scripts run atomically on Redis server

## Running Tests

```bash
go test ./labs/17-rate-limiting/... -v
```

## Production Considerations

1. **Redis as Critical Service**: Implement circuit breaker for Redis client
2. **Failover**: Fall back to local memory when Redis unavailable
3. **Monitoring**: Track bucket deficit, eviction rates, limit violations
4. **Key Naming**: Use consistent key format: `ratelimit:{client_id}`
5. **TTL**: Set appropriate expiration to prevent stale data accumulation

---

## Navigasi

- **Previous**: [Lab 16 — Circuit Breaker](../16-circuit-breaker/)
- **Next**: [Lab 18 — Pessimistic Locking](../18-pessimistic-locking/)