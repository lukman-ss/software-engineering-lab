# Mechanisms and Lease-Based Leader Election

Target Lab: `labs/30-leader-election`  
Last Verified: September 28, 2026  
Addresses: Research Questions 1 & 2

---

## Research Question 1: What are the fundamental mechanisms for leader election in distributed systems?

### Consensus-Based Leader Election (Raft Model)

Leader election in consensus systems like Raft relies on randomized timeouts and quorum agreement.

**Mechanism**: 
- Each server begins as a **follower**, waiting for heartbeat messages from a leader.
- If no heartbeat is received within an **election timeout** (randomized between 150ms–300ms), the follower transitions to **candidate**, increments its term, votes for itself, and sends `RequestVote` RPCs to peers.
- A candidate becomes **leader** if it receives votes from a majority of servers in the same term.
- Leaders send periodic **AppendEntries RPCs** (heartbeats) to maintain authority and prevent new elections.

**Safety Guarantees**:
- At most one leader can be elected per term due to the **majority vote requirement**.
- Split-brain is prevented because network partitions cannot allow disjoint majorities to elect different leaders in the same term.
- Log replication safety is maintained via the **leader-commit invariant**: if an entry is committed in a term, it will be present in the logs of all leaders for higher terms.

**Source Verification**: 
- Ongaro & Ousterhout, "In Search of an Understandable Consensus Algorithm", USENIX ATC '14, Sections 5.2–5.3.
- Confirms randomized timeouts prevent split votes, and the election safety property holds as long as a majority of servers can communicate.

### Lease-Based Leader Election (etcd Model)

Alternative to consensus-based leader election: lease/heartbeat pattern using a strongly consistent key-value store.

**Mechanism**:
- A client attempts to acquire a lease (expiring token) from a coordination service (etcd, ZooKeeper).
- Only one client can hold the lease at a time due to the store's linearizability guarantee.
- The client must **renew the lease periodically** (heartbeat) to maintain leadership.
- If the client crashes, is partitioned, or fails to renew before lease expiration, the coordination service revokes the lease and grants it to another waiter.

**Key Properties**:
- **Mutual Exclusion**: At most one holder of the lease at any instant, guaranteed by the underlying store's linearizability.
- **Failure Detection**: Lease expiration provides bounded failure detection (no need for complex failure detectors).
- **Automatic Recovery**: No manual intervention needed; the lease times out and is re-granted.

**Source Verification**:
- etcd v3 documentation: "How to conduct leader election in etcd cluster" (https://etcd.io/docs/v3.5/tutorials/how-to-conduct-elections/).
- Confirms that etcd uses the `clientv3/concurrency` package to create a campaign for leadership, where the holder must periodically call `KeepAlive` on its lease.
- ZooKeeper documentation similarly describes leader election via ephemeral sequential znodes and watching for deletions.

### Fundamental Trade-Off: Consensus vs. Lease Heartbeat

| Aspect | Consensus-Based (Raft) | Lease/Heartbeat (etcd/ZooKeeper) |
|--------|------------------------|----------------------------------|
| **Leader Election Mechanism** | Randomized timeouts + voting | Lease acquisition + renewal |
| **Failure Detection** | Missing heartbeats + election timeout | Lease TTL expiration |
| **Split-Brain Prevention** | Quorum intersection (majority vote) | Linearizability of lease grant |
| **Recovery Time** | Election timeout (e.g., 150ms–300ms) | Lease TTL (configurable, e.g., 5s) |
| **Operational Complexity** | Requires managing cluster membership | Requires coordinating with external KV store |
| **Fencing Requirement** | Implicit via term numbers | Explicit via lease IDs or revision numbers |

---

## Research Question 2: How do lease/heartbeat mechanisms work to prevent split-brain scenarios?

### Core Principle: Time-Bounded Exclusivity

Lease/heartbeat mechanisms prevent split-brain by coupling mutual exclusion to time validity.

**Workflow**:
1. Client A acquires a lease with TTL = T (e.g., 5 seconds).
2. Client A must send heartbeat/renewal requests before T expires.
3. If Client A fails to renew (due to crash, partition, or pause), the lease expires automatically at time T.
4. After expiration, Client B can acquire the same lease.
5. Since leases cannot overlap in time (mutual exclusion is time-bounded), **two clients cannot simultaneously believe they hold the lease**.

### Critical Dependencies for Safety

To guarantee split-brain prevention, three timing assumptions must hold:
1. **Clock Synchronization (Bounded Drift)**: All nodes must agree on approximate time (e.g., via NTP). Clock jumps can cause premature expiration or overly long leases.
2. **Message Delivery Bounded by TTL**: Network delays must be significantly less than lease TTL, otherwise heartbeats may arrive after expiration despite a healthy leader.
3. **Processing Pausing Bounded by TTL**: Process pauses (e.g., GC stops) must be shorter than lease TTL, otherwise a paused leader may incorrectly believe it still holds the lease after expiration.

### Risk: Unbounded Pauses and Clock Skew

If a client's process pauses longer than the lease TTL (e.g., a stop-the-world GC pause of 10s when TTL=5s):
- The coordination service expires the lease and grants it to Client B.
- Client A, still paused, wakes up and believes it holds the lease (it last successfully renewed before the pause).
- Both Client A and Client B may act as leader concurrently → **split-brain**.

**Real-World Occurrence**: Such pauses are common in JVMs, .NET runtimes, or any language with tracing garbage collectors.

### Mitigation: Fencing Tokens (Generation Numbers)

Lease mechanisms alone do not prevent a paused client from acting on stale authority. To make leader election safe for correctness-critical work (not just efficiency), **fencing tokens** are required.

**How Fencing Works**:
- Each time a lease is granted, the coordination service returns a **monotonically increasing token** (e.g., etcd's `Revision`, ZooKeeper's `zxid` or node version).
- The leader must include this token with every write to the shared resource.
- The storage layer (or resource manager) rejects any write with a token less than or equal to the highest token it has already seen.

**Example**:
- Client A acquires lease, gets token = 100.
- Client A pauses for > TTL.
- Coordination service expires lease, grants to Client B with token = 101.
- Client A wakes up, attempts write with token = 100 → **rejected**.
- Client B writes with token = 101 → **accepted**.

**Necessity of Fencing**: Without fencing, a lease/heartbeat mechanism only provides **timed mutual exclusion** (useful for efficiency optimizations like deduplication), but **not safety** for concurrent updates to shared state.

**Source Verification**:
- Martin Kleppmann, "How to do distributed locking" (https://martin.kleppmann.com/2016/02/08/how-to-do-distributed-locking.html), Sections "Making the lock safe with fencing" and "Breaking Redlock with bad timings".
- etcd documentation: Leases return a globally unique ID; revisions increase monotonically and can be used as fencing tokens.
- ZooKeeper documentation: `zxid` (transaction ID) and node version numbers are strictly increasing and suitable for fencing.
- Raft paper: Term numbers serve as implicit fencing tokens; higher-term leaders invalidate lower-term leaders' authority.

---

## Summary Answers

### Q1: Fundamental Mechanisms
Leader election occurs via two primary patterns:
1. **Consensus-based** (Raft): Randomized election timeouts + majority voting ensure at most one leader per term.
2. **Lease/heartbeat-based** (etcd, ZooKeeper): Linearizable lease acquisition + periodic renewal; only one holder at a time due to store consistency.

### Q2: Lease/Heartbeat and Split-Brain Prevention
Lease/heartbeat prevents split-brain by binding leadership to time validity (TTL). However, safety requires:
- Bounded clock drift (to avoid premature/late expiration).
- Network delays ≪ TTL (to ensure heartbeats reflect liveness).
- Process pauses ≪ TTL (to prevent stale-leader belief after expiration).

For correctness-critical systems (not just efficiency), **fencing tokens** are essential to reject writes from paused or delayed clients who falsely believe they hold the lease.