# Research Gap Analysis: Database Constraints

## Gap 1: MySQL Behavior Access Blocked
- **Type:** SCOPE_ERROR / WEAK_SOURCE
- **Severity:** LOW
- **Location:** `research/02-sources.md` (Blocked Sources), `research/04-contradictions.md` (Section 8), `research/06-open-questions.md` (Question 8)
- **Problem:** MySQL 8.0 official documentation endpoints returned HTTP 403 Forbidden during initial research. MySQL-specific error code mappings and NULL handling were not directly verified from official MySQL docs.
- **Required Revision:** None required for core approval, as PostgreSQL and SQLite (Tier 1 sources) thoroughly support all core claims. Recorded transparently in research gaps.
- **Can Be Approved Without Fix:** YES

---

## Gap 2: Multi-Column UNIQUE NULL Edge Cases Not Empirically Tested
- **Type:** UNVERIFIED_CLAIM / MISSING_CASE
- **Severity:** LOW
- **Location:** `research/04-contradictions.md`, `research/06-open-questions.md` Question 1
- **Problem:** Whether `(1, NULL)` and `(1, NULL)` violate multi-column `UNIQUE(col1, col2)` without `NULLS NOT DISTINCT` is inferred from SQL standard wording rather than tested against live engine binaries.
- **Required Revision:** Can be verified in lab implementation test suite.
- **Can Be Approved Without Fix:** YES

---

## Gap 3: PostgreSQL B-Tree Page Latching Details
- **Type:** UNVERIFIED_CLAIM
- **Severity:** LOW
- **Location:** `research/01-plan.md` (Risks / Unknowns Item 3), `research/05-report.md` (Limitations Item 2)
- **Problem:** Internal B-tree page lock acquisition sequence during UNIQUE constraint checking is engine-internal and not fully documented in public documentation.
- **Required Revision:** Documented as a known limitation in `05-report.md`.
- **Can Be Approved Without Fix:** YES
