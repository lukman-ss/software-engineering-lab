# Engineering Audit Plan

Target Lab: labs/15-load-testing
Implementation Files:
- `internal/server/server.go`
- `internal/loadtest/metrics.go`
- `internal/loadtest/runner.go`
- `cmd/demo/main.go`
- `go.mod`
Tests:
- `internal/loadtest/metrics_test.go`
- `tests/loadtest_test.go`
Executable/Demo:
- `cmd/demo/main.go`
Approved Research Inputs:
- `labs/15-load-testing/research/05-report.md`
- `labs/15-load-testing/research-audit/07-verdict.md`
Main Claims To Verify:
1. Smoke load with low VUs within capacity yields baseline latency with low tail deviation (P95 ≈ Avg).
2. Stress load exceeding capacity causes queue saturation and tail latency explosion (P95, P99 >> Avg).
3. Metric calculations (Min, Max, Avg, P50, P90, P95, P99, RPS) are mathematically sound and handle edge cases (empty inputs, zero duration).
4. Concurrency runner operates safely under high VUs without race conditions or deadlocks.
5. Error accounting tracks HTTP 4xx/5xx and dial/network errors correctly without miscounting context cancellations.
Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`
Primary Risks:
- Race conditions during concurrent metric aggregation or request execution.
- Flaky tests if assertion bounds are too narrow under variable CPU load.
- Inaccurate percentile computation or slice index out-of-bounds on edge cases.
