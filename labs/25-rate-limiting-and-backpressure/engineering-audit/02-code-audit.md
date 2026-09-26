# Code Audit

Target Lab: `labs/25-rate-limiting-and-backpressure`

## Finding 1

Location: `internal/ratelimit/bucket.go:30-46`
Claimed Behavior: Thread-safe token bucket algorithm with continuous rate refill.
Observed Implementation: Uses `sync.Mutex` to guard token count calculation `tb.tokens + elapsed*tb.refillRate` bounded by `capacity`. Continuous elapsed time computed via `now.Sub(tb.lastRefill).Seconds()`.
Assessment: PASS
Severity: LOW
Notes: Correct continuous refill logic using Go's monotonic time.

## Finding 2

Location: `internal/ratelimit/bucket.go:98-115`
Claimed Behavior: Thread-safe leaky bucket traffic smoothing algorithm.
Observed Implementation: Uses `sync.Mutex` to subtract leaked water based on elapsed time `elapsed*lb.leakRate`. Increments water level by 1 if capacity is not exceeded.
Assessment: PASS
Severity: LOW
Notes: Accurately smooths incoming rate and sheds bursts when water reaches capacity.

## Finding 3

Location: `internal/ratelimit/registry.go:21-37`
Claimed Behavior: Per-tenant rate limiter registry guarding against CGNAT IP collisions (RFC 6598).
Observed Implementation: Uses `sync.RWMutex` with double-check locking pattern to return or initialize per-tenant `TokenBucket` instances mapped by tenant string keys.
Assessment: PASS
Severity: LOW
Notes: Guarantees tenant isolation under concurrent access.

## Finding 4

Location: `internal/backpressure/queue.go:61-70`
Claimed Behavior: Bounded queue backpressure returning fast rejection when full.
Observed Implementation: `TrySubmit` uses non-blocking `select` with `default:` block on a buffered Go channel. Increments atomic stats and returns `ErrQueueFull` immediately without blocking caller thread.
Assessment: PASS
Severity: LOW
Notes: Perfect implementation of non-blocking load shedding.

## Finding 5

Location: `internal/retry/backoff.go:24-63`
Claimed Behavior: AWS backoff strategies (NoJitter, FullJitter, EqualJitter, DecorrelatedJitter) matching Marc Brooker specs.
Observed Implementation: Implements exact AWS mathematical formulas:
- NoJitter: $\min(cap, base \cdot 2^{attempt})$
- FullJitter: $rand(0, \min(cap, base \cdot 2^{attempt}))$
- EqualJitter: $temp/2 + rand(0, temp/2)$
- DecorrelatedJitter: $\min(cap, base + rand(0, prev \cdot 3 - base))$
Assessment: PASS
Severity: LOW
Notes: Correct implementation of randomized exponential backoff jitter algorithms.

## Finding 6

Location: `internal/httputil/middleware.go:17-39`
Claimed Behavior: HTTP middleware returning RFC 6585 standard HTTP status 429 and `Retry-After` header.
Observed Implementation: Reads `X-API-Key` header (fallback `anonymous`), evaluates tenant bucket `Allow()`, sets `Content-Type: application/json` and `Retry-After` header with integer seconds, writes status `429 Too Many Requests`, and encodes JSON error payload.
Assessment: PASS
Severity: LOW
Notes: Fully complies with RFC 6585.
