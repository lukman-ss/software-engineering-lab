# Engineering Audit — Test Audit

Target: labs/24-slo-sli-error-budget

## Tests Reviewed

tests/slo_test.go — 6 tests.

### TestMetricsWindowTracker

Covers: happy-path recording, eviction of stale buckets, good/bad/latency classification.
PASS. Edge: total=12, good=10, bad=2 confirmed. Future=20s evicts all.

### TestSLOEvaluator

Covers: happy path (99%) + boundary (budget exhausted flips CanDeploy).
PASS. Good: exercises freeze gate.

### TestAlertEngineBurnRate

Covers: happy alert trigger (2/100 -> 20x burn -> alert); negative transient spike rejection (short 10% / long 0.01% -> no alert).
PASS. Good: proves transient-spike suppression.

### TestOutOfOrderTimestamps

Covers: out-of-order recording and ordered eviction.
PASS.

### TestEvaluatorZeroTraffic

Covers: zero-traffic → SLI=1, CanDeploy=true.
PASS. Good edge case.

### TestConcurrencyMetrics

Covers: 20 goroutines x100 reqs concurrently → total exact, good+bad=total.
PASS. Race-free.

## Coverage Matrix

| Category       | Test                          | Status |
|----------------|-------------------------------|--------|
| Happy path     | TestMetricsWindowTracker       | PASS   |
| Happy path     | TestSLOEvaluator               | PASS   |
| Happy path     | TestAlertEngineBurnRate        | PASS   |
| Failure path   | TestAlertEngineBurnRate (neg)   | PASS   |
| Edge (zero)    | TestEvaluatorZeroTraffic       | PASS   |
| Edge (100% err)| not explicitly                | MISS   |
| Transitions    | CanDeploy flip                | PASS   |
| Eviction       | WindowTracker, OutOfOrder     | PASS   |
| Concurrency    | TestConcurrencyMetrics          | PASS (race clean) |

## Observed Command Results

$ go test -count=1 -v ./...
... all 6 PASS (0.084s)

$ go test -count=1 -race ./...
... ok (race detector clean, 1.114s)

## Gaps

- MISSING_TEST: No 100%-error edge case (all bad requests) for Evaluator/AlertEngine. Not present in tests.
- MISSING_TEST: No latency-only-bad case (slow-but-200 responses); isGood in tests only checks StatusCode. Latency-threshold logic in demo's isGood is not exercised by unit tests.

Severity: LOW for both (behavior exists and is covered by demo, but unit-level missing).
