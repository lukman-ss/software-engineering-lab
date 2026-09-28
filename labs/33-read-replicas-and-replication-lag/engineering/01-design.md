# Engineering Design

Target Lab: `labs/33-read-replicas-and-replication-lag`
Research Status: `APPROVED`

## Concept To Prove
Under asynchronous master-replica database replication, replica reads are eventually consistent and susceptible to replication lag. To prevent stale reads after writes while maintaining read scale, applications can employ:
1. **Read-Your-Own-Writes via LSN / Causal Token Tracking**: Tracking write commit LSN / position token and verifying replica catch-up before serving reads (or waiting / falling back to primary).
2. **Time-Based Sticky Routing**: Routing reads to the primary for a configurable window ($T_{sticky}$) following a write by that session.
3. **Lag-Aware Dynamic Load Balancing**: Routing generic reads only to replicas whose replication lag is within an acceptable SLA threshold, routing to primary or failing over if all replicas exceed lag threshold.
4. **Replication Mode Simulation**: Asynchronous vs Synchronous (`remote_apply`) replication semantics.

## Expected Behavior
- **Asynchronous Replication**: Writes commit immediately on primary; replication worker streams WAL records with configurable network/processing delay.
- **Direct Replica Reads**: Can return stale data or not-found errors if query hits lagging replica before WAL is applied.
- **Session / Token-Based Reads (Read-Your-Own-Writes)**: Reads specifying minimum required LSN either select a replica that has caught up, wait with timeout for replica catch-up, or route to primary, guaranteeing freshness.
- **Sticky Session Routing**: Reads occurring within `stickyDuration` of a write by the same session route directly to primary; after expiry, reads route to replica pool.
- **Lag Monitoring & Filtering**: Replicas exceeding max allowable lag (in LSN difference or delay) are excluded from replica read balancing.

## Failure Scenario
- A user updates their profile or creates an item on primary, then immediately requests a read routed to a lagging replica without session guarantees (stale read anomaly / missing record).
- Replicas experiencing high lag or failure without router failover causing degraded availability or unbounded stale reads.

## Success Criteria
- Deterministic simulation of master-replica architecture with configurable lag.
- Demonstrated occurrence of stale read anomaly under naive round-robin read-split.
- Demonstrated prevention of stale reads using LSN session token routing and time-based sticky routing.
- Demonstration of synchronous replication (`remote_apply`) latency vs freshness trade-off.
- 100% passing tests with Go race detector enabled.

## Architecture
```
                   +-----------------------------------+
                   |           Client / App            |
                   +-----------------+-----------------+
                                     |
                                     v
                   +-----------------------------------+
                   |       Dynamic Query Router        |
                   | - Read/Write Split                |
                   | - Session / LSN Token Tracker     |
                   | - Sticky Routing Table            |
                   | - Replica Health & Lag Monitor    |
                   +---------+---------------+---------+
                             |               |
               (Writes /     |               | (Lag-checked
             Primary Reads)  |               |  Replica Reads)
                             v               v
                   +---------------+   +---------------+
                   | Primary Node  |   | Replica Nodes |
                   | (WAL master)  |   | (Async / Sync)|
                   +-------+-------+   +-------+-------+
                           |                   ^
                           +--- WAL Stream ----+
                                (with Lag)
```

## Components
1. `cluster.Cluster`: In-memory multi-node database cluster with a primary and $N$ replicas, supporting sync/async replication, configurable lag, WAL streaming, and LSN tracking.
2. `cluster.Node`: Individual node storing key-value pairs with last committed LSN and applied LSN.
3. `router.Router`: Smart client router implementing read/write splitting, session LSN tracking (`ReadWithToken`), sticky post-write routing (`ReadWithStickySession`), and replica lag filtering.
4. `cmd/demo`: Interactive demonstration CLI displaying async replication lag anomaly, sticky routing, LSN token wait/routing, and synchronous replication comparison.

## Test Strategy
- Unit & integration tests in `tests/`:
  - `TestNaiveReplicationLag_StaleRead`: Proves stale read occurs when reading from lagging replica.
  - `TestReadYourOwnWrites_LSNToken`: Proves session token routing waits/routes to caught-up replica or primary.
  - `TestStickySessionRouting`: Proves time-based sticky routing directs reads to primary within TTL, and replicas after TTL.
  - `TestReplicaLagThreshold_Fallback`: Proves queries bypass replicas that exceed configured lag SLA.
  - `TestSynchronousReplication_Freshness`: Proves sync replication ensures immediate visibility on replica at the cost of write time.
  - `TestConcurrentAccess_RaceFree`: Concurrency tests under heavy read/write load with race detection.

## Execution Plan
1. Implement `go.mod` (module `labs/33-read-replicas-and-replication-lag`).
2. Implement `internal/cluster` and `internal/router`.
3. Implement `cmd/demo/main.go`.
4. Implement comprehensive test suite in `tests/`.
5. Run tests, race detector, demo, and verify zero errors/warnings.
6. Generate engineering docs (`02-implementation-notes.md`, `03-execution-result.md`, `README.md`).

## Implementation Decisions
- **In-Memory Cluster Simulation**: Implemented with pure Go standard library (`sync`, `time`, `context`) to allow high-precision deterministic replication lag simulation without requiring heavyweight external Docker DB instances during automated unit testing.
- **LSN Monotonic Ordering**: LSNs are represented as `uint64` monotonically increasing numbers corresponding to WAL log offsets.
- **Configurable Sticky Duration**: Sticky session window is exposed as a configurable parameter (`time.Duration`), defaulting to 5s per research recommendations.
