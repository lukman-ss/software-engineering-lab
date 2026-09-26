# Engineering Audit Plan

Target Lab: `labs/15-load-testing`
Implementation Files:
- `internal/server/server.go`
- `internal/loadtest/runner.go`
- `internal/loadtest/metrics.go`
- `cmd/demo/main.go`

Tests:
- `internal/loadtest/metrics_test.go`
- `tests/loadtest_test.go`

Executable/Demo:
- `cmd/demo/main.go`

Approved Research Inputs:
- `research/05-report.md`
- `research-audit/07-verdict.md`
- `engineering/01-design.md`
- `engineering/02-implementation-notes.md`

Main Claims To Verify:
1. Server simulates resource constraint (connection pool capacity) via semaphore.
2. Runner drives concurrent virtual users (VUs) without internal mutex contention or race conditions.
3. Percentile computation (P50, P90, P95, P99) and summary metrics (RPS, min, max, avg) are mathematically sound.
4. Smoke test under capacity shows near-baseline latency (~20ms), while Stress test over capacity demonstrates queueing and tail latency spikes (~200ms+).
5. Error accounting tracks non-2xx status codes and network errors while filtering out benign context cancellation upon test completion.
6. Execution results in `engineering/03-execution-result.md` reflect actual reproducible runs.

Commands To Run:
```bash
go test -v ./...
go test -race ./...
go run ./cmd/demo
```

Primary Risks:
- Race conditions during concurrent metrics collection or request dispatching.
- Connection leaks or port exhaustion under stress load due to HTTP client configuration.
- Flaky latency assertions dependent on OS scheduling jitters.
