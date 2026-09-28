# Code Audit

## Finding 1

Location: internal/metrics/tracker.go:109
Claimed Behavior: Sliding window evicts events older than windowSize.
Observed Implementation: evictStaleLocked uses `w.buckets[idx].StartTime.Before(cutoff)` (strict before). A bucket with start time exactly equal to cutoff (`now - windowSize`) is not evicted.
Assessment: WARNING
Severity: LOW
Notes: Boundary condition may cause slightly larger effective window. Tests do not hit this edge case. Behavior is consistent but may deviate from strict window semantics.

## Finding 2

Location: internal/metrics/tracker.go:67-92
Claimed Behavior: Out-of-order timestamps insert bucket in correct chronological order.
Observed Implementation: Insertion loop iterates from index 0 and inserts before first bucket with start time after the new bucket. Works correctly for all observed cases.
Assessment: PASS
Severity: N/A
Notes: Logic is correct and passes all tests including out-of-order scenarios. Performance is O(n) per insert but acceptable.

## Finding 3

Location: internal/slo/evaluator.go
Claimed Behavior: Config.LatencyThreshold influences SLI calculation.
Observed Implementation: Config.LatencyThreshold field is stored but never used. SLI determination relies solely on the tracker's isGood function, which is set by the caller.
Assessment: FAIL
Severity: MEDIUM
Notes: Dead field creates misleading API. The Config suggests latency-based SLI support, but the evaluator ignores it. Demo works because caller's isGood uses latency directly.

## Finding 4

Location: internal/alerting/engine.go
Claimed Behavior: BurnRateRule.LongWindow and ShortWindow configure alert windows.
Observed Implementation: These fields are never used. Alert windows are determined solely by the shortTracker and longTracker passed to NewAlertEngine.
Assessment: WARNING
Severity: LOW
Notes: Unused fields indicate design drift. The engine correctly implements multi-window logic via separate trackers, but the rule struct contains dead code.

## Finding 5

Location: internal/metrics/tracker.go:118
Claimed Behavior: Summary provides thread-safe snapshot.
Observed Implementation: Summary takes a write lock (mu.Lock()) because it calls evictStaleLocked which mutates buckets.
Assessment: PASS
Severity: N/A
Notes: Correctly handles concurrent modification via eviction. Race detector passes, confirming safety.

## Finding 6

Location: internal/alerting/engine.go:63-89
Claimed Behavior: Multi-window burn-rate alert triggers only when both windows exceed threshold.
Observed Implementation: Check implements `if shortBurn >= factor && longBurn >= factor`. Matches specification.
Assessment: PASS
Severity: N/A
Notes: Logic aligns with Google SRE Workbook recommendations. Demo confirms correct triggering.