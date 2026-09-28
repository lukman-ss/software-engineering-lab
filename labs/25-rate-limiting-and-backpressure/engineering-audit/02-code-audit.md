# Code Audit

## Finding 1

Location: `internal/ratelimit/bucket.go:30-46` (TokenBucket.AllowN)
Claimed Behavior: Token bucket refills fractional tokens based on elapsed wall clock time and caps at capacity $B$, allowing atomic token subtraction.
Observed Implementation: Mutex-guarded calculation using `now.Sub(tb.lastRefill).Seconds() * refillRate`, capped at `tb.capacity`. Correctly validates remaining tokens $\ge n$.
Assessment: PASS
Severity: LOW
Notes: Clean standard library implementation with monotonic time safety.

## Finding 2

Location: `internal/ratelimit/bucket.go:98-115` (LeakyBucket.Allow)
Claimed Behavior: Leaky bucket drains water continuously at constant rate and admits new requests only if adding 1 unit does not exceed capacity.
Observed Implementation: Mutex-guarded calculation decrementing water by `elapsed * leakRate`, bounded below by 0. Rejects when `lb.water + 1.0 > lb.capacity`.
Assessment: PASS
Severity: LOW
Notes: Accurately demonstrates traffic smoothing without burst capacity above ceiling.

## Finding 3

Location: `internal/ratelimit/registry.go:21-37` (Registry.Get)
Claimed Behavior: Thread-safe per-tenant limiter retrieval preventing shared IP / CGNAT collisions (RFC 6598).
Observed Implementation: Uses double-checked locking with `RWMutex` (read lock fast path, upgrade to write lock if missing). Prevents race condition during bucket instantiation.
Assessment: PASS
Severity: LOW
Notes: Verified by `TestRegistry_ConcurrentSameKeyGet`.

## Finding 4

Location: `internal/backpressure/queue.go:67-85` (BoundedQueue.TrySubmit)
Claimed Behavior: Non-blocking enqueue that drops excess requests immediately with `ErrQueueFull` when capacity is reached.
Observed Implementation: Uses a non-blocking `select` on buffered channel with `default` fallback updating atomic counters (`accepted` vs `rejected`). Includes `stopMu` guard against submission during queue shutdown.
Assessment: PASS
Severity: LOW
Notes: Implements fast load shedding preventing unbounded latency accumulation.

## Finding 5

Location: `internal/backpressure/queue.go:91-103` (BoundedQueue.Stop)
Claimed Behavior: Graceful shutdown draining/canceling workers without deadlocks or double-close panics.
Observed Implementation: Protected by `stopMu` with `stopped` bool flag. Calls `cancel()`, closes channel, and waits for `sync.WaitGroup` worker completion.
Assessment: PASS
Severity: LOW
Notes: Robust lifecycle management.

## Finding 6

Location: `internal/retry/backoff.go:24-63` (ComputeBackoff)
Claimed Behavior: Computes backoff duration matching AWS Architecture blog (Marc Brooker) for NoJitter, FullJitter, EqualJitter, and DecorrelatedJitter.
Observed Implementation: Accurately evaluates `temp = min(cap, base * 2^attempt)`. Full jitter generates uniform random in $[0, temp]$; equal jitter generates $temp/2 + [0, temp/2]$; decorrelated jitter generates $[base, prevSleep * 3]$.
Assessment: PASS
Severity: LOW
Notes: Implementation strictly mirrors AWS specifications.

## Finding 7

Location: `internal/httputil/middleware.go:17-40` (RateLimitMiddleware)
Claimed Behavior: HTTP middleware returning standard RFC 6585 status 429 and `Retry-After` header when rate limit is exceeded.
Observed Implementation: Extracts tenant key from header `X-API-Key` (falls back to `anonymous`), calls `bucket.Allow()`, and upon failure sets `Retry-After: <seconds>` header and writes HTTP 429 JSON response.
Assessment: PASS
Severity: LOW
Notes: Verified by HTTP handler unit tests.
