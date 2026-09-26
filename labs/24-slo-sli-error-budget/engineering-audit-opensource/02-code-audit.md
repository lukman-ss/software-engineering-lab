# Code Audit — labs/24-slo-sli-error-budget

Files: internal/metrics/tracker.go, internal/slo/evaluator.go, internal/alerting/engine.go, cmd/demo/main.go
Toolchain: go1.26.7. `go build ./...` PASS. `go vet ./...` PASS. `gofmt -l` flags 3 files (whitespace only).

## Finding 1

Location: internal/metrics/tracker.go:46-76 `Record`
Claimed Behavior: Sliding-window bucketed aggregation.
Observed Implementation: Merges only into last bucket when `StartTime` equal; out-of-order timestamps append a new bucket out of chronological order.
Assessment: WARNING
Severity: MEDIUM
Notes: Demo/tests record in order so behavior correct there. Out-of-order ingest breaks eviction-order assumption in `evictStaleLocked` (idx-scan from front). No input reordering or sorted insert.

## Finding 2

Location: internal/metrics/tracker.go:30-44 `NewWindowTracker`
Claimed Behavior: Windowed tracker construction.
Observed Implementation: No nil guard on `isGood`; nil callback panics on `Record`. `bucketSize<=0` defaults to 1s; `numBuckets` only sizes capacity, slice unbounded until eviction.
Assessment: WARNING
Severity: LOW
Notes: All call sites pass non-nil `isGood`. Design says "ring buffer"; actual is append+evict slice, bounded by window in practice.

## Finding 3

Location: internal/metrics/tracker.go:89-100 `Summary`
Claimed Behavior: Thread-safe window summary.
Observed Implementation: Takes `mu.Lock` (not `RLock`) then evicts and sums. Safe under `-race` (verified PASS).
Assessment: PASS
Severity: LOW
Notes: Lock instead of RLock is conservative, not a bug. Mutating eviction inside read path requires write lock — correct choice.

## Finding 4

Location: internal/slo/evaluator.go:41-71 `Evaluate`
Claimed Behavior: SLI = good/total; budget = (1-target)*total - bad; freeze when exhausted.
Observed Implementation: Matches. Empty window returns SLI=1.0, budget 0, CanDeploy=true. `budgetRemaining<=0 && total>0` freezes. Demo math reproduced: 1100 total, 99.9% target → budget 1.1, consumed 10, remaining -8.90. Rounding SLI 4dp, budget 2dp.
Assessment: PASS
Severity: LOW
Notes: `LatencyThreshold` in Config unused by evaluator (classification lives in caller `isGood`). Scoped correctly, no miscalc.

## Finding 5

Location: internal/alerting/engine.go:63-89 `Check`
Claimed Behavior: Multi-window burn-rate alerting, short+long evaluated concurrently.
Observed Implementation: Fires only when BOTH shortBurn and longBurn >= factor. `CalculateBurnRate` guards total==0 and allowedErrorRate<=0 → 0.0, no div-zero.
Assessment: PASS
Severity: LOW
Notes: AND semantics matches Google SRE fast/slow-burn paging pattern and demo. Verified: 9.09x fires 6.0x TICKET, not 14.4x PAGE.

## Finding 6

Location: internal/alerting/engine.go:17-24 `BurnRateRule`
Claimed Behavior: Rule carries windows and budget-consumed pct.
Observed Implementation: `LongWindow`, `ShortWindow`, `BudgetConsumedPct` stored but never read by `Check`; windows come from injected trackers.
Assessment: WARNING
Severity: LOW
Notes: Dead config fields. No behavioral effect. Either wire them or drop them.

## Finding 7

Location: cmd/demo/main.go
Claimed Behavior: Baseline → incident → burn-rate alert demo.
Observed Implementation: Deterministic; re-ran output byte-identical to engineering/03-execution-result.md (Phase1 1000/1000/0 budget 1.00; Phase2 1100/1090/10 SLI 99.09% budget -8.90; one TICKET 9.09x alert).
Assessment: PASS
Severity: LOW
Notes: Real execution, no fabrication. Demo comment "burn = 0.10/0.001 = 100x" describes incident-only rate; engine reports windowed 9.09x — display correct, comment scope ambiguous only.
