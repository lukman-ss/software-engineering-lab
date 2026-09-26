# Docs vs Code Audit

## Overview

This audit compares README.md, research/05-report.md, engineering design notes, and the actual code implementation and test execution.

---

## Finding 1: Go Version

Location: go.mod vs README.md

Claimed: README states Go 1.22+

Observed: go.mod declares `go 1.22`

Assessment: PASS
Severity: LOW
Notes: Version matches.

---

## Finding 2: Token Bucket Documentation vs Implementation

Location: README.md:41 vs internal/ratelimit/bucket.go

Claimed: TokenBucket allows bursts up to capacity B, refilling at rate R.

Observed: Implementation matches exactly. `AllowN` adds tokens at `elapsed * refillRate`, caps at capacity, decrements on success.

Assessment: PASS
Severity: LOW

---

## Finding 3: Leaky Bucket Documentation vs Implementation

Location: README.md:42 vs internal/ratelimit/bucket.go:81-115

Claimed: LeakyBucket enforces constant drain rate R, rejecting bursts when water reaches capacity.

Observed: Implementation drains at `elapsed * leakRate`, checks `water + 1.0 <= capacity`.

Assessment: PASS
Severity: LOW

---

## Finding 4: Per-Tenant Registry Documentation vs Implementation

Location: README.md:43 vs internal/ratelimit/registry.go

Claimed: Per-tenant registry key isolation guarding against CGNAT IP collisions (RFC 6598).

Observed: Registry creates separate TokenBucket per tenant key. Uses API key header, not IP.

Assessment: PASS
Severity: LOW

---

## Finding 5: Bounded Queue Documentation vs Implementation

Location: README.md:45 vs internal/backpressure/queue.go:61-70

Claimed: Non-blocking submission (TrySubmit) returns fast ErrQueueFull when buffer capacity is reached.

Observed: Implementation uses select-default pattern for non-blocking submission.

Assessment: PASS
Severity: LOW

---

## Finding 6: HTTP 429 Middleware Documentation vs Implementation

Location: README.md:49 vs internal/httputil/middleware.go

Claimed: Enforces rate limits returning standard RFC 6585 429 Too Many Requests response with Retry-After header.

Observed: Middleware returns http.StatusTooManyRequests (429), sets Retry-After header, writes JSON body.

Assessment: PASS
Severity: LOW

---

## Finding 7: Retry Strategies Documentation vs Implementation

Location: README.md:47 vs internal/retry/backoff.go

Claimed: Full Jitter, Equal Jitter, No Jitter, and Decorrelated Jitter algorithms matching AWS Architecture specifications (Marc Brooker).

Observed: All four strategies implemented in ComputeBackoff:
- NoJitter: min(cap, base * 2^attempt)
- FullJitter: random(0, min(cap, base * 2^attempt))
- EqualJitter: temp/2 + random(0, temp/2)
- DecorrelatedJitter: min(cap, random(base, prevSleep * 3))

Assessment: PASS
Severity: LOW

---

## Finding 8: Execution Result vs Actual Test Run

Location: engineering/03-execution-result.md vs actual go test run

Claimed: go test -v ./... passes with 10 tests

Observed: Actual run passes with 10 tests:
- backpressure: 2 tests (RejectionUnderLoad, ConcurrencySafety)
- httputil: 1 test (RFC6585)
- ratelimit: 5 tests (BurstAndRefill, LeakyBucket_LeakRate, TenantIsolation, RetryAfterSeconds, ConcurrencyRace)
- retry: 2 tests (Bounds, DecorrelatedJitter_Bounds)
- Total: 10 tests, all PASS

Assessment: PASS
Severity: LOW
Notes: Test count and names match exactly.

---

## Finding 9: Race Detector Result vs Actual

Location: engineering/03-execution-result.md vs actual go test -race run

Claimed: go test -race ./... passes on all packages

Observed: Actual run passes with no race conditions detected.

Assessment: PASS
Severity: LOW

---

## Finding 10: Demo Output vs Actual Execution

Location: engineering/03-execution-result.md:55-83 vs actual go run ./cmd/demo

Claimed: Demo should produce specific output showing token bucket burst, leaky bucket smoothing, queue rejection, and retry backoff strategies.

Observed: Actual demo output matches expected patterns:
- Token bucket: 3 allowed, 2 rejected, then 1 refilled after 300ms → ✅
- Leaky bucket: 3 allowed, 2 rejected (capacity 3) → ✅
- Bounded queue: 3 accepted, 3 rejected with capacity 3 → ✅
- Retry strategies: Deterministic for NoJitter, random for FullJitter/EqualJitter → ✅

Assessment: PASS
Severity: LOW
Notes: Exact numeric values differ (jitter is random), but behavior is consistent. The demo output in 03-execution-result.md shows representative output. The NoJitter, FullJitter, and EqualJitter values will vary between runs.

---

## Finding 11: Design Document vs Implementation

Location: engineering/01-design.md vs actual code

Claimed Components:
1. ratelimit: TokenBucket, LeakyBucket, Registry → All implemented ✅
2. backpressure: BoundedQueue → Implemented ✅
3. retry: NoJitter, FullJitter, EqualJitter, DecorrelatedJitter → All implemented ✅
4. http middleware: RFC 6585 429 → Implemented ✅
5. cmd/demo → Implemented ✅

Assessment: PASS
Severity: LOW
Notes: All claimed components are implemented.

---

## Finding 12: "What Is Not Demonstrated"

Location: engineering/02-implementation-notes.md:44-46

Claimed Not Demonstrated:
- Distributed rate limiting across multi-region clusters
- Dynamic queue capacity autoscaling

Observed: These features are not present in code. Correctly documented as not demonstrated.

Assessment: PASS
Severity: LOW

---

## Finding 13: Design Decision — Monotonic Clock

Location: engineering/01-design.md:93 vs internal/ratelimit/bucket.go

Claimed: "Monotonic clock comparisons via Go's time.Now() for rate limit refills."

Observed: Implementation uses `time.Now().Sub(lastRefill)` which relies on `time.Now()`. Go's `time.Now()` returns time since January 1, year 1 in a monotonic clock reference. When using `.Sub()`, the result is based on the monotonic clock reading, which is immune to wall clock adjustments. However, this is a property of how Go represents time, not an explicit monotonic clock API.

Assessment: WARNING
Severity: MEDIUM
Notes: The claim says "monotonic clock" but the implementation relies on `time.Now()` which, while Go's implementation does include monotonic readings internally, is not explicitly using `time.Now().Sub()` in a way that is guaranteed to be monotonic-only. In Go, when you store a `time.Time` value and later call `.Sub()`, the monotonic component is preserved IF the times were created close together in time. However, if times are persisted and restored (e.g., after process restart), the monotonic component is lost. For this demo, this is acceptable but the documentation overstates the guarantee. The implementation is correct for in-memory use.

---

## DOC_CODE_MISMATCH Summary

| Finding | Type | Severity | Status |
|---------|------|----------|--------|
| Go version | DOC_CODE_MISMATCH | LOW | Resolved |
| Token bucket | DOC_CODE_MISMATCH | LOW | None |
| Leaky bucket | DOC_CODE_MISMATCH | LOW | None |
| Registry | DOC_CODE_MISMATCH | LOW | None |
| Bounded queue | DOC_CODE_MISMATCH | LOW | None |
| HTTP 429 | DOC_CODE_MISMATCH | LOW | None |
| Retry strategies | DOC_CODE_MISMATCH | LOW | None |
| Test results | DOC_CODE_MISMATCH | LOW | None |
| Race detector | DOC_CODE_MISMATCH | LOW | None |
| Demo output | DOC_CODE_MISMATCH | LOW | Representative |
| Design claim | DOC_CODE_MISMATCH | MEDIUM | Minor inaccuracy in "monotonic clock" |

## Conclusion

The README, engineering design notes, and execution results documents all align with the actual implementation. The only minor discrepancy is the overstatement of "monotonic clock" usage, which is technically true in Go's implementation details but not using a dedicated monotonic clock API like `time.Ticker`.