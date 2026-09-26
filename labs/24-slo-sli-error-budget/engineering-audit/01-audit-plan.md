# Engineering Audit Plan

Target Lab: labs/24-slo-sli-error-budget
Implementation Files:
- internal/metrics/tracker.go
- internal/slo/evaluator.go
- internal/alerting/engine.go
- cmd/demo/main.go

Tests:
- tests/slo_test.go

Executable/Demo:
- cmd/demo/main.go

Approved Research Inputs:
- research/05-report.md
- research/03-evidence.md

Main Claims To Verify:
1. Sliding-window time-bucketed event tracker handles event recording and stale bucket eviction accurately.
2. SLO Evaluator calculates SLI ratios, error budgets, and enforces release freeze (`CanDeploy=false`) when budget is exhausted.
3. Multi-window multi-burn-rate alerting engine correctly detects fast and slow burn rate threshold breaches.
4. Implementation is concurrency safe across concurrent metric recording.
5. README instructions compile and run cleanly (`go test ./...`, `go test -race ./...`, `go run ./cmd/demo`).

Commands To Run:
- `go test ./...`
- `go test -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Race conditions during concurrent `Record()` and `Summary()` calls on `WindowTracker`.
- Incorrect burn rate math or SLI calculation rounding logic.
- Misalignment between README documentation and actual Go package structure/outputs.
