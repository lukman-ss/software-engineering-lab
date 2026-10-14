---
title: "Read Replicas & Replication Lag: Mitigating Stale Reads and Enforcing Read-Your-Own-Writes Guarantees"
published: false
description: "Explore read replicas, replication lag anomalies, and session guarantees (read-your-own-writes) in Go with causal tokens and sticky routing."
tags: database, go, systemdesign, backend
---

# Read Replicas & Replication Lag: Mitigating Stale Reads and Enforcing Read-Your-Own-Writes Guarantees

A production system must remain resilient under load. Lab 33 makes that behavior visible with a focused Go simulation and deterministic test assertions.

> Replication lag is the distance between primary commit and replica apply.

## Stale Read Anomalies on Read Replicas

When scaling database reads across secondary replicas, asynchronous replication creates a propagation window where replicas lag behind the primary. An immediate subsequent read by the writer can hit a replica that has not yet applied the latest write-ahead log entry, returning stale records or false not-found errors.

## Anatomy of Asynchronous Replication Streams

In asynchronous replication mode, the primary commits transactions locally, assigns a monotonically increasing Log Sequence Number, and immediately acknowledges the client. Write-ahead log records are queued and streamed to replicas in the background, making replication lag structurally inevitable under heavy traffic.

## The Read-Your-Own-Writes Consistency Guarantee

Read-your-own-writes ensures that an individual user session immediately observes its own updates, even if global readers observe eventual consistency. Formally proven by Terry et al., this session-scoped constraint prevents user-facing anomalies without the crippling overhead of global consensus.

## Why Naive Read-Splitting Strategies Fail

A naive read-write splitter routes all writes to the primary and distributes all reads round-robin across replicas without inspecting replication status. In local development environments with zero network delay, this appears flawless. In production under real I/O pressure, it causes intermittent not-found and stale read regressions.

```go
func (r *Router) ReadNaive(key string) {
    target := replicas[rrIndex.Add(1) % len(replicas)]
    return target.Read(key)
}
```

## Architecture Overview: A Three-Layer Simulation System

The lab architecture isolates cluster mechanics from client-side routing intelligence. The Cluster layer handles concurrent WAL broadcasting and replica worker routines, while the Router layer maintains session metadata and evaluates freshness policies before executing physical reads.

## Primary Node Execution and Atomic LSN Generation

To establish a rigorous order of mutations, the primary node utilizes an atomic 64-bit unsigned integer counter representing the Log Sequence Number. Every write acquires the lock, mutates local state, assigns the new LSN, and packages a WALEntry for asynchronous replica distribution.

```go
type WALEntry struct {
    LSN       uint64
    Key       string
    Value     string
    Timestamp time.Time
}
```

```go
lsn := c.currentLSN.Add(1)
c.primary.mu.Lock()
c.primary.data[key] = value
c.primary.appliedLSN = lsn
c.primary.mu.Unlock()
```

## Replica Worker Goroutine and Lag Ingestion Pipeline

Each replica runs a dedicated worker goroutine consuming from a buffered WAL channel. Artificial delay simulates network transit and disk synchronization latencies. Once applied, sync.Cond broadcasts trigger waiting causal read queries to wake up deterministically.

```go
for entry := range replica.walChannel {
    if lag > 0 { time.Sleep(lag) }
    replica.mu.Lock()
    replica.data[entry.Key] = entry.Value
    replica.appliedLSN = entry.LSN
    replica.cond.Broadcast()
    replica.mu.Unlock()
}
```

## Strategy 1: Time-Based Sticky Session Routing

Time-based sticky routing tracks the timestamp of each session's last write. Any subsequent read within a configurable TTL (e.g., 5 seconds) routes unconditionally to the primary node. It eliminates stale reads for the writer with zero coordination overhead, though it relies on empirical latency assumptions.

```go
if stateVal, ok := r.sessions.Load(sessionID); ok {
    state := stateVal.(SessionState)
    if time.Since(state.LastWriteTime) < r.config.StickyDuration {
        return r.cluster.Primary().Read(key)
    }
}
return r.ReadLagAware(key)
```

## Strategy 2: Causal Token and Minimum LSN Routing

Causal token routing binds each write's commit LSN to subsequent read requests. The router evaluates replicas for applied LSN progress. If no replica has reached the target LSN, it executes an event-driven wait with a deadline, falling back to the primary if timeout thresholds expire.

```go
for _, replica := range replicas {
    if replica.AppliedLSN() >= minLSN {
        return replica.Read(key)
    }
}
```

```go
waitCtx, cancel := context.WithTimeout(ctx, r.config.WaitTimeout)
defer cancel()
if err := targetReplica.WaitForLSN(waitCtx, minLSN); err == nil {
    return targetReplica.Read(key)
}
return r.cluster.Primary().Read(key) // Fallback aman
```

## Event-Driven Catch-Up with Condition Variables

Rather than wasting CPU cycles on repetitive polling, WaitForLSN suspends the reader goroutine on a sync.Cond condition variable. When the background replica worker applies the transaction, it broadcasts a wake-up signal, immediately resuming the reader with sub-millisecond precision.

```go
func (n *Node) WaitForLSN(ctx context.Context, target uint64) error {
    go func() {
        n.mu.Lock(); defer n.mu.Unlock()
        for n.appliedLSN < target {
            n.cond.Wait() // Dibangunkan oleh Broadcast worker
        }
        close(ch)
    }()
    select {
    case <-ctx.Done(): return ctx.Err()
    case <-ch: return nil
    }
}
```

## Strategy 3: Lag-Aware Dynamic Routing and Primary Fallback

For general reads without explicit session constraints, the lag-aware routing strategy evaluates the delta between primary LSN and replica applied LSN. Nodes exceeding the maximum tolerable difference are temporarily removed from rotation, triggering automated fallback to primary if all replicas degrade.

```go
primaryLSN := r.cluster.CurrentLSN()
for _, replica := range replicas {
    diff := primaryLSN - replica.AppliedLSN()
    if diff <= r.config.MaxLSNDiff {
        eligible = append(eligible, replica)
    }
}
```

```go
if len(eligible) == 0 {
    return r.cluster.Primary().Read(key) // Fallback-lag
}
idx := r.rrIndex.Add(1) % len(eligible)
return eligible[idx].Read(key)
```

## Strategy 4: Synchronous Replication and the Freshness Tax

Synchronous replication halts the primary write process until all designated replica nodes acknowledge full local execution. While this completely eliminates stale reads across the entire cluster, it imposes a severe latency penalty on every write operation proportional to network round-trip delays.

```go
if c.replMode == SyncReplication {
    for _, replica := range c.replicas {
        if lag > 0 { time.Sleep(lag) }
        replica.data[key] = value
        replica.appliedLSN = lsn
    }
}
```

## Quantitative Strategy Matrix: Trade-Off Evaluation

Choosing a mitigation strategy requires balancing consistency guarantees against latency overhead and primary load. While sticky routing offers the lowest operational friction, causal token routing guarantees mathematical correctness without permanently overloading the primary cluster.

Strategi | Garansi | Latensi Read | Overhead Primary | Kompleksitas
|---|---|---|---|---|
Naive Read | Eventual | Terendah | Nol | Sangat Rendah
Sticky Session | Read-Own-Writes | Terendah | Sedang (T_sticky) | Rendah
Causal Token | Causal Session | Sedang (Wait) | Rendah (Fallback) | Sedang
Lag-Aware | Bounded Lag | Rendah | Kondisional | Sedang
Sync Mode | Linearizable | Terendah | Nol (Write Mahal) | Tinggi

## Interactive CLI Demonstration Execution Logs

Running the lab demo CLI provides reproducible evidence of every failure mode and mitigation strategy. Naive reads fail immediately under asynchronous delay, sticky routing shields the author during post-write navigation, causal tokens block cleanly until catch-up, and synchronous writes record a measurable 102ms latency overhead.

```go
[Write Primary] key='user:profile:123', LSN: 1
[Naive Read] Node: replica-2, Applied LSN: 0, Err: key not found
[Sticky Read (within 500ms)] Node: primary, LSN: 1, Val: Alice
[Sticky Read (after TTL)] Node: replica-1, LSN: 1, Val: Alice
```

```go
[Write Primary] key='order:101', LSN: 2
[Token Read (MinLSN=2)] Node: replica-1, LSN: 2, Val: completed
[Sync Write] LSN: 1, Write Duration: 102.19ms
[Sync Naive Read] Node: replica-2, LSN: 1, Val: Bob, Err: <nil>
```

## Comprehensive Unit and Concurrency Test Verification

All seven test cases pass under Go's race detector. The test suite verifies exact stale read occurrences, sticky session TTL hand-offs, causal LSN assertions, threshold fallback to primary, immediate synchronous visibility, and rigorous thread safety under concurrent writer-reader pressure.

```go
TestNaiveReplicationLag_StaleRead: PASS (0.60s)
TestStickySessionRouting: PASS (0.55s)
TestReadWithToken_LSN: PASS (0.20s)
TestReplicaLagThreshold_Fallback: PASS (0.00s)
TestSynchronousReplication_Freshness: PASS (0.10s)
TestConcurrentAccess_RaceFree: PASS (1.10s)
TestWaitForLSN_ContextTimeout: PASS (0.05s)
```

## Mapping In-Memory Concepts to Production PostgreSQL and MySQL

The concepts implemented in this lab directly map to production database mechanisms. PostgreSQL exposes Write-Ahead Log positions via LSN functions, while MySQL leverages Global Transaction Identifiers (GTID) and WAIT_FOR_EXECUTED_GTID_SET to implement deterministic causal synchronization.

## Cloud Managed Database Replicas: Aurora vs Traditional RDS

Managed cloud databases approach replication differently. AWS Aurora utilizes a distributed log-structured storage volume where replicas observe minimal cache invalidation lag (<20ms). Standard AWS RDS and Azure PostgreSQL rely on traditional network streaming where I/O bottlenecks can induce significant replication lag.

## Scaling Session State: From In-Memory Maps to Distributed Systems

While the lab models session tracking via an in-memory sync.Map, horizontal production environments require distributed state coordination. Production systems either store session write tokens in Redis or encode minimum LSN requirements directly into signed client cookies and HTTP headers.

## Mitigating LSN Polling Overhead with Asynchronous Heartbeats

Querying replica positions on every client request introduces severe database round-trip overhead. High-throughput architectures deploy background health workers that poll replica status periodically or read dedicated heartbeat tables, allowing routers to make routing decisions entirely in memory.

## Production Checklist for Database Read Replication

Before distributing read traffic across secondary replicas, engineering teams must establish clear operational guardrails. Enforce automated primary fallbacks, monitor p99 replication lag metrics, and ensure post-write query pathways explicitly request session guarantees.

## Implementation Scope and Architectural Boundaries

This lab focuses strictly on validating routing logic and session guarantees. By implementing an in-memory key-value cluster with atomic LSN sequencing, it demonstrates replication lag phenomena deterministically without requiring heavyweight database infrastructure or multi-node clustering.

## Core Architectural Principles for Replication Management

Asynchronous replication remains the dominant industry standard for high-throughput database systems. Enforcing read-your-own-writes requires thoughtful client routing rather than hoping for zero-lag replicas. Combining lightweight sticky routing with causal token fallbacks ensures both performance and consistency.

## Source

[Software Engineering Lab 33 - Read Replicas & Replication Lag](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/33-read-replicas-and-replication-lag), snapshot commit `cfba3b76`. Source: README, implementation, tests, and execution records. Code excerpts belong in their original package context.
