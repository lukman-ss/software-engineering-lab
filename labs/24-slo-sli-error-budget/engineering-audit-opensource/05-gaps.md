# Gap Analysis — labs/24-slo-sli-error-budget

## Identified Gaps

1. MISSING_TEST — `alerting.CalculateBurnRate`: total==0 and allowedErrorRate<=0 branches return 0.0 but no test asserts. Location: internal/alerting/engine.go:51-61.
2. MISSING_TEST — `alerting.Check`: rule with BurnRateFactor==0 always passes `shortBurn>=0 && longBurn>=0` → guaranteed fire. No guard, no test. Location: engine.go:70-88.
3. MISSING_TEST — `metrics.NewWindowTracker`: bucketSize<=0 guarded, windowSize<=0 not guarded (would zero evict). Nil isGood not guarded. Location: tracker.go:30-44.
4. MISSING_TEST — `slo.Evaluate` exact budgetRemaining==0 boundary (float-sensitive CanDeploy flip). Location: evaluator.go:54-57.
5. MISSING_TEST / MISSING_EDGE_CASE — recovery: no demo phase or test re-records good traffic after incident to prove budget rebounds and CanDeploy restores true. Design §Failure&Success criteria mentions "recovery" only loosely via compressed windows.
6. DOC_CODE_MISMATCH — design §Components says SLOEvaluator computes burn rate; burn rate lives in alerting engine, not slo. Severity LOW.
7. DOC_CODE_MISMATCH — design §Architecture §components + 03-exec imply recovery demonstrated; demo has no Phase 5. Severity LOW.
8. IMPLEMENTATION_SPECIFIC — `BurnRateRule` fields LongWindow/ShortWindow/BudgetConsumedPct declared in alerting engine struct but never wired; all rules share single shortTracker/longTracker pair. Undermines "multi-window per rule" semantics documented in research (not audited) and design §Components ("evaluating short and long windows"). Severity MEDIUM per design intent. No functional bug for current single-pair demo.
9. UNHANDLED_ERROR — `isGood` nil in WindowTracker.Record → nil-pointer panic. Only reachable via bad caller (demo/tests always supply closure). Severity LOW; not exploitable in current usage.

## Classification
- HIGH: 0
- MEDIUM: 1 (gap 8)
- LOW: 8
- UNVERIFIED_RESULT/FAKE_* : 0 — no fabricated artifacts; demo math independently reproduced.
