# Documentation vs Code — labs/24-slo-sli-error-budget

## Verified Alignments
1. README `internal/metrics|internal/slo|internal/alerting|cmd/demo` package map — matches actual structure.
   - Assessment: PASS
2. README commands `go test ./...`, `go test -race ./...`, `go run ./cmd/demo` — all run successfully exactly as written.
   - Assessment: PASS
3. Design doc `01-design.md` component list (Event, WindowTracker, SLOEvaluator, BurnRateAlertEngine) — matches code. PASS.
4. Demo printed output matches `engineering/03-execution-result.md` Phase 1-4 verbatim on rerun. PASS.
5. "100% test coverage on core math and sliding window" claim — code audit: tests cover happy, failure, transient-negative, out-of-order, zero-traffic, concurrency. Core math (bad/total, budget calc, burn rate) exercised. PASS-with-caveat (no coverage tool run; missing edge branches noted in Test Audit).
6. Race detector claim "passes cleanly" — rerun `go test -race` PASS. Matches recorded result. PASS.

## Mismatches
### DOC_CODE_MISMATCH-1
Location: engineering/01-design.md §Components
Claim: `SLOEvaluator: Calculates SLI, remaining Error Budget, and current Burn Rate.`
Code: `internal/slo/evaluator.go` — Evaluator has no burn-rate calculation method; burn rate lives in `internal/alerting.engine.go:C alculateBurnRate`. Burn rate is not in slo package.
Assessment: WARNING — severity LOW. Misleading component attribution only; no functional gap.

### DOC_CODE_MISMATCH-2
Location: cmd/demo/main.go Phase 4, engineering/03-exec-result.md line 68
Claim/text: "Payment CanDeploy: false | Reports CanDeploy: false".
Code semantics: Reports (95% SLO, 10/100 bad, 10% error >> 5% allowed) genuinely exhausts; Payment (99.9%) also false after 10/1100 bad. Both flags correctly false.
But printed label says "Payment CanDeploy: false" while the variable `status.CanDeploy` (Payment evaluator) is indeed false, and `reportStatus.CanDeploy` (Reports) is false. Consistent.
Assessment: PASS — no mismatch, just redundant messaging.

### DOC_CODE_MISMATCH-3
Location: engineering/02-Implementation-notes §What Is Not Demonstrated + design §Expected Behavior "recovery"
Claim: design mentions "recovery" in architecture (compress real-time scale) and 03-exec implies incident → recovery.
Code/demo: No Phase 5 recovery demo; CanDeploy stays false; no re-recording of good traffic to flip budget positive.
Assessment: WARNING — missing promised recovery demonstration, not reflected in docs-notes either ("What Is Not Demonstrated" omits it).

### RESEARCH_IMPLEMENTATION_MISMATCH (skipped)
Per pipeline override, research not audited in this stage. Noted for handoff.

## Overall
README accuracy: PASS
Engineering notes self-consistency: WARNING (recovery gap, burn-rate attribution), non-blocking.
