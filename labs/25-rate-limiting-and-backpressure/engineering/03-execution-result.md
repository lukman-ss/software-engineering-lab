# Execution Result

## Build
Command: `go build ./...`
Result: SUCCESS

## Tests
Command: `go test -v ./...`
Result:
```text
?   	labs/25-rate-limiting-and-backpressure/cmd/demo	[no test files]
=== RUN   TestBoundedQueue_RejectionUnderLoad
--- PASS: TestBoundedQueue_RejectionUnderLoad (0.00s)
=== RUN   TestBoundedQueue_ConcurrencySafety
--- PASS: TestBoundedQueue_ConcurrencySafety (0.02s)
PASS
ok  	labs/25-rate-limiting-and-backpressure/internal/backpressure	(cached)
=== RUN   TestRateLimitMiddleware_RFC6585
--- PASS: TestRateLimitMiddleware_RFC6585 (0.00s)
PASS
ok  	labs/25-rate-limiting-and-backpressure/internal/httputil	0.184s
=== RUN   TestTokenBucket_BurstAndRefill
--- PASS: TestTokenBucket_BurstAndRefill (0.20s)
=== RUN   TestLeakyBucket_LeakRate
--- PASS: TestLeakyBucket_LeakRate (0.25s)
=== RUN   TestRegistry_TenantIsolation
--- PASS: TestRegistry_TenantIsolation (0.00s)
=== RUN   TestTokenBucket_ConcurrencyRace
--- PASS: TestTokenBucket_ConcurrencyRace (0.00s)
PASS
ok  	labs/25-rate-limiting-and-backpressure/internal/ratelimit	(cached)
=== RUN   TestComputeBackoff_Bounds
--- PASS: TestComputeBackoff_Bounds (0.00s)
=== RUN   TestDecorrelatedJitter_Bounds
--- PASS: TestDecorrelatedJitter_Bounds (0.00s)
PASS
ok  	labs/25-rate-limiting-and-backpressure/internal/retry	0.500s
```

## Race Detector
Command: `go test -race ./...`
Result:
```text
?   	labs/25-rate-limiting-and-backpressure/cmd/demo	[no test files]
ok  	labs/25-rate-limiting-and-backpressure/internal/backpressure	1.378s
ok  	labs/25-rate-limiting-and-backpressure/internal/httputil	1.433s
ok  	labs/25-rate-limiting-and-backpressure/internal/ratelimit	1.892s
ok  	labs/25-rate-limiting-and-backpressure/internal/retry	1.477s
```

## Demo
Command: `go run ./cmd/demo`
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
Job #5: REJECTED (Backpressure Shedding: backpressure: queue capacity exceeded)
Job #6: REJECTED (Backpressure Shedding: backpressure: queue capacity exceeded)
Stats: Accepted=3, Rejected=3, Processed=1

=== 4. AWS Retry Backoff Strategies (Attempts 0..3) ===
Attempt 0 -> NoJitter: 100ms  | FullJitter: 99ms   | EqualJitter: 92ms  
Attempt 1 -> NoJitter: 200ms  | FullJitter: 109ms  | EqualJitter: 183ms 
Attempt 2 -> NoJitter: 400ms  | FullJitter: 111ms  | EqualJitter: 202ms 
Attempt 3 -> NoJitter: 800ms  | FullJitter: 125ms  | EqualJitter: 709ms 
```

## Final Engineering Status
READY_FOR_ENGINEERING_AUDIT
