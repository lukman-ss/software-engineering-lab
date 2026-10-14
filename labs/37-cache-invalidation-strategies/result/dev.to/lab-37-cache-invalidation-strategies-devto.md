---
title: "Cache Invalidation Strategies: Write Patterns, Stampede Mitigations, and TTL Jitter in Go"
published: false
description: "Explore Cache Invalidation Strategies in Go: mechanics, architecture, caveats, and testable invariants."
tags: go, testing, performance, softwareengineering
---

# Cache Invalidation Strategies: Write Patterns, Stampede Mitigations, and TTL Jitter in Go

A production system must remain resilient under load. Lab 37 makes that behavior visible with a focused Go simulation and deterministic test assertions.

> Cache invalidation is the hardest problem because consistency is eventual, never instant.

## Dilema Terbesar Sistem Terdistribusi

Cache invalidation remains one of the hardest challenges in distributed systems. While caching shields persistent databases from excessive traffic, maintaining freshness during underlying database updates is critical. Without disciplined invalidation policies, stale reads and data corruption will degrade application reliability and consistency.

## Taksonomi Lab 37

Lab 37 categorizes cache engineering into write policies and stampede mitigations. Three write strategies dictate latency and consistency trade-offs, while four stampede mitigations protect databases from thundering herds during hot key expiration. All patterns are implemented natively in Go without external infrastructure.

Pola | Write Flow | Read Flow | Karakteristik
|---|---|---|---|
Cache-Aside | DB write → cache Delete | Cache miss → DB query → Set | Lazy, default pragmatic
Write-Through | DB write → cache Set (sync) | Cache hit | Strong read-after-write
Write-Behind | Cache Set → async queue flush | Cache hit | High write throughput, durability risk

Mitigasi | Teknik | Scope | Formula / Mekanisme
|---|---|---|---|
Singleflight | Request coalescing | Intra-process Go | singleflight.Group.Do(key, fn)
XFetch | Probabilistic early refresh | Per-request | -Δ·β·ln(U) > TTL_remaining
SWR | Stale serve + async revalidate | Per-key | ExpiresAt + staleDelta window
TTL Jitter | Random TTL offset | Cross-key | TTLWithJitter(base, maxJitter)

## Cache-Aside (Lazy Loading) - Mekanisme & Flow

Cache-Aside operates lazily: reads populate the cache upon a miss, and updates commit to the database before invalidating the cached key. This lazy population saves memory by caching only requested records, while making write-then-invalidate ordering essential for correctness.

```go
func (s *CacheAsideService) Update(ctx, key, val) error {
    s.db.Write(ctx, key, val) // 1. DB dulu
    s.cache.Delete(key)      // 2. baru invalidate
    return nil
}
```

## Cache-Aside - Anomali & Race Condition

Cache-Aside exhibits a subtle race window under concurrent access. If a reader suffers a cache miss while a writer commits an update, a slow-completing reader can overwrite fresh invalidation with stale data, keeping invalid records in cache until TTL expiration.

## Write-Through - Konsistensi Sinkron

Write-Through enforces strong read-after-write consistency by synchronizing database commits directly with cache updates. Subsequent readers immediately hit the cache and receive the freshest data without querying the database, eliminating the stale window inherent in lazy caching.

```go
func (s *WriteThroughService) Update(ctx, key, val) error {
    start := time.Now()
    if err := s.db.Write(ctx, key, val); err != nil {
        return fmt.Errorf("db write failed: %w", err)
    }
    delta := time.Since(start)
    s.cache.Set(key, val, s.ttl, delta)
    return nil
}
```

## Write-Through - Biaya Latensi & Partial Failure

Write-Through incurs a latency penalty by executing two consecutive synchronous writes for every update. Furthermore, partial failures—such as a committed database write followed by a failed cache set—demand TTL safety nets to prevent prolonged cache inconsistency across application instances.

Failure Point | Deteksi | Recovery | Impact
|---|---|---|---|
DB write fail | Error return | Retry caller | Write lost
Cache set fail | Silent (swallowed) | TTL safety net | Temporary miss
Both succeed | - | - | 0 impact, consistent

## Write-Behind (Write-Back) - Throughput Ekstrem

Write-Behind maximizes write throughput by updating memory cache immediately and delegating database persistence to an asynchronous queue drained by background workers. Write operations return instantly with near-zero latency, making it ideal for high-velocity telemetry and write-heavy workloads.

```go
func (s *WriteBehindService) Update(key, val string) {
    s.cache.Set(key, val, s.ttl, 1*time.Millisecond)
    select {
    case s.writeQueue <- WriteRequest{Key: key, Value: val}:
    default: // queue full — drop (demo limitation)
    }
}
```

## Write-Behind - Risiko Durabilitas & Data Loss

Write-Behind trades durability for raw execution speed. Unflushed queue items are permanently lost if the process crashes unexpectedly. In high-concurrency bursts, overflowing queues can silently discard writes unless backed by persistent message brokers like Kafka or Write-Ahead Logs.

## Thundering Herd / Cache Stampede

A cache stampede occurs when a high-traffic hot key expires suddenly. Hundreds of concurrent requests experience simultaneous cache misses and flood the database with identical expensive queries, exhausting connection pools and causing systemic congestion collapse across the infrastructure.

Workload | Concurrency N | DB Query Count (Naive)
|---|---|---|
Local dev | 3 | 3
Staging | 50 | 50
Production hot key | 10,000 | 10,000 (DB overload)

## Singleflight - Request Coalescing

Singleflight provides intra-process request coalescing via Go official sync library. When multiple goroutines request an expired hot key, Singleflight executes the underlying database fetch exactly once while suspending other callers until the shared result is broadcast to all waiters.

```go
import "golang.org/x/sync/singleflight"

type SingleFlightService struct {
    cache  *MemoryCache
    db     *MockDB
    ttl    time.Duration
    flight singleflight.Group
}
```

## Singleflight - Kode & Mekanika Do()

Inside the Singleflight closure, a double-check pattern inspects the cache once more before triggering the database fetch, ensuring optimal efficiency. While Singleflight eliminates intra-process stampedes, multi-node deployments still require distributed coordination such as Redis distributed locks.

```go
func (s *SingleFlightService) Get(ctx, key) (string, error) {
    item, err := s.cache.Get(key) // outer cache check
    if err == nil { return item.Value, nil }
    res, err, _ := s.flight.Do(key, func() (interface{}, error) {
        // Double-check cache inside flight execution
        item, err := s.cache.Get(key)
        if err == nil { return item.Value, nil }
        start := time.Now()
        val, err := s.db.Query(ctx, key)
        if err != nil { return "", err }
        delta := time.Since(start)
        s.cache.Set(key, val, s.ttl, delta)
        return val, nil
    })
    if err != nil { return "", err }
    return res.(string), nil
}
```

## XFetch - Probabilistic Early Expiration

XFetch employs probabilistic early expiration derived from academic research by Vattani et al. By evaluating background fetch duration and uniform random draws against remaining TTL, high-traffic keys trigger proactive refreshes before official expiration, smoothing out concurrent rebuild spikes.

```go
func ShouldRecompute(delta time.Duration, beta float64,
    ttlRemaining time.Duration, u float64) bool {
    if u <= 0 || u >= 1 { return false }
    deltaSec := delta.Seconds()
    ttlRemainingSec := ttlRemaining.Seconds()
    expiryCompute := -deltaSec * beta * math.Log(u)
    return expiryCompute > ttlRemainingSec
}
```

## XFetch - Tanda Negatif Matematika

Mathematical precision in XFetch requires retaining the negative sign in front of the natural logarithm. Because ln(U) is strictly negative for random values between 0 and 1, omitting the minus sign forces the expression negative, silently disabling proactive early recomputation.

## Stale-While-Revalidate (SWR) - Toleransi Data Basi

Stale-While-Revalidate, inspired by RFC 5861, balances latency and freshness across three lifecycle phases. Within the stale grace window, clients receive immediate cached responses while background revalidation fetches updated records asynchronously, guaranteeing low end-user latency.

State | Kondisi | Action | Return
|---|---|---|---|
Fresh | now < ExpiresAt | Return cached | Cached value
Stale-1 | ExpiresAt ≤ now < staleUntil | Trigger async revalidate | Stale value
Stale-2 | now ≥ staleUntil (past grace) | Synchronous fetch | New DB value (blocking)
Miss | Key not in cache | Synchronous fetch | New DB value

```go
item, ok := s.cache.GetRaw(key)
now := time.Now()
if ok {
    if now.Before(item.ExpiresAt) { return item.Value, nil }       // Fresh
    staleUntil := item.ExpiresAt.Add(s.staleDelta)
    if now.Before(staleUntil) {
        s.triggerRevalidate(key)                                   // Stale-1
        return item.Value, nil
    }
}
// Hard miss / Stale-2 → synchronous fetch
val, err := s.db.Query(ctx, key)
s.cache.Set(key, val, ttl, delta)
return val, nil
```

## SWR - Async Revalidation & Deduplikasi

SWR requires robust deduplication to prevent background revalidation storms. Using a mutex-protected tracking map, SWRService guarantees that only a single background goroutine refreshes the database per key at any given time, gracefully falling back to stale data upon upstream errors.

```go
func (s *SWRService) triggerRevalidate(key string) {
    s.mu.Lock()
    if s.revalidating[key] { s.mu.Unlock(); return }
    s.revalidating[key] = true
    s.mu.Unlock()
    s.revalCount.Add(1)
    go func() { // async fetch & set
        defer func() {
            s.mu.Lock()
            delete(s.revalidating, key)
            s.mu.Unlock()
        }()
        ...db.Query → s.cache.Set...
    }()
}
```

## TTL Jitter - Anti-Sinkronisasi Expiry

TTL Jitter adds a random time offset to base TTL values, preventing cross-key expiration synchronization. By desynchronizing bulk-primed keys across a bounded time window, Jitter eliminates simultaneous mass expirations that often overwhelm backend storage following application deployments.

```go
func TTLWithJitter(base, maxJitter time.Duration) time.Duration {
    if maxJitter <= 0 { return base }
    jitter := time.Duration(rand.Int63n(int64(maxJitter)))
    return base + jitter
}
```

## Pengujian: Cache-Aside vs Write-Through vs Write-Behind

Unit testing via TestCachePatterns verifies all three write policies deterministically using MockDB query counters. Invalidation cycles, synchronous write-through hits, and asynchronous write-behind flushes are validated against exact query count assertions, confirming predictable IO behavior.

Test | Pattern | Assertion | Expected DB Queries
|---|---|---|---|
Cache-Aside Read | Cache-Aside | Miss → 1 query | 1
Cache-Aside Hit | Cache-Aside | Hit → 0 query | 1 (stays)
Cache-Aside Update | Cache-Aside | Write+Delete → next miss | 2
Write-Through Write | Write-Through | Update sync set cache | 1 (no extra query)
Write-Behind Flush | Write-Behind | Async flush after sleep | 1 write async

```go
[Cache-Aside]
  1st GET (cache miss)  -> DB queries: 1, val: Laptop
  2nd GET (cache hit)   -> DB queries: 1, val: Laptop
  After update & GET    -> DB queries: 2, val: Laptop Pro

[Write-Through]
  1st GET (cache miss)  -> DB queries: 1
  After update & GET    -> DB queries: 1, val: Phone Pro

[Write-Behind]
  After Update, GET from cache -> DB writes: 0, val: Tablet
  After 50ms flush delay       -> DB writes: 1
```

## Pengujian: Naive Stampede vs Singleflight

Concurrency tests empirically contrast Naive Stampede against Singleflight. Under twenty concurrent goroutines requesting a cold hot-key, the naive service triggers twenty redundant database queries, whereas Singleflight successfully collapses all twenty requests into exactly one database invocation.

```go
[Naive] 20 concurrent goroutines -> 20 DB queries (stampede!)
[SF]    20 concurrent goroutines -> 1 DB queries (singleflight), 20 correct results
```

## Pengujian: XFetch Deterministic Simulation

Deterministic XFetch testing leverages injected random generators to validate decision paths reliably. High uniform draws preserve cache hits, while simulated low draws reliably satisfy the probabilistic formula, triggering early proactive recomputations before hard TTL expiration.

```go
Initial GET       -> DB queries: 1, val: initial-val
Low rand draw     -> DB queries: 2, val: refreshed-val (early proactive refresh)
High rand draw    -> DB queries: 1, val: fresh-data (served from cache, no early refresh)
```

## Pengujian: SWR & Race Detection

SWR test suites confirm immediate delivery of stale data alongside single-goroutine background revalidation. Under ten concurrent readers, deduplication prevents revalidation storms, while Go race detector passes cleanly, validating strict thread-safety across concurrent cache accesses.

```go
Initial GET       -> val: v1
GET (stale)       -> val: v1  (stale served immediately)
GET (fresh)       -> val: v2  (async revalidation complete)
```

## Matriks Komparasi 7 Strategi

The comprehensive comparison matrix evaluates all seven strategies across write latency, read latency, consistency guarantees, and operational failure risks. Engineering teams must weigh strong read-after-write consistency against write latency and durability risks to identify optimal caching architectures.

Strategi | Latensi Tulis | Latensi Baca | Konsistensi | Risiko Utama
|---|---|---|---|---|
Cache-Aside | Sync DB Write | Cache hit (0 RTT) | Eventual read-after-write | Stampede pada miss massal
Write-Through | DB+Cache (2 sync ops) | Cache hit (0 RTT) | Strong read-after-write | Latency penulisan tinggi
Write-Behind | Cache Set + async queue | Cache hit (0 RTT) | Eventual consistency | Data loss pada crash/overflow
Singleflight | Sama | Cache hit setelah 1st miss | Eventual | Tidak proteksi cross-process
XFetch | Sama | Cache hit/proactive refresh | Eventual | Sign error bug (formulasi salah)
SWR | Sama | Stale serve (0 blocking) | Eventual | Staleness window, revalidation storm
TTL Jitter | Sama | Cache hit (tersebar waktu) | Eventual | Tidak proteksi single hot key

Kolom tambahan | Cache-Aside | Write-Through | Write-Behind
|---|---|---|---|
Memory overhead | Lazy, hanya data yang dibaca | Semua key yang ditulis | Semua key yang ditulis
Write throughput | Medium | Low | High (async)
Read-after-write guarantee | No | Yes | No (eventual)
Durability | DB committed | DB committed | Risk of loss before flush

## Panduan Adopsi & Rekomendasi Lapangan

Adoption guidelines provide a decision tree mapping workload characteristics to caching strategies. Cache-Aside with Jitter serves as the pragmatic baseline, Singleflight guards hot keys, SWR satisfies ultra-low latency requirements, and Write-Behind suits non-critical high-throughput ingestion pipelines.

## Source

[Software Engineering Lab 37 - Cache Invalidation Strategies](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/37-cache-invalidation-strategies), snapshot commit `eac3d80a`. Source: README, implementation, tests, and execution records. Code excerpts belong in their original package context.
