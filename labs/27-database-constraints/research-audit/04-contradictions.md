# Contradiction Audit: Lab 27 — Database Constraints

**Target Lab:** `labs/27-database-constraints`  
**Audit Date:** 2026-09-26  

---

## Evaluation Summary

Research documents (`01-plan.md` through `06-open-questions.md`) were compared against one another for internal conflicts, citation discrepancies, and inconsistent claims.

### Internal Consistency
- **Evaluation Order Statement:** In `04-contradictions.md:53`, the note mentions CHECK constraints evaluated alphabetically, and in `06-open-questions.md:29`, the research notes that CHECK constraints run after NOT NULL checks. This matches PostgreSQL DDL documentation.
- **Partial Unique Index vs Constraint:** In `04-contradictions.md:30-35` and `05-report.md:93-97`, the distinction between a `UNIQUE` constraint (cataloged in `table_constraints`, auto-named) and a partial unique index (built via `CREATE UNIQUE INDEX ... WHERE`, not a formal ANSI SQL constraint catalog entry) is consistently maintained.
- **NULL Handling:** `03-evidence.md`, `04-contradictions.md`, and `05-report.md` consistently state that PostgreSQL treats NULLs as distinct by default in UNIQUE indexes, while noting `NULLS NOT DISTINCT` as a PostgreSQL 15+ feature.

### Source Conflicts
- No conflicting statements were made between cited PostgreSQL documentation sources.
- Blocked sources (MySQL 403, 404 links) were not invented or resolved with unsubstantiated guesses; they were properly cordoned off as open questions / unverified.

---

## Verdict

No material contradictions found.
