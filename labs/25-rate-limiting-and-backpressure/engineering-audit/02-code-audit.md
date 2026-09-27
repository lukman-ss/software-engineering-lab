# Code Audit Report

## Target Lab
`labs/25-rate-limiting-and-backpressure`

## Finding 1

Location: `internal/ratelimit/bucket.go:8-46`
Claimed Behavior: Token Bucket allows burst up to capacity and refills at constant rate without background tickers.
Observed Implementation: State is calculated lazily via `now.Sub(tb.lastRefill).Seconds() * refillRate`, capped at capacity. Thread-safety is guarded with `sync.Mutex`.
Assessment: PASS
Severity: LOW
Notes: Lazy calculation avoids background goroutines and ticker leaks. Correct handling of token clamping to capacity.

## Finding 2

Location: `internal/ratelimit/bucket.go:55-79`
Claimed Behavior: `RetryAfterSeconds` calculates necessary wait time in whole integer seconds for client retry guidance.
Observed Implementation: Calculates needed tokens divided by `refillRate`, rounding up to the next ceiling integer. Returns 0 when tokens are already available.
Assessment: PASS
Severity: LOW
Notes: Complies with HTTP RFC 6585 `Retry-After` seconds integer format.

## Finding 3

Location: `internal/ratelimit/bucket.go:81-121`
Claimed Behavior: Leaky Bucket acts as a traffic smoother enforcing leak rate $R$, rejecting bursts exceeding water capacity.
Observed Implementation: Calculates drained water via elapsed time multiplied by `leakRate`, clamped at zero. Accepts requests if `water + 1 <= capacity`, otherwise returns false.
Assessment: PASS
Severity: LOW
Notes: Correctly implements the leaky bucket as a meter/traffic smoother.

## Finding 4

Location: `internal/ratelimit/registry.go:21-37`
Claimed Behavior: Registry manages per-tenant rate limiters to avoid CGNAT / shared IP cross-tenant starvation (RFC 6598).
Observed Implementation: Double-checked locking pattern using `sync.RWMutex` to retrieve or lazily create a `TokenBucket` per tenant key.
Assessment: PASS
Severity: LOW
Notes: Safe against concurrent reads and initializations. In a high-churn multi-tenant production environment, an eviction strategy (LRU/TTL) would be needed to prevent unbounded memory growth, but for this lab scope it is clean and correct.

## Finding 5

Location: `internal/backpressure/queue.go:30-94`
Claimed Behavior: Bounded Queue enforces backpressure through fast drop rejection (`TrySubmit`), graceful worker loop shutdown, and telemetry tracking.
Observed Implementation: `TrySubmit` uses non-blocking `select` channel send with a `default` case returning `ErrQueueFull`. Shutdown uses `atomic.Bool.CompareAndSwap`, context cancellation, channel close, and `sync.WaitGroup.Wait`.
Assessment: PASS
Severity: LOW
Notes: No blocking on full buffer. Atomic counters accurately track accepted, rejected, and processed counts.

## Finding 6

Location: `internal/httputil/middleware.go:17-39`
Claimed Behavior: Middleware extracts tenant key from header, evaluates rate limit, returns HTTP 429 with JSON payload and `Retry-After` header when limited.
Observed Implementation: Reads `X-API-Key` (falling back to "anonymous"), evaluates bucket, sets `Retry-After` header and returns status code `429 Too Many Requests`.
Assessment: PASS
Severity: LOW
Notes: Matches RFC 6585 requirements.

## Finding 7

Location: `internal/retry/backoff.go:23-63`
Claimed Behavior: AWS Architecture exponential backoff strategies (NoJitter, FullJitter, EqualJitter, DecorrelatedJitter) matching Marc Brooker specifications.
Observed Implementation: Implements exact formulas bounded by `cfg.Cap` and `cfg.Base`.
Assessment: PASS
Severity: LOW
Notes: Mathematical definitions match AWS Architecture reference.
