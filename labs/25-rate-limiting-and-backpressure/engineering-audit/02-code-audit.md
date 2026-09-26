# Code Audit

## Finding 1

Location: `internal/ratelimit/bucket.go:8-79`
Claimed Behavior: Token Bucket allows bursts up to capacity and fractional replenishments using monotonic elapsed time.
Observed Implementation: `TokenBucket` uses `sync.Mutex` synchronization, calculates elapsed delta via `time.Now().Sub(...)`, caps at capacity, and implements `RetryAfterSeconds` for RFC 6585 compliance.
Assessment: PASS
Severity: LOW
Notes: Implementation correctly handles time deltas and float64 token balance.

## Finding 2

Location: `internal/ratelimit/bucket.go:81-121`
Claimed Behavior: Leaky Bucket smooths inbound bursts by tracking accumulated water volume and draining at constant leak rate.
Observed Implementation: `LeakyBucket` synchronizes on `sync.Mutex`, computes leaked volume over elapsed time, and rejects incoming requests when current water volume plus 1 exceeds capacity.
Assessment: PASS
Severity: LOW
Notes: Thread-safe and accurately models leaky bucket smoothing behavior.

## Finding 3

Location: `internal/ratelimit/registry.go:6-37`
Claimed Behavior: Per-tenant rate limiter registry to ensure tenant isolation and guard against RFC 6598 CGNAT IP collisions.
Observed Implementation: `Registry` wraps a sync RWMutex guarded map mapping tenant API keys to isolated `TokenBucket` instances with double-checked locking on bucket creation.
Assessment: PASS
Severity: LOW
Notes: Clean thread-safe tenant isolation.

## Finding 4

Location: `internal/backpressure/queue.go:14-80`
Claimed Behavior: Bounded Queue Backpressure with fast rejection via non-blocking submit.
Observed Implementation: `BoundedQueue` uses a buffered channel with non-blocking `select` (`TrySubmit`) returning `ErrQueueFull` immediately on saturation, tracking stats with atomic counters, and cleanly terminating workers on `Stop()`.
Assessment: PASS
Severity: LOW
Notes: Avoids blocking callers or memory growth during overload.

## Finding 5

Location: `internal/retry/backoff.go:24-63`
Claimed Behavior: AWS Exponential Backoff strategies with Full Jitter, Equal Jitter, No Jitter, and Decorrelated Jitter.
Observed Implementation: Implements Marc Brooker's mathematical formulas for exponential backoff with randomized jitter variants bounded by configured `Base` and `Cap`.
Assessment: PASS
Severity: LOW
Notes: Correctly computes jittered backoff bounds.

## Finding 6

Location: `internal/httputil/middleware.go:17-40`
Claimed Behavior: HTTP middleware returning RFC 6585 `429 Too Many Requests` with `Retry-After` header.
Observed Implementation: Extracts tenant key (`X-API-Key`), checks token availability, and returns HTTP status 429 with integer `Retry-After` header and structured JSON payload when quota exhausted.
Assessment: PASS
Severity: LOW
Notes: Fully aligns with RFC 6585 specification.
