# Docs vs Code Audit

## 1. README vs Code

README says structure:
- `internal/metrics`: "Sliding-window time-bucketed event tracker for recording requests and measuring good vs. total events." → matches `internal/metrics/tracker.go`. PASS
- `internal/slo`: "Evaluator calculating SLI ratios, remaining Error Budget, and release freeze policy enforcement." → matches `internal/slo/evaluator.go` (`Evaluate`, `CanDeploy`). PASS
- `internal/alerting`: "Multi-window burn-rate alert calculator evaluating fast and slow budget burn rates against SLO thresholds." → matches `internal/alerting/engine.go` (short+long trackers; 14.4x fast / 6.0x slow rules). PASS
- `cmd/demo`: "Executable demonstration illustrating baseline SLO tracking, error budget depletion during an incident, and burn rate alert triggering." → code implements all three AND a Phase 4 endpoint-criticality comparison (`Payment 99.9%` vs `Reports 95.0%`) that README does **not** mention. PASS (not incorrect, but under-described) → **MEDIUM-LOW mismatch**

README "Running Demo" section lists `go run ./cmd/demo` → matches entry point `cmd/demo/main.go`. PASS

## 2. Research Claims vs Code

Research-approved claims (from `research/05-report.md`) vs implementation:
- "Error budget is the complement of SLO (100% - SLO)" → code: `allowedFailureRate := 1.0 - e.config.TargetUptime`. PASS
- "Burn rate drives release velocity; high budget = ship faster; depleted = halt" → code: `CanDeploy` flips false on budget exhaustion. PASS
- Burn-rate factors 14.4x (page) / 6.0x (ticket) → research did NOT mandate these specific numbers (it flagged Datadog's 1-6/6+ icon thresholds as vendor-specific). The lab uses canonical Google SRE factors, not Datadog's. This is a deliberate alignment with the *primary* SRE source, not a contradiction. PASS
- Datadog's `error budget remaining = 100 * (current - target) / (100 - target)` formula → code uses absolute-count budget (`allowedFailureRate*total - bad`) instead. Research explicitly classified this formula as Datadog implementation-specific. The lab's alternative is mathematically equivalent in concept. Engineering notes (`02-implementation-notes.md:17`) document the choice. PASS (no contradiction; documented alternative)

## 3. Engineering Design (`engineering/01-design.md`) vs Code
- Claim: "Architecture" lists `internal/metrics`, `internal/slo`, `internal/alerting`, `cmd/demo` — all present. PASS
- Claim: "Components" maps to `Event`, `WindowTracker`, `SLOEvaluator`, `BurnRateAlertEngine`. Naming differs slightly (struct is `Evaluator`/`AlertEngine`, not `SLOEvaluator`/`BurnRateAlertEngine`) but behavior maps exactly. PASS (naming is descriptive, not a defect)

## 4. Execution Result (`engineering/03-execution-result.md`) vs Actual Code/Demo — **MISMATCH**
`03-execution-result.md` records the `go run ./cmd/demo` output ending at:
```
[PHASE 3] Checking Multi-Window Burn Rate Alerts...
>>> ALERT TRIGGERED: [TICKET] Slow Burn Alert ...
================================================================
  DEMO COMPLETE
================================================================
```
It **omits Phase 4**. The current `cmd/demo/main.go` (lines 117-147) produces an additional:
```
[PHASE 4] Endpoint Criticality Comparison (Payment 99.9% vs Reports 95.0%)...
Reports Target SLO: 95.0% | Current SLI: 90.0% | Budget Remaining: -5.00
Payment CanDeploy: false | Reports CanDeploy: false ...
```
This Phase 4 was added per `engineering-revision/02-changes-made.md` (Revision 3) but `03-execution-result.md` was never refreshed.

**Assessment: DOC_CODE_MISMATCH / stale recorded result.** Severity: **MEDIUM**.
- The recorded Phase 1-3 numbers were independently reproduced and are accurate (not fabricated). So this is NOT `FAKE_DEMO`.
- However, a reader relying on `03-execution-result.md` would be misled about what the current demo demonstrates (it claims the demo ends at Phase 3 and would not reveal the endpoint-criticality comparison).
- Required revision: refresh `03-execution-result.md` to include the full current demo output (Phases 1-4).

## 5. Success-Criteria vs Test Audit
`01-design.md` Success Criteria: "100% test coverage on core math and sliding window calculations."
- Actual measured coverage: 94.1%. Core math (`Evaluate`) = 100%; sliding-window helpers = 96.9–66.7%.
**Assessment: TEST_CLAIM_MISMATCH** (overclaim). Severity: **LOW** — core math is fully covered; the gap is defensive/validation branches only.

## Summary of Doc/Code Mismatches
| Type | Location | Severity |
|---|---|---|
| DOC_CODE_MISMATCH | `03-execution-result.md` omits Phase 4 (stale vs code) | MEDIUM |
| TEST_CLAIM_MISMATCH | `01-design.md` "100% coverage" claim vs 94.1% actual | LOW |
| DOC_CODE_MISMATCH | README `cmd/demo` description omits Phase 4 | LOW |
| RESEARCH_ALIGNMENT | None (Datadog-specific formula explicitly optional in research) | — |
