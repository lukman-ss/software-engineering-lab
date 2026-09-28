# Implementation Notes

## Files Added
- `go.mod`: Module definition (`labs/30-leader-election`, Go 1.22).
- `internal/coordinator/coordinator.go`: Thread-safe coordinator managing linearizable leases, monotonic fencing tokens, and TTL expirations.
- `internal/candidate/candidate.go`: Node election participant handling campaign lifecycle, heartbeats, simulated GC pauses, and fenced writes.
- `internal/storage/storage.go`: Shared storage validating fencing tokens (`token > lastSeenToken`) to prevent split-brain writes.
- `tests/election_test.go`: Unit and integration test suite covering lease acquisition, expiration, fencing validation, failover under GC pause, and race condition resistance.
- `cmd/demo/main.go`: Interactive runnable CLI demo showcasing leader election, failover on GC pause, and safe rejection of stale leader writes.
- `README.md`: Concise operational documentation.

## Core Design Decisions
- **Monotonic Fencing Tokens**: Uses an atomically incrementing `revision` counter representing term/grant epoch.
- **Resource-Layer Fencing**: Fencing check is performed atomically inside `FencedStorage.Write` via check-and-set semantics.
- **Pure Go Standard Library**: All synchronization uses standard `sync.Mutex` and `sync.RWMutex`, with zero external dependencies.

## Implementation-Specific Choices
- **Simulated GC Pause**: `Node.SimulatePause` bypasses heartbeat loop ticks until a specified future monotonic timestamp, simulating stop-the-world garbage collection or client network disconnection.
- **In-Memory Coordinator**: Used an in-memory coordinator rather than spawning external etcd or ZooKeeper processes, preserving zero-dependency reproducibility.

## Known Limitations
- Coordinator is single-instance in-memory; does not distribute coordinator state across multiple Raft replicas.
- Network transport is simulated in-process via memory calls rather than gRPC/TCP.

## Trade-offs
- Simplicity and test velocity vs. full distributed network fault injection.
- Monotonic in-memory counter vs. cluster-wide Raft log index.

## What Is Demonstrated
- Lease acquisition and automatic renewal (`KeepAlive`).
- Bounded failure detection via lease TTL expiration.
- Promotion of standby candidate upon leader pause/failure.
- Protection against split-brain writes when a paused leader resumes using an expired token.

## What Is Not Demonstrated
- Byzantine faults or clock skew jumps on distributed coordinator nodes.
- Storage disk persistence / write-ahead logging.
