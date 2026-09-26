# Audit Plan: Lab 27 — Database Constraints

**Target Lab:** `labs/27-database-constraints`  
**Audit Date:** 2026-09-26  
**Auditor Mode:** PIPELINE OVERRIDE — Research Audit Only

---

## 1. Scope & Strategy

As instructed by PIPELINE OVERRIDE, this audit evaluates **research files only** (`research/01-plan.md` through `06-open-questions.md`). No implementation/code files were inspected or executed. No research files were modified.

The audit checks:
1. Validity and reachability of cited sources.
2. Verification of all technical claims against official PostgreSQL documentation.
3. Logical consistency and internal contradictions across research documents.
4. Completeness of research gap analysis.

---

## 2. Target Files Reviewed

- `research/01-plan.md`
- `research/02-sources.md`
- `research/03-evidence.md`
- `research/04-contradictions.md`
- `research/05-report.md`
- `research/06-open-questions.md`

---

## 3. Major Technical Claims To Verify

1. **CHECK constraints & NULL evaluation:** CHECK condition passes when expression evaluates to NULL (unknown).
2. **CHECK constraints & immutability / cross-row restriction:** PostgreSQL CHECK constraints cannot contain subqueries or reference other rows; PostgreSQL assumes condition immutability and only checks INSERT/UPDATE.
3. **UNIQUE constraint & NULL handling:** `NULL` values are treated as distinct by default in PostgreSQL UNIQUE constraints unless `NULLS NOT DISTINCT` is specified (PostgreSQL 15+).
4. **FOREIGN KEY & MATCH SIMPLE:** Foreign keys with `MATCH SIMPLE` allow any FK column to be NULL, bypassing validation if any column is NULL.
5. **Partial Unique Indexes:** Partial indexes (`WHERE` clause) enforce uniqueness on filtered subsets; query planner requires mathematical implication to utilize partial indexes; parameterized queries do not imply constant predicates.
6. **SQLSTATE Error Mapping:** Integrity violations map to Class 23 (`23502` NOT NULL, `23503` FK, `23505` UNIQUE, `23514` CHECK, `23P01` EXCLUDE).
7. **Partitioning Limitations:** Unique and EXCLUDE constraints on partitioned tables must include all partition key columns.

---

## 4. Source Audit Strategy

- Verify PostgreSQL documentation URLs directly against official PG docs syntax and content.
- Assess handling of unreachable sources (MySQL 403, Citus Data 404, Martin Fowler 404).

---

## 5. Primary Risks Identified

- Over-reliance on PostgreSQL-specific semantics presented as universal relational database facts (mitigated partially in research notes, but needs auditing).
- Verification of claimed planner behavior regarding partial index implication.
