# Engineering Audit Plan

Target Lab: labs/33-read-replicas-and-replication-lag
Implementation Files:
- internal/cluster/cluster.go
- internal/router/router.go
Tests:
- tests/replication_test.go
Executable/Demo:
- cmd/demo/main.go
Approved Research Inputs:
- engineering/01-design.md
- engineering/02-implementation-notes.md
- engineering/03-execution-result.md
Main Claims To Verify:
- Asynchronous replication lag causes stale reads under naive round-robin read routing.
- Time-based sticky session routing guarantees read-your-own-writes by routing to primary for configured duration.
- Causal token / LSN wait blocks or routes to catch-up node to ensure read freshness.
- Lag-aware routing filters out lagging replicas exceeding max LSN threshold and falls back to primary.
- Synchronous replication (`remote_apply`) guarantees replica freshness at the cost of write duration.
- Code compiles, tests pass, race detector passes, demo works, and no fake results exist.
Commands To Run:
- go test -count=1 -v ./...
- go test -count=1 -race ./...
- go run ./cmd/demo
Primary Risks:
- Race conditions or deadlocks in node synchronization (`sync.Cond`, `walChannel`, channel select drop).
- Non-deterministic channel buffer overflow behavior (`select default:` drop WAL entry).
- Unhandled context timeout / cancellation in `WaitForLSN` leading to leaked goroutines.
