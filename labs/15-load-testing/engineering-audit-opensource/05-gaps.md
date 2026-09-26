# Gap Analysis

Lab: labs/15-load-testing

Allowed gap types enumerated: MISSING_TEST, BROKEN_IMPLEMENTATION, DOC_CODE_MISMATCH, RACE_CONDITION, UNHANDLED_ERROR, MISSING_EDGE_CASE, IMPLEMENTATION_OVERCLAIM, RESEARCH_MISMATCH, FAKE_DEMO, FAKE_BENCHMARK, UNVERIFIED_RESULT.

## Gap 1 — MISSING_TEST
Location: tests/loadtest_test.go (absent)
Severity: MEDIUM
Description: No test directly asserts that the server's semaphore enforces the `MaxDBConnections` hard bound (peak concurrency never exceeds it). Coverage is indirect via latency growth in `TestLoadTest_SmokeVsStress`.
Why it matters: The design doc explicitly claims "BookingServer: HTTP handler with a configurable concurrency limit (semaphore) simulating a database connection pool." The existence of the limit is proven structurally (semaphore capacity), but no test measures it. A dedicated test that counts concurrent in-flight handlers against the semaphore capacity would close this gap.
Upgrade trigger: When adding stress-test suites, instrument peak concurrency directly.

## Gap 2 — MISSING_TEST
Location: internal/loadtest/metrics_test.go (absent)
Severity: LOW
Description: No unit test asserts `RPS = TotalRequests / Duration` against a known duration and sample count.
Why it matters: RPS is shown in the demo but never unit-validated; it is computed by `CalculateMetrics` from `totalDuration` passed in by the runner.

## Gap 3 — MISSING_TEST
Location: internal/loadtest/metrics_test.go (absent)
Severity: LOW
Description: No single-sample edge case for `percentile`: `len(sorted)==1` yields idx 0 for any pct. The math is covered at 100, 12, and 0 samples, but not exactly at 1.
Why it matters: Smallest valid input that could expose an off-by-one; currently relies on the `len==0` guard.

## Gap 4 — MISSING_TEST
Location: tests/loadtest_test.go (absent)
Severity: LOW
Description: No test asserts the invariant `SuccessCount + ErrorCount == TotalRequests` end-to-end.
Why it matters: This invariant is what makes RPS/errors trustworthy. It holds in all observed runs (errors counted, latencies recorded only on success), but is never asserted.

## Gap 5 — DOC_CODE_MISMATCH
Location: engineering/01-design.md:32
Severity: LOW
Description: Component #2 is described as generating traffic "with specified virtual users (VUs) and iterations." The implementation has no iteration-count config; load runs for a `Duration`. Terminology mismatch only; no functional impact.
Resolution: Align wording to "for a bounded Duration" (or add an optional iteration field and document it).

## Gap 6 — DOC_CODE_MISMATCH
Location: engineering/01-design.md:33
Severity: LOW
Description: Component #3 is described as a "MetricsAggregator: Thread-safe latency collector." There is no aggregator struct; collection is per-VU in `runner.go` and aggregation is the pure function `CalculateMetrics`. The description implies a single thread-safe collector with locking, which is not what the code does. The system is nonetheless race-free.
Resolution: Reword to "Per-VU latency buffers aggregated into `CalculateMetrics`."

## Gap 7 — DOC_CODE_MISMATCH
Location: engineering/02-implementation-notes.md:16
Severity: LOW
Description: "Wait durations in server mock are fixed (20ms)." `server.New` defaults to 10ms if unset and is fully configurable; 20ms is only the demo's value. The statement is context-true only within the demo.
Resolution: Clarify "the demo uses a fixed 20ms; the server accepts a configurable `DBQueryDuration`."

## Non-Issues Verified (no gap present)

- BROKEN_IMPLEMENTATION: none. All code paths (semaphore acquire/release, context cancellation, error paths, percentile math) behave as designed.
- RACE_CONDITION: none. `go test -race ./...` passes; per-VU isolation eliminates shared mutation.
- UNHANDLED_ERROR: none. Server errors are caught at `handleBooking` via status codes; client errors via `ctx.Err()` filtering; body close always reached.
- MISSING_EDGE_CASE: none material. Empty input (0 samples) and unordered long-tail spikes are tested; percentile ordering invariants are asserted.
- IMPLEMENTATION_OVERCLAIM: none. The demo is labeled a sample run and explicitly notes values vary by environment; the invariant (stress tail >> smoke tail) holds across runs.
- RESEARCH_IMPLEMENTATION_MISMATCH: not assessed (per pipeline override, research not audited in this stage).
- FAKE_DEMO: none. Demo was executed; output reproduced the claimed Smoke-vs-Stress contrast.
- FAKE_BENCHMARK: none. No benchmarks exist; the lab does not claim benchmark numbers beyond the demo's own throughput.
- UNVERIFIED_RESULT: none. Build, tests, race detector, and demo were all executed and recorded.

## Summary

4 missing-test gaps (2 LOW, 2 MEDIUM-at-largest) and 3 low-severity doc/code terminology mismatches. No HIGH/CRITICAL gaps. No fabrication, no race conditions, no broken implementation. The MISSING_TEST for the connection-pool hard bound (Gap 1) is the most significant coverage shortfall but is an omission rather than a defect.