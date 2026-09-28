# Test Audit

Evidence captured live (not quoted from pre-existing file):

```
$ go test ./...
ok  	labs/25-rate-limiting-and-backpressure/internal/backpressure   0.384s
ok  	labs/25-rate-limiting-and-backpressure/internal/httputil       0.406s
ok  	labs/25-rate-limiting-and-backpressure/internal/ratelimit     1.452s
ok  	labs/25-rate-limiting-and-backpressure/internal/retry           0.373s
?   	labs/25-rate-limiting-and-backpressure/cmd/demo                 [no test files]

$ go test -race ./internal/... ./cmd/...
ok  	.../internal/backpressure   1.454s
ok  	.../internal/httputil       1.646s
ok  	.../internal/ratelimit      2.680s
ok  	.../internal/retry          1.620s
?   	labs/25-rate-limiting-and-backpressure/cmd/demo  [no test files]

$ go run ./cmd/demo   (live output below; matches claims)
=== 1. Token Bucket Burst & Rate Limiting ===
Request #1: Allowed=true (Remaining Tokens: 2.0)
Request #2: Allowed=true (Remaining Tokens: 1.0)
Request #3: Allowed=true (Remaining Tokens: 0.0)
Request #4: Allowed=false (Remaining Tokens: 0.0)
Request #5: Allowed=false (Remaining Tokens: 0.0)
After 300ms pause: Allowed=true (Remaining Tokens: 0.5)

=== 2. Leaky Bucket Traffic Smoothing ===
Request #1: Allowed=true (Current Water Level: 1.0)
Request #2: Allowed=true (Current Water Level: 2.0)
Request #3: Allowed=true (Current Water Level: 3.0)
Request #4: Allowed=false (Current Water Level: 3.0)
Request #5: Allowed=false (Current Water Level: 3.0)

=== 3. Bounded Queue Backpressure (Load Shedding) ===
Job #1: ACCEPTED into bounded buffer
Job #2: ACCEPTED into bounded buffer
Job #3: ACCEPTED into bounded buffer
Job #4: REJECTED (Backpressure Shedding: backpressure: queue capacity exceeded)
Job #5: ACCEPTED into bounded buffer
Job #6: REJECTED (Backpressure Shedding: backpressure: queue capacity exceeded)
Stats: Accepted=4, Rejected=2, Processed=2

=== 4. AWS Retry Backoff Strategies (Attempts 0..3) ===
Attempt 0 -> NoJitter: 100ms  | FullJitter: 20ms    | EqualJitter: 92ms
Attempt 1 -> NoJitter: 200ms  | FullJitter: 47ms    | EqualJitter: 103ms
Attempt 2 -> NoJitter: 400ms  | FullJitter: 270ms   | EqualJitter: 325ms
Attempt 3 -> NoJitter: 800ms  | FullJitter: 79ms    | EqualJitter: 459ms
```

## Coverage Analysis

### ratelimit (bucket_test.go — 8 tests)
- TestTokenBucket_BurstAndRefill: burst exhaustion + refill. PASS. Time-based (sleep).
- TestLeakyBucket_LeakRate: capacity + drain. PASS. Time-based.
- TestRegistry_TenantIsolation: isolation between tenants. PASS.
- TestTokenBucket_RetryAfterSeconds: retry-after semantics. PASS. Time-based.
- TestTokenBucket_ConcurrencyRace / TestLeakyBucket_ConcurrencyRace: 500 concurrent calls. PASS under -race.
- TestRegistry_ConcurrentSameKeyGet: pointer-identity under contended same-key Get. PASS.
- TestTokenBucket_EdgeCases: AllowN multi-token. PASS.

Gaps: no negative-capacity / zero-rate test (Finding 8); no refillRate=0 panic guard; no assertion on exact token math correctness beyond monotonicity.

### backpressure (queue_test.go — 4 tests)
- TestBoundedQueue_RejectionUnderLoad: fill buffer -> ErrQueueFull. PASS.
- TestBoundedQueue_ConcurrencySafety: counters consistent. PASS.
- TestBoundedQueue_SubmitAfterStop: idempotent Stop + ErrQueueStopped. PASS.
- TestBoundedQueue_ConcurrentStopAndSubmit: stop concurrent with submits, no crash this run. PASS under -race.

Gaps: does NOT assert the latent send-on-closed-channel panic (Finding 4); no job-panic test (Finding 5); processed count semantics unverified against spec; no memory-exhaustion pressure test.

### httputil (middleware_test.go — 2 tests)
- TestRateLimitMiddleware_RFC6585: 200 then 429 + Retry-After present. PASS.
- TestRateLimitMiddleware_AnonymousFallbackAndBody: anonymous fallback + Content-Type. PASS.

Gaps: no assertion Retry-After value matches body field; no test for unknown strategy in retry (present); no malformed-header handling.

### retry (backoff_test.go — 3 tests)
- TestComputeBackoff_Bounds: [0,cap] and [base,cap] and [0,cap] for 10 attempts. PASS.
- TestDecorrelatedJitter_Bounds: chain prev across 5 iterations. PASS.
- TestComputeBackoff_UnknownStrategy: fallback to temp. PASS.

Gaps: no statistical check of jitter distribution uniformity; no attempt overflow test; no prevSleep < base assertion for DecorrelatedJitter branch.

## Overall Test Quality

Coverage: Happy path covered. Failure path (queue full, rate exhausted, unknown strategy) covered. Concurrency with race detector clean. Negative cases minimal. Recovery/rollback: not applicable (stateless). Edge cases: AllowN, idempotent Stop — partial.

Assessment: suite proves core behavior; race-clean. Weakness is absence of negative-input and job-failure tests for backpressure.
