# Test Audit

## Coverage Mapping

### Unit Tests (internal/loadtest/metrics_test.go)
- TestCalculateMetrics: PASS
  - Tests 100 latencies from 1ms to 100ms
  - Verifies TotalRequests, Min, Max, Average, P50, P95, P99
  - All assertions correct and passing
- TestCalculateMetrics_Empty: PASS
  - Tests empty latency list edge case
  - Verifies zero values, no panic
- TestCalculateMetrics_Invariants: PASS
  - Tests 12 unsorted latencies with tail spike (500ms)
  - Verifies Min <= P50 <= P90 <= P95 <= P99 <= Max ordering invariant
  - Validates sorting correctness on unsorted input

### Integration Tests (tests/loadtest_test.go)
- TestLoadTest_SmokeVsStress: PASS
  - Smoke test: 1 VU, 2 max DB connections → 0 errors, total requests non-zero
  - Stress test: 10 VUs, 2 max DB connections → 0 errors, total requests non-zero
  - Asserts stress P95 > smoke P95 (proves latency degradation under load)
  - Asserts stress P95 > stress Avg (proves tail latency masking effect)
  - Covers: happy path, smoke vs stress comparison, percentile accuracy
- TestLoadTest_ErrorCount: PASS
  - Mock server returning HTTP 500
  - 2 VUs for 100ms
  - Verifies ErrorCount == TotalRequests (all errors counted)
  - Verifies SuccessCount == 0
  - Covers: failure path (HTTP 500), error counting
- TestServer_MethodNotAllowed: PASS
  - GET request to POST-only /booking endpoint
  - Verifies 405 status returned
  - Covers: HTTP method validation
- TestLoadTest_DialError: PASS
  - VU targets http://127.0.0.1:1 (unreachable)
  - Verifies all requests counted as errors, no successes
  - Covers: network failure path, dial error handling
- TestServer_ContextCanceled: PASS
  - Server with 1 DB slot, 100ms query duration
  - Pre-cancelled context
  - Verifies request aborts (not 201 Created)
  - Covers: context cancellation in server, early abort

## Test Audit Summary

| Test Category | Coverage | Notes |
|---|---|---|
| Happy path | Covered (SmokeVsStress, CalculateMetrics) | |
| Failure path | Covered (ErrorCount, DialError) | |
| Edge cases | Covered (Empty metrics, MethodNotAllowed) | |
| Transitions | Covered (Smoke → Stress comparison) | |
| Recovery | NOT covered | No test for connection recovery, queue draining |
| Rollback | NOT applicable | No transactional state in server |
| Concurrency | Covered (race detector passes) | All tests run with -race |
| Negative cases | Covered (DialError, MethodNotAllowed) | |
| Error propagation | Covered (ErrorCount, DialError) | Errors correctly tallied |

## Test Results

```text
go test -v ./...

=== RUN   TestCalculateMetrics
--- PASS: TestCalculateMetrics (0.00s)
=== RUN   TestCalculateMetrics_Empty
--- PASS: TestCalculateMetrics_Empty (0.00s)
=== RUN   TestCalculateMetrics_Invariants
--- PASS: TestCalculateMetrics_Invariants (0.00s)
PASS
ok  	github.com/lukman/software-engineering-lab/labs/15-load-testing/internal/loadtest

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
ok  	github.com/lukman/software-engineering-lab/labs/15-load-testing/tests

go test -race ./...
ok  	github.com/lukman/software-engineering-lab/labs/15-load-testing/internal/loadtest
ok  	github.com/lukman/software-engineering-lab/labs/15-load-testing/tests
```

## Assessment: PASS with minor gaps

The test suite adequately covers the core claims:
- Percentile calculations are verified with exact expectations
- Smoke vs stress comparison validates latency degradation theory
- Error handling paths are covered (500, dial errors, context cancel)

Missing coverage:
- No test verifying the 10% degradation probability in server
- No test for graceful shutdown under load
- No test for queue draining when requests release semaphore