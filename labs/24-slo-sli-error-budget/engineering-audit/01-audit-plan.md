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
- `research-audit/07-verdict.md` (APPROVED)

Main Claims To Verify:
1. SLI calculation implements standard ratio model: `good_events / total_events`.
2. Error budget dynamically tracks `(1 - SLO) * total_events` and remaining budget correctly triggers deployment freeze (`CanDeploy = false`).
3. Multi-window multi-burn-rate alerting monitors both short and long rolling windows before alerting, preventing false positives from transient spikes.
4. Concurrency safety of `WindowTracker` under parallel write/read workloads.
5. Realistic demo simulating baseline traffic, severe incident budget burn, multi-window burn rate alert triggering, and endpoint criticality comparison.

Commands To Run:
- `go test -v -count=1 ./...`
- `go test -race -v -count=1 ./...`
- `go run ./cmd/demo`

Primary Risks:
- Data races during window eviction and bucket updates under concurrent request ingestion.
- Arithmetic edge cases (zero traffic division by zero, float precision errors in budget subtraction).
- Discrepancy between demo output and recorded execution results.
- Unhandled out-of-order timestamp event insertion.
