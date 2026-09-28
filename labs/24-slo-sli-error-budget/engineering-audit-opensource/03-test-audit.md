# Test Audit — labs/24-slo-sli-error-budget/tests/slo_test.go

Coverage matrix (6 tests):

## TestMetricsWindowTracker
- happy path, eviction. PASS.
- edge: 0 after window expiry. covered.

## TestSLOEvaluator
- happy path 99/1 SLI=0.99 CanDeploy. PASS.
- failure path 1 extra bad → CanDeploy=false. covered.

## TestAlertEngineBurnRate
- happy path: 2% error ×100x? actually 2/100 = 2% → burn 20× vs 14.4 factor → triggered. correct.
- negative/transient: short 10% / long 0.01% → 0 alerts. PASS.

## TestOutOfOrderTimestamps
- edge: out-of-order insert + partial eviction. PASS.

## TestEvaluatorZeroTraffic
- edge: total=0 SLI=1.0 CanDeploy=true. covered.

## TestConcurrencyMetrics
- concurrency: 20×100 = 2000 ops concurrent, race-clean. PASS.

## Gaps
- MISSING_TEST: `CalculateBurnRate` with total=0 or targetSLO such that allowedErrorRate<=0 — returns 0 but untested.
- MISSING_TEST: `BurnRateRule` zero `BurnRateFactor` would always trigger (shortBurn>=0 and longBurn>=0) — no guard.
- MISSING_TEST: tracker with zero/negative windowSize or nil isGood — no construction validation test.
- MISSING_TEST: SLO evaluator exact-boundary budgetRemaining == 0 (not <= due to float) — no explicit test.
- MISSING_TEST: recovery path — design promises recovery demo; tests never re-record good traffic after incident to show CanDeploy flip back to true. Not implemented in demo either.
- MISSING_TEST: 100% error / 0% error edge for SLO (division-safe) — partially implied by zero-traffic test only.

## Race Detector
`go test -race ./...` → PASS (tests ok). No data races reported. PASS.

## Overall Test Quality: PASS-STRONG-WEAKS
Core math happy+failure paths and concurrency proven. Negative/transient suppression works. Missing pure edge tests and no failure-rollback/recovery coverage.
