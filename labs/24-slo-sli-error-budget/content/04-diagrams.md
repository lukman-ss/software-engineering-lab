# Diagrams

Semua diagram berdasar implementasi asli lab. Tidak ada komponen yang tidak ada di kode.

---

## D1 — High-Level Architecture

```
┌────────────────────────────────────────────────────────────┐
│                         CMD/DEMO                           │
│   (Simulates traffic, incidents, checks results)           │
└──────────┬──────────────────┬────────────────┬───────────┘
           │                  │                │
           ▼                  ▼                ▼
  ┌─────────────┐   ┌──────────────┐  ┌──────────────┐
  │ SLO EVALUATOR │   │ SHORT WINDOW │  │ LONG WINDOW  │
  │   (slo/)      │   │ (5 min)      │  │ (60 min)     │
  └──────┬──────┘   └──────┬───────┘  └──────┬───────┘
         │ GetSummary()  │                 │
         │          └────┘                 │
         ▼              ▼                  ▼
  ┌─────────────┐  ┌──────────────┐  ┌──────────────┐
  │ WINDOW-     │  │ SLIDING-     │  │ SLIDING-     │
  │ TRACKER     │  │ BUCKETED     │  │ BUCKETED     │
  │   (metrics) │  │ TRACKER      │  │ TRACKER      │
  └──────┬──────┘  └──────┬───────┘  └──────┬───────┘
         │                   │               │
         └───────────────────┼───────────────┘
                             │ Shared isGood predicate
                             ▼
                    ┌──────────────────┐
                    │  ALERT ENGINE    │
                    │ (alerting/)      │
                    │ multi-window     │
                    │ burn-rate check  │
                    └──────────────────┘
                             │
                             ▼
                     ┌──────────────┐
                     │ Alert Result │
                     │ (Page/Ticket)│
                     └──────────────┘
```

Komponen aktual:
- `internal/metrics/tracker.go`: `WindowTracker` yang digunakan untuk SLO evaluator, short tracker, dan long tracker.
- `internal/slo/evaluator.go`: `Evaluator` menggunakan satu `WindowTracker` untuk SLI & error budget.
- `internal/alerting/engine.go`: `AlertEngine` menggunakan dua `WindowTracker` (short & long).

---

## D2 — WindowTracker Bucket Lifecycle

```
Event arrives at time T
       │
       ▼
┌──────────────────────────────────────┐
│ evictStaleLocked(T)                 ← │ Removes buckets older than (T - windowSize)
└──────────────────────────────────────┘
       │
       ▼
┌──────────────────────────────────────┐
│ bucketStart = T.Truncate(bucketSize) │ (mis. bucketSize = 10s → T = 12:34:25 → 12:34:20)
└──────────────────────────────────────┘
       │
       ▼
┌──────────────────────────────────────┐
│ isGood = isGoodEvent(e)              │ (StatusCode < 500 && Duration ≤ threshold)
└──────────────────────────────────────┘
       │
       ▼
┌──────────────────────────────────────┐
│ Buckets exist?                      │
├────────────── NO ────────────────────┤
│ → Append new bucket: {StartTime: bucketStart, Total: 1, Good/Bad: 1/0}
├────────────── YES ───────────────────┤
│ Last bucket.StartTime == bucketStart?
├────────────── YES ───────────────────┤
│ → Increment last bucket (fast path)
├────────────── NO ────────────────────┤
│ bucketStart BEFORE last bucket.StartTime?
├────── YES ───────────────────────────┤
│ → Search sorted slice:
│   - Found matching bucket start → increment
│   - Found gap (bucket should be inserted) → insert new bucket in order
├────── NO (future timestamp) ─────────┤
│ → Append as new bucket (will be evicted when stale)
└──────────────────────────────────────┘
```

Dipakai untuk test `TestOutOfOrderTimestamps`: event dengan timestamp lebih awal dapat dicatat setelah event lebih baru.

---

## D3 — Error Budget & Burn Rate Relationship

```
                      ┌──────────────────────────────────────┐
                      │  SLI Evaluator  (internal/slo/)       │
                      │                                      │
Target SLO: 0.999     │  total = 1100                         │
Allowed error: 0.1%   │  good = 1090                          │
                      │  bad  = 10                              │
                      └─────────────────┬────────────────────┘
                                        │
                      SLI = good/total = 1090/1100 = 99.09%
                                        │
                      ┌─────────────────┴────────────────────┐
                      │  Error Budget Calculation             │
                      │                                       │
                      │  totalBudget = 0.1% × 1100 = 1.1      │
                      │  consumed     = bad       = 10        │
                      │  remaining    = budget - consumed = -8.9  ← EXHAUSTED
                      │                                       │
                      │  CanDeploy = false ✓                  │
                      └─────────────────┬────────────────────┘
                                        │
                      ┌─────────────────┴────────────────────┐
                      │  Alert Engine  (internal/alerting/)   │
                      │                                      │
                      │  shortBurn = actual/allowed           │
                      │  longBurn  = actual/allowed           │
                      │                                      │
                      │  triggered = shortBurn ≥ 14.4    AND   │
                      │             longBurn  ≥ 14.4    (PAGE) │
                      │          OR                          │
                      │             shortBurn ≥ 6.0     AND   │
                      │             longBurn  ≥ 6.0    (TICKET)│
                      │                                      │
                      │  Result: 9.09x ≥ 6.0 → TICKET (✓)     │
                      │          9.09x < 14.4 → PAGE (✗)      │
                      └──────────────────────────────────────┘
```

---

## D4 — Multi-Window Alert: True Positive vs False Positive

```
                    ┌──────────────────────────────────┐
                    │    INCIDENT: 10% ERROR RATE      │
                    └──────────────────┬─────────────┘
                                       │
              ┌────────────────────────┼────────────────────────┐
              │                        │                        │
              ▼                        ▼                        ▼
    ┌──────────────────┐   ┌──────────────────┐   ┌──────────────────┐
    │   TRUE POSITIVE  │   │  FALSE POSITIVE  │   │  TRUE NEGATIVE   │
    │  (alert fires)   │   │  (no alert)      │   │   (no alert)     │
    └──────────────────┘   └──────────────────┘   └──────────────────┘

  Short Window (5m):        Short Window (5m):          Short Window (5m):
  98 OK + 2 ERR             90 OK + 10 ERR              98 OK + 2 ERR
  2% error → 20x burn        10% error → 100x burn      2% error → 20x burn

  Long Window (60m):         Short Window (5m):          Long Window (60m):
  9800 OK + 200 ERR          90 OK + 10 ERR              9999 OK + 1 ERR
  2% error → 20x burn        10% error → 100x burn       0.01% → 0.1x burn

  BOTH ≥ 6.0 → TICKET  ✓   BOTH ≥ 6.0 → TICKET  ✓      Short ≥ 6.0  BUT  Long < 6.0
                           Short ≥ 14.4 → PAGE  ✓                       → NO ALERT ✓
                                                      (transient spike filtered)
```

Validasi: `TestAlertEngineBurnRate` (true positive) + negative test pada line 130-151 (false positive prevention).

---

## D5 — Test Coverage Map

```
┌─────────────────────────────────────────────────────────────────┐
│                    TEST FILE: tests/slo_test.go                  │
└─────────────────────────────┬───────────────────────────────────┘
                              │
    ┌─────────┬───────────────┼───────────────┬─────────┬─────────┐
    ▼         ▼               ▼               ▼         ▼         ▼
┌─────────┐ ┌────────┐ ┌────────────┐ ┌──────────┐ ┌─────────┐ ┌────────────┐
│ Metrics │ │ SLO    │ │ Alert      │ │ Out-of- │ │ Zero    │ │ Concurrency│
│Window   │ │ Eval   │ │ Burn Rate  │ │ Order   │ │Traffic  │ │ (20×100)   │
│Tracker  │ │        │ │            │ │Timestamps│ │         │ │            │
└─────────┘ └────────┘ └────────────┘ └──────────┘ └─────────┘ └────────────┘
    │         │         │           │         │         │
    ▼         ▼         ▼           ▼         ▼         ▼
  Bucket    SLI =    Burn Rate    Bucket    SLI=1.0    Good+Bad =
  agg,      good/   = actual/    ordering  CanDeploy  Total
  eviction  total   allowed      + eviction (zero)    (race clean)
            budget  ratio
```

Semua 6 unit test + 1 concurrency test **PASSED**, race detector **PASSED**.

---

## D6 — Demo Flow State Transition

```
[INITIAL STATE]
  targetSLO = 0.999 (Payment) / 0.950 (Reports)
  budget = +0.1 (Payment) / +5.0 (Reports)
  CanDeploy = true

        │
        ▼
[PHASE 1: Baseline traffic — 1000 good]
  SLI = 100.00%
  Budget Remaining = +0.1 (Payment)
  CanDeploy = true  ✓
        │
        ▼
[PHASE 2: Incident — 100 requests, 10 errors (10%)]
  SLI = 99.09%
  Budget Remaining = -8.90 (Payment)  ← EXHAUSTED
  CanDeploy = false  ✗
        │
        ▼
[PHASE 3: Burn Rate Check]
  Short = 9.09x ≥ 6.0  ✓
  Long  = 9.09x ≥ 6.0  ✓
  → TICKET alert triggered  ⚠
        │
        ▼
[PHASE 4: Criticality Comparison]
  Payment  (99.9%)  → CanDeploy = false  (budget -8.90)
  Reports  (95.0%)  → CanDeploy = false  (budget -5.00, but 5% tolerance)
  Both false, tapi Payment jauh lebih rentan
```
