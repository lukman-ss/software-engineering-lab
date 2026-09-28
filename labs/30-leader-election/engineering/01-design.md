# Engineering Design

Target Lab: `labs/30-leader-election`
Research Status: APPROVED

## Concept To Prove
1. Lease-based leader election with heartbeat renewals providing timed mutual exclusion and bounded failure detection.
2. Split-brain vulnerability under unbounded process pauses (e.g. simulated GC pause or network partition) where a stale leader erroneously acts after lease expiration.
3. Fencing token safety pattern: monotonically increasing generation/revision tokens enforced by the shared storage/resource write gate to reject stale leader writes.

## Expected Behavior
- Candidates contest leader election via an in-memory linearizable coordination service (coordinator).
- Coordinator grants a lease with a configured TTL (time-to-live) and issues a strictly monotonically increasing fencing token (`token = revision`).
- The leader actively renews the lease via periodic heartbeats (`KeepAlive`).
- When a leader stops renewing or pauses longer than TTL:
  - The lease expires.
  - A follower detects expiration and acquires leadership, receiving a higher fencing token (`token + 1`).
- Shared resource accepts writes with strictly greater fencing tokens (`token > lastSeenToken`) and records write history.

## Failure Scenario
- **Stale Leader / Split-Brain Attempt**: Leader A acquires lease with token 1.
- Leader A experiences a simulated pause (or disconnect) longer than the lease TTL.
- Coordinator expires Leader A's lease.
- Leader B campaigns, acquires the lease, receives token 2, and performs a write to the shared resource.
- Leader A wakes up, unaware its lease expired, and attempts to write to the shared resource using token 1.
- Without fencing, Leader A's write would overwrite Leader B's state (data corruption).
- With fencing, the shared resource rejects Leader A's write with `ErrStaleFencingToken`, preserving consistency.

## Success Criteria
1. Single leader invariant: At any given instant, only one node holds a valid, unexpired lease.
2. Failover liveness: When the active leader terminates or ceases heartbeats, another node successfully acquires leadership within TTL + renewal jitter.
3. Split-brain defense: Any write attempt by a deposed/stale leader with `token <= lastSeenToken` is rejected with an error.
4. Concurrency safety: Zero race conditions detected by `go test -race ./...`.
5. Standalone execution: Pure Go standard library implementation without external infrastructure dependencies (etcd/Redis servers not required to run tests or demo).

## Architecture
```text
┌─────────────────┐       ┌─────────────────┐
│ Candidate /     │       │ Candidate /     │
│ Node A          │       │ Node B          │
└────────┬────────┘       └────────┬────────┘
         │                         │
         │ Acquire / KeepAlive     │ Acquire / KeepAlive
         ▼                         ▼
  ┌─────────────────────────────────────────┐
  │         Coordinator (In-Memory)         │
  │ - Linearizable key-lease mapping        │
  │ - Monotonic revision / fencing token    │
  │ - Lease TTL & expiration monitoring     │
  └────────────────────┬────────────────────┘
                       │
                       │ Fenced Write (token, payload)
                       ▼
  ┌─────────────────────────────────────────┐
  │         Shared Resource (Storage)       │
  │ - Tracks last_seen_fencing_token        │
  │ - Rejects write if token <= last_seen   │
  │ - Accepts and updates state otherwise   │
  └─────────────────────────────────────────┘
```

## Components
1. `coordinator.Coordinator`: In-memory thread-safe coordination store simulating linearizable lease grants, renewals, releases, and monotonic revision increments.
2. `candidate.Node`: Distributed node candidate that attempts lease acquisition, maintains periodic heartbeats, observes leadership loss, and executes fenced work.
3. `storage.FencedStorage`: Shared storage engine enforcing check-and-set semantics on incoming fencing tokens.
4. `demo`: CLI application executing a full lifecycle demonstration: election -> active renewals -> simulated GC pause -> failover -> stale write rejection -> valid write acceptance.

## Test Strategy
- **Unit Tests**:
  - Coordinator lease acquisition, contention, renewal, and expiration.
  - Fenced storage token monotonic increment acceptance and stale token rejection.
- **Integration Tests**:
  - Normal failover upon voluntary release.
  - Automatic failover upon heartbeat timeout.
  - Stale write prevention under simulated pause.
- **Race Condition Tests**:
  - Concurrent campaign attempts by multiple nodes (`go test -race`).

## Execution Plan
1. Initialize Go module `labs/30-leader-election`.
2. Implement `internal/coordinator`, `internal/storage`, and `internal/candidate`.
3. Implement tests in `tests/` covering happy path, failover, GC pause simulation, and race detection.
4. Implement `cmd/demo/main.go` showcasing output logs clearly tracing terms, tokens, and fencing rejections.
5. Execute verification commands (`go test ./...`, `go test -race ./...`, `go run ./cmd/demo`).
6. Document implementation notes and execution results.

## Implementation Decisions
- **Decision 1: In-Memory Coordinator**: To ensure the lab is self-contained and reproducible anywhere without external daemon dependencies (etcd/Redis), the linearizable coordination engine is implemented in Go using monotonic clock (`time.Now()`), mutex synchronization, and revision counting.
- **Decision 2: Monotonic Time**: In accordance with research best practices (antirez & Kleppmann resolution), duration checks use Go monotonic time to guard against wall-clock skew.
- **Decision 3: Fencing Token Rejection Rule**: Strict inequality `incomingToken > lastSeenToken` required for write acceptance. Any `incomingToken <= lastSeenToken` is rejected.
