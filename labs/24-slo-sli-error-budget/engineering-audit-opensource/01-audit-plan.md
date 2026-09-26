# Engineering Audit Plan

Target Lab: labs/24-slo-sli-error-budget

Implementation Files:
- internal/metrics/tracker.go (WindowTracker, Bucket, Event)
- internal/slo/evaluator.go (Evaluator, Status)
- internal/alerting/engine.go (AlertEngine, BurnRateRule)
- cmd/demo/main.go (4-phase executable demo)

Tests:
- tests/slo_test.go (6 tests: TestMetricsWindowTracker, TestSLOEvaluator, TestAlertEngineBurnRate, TestOutOfOrderTimestamps, TestEvaluatorZeroTraffic, TestConcurrencyMetrics)

Executable/Demo:
- cmd/demo: 4-phase simulation (baseline, incident, burn-rate alerts, endpoint criticality comparison)

Approved Research Inputs:
- Per pipeline override: research/content NOT audited in this stage. Only implementation and tests. Engineering docs (01-design.md, 02-implementation-notes.md, 03-execution-result.md) consulted as claimed-behavior reference, not as research verdict.

Main Claims To Verify:
1. SLI = good/total ratio tracked over sliding window
2. Error budget = (1 - SLO) * total; CanDeploy=false when budget exhausted
3. Multi-window multi-burn-rate alerting (fast 14.4x page, slow 6x ticket) requiring both windows above threshold
4. Endpoint criticality (Payment 99.9% vs Reports 95.0%)
5. Thread-safe concurrent recording (race detector clean)
6. Demo shows real budget depletion and alert triggering

Commands To Run:
- go build ./...
- go test ./... -v
- go test -race ./...
- go run ./cmd/demo
- go test -coverpkg=./internal/... ./tests/

Primary Risks:
- AlertEngine uses single global short/long burn for all rules (not per-rule windows); BudgetConsumedPct field never enforced
- WindowTracker out-of-order insertion appends unbounded (no capacity cap despite numBuckets)
- Summary takes write Lock (not RLock) — correctness fine, minor contention
- Recorded execution-result.md omits 2 tests and Phase 4 output — stale doc, not fake code
