# Gap Analysis — labs/24-slo-sli-error-budget

## Identified Gaps

### 1. DOC_CODE_MISMATCH (execution-result.md demo record omits Phase 4)
Location: engineering/03-execution-result.md lines 41-65 (recorded demo output)
Gap Type: DOC_CODE_MISMATCH
Description: The recorded demo output omits Phase 4 entirely ("Endpoint Criticality Comparison..."). The actual `go run ./cmd/demo` prints 4 phases; the engineering doc shows only up to Phase 3 then "DEMO COMPLETE".
Severity: MEDIUM
Notes: Not a fake result; the Phase 3 lines match byte-for-byte. The gap is one of stale documentation (the execution result was recorded before Phase 4 was added, or the capture was truncated).

### 2. DOC_CODE_MISMATCH (execution-result.md test list omits 2 of 6 tests)
Location: engineering/03-execution-result.md lines 10-28 (Tests section)
Gap Type: DOC_CODE_MISMATCH
Description: The recorded test output lists only:
   TestMetricsWindowTracker
   TestSLOEvaluator
   TestAlertEngineBurnRate
   TestConcurrencyMetrics
Actual `go test -v` runs 6 tests; missing TestOutOfOrderTimestamps and TestEvaluatorZeroTraffic from the record.
Severity: LOW
Notes: Again, not a fabrication; the listed tests do pass. The record is incomplete.

### 3. BROKEN_IMPLEMENTATION? No — none of the core claims are broken.
- SLI math correct (verified by demo Phase 1-2 and test).
- Error budget depletion correct (verified by demo Phase 2: -8.90 = (1-0.001)*1100 - 10).
- Burn rate alert triggers correct (verified by demo Phase 3: 9.09x > 6x TICKET; code uses `shortBurn >= factor && longBurn >= factor`).
- No race conditions (`go test -race` passes).

### 4. DEAD_CODE (engineering: BurnRateRule fields never used)
Location: internal/alerting/engine.go lines 17-24 (BurnRateRule fields LongWindow, ShortWindow, BudgetConsumedPct) and Check (lines 71-87)
Gap Type: UNUSED_FIELD / DOC_CODE_MISMATCH (design claimed per-rule windows)
Description: BurnRateRule defines LongWindow, ShortWindow, and BudgetConsumedPct but they are never read. The AlertEngine uses a single pair of shared trackers (constructed at NewAlertEngine time) and applies the SAME BurnRateFactor AND condition to all rules. The per-rule window and budget-consumed gating advertised by the field names is dead code.
Severity: LOW
Notes: Not a correctness bug; it is a deviation from the documented flexibility. Each rule could have had its own window pair + BudgetConsumedPct threshold; today all rules inherit the engine-level windows and use identical short/long thresholds.

### 5. CAPACITY_GAP (WindowTracker unbounded growth on out-of-order)
Location: internal/metrics/tracker.go Record (lines 67-92) and NewWindowTracker (lines 30-44)
Gap Type: MISSING_EDGE_CASE / BROKEN_IMPLEMENTATION under adversarial input
Description: NewWindowTracker computes `numBuckets` as windowSize/bucketSize but only uses it as a slice capacity hint, not a hard limit. The Record function will append a new bucket for every distinct out-of-order bucketStart it sees (if bucketStart.Before(last.bucketStart) and no exact match). Under a workload with widely scattered old timestamps, the slice can grow without bound (no eviction of future/past buckets beyond the simple time-window evictStaleLocked). The data structure is not a ring buffer despite comments.
Severity: MEDIUM
Notes: Under the tested monotonic input (all timestamps increasing) eviction keeps the slice bounded. No test exercises adversarial out-of-order streams. A malicious or mis-configured producer could cause memory growth. The documented "ring/time-bucketed window tracking" promise is not fully honored.

### 6. DOC_CODE_MISMATCH (design.md claims "histogram latency buckets")
Location: engineering/01-design.md line 26
Gap Type: DOC_CODE_MISMATCH
Description: Design claims internal/metrics is a "Sliding-window event recorder (histogram latency buckets & success counts)." The implementation records only per-bucket success/total counts; latency is only used as a predicate in the isGood closure (e.g. `e.Duration <= latencyThreshold`). No histogram of latency values is accumulated.
Severity: LOW
Notes: The core SLO/SLI math does not require latency histograms; the lab still demonstrates the latency-threshold-gated good/bad split. This is an aspiration in the design not realized in code.

### 7. TEST_GAP (missing recovery/rollback test)
Location: tests/slo_test.go
Gap Type: MISSING_TEST
Description: No test asserts that once CanDeploy becomes false (budget exhausted), injecting sufficient good traffic can restore the budget and flip CanDeploy back to true. The design's "success criteria" listed "rollback" as something to prove, but no test covers budget recovery.
Severity: LOW
Notes: The monotonic demo never restores budget; tests are also one-way traffic only. Not a correctness issue today, but an incomplete proof of the claimed release-freeze policy behavior.

### 8. TEST_GAP (missing SLO=1.0 / zero-burn edge case)
Location: internal/alerting/engine.go CalculateBurnRate (lines 51-61)
Gap Type: MISSING_TEST
Description: If targetSLO == 1.0 (allowedErrorRate == 0.0), CalculateBurnRate returns 0 to avoid division by zero. No test asserts this guard (e.g. targetSLO=1.0, some errors => burn 0.0, no alert despite errors).
Severity: LOW
Notes: Low risk because SLO=1.0 is unrealistic; but the code path is untested.

## Summary of Gap Severity

LOW: 
- Dead fields in BurnRateRule (design flexibility not implemented)
- Missing histogram latency buckets (design aspiration)
- Missing SLO=1.0 test
- Missing recovery/rollback test
- Incomplete execution-result.md test list

MEDIUM:
- Incomplete execution-result.md demo record (omitted Phase 4)
- WindowTracker unbounded growth on adversarial out-of-order input

HIGH: 
- (none)

CRITICAL: 
- (none)

## Final Note
All core behavior (SLI, error budget, burn rate alerting) is proven by re-runnable demo and test suite. The gaps are documentation staleness, dead fields, missing edge-case hardening, and incomplete test coverage — none invalidate the demonstrated correctness.