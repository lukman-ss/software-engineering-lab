# Diagrams

## Diagram 1 — Lost Update Anomaly Timeline

```
Time →
T1:  [READ stock=100] → [CALC: 100-1=99] → [WRITE stock=99] → [COMMIT]
T2:  [READ stock=100] ————————→ [CALC: 100-1=99] → [WRITE stock=99] → [COMMIT]

Result: stock=99 ❌   Correct: stock=98 ❌
```

The lost update occurs because T2 reads the value before T1 commits, then both write 99. T1's change is silently overwritten.

---

## Diagram 2 — Pessimistic Locking (SELECT ... FOR UPDATE)

```
T1:  [BEGIN] → [SELECT ... FOR UPDATE] → [LOCK ACQUIRED] → [READ stock=100] → [WRITE stock=99] → [COMMIT] → [LOCK RELEASED]
T2:  [BEGIN] → [SELECT ... FOR UPDATE] ……… BLOCKS ………→ [LOCK AVAILABLE] → [READ] → [WRITE] → [COMMIT]

Only one thread can hold a row lock at a time.
```

---

## Diagram 3 — Optimistic Locking with Conflict Detection

```
T1:  [READ: stock=100, ver=5] → [CALC] → [UPDATE WHERE id=1 AND version=5] → [1 row affected] → [COMMIT] (ver becomes 6)
T2:  [READ: stock=100, ver=5] → [CALC] → [UPDATE WHERE id=1 AND version=5] → [0 rows affected] → [ABORT] → [RETRY] → success

Conflict detected when `affected_rows == 0`. T2 must retry or return an error.
```

---

## Diagram 4 — Atomic Update Flow

```
G1:  [LOCK] → [CHECK stock>=1] → [stock -= 1] → [UNLOCK]
G2:  […BLOCKED…LOCK…BLOCKED…] → [proceed after unlock]

Each goroutine is serialized by the lock. No separate read-modify-write window is exposed.
```

---

## Diagram 5 — Retry Convergence

```
G-A:  [READ ver=1] → [CALC] → [UPDATE ver=1] → [SUCCESS: ver=2]
G-B:  [READ ver=1] → [CALC] → [UPDATE ver=1] → [0 rows → ERR] → [SLEEP 2ms] → [READ ver=2] → [UPDATE ver=2] → [SUCCESS]
G-C:  [READ ver=1] → [CALC] → [UPDATE ver=1] → [0 rows → ERR] → [SLEEP 4ms] → [READ ver=2] … wait for A/B → [SUCCESS]
G-D:  [READ ver=2] → [CALC] → [UPDATE ver=2] → [SUCCESS]

Jittered backoff staggers retries, preventing thundering-herd.
```

---

## Diagram 6 — Store Architecture (Simulated Engine)

```
┌─────────────────────────────────────────────────────────┐
│                        Store                            │
├─────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │  mu (global) │  │  rowLocks    │  │   products   │  │
│  │  sync.Mutex  │  │  map[id]Mutex│  │  map[id]Prod │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  │
│                                                         │
│  Get(id)      → lock/unlock map briefly, copy product  │
│  GetRowLock(id) → get or create per-row mutex         │
│  Seed(id,name,s) → create product + row mutex          │
└─────────────────────────────────────────────────────────┘
```

The global `mu` protects map access. Each row has its own mutex in `rowLocks` for pessimistic-style serialization.

---

## Diagram 7 — Version Guard Pattern (Optimistic SQL)

```sql
-- Initial read (snapshot)
SELECT id, stock, version FROM products WHERE id = 1;
-- => id=1, stock=100, version=5

-- Update with guard
UPDATE products
SET stock = 99, version = 6
WHERE id = 1 AND version = 5;  -- key condition!

-- Interpretation:
-- 1 row affected → version was 5, now updated
-- 0 rows affected → version changed to something other than 5 → conflict
```

---

## Diagram 8 — Service Retry Loop

```
          ┌─────────────────────────┐
          │ DeductOptimisticWithRetry │
          │ (id, qty, maxRetries)    │
          └────────────┬──────────────┘
                       │
              ┌────────┴────────┐  YES
              │ attempt == 0?      ├───────┐
              │ (sleep 0-5ms)     │       │
              └────────┬──────────┘       │
                       │ NO               │
              ┌────────┴────────┐  SUCCESS?├──────► return nil
              │ OptimisticDeduct OK?├─────┘
              └────────┬────────┘
                       │ FAIL
              ┌────────┴────────┐
              │ err == OptimisticLock?│
              └────────┬────────┘
         YES │        │ NO
               │        │
        ┌──────┴────┐   │      ┌───────────┐
        │ attempt < max?│   │      │ return err│
        └──────┬────┘   │      │ (non-retry) │
               │ NO    │ YES  └───────────┘
       ┌───────┴─────┐ │
       │ return Err  │ │
       │ Optimistic  │ │
       │ Lock        │ │
       └─────────────┘ │
                       │
                       │ increment attempt
                       │ sleep (1<<attempt)ms + jitter
                       └───── LOOP ──────►
```

---

## Diagram 9 — Error Flow

```
DeductNaive(id, qty)
        │
 ┌──────┴──────┐
 │ GET(id) fail?│──YES──► ErrNotFound
 └──────┬──────┘
        │ NO
 ┌──────┴──────┐
 │ STOCK < qty?│──YES──► ErrInsufficientStock
 └──────┬──────┘
        │ NO
 ┌──────┴──────┐
 │ time.Sleep │ (widens race)
 └──────┬──────┘
        │
    [WRITE stale stock]

DeductOptimistic(id, qty)
        │
 ┌──────┴──────┐
 │ Version match?│──NO──► ErrOptimisticLock
 └──────┬──────┘
        │ YES
 ┌──────┴──────┐
 │ WRITE & inc │
 │ version     │
 └─────────────┘

DeductAtomic(id, qty)
        │
 ┌──────┴──────┐
 │ STOCK < qty?│──YES──► ErrInsufficientStock
 └──────┬──────┘
        │ NO
 ┌──────┴──────┐
 │ STOCK -= qty│
 └─────────────┘
```

---

## Diagram 10 — Counter Invariants

```
Store counters:
┌────────────────┬─────────────────────────────┐
│ Field           │ Purpose                     │
├────────────────┼─────────────────────────────┤
│ NaivelyDrawn    │ count of naive attempts     │
│ Pessimistically │ count of pessimistic success│
│ Optimistically  │ count of optimistic success │
│ OptimisticFails │ count of optimistic conflicts│
│ Atomically      │ count of atomic success     │
└────────────────┴─────────────────────────────┘

Invariant definition:
   FinalStock = InitialStock − TotalSuccessfulDeductions
   = Initial − (Pessimistically + Optimistically + Atomically)

For naive/optimistic without retry:
   FinalStock = InitialStock − Optimistically
   AND
   OptimisticFails > 0 (conflicts existed)
```