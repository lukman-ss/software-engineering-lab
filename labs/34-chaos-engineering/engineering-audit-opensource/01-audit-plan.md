# Engineering Audit Plan

Target Lab: labs/34-chaos-engineering
Implementation Files:
- internal/fault/injector.go
- internal/circuitbreaker/circuitbreaker.go
- internal/monitor/monitor.go
- internal/experiment/runner.go
- cmd/demo/main.go
Tests:
- tests/chaos_test.go
Executable/Demo:
- cmd/demo/main.go (run with `go run ./cmd/demo`)
Approved Research Inputs:
- research/*.md (plan, evidence, report)
Main Claims To Verify:
1. Fault injection latency and forced error work.
2. Circuit breaker state transitions and fallback.
3. Monitor tracks error rate and triggers auto‑abort.
4. Demo reflects these behaviors.
Commands To Run:
- `go test ./...`
- `go test -race ./...`
- `go run ./cmd/demo`
Primary Risks:
- Concurrency bugs hidden from tests.
- Metric calculation simplified (no sliding window).
- Demo output could be fabricated.
