# Execution Result

## Build
Command: `go build ./...`
Result: SUCCESS

## Tests
Command: `go test -count=1 -v ./...`
Result:
```text
?   	labs/25-rate-limiting-and-backpressure/cmd/demo	[no test files]
=== RUN   TestBoundedQueue_RejectionUnderLoad
--- PASS: TestBoundedQueue_RejectionUnderLoad (0.00s)
=== RUN   TestBoundedQueue_ConcurrencySafety
--- PASS: TestBoundedQueue_ConcurrencySafety (0.01s)
=== RUN   TestBoundedQueue_SubmitAfterStop
--- PASS: TestBoundedQueue_SubmitAfterStop (0.00s)
=== RUN   TestBoundedQueue_ConcurrentStopAndSubmit
--- PASS: TestBoundedQueue_ConcurrentStopAndSubmit (0.01s)
PASS
ok  	labs/25-rate-limiting-and-backpressure/internal/backpressure	0.323s
=== RUN   TestRateLimitMiddleware_RFC6585
--- PASS: TestRateLimitMiddleware_RFC6585 (0.00s)
=== RUN   TestRateLimitMiddleware_AnonymousFallbackAndBody
--- PASS: TestRateLimitMiddleware_AnonymousFallbackAndBody (0.00s)
PASS
ok  	labs/25-rate-limiting-and-backpressure/internal/httputil	0.320s
=== RUN   TestTokenBucket_BurstAndRefill
--- PASS: TestTokenBucket_BurstAndRefill (0.20s)
=== RUN   TestLeakyBucket_LeakRate
--- PASS: TestLeakyBucket_LeakRate (0.25s)
=== RUN   TestRegistry_TenantIsolation
--- PASS: TestRegistry_TenantIsolation (0.00s)
=== RUN   TestTokenBucket_RetryAfterSeconds
--- PASS: TestTokenBucket_RetryAfterSeconds (0.60s)
=== RUN   TestTokenBucket_ConcurrencyRace
--- PASS: TestTokenBucket_ConcurrencyRace (0.00s)
=== RUN   TestLeakyBucket_ConcurrencyRace
--- PASS: TestLeakyBucket_ConcurrencyRace (0.00s)
=== RUN   TestRegistry_ConcurrentSameKeyGet
--- PASS: TestRegistry_ConcurrentSameKeyGet (0.00s)
=== RUN   TestTokenBucket_EdgeCases
--- PASS: TestTokenBucket_EdgeCases (0.00s)
PASS
ok  	labs/25-rate-limiting-and-backpressure/internal/ratelimit	1.354s
=== RUN   TestComputeBackoff_Bounds
--- PASS: TestComputeBackoff_Bounds (0.00s)
=== RUN   TestDecorrelatedJitter_Bounds
--- PASS: TestDecorrelatedJitter_Bounds (0.00s)
=== RUN   TestComputeBackoff_UnknownStrategy
--- PASS: TestComputeBackoff_UnknownStrategy (0.00s)
PASS
ok  	labs/25-rate-limiting-and-backpressure/internal/retry	0.308s
```

## Race Detector
Command: `go test -race -count=1 ./...`
Result:
```text
?   	labs/25-rate-limiting-and-backpressure/cmd/demo	[no test files]
ok  	labs/25-rate-limiting-and-backpressure/internal/backpressure	1.366s
ok  	labs/25-rate-limiting-and-backpressure/internal/httputil	1.341s
ok  	labs/25-rate-limiting-and-backpressure/internal/ratelimit	2.380s
ok  	labs/25-rate-limiting-and-backpressure/internal/retry	1.322s
```

## Demo (Sample Execution)
Command: `go run ./cmd/demo`
Note: Section 3 and Section 4 produce timing- and randomized-dependent outputs across runs by design.
Result:
```text
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
Stats: Accepted=4, Rejected=2, Processed=1

=== 4. AWS Retry Backoff Strategies (Attempts 0..3) ===
Attempt 0 -> NoJitter: 100ms  | FullJitter: 25ms   | EqualJitter: 82ms  
Attempt 1 -> NoJitter: 200ms  | FullJitter: 31ms   | EqualJitter: 107ms 
Attempt 2 -> NoJitter: 400ms  | FullJitter: 311ms  | EqualJitter: 231ms 
Attempt 3 -> NoJitter: 800ms  | FullJitter: 788ms  | EqualJitter: 562ms 
```

## Final Engineering Status
READY_FOR_ENGINEERING_AUDIT
