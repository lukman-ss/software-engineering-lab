# Lab 04 — Caching: Mengurangi Latency, Tapi Apa Biayanya?

> **Mental Model**: Caching adalah trade-off antara **latency** (cepat) dan **consistency** (benar). Senior engineer memilih teknik caching yang tepat untuk workload-nya, bukan menerapkan caching di mana-mana.

---

## Learning Objectives

Setelah menyelesaikan lab ini, Anda akan memahami:

- **Kapan** caching layak digunakan dan **kapan tidak**
- Cache Aside pattern & implementasinya
- Cara mengatasi *Cache Stampede* dengan Single Flight
- TTL strategy & jitter untuk stampede mitigation
- Cache Key design untuk multi-tenant & versioning
- Graceful degradation saat cache down
- Source of Truth: Database vs Cache role
- Cache Invalidation strategy
- Write Through pattern & failure modes
- Stale data & read-your-writes consistency
- Cache vs Session
- Heat map bottleneck: Database → Redis tradeoff
- Time-zone aware caching untuk business day

---

## Problem

Dashboard workshop menampilkan statistik dilihat oleh banyak user. Tanpa cache, setiap request menghasilkan *heavy aggregation* ke database:

**Dalam simulasi lab ini:**

```
[500 concurrent users] → Dashboard Request
                              ↓
                    Database: 6 queries + join/aggregation per request
                              ↓
                           3000 total DB queries
```

---

## Why Cache Exists

```
┌────────┐     ┌──────────┐     ┌───────┐
│ Client │ ──► │  YourApp │ ──► │  DB   │
└────────┘     └──────────┘     └───────┘
                     │               ✗
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

## Source of Truth

| Layer | Technology | Role |
|-------|------------|------|
| **Primary** | PostgreSQL | Permanent, ACID, can rebuild cache |
| **Cache** | Redis | Derived, TTL-bound, on-demand rebuild |

Redis down → Cache miss → Fallback to DB → Request sukses (degraded mode).

---

## Architecture

```
                    ┌─────────────────┐
                    │     Client      │
                    └────────┬────────┘
                             │ Request
                             ▼
┌────────────────────────────────────────┐
│        Application (Go + context)      │
│  ┌──────────────┐ ┌─────────────────┐ │
│  │  CacheAside  │ │  SingleFlight   │ │
│  └──────┬───────┘ └────────┬────────┘ │
└─────────┼──────────────────┼──────────┘
          │                  │
          ▼                  │
    ┌─────┴─────┐            │
    │   Redis   │            │ Cache Hit
    └───────────┘            │
          │                  │
          ▼                ┌─┴─┐
    ┌─────┴─────┐        │   │
    │   Cache   │ ◄──────┘   │ MISS
    └─────┬─────┘           │
          │                 │
          ▼                 │
    ┌─────┴─────┐           │
    │  Database │ ──────────┘ (rebuild)
    │ (PostgreSQL)          │
    └───────────┘           │
```

---

## Without Cache

```
[Client]
    │
    ▼
[Application]
    │
    ▼
[Database Query] → Compute Aggregation → Return
    │                    ↑
    └────────────────────┘ (50ms + CPU + network)
```

**Result:** Setiap request = heavy DB query.

---

## Apa yang Sebaiknya Di-cache?

Prinsip utama: **bisa di-cache** ≠ **harus di-cache**. Cache hanya layak jika semua empat syarat berikut terpenuhi:

1. **Mahal dihitung** — Query melibatkan join, aggegasi, atau I/O berulang yang signifikan
2. **Sering dibaca** — Volume read jauh melebihi volume write (read-heavy workload)
3. **Jarang berubah vs. frekuensi baca** — Rasio read:change jauh lebih tinggi dari 1
4. **Toleran terhadap sedikit staleness** — Sesaat data ketinggalan masih dapat diterima
5. **Data identik untuk request berulang** — Hasil yang sama untuk parameter yang sama

### Contoh pada Sistem CMMS / Bengkel

| Kategori | Contoh | Mengapa Cocok Di-cache | Alasan Engineering |
|---------|--------|----------------------|-------------------|
| **Dashboard** | Statistik harian bengkel (jumlah invoice, revenue, mekanik teratas) | Diakses oleh banyak user sekaligus, computed lewat 6+ query join | 1 read batch menggantikan 3000 query; user toleran data 30s–2min |
| **Statistik** | Laporan penjualan per minggu/bulan | Aggregation result jarang berubah, query expensive | Aggregate query di DB mahal; cache sekali compute |
| **Laporan Agregasi** | Laporan sparepart terjual, omzet per mekanik | Hasil agregasi, bersifat idempoten untuk jendela waktu | Same data for range; recompute hanya periode berubah |
| **Konfigurasi** | Setting sistem, tarif service, pajak | Dibaca setiap request, diupdate kali/kadang | Config reads >> writes; invalidate on update |
| **Master Data** | Daftar cabang, jenis service, master supplier, daftar mekanik | Read-heavy, write sangat jarang | Static reference data; cache months + version for migration |
| **Template Invoice** | Template HTML untuk pencetakan invoice | Static template, digenerate kali/kadang | Template rarely changes; cache long TTL |
| **Menu** | Menu navigasi per role | Dibaca setiap page load, berubah kali/kadang | Menu rendering overhead eliminated |
| **Permission** | Role-permission lookup | Read pada setiap request otorisasi | ⚠️ **Perhatian:** Security-sensitive. Gunakan TTL sangat pendek (5–10s) + active invalidation saat role berubah. Pertimbangkan jangan cache jika consistency requirement ketat. |
| **Exchange Rate** | Kurs mata uang | Diupdate periodik (mis. tiap 5 menit), dipakai banyak transaksi | Cache with TTL < update interval; stale rate = acceptable selama TTL |

### Engineering Rationale

- **Dashboard & statistik**: 6 query dengan join + aggregation dapat dikonsolidasikan menjadi satu compute. Cache mengubah O(N×6) query menjadi O(1) query per periode. Cost: sedikit staleness (30s–2min) yang acceptable untuk operational metrics.
- **Master data**: Sifatnya immutable dalam jangka pendek. Cache hit mengurangi join ke lookup table. Key versioning (`entity:id:v1`) memungkinkan invalidation via key rotation saat data berubah.
- **Permission**: Trade-off paling sensitif. Stale permission = security vulnerability. TTL 5–10s + invalidation pada role change meminimalkan exposure window, tapi bukan eliminasi total. Evaluasi: apakah exposure window < business risk tolerance?

---

## Apa yang Sebaiknya Tidak Di-cache?

Berikut data yang **umumnya tidak layak di-cache secara langsung**, karena membutuhkan **strong consistency** atau berubah sangat cepat:

| Kategori | Contoh | Mengapa Berbahaya |
|---------|--------|-------------------|
| **Saldo** | Saldo wallet, saldo kas | Inconsistency berdampak langsung pada keuangan. Over/under-billing. |
| **Status Transaksi** | Status pembayaran, status service | Race condition dapat menyebabkan duplikat pembayaran atau service gratis. |
| **OTP / Token Sekali Pakai** | Kode OTP, token reset password | Stale OTP dapat menyebabkan security bypass. |
| **Progress** | Progress upload, progress download | User experience terputus jika stale. |
| **Live Data** | Live tracking kendaraan, stock real-time | Business decision basis data real-time. Stale = keputusan salah. |
| **Inventory During Transaksi** | Ketersediaan sparepart saat checkout | Race condition dapat menyebabkan overselling. |

### Nuance Penting: "Tidak Boleh Di-cache" ≠ "Tidak Boleh Di-cache Sama Sekali"

Beberapa data real-time **masih dapat menggunakan cache sebagai optimization layer**, asalkan:

1. **Database tetap source of truth (canonical source)**
2. **Correctness tidak bergantung pada cache** — critical validation selalu ke database
3. **TTL sangat pendek** — di bawah ambang kehausan bisnis (mis. 1–5 detik)
4. **Invalidation yang benar** — invalidate segera saat data berubah
5. **Critical validation tetap ke authoritative storage** — database dengan concurrency control

#### Contoh: Stock Real-time vs Cache

```markdown
Use Case: "Cek ketersediaan sparepart"

❌ Salah: 
   Cache stock = 5 untuk 30 menit
   User A dan B lihat stock = 5
   User A checkout → cache masih 5
   User B checkout → SOLD OUT tapi system keliru accept
   
✅ Benar:
   - Cache menampilkan "approximate stock" (5) untuk UI
   - TTL 1–5 detik (hanya untuk tampilan)
   - Di backend, checkout selalu:
     1. Query DB dengan stock check + transaction:
        SELECT stock FROM parts WHERE id = $1 FOR UPDATE
     2. Jika stock > 0: kurangi (atomic)
     3. Commit → Invalidate cache
   - Cache boleh stale untuk display, tapi DB validation final selalu otoritas
```

### Engineering Rationale

- **Saldo / status transaksi**: Setiap nilai harus konsisten. Cache menambahkan failure mode (invalidation lag, crash sebelum invalidate) yang bisa menyebabkan kerusakan finansial. Lebih baik bayar latency cost ke DB.
- **Stock during transaksi**: Cache boleh untuk approximate display (UX), tapi ketersediaan final harus dicek ke DB dengan `SELECT FOR UPDATE` atau comparable concurrency control. Cache di sini = **read-through optimization**, bukan authoritas.

---

## Cache Patterns

### Cache Aside

```
[Client]
    │
    ▼
[App: Redis GET]  → HIT → [Return Cached]
    │ MISS
    ▼
[DB Query] → Compute → [Redis SET w/ TTL] → [Return]
```

**Flujo:** Read-through cache. Cache cek dulu; jika miss, app query DB dan populate cache.

**Kode (lihat `cache_aside.go`, `dashboard_cache.go`):**

```go
// 1. Check cache
cached, err := cache.Get(ctx, key)
if err == nil && cached != "" {
    var p Product
    if err := json.Unmarshal([]byte(cached), &p); err == nil {
        return p, nil // cache hit
    }
    _ = cache.Delete(ctx, key) // corrupt → delete
}

// 2. Cache miss → query DB
p, err := repo.GetProduct(ctx, id)
if err != nil {
    return Product{}, err
}

// 3. Populate cache
data, _ := json.Marshal(p)
_ = cache.Set(ctx, key, string(data), ttl)
return p, nil
```

**Pros:** Sederhana, fleksibel (bisa populate cache setelah berhasil), tidak butuh dependency tambahan.
**Cons:** Stale data sampai TTL/invalidation; cache stampede pada mass expiration.

### Cache Invalidation

Invalidation adalah mekanisme **delete** atau **update** cache entry setelah data berubah, sehingga request berikutnya tidak lagi membaca data stale.

#### Strategi Invalidation

1. **TTL Expiration**
   - Cara paling sederhana: cache entry otomatis expire setelah TTL.
   - Cocok untuk data dengan toleransi staleness tinggi (dashboard, statistik).
   - Trade-off: jika TTL terlalu panjang, data bisa sangat stale. Terlalu pendek = hit rate rendah.

2. **Explicit Delete Setelah Write (Cache-Aside + Invalidate)**
   ```
   [Client]: Pay invoice → [App: UPDATE DB] → [App: COMMIT ✅] → [App: DELETE cache] → [Return]
   ```
   **Invariant:** DELETE cache **hanya setelah** COMMIT DB sukses. Jika delete sebelum commit dan commit gagal, cache akan mengacu pada data yang sebenarnya tidak ada (stale miss / data loss).

3. **Versioned Cache Key**
   Bump versi key saat schema/data berubah signifikan.
   ```go
   // Sebelum: cmms:dashboard:v1:branch:12:2026-08-29
   // Setelah update: cmms:dashboard:v2:branch:12:2026-08-29
   // Cache lama (v1) akan expire naturally; tidak perlu DELETE.
   ```
   **Pros:** Atomic, tidak ada race window antara delete dan read.
   **Cons:** Temporary duplicated data volume (v1 + v2 sampai v1 expire).

4. **Event-Driven Invalidation (Pub/Sub)**
   - Domain event (`product.updated`) di-publish setelah DB commit.
   - Semua instance yang subscribe event akan invalidate cache mereka.
   - Cocok untuk arsitektur microservices dengan banyak instance.

5. **Tag/Group Invalidation (Jika Storage Mendukung)**
   - Redis tidak native mendukung tag invalidation, tapi beberapa library (mis. `redis-py` dengan `SCAN` + pattern delete) bisa mengelompokkan key dengan prefix/tag.
   - Contoh: `tag:product:123` menyimpan reference ke `product:123:v1`, `products:list`, `products:popular`. Saat tag dihapus, semua key terkait ikut invalid.

6. **RememberForever Anti-Pattern**
   ```go
   // ⚠️ DANGER: cache key tanpa TTL
   cache.Set(ctx, key, data, 0) // atau TTL = infinity
   ```
   **Mengapa berbahaya:**
   - Data tidak pernah expire otomatis.
   - Jika admin update harga produk:
     - Database = harga baru ✅
     - Redis = harga lama (stale sampai di-delete manual)
   - User checkout: harga $100 → Revenue loss ❌

   **Contoh Bug:**
   ```
   T1: Admin update Produk A: $100 → $150 (DB ✅)
   T2: Redis cache: $100 (stale, never expires)
   T3: User checkout pakai harga lama → loss
   ```

   **Only safe if:**
   - Deterministic invalidation (setiap update selalu trigger invalidate)
   - Event-driven invalidation (publish-subscribe)
   - Data hampir immutable (kode negara, enum status)

7. **Background Refresh (Stale-While-Revalidate)**
   - Saat TTL hampir habis, satu worker (singleflight) mem-fetch data fresh dan update cache.
   - User yang datang saat refresh tetap dapat data stale (slightly) tanpa menunggu.

#### Dependency Invalidation: Masalah yang Paling Sulit

```markdown
Satu Product 123 bisa muncul di banyak cache key:
- product:123
- products:list
- products:popular
- products:category:5
- products:by-supplier:42

Jika Product 123 berubah, idealnya kita invalidate semua key di atas.
Tapi bagaimana kita tahu key mana yang menyimpan Product 123?

Options:
1. Track reverse index (product:123 → [products:list, products:popular, ...])
   - Cons: Reverse index harus dijaga sendiri (overhead write)
2. Pattern-based delete (SCAN product:*) 
   - Cons: SCAN expensive, bukan atomic
3. Versioned key untuk semua affected keys
   - Cons: setiap key = satu set baru, memory naik
4. Tag-based invalidation dengan proxy
   - Cons: Tergantung dukungan library/storage
```

**Inilah salah satu alasan cache invalidation dianggap salah satu masalah tersulit dalam computer science (bersama naming things dan off-by-one errors).** Tidak ada solusi universal; setiap strategi punya trade-off.

#### Invalidation Checklist

| Scenario | Strategi | TTL Safety Net |
|----------|----------|----------------|
| Update invoice (dashboard affected) | Invalidate dashboard key | 30s |
| Update product master | Invalidate product key + dependency key | 5m |
| Update permission | Invalidate permission key + TTL 10s | 10s |
| Config change | Invalidate config key + versioning | N/A (manual invalidation) |

---

## TTL (Time-To-Live)

### Strategi TTL Berdasarkan Volatilitas

| Field | TTL | Reason |
|-------|-----|--------|
| InvoiceCount | 30s | Near transaction time |
| Revenue | 30s | Same as invoice count |
| TopMechanic | 2min | Changes slowly |
| TopSparepart | 2min | Changes slowly |
| VehicleCount | 15s | Usually at day start |
| ActiveCustomer | 1min | Moderate changes |

### TTL Jitter untuk Stampede Mitigation

```
Without Jitter:
100 keys created at T0
All expire at T0 + 60s  ← Stampede!

With Jitter (baseTTL + random(0..15s)):
Keys expire between T0 + 60s to T0 + 75s
```

**Helper:**
```go
func TTLWithJitter(base, maxJitter time.Duration) time.Duration {
    jitter := time.Duration(rand.Int63n(int64(maxJitter)))
    return base + jitter
}
```

### Probabilistic Early Refresh

TTL saja tidak cukup. Saatnya pertimbangkan **probabilistic early refresh**: ketika TTL masuk ke window terakhir (20% dari total TTL), refresh secara probabilistik untuk menghindari masih terjebak pada stale data hingga benar-benar expire.

```go
func ShouldRefreshEarly(now, expiry time.Time, originalTTL time.Duration, probability float64, randomFloat func() float64) bool {
    if expiry.IsZero() { return false }
    rem := expiry.Sub(now)
    if rem < 0 { return false } // expired, handled elsewhere
    refreshWindow := originalTTL - (originalTTL * 80 / 100) // 20% of TTL
    if rem <= refreshWindow {
        return randomFloat() < probability
    }
    return false
}
```

---

## Stale Data & Read-Your-Writes Consistency

### Apakah Stale Data Masalah?

Pertanyaan kuncinya bukan "apakah data berubah?" tetapi **"berapa lama stale data dapat diterima?"**

| Data | Max Staleness | Reasonable? |
|------|--------------|-------------|
| Dashboard statistik | 30s–2min | ✅ Ya, operational metrics |
| Stock display | 1–5s | ✅ Ya, untuk UI/UX saja |
| Saldo wallet | 0s (immediate) | ❌ Tidak, keuangan |
| Status pembayaran | 0s | ❌ Tidak, dapat duplikat/refund |
| Permission | 5–10s | ⚠️ Hanya dengan invalidation + short TTL |

### Read-Your-Writes Problem

Cache-aside sederhana tidak menjamin **read-your-writes**: user yang baru saja melakukan update mungkin masih membaca data lama dari cache.

```markdown
Timeline:
T1: User A update profile (DB commit sukses)
T2: User A immediately view profile
T3: App read dari cache → masih stale (cache belum diinvalidate)
T4: User A lihat data lama ❌
```

**Solusi:**
1. **Invalidate-on-write**: Setelah DB commit, invalidate cache key segera.
2. **Write-through**: Update cache sekaligus dengan DB (lihat section berikut).
3. **Read-through dengan staleness budget**: Jika tidak memungkinkan read-your-writes, ini *by design* — dokumentasikan sebagai SLA.
5. **Hybrid**: Kombinasi invalidate + very short TTL (5–10s) + versioning.

---

## Cache Stampede

### Stampede Scenario (Unprotected)

```
             Cache EXPIRES
                    │
           100 Requests
                    │
        ┌─────────┴──────────┐
        │                    │
        ▼                    ▼
   DB Query              DB Query
        │                    │
        └────────────────────┘
        │
   ┌────┴────┐
   │  DB Overload!  │
   └──────────┘
```

### Protected (Single Flight)

```
             Cache EXPIRES
                    │
           100 Requests
                    │
        ┌──────────┴───────────┐
        │ singleflight.DoChan()│
        └──────────────────────┘
                    │
                    ▼
             1 DB Query Only
                    │
                    ▼
          Shared Result to ALL
```

> **Note**: `singleflight` bekerja **dalam satu proses Go (intra-process)**. Untuk koordinasi antar instance aplikasi (lintas server), gunakan Distributed Lock.

### Single Flight

Gunakan `DoChan` agar panggilan mendukung *context cancellation* dengan aman.

```go
var flight singleflight.Group{}

func GetData(ctx context.Context, key string) (Dashboard, error) {
    ch := flight.DoChan(key, func() (interface{}, error) {
        return fetchFromDB()
    })

    select {
    case <-ctx.Done():
        return Dashboard{}, ctx.Err()
    case res := <-ch:
        if res.Err != nil {
            return Dashboard{}, res.Err
        }
        return res.Val.(Dashboard), nil
    }
}
```

Kode lengkap di `stampede.go` — meliputi `BrokenStampedeService` (bukti stampede terjadi) dan `ProtectedStampedeService` (dengan singleflight + double-check cache di dalam gate).

---

## Write Through

### Flow

```
[Request update]
        ↓
   Update Database      ← Source of Truth (authoritative)
        ↓
   Update Cache         ← Keep cache in sync
        ↓
    Return
```


Distributed lock mengkoordinasi akses antar instance aplikasi berbeda.

```
Instance A                          Redis                   Instance B
    │                                 │                         │
    ├── SET lock:1 tokenA NX PX 10k ─►│                         │
    │                                 ├── success              │
    │◄─ OK ───────────────────────────┤                         │
    │                                 │                         │
    ├── Rebuild Cache                 ├── SET lock:1 tokenB NX ─┤
    │                                 │                         │
    │                                 ├── false                │
    │◄───────────────────────────────┼────────────────────────►│
    │                                 │                         │ wait/retry/fallback
    │                                 │
    ├── Lua: if GET(lock) == tokenA   │
    ├── then DEL(lock) ──────────────►│
    │                                 │
    │◄─ OK (lock released) ───────────┤
```

*Note: Token digunakan untuk membuktikan ownership lock, sehingga Instance A tidak secara tidak sengaja menghapus lock milik instance lain jika A terlambat mengeksekusi.*

Kode lengkap di `distributed_lock.go`.

---

## Cache Key Design

### Format

```
namespace:entity:vVersion:tenant:identifier:date
```

### Examples

- Dashboard: `cmms:dashboard:v1:branch:12:2026-08-29`
- Product: `product:123:v1` (versioned for migration)

### Multi-Tenant Isolation

| Tenant | Key |
|--------|-----|
| Tenant A, Branch 1 | `cmms:dashboard:v1:tenant:1:branch:1:2026-08-29` |
| Tenant B, Branch 1 | `cmms:dashboard:v1:tenant:2:branch:1:2026-08-29` |

### Versioning untuk Migration

Bump versi key saat schema berubah. Key lama akan expire naturally; tidak perlu mass-delete.

```go
keyV1 := NewDashboardKey(42).Build()        // cmms:dashboard:v1:...
keyV2 := NewDashboardKey(42).WithVersion(2).Build() // cmms:dashboard:v2:...
```

Kode lengkap di `dashboard_key.go`.

---

## Cache vs Session

| Aspect | Cache | Session |
|--------|-------|---------|
| **Purpose** | Cache data yang sama untuk semua user (read optimization) | Simpan state per-user (whoami, auth, UI preferences) |
| **Scope** | Global / shared | User-specific |
| **TTL** | Ditentukan oleh volatilitas data | Ditentukan oleh idle timeout / session lifetime |
| **Key** | `entity:id:version` | `session:{uuid}` |
| **Contoh** | Dashboard statistik semua bengkel | User login state, shopping cart |

**Overlap yang sering disalahartikan:**

```markdown
❌ Salah: Cache permission user di Redis shared cache
   Key: user:123:perms
   Masalah: Jika ada bug key isolation, User A bisa baca perms User B

✅ Benar: Simpan di session store (bisa Redis, tapi session-scoped)
   Key: session:{sid}:perms
   Session ID di-set via cookie, otomatis terisolasi per-user
```

**Aturan:** Jika data bersifat **per-user**, simpan di session. Jika data bersifat **global/bisa dibagi**, pertimbangkan cache. Jika data **per-user tapi read berulang**, gunakan **session dengan TTL manager** (bukan cache key manual).

---

## Write Through

### Flow

```
Request update
      ↓
Update Database      ← Source of Truth (authoritative)
      ↓
Update Cache         ← Keep cache in sync
      ↓
   Return
```

Setelah DB transaction COMMIT, update cache secara **atomik dalam aplikasi layer**. Request selanjutnya mendapatkan data terbaru dari cache — tidak ada window stale antara write dan next read.

### Tujuan

Write-through memastikan **cache dan database selalu konsisten setelah write** — tidak ada window di mana cache menjadi stale. Berguna untuk data yang:

- Sering dibaca tapi butuh konsistensi relatif tinggi
- Read-your-writes penting (user update → melihat hasilnya segera)
- Cache miss storm mahal (contoh: product catalog yang dipakai checkout)

### Perbedaan dengan Cache-Aside

| Aspect | Cache-Aside | Write-Through |
|--------|-------------|---------------|
| **Read path** | App check cache → miss → DB | App check cache → miss → DB (same) |
| **Write path** | Update DB → invalidate cache | Update DB → update cache |
| **Staleness** | Window stale sampai invalidation | Tidak ada (cache sync dengan DB) |
| **Failure: cache gagal on read** | Fetch from DB, populate cache | Sama |
| **Failure: cache gagal on write** | Stale sampai TTL — user baca data lama | DB updated, cache stale — hapus/invalidate |

### Kelebihan

1. **Read-your-writes guaranteed** — setelah update, semua request membaca data baru dari cache.
2. **No invalidation needed** — cache update atomik dengan write.
3. **Tidak ada stale window** — berbeda dengan invalidate-after-write yang punya race window.

### Kekurangan

1. **Complexity pada failure mode** — DB commit sukses tapi cache gagal (Redis down). Perlu strategy fallback.
2. **Cache write load** — setiap update product = 2 operasi (DB + cache). Jika update sering, ini menambah write amplification.
3. **Tidak semua write path konsisten** — jika ada path lain yang update DB langsung (bypass WriteThroughService), cache di sana tidak akan terupdate.

### Risiko: Partial Failure

```markdown
Scenario: Update harga product sparepart

1. DB transaction COMMIT → SUCCESS ✅
2. Redis SET → FAILURE ❌ (Redis down / network partition)

Result:
- Database punya data terbaru ✅
- Cache masih punya data lama (atau key tidak ada)
- Request berikutnya: cache hit → stale data ❌

Mitigation:
- Pada WriteThroughService, jika cache SET gagal:
  1. Log error (structured logging, bukan fmt.Printf — lihat Common Mistakes)
  2. Invalidate (delete) key → memaksa cache miss berikutnya fetch fresh dari DB
  3. Request tetap return success (DB adalah authoritas)
- TTL tetap berlaku sebagai safety net untuk kasus crash antara DB commit dan cache update.
```

### Cache Tidak Boleh Dianggap Source of Truth

Kunci: **Cache is a optimization layer, not the system of record.** Prinsip ini penting untuk failure mode:

```markdown
❌ Salah:
   - "Redis SET gagal → return error ke user"
   - User bayar tapi sistem bilang error karena cache write gagal

✅ Benar:
   - "Redis SET gagal → log error, invalidate key, return success"
   - User bayar, DB terupdate. Cache akan catch up via next read.
```

### Urutan Operasi Database / Cache

Wajib: **Commit → Cache Update, JANGAN pernah Cache Update → Commit.**

```
Safe order:
┌────┐    ┌────┐    ┌────┐
│ DB │ →  │Commit│ →  │Cache│
└────┘    └────┘    └────┘
  ↑                  ↓
  ←── invalidate on failure

Unsafe order:
┌────┐    ┌────┐    ┌────┐
│Cache│ →  │ DB  │ →  │Commit│
└────┘    └────┘    └────┘
  (cache terupdate, tapi DB rollback?)
  → cache punya data tidak sinkron dengan DB!
```

### Apa yang Terjadi Jika Database Berhasil Tapi Redis Gagal?

| Layer | Status | Action |
|-------|--------|--------|
| Database | Terupdate ✅ | Source of truth berubah |
| Redis | Gagal set ❌ | Cache stale atau tidak ada |
| App | Return success ✅ | Log error, invalidate cache key |
| Request berikutnya | Cache miss → DB | Populate cache fresh |
| Observability | Error counter + log | Alert bila error rate tinggi |

### Kapan Invalidation Lebih Aman Daripada Update Cache Langsung?

| Scenario | Lebih Baik dengan Invalidation | Reason |
|----------|-------------------------------|--------|
| Data berubah sering | ✅ Invalidate + short TTL | Cache write load meningkat |
| Data konsisten wajib (keuangan) | ✅ Invalidate | Avoid cache corruption risk |
| Cache di proses berbeda | ✅ Invalidate | No shared memory of cache value |
| Cache write reliability rendah | ✅ Invalidate | Redis flapping → write-through error |
| Cross-table denormalisasi | ✅ Invalidate | Update satu table butuh update banyak cache key |

### Contoh: Update Produk di CMMS

```go
// write_through.go
func (s *WriteThroughService) UpdateProduct(ctx context.Context, p Product) error {
    // 1. Update DB (Source of Truth)
    res, err := s.db.ExecContext(ctx,
        "UPDATE products SET name = $1, price = $2 WHERE id = $3",
        p.Name, p.Price, p.ID)
    if err != nil {
        return fmt.Errorf("update DB: %w", err)
    }
    rows, _ := res.RowsAffected()
    if rows == 0 {
        return fmt.Errorf("product not found")
    }

    // 2. Write-through ke Cache
    key := CacheKey("product", p.ID, 1)
    data, _ := json.Marshal(p)
    jitteredTTL := TTLWithJitter(s.ttl, s.jitterTTL)
    if err := s.cache.Set(ctx, key, string(data), jitteredTTL); err != nil {
        // DB sudah sukses, tapi cache gagal.
        // Log + invalidate (delete stale key) sebagai safety fallback.
        fmt.Printf("warn: write-through cache set failed for %s: %v\n", key, err)
        _ = s.cache.Delete(ctx, key) // safety: force cache miss on next read
    }
    return nil
}
```

Demo:
```bash
go run ./cmd/demo -scenario=write-through
```

---

## Cache vs Session

Redis dapat dipakai untuk **kedua** tujuan ini, tetapi semantics-nya sepenuhnya berbeda. Teknologi simpanan sama ≠ konsep sama.

| Aspect | Cache | Session |
|--------|-------|---------|
| **Purpose** | Optimization read (redundansi data yang sama untuk banyak request/user) | Simpan state per-user (auth, cart, preferences) |
| **Ownership** | Shared/global — dimiliki oleh sistem, bukan user | User/session-specific — dimiliki oleh satu user session |
| **Source of Truth** | Database / service lain | Database (token/session table) |
| **Consequence if Redis flushed** | Slower responses, DB load naik — **service tetap berfungsi** | Semua user terlogout — **service terputus secara fungsional** |
| **Key scope** | `product:123:v1`, `cmms:dashboard:v1:branch:12:date:2026-08-29` | `session:{uuid}`, `user:123:cart` |
| **TTL** | Berdasarkan volatilitas data (30s–24h) | Lifecycle user session (30min–24h, idle timeout) |
| **Data loss if evicted** | Acceptabel — di-recompute dari DB | Tidak dapat diterima — auth/cart hilang |
| **Read pattern** | Many readers, many writers, read-heavy | One session → one user → sequential access |

### Engineering Guidance

```markdown
Cache data per-user di Redis shared cache?
❌ Salah: Key = user:123:perms (bisa bentrok, user A baca perms user B)
✅ Benar: Key = session:{sid}:perms (session-scoped, cookie-based isolation)

Gunakan session store untuk:
- Authentication state
- Shopping cart
- UI preferences sementara
- Data yang benar-benar per-user dan harus persist over request lifetime

Gunakan cache untuk:
- Data read-heavy yang sama untuk banyak user
- Hasil aggregation yang expensive
- Master data yang jarang berubah
- Konfigurasi statis
```

---

## Contoh Keputusan Caching pada Bengkel / CMMS

Berikut keputusan cache untuk domain bengkel/CMMS. TTL ditentukan berdasarkan **acceptable staleness** dan **karakteristik update data**, bukan angka universal.

| Data | Cache? | Alasan | Contoh TTL | Invalidation Trigger |
|------|--------|--------|------------|----------------------|
| Daftar Cabang | ✅ Ya | Read-heavy, static reference data, rarely changes | 24h | Update master cabang (admin) |
| Jenis Service | ✅ Ya | Static, dipakai di setiap transaksi service | 24h | Update master service type |
| Daftar Mekanik | ✅ Ya | Read-heavy, update kala/kadang | 5m | Hire/fire mekanik, update profile |
| Master Supplier | ✅ Ya | Read-heavy untuk PO/purchase, rarely change | 1h | Update supplier profile |
| Template Invoice | ✅ Ya | Static template, cached long TTL | 1h | Update template (version bump) |
| Konfigurasi Pajak/PPN | ✅ Ya | Static setting, read di setiap invoice | 10m | Update tax configuration |
| **Permission** (per-user) | ⚠️ Conditional | Security-sensitive | **5–10s** | Role change, permission revocation (immediate invalidate) |
| Dashboard Pendapatan Hari Ini | ✅ Ya | Aggregation expensive, 6 query | 30s | New invoice, new payment |
| Top Mekanik | ✅ Ya | Aggregation, changes slowly | 2min | New service record (invalidate today's dashboard) |
| Top Sparepart | ✅ Ya | Aggregation, changes slowly | 2min | New parts usage in invoice |
| **Invoice Baru** (create) | ❌ Tidak | Write path, bukan read optimization | N/A | N/A |
| **Status Service** | ❌ Tidak (critical) | Read-your-writes penting, race condition risk | N/A | N/A |
| **Pembayaran** | ❌ Tidak (critical) | Keuangan, konsistensi kuat dibutuhkan | N/A | N/A |
| **Saldo Kas** | ❌ Tidak (critical) | Keuangan, tiap transaksi harus akurat | N/A | N/A |
| Stock Sparepart | ⚠️ Display only | Approximate OK for UI, final validation ke DB | 1–5s | New purchase, parts usage |

### Catatan pada TTL

- **30s dashboard**: Operational metrics — 30s stale masih acceptable untuk monitoring live.
- **5–10s permission**: Security boundary — setiap perubahan role harus propagated dalam ≤10s.
- **1–5s stock display**: UI/UX optimization, **bukan** untuk transaksi. Final stock validation selalu ke DB dengan `SELECT FOR UPDATE`.
- **5m mekanik**: Update kala/kadang (shift change, training). Stale 5m tidak berdampak pada service quality.
- **24h cabang/service**: Master data hampir immutable — cache months. Gunakan versioning (`v2`) untuk forced eviction.

### Data Kritis: Authoritative Read/Write Tetap ke Database

| Kritis Data | Cache Purpose | Database (Authoritative) |
|-------------|---------------|--------------------------|
| Stock Sparepart saat checkout | Approximate display | `SELECT FOR UPDATE` + decrement + commit |
| Payment amount | N/A | Payment gateway → DB → invalidate |
| Saldo kas | N/A | DB transaction |
| Status service | N/A | DB state machine |

---

## Memory & Cardinality: Key Explosion Anti-Pattern

Cache yang dikelola dengan baik juga harus **memory-efficient**. Salah satu masalah serius di Redis adalah **key explosion** — membuat terlalu banyak key yang tidak dibutuhkan.

### Contoh Anti-Pattern: Key Explosion

```markdown
❌ Salah: cache:user:{user_id}:dashboard

Masalah:
- 1 juta user → 1 juta dashboard key (padahal dashboard hanya 50 cabang!)
- Setiap user punya data yang sama → duplikasi data
- Memory usage tinggi → Redis OOM
- Eviction yang tidak terduga → cache miss tiba-tiba
- Tidak ada benefit reuse antar user

✅ Benar: cache:branch:{branch_id}:dashboard:{date}

Keuntungan:
- 50 cabang → 50 dashboard key per hari
- Semua user di cabang yang sama pakai dashboard yang sama
- Reuse tinggi → memory efficient
- Invalidation tepat sasaran (hanya cabang yang berubah)
```

### Prinsip Cardinality

> **Semakin spesifik cache key, semakin kecil sharing/reuse-nya dan semakin besar kemungkinan cardinality meningkat.**

| Key Pattern | Jumlah Key (1M user) | Reuse | Memory | Risk |
|-------------|---------------------|-------|--------|------|
| `dashboard:{user}` | 1,000,000 | Sangat rendah | Tinggi | OOM, eviction berlebih |
| `dashboard:{branch}` | 100 | Tinggi | Rendah | Aman |
| `products` | 1 | Sangat tinggi | Sangat rendah | Aman |
| `products_by_category:{category}` | 10 | Tinggi | Rendah | Aman |
| `user:{id}:recommendation:{date}` | 30M (30 hari) | Sangat rendah | Sangat tinggi | Anti-pattern |

### Cardinality vs Reuse

- **Personalized data** (recommendation, history): boleh high cardinality, tapi pertahankan TTL pendek + eviction
- **Shared data** (dashboard, master data, config): gunakan key yang agregat untuk maksimal reuse

### Memory Pressure Indicators

| Metric | Alert Threshold | Action |
|--------|-----------------|--------|
| `used_memory_pct > 80%` | Critical | Review key cardinality |
| `evicted_keys > 0` | Warning | TTL terlalu pendek / memory kurang |
| `keyspace_hits` vs `keyspace_misses` | Low hit ratio | Key spesifik terlalu banyak |

---

## Cache::remember (Laravel) — Caveats

```php
Cache::remember($key, $ttl, fn () => ...);
```

Framework seperti Laravel menyediakan helper `Cache::remember` yang memudahkan caching. Namun, penting memahami **limitasi**nya:

### Peringatan penting:

1. **Key harus merepresentasikan semua input**
   - ✅ `Cache::remember("user:{$id}", 300, fn() => ...)`
   - ❌ `Cache::remember("user", 300, fn() => ...)` (brute data user semua!)

2. **TTL bukan strategi invalidation satu-satunya**
   - TTL adalah safety net, bukan mekanisme utama
   - Jika data penting berubah, gunakan explicit invalidation
   - TTL-only = stale data mungkin menempel lebih lama

3. **Closure dapat dieksekusi bersamaan oleh beberapa request pada cache miss**
   - `Cache::remember` tidak otomatis mencegah stampede
   - Beberapa request masuk, semua lihat cache miss, semua eksekusi closure
   - **Potensi DB overload!**

4. **Cache::remember sendiri belum tentu mencegah stampede**
   - Behavior locking bergantung implementasi cache driver
   - `file`, `database`, `array`: biasanya tidak ada locking
   - `redis` dengan lua script atau mutex: bisa ada stampede protection

5. **Implementation-dependent behavior**
   - Laravel menulis ke cache setelah closure selesai
   - Jika Redis down, closure tetap dieksekusi (ngaco)
   - Tidak ada retry/rollback mekanisme khusus

### Contoh Stampede dengan remember

```php
// 100 request datang bersamaan, semua cache miss
Cache::remember('popular_products', 3600, fn() => Product::popular()->get());
// Tanpa stampede protection, semua 100 request eksekusi closure
// → 100 query ke database → DB overload
```

### Solusi:

- Gunakan Laravel `cache:remember` **dengan stampede protection** (middleware, lock, atau package seperti `predis/predis` + lua script)
- Atau gunakan `Cache::rememberForever` **hanya** untuk data yang hampir immutable + ada event-driven invalidation

---

## Senior Engineer Mindset: Caching Decisions

### Junior Engineer:
> "Query lambat, cache saja"

### Senior Engineer:

```
Apakah ada bottleneck?
  ↓ Tidak → Jangan cache
  ↓ Ya → Optimalkan dulu:
  1. Profil latency → mana yang lambat?
  2. Cek query → N+1? missing index?
  3. Optimalkan query sebelum cache
  ↓ Query tetap lambat → Cache kandidat
  ↓ Data konsistensi penting?
    ↓ Ya → gunakan write-through + invalidate
    ↓ Tidak → cache-aside + ttl
  ↓ TTL apa yang acceptable?
  ↓ Invalidation strategy jelas?
  ↓ Cache down bagaimana?
  ↓ Memory/cardinality mana yang masuk akal?
  ↓ Benefit > complexity?
    ↓ Ya → implementasi
    ↓ Tidak → batal atau cari alternatif
```

### Decision Framework: Does This Need to Be Cached?

| Question | Answer: NO | Answer: YES |
|----------|------------|-------------|
| Query heavy? | Skip cache | Cache candidate |
| Read-heavy (5:1+)? | Skip cache | Cache candidate |
| Stale acceptable? | Skip cache | Consider |
| Clear invalidation? | Skip cache | Strong candidate |
| Redis down fallback? | Skip cache | Ensure graceful degradation |
| Memory/cardinality okay? | Skip cache | Verify capacity |

### Prinsip Utama

1. **Cache yang baik mengurangi pekerjaan berulang** — tidak mengorbankan correctness.
2. **Database = Authoritative** — cache dapat hancur tanpa menghancurkan bisnis.
3. **Memory terbatas** — rancangan key harus mempertimbangkan cardinality.
4. **Failure adalah bukti desain** — jika Redis down crash aplikasi, cache keyinya buruk.
5. **Staleness adalah trade-off** — cari keseimbangan antara latency dan consistency.

---

## Mindset Senior Engineer

### Flow Ideal untuk Read

```
Client Request
      │
      ▼
   [Redis] ─X─ ERROR (Redis down / timeout)
      │
      ▼
[Fallback to Database] ← Source of Truth (always available)
      │
      ▼
  Request succeeds ✓ (slower, but works)
```

Cache adalah **optimization layer**, bukan dependency correctness utama.

### Scenarios Cache Gagal

| Scenario | Behavior | Handle |
|----------|----------|--------|
| Redis timeout | App block/slow | Circuit breaker → fallback |
| Redis connection refused | Immediate fail | Fallback to DB |
| Redis restart (transient) | Spike latency | Retry with backoff |
| Cache flush (accidental) | Mass miss → DB load spike | Monitor hit ratio drop, alert |
| Memory OOM | Redis crash | Detect → graceful degradation |

### Perbedaan Error Types untuk Observability

| Event | Metric | Action |
|-------|--------|--------|
| Cache miss (normal TTL expiry) | `cache_miss++` | Expected |
| Cache error (network, timeout) | `cache_error++` | Alert if rate > threshold |
| Read timeout | `redis_timeout` | Circuit breaker open |
| Fallback to DB | `db_fallback++` | Monitor rate, alert on spike |

### Engineering Guidance

```markdown
❌ Salah:
   - Log semua error sebagai "cache miss" → tidak bisa bedakan kegagalan teknis dari behavior normal
   - Menganggap "error" = "cache miss" → observability menjadi meaningless

✅ Benar:
   - `cache_miss`: Key tidak ada (TTL expire, key baru)
   - `cache_error`: Redis down, timeout, corrupt value, Set/Delete gagal
   - Monitoring: hit_ratio < 50% + error_rate > 1% = investigate cache tier
```

### Contoh Implementasi Graceful Degradation

```go
func GetDashboard(ctx context.Context) (Dashboard, error) {
    cached, err := cache.Get(ctx, key)
    if err != nil {
        // Redis error (down, timeout, corrupt) → fallback
        metrics.IncError()  // bukan incMiss
        return fetchFromDB(ctx)
    }
    if cached == "" {
        metrics.IncMiss()
        return fetchFromDB(ctx)
    }
    metrics.IncHit()
    // ... unmarshal ...
}
```

---

## Metrics

| Counter | Meaning |
|---------|---------|
| `cache_hit` | Cache returned valid data |
| `cache_miss` | Expected state when cache is empty |
| `cache_error` | Redis network error, down, corrupt value, or Set/Delete failure |
| `database_query` | Query successfully reached PostgreSQL |
| `cache_rebuild_attempt` | Attempted to populate cache from DB |
| `cache_rebuild_success` | Successfully populated cache from DB |

---

## Query Optimization → Index → Caching

Pipeline keputusan untuk performance optimization:

```
1. Query lambat?
   │
   ├── Profile query (EXPLAIN ANALYZE)
   │   ├── Table scan? → Add index
   │   ├── Index ada tapi tidak terpakai? → Perbaiki query (SARGability)
   │   └── Join mahal? → Consider denormalization
   │
2. Setelah indexed, query tetap lambat?
   │   ├── Hasil aggregation expensive? → Cache aggregation result
   │   ├── Data read-heavy? → Cache entity
   │   └── Data read-write balanced? → Pertimbangkan invalidate vs write-through
   │
3. Cache strategy selection:
   │   ├── Cache-Aside + invalidate: General purpose
   │   ├── Cache-Aside + write-through: Read-your-writes critical
   │   └── Write-through: Strong consistency needed
```

**Principle:** Caching adalah optimization **terakhir**, bukan pertama. Selalu exhaust index-based optimization sebelum cache. Cache justru dapat mem-perburuk performance jika query underlying lambat (cache miss storm → DB overload).

---

## Jangan Gunakan Cache untuk Menutupi Query yang Buruk

Cache **bukan pengganti database optimization**. Sebuah query yang buruk tidak akan diperbaiki dengan cache; cache hanya menyembunyikan gejala.

### Diagnostic Checklist (Sebelum Cache)

```
1. Ukur dulu → Profile query latency
2. Cek query → Gunakan EXPLAIN ANALYZE di Lab 02
3. Cek N+1 → Batch query, gunakan join/IN clause
4. Cek execution plan → Table scan? Index missing?
5. Tambah/fix index jika diperlukan → Kembali ke #1
6. Kurangi data yang dibaca → SELECT kolom perlu, bukan SELECT *
7. Optimalkan query → CTE, subquery, atau denormalize jika perlu
8. Setelah itu → Evaluasi caching untuk workload yang repetitif
```

### Contoh: SELECT * FROM branch

- Hanya ≈ 12 baris
- Latency < 1ms tanpa index (table small)
- **Tidak butuh cache** — overhead Redis GET + network > DB query

### Contoh: Dashboard Agregasi

- Query 6 buah aggregation + join
- Dijalankan ratusan kali per menit
- Latency 50ms+ → **Candidate wajib di-cache**

### Prinsip

> Cache adalah optimization **terakhir**, bukan cara memaksa query lambat tetap cepat. Jika query Anda tidak bisa diperbaiki tanpa cache, pertanyaan berikutnya: apakah itu data yang seharusnya di-cache?

---

## Multi-Tenant Cache Key Design

Untuk aplikasi multi-tenant seperti Bengkel/CMMS, cache key **harus mencakup semua dimensi yang mempengaruhi hasil**:

### Pola Key Design

```
{namespace}:{tenant_id}:{branch_id}:{resource}:{dimension}:{modifier}
```

### Contoh Kunci untuk CMMS

| Key | Struktur | Contoh |
|-----|----------|--------|
| Dashboard hari ini | `{app}:tenant:{t}:branch:{b}:dashboard:{date}` | `cmms:tenant:42:branch:7:dashboard:2026-09-01` |
| Daftar cabang | `{app}:tenant:{t}:branches:list` | `cmms:tenant:42:branches:list` |
| Detail produk | `{app}:tenant:{t}:product:{id}` | `cmms:tenant:42:product:123` |
| Statistik service | `{app}:tenant:{t}:branch:{b}:service:stats:{period}` | `cmms:tenant:42:branch:7:service:stats:2026-W36` |

### Dimensi yang Harus Di-Sertakan

| Dimensi | Contoh Penggunaan | Mengapa Perlu |
|---------|-------------------|--------------|
| **tenant/company** | `tenant:42` | Isolasi data antar konsumen |
| **branch** | `branch:7` | Setiap bengkel punya dashboard terpisah |
| **user** | Jika hasil personalized (recommendation) | Bukan global |
| **role/permission context** | Jika query tergantung level akses | Bisa bocor data sensitif |
| **timezone** | `timezone:Asia/Jakarta` | "Today" bersentuhan dengan zona waktu |
| **date** | `2026-09-01` | Business day, bukan UTC now |
| **query/filter** | `filter:status:active` | Hasil tergantung filter |
| **version** | `v2` | Migration schema, bukan invalidation |

### Contoh Bug Cross-Tenant

```markdown
❌ Salah: Key = "dashboard" (global tanpa tenant)
Tenant A, Branch 1 memiliki dashboard: invoice=100, revenue=10J
Tenant B, Branch 1 memiliki dashboard: invoice=50, revenue=5J
Cache global menyimpan salah satu: {"branch":1, "invoice":100, "revenue":10000000}

Request Tenant B → dapat data Tenant A ❌ Data leak!

✅ Benar: Key = "cmms:tenant:1:branch:1:2026-09-01"
Tenant A key berbeda dari Tenant B key → isolasi proper
```

### Cache Key Versioning Pattern

```go
// Untuk migrations: jangan pakai DELETE-heavy invalidation
keyV1 := fmt.Sprintf("cmms:dashboard:v1:tenant:%d:branch:%d:date:%s", tenantID, branchID, date)
keyV2 := fmt.Sprintf("cmms:dashboard:v2:tenant:%d:branch:%d:date:%s", tenantID, branchID, date)
// Key lama expire, key baru fetch fresh
```

---

## Permission Caching: Security Considerations

Permission adalah kandidat cache **bisa**, tapi sangat rentan security risk bila tidak di-handle dengan benar.

### When Permission Caching is OK

- **Read-heavy**: Permission dicek pada setiap request API
- **Short TTL**: 5–10 detik sebagai safety net
- **Active Invalidation**: Trigger delete on role change

### Contoh Bug Permission Cache

```
Timeline:
T1: User A login → cache:user:123:permission = ["admin"]
T2: Admin mencabut role User A → DB permission = [] (empty)
T3: User A refresh halaman → cache masih ["admin"] (TTL 10m lagi)
T4: User A tetap dapat akses admin! ❌ Security breach
```

### Mitigasi Security

1. **TTL sangat pendek** (5–10s) untuk authorization data.
2. **Active Invalidation** — publish event `user.role.changed` → semua instance invalidate `user:{id}:permission`.
3. **Key Versioning** — `user:123:permission:v{n}` — bukan hanya `user:123:permission`.
4. **Jangan pernah cache read-your-writes kritikal** — permission check penting pada endpoint write harus query DB.

### Permission Caching Strategy

| Data | Cache? | TTL | Invalidation |
|------|--------|-----|--------------|
| Menu access (read UI) | ✅ Ya | 10s | Role change event |
| API permission (critical write) | ❌ Tidak | N/A | Always query DB |
| Permission list (admin view) | ✅ Ya | 30s | Role assignment |

### Engineering Guidance

```markdown
✅ Benar: Cache untuk read-heavy UI data (menu, nav, read-only endpoints)
   - Key: session:{sid}:menu_cache
   - TTL: 10-30s
   - Invalidate: session logout

❌ Salah: Cache untuk authorization decision pada transaksi
   - Jangan gunakan cached permission untuk approve payment
   - Selalu query DB untuk operasi kritikal
```

---

## Rule of Thumb Sebelum Menggunakan Cache

1. **Apakah data mahal dihitung?** Cache hasil aggregation expensive query, bukan lookup murah.
2. **Read:Write ratio > 10:1?** Cache hanya berguna jika read jauh lebih banyak dari write.
3. **Berapa stale data yang dapat diterima?** Jika 0, jangan cache. Jika ada toleransi, cache bisa.
4. **Apakah ada invalidation strategy?** Tanpa invalidation, cache = stale sampai TTL (bisa sangat lama).
5. **Apakah DB bisa handle cache miss rate saat cache down?** Jadwalkan capacity: cache down = full traffic ke DB.
6. **Apakah multi-tenant?** Key harus isolate tenant dari hari pertama.
7. **Apakah user perlu read-your-writes?** Pertimbangkan write-through atau invalidate-on-write.
8. **Apakah ada monitoring untuk hit/miss/error?** Tanpa metrik, tidak bisa memvalidasi keputusan cache.

---

## Latihan Menentukan Cache Strategy

### Scenario 1: Update Harga Produk

Seorang admin mengupdate harga sparepart. Produk ini ditampilkan di katalog (dibaca ribuan user per hari) dan di invoice (dibaca saat checkout).

**Pertanyaan:**
- Strategi cache apa yang tepat?
- Berapa TTL yang reasonable?
- Bagaimana invalidation?

<details>
<summary>💡 Solusi</summary>

- **Strategi:** Write-through (harga produk penting untuk checkout)
- **TTL:** 5–10 menit (katalog tetap cepat, harga tidak perlu real-time sempurna)
- **Invalidation:** Setelah update harga, invalidate product key + gunakan versioning (`product:id:v2`) untuk forced eviction, OR write-through update cache langsung
- **Critical path:** Checkout selalu validate final price ke DB (authoritativeness over cache)
</details>

---

### Scenario 2: Live Tracking Kendara

Sistem live tracking posisi kendara sedang diperbaiki. Update GPS setiap 2 detar. Ditampilkan di dashboard operator.

**Pertanyaan:**
- Strategi cache apa yang tepat?
- Berapa TTL yang reasonable?

<details>
<summary>💡 Solusi</summary>

- **Strategi:** Cache **bukan pilihan** untuk live position. Gunakan:
  - WebSocket push dari GPS device → operator UI
  - Atau pub/sub (Redis) → stream position update
- **TTL:** 2–5 detik (untuk approximate position di map)
- **Alternative:** Jika harus cache, gunakan cache sebagai **buffer** dengan TTL < 2 detik, dan DB sebagai source of truth untuk historical replay
</details>

---

### Scenario 3: Laporan Bulanan

Laporan omzet bulanan di-generate sekali sehari pukul 01:00. Dibaca oleh manajer dan stakeholder (ratusan user per hari).

**Pertanyaan:**
- Strategi cache apa yang tepat?
- Berapa TTL yang reasonable?

<details>
<summary>💡 Solusi</summary>

- **Strategi:** Cache-aside + scheduled precompute
- **TTL:** 24 jam (sampai laporan bulan depan di-generate)
- **Invalidation:** Manual via cron job setelah batch process, atau versioning dengan month:year (`report:2025-08:omzet`)
- **Precompute:** Generate dan populate cache via background job, bukan on-demand
</details>

---

## Common Mistakes

| Mistake | Correct |
|---------|---------|
| Cache as single point of failure | Cache as optimization layer |
| No invalidation after write | Invalidate on mutation |
| TTL never reviewed | TTL based on data volatility |
| Global keys (no tenant) | Tenant-scoped keys |
| No monitoring | Track hit/miss/error rates |

---

## Running the Lab

### Prerequisites

- Go 1.22+
- Docker / Docker Compose (untuk Redis integration tests)

Untuk integration tests, pastikan Redis berjalan:
```bash
docker compose up -d redis
```

### Run Unit Tests

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

### Expected Results (Run with `go run ./cmd/demo -scenario=...`)

```
Scenario: without-cache
Requests: 100
Repository Calls: 100
Cache Hits: 0

Scenario: cache-aside  
Requests: 100
Repository Calls: 1
Cache Hits: 99

Scenario: stampede-unprotected
Concurrent Requests: 100
Repository Calls: 100

Scenario: stampede-protected
Concurrent Requests: 100
Repository Calls: 1    ← Singleflight!

Scenario: write-through
Requests: 100
Repository Calls: 100 (update path)
Cache Hits after update: 100 (cache synced)
```

---

## Exercises

1. **Write-Through** — Implement `UpdateProduct` dengan Write-Through pattern (lihat `write_through.go`)
2. **Add Negative Caching** - Cache 404 untuk 5 detik
3. **Multi-Tenant Key Builder** - Create factory function for tenant-scoped keys (lihat `dashboard_key.go`)
4. **Time-Zone aware** - Implement `TodayInLocation()` untuk branch timezone support
5. **Benchmark** - Measure latency with/without cache under load
6. **Lock Contention** - Add distributed lock untuk critical update path
7. **Permission Invalidation** - Implement event-driven invalidation (publish/subscribe) untuk permission cache saat role berubah
8. **Multi-Tenant Dashboard** - Buat service yang meng-handle 10 tenant secara benar dengan key isolation

---

## Design Exercise — Dashboard Bengkel

Dashboard memiliki:

- Jumlah Invoice Hari Ini
- Pendapatan Hari Ini
- Top Mekanik
- Top Sparepart
- Jumlah Kendaraan Baru
- Jumlah Customer Aktif

Untuk setiap data, tentukan:

1. Perlu cache atau tidak?
2. Mengapa?
3. TTL berapa dan berdasarkan asumsi apa?
4. Cache key seperti apa?
5. Event apa yang menyebabkan invalidation?
6. Apa acceptable staleness?
7. Bagaimana behavior jika Redis down?
8. Bagaimana mencegah cache stampede?
9. Bagaimana memastikan tenant/cabang tidak saling membaca cache?
10. Metrics apa yang perlu dimonitor?

### Expected Reasoning / Rubric

| Data | Cache? | TTL | Key | Invalidation | Staleness | Redis Down | Stampede Mitigation | Tenant Isolation |
|------|--------|-----|-----|--------------|-----------|------------|---------------------|------------------|
| InvoiceCount | ✅ Ya | 30s | `tenant:{t}:branch:{b}:invoice:count` | New invoice/payment | 30s acceptable | Fallback DB | Singleflight + jitter | Tenant + branch |
| Revenue | ✅ Ya | 30s | `tenant:{t}:branch:{b}:revenue` | Payment record | 30s acceptable | Fallback DB | Singleflight + jitter | Tenant + branch |
| TopMekanik | ✅ Ya | 2min | `tenant:{t}:branch:{b}:top:mechanic` | New service record | 2min acceptable | Fallback DB | Singleflight + jitter | Tenant + branch |
| TopSparepart | ✅ Ya | 2min | `tenant:{t}:branch:{b}:top:part` | New parts usage | 2min acceptable | Fallback DB | Singleflight + jitter | Tenant + branch |
| VehicleCount | ✅ Ya | 15s | `tenant:{t}:branch:{b}:vehicle:count` | New vehicle created | 15s acceptable | Fallback DB | Singleflight + jitter | Tenant + branch |
| ActiveCustomer | ✅ Ya | 1min | `tenant:{t}:branch:{b}:customer:active` | New invoice (30 days) | 1min acceptable | Fallback DB | Singleflight + jitter | Tenant + branch |

### Key Patterns

- Multi-tenant: **Always include tenant_id in key**.
- Date-sensitive: **Include business date** (respect timezone).
- Real-time critical: **Don't cache** (status pembayaran, stok saat checkout).
- Display-only: **Cache with short TTL + DB validation for transactions**.

---

## Bagian 21 — Separation of Concerns: Caching vs Optimistic Locking

Lab ini **tidak mencampur** Caching dengan Optimistic Locking.

| Topik | Fokus Utama |
|-------|-------------|
| **Optimistic Locking (Lab 05)** | Concurrent writes, Lost update, Version column (`version = version + 1`), Compare-and-swap |
| **Caching (Lab 04)** | Redundant reads/computation, Cache Aside, TTL, Invalidation, Stampede |

Jangan pernah menggunakan Optimistic Locking sebagai mekanisme cache consistency. Optimistic locking adalah pattern database untuk mencegah race condition pada update, sedangkan caching adalah pattern read-side untuk mengurangi latency dan load.

---

## Bagian 22 — Security, PII, & Permission Caching

### Sensitive Data (PII & Credentials)
- **Jangan pernah cache** password, credential, token akses, atau PII sensitif (seperti nomor kartu kredit, nomor KTP) secara sembarangan di shared cache.
- Pastikan cache key mempertahankan **authorization boundary** sehingga User A tidak pernah membaca data milik User B.

### Permission Caching & Risks
Permission sering di-cache untuk mengurangi query ke DB authorization. Namun ini membawa **risiko besar**:

```
[Admin mencabut permission User A]
               ↓
[Database terupdate]
               ↓
[Cache permission User A MASIH ADA (TTL 10 menit)]
               ↓
[User A tetap bisa akses resource terlarang selama 10 menit!]
```

**Mitigasi:**
1. Gunakan **TTL sangat pendek** (misal 5-10 detik) untuk authorization data.
2. Lakukan **Active Invalidation** (hapus cache permission segera saat role/permission berubah).
3. Gunakan **Key Versioning** pada permission (misal `user:123:perms:v5`).
4. **Jika consistency requirement terlalu ketat** (misal sistem finansial core), jangan cache permission — query DB setiap request.

---

## Bagian 23 — Common Mistakes (Kesalahan Umum)

1. **Cache Everything**: Meng-cache semua hal tanpa analisis akses pattern.
2. **Cache Cheap Queries**: Meng-cache query `SELECT NOW()` atau key-value lookup yang sudah < 0.1ms di DB.
3. **No Invalidation Strategy**: Mengandalkan 100% pada TTL tanpa aktif invalidation saat data berubah.
4. **Remember Forever**: Menyimpan data tanpa TTL atau eviction policy, berisiko memory leak/OOM.
5. **Global Key on Multi-Tenant**: Menggunakan key generik (`products:list`) sehingga data antar tenant tercampur.
6. **TTL Terlalu Panjang / Pendek**: Terlalu panjang = stale data; terlalu pendek = cache miss rate tinggi.
7. **Synchronized Expiration**: Ribuan key expire di detik yang sama (stampede tanpa jitter).
8. **Cache = Source of Truth**: Menganggap Redis sebagai penyimpanan utama dan mengabaikan DB backup.
9. **Redis Failure = App Crash**: Tidak menerapkan fallback (graceful degradation) saat Redis down.
10. **Singleflight vs Distributed Lock**: Menganggap singleflight (`golang.org/x/sync/singleflight`) menyelesaikan stampede lintas instance (padahal hanya single-process).
11. **Unsafe Permission Caching**: Meng-cache permission tanpa invalidation, menyebabkan eskalsasi privilese setelah hak akses dicabut.
12. **Print/Console Log for Cache Errors**: Production system harus memakai structured logging yang benar, bukan `fmt.Printf` (hanya dipakai di lab ini untuk kesederhanaan).
13. **Write-Through Without Failure Handling**: Menganggap cache update selalu berhasil; seharusnya handle DB commit sukses tapi cache failure dengan invalidate sebagai fallback.
14. **Cache Session Data**: Menyimpan per-user session data di shared cache tanpa sesi isolation; seharusnya gunakan session store yang tepat.

---

## Rule of Thumb Sebelum Membuat Cache

Jika engineer belum bisa menjelaskan strategi invalidation, **jangan buru-buru menambahkan cache**. Cache tanpa invalidation = stale data yang menempel sampai TTL habis.

### Lima Pertanyaan Utama

1. **Berapa lama data ini boleh stale?**
   - Jika jawaban = "0 detik" → **jangan cache**. Database adalah satu-satunya answer yang valid.
   - Jika boleh 5–60s → **boleh cache dengan TTL pendek + invalidation.**
   - Jika boleh 5–60 menit → **cache is ideal** (dashboard, statistik, master data).

2. **Seberapa mahal mendapatkan atau menghitung data ini?**
   - Jika query < 1ms → **jangan cache**. Cache lookup (Redis round-trip) mungkin lebih lambat dari DB.
   - Jika query > 50ms (join, aggregation, network) → **serius pertimbangkan cache.**

3. **Apa dampaknya jika user melihat data 5, 30, atau 60 detik terlambat?**
   - Dashboard stats 30s stale? → Tidak masalah.
   - Saldo wallet 10s stale? → Financial loss / audit issue. **Jangan cache.**

4. **Bagaimana cache di-invalidasi saat source data berubah?**
   - Harus ada trigger yang eksplisit (after DB commit → invalidate).
   - Jika tidak punya trigger → gunakan TTL-only (risky) **atau jangan cache.**

5. **Apakah manfaat performanya lebih besar daripada complexity yang ditambahkan?**
   - Hitungan sederhana: Jika cache mengurangi 99% dari 1000 req/s DB query, manfaat jelas.
   - Jika hanya mengurangi 10%, complexity invalidation + failure handling mungkin tidak worth it.

### Senior-Level Pertanyaan Tambahan

- **Apa source of truth?** — Jika jawaban ambigu, cache akan menjadi source of truth yang salah.
- **Apa behavior saat Redis unavailable?** — Harus ada fallback ke DB (graceful degradation).
- **Apakah cache key aman untuk multi-tenant?** — Key harus meng-embed tenant ID. Collision = data leak.
- **Bagaimana observability hit/miss/error rate?** — Tanpa metrics, tidak bisa memvalidasi ROI cache. Hit ratio < 50% biasanya berarti cache tidak efektif.
- **Apa risiko stampede?** — Pastikan ada singleflight (intra-process) atau distributed lock (cross-instance).
- **Apa memory footprint-nya?** — Cache semua data = OOM. Pertimbangkan eviction policy + capacity planning.

---

## Senior Engineer Takeaways

1. **Cache is not free** — TTL + Eviction planning are required
2. **Hot keys exist** — Cache moves bottleneck from DB to Redis (capacity planning matters)
3. **Graceful degradation** — Cache down ≠ system down
4. **Invalidation hard** — Jika DB COMMIT berhasil tapi proses crash sebelum cache DELETE, *stale cache remains until TTL*. Oleh karena itu TTL = safety net. (Untuk reliability absolut, pertimbangkan Outbox/event-driven invalidation - advanced).
5. **Multi-tenancy** — Key design must isolate tenants from day one
6. **Timezone matters** — "Today" depends on business timezone
7. **Staleness trade-off** — Ask "How stale is acceptable?" not "Is data changed?"
8. **Redis ≠ Cache** — Redis adalah teknologi/datastore, cache adalah pola arsitektur. Redis juga bisa digunakan untuk session, distributed lock, queue, atau ephemeral data; tidak semua data di Redis otomatis bermakna "cache".
9. **Cache is a trade-off, not a default** — Measure read:write ratio, query cost, and staleness tolerance before implementing.
10. **Write-through ≠ always better** — Pilih berdasarkan consistency requirement. Cache-aside + invalidate adalah default yang aman untuk sebagian besar kasus.
11. **Rule of thumb** — Jika belum punya invalidation strategy, jangan tambahkan cache.

---

## Navigasi

- **Previous**: [Lab 03 — Distributed Transaction](../03-database-transaction/)
- **Next**: [Lab 05 — Race Condition](../05-race-condition/)
- **All Labs**: [Labs](../)

---

## Menggunakan Lab Ini

### Run Experiments

```bash
cd labs/04-caching

# Run tests
go test -v -count=1 ./...

# Run specific scenarios
go run ./cmd/demo -scenario=without-cache
go run ./cmd/demo -scenario=cache-aside
go run ./cmd/demo -scenario=stampede-unprotected
go run ./cmd/demo -scenario=stampede-protected
go run ./cmd/demo -scenario=write-through
```