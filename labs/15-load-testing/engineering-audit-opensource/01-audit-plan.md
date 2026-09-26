# Engineering Audit Plan

Target Lab: `labs/15-load-testing` (module `github.com/lukman/software-engineering-lab/labs/15-load-testing`)
Implementation Files:
- `internal/loadtest/runner.go`
- `internal/loadtest/metrics.go`
- `internal/server/server.go`
- `cmd/demo/main.go`
Tests:
- `internal/loadtest/metrics_test.go`
- `tests/loadtest_test.go`
Executable/Demo:
- `cmd/demo` (`go run ./cmd/demo`)
Approved Research Inputs:
- `research/01-plan.md`, `02-sources.md`, `03-evidence.md`, `04-contradictions.md`, `05-report.md` (research phase; implementation alignment verified, not content depth)
- `engineering/01-design.md`, `02-implementation-notes.md`
Main Claims To Verify:
1. `go build ./...` compiles cleanly.
2. `go test ./...` passes (smoke/stress, error path, method-not-allowed, dial error, metrics).
3. `go test -race ./...` race-free.
4. Demo (`go run ./cmd/demo`) runs and produces realistic latency divergence between smoke (≤ DB connections) and stress (≫ DB connections) load.
5. Stress P95 latency exceeds smoke P95 (proves server queuing behavior is measurable).
6. Metrics (RPS, percentiles, error count) computed correctly and match observed output.
7. README matches the implementation and demo invocation.
Commands To Run:
- `go build ./...`
- `go test -race ./...`
- `go run ./cmd/demo`
Primary Risks:
- Latency math: P95/P99 use index = `(len-1)*pct/100`; verify boundary expectations in tests.
- Concurrency: per-VU result slices written only by owning goroutine; aggregate after `wg.Wait` — must be race-free.
- Stress-vs-smoke ordering assertion depends on server DB-connection semaphore actually throttling at 50 VUs.
- Demo output timing-dependent; re-run to confirm reproducibility of latency divergence.
