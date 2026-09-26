# Code Audit

## Finding 1

Location: `internal/ratelimit/bucket.go:30-46` (`TokenBucket.AllowN`)
Claimed Behavior: Thread-safe token bucket replenishment and consumption.
Observed Implementation: Mutex lock guards `tokens` calculation based on `now.Sub(tb.lastRefill).Seconds()`.
Assessment: PASS
Severity: LOW
Notes: Correct continuous refill logic using standard stdlib `sync.Mutex` and monotonic time delta.

## Finding 2

Location: `internal/ratelimit/bucket.go:55-79` (`TokenBucket.RetryAfterSeconds`)
Claimed Behavior: Calculates wait time in seconds for `n` tokens.
Observed Implementation: Calculates needed tokens, divides by refill rate, and rounds up to whole seconds (`rounded++`).
Assessment: PASS
Severity: LOW
Notes: Complies with integer RFC 6585 `Retry-After` header requirements.

## Finding 3

Location: `internal/ratelimit/bucket.go:98-115` (`LeakyBucket.Allow`)
Claimed Behavior: Enforces continuous drain rate and limits maximum water capacity.
Observed Implementation: Leaks water based on elapsed time delta, caps minimum water at 0, adds 1 unit if `water + 1.0 <= capacity`.
Assessment: PASS
Severity: LOW
Notes: Implements meter/leaky bucket traffic smoothing correctly with proper mutex synchronization.

## Finding 4

Location: `internal/ratelimit/registry.go:21-37` (`Registry.Get`)
Claimed Behavior: Multi-tenant rate limit bucket registry with double-check locking.
Observed Implementation: Uses `sync.RWMutex`. Checks with RLock, upgrades to Lock on cache miss, re-checks existence before creating bucket.
Assessment: PASS
Severity: LOW
Notes: Safe against multi-tenant race conditions and prevents CGNAT IP sharing bottlenecks.

## Finding 5

Location: `internal/backpressure/queue.go:61-70` (`BoundedQueue.TrySubmit`)
Claimed Behavior: Non-blocking enqueue with fast rejection on capacity overflow.
Observed Implementation: `select` statement attempts write to `queue` buffered channel; falls through to `default` immediately returning `ErrQueueFull` and atomic metric update.
Assessment: PASS
Severity: LOW
Notes: Perfect non-blocking backpressure shedding pattern without blocking caller thread.

## Finding 6

Location: `internal/backpressure/queue.go:44-58` (`BoundedQueue.workerLoop`), `76-79` (`BoundedQueue.Stop`)
Claimed Behavior: Safe shutdown of worker pool and queue drain.
Observed Implementation: `Stop()` cancels context, closes channel, and waits on `sync.WaitGroup`.
Assessment: PASS
Severity: LOW
Notes: Worker loop handles context cancellation and channel closure safely without goroutine leak.

## Finding 7

Location: `internal/retry/backoff.go:24-63` (`ComputeBackoff`)
Claimed Behavior: Exponential backoff with AWS Jitter variants.
Observed Implementation: Implements NoJitter, FullJitter, EqualJitter, and DecorrelatedJitter per AWS Marc Brooker spec.
Assessment: PASS
Severity: LOW
Notes: Range bounds accurately guarded using stdlib `math` and `math/rand`.

## Finding 8

Location: `internal/httputil/middleware.go:17-40` (`RateLimitMiddleware`)
Claimed Behavior: HTTP middleware returning standard HTTP 429 and `Retry-After` header.
Observed Implementation: Checks header `X-API-Key`, gets bucket, returns 429 with JSON body and `Retry-After` header when limit exceeded.
Assessment: PASS
Severity: LOW
Notes: Standard RFC 6585 compliant middleware. Standard stdlib HTTP implementation.
