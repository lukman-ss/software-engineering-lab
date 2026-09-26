# Changes Made

## Revision 1

Audit Issue:
HIGH — RFC 2119 Misattribution (PSR-11 `MUST NOT` vs `SHOULD NOT`)

Files Changed:
- research/04-contradictions.md (Divergence 5)

Action:
- Changed `MUST NOT` to `SHOULD NOT` in PSR-11 quote
- Added interpretation note: strong recommendation, not strict prohibition

Verification:
- Verified PSR-11 spec text at https://www.php-fig.org/psr/psr-11/ section 1.3
- Confirmed text: "Users SHOULD NOT pass a container into an object..."
- No other `MUST NOT` usage in research files

Status:
RESOLVED

---

## Revision 2

Audit Issue:
MEDIUM — Heuristic Overreach (12-parameter threshold and value-object list)

Files Changed:
- research/05-report.md (Findings 11 and 12)

Action:
- Finding 11: Clarified 12-parameter rule as "lab-specific heuristic, not an industry standard"
- Finding 12: Clarified DateTime/Money/Address examples as "lab heuristics, not a universal standard"
- Added "(Lab Requirement/Heuristic)" tags to source references

Verification:
- Internal lab prompt heuristic preserved
- External authoritative sources remain unchanged
- Confidence rating unchanged (MEDIUM, number 12 not in primary source)

Status:
RESOLVED

---

## Revision 3

Audit Issue:
LOW — Source Tiering (Wikipedia classified as Tier 1)

Files Changed:
- research/02-sources.md (Sources 2 and 3)

Action:
- Downgraded Wikipedia tiering from "Tier 1 (encyclopedia)" to "Tier 2 (secondary encyclopedia)"
- Source 2 (Inversion of Control): updated line 18
- Source 3 (Dependency Injection): updated line 27

Verification:
- Wikipedia is community-edited secondary/tertiary source, not Tier 1 authoritative specification
- Cross-referenced facts verified against primary sources (Fowler, Spring, Laravel, PSR-11)
- No factual changes required; only tier classification corrected

Status:
RESOLVED
