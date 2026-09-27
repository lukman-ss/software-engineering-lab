## Finding 1

Location: internal/ratelimit/bucket.go:29-46
Claimed Behavior: TokenBucket enforces burst capacity then refill rate with fractional precision.
Observed Implementation: Mutex-guarded tokens + time-based refill (`elapsed*refillRate`), clamped to capacity. `AllowN(n)` deducts if sufficient; else returns false.
Assessment: PASS
Severity: LOW
Notes: `time.Now()` inside lock — low contention, no correctness issue.

## Finding 2

Location: internal/ratelimit/bucket.go:55-79
Claimed Behavior: HTTP `Retry-After` seconds calculation.
Observed Implementation: `RetryAfterSeconds(n)` computes refill needed, rounds up (`ceil`) to int. Verified by `TestTokenBucket_RetryAfterSeconds`.
Assessment: PASS
Severity: LOW
Notes: Uses separate lock copy without mutating state — correct.

## Finding 3

Location: internal/ratelimit/bucket.go:98-115
Claimed Behavior: LeakyBucket smooths traffic to leak rate `R`, rejecting burst overflow.
Observed Implementation: Mutex-guarded water level; leak applied per `time.Now()` delta, floor clamped at 0; admits if `water+1 <= capacity`.
Assessment: PASS
Severity: LOW
Notes: Boundary allows refill bursts up to capacity — consistent with design.

## Finding 4

Location: internal/ratelimit/registry.go:21-37
Claimed Behavior: Per-tenant key isolation to avoid CGNAT IP collision.
Observed Implementation: RWMutex double-checked locking; buckets map per tenantKey.
Assessment: PASS
Severity: LOW
Notes: Correct concurrent registration; race test covers.

## Finding 5

Location: internal/backpressure/queue.go:66-81
Claimed Behavior: Bounded queue fast-rejects without blocking caller.
Observed Implementation: `select { case queue<-job: ...; default: rejected++ }` — true non-blocking drop.
Assessment: PASS
Severity: LOW
Notes: StoppedCtx + closed channel handling correct; `Stop()` idempotent.

## Finding 6

Location: internal/backpressure/queue.go:48-61
Claimed Behavior: Worker pool consumes jobs.
Observed Implementation: N workers select ctx/shutdown vs job channel; `processed++` per job.
Assessment: PASS
Severity: LOW
Notes: Job errors ignored (`_ = job`) by design — no error propagation; acceptable for demo.

## Finding 7

Location: internal/retry/backoff.go:24-62
Claimed Behavior: AWS Full/Equal/No/Decorrelated jitter per Marc Brooker spec.
Observed Implementation: Formulas match spec: Full `rand*temp`; Equal `temp/2 + rand*temp/2`; Decorrelated `min(cap, rand(base, prev*3))`; Cap via `math.Min`.
Assessment: PASS
Severity: LOW
Notes: Uses global `math/rand` without seed — random output per run; fine for bounds.

## Finding 8

Location: internal/httputil/middleware.go:17-40
Claimed Behavior: RFC 6585 429 + Retry-After + JSON payload.
Observed Implementation: Checks `X-API-Key` (fallback anonymous), `bucket.Allow()`; on exhaust sets Content-Type json, Retry-After header, 429, encodes `{"error":"rate_limit_exceeded","retry_after":n}`.
Assessment: PASS
Severity: LOW
Notes: Minor: body ignores `json.Encode` error; irrelevant for lab.
