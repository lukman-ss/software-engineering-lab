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

## 4. Cache Aside vs Read Through

Lab ini menggunakan pattern **Cache Aside**. Penting untuk tidak menyamakannya dengan Read Through.

### Cache Aside (Pola yang digunakan di lab ini)

Aplikasi sendiri yang mengontrol orkestrasi antara cache dan database:

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

### Read-Your-Writes Problem
Cache-aside sederhana tidak menjamin read-your-writes. User yang baru saja melakukan update mungkin melihat data lama.
- Solusi: TTL sangat pendek, invalidate-on-write (dengan safety net TTL), atau write-through (dengan safety net TTL).

---

## 6. Cache Invalidation

Invalidation adalah mekanisme untuk membersihkan atau mengganti stale data.

### Race Condition pada Invalidation (DELETE)

Timeline yang benar untuk memahami risiko stale cache:

**UNSAFE: DELETE cache sebelum DB COMMIT**

```
Timeline:
T1 Writer: DELETE cache
T2 Reader: cache MISS → baca DB (lama, belum committed)
T3 Reader: SET nilai lama ke cache
T4 Writer: DB COMMIT nilai baru

Hasil:
Database = nilai baru
Cache = nilai lama (stale)

ISNIH: Ini adalah STALE CACHE BUKAN DATA LOSS. Business data di DB tetap aman.
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

Hasil:
Database = baru
Cache = lama (stale)
```

**Kesimpulan:**

Cache-aside adalah **eventually-consistent optimization**, bukan strong consistency primitive.

- **TTL adalah safety net** untuk menghapus stale entry yang terjebak akibat race condition.
- **Correctness-critical operation tidak boleh bergantung pada cache freshness.**
- Validasi akhir harus selalu ke database jika keakuratan data adalah prioritas utama.

### Strategi Invalidation

1. **TTL Expiration** — Data dibiarkan stale hingga waktu habis
2. **Explicit Delete** — DELETE key pasca-commit
3. **Event-Driven Invalidation** — Publish event pasca-commit untuk memicu invalidasi (pub/sub)

### 4. Tag / Index-based Invalidation

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
- Waktu eksekusi tidak terprediksi

#### Tag / Reverse Index Invalidation

Aplikasi/framework maintain mapping manual:

```
tag:product:123
      ↓
product:123
products:list
products:popular
```

Saat product berubah:
- Baca members dari tag index
- Invalidate semua key yang terhubung

**Trade-off:**
- Metadata tambahan yang harus dipelihara
- Synchronization overhead antara data dan index
- Stale reverse index mungkin terjadi
- Konsistensi index tidak otomatis

#### Alternatif Lain

- **Versioned namespace** — `product:123:v1` → `product:123:v2` (manual switch)
- **Explicit dependency invalidation** — Invalidate key secara eksplisit sesuai dependency
- **Event-driven invalidation** — Pub/sub untuk trigger invalidasi real-time

### Cache Key Versioning

Versioning ada dua tujuan yang berbeda:

#### 1. Schema/Deployment Versioning

```
product:v1:123  → format cache lama
product:v2:123  → format cache baru (migration)
```

**Digunakan ketika:**
- Serialization format berubah
- Schema cache berubah
- Deployment baru membutuhkan namespace terpisah

**Karakteristik:**
- Manual deployment-level change
- Tidak ada runtime dependency pada data change
- Bisa flush v1 sekaligus tanpa menyentuh data produk

#### 2. Data Revision/Version (Data Versioning)

```
product:123:revision:42
```

**Digunakan ketika:**
- Revision authoritative berubah bersama state bisnis
- Butuh cara mengetahui revision terbaru

**Penting:** Data-revision cache membutuhkan cara mengetahui revision terbaru. Ini berarti:
- Ada pointer/version metadata yang harus di-coordinasi
- Jika metadata tidak tersinkronisasi, masih ada race window
- Version bump TIDAK selalu atomic invalidation

**Contoh implementasi:**
```
product:123:revision:{current_revision}
```
Dengan `current_revision` ditarik dari:
- Database trigger + pub/sub
- Event stream offset
- Monotonically increasing version table

**Bukan solusi pemanis:**
- Version bump tidak menggantikan kebutuhan invalidation yang kompleks
- Masih perlu TTL sebagai safety net untuk edge case

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

**Penting: Write-Through tidak memberikan atomicitas atau strong consistency secara otomatis.**

Database PostgreSQL dan Redis adalah dua sistem berbeda. Tanpa distributed transaction/coordinated commit, urutan:

```
DB COMMIT
↓
Redis SET
```

**TIDAK atomic.**

#### Failure Window yang Mungkin Terjadi:

**Scenario A: Process Crash (Write Skew)**

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

Hasil:
Database = B
Cache = A (stale)
```

#### Kesimpulan yang Benar:

- **Tidak Atomic**:Tidak ada guarantee urutan eksekusi DB→Cache.
- **Tidak Strong Consistency**: Stale window tetap ada, bahkan setelah DB commit.
- **Read-Your-Writes Tidak Garantir**: Concurrent readers mungkin melihat data stale.
- **TTL Wajib sebagai Safety Net**: Untuk membersihkan stale data yang terjebak.

**Write-Through dapat membantu:**
- Mengurangi stale window dibanding Cache-Aside, namun tidak menghilangkan kebutuhan TTL/validation
- Memungkinkan read-after-write pada happy path (bukan guarantee otomatis)
- Mengurangi, tidak menghilangkan, kebutuhan invalidation manual

**Tapi perlu:**
- TTL sebagai fallback untuk semua failure mode
- Version ordering atau mutex jika concurrent writes relevan
- Desain error handling yang matang untuk partial failure

---

## Tabel Perbandingan: Cache-Aside vs Write-Through

| Aspek | Cache-Aside | Write-Through |
|-------|-------------|---------------|
| **Read Path** | App → Check cache → Miss → DB → SET cache | App → Check cache → HIT → Return |
| **Write Path** | App → DB write → Explicit cache invalidation | App → DB write → Update cache (best-effort) |
| **Stale Window** | Perlu (hingga invalidation TTL) | Lebih pendek, tapi tetap ada |
| **Partial Failure Risk** | Low (DB write saja, cache optional) | Tinggi (cache tidak pernah diupdate mungkin) |
| **Read-Your-Writes** | Tidak garantir tanpa invalidation | Garantir pada happy path, tidak otomatis |
| **Write Amplification** | Tidak ada (lazy population) | Ya (write to both DB and cache) |
| **Complexitas** | Sedang (invalidation logic diperlukan) | Sedang (failure handling diperlukan) |
| **TTL Safety Net** | Wajib | Wajib |
| **Use Case** | Read-heavy, write-rare, toleran stale | Read-heavy, write-moderate, butuh read-after-write |

**Penting:** Tidak ada pola yang "universally better". Pilih berdasarkan:
- Business consistency requirements
- Write frequency
- Failure tolerance
- Operational complexity yang dapat diterima

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

Redis dan database keduanya dapat digunakan untuk menyimpan data, tetapi **tujuan dan semantiknya berbeda total**.

### CACHE

**Tujuan:**
- Optimization / reuse
- Menghindari computation / I/O berulang

**Scope (bisa multiple):**
- **Global** — Data yang sama untuk semua user (contoh: daftar cabang)
- **Tenant scoped** — Isolasi antar konsumen (`tenant:42:*`)
- **Branch scoped** — Berdasarkan resource bisnis (`branch:7:*`)
- **User scoped** — Data khusus per-user (`user:123:*`)
- **Query scoped** — Hasil query tertentu (`query:revenue:2026-09:01`)

**Contoh cache per-user yang VALID:**
```
cmms:tenant:42:user:123:permissions:revision:7
```

**Syarat aman:**
- Key isolation benar (tenant + user scope)
- Invalidation benar (triggers saat role berubah)
- Security requirement dipenuhi (access control di application layer)

---

### SESSION

**Tujuan:**
- Menyimpan state suatu user/session
- Lifecycle mengikuti login/session

**Contoh:**
- Login state (auth token, user_id, expiry)
- CSRF token
- Shopping cart server-side
- Temporary preference (language, theme)

**Implementasi session store dapat menggunakan:**
- Redis (populer untuk scalability)
- Database
- Memory (in-process, tidak scalable)
- Cookie (client-side, terbatas ukuran)
- Datastore lain (memcached, DynamoDB, dll)

**Penting:** Session store tidak harus memiliki database sebagai source of truth. Implementasi tertentu mungkin menggunakan database untuk persistensi session, tapi itu adalah detail implementasi, bukan definisi semantik session.

---

### Perbedaan Utama

| Aspek | Cache | Session |
|-------|-------|---------|
| **Tujuan Utama** | Performance optimization | State persistence |
| **Allowed Data Types** | Query results, computed data, aggregations | Auth state, user preferences, shopping cart |
| **Acceptable Staleness** | Bisa stal (detik-menit) | Harus konsisten |
| **Failure Mode** | Degradasi performa | Login hilang / cart hilang |
| **TTL/Expiration** | Singkat (detik-menit) | Panjang (jam-hari) |

---

### Miskonsepsi Permission Cache di Redis

**Miskonsepsi:** "Permission per-user di Redis shared tidak aman."

**Fakta:** Masalahnya bukan Redis shared, tetapi biasanya:
- Missing tenant scope
- Missing user scope
- Stale authorization
- Revoke latency
- Invalidation failure
- Incorrect version/revision

**Solusi yang benar:**
- Key isolation (`tenant:42:user:123:permissions`)
- Short TTL (5-10 detik) sebagai safety net
- Active invalidation saat role berubah
- Documented Security SLA

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
- **Observability Pembedaan:** Pastikan `cache_miss` dan `cache_error` (backend error) dihitung di metric terpisah. Hit ratio tinggi saja tidak berarti cache berguna; bandingkan dengan query reduction di DB.

**Catatan tentang Hit Ratio:**
Hit ratio tidak universal. Contoh:
- Hit ratio 30% untuk operation yang sangat mahal (100ms) masih sangat valuable
- Hit ratio 99% untuk operation 0.1ms belum tentu worth complexity-nya

---

## 15. Memory / Cardinality / Eviction

**Key Explosion: Cardinality Trade-off**

```
❌ Potensi masalah: cache:tenant:42:user:{user_id}:dashboard
```

Jika 10.000 user membaca dashboard agregat yang identik per-cabang, key berbasis `user_id`:
- Mengorbankan reuse yang tinggi
- Membuat memory footprint menjadi mahal
- Mempercepat eviction saat memory pressure

Namun, key spesifik bukan selalu anti-pattern. Pertimbangkan:
- Apakah data memang berbeda per-user?
- Apakah kebutuhan isolation business?
- Apa benefit reuse yang hilang?

```
✅ Baik: cache:tenant:42:branch:{branch_id}:dashboard
```
Satu branch key bisa di-reuse 10.000 user.

**Prinsip Cardinality:**

Semakin spesifik cache key = semakin kecil sharing/reuse = semakin memory intensive *jika tidak diperlukan*. Namun:
- High-cardinality cache tidak otomatis anti-pattern
- Key `products` tanpa tenant/filter bukan aman otomatis — bisa menyebabkan correctness bug
- Personalized cache bisa diperlukan untuk isolasi data yang benar

**Trade-off yang perlu dipertimbangkan:**
> value gained vs memory footprint vs reuse vs operational complexity

**Low Hit Ratio:**
Low hit ratio bukan berarti key terlalu spesifik saja. Bisa disebabkan oleh:
- Workload cold (jarang diakses)
- TTL terlalu pendek
- Poor reuse design
- Cardinality tinggi
- Aggressive invalidation
- Workload memang tidak cache-friendly

**TTL & Eviction:**
- `evicted_keys > 0` tidak otomatis berarti TTL terlalu pendek
- Eviction biasanya terjadi ketika Redis menghadapi memory pressure/maxmemory sesuai eviction policy
- TTL terlalu pendek lebih mungkin menghasilkan: expiration meningkat, hit ratio turun, cache miss naik

**Memory Utilization:**
`used_memory_pct > 80%` bukan universal "critical" threshold. Nilai kritis tergantung:
- maxmemory
- fragmentation
- replication
- persistence
- workload
- eviction policy
- operational headroom

---

## 16. Laravel Cache::remember Caveats

```php
Cache::remember('key', $ttl, fn () => ...);
```

**Penting:** Helper ini praktis, namun memiliki batasan yang harus disadari oleh Senior Engineer.

### Closure Concurrent Execution (Stampede)

`Cache::remember()` tidak otomatis memberikan stampede protection. Beberapa concurrent request dapat menjalankan callback secara bersamaan pada saat miss. Protection harus dirancang secara eksplisit.

### Fallback Failure

Saat cache backend (Redis) down:
- Cache GET dapat melempar exception
- Callback mungkin tidak pernah dijalankan
- Framework error bisa menggagalkan eksekusi

Graceful degradation ke database harus didesain secara eksplisit dengan try-catch atau circuit breaker.

### Pattern Locking Eksplisit (Cache::lock)

Untuk perlindungan stampede, gunakan `Cache::lock(...)` secara manual:

```
[cache GET]
   ↓
miss
   ↓
[acquire lock]
   ↓
double-check cache
   ↓
[query DB]
   ↓
[populate cache]
   ↓
[release lock]
```

**Jika lock gagal:**
- Bounded wait/retry
- Re-check cache (mungkin sudah ter-populate)
- Fallback sesuai policy (DB query tanpa populate, atau return error)

### Catatan Tambahan

- Driver-dependent: behavior locking tidak konsisten di semua driver
- Key mapping: key string harus merepresentasikan seluruh input variable
- Framework modern Laravel: mendukung stale-while-revalidate / flexible caching sebagai alternatif

### Tentang Library Predis

Predis adalah client Redis untuk PHP. Library ini tidak otomatis memberikan stampede protection. Protection datang dari algorithm/locking yang dibangun di atasnya, bukan dari client library.

---

## 17. Permission Caching

Apakah Role/Permission boleh di-cache? Ya, tergantung requirement *revoke latency* dan threat model.

**Model yang Benar:**
Permission dapat di-cache dengan:
- Key terisolasi (scoped ke `tenant` dan `user`).
- Active Invalidation (delete cache saat Admin mengubah role).
- Short TTL (contoh: 5–10 detik) sebagai safety net revoke latency.
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

Contoh: Query `SELECT * FROM branch` (12 baris, index PK) < 1ms tidak selalu layak di-cache. Untuk request volume tinggi (ratusan ribu per detik), overhead Redis cache hit bisa menurunkan throughput DB secara signifikan. Namun:
- Query 2ms × 100,000 req/s dapat sangat layak di-cache
- Perubahan desain (network latency, DB scaling) bisa mengubah trade-off ini

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
- **TTL:** Contoh: 10-30s untuk operational data (tergantung workload), bukan forever.

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