# Implementation Notes

## Files Added
- `go.mod`: Module definition for `labs/27-database-constraints`.
- `internal/dberr/errors.go`: Defines SQLSTATE error codes (`23502`, `23503`, `23505`, `23514`, `23001`, `23P01`, `40001`), `ConstraintError` type, helper inspection functions, and domain error mapping.
- `internal/model/model.go`: Domain models (`User`, `Order`).
- `internal/engine/engine.go`: Thread-safe relational engine modeling declarative constraint enforcement (NOT NULL, CHECK, UNIQUE, FOREIGN KEY, and PARTIAL UNIQUE INDEX).
- `internal/store/store.go`: Application store demonstrating safe enforcement using engine constraints vs vulnerable patterns.
- `internal/store/store_test.go`: Comprehensive unit and concurrency race condition tests verifying each constraint behavior.
- `cmd/demo/main.go`: Interactive executable demonstrating constraint violations and concurrent race elimination.

## Core Design Decisions
- **Storage Engine Integrity**: Invariants are enforced atomically within the engine mutex/index layer during write operations rather than relying on application read-then-write checks.
- **SQLSTATE Error Standardization**: Follows ANSI SQL and PostgreSQL Class 23 error definitions (`23502` NOT NULL, `23503` FOREIGN KEY, `23505` UNIQUE, `23514` CHECK).
- **Partial Unique Index Semantics**: Models `CREATE UNIQUE INDEX ... WHERE deleted_at IS NULL` where index presence is conditional on row predicate evaluation, allowing multiple soft-deleted duplicates while strictly ensuring a single active row.

## Implementation-Specific Choices
- In-memory simulation: Built using pure Go standard library primitives (`sync.RWMutex`, `atomic.Int64`, map structures) without external database server dependencies to ensure 100% reproducible, zero-setup execution.

## Known Limitations
- Does not connect to live PostgreSQL server instances or external databases (in-memory simulator models SQLSTATE behaviors).
- EXCLUDE constraints with custom GiST distance operators are omitted in favor of core constraint types (NOT NULL, CHECK, UNIQUE, PK, FK, Partial Index).

## Trade-offs
- In-memory engine uses coarse-grained table locks rather than fine-grained B-tree page lock latches, prioritizing code readability and zero-dependency testing.

## What Is Demonstrated
- Enforcing NOT NULL constraints returning `23502`.
- Enforcing CHECK constraints for boundary validation (`age >= 18`, enum status) returning `23514`.
- Enforcing FOREIGN KEY constraints rejecting orphan records returning `23503`.
- Enforcing UNIQUE constraints under 50 concurrent goroutines preventing race conditions, returning `23505` on collisions.
- Partial Unique Index allowing soft-deleted record reuse.
- Domain error mapping decoupling storage SQLSTATE codes from user-facing error messages.

## What Is Not Demonstrated
- Distributed transaction constraints across separate database nodes.
- Full SQL parser or AST interpretation.
