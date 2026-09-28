# Master Draft Audit

## File: content/02-master-draft.md

### Accuracy Assessment: PASS WITH WARNINGS

### Overall:
Main article accurately covers the topic with claims traceable to research, engineering design, and actual code implementation. All technical facts verified. No hallucinated claims found.

### Detailed Findings:

#### Verified Claims — PASS
- **Problem statement**: Cascading failures from unbounded queues without rate limiting/backpressure — verified against research findings 2, 5, 10
- **Token Bucket vs Leaky Bucket**: Algorithm descriptions match implementation in `internal/ratelimit/bucket.go` — verified
- **Failure Scenario math**: 5,000,000 / 2,000 = 2,500s ≈ 41m40s — verified against Little's Law and research F5
- **AWS Full Jitter formula**: `delay = random(0,1) × min(cap, base × 2^attempt)` — matches code in `internal/retry/backoff.go:36-41` and research F3
- **NGINX leaky bucket**: Uses leaky bucket as meter via `limit_req_module`, default 503 — verified against research F1
- **RFC 6585 429 MUST NOT be cached**: Matches RFC source requirement
- **Stripe's four limiters**: Matches research F6 — verified
- **Google SRE CPU quotas and criticality levels**: Matches research F8 — verified
- **Retry budgets (per-request 3, per-client 10%)**: Matches research F9 — verified

#### Code Walkthrough Accuracy — PASS
- Token bucket demo (cap 3, refill 5/s): First 3 allowed, 4th/5th rejected, 300ms pause yields ~0.5 tokens — matches actual `cmd/demo/main.go:15-22` and code
- Leaky bucket demo (cap 3, leak 10/s): Fills to 3, rejects when full — matches code
- Bounded queue demo (cap 3, 1 worker): accepted=4, rejected=2, processed=1 — matches code behavior
- HTTP middleware: Returns 429 with `Retry-After` header and JSON body `{error: "rate_limit_exceeded", retry_after: N}` — matches `internal/httputil/middleware.go:17-40`
- Implementation details (stdlib, monotonic time, floating-point tokens, select-default pattern, mutex vs CAS) — all match code

#### Verified Claims — PASS
- **Test descriptions**: All test names listed match actual test functions in test files
  - `TestTokenBucket_BurstAndRefill` (bucket_test.go:9) ✓
  - `TestLeakyBucket_LeakRate` (bucket_test.go:32) ✓
  - `TestRegistry_TenantIsolation` (bucket_test.go:51) ✓
  - `TestTokenBucket_RetryAfterSeconds` (bucket_test.go:69) ✓
  - `TestTokenBucket_ConcurrencyRace` (bucket_test.go:83) ✓
  - `TestLeakyBucket_ConcurrencyRace` (bucket_test.go:100) ✓
  - `TestRegistry_ConcurrentSameKeyGet` (bucket_test.go:117) ✓
  - `TestTokenBucket_EdgeCases` (bucket_test.go:140) ✓
  - All BoundedQueue tests (queue_test.go) ✓
  - All Backoff tests (backoff_test.go) ✓
  - All Middleware tests (middleware_test.go) ✓

### Issues Found — WARNINGS:

#### WARNING 1: Architecture Diagram — 503 vs ErrQueueFull (MEDIUM)
- **Location**: Lines 42-58, Diagram section
- **Content**: Architecture diagram shows "Queue Full? → 503 Overloaded"
- **Actual Code**: `BoundedQueue.TrySubmit()` returns `ErrQueueFull` (an application-layer error), NOT HTTP 503. The HTTP middleware returns 429, not 503.
- **Impact**: The diagram suggests an HTTP 503 response that is not implemented. The actual code returns `ErrQueueFull` to the caller, who must decide how to map it. The 503 is a design intent, not an implemented behavior.
- **Severity**: LOW — conceptual diagram uses design intent rather than exact implementation status
- **Note**: This originates from the engineering design doc `engineering/01-design.md:49`

#### WARNING 2: Backlog Growth Formula Notation (LOW)
- **Location**: Line 107
- **Content**: "Backlog growth follows the formula: `(arrival_rate - processing_rate) × time`"
- **Note**: This is consistent with Little's Law but is an approximation (assumes constant rates). Acceptable for educational context.

#### WARNING 3: Sources Section — Inline Source List (LOW)
- **Location**: Lines 152-167
- **Content**: Uses Wikipedia URLs directly rather than the [number]-style reference system used in the Sources section of03-master-draft.md
- **Note**: Functionally correct but inconsistent citation style
- **Severity**: LOW

### Test Coverage Descriptions — VERIFIED
- TokenBucket tests: burst capacity, refill timing, depletion, 50-goroutine concurrency — all verified against actual test file
- LeakyBucket tests: water tracking, capacity rejection, concurrency — verified
- Registry tests: tenant isolation and single-pointer guarantee for same key — verified
- BoundedQueue tests: rejection under load, concurrency stress, stop cleanup — verified
- Retry tests: bounds checking for Full/Equal/ Decorrelated, unknown strategy fallback — verified
- HTTP middleware tests: 429 status, Retry-After header, JSON body, anonymous fallback — verified

### Verdict: PASS WITH WARNINGS