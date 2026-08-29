# Lab 04 — Caching: Optimalkan Latency, Tapi Apa Kostsnya?

> **Mental Model**: Caching adalah trade-off antara **latency** (cepat) dan **consistency** (benar). Senior engineer memilih teknik caching yang tepat untuk workadje nya.

---

## 1. Mengapa Caching Diperlukan?

```
Timeline:
T0: Client query inventory item
T1: App query DB (50ms)
T2: App return ke client
[Total latency: ~50ms]

Dengan Caching:
T0: Client query inventory item  
T1: App check cache (Redis, 0.5ms) → HIT
T2: App return ke client
[Total latency: ~1ms]
```

**Keuntungan**: Latency turun 50x. Namun...

---

## 2. Cache Aside Pattern (Standar)

### Flow

```
[Client] → [App] → [Cache] → (MISS) → [DB] → [Cache] → [App] → [Client]
                         → (HIT) →  [App] → [Client]
```

### Implementasi

```go
func GetProductWithCache(ctx context.Context, cache *RedisClient, db *sql.DB, id string) (Product, error) {
    // 1. Check cache dulu
    key := cacheKey("product", id)
    cached, err := cache.Get(ctx, key)
    if err == nil && cached != "" {
        var p Product
        if json.Unmarshal([]byte(cached), &p) == nil {
            return p, nil  // CACHE HIT
        }
    }
    
    // 2. Cache miss → query DB
    p, err := db.GetProduct(ctx, id)
    if err != nil {
        return Product{}, err
    }
    
    // 3. Populate cache dengan TTL
    data, _ := json.Marshal(p)
    cache.Set(ctx, key, string(data), 5*time.Minute)
    
    return p, nil
}
```

### Kapan Package?

- ✅ Data **aba-aba** (profil, konfigurasi)
- ✅ Data yang **jarang berubah**
- ✅ Query yang **mahal/dangerous** di DB

### Kapan Jangan Package?

- ❌ Data **transaksi** (payment amount)
- ❌ Data yang **kritis keakuratan**
- ❌ Data yang **sering berubah**

---

## 3. Cache Stampede (Thundering Herd)

### Masalah

```
Timeline:
T0: Cache key EXPIRES
T1: 1000 request masuk, semua CACHE MISS
T2: Semua query DB → overload database
T3: DB timeout → semua request gagal
```

### Solusi: Cache dengan Probabilistic Early Expiration

```go
func shouldRefresh(expiredAt time.Time) bool {
    age := time.Since(expiredAt)
    ttl := 5 * time.Minute
    
    // Jika sudah 80% TTL berlalu, refresh randomly
    if age > ttl*8/10 {
        return rand.Float64() < 0.5  // 50% chance
    }
    return false
}
```

---

## 4. Single Flight Pattern

### Ide

Jika banyak request untuk **same key**, jalankan **hanya satu** query DB. Result dialokasikan ke semua request.

```
Request A ──┐
Request B ──┼──► SingleFlight ──► DB ──► Cache
Request C ──┘
```

### Implementasi

```go
var flight = singleflight.Group{}

func GetWithSingleFlight(ctx context.Context, key string) (Product, error) {
    result, err, _ := flight.Do(key, func() (interface{}, error) {
        return db.GetProduct(ctx, key)
    })
    return result.(Product), err
}
```

### Keuntungan

- ⭐ Mengurangi load DB 90%+
- ⭐ Mengatasi stampede tanpa kompleksitas lock

---

## 5. Distributed Lock (Redis-based)

### Masalah

Single flight tidak aman bila **multiple app instances** (horizontal scaling). Butuh lock yang **distribusi**.

### Redlock Algorithm (Sederhana)

```go
const LOCK_KEY = "lock:item:" + id
const LOCK_VALUE = uuid.New().String()
const LOCK_TTL = 10 * time.Second

// SET key value NX PX ttl ( dengan random value untuk uniqueness )
func acquireLock(ctx context.Context, client *redis.Client, key string) (bool, string) {
    value := uuid.New().String()
    acquired := client.SetNX(ctx, key, value, LOCK_TTL).Val()
    return acquired, value
}

func releaseLock(ctx context.Context, client *redis.Client, key, value string) {
    // Lua script: Hanya delete jika value cocok (prevent releasing someone else's lock)
    script := redis.NewScript(`
        if redis.call("GET", KEYS[1]) == ARGV[1] then
            return redis.call("DEL", KEYS[1])
        else
            return 0
        end
    `)
    script.Run(ctx, client, key, value)
}
```

### Catatan Penting

- Lock harus memiliki **TTL** (auto-release bila process crash)
- Value harus **unique** (UUID) supaya hanya yang mengunci yang bisa lepaskan
- Jangan gunakan lock untuk query read-only

---

## 6. Cache Key Design

### Prinsip

1. **Prefix untuk namespace**: `product:123`, `user:456`
2. **Version untuk invalidation**: `product:123:v2`
3. **Tag-based**: `user:456:posts` → semua key yang berhubungan

### Contoh

```go
func cacheKey(entity, id string) string {
    return fmt.Sprintf("%s:%s:v2", entity, id)
}

// Untuk multi-level cache
func productDetailKey(id string) string {
    return fmt.Sprintf("product:%s:detail:v1", id)
}
```

---

## 7. Cache Invalidation Strategi

### 7.1 Time-based (TTL)

```go
redis.Set(ctx, key, value, 5*time.Minute)
```

- ✅ Simple
- ❌ Data bisa stale sampai TTL habis

### 7.2 Write-through

```go
db.Update(entity)
cache.Set(key, newValue, ttl)  // Update cache setelah DB success
```

- ✅ Consistent
- ❌ Latency write naik

### 7.3 Write-behind

```go
cache.Set(key, newValue, ttl)
asyncQueue.Enqueue(updateDB)  // DB update nanti
```

- ✅ Latency write rendah
- ❌ Risiko data hilang bila crash

---

## 8. Cache Metrics yang Harus Dipantau

| Metric | Tujuan |
|--------|--------|
| `cache_hit_rate` | Hitrate harus > 80% |
| `cache_miss_latency` | Bandingkan latency miss vs hit |
| `eviction_count` | Jika banyak eviction, TTL terlalu pendek |
| `cache_size_bytes` | Memory usage |
| `cached_item_ttl_seconds` | Distribusi TTL yang diterima |

---

## 9. Cache Invalidation (Bagian 8)

Cache invalidation adalah masalah nomor dua paling sulit dalam computer science (No. 1: naming, No. 3: off-by-one).

### Masalah

Cache menyimpan snapshot. Jika data berubah tapi cache tidak di-invalidate, client melihat data stale.

### Flow yang Benar: COMMIT → INVALIDATE

```
[Aplikasi]
    │
    ├── BEGIN TRANSACTION
    ├── UPDATE invoices SET status='paid' WHERE id=100
    ├── COMMIT  ✅ (DB committed)
    │
    └── INVALIDATE dashboard cache
            │
            └── SET cache_key "" (atau DEL cache_key)
```

### Mengapa TIDAK: INVALIDATE → COMMIT?

```
[Aplikasi]
    │
    ├── INVALIDATE cache  ❌ (cache kosong)
    │
    ├── BEGIN TRANSACTION
    ├── UPDATE invoices...
    ├── COMMIT (GAGAL! DB crash/timeout)
    │
    └── Hasil: Cache kosong + DB tidak berubah
              Client Baca → MISS → Query DB → Dapat OLD data (stale!)
```

**Race condition:** Jika delete cache sebelum commit, dan commit gagal, cache akan tetap kosong sampai TTL. Client yang read akan dapat stale data dari DB yang tidak jadi commit.

### Implementasi Invalidation

```go
// Setelah DB commit berhasil
func (s *DashboardService) CreateInvoice(ctx context.Context, branchID int64) error {
    // BEGIN, INSERT, COMMIT di sini
    tx, _ := s.db.BeginTx(ctx, nil)
    _, err := tx.ExecContext(ctx, "INSERT INTO invoices ...")
    if err != nil {
        tx.Rollback()
        return err
    }
    tx.Commit() // ✅ COMMIT success

    // INVALIDATE setelah commit
    err = s.InvalidateBranchDashboard(ctx, branchID)
    if err != nil {
        // Log error tapi jangan fail transaction
        // TTL adalah safety net jika invalidation gagal
        log.Println("cache invalidation failed:", err)
    }
    return nil
}
```

### Failure Window: COMMIT → Crash → Invalidate

```
T1: DB COMMIT success
T2: Process crash (sebelum invalidate)
T3: Client baca → cache HIT → data stale (lama)
T4: TTL expire (30s) → cache miss → query DB → data baru
```

**Mitigation:**
- TTL adalah safety net (30s dalam contoh ini)
- Untuk sistem kompleks: gunakan Outbox pattern (lihat Lab 03) untuk reliable invalidation

### Events yang Perlu Di-invalidate

| Event | Action |
|-------|--------|
| Invoice dibuat | `InvalidateBranchDashboard(branchID)` |
| Invoice dibayar | `InvalidateBranchDashboard(branchID)` |
| Service selesai | `InvalidateBranchDashboard(branchID)` |
| Sparepart digunakan | `InvalidateBranchDashboard(branchID)` |
| Customer dibuat/diubah | `InvalidateBranchDashboard(branchID)` |
| Vehicle dibuat | `InvalidateBranchDashboard(branchID)` |

### Test: Invalidation Membuktikan Cache Baru

```go
// Step 1: Request dashboard (cache miss)
d1, _ := svc.GetDashboard(ctx, branchID)
// Step 2: Mutate DB (invoice created)
svc.CreateInvoice(ctx, branchID, 5000.0)
// Step 3: Invalidate cache
svc.InvalidateBranchDashboard(ctx, branchID)
// Step 4: Request berikutnya (cache miss, fresh data)
d2, _ := svc.GetDashboard(ctx, branchID)
// d2 harus berbeda dari d1 (menunjukkan invalidation berhasil)
```

**Test harus gagal jika `InvalidateBranchDashboard()` dihapus:**
- Tanpa invalidation, Step 4 akan mengembalikan d1 (stale)
- With invalidation, Step 4 mengembalikan d2 (fresh)

---

## 10. Stale Data (Bagian 9)

Pertanyaan utama caching bukan:

> "Apakah datanya berubah?"

tetapi:

> "Berapa lama stale data masih dapat diterima?"

### Contoh

| Data | Stale 60 detik | Acceptable? |
|------|----------------|-------------|
| `TopMechanic` | Terlambat 1 menit | ✅ Biasanya OK |
| `InvoiceCountToday` | Terlambat 1 menit | ✅ Operasional OK |
| `SaldoWallet` | Terlambat 1 menit | ❌ Mungkin tidak OK |
| `StockRealTime` | Terlambat 1 menit | ❌ Race condition |

### Miskonsepsi: "Data sering berubah = tidak boleh cache"

Ini **SALAH**. Data real-time tetap dapat memiliki cache/read model tertentu.

**Contoh:** Stock gudang
- Tidak boleh cache hasil query langsung (race condition)
- Tetapi boleh cache: "Stock snapshot per 5 detik" jika business menerima 5s delay
- Atau: "Reserved stock" vs "Available stock" (dua model berbeda)

**Key insight:** Cache bukan about "data berubah atau tidak", tapi about **consistency requirement**.

### TTL Bukan Correctness Mechanism

TTL adalah safety net, bukan source of truth.

```
Tanpa invalidation aktif:
- Data bisa stale sampai TTL habis (max 30s dalam contoh)
- Jika DB crash setelah commit, cache masih serve stale data

Dengan invalidation aktif:
- Data fresh segera setelah mutation
- TTL hanya backup jika invalidation gagal
```

---

## 11. Redis Bukan Selalu Cache (Bagian 10)

### Miskonsepsi Umum

> "OTP tidak boleh disimpan di Redis karena bukan cacheable"

**SALAH.** Redis dapat digunakan untuk banyak hal:

| Use Case | Role Redis |
|----------|-----------|
| Cache dashboard | **Cache** (snapshot dari DB) |
| OTP verification | **Ephemeral state store** (TTL-bound) |
| Rate limiter | **Counter storage** (atomic INCR) |
| Distributed lock | **Lock storage** (SETNX) |
| User session | **Session storage** (TTL-bound) |
| Job queue | **Queue storage** (LIST/BRPOP) |

### Bedanya

- **Redis** = teknologi / data store
- **Cache** = architectural responsibility (mengurangi DB load)

OTP di Redis:
- ✅ Valid (Redis sebagai ephemeral state store)
- ❌ Bukan cache (tidak menyimpan snapshot dari DB)
- TTL = expiration (bukan safety net)

### Kapan Redis = Cache?

Redis adalah cache jika:
1. Data berasal dari database (snapshot)
2. Tujuannya mengurangi DB load
3. Stale data acceptable (TTL-based)

Redis bukan cache jika:
1. Data tidak ada di database (OTP, session)
2. Tujuannya state management (lock, rate limit)
3. Data harus hilang setelah TTL (ephemeral)

---

## 9. Trade-off: Staleness vs Performance

```
Tanpa Cache:
- Latency: 50ms
- Consistency: Strong (data selalu baru)

Dengan Cache (5min TTL):
- Latency: 1ms
- Consistency: Eventual (data mungkin 5 menit beda)

Pertanyaan: "Data berubah seberapa sering? Apakah 5 menit stale OK?"
```

---

## 10. Studi Kasus: CMMS/Workshop Dashboard

### Use Case

Dashboard menampilkan statistik workshop/bengkel:

- `InvoiceCountToday` — Jumlah invoice hari ini
- `TotalRevenueToday` — Pendapatan hari ini  
- `TopMechanic` — Mekanik paling banyak service
- `TopSparepart` — Sparepart terlaris
- `VehicleCountToday` — Kendaraan baru hari ini
- `ActiveCustomer` — Customer aktif 30 hari terakhir

### Masalah: Cache Miss on Every Request

```
[500 concurrent users request dashboard]
          ↓
      Cache Check (MISS)
          ↓
      Database Query (6 queries per request)
          ↓
   Aggregation + Join
          ↓
       Return to client
          ↓
      All 500 requests... same DB work!
```

**Baseline (Naive Service):**
- Setiap request dashboard = 6 query DB + aggregation
- 500 request = ~3000 query DB
- Latensi tinggi, beban DB meningkat

### Solusi: Cache Aside + Granular TTL

#### Cache Key Design

Key: `cmms:dashboard:v1:branch:{branchID}:date:{YYYY-MM-DD}`

**Alasan:**
- `cmms` = namespace
- `dashboard` = entity type
- `v1` = versioning (invalidasi saat format berubah)
- `branch:{branchID}` = multi-tenant support
- `date:{YYYY-MM-DD}` = daily granularity (reset otomatis tiap hari)

#### TTL Strategy

| Field | TTL | Alasan |
|-------|-----|--------|
| InvoiceCountToday | 30 detik | Dekat transaction time, butuh fresh |
| TotalRevenueToday | 30 detik | Mirip invoice count |
| TopMechanic | 2 menit | Perubahan lebih lambat |
| TopSparepart | 2 menit | Perubahan paling lambat |
| VehicleCountToday | 15 detik | Biasanya di awal hari |
| ActiveCustomer | 1 menit | Perubahan sedang |

**Trade-off:** Cache teraggregasi sebagai satu object dengan TTL 30 detik. TTL bukan correctness mechanism utama—ini safety net. Untuk invalidation aktif, panggil `InvalidateDashboardCache()` setelah:
- POST /invoices → invalidate
- POST /vehicles → invalidate

---

## Bagian 8 — Cache Stampede (Thundering Herd)

### Masalah

```
Timeline:
T0: Cache key EXPIRES
T1: 1000 request masuk, semua CACHE MISS
T2: Semua query DB → overload database
T3: DB timeout → semua request gagal
```

### Demonstrasi: Broken vs Protected

| Service | DB Rebuild Count (100 concurrent) |
|---------|-----------------------------------|
| `BrokenStampedeService` | ~100 |
| `ProtectedStampedeService` | 1 |

**Broken (Tanpa singleflight):**
```go
// Setiap goroutine langsung query DB
for i := 0; i < 100; i++ {
    cache.Get() // MISS
    db.Query()   // SEMUA query DB !
}
```

**Protected (Dengan singleflight):**
```go
// SingleFlight deduplicates concurrent DB queries
flight.Do(key, func() {
    return cache.Get() // HANYA SATU eksekusi
})
// 99 goroutine lain menunggu hasil yang sama
```

---

## Bagian 9 — Distributed Stampede & Single Flight Limitation

### Limitation Single Flight

Single flight hanya melindungi **satu process/application instance**.

```
App Instance A ──► SingleFlight ──► DB
App Instance B ──► SingleFlight ──► DB
App Instance C ──► SingleFlight ──► DB

Semua instance punya SingleFlight sendiri → Stampede masih terjadi!
```

### Solusi: Distributed Lock dengan Redis

Gunakan Redis `SET key unique-token NX PX ttl` untuk mutual exclusion lintas instance.

```go
// Acquire lock
token := uuid.New().String()
acquired := redis.SetNX(ctx, "lock:dashboard:branch:1", token, 10*time.Second)

// Release dengan atomic compare-and-delete (Lua script)
// if redis.call("GET", KEYS[1]) == ARGV[1] then
//     return redis.call("DEL", KEYS[1])
// else
//     return 0
// end
```

### Trade-off

- ✅ Single flight: Sangat efisien (single-process)
- ✅ Distributed lock: Melindungi multi-instance
- ❌ Distributed lock: Network round-trip + Lua script overhead

---

## Bagian 10 — TTL Jitter

### Masalah: Synchronized Expiration

```
100 cache key dibuat bersamaan
TTL semuanya 60 detik
↓
60 detik kemudian semuanya expired bersamaan
↓
STAMPEDE! Semua request miss cache
```

### Solusi: baseTTL + random jitter

```go
// 60s + random(0..15s) → TTL antara 60-75 detik
func TTLWithJitter(base, maxJitter time.Duration) time.Duration {
    jitter := time.Duration(rand.Int63n(int64(maxJitter)))
    return base + jitter
}
```

**Result:** Expiration tersebar, stampede mitigasi.

---

## Bagian 11 — Negative Caching

### Scenario

```
GET /product/999999
Product tidak ada di database
```

### Tanpa Negative Cache

- Request bot berulang
- Setiap request = database lookup
- DB tertekan oleh traffic pencarian tidak-valid

### Dengan Negative Cache

```go
cache.Set("product:999999", "NULL_NOT_FOUND", 30*time.Second)
```

- Subsequent request → cache hit (not-found)
- DB protection untuk key tidak ada

### Trade-off

Jika object baru dibuat selama negative TTL:
- User masih dapat `404` sampai TTL habis
- Business decision: TTL pendek (5-10 detik) untuk balance

---

## Bagian 12 — Cache Key Design

### Format

```
namespace:entity:vVersion:tenant:identifier:date
```

Contoh: `cmms:dashboard:v1:branch:12:2026-08-28`

### Komponen

| Komponen | Tujuan |
|----------|--------|
| `cmms` | Namespace aplikasi |
| `dashboard` | Entity/resource type |
| `v1` | Version (invalidation saat schema berubah) |
| `branch:12` | Tenant/owner isolation |
| `2026-08-28` | Date/window granularity |

### Mengapa Version Berguna?

```go
// v1 → v2 (schema change)
"product:123:v1" → "product:123:v2"

// Tanpa version: harus bulk delete semua key
// Dengan version: cuma buat key baru, lama auto-expire
```

### Multi-Tenant Isolation

```
Tenant A: cmms:dashboard:v1:branch:1:2026-08-28
Tenant B: cmms:dashboard:v1:branch:2:2026-08-28
```

Setiap tenant punya key terpisah, tidak ada collision/cache poisoning.

### Evolvable Key Patterns

| Pattern | Use Case |
|---------|----------|
| `resource:id` | Single resource |
| `resource:id:field` | Nested field |
| `resource:list:tenant:status` | Filtered list |
| `resource:search:query:hash` | Query result |

---

## Menjalankan Semua Test

```bash
cd labs/04-caching && go test -v -count=1 ./...
```

| Test | Bukti |
|------|-------|
| `TestNaiveNoCache` | Latency tidak ada cache |
| `TestCacheAsideHit` | Cache hit mengurangi latency |
| `TestCacheStaleRead` | Data cache bisa stal |
| `TestSingleFlightConcurrentRequests` | Single flight menghindari duplicate query |
| `TestCacheStampedeMitigation` | Probabilistic refresh mengurangi stampede |
| `TestDistributedLockMutualExclusion` | Lock only satu yang dapat lock |
| `TestCacheKeyIncludesVersion` | Version di key untuk invalidation |
| `TestStampedeBrokenVersion` | Broken service = stampede |
| `TestStampedeProtectedVersion` | Protected service = single query |
| `TestTTLWithJitter` | TTL jitter mengurangi synchronized expirations |
| `TestNegativeCache` | Negative cache mengurangi DB load untuk 404s |
| `TestDashboardCacheInvalidation` | Invalidation setelah mutation |