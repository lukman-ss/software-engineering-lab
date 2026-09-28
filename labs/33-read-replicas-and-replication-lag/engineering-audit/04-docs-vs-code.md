# Docs vs Code Analysis

Target Lab: `labs/33-read-replicas-and-replication-lag`

## Comparison Matrix

| Claim / Documentation Item | Location in Docs | Code / Test Implementation | Alignment Status | Notes |
|---|---|---|---|---|
| Asynchronous vs Synchronous Replication (`remote_apply`) | `README.md:6`, `01-design.md:11` | `internal/cluster/cluster.go:98-103, 203-226` | MATCH | Implemented via `AsyncReplication` and `SyncReplication` modes. |
| Time-Based Sticky Routing | `README.md:8`, `01-design.md:9` | `internal/router/router.go:80-90` | MATCH | `ReadWithStickySession` routes to primary if within `StickyDuration`. |
| Causal Token / Minimum LSN Routing | `README.md:9`, `01-design.md:8` | `internal/router/router.go:92-115` | MATCH | `ReadWithToken` checks `AppliedLSN() >= minLSN` and blocks on `WaitForLSN`. |
| Lag-Aware Dynamic Routing & Fallback | `README.md:10`, `01-design.md:10` | `internal/router/router.go:117-140` | MATCH | Filters replicas where `primaryLSN - appliedLSN > MaxLSNDiff`, falls back to primary. |
| Demo Execution Output | `engineering/03-execution-result.md` | `cmd/demo/main.go` | MATCH | Recorded output matches real execution results closely (timing variance within normal range). |
| Race detector clean execution | `engineering/03-execution-result.md` | `tests/replication_test.go` | MATCH | `go test -race ./...` executes and passes with 0 races. |

## Discrepancies Found
- None. All claims documented in `README.md` and `engineering/01-design.md` correspond directly to functional code in `internal/cluster`, `internal/router`, and tests in `tests/replication_test.go`.
