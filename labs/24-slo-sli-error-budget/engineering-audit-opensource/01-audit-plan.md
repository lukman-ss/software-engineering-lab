# Engineering Audit Plan

Target Lab:
labs/24-slo-sli-error-budget

Implementation Files:
- internal/metrics/tracker.go
- internal/slo/evaluator.go
- internal/alerting/engine.go
- cmd/demo/main.go
- tests/slo_test.go

Tests:
- go test ./...
- go test -race ./...

Executable/Demo:
- go run ./cmd/demo

Approved Research Inputs:
- SLI = good/total ratio (Google SRE Book, SRE Workbook)
- Error budget = 1 - SLO target (Google SRE Workbook, Datadog)
- Burn rate = actual error rate / allowed error rate (SRE Workbook Table 5-8)
- Multi-window burn-rate alerting fires when both short and long windows exceed threshold

Main Claims To Verify:
- SLI calculation accuracy across normal and incident traffic
- Error budget depletion correctly triggers CanDeploy=false
- Burn rate alert triggering with multi-window configuration
- Concurrency safety of WindowTracker under parallel Record calls
- Demo output matches claimed behavior (baseline, incident, alert, comparison)

Commands To Run:
- go test ./...
- go test -race ./...
- go run ./cmd/demo

Primary Risks:
- Race condition in WindowTracker bucket insertion under high concurrency
- SLI/budget off-by-one when total events change between Record and Evaluate
- Alert engine trigger logic requiring both windows to exceed threshold (may mask short-window spikes)