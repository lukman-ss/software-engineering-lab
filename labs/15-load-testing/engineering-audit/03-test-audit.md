# Test Audit

Target Lab: labs/15-load-testing

## Test Suite Execution Results

### 1. `go test -v ./...`
```text
=== RUN   TestCalculateMetrics
--- PASS: TestCalculateMetrics (0.00s)
=== RUN   TestCalculateMetrics_Empty
--- PASS: TestCalculateMetrics_Empty (0.00s)
=== RUN   TestCalculateMetrics_Invariants
--- PASS: TestCalculateMetrics_Invariants (0.00s)
PASS
ok  	github.com/lukman/software-engineering-lab/labs/15-load-testing/internal/loadtest	0.00s
=== RUN   TestLoadTest_SmokeVsStress
--- PASS: TestLoadTest_SmokeVsStress (1.36s)
=== RUN   TestLoadTest_ErrorCount
--- PASS: TestLoadTest_ErrorCount (0.10s)
=== RUN   TestServer_MethodNotAllowed
--- PASS: TestServer_MethodNotAllowed (0.00s)
=== RUN   TestLoadTest_DialError
--- PASS: TestLoadTest_DialError (0.10s)
=== RUN   TestServer_ContextCanceled
--- PASS: TestServer_ContextCanceled (0.00s)
PASS
ok  	github.com/lukman/software-engineering-lab/labs/15-load-testing/tests	1.56s
```

### 2. `go test -race ./...`
```text
ok  	github.com/lukman/software-engineering-lab/labs/15-load-testing/internal/loadtest	1.150s
ok  	github.com/lukman/software-engineering-lab/labs/15-load-testing/tests	2.592s
```

### 3. `go run ./cmd/demo`
```text
Starting Load Test Demo (Booking Bengkel)
Server simulated DB connections: 5
Simulated DB query duration: 20ms

--- Running Smoke Test (2 VUs) ---
Total Requests:  162
Success:         162
Errors:          0
RPS:             80.98
Average:         24.502425ms
P50:             22.486792ms
P95:             39.610542ms
P99:             48.615041ms

--- Running Stress Test (50 VUs) ---
Total Requests:  188
Success:         188
Errors:          0
RPS:             93.95
Average:         474.9686ms
P50:             545.429458ms
P95:             736.484916ms
P99:             856.42775ms
```

## Coverage Verification
- Unit coverage: Metric calculations, empty latencies, percentile invariants tested.
- Integration coverage: Real HTTP traffic under concurrency, context cancellation, error counts, dial errors tested.
- Assertions prove claims: Tail latency growth and P95 > Avg rigorously asserted.
- All tests PASS without race conditions.