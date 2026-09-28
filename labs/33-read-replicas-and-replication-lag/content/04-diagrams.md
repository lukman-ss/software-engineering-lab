## Diagram 1 — High-Level Architecture

```
┌──────────────────────────────────────┐
│           Client / App               │
└──────────────────┬───────────────────┘
                   │
                   ▼
┌──────────────────────────────────────┐
│         Dynamic Query Router          │
│                                      │
│  ┌─────────────────────────────────┐ │
│  │ Read/Write Split                │ │
│  │ - Write → always primary        │ │
│  │ - Read → router strategy below  │ │
│  └─────────────────────────────────┘ │
│  ┌─────────────────────────────────┐ │
│  │ Strategy Selection              │ │
│  │ 1. ReadNaive       → replica   │ │
│  │ 2. StickySession   → primary   │ │
│  │    (within TTL)                │ │
│  │ 3. Token (minLSN)  → replica   │ │
│  │    (wait + fallback primary)   │ │
│  │ 4. LagAware        → replica   │ │
│  │    (fallback primary if empty) │ │
│  └─────────────────────────────────┘ │
└──────────┬──────────────┬────────────┘
           │              │
     Writes /         Reads
     Primary Read     (lag-checked)
           │              │
           ▼              ▼
┌────────────────┐ ┌────────────────┐
│  Primary Node  │ │ Replica Pool   │
│  appliedLSN: N │ │ appliedLSN: ?  │
│  (WAL master)  │ │ (async / sync) │
└───────┬────────┘ └───────┬────────┘
        │                  ▲
        │  WAL Stream      │
        └──────────────────┘
            (with lag)
```

---

## Diagram 2 — Write Path: Async vs Sync

```
                    ┌─────────────────────┐
                    │  cluster.Write()     │
                    └─────────┬───────────┘
                              │
                    ┌─────────▼───────────┐
                    │  primary.data[key]  │
                    │  primary.LSN += 1   │
                    └─────────┬───────────┘
                              │
              ┌───────────────┼───────────────┐
              │                               │
     ┌────────▼────────┐            ┌─────────▼─────────┐
     │  ASYNC MODE     │            │  SYNC MODE        │
     │  (default)      │            │  (remote_apply)   │
     └────────┬────────┘            └─────────┬─────────┘
              │                               │
    Send WAL  │                 For each replica:
    to channel│                   sleep(lag)
    (non-blocking)              apply data + LSN
              │                 broadcast Cond
              │                               │
              ▼                               ▼
    ┌──────────────────┐           ┌──────────────────┐
    │ Worker goroutine │           │ Write blocks     │
    │ reads channel,   │           │ until ALL        │
    │ sleeps(lag),     │           │ replicas applied │
    │ applies data     │           │ return LSN       │
    └──────────────────┘           └──────────────────┘
```

---

## Diagram 3 — Read Path: Four Routing Strategies

```
                  Read(sessionID, key)
                         │
            ┌────────────┼────────────┐
            │            │            │            │
   ┌────────▼───┐ ┌──────▼─────┐ ┌───▼────────┐ ┌▼───────────┐
   │  Naive     │ │  Sticky    │ │  Token     │ │  LagAware  │
   │  (None)    │ │  Session   │ │  (minLSN)  │ │            │
   └────────┬───┘ └──────┬─────┘ └───┬────────┘ └┬───────────┘
            │            │            │            │
    round-robin   check lastWrite   find replica  filter by
    ALL replicas  time vs StickyDuration  appliedLSN≥minLSN  ΔLSN ≤ MaxLSNDiff
            │            │            │            │
            │     ┌──────┴──────┐    │       ┌────┴────┐
            │     │             │    │       │         │
            │  within TTL   after TTL    found    no replica
            │     │             │        │       eligible
            │  PRIMARY    fallback to   │          │
            │             ReadLagAware  │      PRIMARY
            │                          │     (fallback-lag)
            │                 ┌────────┘
            │                 │
            │          replica found?
            │         ┌──┴──┐
            │         │     │
            │        yes   no (timeout)
            │         │     │
            │      REPLICA  PRIMARY
            │              (fallback)
            │
     ┌──────▼──────┐
     │ REPLICA     │
     └─────────────┘
```

---

## Diagram 4 — Sticky Session Routing State Diagram

```
Session State
─────────────

[After Write]
     │
     ▼
┌─────────────────────────────────────┐
│  SessionState.LastWriteTime = now   │
│  SessionState.LastWriteLSN  = lsn  │
└───────────────┬─────────────────────┘
                │
                ▼
┌──────────────────────────────────┐     time.Since(lastWrite) >= StickyDuration
│ ReadWithStickySession(session)   │ ──────────────────────────────────────────────►
│   check sessions.Load(sessionID) │     │
│   if within StickyDuration:      │     ▼
│       → PRIMARY                  │   ReadLagAware(key) → replica with ΔLSN ≤ MaxLSNDiff
│   else:                          │   (fallback to primary if no eligible replica)
└──────────────────────────────────┘
```

---

## Diagram 5 — Causal Token (ReadWithToken) Flow

```
ReadWithToken(ctx, minLSN=47, key)
        │
        ▼
┌──────────────────────────────┐
│ Iterate replica list         │
│ Find replica.AppliedLSN ≥ 47 │
└──────┬───────────┬───────────┘
       │           │
    found       not found
       │           │
       ▼           ▼
  ┌─────────┐  ┌────────────────────────────┐
  │ REPLICA │  │ WaitForLSN(ctx, 47)        │
  │ (fresh) │  │   sync.Cond.Wait() loop    │
  └─────────┘  │   blocks until appliedLSN  │
               │   ≥ 47 or context timeout  │
               └──────┬───────────┬──────────┘
                      │           │
                   reached     timeout
                   targetLSN  (ctx.Err)
                      │           │
                      ▼           ▼
                ┌──────────┐  ┌──────────┐
                │ REPLICA  │  │ PRIMARY  │
                │ (fresh)  │  │(fallback)│
                └──────────┘  └──────────┘
```

---

## Diagram 6 — Test Matrix

```
┌───────────────────────────────────────────────────────────────────┐
│  Test                              │  Lag   │ Mode │ Assertion   │
├───────────────────────────────────────────────────────────────────┤
│  StaleRead                         │  500ms │async │ ErrNotFound │
│  StickySession                     │  500ms │async │ primary→replica│
│  ReadWithToken                     │  200ms │async │ wait ~200ms │
│  ReplicaLagThreshold_Fallback      │  10s   │async │ fallback    │
│  SynchronousReplication_Freshness  │  50ms  │sync  │ fresh       │
│  ConcurrentAccess                  │  10ms  │async │ race-free   │
│  WaitForLSN_ContextTimeout         │  10s   │async │ deadline    │
└───────────────────────────────────────────────────────────────────┘
```

Diagram-diagram di atas merupakan representasi langsung dari source code lab dan arsitektur seperti yang diimplementasikan. Komponen `replicaWorker`, `sync.Cond`, `walChannel`, dan `sessions sync.Map` tidak ditambahkan dari luar — semuanya ada dalam implementasi.
