# Code Audit

## Finding 1
Location: `internal/metrics/tracker.go` — `Record`, `evictStaleLocked`, `Summary`; `sync.RWMutex` guarding `buckets`.
Claimed Behavior: Thread-safe sliding-window event aggregation with eviction of stale buckets on every read/write.
Observed Implementation: `Record` and `Summary` both acquire `w.mu.Lock()` (exclusive). `evictStaleLocked` is called under the lock in both `Record` and `Summary`, mutating state. Out-of-order timestamps are handled by scanning from the front to find an equal or insertion point, preserving an ascending `StartTime` invariant; eviction slices from the front only.
Assessment: PASS
Severity: LOW (informational)
Notes: `Summary` uses a write-lock even though it is a read-only aggregation — this is correct, not a bug, because it calls `evictStaleLocked` which mutates. No data race; `go test -race` clean. The slice invariant is preserved by every mutating path. The `bucketStart.Before(cutoff)` eviction boundary is inclusive of buckets exactly at `cutoff` (not evicted), a defensible interpretation.

## Finding 2
Location: `internal/metrics/tracker.go:88` — `w.buckets = append(w.buckets[:i], append([]Bucket{b}, w.buckets[i:]...)...)`
Claimed Behavior: O(n) in-slice insertion for out-of-order events.
Observed Implementation: Correct two-append insert that shifts elements. Holds the write lock, so safe.
Assessment: PASS
Severity: LOW
Notes: O(n) per out-of-order insert; acceptable for a demo/educational tracker. Not a correctness issue.

## Finding 3
Location: `internal/slo/evaluator.go:41` — `Evaluate`.
Claimed Behavior: SLI = good/total; total error budget = (1 - SLO) * total; consumed = bad; remaining = budget - consumed; release frozen when remaining <= 0 and total > 0.
Observed Implementation: Matches exactly. SLI rounded to 4 decimals for display; budget math uses unrounded values (more accurate than display). Zero-traffic → SLI defaults to 1.0, `CanDeploy=true`.
Assessment: PASS
Severity: LOW
Notes: Rounding of `CurrentSLI`/`TotalErrorBudget`/`BudgetRemaining` is display-only; policy decision uses raw floats. Verified against demo numbers: Phase 2 (1090/1100, SLO 0.999) → budget 1.1, remaining -8.90, CanDeploy=false. Correct.

## Finding 4
Location: `internal/alerting/engine.go:51` — `CalculateBurnRate`.
Claimed Behavior: Burn rate = actualErrorRate / allowedErrorRate, with divide-by-zero guards (total == 0, allowedErrorRate <= 0 return 0.0).
Observed Implementation: Correct formula and guards. Phase 3 demo: (10/1100)/0.001 = 9.09x — matches output. 14.4x rule not met (9.09 < 14.4); 6.0x rule met (9.09 >= 6.0) → only TICKET triggered. Matches demo.
Assessment: PASS
Severity: LOW
Notes: The `>=` comparison (not `>`) means a burn rate exactly equal to the factor triggers — reasonable.

## Finding 5
Location: `internal/alerting/engine.go:63` — `Check`.
Claimed Behavior: Multi-window (short+long) burn-rate alert; rule fires only when BOTH short and long burn rates meet the factor (suppresses transient spikes).
Observed Implementation: `shortBurn >= factor && longBurn >= factor`. Uses the single injected short/long tracker pair; does NOT use the per-rule `LongWindow`/`ShortWindow`/`BudgetConsumedPct` fields (see Finding 6). Negative transient test confirms suppression when long window is clean.
Assessment: PASS (behavior is self-consistent and tested)
Severity: LOW
Notes: The "both windows" condition is a conservative (simplified) variant of canonical SRE multi-window alerting; it is not a crash/data bug and all demonstrated behavior is reproducible.

## Finding 6
Location: `internal/alerting/engine.go:17-24` — `BurnRateRule` struct fields `LongWindow`, `ShortWindow`, `BudgetConsumedPct`.
Claimed Behavior: (Implied by field presence) per-rule window configuration for multi-window multi-burn-rate alerting.
Observed Implementation: These fields are never read in `Check`; windows come from the trackers injected into `AlertEngine`. In the demo these fields are zero-valued. They are dead configuration.
Assessment: WARNING
Severity: LOW
Notes: Functional simplification, documented as in-memory/demo-scoped. Does not break the demo or tests, but over-declares per-rule configurability. Upgrade path: route each rule through its own tracker pair or compute burn rates per `ShortWindow`/`LongWindow`.

## Finding 7
Location: `internal/slo/evaluator.go:13` — `Config.LatencyThreshold`.
Claimed Behavior: Latency threshold used by the SLO evaluator to classify events.
Observed Implementation: `LatencyThreshold` is stored in `Config` but never read in `Evaluate`. Latency classification is performed by the tracker's `isGood` closure (see `cmd/demo/main.go:20-22`), not by the evaluator.
Assessment: WARNING
Severity: LOW
Notes: Redundant dead field. Correct behavior still occurs (via the closure), so no functional defect.

## Finding 8 (no issue)
Location: concurrency / locking across `Evaluator` and `AlertEngine`.
Claimed Behavior: No cross-tracker nested locking.
Observed Implementation: `Evaluate` locks only the SLO tracker; `Check` locks short then long tracker sequentially (each lock released before the next is acquired). Single mutexes, no nesting, no deadlock surface.
Assessment: PASS
Severity: —
Notes: `go test -race` passes (see Test Audit).
