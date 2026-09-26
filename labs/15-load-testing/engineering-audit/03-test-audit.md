# Engineering Test Audit

Target Lab: labs/15-load-testing

## Coverage

- Happy path: Covered by `TestLoadTest_SmokeVsStress` and unit tests.
- Failure path: Covered by `TestLoadTest_ErrorCount` and `TestLoadTest_DialError`.
- Edge cases: Empty metrics calculation covered by `TestCalculateMetrics_Empty`.
- Concurrency: Proven safe via `go test -race`.

## Findings

1. `TestLoadTest_SmokeVsStress`:
   - Checks that `smokeRes.ErrorCount == 0`.
   - Checks that `stressRes.P95Latency > smokeRes.P95Latency`.
   - Does NOT verify that `stressRes.P95Latency` diverges significantly from `stressRes.AvgLatency`. The test suite fails to enforce the core lab claim (the masking effect of averages).

2. `TestCalculateMetrics`:
   - Validates metrics math effectively on a 1-100ms artificial latency slice.

3. Context Cancellation:
   - `TestServer_ContextCanceled` proves the server abandons semaphore queues accurately when requests timeout.

Overall Assessment: The tests pass and have good mechanical coverage, but lack assertion of the primary architectural claim.
