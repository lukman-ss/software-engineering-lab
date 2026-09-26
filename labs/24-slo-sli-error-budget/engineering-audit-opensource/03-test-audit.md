# Test Audit

## Coverage Measurement

Command: `go test -coverprofile=/tmp/cover.out -coverpkg=labs/24-slo-sli-error-budget/internal/... ./...`
Result: 94.1% of internal/... statements covered (NOT the claimed 100% per engineering/01-design.md).

Per-function coverage:
- internal/alerting/engine.go: NewAlertEngine 100%, CalculateBurnRate 71.4%, Check 100%
- internal/metrics/tracker.go: NewWindowTracker 66.7%, Record 96.9%, evictStaleLocked 100%, Summary 100%
- internal/slo/evaluator.go: NewEvaluator 100%, Evaluate 100%

Uncovered branches:
- CalculateBurnRate: `allowedErrorRate <= 0` branch (targetSLO >= 1.0) — never tested.
- NewWindowTracker: default-bucket-size fallback (`bucketSize <= 0`) and `numBuckets < 1` fallback — never tested.

## Test Cases Reviewed (tests/slo_test.go)

| Test | Area | Assessment |
|------|------|-----------|
| TestMetricsWindowTracker | happy path + eviction + latency breach | PASS |
| TestSLOEvaluator | budget math, SLI ratio, CanDeploy on exhaustion | PASS |
| TestAlertEngineBurnRate | happy path + transient-spike negative case | PASS |
| TestOutOfOrderTimestamps | ordering edge case | PASS |
| TestEvaluatorZeroTraffic | zero-traffic edge case | PASS |
| TestConcurrencyMetrics | concurrent writes, race detector | PASS |

## Gap Analysis

- happy path: covered
- failure path: covered (budget exhaustion)
- edge cases: partially covered (zero traffic, out-of-order). Missing: all-events-bad (100% error), exact-zero-budget boundary, default-bucket fallback.
- transitions: covered (window expiry/eviction)
- recovery: NOT covered — no test for budget recovery or release policy re-enabling. The design/implementation notes mention "recovery" as demonstrated, but no test asserts it.
- rollback: NOT applicable (in-memory only); not tested.
- concurrency: covered and race-clean.
- negative cases: covered (transient spike).

## Assessment
Tests are strong for the core math and concurrency. Two uncovered branches in alerting/metrics constructors. No test exercises recovery. Coverage claim of 100% is false; actual is 94.1%.

## Missing Test Cases
- MISSING_TEST: recovery / budget-replenishment (re-population after incident)
- MISSING_EDGE_CASE: 100% error rate (all-bad)
- MISSING_EDGE_CASE: exact zero budget boundary (budgetRemaining == 0)
- MISSING_EDGE_CASE: NewWindowTracker default-param fallbacks
- MISSING_EDGE_CASE: CalculateBurnRate targetSLO == 1.0 (allowedErrorRate == 0)