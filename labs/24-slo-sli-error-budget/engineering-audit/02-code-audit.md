# Code Audit

Target Lab: `labs/24-slo-sli-error-budget`

## Finding 1

Location: `internal/metrics/tracker.go:22-128`
Claimed Behavior: Thread-safe sliding-window event aggregation with time-bucket eviction and out-of-order insertion.
Observed Implementation: `WindowTracker` uses a `sync.RWMutex` protecting `buckets []Bucket`. On `Record(e Event)`, it locks, evicts stale buckets relative to `e.Timestamp - windowSize`, truncates to `bucketSize`, and either appends or inserts sorted by `StartTime`. On `Summary(now)`, it evicts relative to `now - windowSize` and sums total, good, and bad counts.
Assessment: PASS
Severity: LOW
Notes: Correctly handles chronological and reverse-chronological event arrival; race detector verified with 20 concurrent goroutines.

## Finding 2

Location: `internal/slo/evaluator.go:41-71`
Claimed Behavior: Calculates SLI ratio `good/total`, remaining budget `((1 - Target) * total) - bad`, and determines deployment freeze `CanDeploy = false` when budget is depleted.
Observed Implementation: `Evaluate(now)` handles `total == 0` safely with default SLI `1.0` and `canDeploy = true`. When `total > 0`, correctly calculates remaining budget and flags `canDeploy = false` when `budgetRemaining <= 0`.
Assessment: PASS
Severity: LOW
Notes: Output numbers are cleanly rounded using `math.Round` for stability.

## Finding 3

Location: `internal/alerting/engine.go:51-88`
Claimed Behavior: Multi-window multi-burn-rate alerting requiring both short and long window burn rates to breach the factor threshold.
Observed Implementation: `CalculateBurnRate` returns `(bad / total) / (1 - targetSLO)`. `Check(now)` fetches summaries from `shortTracker` and `longTracker` and evaluates `shortBurn >= rule.BurnRateFactor && longBurn >= rule.BurnRateFactor`.
Assessment: PASS
Severity: LOW
Notes: Handles `total == 0` and zero allowed error rate edge cases cleanly.

## Finding 4

Location: `cmd/demo/main.go:1-151`
Claimed Behavior: Executable demonstration covering baseline traffic, severe incident budget burn, burn-rate alerting, and multi-endpoint criticality comparison.
Observed Implementation: 4 distinct phases executed with realistic request events, demonstrating SLI degradation, error budget depletion below zero, deployment blocking, slow burn alert triggering, and endpoint comparison between Payment (99.9%) and Reports (95.0%).
Assessment: PASS
Severity: LOW
Notes: Output is fully reproducible and matches documented execution logs.
