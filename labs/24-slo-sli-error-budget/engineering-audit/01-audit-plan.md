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
- `research-audit/07-verdict.md`

Main Claims To Verify:
1. SLI calculation implements the `good_events / total_events` ratio over sliding windows correctly.
2. Error Budget calculation adheres to `(1.0 - target_slo) * total_events` and properly tracks consumed/remaining budget.
3. Release freeze policy (`CanDeploy`) correctly evaluates to `false` when remaining budget is exhausted (`budgetRemaining <= 0`).
4. Multi-window burn rate alert engine checks both short and long windows against specified burn rate factors (`shortBurn >= factor && longBurn >= factor`).
5. Window tracker is concurrency safe and handles out-of-order timestamp events without panicking or corrupting state.
6. Documentation (`README.md`, `engineering/*.md`) accurately matches actual implementation and test commands.

Commands To Run:
- `go test ./...`
- `go test -race ./...`
- `go test -count=1 -v -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Floating point precision rounding in `math.Round` causing boundary calculation anomalies for `CanDeploy`.
- Out-of-order slice insertion logic in `WindowTracker` causing slice corruption or improper time sorting during eviction.
- Divergence between demo output claims and actual runtime execution results.
