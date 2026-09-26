# Test Audit

Target Lab: `labs/24-slo-sli-error-budget`

## Test Execution Results

Command executed:
```bash
go test -v -count=1 ./tests
```
Output:
```text
=== RUN   TestMetricsWindowTracker
--- PASS: TestMetricsWindowTracker (0.00s)
=== RUN   TestSLOEvaluator
--- PASS: TestSLOEvaluator (0.00s)
=== RUN   TestAlertEngineBurnRate
--- PASS: TestAlertEngineBurnRate (0.00s)
=== RUN   TestConcurrencyMetrics
--- PASS: TestConcurrencyMetrics (0.00s)
PASS
ok  	labs/24-slo-sli-error-budget/tests	0.205s
```

Race detector command:
```bash
go test -race -count=1 ./tests
```
Output:
```text
PASS
ok  	labs/24-slo-sli-error-budget/tests	0.448s
```

## Coverage of Test Scenarios

### 1. Happy Path
- `TestMetricsWindowTracker`: Verifies recording of successful events and correct total/good/bad counts within window.
- `TestSLOEvaluator`: Verifies 99% SLI calculation with 99 good and 1 bad event.
- Status: COVERED (PASS)

### 2. Failure Path & Policy Enforcement
- `TestSLOEvaluator`: Injects an additional bad event to deplete budget below SLO, asserting `CanDeploy == false`.
- `TestAlertEngineBurnRate`: Verifies 2% error rate against 99.9% SLO triggers a 20x burn rate alert with `SeverityPage`.
- Status: COVERED (PASS)

### 3. Edge Cases & Eviction
- `TestMetricsWindowTracker`: Tests advancing timestamp past window (`now.Add(20 * time.Second)`), verifying that all stale buckets are evicted and total/good/bad counts return to 0.
- Missing: Explicit test verifying behaviour when `total == 0` on `Evaluator.Evaluate()` (though code inspection shows it returns 1.0 and `CanDeploy=true`).
- Status: COVERED WITH MINOR GAP (PASS)

### 4. Concurrency Safety
- `TestConcurrencyMetrics`: Spawns 20 goroutines running 100 requests each concurrently writing to the `WindowTracker` with alternating status codes.
- Asserts that `total == 2000` and `good + bad == total` with zero data races detected under `-race`.
- Status: COVERED (PASS)

### 5. Multi-Window Alerting Thresholds
- `TestAlertEngineBurnRate`: Covers threshold exceeding condition.
- Missing: Negative test asserting that an alert is NOT triggered when only short window exceeds threshold while long window does not.
- Status: PARTIAL (WARNING)
