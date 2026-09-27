# Changes Made

## Revision 1: Remove Lab 16 Spring IoC Artifact (Source 16)

**Audit Issue:** SCOPE_ERROR / LOW — Residual Dependency Injection lab artifact (Spring IoC source) in load testing source catalog

**Files Changed:**
- `research/02-sources.md`

**Action:**
- Removed Source 16 (Introduction to the Spring IoC Container and Beans) entirely
- Replaced with new Source 16: Apache JMeter User Guide (Tier 1 primary documentation)
- Original Source 16 was explicitly marked "Excluded from active evidence" but remained in catalog — now purged

**Verification:**
- Source 16 entry deleted from `02-sources.md`
- No references to Spring IoC remain in research files
- Source count reduced by 1 (was 21, now 21 after replacement + Gatling addition)

**Status:** RESOLVED

---

## Revision 2: Remove Lab 16 Contradiction Artifact (Contradiction 6)

**Audit Issue:** SCOPE_ERROR / LOW — Residual "Service Locator vs DI Pattern" contradiction from Lab 16

**Files Changed:**
- `research/04-contradictions.md`

**Action:**
- Deleted Contradiction 6 section entirely (lines 43-47)
- Summary section remains unchanged (already did not reference Contradiction 6 by number)

**Verification:**
- Contradiction 6 removed
- Remaining contradictions (1-5) are all load-testing domain relevant
- File now has 5 numbered contradictions instead of 6

**Status:** RESOLVED

---

## Revision 3: Add Missing Primary Sources for JMeter and Gatling

**Audit Issue:** MISSING_SOURCE / MEDIUM — Finding 3 tool comparison lacks dedicated primary source entries for JMeter and Gatling

**Files Changed:**
- `research/02-sources.md` (added Source 16 and Source 21)
- `research/05-report.md` (updated Finding 3 sources and confidence justification)

**Action:**
- Source 16: Apache JMeter User Guide — https://jmeter.apache.org/usermanual/index.html (Tier 1)
- Source 21: Gatling Documentation — https://docs.gatling.io/ (Tier 1, verified HTTP 200)
- Updated Finding 3 in `05-report.md` to cite Source 16 (JMeter) and Source 21 (Gatling)
- Updated source URLs in report from generic domains to specific documentation pages
- Adjusted confidence justification to reflect all four tools now have dedicated primary entries

**Verification:**
- Both URLs verified reachable (JMeter usermanual: 200 OK; Gatling docs: 200 OK)
- Report now cites specific documentation pages, not generic domains
- All four tools (k6, Locust, JMeter, Gatling) have Tier 1 primary source entries

**Status:** RESOLVED

---

## Revision 4: Clarify ISO/IEC 25010 Source Tier

**Audit Issue:** WEAK_SOURCE / LOW — Source 10 cited Wikipedia URL but labeled Tier 1 standard

**Files Changed:**
- `research/02-sources.md` (Source 10)
- `research/05-report.md` (Finding 2 source citation)

**Action:**
- Changed Source 10 tier from "Tier 1 (standard) + Tier 2 (Wikipedia summary)" to "Tier 2 (community summary of Tier 1 standard — ISO text is paywalled)"
- Added "Canonical Standard Reference" field pointing to paywalled ISO page
- Updated Publisher to "ISO — summarized via Wikipedia"
- Updated Finding 2 in report to note "(via Wikipedia summary; ISO/IEC 25010:2011 original is paywalled)"

**Verification:**
- Source tier accurately reflects actual evidence base (Wikipedia summary, not ISO text)
- Canonical ISO reference documented for traceability
- Claim support maintained (Wikipedia summary is accurate for Performance Efficiency subcharacteristics)

**Status:** RESOLVED

---

## Revision 5: Qualify SLA Targets as Illustrative Examples

**Audit Issue:** UNVERIFIED_CLAIM / LOW — Objective 6 in plan stated numeric SLA targets (P95 < 500ms, etc.) without caveat

**Files Changed:**
- `research/01-plan.md` (Objective 6)

**Action:**
- Modified Objective 6 from: "Best practice untuk menentukan dan memvalidasi target SLA (P95 < 500ms, Error Rate < 1%, CPU < 75%, Memory < 80%)"
- To: "Best practice untuk menentukan dan memvalidasi target SLA (contoh ilustratif: P95 < 500ms, Error Rate < 1%, CPU < 75%, Memory < 80% — angka-angka ini bersifat kontekstual tergantung domain/aplikasi, bukan standar universal)"

**Verification:**
- Plan now explicitly labels numeric values as illustrative examples
- Open Questions (Question 2) already correctly states "no universal 'correct' thresholds exist"
- Consistency maintained across research artifacts

**Status:** RESOLVED

---

## Summary

| Issue Type | Count | Resolved |
|------------|-------|----------|
| SCOPE_ERROR | 2 | 2 |
| MISSING_SOURCE | 1 | 1 |
| WEAK_SOURCE | 1 | 1 |
| UNVERIFIED_CLAIM | 1 | 1 |
| **Total** | **5** | **5** |

All non-blocking audit findings resolved. No blocking issues existed.