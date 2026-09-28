# Gap Analysis

## MISSING_EDGE_CASE
- Location: internal/metrics/tracker.go:109
- Description: No test for bucket start time exactly at cutoff boundary (eviction uses strict Before).
- Severity: LOW

## IMPLEMENTATION_OVERCLAIM
- Location: internal/slo/evaluator.go (Config.LatencyThreshold)
- Description: Config struct advertises LatencyThreshold field that is never used by Evaluator.
- Severity: MEDIUM

## IMPLEMENTATION_OVERCLAIM
- Location: internal/alerting/engine.go (BurnRateRule.LongWindow, ShortWindow, BudgetConsumedPct)
- Description: Struct fields declared but never used in alert logic.
- Severity: LOW

## MISSING_TEST
- Location: internal/metrics/tracker.go
- Description: No test for CalculateBurnRate edge cases (zero total, zero allowed error rate).
- Severity: LOW

## MISSING_TEST
- Location: internal/alerting/engine.go
- Description: No test for empty rule slice, threshold equality, or asymmetric window triggers.
- Severity: LOW