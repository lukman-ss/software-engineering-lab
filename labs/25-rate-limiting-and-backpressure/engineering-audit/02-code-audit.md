# Code Audit

Target Lab: `labs/25-rate-limiting-and-backpressure`

## Finding 1

Location: `internal/ratelimit/bucket.go:30-46`
Claimed Behavior: Token bucket allows burst up to capacity and refills fractional tokens based on elapsed monotonic time.
Observed Implementation: State protected by `sync.Mutex`. Elapsed time calculated with `time.Now().Sub(tb.lastRefill).Seconds()`. Token count clamped to `tb.capacity`. Deducts requested tokens and returns true if available.
Assessment: PASS
Severity: LOW
Notes: Correct state transition and mutex protection.

## Finding 2

Location: `internal/ratelimit/bucket.go:98-115`
Claimed Behavior: Leaky bucket enforces leak rate smoothing, rejecting bursts when water reaches capacity.
Observed Implementation: Mutex protected. Drains `elapsed * leakRate` from water down to minimum 0. Adds 1 unit if `water + 1.0 <= capacity`. Rejects when full.
Assessment: PASS
Severity: LOW
Notes: Accurately models continuous leaking bucket traffic smoother.

## Finding 3

Location: `internal/ratelimit/registry.go:21-37`
Claimed Behavior: Multi-tenant registry maintains separate token buckets per tenant key (API key / tenant ID) to avoid CGNAT IP collisions (RFC 6598).
Observed Implementation: Double-checked locking with `sync.RWMutex` around `buckets map[string]*TokenBucket`. Instantiates new `TokenBucket` for unseen keys.
Assessment: PASS
Severity: LOW
Notes: Tenant isolation correctly implemented and thread-safe.

## Finding 4

Location: `internal/backpressure/queue.go:61-70`
Claimed Behavior: Bounded queue sheds excess load fast with non-blocking rejection.
Observed Implementation: `TrySubmit` uses `select` with `case bq.queue <- job` and `default: return ErrQueueFull`. Atomic counters maintain statistics.
Assessment: PASS
Severity: LOW
Notes: Zero-blocking channel drop pattern implemented cleanly.

## Finding 5

Location: `internal/retry/backoff.go:24-62`
Claimed Behavior: Exponential backoff with AWS Jitter strategies (NoJitter, FullJitter, EqualJitter, DecorrelatedJitter) following Marc Brooker's formulation.
Observed Implementation: Implements exact formulas: FullJitter uniformly distributed in $[0, temp]$, EqualJitter in $[temp/2, temp]$, DecorrelatedJitter bounded in $[base, prev \cdot 3]$.
Assessment: PASS
Severity: LOW
Notes: Bounds and formulas strictly align with AWS Architecture Blog specifications.

## Finding 6

Location: `internal/httputil/middleware.go:17-40`
Claimed Behavior: HTTP middleware checks tenant key, returns RFC 6585 status 429 and `Retry-After` header when exhausted.
Observed Implementation: Reads `X-API-Key` (defaults to "anonymous"), queries registry, checks `bucket.Allow()`. If false, writes `Retry-After` header, status code 429, and JSON body.
Assessment: PASS
Severity: LOW
Notes: Standard-compliant header and status code propagation.
