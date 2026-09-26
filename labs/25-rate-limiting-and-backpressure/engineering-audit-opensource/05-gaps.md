# Engineering Audit Gaps

## Gap 1: BoundedQueue Submission After Stop()

Type: UNHANDLED_ERROR
Severity: MEDIUM
Location: internal/backpressure/queue.go:61-69

Description: The `TrySubmit` method does not guard against submissions after `Stop()` has been called. After `Stop()` closes the channel, any call to `TrySubmit` will panic with "send on closed channel". The select-default pattern prevents blocking but not panics on closed channels.

Why It Matters: In production use, if a worker pool is shut down and then a late-arriving request tries to submit, the application would crash. In tests/demo this doesn't happen because `Stop()` is always called via `defer` and no submissions occur after.

Recommendation: Add a `stopped` atomic flag or use a `sync.Once` for Stop, and check in TrySubmit before channel send. Or use a separate goroutine to handle submissions with context cancellation.

---

## Gap 2: Registry Memory Growth

Type: MISSING_EDGE_CASE
Severity: LOW
Location: internal/ratelimit/registry.go:6-37

Description: The Registry never cleans up expired or unused tenant buckets. In a long-running system with high tenant churn, this could lead to unbounded memory growth.

Why It Matters: This is a known limitation documented in implementation notes as "in-memory rate limiting state without distributed synchronization." For a production system, TTL-based cleanup would be needed. For this lab, it is acceptable.

Recommendation: Not required for lab scope. Document as a known limitation (already done).

---

## Gap 3: Leaky Bucket Has No AllowN

Type: IMPLEMENTATION_OVERCLAIM
Severity: LOW
Location: internal/ratelimit/bucket.go:98-115

Description: The design document describes LeakyBucket with Allow() taking unit 1.0. The TokenBucket has AllowN for arbitrary float amounts. This asymmetry is fine for the lab but could be noted.

Why It Matters: Not a bug, but a minor inconsistency in API surface.

Recommendation: Acceptable design choice for the lab scope.

---

## Gap 4: Retry Backoff Race Condition

Type: RACE_CONDITION
Severity: LOW
Location: internal/retry/backoff.go:5-7

Description: `math/rand` uses a global source that is seeded only if the program calls `rand.Seed()`. In Go 1.20+, the global source is auto-seeded. Go 1.22 (used in this lab) auto-seeds, so concurrent calls to `rand.Float64()` are safe due to the global mutex but serialized. For higher throughput, a per-caller `rand.Rand` would be preferable.

Why It Matters: The code is safe (math/rand's Float64 is thread-safe), but performance could be improved with a local source. For this lab the global source is acceptable.

Recommendation: No change needed for lab scope. Auto-seeding in Go 1.22 handles this.

---

## Gap 5: Missing Edge Case Tests for Retry Backoff

Type: MISSING_TEST
Severity: LOW
Location: internal/retry/backoff_test.go

Description: No tests for edge cases:
- Cap < Base (should cap at cap)
- Base = 0
- Attempt = 0 (already partially tested)
- Very large attempt (integer overflow in math.Pow)

Why It Matters: The implementation uses `math.Pow(2, float64(attempt))` which is bounded by float64 precision. For very large attempt values, this produces infinity, which `math.Min(capFloat, expBackoff)` handles by returning capFloat. No overflow panic.

Recommendation: Add edge case tests for robustness.

---

## Gap 6: BoundedQueue Zero Workers

Type: MISSING_EDGE_CASE
Severity: LOW
Location: internal/backpressure/queue.go:26-42

Description: If `NewBoundedQueue` is called with `workers = 0`, no worker goroutines start. Submitted jobs would sit in the channel forever (until Stop is called). If capacity is also 0, jobs would be rejected immediately.

Why It Matters: The demo and tests always use workers >= 1. No panic would occur.

Recommendation: Not required for lab scope. Documented behavior acceptable.

---

## Gap 7: Demo BoundedQueue Stats Timing

Type: MISSING_EDGE_CASE
Severity: LOW
Location: cmd/demo/main.go:49

Description: The demo calls `Stats()` after `time.Sleep(100 * time.Millisecond)`. Each job sleeps 50ms. With 1 worker, by the time Stats is called, the worker may have processed 1 job (from the first 50ms) plus another starting. The execution-result.md shows "Processed=1", which matches.

Why It Matters: The timing is correct and consistent.

Recommendation: No gap. Verified correct.

---

## Gap 8: TokenBucket RetryAfterSeconds Race Window

Type: RACE_CONDITION
Severity: LOW
Location: internal/ratelimit/bucket.go:55-79

Description: `RetryAfterSeconds` reads and computes tokens independently from `AllowN`. If another goroutine calls `AllowN` between the time `RetryAfterSeconds` reads `lastRefill` and returns, the computed retry time may be slightly inaccurate. This is a TOCTOU issue but is acceptable for a rate-limiting hint (not a hard guarantee).

Why It Matters: The Retry-After header is advisory. Small inaccuracies are acceptable per HTTP spec.

Recommendation: Acceptable for lab scope.

---

## Gap Summary

| Type | Count | Severity |
|------|-------|----------|
| UNHANDLED_ERROR | 1 | MEDIUM |
| MISSING_TEST | 1 | LOW |
| MISSING_EDGE_CASE | 3 | LOW |
| RACE_CONDITION | 2 | LOW |
| IMPLEMENTATION_OVERCLAIM | 1 | LOW |

## Overall Risk Assessment

The implementation has no critical design flaws. The single MEDIUM issue (submission after Stop) is a lifecycle guard that doesn't affect the current code paths. All LOW issues are either known/accepted limitations or test coverage gaps that don't affect correctness.

No fabricated results or fake benchmarks found.</think><tool_call>write<arg_key>content</arg_key><arg_value># Engineering Audit Gaps Analysis

## Identified Gaps

### Gap 1: BoundedQueue Submission After Stop()
- Type: UNHANDLED_ERROR
- Severity: MEDIUM
- Location: internal/backpressure/queue.go:61-69
- Description: `TrySubmit` does not guard against submissions after `Stop()` closes the channel. A post-Stop submission would panic with "send on closed channel".
- Mitigation: Not triggered in current demo/test paths. Should add a stopped flag for production use.

### Gap 2: Registry Memory Growth
- Type: MISSING_EDGE_CASE
- Severity: LOW
- Location: internal/ratelimit/registry.go:6-37
- Description: No cleanup of unused/expired tenant buckets. Could cause memory growth in long-running systems with high tenant churn.
- Mitigation: Documented as known limitation. Acceptable for in-memory lab scope.

### Gap 3: Missing Edge Case Tests for Retry Backoff
- Type: MISSING_TEST
- Severity: LOW
- Location: internal/retry/backoff_test.go
- Description: No tests for cap < base, base = 0, large attempt values.
- Mitigation: Implementation handles these via math.Min. Acceptable coverage for lab scope.

### Gap 4: BoundedQueue Zero Workers / Zero Capacity
- Type: MISSING_EDGE_CASE
- Severity: LOW
- Location: internal/backpressure/queue.go:26-42
- Description: NewBoundedQueue with workers=0 produces no workers. Jobs never processed.
- Mitigation: Not exercised in tests. No panic. Documented behavior acceptable.

### Gap 5: TokenBucket RetryAfterSeconds TOCTOU
- Type: RACE_CONDITION
- Severity: LOW
- Location: internal/ratelimit/bucket.go:55-79
- Description: RetryAfterSeconds reads state independently from AllowN. Possible small inaccuracy if concurrent mutation.
- Mitigation: Retry-After is advisory per HTTP spec. Acceptable.

### Gap 6: DecorrelatedJitter PrevSleep < Base Edge Case
- Type: MISSING_TEST
- Severity: LOW
- Location: internal/retry/backoff.go:50-58
- Description: If prevSleep is zero (as in tests with attempt=0), it defaults to base. But this edge case isn't tested.
- Mitigation: Code handles correctly. Minor test gap.

### Gap 7: LeakyBucket No AllowN Method
- Type: IMPLEMENTATION_OVERCLAIM
- Severity: LOW
- Location: internal/ratelimit/bucket.go:98-115 vs 01-design.md
- Description: Design doc does not require AllowN for LeakyBucket. Only Allow exists. This is fine but creates API asymmetry.
- Mitigation: Acceptable design choice.

### Gap 8: LeakyBucket No RetryAfter/Water Method Naming Inconsistency
- Type: MISSING_EDGE_CASE
- Severity: LOW
- Location: internal/ratelimit/bucket.go
- Description: LeakyBucket has Water() but no retry-after calculation method, unlike TokenBucket.
- Mitigation: Leaky bucket is a traffic shaper, not typically used with Retry-After. Acceptable.

### Gap 9: HTTP Middleware No Concurrent Same-Tenant Test
- Type: MISSING_TEST
- Severity: LOW
- Location: internal/httputil/middleware_test.go
- Description: Middleware test only does sequential requests, doesn't stress concurrent same-tenant requests under rate limiting.
- Mitigation: Registry and TokenBucket have concurrency tests. Good enough for lab scope.

### Gap 10: No Test Coverage Measurement
- Type: MISSING_TEST
- Severity: LOW
- Location: All test files
- Description: No `go test -cover` reported or enforced. Coverage percentage unknown.
- Mitigation: Could be added but not required for lab.

## Gap Summary

| Type | Count |
|------|-------|
| UNHANDLED_ERROR | 1 |
| MISSING_TEST | 1 |
| MISSING_EDGE_CASE | 4 |
| RACE_CONDITION | 1 |
| IMPLEMENTATION_OVERCLAIM | 1 |

**MEDIUM severity: 1** — BoundedQueue post-Stop submission panic
**LOW severity: 9** — Various edge cases and coverage gaps