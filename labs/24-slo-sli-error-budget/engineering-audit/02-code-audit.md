# Code Audit

Target Lab: `labs/24-slo-sli-error-budget`

## Finding 1

Location: `internal/metrics/tracker.go:46-104` (`WindowTracker.Record`)
Claimed Behavior: Thread-safe recording of events into time buckets with support for in-order and out-of-order timestamps and safe window eviction.
Observed Implementation: Uses `sync.RWMutex` (`mu.Lock()` on write). Calls `evictStaleLocked(e.Timestamp)` prior to bucket append/lookup. Handles current bucket match, chronological out-of-order insertion via binary/linear search and slice splicing, and monotonic append.
Assessment: PASS
Severity: LOW
Notes: Linear insertion into slice is bounded by window bucket count (typical window / bucketSize < 10,000 buckets), safe for memory and CPU.

## Finding 2

Location: `internal/metrics/tracker.go:106-128` (`evictStaleLocked`, `Summary`)
Claimed Behavior: Thread-safe window eviction and aggregation of metrics across active buckets.
Observed Implementation: `Summary` acquires `mu.Lock()` and invokes `evictStaleLocked(now)` before summing total, good, and bad counters. Buckets older than `cutoff := now.Add(-w.windowSize)` are sliced out.
Assessment: PASS
Severity: LOW
Notes: Correctly purges expired buckets based on query time.

## Finding 3

Location: `internal/slo/evaluator.go:41-71` (`Evaluator.Evaluate`)
Claimed Behavior: Accurate calculation of SLI ratio, total error budget, budget consumed, budget remaining, and deployment freeze determination.
Observed Implementation: Defaults SLI to 1.0 when `total == 0` avoiding division by zero. Calculates `totalErrorBudget = (1 - targetUptime) * total`, `budgetRemaining = totalErrorBudget - bad`, and flags `canDeploy = false` only when `total > 0 && budgetRemaining <= 0`.
Assessment: PASS
Severity: LOW
Notes: Properly enforces deployment freeze when budget is depleted.

## Finding 4

Location: `internal/alerting/engine.go:51-88` (`AlertEngine.CalculateBurnRate`, `AlertEngine.Check`)
Claimed Behavior: Evaluates multi-window multi-burn-rate alerting requiring both short-window and long-window burn rates to exceed the configured burn rate factor before firing.
Observed Implementation: `CalculateBurnRate` returns `(bad/total) / (1 - targetSLO)`. Handles `total == 0` and invalid `targetSLO >= 1.0` gracefully. `Check` evaluates rules with `shortBurn >= rule.BurnRateFactor && longBurn >= rule.BurnRateFactor`.
Assessment: PASS
Severity: LOW
Notes: Accurately mirrors Google SRE multi-window alerting logic to suppress alerts on transient blips.

## Finding 5

Location: `cmd/demo/main.go:1-151`
Claimed Behavior: Demonstrates baseline traffic, incident impact, multi-window burn rate alert firing, and tiered SLO evaluation between critical and non-critical endpoints.
Observed Implementation: Real simulation using genuine `internal/metrics`, `internal/slo`, and `internal/alerting` types. No hardcoded or fabricated mock outputs.
Assessment: PASS
Severity: LOW
Notes: Output matches the mathematical execution of the engine.
