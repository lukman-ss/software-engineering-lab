# Code Audit — labs/24-slo-sli-error-budget

Files: internal/metrics/tracker.go, internal/slo/evaluator.go, internal/alerting/engine.go, cmd/demo/main.go

## Finding 1

Location: internal/metrics/tracker.go:117-128 (Summary)
Claimed Behavior: Sliding-window aggregation over [now-window, now].
Observed Implementation: Evicts buckets with StartTime < now-window, sums ALL remaining buckets including StartTime > now. No upper-bound filter.
Assessment: WARNING
Severity: MEDIUM
Notes: Future-timestamped buckets leak into earlier reads. No impact on claimed monotonic-time flows (tests/demo use now >= event times). Fix: skip buckets with StartTime.After(now).

## Finding 2

Location: internal/slo/evaluator.go:54-57
Claimed Behavior: Freeze deploy when budget exhausted.
Observed Implementation: `budgetRemaining <= 0` on unrounded float. Boundary exact-zero depends on float artifact (1-0.99=0.010000000000000009; test's CanDeploy=true passes only via 8.9e-16 residue).
Assessment: WARNING
Severity: MEDIUM
Notes: Exact-exhaustion verdict nondeterministic by representation. Use integer counts or epsilon for decision; round only for display.

## Finding 3

Location: internal/alerting/engine.go:17-24 (BurnRateRule)
Claimed Behavior: Multi-window burn-rate alerting with configured windows.
Observed Implementation: `LongWindow`, `ShortWindow`, `BudgetConsumedPct` stored, never read. Engine uses pre-wired trackers; rule durations unenforced. Mismatched tracker/rule windows fail silently.
Assessment: WARNING
Severity: MEDIUM
Notes: Dead config fields. Either wire rule windows to tracker selection/validation or remove fields.

## Finding 4

Location: internal/slo/evaluator.go:10-14 (Config.LatencyThreshold)
Claimed Behavior: Config carries latency threshold.
Observed Implementation: Field stored, never read. Good/bad decided by tracker's isGood closure. Demo sets both independently; divergence possible.
Assessment: WARNING
Severity: LOW
Notes: Dead field. Remove or enforce consistency.

## Finding 5

Location: internal/metrics/tracker.go:30-44 (NewWindowTracker), internal/slo/evaluator.go:41, internal/alerting/engine.go:51-61
Claimed Behavior: None explicit (no validation claimed).
Observed Implementation: nil isGoodEvent panics on Record; windowSize<=0 unhandled; TargetUptime outside [0,1] yields inverted budget; target=1.0 forces burn rate 0 (never alerts).
Assessment: PASS (with note)
Severity: LOW
Notes: No input validation. Acceptable for lab scope; callers in tests/demo always valid. Document preconditions.

## Finding 6

Location: internal/metrics/tracker.go:46-48, 117-119; tests/slo_test.go:TestConcurrencyMetrics
Claimed Behavior: Thread-safe tracker.
Observed Implementation: Record takes write Lock; Summary takes write Lock (mutates via eviction). `go test -race` clean.
Assessment: PASS
Severity: LOW
Notes: Correct locking. Concurrency test covers Record/Record contention only (see test audit).

## Finding 7

Location: internal/metrics/tracker.go:106-115 (evictStaleLocked)
Claimed Behavior: Stale eviction.
Observed Implementation: Whole-bucket eviction, `Before(cutoff)` keeps boundary bucket. Standard bucketed approximation.
Assessment: PASS
Severity: LOW
Notes: Expected granularity tradeoff, documented bucket-size tradeoff in implementation notes.

## Finding 8

Location: internal/alerting/engine.go:63-89 (Check)
Claimed Behavior: Alert when short AND long burn exceed factor (reset/suppress transient spikes).
Observed Implementation: Both must exceed; negative transient test proves suppression.
Assessment: PASS
Severity: LOW
Notes: Matches Google multiwindow semantics.
