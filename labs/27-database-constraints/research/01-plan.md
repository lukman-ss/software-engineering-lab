# Research Plan: Database Constraints — Enforcing Data Integrity at the Storage Layer

## Research Topic

Database Constraints — Jangan Serahkan Integritas Data Hanya ke Application Code. Why data integrity enforcement should live at the database layer, not just in application validation. Covers UNIQUE, NOT NULL, FOREIGN KEY, CHECK, partial unique indexes, idempotency, race condition prevention, error handling, and production migration strategies.

## Objective

Investigate and verify the claims made in the lab topic specification:

1. Application-only validation (exists() + insert) is vulnerable to race conditions
2. Database UNIQUE constraint atomically prevents duplicate inserts
3. Each constraint type (NOT NULL, FOREIGN KEY, CHECK, UNIQUE) enforces specific invariants
4. Partial unique indexes solve soft-delete uniqueness patterns
5. Idempotency keys + UNIQUE constraints prevent duplicate processing in webhook scenarios
6. Constraint violations must be mapped to proper HTTP status codes, not generic 500
7. Constraints cannot handle all business rules — knowing the boundary is critical
8. Production constraint addition requires careful migration strategy (NOT VALID + VALIDATE)

## Research Questions

### RQ1: Race Conditions with Check-Then-Act Pattern
- How does SELECT → exists() → INSERT fail under concurrency?
- What is the exact mechanism that causes two concurrent inserts to succeed?
- At what isolation level does the race still occur?

### RQ2: UNIQUE Constraint as Race Condition Prevention
- How does the database atomically enforce uniqueness during concurrent inserts?
- What lock mechanism prevents duplicate key insertion?
- What SQLSTATE code is returned on violation?
- How does this compare to application-level exists() check?

### RQ3: Constraint Types and Their Semantics
- What does each constraint type enforce (NOT NULL, CHECK, UNIQUE, FK, EXCLUDE)?
- NULL handling differences across constraint types
- Which constraints automatically create indexes?
- What are the deferrability options?

### RQ4: Partial Unique Indexes for Soft Delete
- How does CREATE UNIQUE INDEX ... WHERE deleted_at IS NULL work?
- How does NULLS NOT DISTINCT interact with partial indexes?
- What are the planner implication rules?

### RQ5: Error Handling and Mapping
- What SQLSTATE codes map to each constraint type?
- How should applications translate constraint violations to HTTP responses?
- What information does PostgreSQL include in constraint violation error reports?

### RQ6: Idempotency with UNIQUE + Transaction
- How does a UNIQUE constraint on reference_number prevent duplicate webhook processing?
- How does Stripe's idempotency model compare to database-level idempotency?
- What is the pattern: UNIQUE + idempotency key + transaction?

### RQ7: When Constraints Cannot Help
- What business rules cannot be expressed as CHECK constraints?
- Performance implications of constraints
- Production migration strategy: NOT VALID + VALIDATE CONSTRAINT
- Partitioned table constraint limitations

## Search Strategy

### Primary Sources (Tier 1)
- PostgreSQL 18 official documentation (constraints, indexes, locking, error codes, partitioning)
- SQLite official documentation (independent cross-check for NULL handling and constraint enforcement)
- Stripe API documentation (idempotency pattern from industry reference)

### Secondary Sources (Tier 2)
- Industry articles on database constraints and data integrity

### Method
1. Fetch PostgreSQL docs for each constraint type — extract exact quotes with section references
2. Cross-check NULL handling claims against SQLite docs (different implementation)
3. Verify idempotency pattern against Stripe's official documentation
4. Verify ALTER TABLE NOT VALID / VALIDATE CONSTRAINT for production migration
5. Compare claims against SQLite's independent implementation for cross-database validation

## Expected Primary Sources

| Source | Topic | URL |
|--------|-------|-----|
| PostgreSQL 18 Docs - DDL Constraints | CHECK, NOT NULL, UNIQUE, PK, FK, EXCLUDE | https://www.postgresql.org/docs/current/ddl-constraints.html |
| PostgreSQL 18 Docs - Partial Indexes | Partial unique indexes, planner implication | https://www.postgresql.org/docs/current/indexes-partial.html |
| PostgreSQL 18 Docs - Unique Indexes | NULLS NOT DISTINCT, auto-index creation | https://www.postgresql.org/docs/current/indexes-unique.html |
| PostgreSQL 18 Docs - Error Codes | SQLSTATE class 23, constraint violation codes | https://www.postgresql.org/docs/current/errcodes-appendix.html |
| PostgreSQL 18 Docs - Explicit Locking | Row-level locks, INSERT lock behavior | https://www.postgresql.org/docs/current/explicit-locking.html |
| PostgreSQL 18 Docs - App-Level Consistency | Serializability, SSI, read/write conflicts | https://www.postgresql.org/docs/current/applevel-consistency.html |
| PostgreSQL 18 Docs - ALTER TABLE | NOT VALID, VALIDATE CONSTRAINT | https://www.postgresql.org/docs/current/sql-altertable.html |
| PostgreSQL 18 Docs - Partitioning | Unique constraint limitations on partitions | https://www.postgresql.org/docs/current/ddl-partitioning.html |
| SQLite Docs - CREATE TABLE | Cross-check: UNIQUE, NULL behavior, CHECK enforcement | https://www.sqlite.org/lang_createtable.html |
| Stripe API Docs - Idempotent Requests | Idempotency key pattern | https://docs.stripe.com/api/idempotent_requests |

## Risks / Unknowns

1. **MySQL behavior unverifiable**: MySQL documentation returned HTTP 403 during research; MySQL-specific NULL handling, DEFERRABLE support, and error codes not directly verified
2. **NULL semantics subtlety**: Different databases handle NULL in UNIQUE constraints differently (PostgreSQL: NULLs distinct; SQLite: NULLs distinct; SQL Standard: implementation-defined)
3. **Locking internals**: PostgreSQL doesn't fully document internal B-tree page lock behavior during UNIQUE constraint enforcement — the exact lock acquisition sequence during concurrent inserts is inferred from lock mode documentation
4. **Stripe idempotency vs database idempotency**: Stripe's implementation is application-level; the lab's exercise asks about database-level idempotency via UNIQUE constraints, which is a different mechanism
