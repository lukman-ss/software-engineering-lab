# Changes Made

## Revision 1

Audit Issue: Gap 1 — RESEARCH_MISMATCH (Mock server uniform queuing failure to demonstrate masking effect of averages)
Severity: HIGH
Files Changed: `internal/server/server.go`
Action: Modified database query execution to introduce tail contention when active requests exceed connection capacity. Under saturated conditions, a 10% probability triggers a 25x query duration penalty, inducing long-tail latency spikes that are concealed by average metrics.
Verification: Executed `go run ./cmd/demo`. Smoke test shows Avg ~21ms and P95 ~21ms; Stress test demonstrates Avg ~739ms versus P95 ~1356ms (P95 is 1.83x higher than Avg).
Status: RESOLVED

## Revision 2

Audit Issue: Test Audit Finding 1 & Gap 2 — TEST_ASSERTION_DEFICIENCY (Test fails to assert divergence between P95 and Average)
Severity: MEDIUM
Files Changed: `tests/loadtest_test.go`
Action: Added explicit assertion in `TestLoadTest_SmokeVsStress` ensuring `stressRes.P95Latency > stressRes.AvgLatency`.
Verification: Executed `go test -v ./tests -run TestLoadTest_SmokeVsStress` across multiple iterations with zero failures.
Status: RESOLVED

## Revision 3

Audit Issue: Gap 3 — MISSING_TEST (No invariant validation for percentile monotonicity)
Severity: LOW
Files Changed: `internal/loadtest/metrics_test.go`
Action: Added `TestCalculateMetrics_Invariants` checking `Min <= P50 <= P90 <= P95 <= P99 <= Max` on an unsorted, non-monotonic latency dataset with tail outliers.
Verification: Executed `go test -v ./internal/loadtest -run TestCalculateMetrics_Invariants`. Test passes cleanly.
Status: RESOLVED

## Revision 4

Audit Issue: Gap 2 & Docs Audit Finding 1 — DOC_CODE_MISMATCH
Severity: MEDIUM
Files Changed: `engineering/03-execution-result.md`
Action: Refreshed execution logs with newly recorded runs for unit tests, race detection, and demo outputs demonstrating the long-tail latency divergence.
Verification: Inspected `engineering/03-execution-result.md` against live console output.
Status: RESOLVED
