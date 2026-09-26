# Test Audit

Lab: labs/15-load-testing

## Test Inventory
- `internal/loadtest/metrics_test.go`
  - `TestCalculateMetrics` — happy path: 100 known latencies (1ms..100ms), asserts TotalRequests, Min, Max, Avg, P50, P95, P99.
  - `TestCalculateMetrics_Empty` — edge case: nil latencies + 0 errors.
  - `TestCalculateMetrics_Invariants` — asserts ordering Min <= P50 <= P90 <= P95 <= P99 <= Max on a long-tail spike sample.
- `tests/loadtest_test.go`
  - `TestLoadTest_SmokeVsStress` — integration: 1 VU smoke vs 10 VU stress against 2-slot server, asserts stress P95 > smoke P95 and stress P95 > stress Avg.
  - `TestLoadTest_ErrorCount` — failure path: HTTP 500 server, asserts all responses counted as errors, 0 successes.
  - `TestLoadTest_DialError` — negative path: unreachable port, asserts all attempts errors, 0 successes.
  - `TestServer_MethodNotAllowed` — failure path: GET on POST-only endpoint, asserts 405.
  - `TestServer_ContextCanceled` — recovery/abort path: pre-canceled context, asserts no 201 returned.

## Execution Results (recorded verbatim)
```
go build ./...        -> success, no output
go test -v ./...      -> PASS (all 8 tests, see 02-code-audit for full block)
go test -race ./...   -> PASS (no races detected)
go run ./cmd/demo     -> runs; stress P95/P99 >> smoke P95/P99 (invariant holds across 2 sample runs)
```
All commands executed against the checkout; no repairs made.

---

## Coverage Matrix

| Category            | Covered? | Location | Notes |
|---------------------|----------|----------|-------|
| Happy path          | PASS     | metrics_test.go:8, loadtest_test.go:14 | Metrics happy path + successful request flow |
| Failure path        | PASS     | loadtest_test.go:60, :87 | HTTP 500 and 405 |
| Edge case (empty)   | PASS     | metrics_test.go:43 | Zero samples |
| Edge case (ordering invariants) | PASS | metrics_test.go:53 | Long-tail spikes |
| Transitions         | PARTIAL  | loadtest_test.go:14 | Smoke->Stress contrast, but no explicit state transition under load |
| Recovery / rollback | NOT_APPLICABLE | — | No persistent state; rollback not in scope |
| Concurrency / races | PASS     | runner.go per-VU isolation + `go test -race` green | No shared mutable state during run |
| Negative cases      | PASS     | loadtest_test.go:103 | Dial error against dead port |
| Cancellation        | PASS     | loadtest_test.go:122, runner.go ctx filtering | Context abort, timeout, cancellation filtering |

## Detailed Findings

## Finding 1
Location: internal/loadtest/metrics_test.go
Claimed Coverage: Percentile accuracy for P50, P95, P99.
Observed: `TestCalculateMetrics` uses an arithmetic sequence 1..100ms and asserts exact values (P50=50ms, P95=95ms, P99=99ms) consistent with the nearest-rank implementation.
Assessment: PASS
Severity: LOW
Notes: Strong, deterministic validation of the percentile math.

## Finding 2
Location: tests/loadtest_test.go:14 (TestLoadTest_SmokeVsStress)
Claimed Coverage: Stress-induced latency growth (P95 divergence).
Observed: 1 VU smoke vs 10 VU stress against a 2-slot server over 500ms. Asserts stress P95 > smoke P95 and stress P95 > stress Avg.
Assessment: PASS, with caveat
Severity: MEDIUM
Notes: This is the lab's central behavioral claim and it is proven. The assertion `stressRes.P95Latency <= smokeRes.P95Latency == false` could theoretically be flaky on a heavily loaded CI box, but 10x oversubscription over 500ms makes queueing deterministic in practice. Flagged as MEDIUM risk because the assertion is timing-based rather than structural (see 05-gaps).

## Finding 3
Location: tests/loadtest_test.go:60 (TestLoadTest_ErrorCount)
Claimed Coverage: Error counting on non-2xx responses.
Observed: Mock returning 500; asserts ErrorCount == TotalRequests and SuccessCount == 0.
Assessment: PASS
Severity: LOW
Notes: Confirms the `StatusCode >= 400` branch (metrics not recorded, errors incremented).

## Finding 4
Location: tests/loadtest_test.go:103 (TestLoadTest_DialError)
Claimed Coverage: Network/error-path handling.
Observed: Targets `127.0.0.1:1` (closed port); asserts all attempts are errors, 0 successes.
Assessment: PASS
Severity: LOW
Notes: Exercises `ctx.Err() == nil` filtering on real connection failures.

## Finding 5
Location: tests/loadtest_test.go:87, :122
Claimed Coverage: Request cancellation / context handling.
Observed: 405 on GET; pre-canceled context yields no 201.
Assessment: PASS
Severity: LOW

## Finding 6
Location: missing
Claimed Coverage: Hard concurrency limit (semaphore enforces MaxDBConnections).
Observed: No test asserts that concurrent in-flight requests to the server never exceed `MaxDBConnections`. Coverage is indirect (latency growth implies queuing) but the hard bound is never measured.
Assessment: FAIL (gap)
Severity: MEDIUM
Notes: A test that injects a counter via the semaphore / waits on the server to observe peak concurrency would directly prove the documented "configurable concurrency limit." Absent.

## Finding 7
Location: missing
Claimed Coverage: RPS accuracy and latency/throughput invariants.
Observed: No unit test for `CalculateMetrics.RPS` against a known duration/sample count.
Assessment: FAIL (gap)
Severity: LOW

## Finding 8
Location: missing
Claimed Coverage: Single-sample percentile edge case.
Observed: No test with len(latencies)==1; `percentile` index math is exercised only at scale (100) and small scale (12).
Assessment: FAIL (gap)
Severity: LOW

## Finding 9
Location: missing
Claimed Coverage: SuccessCount + ErrorCount == TotalRequests invariant.
Observed: No test asserts the composition invariant explicitly.
Assessment: FAIL (gap)
Severity: LOW

## Summary

The test suite proves the core claim (stress P95 degrades relative to smoke) and validates the percentile math and error paths. It executes cleanly under the race detector. Gaps are coverage holes, not correctness bugs: no direct assertion of the connection-pool hard bound, no RPS unit test, and no single-sample edge case. These are documented as MISSING_TEST in 05-gaps.md.
