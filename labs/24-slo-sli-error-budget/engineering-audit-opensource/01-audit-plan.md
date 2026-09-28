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
Approved Research Inputs: (research files omitted per pipeline)
Main Claims To Verify:
- SLI = good/total ratio
- Error budget calculation and CanDeploy policy
- Multi‑window burn‑rate alerting (fast & slow thresholds)
- Concurrency safety of WindowTracker
- Demo output reflects real calculations
Commands To Run:
- go test ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Concurrency race conditions
- Eviction correctness
- Division‑by‑zero handling
- Rounding masking negative budget
