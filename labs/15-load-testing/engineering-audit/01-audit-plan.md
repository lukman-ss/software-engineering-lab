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
- `research/01-plan.md`
- `research/02-sources.md`
- `research/03-evidence.md`
- `research/04-contradictions.md`
- `research/05-report.md`
- `research/06-open-questions.md`

Main Claims To Verify:
1. Load tester correctly measures throughput (RPS), total requests, error counts, and latency percentiles (P50, P90, P95, P99).
2. Stress testing against resource-constrained server (DB connection pool semaphore limit) causes queueing and tail latency explosion (P95/P99 latency growth).
3. Concurrency runner correctly handles context cancellation, connection pool reuse tuning, and thread-safe metric aggregation without data races.
4. Documentation (README.md) accurately reflects implementation, test suite execution, and structure.

Commands To Run:
- `go test -v ./...`
- `go test -race -v ./...`
- `go run ./cmd/demo`

Primary Risks:
- Data race conditions in concurrent request execution and per-VU latency aggregation.
- Inaccurate percentile calculation logic (off-by-one or sorting errors).
- Mismatch between mock server queueing/tail latency simulation and reported metrics.
- Unhandled HTTP request context cancellation or goroutine leaks.
