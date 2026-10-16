---
title: "Cache Invalidation Strategies: Write Patterns, Stampede Mitigations, and TTL Jitter in Go"
published: false
description: "Explore Cache Invalidation Strategies in Go: mechanics, architecture, caveats, and testable invariants through in-memory simulations."
tags: go, testing, performance, softwareengineering
---

# Cache Invalidation Strategies: Write Patterns, Stampede Mitigations, and TTL Jitter in Go

A production system must remain resilient under load. Lab 37 makes that behavior visible with a focused Go simulation and deterministic test assertions.

> Cache invalidation is the hardest problem because consistency is eventual, never instant.

## Distributed Systems' Biggest Dilemma

Cache invalidation remains one of the hardest challenges in distributed systems. While caching shields persistent databases from excessive traffic, maintaining freshness during underlying database updates is critical. Without disciplined invalidation policies, stale reads and data corruption will degrade application reliability and consistency.

## Lab 37 Taxonomy

Lab 37 categorizes cache engineering into write policies and stampede mitigations. Three write strategies dictate latency and consistency trade-offs, while four stampede mitigations protect databases from thundering herds during hot key expiration. All patterns are implemented natively in Go without external infrastructure.

| Strategy | Write Flow | Read Flow | Characteristics |
|---|---|---|---|
| Cache-Aside | DB write → cache Delete | Cache miss → DB query → Set | Lazy, default pragmatic |
| Write-Through | DB write → cache Set (sync) | Cache hit | Strong read-after-write |
| Write-Behind | Cache Set → async queue flush | Cache hit | High write throughput, durability risk |

| Mitigation | Technique | Scope | Mechanism |
|---|---|---|---|
| Singleflight | Request coalescing | Intra-process Go | `singleflight.Group.Do(key, fn)` |
| XFetch | Probabilistic early refresh | Per-request | `-Δ·β·ln(U) > TTL_remaining` |
| SWR | Stale serve + async revalidate | Per-key | `ExpiresAt + staleDelta window` |
| TTL Jitter | Random TTL offset | Cross-key | `TTLWithJitter(base, maxJitter)` |

## Cache-Aside (Lazy Loading) Mechanics & Flow

Cache-Aside operates lazily: reads populate the cache upon a miss, and updates commit to the database before invalidating the cached key. This lazy population saves memory by caching only requested records, while making write-then-invalidate ordering essential for correctness.

```go
func (s *CacheAsideService) Update(ctx, key, val) error {
    s.db.Write(ctx, key, val) // 1. DB first
    s.cache.Delete(key)       // 2. then invalidate
    return nil
}
```

## Cache-Aside Anomalies & Race Conditions

Cache-Aside exhibits a subtle race window under concurrent access. If a reader suffers a cache miss while a writer commits an update, a slow-completing reader can overwrite fresh invalidation with stale data, keeping invalid records in cache until TTL expiration.

## Write-Through Synchronous Consistency

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

## Write-Through Latency Cost & Partial Failures

Write-Through incurs a latency penalty by executing two consecutive synchronous writes for every update. Furthermore, partial failures—such as a committed database write followed by a failed cache set—demand TTL safety nets to prevent prolonged cache inconsistency across application instances.

| Failure Point | Detection | Recovery | Impact |
|---|---|---|---|
| DB write fail | Error return | Retry caller | Write lost |
| Cache set fail | Silent (swallowed) | TTL safety net | Temporary miss |
| Both succeed | - | - | 0 impact, consistent |

## Write-Behind (Write-Back) Extreme Throughput

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

## Write-Behind Durability Risks & Data Loss

Write-Behind trades durability for raw execution speed. Unflushed queue items are permanently lost if the process crashes unexpectedly. In high-concurrency bursts, overflowing queues can silently discard writes unless backed by persistent message brokers like Kafka or Write-Ahead Logs.

## Thundering Herd / Cache Stampede

A cache stampede occurs when a high-traffic hot key expires suddenly. Hundreds of concurrent requests experience simultaneous cache misses and flood the database with identical expensive queries, exhausting connection pools and causing systemic congestion collapse across the infrastructure.

| Workload | Concurrency N | DB Query Count (Naive) |
|---|---|---|
| Local dev | 3 | 3 |
| Staging | 50 | 50 |
| Production hot key | 10,000 | 10,000 (DB overload) |

## Singleflight Request Coalescing

Mitigasi pertama terhadap cache stampede adalah request coalescing menggunakan paket resmi Go `golang.org/x/sync/singleflight`. Pola ini memastikan bahwa untuk kunci yang sama, fungsi pengambilan data ke database hanya dijalankan tepat satu kali oleh satu goroutine pemenang. Goroutine-goroutine lain yang datang bersamaan akan ditahan pada barrier internal dan menunggu selesainya eksekusi goroutine pertama.

```go
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
```

## Singleflight Mechanics & Do() Implementation

Implementasi SingleFlightService menerapkan teknik double-check locking di dalam closure fungsi Do untuk efisiensi maksimal. Sebelum memulai query database yang mahal, goroutine memeriksa kembali isi cache untuk mengantisipasi kemungkinan kunci tersebut telah berhasil diisi oleh eksekusi paralel lain. Singleflight mengelola tabel pemetaan pemanggilan aktif secara internal dan membagikan return value melalui channel bersama. Namun, perlindungan ini terbatas pada level satu proses; sistem multi-pod tetap membutuhkan distributed lock seperti Redis SETNX.

## XFetch Probabilistic Early Expiration

XFetch adalah algoritma mitigasi stampede probabilistik berbasis riset akademis Vattani et al. Setiap kali request membaca cache, algoritma menghitung ambang batas kedaluwarsa dini berdasarkan biaya komputasi query database terakhir, parameter tuning beta, dan bilangan acak seragam. Jika nilai komputasi probabilistik melebihi sisa TTL, request tersebut akan berinisiatif merefresh cache secara proaktif sebelum waktu kedaluwarsa resmi tiba.

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

## XFetch Mathematical Sign Error

Penerapan matematis formula XFetch menuntut pemahaman mendalam terhadap tanda minus logaritma natural. Karena variabel U berada pada rentang nol hingga satu, nilai logaritma natural ln(U) selalu menghasilkan bilangan negatif. Penambahan tanda minus di depan rumus membuat hasil kalkulasi bernilai positif, sehingga memungkinkan perbandingan yang valid terhadap sisa TTL. Kesalahan umum dalam mengabaikan tanda minus akan menghasilkan nilai negatif permanen, menyebabkan kondisi penyegaran awal tidak pernah terpenuhi dan stampede tetap terjadi.

## Stale-While-Revalidate (SWR) Stale Tolerance

Stale-While-Revalidate, diadaptasi dari RFC 5861, menawarkan solusi performa tinggi dengan memperluas toleransi terhadap data kedaluwarsa. Sistem membagi siklus hidup cache menjadi tiga jendela: fresh di mana data disajikan langsung, stale window di mana data usang segera disajikan ke klien sementara revalidasi dipicu di latar belakang, serta hard miss di mana query database harus dieksekusi secara sinkron.

| State | Condition | Action | Return |
|---|---|---|---|
| Fresh | now < ExpiresAt | Return cached | Cached value |
| Stale-1 | ExpiresAt ≤ now < staleUntil | Trigger async revalidate | Stale value |
| Stale-2 | now ≥ staleUntil (past grace) | Synchronous fetch | New DB value (blocking) |
| Miss | Key not in cache | Synchronous fetch | New DB value |

```go
item, ok := s.cache.GetRaw(key)
now := time.Now()
if ok {
    if now.Before(item.ExpiresAt) { return item.Value, nil }
    staleUntil := item.ExpiresAt.Add(s.staleDelta)
    if now.Before(staleUntil) {
        s.triggerRevalidate(key)
        return item.Value, nil
    }
}
// Hard miss / Stale-2
val, err := s.db.Query(ctx, key)
s.cache.Set(key, val, s.ttl, delta)
return val, nil
```

## SWR Async Revalidation & Deduplication

Keberhasilan SWR sangat bergantung pada mekanisme deduplikasi revalidasi latar belakang agar tidak memicu badai permintaan baru ke database. Tanpa pembatas, ratusan request pada data basi akan masing-masing meluncurkan goroutine revalidasi baru. SWRService mengatasi ancaman ini dengan memelihara tabel penanda revalidasi berbasis mutex; hanya satu goroutine yang diizinkan memperbarui database untuk kunci yang sama pada satu waktu.

```go
func (s *SWRService) triggerRevalidate(key string) {
    s.mu.Lock()
    if s.revalidating[key] { s.mu.Unlock(); return }
    s.revalidating[key] = true
    s.mu.Unlock()
    s.revalCount.Add(1)
    go func() {
        defer func() {
            s.mu.Lock()
            delete(s.revalidating, key)
            s.mu.Unlock()
        }()
        ...db.Query → s.cache.Set...
    }()
}
```

## TTL Jitter Anti-Expiration Synchronization

TTL Jitter adalah teknik pertahanan sederhana namun ampuh untuk mencegah sinkronisasi waktu kedaluwarsa lintas kunci cache. Ketika ribuan kunci diisi secara bersamaan saat startup sistem atau deploy aplikasi, penggunaan TTL statis akan menyebabkan seluruh kunci tersebut kedaluwarsa pada detik yang sama persis.

```go
func TTLWithJitter(base, maxJitter time.Duration) time.Duration {
    if maxJitter <= 0 { return base }
    jitter := time.Duration(rand.Int63n(int64(maxJitter)))
    return base + jitter
}
```

## Testing: Cache-Aside vs Write-Through vs Write-Behind

Pengujian TestCachePatterns memvalidasi perilaku ketiga write pattern menggunakan MockDB dengan counter query atomik yang presisi. Pada Cache-Aside, query database hanya bertambah pada cache miss pertama dan tetap stabil pada hit kedua, lalu bertambah lagi setelah operasi update menghapus cache. Pada Write-Through, pembacaan pasca-update langsung mendapati cache hit tanpa query database. Pada Write-Behind, update cache terjadi seketika, dan penulisan database terbukti dieksekusi secara asinkron beberapa milidetik kemudian.

| Test | Pattern | Assertion | Expected DB Queries |
|---|---|---|---|
| Cache-Aside Read | Cache-Aside | Miss → 1 query | 1 |
| Cache-Aside Hit | Cache-Aside | Hit → 0 query | 1 (stays) |
| Cache-Aside Update | Cache-Aside | Write+Delete → next miss | 2 |
| Write-Through Write | Write-Through | Update sync set cache | 1 (no extra query) |
| Write-Behind Flush | Write-Behind | Async flush after sleep | 1 write async |

```bash
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

## Testing: Naive Stampede vs Singleflight

Pengujian mitigasi stampede membuktikan efektivitas Singleflight secara empiris di bawah beban konkurensi tinggi. Ketika dua puluh goroutine diluncurkan secara simultan untuk mengakses kunci panas yang belum ada di cache, implementasi Naive menghasilkan dua puluh query database paralel—memperlihatkan stampede nyata. Sebaliknya, SingleFlightService berhasil menggabungkan seluruh dua puluh panggilan tersebut menjadi tepat satu query database tunggal.

```bash
[Naive] 20 concurrent goroutines -> 20 DB queries (stampede!)
[SF]    20 concurrent goroutines -> 1 DB queries (singleflight), 20 correct results
```

## Testing: XFetch Deterministic Simulation

Pengujian logika XFetch dirancang secara deterministik melalui injeksi fungsi pengacak randFunc tiruan untuk menguji kedua percabangan algoritma. Ketika generator acak menghasilkan nilai tinggi mendekati satu, formula matematis memutuskan tidak perlu recompute, sehingga data disajikan langsung dari cache. Sebaliknya, saat diinjeksi nilai acak yang sangat kecil, rumus menghasilkan nilai komputasi yang melampaui sisa TTL, memicu eksekusi query pembaruan dini ke database secara terbukti dan terukur.

```bash
delta = 100ms, beta = 1.0, ttl = 1s
Initial GET       -> DB queries: 1, val: initial-val
Low rand draw     -> DB queries: 2, val: refreshed-val (early proactive refresh)
High rand draw    -> DB queries: 1, val: fresh-data (served from cache)
```

## Testing: SWR & Race Detection

Pengujian Stale-While-Revalidate memverifikasi bahwa data usang langsung dikembalikan kepada pemanggil tanpa jeda latensi begitu TTL kedaluwarsa namun masih dalam rentang toleransi stale. Sepuluh goroutine konkuren yang membaca kunci basi secara simultan terbukti hanya memicu tepat satu proses revalidasi latar belakang berkat mekanisme deduplikasi. Pengujian race detector Go dengan flag `-race` lulus sempurna tanpa kebocoran data, membuktikan keandalan thread-safety implementasi SWR di lingkungan produksi.

```bash
Initial GET       -> val: v1
GET (stale)       -> val: v1  (stale served immediately)
GET (fresh)       -> val: v2  (async revalidation complete)
```

## 7-Strategy Comparison Matrix

Matriks komparasi tujuh strategi menyajikan panduan evaluasi menyeluruh berdasarkan profil latensi tulis, latensi baca, tingkat konsistensi, dan risiko arsitektural utama. Tidak ada satu strategi tunggal yang sempurna untuk semua skenario; setiap pilihan merupakan kompromi teknik yang sadar.

| Strategy | Write Latency | Read Latency | Consistency | Main Risk |
|---|---|---|---|---|
| Cache-Aside | Sync DB Write | Cache hit (0 RTT) | Eventual | Stampede on mass miss |
| Write-Through | DB+Cache (2 sync ops) | Cache hit (0 RTT) | Strong read-after-write | High write latency |
| Write-Behind | Cache Set + async queue | Cache hit (0 RTT) | Eventual | Data loss on crash/overflow |
| Singleflight | Same | Cache hit after 1st miss | Eventual | No cross-process protection |
| XFetch | Same | Cache hit/proactive refresh | Eventual | Sign error bug |
| SWR | Same | Stale serve (0 blocking) | Eventual | Staleness window |
| TTL Jitter | Same | Cache hit (spread time) | Eventual | No single hot-key protection |

| Additional Column | Cache-Aside | Write-Through | Write-Behind |
|---|---|---|---|
| Memory overhead | Lazy, only read data | All keys written | All keys written |
| Write throughput | Medium | Low | High (async) |
| Read-after-write guarantee | No | Yes | No (eventual) |
| Durability | DB committed | DB committed | Risk before flush |

## Adoption Guidelines & Field Recommendations

Panduan adopsi menetapkan pohon keputusan praktis untuk mengimplementasikan strategi cache di lingkungan produksi. Untuk sebagian besar API read-heavy, Cache-Aside yang dipadukan dengan TTL Jitter adalah pilihan default terbaik. Ketika beban lalu lintas tinggi memicu stampede pada kunci panas, pasang Singleflight untuk layanan monolitik atau kunci terdistribusi untuk arsitektur multi-pod. Terapkan SWR jika sistem menoleransi data usang demi latensi sub-milidetik, dan gunakan Write-Behind hanya jika kehilangan data non-kritis dapat diterima.

| Recommendation | Scenario |
|---|---|
| Cache-Aside + TTL Jitter | Read-heavy APIs (default) |
| Singleflight | In-process hot-key stampede mitigation |
| XFetch beta ≥ 1.0 | Variable compute costs |
| SWR | Ultra-low latency requirements |
| Write-Behind + persistent queue | Non-critical high-throughput writes |
| Write-Through | Financial/audit-sensitive data |

## Closing: Where Performance Meets Data Truth

Cache invalidation: di mana performance bertemu kebenaran data.

| Module Files | |
|---|---|
| internal/cache/store.go | In-memory store implementation |
| internal/cache/repo.go | Repository pattern abstraction |
| internal/cache/patterns.go | Write patterns: Cache-Aside, Write-Through, Write-Behind |
| stampede.go | Stampede mitigations: Singleflight, XFetch, SWR, Jitter |
| tests/cache_test.go | Deterministic unit tests with race detection |
| cmd/demo/main.go | Interactive demonstration CLI |

**Verification & Status**: Source: commit `eac3d80a`. Test passing and verified (Race-Clean Go 1.22). Invariants proven: 20 concurrent goroutines naive → 20 DB queries; with Singleflight → 1 DB query.

**Takeaways**: Cache-aside remains the pragmatic default; write-then-invalidate ordering is critical. Write-through guarantees read-after-write synchronously. Write-behind sacrifices durability for throughput. Singleflight, XFetch, SWR, and TTL Jitter together form defense-in-depth against stampede.

[Software Engineering Lab 37 - Cache Invalidation Strategies](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/37-cache-invalidation-strategies)