# Revision Plan

Target Lab: labs/15-load-testing

Previous Audit Status: APPROVED_WITH_WARNINGS

## Blocking Issues

None.

## Non-Blocking Issues

1. **SCOPE_ERROR**: Residual Lab 16 Spring IoC / Service Locator artifacts remain in `research/02-sources.md` (Source 16) and `research/04-contradictions.md` (Contradiction 6).
2. **MISSING_SOURCE**: No dedicated primary documentation URLs for Apache JMeter and Gatling in `research/02-sources.md` (only generic domain links in report).
3. **WEAK_SOURCE**: ISO/IEC 25010 cited via Wikipedia (Source 10) — should clarify Tier 2 status and add ISO standard reference.
4. **UNVERIFIED_CLAIM**: Numeric SLA targets (P95 < 500ms, Error Rate < 1%, CPU < 75%, Memory < 80%) in `research/01-plan.md` objective 6 stated without caveat that they are illustrative/context-dependent.

## Files To Modify

- `research/02-sources.md` — remove Source 16, add JMeter/Gatling entries, update Source 10 tier
- `research/04-contradictions.md` — remove Contradiction 6
- `research/01-plan.md` — add illustrative caveat to SLA targets in objective 6

## Verification Plan

- source verification: Confirm JMeter/Gatling URLs reachable and primary
- documentation consistency: Verify all claims properly qualified
- audit traceability: Each revision addresses specific audit finding