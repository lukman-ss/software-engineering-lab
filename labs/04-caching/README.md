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
17. Permission Caching
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

Cache-Aside adalah pola di mana **application** (bukan caching layer abstraction) mengontrol cache interaction:

```
[Client]
    │
    ▼
[App: GET Cache]  → HIT → [Return Cached]
    │ MISS
    ▼
[App: GET Database] → Compute → [App: SET Cache w/ TTL] → [Return]
```

**Karakteristik:**
- Cache miss diselesaikan oleh aplikasi, bukan oleh cache provider.
- BUKAN "Read-Through" (di mana cache service secara otomatis load ke database pada saat miss).

---

## 5. TTL & Staleness

### Apakah Stale Data Masalah?
Pertanyaan kuncinya bukan "apakah data berubah?" tetapi **"berapa lama stale data dapat diterima?"**

| Data | Max Staleness | Reasonable? |
|------|--------------|-------------|
| Dashboard statistik | 30s–2min | ✅ Ya, operational metrics |
| Stock display | 1–5s | ✅ Ya, untuk UI/UX saja |
| Saldo wallet | 0s | ❌ Tidak, risiko audit |

### Read-Your-Writes Problem
Cache-aside sederhana tidak menjamin read-your-writes. User yang baru saja melakukan update mungkin melihat data lama.
- Solusi: TTL sangat pendek, invalidate-on-write, atau write-through.

---

## 6. Cache Invalidation

Invalidation adalah mekanisme untuk membersihkan atau mengganti stale data.

### Race Condition pada Invalidation (DELETE)

**Unsafe: DELETE cache sebelum DB COMMIT**
```
T1 Writer DELETE cache
T2 Reader cache MISS
T3 Reader membaca DB (lama)
T4 Reader SET nilai lama ke cache
T5 Writer COMMIT nilai baru

Hasil: DB = baru, Cache = lama (stale)
```

**Lebih aman: COMMIT → DELETE**
```
T1 Writer COMMIT nilai baru
T2 Writer DELETE cache
T3 Reader cache MISS → baca DB baru → SET cache baru

Hasil: DB = baru, Cache = baru
```

**TAPI: COMMIT → DELETE masih tidak perfect**
```
T1 Reader cache MISS
T2 Reader membaca DB (lama)
T3 Writer COMMIT nilai baru
T4 Writer DELETE cache
T5 Reader SET hasil lama ke cache

Hasil akhir: Cache dapat tetap stale.
```

**Kesimpulan:** Cache-aside adalah **eventually consistent optimization**, bukan strong consistency primitive. **TTL menjadi safety net** untuk menghapus stale entry yang terjebak akibat race condition. Untuk correctness-critical data, jangan menggantungkan correctness pada cache.

### Strategi Invalidation

1. **TTL Expiration** — Data dibiarkan stale hingga waktu habis
2. **Explicit Delete** — DELETE key pasca-commit
3. **Event-Driven Invalidation** — Publish event pasca-commit untuk memicu invalidasi (pub/sub)
4. **Tag / Index-based Invalidation** 
   - **Pattern deletion (SCAN)**: Mencari prefix `product:*` lalu DELETE. (SCAN bukan `KEYS` dan production-friendly, tapi non-atomic dan memiliki operational cost).
   - **Tag index mapping**: Aplikasi menyimpan metadata key mana yang terhubung dengan tag `product_list`. (Butuh state maintenance yang kompleks).
   - **Note:** Redis core tidak menyediakan cache-tag abstraction secara native.

### Cache Key Versioning

Versioning berguna untuk dua use case berbeda:

**1. Schema/Deployment Versioning**
```
product:v1:123  → format cache lama
product:v2:123  → format cache baru (migration)
```
Cocok untuk migrasi payload cache; tidak membutuhkan invalidation.

**2. Data Revision/Version (Data Versioning)**
```
product:123:revision:42
```
Bisa digunakan untuk invalidate entries asalkan *revision authoritative* dapat ditarik dengan cepat tanpa database bottleneck yang sama. Sekadar mengganti hard-coded `v1` menjadi `v2` dalam code base BUKAN solusi invalidation untuk perubahan harga produk.

### RememberForever Anti-Pattern

Menyimpan cache tanpa batas waktu:
```php
Cache::rememberForever('product_123', fn() => ...); // atau TTL = 0
```
**Mengapa ini anti-pattern?**
Jika developer lupa mendesain event-driven invalidation, data stale tidak akan pernah hilang. TTL memberikan safety net. Penggunaan ini hanya aman jika:
- Data hampir immutable (mis: kode area).
- Perubahan source data dijamin 100% selalu memicu invalidate cache.

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

Write-Through di sini berarti aplikasi meng-update DB dan Cache pada response path yang sama.

### Karakteristik & Risiko Partial Failure

- **TIDAK Atomic**: Database dan Redis adalah sistem terpisah. Tidak ada atomic commit lintas-sistem.
- **Bukan Strong Consistency**: Writer race conditions masih bisa terjadi jika tidak dilindungi version order.
- **Failure Mode Window**:
  - Jika proses crash setelah DB COMMIT sukses tetapi sebelum Redis SET, data cache menjadi stale.
  - Jika Cache SET gagal (Redis timeout), operasi bisnis (DB) harus tetap sukses. Kita dapat mencoba fallback DELETE cache, namun jika DELETE juga gagal, stale cache tetap bertahan hingga TTL.

Oleh karena itu, Write-Through hanya *memperpendek* stale window dan memberikan jaminan *read-your-writes* pada best-case scenario, namun **TTL tetap wajib digunakan sebagai safety net**.

---

## 8. Cache Stampede

### Problem Flow

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

Stampede terjadi ketika key populer expire dan traffic besar masuk di saat yang sama.

---

## 9. Singleflight

Mengatasi stampede **dalam satu instance/process**.

- Menggunakan `golang.org/x/sync/singleflight`.
- Request concurrent ditahan, satu goroutine jalan ke DB, hasilnya dibagikan ke caller lain.
- Gunakan `DoChan` agar mendukung context cancellation dan tidak deadlock.
- **Batasan Penting:** Singleflight lokal tidak mencegah stampede lintas instance (multi-deployment).

---

## 10. Distributed Lock

Mengatasi stampede **multi-instance**.

Menggunakan `SETNX` (Set if Not Exists) untuk menahan instance lain saat satu instance rebuild cache.

**Persyaratan keamanan lock:**
- **Harus memiliki expiration (TTL):** Mencegah deadlock permanen jika lock holder mati (OOM/Crash).
- **Menggunakan unique owner/token:** Menyimpan UUID saat memegang lock.
- **Tidak sembarang delete:** Lock hanya dilepas jika token di Redis cocok dengan UUID holder (melalui Lua script / CompareAndDel).

---

## 11. TTL Jitter & Background Refresh

### TTL Jitter
Tambahkan pengacak kecil ke TTL (misal: 60s ± 15s) agar batch cache keys yang di-set bersamaan tidak kadaluarsa tepat di detik yang bersamaan.

### Background Refresh (Stale-While-Revalidate)
Mem-refresh data secara background sebelum key benar-benar expired. Client tetap dilayani dengan data cache, tetapi menghindari delay synchronous hit.

---

## 12. Cache Key Design / Multi-Tenant

Format canonical untuk aplikasi B2B/CMMS:

```
{app}:{tenant}:{branch}:{resource}:{dimension}
```

Contoh: `cmms:tenant:42:branch:7:dashboard:2026-09-01`

**Dimensi yang Wajib Masuk Key:**
- `tenant:42` — Isolasi data antar konsumen (hindari cross-tenant data leak).
- `branch:7` — Resource scoping.
- `date:2026-09-01` — Scope business date (harus timezone-aware, bukan UTC).

---

## 13. Cache vs Session

Redis adalah storage engine yang bisa menyimpan keduanya, tetapi semantics aplikasinya sangat berbeda.

| Aspect | Cache | Session |
|--------|-------|---------|
| **Purpose** | Optimization (read reuse) | Conversational state (cart, auth state) |
| **Ownership** | Bisa shared/global atau per-user | Per-user session |
| **Source of Truth** | Derived (biasanya bisa di-rebuild dari DB) | Authoritative (contoh: cart sblm checkout) |
| **Flush Impact** | Latency naik (DB di-hit) | User ter-logout / Cart hilang |
| **TTL** | Volatilitas data (detik/menit) | Lifecycle login (jam/hari) |

**Miskonsepsi:** "Data per-user harus masuk session".
**Koreksi:** Data per-user BISA dan SAH di-cache di layer caching asalkan **key isolation** dilakukan dengan benar (e.g., `tenant:42:user:123:dashboard`). Session storage digunakan untuk *state lifecycle*, caching digunakan untuk *query optimization*.

---

## 14. Cache Failure & Redis Unavailable

Prinsip: **Cache adalah optimization layer, bukan dependency correctness utama.**

**Flow Read Ideal:**
```
Redis Error (Timeout / Connection Refused)
      ↓
Fallback ke Database / Source of Truth
      ↓
Request tetap berjalan jika DB capacity memungkinkan
```

- **Degraded Performance:** Aplikasi harus tetap berfungsi meskipun lambat.
- **Circuit Breaker:** Mencegah aplikasi stuck menunggu Redis Timeout terus-menerus.
- **Observability Pembedaan:** Pastikan `cache_miss` dan `cache_error` (backend error) dihitung di metric terpisah. Hit ratio tinggi saja tidak berarti cache berguna, bandingkan dengan query reduction di DB.

---

## 15. Memory / Cardinality / Eviction

**Anti-Pattern:** Key Explosion

```
❌ Salah: cache:tenant:42:user:{user_id}:dashboard
```
Jika 10.000 user dari tenant yang sama membaca dashboard yang agregat datanya identik per-cabang, membuat key berbasis `user_id` menyebabkan:
- Duplikasi memory ekstrim
- Redis Eviction Rate meroket
- Cache reuse sangat rendah

```
✅ Benar: cache:tenant:42:branch:{branch_id}:dashboard
```
Satu branch key di-reuse 10.000 user. Memory footprint rendah, eviction aman.

**Prinsip Cardinality:** Semakin spesifik cache key (misal: scope user), semakin kecil sharing/reuse-nya, semakin memory intensive caching-nya. Cache only what provides meaningful reuse.

---

## 16. Laravel Cache::remember Caveats

```php
Cache::remember('key', $ttl, fn () => ...);
```

Helper ini praktis, namun memiliki batasan yang harus disadari oleh Senior Engineer:

1. **Closure dapat dieksekusi bersamaan oleh beberapa request.**
   `Cache::remember` BUKAN jaminan perlindungan terhadap cache stampede.
2. **Behavior locking bergantung driver.**
   Tidak semua backend secara native melindungi concurrency. Double-check pattern / explicit locking (`Cache::lock`) diperlukan untuk query yang rentan overload.
3. **Fallback failure:**
   Jika Redis (atau cache provider) throw connection exception, eksekusi closure bisa digagalkan oleh framework error. Harus ada try-catch atau circuit breaker fallback khusus.
4. **Key mapping:**
   Key string harus merepresentasikan seluruh input variable yang masuk ke dalam closure.

---

## 17. Permission Caching

Apakah Role/Permission boleh di-cache? Ya, tergantung requirement *revoke latency* dan threat model.

**Model yang Benar:**
Permission dapat di-cache dengan:
- Key terisolasi (scoped ke `tenant` dan `user`).
- Active Invalidation (delete cache saat Admin mengubah role).
- Short TTL (5–10 detik) sebagai safety net revoke latency.
- Documented Security SLA.

```
Key: cmms:tenant:42:user:123:permission:revision:17
```

**Miskonsepsi:** "Permission di shared Redis selalu tidak aman."
**Koreksi:** Yang tidak aman adalah missing tenant scope, stale authorization tanpa invalidation fallback, dan version mismatch. Jika sistem butuh *immediate revocation absolut*, gunakan authoritative lookup ke DB di setiap request (jangan cache untuk operasi transaksional kritikal).

---

## 18. Query Optimization → Index → Cache

Cache **bukan pengganti database optimization**.

**Urutan Diagnosis:**
1. Ukur endpoint (Profile query).
2. Cek N+1 query.
3. Cek Execution Plan (EXPLAIN ANALYZE).
4. Tambahkan atau optimalkan index.
5. Kurangi column yang dibaca (jangan `SELECT *`).
6. Optimalkan Join/Subquery.
7. **Baru Evaluasi Caching** jika workload agregasi masih memakan I/O mahal dan repetitif.

Contoh: Query `SELECT * FROM branch` (12 baris, index PK) < 1ms tidak layak di-cache, karena overhead hit ke Redis (network round-trip) mungkin lebih mahal dari query native.

---

## 19. Decision Framework / Rule of Thumb

### 5 Pertanyaan Rule of Thumb Sebelum Membuat Cache:

1. Berapa lama data ini boleh stale?
2. Seberapa mahal mendapatkan/menghitung data ini vs cache overhead?
3. Apa dampaknya jika user melihat data 5-30 detik terlambat?
4. Bagaimana cache di-invalidasi saat source berubah?
5. Apakah manfaat performanya lebih besar daripada complexity tambahan?

**Aturan Emas:** Jika Anda tidak bisa mendefinisikan strategi invalidation, JANGAN buru-buru menambahkan cache.

### Pertanyaan Senior-Level Tambahan:
- Apa source of truth data ini?
- Apa behavior saat Redis unavailable?
- Apakah cache key aman untuk multi-tenant?
- Bagaimana observability (hit/miss) dirancang?
- Apa risiko stampede?
- Apa memory footprint/cardinality impact-nya?

---

## 20. CMMS Examples

Tabel Keputusan Caching untuk Domain Bengkel:

| Data | Cache? | Alasan & Acceptable Staleness | Invalidation Trigger |
|------|--------|-------------------------------|----------------------|
| Daftar Cabang | ✅ Ya | Read-heavy, rarely changes | Update master cabang |
| Template Invoice | ✅ Ya | Static document format | Versi template berubah |
| Dashboard Revenue | ✅ Ya | Expensive aggregation | New payment |
| Permission | ⚠️ Conditional | Security, but read-heavy | Revoke role, timeout TTL |
| Saldo Kas | ❌ Tidak | Keuangan, butuh ACID | N/A (Always DB) |
| Ketersediaan Stock | ⚠️ UI Only | Optimasi tampilan | Purchase part (Validasi transaksional tetap ke DB) |

---

## Exercises

### Design Exercise — Dashboard Bengkel

Dashboard menampilkan: Jumlah Invoice, Pendapatan, Top Mekanik, Kendaraan Baru.

**Soal:** Untuk setiap data di atas, jelaskan:
1. Perlu cache atau tidak?
2. Mengapa?
3. TTL berapa dan asumsinya?
4. Format Cache Key?
5. Event penyebab invalidasi?
6. Acceptable staleness?
7. Behavior jika Redis down?
8. Mencegah cache stampede?
9. Mencegah tenant cross-reading?
10. Metric monitoring?

**Expected Reasoning:**
- **Key:** Wajib `cmms:tenant:{id}:branch:{id}:metric`
- **Stampede:** Gunakan In-Process Singleflight + Jitter
- **Redis Down:** Fallback to DB (degraded performance)
- **TTL:** 10-30s untuk operational data, bukan forever.

---

## Common Mistakes

1. **Cache Everything** — Tidak mengukur cost vs benefit.
2. **Missing Invalidation Strategy** — Remember forever anti-pattern.
3. **Missing Isolation** — Global key bocor ke tenant lain.
4. **Cache as Single Point of Failure** — Tidak mendesain database fallback saat Redis network timeout.
5. **Key Explosion** — Menyimpan cache agregat berdasarkan per-user ID, padahal datanya identik per-cabang.
6. **Assume Cache = Strong Consistency** — Menganggap write-through menyelesaikan concurrency.

---

## Senior Engineer Takeaways

Cache yang baik mengurangi pekerjaan berulang tanpa menjadikan correctness sistem bergantung pada data yang mungkin stale.

**Mindset Junior:**
> "Query lambat, cache saja."

**Mindset Senior:**
> 1. Ukur bottleneck (latency network vs DB cpu)
> 2. Optimalkan query di level storage / index
> 3. Pahami business consistency requirement
> 4. Tentukan acceptable staleness
> 5. Desain fallback & invalidation flow
> 6. Hitung operational complexity
> 7. Baru tambahkan cache jika benefit throughput sangat signifikan.

---

## Running the Lab

### Prerequisites
- Go 1.22+
- Docker Compose (untuk Redis)
```bash
docker compose up -d redis
```

### Run Tests
```bash
cd labs/04-caching
go test -v -count=1 ./...
```

### Run Experiments
```bash
go run ./cmd/demo -scenario=without-cache
go run ./cmd/demo -scenario=cache-aside
go run ./cmd/demo -scenario=stampede-unprotected
go run ./cmd/demo -scenario=stampede-protected
go run ./cmd/demo -scenario=write-through
```