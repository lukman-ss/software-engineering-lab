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
- `research/03-evidence.md`

Main Claims To Verify:
1. SLI Calculation as Good/Total ratio with sliding-window time-bucketed event tracking.
2. Error Budget calculation (`1 - SLO`) and remaining budget tracking with deployment freeze enforcement when budget <= 0.
3. Multi-window multi-burn-rate alerting logic detecting threshold breaches (fast burn vs slow burn).
4. Thread-safety under concurrent metric recording.
5. README accuracy against implemented code and execution output.

Commands To Run:
- `go test -count=1 -v ./...`
- `go test -count=1 -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Mathematical rounding or window calculation errors in SLI / Error budget evaluation.
- Multi-window burn-rate calculation mismatch (e.g. short vs long window tracker inputs).
- Discrepancies between demo output / design claims and code implementation.
