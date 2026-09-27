# Engineering Audit Plan

Target Lab: labs/24-slo-sli-error-budget
Implementation Files:
- internal/metrics/tracker.go (WindowTracker, sliding time-bucketed aggregation)
- internal/slo/evaluator.go (SLI ratio, error budget, CanDeploy freeze)
- internal/alerting/engine.go (burn-rate math, dual-window AND check)
- cmd/demo/main.go (4-phase executable demo)
- tests/slo_test.go (6 tests)
Tests: tests/slo_test.go — TestMetricsWindowTracker, TestSLOEvaluator, TestAlertEngineBurnRate, TestOutOfOrderTimestamps, TestEvaluatorZeroTraffic, TestConcurrencyMetrics
Executable/Demo: cmd/demo (go run ./cmd/demo)
Approved Research Inputs: NOT AUDITED per pipeline override (implementation + tests only)
Main Claims To Verify:
1. SLI = good/total ratio over sliding window
2. ErrorBudget = (1-target)*total; freeze when remaining <= 0
3. BurnRate = actualErrorRate/allowedErrorRate; alert only when short AND long >= threshold (transient suppression)
4. Criticality bucketing (99.9% vs 95% SLOs behave differently)
5. Thread-safe concurrent recording; race detector clean
6. Demo output real; execution-result matches reality
Commands To Run:
- go vet ./...
- go test -count=1 -v ./...
- go test -count=1 -race ./...
- go run ./cmd/demo
Primary Risks:
- Per-rule window fields (LongWindow/ShortWindow/BudgetConsumedPct) ignored by Check — shared-tracker thresholds only
- Demo Phase 4 uses 10% error rate against both SLOs so both freeze; differentiation claim unproven by demo
- Nil isGoodEvent / nil tracker panics unguarded; invalid TargetUptime unvalidated
- Suite passes but lacks 100%-errors, burn-edge (total=0, SLO=1.0), concurrent-reader cases behind "100% coverage" claim
