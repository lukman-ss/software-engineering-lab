# Engineering Audit Plan

Target Lab: `labs/24-slo-sli-error-budget`
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

Main Claims To Verify:
1. Ratio-based SLI calculation (`good / total`) across time-windowed event tracker.
2. Error budget computation `(1 - SLO) * total` and depletion logic leading to release freeze policy enforcement (`CanDeploy = false`).
3. Multi-window multi-burn-rate alerting requiring both short and long window burn rates to breach threshold before firing.
4. Concurrency safety of metrics tracker under parallel reads and writes.
5. Out-of-order timestamp handling and bucket eviction correctness.
6. Real demo execution and output match documented execution.

Commands To Run:
- `go test ./...`
- `go test -race ./...`
- `go test -v -count=1 ./tests`
- `go run ./cmd/demo`

Primary Risks:
- Memory leaks or unbounded slice growth in metrics tracker.
- Race conditions during concurrent `Record` and `Summary` calls.
- Division by zero on zero traffic or 100% SLO target.
- False positive alerts on transient spikes if multi-window logic is flawed.
- Documentation divergence from actual executable output.
