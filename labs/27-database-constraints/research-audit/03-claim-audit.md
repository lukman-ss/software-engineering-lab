# Claim Audit: Lab 27 — Database Constraints

**Target Lab:** `labs/27-database-constraints`  
**Audit Date:** 2026-09-26  

---

## Claim 1: CHECK constraint condition passing with NULL results
- **Claim:** A CHECK constraint evaluates boolean expressions where `NULL` is treated as passing (the condition must evaluate to `FALSE` to fail).
- **Location:** `research/05-report.md:14`, `research/03-evidence.md:24-28`
- **Evidence Provided:** Cited from PG Docs Chapter 5.5.1 ("CHECK constraints pass if the condition evaluates to true or null").
- **Source:** PostgreSQL Chapter 5.5.1
- **Source Actually Supports Claim:** YES
- **Classification:** FACT
- **Severity:** LOW
- **Notes:** Standard SQL 3-valued logic behavior accurately documented.

---

## Claim 2: CHECK constraints cannot reference other rows or subqueries
- **Claim:** CHECK constraints in PostgreSQL cannot contain subqueries or reference columns from rows other than the current row being inserted/updated.
- **Location:** `research/05-report.md:22`, `research/03-evidence.md:25`, `research/05-report.md:135`
- **Evidence Provided:** PostgreSQL Doc 5.5.1 Note: "CHECK expressions cannot contain subqueries nor refer to variables other than columns of the current row".
- **Source:** PostgreSQL Chapter 5.5.1
- **Source Actually Supports Claim:** YES
- **Classification:** FACT
- **Severity:** LOW
- **Notes:** Crucial distinction between row constraints and triggers/assertions.

---

## Claim 3: UNIQUE constraints create B-tree indexes and treat NULLs as distinct by default
- **Claim:** Adding a UNIQUE or PRIMARY KEY constraint automatically creates a unique B-tree index. NULL values are considered distinct by default unless `NULLS NOT DISTINCT` is explicitly defined (PostgreSQL 15+).
- **Location:** `research/05-report.md:15`, `research/03-evidence.md:6-9`
- **Evidence Provided:** PostgreSQL Chapter 5.5.3 and 11.6 citations.
- **Source:** PostgreSQL Chapter 5.5.3 & 11.6
- **Source Actually Supports Claim:** YES
- **Classification:** FACT
- **Severity:** LOW
- **Notes:** Accurate across PostgreSQL version lineage.

---

## Claim 4: Foreign Key index creation asymmetry
- **Claim:** A Foreign Key constraint automatically requires/uses an index on the referenced (parent) table, but does NOT automatically create an index on the referencing (child) column.
- **Location:** `research/05-report.md:17,28`
- **Evidence Provided:** PostgreSQL 5.5 DDL Constraints FK section.
- **Source:** PostgreSQL Chapter 5.5.5
- **Source Actually Supports Claim:** YES
- **Classification:** FACT
- **Severity:** MEDIUM
- **Notes:** Real-world performance pitfall (child index needed for fast CASCADE and foreign key check updates) correctly highlighted.

---

## Claim 5: Elimination of read-then-write race conditions by UNIQUE constraints
- **Claim:** Concurrent inserts with identical unique keys are resolved atomically at the index/row lock level, converting race conditions into `23505 unique_violation` errors without requiring manual application-level mutexes or `SELECT ... FOR UPDATE`.
- **Location:** `research/05-report.md:36-56`
- **Evidence Provided:** PostgreSQL Chapter 13.3 & 13.4 citations.
- **Source:** PostgreSQL Chapter 13.3 & 13.4
- **Source Actually Supports Claim:** YES
- **Classification:** FACT
- **Severity:** LOW
- **Notes:** Accurate explanation of relational database write serialization semantics.

---

## Claim 6: Partial index query planner implication requirements
- **Claim:** Partial indexes (`CREATE UNIQUE INDEX ... WHERE <pred>`) require the query's `WHERE` clause to mathematically imply the index predicate at plan time. Prepared/parameterized statements with bind parameters (e.g. `x < $1`) cannot imply constant predicates.
- **Location:** `research/05-report.md:88-92`, `research/03-evidence.md:63-65`
- **Evidence Provided:** PostgreSQL Chapter 11.8 citation.
- **Source:** PostgreSQL Chapter 11.8
- **Source Actually Supports Claim:** YES
- **Classification:** FACT
- **Severity:** LOW
- **Notes:** Important constraint limitation correctly documented.

---

## Claim 7: Partitioned table unique constraints must include partition keys
- **Claim:** UNIQUE and EXCLUDE constraints on partitioned tables must include all partition key columns because index enforcement is local to each partition.
- **Location:** `research/05-report.md:150-153`, `research/03-evidence.md:103-105`
- **Evidence Provided:** PostgreSQL Chapter 5.12.2.3 citation.
- **Source:** PostgreSQL Chapter 5.12
- **Source Actually Supports Claim:** YES
- **Classification:** FACT / IMPLEMENTATION-SPECIFIC
- **Severity:** LOW
- **Notes:** Correctly identifies architectural boundary of partitioned indexes in PostgreSQL.
