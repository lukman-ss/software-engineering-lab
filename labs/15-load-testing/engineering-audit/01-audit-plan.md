# Engineering Audit Plan

Target Lab: labs/15-load-testing
Implementation Files:
- `internal/loadtest/metrics.go`
- `internal/loadtest/runner.go`
- `internal/server/server.go`

Tests:
- `internal/loadtest/metrics_test.go`
- `tests/loadtest_test.go`

Executable/Demo:
- `cmd/demo/main.go`

Approved Research Inputs:
- `research/05-report.md`
- `research-audit/07-verdict.md` (APPROVED)

Main Claims To Verify:
1. Load test runner simulates VUs concurrently and records latencies/errors.
2. Metrics module calculates P50, P90, P95, and P99 percentiles correctly.
3. Server bottleneck simulation via Semaphore accurately reflects connection pool contention and latency degradation.
4. Concurrency safety (race detector clean).
5. README claims match executable demo and test execution instructions.

Commands To Run:
- `go test ./...`
- `go test -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Race conditions during concurrent slice appending or metric calculation in loadtest runner.
- Miscalculation of percentiles (off-by-one or non-sorted input handling).
- Discrepancy between README documented behavior and actual runtime results.
