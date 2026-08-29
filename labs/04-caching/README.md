# Lab 04 — Caching: Mengurangi Latency, Tapi Apa Biayanya?

> **Mental Model**: Caching adalah trade-off antara **latency** (cepat) dan **consistency** (benar). Senior engineer memilih teknik caching yang tepat untuk workload nya.

---

## Learning Objectives

Setelah menyelesaikan lab ini, Anda akan memahami:

- Mengapa caching diperlukan & kapan tidak boleh
- Cache Aside pattern & implementasinya
- Cara mengatasi *Cache Stampede* dengan Single Flight
- TTL strategy & jitter untuk stampede mitigation
- Cache Key design untuk multi-tenant & versioning
- Graceful degradation saat cache down
- Source of Truth: Database vs Cache role
- Heat map bottleneck: Database → Redis tradeoff
- Time-zone aware caching untuk business day

---

## Problem

Dashboard workshop menampilkan statistik dilihat oleh banyak user. Tanpa cache, setiap request menghasilkan *heavy aggregation* ke database:

**Traffic Simulation:**

```
[500 concurrent users] → Dashboard Request
                              ↓
                     Database: 6 queries + join/aggregation
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

## Cache Aside (Hit & Miss)

### Cache Hit Path

```
[Client]
    │
    ▼
[Application]
    │
    ▼
[Redis GET]
    │ HIT
    ▼
[Return Cached JSON]
    │
    ▼
[Client]
```

### Cache Miss Path

```
[Client]
    │
    ▼
[Application]
    │
    ▼
[Redis GET]
    │ MISS
    ▼
[Database Query] → Compute
    │
    ▼
[Redis SET (TTL)] ← Populate cache
    │
    ▼
[Return JSON]
```

---

## TTL & Cache Invalidation

### TTL Strategy (CMMS Dashboard)

| Field | TTL | Reason |
|-------|-----|--------|
| InvoiceCount | 30s | Near transaction time |
| Revenue | 30s | Same as invoice count |
| TopMechanic | 2min | Changes slowly |
| TopSparepart | 2min | Changes slowly |
| VehicleCount | 15s | Usually at day start |
| ActiveCustomer | 1min | Moderate changes |

### Invalidation Flow

```
[Aplikasi]
    │
    ├── BEGIN
    ├── UPDATE invoices SET status='paid' WHERE id=100
    ├── COMMIT ✅
    │
    └── INVALIDATE cache key
            │
            ├── (or TTL asafety net)
```

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
        ┌────────────────────┘
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
        ┌──────────┴──────────┐
        │  singleflight.Do()  │
        └─────────────────────┘
                    │
                    ▼
             1 DB Query Only
                    │
                    ▼
          Shared Result to ALL
```

---

## Single Flight

```go
var flight singleflight.Group{}

func GetData(ctx context.Context, key string) (Dashboard, error) {
    result, _, _ := flight.Do(key, func() (interface{}, error) {
        return fetchFromDB()
    })
    return result.(Dashboard), nil
}
```

---

## Distributed Lock (Redis SETNX)

```
App Instance A          Redis           App Instance B
      │                   │                   │
      ├── SET lock:1 token NX PX 10000 ──► Accepted
      │                   │                   │
      │            ┌──────┴──────┐            │
      │            │  Lock held  │            │
      │            └─────────────┘            │
      │                   │                   │
      └── GET lock:1 ─────┼───► Returns old token
                          │                   │
      ◄─── Wait ──────────┘                   │
```

---

## TTL Jitter

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

---

## Cache Failure (Graceful Degradation)

```
Client Request
      │
      ▼
   [Redis] ─X─ ERROR (Redis down)
      │
      ▼
[Fallback to Database] ← Source of Truth
      │
      ▼
  Request succeeds ✓ (slower)
```

---

## Metrics

| Counter | Meaning |
|---------|---------|
| `cache_hit` | Cache returned valid data |
| `cache_miss` | Cache empty |
| `cache_error` | Redis error/down |
| `database_query` | Query to PostgreSQL |
| `cache_rebuild` | Cache populated from DB |

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

### Run Unit Tests

```bash
cd labs/04-caching
go test -v -count=1 ./...
```

### Run Experiments

```bash
go run . -scenario=without-cache
go run . -scenario=cache-aside
go run . -scenario=stampede-unprotected
go run . -scenario=stampede-protected
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
```

---

## Exercises

1. **Implement Write-Through** - Add `PayInvoice(ctx, invoiceID)` that updates both DB and cache
2. **Add Negative Caching** - Cache 404 for 5 seconds
3. **Multi-Tenant Key Builder** - Create factory function for tenant-scoped keys
4. **Time-Zone aware** - Implement `TodayInLocation()` for branch timezone support
5. **Benchmark** - Measure latency with/without cache under load
6. **Lock Contention** - Add distributed lock for invoice updates

---

## Bagian 21 — Separation of Concerns: Caching vs Optimistic Locking

Lab ini **tidak mencampur** Caching dengan Optimistic Locking.

| Topik | Fokus Utama |
|-------|-------------|
| **Optimistic Locking (Lab 12)** | Concurrent writes, Lost update, Version column (`version = version + 1`), Compare-and-swap |
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
4. **Jangan cache** permission jika consistency requirement terlalu ketat (misal sistem finansial core).

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
11. **Unsafe Permission Caching**: Meng-cache permission tanpa invalidation, menyebabkan ekskalasi privilege setelah hak akses dicabut.

---

## Senior Engineer Takeaways

1. **Cache is not free** — TTL + Eviction planning are required
2. **Hot keys exist** — Cache moves bottleneck from DB to Redis (capacity planning matters)
3. **Graceful degradation** — Cache down ≠ system down
4. **Invalidation hard** — TTL as safety net, not correctness
5. **Multi-tenancy** — Key design must isolate tenants from day one
6. **Timezone matters** — "Today" depends on business timezone
7. **Staleness trade-off** — Ask "How stale is acceptable?" not "Is data changed?"
8. **Redis ≠ Cache** — Redis is the technology, cache is the architectural pattern

---

## Navigasi

- **Previous**: [Lab 03 — Distributed Transaction](../03-database-transaction/)
- **Next**: [Lab 05 — Pessimistic Locking](../05-pessimistic-locking/)
- **All Labs**: [](../)

---

## Menggunakan Lab Ini

### Run Experiments

```bash
cd labs/04-caching

# Run tests
go test -v -count=1 ./...

# Run specific scenarios
go run . -scenario=without-cache
go run . -scenario=cache-aside
go run . -scenario=stampede-unprotected
go run . -scenario=stampede-protected
```