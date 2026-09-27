# Test Audit

Target Lab: `labs/24-slo-sli-error-budget`

## Test Suite Overview

Test file: `tests/slo_test.go`
Execution Commands:
- `go test -count=1 -v ./...`
- `go test -count=1 -race ./...`

## Test Case Breakdown

### 1. Happy Path & Window Eviction
- Function: `TestMetricsWindowTracker` (`tests/slo_test.go:13-59`)
- Verified Behavior: Records 10 good events and 2 bad events (1 slow latency, 1 HTTP 500). Verifies total 12, good 10, bad 2. Verifies complete eviction when timestamp advances beyond sliding window.
- Status: PASS

### 2. SLO Evaluation & Release Policy Enforcement
- Function: `TestSLOEvaluator` (`tests/slo_test.go:61-91`)
- Verified Behavior: 99 good, 1 bad on 99% SLO maintains `CanDeploy = true`. Adding 1 more bad event violates budget and asserts `CanDeploy = false`.
- Status: PASS

### 3. Multi-Window Multi-Burn-Rate Alerting & Transient Negative Case
- Function: `TestAlertEngineBurnRate` (`tests/slo_test.go:93-152`)
- Verified Behavior:
  - Both windows elevated (20x burn rate > 14.4x threshold): Fires `SeverityPage` alert.
  - Negative test: Transient spike in short window (100x burn rate) with quiet long window (0.1x burn rate < 14.4x threshold): Verified 0 alerts triggered.
- Status: PASS

### 4. Out-of-Order Timestamp Handling & Sorted Eviction
- Function: `TestOutOfOrderTimestamps` (`tests/slo_test.go:154-177`)
- Verified Behavior: Out-of-order events are inserted into existing or sorted intermediate buckets. Partial eviction drops older bucket while keeping newer bucket intact.
- Status: PASS

### 5. Edge Case: Zero Traffic
- Function: `TestEvaluatorZeroTraffic` (`tests/slo_test.go:179-198`)
- Verified Behavior: Zero traffic returns 0 events, default SLI 1.0, and `CanDeploy = true`.
- Status: PASS

### 6. Concurrency Safety
- Function: `TestConcurrencyMetrics` (`tests/slo_test.go:200-236`)
- Verified Behavior: 20 goroutines x 100 requests (2,000 total events) recorded concurrently. Verified total = 2,000 and good + bad == total with zero data races.
- Status: PASS

## Test Execution Results

```text
=== RUN   TestMetricsWindowTracker
--- PASS: TestMetricsWindowTracker (0.00s)
=== RUN   TestSLOEvaluator
--- PASS: TestSLOEvaluator (0.00s)
=== RUN   TestAlertEngineBurnRate
--- PASS: TestAlertEngineBurnRate (0.00s)
=== RUN   TestOutOfOrderTimestamps
--- PASS: TestOutOfOrderTimestamps (0.00s)
=== RUN   TestEvaluatorZeroTraffic
--- PASS: TestEvaluatorZeroTraffic (0.00s)
=== RUN   TestConcurrencyMetrics
--- PASS: TestConcurrencyMetrics (0.00s)
PASS
ok  	labs/24-slo-sli-error-budget/tests	0.309s
```

Race detector output:
```text
ok  	labs/24-slo-sli-error-budget/tests	1.113s
```
Zero race warnings or memory leaks detected.
