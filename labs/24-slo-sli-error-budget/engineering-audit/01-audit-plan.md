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
- `engineering/01-design.md`
- `engineering/02-implementation-notes.md`
- `engineering/03-execution-result.md`
Main Claims To Verify:
1. SLI is evaluated as a quantitative `good_events / total_events` ratio across sliding time windows.
2. Error budget is calculated as `(1 - SLO) * total_events - bad_events`, and triggers deployment freeze (`CanDeploy = false`) when depleted (`<= 0`).
3. Multi-window multi-burn-rate alerting triggers only when both short and long rolling windows breach burn-rate thresholds simultaneously.
4. Concurrency safety across concurrent metric ingestion and evaluations.
5. Robust handling of out-of-order timestamps and zero traffic states.
6. Execution outputs in demo match mathematical models and no fabricated logs/benchmarks exist.

Commands To Run:
```bash
go test -v ./...
go test -race ./...
go run ./cmd/demo
```

Primary Risks:
- Race conditions or deadlocks under concurrent metric ingestion and window eviction.
- Slicing and pointer mutations during sorted insertion or eviction in `tracker.go`.
- Floating-point precision truncation causing false deployments or false alerts.
- Discrepancies between demo console output claims and code-evaluated values.
