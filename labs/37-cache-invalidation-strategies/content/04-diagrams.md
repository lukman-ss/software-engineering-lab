# Diagrams — Cache Invalidation Strategies

Diagrams berikut diturunkan dari struktur implementasi lab (`internal/cache/*`). Tidak ada komponen tambahan yang diciptakan di luar yang ada di source.

---

## Diagram 1 — Arsitektur Package

```
cmd/demo/main.go
        │
        ▼
internal/cache/
 ├── store.go       MemoryCache + TTLWithJitter
 ├── repo.go        MockDB (query/write counters)
 ├── patterns.go    CacheAsideService, WriteThroughService, WriteBehindService
 └── stampede.go    NaiveStampedeService, SingleFlightService,
                     XFetchService, SWRService, ShouldRecompute
        │
        ▼
tests/cache_test.go
```

---

## Diagram 2 — Cache-Aside Flow

```
Reader
  │
  ├─ GET ─► MemoryCache.Get(key)
  │            │
  │            ├─ HIT ─► return value
  │            │
  │            └─ MISS ─► MockDB.Query(key)
  │                            │
  │                            ├─ found ─► MemoryCache.Set(key, val, ttl, delta)
  │                            │              │
  │                            │              └─ return val
  │                            └─ not found ─► error
  │
  └─ UPDATE(key, val) ─► MockDB.Write(key, val)
                              │
                              └─ MemoryCache.Delete(key)
```

---

## Diagram 3 — Write-Through vs Write-Behind

```
WRITE(key, val)
  │
  ├── Write-Through
  │     ├─ MockDB.Write(key, val)
  │     └─ MemoryCache.Set(key, val, ttl, delta)
  │
  └── Write-Behind
        ├─ MemoryCache.Set(key, val, ttl, ...)
        └─ select { case writeQueue <- {key,val}:
                  default: drop (overflow) }
              │
              ▼
        flushWorker (goroutine)
              │
              └─ for req in writeQueue: MockDB.Write(req.Key, req.Value)
```

---

## Diagram 4 — Stampede Mitigation Paths

```
GET(key) — cache miss
  │
  ├── NaiveStampedeService
  │     └─ MockDB.Query(key)   (each goroutine does its own query)
  │
  ├── SingleFlightService
  │     └─ singleflight.Group.Do(key, fn)
  │           │
  │           ├─ first caller ─► fn runs once:
  │           │                   MockDB.Query(key)
  │           │                   MemoryCache.Set(...)
  │           │
  │           └─ waiting callers ─► receive same result (shared=true)
  │
  ├── XFetchService
  │     ├─ GetRaw(key)
  │     ├─ compute remaining
  │     └─ ShouldRecompute(delta, beta, remaining, u)
  │           ├─ true  ─► MockDB.Query(key) + MemoryCache.Set
  │           └─ false ─► return stale item immediately
  │
  └── SWRService
        ├─ GetRaw(key)
        ├─ if expired && now < ExpiresAt+staleDelta
        │     └─ triggerRevalidate(key) → goroutine → MockDB.Query
        │     return item.Value (stale)
        └─ else → MockDB.Query(key) + MemoryCache.Set
```

---

## Diagram 5 — XFetch Condition Geometry

```
condition: -Δ · β · ln(U) > TTL_remaining

U ∈ (0,1)  →  ln(U) < 0  →  -ln(U) > 0
              │
              ▼
offset = Δ·β·(-ln(U))   (positive, exponential tail)
              │
              ├── small U (near 0)  → large offset → early refresh
              │
              └── large U (near 1)  → small offset → likely no refresh
```

---

## Diagram 6 — SWR Timeline

```
Time ─────────────────────────────────────────────►
       │←── TTL ──→│←── stale-while-revalidate ──→│
       │            │                             │
       │  fresh     │  stale served               │  hard miss
       │  (cache)   │  + async revalidation        │  (refetch)
       │            │                             │
       ├────────────┼─────────────────────────────┤
                    │◄───── revalidation goroutine ───►
                         (single per key, deduped)
```

---

## Diagram 7 — Jitter Distribution

```
Base TTL = 5m, MaxJitter = 30s
Result range: [5m0s, 5m30s)

Sample outputs (from demo):
  5m13.339s
  5m26.959s
  5m18.809s
  5m28.543s
  5m12.312s

Effect: keys do not expire at the same instant
Limitation: single hot key still vulnerable to concurrent misses
```
