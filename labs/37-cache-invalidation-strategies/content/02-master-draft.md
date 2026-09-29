# Cache Invalidation Strategies

## Problem

Cache mempercepat backend dengan menyimpan hasil query database di memori agar request berikutnya tidak perlu hit DB lagi. Tetapi cache mendatangkan masalah baru:

1. **Staleness** — cache menyimpan data lama; jika DB diperbarui, cache bisa tetap memegang nilai lama.
2. **Cache Stampede (Thundering Herd)** — ketika item cache populer (hot key) kadaluarsa, banyak reader secara bersamaan melihat miss dan merecompute nilai dari DB. Jika setiap recompute butuh waktu `T` detik dan arrival rate `λ`, maka sekitar `λ·T` worker akan bereksekusi bersamaan. Di bawah beban tinggi, ini bisa menyebabkan *congestion collapse*: DB kewalahan, latency melonjak, dan hit rate jatuh ke nol hingga beban turun.
3. **Trade-off consistency vs throughput** — strategi tulis (write-through, write-behind, cache-aside) memilih antara freshness read-after-write dan kecepatan tulis.

Artikel ini menjelaskan tiga pola cache invalidation utama dan empat mitigasi stampede yang telah diverifikasi melalui implementasi Go beserta tes konkurensinya.

## Why This Matters

Tanpa invalidasi yang tepat, cache menjadi sumber inkonsistensi. Tanpa mitigasi stampede, cache dapat berubah dari penyelamat latency menjadi pemicu downtime massal. Pilihan pola cache dan mitigasi stampede merupakan keputusan arsitektural yang memengaruhi latency pembacaan, throughput penulisan, dan ketahanan sistem terhadap spike traffic.

## Mental Model

Pikirkan cache sebagai "memori kerja" di depan database:

- **Cache-Aside** adalah "tanya dulu ke cache, baru ke DB kalau kosong". Ini lazy — hanya data yang dibaca yang masuk cache.
- **Write-Through** adalah "sync tulis ke DB, langsung juga ke cache".读者 melihat nilai terbaru seketika.
- **Write-Behind** adalah "tulis cache dulu, DB nanti". Penulis cepat, tetapi jika cache mati sebelum sinkron, write hilang.
- **Stampede** terjadi ketika banyak orang "serentak tanya ke DB" karena cache telah expired bersamaan.
- **Mitigasi** = singleflight (koalesce satu query), XFetch (refresh awal secara probabilistik), SWR (saat stale, layani data lama sambil revalidate async), jitter (acak TTL agar expiry tidak bareng).

## Core Concept: Three Write Policies

### Cache-Aside (Lazy Loading on Miss)

Pola paling umum untuk read-heavy workload:

1. **Read**: cek cache → miss → query DB → simpan ke cache dengan TTL.
2. **Update**: tulis DB dulu → hapus key dari cache. Urutan penting: DB sebelum cache delete, untuk menghindari window di mana reader bisa re-cache nilai lama setelah cache dihapus.

Kelemahan: tidak menjamin read-after-write consistency. Setelah write+invalidate, reader yang datang tepat pada momen itu bisa melihat cache miss atau nilai stale sementara sebelum key di-populate ulang.

### Write-Through

Sama seperti cache-aside di read, tetapi pada write:
1. Update DB.
2. Update cache dengan nilai baru secara sinkron.

Reader berikutnya akan hit cache dan mendapat nilai terbaru. Trade-off: latency write lebih tinggi karena dua operasi (DB + cache) dilakukan berurutan, dan cache menyimpan semua key yang ditulis (bisa jadi kurang efisien memory jika write jarang dibaca).

Catatan penting: "same write operation" dalam dokumentasi vendor (Microsoft) berarti urutan aplikasi, bukan transaksi ACID lintas sistem. DB dan cache tetap sistem terpisah; jika crash terjadi di antara DB commit dan cache set, TTL bertindak sebagai safety net.

### Write-Behind (Write-Back)

1. Update cache segera.
2. Enqueue request ke queue async; background worker flush ke DB.

Read berikutnya selalu hit cache (kecepatan tinggi). Write throughput meningkat karena aplikasi tidak menunggu DB. Risiko: jika process crash sebelum queue tereksekusi, write yang belum di-flush hilang. Implementasi lab menggunakan `select { case <-queue; default }` untuk drop writes bila queue penuh (demo limitation, bukan rekomendasi produksi).

## Core Concept: Cache Stampede & Mitigations

### Stampede (Dog-Piling)

Definisi resmi (Wikipedia, merujuk Vattani et al. 2015): cascading failure saat item cache populer expired dan banyak goroutine concurrently melihat miss, lalu masing-masing query DB. Contoh ilustratif: halaman butuh 3 detik dirender ulang, arrival rate 10 req/s → 30 concurrent rebuild. Pada load lebih tinggi, setiap attempt timeout → hit rate nol.

### Mitigasi 1: Single-Flight Coalescing

`golang.org/x/sync/singleflight` memberikan koalesensi intra-process: `Do(key, fn)` menjalankan `fn` sekali per key; pendatang yang request key sama akan menunggu hasil dari goroutine pertama dan menerima return value yang sama.

Implementasi lab (`SingleFlightService`):
- Double-check cache hit di dalam closure (`singleflight.Group.Do`), untuk menghindari re-fetch jika key baru saja di-populate.
- Menggunakan `context.WithoutCancel` pada closure agar timeout caller tidak membatalkan in-flight build.
- Batasan: hanya proteksi di dalam satu proses. Untuk multi-pod deployment dibutuhkan distributed lock (mis. Redis SET NX PX) — biaya ekstra write per lock.

### Mitigasi 2: XFetch — Probabilistic Early Expiration

Dari Vattani, Chierichetti & Lowenstein (PVLDB 2015). Setiap request menghitung offset probabilistik:

```
offset = -Δ · β · ln(U)
```

dengan:
- `Δ` = durasi recompute (time.Since start pada DB query terakhir),
- `β` = parameter tuning (praktis default = 1),
- `U` ~ Uniform(0,1).

Recompute dijadwalkan sebelum TTL resmi ketika `now + offset ≥ expiry`, ekuivalen:

```
-Δ · β · ln(U) > TTL_remaining
```

Efeknya: request dengan traffic tinggi (yang mengukur `Δ` besar) cenderung early refresh lebih awal; request jarang justru sering melewati kondisi dan pakai cache stale sampai TTL resmi. Ini menggeser puncak rebuild dari satu titik expiry ke distribusi exponential sepanjang TTL, menghilangkan sinkronisasi.

**Peringatan tanda minus kritis:** formula pedagogis lab awal (`Δ · β · ln(U) > TTL_remaining`) selalu FALSE untuk U ∈ (0,1) karena `ln(U)` negatif, sehingga LHS negatif sementara RHS positif. Implementasi lab telah dikoreksi oleh research revision menjadi `-Δ · β · ln(U) > TTL_remaining` — sesuai temuan research audit C1.

Optimality claim paper diterima dengan kepercayaan bibliografis (DOI metadata OK); bukti matematis dan benchmark empiris paper primer belum diverifikasi karena PDF tidak dapat diparsenoleh tool fetch (render binary stream).

### Mitigasi 3: Stale-While-Revalidate (SWR)

Didefinisikan dalam RFC 5861 (Independent Submission, 2010) sebagai Cache-Control extension: cache dapat melayani response setelah stale sampai `delta` detik, sambil melakukan revalidasi async (non-blocking) ke origin.

Implementasi aplikasi-level (`SWRService`):
- Raw item dikembalikan baik yang sudah expired (selama masih dalam `staleDelta`).
- `triggerRevalidate` dispatch satu goroutine background per key; `revalidating` map + mutex mencegah multiple concurrent revalidation untuk key yang sama (RFC §5 mengamanatkan revalidation request-triggered untuk menghindari amplifikasi attack).
- Jika revalidate gagal, client tetap menerima stale value (fallback).

Perlu dicatat: RFC 5861 adalah Independent Submission (informational), bukan standards-track. Konsepnya valid dan banyak diterapkan oleh CDN/proxy modern, tetapi penerapan di application cache layer adalah analogi, bukan kepatuhan spesifik RFC.

### Mitigasi 4: TTL Jitter

Menambahkan random offset `[0, maxJitter)` ke base TTL (`TTLWithJitter`) — mencegah banyak key expire pada timestamp yang sama (mis. setelah deploy serentak). Manfaat utama: mengurangi sinkronisasi cross-key, bukan single-key contention.

Batasan (inferensi logis): jitter tidak mencegah stampede pada satu hot key yang expired — selama interval TTL masih ada multiple concurrent miss, rebuild akan tetap terjadi (jitter hanya menggeser waktu expiry, bukan jumlah request).

## Architecture

Lab menggunakan struktur package `internal/cache` dengan komponen:

```
internal/cache/
├── store.go     → MemoryCache, Item, TTLWithJitter
├── repo.go      → MockDB (query/write counters, configurable latency)
├── patterns.go  → CacheAsideService, WriteThroughService, WriteBehindService
└── stampede.go  → NaiveStampedeService, SingleFlightService,
                   XFetchService, SWRService, ShouldRecompute
```

`MemoryCache` adalah map[string]Item dengan `sync.RWMutex`. `Get` mengecek expiry via `ExpiresAt`; `GetRaw` mengembalikan item tanpa peduli expiry (diperlukan untuk SWR dan XFetch inspection).

`MockDB` mensimulasikan database dengan `queryDelay` dan counter `queryCount` / `writeCount` atomik.

Implementasi sengaja menghindari Redis/PostgreSQL eksternal agar lab dapat dijalankan standalone.

## Implementation Walkthrough

### Cache-Aside

```go
func (s *CacheAsideService) Get(ctx, key) (string, error) {
    item, err := s.cache.Get(key)
    if err == nil { return item.Value, nil }
    val, err := s.db.Query(ctx, key)
    s.cache.Set(key, val, s.ttl, delta)
    return val, err
}
func (s *CacheAsideService) Update(ctx, key, val) error {
    s.db.Write(ctx, key, val)   // 1. DB dulu
    s.cache.Delete(key)         // 2. lalu invalidate
    return nil
}
```

Urutan write-then-invalidate penting untuk mencegah race di mana reader mengambil nilai lama setelah cache dihapus namun sebelum DB terupdate.

### Write-Through

```go
func (s *WriteThroughService) Update(ctx, key, val) error {
    s.db.Write(ctx, key, val)   // 1. DB write
    s.cache.Set(key, val, s.ttl, delta) // 2. cache update
    return nil
}
```

Read berikutnya selalu hit cache.

### Write-Behind

```go
func (s *WriteBehindService) Update(key, val) {
    s.cache.Set(key, val, s.ttl, ...)
    select {
    case s.writeQueue <- WriteRequest{key, val}:
    default: // queue full — drop (demo limitation)
    }
}
```

`flushWorker` background loop membaca dari `writeQueue` dan memanggil `db.Write`. `Close()` menutup quit channel dan draining queue sebelum keluar.

### SingleFlight

```go
res, err, _ := s.flight.Do(key, func() (interface{}, error) {
    item, err := s.cache.Get(key) // double-check
    if err == nil { return item.Value, nil }
    val, err := s.db.Query(ctx, key)
    s.cache.Set(key, val, s.ttl, delta)
    return val, err
})
```

Hanya satu closure dieksekusi per key; waiter menerima `shared = true`.

### XFetch

```go
func ShouldRecompute(delta time.Duration, beta float64, ttlRemaining time.Duration, u float64) bool {
    if u <= 0 || u >= 1 { return false }
    return -delta.Seconds() * beta * math.Log(u) > ttlRemaining.Seconds()
}
```

`XFetchService.Get` memanggil `GetRaw`, hitung remaining, panggil `ShouldRecompute`, recompute bila true.

### SWR

```go
staleUntil := item.ExpiresAt.Add(s.staleDelta)
if now.Before(staleUntil) {
    s.triggerRevalidate(key)
    return item.Value, nil // serve stale immediately
}
```

Revalidate dilakukan satu goroutine per key, guard via `revalidating` map.

## What the Tests Prove

Run command: `go test -v ./...` (PASS, 0.693s), `go test -race ./...` (PASS, 1.553s). Demo: `go run ./cmd/demo`.

| Test | Claim Diverifikasi |
|------|-------------------|
| `TestCachePatterns/Cache-Aside_Read_&_Write` | Miss→1 query; hit→0 query; write+delete→next miss reload DB |
| `TestCachePatterns/Write-Through_Read_&_Write` | Write sync DB+cache; subsequent read hit cache |
| `TestCachePatterns/Write-Behind_Asynchronous_Flush` | Cache read immediate; DB write muncul setelah sleep |
| `TestStampedeMitigation/Naive_Stampede` | 20 goroutine → queryCount > 1 (stampede) |
| `TestStampedeMitigation/SingleFlight_Coalesces` | 20 goroutine → queryCount == 1; semua return val sama |
| `TestXFetchLogic` | Formula `-delta*beta*ln(u) > remaining` benar; sign error menghasilkan false |
| `TestStaleWhileRevalidate` | Stale returned immediately; async revalidate replaces cache |
| `TestJitter` | `base + [0, maxJitter)` for 100 iterations |

Demo output (`go run ./cmd/demo`):
- Naive 20 goroutines → 20 DB queries.
- SingleFlight 20 goroutines → 1 DB query, 20 correct results.
- XFetch low rand draw → early refresh; high rand draw → no early refresh.
- SWR: v1 served stale, then v2 after async revalidate.

## Failure Scenarios

1. **Naive stampede**: Pada expiry hot key dengan N concurrent readers, DB menerima N query identik. Dengan `queryDelay=20ms`, total waktu tunggu naik linear (N×20ms wall-clock karena goroutines konkuren menunggu I/O).
2. **Write-behind data loss**: Crash sebelum `flushWorker` mengeksekusi `writeQueue` item → write hilang. Demo tidak mensimulasikan ini karena in-memory.
3. **Formula sign error**: Tanpa minus sign, `-ln(U)` tidak terbentuk; kondisi selalu false; early refresh tidak pernah terjadi.
4. **SWR revalidation storm**: Jika tidak ada `revalidating` guard, setiap request pada stale key akan spawn goroutine baru → amplifikasi load pada DB.
5. **Queue overflow write-behind**: `select-default` drop silently; aplikasi lain tidak tahu write hilang.

## Production Considerations

- **Singleflight scope**: hanya efektif intra-process. Untuk multi-replica, pertimbangkan distributed lock (Redis SET NX PX) atau shared cache lock, dengan trade-off extra write dan complexity.
- **Write-behind durability**: gunakan persistent queue (Kafka, WAL, file log) sebelum flush ke DB. Drop-on-full bukan solusi produksi.
- **XFetch beta tuning**: `beta=1` praktis default; `beta>1` mendorong earlier refresh (lebih agresif), `beta<1` sebaliknya.
- **SWR stale window sizing**: window terlalu kecil membuat revalidate jarang tersentuh (beberapa request tetap blocking); terlalu besar memperpanjang staleness.
- **TTL jitter range**: jitter yang terlalu kecil tidak efektif deskronisasi; terlalu besar menambah variability pada staleness SLA.
- **Monitor counters**: `MockDB.queryCount` / `writeCount` harus diganti metric aktual (Prometheus histogram/counter) di produksi.

## Common Mistakes

1. **Implementasi XFetch tanpa minus sign** → formula salah arah, early refresh tidak pernah trigger.
2. **Delete cache sebelum DB write pada cache-aside** → window inkonsistensi di mana reader re-cache nilai lama.
3. **Menganggap singleflight menyelesaikan stampede di multi-process** → perlu distributed coordination.
4. **Memanfaatkan background job unconditional untuk SWR revalidation** → melanggar saran RFC 5861 §5 (amplification risk).
5. **Menganggap jitter menyelesaikan stampede single hot key** → jitter deskronisasi cross-key, bukan koalesensi single-key.
6. **Menganggap write-through menjamin ACID cross-system** → hanya sequential application-level write.

## Checklist

- [ ] Pilih write policy berdasarkan SLA freshness vs throughput.
- [ ] Untuk read-heavy + hot keys, pasang singleflight atau XFetch.
- [ ] Set TTL dengan jitter jika banyak key di-prime serentak.
- [ ] Terapkan SWR jika toleransi staleness ada.
- [ ] Uji stampede mitigation dengan benchmark konkurensi sebelum deploy.
- [ ] Instrumentasi counter (queries, writes, revalidations) untuk observability.
- [ ] Dokumentasikan bounded drop / overflow behavior write-behind; ponytail: expose overflow metric for production.

## Key Takeaways

1. Cache-aside adalah default pragmatis; write-then-invalidate urutannya kritis.
2. Write-through memberi read-after-write freshness sinkron; write-behind mengorbankan durability demi throughput.
3. Cache stampede adalahFailure mode yang dapat menurunkan hit rate ke nol; mitigasi wajib untuk hot keys di bawah load tinggi.
4. Singleflight intra-process cukup untuk single-process Go services; distributed lock diperlukan untuk multi-pod.
5. XFetch menggeser distribusi rebuild ke forward-time dengan probabilitas proporsional terhadap `Δ` dan `ln(1/U)`.
6. SWR memanfaatkan window staleness untuk melayani latency rendah sambil merevalidasi async request-triggered.
7. Jitter efektif deskronisasi TTL cross-key, bukan penyelesaian stampede single key.
8. Implementasi lab bersifat pedagogis (in-memory, synthetic concurrency) — transfer ke produksi memerlukan pertimbangan durability, observability, dan distributed coordination.

## Sources

**Research:**
- `research/05-report.md` — ringkasan temuan, Finding 1-8.
- `research/02-sources.md` — daftar 11 sumber.
- `research/03-evidence.md` — Evidence 1-20 detail.
- `research/04-contradictions.md` — C1-C6 tensions.
- `research-audit/07-verdict.md` — APPROVED.

**Implementation:**
- `internal/cache/store.go` — MemoryCache, Item, TTLWithJitter.
- `internal/cache/repo.go` — MockDB.
- `internal/cache/patterns.go` — CacheAside, WriteThrough, WriteBehind.
- `internal/cache/stampede.go` — Naive, SingleFlight, XFetch, SWR, ShouldRecompute.
- `cmd/demo/main.go` — runnable comparison.

**Tests:**
- `tests/cache_test.go` — 5 test functions, 8 sub-tests.

**Audit:**
- `engineering-audit/06-verdict.md` — APPROVED.
- `engineering-audit/02-code-audit.md` — 12 findings (all LOW).
- `engineering-audit/03-test-audit.md` — test matrix.
- `engineering-audit/04-docs-vs-code.md` — MATCH across all rows.
- `engineering-revision/03-revision-result.md` — no revisions needed.

**External:**
- Microsoft Learn — Cache-Aside Pattern (Source 01).
- IETF RFC 5861 — HTTP Cache-Control Extensions for Stale Content (Source 02).
- Go singleflight pkg.go.dev (Source 03).
- Wikipedia — Cache stampede; Thundering herd; Cache (computing) — Write policies; Cache invalidation (Sources 04-07).
- Vattani, Chierichetti, Lowenstein — "Optimal Probabilistic Cache Stampede Prevention", PVLDB 8(8):886-897, 2015 (Sources 08-09, DOI 10.14778/2757807.2757813).
