# Engineering Audit Plan

Target Lab: labs/15-load-testing
Implementation Files:
- `internal/server/server.go`
- `internal/loadtest/runner.go`
- `internal/loadtest/metrics.go`
Tests:
- `internal/loadtest/metrics_test.go`
- `tests/loadtest_test.go`
Executable/Demo:
- `cmd/demo/main.go`
Approved Research Inputs:
- `research/05-report.md`
- `engineering/01-design.md`
- `engineering/02-implementation-notes.md`
- `engineering/03-execution-result.md`
Main Claims To Verify:
1. Smoke vs Stress differentiation: Smoke (low VUs) exhibits low latency and zero queuing; Stress (high VUs) saturates the mock DB connection pool and drives non-linear tail latency inflation (P95/P99).
2. Percentile accuracy: `CalculateMetrics` correctly derives Min, Max, Average, P50, P90, P95, and P99 latencies and RPS across sample sets.
3. Concurrency safety: Load runner executes parallel VUs and aggregates metrics without race conditions (`go test -race ./...`).
4. Resource bound enforcement: Mock server constrains concurrent DB operations to `MaxDBConnections` via a semaphore channel.
5. Error handling and propagation: Network failure/dial error, non-2xx status codes, and HTTP method rejections are counted and recorded accurately without panic or lockup.
6. Honest demo execution: `cmd/demo/main.go` dynamically exercises the mock server and prints genuine runtime statistics matching documented behaviors.
Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`
Primary Risks:
- Race conditions during concurrent metric collection across goroutines.
- Flaky tests caused by nondeterministic timing or jitter under heavy stress load.
- Inaccurate percentile calculation logic (off-by-one or non-monotonic ordering).
- Unhandled request timeouts or connection leaks in custom HTTP client transports.
