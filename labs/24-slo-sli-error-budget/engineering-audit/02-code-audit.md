# Code Audit

Target Lab: `labs/24-slo-sli-error-budget`

## Finding 1: WindowTracker Bucket Append Assumption on Timestamp Monotonicity

Location: `internal/metrics/tracker.go:55-75`
Claimed Behavior: Events are recorded into time buckets and aggregated across rolling windows.
Observed Implementation:
In `Record(e Event)`:
```go
n := len(w.buckets)
if n > 0 && w.buckets[n-1].StartTime.Equal(bucketStart) {
    w.buckets[n-1].TotalCount++
    ...
    return
}
b := Bucket{ StartTime: bucketStart, ... }
w.buckets = append(w.buckets, b)
```
If events arrive out of order (e.g., an earlier timestamp arrives after a later timestamp), new buckets are appended at the end with a timestamp before the tail bucket. In `evictStaleLocked`, eviction assumes buckets are sorted by `StartTime` and slices `w.buckets[idx:]`. Out-of-order events could cause stale buckets to persist or premature eviction.
Assessment: WARNING
Severity: LOW
Notes: In typical streaming logging or the test/demo harness, timestamps are monotonic. However, for a general sliding window tracker, out-of-order arrivals should either be inserted at the sorted position or rejected. In this lab context, timestamps are generated sequentially.

## Finding 2: Division by Zero and Edge Case Handling in Evaluator & AlertEngine

Location: `internal/slo/evaluator.go:44-57`, `internal/alerting/engine.go:51-61`
Claimed Behavior: Correct computation of SLI, Error Budget, and Burn Rates under empty/zero traffic conditions.
Observed Implementation:
- In `evaluator.go`:
  ```go
  var sli float64 = 1.0
  if total > 0 {
      sli = float64(good) / float64(total)
  }
  ...
  canDeploy := true
  if total > 0 && budgetRemaining <= 0 {
      canDeploy = false
  }
  ```
  Safely defaults SLI to 1.0 and `canDeploy = true` when total == 0.
- In `engine.go`:
  ```go
  if total == 0 {
      return 0.0
  }
  if allowedErrorRate <= 0 {
      return 0.0
  }
  ```
  Safely returns 0.0 when total is 0 or targetSLO is 1.0 (allowedErrorRate <= 0).
Assessment: PASS
Severity: LOW
Notes: No panics or NaN values possible on zero traffic.

## Finding 3: Concurrency Safety of WindowTracker

Location: `internal/metrics/tracker.go:46-49`, `internal/metrics/tracker.go:89-92`
Claimed Behavior: Thread-safe metric collection across multiple concurrent goroutines.
Observed Implementation:
`Record` and `Summary` both acquire `w.mu.Lock()` (exclusive lock) and release via `defer w.mu.Unlock()`. Although `w.mu` is defined as `sync.RWMutex`, `Summary` mutates `w.buckets` via `evictStaleLocked(now)`, so taking a write lock (`w.mu.Lock()`) in `Summary` is necessary and correctly implemented.
Assessment: PASS
Severity: LOW
Notes: Verified with `go test -race ./...`. No race condition detected.

## Finding 4: Error Budget Math and Deployment Gate Enforcement

Location: `internal/slo/evaluator.go:49-57`
Claimed Behavior: Error budget equals allowed failure rate times total events; releases are blocked when budget is exhausted (`budgetRemaining <= 0`).
Observed Implementation:
```go
allowedFailureRate := 1.0 - e.config.TargetUptime
totalErrorBudget := allowedFailureRate * float64(total)
budgetConsumed := float64(bad)
budgetRemaining := totalErrorBudget - budgetConsumed
```
When total = 1100, target = 0.999 (0.1% budget), `totalErrorBudget = 1.10`.
With 10 bad events, `budgetConsumed = 10`, `budgetRemaining = 1.10 - 10 = -8.90`.
Since `budgetRemaining <= 0`, `canDeploy` evaluates to `false`.
Assessment: PASS
Severity: LOW
Notes: Conforms precisely to Google SRE Error Budget definitions.

## Finding 5: AlertEngine Multi-Window Burn Rate Logic

Location: `internal/alerting/engine.go:63-88`
Claimed Behavior: Multi-window alerting requires both short window and long window burn rates to meet or exceed the threshold factor before triggering an alert.
Observed Implementation:
```go
if shortBurn >= rule.BurnRateFactor && longBurn >= rule.BurnRateFactor {
    triggered = true
}
```
Requires both conditions to be true, preventing transient spikes (short window only) and historical resets (long window only) from triggering alerts erroneously.
Assessment: PASS
Severity: LOW
Notes: Exactly matches Google SRE Workbook Chapter 5 specification for multi-window burn rate alerts.
