# Code Audit — labs/24-slo-sli-error-budget

## Finding 1 — BudgetConsumedPct, LongWindow, ShortWindow fields are dead code
Location: internal/alerting/engine.go:17-24 (BurnRateRule), NewAlertEngine:42, Check:71-87
Claimed Behavior: BurnRateRule defines LongWindow, ShortWindow, and BudgetConsumedPct to drive per-rule multi-window / budget-consumed alerting (per design 01-design.md: "BurnRateAlertEngine monitors short and long windows for fast/slow burn threshold breaches").
Observed Implementation: Check iterates over `a.rules` using only `rule.BurnRateFactor`. The engine's short/long windows are fixed at construction time and shared across ALL rules; LongWindow, ShortWindow, and BudgetConsumedPct are never read.
Assessment: WARNING
Severity: LOW
Notes: Not a correctness bug (the demo still works). It is a deviation from the documented "multi-burn-rate" flexibility (each rule was meant to carry its own window pair + consumed-budget threshold). Every rule currently inherits the same engine-level window pair and the same short/long AND condition, so the "multi" in multi-burn-rate collapses to "same threshold on two fixed windows". Skipped design intent becomes dead fields. Upgrade path: pass per-rule window configs and evaluate short vs long against the rule's own BudgetConsumedPct.

## Finding 2 — WindowTracker has no capacity cap; numBuckets is advisory
Location: internal/metrics/tracker.go:30-44 (NewWindowTracker), Record:46-104
Claimed Behavior: windowSize/bucketSize define a bounded sliding window; stale buckets are evicted on read/write via evictStaleLocked.
Observed Implementation: buckets grows via append in Record; eviction only removes buckets whose StartTime < (timestamp - windowSize). There is no cap enforcing that len(buckets) <= numBuckets. A steady stream of out-of-order (older-bucket) inserts appends forever.
Assessment: WARNING
Severity: MEDIUM
Notes: Under monotonic increasing timestamps (the only happy path tested) eviction keeps the slice bounded; no bug today. But the documented "ring/time-bucketed" bounded behavior is not actually enforced. A malicious or replayed workload with widely scattered old timestamps would grow the slice unboundedly -> memory leak. The numBuckets capacity on the slice (`make([]Bucket, 0, numBuckets)`) is a hint, not a limit. Fix: cap slice length after eviction to numBuckets by trimming oldest; or enforce per-bucket count.

## Finding 3 — Record() out-of-order insert is O(n) and only handles "earlier existing bucket"
Location: internal/metrics/tracker.go:67-92
Claimed Behavior: out-of-order events insert into an existing older bucket (covered by TestOutOfOrderTimestamps: +2s event, +2s event with +100ms, +5s event).
Observed Implementation: The `if n>0 && bucketStart.Before(w.buckets[n-1].StartTime)` branch scans from the front to find an exact matching StartTime, inserting a new bucket in order if none matches. This is O(n) per record and O(n^2) for n events. It also only fires when the new event is *earlier* than the last bucket's start.
Assessment: WARNING
Severity: LOW
Notes: TestOutOfOrderTimestamps exercises same-bucket accumulation and a partial-eviction read. Correctness is fine for the tested scenario. Performance degrades on large out-of-order streams; acceptable for a lab but worth a `ponytail:` comment in real code. No failure path exercised.

## Finding 4 — Summary() takes a full write Lock instead of RLock
Location: internal/metrics/tracker.go:117-128
Claimed Behavior: concurrent reads are safe (TestConcurrencyMetrics asserts good+bad==total under concurrency).
Observed Implementation: Summary acquires `w.mu.Lock()` (exclusive) rather than `RLock()`, so concurrent Summarize calls serialize unnecessarily; no correctness risk, just contention.
Assessment: PASS
Severity: LOW
Notes: Safe but suboptimal. The test still proves no data loss under concurrent writers.

## Finding 5 — SLI/SLI rounding is inconsistent between fields
Location: internal/slo/evaluator.go:41-71
Claimed Behavior: SLI accuracy and budget math are correct.
Observed Implementation: CurrentSLI is rounded to 4 decimals; BudgetConsumed is NOT rounded (raw float64(bad)); TotalErrorBudget and BudgetRemaining are rounded to 2 decimals; CanDeploy uses the UNROUNDED budgetRemaining. Rounding is purely presentational in Status — logic is unaffected.
Assessment: PASS
Severity: LOW
Notes: Demo Phase 3 shows SLI 99.0900% (4-decimal). BudgetRemaining -8.90 (2-decimal). Consistent with code. No bug.

## Finding 6 — SLI defined as 1.0 on zero traffic
Location: internal/slo/evaluator.go:44-47
Claimed Behavior: sensible default when no events have occurred.
Observed Implementation: `var sli float64 = 1.0; if total > 0 { sli = good/total }` -> returns 1.0 for zero traffic; CanDeploy stays true. Matches TestEvaluatorZeroTraffic.
Assessment: PASS
Severity: LOW
Notes: Reasonable convention (no measured failures => assume ok). Validated by test.

## Finding 7 — Future-bucket events counted in Summary
Location: internal/metrics/tracker.go:106-127 (evictStaleLocked + Summary)
Claimed Behavior: sliding window covers [now - window, now].
Observed Implementation: evictStaleLocked evicts only buckets with StartTime.Before(now - windowSize). Buckets with StartTime *after* `now` (future-dated events) are retained and summed in Summary. A Summary taken at time T includes any bucket whose StartTime is >= T-window (including future).
Assessment: WARNING
Severity: LOW
Notes: Not exercised by tests (demo uses simTime/evalTime that are ahead of all recorded timestamps). In a production system with clock skew, future-dated events could inflate counts. Low risk for the lab's monotonic-time simulation but it is a latent correctness gap. Fix: evict buckets whose StartTime.After(now) too, i.e. clamp window to [now-window, now].

## Finding 8 — allowedErrorRate<=0 guard returns 0 burn (SLO=1.0 would divide by zero)
Location: internal/alerting/engine.go:51-61
Claimed Behavior: burn rate = actualErrorRate / allowedErrorRate.
Observed Implementation: if allowedErrorRate <= 0 (i.e. SLO==1.0, no budget), returns 0.0 to avoid div-by-zero. Correct but silently masks a degenerate SLO=100% config by never alerting.
Assessment: PASS
Severity: LOW
Notes: No test asserts SLO=1.0 behavior. Acceptable. A warning log would be preferable but out of scope.

## Finding 9 — alertRules LongWindow/ShortWindow fields omit (gofmt alignment)
Location: internal/alerting/engine.go:17-24
Claimed Behavior: n/a
Observed Implementation: struct field `BudgetConsumedPct float64` is not gofmt-aligned with other float fields, but otherwise unused. gofmt does not enforce intra-struct spacing across separate groups here.
Assessment: PASS
Severity: LOW
Notes: Cosmetic; build is clean. Not an issue.

## Finding 10 — demo hard-codes window scale, not real 30-day windows
Location: cmd/demo/main.go:24-30
Claimed Behavior: design says "in-memory ring/time-bucketed window tracking to simulate 30-day/rolling windows in compressed real-time."
Observed Implementation: window30d = 30 * time.Minute (30 minutes, NOT 30 days), short=5m, long=60m. The variable name `window30d` is misleading but this is an intentional compressed-time simulation, explicitly a known implementation decision.
Assessment: PASS
Severity: LOW
Notes: Matches design decision. Verified by running the demo; output is internally consistent (Phase 2 budget -8.90 = (1-0.001)*1100 - 10 = 1.1 - 10 = -8.9).
