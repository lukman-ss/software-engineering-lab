# Gap Analysis

## MISSING_EDGE_CASE
Location: internal/slo/evaluator.go:41-71
Description: SLOEvaluator.Evaluate has no test branch for `total == 0`. While the guard `if total > 0` exists, no test verifies behavior with zero traffic. Division-by-zero avoidance is guarded but unverified under test.
Severity: LOW

## MISSING_EDGE_CASE
Location: internal/alerting/engine.go:51-61
Description: CalculateBurnRate branches `total==0` and `allowedErrorRate<=0` are not exercised by tests.
Severity: LOW

## MISSING_EDGE_CASE
Location: internal/metrics/tracker.go:30-37
Description: NewWindowTracker default branch (bucketSize<=0, numBuckets<1) not covered.
Severity: LOW

## TEST_CLAIM_MISMATCH
Location: engineering/01-design.md:21
Description: Design documents state "100% test coverage on core math and sliding window calculations." Measured coverage: CalculateBurnRate 71.4%, NewWindowTracker 66.7%, other functions 100%. Statement is inaccurate.
Severity: MEDIUM

## DOC_CODE_MISMATCH
Location: internal/alerting/engine.go:17-24
Description: BurnRateRule defines LongWindow, ShortWindow, BudgetConsumedPct fields intended to parameterize per-rule alert windows. These are NOT used in Check(); the engine uses the fixed shortTracker/longTracker passed to NewAlertEngine. Design implies per-rule windows; code does not implement them.
Severity: MEDIUM

## RESEARCH_MISMATCH
Location: engineering/01-design.md:10
Description: Design claims endpoint criticality bucketing with distinct SLOs per endpoint (Payment 99.9%, Reports 95.0%). Code has Event.Endpoint field but no per-endpoint SLO logic; demo uses a single fixed SLO. Feature missing.
Severity: HIGH

## IMPLEMENTATION_OVERCLAIM
Location: slo.Config.LatencyThreshold (internal/slo/evaluator.go:13)
Description: LatencyThreshold stored on SLO Config but never consulted by Evaluator.Evaluate; latency enforcement lives in the metrics isGood closure. The SLO component claims latency awareness it does not exercise.
Severity: MEDIUM

## CONCURRENCY_GAP
Location: tests/slo_test.go:131-166
Description: TestConcurrencyMetrics calls Record concurrently but invokes Summary only AFTER `wg.Wait()`. No test exercises concurrent Record + Summary, which is the unsafe path the race detector would flag.
Severity: MEDIUM