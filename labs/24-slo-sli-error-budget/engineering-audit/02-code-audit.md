# Code Audit

Target Lab: `labs/24-slo-sli-error-budget`

## Finding 1

Location: `internal/metrics/tracker.go:46-104`
Claimed Behavior: Thread-safe, ordered bucket insertion for sliding window metrics.
Observed Implementation: `Record(e Event)` acquires write lock `w.mu.Lock()`, runs `evictStaleLocked(e.Timestamp)`, truncates bucket start time, updates existing matching bucket or inserts out-of-order buckets at the correct index maintaining sorted order by `StartTime`.
Assessment: PASS
Severity: LOW
Notes: Correctly handles both in-order and out-of-order timestamps without panics or index corruption.

## Finding 2

Location: `internal/metrics/tracker.go:106-115`
Claimed Behavior: Sliding window eviction of stale buckets.
Observed Implementation: `evictStaleLocked(now)` computes `cutoff := now.Add(-w.windowSize)`. Since buckets are maintained in chronological order, it slices out all buckets where `StartTime.Before(cutoff)`.
Assessment: PASS
Severity: LOW
Notes: Linear scan on sorted bucket list slices out expired buckets cleanly.

## Finding 3

Location: `internal/slo/evaluator.go:41-70`
Claimed Behavior: SLI ratio and remaining error budget calculation with deployment freeze trigger.
Observed Implementation: Evaluates `good / total` ratio. Computes `allowedFailureRate := 1.0 - targetUptime`, `totalErrorBudget := allowedFailureRate * total`, and `budgetRemaining := totalErrorBudget - bad`. If `total > 0 && budgetRemaining <= 0`, flags `CanDeploy = false`. Handles zero traffic safely with default SLI = 1.0 and `CanDeploy = true`.
Assessment: PASS
Severity: LOW
Notes: Mathematical definitions match Google SRE Handbook principles.

## Finding 4

Location: `internal/alerting/engine.go:51-88`
Claimed Behavior: Multi-window multi-burn-rate alerting logic.
Observed Implementation: `CalculateBurnRate` computes `actualErrorRate / allowedErrorRate`. `Check(now)` samples both short window and long window trackers. Alert triggers if and only if both `shortBurn >= rule.BurnRateFactor && longBurn >= rule.BurnRateFactor`.
Assessment: PASS
Severity: LOW
Notes: Prevents false alarms on brief transient spikes as required by Google SRE alerting design.

## Finding 5

Location: `internal/metrics/tracker.go:117-128`
Claimed Behavior: Thread-safe summary computation.
Observed Implementation: `Summary(now)` locks mutex `w.mu.Lock()`, executes stale bucket eviction against `now`, and aggregates total, good, and bad event counters across remaining buckets.
Assessment: PASS
Severity: LOW
Notes: Thread-safety verified with `-race` flag.
