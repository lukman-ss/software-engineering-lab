# Engineering Audit Plan

Target Lab: labs/24-slo-sli-error-budget
Implementation Files:
- internal/metrics/tracker.go (WindowTracker)
- internal/slo/evaluator.go (Evaluator)
- internal/alerting/engine.go (AlertEngine)
- cmd/demo/main.go
Tests:
- tests/slo_test.go (6 tests: TestMetricsWindowTracker, TestSLOEvaluator, TestAlertEngineBurnRate, TestOutOfOrderTimestamps, TestEvaluatorZeroTraffic, TestConcurrencyMetrics)
Executable/Demo: cmd/demo/main.go
Approved Research Inputs: SKIPPED per PIPELINE OVERRIDE (implementation + tests only)
Main Claims To Verify:
1. SLI = good/total ratio over sliding window
2. ErrorBudget = (1-target)*total; CanDeploy=false when exhausted
3. Multi-window burn-rate alerting (short AND long >= factor)
4. Window expiry/eviction correctness + out-of-order handling
5. Zero-traffic SLI=1.0, CanDeploy=true
6. Concurrency safety under parallel Record
7. Demo output real, matches engineering/03-execution-result.md
Commands To Run:
- go build ./...
- go test -count=1 -v ./...
- go test -count=1 -race ./...
- go run ./cmd/demo
Primary Risks:
- BurnRateRule window fields unused; per-rule windows decorative
- Config.LatencyThreshold dead in evaluator (goodness via closure)
- Nil isGoodEvent panic; zero/negative windowSize unguarded
