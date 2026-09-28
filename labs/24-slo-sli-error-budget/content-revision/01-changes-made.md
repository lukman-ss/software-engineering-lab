# Content Revision Record

Target Lab: `labs/24-slo-sli-error-budget`
Revision Date: 2026-09-28
Reviser: Technical Content Reviser

## Changes Applied

### 1. Fixed Burn Rate Calculation in Code Snippet 11 (`03-code-snippets.md:442`)
- **Before**: Claimed burn rate = 100x based on isolated batch (10% error / 0.1% allowed)
- **After**: Corrected to cumulative calculation: 10/1100 = 0.91% error rate → 9.09x burn rate
- **Rationale**: AlertEngine computes burn rate on cumulative data (baseline + incident), confirmed by demo output `ShortBurn: 9.09x`
- **Impact**: Aligns with audit finding HIGH-1

### 2. Fixed Formula Example in Master Draft (`02-master-draft.md:108-114`)
- **Before**: Comment `// mis. 10% / 1000` and `// 10% / 0.1% = 100x` implied 100x burn rate
- **After**: Example uses clear 20 error / 1000 total = 2% → 20x burn rate. Added note about demo Fase 2 actual 9.09x
- **Rationale**: Clarifies formula with consistent units; separates example from demo scenario
- **Impact**: Resolves audit finding MEDIUM-3

### 3. Fixed FALSE POSITIVE Diagram in D4 (`04-diagrams.md:154-165`)
- **Before**: Both windows showed 90 OK + 10 ERR (100x burn) with label `Short ≥ 6.0 BUT Long < 6.0 → NO ALERT` — internal contradiction
- **After**: FALSE POSITIVE column shows short 100x (transient) vs long 0.1x (clean) → NO ALERT. TRUE NEGATIVE column shows both windows below threshold.
- **Rationale**: Correctly illustrates transient spike (short high, long low) filtered by multi-window logic, matching `tests/slo_test.go:130-151`
- **Impact**: Fixes audit finding HIGH-2

### 4. Fixed Typos in Key Takeaways (`05-key-takeaways.md`)
- `perbaiken` → `perbaikan`
- `semaakin` → `semakin`
- **Impact**: Improves readability, non-blocking issue from verdict

## Verification

- Burn rate examples now consistently use cumulative calculation (0.91% / 0.1% = 9.09x)
- D4 FALSE POSITIVE column correctly shows transient-spike scenario
- All code snippets remain verbatim from approved implementation
- All formulas and diagrams now consistent with `internal/alerting/engine.go:63-75` and `tests/slo_test.go`

## Remaining Items (Non-blocking per verdict)

- Phase 4 demo could show differentiation at lower error rates (e.g., 2%) for clearer endpoint criticality illustration
- Release policy simplified vs Google Appendix B (postmortem >20%, P0 exception) — documented but not corrected as it is a simplification, not inaccuracy

## Final Status

REVISIONS COMPLETE — Content now consistent with approved research and engineering implementation.
Ready for re-audit.
