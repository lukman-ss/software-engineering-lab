# Test Audit

## Coverage Summary

Measured independently (fresh run, `-count=1`):

```text
goos: darwin  goarch: arm64  pkg: labs/24-slo-sli-error-budget/tests
coverage: 94.1% of statements in ./internal/...

Function                      Coverage
NewAlertEngine                100.0%
CalculateBurnRate             71.4%
Check                         100.0%
NewWindowTracker              66.7%
Record                        96.9%
evictStaleLocked              100.0%
Summary                       100.0%
NewEvaluator                  100.0%
Evaluate                      100.0%
```

## Test Inventory (`tests/slo_test.go`)

| Test | Target | Path | Pass |
|---|---|---|---|
| `TestMetricsWindowTracker` | happy path + eviction | metrics/tracker.go | PASS |
| `TestSLOEvaluator` | SLO math + budget exhaustion (CanDeploy flip) | slo/evaluator.go | PASS |
| `TestAlertEngineBurnRate` | burn-rate threshold + transient negative case | alerting/engine.go | PASS |
| `TestOutOfOrderTimestamps` | out-of-order insert + partial eviction | metrics/tracker.go | PASS |
| `TestEvaluatorZeroTraffic` | zero-traffic edge case | slo/evaluator.go | PASS |
| `TestConcurrencyMetrics` | concurrent writers + invariant (good+bad==total) | metrics/tracker.go | PASS |

## Coverage of Required Test Categories

- Happy path — PASS (`TestMetricsWindowTracker`, `TestSLOEvaluator`, `TestAlertEngineBurnRate`).
- Failure path — PASS (`TestSLOEvaluator` budget exhaustion → `CanDeploy=false`; `TestAlertEngineBurnRate` transient-negative suppression).
- Edge cases — PASS (zero traffic, out-of-order timestamps, partial eviction).
- Transitions — PASS (release-freeze transition via second bad event in `TestSLOEvaluator`).
- Recovery/rollback — NOT COVERED. The lab models error-budget *depletion* and a release *freeze*, but defines no recovery/replenishment path; there is no test for budget recovering (e.g., good traffic restoring deploy permission). This is a conceptual gap, not a code defect — the design has no recovery mechanism to test.
- Concurrency — PASS (`TestConcurrencyMetrics` + `go test -race` clean).
- Negative cases — PASS (transient spike alert suppression: short window over threshold, long window below threshold → no alert).

## Strengths

- Burn-rate math independently verified: 10 bad / 1100 total / allowed 0.001 = 9.09x; correctly NOT tripping the 14.4x rule, correctly tripping the 6.0x rule.
- The transient negative test (`engineTransient`) meaningfully exercises the multi-window AND condition and proves the short-only spike does not fire.
- Concurrency test uses 20 goroutines × 100 writes on overlapping 100 ms buckets, then asserts the exact total and the `good+bad==total` invariant.

## Weaknesses / Gaps

1. `CalculateBurnRate` at 71.4% — the `total == 0` (zero-division guard) and `allowedErrorRate <= 0` (SLO >= 1.0, i.e., zero error budget) branches are never exercised. These are defensive guards, not demonstrated failure paths, but they are untested.
2. `NewWindowTracker` at 66.7% — the `bucketSize <= 0` fallback and `numBuckets < 1` fallback validation branches are untested. Defensive input-validation defaults.
3. `Record` at 96.9% — one rare branch uncovered (see code audit).
4. The design's "100% test coverage on core math and sliding window calculations" success criterion is **not met**: overall 94.1%. (Core math `Evaluate` is 100%; sliding-window helpers range 96.9–66.7%.) The shortfall is confined to defensive guard branches, not core aggregation logic.
5. No test exercises burn-rate alerting with SLO other than 0.999; the `allowedErrorRate <= 0` guard (SLO == 1.0) and the `CalculateBurnRate(total==0)` path have no test.

## Race Detector

`go test -race ./...` → clean (ok tests 1.386s). No data races.

## Verdict on Test Suite

Assessment: PASS (suite is solid and proves the demonstrated behaviors) but with coverage overclaim. The suite *proves* the core math and the demo's behavior; it does not fully exercise defensive guards. A "passing suite that is weak" — not the case here for proven behavior, but coverage is overstated.
