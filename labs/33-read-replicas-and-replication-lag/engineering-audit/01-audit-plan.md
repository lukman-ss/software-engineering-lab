# Engineering Audit Plan

Target Lab: `labs/33-read-replicas-and-replication-lag`
Implementation Files:
- `internal/cluster/cluster.go`
- `internal/router/router.go`
- `cmd/demo/main.go`
Tests:
- `tests/replication_test.go`
Executable/Demo:
- `cmd/demo/main.go`
Approved Research Inputs:
- `labs/33-read-replicas-and-replication-lag/research/05-report.md`
- `labs/33-read-replicas-and-replication-lag/research-audit/07-verdict.md`
Main Claims To Verify:
1. Asynchronous replication produces replication lag and stale reads when reading directly from replicas before WAL is applied.
2. Read-your-own-writes session guarantee is achieved via time-based sticky routing to primary within TTL window.
3. Read-your-own-writes session guarantee is achieved via causal token / minimum LSN tracking and replica catch-up waiting.
4. Lag-aware dynamic routing excludes replicas exceeding maximum acceptable LSN lag and falls back to primary.
5. Synchronous replication (`remote_apply`) guarantees immediate visibility on replicas at the cost of higher write latency.
6. Concurrency safety under concurrent reads and writes with Go race detector.
7. README, engineering design, implementation notes, and execution results accurately reflect the code.
Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`
Primary Risks:
- Thread-safety / race conditions during replica WAL propagation and conditional wait (`sync.Cond`).
- Goroutine leak or deadlock on `WaitForLSN` if replica channel is closed or context cancels.
- Inconsistencies between README documentation and router / cluster APIs.
- Fabricated demo or test output.
