# Test Audit

## Finding 1

Location: internal/ratelimit/bucket_test.go:9-30 (TestTokenBucket_BurstAndRefill)
Covered: Happy path burst to capacity, exhaustion, timed refill.
Missing: Over-consumption (AllowN >1), refillRate=0 edge, large n > capacity.
Assessment: PASS
Notes: Core behavior verified; fractional tokens exercised.

## Finding 2

Location: internal/ratelimit/bucket_test.go:32-49 (TestLeakyBucket_LeakRate)
Covered: Two allows, third blocked, time-based leak enabling third.
Missing: LeakRate=0 (never leaks), fractional water tracking not asserted.
Assessment: PASS
Notes: Basic leak validated.

## Finding 3

Location: internal/ratelimit/bucket_test.go:51-67 (TestRegistry_TenantIsolation)
Covered: Two tenants, isolated quotas.
Missing: More than two tenants, registry update after creation, zero-capacity bucket.
Assessment: PASS
Notes: Isolation proven.

## Finding 4

Location: internal/ratelimit/bucket_test.go:69-81 (TestTokenBucket_RetryAfterSeconds)
Covered: Empty bucket → positive retry-after; refill → zero.
Missing: n>1 token requests, refillRate=0 edge (division by zero), large n.
Assessment: WARNING
Severity: MEDIUM
Notes: Only n=1.0 tested; function signature allows fractional n.

## Finding 5

Location: internal/ratelimit/bucket_test.go:83-98 (TestTokenBucket_ConcurrencyRace)
Covered: No race under 500 Allow calls.
Missing: Behavioral assertion (total allowed vs expected); mixed AllowN; stress with stops/starts.
Assessment: WARNING
Severity: MEDIUM
Notes: Race pass but weak validity; no post-condition checked.

## Finding 6

Location: internal/backpressure/queue_test.go:10-45 (TestBoundedQueue_RejectionUnderLoad)
Covered: Queue fills capacity+1 → immediate ErrQueueFull on next TrySubmit.
Missing: Mixed submit/consumer rates, worker count >1, stress with variable job times.
Assessment: PASS
Notes: Classic bounded-buffer rejection proven.

## Finding 7

Location: internal/backpressure/queue_test.go:47-67 (TestBoundedQueue_ConcurrencySafety)
Covered: 30 concurrent TrySubmit with 4 workers, capacity 20; accepted+rejected=30.
Missing: Latency or order assertions; worker liveness under sustained load.
Assessment: PASS
Notes: Counters accurate under race.

## Finding 8

Location: internal/httputil/middleware_test.go:11-43 (TestRateLimitMiddleware_RFC6585)
Covered: First 200, second 429 with Retry-After header.
Missing: Multiple tenants, different keys, Retry-After calculation validation, success path passthrough with body.
Assessment: WARNING
Severity: MEDIUM
Notes: Single-tenant test; header value not checked for correctness.

## Finding 9

Location: internal/retry/backoff_test.go:8-33 (TestComputeBackoff_Bounds)
Covered: Full/Equal/No jitter within [base*2^attempt, cap] for attempts 0-9.
Missing: Distribution shape, actual randomness, prevSleep usage (NotJitter/EqualJitter ignore prevSleep).
Assessment: PASS
Notes: Bounds correct per AWS formulas.

## Finding 10

Location: internal/retry/backoff_test.go:35-48 (TestDecorrelatedJitter_Bounds)
Covered: DecorrelatedJitter within [base, cap] with chained prevSleep.
Missing: Lower bound > base? (formula allows base exactly); distribution over interval.
Assessment: PASS
Notes: Bound adherence verified.

## Overall Test Coverage Assessment
- Happy path: covered for all major functions.
- Failure paths: partial (queue full, rate limit exceeded, bucket empty).
- Edge cases: sparse (zero rates, huge bursts, clock skew, Stop vs TrySubmit races).
- Concurrency: smoke tests present but weak assertions; race detector clean.
- Negative cases: missing (e.g., invalid inputs, panic conditions).

Testing proves mechanisms work but does not stress all claimed properties under duress.