# Code Audit — labs/24-slo-sli-error-budget

## Finding 1
Location: internal/metrics/tracker.go:46-104 (Record), 106-115 (evictStaleLocked)
Claimed Behavior: Sliding-window bucketed good/total tracking with stale eviction.
Observed Implementation: RWMutex-guarded slice of Buckets; evict on Record+Summary; out-of-order insert via sorted splice; same-bucket fast path.
Assessment: PASS
Severity: LOW
Notes: O(n) splice insert fine at lab scale. Boundary `Before(cutoff)` retains bucket exactly at cutoff — reasonable semantics.

## Finding 2
Location: internal/metrics/tracker.go:30-44 (NewWindowTracker), 46 (Record)
Claimed Behavior: Robust tracker construction.
Observed Implementation: No nil-guard on `isGood`; nil func panics on Record. No validation of `windowSize <= 0` (zero window evicts everything except exact-now bucket).
Assessment: WARNING
Severity: LOW
Notes: No test hits these paths; normal construction unaffected.

## Finding 3
Location: internal/metrics/tracker.go:117-128 (Summary)
Claimed Behavior: Thread-safe read.
Observed Implementation: Takes write `mu.Lock` (not RLock) because eviction mutates slice. Correct.
Assessment: PASS
Severity: LOW
Notes: `go test -race` clean on 20×100 concurrent Record test.

## Finding 4
Location: internal/slo/evaluator.go:41-71 (Evaluate)
Claimed Behavior: SLI=good/total; budget=(1-target)*total-bad; freeze when exhausted.
Observed Implementation: Matches. Zero-traffic SLI=1.0, CanDeploy=true. Rounding SLI 4dp, budget 2dp.
Assessment: PASS
Severity: LOW
Notes: Verified math against demo: 1100 total/10 bad @99.9% → budget 1.1-10=-8.9, SLI 0.9909. Correct.

## Finding 5
Location: internal/slo/evaluator.go:10-14 (Config.LatencyThreshold), 54-57 (freeze `<= 0`)
Claimed Behavior: Latency threshold enforced; freeze when budget exhausted.
Observed Implementation: `LatencyThreshold` stored, never read — goodness delegated to tracker `isGood` closure. Freeze uses `budgetRemaining <= 0`; exact-zero boundary depends on float artifact (1-0.99=0.010000000000000009 saves the TestSLOEvaluator CanDeploy=true case).
Assessment: WARNING
Severity: LOW
Notes: Dead field, fragile boundary. No behavioral failure observed.

## Finding 6
Location: internal/alerting/engine.go:51-89 (CalculateBurnRate, Check)
Claimed Behavior: Multi-window multi-burn-rate alerting; fast 14.4x page, slow 6x ticket.
Observed Implementation: Dual-window check real: requires BOTH shortBurn AND longBurn >= factor; transient-spike negative case proven by test. BUT per-rule `LongWindow`/`ShortWindow`/`BudgetConsumedPct` fields declared, never used — all rules evaluated against single shared tracker pair, differing only by `BurnRateFactor`.
Assessment: WARNING
Severity: MEDIUM
Notes: Core false-positive suppression works; rule-specific windows unimplemented. Zero-traffic returns 0 (no alert) — sane. Zero-factor rule would always trigger — untested edge.

## Finding 7
Location: internal/alerting/engine.go:63-89; internal/metrics/tracker.go (locks)
Claimed Behavior: Concurrency-safe alerting.
Observed Implementation: Engine holds no own lock; safety inherited from tracker locks. No shared mutable engine state.
Assessment: PASS
Severity: LOW
Notes: Race detector clean.

## Finding 8
Location: cmd/demo/main.go (all phases)
Claimed Behavior: Baseline → incident → alerts → criticality comparison.
Observed Implementation: Real computation via same packages; rerun output byte-identical to engineering/03-execution-result.md. Phase 3 honestly shows only TICKET (9.09x < 14.4x page threshold), not fabricated page alert.
Assessment: PASS
Severity: LOW
Notes: No recovery phase despite design mentioning recovery. Demo timescales compressed (seconds vs 30min windows) — disclosed.
