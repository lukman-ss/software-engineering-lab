# Engineering Code Audit

Target Lab: `labs/25-rate-limiting-and-backpressure`

## Finding 1

Location: `internal/ratelimit/bucket.go:8-52` (`TokenBucket`)
Claimed Behavior: Token bucket tracks token availability with capacity cap and fractional refill based on elapsed time; thread-safe.
Observed Implementation: Uses `sync.Mutex` around state mutations, calculates elapsed monotonic time via `time.Now().Sub(lastRefill).Seconds()`, clamps to `capacity`, deducts tokens if available.
Assessment: PASS
Severity: LOW
Notes: Clean implementation with standard token bucket state transition mechanics.

## Finding 2

Location: `internal/ratelimit/bucket.go:54-79` (`TokenBucket.RetryAfterSeconds`)
Claimed Behavior: Calculates estimated waiting seconds until next token is available.
Observed Implementation: Correctly computes `needed / refillRate`, rounding up to the nearest integer second for RFC 6585 compliance.
Assessment: PASS
Severity: LOW
Notes: Handles edge cases where tokens are already sufficient (`return 0`).

## Finding 3

Location: `internal/ratelimit/bucket.go:81-121` (`LeakyBucket`)
Claimed Behavior: Leaky bucket enforces leak rate, allows unit additions when current water + 1 <= capacity, drains over time.
Observed Implementation: Uses `sync.Mutex`, computes leaked volume over elapsed time, floors water level at 0, admits or rejects additions.
Assessment: PASS
Severity: LOW
Notes: Correctly distinguishes leaky bucket traffic smoothing from token bucket burst allowance.

## Finding 4

Location: `internal/ratelimit/registry.go:6-37` (`Registry`)
Claimed Behavior: Manages per-tenant rate limiters via `X-API-Key` or tenant ID to prevent carrier-grade NAT (RFC 6598) collisions.
Observed Implementation: Double-checked locking with `sync.RWMutex` protecting `map[string]*TokenBucket`.
Assessment: PASS
Severity: LOW
Notes: Prevents race conditions during dynamic tenant initialization.

## Finding 5

Location: `internal/backpressure/queue.go:17-94` (`BoundedQueue`)
Claimed Behavior: Fixed-capacity channel worker pool rejecting jobs fast when full, providing backpressure load shedding.
Observed Implementation: Non-blocking channel write via `select-default` inside `TrySubmit`. Thread-safe atomic counters for metrics (`accepted`, `rejected`, `processed`). Clean lifecycle shutdown via `atomic.Bool`, `context.CancelFunc`, channel closure, and `sync.WaitGroup`.
Assessment: PASS
Severity: LOW
Notes: Bounded channel prevents unbounded memory growth under overload.

## Finding 6

Location: `internal/retry/backoff.go:24-63` (`ComputeBackoff`)
Claimed Behavior: Calculates exponential backoff with AWS Marc Brooker jitter algorithms (NoJitter, FullJitter, EqualJitter, DecorrelatedJitter).
Observed Implementation: Formulations match AWS Architecture blog specifications:
- `temp = min(cap, base * 2^attempt)`
- `FullJitter`: `random(0, temp)`
- `EqualJitter`: `temp/2 + random(0, temp/2)`
- `DecorrelatedJitter`: `min(cap, random(base, prevSleep * 3))`
Assessment: PASS
Severity: LOW
Notes: Mathematically sound and adheres to research evidence.

## Finding 7

Location: `internal/httputil/middleware.go:17-40` (`RateLimitMiddleware`)
Claimed Behavior: Enforces tenant-based rate limits returning RFC 6585 status `429 Too Many Requests` with `Retry-After` header.
Observed Implementation: Extracts `X-API-Key` (defaults to `anonymous`), queries registry, inspects token availability, sets `Retry-After` and `Content-Type: application/json`, responds with `http.StatusTooManyRequests` on limit exceeded.
Assessment: PASS
Severity: LOW
Notes: Compliant with RFC 6585.
