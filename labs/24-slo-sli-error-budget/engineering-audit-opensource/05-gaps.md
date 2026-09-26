# Gap Analysis

## Identified Gaps

| # | Gap Type | Location | Severity |
|---|----------|----------|----------|
| 1 | UNHANDLED_ERROR | internal/slo/evaluator.go:55 | MEDIUM |
| 2 | DOC_CODE_MISMATCH | cmd/demo/main.go:24 | LOW |
| 3 | UNVERIFIED_RESULT | tests/slo_test.go:81-83 | MEDIUM |
| 4 | MISSING_EDGE_CASE | internal/slo/evaluator.go | LOW |
| 5 | MISSING_EDGE_CASE | internal/metrics/tracker.go:30-34 | LOW |
| 6 | MISSING_EDGE_CASE | internal/metrics/tracker.go:4-13 | LOW |
| 7 | MISSING_EDGE_CASE | internal/alerting/engine.go | LOW |
| 8 | UNHANDLED_ERROR | internal/slo/evaluator.go:11-14 | LOW |
| 9 | RACE_CONDITION | internal/metrics/tracker.go:89-100 | LOW (latent) |

## Detailed Gap Descriptions

### Gap 1: UNHANDLED_ERROR — CanDeploy boundary relies on floating-point precision
**Severity**: MEDIUM
**Location**: internal/slo/evaluator.go:55
**Description**: The CanDeploy boundary check `budgetRemaining <= 0` is defeated by IEEE 754 floating-point arithmetic. When the budget is exactly consumed, `budgetRemaining` is a tiny positive number (~8.88e-16) instead of 0.0. The test passes by accident, but if floating-point semantics change or if `TargetUptime` is set differently, the boundary behavior could silently flip. The code does not explicitly handle the exact-zero case.

### Gap 2: DOC_CODE_MISMATCH — Misleading variable name `window30d`
**Severity**: LOW
**Location**: cmd/demo/main.go:24
**Description**: Variable `window30d` is assigned `30 * time.Minute`, not 30 days. The name contradicts its assigned value, which could mislead readers about the actual window being simulated.

### Gap 3: UNVERIFIED_RESULT — CanDeploy boundary test relies on floating-point quirk
**Severity**: MEDIUM
**Location**: tests/slo_test.go:81-83
**Description**: The test `TestSLOEvaluator` expects `CanDeploy=true` when SLI equals 99% (the SLO target). This test only passes because `1.0 - 0.99 = 0.010000000000000009` in Go, making the error budget slightly larger than the consumed budget. If exact decimal arithmetic were used, the test would FAIL (budgetRemaining would be exactly 0, triggering CanDeploy=false). The test does not prove the intended behavior — it accidentally proves the floating-point behavior.

### Gap 4: MISSING_EDGE_CASE — No test for zero-event tracker
**Severity**: LOW
**Location**: internal/slo/evaluator.go, internal/metrics/tracker.go
**Description**: Neither the evaluator nor the tracker tests cover the case where no events have been recorded. `Summary()` returns (0,0,0), `Evaluate` sets SLI=1.0 (default), budgetRemaining=0.0, and CanDeploy=true (since `total > 0` is false). This edge case is untested.

### Gap 5: MISSING_EDGE_CASE — No validation of negative window size
**Severity**: LOW
**Location**: internal/metrics/tracker.go:30-34
**Description**: `NewWindowTracker` validates `bucketSize <= 0` but does not validate `windowSize < 0`. A negative `windowSize` would cause `evictStaleLocked` to compute a positive cutoff offset, preventing any eviction and causing unbounded memory growth.

### Gap 6: MISSING_EDGE_CASE — No out-of-order event handling test
**Severity**: LOW
**Location**: internal/metrics/tracker.go:4-13
**Description**: The tracker does not support out-of-order events. There is no test validating this limitation, and no documentation (code comment) explaining the assumption.

### Gap 7: MISSING_EDGE_CASE — No test for burn rate < threshold
**Severity**: LOW
**Location**: internal/alerting/engine.go
**Description**: `TestAlertEngineBurnRate` only tests the triggering case (burn rate > threshold). There is no test verifying that no alerts fire when burn rate is below the threshold, or when short/long windows give different results.

### Gap 8: UNHANDLED_ERROR — LatencyThreshold config field is unused
**Severity**: LOW
**Location**: internal/slo/evaluator.go:11-14
**Description**: The `Config.LatencyThreshold` field is defined but never read. The evaluator delegates good/bad classification to the tracker's `isGood` callback. This is a dead field that creates a misleading API surface.

### Gap 9: RACE_CONDITION — Latent race in Summary's implicit write
**Severity**: LOW (latent)
**Location**: internal/metrics/tracker.go:89-100
**Description**: `Summary()` uses `w.mu.Lock()` (write lock) but performs `evictStaleLocked(now)` which reassigns the `buckets` slice. Under concurrent access, `Record()` calls `evictStaleLocked(e.Timestamp)` with the event's own timestamp, while `Summary()` calls it with `now`. If two calls interleave, the `buckets` slice reassignment is atomic under the mutex, so there's no data race. However, the use of a write lock for what appears to be a read operation is a design smell. The `-race` detector does not flag this as a true data race (the mutex protects the slice), but the semantics are confusing.
