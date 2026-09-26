# Research Gap Analysis: Lab 27 — Database Constraints

**Target Lab:** `labs/27-database-constraints`  
**Audit Date:** 2026-09-26  

---

## Gap 1: Multi-Database Comparison (MySQL / SQLite / SQL Server)
- **Type:** SCOPE_ERROR / OUTDATED_SOURCE
- **Severity:** LOW
- **Location:** `research/02-sources.md:79-86`, `research/04-contradictions.md:70-77`, `research/06-open-questions.md:82-91`
- **Problem:** MySQL documentation was blocked by HTTP 403. Consequently, findings and recommendations in `05-report.md` reflect PostgreSQL behavior almost exclusively (e.g. `NULLS NOT DISTINCT`, `23xxx` codes, SSI serializable implementation).
- **Required Revision:** None strictly required for a PostgreSQL-focused lab. If the lab is intended to be general SQL, clarify in `05-report.md` title/scope that claims reflect PostgreSQL engine behavior.
- **Can Be Approved Without Fix:** YES

---

## Gap 2: Multi-Column UNIQUE NULL Combinations
- **Type:** MISSING_CASE
- **Severity:** LOW
- **Location:** `research/06-open-questions.md:5-12`
- **Problem:** Exact behavior of `(1, NULL)` vs `(1, NULL)` under multi-column unique constraints is flagged as an open question needing a live test.
- **Required Revision:** Document standard SQL behavior during lab implementation: by SQL standard and PG implementation, multiple `(1, NULL)` rows are allowed under default distinct NULL semantics because `(1, NULL) = (1, NULL)` evaluates to NULL (unknown), unless `NULLS NOT DISTINCT` is defined.
- **Can Be Approved Without Fix:** YES (Appropriately tracked in open questions for lab verification).

---

## Gap 3: Industry Best Practice Case Studies
- **Type:** WEAK_SOURCE
- **Severity:** LOW
- **Location:** `research/02-sources.md:87-93`
- **Problem:** Cited industry blog articles resulted in 404s. The research relies purely on primary database documentation.
- **Required Revision:** Not blocking. Primary database documentation is authoritative and sufficient for core technical mechanisms.
- **Can Be Approved Without Fix:** YES
