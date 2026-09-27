# Evidence Gathered

## 1. Race Condition in Check-Then-Act Pattern

### Claim
Application-level validation using `SELECT → exists() → INSERT` is vulnerable to race conditions under concurrent requests, leading to duplicate inserts.

### Evidence
```
-- Without constraint (vulnerable):
tx1: SELECT count(*) FROM items WHERE code = 'x' → 0
tx2: SELECT count(*) FROM items WHERE code = 'x' → 0
tx1: INSERT INTO items (code) VALUES ('x') → INSERT succeeds
tx2: INSERT INTO items (code) VALUES ('x') → INSERT succeeds (BUG)
```

### Source
Topic specification itself demonstrates the race condition; verified by PostgreSQL documentation showing constraints enforce at INSERT time via atomic index update.

### URL
https://www.postgresql.org/docs/current/ddl-constraints.html  
https://www.postgresql.org/docs/current/indexes-unique.html  

### Confidence
HIGH

### Corroborated By
PostgreSQL docs (INSERT/UPDATE constraint checking), SQLite docs (constraint enforcement on INSERT/UPDATE), PostgreSQL locking docs (ROW EXCLUSIVE lock on INSERT)

### Notes
The race occurs because two concurrent transactions both see zero existing rows before either inserts. Database constraint prevents this by checking uniqueness during the INSERT operation itself under row-level locks.

## 2. UNIQUE Constraint Prevents Race Conditions Atomically

### Claim
UNIQUE constraint prevents duplicate key insertion by atomically checking for duplicates during INSERT/UPDATE under row-level locks, ensuring only the first concurrent insert succeeds.

### Evidence
```
-- With UNIQUE constraint (safe):
tx1: INSERT INTO items (code) VALUES ('x') → INSERT succeeds
tx2: INSERT INTO items (code) VALUES ('x') → ERROR 23505 unique_violation
```

### Source
PostgreSQL docs explain: UNIQUE constraints use B-tree indexes; concurrent inserts with same key conflict during index update under ROW EXCLUSIVE lock, with conflict detection raising unique_violation.

### URL
https://www.postgresql.org/docs/current/ddl-constraints.html  
https://www.postgresql.org/docs/current/indexes-unique.html  
https://www.postgresql.org/docs/current/explicit-locking.html  

### Confidence
HIGH

### Corroborated By
PostgreSQL docs (unique index auto-creation, B-tree locking), PostgreSQL locking docs (ROW EXCLUSIVE lock modes), SQLite docs (constraint enforcement on INSERT)

### Notes
INSERT acquires ROW EXCLUSIVE lock; concurrent INSERT with same key value will block; conflict detected during B-tree index update; second INSERT fails with 23505 unique_violation.

## 3. UNIQUE Constraint Semantics: NULL Handling and Index Creation

### Claim
UNIQUE constraint prevents duplicate non-NULL values; by default treats NULLs as distinct (multiple NULLs allowed); NULLS NOT DISTINCT treats NULLs as equal; automatically creates unique B-tree index.

### Evidence
"By default, two null values are not considered equal in this comparison"  
"NULLS NOT DISTINCT option modifies this and causes the index to treat nulls as equal"  
"Adding a unique constraint will automatically create a unique btree index on the column or group of columns used in the constraint"

### Source
PostgreSQL DDL Constraints section 5.5.3 and Unique Indexes section 11.6

### URL
https://www.postgresql.org/docs/current/ddl-constraints.html#DDL-CONSTRAINTS-UNIQUE-CONSTRAINTS  
https://www.postgresql.org/docs/current/indexes-unique.html  

### Confidence
HIGH

### Corroborated By
PostgreSQL docs (DDL Constraints + Unique Indexes), SQLite docs (NULLs distinct in UNIQUE constraints)

### Notes
Multi-column UNIQUE: combination of values must be unique. Index auto-creation means no manual index needed for constraint enforcement.

## 4. NOT NULL Constraint: More Efficient Alternative to CHECK

### Claim
NOT NULL constraint is functionally equivalent to CHECK (col IS NOT NULL) but more efficient; majority of columns should be NOT NULL.

### Evidence
"A not-null constraint is functionally equivalent to creating a check constraint CHECK (column_name IS NOT NULL), but in PostgreSQL creating an explicit not-null constraint is more efficient"  
"In most database designs the majority of columns should be marked not null"

### Source
PostgreSQL DDL Constraints section 5.5.2

### URL
https://www.postgresql.org/docs/current/ddl-constraints.html#DDL-CONSTRAINTS-NOT-NULL  

### Confidence
HIGH

### Corroborated By
PostgreSQL docs (functional equivalence + efficiency tip), SQLite docs (NOT NULL constraint verification on INSERT/UPDATE)

### Notes
NOT NULL verified during INSERT/UPDATE; SQLite shows attempts to insert NULL into NOT NULL column cause constraint violation.

## 5. FOREIGN KEY Constraint: Referential Integrity Details

### Claim
FOREIGN KEY enforces referential integrity; MATCH FULL/SIMPLE determine NULL handling; ON DELETE actions include NO ACTION, RESTRICT, CASCADE, SET NULL, SET DEFAULT; does NOT automatically index referencing columns.

### Evidence
"MATCH FULL: will not allow one column of a multicolumn foreign key to be null unless all foreign key columns are null"  
"MATCH SIMPLE: allows any of the foreign key columns to be null"  
ON DELETE actions: NO ACTION (default), RESTRICT, CASCADE, SET NULL, SET DEFAULT  
"A foreign key must reference columns that either are a primary key or form a unique constraint, or are columns from a non-partial unique index"  
"Because this is not always needed, and there are many choices available on how to index, the declaration of a foreign key constraint does not automatically create an index on the referencing columns"

### Source
PostgreSQL DDL Constraints section 5.5.5

### URL
https://www.postgresql.org/docs/current/ddl-constraints.html#DDL-CONSTRAINTS-FK  

### Confidence
HIGH

### Corroborated By
PostgreSQL docs (all FK semantics), SQLite docs (FOREIGN KEY constraint enforcement)

### Notes
Referenced table columns must have PK/unique/index; referencing table does NOT get auto-index (important for performance with CASCADE).

## 6. CHECK Constraint: Row-Scoped and Immutable Assumption

### Claim
CHECK enforces boolean expression on row values; cannot contain subqueries; satisfied if expression evaluates to TRUE or NULL; PostgreSQL assumes immutability; cross-row validation requires triggers.

### Evidence
"A check constraint can also refer to several columns"  
"CHECK expressions cannot contain subqueries nor refer to variables other than columns of the current row"  
"PostgreSQL assumes that CHECK constraints' conditions are immutable"  
"PostgreSQL does not support CHECK constraints that reference table data other than the new or updated row being checked"  
"If what you desire is a one-time check against other rows at row insertion... a custom trigger can be used"

### Source
PostgreSQL DDL Constraints section 5.5.1

### URL
https://www.postgresql.org/docs/current/ddl-constraints.html#DDL-CONSTRAINTS-CHECK-CONSTRAINTS  

### Confidence
HIGH

### Corroborated By
PostgreSQL docs (CHECK limitations), SQLite docs (CHECK only verified on INSERT/UPDATE)

### Notes
CHECK verification timing: only on INSERT/UPDATE, not on SELECT. Immutability assumption justifies skipping re-validation. Cross-row CHECK appears to work in simple tests but fails under concurrent load.

## 7. Partial Unique Index for Soft-Delete Pattern

### Claim
CREATE UNIQUE INDEX ... WHERE deleted_at IS NULL enforces uniqueness only for active records, allowing reuse of values after soft-delete.

### Evidence
"CREATE UNIQUE INDEX tests_success_constraint ON tests (subject, target) WHERE success"  
"It is also possible to allow only one null in a column by creating a unique partial index with an IS NULL restriction"

### Source
PostgreSQL Partial Indexes section 11.8

### URL
https://www.postgresql.org/docs/current/indexes-partial.html  

### Confidence
HIGH

### Corroborated By
PostgreSQL docs (partial unique index examples), PostgreSQL locking docs (index update during INSERT)

### Notes
Example: `CREATE UNIQUE INDEX users_active_email_idx ON users (email) WHERE deleted_at IS NULL;` ensures only one active record per email.

## 8. Partial Index Planner Implication Rules

### Claim
Partial index can be used only if query's WHERE clause mathematically implies the index's predicate; parameterized queries don't work; only simple inequality implications recognized.

### Evidence
"A partial index can be used in a query only if the system can recognize that the WHERE condition of the query mathematically implies the predicate of the index"  
"PostgreSQL does not have a sophisticated theorem prover that can recognize mathematically equivalent expressions that are written in different forms"  
"The system can recognize simple inequality implications, for example 'x < 1' implies 'x < 2'"  
"Parameterized query clauses do not work with a partial index"

### Source
PostgreSQL Partial Indexes section 11.8

### URL
https://www.postgresql.org/docs/current/indexes-partial.html  

### Confidence
HIGH

### Corroborated By
PostgreSQL docs (predicate matching + parameterized query limitation)

### Notes
Example: Index `WHERE x < 2` can be used by query `WHERE x < 1` but NOT by `WHERE x < ?` (parameterized).

## 9. Constraint Violation Error Codes and Structured Fields

### Claim
Constraint violations return specific SQLSTATE class 23 error codes; constraint names, table names, column names supplied in separate fields of error report; applications should test SQLSTATE, not error text.

### Evidence
Class 23 - Integrity Constraint Violation:  
`23502` - not_null_violation  
`23503` - foreign_key_violation  
`23505` - unique_violation  
`23514` - check_violation  
`23P01` - exclusion_violation  
`23001` - restrict_violation  
"such names are supplied in separate fields of the error report message"

### Source
PostgreSQL Error Codes Appendix A

### URL
https://www.postgresql.org/docs/current/errcodes-appendix.html  

### Confidence
HIGH

### Corroborated By
PostgreSQL docs (error code definitions), PostgreSQL app-level consistency docs (error field availability)

### Notes
Error codes stable across versions and not localized; error text may vary by locale. Error detail fields provide structured access to constraint/table/column information.

## 10. Serializable Isolation for Multi-Row Invariants

### Claim
SERIALIZABLE isolation prevents multi-row race conditions by detecting dangerous read/write conflict patterns and rolling back one transaction.

### Evidence
"Serializable transactions are just Repeatable Read transactions which add nonblocking monitoring for dangerous patterns of read/write conflicts"  
"When a pattern is detected which could cause a cycle in the apparent order of execution, one of the transactions involved is rolled back to break the cycle"

### Source
PostgreSQL App-Level Consistency section 13.4.1

### URL
https://www.postgresql.org/docs/current/applevel-consistency.html  

### Confidence
HIGH

### Corroborated By
PostgreSQL docs (SERIALIZABLE definition + SSI), PostgreSQL error docs (40001 serialization_failure)

### Notes
Alternative to explicit locking (SELECT FOR UPDATE) for complex invariants; applications should retry on 40001 error.

## 11. Production Constraint Migration: NOT VALID + VALIDATE

### Claim
Adding constraint to production table: use NOT VALID to skip initial table scan, then VALIDATE CONSTRAINT to verify existing data with minimal locking.

### Evidence
"ADD CONSTRAINT NOT VALID: skips table scan; constraint still enforced on inserts/updates"  
"VALIDATE CONSTRAINT: acquires only SHARE UPDATE EXCLUSIVE lock on the table being altered"  
"If the constraint is a foreign key then a ROW SHARE lock is also required on the table referenced by the constraint"

### Source
PostgreSQL ALTER TABLE documentation

### URL
https://www.postgresql.org/docs/current/sql-altertable.html  

### Confidence
HIGH

### Corroborated By
PostgreSQL docs (NOT VALID + VALIDATE CONSTRAINT behavior)

### Notes
NOT VALID allows immediate constraint enforcement on new data; VALIDATE CONSTRAINT can run later with low lock contention to verify existing rows.

## 12. Cross-Database Verification: NULL Handling in UNIQUE Constraints

### Claim
PostgreSQL and SQLite both treat NULLs as distinct in UNIQUE constraints by default, allowing multiple NULLs; SQL Standard considers this implementation-defined.

### Evidence
PostgreSQL: "By default, two null values are not considered equal in this comparison"  
SQLite: "For the purposes of UNIQUE constraints, NULL values are considered distinct from all other values, including other NULLs"

### Source
PostgreSQL DDL Constraints 5.5.3 + SQLite CREATE TABLE documentation

### URL
https://www.postgresql.org/docs/current/ddl-constraints.html#DDL-CONSTRAINTS-UNIQUE-CONSTRAINTS  
https://www.sqlite.org/lang_createtable.html  

### Confidence
HIGH

### Corroborated By
PostgreSQL docs + SQLite docs (independent Tier 1 sources)

### Notes
Cross-verification shows consistent NULL handling between two major SQL implementations despite standard leaving it implementation-defined.

## 13. Stripe Idempotency Pattern vs Database-Level Idempotency

### Claim
Stripe's idempotency saves first request result per key; database-level idempotency for webhook deduplication uses UNIQUE constraint on reference_number + transaction.

### Evidence
Stripe: "Stripe's idempotency works by saving the resulting status code and body of the first request made for any given idempotency key"  
Database pattern: UNIQUE(reference_number) + transaction + idempotency key prevents duplicate webhook processing

### Source
Stripe API Documentation + Topic specification exercise

### URL
https://docs.stripe.com/api/idempotent_requests  

### Confidence
MEDIUM (Stripe verified, database pattern inferred from lab spec)

### Corroborated By
Stripe docs verified; database pattern logically follows from UNIQUE constraint atomicity

### Notes
Stripe is application-level idempotency; database UNIQUE provides storage-level deduplication. Combined pattern: idempotency key in application + UNIQUE constraint + transaction.

## 14. Partitioned Table Unique Constraint Limitation

### Claim
Unique constraint on partitioned table must include ALL partition key columns; cannot enforce cross-partition uniqueness without full partition key in constraint.

### Evidence
"To create a unique or primary key constraint on a partitioned table, the partition keys must not include any expressions or function calls and the constraint's columns must include all of the partition key columns"  
"This limitation exists because the individual indexes making up the constraint can only directly enforce uniqueness within their own partitions; therefore, the partition structure itself must guarantee that there are not duplicates in different partitions"

### Source
PostgreSQL Table Partitioning section 5.12.2.3

### URL
https://www.postgresql.org/docs/current/ddl-partitioning.html  

### Confidence
HIGH

### Corroborated By
PostgreSQL docs (partitioning limitations)

### Notes
Critical for partitioned table design: if partition key not in unique constraint, duplicates can exist across partitions even if each partition individually unique.

## 15. Constraint Cannot Replace Application Validation for UX

### Claim
Application validation remains necessary for user-friendly error messages; database constraint provides last-line defense for correctness.

### Evidence
"Application validation is useful for UX: 'Email already used'"  
"Database constraint is useful for correctness: 'Impossible to have two emails that are same'"  
"Use both: Application Validation → Friendly error, Database Constraint → Absolute protection"

### Source
Topic specification (best practices section)

### URL
N/A (reasoning from lab spec)

### Confidence
HIGH

### Corroborated By
Topic specification + PostgreSQL error handling docs (mapping constraint violations to domain errors)

### Notes
Constraint violations should be mapped to HTTP 409 Conflict with business-meaningful messages, not exposed as raw 500 errors.