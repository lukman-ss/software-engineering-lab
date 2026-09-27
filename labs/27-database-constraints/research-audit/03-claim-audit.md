# Claim Audit: Database Constraints

## Claim 1: Race Condition in Check-Then-Act Pattern
- **Claim:** Application-level validation using `SELECT → exists() → INSERT` is vulnerable to race conditions under concurrent requests, leading to duplicate inserts.
- **Location:** `research/05-report.md` Finding 1; `research/03-evidence.md` Item 1
- **Evidence Provided:** Concurrency diagram, PostgreSQL explicit locking / constraint checking timing.
- **Source:** PostgreSQL DDL Constraints (5.5), Explicit Locking (13.3).
- **Source Actually Supports Claim:** YES
- **Classification:** FACT
- **Severity:** HIGH (Core thesis)
- **Notes:** Check-then-act is fundamentally non-atomic without serializable isolation or locking.

---

## Claim 2: UNIQUE Constraint Prevents Race Conditions Atomically
- **Claim:** A `UNIQUE` constraint ensures only the first concurrent INSERT succeeds; subsequent inserts fail with SQLSTATE `23505` (`unique_violation`).
- **Location:** `research/05-report.md` Finding 2; `research/03-evidence.md` Item 2
- **Evidence Provided:** Unique index auto-creation, row-level exclusive locks during INSERT, error code `23505`.
- **Source:** PostgreSQL DDL Constraints (5.5.3), Error Codes (Appendix A), Explicit Locking (13.3).
- **Source Actually Supports Claim:** YES
- **Classification:** FACT
- **Severity:** HIGH
- **Notes:** Verified by PostgreSQL and SQLite specifications.

---

## Claim 3: NOT NULL Constraint Efficiency vs CHECK
- **Claim:** `NOT NULL` constraint is functionally equivalent to `CHECK (column IS NOT NULL)` but PostgreSQL implements it more efficiently.
- **Location:** `research/05-report.md` Finding 3; `research/03-evidence.md` Item 4
- **Evidence Provided:** Direct quotation from PostgreSQL DDL Constraints section 5.5.2.
- **Source:** PostgreSQL DDL Constraints 5.5.2.
- **Source Actually Supports Claim:** YES
- **Classification:** FACT
- **Severity:** LOW
- **Notes:** Direct quote from official documentation.

---

## Claim 4: FOREIGN KEY Constraint Does Not Auto-Index Referencing Columns
- **Claim:** `FOREIGN KEY` enforces referential integrity but does NOT automatically create an index on the referencing table columns in PostgreSQL.
- **Location:** `research/05-report.md` Finding 4; `research/03-evidence.md` Item 5
- **Evidence Provided:** Direct quotation from PostgreSQL DDL Constraints section 5.5.5 explaining why referencing columns are not automatically indexed.
- **Source:** PostgreSQL DDL Constraints 5.5.5.
- **Source Actually Supports Claim:** YES
- **Classification:** FACT / IMPLEMENTATION-SPECIFIC (PostgreSQL behavior)
- **Severity:** MEDIUM
- **Notes:** Accurate; correctly flags performance implications on cascading deletes.

---

## Claim 5: CHECK Constraints Are Row-Scoped and Assume Immutability
- **Claim:** `CHECK` constraints can only reference columns of the current row, cannot contain subqueries, and PostgreSQL assumes expressions are immutable.
- **Location:** `research/05-report.md` Finding 5; `research/03-evidence.md` Item 6
- **Evidence Provided:** Direct quotes from PostgreSQL DDL Constraints section 5.5.1 and SQLite CREATE TABLE docs.
- **Source:** PostgreSQL DDL Constraints 5.5.1, SQLite Docs.
- **Source Actually Supports Claim:** YES
- **Classification:** FACT
- **Severity:** HIGH
- **Notes:** Explains the critical boundary of where database constraints end and application/trigger logic begins.

---

## Claim 6: Partial Unique Indexes for Soft-Delete Patterns
- **Claim:** `CREATE UNIQUE INDEX ... WHERE deleted_at IS NULL` enforces uniqueness only for active records, allowing value reuse after soft deletion.
- **Location:** `research/05-report.md` Finding 6; `research/03-evidence.md` Item 7
- **Evidence Provided:** Partial index examples and predicate rules from PostgreSQL Indexes Partial section 11.8.
- **Source:** PostgreSQL Indexes Partial (11.8).
- **Source Actually Supports Claim:** YES
- **Classification:** FACT / PATTERN
- **Severity:** HIGH
- **Notes:** Standard industry and PostgreSQL solution for soft-delete uniqueness.

---

## Claim 7: Partial Index Planner Implication and Parameterized Query Limitations
- **Claim:** A partial index can be used by a query only if the query's WHERE clause mathematically implies the index predicate; parameterized query clauses do not match partial index predicates.
- **Location:** `research/05-report.md` Finding 6; `research/03-evidence.md` Item 8
- **Evidence Provided:** Direct quotation from PostgreSQL Partial Indexes section 11.8.
- **Source:** PostgreSQL Partial Indexes (11.8).
- **Source Actually Supports Claim:** YES
- **Classification:** FACT
- **Severity:** MEDIUM
- **Notes:** Important practical nuance for engineers using partial unique indexes.

---

## Claim 8: Constraint Violations Return Class 23 SQLSTATE and Structured Fields
- **Claim:** Constraint violations return SQLSTATE class 23 codes (e.g. `23505`, `23502`, `23503`, `23514`), and applications should inspect error codes and structured fields rather than localized error text.
- **Location:** `research/05-report.md` Finding 7; `research/03-evidence.md` Item 9
- **Evidence Provided:** PostgreSQL Error Codes Appendix A error listing and explanation of separate error report fields.
- **Source:** PostgreSQL Error Codes Appendix A.
- **Source Actually Supports Claim:** YES
- **Classification:** FACT / BEST PRACTICE
- **Severity:** HIGH
- **Notes:** Strongly supported by documentation.

---

## Claim 9: Multi-Row Invariants Require SERIALIZABLE Isolation or Triggers
- **Claim:** Invariants involving multiple rows cannot be expressed via CHECK constraints and require `SERIALIZABLE` isolation with retry handling (`40001` `serialization_failure`) or custom triggers.
- **Location:** `research/05-report.md` Finding 8; `research/03-evidence.md` Item 10
- **Evidence Provided:** PostgreSQL App-Level Consistency 13.4 and error code `40001`.
- **Source:** PostgreSQL App-Level Consistency 13.4.
- **Source Actually Supports Claim:** YES
- **Classification:** FACT / ARCHITECTURAL PATTERN
- **Severity:** HIGH
- **Notes:** Correctly identifies that Read Committed transactions suffer from shifting views.

---

## Claim 10: Production Constraint Migration Pattern (NOT VALID + VALIDATE)
- **Claim:** Adding constraints on large live tables should use `ADD CONSTRAINT ... NOT VALID` followed by `VALIDATE CONSTRAINT` to minimize table lock duration.
- **Location:** `research/05-report.md` Finding 10; `research/03-evidence.md` Item 11
- **Evidence Provided:** Lock mode analysis (`SHARE UPDATE EXCLUSIVE` vs `ACCESS EXCLUSIVE`) from ALTER TABLE documentation.
- **Source:** PostgreSQL ALTER TABLE.
- **Source Actually Supports Claim:** YES
- **Classification:** FACT / OPERATIONAL PATTERN
- **Severity:** HIGH
- **Notes:** Accurate PostgreSQL operational procedure.

---

## Claim 11: Partitioned Table Unique Constraint Limitations
- **Claim:** Unique or primary key constraints on partitioned tables must include all partition key columns because individual indexes enforce uniqueness only within their own partitions.
- **Location:** `research/05-report.md` Finding 11; `research/03-evidence.md` Item 14
- **Evidence Provided:** Direct quotation from PostgreSQL DDL Partitioning section 5.12.2.3.
- **Source:** PostgreSQL DDL Partitioning 5.12.2.3.
- **Source Actually Supports Claim:** YES
- **Classification:** FACT / LIMITATION
- **Severity:** MEDIUM
- **Notes:** Directly documented in PostgreSQL reference.
