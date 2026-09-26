# Code Audit

Target Lab: `labs/24-slo-sli-error-budget`

## Finding 1

Location: `internal/metrics/tracker.go:46-104`
Claimed Behavior: Thread-safe recording of events into time-bucketed sliding windows, with support for in-order and out-of-order timestamps.
Observed Implementation:
- Mutex locked at entry (`w.mu.Lock()`).
- Calls `evictStaleLocked(e.Timestamp)` on write to keep buckets within window bounds.
- Checks if bucket matches existing latest bucket; if earlier, searches and inserts at sorted index or increments matching bucket.
- Correctly updates `TotalCount`, `GoodCount`, and `BadCount`.
Assessment: PASS
Severity: LOW
Notes: Linear search/insert on out-of-order events is efficient for moderate window bucket sizes (e.g. seconds to hours).

## Finding 2

Location: `internal/metrics/tracker.go:106-115`
Claimed Behavior: Proper sliding window eviction of stale buckets.
Observed Implementation:
- Computes `cutoff := now.Add(-w.windowSize)`.
- Scans sorted buckets from head and slices out buckets where `StartTime.Before(cutoff)`.
- Eviction called safely under lock in both `Record` and `Summary`.
Assessment: PASS
Severity: LOW
Notes: Buckets are maintained in chronological order, allowing slice-reslicing eviction.

## Finding 3

Location: `internal/slo/evaluator.go:41-71`
Claimed Behavior: Evaluates SLI ratio, total error budget, consumed budget, remaining budget, and release freeze policy (`CanDeploy`).
Observed Implementation:
- If `total == 0`, defaults `sli` to `1.0` (100%), preventing divide-by-zero panics.
- Computes `totalErrorBudget = (1.0 - TargetUptime) * total`.
- Computes `budgetRemaining = totalErrorBudget - budgetConsumed`.
- If `total > 0` and `budgetRemaining <= 0`, sets `CanDeploy = false`.
Assessment: PASS
Severity: LOW
Notes: Mathematical implementation matches Google SRE formula.

## Finding 4

Location: `internal/alerting/engine.go:51-61`
Claimed Behavior: Burn rate calculation per window with zero-traffic safeguard.
Observed Implementation:
- Handles `total == 0` by returning `0.0`.
- Safeguards against `allowedErrorRate <= 0` returning `0.0`.
- Correctly calculates `(bad / total) / (1 - targetSLO)`.
Assessment: PASS
Severity: LOW
Notes: Zero division completely guarded.

## Finding 5

Location: `internal/alerting/engine.go:63-89`
Claimed Behavior: Multi-window multi-burn-rate alert triggering logic requiring both short and long window conditions.
Observed Implementation:
- Evaluates `shortBurn` and `longBurn`.
- Triggers alert if and only if `shortBurn >= rule.BurnRateFactor && longBurn >= rule.BurnRateFactor`.
- Accurately captures severity and rate metrics into `AlertResult`.
Assessment: PASS
Severity: LOW
Notes: Implements the Google SRE multi-window burn rate requirement preventing false alerts on single transient spikes.
