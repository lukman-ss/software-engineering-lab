# Lab 04 — Caching: Mengurangi Latency, Tapi Apa Biayanya?

> **Mental Model**: Caching adalah trade-off antara **latency** (cepat) dan **consistency** (benar). Senior engineer memilih teknik caching yang tepat untuk workload-nya, bukan menerapkan caching di mana-mana.

---

## Learning Objectives

Setelah menyelesaikan lab ini, Anda akan memahami:

1. Problem / Why Cache Exists
2. Source of Truth
3. Kapan Cache Layak / Tidak Layak
4. Cache Aside
5. TTL & Acceptable Staleness
6. Cache Invalidation
7. Write-Through Strategy
8. Cache Stampede
9. Singleflight
10. Distributed Lock
11. TTL Jitter / Background Refresh
12. Cache Key Design
13. Multi-Tenant Isolation
14. Cache vs Session
15. Cache Failure / Graceful Degradation
16. Memory / Cardinality / Eviction
17. Laravel Cache::remember Caveats
18. Permission & Security Caching
19. Metrics / Observability
20. Query Optimization → Index → Cache
21. Decision Framework
22. CMMS Examples
23. Exercises
24. Common Mistakes
25. Senior Engineer Takeaways
26. Running the Lab
27. Navigation

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
| **Primary** | PostgreSQL | Durable, persistent, authoritative storage yang dapat digunakan untuk rebuild cache |
| **Cache** | Redis | Derived data, TTL-bound, on-demand rebuild |

**Prinsip Utama:** Dalam lab ini, PostgreSQL adalah source of truth untuk business data. Cache correctness harus dibangun dengan asumsi cache dapat kosong atau di-flush kapan saja.

**General principle (beyond this lab):** Cache umumnya dapat berasal dari:
- Database (authoritative source)
- External API
- Computation
- Filesystem
- Another service

Yang penting: cache harus dapat direbuild atau divalidasi terhadap authoritative source sesuai architecture.

---

## 3. Kapan Cache Layak / Tidak Layak?

### Apa yang Sebaiknya Di-cache?

Cache cocok untuk data yang:
1. **Mahal dihitung** — Query expensive (join, aggregation)
2. **Sering dibaca** — Read-heavy workload
3. **Jarang berubah dibanding frekuensi baca**
4. **Toleran terhadap stale data**

Contoh: Dashboard operasional, laporan agregasi, master supplier, katalog jenis service.

### Apa yang Sebaiknya Tidak Di-cache?

Data yang **butuh strong consistency** atau berubah **sangat cepat**:
- Saldo wallet / kas
- Status pembayaran

**Nuance:** Data real-time **bisa** menggunakan cache sebagai optimization layer asalkan:
- Database tetap source of truth
- Correctness (validasi akhir) tidak bergantung pada cache
- TTL sangat pendek

Contoh: Stock display cache 3s untuk UI, tetapi validasi final ke DB dengan `SELECT FOR UPDATE` saat transaksi checkout.

### OTP & Ephemeral State

OTP/token sekali pakai **bukan cache** tetapi sering cocok disimpan di Redis sebagai **ephemeral state store**:
- TTL native
- Fast lookup
- Atomic commands
- Ephemeral lifecycle (expire setelah use/expiry)

**Penting:** Redis ≠ cache. Redis adalah datastore/technology. "Cache" adalah semantic pattern.

| Semantic | Redis Usage |
|----------|-------------|
| Cache | Derived/reconstructable optimization |
| Session Store | User/session state |
| OTP Store | Security-sensitive ephemeral state |
| Lock Backend | Coordination state |
| Queue | Messaging state |

OTP tetap tidak boleh digunakan sebagai **stale reusable cache** atau data yang aman direplay untuk retry.

---

## 4. Cache Aside

Lab ini menggunakan pattern **Cache Aside** untuk read operations.

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

## 5. TTL & Acceptable Staleness

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
T1 Reader: cache MISS
T2 Reader: baca DB lama, tetapi belum menulis cache
T3 Writer: DB COMMIT nilai baru
T4 Writer: DELETE cache
T5 Reader: SET hasil DB lama ke cache

Hasil: Database = nilai baru, Cache = nilai lama (stale)
```

**Kesimpulan:** Urutan `COMMIT → DELETE` lebih aman daripada `DELETE → COMMIT` karena
menghapus cache *setelah* nilai baru committed, sehingga mengurangi window stale write.
Namun, cache-aside tetap **eventually-consistent optimization**, bukan strong consistency primitive.
TTL tetap menjadi safety net untuk window residual di mana reader sempat mengambil nilai lama
sebelum DELETE tereksekusi.

### Invalidation Strategies

1. **TTL Expiration** — Data dibiarkan stale hingga waktu habis
2. **Explicit Delete** — DELETE key pasca-commit
3. **Event-Driven Invalidation** — Publish event pasca-commit untuk memicu invalidasi (pub/sub)

**Reliability Note:** Redis Pub/Sub bersifat ephemeral. Jika subscriber disconnect, event dapat hilang. Untuk cache invalidation biasa yang memiliki TTL safety net, loss event mungkin acceptable tergantung requirement.

Jika invalidation harus reliable (lost invalidation event):
- **Durable queue** (RabbitMQ)
- **Redis Streams** dengan consumer groups
- **Kafka / NATS JetStream**
- **Transactional outbox + CDC consumer**

TTL tetap berguna sebagai recovery/safety net untuk event loss.

### Tag / Index-based Invalidation

**Penting:** Redis core tidak menyediakan generic cache tagging abstraction. Banyak penyedia menyamakan SCAN+pattern delete dengan tag implementation. Mereka adalah hal yang berbeda.

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

## 7. Write-Through Strategy

**Terminologi:** Dalam lab ini istilah "Write-Through" digunakan sebagai simplifikasi untuk _application-managed update-on-write_. Ini berbeda dari formal cache-managed write-through architecture di mana cache layer secara sinkron persist ke backing store.

Flow implementation kami:
```
DB update (RETURNING) → get authoritative value → best-effort cache update
```

Write-Through di sini berarti aplikasi meng-update DB dan Cache pada urutan yang sama setelah DB commit sukses.

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

### Perbandingan: Read Strategy vs Write Strategy

**Catatan penting:** Cache Aside dan Write-Through biasanya digunakan untuk menyelesaikan masalah yang berbeda:

- **Cache Aside** menggambarkan **read strategy**: cara mengambil data dari cache
- **Write-Through / Invalidate-on-Write** menggambarkan **write strategy**: cara memelihara konsistensi setelah write

#### Read Strategy (Cache Aside)

| Aspek | Description |
|-------|-------------|
| Read Path | App → Check cache → Miss → DB read → Populate cache |
| Miss Handling | Ditangani secara explicit oleh aplikasi |
| Error Mode | Graceful degradation ke DB |

#### Write Strategy Options

| Strategi | Flow | Trade-off |
|----------|------|-----------|
| **Invalidate** | DB commit → delete cache | Simple, stale cache mungkin bertahan sampai TTL |
| **Write-Through** | DB commit → update cache (best-effort) | Cache cenderung warm, partial failure diperlukan handle |

**Partial Failure pada Write Strategy:**
```
DB COMMIT berhasil
↓
cache DELETE/SET gagal
↓
stale cache masih hidup sampai TTL
```

TTL adalah safety net yang sangat berguna, tetapi bukan "wajib" sebagai hukum universal.

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

**Double-check pattern:**
```
[initial cache GET] → miss → singleflight.Do → [check cache lagi] → query DB → populate cache → share result
```

**Alasan:** Entry mungkin sudah dipopulate di antara initial miss dan saat caller menjadi leader.

---

## 10. Distributed Lock

Mengatasi stampede **multi-instance**.

**Primitive (implemented in source code):**
```
WithLock() = try-once lock primitive (no retry/backoff)
```

**Production cache-regeneration strategy (conceptual pattern):**
```
[cache GET] → miss → [acquire distributed lock] → [check cache lagi] → query DB → populate cache → [safe release]

NON-HOLDER: buffered wait/retry → re-check cache
```

**Nota:** Diagram di atas menunjukkan pattern konseptual untuk cache regeneration. Implementasi `WithLock()` hanya menyediakan try-once primitive. Untuk pola non-holder bounded wait/retry, dibutuhkan coordinator khusus yang belum diimplementasikan.

**Lock Requirements:**
- Menggunakan unique token/owner
- Memiliki TTL (mencegah deadlock permanen)
- Release menggunakan atomic compare-and-delete
- Tidak dapat menghapus lock holder lain

**Lock Lease Expiration Problem:**
Jika rebuild lebih lama dari lock TTL:
```
Lock expires
↓
Instance B memperoleh lock baru
↓
Instance A masih bekerja
↓
Duplicate rebuild dapat berjalan
```

Lock pada Lab 04 digunakan untuk mengurangi duplicate cache rebuild, bukan sebagai correctness primitive transaksi bisnis.

---

## 11. TTL Jitter & Background Refresh

### TTL Jitter
Tambahkan pengacak kecil ke TTL (misal: 60s + random 0–15s) agar batch cache keys tidak kadaluarsa hampir bersamaan.
Model implementasi menghasilkan TTL dalam rentang [base, base + maxJitter], tidak pernah lebih rendah dari base.

### Background Refresh (Stale-While-Revalidate)
Mem-refresh data secara background sebelum key expired. Client tetap dilayani dengan data cache.

---

## 12. Cache Key Design

Format canonical: `{app}:{tenant}:{branch}:{resource}:{dimension}`

Contoh: `cmms:tenant:42:branch:7:dashboard:2026-09-01`

**Prinsip utama:** Cache key harus mencakup semua input/dimensi yang mempengaruhi result.

**Dimensi pada contoh Dashboard:**
- `tenant:42` — Isolasi data antar konsumen
- `branch:7` — Resource scoping  
- `date:...` — Scope business date (timezone-aware)

**Contoh lain:**

| Entity | Key Dimensions |
|--------|----------------|
| Product global | product ID + schema version |
| Product tenant-specific | tenant + product ID + version |
| Permission | tenant + user + permission revision |
| Session | user ID + session ID |

---

## 13. Multi-Tenant Isolation

Key isolation merupakan fondasi keamanan multi-tenant. tanpa isolation yang benar, data tenant lain dapat bocor.

**Prinsip:** Setiap key yang mengandung data sensitif harus mengandung tenant scope.

---

## 14. Cache vs Session

### Perbedaan Fundamental

| Aspect | Cache | Session |
|--------|-------|---------|
| **Purpose** | Optimization / reuse, menghindari computation / I/O berulang | State persistence, merepresentasikan user/session lifecycle |
| **Ownership/Scope** | Bisa global, tenant-scoped, user-scoped, query-scoped | Per-user / per-session |
| **Data Semantics** | Biasanya derived/reconstructable dari source of truth | Bisa authoritative untuk state session, bukan selalu business data |
| **Lifecycle** | Data-dependen, expired berdasarkan data relevance | Login/session-bound, clear pada logout/expiry |
| **Failure Impact** | Degradasi performa (DB di-hantam) | Logout / cart hilang / state workflow reset |
| **Typical TTL** | Workload-dependent (detik sampai hari) | Policy-dependent (menit sampai jam/hari) |

### Semantik yang Penting

**CACHE:** Tujuan utama optimization/reuse. Data biasanya derived/reconstructable dari source of truth. Cache dapat global, tenant-scoped, user-scoped, atau query-scoped. Cache correctness harus dibangun dengan asumsi cache dapat kosong/kadaluarsa kapan saja.

**SESSION:** Merepresentasikan state lifecycle session/user. Session storage sendiri dapat menjadi authoritative untuk state session (seperti shopping cart), bukan selalu authoritative business data. Kehilangan session biasanya menyebabkan logout/reset workflow, bukan corruption business data.

**Perbedaan utama** adalah semantics (tujuan data) dan lifecycle, bukan sekadar stale-allowed vs tidak-stale-allowed. Session pun dapat memiliki stale data tergantung implementasinya.

### OTP & Ephemeral State

OTP/token sekali pakai **bukan cache** tetapi sering cocok disimpan di Redis sebagai **ephemeral state store**:
- TTL native
- Fast lookup
- Atomic commands
- Ephemeral lifecycle (expire setelah use/expiry)

Redis adalah **datastore/technology**. "Cache" adalah **semantic/pattern**.

Pemakaian Redis tidak otomatis berarti tersebut adalah cache.

---

## 15. Cache Failure / Graceful Degradation

Prinsip: **Cache adalah optimization layer, bukan dependency correctness utama.**

Jika cache backend gagal, aplikasi dapat fallback ke authoritative DB selama:
- DB tersedia
- capacity masih mencukupi
- timeout/backpressure policy mengizinkan

**Peringatan:** Cache outage bisa menyebabkan:

```
traffic yang sebelumnya diserap Redis
      ↓
langsung masuk DB
      ↓
load spike / cache-failure amplification
```

Graceful degradation harus mempertimbangkan:
- Redis timeout pendek
- DB capacity
- request timeout
- circuit breaker jika relevan

**Degraded Performance:** Aplikasi harus tetap berfungsi meski lambat.
- **Circuit Breaker:** Mencegah aplikasi stuck menunggu Redis Timeout.
- **Observability:** Log `cache_miss` dan `cache_error` terpisah.

**Catatan tentang Hit Ratio:** Hit ratio tinggi tidak otomatis berarti cache memberikan ROI tinggi. Hit ratio 30% untuk operation 100ms masih valuable. Hit ratio 99% untuk operation 0.1ms belum tentu worth complexity.

---

## 16. Memory / Cardinality / Eviction

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

Satu branch key bisa di-reuse 10.000 user.

**Prinsip Cardinality:**
- Semakin spesifik cache key = semakin kecil sharing/reuse = memory intensive
- High-cardinality tidak otomatis anti-pattern
- Key `products` tanpa tenant/filter bukan aman otomatis

`evicted_keys > 0` tidak otomatis berarti TTL terlalu pendek. Eviction tergantung pada memory pressure, maxmemory config, dan eviction policy.

---

## 17. Laravel Cache::remember Caveats

```php
Cache::remember('key', $ttl, fn () => ...);
```

`Cache::remember()` tidak otomatis memberikan stampede protection. Beberapa concurrent request dapat menjalankan callback secara bersamaan pada saat miss.

Untuk protection, gunakan `Cache::lock(...)` dengan double-check pattern:

```
[cache GET] → miss → [acquire lock] → double-check → [query DB] → [populate cache] → [release lock]
```

**Jika lock gagal:** Bounded wait/retry, re-check cache, fallback sesuai policy.

---

## 18. Permission & Security Caching

Permission dapat di-cache jika security model mengizinkan:

- **Key isolation** (`tenant:42:user:123:permissions`)
- **Active Invalidation** saat role berubah
- **Short TTL** (contoh: beberapa detik) sebagai safety net revoke latency
- **Documented Security SLA**

TTL permission harus ditentukan berdasarkan:
- **Security SLA**: seberapa cepat revoke harus effective?
- **Revoke latency requirement**: berapa lama window approve untuk privilege escalation?
- **Threat model**: risk dari stale authorization?
- **Invalidation reliability**: apakah invalidation reliable atau perlu TTL?

Contoh 5–10 detik hanyalah illustrative. Beberapa skenario memerlukan TTL lebih pendek (0-1 detik), yang lain dapat bertahan lebih lama dengan durable invalidation.

Untuk immediate revocation, pertimbangkan:
- Authoritative lookup
- Centralized policy service
- Short-lived authorization credentials

**Miskonsepsi:** "Permission di shared Redis tidak aman."

**Fakta:** Masalahnya biasanya:
- Missing tenant scope
- Stale authorization
- Revoke latency
- Invalidation failure
- Version mismatch

---

## 19. Metrics / Observability

Cache value harus dievaluasi bersama:
- Avoided query cost
- Cache latency
- Memory cost
- Invalidation complexity
- Failure amplification

**Hit ratio tinggi tidak otomatis berarti cache memberikan ROI tinggi.**

Production implementation biasanya memakai Prometheus, OpenTelemetry, atau Redis monitoring commands (INFO, LATENCY).

---

## 20. Query Optimization → Index → Cache

Cache **bukan pengganti** database optimization.

**Urutan Diagnosa:**
1. Ukur endpoint
2. Cek N+1 query
3. Cek Execution Plan
4. Tambahkan/index optimalkan
5. Kurangi column yang dibaca
6. Optimalkan Join/Subquery
7. Baru Evaluasi Caching jika workload membutuhkan

**Penting:** untuk query sangat murah (< 1ms) dan workload tertentu, overhead Redis (network hop + latency + cache complexity) mungkin tidak memberikan benefit end-to-end yang berarti. Namun Redis tidak secara otomatis "menurunkan throughput DB" - traffic yang memakai cache justru mengurangi beban DB.

---

## 21. Decision Framework

### 5 Pertanyaan Sebelum Membuat Cache:

1. Berapa lama data boleh stale?
2. Seberapa mahal mendapatkan/menghitung data vs cache overhead?
3. Seberapa besar dampak penundaan 5-30 detik?
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

## 22. CMMS Examples

| Data | Cache? | Alasan | Contoh TTL | Acceptable Staleness | Invalidation Trigger |
|------|--------|--------|------------|---------------------|----------------------|
| Daftar Cabang | ✅ Ya | Read-heavy master data, relatif jarang berubah | Beberapa menit | Menit | Update master |
| Daftar Mekanik | ✅ Ya | Read-heavy master data, biasanya stable | 5–15 menit (illustrative) | Detik-menit | Mechanic master berubah / active invalidation |
| Daftar Sparepart | ⚠️ Conditional | Read-heavy master data, display UI boleh cache | 1–5 detik (illustrative) | Detik | Sparepart master berubah |
| Jenis Service | ✅ Ya | Relatively stable master/reference data | Long-lived (illustrative) | Menit-pertama | Versi baru / deployment |
| Master Supplier | ✅ Ya | Infrequent update | 5–15 menit | Menit | Update supplier |
| Konfigurasi Pajak | ⚠️ Conditional | Butuh akurat untuk transaksi | 10–30 detik | Detik | Update konfigurasi |
| Template Invoice | ✅ Ya | Static document | Versi-tag | Versi | Versi template berubah |
| Permission | ⚠️ Conditional | Read-heavy, security-bound | 5–10 detik (illustrative) | Detik | Revoke role |
| Dashboard Revenue | ✅ Ya | Expensive aggregation | 10–30 detik | Detik | New payment |
| Top Mekanik | ✅ Ya | Aggregation heavy | 30 detik | Detik | Job completion |
| Top Sparepart | ⚠️ Conditional | Aggregation heavy, ranking berasal dari transaksi | 30 detik | Detik | Invoice/service part line posted, transaction completed |
| Status Service | ⚠️ Conditional | Caching boleh untuk UI display, tapi validation transaksi harus DB | 5 detik | Detik | Service update |
| Invoice Baru | ⚠️ Conditional | Bisa cache read projection untuk UI, tapi tidak untuk authoritative state | 1–5 detik | Detik | Invoice status berubah |
| Status Pembayaran | ⚠️ Conditional | Bisa cache read projection untuk UI, tapi tidak untuk authoritative state | 1–5 detik | Detik | Payment status berubah |
| Saldo Kas | ⚠️ Conditional | Bisa cache read projection untuk display, tapi tidak untuk transactional operation | 1 detik | Detik | Transaction completed |
| Saldo Wallet | ⚠️ Conditional | Bisa cache read projection untuk UI, tapi tidak untuk transactional operation | 1–5 detik | Detik | Transaction completed |
| Stock Sparepart | ⚠️ Conditional | Display UI boleh, validation final harus DB + concurrency control | 1–5 detik | Detik | Purchase part, invoice/service part line posted |

**Nuance:** Cache? bukan sekadar YES/NO; beberapa kasus adalah conditional dengan trade-off yang complex. Semua entri dapat menggunakan read projection cache untuk UI display, selama tidak digunakan untuk transaction authorization.

**Catatan TTL:** Nilai TTL pada contoh hanyalah illustrative. TTL yang sebenarnya bergantung pada business requirement dan invalidation mechanism yang tersedia. Active invalidation yang reliable dapat memungkinkan TTL jauh lebih panjang untuk master data.

---

## 23. Exercises

### Design Exercise — Dashboard Bengkel

Dashboard menampilkan: Jumlah Invoice, Pendapatan, Top Mekanik, Top Sparepart, Kendaraan Baru, Customer Aktif.

**Setiap data, tentukan:**

1. **Layak di-cache atau tidak?** (Ya/No/Conditional)
2. **Mengapa?** (Workload, cost, consistency need)
3. **Acceptable staleness?** (0s, detik, menit, lebih lama)
4. **Contoh TTL yang masuk akal?** (berikan rentang, bukan angka tunggal)
5. **Cache key seperti apa?** (format spesifik)
6. **Dimensi apa yang wajib masuk key?** (tenant, branch, date, dll)
7. **Kapan cache harus di-invalidasi?** (event bisnis)
8. **Event bisnis yang memicu invalidation?** (spesifik)
9. **Apa yang terjadi jika Redis down?** (fallback behavior)
10. **Bagaimana mencegah cache stampede?** (singleflight, lock, jitter)
11. **Bagaimana menjaga tenant/branch isolation?** (key format)
12. **Metrics apa yang perlu diamati?** (hit/miss, latency, error)

**Rubric / Expected Reasoning:**

Untuk setiap data:

| Item | Expected Content |
|------|------------------|
| Cache Decision | Berdasarkan read frequency, cost, staleness acceptability |
| TTL | Diberi konteks (operational vs transactional vs static) |
| Invalidation | Spesifik event bisnis, bukan sekadar "update" |
| Redis Down | Fallback strategy yang jelas |
| Stampede | Salah satu defense mechanism |

**Contoh Penyelesaian — Dashboard Pendapatan Hari Ini:**

- **Cache?:** Bisa layak karena aggregation mahal
- **Mengapa?:** Read-heavy dashboard, query join/aggregation, toleran detik-stale untuk operational metric
- **Acceptable Staleness?:** 10–30 detik (operational reporting)
- **Contoh TTL?:** Contoh: 15-30 detik (illustrative, tergantung replay tolerance)
- **Cache Key?:** `dashboard:tenant:{id}:branch:{id}:revenue:2026-09-01`
- **Dimensi?:** tenant, branch, date
- **Invalidation?:** Setelah payment baru, invoice diterima
- **Event?:** Payment success, Invoice status update
- **Redis Down?:** Fallback to DB (degraded performance), cache nilai untuk beberapa cycle
- **Stampede?:** Singleflight + jitter TTL
- **Isolation?:** tenant+branch scope
- **Metrics?:** cache_hit, cache_miss, db_fallback, latency

**Catatan:** TTL contoh 10-30 detik hanya illustrative. TTL yang sebenarnya bergantung pada:
- Acceptable staleness business
- Update frequency
- Invalidation mechanism yang tersedia
- Replay tolerance user

---

## 24. Common Mistakes

1. **Cache Everything** — Tidak mengukur cost vs benefit
2. **Missing Invalidation Strategy** — Remember forever anti-pattern
3. **Missing Isolation** — Global key bocor tenant
4. **Cache as Single Point of Failure** — Tidak desain fallback DB
5. **Key Explosion** — Cache berdasarkan per-user ID bila data identik
6. **Assume Cache = Strong Consistency** — Write-through tidak menghilangkan race condition
7. **Ignoring Double-Check Pattern** — Bisa miss optimization opportunity

### Anti-Pattern: Cache::rememberForever

```php
Cache::rememberForever('key', fn () => ...);
```

**Masalah:** Tanpa TTL, stale data tidak pernah autoreset.

```
Admin update harga
↓
Database = harga baru
Redis = harga lama
↓
tidak ada TTL
↓
user terus melihat harga lama
```

**Remember forever bukan selalu salah** — Ia dapat digunakan jika:
- **Deterministic invalidation tersedia**: setiap mutation yang mempengaruhi result memicu invalidation
- **Cache memang layak long-lived**: data yang relatif stable/master/reference data
- **Invalidation reliability sesuai requirement**: bisa berupa synchronous delete, versioned namespace, revision-based key, atau durable event

Rule yang benar: `rememberForever` aman hanya jika freshness requirement dan invalidation mechanism memiliki reliability yang sesuai.

**Tanpa invalidation yang reliable, rememberForever adalah high-risk pattern** karena:
- Loss event = data stale selamanya
- Debugging stale data bisa berjam-jam

Bounded TTL sering menjadi safety net yang lebih forgiving jika perfect invalidation sulit dijamin, tidak harus "better = short TTL" secara universal.

**Better alternative:** gunakan TTL dengan duration yang sesuai business requirement + background refresh / proactive invalidation.

---

## 25. Senior Engineer Takeaways

Cache yang baik mengurangi pekerjaan berulang tanpa bergantung pada data yang mungkin stale.

**Mindset Junior:** "Query lambat, cache saja."

**Mindset Senior:**
1. Ukur bottleneck (latency network vs DB CPU)
2. Optimalkan query di level storage / index
3. Pahami business consistency requirement
4. Tentukan acceptable staleness
5. Desain fallback & invalidation flow
6. Hitung operational complexity
7. Tambahkan cache jika benefit throughput signifikan

---

## 26. Running the Lab

### Prerequisites
- Go 1.22+
- Docker Compose

### From Repository Root

```bash
# Start infrastructure (Redis)
docker compose up -d redis

# Run tests
make lab-04-test
make lab-04-test-race
make lab-04-vet
make lab-04-demo
make lab-04-integration
```

### Manual Commands

```bash
# Start Redis (from repo root)
docker compose up -d redis

# Navigate to lab directory
cd labs/04-caching

# Run unit tests
go test -v -count=1 ./...

# Run tests with race detector
go test -race -count=1 ./...

# Run vet/lint
go vet ./...

# Run demo scenarios
go run ./cmd/demo -scenario=without-cache
go run ./cmd/demo -scenario=cache-aside
go run ./cmd/demo -scenario=stampede-unprotected
go run ./cmd/demo -scenario=stampede-protected
go run ./cmd/demo -scenario=write-through
```

Catatan: `docker compose` dijalankan dari repository root, bukan di dalam `labs/04-caching`.

---

## Navigation

- **Previous:** [Lab 03 — Distributed Transaction](../03-distributed-transaction/)
- **Next:** [Lab 05 — Race Condition](../05-race-condition/)