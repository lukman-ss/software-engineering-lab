# Code Audit

## Finding 1

Location: `internal/ratelimit/bucket.go:29-46` (`TokenBucket.AllowN`)
Claimed Behavior: Continuous refilling of tokens based on elapsed time up to capacity, thread-safe token deduction.
Observed Implementation: Mutex locking guards `lastRefill` and `tokens` calculation. Monotonic clock arithmetic (`time.Now().Sub(...)`) updates tokens up to `capacity`.
Assessment: PASS
Severity: LOW
Notes: Correct continuous refill logic.

## Finding 2

Location: `internal/ratelimit/bucket.go:98-115` (`LeakyBucket.Allow`)
Claimed Behavior: Leak rate deduction based on elapsed time, rejecting requests when water exceeds capacity.
Observed Implementation: Mutex locking guards `lastLeak` and `water` level updates. Water drains over time down to zero, incremented by 1 per allowed request.
Assessment: PASS
Severity: LOW
Notes: Properly models traffic smoothing.

## Finding 3

Location: `internal/ratelimit/registry.go:21-37` (`Registry.Get`)
Claimed Behavior: Thread-safe double-checked locking for tenant token bucket retrieval/initialization.
Observed Implementation: Initial `RLock` read check, followed by full `Lock` and double check before allocation.
Assessment: PASS
Severity: LOW
Notes: Prevents tenant IP/key CGNAT collision issues (RFC 6598 context).

## Finding 4

Location: `internal/backpressure/queue.go:66-81` (`BoundedQueue.TrySubmit`)
Claimed Behavior: Non-blocking job submission returning `ErrQueueFull` when queue buffer is full.
Observed Implementation: Select block with non-blocking send to channel (`bq.queue <- job`) and `default:` fallback incrementing atomic `rejected` counter and returning `ErrQueueFull`.
Assessment: PASS
Severity: LOW
Notes: Zero-allocation fast rejection prevents unbounded queue depth.

## Finding 5

Location: `internal/retry/backoff.go:24-63` (`ComputeBackoff`)
Claimed Behavior: Exponential backoff supporting AWS Full Jitter, Equal Jitter, No Jitter, and Decorrelated Jitter matching Marc Brooker specification.
Observed Implementation: Exponential ceiling cap calculation (`math.Min(capFloat, expBackoff)`), followed by correct randomized interval formulas.
Assessment: PASS
Severity: LOW
Notes: Accurately implements AWS Architecture specifications.

## Finding 6

Location: `internal/httputil/middleware.go:17-40` (`RateLimitMiddleware`)
Claimed Behavior: HTTP 429 response generation compliant with RFC 6585 with `Retry-After` header.
Observed Implementation: Extracts `X-API-Key` (defaults to "anonymous"), checks tenant bucket, sets `Content-Type: application/json`, `Retry-After` header, and status code 429.
Assessment: PASS
Severity: LOW
Notes: Full RFC 6585 compliance.
