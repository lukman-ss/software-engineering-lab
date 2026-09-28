# Engineering Audit Plan

Target Lab: labs/24-slo-sli-error-budget
Implementation Files:
- internal/metrics/tracker.go
- internal/slo/evaluator.go
- internal/alerting/engine.go
- cmd/demo/main.go
Tests: tests/slo_test.go
Executable/Demo: go run ./cmd/demo
Approved Research Inputs: engineering/01-design.md, engineering/02-implementation-notes.md (research claims summarized)
Main Claims To Verify:
- SLI calculation as good/total ratio
- Error budget = (1-SLO)*total; budget consumed = failures; remaining = allowed - consumed
- Multi-window burn-rate alerting (fast/slow windows vs thresholds)
- Release freeze policy: CanDeploy = true if budget > 0, false if exhausted
- Endpoint criticality: stricter SLO for critical services (e.g. Payment 99.9% vs Reports 95.0%)
- Concurrency safety under load
- Demo illustrates baseline, incident, alerting
Commands To Run:
- go build ./...
- go test ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Floating-point boundary fragility at exact budget exhaustion
- Unused Config fields and rule window parameters (design vs implementation divergence)
- Missing recovery demonstration (claimed in design doc)