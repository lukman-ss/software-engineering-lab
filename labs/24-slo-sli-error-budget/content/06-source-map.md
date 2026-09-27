# Source Map

Pemetaan setiap seksi master draft ke sumber penelitian, file implementasi, dan tes.

---

## Architecture Overview

Research:
- research/01-plan.md (Concept to Prove)
- research/05-report.md (Findings 1-14, Areas of Agreement)

Implementation:
- internal/metrics/tracker.go (WindowTracker)
- internal/slo/evaluator.go (Evaluator)
- internal/alerting/engine.go (AlertEngine)

Demo:
- cmd/demo/main.go (PHASE 1-4 setup)

Tests:
- tests/slo_test.go (TestConcurrencyMetrics)

---

## SLI (Service Level Indicator)

Research:
- research/03-evidence.md (Evidence 1: Definition, Evidence 3: Good/Total ratio, Evidence 2: Common types)
- research/05-report.md (Finding 1: Definition, Finding 2: Ratio form, Finding 10: Percentile)

Implementation:
- internal/slo/evaluator.go:41-47 (SLI = good/total computation)
- cmd/demo/main.go:20-22 (isGood predicate)

Tests:
- tests/slo_test.go:13-58 (TestMetricsWindowTracker: good/bad aggregation)
- tests/slo_test.go:179-198 (TestEvaluatorZeroTraffic: SLI=1.0 on zero traffic)

---

## SLO (Service Level Objective)

Research:
- research/03-evidence.md (Evidence 4: SLA vs SLO, Evidence 3: SLO definition)
- research/05-report.md (Finding 3: SLO definition, Finding 4: SLA vs SLO)

Implementation:
- internal/slo/evaluator.go:10-14 (Config: Name, TargetUptime, LatencyThreshold)

Demo:
- cmd/demo/main.go:32-36 (Payment Service 99.9% SLO)
- cmd/demo/main.go:117-123 (Reports Service 95.0% SLO)

Tests:
- tests/slo_test.go:61-91 (TestSLOEvaluator: target 99%, budget consumption)

---

## Error Budget

Research:
- research/03-evidence.md (Evidence 7: Definition, Evidence 8: Decision Tool)
- research/05-report.md (Finding 7: Error Budget = 1 − SLO, Finding 8: Decision Tool)
- research/05-report.md (Finding 13: Datadog formula — MEDIUM confidence, vendor-specific)

Implementation:
- internal/slo/evaluator.go:49-57 (totalErrorBudget, budgetRemaining, CanDeploy)

Demo:
- cmd/demo/main.go:71-72 (Phase 1: budget +1.00, CanDeploy=true)
- cmd/demo/main.go:103-104 (Phase 2: budget -8.90, CanDeploy=false)

Tests:
- tests/slo_test.go:77-90 (TestSLOEvaluator: budget exhausted → CanDeploy=false)

---

## Burn Rate Alerting

Research:
- research/03-evidence.md (Evidence 9: Burn rate alerting, multi-window)
- research/05-report.md (Finding 9: Multi-window thresholds)
- research/04-contradictions.md (Contradiction 4: Google vs Datadog thresholds)

Implementation:
- internal/alerting/engine.go:51-61 (CalculateBurnRate)
- internal/alerting/engine.go:63-89 (Check: multi-window AND logic)

Demo:
- cmd/demo/main.go:38-49 (alertRules: 14.4x PAGE, 6.0x TICKET)
- cmd/demo/main.go:107-115 (Phase 3: alert triggered)

Tests:
- tests/slo_test.go:93-152 (TestAlertEngineBurnRate: true positive + transient spike filter)

---

## Multi-Window Concept

Research:
- research/03-evidence.md (Evidence 9: short=1/12 long window)
- research/04-contradictions.md (Window 28 vs 30 days)

Implementation:
- internal/alerting/engine.go:63-75 (shortTracker + longTracker AND condition)

Demo:
- cmd/demo/main.go:28-30 (shortWin=5min, longWin=60min)

Tests:
- tests/slo_test.go:130-151 (Negative test: short spike, long clean → no alert)

---

## Criticality Bucketing

Research:
- research/03-evidence.md (Evidence 12: Endpoint-specific SLO)
- research/05-report.md (Finding 12: Different SLO per endpoint)

Demo:
- cmd/demo/main.go:117-146 (Phase 4: Payment 99.9% vs Reports 95%)

Tests:
- (Conceptual — demo, not unit tested directly)

---

## WindowTracker Eviction & Out-of-Order

Research:
- research/03-evidence.md (Evidence 3: rolling window concept)
- research/04-contradictions.md (Window 28 vs 30 days)

Implementation:
- internal/metrics/tracker.go:106-127 (evictStaleLocked, Summary)
- internal/metrics/tracker.go:46-104 (Record with out-of-order handling)

Tests:
- tests/slo_test.go:154-177 (TestOutOfOrderTimestamps: bucket ordering + eviction)

---

## Concurrency & Thread-Safety

Research:
- (Engineering design concept)

Implementation:
- internal/metrics/tracker.go:22-28, 46-48 (sync.RWMutex)

Tests:
- tests/slo_test.go:200-236 (TestConcurrencyMetrics: 20 goroutine × 100 requests)

---

## Release Freeze Policy

Research:
- research/03-evidence.md (Evidence 8, 15: deploy decisions)
- research/05-report.md (Finding 8: Budget → Deploy decision)

Implementation:
- internal/slo/evaluator.go:54-57 (CanDeploy logic)

Demo:
- cmd/demo/main.go:72 (Phase 1: deployment allowed)
- cmd/demo/main.go:104 (Phase 2: deployment blocked)

---

## Failure Scenario / Incident Response

Research:
- research/03-evidence.md (Evidence 15: SLO Miss Policy — halt, postmortem)
- research/05-report.md (Finding 14: Policy template)

Demo:
- cmd/demo/main.go:74-98 (Phase 2: simulated incident, 10% error rate)
- cmd/demo/main.go:68 (Recovery context: budget must recover)

---

## Golden Signals (Context)

Research:
- research/03-evidence.md (Evidence 13: Four golden signals)
- research/05-report.md (Finding 15: Golden signals)

Implementation:
- Not implemented (lab uses isGood predicate; documented in engineering notes as simplification)

---

## Demo Output (Verified)

Research:
- N/A (demonstration, not research claim)

Implementation:
- cmd/demo/main.go (entire file)

Execution Result:
- engineering/03-execution-result.md:48-73 (captured demo output)

---

## Key Claims Audit Trail

| Claim | Research Evidence | Implement Proof | Test Case |
|---|---|---|---|
| SLI = good/total | Evidence 3, Finding 2 | evaluator.go:44-47 | TestMetricsWindowTracker |
| Error Budget = 1 − SLO | Evidence 7, Finding 7 | evaluator.go:49-52 | TestSLOEvaluator |
| Burn Rate = actual/allowed | Evidence 9, Finding 9 | engine.go:51-61 | TestAlertEngineBurnRate |
| Multi-window prevents false positive | Evidence 9, Test design | engine.go:73-75 | TestAlertEngineBurnRate (negative) |
| CanDeploy = false when budget ≤ 0 | Evidence 8, Finding 8 | evaluator.go:55-57 | TestSLOEvaluator |
| 100% not realistic | Evidence 6, Finding 6 | Config target ≤ 0.999 | Demo: SLO 99.9% |
| Percentile over average | Evidence 10, Finding 10 | — (simplified to boolean) | TestMetricsWindowTracker (latency as good/bad) |
| CPU/RAM not user SLO | Evidence 11, Finding 11 | isGood predicate = user-facing | — |
| Criticality bucketing | Evidence 12, Finding 12 | Demo Phase 4 | — |
| "70% outages from change" | Evidence 16 (LOW conf) | NOT claimed | — (flagged as Google-internal) |
| 100x cost per nine | Evidence 6 (heuristic) | NOT quantified | — (documented as heuristic) |
