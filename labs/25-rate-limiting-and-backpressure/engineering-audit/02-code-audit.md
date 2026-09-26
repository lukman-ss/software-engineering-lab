# Code Audit

Target Lab: `labs/25-rate-limiting-and-backpressure`

## Finding 1

Location: `internal/ratelimit/bucket.go:30-46`
Claimed Behavior: Thread-safe token bucket with fractional token replenishment and burst capacity allowance.
Observed Implementation: Mutex locking (`tb.mu.Lock()`) protects elapsed time calculations and token replenishment. `AllowN` correctly calculates token refill based on elapsed duration and caps at capacity.
Assessment: PASS
Severity: LOW
Notes: Monotonic time subtraction via Go `time.Time` avoids drift.

## Finding 2

Location: `internal/ratelimit/bucket.go:98-115`
Claimed Behavior: Thread-safe leaky bucket enforcing constant drain rate and rejecting burst overflow.
Observed Implementation: `Allow()` acquires lock, calculates leaked water based on elapsed time, floors water level at 0, and accepts incoming request if current water + 1 <= capacity.
Assessment: PASS
Severity: LOW
Notes: Clean implementation matching leaky bucket specifications.

## Finding 3

Location: `internal/backpressure/queue.go:61-70`
Claimed Behavior: Non-blocking bounded queue load shedding returning fast rejection (`ErrQueueFull`).
Observed Implementation: `TrySubmit` uses non-blocking select statement with `default` case to reject excess items immediately and increment atomic counters.
Assessment: PASS
Severity: LOW
Notes: Prevents channel blocking or memory accumulation.

## Finding 4

Location: `internal/retry/backoff.go:24-63`
Claimed Behavior: AWS Marc Brooker exponential backoff algorithms (NoJitter, FullJitter, EqualJitter, DecorrelatedJitter).
Observed Implementation: Implements exact formulas: Full Jitter uses `rand.Float64() * temp`, Equal Jitter splits base exponentially into equal deterministic and random halves, Decorrelated Jitter scales relative to previous sleep time.
Assessment: PASS
Severity: LOW
Notes: Uses standard math/rand.

## Finding 5

Location: `internal/httputil/middleware.go:17-39`
Claimed Behavior: HTTP middleware returning RFC 6585 `429 Too Many Requests` status code and `Retry-After` header.
Observed Implementation: Extracts tenant key from `X-API-Key` header, queries registry bucket, and sets status 429 along with `Retry-After` header and JSON error payload when `Allow()` is false.
Assessment: PASS
Severity: LOW
Notes: Fully complies with HTTP 429 semantics and RFC 6598 tenant isolation.
