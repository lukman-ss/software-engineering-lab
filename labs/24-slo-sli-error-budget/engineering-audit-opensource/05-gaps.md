# Gaps — labs/24-slo-sli-error-budget

## GAP-1: IMPLEMENTATION_OVERCLAIM (MEDIUM)
Location: internal/alerting/engine.go:17-24 vs Check 63-89. Rule LongWindow/ShortWindow/BudgetConsumedPct never read; all rules share two construction trackers. Core AND-gating proven; per-rule window semantics absent.

## GAP-2: DOC_CODE_MISMATCH (LOW)
Location: internal/slo/evaluator.go:10-14. Config.LatencyThreshold never read by Evaluator; latency judged by caller isGood closure only.

## GAP-3: MISSING_TEST / MISSING_EDGE_CASE (MEDIUM)
No tests for: TargetUptime 0/>1/negative, nil isGoodEvent, nil tracker (all panic/nonsense unguarded); burn guards total==0 and SLO==1.0 unasserted; 100%-errors case absent; far-future timestamp wipe untested; no min-sample guard (1/1 → 1000x pages). "100% coverage" claim unproven.

## GAP-4: UNHANDLED_ERROR (LOW)
Nil tracker / nil isGoodEvent panic; invalid TargetUptime silent. Defensive-validation gap, no observed failure in tested paths.

## GAP-5: IMPLEMENTATION_OVERCLAIM (LOW)
Demo Phase 4 feeds identical 10% errors to 99.9% and 95% SLOs; both freeze. Budget math per SLO correct; differentiation scenario unproven.

No RACE_CONDITION. No FAKE_DEMO. No FAKE_BENCHMARK. No BROKEN_IMPLEMENTATION. No HIGH/CRITICAL gaps.
