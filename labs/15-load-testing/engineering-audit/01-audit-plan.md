# Engineering Audit Plan

Target Lab: labs/15-load-testing
Implementation Files:
- internal/server/server.go
- internal/loadtest/runner.go
- internal/loadtest/metrics.go
Tests:
- internal/loadtest/metrics_test.go
- tests/loadtest_test.go
Executable/Demo:
- cmd/demo/main.go
Approved Research Inputs:
- research/01-plan.md
- research/05-report.md
- engineering/01-design.md
- engineering/02-implementation-notes.md
- engineering/03-execution-result.md
Main Claims To Verify:
- Load runner accurately measures latency percentiles (P50, P90, P95, P99) and RPS.
- Mock server bounds DB connection capacity via semaphore.
- Latency degrades predictably during saturation (Stress test P95 > Smoke test P95).
- Error counts and total requests strictly satisfy `TotalRequests == SuccessCount + ErrorCount`.
- Go test suite passes including `-race` detector check.
- Executable demo runs and outputs real comparative load metrics.
Commands To Run:
- `go test -v ./...`
- `go test -race -v ./...`
- `go run ./cmd/demo`
Primary Risks:
- Percentile calculation index bounds under small sample sizes.
- Race conditions during concurrent sample aggregation or server metrics polling.
- Inaccurate total request or error tracking during context timeout cancellation.
