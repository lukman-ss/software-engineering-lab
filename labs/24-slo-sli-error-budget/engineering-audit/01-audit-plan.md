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
- Research phase completed with focus on Google SRE SLI/SLO ratios, Error Budget, and Multi-Window Multi-Burn-Rate alerting.

Main Claims To Verify:
1. SLI correctly calculated as `good_events / total_events`.
2. Error Budget calculation and depletion (`TotalErrorBudget = (1 - SLO) * total_events`, `BudgetRemaining = TotalErrorBudget - bad_events`).
3. Release freeze policy (`CanDeploy = false`) when Error Budget is exhausted.
4. Multi-window multi-burn-rate alerting logic requiring both short and long windows to breach thresholds.
5. Code compiles cleanly and passes all test suites including race detector (`go test -race ./...`).
6. Demo output is real and reproducible via `go run ./cmd/demo`.
7. Documentation in `README.md` accurately reflects code structure and instructions.

Commands To Run:
- `go test -v -count=1 ./...`
- `go test -v -count=1 -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Race conditions during concurrent `Record` or `Summary` calls in `metrics.WindowTracker`.
- Incorrect out-of-order event ingestion or sliding window bucket eviction logic.
- Miscalculated burn rates under zero traffic or edge-case window boundaries.
