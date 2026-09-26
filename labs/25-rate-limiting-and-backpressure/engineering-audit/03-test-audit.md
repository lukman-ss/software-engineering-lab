# Test Audit

## Coverage Summary

- Happy Path: Covered (Token bucket allow, leaky bucket allow, queue submit accept, HTTP 200 OK).
- Failure Path / Rejection: Covered (Token bucket exhaustion, leaky bucket capacity rejection, queue `ErrQueueFull`, HTTP 429 response).
- Edge Cases: Covered (Refill sleep recovery, RetryAfterSeconds estimation, tenant key isolation).
- Concurrency & Race Safety: Covered (`TestTokenBucket_ConcurrencyRace`, `TestBoundedQueue_ConcurrencySafety`, verified with `go test -count=1 -race ./...`).

## Execution Verification

### Command 1: `go test -v ./...`
```text
=== RUN   TestBoundedQueue_RejectionUnderLoad
--- PASS: TestBoundedQueue_RejectionUnderLoad (0.00s)
=== RUN   TestBoundedQueue_ConcurrencySafety
--- PASS: TestBoundedQueue_ConcurrencySafety (0.02s)
PASS
ok  	labs/25-rate-limiting-and-backpressure/internal/backpressure	0.031s
=== RUN   TestRateLimitMiddleware_RFC6585
--- PASS: TestRateLimitMiddleware_RFC6585 (0.00s)
PASS
ok  	labs/25-rate-limiting-and-backpressure/internal/httputil	0.015s
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
PASS
ok  	labs/25-rate-limiting-and-backpressure/internal/ratelimit	1.056s
=== RUN   TestComputeBackoff_Bounds
--- PASS: TestComputeBackoff_Bounds (0.00s)
=== RUN   TestDecorrelatedJitter_Bounds
--- PASS: TestDecorrelatedJitter_Bounds (0.00s)
PASS
ok  	labs/25-rate-limiting-and-backpressure/internal/retry	0.012s
```

### Command 2: `go test -count=1 -race ./...`
```text
?   	labs/25-rate-limiting-and-backpressure/cmd/demo	[no test files]
ok  	labs/25-rate-limiting-and-backpressure/internal/backpressure	1.148s
ok  	labs/25-rate-limiting-and-backpressure/internal/httputil	1.098s
ok  	labs/25-rate-limiting-and-backpressure/internal/ratelimit	2.131s
ok  	labs/25-rate-limiting-and-backpressure/internal/retry	1.074s
```

### Command 3: `go run ./cmd/demo`
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
Job #5: REJECTED (Backpressure Shedding: backpressure: queue capacity exceeded)
Job #6: REJECTED (Backpressure Shedding: backpressure: queue capacity exceeded)
Stats: Accepted=3, Rejected=3, Processed=1

=== 4. AWS Retry Backoff Strategies (Attempts 0..3) ===
Attempt 0 -> NoJitter: 100ms  | FullJitter: 43ms   | EqualJitter: 72ms  
Attempt 1 -> NoJitter: 200ms  | FullJitter: 178ms  | EqualJitter: 172ms 
Attempt 2 -> NoJitter: 400ms  | FullJitter: 93ms   | EqualJitter: 231ms 
Attempt 3 -> NoJitter: 800ms  | FullJitter: 149ms  | EqualJitter: 642ms 
```

Test Audit Assessment: PASS
No weak or fake tests found. All claimed behaviors are validated programmatically and race-detector verified.
