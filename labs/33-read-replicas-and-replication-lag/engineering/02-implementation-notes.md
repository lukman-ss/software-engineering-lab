# Implementation Notes

## Files Added
- `go.mod`: Module definition for `labs/33-read-replicas-and-replication-lag` targeting Go 1.22+.
- `internal/cluster/cluster.go`: In-memory multi-node primary/replica database simulation with WAL streaming, configurable replication lag, LSN progression, sync vs async modes, and conditional broadcast notifications.
- `internal/router/router.go`: Smart client router supporting naive read splitting, time-based sticky session routing, LSN causal token tracking (`ReadWithToken`), and lag-aware routing with primary fallback.
- `cmd/demo/main.go`: End-to-end runnable demonstration executing realistic write-then-read flows across naive routing, sticky routing, LSN token waiting, and synchronous replication comparison.
- `tests/replication_test.go`: Complete automated test suite testing stale read anomalies, sticky session routing, causal LSN waiting, lag SLA threshold fallback, synchronous freshness, and race detector concurrency.
- `README.md`: Clear lab documentation describing architecture, components, and execution commands.
- `engineering/01-design.md`: Engineering design document.
- `engineering/02-implementation-notes.md`: Implementation design decisions, trade-offs, and boundaries.
- `engineering/03-execution-result.md`: Recorded build, test, race detector, and demo execution outputs.

## Core Design Decisions
1. **In-Memory Cluster Emulation**:
   - Built a lightweight concurrency-safe in-memory database simulation using Go's standard library (`sync.RWMutex`, `sync.Cond`, channels).
   - Allows sub-millisecond precision testing of WAL delivery lag without third-party database daemon overhead or flaky network environments.
2. **LSN Representation as Monotonic `uint64`**:
   - Modeled Log Sequence Numbers (LSN) as monotonically increasing 64-bit unsigned integers representing commit order and log offset.
3. **Session Guarantees**:
   - `ReadWithStickySession`: Tracks session last-write timestamp in thread-safe map. Routes reads to primary while `now - lastWrite < StickyDuration`; routes to replica pool once expired.
   - `ReadWithToken`: Accepts client causal token (`minLSN`). First checks if any replica has already applied the LSN; if not, blocks with timeout on a replica's conditional variable (`sync.Cond`) waiting for the WAL entry, falling back to primary on timeout.
4. **Lag-Aware Replica Filtering**:
   - Evaluates replica staleness using $\Delta LSN = LSN_{primary} - LSN_{replica}$. Replicas exceeding `MaxLSNDiff` are excluded from read load balancing, triggering primary fallback if all replicas lag.

## Implementation-Specific Choices
- Default sticky session duration exposed in `router.Config` and defaults to 5s per research recommendations, overridden to shorter values in automated tests for fast deterministic runs.
- `sync.Cond` used for event-driven replica catch-up notification instead of busy polling.

## Known Limitations
- WAL entry payload is a key-value store, not full relational SQL schemas or multi-statement transactions.
- In-memory node states do not persist to disk across process restarts.

## Trade-offs
- Synchronous replication (`SyncReplication` / `remote_apply`) guarantees immediate read freshness on all nodes, but adds write latency proportional to replica apply lag.
- Token-based waiting trades read latency (waiting for replica catch-up) to preserve read offloading from the primary.

## What Is Demonstrated
- Stale read anomaly under naive asynchronous read replica routing.
- Read-your-own-writes guarantee via time-based sticky routing to primary.
- Read-your-own-writes guarantee via LSN causal token tracking and replica catch-up wait.
- Dynamic lag filtering protecting queries against lagging replicas.
- Latency vs consistency trade-offs in asynchronous vs synchronous replication.

## What Is Not Demonstrated
- Cascading replication topologies (replica of replica).
- Network partitioning split-brain elections (handled by consensus systems like Raft/Paxos).
