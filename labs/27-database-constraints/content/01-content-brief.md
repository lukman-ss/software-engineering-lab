# Content Brief

Topic:
Database constraints as the foundation for data integrity, race condition prevention, and schema invariants.

Target Reader:
Software engineers, backend developers, and database administrators who design data models and need to understand when to use database constraints vs application-level validation.

Problem:
Application-level validation alone is vulnerable to race conditions and cannot provide the same strong correctness guarantees as database-level constraints. Teams often rely solely on application validation, leading to data inconsistencies under concurrent workloads.

Core Mental Model:
Database constraints provide declarative, ACID-guaranteed enforcement that happens atomically during INSERT/UPDATE operations, eliminating the "check-then-act" vulnerability. Application validation serves as the first line for user experience, while database constraints serve as the last line of defense for data correctness.

Approved Research Status:
APPROVED

Approved Engineering Status:
APPROVED

Main Concepts:
- NOT NULL constraint (SQLSTATE 23502)
- CHECK constraint (SQLSTATE 23514)
- UNIQUE constraint (SQLSTATE 23505)
- FOREIGN KEY constraint (SQLSTATE 23503)
- PARTIAL UNIQUE INDEX (conditional uniqueness)
- Constraint violation error mapping
- Race condition prevention via atomic constraint enforcement

Verified Behaviors:
- NOT NULL rejects null/empty values for mandatory fields
- CHECK evaluates boolean predicates on row data
- UNIQUE prevents duplicate keys under concurrent inserts
- FOREIGN KEY prevents orphaned records
- PARTIAL UNIQUE INDEX allows value reuse after soft-delete while maintaining active uniqueness
- Concurrent inserts against UNIQUE constraint result in exactly one success and N-1 SQLSTATE 23505 errors
- Constraint violations map to specific SQLSTATE codes for proper error handling

Available Case Studies:
- User registration with email uniqueness (race condition prevention)
- Soft-delete pattern with partial unique index (email reuse)
- Order creation with foreign key integrity
- Age and status validation with CHECK constraints
- Mandatory field validation with NOT NULL constraints

Warnings:
- Constraints are limited to row-scoped logic; cross-row validation requires SERIALIZABLE isolation or triggers
- MySQL-specific behavior was not verified due to inaccessible documentation during research (documented as open question)
- Multi-column UNIQUE NULL handling edge cases require live testing (documented as open question)
- Implementation uses in-memory Go simulator, not live PostgreSQL (but faithfully reproduces SQLSTATE behavior)