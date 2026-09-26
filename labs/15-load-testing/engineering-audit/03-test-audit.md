# Test Audit

## Coverage Areas

- **Happy Path:** `TestCalculateMetrics` verifies basic latency aggregation and percentile calculations. `TestLoadTest_SmokeVsStress` verifies standard Smoke execution.
- **Failure Path:** `TestLoadTest_ErrorCount` forces 500 status codes and asserts error tally logic. `TestServer_MethodNotAllowed` forces 405 error on the mock server.
- **Edge Cases:** `TestCalculateMetrics_Empty` verifies division-by-zero prevention when zero requests complete.
- **Transitions / Degradation:** `TestLoadTest_SmokeVsStress` proves that stress-induced queuing causes `stress P95 > smoke P95`.
- **Concurrency:** Fully verified under `-race`. `runner.go` prevents data races organically via index segregation.

## Actual Results

```bash
$ go test -v -count=1 ./...
=== RUN   TestCalculateMetrics
--- PASS: TestCalculateMetrics (0.00s)
=== RUN   TestCalculateMetrics_Empty
--- PASS: TestCalculateMetrics_Empty (0.00s)
PASS
ok      github.com/lukman/software-engineering-lab/labs/15-load-testing/internal/loadtest      0.134s
=== RUN   TestLoadTest_SmokeVsStress
--- PASS: TestLoadTest_SmokeVsStress (1.05s)
=== RUN   TestLoadTest_ErrorCount
--- PASS: TestLoadTest_ErrorCount (0.10s)
=== RUN   TestServer_MethodNotAllowed
--- PASS: TestServer_MethodNotAllowed (0.00s)
PASS
ok      github.com/lukman/software-engineering-lab/labs/15-load-testing/tests  1.155s
```

```bash
$ go test -race ./...
ok      github.com/lukman/software-engineering-lab/labs/15-load-testing/internal/loadtest      1.045s
ok      github.com/lukman/software-engineering-lab/labs/15-load-testing/tests  2.324s
```

```bash
$ go run ./cmd/demo
Starting Load Test Demo (Booking Bengkel)
Server simulated DB connections: 5
Simulated DB query duration: 20ms

--- Running Smoke Test (2 VUs) ---
Total Requests:  186
Success:         186
Errors:          0
RPS:             92.96
Average:         21.453647ms
P50:             21.43425ms
P95:             21.599334ms
P99:             24.139542ms

--- Running Stress Test (50 VUs) ---
Total Requests:  473
Success:         473
Errors:          0
RPS:             236.35
Average:         200.506762ms
P50:             209.993916ms
P95:             210.767ms
P99:             210.999708ms
```

## Assessment
The tests are robust and directly assert the behavioral claims of the research (tail latency growth under resource saturation). Execution is deterministic and race-free. Output metrics from `go run ./cmd/demo` match the `engineering/03-execution-result.md` claims completely.
