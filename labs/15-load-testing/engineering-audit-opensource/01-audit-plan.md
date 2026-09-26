# Engineering Audit Plan

Target Lab: labs/15-load-testing
Implementation Files: 
- internal/loadtest/runner.go
- internal/loadtest/metrics.go
- internal/server/server.go
- cmd/demo/main.go
Tests: 
- internal/loadtest/metrics_test.go
- tests/loadtest_test.go
Executable/Demo: cmd/demo/main.go
Approved Research Inputs: 
- engineering/01-design.md (APPROVED)
- engineering/02-implementation-notes.md
- engineering/03-execution-result.md
Main Claims To Verify:
1. Load generator computes Min, Max, Average, P50, P90, P95, P99 latencies accurately.
2. Smoke test (low VUs) shows low latency and zero errors.
3. Stress test (high VUs) shows increased tail latency (P95, P99) due to resource exhaustion.
4. No race conditions in concurrent load generation.
5. Demo output matches expected behavior (smoke vs stress comparison).
Commands To Run:
- go build ./...
- go test -v ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Race condition in metrics aggregation (mitigated by per-VU slicing).
- Incorrect percentile calculation due to off-by-one errors.
- Demo may not reflect actual stress if VUs insufficient to saturate resources.
- Test may be flaky due to timing dependencies.