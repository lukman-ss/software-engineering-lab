# Gap Analysis — labs/24-slo-sli-error-budget

## MISSING_EDGE_CASE (LOW)

internal/slo/evaluator.go:62 rounding uses math.Round(sli*10000)/10000. At exactly 0.99995 boundary, rounds to 1.0. No test at boundary.

## MISSING_TEST (LOW)

- Recovery path: no test for CanDeploy returning to true after stale buckets evict.
- Alert clearing: no test for alerts disappearing once burn rate drops below threshold.
- 100%-error edge: 0 good events → SLI=0, budget=0; not asserted explicitly.

## UNHANDLED_ERROR (MEDIUM)

internal/metrics/tracker.go: no nil-isGoodEvent guard → Record or Summary panics on nil.

internal/metrics/tracker.go: NewWindowTracker does not validate windowSize; negative/zero windowSize silently treated as 1 bucket, disabling sliding window.

## IMPLEMENTATION_OVERCLAIM (MEDIUM)

internal/alerting/engine.go: BurnRateRule.LongWindow / ShortWindow fields declared but unused (implementation uses engine-level trackers). Docs imply per-rule windows.

internal/slo/evaluator.go: Config LatencyThreshold field never read by evaluator; SLO is configured via tracker isGood closure only.

engineering/01-design.md: claims "histogram latency buckets" in metrics — code stores counts only, latency encoded via isGood predicate.

## DOC_CODE_MISMATCH (LOW)

engineering/01-design.md §Architecture: "histogram latency buckets" vs code: none.

engineering/03-execution-result.md Phase 4 caption "wider 5% error tolerance" while both services show CanDeploy=false and Reports below 95%.

engineering/03-execution-result.md §Tests: records PASS; verified independently.

engineering/03-execution-result.md §Race Detector: records PASS; verified independently.

## SUMMARY

Core functionality proven via reproducible build, passing race-clean test suite, and demo output reproduced verbatim. Gaps are scope-overclaim, missing recovery/clear tests (non-blocking), and nil/windowSize guard.
