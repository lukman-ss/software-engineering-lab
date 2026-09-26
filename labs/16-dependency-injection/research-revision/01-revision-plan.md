# Revision Plan

Target Lab: labs/16-dependency-injection
Previous Audit Status: APPROVED_WITH_WARNINGS

## Blocking Issues

None.

## Non-Blocking Issues

1. **RFC 2119 Misattribution (HIGH)**: `research/04-contradictions.md` Divergence 5 incorrectly claims PSR-11 uses `MUST NOT` for passing containers; PSR-11 uses `SHOULD NOT`. The evidence in `03-evidence.md` is correct; the error is in the contradiction analysis layer.

2. **Heuristic Overreach (MEDIUM)**: The "12 parameters" numeric threshold and the specific value-object list (DateTime, Money, Address) are internal lab conventions, not universal external standards. Presentation should distinguish lab heuristics from industry norms.

3. **Source Tiering (LOW)**: Wikipedia sources (Source 2, Source 3 in `02-sources.md`) incorrectly classified as Tier 1 authoritative. Should be Tier 2 (secondary/encyclopedia). Cross-referenced facts are safe; tiering classification needs correction.

## Files To Modify

- `research/04-contradictions.md` — Fix Divergence 5 PSR-11 keyword
- `research/05-report.md` — Clarify Findings 11/12 as lab heuristics
- `research/02-sources.md` — Downgrade Wikipedia tiering

## Verification Plan

- Verify PSR-11 text at https://www.php-fig.org/psr/psr-11/ — confirm `SHOULD NOT`
- Confirm no other files use `MUST NOT` in PSR-11 context
- Verify internal consistency of edited files
