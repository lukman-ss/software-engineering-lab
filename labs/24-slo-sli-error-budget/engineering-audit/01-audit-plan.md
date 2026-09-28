# Engineering Audit Plan

Target Lab: labs/24-slo-sli-error-budget
Implementation Files:
- `internal/metrics/tracker.go`
- `internal/slo/evaluator.go`
- `internal/alerting/engine.go`
- `cmd/demo/main.go`

Tests:
- `tests/slo_test.go`

Executable/Demo:
- `cmd/demo/main.go`

Approved Research Inputs:
- `research/05-report.md`
- `research-audit/07-verdict.md`
- `engineering/01-design.md`
- `engineering/02-implementation-notes.md`

Main Claims To Verify:
1. SLI is evaluated as a quantitative ratio of good requests divided by total requests (`good / total`).
2. Error Budget is calculated based on `(1 - SLO) * total` events and correctly decrements upon bad events.
3. Multi-window multi-burn-rate alerting triggers when both short and long windows exceed defined burn rate thresholds.
4. Deployment gating policy (`CanDeploy`) halts releases when remaining error budget is depleted (`<= 0`).
5. Concurrency safety is upheld across all metric tracking components under parallel goroutine traffic.
6. Execution outputs in `engineering/03-execution-result.md` match live command outputs.

Commands To Run:
- `go test -v -count=1 ./...`
- `go test -race -count=1 ./...`
- `go run ./cmd/demo`

Primary Risks:
- Race conditions or data races in sliding window state updates.
- Inaccurate time-bucket aggregation or out-of-order event handling anomalies.
- Incorrect burn rate math causing false positives or missing alerts.
- Discrepancies between demo execution output and documented execution results.
