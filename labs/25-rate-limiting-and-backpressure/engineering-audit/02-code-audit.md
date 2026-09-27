# Engineering Code Audit

## Finding 1

Location: `internal/ratelimit/bucket.go:8-79`
Claimed Behavior: Token bucket tracks token count, refills tokens based on elapsed monotonic duration, allows bursts up to capacity, and computes RFC 6585 retry-after wait duration.
Observed Implementation: `TokenBucket` uses `sync.Mutex` to synchronize fractional token refills calculated from `time.Now().Sub(lastRefill)`. Tokens are capped at capacity. `AllowN` decrements tokens accurately when available. `RetryAfterSeconds` calculates required refill duration and rounds up seconds.
Assessment: PASS
Severity: LOW
Notes: Implementation correctly uses Go monotonic clock timestamps and avoids division by zero.

## Finding 2

Location: `internal/ratelimit/bucket.go:81-121`
Claimed Behavior: Leaky bucket enforces constant drain rate, smoothing burst traffic and rejecting requests when capacity is exceeded.
Observed Implementation: `LeakyBucket` tracks water level with `sync.Mutex`, leaks volume based on elapsed time, floors water level at 0, and permits execution if adding unit volume stays within capacity.
Assessment: PASS
Severity: LOW
Notes: Clean mathematical implementation of leaky bucket as a meter.

## Finding 3

Location: `internal/ratelimit/registry.go:6-37`
Claimed Behavior: Multi-tenant registry maintains separate token buckets keyed by API key or tenant ID to prevent CGNAT IP collision (RFC 6598).
Observed Implementation: `Registry` wraps a bucket map with double-checked locking using `sync.RWMutex` to retrieve or lazily instantiate per-tenant `TokenBucket` instances.
Assessment: PASS
Severity: LOW
Notes: Thread-safe map access; no lock contention on read path.

## Finding 4

Location: `internal/backpressure/queue.go:17-94`
Claimed Behavior: Bounded queue provides non-blocking submission (`TrySubmit`), immediately rejecting excess jobs with `ErrQueueFull` when buffer capacity is reached, and gracefully shutting down workers.
Observed Implementation: `BoundedQueue` uses a buffered channel with non-blocking `select default` in `TrySubmit`. Rejection and acceptance stats are tracked atomically with `sync/atomic`. `Stop` uses `CompareAndSwap` on atomic boolean, cancels context, closes channel, and waits for worker goroutines with `sync.WaitGroup`.
Assessment: PASS
Severity: LOW
Notes: Zero deadlocks; idempotent stop; handles submission after stop with `ErrQueueStopped`.

## Finding 5

Location: `internal/retry/backoff.go:24-63`
Claimed Behavior: Implements AWS Marc Brooker exponential backoff jitter strategies: No Jitter, Full Jitter, Equal Jitter, and Decorrelated Jitter.
Observed Implementation: `ComputeBackoff` calculates exponential backoff bounded by `Cap`, and applies jitter calculations conforming to AWS Architecture Blog formulas.
Assessment: PASS
Severity: LOW
Notes: Formulas strictly adhere to standard reference implementations.

## Finding 6

Location: `internal/httputil/middleware.go:17-40`
Claimed Behavior: HTTP middleware extracts tenant key from `X-API-Key` header, enforces rate limiting, and returns standard RFC 6585 HTTP 429 response with `Retry-After` header.
Observed Implementation: `RateLimitMiddleware` inspects header with fallback to `"anonymous"`, calls bucket `Allow()`, and writes status code `429` with header `Retry-After` and JSON error payload when exhausted.
Assessment: PASS
Severity: LOW
Notes: Complies with RFC 6585.
