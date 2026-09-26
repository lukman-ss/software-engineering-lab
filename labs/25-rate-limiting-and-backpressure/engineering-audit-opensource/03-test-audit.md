# Engineering Test Audit

Target Lab: labs/25-rate-limiting-and-backpressure
Test files examined:
- internal/ratelimit/bucket_test.go
- internal/backpressure/queue_test.go
- internal/httputil/middleware_test.go
- internal/retry/backoff_test.go
- [no test files for cmd/demo]

Test commands executed:
- `go test -v -count=1 ./...`   (recorded in 03-test-audit.md)
- `go test -race ./...`         (race detector passed on shipped tests)

## Coverage Analysis

### TokenBucket / LeakyBucket (bucket_test.go)
Tests (5):
1. `TestTokenBucket_BurstAndRefill` — happy path: initial burst, refill after sleep. ✅
2. `TestLeakyBucket_LeakRate` — happy path: two allows, leak, then allow. ✅
3. `TestRegistry_TenantIsolation` — multi-tenant isolation (registry). ✅
4. `TestTokenBucket_RetryAfterSeconds` — retry-after >0 when empty, ==0 after refill ✅
5. `TestTokenBucket_ConcurrencyRace` — 50 goroutines × 10 ops each, no assert on counts ✅
- **Missing**: exact refill accounting (e.g., tokens after fractional time), stress on rate overflow, refillRate==0 edge, boundary exactness (AllowN > capacity), concurrent Allow + Tokens/Ret
- **Concurrency invariant**: none (only checks no panic/race). MEDIUM gap.

### BoundedQueue (queue_test.go)
Tests (3):
1. `TestBoundedQueue_RejectionUnderLoad` — fill buffer, next submission ErrQueueFull ✅
2. `TestBoundedQueue_ConcurrencySafety` — 30 goroutines submit; assert accepted+rejected==30 ✅ (strong invariant!)
3. `TestBoundedQueue_SubmitAfterStop` — Stop then TrySubmit -> ErrQueueStopped; double-Stop idempotent ✅
- **Missing**: Test for **concurrent submit+Stop** (the actual panic path). The exact defect that ships with the code (see 02-code-audit §Finding 5) is untested. HIGH gap.
- **Missing**: invariants on accepted/rejected/processed under race (current test only checks sum). Could assert processed <= accepted, etc.
- **Missing**: Test worker progress (jobs actually run), ctx propagation.

### Middleware (middleware_test.go)
Test (1):
1. `TestRateLimitMiddleware_RFC6585` — 1-token registry: 1st request 200 OK, 2nd 429 + Retry-After header ✅
- **Missing**: anonymous fallback (`X-API-Key` missing/empty)
- **Missing**: body content validation (exact JSON structure, Retry-After numeric)
- **Missing**: multiple tenants (two keys independent)
- **Missing**: 200 path body passthrough
- **Missing**: error handling (malformed JSON?)
Coverage: happy + one failure path. MEDIUM gap.

### Retry (backoff_test.go)
Tests (2):
1. `TestComputeBackoff_Bounds` — attempt 0..9: Full/Equal/NoJitter within [0, cap]; NoJitter also >= base. ✅
2. `TestDecorrelatedJitter_Bounds` — DecorrelatedJitter within [base, cap] for 5 iterations ✅
- **Missing**: exact formula checks (e.g., NoJitter == min(cap, base*2^attempt) for many attempts)
- **Missing**: EqualJitter exact distribution check (hard but possible via stats over many samples)
- **Missing**: DecorrelatedJitter clamping (prevSleep < base -> use base)
- **Missing**: integration with backoff loop (not just unit function)
Coverage: bounds only; no exact-value validation. LOW/MEDIUM gap.

## Quality of Concurrency Tests

- **TokenBucket concurrency**: 50 goroutine hammer — no invariant checked, only no-race/panic. **WEAK**.
- **Backpressure concurrency**: 30 goroutine hammer + invariant accepted+rejected==30 — **STRONG** for the submission path. However, the **shutdown race** (submit+Stop) has **zero** coverage (see gap above). This is the sole defect that escaped testing.
- **Middleware**: no concurrency test (only sequential HTTP handler). LOW.
- **Retry**: pure function; no shared state; concurrency irrelevant.

## Test Reliability

- No flaky sleeps (sleeps are long enough: 200ms, 250ms, 600ms). Design claim "deterministic simulated time helpers" is FALSE (tests use real time, see 04-docs-vs-code.md). MEDIUM documentation gap.
- Tests are fast (<2s) and deterministic (given wall clock stability). PASS.

## Summary

| Test Suite | Coverage | Quality | Primary Gaps |
|---|---|---|---|
| bucket_test.go | Happy path + isolation + retry-after + race | WEAK concurrency invariant | fractional accounting, refillRate==0 edge, concurrent ops invariant |
| queue_test.go | Boundary rejection, concurrency sum invariant, stop-after | **CRITICAL GAP** — no concurrent submit+Stop test (missed the ship-stopping panic) | submit+Stop race, worker progress, processed<=accepted invariant |
| middleware_test.go | 200/429 single tenant | Missing anonymous, body, multi-tenant, value checks | HTTP negative coverage thin |
| backoff_test.go | Bounds for all strategies | No exact-value checks (only bounds) | exact formula validation (NoJitter, EqualJitter, Decorrelated) |