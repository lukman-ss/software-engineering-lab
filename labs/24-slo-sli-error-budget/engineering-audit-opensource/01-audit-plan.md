# Engineering Audit Plan

Target Lab: labs/24-slo-sli-error-budget
Implementation Files: 
- internal/metrics/tracker.go
- internal/slo/evaluator.go
- internal/alerting/engine.go
- cmd/demo/main.go
Tests: tests/slo_test.go
Executable/Demo: cmd/demo/main.go
Approved Research Inputs: research/ directory (not audited in this stage)
Main Claims To Verify:
1. SLI calculation as good/total ratio
2. Error budget calculation and management
3. Burn rate alerting logic (multi-window)
4. Endpoint criticality with different SLOs
5. Thread safety and concurrency handling
Commands To Run:
- go test ./...
- go test -race ./...
- go test ./...
- go run ./cmd/demo
Primary Risks:
- Race conditions in WindowTracker due to concurrent map/slice access
- Incorrect error budget math leading to false budget remaining
- Burn rate calculation not matching SLO specification
- Demo output not reflecting actual internal state