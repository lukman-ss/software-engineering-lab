# Revision Plan

Target Lab: labs/29-saga-pattern

Previous Audit Status: APPROVED_WITH_WARNINGS (research-audit/07-verdict.md)

## Blocking Issues
None (per audit verdict lines 26-28)

## Non-Blocking Issues
1. MEDIUM — The 3-way transaction taxonomy (compensable, pivot, retryable) relies on a single public source (Microsoft Azure Architecture Center). (research-audit/06-gaps.md Gap 1)
2. MEDIUM — The specific 6-item list of isolation countermeasures relies on Microsoft Azure Architecture Center without external cross-enumeration. (research-audit/06-gaps.md Gap 2)
3. LOW — No automated protocol exists for failure of compensating transactions after retries; operational intervention is required. (research-audit/06-gaps.md Gap 3)
4. LOW — Garcia-Molina & Salem (1987) text was verified via citation chain / ACM references rather than direct OCR/text parsing of the scanned PDF. (research-audit/06-gaps.md Gap 4)

## Files To Modify
- `research/05-report.md` — Qualify pivot/retryable taxonomy and 6 countermeasures as Microsoft-specific; add compensation failure handling note.
- `research/03-evidence.md` — Add source-limitation notes to Evidence 7 (pivot/retryable) and Evidence 12 (6 countermeasures).
- `research/04-contradictions.md` — Add clarification that isolation countermeasure list reflects Microsoft's formalization.

## Verification Plan
- Source verification: Confirm Microsoft Azure Architecture Center pages remain reachable and contain the cited text.
- Content verification: Ensure added notes accurately reflect source limitations without altering supported claims.
- Documentation consistency: Check that research/05-report.md still aligns with evidence and does not introduce unsupported generalizations.
- No code changes required (pipeline override: research only).