# Engineering Audit Plan

Target Lab: `labs/24-slo-sli-error-budget`
Implementation Files:
- `internal/metrics/tracker.go`
- `internal/slo/evaluator.go`
- `internal/alerting/engine.go`
Tests:
- `tests/slo_test.go`
Executable/Demo:
- `cmd/demo/main.go`
Approved Research Inputs:
- `research/05-report.md`
- `research/runs/2026-09-27-slo-sli-error-budget/05-report.md`
Main Claims To Verify:
1. SLI quantitative measurement as `good_events / total_events` ratio over sliding windows.
2. Error Budget calculation (`1 - SLO`) with release freeze policy enforcement (`CanDeploy = false` when exhausted).
3. Multi-window multi-burn-rate alerting logic triggering based on short and long window burn rates exceeding thresholds.
4. Endpoint criticality differentiation (e.g. 99.9% critical vs 95.0% non-critical).
5. Thread-safe concurrent event recording and metric evaluation.
Commands To Run:
- `go test -count=1 ./...`
- `go test -count=1 -race ./...`
- `go run ./cmd/demo`
Primary Risks:
- Race conditions in sliding-window bucket updates or eviction.
- Numerical inaccuracy/rounding errors in error budget and burn rate calculations.
- Discrepancy between demo output, engineering notes, and codebase implementation.
