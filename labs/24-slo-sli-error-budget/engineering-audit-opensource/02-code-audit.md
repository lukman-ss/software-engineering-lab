# Code Audit

## Finding 1: Floating-point boundary fragility in budget calculation

Location: `internal/slo/evaluator.go:50-57`
Claimed Behavior: Release freeze when error budget is exhausted (budgetRemaining <= 0)
Observed Implementation: BudgetRemaining calculation uses floating-point arithmetic; at exact budget exhaustion, IEEE 754 representation may yield slightly positive value due to floating-point imprecision, causing CanDeploy to remain true when it should be false.
Assessment: WARNING
Severity: MEDIUM
Notes: The test `TestSLOEvaluator` passes for 99% SLO with 1 error out of 100 events because budgetRemaining computes to ~8.9e-16 (> 0). However, if the SLO were changed slightly (e.g., 0.989999), the boundary could fail. The Status display rounds BudgetRemaining to 2 decimals, potentially showing "0.00" while CanDeploy=true, causing confusion.

## Finding 2: Unused fields in configuration structs

Location: 
- `internal/slo/evaluator.go:10-14` (Config.LatencyThreshold)
- `internal/alerting/engine.go:17-24` (BurnRateRule.LongWindow, ShortWindow, BudgetConsumedPct)
Claimed Behavior: Configuration structs contain fields for latency thresholds and window sizes
Observed Implementation: These fields are declared but never read or used in the implementation. The `isGood` function passed to `WindowTracker` already handles latency checks, and window sizes are determined by the trackers passed to constructors.
Assessment: WARNING
Severity: LOW
Notes: Dead code increases cognitive load without benefit. The LatencyThreshold in SLO Config is particularly misleading as it suggests the evaluator considers latency, but it only uses the tracker's isGood function.

## Finding 3: Docs vs code mismatch on time window compression

Location: 
- Design doc: `engineering/01-design.md:31` ("30-day/rolling windows in compressed real-time (e.g. 1-second = 1-hour scale for demo)")
- Implementation: `cmd/demo/main.go:24` (`window30d := 30 * time.Minute`)
Claimed Behavior: System implements 30-day SLO windows using time compression for demonstration
Observed Implementation: The demo uses literal 30-minute windows with no time compression demonstrated. Events are recorded in real time over ~50 seconds, fitting within the 30-minute window. No scaling or compression of time units is applied.
Assessment: WARNING
Severity: MEDIUM
Notes: This creates a mismatch between documented design and actual implementation. The windows are fixed durations, not compressed representations of longer periods.

## Finding 4: Inefficient bucket insertion logic in WindowTracker

Location: `internal/metrics/tracker.go:66-94`
Claimed Behavior: Efficient insertion of out-of-order timestamp events into time-bucketed sliding window
Observed Implementation: The bucket search/insertion loop creates temporary slices during insertion: `w.buckets = append(w.buckets[:i], append([]Bucket{b}, w.buckets[i:]...)...)`. This generates unnecessary memory allocations and copies.
Assessment: WARNING
Severity: LOW
Notes: While correct and protected by mutex, this pattern is inefficient under high write load. A more efficient approach would use Go 1.21's `slices.Insert` or pre-allocate buckets. However, given typical event rates, this is unlikely to be a bottleneck.

## Finding 5: Summary method blocks concurrent reads unnecessarily

Location: `internal/metrics/tracker.go:117-127`
Claimed Behavior: Concurrent-safe access to metrics summary
Observed Implementation: `Summary()` uses `w.mu.Lock()` (exclusive lock) instead of `w.mu.RLock()` (read lock) because it calls `evictStaleLocked()` which mutates state via bucket eviction.
Assessment: WARNING
Severity: LOW
Notes: Since eviction only removes old data and doesn't affect summary counts for current window, the mutex could be released after eviction before summing, allowing concurrent reads. Current implementation blocks all readers during summary, reducing throughput under mixed read/write workloads.

## Finding 6: Missing input validation on SLO target

Location: `internal/slo/evaluator.go:10-14, 34-39`
Claimed Behavior: SLO target represents a validity probability (0.0 to 1.0)
Observed Implementation: No validation of `Config.TargetUptime` allows values outside [0,1] range, leading to nonsensical calculations (e.g., negative error budgets for TargetUptime > 1.0).
Assessment: WARNING
Severity: LOW
Notes: While the demo and tests use valid values, the lack of validation could cause silent logic errors if invalid configuration is provided. The alerting engine handles `allowedErrorRate <= 0` gracefully, but the evaluator does not.

## Finding 7: Correctness verified under concurrency and edge cases

Location: Multiple files (verified via testing)
Claimed Behavior: Thread-safe metric recording, correct mathematical calculations, proper alerting logic
Observed Implementation: All tests pass, race detector reports no data races, and manual edge case verification confirms:
- SLO=1.0 boundary works correctly
- 100% error rate handled
- Multiple alert rules trigger appropriately
- Zero-event cases handled
- Partial eviction works as expected
- Concurrent recording from 20 goroutines preserves data integrity
Assessment: PASS
Severity: N/A
Notes: Implementation demonstrates strong correctness under stress conditions. The mutex-protected design ensures data integrity despite complex bucket manipulation logic.

## Finding 8: Demo output matches expected behavior

Location: `cmd/demo/main.go`
Claimed Behavior: Demo shows baseline traffic, incident-induced budget depletion, burn rate alerting, and endpoint criticality comparison
Observed Implementation: Demo execution matches expected output exactly:
- Phase 1: 1000 good requests → SLI=100.00%, Budget Remaining=1.00, CanDeploy=true
- Phase 2: +100 requests (90 good, 10 bad) → SLI=99.09%, Budget Remaining=-8.90, CanDeploy=false
- Phase 3: Burn rate 9.09x triggers Slow Burn alert (threshold 6.0x)
- Phase 4: Reports service (95% SLO) with same incident shows SLI=90.00%, Budget Remaining=-5.00, CanDeploy=false
Assessment: PASS
Severity: N/A
Notes: Demo correctly illustrates all core concepts: SLI calculation, error budget tracking, release freeze policy, multi-window burn-rate alerting, and SLO differentiation by service criticality.