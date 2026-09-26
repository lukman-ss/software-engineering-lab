# Test Audit

Target Lab: labs/15-load-testing
Audit Scope: `internal/loadtest/metrics_test.go`, `tests/loadtest_test.go`

---

## 1. Test Suite Coverage Summary

| Test Case | Scope | Behavior Verified | Result |
|---|---|---|---|
| `TestCalculateMetrics` | Unit | Correct P50, P95, P99, Avg, Min, Max calculation across 100 samples | PASS |
| `TestCalculateMetrics_Empty` | Unit | Graceful handling of empty input slice without panic | PASS |
| `TestCalculateMetrics_SingleSample` | Unit | Invariant equality of Min, Max, P50, P95, P99 on 1-element slice | PASS |
| `TestCalculateMetrics_RPS` | Unit | Accurate RPS arithmetic based on total requests and duration | PASS |
| `TestCalculateMetrics_Invariants` | Unit | Monotonic invariant check `Min <= P50 <= P90 <= P95 <= P99 <= Max` on unsorted tail spike data | PASS |
| `TestLoadTest_SmokeVsStress` | Integration | Stress test P95 latency significantly exceeds Smoke test P95 and Average latency | PASS |
| `TestLoadTest_ErrorCount` | Integration | Non-2xx HTTP responses (HTTP 500) correctly tagged as errors with 0 success count | PASS |
| `TestServer_MethodNotAllowed` | Integration | HTTP GET on POST-only endpoint returns 405 Method Not Allowed | PASS |
| `TestLoadTest_DialError` | Integration | Unreachable host dial failures tracked as errors without panic | PASS |
| `TestServer_ContextCanceled` | Integration | Client context cancellation aborts request before completion | PASS |
| `TestServer_MaxDBConnectionsBound` | Concurrency | Live semaphore monitoring confirms concurrent active DB connections never exceed configured maximum | PASS |
| `TestLoadTest_SuccessAndErrorInvariant` | Property | Verifies invariant `SuccessCount + ErrorCount == TotalRequests` holds under concurrent load | PASS |

---

## 2. Test Execution Verification

### Command: `go test -v ./...`
```text
?   	github.com/lukman/software-engineering-lab/labs/15-load-testing/cmd/demo	[no test files]
=== RUN   TestCalculateMetrics
--- PASS: TestCalculateMetrics (0.00s)
=== RUN   TestCalculateMetrics_Empty
--- PASS: TestCalculateMetrics_Empty (0.00s)
=== RUN   TestCalculateMetrics_SingleSample
--- PASS: TestCalculateMetrics_SingleSample (0.00s)
=== RUN   TestCalculateMetrics_RPS
--- PASS: TestCalculateMetrics_RPS (0.00s)
=== RUN   TestCalculateMetrics_Invariants
--- PASS: TestCalculateMetrics_Invariants (0.00s)
PASS
ok  	github.com/lukman/software-engineering-lab/labs/15-load-testing/internal/loadtest	0.323s
?   	github.com/lukman/software-engineering-lab/labs/15-load-testing/internal/server	[no test files]
=== RUN   TestLoadTest_SmokeVsStress
--- PASS: TestLoadTest_SmokeVsStress (1.23s)
=== RUN   TestLoadTest_ErrorCount
--- PASS: TestLoadTest_ErrorCount (0.10s)
=== RUN   TestServer_MethodNotAllowed
--- PASS: TestServer_MethodNotAllowed (0.00s)
=== RUN   TestLoadTest_DialError
--- PASS: TestLoadTest_DialError (0.10s)
=== RUN   TestServer_ContextCanceled
--- PASS: TestServer_ContextCanceled (0.00s)
=== RUN   TestServer_MaxDBConnectionsBound
--- PASS: TestServer_MaxDBConnectionsBound (1.35s)
=== RUN   TestLoadTest_SuccessAndErrorInvariant
--- PASS: TestLoadTest_SuccessAndErrorInvariant (0.20s)
PASS
ok  	github.com/lukman/software-engineering-lab/labs/15-load-testing/tests	3.312s
```

### Command: `go test -race ./...`
```text
?   	github.com/lukman/software-engineering-lab/labs/15-load-testing/cmd/demo	[no test files]
ok  	github.com/lukman/software-engineering-lab/labs/15-load-testing/internal/loadtest	1.106s
?   	github.com/lukman/software-engineering-lab/labs/15-load-testing/internal/server	[no test files]
ok  	github.com/lukman/software-engineering-lab/labs/15-load-testing/tests	4.548s
```

---

## 3. Test Quality Evaluation

- Happy path covered: Yes (`TestLoadTest_SmokeVsStress`, `TestCalculateMetrics`).
- Edge cases covered: Yes (`TestCalculateMetrics_Empty`, `TestCalculateMetrics_SingleSample`).
- Failure/Negative paths covered: Yes (`TestLoadTest_ErrorCount`, `TestLoadTest_DialError`, `TestServer_MethodNotAllowed`, `TestServer_ContextCanceled`).
- Concurrency & Invariants covered: Yes (`TestServer_MaxDBConnectionsBound`, `TestCalculateMetrics_Invariants`, `TestLoadTest_SuccessAndErrorInvariant`).
- Race detector passes: Yes, clean with 0 warnings or data race reports.
