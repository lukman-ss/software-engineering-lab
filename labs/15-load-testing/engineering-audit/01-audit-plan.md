# Engineering Audit Plan

Target Lab: labs/15-load-testing
Implementation Files: `cmd/demo/main.go`, `internal/server/server.go`, `internal/loadtest/runner.go`, `internal/loadtest/metrics.go`
Tests: `internal/loadtest/metrics_test.go`, `tests/loadtest_test.go`
Executable/Demo: `cmd/demo/main.go`
Approved Research Inputs: `engineering/01-design.md`, `research/05-report.md`

Main Claims To Verify:
1. Load test runner correctly measures P50, P95, and P99 percentiles manually.
2. Under stress load, resource exhaustion (connection pool) causes queuing and latency degradation.
3. Average response time conceals tail latency spikes; percentiles (P95, P99) are necessary to uncover degradation.
4. Stress test average latency degrades less severely than P95 and P99.

Commands To Run:
```bash
go test -v ./...
go test -race ./...
go run ./cmd/demo
```

Primary Risks:
- Race conditions in metrics collection across multiple VUs.
- Math errors in percentile calculation.
- The load generator becoming the bottleneck instead of the server.
- The closed-loop workload model failing to demonstrate the "average masks P95" claim.
