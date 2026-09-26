# Engineering Audit Plan

Target Lab: labs/24-slo-sli-error-budget
Implementation Files:
- cmd/demo/main.go
- internal/metrics/tracker.go
- internal/slo/evaluator.go
- internal/alerting/engine.go
Tests: tests/slo_test.go
Executable/Demo: go run ./cmd/demo
Approved Research Inputs: Research status APPROVED per engineering/01-design.md
Main Claims To Verify:
1. SLI calculation as Good/Total ratio
2. Error Budget management (1 - SLO)
3. Multi-window multi-burn-rate alerting
4. Release freeze policy enforcement
Commands To Run:
- go test ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Concurrency safety in WindowTracker
- Correctness of error budget calculations under window eviction
- Alert logic requiring both short and long window thresholds