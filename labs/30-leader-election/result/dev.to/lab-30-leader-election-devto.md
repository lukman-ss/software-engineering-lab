---
title: "Leader Election with Fencing Tokens: Preventing Split-Brain and State Corruption via Lease TTL and Monotonic Revision"
published: false
description: "Explore Leader Election with Fencing Tokens in Go: mechanics, architecture, caveats, and testable invariants."
tags: go, testing, performance, softwareengineering
---

# Leader Election with Fencing Tokens: Preventing Split-Brain and State Corruption via Lease TTL and Monotonic Revision

A production system must remain resilient under load. Lab 30 makes that behavior visible with a focused Go simulation and deterministic test assertions.

> Time passes invisibly to a paused process; local belief of leadership is never sufficient proof.

## Masalah Nyata: Ilusi Otoritas Tunggal

In production distributed systems, processes rarely fail cleanly or predictably. Runtime anomalies such as stop-the-world garbage collection pauses or transient network partitions can freeze a node for hundreds of milliseconds. During this pause window, the leader's lease TTL expires on the coordinator, triggering an automatic failover to promote a successor. When the old leader resumes execution, it remains completely unaware that time elapsed and continues writing, creating silent data corruption unless strictly intercepted.

## Mengapa Lease Berbasis Waktu Saja Gagal?

The naive approach assumes checking a local boolean flag `if n.state == StateLeader` guarantees mutual exclusion. This assumption fails catastrophically due to the time-of-check to time-of-use (TOCTOU) vulnerability in asynchronous distributed environments. Immediately after the local check evaluates to true, the thread can be suspended by a GC pause or OS preemption for hundreds of milliseconds. When the thread resumes, leadership has already transitioned to another candidate, yet the write proceeds, corrupting storage unless validated at the boundary.

```go
func (n *Node) DoWorkNaive(val string) error {
    if n.state != StateLeader { return ErrNotLeader }
    // GC pause terjadi tepat di sini selama 350ms!
    // Lease TTL 200ms kedaluwarsa di coordinator.
    // Node-B terpilih sebagai leader baru.
    return n.sharedStorage.Save(val) // Menimpa data Node-B!
}
```

## Mental Model: Otoritas Temporal vs Validasi Monotonik

This architectural mental model decouples two distinct concepts: transient leadership authority (lease) and permanent write eligibility (fencing token). The coordinator administers time-bounded leases, issuing a strictly monotonically increasing revision ticket whenever a new leader is recognized. The shared storage layer no longer blindly trusts client identity; instead, it acts as an autonomous gatekeeper that accepts mutations if and only if the presented token strictly exceeds the last recorded token.

```go
type Lease struct {
    Key          string
    HolderID     string
    FencingToken int64         // revision monotonik naik
    ExpiresAt    time.Time
    TTL          time.Duration
}
```

## Alur Koordinasi & Penolakan Mutasi Basi

The coordination pipeline demonstrates how temporal anomalies are handled deterministically. Stages 1 through 2 capture normal execution where Node-A wins election and writes successfully under generation token 1. During stages 3 to 5, Node-A pauses beyond its TTL window, enabling the coordinator to expire the lease and promote Node-B with token 2. When Node-A wakes up in stage 7 and attempts to commit with stale token 1, the storage layer evaluates the fencing invariant and rejects the write.

## Implementasi Gerbang Fencing di Shared Storage

This code block constitutes the primary defense against split-brain corruption across the architecture. The Write method on FencedStorage acquires an exclusive mutex lock to ensure linearizable check-and-set semantics. The validation logic verifies whether the incoming token is less than or equal to lastSeenToken. If so, storage immediately aborts the write with ErrStaleFencingToken. Only strictly greater tokens update the internal watermark and append mutation records.

```go
func (s *FencedStorage) Write(author string, token int64, value string) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    if token <= s.lastSeenToken {
        return fmt.Errorf("%w: token %d <= last seen %d (author: %s)",
            ErrStaleFencingToken, token, s.lastSeenToken, author)
    }
    s.lastSeenToken = token
    s.records = append(s.records, Record{FencingToken: token, Author: author, Value: value, Timestamp: time.Now()})
    return nil
}
```

## Failure Path: Penolakan Stale Write Pascajeda GC

This failure path is extracted directly from the executable lab demo. In steps 4 through 5, Node-A is paused for 350ms, comfortably exceeding the 200ms TTL window and triggering lease expiration. Crucially, the cluster momentarily enters a transient dual-leader state where both nodes view themselves as active leaders. However, when Node-A wakes up and attempts a write with token 1, storage intercepts the call and rejects it because token 1 is less than Node-B's token 2.

## Mekanisme vs Gejala: Evaluasi Solusi Koordinasi

This table evaluates four distributed coordination patterns under network partitions and runtime pauses. Naive heartbeats and bare TTL leases fail to preserve data integrity because both remain susceptible to execution pauses and split-brain phenomena. Distributed locks devoid of generation ordering only suffice for efficiency optimizations like job deduplication. Combining time-bounded leases with monotonic fencing tokens guarantees absolute write safety at the storage boundary.

Pendekatan | Mekanisme Koordinasi | Gejala / Failure Mode | Garansi Split-Brain
|---|---|---|---|
Heartbeat Naif | Ping status via interval | Network partition memicu dual-leader | Tidak Aman (korupsi data)
Lease TTL Murni | Kunci kedaluwarsa otomatis | GC pause memicu stale write pasca-wake | Tidak Aman (TOCTOU gap)
Distributed Lock | Mutex terpusat (Redis SETNX) | Lock lepas sebelum job tuntas | Sebagian (hanya efisiensi)
Lease + Fencing | Monotonic revision + check-and-set | Stale write dihentikan di resource | Aman Mutlak (linierizable)

## 3 Keputusan Desain Kritis Leader Election

Building secure coordination hinges upon three foundational architectural decisions. First, lease lifetime computations must rely strictly on monotonic clock readings to eliminate clock jumps from NTP synchronization. Second, every leadership handover must allocate a strictly monotonic revision token; random UUIDs are unacceptable because they lack causal ordering. Third, the validation gate at the shared resource layer must enforce atomic check-and-set semantics that reject non-increasing tokens.

## Bukti Eksekusi: Siklus Lengkap Election & Fencing

This terminal execution trace verifies that the full leadership lifecycle operates precisely as designed. Node-A initially attains leadership with Token 1 and commits the first record. When Node-A is paused for 350ms, Node-B succeeds it with Token 2 and writes the second record. When Node-A resumes and attempts to write using Token 1, the storage layer rejects the request. The final audit history confirms exactly two legitimate records exist, completely clean of any overwritten state.

## YANG TIDAK DIBUKTIKAN

The technical claims established in this lab are explicitly bounded by the in-process test harness. The coordinator runs in-memory within a single Go runtime rather than across a multi-node Raft consensus cluster such as etcd or ZooKeeper. Real-world physical network degradation like packet loss, storage corruption, or Byzantine behaviors were not exercised. Furthermore, illustrative timing constants (200ms TTL, 50ms renewal) were chosen for rapid test turnaround and must be re-calibrated against production p99 network latency and GC statistics.

## Trade-off Arsitektur: Fencing vs Kompleksitas

Adopting fencing tokens introduces trade-offs between guaranteed data safety and overall system complexity. The primary benefit is absolute immunity from split-brain data corruption and the freedom to promote new leaders without waiting for explicit old-leader termination. Conversely, the operational overhead requires modifying backend storage engines to support check-and-set semantics, accepting failover dead-time proportional to the lease TTL, and engineering client resilience when writes are rejected.

## Validasi Konkurensi: 5 Node Berlomba Rebut Lease

Concurrency stress testing validates implementation correctness when five candidate nodes compete for lease ownership simultaneously. The coordinator leverages an internal sync.Mutex to guarantee linearizable evaluation of lease acquisition and renewal requests. In TestConcurrentElectionRace, after 100ms of concurrent campaigning, exactly one node transitions to LEADER while the remaining four stably remain as FOLLOWERS. Running under Go's race detector flags zero data races across all components.

```go
func TestConcurrentElectionRace(t *testing.T) {
    coord := coordinator.New()
    store := storage.New()
    nodes := make([]*candidate.Node, 5)
    for i := 0; i < 5; i++ {
        nodes[i] = candidate.NewNode(string(rune('A'+i)), "leader-lock", 100*time.Millisecond, 20*time.Millisecond, coord, store)
        nodes[i].Start(ctx)
    }
    time.Sleep(100 * time.Millisecond)
    // Assert: leaderCount == 1 di antara 5 node yang berlomba
}
```

## Before vs After: Evolusi Sistem Koordinasi

This before-and-after comparison highlights the architectural paradigm shift in distributed coordination. Prior to adopting fencing tokens, the topology was exceptionally fragile under execution pauses because storage passively trusted client assertions. As a result, paused legacy leaders overwrote fresh state written by active successors. With the fencing token pattern in place, shared storage functions as an active verification barrier, ensuring mutations from obsolete leadership epochs are systematically rejected.

## Test 1: Akuisisi, Renewal, dan Kedaluwarsa Lease

The first unit test validates the operational lifecycle of leases managed by the coordinator module. Node-1 successfully acquires the initial lease, receiving fencing token 1. While this lease remains valid, Node-2's concurrent acquisition attempt is blocked with ErrLeaseHeld. Heartbeat renewal operations demonstrably extend the lease expiration without perturbing the token value. After pausing for 250ms to exceed the 200ms TTL window, the lease expires cleanly, enabling Node-2 to secure ownership with token 2.

```go
func TestCoordinatorLeaseAcquisitionAndRenewal(t *testing.T) {
    coord := coordinator.New()
    ttl := 200 * time.Millisecond
    lease1, err := coord.Acquire("leader-lock", "node-1", ttl)
    // Assert: lease1.FencingToken == 1
    _, err = coord.Acquire("leader-lock", "node-2", ttl)
    // Assert: errors.Is(err, coordinator.ErrLeaseHeld)
    time.Sleep(250 * time.Millisecond) // Lampaui TTL 200ms
    lease2, err := coord.Acquire("leader-lock", "node-2", ttl)
    // Assert: lease2.FencingToken == 2
}
```

## Test 2: Storage Menolak Token Usang & Duplikat

The second unit test rigorously examines the validation invariants implemented inside FencedStorage. The storage engine sequentially admits tokens 10 and 11. When Node-1 subsequently re-submits a write under token 10, the call is rejected since 10 is strictly less than 11. Similarly, when Node-2 attempts a duplicate write under token 11, storage rejects it because the contract requires strict inequality. Inspecting History() confirms exactly two legitimate records exist, guaranteeing invalid writes never pollute state.

```go
func TestFencedStorageRejectsStaleTokens(t *testing.T) {
    store := storage.New()
    _ = store.Write("node-1", 10, "state-1") // Sukses: 10 > 0
    _ = store.Write("node-2", 11, "state-2") // Sukses: 11 > 10
    err1 := store.Write("node-1", 10, "stale-state")
    // Assert: errors.Is(err1, storage.ErrStaleFencingToken)
    err2 := store.Write("node-2", 11, "duplicate-token")
    // Assert: errors.Is(err2, storage.ErrStaleFencingToken)
    // Assert: len(store.History()) == 2 (mutasi ilegal tidak dicatat)
}
```

## Test 3: Failover Otomatis & Pertahanan Split-Brain

The third test orchestrates the pivotal split-brain scenario solved by this lab. The primary leader is subjected to a simulated stop-the-world pause of 250ms, deliberately exceeding the 150ms TTL threshold. The standby candidate detects the expired lease, claims leadership, and commits state under token 2. When the original leader resumes and attempts to flush data using its stale token 1, the write is stopped cold by ErrStaleFencingToken. This verifies end-to-end split-brain defense under real execution pauses.

```go
// 1. Leader-1 terpilih dan menulis sukses dengan token1
leader1.SimulatePause(250 * time.Millisecond) // Lebih lama dari TTL 150ms
time.Sleep(200 * time.Millisecond)
// 2. Leader-2 mengambil alih kepemimpinan dengan token2 > token1
err = leader2.PerformFencedWrite("leader-2-write") // SUKSES
// 3. Leader-1 bangun dan mencoba menulis dengan token1 usang
err = store.Write(leader1.ID(), token1, "stale-write-after-gc-pause")
// Assert: errors.Is(err, storage.ErrStaleFencingToken)
```

## Studi Kasus Industri: Mekanisme Fencing Nyata

Industry research reveals widespread consensus regarding the necessity of fencing tokens. etcd leverages its 64-bit Raft Revision counter as a native, strictly monotonic fencing token across transactional mutations. Apache ZooKeeper provides zxid timestamps and sequential ephemeral node numbers. The famous debate between Martin Kleppmann and Salvatore Sanfilippo (antirez) over Redis Redlock established a lasting maxim: distributed locks lacking monotonic tokens serve efficiency purposes, never safety-critical correctness.

## Mengapa Lease Wajib Menggunakan Monotonic Clock?

Selecting the correct clock abstraction separates resilient distributed architectures from fragile prototypes. System wall-clocks track human calendar time and can step backward or forward abruptly during NTP synchronizations. If lease calculations depend on wall-clocks, time jumps cause severe anomalies: phantom failovers or indefinite lease starvation. The Go standard library eliminates this vulnerability by transparently embedding hardware monotonic clock readings inside time.Now(), guaranteeing linear duration calculations.

## Checklist Implementasi Leader Election & Fencing

This checklist codifies nine critical rules required before deploying distributed leader election systems to production. Lease TTLs and heartbeat renewal intervals must never be guessed; they require empirical calibration against observed p99 network RTT and runtime GC pause distributions. Coordinators must issue monotonically increasing generation numbers, and shared resources must isolate token validations within atomic check-and-set transactions. Any rejection spike must trigger alerts indicating node stalls or network partitions.

## Anti-Pattern: Kesalahan Desain yang Harus Dihindari

Violating distributed coordination fundamentals compromises enterprise data consistency. The most prevalent anti-pattern is executing leadership validation exclusively in client memory instead of at the storage gate, inviting catastrophic TOCTOU race conditions. Accepting equal tokens by substituting strict inequality with greater-than-or-equal operators likewise shatters epoch exclusivity. This can be reproduced in the lab by relaxing storage.go line 37, causing the test suite to immediately catch the security violation.

## Verifikasi Mandiri: Uji & Jalankan Lab

This laboratory is engineered for immediate, frictionless local verification. The entire codebase is implemented with the Go standard library, eliminating external dependencies like Redis clusters or ZooKeeper daemons. Running go test -v executes four exhaustive test scenarios validating lease lifecycles and split-brain defenses. Executing go test -race proves thread safety under multi-goroutine contention, while go run demo delivers an interactive terminal demonstration of stale write rejection.

```sh
cd labs/30-leader-election
go test -v ./...
go test -race ./...
go run ./cmd/demo
```

## Prinsip Arsitektur: Desain untuk Asinkronitas

These three architectural principles transcend this specific Go implementation. First, recognize that time-bounded leases merely provide timed mutual exclusion ideal for efficiency optimization, never absolute safety. Second, position the enforcement gate at the shared resource layer because paused clients remain fundamentally blind to elapsed time. Third, replace fragile temporal assumptions with explicit generation numbers; strictly monotonic tokens provide deterministic causal ordering immune to distributed execution pauses.

## Source

[Software Engineering Lab 30 - Leader Election with Fencing Tokens](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/30-leader-election), snapshot commit `ba4b6216`. Source: README, implementation, tests, and execution records. Code excerpts belong in their original package context.
