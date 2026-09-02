# Lab 04 — Caching: Mengurangi Latency, Tapi Apa Biayanya?

> **Mental Model**: Caching adalah trade-off antara **latency** (cepat) dan **consistency** (benar). Senior engineer memilih teknik caching yang tepat untuk workload-nya, bukan menerapkan caching di mana-mana.

---

## Learning Objectives

Setelah menyelesaikan lab ini, Anda akan memahami:

1. Problem / Why Cache Exists
2. Source of Truth
3. Kapan Cache Layak / Tidak Layak
4. Cache Aside
5. TTL & Staleness
6. Cache Invalidation
7. Write Through
8. Cache Stampede
9. Singleflight
10. Distributed Lock
11. TTL Jitter / Background Refresh
12. Cache Key Design / Multi-Tenant
13. Cache vs Session
14. Cache Failure
15. Memory / Cardinality / Eviction
16. Laravel Cache::remember Caveats
17. Permission & Security Caching
18. Query Optimization → Index → Cache
19. Decision Framework / Rule of Thumb
20. CMMS Examples

---

## 1. Problem / Why Cache Exists

Dashboard workshop menampilkan statistik dilihat oleh banyak user. Tanpa cache, setiap request menghasilkan *heavy aggregation* ke database:

```
[500 concurrent users] → Dashboard Request
                              ↓
                    Database: 6 queries + join/aggregation per request
                              ↓
                           3000 total DB queries
```

```
┌────────┐     ┌──────────┐     ┌───────┐
│ Client │ ──► │  YourApp │ ──► │  DB   │
└────────┘     └──────────┘     └───────┘
                     │               ✗ Heavy query every request
                     │               ✗ Latency ~50ms+
                     ▼
                     │
           ┌─────┴─────┐
           │   Redis   │  ← Cache
           └───────────┘
                     │
                     ▼
                  ~1ms latency
```

---

## 2. Source of Truth

| Layer | Technology | Role |
|-------|------------|------|
| **Primary** | PostgreSQL | Permanent, ACID, can rebuild cache |
| **Cache** | Redis | Derived, TTL-bound, on-demand rebuild |

**Prinsip Utama:** Cache correctness harus dibangun dengan asumsi cache sewaktu-waktu dapat kosong atau di-flush. Database selalu menjadi otoritas final.

---

## 3. Kapan Cache Layak / Tidak Layak?

### Apa yang Sebaiknya Di-cache?

Cache cocok untuk data yang:
1. **Mahal dihitung** — Query expensive (join, aggregation)
2. **Sering dibaca** — Read-heavy workload
3. **Jarang berubah dibanding frekuensi baca**
4. **Toleran terhadap staleness**

Contoh: Dashboard operasional, laporan agregasi, master supplier, katalog jenis service.

### Apa yang Sebaiknya Tidak Di-cache?

Data yang **butuh strong consistency** atau berubah **sangat cepat**:
- Saldo wallet / kas
- Status pembayaran
- OTP / Token sekali pakai

**Nuance:** Data real-time **bisa** menggunakan cache sebagai optimization layer asalkan:
- Database tetap source of truth
- Correctness (validasi akhir) tidak bergantung pada cache
- TTL sangat pendek

Contoh: Stock display cache 3s untuk UI, tetapi validasi final ke DB dengan `SELECT FOR UPDATE` saat transaksi checkout.

---

## 4. Cache Aside

Lab ini menggunakan pattern **Cache Aside**.

```
[App: GET cache]
      ↓ miss
[App: Query database]
      ↓
[App: SET cache]
      ↓
[App: Return data]
```

**Karakteristik Cache Aside:**
- Cache miss diselesaikan oleh kode aplikasi kita.
- Aplikasi tahu bahwa ada dua storage layer (Redis & DB).
- Memberikan kontrol lebih besar atas penanganan error (misal: Redis mati, aplikasi masih bisa query DB).

---

## 5. TTL & Staleness

### Apakah Stale Data Masalah?

Pertanyaan kuncinya bukan "apakah data berubah?" tetapi **"berapa lama stale data dapat diterima?"**

| Data | Max Staleness | Reasonable? |
|------|--------------|-------------|
| Dashboard statistik | 30s–2min | ✅ Ya, operational metrics |
| Stock display | 1–5s | ✅ Ya, untuk UI/UX saja |
| Saldo wallet | 0s | ❌ Tidak, risiko audit |

---

## 6. Cache Invalidation

### Race Condition pada Invalidation (DELETE)

**UNSAFE: DELETE cache sebelum DB COMMIT**

```
Timeline:
T1 Writer: DELETE cache
T2 Reader: cache MISS → baca DB (lama, belum committed)
T3 Reader: SET nilai lama ke cache
T4 Writer: DB COMMIT nilai baru

Hasil: Database = nilai baru, Cache = nilai lama (stale)
Catatan: Ini adalah STALE CACHE BUKAN DATA LOSS. Business data di DB tetap aman.
```

**Lebih aman: COMMIT → DELETE**

```
Timeline:
T1 Writer: DB COMMIT nilai baru
T2 Writer: DELETE cache
T3 Reader: cache MISS → baca DB baru → SET cache baru

Hasil: Database = baru, Cache = baru
```

**NAMUN: COMMIT → DELETE tetap tidak memberikan strong consistency**

```
Timeline:
T1 Reader: cache MISS → baca DB lama
T2 Reader: SET nilai lama ke cache
T3 Writer: DB COMMIT nilai baru
T4 Writer: DELETE cache
T5 Reader: SET hasil lama ke cache

Hasil: Database = baru, Cache = lama (stale)
```

**Kesimpulan:** Cache-aside adalah **eventually-consistent optimization**, bukan strong consistency primitive.

### Strategi Invalidation

1. **TTL Expiration** — Data dibiarkan stale hingga waktu habis
2. **Explicit Delete** — DELETE key pasca-commit
3. **Event-Driven Invalidation** — Publish event pasca-commit untuk memicu invalidasi (pub/sub)

### Tag / Index-based Invalidation

**Penting:** Redis core tidak menyediakan generic cache tagging abstraction. Banyak perbufet menyamakan SCAN+pattern delete dengan tag implementation. Mereka adalah hal yang berbeda.

#### Pattern Invalidation (SCAN-based)

```
SCAN namespace/pattern
      ↓
temukan matching keys
      ↓
DELETE / UNLINK
```

**Karakteristik:**
- SCAN lebih production-friendly daripada KEYS
- Tetap memiliki operational cost
- **Bukan atomic operation** atas seluruh matching dataset

#### Tag / Reverse Index Invalidation

```
tag:product:123
      ↓
product:123
products:list
products:popular
```

Saat product berubah: baca members dari tag index, invalidate semua key yang terhubung.

**Trade-off:**
- Metadata tambahan yang harus dipelihara
- Synchronization overhead
- Stale reverse index mungkin terjadi

### Cache Key Versioning

**Schema/Deployment Versioning:**
```
product:v1:123  → format cache lama
product:v2:123  → format cache baru (migration)
```
- Gunakan untuk: serialization format berubah, schema cache berubah, deployment baru membutuhkan namespace terpisah

**Data Revision/Version (Data Versioning):**
```
product:123:revision:42
```
- Butuh cara mengetahui revision terbaru (metadata pointer)
- Version bump TIDAK selalu atomic invalidation

---

## 7. Write Through

### Flow

```
Request update
      ↓
Update Database      ← Authoritative
      ↓
DB Commit Success
      ↓
Update Cache         ← Best-effort sync
      ↓
   Return success
```

### Karakteristik & Risiko Partial Failure

**Tidak Atomic:** Database PostgreSQL dan Redis adalah sistem berbeda. Urutan DB COMMIT → Redis SET tidak atomic.

**Scenario A: Process Crash**
```
DB COMMIT sukses
↓
process crash
↓
Redis SET tidak pernah terjadi
↓
cache lama masih tersimpan
```

**Scenario B: Concurrent Writer Race Condition**
```
Writer A commit value A
Writer B commit value B
Writer B Redis SET B
Writer A terlambat Redis SET A

Hasil: Database = B, Cache = A (stale)
```

### Perbandingan: Cache-Aside vs Write-Through

| Aspek | Cache-Aside | Write-Through |
|-------|-------------|---------------|
| Read Path | App → Check cache → Miss → DB → SET cache | App → Check cache → HIT → Return |
| Write Path | App → DB write → Explicit cache invalidation | App → DB write → Update cache (best-effort) |
| Stale Window | Perlu (hingga invalidation TTL) | Lebih pendek, tapi tetap ada |
| Partial Failure Risk | Low (DB write saja, cache optional) | Tinggi (cache tidak pernah diupdate mungkin) |
| TTL Safety Net | Wajib | Wajib |
| Use Case | Read-heavy, write-rare, toleran stale | Read-heavy, write-moderate, butuh read-after-write |

**Tidak ada pola yang "universally better".** Pilih berdasarkan business requirements.

---

## 8. Cache Stampede

```
cache expires
      ↓
1000 concurrent requests datang
      ↓
1000 request cache miss
      ↓
1000 parallel DB queries
      ↓
Database overload / crash
```

Stampede terjadi ketika key populer expire dan traffic besar masuk bersamaan.

---

## 9. Singleflight

Mengatasi stampede **dalam satu instance/process** menggunakan `golang.org/x/sync/singleflight`.

---

## 10. Distributed Lock

Mengatasi stampede **multi-instance** menggunakan `SETNX` dengan TTL expiration dan unique owner/token.

---

## 11. TTL Jitter & Background Refresh

### TTL Jitter
Tambahkan pengacak kecil ke TTL (misal: 60s ± 15s) agar batch cache keys tidak kadaluarsa tepat bersamaan.

### Background Refresh (Stale-While-Revalidate)
Mem-refresh data secara background sebelum key expired. Client tetap dilayani dengan data cache.

---

## 12. Cache Key Design / Multi-Tenant

Format canonical: `{app}:{tenant}:{branch}:{resource}:{dimension}`

Contoh: `cmms:tenant:42:branch:7:dashboard:2026-09-01`

**Dimensi yang Wajib Masuk Key:**
- `tenant:42` — Isolasi data antar konsumen
- `branch:7` — Resource scoping
- `date:...` — Scope business date (timezone-aware)

---

## 13. Cache vs Session

### CACHE

**Tujuan:**
- Optimization / reuse
- Menghindari computation / I/O berulang

**Scope:**
- **Global** — Data untuk semua user
- **Tenant scoped** — `tenant:42:*`
- **Branch scoped** — `branch:7:*`
- **User scoped** — `user:123:*`
- **Query scoped** — Hasil query tertentu

**Contoh cache per-user yang VALID:**
```
cmms:tenant:42:user:123:permissions:revision:7
```

**Syarat aman:**
- Key isolation benar
- Invalidation benar
- Security requirement dipenuhi

### SESSION

**Tujuan:**
- Menyimpan state user/session
- Lifecycle mengikuti login/session

**Contoh:** Login state, CSRF token, shopping cart server-side, temporary preference

**Implementasi dapat menggunakan:** Redis, Database, Memory, Cookie, datastore lain.

### Perbedaan Utama

| Aspek | Cache | Session |
|-------|-------|---------|
| Tujuan Utama | Performance optimization | State persistence |
| Acceptable Staleness | Bisa stal (detik-menit) | Harus konsisten |
| Failure Mode | Degradasi performa | Login hilang / cart hilang |
| TTL | Singkat | Panjang |

---

## 14. Cache Failure / Graceful Degradation

Prinsip: **Cache adalah optimization layer, bukan dependency correctness utama.**

- **Degraded Performance:** Aplikasi harus tetap berfungsi meski lambat.
- **Circuit Breaker:** Mencegah aplikasi stuck menunggu Redis Timeout.
- **Observability:** Log `cache_miss` dan `cache_error` terpisah.

**Catatan tentang Hit Ratio:** Hit ratio tidak universal:
- Hit ratio 30% untuk operation 100ms masih valuable
- Hit ratio 99% untuk operation 0.1ms belum tentu worth complexity

---

## 15. Memory / Cardinality / Eviction

### Key Explosion: Cardinality Trade-off

```
❌ Potensi masalah: cache:tenant:42:user:{user_id}:dashboard
```

Jika 10.000 user membaca dashboard agregat yang identik per-cabang, key berbasis `user_id`:
- Mengorbankan reuse yang tinggi
- Mempercepat eviction saat memory pressure

```
✅ Baik: cache:tenant:42:branch:{branch_id}:dashboard
```

**Prinsip Cardinality:**
- Semakin spesifik cache key = semakin kecil sharing/reuse = memory intensive
- High-cardinality tidak otomatis anti-pattern
- Key `products` tanpa tenant/filter bukan aman otomatis

`evicted_keys > 0` tidak otomatis berarti TTL terlalu pendek.

---

## 16. Laravel Cache::remember Caveats

```php
Cache::remember('key', $ttl, fn () => ...);
```

`Cache::remember()` tidak otomatis memberikan stampede protection. Beberapa concurrent request dapat menjalankan callback secara bersamaan pada saat miss.

Untuk protection, gunakan `Cache::lock(...)` dengan pattern:

```
[cache GET] → miss → [acquire lock] → double-check → [query DB] → [populate cache] → [release lock]
```

**Jika lock gagal:** Bounded wait/retry, re-check cache, fallback sesuai policy.

---

## 17. Permission & Security Caching

Permission dapat di-cache jika security model mengizinkan:

- **Key isolation** (`tenant:42:user:123:permissions`)
- **Active Invalidation** saat role berubah
- **Short TTL** (contoh: 5–10 detik) sebagai safety net
- **Documented Security SLA**

**Miskonsepsi:** "Permission di shared Redis tidak aman."

**Fakta:** Masalahnya biasanya:
- Missing tenant scope
- Stale authorization
- Revoke latency
- Invalidation failure
- Version mismatch

Untuk immediate revocation: authoritative lookup, centralized policy service, atau short-lived authorization credentials.

---

## 18. Query Optimization → Index → Cache

Cache **bukan pengganti** database optimization.

**Urutan Diagnosis:**
1. Ukur endpoint
2. Cek N+1 query
3. Cek Execution Plan
4. Tambahkan/index optimalkan
5. Kurangi column yang dibaca
6. Optimalkan Join/Subquery
7. **Baru Evaluasi Caching** jika workload membutuhkan

Contoh: Query < 1ms tidak selalu tidak layak di-cache. Untuk request volume tinggi (100k+ req/s), overhead Redis bisa menurunkan throughput DB.

---

## 19. Decision Framework / Rule of Thumb

### 5 Pertanyaan Sebelum Membuat Cache:

1. Berapa lama data boleh stale?
2. Seberapa mahal mendapatkan/menghitung data vs cache overhead?
3. At what cost is 5-30 detik delay acceptable?
4. Bagaimana cache di-invalidasi?
5. Apakah benefit performanya lebih besar daripada complexity?

### Senior Level Ask:
- Apa source of truth?
- Behavior saat Redis unavailable?
- Cache key aman untuk multi-tenant?
- Observability (hit/miss) dirancang?
- Risiko stampede?
- Memory/cardinality impact?

---

## 20. CMMS Examples

| Data | Cache? | Alasan & Acceptable Staleness | Invalidation Trigger |
|------|--------|-------------------------------|----------------------|
| Daftar Cabang | ✅ Ya | Read-heavy, rarely changes | Update master |
| Template Invoice | ✅ Ya | Static document | Version change |
| Dashboard Revenue | ✅ Ya | Expensive aggregation | New payment |
| Permission | ⚠️ Conditional | Read-heavy, security-bound | Revoke/TTL |
| Saldo Kas | ❌ Tidak | Keuangan, butuh ACID | N/A |
| Ketersediaan Stock | ⚠️ UI Only | Optimasi tampilan | Purchase part |

---

## Exercises

### Design Exercise — Dashboard Bengkel

TTL contoh: 10-30s untuk operational data (tergantung workload), bukan forever.

---

## Common Mistakes

1. Cache Everything — Tidak mengukur cost vs benefit
2. Missing Invalidation Strategy — Remember forever anti-pattern
3. Missing Isolation — Global key bocor tenant
4. Cache as Single Point of Failure — Tidak desain fallback DB
5. Key Explosion — Cache berdasarkan per-user ID bila data identik
6. Assume Cache = Strong Consistency — Write-through tidak selesai concurrency

---

## Senior Engineer Takeaways

Cache yang baik mengurangi pekerjaan berulang tanpa bergantung pada data yang mungkin stale.

**Mindset Junior:** "Query lambat, cache saja."

**Mindset Senior:**
1. Ukur bottleneck
2. Optimalkan query/index
3. Pahami business consistency
4. Tentukan acceptable staleness
5. Desain fallback/invalidation
6. Hitung operational complexity
7. Tambahkan cache jika benefit throughput signifikan

---

## Running the Lab

### Prerequisites
- Go 1.22+
- Docker Compose

```bash
docker compose up -d redis
cd labs/04-caching
go test -v -count=1 ./...

# Run Experiments
go run ./cmd/demo -scenario=without-cache
go run ./cmd/demo -scenario=cache-aside
go run ./cmd/demo -scenario=stampede-unprotected
go run ./cmd/demo -scenario=stampede-protected
go run ./cmd/demo -scenario=write-through
```

---

## Navigation

- **Previous:** [Lab 03 — Distributed Transaction](../03-distributed-transaction/)
- **Next:** [Lab 05 — Race Condition](../05-race-condition/)
