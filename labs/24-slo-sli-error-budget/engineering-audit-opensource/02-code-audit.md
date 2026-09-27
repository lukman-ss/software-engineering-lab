# Code Audit — labs/24-slo-sli-error-budget

Files read: internal/metrics/tracker.go, internal/slo/evaluator.go, internal/alerting/engine.go, cmd/demo/main.go, tests/slo_test.go, go.mod. No code modified.

## Finding 1

Location: internal/metrics/tracker.go:46-115 (Record, evictStaleLocked)
Claimed Behavior: Sliding-window bucketed aggregation; stale buckets evicted; out-of-order events handled.
Observed Implementation: Same-bucket fast path; ordered insert for late events; eviction by `now.Add(-windowSize)` with `Before(cutoff)` (boundary bucket kept). `Summary` evicts then sums under write lock.
Assessment: PASS
Severity: LOW
Notes: Eviction keyed to caller-supplied event timestamp, so one far-future event wipes current history. No test covers it. Acceptable for in-memory lab; document.

## Finding 2

Location: internal/metrics/tracker.go:117-128 (Summary takes mu.Lock, not RLock)
Claimed Behavior: Thread-safe concurrent recording.
Observed Implementation: Full mutex on both Record and Summary; race detector clean on 20x100 concurrent test.
Assessment: PASS
Severity: LOW
Notes: Correct, only conservative (readers serialized). No race.

## Finding 3

Location: internal/slo/evaluator.go:41-70
Claimed Behavior: SLI = good/total; budget = (1-target)*total; freeze (CanDeploy=false) when budget exhausted.
Observed Implementation: Matches exactly. Zero-traffic yields SLI=1.0, CanDeploy=true. Boundary `budgetRemaining <= 0` freezes at exactly zero. Rounding (SLI 4dp, budget 2dp) display-only.
Assessment: PASS
Severity: LOW
Notes: `TargetUptime` unvalidated (0, >1, negative silently produce nonsense budgets); nil tracker panics. No test pins invalid config. Defensive-validation gap only.

## Finding 4

Location: internal/alerting/engine.go:51-89
Claimed Behavior: Multi-window multi-burn-rate alerting; transient spikes suppressed (short AND long must breach).
Observed Implementation: Burn math `actual/allowed` correct with total==0 and allowed<=0 guards. AND-gating proven by negative test (short 100x + long 0.1x → no alert). Demo burn 9.09x fires only the 6x rule, not 14.4x — differential thresholds work.
Assessment: WARNING
Severity: MEDIUM
Notes: `BurnRateRule.LongWindow/ShortWindow/BudgetConsumedPct` fields are never read; every rule shares the two construction-time trackers. Per-rule window semantics claimed by field names do not exist. No minimum-sample guard: 1 bad / 1 total at 99.9% SLO = 1000x burn and pages. Works as tested, less general than struct implies.

## Finding 5

Location: cmd/demo/main.go (4 phases)
Claimed Behavior: Baseline → incident depletion + freeze → burn alert → criticality comparison.
Observed Implementation: Ran `go run ./cmd/demo`; output matches engineering/03-execution-result.md numerically (1000/0 baseline, 1100/1090/10 incident, -8.90 budget, 9.09x TICKET-only alert). Real execution, no fabrication.
Assessment: WARNING
Severity: MEDIUM
Notes: Phase 4 feeds identical 10% error rate to both 99.9% and 95% SLOs, so both freeze (Reports -5.00). Demonstrates budget math per SLO, not usable differentiation; criticality claim rests on arithmetic, unproven by demo scenario.

## Finding 6

Location: cross-cutting (timeouts, recovery, rollback, cleanup)
Claimed Behavior: None claimed (in-memory lab, no goroutines/timers/IO).
Observed Implementation: No leaked resources; nothing to clean up.
Assessment: PASS
Severity: LOW
Notes: Design doc mentions "recovery" demo phase that does not exist — tracked as docs mismatch, not code defect.
