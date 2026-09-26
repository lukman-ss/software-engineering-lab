# Code Audit

## Finding 1: Summary() uses write lock instead of read lock

Location: internal/metrics/tracker.go:89-91
Claimed Behavior: Summary should be a read-only operation using RLock
Observed Implementation:
```go
func (w *WindowTracker) Summary(now time.Time) (total int64, good int64, bad int64) {
    w.mu.Lock()
    defer w.mu.Unlock()
    w.evictStaleLocked(now)
    ...
}
```
Assessment: WARNING
Severity: LOW
Notes: The `Summary` method calls `evictStaleLocked` which modifies the `buckets` slice by reassigning it. Since it modifies state, using `Lock()` is technically correct. However, this means `Summary` is NOT a pure read — it performs implicit writes (garbage collection of stale buckets). This reduces concurrency. If `evictStaleLocked` were removed from `Summary` (and only called in `Record`), a `RLock` could be used, allowing concurrent reads. The current design works but has suboptimal concurrency. The `RWMutex` is used but only for write-side semantics; read-side locking is never exercised.

## Finding 2: Floating-point boundary fragility in CanDeploy

Location: internal/slo/evaluator.go:49-56
Claimed Behavior: CanDeploy should be true when budget remaining > 0, false when exhausted (<=0)
Observed Implementation:
```go
allowedFailureRate := 1.0 - e.config.TargetUptime
totalErrorBudget := allowedFailureRate * float64(total)
budgetConsumed := float64(bad)
budgetRemaining := totalErrorBudget - budgetConsumed

canDeploy := true
if total > 0 && budgetRemaining <= 0 {
    canDeploy = false
}
```
Assessment: WARNING
Severity: MEDIUM
Notes: The test `TestSLOEvaluator` expects `CanDeploy=true` when SLI exactly equals target SLO (99% SLI, 99% target, 1 error out of 100). Mathematically, `budgetRemaining` should be exactly 0.0, which would trigger `budgetRemaining <= 0` → CanDeploy=false. However, due to IEEE 754 floating-point arithmetic, `1.0 - 0.99` yields `0.010000000000000009` (not exactly 0.01), making `totalErrorBudget` slightly exceed 1.0, and `budgetRemaining` a tiny positive value (~8.88e-16). Thus `budgetRemaining <= 0` evaluates to false, and CanDeploy stays true. The test passes "by accident" of floating-point precision. If the arithmetic were exact, the test would fail. This is fragile: any future refactoring (e.g., using exact decimal arithmetic or different target values) could break this behavior silently.

## Finding 3: No handling for events recorded out of order

Location: internal/metrics/tracker.go:46-76
Claimed Behavior: Tracker should handle events in any timestamp order
Observed Implementation: `Record` checks only if the new event matches the LAST bucket (index n-1). If an out-of-order event arrives that falls into a non-last bucket, it creates a new duplicate bucket.
Assessment: WARNING
Severity: LOW
Notes: The tracker assumes monotonically increasing event timestamps. Out-of-order events would create duplicate buckets and inflate counts. While this is consistent with typical use cases (recording events as they occur), there is no validation or handling. This is acceptable for the stated scope but is a limitation.

## Finding 4: Missing validation for negative or zero window/bucket sizes

Location: internal/metrics/tracker.go:30-44
Claimed Behavior: Constructor should validate inputs
Observed Implementation: `NewWindowTracker` only guards against `bucketSize <= 0` (defaulting to 1 second) and `numBuckets < 1` (defaulting to 1). No validation for negative `windowSize`.
Assessment: WARNING
Severity: LOW
Notes: If a caller passes negative `windowSize`, `evictStaleLocked` computes `cutoff := now.Add(-w.windowSize)`. A negative windowSize with negation produces a positive offset, meaning ALL buckets would appear "fresh" (never evicted). This could cause unbounded memory growth. The constructor should reject negative window sizes.

## Finding 5: Demo variable naming inconsistency

Location: cmd/demo/main.go:24
Claimmed Behavior: Variable names should be accurate
Observed Implementation:
```go
window30d := 30 * time.Minute
```
Assessment: WARNING
Severity: LOW
Notes: `window30d` implies a 30-day window, but is set to 30 minutes. This is misleading and inconsistent with the README claim of `internal/metrics` using sliding-window tracking. The name does not reflect reality. Should be renamed to `window30m`.

## Finding 6: LatencyThreshold in Config is unused

Location: internal/slo/evaluator.go:11-14
Claimed Behavior: Config should use LatencyThreshold
Observed Implementation:
```go
type Config struct {
    Name            string
    TargetUptime    float64
    LatencyThreshold time.Duration
}
```
The `LatencyThreshold` field in `Config` is never read by `Evaluator.Evaluate`. Good/bad determination is delegated to the `metrics.Event` -> `isGood` callback in `WindowTracker`. The `LatencyThreshold` is passed to `Config` but not used.
Assessment: WARNING
Severity: LOW
Notes: The field is defined but dead. Either the evaluator should use it for latency-based SLI calculation, or it should be removed. The demo passes `latencyThreshold` in the config but it has no effect on the evaluator's logic — the `isGood` function in the demo does the latency check instead.

## Finding 7: No error handling or panic recovery

Location: All implementation files
Claimed Behavior: The code should handle edge cases gracefully
Observed Implementation: No input validation, no nil checks, no panic recovery. E.g., `NewEvaluator` does not check if `tracker` is nil.
Assessment: PASS (with caveats)
Severity: LOW
Notes: The absence of defensive checks is consistent with Go idioms for internal packages where invariants are maintained at the call site. For this lab scope, this is acceptable. Not a defect but a design choice.
