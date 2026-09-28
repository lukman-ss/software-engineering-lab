# System Comparison and Production Best Practices

Target Lab: `labs/30-leader-election`  
Last Verified: September 28, 2026  
Addresses: Research Question 5

---

## Research Question 5: What are the best practices for implementing leader election in production systems?

### Synthesis of Trade-offs from Prior Sections

From the analysis of mechanisms (03-mechanisms-and-leases.md), fencing and Redlock debate (04-fencing-and-redlock.md), and source verification (02-sources.md), we derive actionable best practices grounded in the realities of distributed systems.

### General Production Best Practices

#### 1. Match the Guarantee to the Use Case

| Use Case | Locking Guarantee Required | Recommended Approach |
|----------|----------------------------|----------------------|
| **Efficiency Optimization** (deduplication, rate-limiting, cache warming) | Best-effort mutual exclusion (timed, not absolute) | Single-node atomic operations (Redis `SET NX EX`, database `SELECT FOR UPDATE SKIP LOCKED`) |
| **Correctness-Critical Coordination** (leader election for primaries, distributed transactions, schema migrations) | Strong mutual exclusion + fencing | Consensus-based (Raft) or strongly consistent KV store (etcd, ZooKeeper) with fencing tokens |
| **Service Discovery + Health-Based Leadership** | Session TTL + health checks + automatic failover | Consul with session TTL and service health checks |

> **Never** use Redis Redlock for correctness-critical operations without explicit fencing at the resource layer.

#### 2. Always Use Fencing for Correctness-Critical Locks

Even if the coordination mechanism provides leases (Redis Redlock, etcd, ZooKeeper, Consul), **the application must implement fencing** if:
- The shared resource is not inherently linearizable (e.g., external database, object storage, legacy API).
- Write operations are not idempotent.
- Concurrent writes could cause data loss, corruption, or financial error.

**Fencing Pattern**:
```python
# Pseudo-code for fenced critical section
token, valid = acquire_lease()  # Returns (token, is_valid)
if not valid:
    raise LockNotAcquired()

try:
    result = resource.perform_fenced_operation(data, token)
    # resource checks: if token <= last_seen_token: reject; else: accept and update last_seen_token
finally:
    release_lease()
```

**Sources**:
- Kleppmann (2016): "you need to include a fencing token with every write request... only the client with the greatest lock number will be able to write to the database."
- ZooKeeper recipes: use `zxid` for fencing in two-phase commit.
- etcd documentation: `Revision` numbers are monotonic and suitable for fencing.

#### 3. Prefer Monotonic Clocks for Timing-Sensitive Logic

Use `clock_gettime(CLOCK_MONOTONIC)` or language equivalents (`time.monotonic()` in Python, `time.Now().UnixNano()` in Go) for:
- Calculating elapsed time during lock acquisition.
- Implementing heartbeat/renewal logic.
- Guarding against wall-clock jumps (NTP steps, manual changes, virtualization pauses).

**Source**: antirez's acknowledgment in "Is Redlock safe?" that monotonic APIs should be used to avoid wall-clock issues.

#### 4. Set Timeouts Based on Empirical Observations, Not Guesses

Determine lease TTL, election timeout, and heartbeat intervals via:
- Measuring 99th percentile GC pause duration in your runtime.
- Measuring round-trip time (RTT) between nodes under load.
- Setting TTL = 3 × (max GC pause + max RTT) to allow for retransmission.

**Source**: etcd documentation recommends TTL >> expected network delay; ZooKeeper recommends session timeout > 2× RTT to avoid spurious expirations.

#### 5. Design for Graceful Degradation and Observability

- Emit metrics: lock acquisition latency, renewal success rate, fencing rejections.
- Log fencing rejections as WARNING-level events—they indicate a paused or delayed client.
- Implement circuit breakers if fencing rejections exceed a threshold (possible systemic pause).
- Provide manual override for leader election in case of total coordination service outage (use with extreme caution).

#### 6. Avoid Herd Effects During Failover

When using watch-based notification (etcd, ZooKeeper, Consul), **do not have all candidates watch the same object**. Instead:
- Watch only the immediate predecessor (ZooKeeper leader election recipe).
- Use `?prev` parameter in etcd watch to avoid thundering herd.
- Randomize renewal jitter (±10%) to desynchronize heartbeats.

**Source**: ZooKeeper recipes: "To avoid the herd effect, it is sufficient to watch for the next znode down on the sequence."

---

### Detailed Comparison: Production-Ready Systems

| Characteristic | etcd (v3.5+) | ZooKeeper (3.8+) | Consul (1.12+) | Single-Node Redis |
|----------------|--------------|------------------|----------------|-------------------|
| **Consensus Protocol** | Raft | Zab (Paxos variant) | Raft | N/A (single node) |
| **Linearizability** | ✅ Strong | ✅ Strong | ✅ Strong | ❌ Weak (no quorum) |
| **Fencing Token Source** | ✅ Lease ID + Revision | ✅ zxid + Node Version | ✅ Session ID + Raft Term | ❌ None |
| **Automatic Leader Eviction** | ✅ Lease TTL expiry | ✅ Ephemeral node deletion | ✅ Session TTL + health checks | ❌ Manual failover |
| **Failure Detection** | Lease TTL + TCP keepalive | Session timeout + heartbeat | Health checks + TTL | Application-level ping |
| **Clock Sensitivity** | Low (TTL vs. network delay) | Very low (ephemeral via session) | Low (health-based TTL) | High (Redis clock for expiry) |
| **Operational Overhead** | Moderate (3+ node cluster) | Moderate (3+ node cluster) | Moderate-High (agent per node + gossip) | Low (single instance) |
| **Throughput (Locks/sec)** | 10K–50K | 5K–20K | 1K–5K (agent-bound) | 50K–100K (but unsafe) |
| **Best For** | Correctness-critical services (Kubernetes, cloud infra) | Battle-tested infra (HDFS, HBase, Kafka) | Service mesh, multi-datacenter, intent-based networking | Prototyping, caching, non-critical deduplication |
| **Production Adoption** | CNCF Graduated (Kubernetes core) | Apache Top-Level (Hadoop ecosystem) | HashiCorp stack (Nomad, Vault, Terraform) | Ubiquitous but misused |

#### etcd-Specific Production Notes
- Use `clientv3/concurrency` package for leader election and mutex patterns.
- Prefer `LeaseGrant` + `KeepAlive` over raw KV TTL for automatic cleanup.
- Monitor `etcd_server_leases_total` and `etcd_network_client_grpc_*` metrics.
- Set `--heartbeat-interval=100ms` and `--election-timeout=500ms` for WAN deployments.

#### ZooKeeper-Specific Production Notes
- Use Apache Curator framework (`LeaderSelector`, `Locks`, `Barriers`) to avoid low-level pitfalls.
- Disable symmetric encryption (`skipACL=no`) unless ACLs are strictly required (adds latency).
- Monitor `zk_avg_latency` and `zk_outstanding_requests` via JMX.
- Use `autopurge.snapRetainCount=3` and `autopurge.purgeInterval=1` to manage disk usage.

#### Consul-Specific Production Notes
- Use Consul Sessions (`ttl` + `lock_delay=0`) for simple leader election.
- Rely on service health checks (`status=passing`) to trigger automatic failover.
- Monitor `consul_raft_apply_lag` and `consul_session_num` telemetry.
- Avoid long poll on KV for leader election; use session-based locks instead.

#### Single-Node Redis-Specific Production Notes
- Only use for:
  - Rate limiting with `INCR` + `EXPIRE` (approximate counters).
  - Idempotent job deduplication where duplicate execution is acceptable (e.g., sending welcome email).
  - Never use `SET NX EX` for locking financial updates, inventory decrements, or schema migrations.
- If Redis is already deployed, consider Redisson library which provides fencing via `RedLock` + external token store (advanced).

---

### Recommended Production Architecture

For most microservices and cloud-native applications requiring correctness-critical leader election:

```
┌─────────────┐    ┌────────────────┐    ┌────────────────────┐
│   Service   │───▶│   etcd Cluster │───▶│  Shared Resource   │
│ (Candidate) │    │  (3–5 nodes)   │    │  (Database, Queue) │
└─────────────┘    └────────────────┘    └────────────────────┘
                   │        ▲
                   │        │ Lease + KeepAlive (TTL=15s)
                   │        │ Fencing Token: LeaseID:Revision
                   ▼        │
           ┌────────────────┐
           │   Watcher      │◄─────────────────────┐
           │ (etcd client)  │                      │
           └────────────────┘                      │
                                                   │ Fencing Check: 
                   ┌────────────────┐             │ if token ≤ last_seen: REJECT
                   │  Resource      │             │ else: ACCEPT + update
                   │  (write gate)  │             └─────────────────────────────
                   └────────────────┘
```

**Why This Works**:
1. **Strong Consistency**: etcd provides linearizable lease operations via Raft.
2. **Automatic Failover**: Leader loses lease on crash/network partition → fencing token increases.
3. **Fencing Enforcement**: Resource layer validates monotonically increasing tokens.
4. **Observability**: Metrics on lease renewals, fencing rejections, and acquisition latency.
5. **No Herd Effect**: Only the current leader renews; watchers are passive until notified.

---

### Summary Answers

#### Q5: Production Best Practices for Leader Election
1. **Match guarantee to use case**: Use best-effort locks (Redis single node) only for efficiency; use strongly consistent systems (etcd/ZooKeeper/Consul) with fencing for correctness.
2. **Always implement fencing** for shared resources where concurrent writes could cause corruption—rely on monotonic tokens (lease ID, revision, zxid) from the coordination service.
3. **Prefer monotonic clocks** for all timing measurements to avoid NTP/jump hazards.
4. **Set timeouts empirically**: Base TTL and election timeouts on measured GC pause and RTT, not guesswork.
5. **Avoid herd effects** by watching immediate predecessors, not the head of a queue.
6. **Design for observability**: Track fencing rejections as early warning signals of systemic pauses or clock issues.
7. **Choose the tool for the job**:
   - **etcd**: Cloud-native, Kubernetes-adjacent, high performance.
   - **ZooKeeper**: Mature, battle-tested, lower-level control.
   - **Consul**: Service mesh integration, intention-based networking.
   - **Redis single node**: Prototyping, caching, non-critical deduplication only.

This approach ensures that leader election mechanisms in production systems are not only theoretically sound but empirically validated against the failure modes observed in real-world distributed systems: clock jumps, garbage collection pauses, network partitions, and process suspensions.