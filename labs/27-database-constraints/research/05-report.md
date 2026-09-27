# Research Report: Database Constraints — Jangan Serahkan Integritas Data Hanya ke Application Code

## Research Question

When business requirements demand that certain conditions must always hold true for application data, should enforcement live in application code alone, or should the database itself enforce these invariants? The research investigates whether database constraints provide stronger correctness guarantees than application-only validation, with particular focus on race condition prevention, partial indexes for soft-delete patterns, idempotency via UNIQUE constraints, and the boundary of when constraints are appropriate.

## Executive Summary

Database constraints provide **declarative, ACID-guaranteed** data integrity enforcement that is strictly stronger than application-level validation alone. The key finding: constraints prevent race conditions atomically at INSERT time via B-tree index updates under row-level locks, eliminating the "check-then-act" vulnerability that affects application-only patterns. Constraints are limited to **row-scoped logic** — cross-row validation and external service checks require triggers or SERIALIZABLE isolation. The evidence strongly supports the thesis that **database should be the last line of defense for data integrity**, while application validation serves as the first line for user experience.

**Critical limitation identified:** Constraints cannot express business rules involving other tables (cross-row), mutable functions, or external services — these require alternative approaches.

---

## Findings

### Finding 1: Application-Only Validation Is Vulnerable to Race Conditions

**Claim:** The check-then-act pattern (`SELECT EXISTS() → INSERT`) fails under concurrent requests because two transactions can both see "no existing record" before either inserts.

**Evidence:**
The race condition proceeds as follows:
```
Request A                    Request B
  │                           │
  ├─ SELECT → 0 rows         │
  │                           ├─ SELECT → 0 rows
  │                           │
  ├─ INSERT succeeds          │
  │                           ├─ INSERT succeeds (DUPLICATE)
```

PostgreSQL documentation confirms that constraints are checked "when rows are inserted or updated" — the enforcement happens atomically during the INSERT operation, not during a preceding SELECT. The UNIQUE constraint uses a B-tree index that is updated atomically with the row insert under `ROW EXCLUSIVE` lock [PostgreSQL DDL Constraints 5.5, Explicit Locking 13.3].

**Source:** PostgreSQL DDL Constraints (5.5), Explicit Locking (13.3)
**URL:** https://www.postgresql.org/docs/current/ddl-constraints.html
**Confidence:** HIGH — cross-verified against PostgreSQL docs and SQLite docs (constraint enforcement timing)

---

### Finding 2: UNIQUE Constraint Atomically Prevents Duplicate Inserts

**Claim:** A `UNIQUE` constraint on `(user_id, voucher_id)` ensures only the first concurrent INSERT succeeds; subsequent inserts fail with SQLSTATE `23505` (`unique_violation`).

**Evidence:**
```
-- With UNIQUE constraint (safe):
tx1: INSERT INTO voucher_claims (user_id, voucher_id) VALUES (7, 10) → INSERT succeeds
tx2: INSERT INTO voucher_claims (user_id, voucher_id) VALUES (7, 10) → ERROR 23505 unique_violation
```

Mechanism:
1. INSERT acquires `ROW EXCLUSIVE` lock on the table [PostgreSQL Explicit Locking 13.3]
2. The B-tree index for the UNIQUE constraint is updated atomically with the row insert
3. Concurrent INSERT with identical key value blocks, and conflict detected during index update
4. Second INSERT raises `unique_violation` (SQLSTATE `23505`) [PostgreSQL Error Codes Appendix A]

PostgreSQL docs confirm: "Adding a unique constraint will automatically create a unique B-tree index on the column or group of columns used in the constraint" [PostgreSQL DDL Constraints 5.5.3].

**Source:** PostgreSQL DDL Constraints (5.5.3), Error Codes (Appendix A), Explicit Locking (13.3)
**URL:** https://www.postgresql.org/docs/current/indexes-unique.html
**Confidence:** HIGH — verified via PostgreSQL 18 documentation on constraint behavior, locking, and error codes

---

### Finding 3: NOT NULL Constraint Is More Efficient Than CHECK (col IS NOT NULL)

**Claim:** `NOT NULL` constraint is functionally equivalent to `CHECK (column IS NOT NULL)` but PostgreSQL implements it more efficiently.

**Evidence:**
PostgreSQL documentation explicitly states: "A not-null constraint is functionally equivalent to creating a check constraint CHECK (column_name IS NOT NULL), but in PostgreSQL creating an explicit not-null constraint is more efficient" [PostgreSQL DDL Constraints 5.5.2]. The documentation also recommends: "In most database designs the majority of columns should be marked not null" [PostgreSQL DDL Constraints 5.5.2].

Cross-verification via SQLite confirms both databases verify NOT NULL during INSERT/UPDATE and reject NULL values [SQLite CREATE TABLE documentation].

**Source:** PostgreSQL DDL Constraints (5.5.2), SQLite CREATE TABLE docs
**URL:** https://www.postgresql.org/docs/current/ddl-constraints.html
**Confidence:** HIGH — verified in both PostgreSQL and SQLite documentation

---

### Finding 4: FOREIGN KEY Provides Referential Integrity Without Automatic Indexing on Referencing Columns

**Claim:** `FOREIGN KEY` enforces that `customer_id` in `work_orders` must reference a valid `customer`; referencing columns do NOT automatically get indexed.

**Evidence:**
PostgreSQL documentation states: "the declaration of a foreign key constraint does not automatically create an index on the referencing columns" [PostgreSQL DDL Constraints 5.5.5]. The referenced table columns must have a primary key, unique constraint, or non-partial unique index.

This is critical for **orphaned data prevention**: without foreign keys, the database can contain `work_orders.customer_id = 9382` pointing to a non-existent customer. The foreign key prevents this.

Performance note: FK enforcement requires scanning the referencing table during DELETE/UPDATE on the referenced table — indexing referencing columns is recommended for performance but is not automatic.

**Source:** PostgreSQL DDL Constraints (5.5.5)
**URL:** https://www.postgresql.org/docs/current/ddl-constraints.html
**Confidence:** HIGH — verified in PostgreSQL documentation

---

### Finding 5: CHECK Constraints Are Row-Scoped and Cannot Reference Other Tables

**Claim:** `CHECK` constraints can only reference columns of the current row; cross-row validation requires triggers or alternative approaches.

**Evidence:**
PostgreSQL documentation explicitly states: "PostgreSQL does not support CHECK constraints that reference table data other than the new or updated row being checked. While a CHECK constraint that violates this rule may appear to work in simple tests, it cannot guarantee that the database will not reach a state in which the constraint condition is false" [PostgreSQL DDL Constraints 5.5.1].

Additionally: "PostgreSQL assumes that CHECK constraints' conditions are immutable" — meaning the database trusts the condition won't change between INSERTs, justifying the optimization of checking only at INSERT/UPDATE time [PostgreSQL DDL Constraints 5.5.1].

**Cross-verification:** SQLite documentation confirms "CHECK constraints are only verified when the table is written, not when it is read" [SQLite CREATE TABLE].

**Source:** PostgreSQL DDL Constraints (5.5.1), SQLite CREATE Table docs
**URL:** https://www.postgresql.org/docs/current/ddl-constraints.html
**Confidence:** HIGH — verified in both PostgreSQL and SQLite documentation

**Practical implication:** Complex business rules like "Customer VIP gets 20% discount only if total transactions last 12 months > Rp50 juta" cannot be expressed as a CHECK constraint. These belong in the application/domain layer.

---

### Finding 6: Partial Unique Indexes Solve the Soft-Delete Uniqueness Problem

**Claim:** A partial unique index `CREATE UNIQUE INDEX users_email_active_idx ON users(email) WHERE deleted_at IS NULL` enforces uniqueness only for active records, allowing email reuse after soft-delete.

**Evidence:**
PostgreSQL documentation provides an equivalent example: "CREATE UNIQUE INDEX tests_success_constraint ON tests (subject, target) WHERE success" — "This is a particularly efficient approach when there are few successful tests and many unsuccessful ones" [PostgreSQL Partial Indexes 11.8].

The documentation also confirms: "It is also possible to allow only one null in a column by creating a unique partial index with an IS NULL restriction" [PostgreSQL Partial Indexes 11.8].

**Critical limitation:** The partial index predicate must match the query's WHERE clause exactly. PostgreSQL's planner "does not have a sophisticated theorem prover that can recognize mathematically equivalent expressions that are written in different forms" and "parameterized query clauses do not work with a partial index" [PostgreSQL Partial Indexes 11.8].

**Source:** PostgreSQL Partial Indexes (11.8)
**URL:** https://www.postgresql.org/docs/current/indexes-partial.html
**Confidence:** HIGH — verified in PostgreSQL documentation with explicit examples

---

### Finding 7: Constraint Violations Map to Specific SQLSTATE Codes; Applications Should Use Error Codes, Not Text

**Claim:** Constraint violations return structured SQLSTATE error codes with constraint/table/column names in separate error fields; applications should test SQLSTATE codes and map them to user-friendly HTTP responses.

**Evidence:**
PostgreSQL Error Codes documentation [Appendix A] provides the following class 23 codes:
- `23502` — `not_null_violation`
- `23503` — `foreign_key_violation`
- `23505` — `unique_violation`
- `23514` — `check_violation`
- `23P01` — `exclusion_violation`
- `23001` — `restrict_violation`

The documentation states: "such names are supplied in separate fields of the error report message so that applications need not try to extract them from the possibly-localized human-readable text of the message" [PostgreSQL Error Codes Appendix A].

Best practice: Applications should check SQLSTATE, not error text. Named constraints (e.g., `CONSTRAINT email_unique UNIQUE (email)`) appear in error details, enabling mapping to domain errors like HTTP `409 Conflict` with message "Voucher sudah pernah diklaim."

**Source:** PostgreSQL Error Codes (Appendix A)
**URL:** https://www.postgresql.org/docs/current/errcodes-appendix.html
**Confidence:** HIGH — verified in PostgreSQL documentation

---

### Finding 8: SERIALIZABLE Isolation Handles Multi-Row Invariants That Constraints Cannot

**Claim:** For invariants involving multiple rows (e.g., "sum of balances must be zero"), database constraints are insufficient; `SERIALIZABLE` isolation level with retry logic is the appropriate solution.

**Evidence:**
PostgreSQL documentation states: "Serializable transactions are just Repeatable Read transactions which add nonblocking monitoring for dangerous patterns of read/write conflicts. When a pattern is detected which could cause a cycle in the apparent order of execution, one of the transactions involved is rolled back to break the cycle" [PostgreSQL App-Level Consistency 13.4].

The documentation also warns: "It is very difficult to enforce business rules regarding data integrity using Read Committed transactions because the view of the data is shifting with each statement" [PostgreSQL App-Level Consistency 13.4].

`SERIALIZABLE` violations return SQLSTATE `40001` (`serialization_failure`) which applications should handle with retry logic.

**Source:** PostgreSQL App-Level Consistency (13.4), Error Codes Appendix A
**URL:** https://www.postgresql.org/docs/current/applevel-consistency.html
**Confidence:** HIGH — verified in PostgreSQL documentation

---

### Finding 9: Idempotency via UNIQUE Constraint + Transaction Prevents Duplicate Webhook Processing

**Claim:** A `UNIQUE` constraint on `reference_number` combined with a transaction prevents duplicate webhook processing from payment providers.

**Evidence:**
Stripe documentation describes idempotency: "Stripe's idempotency works by saving the resulting status code and body of the first request made for any given idempotency key, regardless of whether it succeeds or fails. Subsequent requests with the same key return the same result" [Stripe API Documentation].

The database-level pattern extends this:
1. Each webhook includes a unique `reference_number`
2. Application creates transaction: `INSERT INTO wallet_transactions (reference_number, ...) VALUES (...)`
3. `UNIQUE(reference_number)` constraint ensures second identical webhook fails with `23505`
4. Application catches `unique_violation` and returns `200 OK` with previous result (idempotent)

This combines Stripe's idempotency key concept (application-level) with database-level UNIQUE constraint (storage-level) for defense-in-depth.

**Source:** Stripe API Documentation (idempotency), PostgreSQL DDL Constraints + Error Codes
**URL:** https://docs.stripe.com/api/idempotent_requests, https://www.postgresql.org/docs/current/ddl-constraints.html
**Confidence:** MEDIUM — Stripe idempotency verified; database pattern inferred from lab specification and constraint atomicity evidence

---

### Finding 10: Production Constraint Addition Requires NOT VALID + VALIDATE CONSTRAINT Strategy

**Claim:** Adding constraints to production tables should use `NOT VALID` to skip initial table scan, followed by `VALIDATE CONSTRAINT` to verify existing data with minimal locking.

**Evidence:**
PostgreSQL ALTER TABLE documentation states: "With NOT VALID, the ADD CONSTRAINT command does not scan the table and can be committed immediately. After that, a VALIDATE CONSTRAINT command can be issued to verify that existing rows satisfy the constraint. The validation step does not need to lock out concurrent updates" [PostgreSQL ALTER TABLE].

`VALIDATE CONSTRAINT` acquires only `SHARE UPDATE EXCLUSIVE` lock — significantly less contention than the `ACCESS EXCLUSIVE` lock required by a normal `ADD CONSTRAINT`.

**Critical warning for production:** "Adding a constraint to a production table without checking for existing violations can cause migration failure" (topic specification claim). The NOT VALID + VALIDATE approach solves this.

**Source:** PostgreSQL ALTER TABLE documentation
**URL:** https://www.postgresql.org/docs/current/sql-altertable.html
**Confidence:** HIGH — verified in PostgreSQL documentation

---

### Finding 11: Partitioned Tables Have Unique Constraint Limitations

**Claim:** Unique constraints on partitioned tables must include ALL partition key columns; cross-partition uniqueness enforcement is impossible without full partition key in constraint.

**Evidence:**
PostgreSQL Partitioning documentation states: "To create a unique or primary key constraint on a partitioned table, the partition keys must not include any expressions or function calls and the constraint's columns must include all of the partition key columns. This limitation exists because the individual indexes making up the constraint can only directly enforce uniqueness within their own partitions; therefore, the partition structure itself must guarantee that there are not duplicates in different partitions" [PostgreSQL Partitioning 5.12.2.3].

**Practical implication:** If `wallet_transactions` is partitioned by `created_at` month, a `UNIQUE(reference_number)` constraint is NOT valid — you must include `created_at` in the constraint: `UNIQUE(reference_number, created_at)`.

**Source:** PostgreSQL Partitioning (5.12.2.3)
**URL:** https://www.postgresql.org/docs/current/ddl-partitioning.html
**Confidence:** HIGH — verified in PostgreSQL documentation

---

## Areas of Agreement

1. **Database constraints provide stronger correctness guarantees than application-only validation** — all sources agree constraints enforce invariants at the storage layer, making them the last line of defense.

2. **UNIQUE constraints prevent race conditions atomically** — PostgreSQL and SQLite both enforce uniqueness during INSERT under row-level locks; application-level exists() check cannot replicate this.

3. **Constraints are limited to row-scoped logic** — CHECK cannot reference other rows; cross-row invariants require SERIALIZABLE isolation or triggers.

4. **Partial unique indexes are the correct solution for soft-delete uniqueness** — PostgreSQL documentation provides explicit examples matching the lab specification.

5. **Application validation and database constraints serve complementary roles** — application for UX (friendly errors), database for correctness (absolute protection).

6. **Named constraints and SQLSTATE codes enable proper error mapping** — applications should test error codes, not error text.

7. **Production migrations require careful data validation before adding constraints** — NOT VALID + VALIDATE CONSTRAINT is the recommended pattern.

## Areas of Disagreement

**No material disagreements discovered** between primary sources (PostgreSQL 18, SQLite). Both implement consistent constraint behavior for the topics investigated.

**MySQL-specific behavior remains unverified** due to HTTP 403 on dev.mysql.com. MySQL may have different NULL handling in UNIQUE constraints and may not support DEFERRABLE constraints — these claims remain uncertain pending verification.

## Limitations

1. **Single-database focus:** All primary evidence is from PostgreSQL 18 documentation. SQLite was used for cross-verification but is a different system. MySQL, Oracle, SQL Server behavior not verified.

2. **Locking internals not fully documented:** PostgreSQL does not fully document the exact B-tree page lock acquisition sequence during UNIQUE constraint enforcement. The behavior is inferred from general locking documentation (ROW EXCLUSIVE on INSERT) and the atomicity guarantee of constraint checking.

3. **Partial index planner implication depth unknown:** Documentation says PostgreSQL recognizes "simple inequality implications" but does not provide the complete list or algorithm.

4. **NULL in multi-column UNIQUE constraints not explicitly documented:** Behavior for `(1, NULL)` vs `(1, NULL)` is strongly implied but not directly stated in documentation.

5. **Stripe idempotency is application-level:** Stripe's mechanism is not database-level; the database pattern inferred from the lab specification.

## Conclusion

The evidence strongly supports the central thesis of the lab specification: **data integrity enforcement should live at the database layer, not only in application code.** Database constraints provide ACID-guaranteed enforcement that eliminates race conditions through atomic index updates under row-level locks — a guarantee that application-level `exists()` checks cannot replicate.

The evidence also clarifies the **boundary of constraint applicability**: constraints are powerful for row-scoped invariants (UNIQUE, NOT NULL, CHECK, FK) but insufficient for cross-row, cross-table, or external-service-dependent rules. The senior engineer's role is to identify which rules belong in constraints and which belong in the application layer — and to design the defense-in-depth architecture where application validation provides UX and database constraints provide absolute correctness protection.

The recommended architecture for the `wallet_transactions` table example from the lab:
- `UNIQUE(reference_number)` — prevents duplicate payment recording
- `CHECK(amount > 0)` — prevents negative amounts
- `FOREIGN KEY(user_id) REFERENCES users(id)` — prevents orphaned records
- `CHECK(type IN ('CREDIT', 'DEBIT'))` — restricts to valid types
- `CHECK(status IN ('PENDING', 'SUCCESS', 'FAILED'))` — restricts to valid statuses
- `NOT NULL` on required columns
- Partial unique index on `(user_id, reference_number) WHERE status = 'PENDING'` if idempotency requires

This design ensures that **even if all application validation is bypassed**, the database rejects data that violates business invariants.
