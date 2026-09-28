# Engineering Audit Plan

Target Lab: labs/24-slo-sli-error-budget
Implementation Files:
- internal/metrics/tracker.go
- internal/slo/evaluator.go
- internal/alerting/engine.go
- cmd/demo/main.go

Tests: tests/slo_test.go

Executable/Demo: cmd/demo/main.go

Approved Research Inputs: research/, research-audit/ (to be inspected but not part of pipeline override)

Main Claims To Verify:
1. Correct SLI/SLO calculation and error budget tracking
2. Multi-window burn-rate alerting (fast/slow) triggers correctly
3. Concurrency safety in WindowTracker
4. Release freeze policy based on exhausted budget
5. Demo matches described behavior

Commands To Run:
- go test ./...
- go test -race ./...
- go run ./cmd/demo

Primary Risks:
- Race conditions in bucket insertion logic (ordered by timestamp)
- Incorrect budget remaining sign handling
- Alert logic requiring both windows to exceed factor simultaneously
- Zero traffic edge cases