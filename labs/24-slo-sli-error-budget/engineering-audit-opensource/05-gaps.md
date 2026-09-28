# Gap Analysis

## GAP-1

Type: MISSING_EDGE_CASE
Severity: MEDIUM
Description: No test for 100% error rate (all events bad). No test for budgetRemaining exactly 0 boundary.
Impact: Core math boundary unproven.

## GAP-2

Type: IMPLEMENTATION_OVERCLAIM
Severity: MEDIUM
Description: engineering/01-design.md claims "100% test coverage" — not met.
Impact: Overclaim in design doc.

## GAP-3

Type: DOC_CODE_MISMATCH
Severity: LOW
Description: README says "Multi-Burn-Rate" but rules use single BurnRateFactor for both windows.
Impact: Minor docs/code mismatch.

## GAP-4

Type: WARNING (no exact gap type; mapped to IMPLEMENTATION_OVERCLAIM)
Severity: MEDIUM
Description: Alert engine uses AND logic (both windows must exceed threshold). Standard SRE uses OR. Can miss fast-burn incidents.
Impact: Alerting may not fire during short-window-only spikes.

## GAP-5

Type: MISSING_EDGE_CASE
Severity: LOW
Description: Evaluator Config.LatencyThreshold field is never used (dead config).
Impact: Minor — latency filtering is handled by isGood callback, but dead field is confusing.