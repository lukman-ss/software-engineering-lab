# Test Audit

Target Lab: labs/15-load-testing

## Test Suite Coverage Overview

### 1. `internal/loadtest/metrics_test.go`
- `TestCalculateMetrics`: Verifies accurate computation of Min, Max, Avg, P50, P90, P95, P99, and RPS across a known 100-element latency slice.
- `TestCalculateMetrics_Empty`: Verifies behavior with empty latency inputs and zero total duration (returns zero-value `Result` without panic or divide-by-zero).
- `TestCalculateMetrics_Invariants`: Verifies mathematical ordering invariants (`Min <= P50 <= P90 <= P95 <= P99 <= Max`).

### 2. `tests/loadtest_test.go`
- `TestLoadTest_SmokeVsStress`: Runs a 1 VU smoke test vs a 10 VU stress test against a 2-slot server. Verifies `P95_stress > P95_smoke` and `P95_stress > Avg_stress`.
- `TestLoadTest_ErrorCount`: Verifies HTTP 500 responses are recorded as errors (`ErrorCount == TotalRequests`).
- `TestServer_MethodNotAllowed`: Verifies server rejects non-POST requests to `/booking` with 405 Method Not Allowed.
- `TestLoadTest_DialError`: Verifies connection failure to invalid address (`127.0.0.1:1`) is recorded as errors.
- `TestServer_ContextCanceled`: Verifies server immediately aborts processing if client cancels context before acquiring connection slot.

## Test Execution Results

```text
=== RUN   TestCalculateMetrics
--- PASS: TestCalculateMetrics (0.00s)
=== RUN   TestCalculateMetrics_Empty
--- PASS: TestCalculateMetrics_Empty (0.00s)
=== RUN   TestCalculateMetrics_Invariants
--- PASS: TestCalculateMetrics_Invariants (0.00s)
=== RUN   TestLoadTest_SmokeVsStress
--- PASS: TestLoadTest_SmokeVsStress (1.28s)
=== RUN   TestLoadTest_ErrorCount
--- PASS: TestLoadTest_ErrorCount (0.10s)
=== RUN   TestServer_MethodNotAllowed
--- PASS: TestServer_MethodNotAllowed (0.00s)
=== RUN   TestLoadTest_DialError
--- PASS: TestLoadTest_DialError (0.10s)
=== RUN   TestServer_ContextCanceled
--- PASS: TestServer_ContextCanceled (0.00s)
PASS
```

Race Detector Verification (`go test -race ./...`):
- All 8 tests pass with 0 data races detected.

## Assessment
PASS — Test coverage addresses happy path, queue saturation under stress, error accounting (HTTP 5xx and Dial errors), edge cases (empty metrics, HTTP 405, context cancellation), and mathematical invariants.
