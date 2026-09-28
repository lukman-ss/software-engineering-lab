# Gap Analysis

## Gap 1
Type: MISSING_TEST
Location: tests/slo_test.go (no recovery test)
Description: No test verifies error budget recovery (CanDeploy false -> true) after window expiry or failure rate drop.
Severity: MEDIUM

## Gap 2
Type: MISSING_TEST
Location: internal/alerting/engine.go:51 (CalculateBurnRate degenerate guards)
Description: total==0 and allowedErrorRate<=0 early-return branches uncovered (function at 71.4%). No test for SLO=1.0 or empty-window burn.
Severity: LOW

## Gap 3
Type: DOC_CODE_MISMATCH
Location: engineering/01-design.md (recovery claim) vs cmd/demo/main.go
Description: Design promises demo illustrates recovery. Demo implements baseline, incident, alerting, criticality comparison only. No recovery phase.
Severity: MEDIUM

## Gap 4
Type: IMPLEMENTATION_OVERCLAIM
Location: engineering/01-design.md (100% coverage claim)
Description: Doc states 100% test coverage on core math and sliding window. Real coverage via -coverpkg: evaluator 100%, metrics Record ~97%, CalculateBurnRate ~71%. Defensive branches uncovered. Math paths covered, blanket claim overstated.
Severity: LOW

## Gap 5
Type: DOC_CODE_MISMATCH
Location: internal/slo/evaluator.go:15 (Config.LatencyThreshold)
Description: LatencyThreshold present in Config but unused by Evaluator; classification owned by tracker's isGood callback. Doc/design imply evaluator enforces latency.
Severity: LOW

## Gap 6
Type: DOC_CODE_MISMATCH
Location: internal/alerting/engine.go:17,26 (BurnRateRule windows + BudgetConsumedPct)
Description: Rule fields LongWindow, ShortWindow, BudgetConsumedPct never referenced. Multi-window behaviour comes from separate short/long tracker instances, not rule-level window pairs. Dead fields risk future misuse.
Severity: LOW

## Gap 7
Type: MISSING_EDGE_CASE
Location: tests/slo_test.go (no 100%-error test)
Description: No test for all-bad traffic (SLI=0, budget deeply negative). Evaluator handles it mathematically, but unproven by test.
Severity: LOW

No BROKEN_IMPLEMENTATION, RACE_CONDITION, UNHANDLED_ERROR, FAKE_DEMO, FAKE_BENCHMARK, or UNVERIFIED_RESULT found. Recorded execution results match actual runs verbatim.