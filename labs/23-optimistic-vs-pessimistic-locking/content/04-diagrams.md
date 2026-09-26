# Diagrams

## Diagram 1 — Lost Update Anomaly Flow

```
Time →
Thread-1:  [READ stock=100] → [CALC: 100-1=99] → [WRITE stock=99] → [COMMIT]
Thread-2:  [READ stock=100] ————————→ [CALC: 100-1=99] → [WRITE stock=99] → [COMMIT]

Result: stock=99, bukan 98
```

---

## Diagram 2 — Pessimistic Locking (SELECT ... FOR UPDATE)

```
Thread-1:  [BEGIN] → [SELECT FOR UPDATE] → [LOCK ACQUIRED] → [READ stock=100] → [WRITE stock=99] → [COMMIT] → [LOCK RELEASED]
Thread-2:  [BEGIN] → [SELECT FOR UPDATE] … BLOCKS … → [LOCK AVAILABLE] → [READ] → [WRITE] → …

Hanya satu thread yang dapat hold lock pada satu waktu.
```

---

## Diagram 3 — Optimistic Locking with Conflict

```
Thread-1:  [READ: stock=100, ver=5] → [CALC] → [UPDATE WHERE ver=5] → [COMMIT]
Thread-2:  [READ: stock=100, ver=5] → [CALC] → [UPDATE WHERE ver=5] → [FAIL: 0 rows] → [RETRY/409]

Thread-2 deteksi konflik lewat affected_rows = 0, lalu retry atau return error.
```

---

## Diagram 4 — Atomic Update Flow

```
Goroutine-1: [LOCK] → [CHECK stock>=1] → [stock -= 1] → [UNLOCK]
Goroutine-2: […BLOCKED…LOCK…BLOCKED…] → [proses setelah unlock]

Setiap goroutine serial secara implicit lewat lock. Tidak ada jendela read-modify-write terpisah.
```

---

## Diagram 5 — Retry Convergence

```
Goroutine A:  [READ ver=1] → [UPDATE ver=1] → SUCCESS (ver=2)
Goroutine B:  [READ ver=1] → [UPDATE ver=1] → CONFLICT → SLEEP(2ms) → [RETRY] → [READ ver=2] → [UPDATE ver=2] → SUCCESS

Goroutine C:  [READ ver=1] → [UPDATE ver=1] → CONFLICT → SLEEP(4ms) → [RETRY] → …

Semua goroutine akhirnya konvergen setelah beberapa percobaan.
```

---

## Diagram 6 — Store Architecture (Simulated DB Engine)

```
┌─────────────────────────────────────────────────────────┐
│                    Store                                 │
├─────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐ │
│  │  mu (global) │  │  rowLocks    │  │   products   │ │
│  │  sync.Mutex  │  │  map[id]Mutex│  │  map[id]Prod │ │
│  └──────────────┘  └──────────────┘  └──────────────┘ │
│                                                          │
│  Get(id)     → baca product (global lock pendek)        │
│  GetRowLock(id) → return mutex khusus baris            │
│  Seed()      → buat product + row lock                  │
└─────────────────────────────────────────────────────────┘
```

---

## Diagram 7 — Version Guard Pattern (Optimistic SQL)

```sql
-- Baca
SELECT id, stock, version FROM products WHERE id = 1;
-- => id=1, stock=100, version=5

-- Update dengan guard
UPDATE products 
SET stock = 99, version = 6
WHERE id = 1 AND version = 5;  -- key condition!

-- Jika satu baris terpengaruh → OK
-- Jika nol baris → ada yang lain mengubah version → konflik
```

---

## Diagram 8 — Service Layer Retry Loop

```
          ┌──────────────────┐
          │ DeductOptimistic │
          │ WithRetry(id,qty)│
          └────────┬─────────┘
                   │
         ┌─────────┴─────────┐
         │ attempt = 0       │
         │   sleep = 0ms     │
         └─────────┬─────────┘
                   │
         ┌─────────┴─────────┐  YES
         │ OptimisticDeduct OK?├───────► RETURN nil
         └─────────┬─────────┘
                   │ NO
         ┌─────────┴─────────┐
         │ err == OptimisticLock?│
         └─────────┬─────────┘
            NO   │      YES
         RETURN   │   ┌───┴──────┐
         err      │   │ attempt < max?│
                  │   └───┬──────┘
                      NO │      YES
               RETURN    │   ┌─────────────┐
               ErrLocking│   │ attempt++   │
                         │   │ backoff     │
                         └───┴─────────────┘
```

---

## Diagram 9 — Error Flow

```
            DeductNaive(id, qty)
                    │
           ┌────────┴────────┐
           │ GET(id) fails?  │──YES──► ErrNotFound
           └────────┬────────┘
                    │ NO
           ┌────────┴────────┐
           │ STOCK < qty?    │──YES──► ErrInsufficientStock
           └────────┬────────┘
                    │ NO
           ┌────────┴────────┐
           │ Write stock     │
           │ (stale value)   │
           └─────────────────┘
                    
            DeductOptimistic(id, qty)
                    │
           ┌────────┴────────┐
           │ Version match?  │──NO──► ErrOptimisticLock
           └────────┬────────┘
                    │ YES
           ┌────────┴────────┐
           │ Write & increment│
           │ version          │
           └──────────────────┘
```

---

## Diagram 10 — Counter Invariants

```
Track counter:
┌────────────────┬─────────────────────┐
│ Field          │ Purpose             │
├────────────────┼─────────────────────┤
│ NaivelyDrawn   │ counter naive       │
│ Pessimistically│ counter pessimistic │
│ Optimistically │ counter optimistic  │
│ OptimisticFails│ conflict counter    │
│ Atomically     │ counter atomic      │
└────────────────┴─────────────────────┘

Invarian:
InitialStock - TotalDrawn == FinalStock (atau accounting conflicts)
```