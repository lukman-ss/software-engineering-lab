# Test Audit — labs/24-slo-sli-error-budget

Coverage (measured): `go test -coverpkg=./internal/... ./tests/` => 94.1% of statements in ./internal/...

## Coverage By Function (per go tool cover -func)

| Package | Function | Status |
|---|---|---|
| internal/metrics | NewWindowTracker | covered |
| internal/metrics | Record | covered |
| internal/metrics | evictStaleLocked | covered |
| internal/metrics | Summary | covered |
| internal/slo | NewEvaluator | covered |
| internal/slo | Evaluate | covered |
| internal/alerting | NewAlertEngine | covered |
| internal/alerting | CalculateBurnRate | covered |
| internal/alerting | Check | covered |
| cmd/demo | main | 0% (no test) |

## Test Inventory

1. TestMetricsWindowTracker — happy path (10 good + 1 slow + 1 HTTP 500) bucketed counts + eviction past window.
2. TestSLOEvaluator — boundary: 99 good/1 bad -> SLI exactly 0.99, CanDeploy true; then +1 bad -> SLI < 0.99, CanDeploy false.
3. TestAlertEngineBurnRate — positive: 98 good/2 bad (2% error, 9.09x burn > 14.4? NO, > 6? YES) asserts 1 alert, Page severity. Wait — rule BurnRateFactor is 14.4 here, but shortBurn/longBurn = (2/100)/0.001 = 20x, 20 >= 14.4 => triggers. Re-check below. Plus negative: transient short-window spike (100x short, 0.1x long) -> 0 alerts.
4. TestOutOfOrderTimestamps — out-of-order insert + partial eviction read.
5. TestEvaluatorZeroTraffic — zero traffic -> SLI 1.0, CanDeploy true.
6. TestConcurrencyMetrics — 20 goroutines × 100 records; asserts total == 2000 and good+bad == total; race-clean.

## Test 3 arithmetic re-check (TestAlertEngineBurnRate positive case)
- 98 good + 2 bad = 100 total in both windows.
- allowedErrorRate = 1 - 0.999 = 0.001.
- actualErrorRate = 2/100 = 0.02.
- burn = 0.02/0.001 = 20x. Rule factor = 14.4. 20 >= 14.4 (short) AND 20 >= 14.4 (long) => triggered, 1 alert. Assertion `len(alerts)==1` PASS. Severity==Page matches rule. PASS.
(Engineering note: the test comment text "20x burn rate (> 14.4x)" is correct; code is correct. PASS with note: test only fires ONE rule so it does not prove multi-rule independence.)

## Coverage Gaps (what is NOT tested)
- NO test for SLO = 1.0 (allowedErrorRate == 0) path in CalculateBurnRate's guard (engine.go:57-59).
- NO test asserting the slow-burn (6.0x TICKET) rule vs fast-burn (14.4x PAGE) rule distinction: the demo does this but tests only assert the single-rule positive case.
- NO test for 100% error traffic (SLI == 0.0, BudgetRemaining fully negative, CanDeploy false) — only TestSLOEvaluator's near-threshold and below cases.
- NO test for latency-only bad events at the SLO layer (good via status 200 but over latency threshold): TestMetricsWindowTracker uses slow duration with StatusCode 200 and isGood checks `Duration <= 100ms`, so it IS covered. PASS.
- NO test that CanDeploy becomes false AND THEN flips back true once budget recovers (no "recovery/rollback" path tested — the lab is monotonic in tests).
- NO negative test that a rule with BudgetConsumedPct set behaves differently (budget-consumed gating is not implemented; see Finding 1).
- NO unit test for eviction-by-time ordering edge (a bucket with StartTime == cutoff boundary: `Before` excludes equal — boundary is half-open [cutoff, now)? Tests do not probe the boundary).

## Test Quality Assessment
- Happy path: PASS (TestMetricsWindowTracker, TestSLOEvaluator, TestAlertEngineBurnRate).
- Failure path: WEAK. No HTTP/5xx-flood + latency-spike-only test; no SLO=100% degenerate test; no recovery-from-false test.
- Edge cases: PARTIAL (zero traffic covered; 100% errors not covered; boundary eviction not covered).
- Transitions/recovery/rollback: NOT TESTED. The design's "success criteria" explicitly listed "rollback" but no test asserts budget recovery or deploy-flag flip.
- Concurrency: PASS (TestConcurrencyMetrics + `-race` clean).
- Race detector: PASS (`go test -race ./...` => ok).

## Verdict on Tests
Sufficient to prove core happy-path math and thread-safety. NOT sufficient to prove the full success criteria claimed in 01-design.md ("100% test coverage" / recovery / rollback / multi-rule). Coverage 94.1%, not 100%.

Assessment: PASS (test suite runs, race-clean, covers core claims) but with MEDIUM caveats on incomplete edge/failure/recovery coverage.
