# Engineering Audit Plan

Target Lab: `labs/33-read-replicas-and-replication-lag`
Implementation Files:
- `internal/cluster/cluster.go`
- `internal/router/router.go`
- `cmd/demo/main.go`
- `go.mod`

Tests:
- `tests/replication_test.go`

Executable/Demo:
- `cmd/demo/main.go`

Approved Research Inputs:
- `research/05-report.md`
- `research-audit/07-verdict.md` (Research verdict: `APPROVED`)

Main Claims To Verify:
1. Asynchronous replication produces observable stale read anomalies under naive replica read routing.
2. Read-Your-Own-Writes is guaranteed via session LSN causal token tracking (`ReadWithToken`) by waiting for replica WAL catch-up or falling back to primary.
3. Read-Your-Own-Writes is guaranteed via time-based sticky routing (`ReadWithStickySession`) during the sticky TTL window.
4. Lag-aware routing dynamically detects replicas exceeding acceptable LSN lag thresholds and routes queries to fresh replicas or primary fallback.
5. Synchronous replication (`SyncReplication` / `remote_apply`) eliminates replica lag visibility anomalies at the expense of higher write latency.
6. Concurrent read and write workloads execute cleanly without data races under Go race detector.

Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Race conditions or deadlocks in cluster synchronization and replica condition variables (`sync.Cond`, `sync.RWMutex`, `atomic.Uint64`).
- Goroutine leaks during `WaitForLSN` or cluster shutdown (`Cluster.Close`).
- Mismatches between documentation, research findings, and actual code implementation.
