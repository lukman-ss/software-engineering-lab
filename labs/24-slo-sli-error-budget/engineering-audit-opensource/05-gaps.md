# Gaps — labs/24-slo-sli-error-budget

## GAP-1: MISSING_TEST (MEDIUM)
Negative burn-rate case untested. No test asserts `Check` returns zero alerts when burn < threshold, or that 14.4x PAGE stays silent at 9.09x. Add: below-threshold scenario expecting `len(alerts)==0`.

## GAP-2: MISSING_TEST (MEDIUM)
Empty-window and all-error evaluator paths untested. SLI==1.0/CanDeploy==true on zero events and SLI==0.0 freeze on 100% errors are implemented but unproven. Add both.

## GAP-3: MISSING_EDGE_CASE (LOW)
`CalculateBurnRate` total==0 and targetSLO>=1.0 guards implemented, uncovered (71.4%). Add two asserts.

## GAP-4: MISSING_EDGE_CASE (LOW)
`NewWindowTracker` defaults (`bucketSize<=0`, `window<bucket`) uncovered (66.7%). Add constructor asserts.

## GAP-5: MISSING_TEST (LOW)
No recovery test: events aging out should restore budget and CanDeploy=true. Implemented via eviction, never exercised end-to-end through Evaluator. Add: exhaust then advance past window, re-evaluate.

## GAP-6: DOC_CODE_MISMATCH (MEDIUM)
"100% test coverage on core math" overclaims; measured <100% on two functions. Fix wording or add GAP-1–4 tests.

## GAP-7: DOC_CODE_MISMATCH (LOW)
Design terms "histogram", "ring buffer", "endpoint criticality bucketing", "recovery demo" exceed implementation. Narrow design wording to bucketed counters, append+evict window, single-SLO demo without recovery phase — or implement the missing bits.

## GAP-8: UNHANDLED_ERROR (LOW)
`NewWindowTracker` nil `isGood` panics on Record; out-of-order timestamps silently misbucket (see 02-code-audit Findings 1–2). Guard or document ordering + non-nil precondition. No production impact in lab scope.

## Explicit non-gaps
No BROKEN_IMPLEMENTATION, RACE_CONDITION, FAKE_DEMO, FAKE_BENCHMARK, UNVERIFIED_RESULT. Build/vet/tests/race/demo all reproduced PASS. No HIGH/CRITICAL gaps.
