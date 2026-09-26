# Test Audit

Target: labs/25-rate-limiting-and-backpressure
Executed: go test ./..., go test -race ./...
Result: PASS (all packages), race-clean

## Coverage By Package

### ratelimit (bucket_test.go)
- TestTokenBucket_BurstAndRefill: burst exhaustion + refill timing. PASS. Timing-sensitive (200ms sleep).
- TestLeakyBucket_LeakRate: capacity admit + drain. PASS. Timing-sensitive (250ms).
- TestRegistry_TenantIsolation: per-tenant isolation. PASS.
- TestTokenBucket_RetryAfterSeconds: RetryAfter>0 empty, ==0 after refill. PASS.
- TestTokenBucket_ConcurrencyRace: 50 goroutines x10 against shared bucket. PASS.

Assessment: Covers happy path, burst, refill timing, isolation. Missing: refill rate accuracy over longer interval, zero-token boundary, negative/invalid constructor input, capacity rounding edge, RetryAfter with refillRate<=0 (Finding 2 div-by-zero). Race test present.

### backpressure (queue_test.go)
- TestBoundedQueue_RejectionUnderLoad: 1 worker blocking via channel, fill capacity+1 -> ErrQueueFull. PASS.
- TestBoundedQueue_ConcurrencySafety: 30 goroutines TrySubmit, assert accepted+rejected==30. PASS. Race-clean.

Assessment: Proves fast-reject under capacity and concurrency accounting. Missing: worker pool processing correctness (processed count asserted only in demo, not tests), Stop() idempotency, submit-after-close, capacity<=0 input, job error propagation (deliberately swallowed).

### httputil (middleware_test.go)
- TestRateLimitMiddleware_RFC6585: 1st -> 200, 2nd -> 429 + Retry-After present. PASS.

Assessment: Proves 429 + Retry-After + tenant gate. Missing: Retry-After value correctness (only presence checked), anonymous fallback path, JSON body schema assertion, multi-tenant isolation via middleware, rate-limit-reset-after-refill behavior.

### retry (backoff_test.go)
- TestComputeBackoff_Bounds: Full/Equal/No bounds vs cap. PASS.
- TestDecorrelatedJitter_Bounds: Decorrelated bounds across prev-sleep sequence. PASS.

Assessment: Bounds verified. MISSING edge cases: EqualJitter lower-bound = temp/2 not asserted (only >0,<cap), Decorrelated minimum = base not strict (test asserts base..cap), NoJitter exactness (asserts range not formula), FullJitter lower bound = 0 not asserted, negative/zero attempt/attempt<0 handling, Cap==Base floor, prevSleep monotonic growth not verified.

## Coverage Summary

| Package | happy/failure/edge | concurrency | race | negative cases |
| --- | --- | --- | --- | --- |
| ratelimit | happy+timing+isolation; no edge | yes | yes | no invalid input tests |
| backpressure | rejection+correctness; no lifecycle | yes | yes | no stop/edge tests |
| httputil | 429+presence; no value/schema | no | n/a | no negative cases |
| retry | bounds; no formula-exact/edge | no | n/a | weak edge tests |

Overall: Tests prove core behavior. Gaps are edge/invalid-input and lifecycle, all non-blocking for the demonstrated claim set.
