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
- `engineering/01-design.md`
- `engineering/02-implementation-notes.md`

Main Claims To Verify:
1. SLI calculation accurately computes `good / total` ratio.
2. Error budget calculation accurately reflects target SLO uptime allowance (`(1 - target) * total - bad`).
3. Release freeze policy dynamically blocks deployments (`CanDeploy = false`) when remaining budget <= 0.
4. Multi-window multi-burn-rate alert engine triggers correctly when short and long burn rate factors exceed thresholds.
5. Code compiles cleanly and passes Go race detector without data races.
6. Demo output is authentic and repeatable.
7. Documentation in README and engineering notes matches the implementation accurately.

Commands To Run:
- `go test -count=1 ./...`
- `go test -count=1 -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Mathematical rounding/floating point inaccuracies causing incorrect `CanDeploy` state or burn rate triggers.
- In-memory sliding window eviction bugs when timestamps arrive out-of-order.
- Concurrency race conditions in `WindowTracker` during simultaneous writes and summary reads.
- Discrepancy between README documentation / engineering notes and actual implementation logic.
