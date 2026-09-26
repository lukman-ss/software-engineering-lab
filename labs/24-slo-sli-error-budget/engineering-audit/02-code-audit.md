# Code Audit

Target Lab: `labs/24-slo-sli-error-budget`

## Finding 1: WindowTracker Eviction and Sliding-Window Bucketing

Location: `internal/metrics/tracker.go:50-87`
Claimed Behavior: Thread-safe sliding time window tracker storing bucketed events, evicting stale events based on `cutoff = now.Add(-w.windowSize)`.
Observed Implementation:
- `Record` takes write lock (`w.mu.Lock()`) and calls `evictStaleLocked(e.Timestamp)`.
- Truncates event timestamp by `bucketSize` and merges counts into existing latest bucket or appends new bucket.
- `Summary` takes write lock and evicts stale buckets before summing `TotalCount`, `GoodCount`, `BadCount`.
- Note: If events arrive out-of-order prior to `buckets[n-1].StartTime`, events create a new bucket appended at end rather than sorted. For in-order synthetic generation this works as intended.
Assessment: PASS
Severity: LOW
Notes: Clean synchronization, proper lock deferral, slice memory reclaimed via reslicing.

## Finding 2: SLI and Error Budget Calculation

Location: `internal/slo/evaluator.go:41-71`
Claimed Behavior: Calculates SLI ratio (`good / total`), total error budget `(1 - TargetUptime) * total`, remaining error budget `totalErrorBudget - BadEvents`, and release freeze boolean `CanDeploy`.
Observed Implementation:
- Zero total events returns `sli = 1.0` and `canDeploy = true` gracefully.
- Formula matches Google SRE book definition: `allowedFailureRate := 1.0 - e.config.TargetUptime`, `totalErrorBudget := allowedFailureRate * float64(total)`, `budgetRemaining := totalErrorBudget - budgetConsumed`.
- `CanDeploy` evaluates `false` strictly when `total > 0 && budgetRemaining <= 0`.
Assessment: PASS
Severity: LOW
Notes: Mathematical precision handled with `math.Round` for clean output.

## Finding 3: Multi-Window Multi-Burn-Rate Alert Evaluation

Location: `internal/alerting/engine.go:51-88`
Claimed Behavior: Multi-window burn rate evaluation checking short and long rolling windows against threshold factors (`BurnRateFactor`).
Observed Implementation:
- `CalculateBurnRate`: `(bad / total) / (1.0 - targetSLO)`. Returns `0.0` when total is 0 or allowedErrorRate <= 0.
- `Check`: Computes `shortBurn` and `longBurn`. Triggers alert only when `shortBurn >= rule.BurnRateFactor && longBurn >= rule.BurnRateFactor`.
- Matches SRE multi-window requirement (both short and long windows must fire).
Assessment: PASS
Severity: LOW
Notes: Simple and correct evaluation logic.

## Finding 4: Concurrency and Thread Safety

Location: `internal/metrics/tracker.go:23, 47, 90`
Claimed Behavior: Safe concurrent access during metric ingestion and summary reads.
Observed Implementation:
- `sync.RWMutex` protects all bucket mutations and summary reads.
- `Summary` acquires full write lock `w.mu.Lock()` because it calls `evictStaleLocked` which mutates the slice `w.buckets`.
- `go test -race ./...` passed with zero race conditions detected.
Assessment: PASS
Severity: LOW
Notes: Thread safe under multi-goroutine access.
