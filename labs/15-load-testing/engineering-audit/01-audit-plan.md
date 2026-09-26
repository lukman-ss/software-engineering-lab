# Engineering Audit Plan

Target Lab: labs/15-load-testing
Implementation Files: `internal/server/server.go`, `internal/loadtest/runner.go`, `internal/loadtest/metrics.go`
Tests: `tests/loadtest_test.go`, `internal/loadtest/metrics_test.go`
Executable/Demo: `cmd/demo/main.go`
Approved Research Inputs: `research/` / `engineering/01-design.md`
Main Claims To Verify:
1. Load testing reveals boundaries/saturation.
2. P95/P99 latency spikes significantly over average in stress test.
3. Smoke test executes without queues; stress test exhausts resources.
4. Calculations for percentiles are accurate.
Commands To Run: `go test -v ./...`, `go test -race ./...`, `go run ./cmd/demo`
Primary Risks:
1. Concurrency bugs in metrics aggregation.
2. Server mock doesn't properly bound connections.
3. Tests don't properly assert latency relationships.
4. Runner fails to clean up connections.