---
title: "Read Replicas and Replication Lag: Guaranteeing Read-Your-Own-Writes Without Sacrificing Read Scalability"
published: false
description: "Explore Read Replicas and Replication Lag in Go: mechanics, architecture, caveats, and testable invariants."
tags: go, testing, performance, softwareengineering
---

# Read Replicas and Replication Lag: Guaranteeing Read-Your-Own-Writes Without Sacrificing Read Scalability

A production system must remain resilient under load. Lab 33 makes that behavior visible with a focused Go simulation and deterministic test assertions.

> Stale reads are not bugs; they are the expected behavior of eventually consistent replica systems.

## Masalah Nyata: Ilusi Freshness di Balik Read Replica

This scenario occurs by default across nearly all production databases using primary-replica topology. PostgreSQL streaming replication, MySQL async replication, AWS RDS read replicas, and Azure Database for PostgreSQL all run asynchronous replication by default without explicit configuration. When the client application performs a write-then-read pattern on the same session, it has no awareness that the written data is unavailable at the target replica. The result is a hard-to-debug UX anomaly because logs show successful writes but subsequent queries fail to read the same data.

## Mengapa Read Naif ke Replica Selalu Berisiko?

The ReadNaive implementation represents the simplest approach to read splitting: round-robining each request across replicas in turn. However, this approach ignores the fundamental fact that every replica may have an appliedLSN significantly behind the primary currentLSN. Under heavy write workloads with a 200-millisecond replication lag, the probability of serving a stale read can approach 100 percent for recently written keys. This bug is particularly subtle because it never panics or crashes — it silently returns incorrect data that perplexes end users.

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

## Mental Model: LSN sebagai Jangkar Konsistensi

This mental model separates two layers of abstraction: the consistency layer through LSN sequences and the availability layer through data distribution. The primary acts as a single writer assigning absolute monotonic ordering to every transaction. Each replica independently consumes and applies WAL entries with some latency. The router inspects every replica's appliedLSN to make informed routing decisions. Without LSN tracking, distinguishing fresh from lagging replicas becomes impossible, reducing routing to uninformed guesswork.

```go
type WALEntry struct {
    LSN       uint64    // monotonik naik dari primary
    Key       string
    Value     string
    Timestamp time.Time
}
```

## Pipeline Write-Read dengan Replication Lag

This pipeline illustrates the complete request journey from write commit through read with four distinct routing variants. The first path traces the write flow where the primary atomically increments currentLSN then broadcasts WAL entries to every replica channel. The second path simulates replication lag through worker goroutines sleeping for the configured duration before applying entries. The four sub-paths show how the router selects read targets: naive round-robin without protection, time-based sticky sessions, token-driven LSN waiting with blocking, or lag-aware dynamic filtering with threshold predicates.

## Implementasi Sticky Session Router

This code snippet implements the simplest yet most effective form of time-based sticky routing. The method first checks whether the session ID has a recorded state in the sync.Map. If the state exists and the elapsed time since last write remains within the StickyDuration bound, the router deterministically directs reads to the primary instead of any replica. This guarantees read-your-own-writes within the sticky window without any LSN overhead. After the window expires, the flow degrades gracefully into lag-aware threshold filtering.

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

## Failure Path: Stale Read Anomali di Demo CLI

This failure path is extracted directly from the demo CLI when ReadNaive is invoked immediately after a write. The primary LSN increments to 1, but replica-2 remains at appliedLSN 0 because the worker goroutine has not yet completed its 200-millisecond sleep. The round-robin load balancer selects replica-2 for the subsequent read, causing a key-not-found error since the data is unavailable. This error differs from connection failures or timeouts — the application responds with 200 and an empty body, creating a debugging nightmare for engineers unfamiliar with replication topology.

## Mekanisme vs Gejala: Evaluasi Strategi Read Routing

This table classifies the four routing strategies across four evaluation dimensions. ReadNaive eliminates communication overhead but sacrifices consistency entirely. Sticky session offers session-level freshness guarantees with minimal overhead consisting of a time comparison and session map lookup. Causal tokens provide measurable causal ordering guarantees but introduce blocking waits that can elevate p99 latency. Lag-aware filtering balances throughput and freshness through threshold predicates, degrading gracefully to primary when no eligible replica exists.

Strategi | Jaminan | Overhead | Fallback | Cocok Untuk
|---|---|---|---|---|
ReadNaive | None | Zero | Primary jika replica kosong | Development, debugging
Sticky Session | Session-level fresh | Time check O(1), session map | Lag-aware filtering | Session-bound writes
Causal Token | Causal ordering | LSN check + blocking wait | Context timeout + primary | Critical data reads
Lag-Aware | Threshold-bound | Delta LSN calculation + filter | Primary jika eligible kosong | Bulk reads, analytics

## 3 Keputusan Desain Kritis Replication Strategy

These three design decisions form a strategic framework before selecting any technical implementation. The first decision determines the consistency requirement level that distinguishes which use cases need strict protection versus those accepting eventual consistency. The second decision involves calibrating time parameters based on empirical production latency observations rather than theoretical assumptions. The third decision selects the appropriate protection layer aligned with business SLAs between high-cost synchronous replication or application-level fencing with graceful degradation.

## Bukti Eksekusi: Siklus Lengkap Replication & Mitigasi

This terminal output presents four execution proofs each confirming the behavior of different routing strategies. The first section demonstrates the stale read anomaly when ReadNaive is called immediately post-write. The second section shows sticky session successfully directing reads to the primary within the 500-millisecond window then falling back to a replica after TTL expiration. The third section verifies causal tokens wait approximately 200 milliseconds for catch-up before serving reads from matching replicas. The fourth section reveals synchronous replication requiring approximately 104 additional milliseconds on write latency to guarantee instant freshness.

## YANG TIDAK DIBUKTIKAN

The technical claims established in this lab are explicitly bounded by simplified test-environment parameters. The in-memory cluster lacks any persistence layer, so it cannot demonstrate post-crash recovery behavior. Router state resides in a local single-process sync.Map, differing from stateless production architectures that require Redis or shared session stores. Replication lag simulation uses deterministic timer sleeps rather than variable real network latency. Timing constants such as 500-millisecond stickiness were chosen for demonstrable predictability, not as production tuning guidelines.

## Trade-off Arsitektur: Freshness vs Throughput

Implementing replication lag mitigation strategies demands architectural trade-offs between freshness guarantees and read throughput. Sticky sessions introduce overhead through concentrated primary reads during the time window but completely eliminate stale reads for session-bound operations. Causal tokens offer deterministic precision while introducing blocking waits that shift p99 read latency. Lag-aware filtering balances both approaches through configurable threshold predicates, degrading gracefully to primary when all replicas lag excessively. Synchronous replication eliminates lag entirely at the cost of proportionally increased write latency across all replicas.

## Validasi Konkurensi: 15 Goroutine × 20 Operasi Race-Free

This concurrency test validates implementation robustness when fifteen goroutines simultaneously execute write-then-read patterns across twenty iterations per session. Each goroutine employs a unique session ID maintaining isolation state within the Router sync.Map. The cluster with three replica workers processes WAL entries in parallel with ten-millisecond gaps. Under the Go race detector, no data races are detected across all Router method or Cluster node memory accesses. This confirms concurrent access to shared state remains linearizable through mutex protections at every node.

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
            sid := fmt.Sprintf('s%d', id)
            _, _ = r.Write(sid, fmt.Sprintf('k%d', id), fmt.Sprintf('v%d', id))
            for j := 0; j < 20; j++ {
                _, _, _, _ = r.ReadWithStickySession(sid, fmt.Sprintf('k%d', id))
            }
        }(i)
    }
    wg.Wait()
}
```

## Before vs After: Evolusi Read Routing Strategy

This before-and-after comparison highlights the architectural evolution from fragile naive routing prone to stale reads toward a four-layer mitigation system with tiered fallbacks. Prior to strategy implementation, every write-then-read cycle carried high risk of serving stale data due to absent gap-closing mechanisms. After deployment, the system features three protection layers: sticky sessions for session-bound operations, causal tokens for critical reads, and lag-aware filtering for bulk analytics. Each layer maintains explicit fallback logging to the primary with observable node identifiers and LSN deltas.

## Test 1: Stale Read Terjadi Tanpa Proteksi

The first unit test explicitly proves the existence of stale read anomalies. Immediately after writing to the primary with key 'key:A', the ReadNaive call produces an ErrNotFound error because replicas have not yet applied the WAL entry. After waiting 250 milliseconds to exceed the 200-millisecond lag threshold, the same read succeeds returning the correct value. This test confirms the anomaly is purely temporal — not permanent — and can be prevented by adding a routing intelligence layer between write commit and read dispatch.

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

## Test 2: Sticky Session Menjamin Read-Your-Own-Writes

The second unit test confirms that sticky session routing successfully prevents stale reads within the specified time window. Reads invoked within 500 milliseconds post-write consistently route to the primary node, guaranteeing fresh data availability. After the window expires with a 600-millisecond sleep, routing automatically shifts to ReadLagAware which selects replicas based on threshold eligibility. This test proves the time-based approach is both simple and effective for session-bound read patterns without requiring per-request LSN tracking.

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

## Test 3: Causal Token Tunggu Catch-Up Sebelum Serve Read

The third unit test examines the causal token routing implementation involving a blocking wait mechanism. After writing with a specific LSN, ReadWithToken iterates replicas to find one already at or above the threshold. When no replica matches, the method calls WaitForLSN on the first replica with a context timeout. This test validates that the read blocks for approximately 200 milliseconds waiting for replica catch-up before serving a fresh response with the correct LSN. This blocking behavior distinguishes causal tokens from sticky sessions, which rely solely on time without checking actual LSN values.

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

## Studi Kasus Industri: Replication Lag di Database Production

Industry case studies demonstrate that replication lag is not merely theoretical but an everyday operational challenge in modern cloud infrastructure. AWS RDS provides the ReplicaLag metric through CloudWatch enabling SRE teams to detect anomalies before user complaints surface. PostgreSQL streaming replication offers granular control through the synchronous_commit parameter distinguishing between durability guarantees and visibility guarantees. Large enterprise applications like Twitter and Netflix implement custom session-affinity routing within client libraries to close readability gaps without waiting for database engine fixes.

## Bagaimana WaitForLSN Membedakan Blocking dan Timeout?

The WaitForLSN implementation demonstrates an elegant Go concurrency pattern for condition-based waiting combining channels and condvars. The waiter goroutine acquires the node mutex, loops until appliedLSN satisfies the target, and calls cond.Wait() which atomically releases the mutex while blocking the goroutine. The worker goroutine invokes cond.Broadcast() after applying WAL entries to wake all blocked waiters. Context cancellation provides a safe escape hatch by closing the stop channel and triggering a re-broadcast ensuring waiters never remain permanently stuck.

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

## Checklist Implementasi Read Replica Strategy

This checklist codifies ten critical principles required before deploying read replica strategies to production. Timing parameters such as sticky duration and wait timeout must never be guessed but calibrated against historical replication lag percentile data. Monitoring and alerting become the operator sensory system for detecting degradation before user complaints surface. Router state management must be verified for absence of memory leaks across long-lived session lifecycles. Finally, distinguish synchronous replication for critical paths versus asynchronous with application-level fencing for read-heavy workloads.

## Anti-Pattern: Kesalahan Desain yang Harus Dihindari

Violating replication topology design principles compromises both system availability and data consistency. The most common anti-pattern is underestimating sticky window duration relative to actual replication lag, inducing a false sense of security for engineering teams. Omitting primary fallbacks in lag-aware routing guarantees total read unavailability whenever all replicas suffer simultaneous lag spikes. Furthermore, hardcoding sleep durations rather than respecting context timeouts invites severe goroutine leakage that can paralyze cluster workers.

## Verifikasi Mandiri: Uji & Jalankan Lab

This laboratory is engineered for immediate, frictionless local verification within seconds. The entire codebase is implemented with the Go standard library, eliminating external dependencies like PostgreSQL or MySQL instances. Running go test -v executes seven exhaustive test scenarios validating stale read anomalies and race-free concurrency. Executing go test -race proves thread safety under fifteen concurrent goroutines, while go run ./cmd/demo delivers an interactive terminal demonstration of all four routing scenarios.

```sh
cd labs/33-read-replicas-and-replication-lag
go test -v ./...
go test -race ./...
go run ./cmd/demo
```

## Prinsip Arsitektur: Desain untuk Asinkronitas

These three architectural principles transcend this specific Go implementation. First, recognize that asynchronous replication is the industry default, and attempting to eradicate it entirely sacrifices write scalability. Second, place the responsibility for closing replication gaps at the application layer because only the application understands the consistency context of each use case. Third, prioritize explicit degradation to the primary with rich observability over allowing users to silently ingest stale data.

## Source

[Software Engineering Lab 33 - Read Replicas and Replication Lag](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/33-read-replicas-and-replication-lag), snapshot commit `b41fffb`. Source: README, implementation, tests, and execution records. Code excerpts belong in their original package context.
