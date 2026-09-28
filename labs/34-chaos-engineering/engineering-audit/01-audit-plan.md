# Engineering Audit Plan

Target Lab: labs/34-chaos-engineering
Implementation Files:
- `internal/fault/injector.go`
- `internal/circuitbreaker/circuitbreaker.go`
- `internal/monitor/monitor.go`
- `internal/experiment/runner.go`
Tests:
- `tests/chaos_test.go`
Executable/Demo:
- `cmd/demo/main.go`
Approved Research Inputs:
- `research/05-report.md` / `research/runs/2026-09-28-chaos-engineering/05-report.md`
- `research-audit/07-verdict.md` (APPROVED)
Main Claims To Verify:
- Fault injection supports controlled latency delays and forced error invocation.
- Circuit breaker transitions cleanly across Closed, Open, Half-Open states with graceful fallback execution.
- Steady state monitor tracks error rates accurately using atomic counters and flags health degradation.
- Chaos experiment runner monitors steady state metrics and automatically aborts / neutralizes active faults upon threshold breach.
- Blast radius is strictly constrained via synchronous fault neutralization on abort or termination.
- Zero race conditions under concurrent executions (`go test -race ./...`).
Commands To Run:
- `go test -v -count=1 ./...`
- `go test -race -v -count=1 ./...`
- `go run ./cmd/demo`
Primary Risks:
- Data races during fault mutation, circuit breaker state checks, or metric recording.
- Unhandled context cancellation during latency injection sleep.
- Inaccurate error rate calculations or false positives before sufficient request sample count.
- Lingering fault state if experiment abort or termination fails to clear the injector.
