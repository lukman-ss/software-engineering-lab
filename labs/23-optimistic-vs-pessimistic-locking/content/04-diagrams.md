# Diagrams

## Diagram 1 — Concurrency Anomaly Comparison

```text
┌─────────────────────┐    ┌──────────────────────┐    ┌─────────────────────┐
│        NAIVE        │    │     PESSIMISTIC      │    │     OPTIMISTIC      │
│   Read-Modify-Write │    │     (SELECT FOR      │    │   (Version Guard)   │
└─────────────────────┘    │    UPDATE)           │    └─────────────────────┘
          │                 └──────────────────────┘             │
          ▼                          │                          ▼
┌─────────────────────┐              ▼              ┌─────────────────────────┐
│  Tx A: READ stock=  │    ┌──────────────────────┐   │ Tx A: READ stock,ver=5  │
│       100           │    │ Tx A: LOCK row       │   │                         │
└─────────────────────┘    │ Tx A: READ stock=100 │   └─────────────────────────┘
          │                 └──────────────────────┘              │
          │                          │                          │
          ▼                          ▼                          ▼
┌─────────────────────┐    ┌──────────────────────┐   ┌─────────────────────────┐
│  Tx B: READ stock=  │    │ Tx B: WAIT (blocked) │   │ Tx B: READ stock,ver=5  │
│       100           │    └──────────────────────┘   └─────────────────────────┘
└─────────────────────┘                          │                          │
          │                                      │                          │
          ▼                                      ▼                          ▼
┌─────────────────────┐    ┌──────────────────────┐   ┌─────────────────────────┐
│  Tx A: STOCK=99    │    │ Tx A: WRITE stock=99 │   │ Tx A: WRITE IF ver=5    │
│  Tx A: WRITE        │    │ Tx A: UNLOCK         │   │  -> SUCCESS (rows=1)    │
│       (stale!)      │    └──────────────────────┘   └─────────────────────────┘
└─────────────────────┘                          │                          │
          │                                      │                          │
          ▼                                      ▼                          ▼
┌─────────────────────┐    ┌──────────────────────┐   ┌─────────────────────────┐
│  Tx B: STOCK=99    │    │ Tx B: LOCK row       │   │ Tx B: WRITE IF ver=5    │
│  Tx B: WRITE        │    │ Tx B: WRITE stock=99 │   │  -> FAIL (rows=0)       │
│  (overwrites A!)    │    │ Tx B: UNLOCK         │   │  -> ErrOptimisticLock   │
└─────────────────────┘    └──────────────────────┘   └─────────────────────────┘
          │                                      │                          │
          ▼                                      ▼                          ▼
┌─────────────────────┐    ┌──────────────────────┐   ┌─────────────────────────┐
│ Final Stock: 99     │    │ Final Stock: 99      │   │ Final Stock: 99         │
│ (LOST UPDATE!)      │    │ (CORRECT)            │   │ (CORRECT, conflict      │
│ Tx A write LOST     │    │ Tx A + B sequential  │   │  detected, retry or 409)│
└─────────────────────┘    └──────────────────────┘   └─────────────────────────┘


┌─────────────────────┐    ┌──────────────────────┐
│      ATOMIC         │    │                    │
│   (Single Statement)│    │                    │
└─────────────────────┘    └──────────────────────┘
          │                         │
          ▼                         ▼
┌─────────────────────┐    ┌──────────────────────┐
│ Tx A: UPDATE        │    │ Tx B: UPDATE         │
│ stock = stock - 1    │    │ stock = stock - 1     │
│ WHERE stock >= 1      │    │ WHERE stock >= 1     │
│ -> SUCCESS (rows=1)   │    │ -> BLOCKED (row       │
└─────────────────────┘    │   locked by Tx A)    │
          │                 └──────────────────────┘
          ▼                          │
┌─────────────────────┐              ▼
│ Tx A COMMIT         │    ┌──────────────────────┐
│ -> Stock = 99      │    │ Tx B: NOW PROCEEDS    │
└─────────────────────┘    │ -> SUCCESS (rows=1)   │
                           │ -> Stock = 98         │
                           └──────────────────────┘
```

## Diagram 2 — Architecture (Code Structure)

```text
┌──────────────────────────────────────────────────────────────────┐
│                          cmd/demo/main.go                         │
│  CLI runner: 5 scenarios side-by-side                              │
│  (50 goroutines each, except optimistic=20)                        │
└──────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌──────────────────────────────────────────────────────────────────┐
│                     internal/inventory/                            │
│                                                                  │
│  ┌──────────────┐         ┌──────────────┐                       │
│  │ model.go     │         │ service.go   │                       │
│  │ ─────────    │◄────────│ ─────────    │                       │
│  │ Product      │         │ Service      │                       │
│  │ struct       │         │ ─ DeductNaive│                       │
│  │ Errors       │         │ ─ DeductPess│                       │
│  └──────────────┘         │ ─ DeductOpt.  │                       │
│                            │ ─ DeductOpt.  │                       │
│                            │    WithRetry  │                       │
│                            │ ─ DeductAtom  │                       │
│                            └──────────────┘                       │
│                                    │                               │
│                                    ▼                               │
│  ┌──────────────┐         ┌──────────────┐                       │
│  │ store.go     │◄────────│ tests/        │                       │
│  │ ─────────    │         │ ─────────    │                       │
│  │ Store       │         │ locking_   │                       │
│  │ struct      │         │ test.go   │                       │
│  │ ── Seed()    │         │ TestNaive-  │                       │
│  │ ── Get()     │         │ LostUpdate    │                       │
│  │ ── NaiveD-   │         │ TestPess-   │                       │
│  │    lect()     │         │ Locking       │                       │
│  │ ── Pess-     │         │ TestOpt-    │                       │
│  │    imisticDeduct()│    │ Locking...  │                       │
│  │ ── Optim-    │         │ TestAtom-   │                       │
│  │    isticDeduct()│    │ ic...        │                       │
│  │ ── AtomicDed-│         │             │                       │
│  │    uct()    │         │             │                       │
│  │ rowLocks    │         │             │                       │
│  │ products    │         │             │                       │
│  │ atomic coun-│         │             │                       │
│  │ ters        │         │             │                       │
│  └──────────────┘         └──────────────┘                       │
└──────────────────────────────────────────────────────────────────┘
```

## Diagram 3 — Optimistic Locking Retry Flow

```text

┌─────────────────────────────────────────────────────────────┐
│                    Optimistic With Retry                     │
│                   (Service layer)                            │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
               ┌────────────────────┐
               │  Attempt 0         │
               │  OptimisticDeduct()  │
               └────────┬───────────┘
                        │
               ┌────────▼────────────┐
               │  Conflict?          │
               │  (ErrOptimisticLock)│
               └────┬─────────┬──────┘
              No   │         │  Yes
                   │         │
              ┌───▼───┐  ┌───▼────────┐
              │Return │  │ Sleep:     │
              │ nil   │  │ (1<<0)+ms +│
              │(succes│  │ jitter 0-5│
              │ s)    │  │ ms        │
              └───────┘  │ = 1-5ms  │
                         └────┬──────┘
                              │
                     ┌────────▼────────┐
                     │  Attempt 1       │
                     │  Optimistic-     │
                     │  Deduct()        │
                     └────────┬────────┘
                              │
                     ┌────────▼────────┐
                     │  Conflict?      │
                     └────┬─────────┬──┘
                No        │         │  Yes
                       ┌──▼──┐    ┌─▼──────────┐
                       │Retur│    │ Sleep:     │
                       │n nil│    │ (1<<1)+ms +│
                       │(suc-│    │ jitter 0-5│
                       │cess)│    │ ms        │
                       └─────┘    │ = 2-7ms  │
                                  └────┬──────┘
                                       │
                              ┌────────▼────────┐
                              │  Attempt 2       │
                              │  Optimistic-     │
                              │  Deduct()        │
                              └────────┬────────┘
                                       │
                              ┌────────▼────────┐
                              │  ...continue...  │
                              │  until maxRetries │
                              │  or success      │
                              └─────────────────┘
```

## Diagram 4 — In-Memory Store Internal State

```text
Store struct
┌─────────────────────────────────┐
│        products (map[int]*Product)   │
├─────────────────────────────────┤
│  key=1: {Stock: 50, Version: 6}  │
│  key=2: {Stock: 50, Version: 51} │
│  key=3: {Stock: 99, Version: 2}  │
│  key=4: {Stock: 80, Version: 21} │
│  key=5: {Stock: 50, Version: 51} │
└─────────────────────────────────┘

┌─────────────────────────────────┐
│        rowLocks (map[int]*Mutex)  │
├─────────────────────────────────┤
│  key=1: [unlocked]             │
│  key=2: [locked: goroutine 42] │
│  key=3: [unlocked]             │
│  key=4: [locked by retry]      │
│  key=5: [unlocked]             │
└─────────────────────────────────┘

Atomic Counters (updated via atomic.AddInt64):
┌──────────────────────────┐
│  NaivelyDrawn:   50      │
│  Pessimistically: 50      │
│  Optimistically: 21      │
│  OptimisticFails: 42      │
│  Atomically:     50      │
└──────────────────────────┘

Note: "Atomic" counters track total qty applied per strategy.
"OptimisticFails" tracks conflict count (not qty).
```

## Diagram 5 — Demo Scenario Results

```text
┌═══════════════════════════════════════════════════════════╗
║                    DEMO SCENARIO RESULTS                   ║
╠═══════════════════════════════════════════════════════════════
       Initial    Concurrent   Final    Status       Notes
       Stock      Requests     Stock

 [1]  100          50           99     FAIL         Lost Update!
      Naive RMW

 [2]  100          50           50     PASS         Fully Synced
      Pessimistic
      (SELECT FOR UPDATE)

 [3]  100          20           99     PASS*       State Guarded
      Optimistic               (1 ok,    Zero Corruption
      (no retry)                19 confl)

 [4]  100          20           80     PASS         All Converged
      Optimistic     (61 retries via
      +Retry          backoff)

 [5]  100          50           50     PASS         Lockless
      Atomic       Single Statement

 [*] = No data corruption, but conflicts detected
       Application must retry or return 409
```

## Diagram 6 — Selection Decision Matrix

```text
┌─────────────────────┬─────────────────────┬─────────────────────┐
│  Contention Level   │   Recommended        │  Reasoning          │
├─────────────────────┼─────────────────────┼─────────────────────┤
│                     │                     │                     │
│  Low Contention     │  Optimistic Locking  │  Few conflicts =    │
│  (many reads,       │  + Retry             │  few retries =      │
│  few writes)        │                     │  high throughput    │
│                     │                     │                     │
├─────────────────────┼─────────────────────┼─────────────────────┤
│                     │                     │                     │
│  High Contention    │  Pessimistic Lock-  │  Prevent conflicts  │
│  (many concurrent    │  ing                │  by blocking =      │
│  writes)            │                     │  no retry needed    │
│                     │                     │                     │
├─────────────────────┼─────────────────────┼─────────────────────┤
│                     │                     │                     │
│  Simple Counter     │  Atomic Single-     │  No read-modify-    │
│  Decrement          │  Statement          │  write window       │
│  (stock, quota)     │                     │  possible = lock-  │
│                     │                     │  less               │
│                     │                     │                     │
├─────────────────────┼─────────────────────┼─────────────────────┤
│                     │                     │                     │
│  Cross-DB /        │  Distributed Lock   │  Database-native    │
│  Microservices      │  (Redis)            │  lock tidak bisa   │
│  Resource           │                     │  span multiple DB   │
│                     │                     │                     │
└─────────────────────┴─────────────────────┴─────────────────────┘
```