# Code Audit

## Finding 1

Location: internal/metrics/tracker.go
Claimed Behavior: SLI tracks good vs. total events across rolling time windows; out-of-order events are handled; stale buckets are evicted.
Observed Implementation: Record() evicts stale buckets by cutoff (now - windowSize), truncates event timestamp to bucket size, and inserts out-of-order events via linear scan + slice insertion. Summary() evicts stale buckets then aggregates counts.
Assessment: PASS
Severity: LOW
Notes: Out-of-order handling in Record() covers both the "later event first, earlier event second" case and the "same bucket" case correctly. Eviction is O(n) per record but acceptable for the lab.

## Finding 2

Location: internal/metrics/tracker.go
Claimed Behavior: Concurrency safety for concurrent Record/Summary callers.
Observed Implementation: Single sync.RWMutex guarding all bucket mutations. Record() uses Lock(); Summary() uses Lock() (write lock) rather than RLock().
Assessment: PASS with minor inefficiency
Severity: LOW
Notes: Summary() uses a write Lock because it mutates state via evictStaleLocked(). This is correct (not a race), but RLock would suffice if eviction were moved out of Summary. Verified by `go test -race` passing with 20 goroutines x 100 requests.

## Finding 3

Location: internal/slo/evaluator.go
Claimed Behavior: Error Budget = (1 - SLO) * total; consumed on bad events; release freeze enforced when budget <= 0.
Observed Implementation: allowedFailureRate = 1 - TargetUptime; totalErrorBudget = allowedFailureRate * total; budgetConsumed = bad; budgetRemaining = totalErrorBudget - budgetConsumed; canDeploy = false when total>0 && budgetRemaining <= 0. Zero-traffic -> SLI=1.0, CanDeploy=true.
Assessment: PASS
Severity: LOW
Notes: Budget math is correct. The CanDeploy decision uses the unrounded budgetRemaining (correct), while the reported field is rounded.

## Finding 4

Location: internal/slo/evaluator.go
Claimed Behavior: LatencyThreshold is part of SLO config.
Observed Implementation: Config.LatencyThreshold is declared but never read by Evaluator. Latency breaches are instead enforced via the isGood callback passed to WindowTracker.
Assessment: WARNING
Severity: LOW
Notes: Not a correctness bug (latency is still enforced), but the Config field is dead and misrepresents where latency is handled. Confusing for maintainers.

## Finding 5

Location: internal/alerting/engine.go
Claimed Behavior: Multi-window burn-rate alerting using ShortWindow/LongWindow/BudgetConsumedPct rule config.
Observed Implementation: BurnRateRule declares LongWindow, ShortWindow, and BudgetConsumedPct fields. None are referenced in Check(). Only BurnRateFactor is used. Trigger condition: shortBurn >= factor AND longBurn >= factor (both windows must breach).
Assessment: PASS (alerting logic correct), WARNING (unused fields)
Severity: LOW
Notes: The unused fields do not affect correctness but imply configuration that is silently ignored. The multi-window conjunction check matches the Google SRE "both windows must fire" intent (reduces false positives).

## Finding 6

Location: internal/metrics/tracker.go, internal/slo/evaluator.go, internal/alerting/engine.go
Claimed Behavior: Standard-library-only implementation; no external dependencies.
Observed Implementation: Only `time`, `sync`, `math`, `fmt` used. go.mod declares no require directives.
Assessment: PASS
Severity: LOW
Notes: Matches design claim (01-design.md item 5).

## Summary
No correctness or safety defects found in the core math, state transitions, eviction, or failure handling. Two low-severity dead-field warnings (Findings 4 and 5). Concurrency is safe under the race detector.