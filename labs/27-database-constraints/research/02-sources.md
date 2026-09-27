# Sources Accessed

## Tier 1: PostgreSQL Official Documentation (v18.6)

### 1. PostgreSQL 18 - Chapter 5.5: Constraints
**URL:** https://www.postgresql.org/docs/current/ddl-constraints.html
**Published:** September 24, 2026 (latest stable)
**Accessed:** September 27, 2026
**Source Tier:** Tier 1 (Official documentation)
**Relevance:** Primary source for CHECK, NOT NULL, UNIQUE, FOREIGN KEY, EXCLUDE constraint semantics, syntax, and limitations

#### Key Evidence Extracted:
- Section 5.5.1 CHECK Constraints:
  - "A check constraint is satisfied if the check expression evaluates to true or the null value"
  - "CHECK expressions cannot contain subqueries nor refer to variables other than columns of the current row"
  - "PostgreSQL assumes that CHECK constraints' conditions are immutable"
  - "PostgreSQL does not support CHECK constraints that reference table data other than the new or updated row being checked"

- Section 5.5.2 NOT NULL Constraints:
  - "A not-null constraint is functionally equivalent to creating a check constraint CHECK (column_name IS NOT NULL), but in PostgreSQL creating an explicit not-null constraint is more efficient"
  - "In most database designs the majority of columns should be marked not null"

- Section 5.5.3 UNIQUE Constraints:
  - "By default, two null values are not considered equal in this comparison"
  - "NULLS NOT DISTINCT option modifies this and causes the index to treat nulls as equal"
  - "Adding a unique constraint will automatically create a unique btree index on the column or group of columns used in the constraint"
  - "Multi-column: This specifies that the combination of values in the indicated columns is unique across the whole table"

- Section 5.5.5 FOREIGN KEY Constraints:
  - "MATCH FULL: will not allow one column of a multicolumn foreign key to be null unless all foreign key columns are null"
  - "MATCH SIMPLE: allows any of the foreign key columns to be null"
  - "MATCH PARTIAL: is not yet implemented"
  - ON DELETE/ON UPDATE actions: NO ACTION, RESTRICT, CASCADE, SET NULL, SET DEFAULT
  - "A foreign key must reference columns that either are a primary key or form a unique constraint, or are columns from a non-partial unique index"

- Section 5.5.6 EXCLUSION Constraints:
  - "Exclusion constraints ensure that if any two rows are compared on the specified columns or expressions using the specified operators, at least one of these operator comparisons will return false or null"
  - "Adding an exclusion constraint will automatically create an index of the type specified in the constraint declaration"

### 2. PostgreSQL 18 - Chapter 11.6: Unique Indexes
**URL:** https://www.postgresql.org/docs/current/indexes-unique.html
**Published:** September 24, 2026
**Accessed:** September 27, 2026
**Source Tier:** Tier 1 (Official documentation)
**Relevance:** Technical details on unique index implementation and NULL handling

#### Key Evidence Extracted:
- "Currently, only B-tree indexes can be declared unique"
- "When an index is declared unique, multiple table rows with equal indexed values are not allowed"
- "By default, null values in a unique column are not considered equal, allowing multiple nulls in the column"
- "NULLS NOT DISTINCT option modifies this and causes the index to treat nulls as equal"
- "PostgreSQL automatically creates a unique index when a unique constraint or primary key is defined for a table"
- "There's no need to manually create indexes on unique columns; doing so would just duplicate the automatically-created index"

### 3. PostgreSQL 18 - Chapter 11.8: Partial Indexes
**URL:** https://www.postgresql.org/docs/current/indexes-partial.html
**Published:** September 24, 2026
**Accessed:** September 27, 2026
**Source Tier:** Tier 1 (Official documentation)
**Relevance:** Primary source for partial unique indexes used in soft-delete patterns

#### Key Evidence Extracted:
- "A partial index is an index built over a subset of a table; the subset is defined by a conditional expression"
- "To create a partial index that suits our example, use: CREATE UNIQUE INDEX tests_success_constraint ON tests (subject, target) WHERE success"
- "This specifies uniqueness only among rows that satisfy the predicate"
- "The predicate must match the conditions used in the queries that are supposed to benefit from the index"
- "PostgreSQL does not have a sophisticated theorem prover that can recognize mathematically equivalent expressions that are written in different forms"
- "A partial index can be used in a query only if the system can recognize that the WHERE condition of the query mathematically implies the predicate of the index"
- "Parameterized query clauses do not work with a partial index"
- "It is possible to allow only one null in a column by creating a unique partial index with an IS NULL restriction"

### 4. PostgreSQL 18 - Appendix A: PostgreSQL Error Codes
**URL:** https://www.postgresql.org/docs/current/errcodes-appendix.html
**Published:** September 24, 2026
**Accessed:** September 27, 2026
**Source Tier:** Tier 1 (Official documentation)
**Relevance:** Authority for SQLSTATE error codes returned by constraint violations

#### Key Evidence Extracted:
- Class 23 - Integrity Constraint Violation:
  - `23502` - not_null_violation
  - `23503` - foreign_key_violation
  - `23505` - unique_violation
  - `23514` - check_violation
  - `23P01` - exclusion_violation
  - `23001` - restrict_violation
- "such names are supplied in separate fields of the error report message" - constraint names, table names, column names all returned as structured fields
- "Applications that need to know which error condition has occurred should usually test the error code, rather than looking at the textual error message"

### 5. PostgreSQL 18 - Chapter 13.4: Data Consistency Checks at the Application Level
**URL:** https://www.postgresql.org/docs/current/applevel-consistency.html
**Published:** September 24, 2026
**Accessed:** September 27, 2026
**Source Tier:** Tier 1 (Official documentation)
**Relevance:** How SERIALIZABLE isolation prevents race conditions, relationship to constraint enforcement

#### Key Evidence Extracted:
- "Serializable transactions are just Repeatable Read transactions which add nonblocking monitoring for dangerous patterns of read/write conflicts"
- "When a pattern is detected which could cause a cycle in the apparent order of execution, one of the transactions involved is rolled back to break the cycle"
- "Serializable...is the recommended isolation level for ensuring consistency across multiple data items"
- "If one transaction reads a row, then another updates the row, then the first reads it again, it might not see the update" (for Read Committed)

### 6. PostgreSQL 18 - Chapter 13.3: Explicit Locking
**URL:** https://www.postgresql.org/docs/current/explicit-locking.html
**Published:** September 24, 2026
**Accessed:** September 27, 2026
**Source Tier:** Tier 1 (Official documentation)
**Relevance:** Lock acquisition during INSERT/UPDATE for constraint enforcement

#### Key Evidence Extracted:
- "Row EXCLUSIVE lock is acquired by UPDATE, DELETE, INSERT, and MERGE"
- "FOR UPDATE causes the rows retrieved by the SELECT statement to be locked as though for update. This prevents them from being locked, modified or deleted by other transactions until the current transaction ends"
- Lock conflict matrix shows ROW EXCLUSIVE conflicts with other writers

### 7. PostgreSQL 18 - Chapter 5.12: Table Partitioning Limitations
**URL:** https://www.postgresql.org/docs/current/ddl-partitioning.html
**Published:** September 24, 2026
**Accessed:** September 27, 2026
**Source Tier:** Tier 1 (Official documentation)
**Relevance:** Constraint limitations on partitioned tables

#### Key Evidence Extracted:
- "To create a unique or primary key constraint on a partitioned table, the partition keys must not include any expressions or function calls and the constraint's columns must include all of the partition key columns"
- "This limitation exists because the individual indexes making up the constraint can only directly enforce uniqueness within their own partitions"

### 8. PostgreSQL 18 - ALTER TABLE: NOT VALID and VALIDATE CONSTRAINT
**URL:** https://www.postgresql.org/docs/current/sql-altertable.html
**Published:** September 24, 2026
**Accessed:** September 27, 2026
**Source Tier:** Tier 1 (Official documentation)
**Relevance:** Production constraint migration strategy

#### Key Evidence Extracted:
- "ADD CONSTRAINT NOT VALID: skips table scan; constraint still enforced on inserts/updates"
- "VALIDATE CONSTRAINT: acquires only SHARE UPDATE EXCLUSIVE lock on the table being altered"
- "VALIDATE CONSTRAINT: scans the table to ensure there are no rows for which the constraint is not satisfied"

## Tier 1 (Cross-Database Verification): SQLite Documentation

### 9. SQLite - CREATE TABLE Documentation
**URL:** https://www.sqlite.org/lang_createtable.html
**Published:** April 30, 2025 (last revised)
**Accessed:** September 27, 2026
**Source Tier:** Tier 1 (Official documentation)
**Relevance:** Independent verification of constraint behavior, particularly NULL handling

#### Key Evidence Extracted:
- "SQLite supports UNIQUE, NOT NULL, CHECK and FOREIGN KEY constraints"
- "For the purposes of UNIQUE constraints, NULL values are considered distinct from all other values, including other NULLs"
- "CHECK constraints are only verified when the table is written, not when it is read"
- "UNIQUE and PRIMARY KEY constraints are implemented by creating a unique index in the database"
- "If an INSERT or UPDATE statement attempts to modify the table content so that two or more rows have identical primary key values, that is a constraint violation"
- "For the purposes of PRIMARY KEY constraints, NULL values are considered distinct from all other values, including other NULLs"

## Tier 2: Industry Reference (Stripe)

### 10. Stripe API Documentation - Idempotent Requests
**URL:** https://docs.stripe.com/api/idempotent_requests
**Published:** September 2026 (current)
**Accessed:** September 27, 2026
**Source Tier:** Tier 2 (Company documentation - industry standard)
**Relevance:** Reference for idempotency pattern, comparison with database-level idempotency

#### Key Evidence Extracted:
- "Stripe's idempotency works by saving the resulting status code and body of the first request made for any given idempotency key, regardless of whether it succeeds or fails"
- "Subsequent requests with the same key return the same result"
- "A client generates an idempotency key, which is a unique key that the server uses to recognize subsequent retries of the same request"
- "We suggest using V4 UUIDs, or another random string with enough entropy to avoid collisions"
- "Keys can be removed from the system automatically after they're at least 24 hours old"

## Tier 3: Secondary Sources (Not Yet Accessed)

### Blocked Sources
- MySQL 8.0 Reference Manual (constraints.html, create-table.html, unique-constraints.html) - HTTP 403 Forbidden
  - Cannot verify MySQL-specific constraint behavior, error codes, DEFERRABLE support

### Unreachable Sources
- CitusData blog post on race conditions - 404 Not Found
- Martin Fowler article on DB constraints - 404 Not Found