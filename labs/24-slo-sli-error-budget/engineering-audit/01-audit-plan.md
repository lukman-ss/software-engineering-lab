# Engineering Audit Plan

Target Lab: `labs/24-slo-sli-error-budget`
Implementation Files:
- `internal/metrics/tracker.go`
- `internal/slo/evaluator.go`
- `internal/alerting/engine.go`
- `go.mod`
Tests:
- `tests/slo_test.go`
Executable/Demo:
- `cmd/demo/main.go`
Approved Research Inputs:
- `research/runs/2026-09-26-slo-sli-error-budget/05-report.md`
- `research-audit/07-verdict.md` (APPROVED)
Main Claims To Verify:
1. Ratio-based SLI calculation (`good_events / total_events`).
2. Error budget management (`(1.0 - SLO) * total_events - bad_events`) and release policy gating.
3. Multi-window multi-burn-rate alerting logic evaluation.
4. Concurrency thread safety under parallel request ingestion.
5. Exact matching of real execution output with recorded execution docs and README.
Commands To Run:
- `go test ./...`
- `go test -race ./...`
- `go run ./cmd/demo`
Primary Risks:
- Race conditions in metrics bucket window eviction or slice resizing.
- Inconsistencies in floating-point calculations or rounding for budget/SLI.
- Divergence between claims in README/design vs implementation code.
