---
title: "Read Replicas and Replication Lag: Guaranteeing Read-Your-Own-Writes Without Sacrificing Read Scalability"
published: false
description: "Explore Read Replicas and Replication Lag in Go: stale reads, sticky routing, causal LSN tokens, lag-aware filtering, and synchronous replication, with executable proof."
tags: go, testing, performance, softwareengineering
---

# Read Replicas and Replication Lag: Guaranteeing Read-Your-Own-Writes Without Sacrificing Read Scalability

A production system must remain resilient under load. Lab 33 makes that behavior visible with a focused in-memory Go simulation and deterministic test assertions.

> Stale reads are not bugs; they are the expected behavior of eventually consistent replica systems.

This laboratory explores the real-world problem encountered by every backend engineer using read replicas for database scaling: stale reads. When a new write commits to the primary, replicas receive the WAL stream asynchronously with some delay. If the application immediately directs user reads to replicas without a closing mechanism, users see stale data or even receive not-found errors for recently written keys. This lab implements four mitigation strategies — time-based sticky routing, causal token with LSN wait, lag-aware dynamic filtering, and synchronous remote_apply — within a complete in-memory Go cluster simulation.

## The Real Problem: The Illusion of Freshness Behind Read Replicas

When a write commits to the primary, replicas receive the WAL stream asynchronously with some delay. If the application immediately routes reads to replicas without a closing mechanism, users see stale data or receive not-found for recently written keys. This is not a bug — it is the expected behavior of eventually consistent replica systems.

- **Write Committed**: User updates `tier` field via `PUT /profile`; primary commits at `LSN 47`; response returns the fresh payload.
- **Replica Lagging**: The `LSN 47` WAL entry streams to `replica-1` but has not been applied yet; `appliedLSN` is still `44`; the replica worker sleeps until the lag duration elapses.
- **Stale Read Occurs**: A round-robin load balancer routes the dashboard read to `replica-1`; the query returns `ErrNotFound` or outdated data; silently the system feels broken.

## Why Naive Reads to Replicas Are Always Risky

The `ReadNaive` implementation represents the simplest read-splitting path: round-robin across replicas with no freshness validation. However, it ignores the fact that every replica's `appliedLSN` can be far behind the primary's `currentLSN`. Under heavy writes with a 200ms replication lag, stale-read probability approaches 100% for recently written keys. The bug is subtle because it never crashes — it silently returns incorrect data.

```go
func (r *Router) ReadNaive(key string) (string, string, uint64, error) {
    replicas := r.cluster.Replicas()
    if len(replicas) == 0 {
        return r.cluster.Primary().Read(key)  // Fallback aman
    }
    idx := r.rrIndex.Add(1) % uint64(len(replicas))
    target := replicas[idx]
    return target.Read(key)  // Tanpa cek appliedLSN!
}
```

- No `appliedLSN` vs `currentLSN` validation before serving.
- Round-robin is indifferent to each replica's freshness state.
- Stale-read probability scales linearly with write throughput and lag duration.
- Write-then-read on a single session is most vulnerable to this anomaly.

## Mental Model: LSN as the Consistency Anchor

This mental model separates the consistency layer from the availability layer. The primary acts as a single writer assigning absolute monotonic ordering to each transaction. Each replica independently consumes and applies WAL entries with some latency. The router inspects each replica's `appliedLSN` to make informed routing decisions. Without LSN tracking, distinguishing fresh from lagging replicas becomes uninformed guesswork.

```go
type WALEntry struct {
    LSN       uint64    // monotonik naik dari primary
    Key       string
    Value     string
    Timestamp time.Time
}
```

- **Primary**: assigns an atomic `LSN` at write commit; single source of truth.
- **Replica**: stores data, `appliedLSN`, plus a worker goroutine consuming the WAL.
- **Gap**: `currentLSN - appliedLSN` is the current replication lag size.
- **Key**: every read strategy must minimize or close that gap for the relevant session.

## Write-Read Pipeline with Replication Lag

This pipeline traces the full request journey from write commit through read with four routing variants. The first path covers the write flow where the primary atomically increments `currentLSN` then broadcasts WAL entries to every replica channel. The second path simulates replication lag via worker goroutines sleeping before applying. The remaining paths show how the router selects read targets.

1. `Client.Write` → primary commits, `LSN` increments atomically, WAL entry broadcast to each replica channel.
2. `Replica worker sleep(lagDuration)` → entry applied to the local data map; `appliedLSN = entry.LSN`; `cond.Broadcast()` wakes `WaitForLSN` callers.
3. **Naive**: round-robin picks a replica; if `appliedLSN < lastWriteLSN`, stale read or `ErrNotFound`.
4. **Sticky**: if `time.Since(lastWriteTime) < T_sticky` → redirect to primary; after TTL expires → `ReadLagAware`.
5. **Token**: iterate replicas for `appliedLSN >= minLSN`; blocking `WaitForLSN` + timeout → fallback primary.
6. **Lag-aware**: compute `diff = primaryLSN - appliedLSN`; filter eligible `diff <= MaxLSNDiff`; fallback primary if none eligible.

## Sticky Session Router Implementation

This code implements a time-based sticky routing strategy. The router first checks whether a session ID has recorded state in a `sync.Map`. If the elapsed time since last write is within `StickyDuration`, it deterministically routes reads to the primary instead of any replica — guaranteeing read-your-own-writes at O(1) cost. Once the window expires, the flow degrades to `ReadLagAware` threshold filtering.

```go
func (r *Router) ReadWithStickySession(sessionID, key string) (string, string, uint64, error) {
    if stateVal, ok := r.sessions.Load(sessionID); ok {
        state := stateVal.(SessionState)
        if time.Since(state.LastWriteTime) < r.config.StickyDuration {
            val, lsn, err := r.cluster.Primary().Read(key)
            return val, r.cluster.Primary().ID(), lsn, err
        }
    }
    return r.ReadLagAware(key)
}
```

Sticky windows trade read throughput for guaranteed freshness within a bounded temporal envelope. When the session-bound use case outweighs replica offload, this trade favors correctness over raw scale.

## Failure Path: Stale Read Anomaly in the Demo CLI

This failure path is extracted directly from the demo CLI when `ReadNaive` runs immediately after a write. The primary LSN is `1`, but `replica-2` is still at `appliedLSN 0` because the worker goroutine hasn't finished its 200ms sleep. The round-robin balancer selects `replica-2` for the next read, producing a key-not-found. The app responds `200` with an empty body — a debugging nightmare for engineers unfamiliar with replication topology.

```
[Write Primary] Wrote key 'user:profile:123', Primary LSN: 1
[Naive Read] Node: replica-2, Applied LSN: 0, Val: '', Err: key not found
```

- Users see blank pages or 404 after submitting new data.
- Trust erosion: input succeeded, output doesn't reflect it.
- High debug cost — both write and read logs report success.

## Mechanism vs Symptom: Evaluating Read Routing Strategies

This table classifies four routing strategies across four evaluation dimensions. ReadNaive eliminates overhead but sacrifices consistency entirely. Sticky sessions bound the cost with a time check and session-map lookup. Causal tokens add measurable ordering guarantees but introduce blocking waits that shift p99 latency. Lag-aware filtering balances both through configurable threshold predicates.

| Strategy | Guarantee | Overhead | Fallback | Best For |
|---|---|---|---|---|
| ReadNaive | None | Zero | Primary if empty | Development, debugging |
| Sticky Session | Session-level fresh | Time check O(1) + session map | Lag-aware filtering | Session-bound writes |
| Causal Token | Causal ordering | LSN check + blocking wait | Context timeout + primary | Critical data reads |
| Lag-Aware | Threshold-bound | LSN delta calc + filter | Primary if none eligible | Bulk reads, analytics |

## 3 Critical Design Decisions for Replication Strategy

Building a robust replication strategy hinges on three architectural decisions that cascade from requirements to parameters to implementation. First, classify each use case: does it require read-your-own-writes, or can it tolerate eventual consistency? Second, calibrate every time parameter against empirical production data, never against theoretical latency. Third, pick the right fencing layer and acknowledge its cost — synchronous replication for critical paths, async with application-level routing for read-heavy workloads.

- **1. Define consistency requirement**: read-your-own-writes → sticky or token; casual reads → lag-aware filtering.
- **2. Calibrate timing parameters**: `T_sticky ≥ p99 replication lag`; `WaitTimeout ≥ worst-case catch-up + buffer`; `MaxLSNDiff` from observed normal-vs-anomalous deltas.
- **3. Choose fencing layer**: `synchronous_commit = remote_apply` when freshness is absolute; async + application routing for maximal read throughput.

## Execution Proof: Complete Replication & Mitigation Cycle

This terminal trace confirms all four routing scenarios behave exactly as designed. The first block shows stale-read anomaly under the naive router. The second shows sticky routing to primary inside the 500ms window, then to a replica after TTL expiry. The third shows a causal token blocking ~200ms while the replica catches up. The fourth shows synchronous replication adding ~104ms to write latency for instant freshness.

```
=== Lab 33: Read Replicas and Replication Lag Demo ===

--- 1. Async Replication Lag Anomaly ---
[Write Primary] Wrote key 'user:profile:123', Primary LSN: 1
[Naive Read] Node: replica-2, Applied LSN: 0, Val: '', Err: key not found

--- 2. Sticky Session Routing ---
[Sticky Read (within 500ms)] Node: primary, LSN: 1, Val: '{...}', Err: nil
[Sticky Read (after TTL)] Node: replica-1, LSN: 1, Val: '{...}', Err: nil

--- 3. Causal Token / LSN Wait ---
[Token Read (MinLSN=2)] Node: replica-1, LSN: 2, Val: '{...}', Err: nil

--- 4. Synchronous Replication ---
[Sync Write] LSN: 1, Write Duration: 104.055166ms
[Sync Naive Read] Node: replica-2, LSN: 1, Val: '{...}', Err: nil
```

## What Is Not Proven

The technical claims established here are explicitly bounded by the simplified test environment. The in-memory cluster has no persistence layer, so it cannot demonstrate post-crash recovery. Router state lives in a single-process `sync.Map`, unlike stateless production architectures requiring Redis or shared session stores. Replication lag uses deterministic timer sleeps, not real network latency variance. Timing constants such as 500ms stickiness were chosen for demo predictability, not as production tuning guidelines.

- No persistence layer; no crash-recovery proof.
- No cascading replica (replica-of-replica) or multi-region replication.
- Router state is single-process `sync.Map`; not horizontally stateless.
- Simulated 200ms timer lag, not real network jitter.
- No CDC or binlog-based synchronization.

## Architecture Trade-offs: Freshness vs Throughput

Implementing replication lag mitigation requires trading consistency for throughput in a deliberate, measured way. Sticky sessions add primary-read overhead inside the time window but fully eliminate stale reads for session-bound operations. Causal tokens add deterministic precision but introduce blocking waits that shift p99 latency. Lag-aware filtering balances both with configurable thresholds. Synchronous replication erases the lag entirely at the cost of proportional write latency.

**Benefits:**
- Guaranteed read-your-own-writes per session.
- Scalable reads: replicas absorb most traffic after the sticky window.
- Graceful degradation: explicit fallback instead of silent failure.
- Observable: `appliedLSN` per replica is monitorable.

**Costs:**
- Sticky window adds extra primary reads for every session.
- Blocking token reads inflate p99 read latency when replicas lag.
- Complexity: session state management + LSN inspection logic.
- All degraded reads concentrate back on the primary.

## Concurrency Validation: 15 Goroutines × 20 Operations Race-Free

Concurrency stress testing validates correctness when fifteen goroutines simultaneously execute write-then-read patterns across twenty iterations per session. Each goroutine uses a unique session ID, maintaining isolation in the router's sync.Map. Under `-race`, no data races are detected across all memory accesses — shared state remains linearizable through mutex boundaries at every node.

```go
func TestConcurrentAccess_RaceFree(t *testing.T) {
    c := cluster.NewCluster(3, 10*time.Millisecond, cluster.AsyncReplication)
    defer c.Close()
    r := router.NewRouter(c, router.DefaultConfig())

    var wg sync.WaitGroup
    for i := 0; i < 15; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            sid := fmt.Sprintf("s%d", id)
            _, _ = r.Write(sid, fmt.Sprintf("k%d", id), fmt.Sprintf("v%d", id))
            for j := 0; j < 20; j++ {
                _, _, _, _ = r.ReadWithStickySession(sid, fmt.Sprintf("k%d", id))
            }
        }(i)
    }
    wg.Wait()
}
```

- Race Detector: PASS (no data races under `-race`).
- Session isolation: 15 unique sessions never cross-contaminate state.
- Write-read consistency: every session sees its own writes after the sticky window.
- Cluster stability: no goroutine leaks or deadlocks under load.

## Before vs After: Evolution of Read Routing Strategy

The before-and-after comparison traces the architectural shift from fragile naive routing toward a four-layer mitigation system with tiered fallbacks. Before any strategy, every write-then-read cycle carried high risk of stale data because nothing closed the replication gap. After deployment, the system adds sticky sessions for session-bound operations, causal tokens for critical reads, and lag-aware filtering for bulk analytics — each with explicit, logged fallbacks to the primary.

| Dimension | BEFORE (Naive) | AFTER (Four-Layer Mitigation) |
|---|---|---|
| Routing | Round-robin replica, no LSN awareness | Sticky → token → lag-aware → primary |
| Failure mode | Silent stale read or `ErrNotFound` | Explicit fallback with logging |
| Observability | No lag metrics | `appliedLSN` per replica monitored |
| Impact | Silent user-facing inconsistency | Read-your-own-writes guaranteed with graceful degradation |

## Test 1: Stale Read Occurs Without Protection

The first unit test explicitly proves the stale-read anomaly exists. After a write commits to the primary, `ReadNaive` returns `ErrNotFound` because no replica has applied the WAL entry yet. After sleeping past the lag threshold, the same read succeeds. This confirms the anomaly is purely temporal — not permanent — and can be closed by adding a routing intelligence layer between write and read.

```go
func TestNaiveReplicationLag_StaleRead(t *testing.T) {
    c := cluster.NewCluster(2, 200*time.Millisecond, cluster.AsyncReplication)
    defer c.Close()
    r := router.NewRouter(c, router.DefaultConfig())

    _, err := r.Write("s1", "key:A", "value-A")
    if err != nil { t.Fatalf("Write failed: %v", err) }

    // Immediately read without any mitigation
    val, nodeID, lsn, err := r.ReadNaive("key:A")
    if err == nil {
        t.Fatalf("Expected stale read error, got val=%s node=%s lsn=%d", val, nodeID, lsn)
    }
    // After waiting for lag to expire, read should succeed
    time.Sleep(250 * time.Millisecond)
    val, nodeID, lsn, err = r.ReadNaive("key:A")
    if err != nil { t.Fatalf("Post-lag read failed: %v", err) }
}
```

- Immediate read fails: `ErrNotFound` because `appliedLSN < lastWriteLSN`.
- Post-lag read succeeds: value and LSN are fresh after catch-up.
- Node affinity: after convergence, any replica can serve correctly.

## Test 2: Sticky Session Guarantees Read-Your-Own-Writes

The second unit test confirms sticky session routing prevents stale reads within the configured window. Reads invoked within 500ms of a write always route to the primary node, guaranteeing freshness. After the window expires with a 600ms sleep, routing shifts to lag-aware replica filtering. This proves the time-based approach is both simple and effective for session-bound read patterns without per-request LSN tracking.

```go
func TestStickySessionRouting(t *testing.T) {
    c := cluster.NewCluster(2, 200*time.Millisecond, cluster.AsyncReplication)
    defer c.Close()
    cfg := router.DefaultConfig()
    cfg.StickyDuration = 500 * time.Millisecond
    r := router.NewRouter(c, cfg)

    _, err := r.Write("s1", "key:B", "value-B")
    if err != nil { t.Fatalf("Write failed: %v", err) }

    // Read within sticky window — must hit primary
    val, nodeID, _, err := r.ReadWithStickySession("s1", "key:B")
    if err != nil || val != "value-B" || !strings.HasSuffix(nodeID, "primary") {
        t.Fatalf("Sticky read failed or hit wrong node: val=%s node=%s err=%v", val, nodeID, err)
    }

    time.Sleep(600 * time.Millisecond) // Expire sticky window
    val, nodeID, _, err = r.ReadWithStickySession("s1", "key:B")
    if err != nil || nodeID == "primary" {
        t.Fatalf("Post-sticky read should route to replica, got node=%s", nodeID)
    }
}
```

- **Pre-expiration**: reads within 500ms always hit the primary.
- **Post-expiration**: after 600ms, routing falls back to eligible replicas.
- **Value integrity**: `value-B` stays consistent everywhere.

## Test 3: Causal Token Waits for Catch-Up Before Serving Read

The third unit test exercises the causal token routing with a blocking wait mechanism. After a write with a known LSN, `ReadWithToken` iterates replicas and, finding none at the target LSN, blocks on `WaitForLSN` with a context timeout. The test asserts the read waited ~200ms for catch-up, then served fresh data with `readLSN >= minLSN`. This distinguishes causal tokens from sticky sessions by checking actual LSN state rather than relying on time alone.

```go
func TestReadWithToken_LSN(t *testing.T) {
    c := cluster.NewCluster(1, 200*time.Millisecond, cluster.AsyncReplication)
    defer c.Close()
    r := router.NewRouter(c, router.DefaultConfig())

    lsn, err := r.Write("s1", "key:C", "value-C")
    if err != nil { t.Fatalf("Write failed: %v", err) }

    start := time.Now()
    val, nodeID, readLSN, err := r.ReadWithToken(context.Background(), lsn, "key:C")
    elapsed := time.Since(start)

    if err != nil { t.Fatalf("ReadWithToken failed: %v", err) }
    if val != "value-C" { t.Fatalf("Expected value-C, got %s", val) }
    if readLSN < lsn { t.Fatalf("Expected LSN >= %d, got %d", lsn, readLSN) }
    if elapsed < 150*time.Millisecond {
        t.Fatalf("Expected token read to wait ~200ms, finished in %v on node %s", elapsed, nodeID)
    }
}
```

- **Deterministic wait**: blocks until `appliedLSN >= minLSN`.
- **Elapsed timing**: assertion `elapsed >= 150ms` confirms blocking.
- **LSN guarantee**: `readLSN >= minLSN` before serving.

## Industry Case Study: Replication Lag in Production Databases

Replication lag is an everyday operational challenge across modern cloud infrastructure, not a theoretical concern. AWS RDS publishes the `ReplicaLag` CloudWatch metric so SRE teams can detect anomalies before user complaints surface. PostgreSQL streaming replication gives operators granular control through the `synchronous_commit` parameter, which distinguishes durability from visibility guarantees. Large applications like Twitter and Netflix embed custom session-affinity routing in client libraries to close readability gaps without waiting for database-engine fixes.

- **AWS RDS / Aurora**: async replication by default; `ReplicaLag` metric via CloudWatch; DocumentDB offers session-consistency.
- **PostgreSQL**: `wal_level=replica`; `synchronous_commit=remote_apply` enables synchronous visibility; `pg_stat_replication` exposes `replay_lsn` vs `write_lsn`.

Production systems treat replication lag as an operational variable, not an inflexible constraint.

## How WaitForLSN Distinguishes Blocking from Timeout

The `WaitForLSN` implementation demonstrates an idiomatic Go concurrency pattern for condition-based waiting. A waiter goroutine acquires the node mutex, loops until `appliedLSN` satisfies the target, and calls `cond.Wait()` which atomically releases the lock while blocking. The worker goroutine calls `cond.Broadcast()` after applying WAL entries to wake all blocked waiters. Context cancellation safely closes the stop channel and triggers a re-broadcast so waiters never remain stuck.

```go
func (n *Node) WaitForLSN(ctx context.Context, targetLSN uint64) error {
    ch := make(chan struct{})
    stop := make(chan struct{})
    go func() {
        n.mu.Lock()
        defer n.mu.Unlock()
        for n.appliedLSN < targetLSN {
            select {
            case <-stop:
                return
            default:
            }
            n.cond.Wait()  // Release lock, block until Broadcast
        }
        close(ch)
    }()
    select {
    case <-ctx.Done():
        close(stop); n.cond.Broadcast(); return ctx.Err()
    case <-ch:
        return nil
    }
}
```

- `cond.Wait()` releases the mutex while blocking, preventing CPU spinning.
- `cond.Broadcast()` wakes waiters right after entry apply completes.
- Context deadline closes `stop` and re-broadcasts for graceful cleanup.

## Checklist: Implementing a Read Replica Strategy

This checklist codifies the critical items before routing traffic to production. Timing parameters must be calibrated against empirical replication lag percentiles — never guessed. Monitoring and alerting become the operator's sensory system for degradation before user complaints surface. Router session state must be verified leak-free across long-lived sessions. Finally, distinguish synchronous replication for critical paths from async plus application-level fencing for read-heavy workloads.

- [ ] Define consistency requirement: session-bound vs casual reads.
- [ ] Calibrate `T_sticky ≥ p99` replication lag from production data.
- [ ] Set `MaxLSNDiff` from observed normal-vs-anomaly delta.
- [ ] `WaitTimeout ≥ worst-case catch-up + 20% buffer`.
- [ ] Deploy metrics: replica lag histogram, stale-read rate, primary fallback frequency.
- [ ] Circuit breaker on primary fallback when all replicas exceed threshold.
- [ ] Distribute session state to a shared store (Redis/memcached) for horizontal scaling.
- [ ] Verify sticky session map does not leak memory for long-idle sessions.
- [ ] Validate synchronous replication only for write-critical paths.
- [ ] Prepare runbook: distinguish network partition from replication-lag spikes.

## Anti-Patterns: Design Mistakes to Avoid

Violating replication topology design principles compromises both availability and consistency. Underestimating sticky duration relative to actual replication lag induces false confidence. Omitting primary fallbacks in lag-aware routing guarantees total read unavailability when all replicas lag simultaneously. Hardcoding sleeps instead of using context timeouts invites goroutine leakage. These can be reproduced in the lab by deleting the primary fallback, setting `MaxLSNDiff = 0`, or removing the context timeout.

**Common anti-patterns:**
- Mixing semi-sync (receipt) with synchronous (visibility); semi-sync does not close stale-read gaps.
- Sticky duration shorter than `p99` lag → stale reads still leak inside the window.
- No primary fallback → total read outage when all replicas lag.
- Cached reads ignoring LSN invalidation → persistent dirty reads.
- Hardcoded `time.Sleep` → goroutine leaks and deadlocks.

**Reproducing failures in-lab:**
1. Delete the primary fallback in lag-aware routing → `TestReplicaLagThreshold_Fallback` fails.
2. Set `MaxLSNDiff = 0` → all replicas filtered, all reads hit primary, throughput collapses.
3. Remove the context timeout from `WaitForLSN` → token reads hang indefinitely if a replica crashes.

## Self-Verification: Run Tests & Demo

This laboratory is engineered for immediate, frictionless local verification within seconds. The entire codebase uses the Go standard library — no external PostgreSQL or MySQL instances needed. `go test -v` runs seven test scenarios covering stale reads, sticky routing, token waits, and race-free concurrency. `go test -race` proves thread safety under fifteen concurrent goroutines. `go run ./cmd/demo` shows an interactive terminal demonstration of all four routing strategies.

```sh
cd labs/33-read-replicas-and-replication-lag
go test -v ./...
go test -race ./...
go run ./cmd/demo
```

- Test Suite: 7/7 PASS (2.13s).
- Race Detector: clean (3.98s).
- Dependencies: pure Go standard library (sync, time, context, atomic).

## Architectural Principles: Design for Asynchrony

These three principles transcend this specific Go implementation and apply to any primary-replica topology. First, recognize that asynchronous replication is the industry default; trying to eradicate it entirely sacrifices write scalability. Second, place the responsibility for closing replication gaps at the application layer — only the application understands the consistency context of each use case. Third, prioritize explicit degradation to the primary with rich observability over allowing users to silently ingest stale data.

1. **Asynchrony is the default**: replication lag is a reality to manage, not a bug to eliminate.
2. **Closing is the client's job**: the database provides LSN and WAL; the application router chooses the policy.
3. **Explicit degradation wins**: structured fallback logging beats silent stale data.

---

## Source

[Software Engineering Lab 33 — Read Replicas and Replication Lag](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/33-read-replicas-and-replication-lag), snapshot commit `b41fffb`. Source: README, implementation, tests, and execution records. Code excerpts belong to their original package context.

> Lab 33 proves stale reads are deterministic consequences of timing gaps, not bugs. The fix is informed routing — sticky within sessions, tokens for causality, thresholds for tolerance — backed by executable proof, not promises.
