# Test Audit

Target Lab: `labs/24-slo-sli-error-budget`
Test File: `tests/slo_test.go` (6 tests, 236 lines)

## Coverage Matrix

| Requirement | Test | Result |
|---|---|---|
| Happy path: bucket aggregation (10 good + 2 bad) | TestMetricsWindowTracker | PASS |
| Window expiry / eviction to zero | TestMetricsWindowTracker | PASS |
| SLO boundary: 99/100 deploys, 2nd bad freezes | TestSLOEvaluator | PASS |
| Burn-rate fires at 20x > 14.4x PAGE | TestAlertEngineBurnRate | PASS |
| Negative: short-only spike, clean long window → no alert | TestAlertEngineBurnRate | PASS |
| Out-of-order insert + partial eviction | TestOutOfOrderTimestamps | PASS |
| Zero traffic → SLI 1.0, CanDeploy=true | TestEvaluatorZeroTraffic | PASS |
| Concurrency: 20x100 records, no loss, race clean | TestConcurrencyMetrics | PASS |

## Execution (lab dir, actual)

```text
go build ./...              → BUILD_EXIT=0 (no output, no errors)
go test -count=1 -v ./...   → all 6 tests PASS (tests package ok)
go test -count=1 -race ./...→ ok labs/24-slo-sli-error-budget/tests (no races)
go run ./cmd/demo           → matches engineering/03-execution-result.md numerically
```

Note: `./...` only matches packages when run from the lab dir (module-scoped). From outside
(e.g. `scripts/orchestrator`) it correctly matches nothing — not a repo defect.

## Strengths

- Boundary test is exact (budget == 0 still deploys; first negative freezes) — proves the `<= 0` policy.
- Negative alert test is the strongest asset: proves dual-window gating, not just threshold firing.
- Out-of-order + partial-eviction test proves the sorted-insert path, the trickiest code in the tracker.
- Concurrency test uses identical timestamps across goroutines (same-bucket contention) under `-race`.

## Weaknesses

1. No 100%-error test, although `engineering/01-design.md` Test Strategy explicitly lists "100% errors"
   as a planned edge case. Only mixed ratios are exercised. → MISSING_TEST (LOW).
2. `CalculateBurnRate` guards (`total == 0`, `targetSLO == 1.0`) have no direct unit test; the zero case
   is reached only indirectly via `TestEvaluatorZeroTraffic` (evaluator, not burn rate). → MISSING_TEST (LOW).
3. Concurrency test asserts `total == 2000` and `good+bad == total` but not the exact deterministic split
   (1800 good / 200 bad, since `i%10==0` fails). A miscount preserving the total would pass. → weak
   assertion (LOW).
4. No test touches `BurnRateRule.LongWindow/ShortWindow/BudgetConsumedPct` — unsurprising, since the
   engine ignores them (see code Finding 5). Untestable dead fields. → consequence of WARNING (MEDIUM).

## Assessment

Suite is honest and above-average for scope: happy path, failure path, eviction, out-of-order,
zero-traffic, concurrency, and a true negative alert case all execute green under `-race`.
Gaps are additive (missing 100%-error case, weak exact-count assertion), none invalidate what is proven.
A passing suite here is backed by boundary-exact and negative-case tests, not just happy paths.
