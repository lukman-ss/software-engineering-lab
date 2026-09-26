# Engineering Design

Target Lab: labs/27-database-constraints
Research Status: APPROVED

## Concept To Prove
1. Declarative database constraints (NOT NULL, CHECK, UNIQUE, PRIMARY KEY, FOREIGN KEY, and PARTIAL UNIQUE INDEX) enforce relational integrity and invariants at the storage engine level.
2. Under concurrent write workloads (read-then-write race condition), application-only validation fails with duplicate key creation, while storage engine UNIQUE constraints guarantee consistency and reject conflicts with predictable SQLSTATE `23505` (`unique_violation`).
3. Partial unique index mechanics enforce conditional uniqueness (e.g., active vs soft-deleted records) allowing multiple logically deleted rows while strictly allowing only one active record per key.
4. Error mapping bridges standard database error codes (`23502`, `23503`, `23505`, `23514`) to domain-friendly validation errors.

## Expected Behavior
- NOT NULL constraint rejects insertions/updates with missing mandatory fields (`23502`).
- CHECK constraint evaluates row values against boolean expressions (e.g., `balance >= 0`, `status IN (...)`) and rejects invalid values (`23514`), while passing NULL if nullable.
- UNIQUE / PRIMARY KEY constraints reject duplicate keys (`23505`), preventing race condition concurrency bugs.
- FOREIGN KEY constraint rejects orphan records referencing non-existent parent rows (`23503`).
- Partial Unique Index allows duplicate records when inactive (`deleted_at IS NOT NULL`), but forbids duplicate active records (`deleted_at IS NULL`).
- Concurrency simulation: Multiple simultaneous goroutines trying to insert the same username/email without DB constraints lead to duplicate records. With DB constraints, exactly one succeeds and others receive `23505`.

## Failure Scenario
- Without database-level UNIQUE constraints, concurrent registration requests read `count == 0` simultaneously, pass app-level validation, and insert duplicate user records.
- With database-level constraints active, conflicting concurrent inserts are rejected by the engine lock/index layer and return `ErrUniqueViolation` (`SQLSTATE 23505`).

## Success Criteria
1. Full test coverage demonstrating NOT NULL, CHECK, UNIQUE, PRIMARY KEY, FOREIGN KEY, and Partial Index behavior.
2. Concurrency test proving that unsafe store suffers race condition duplicates while safe store (with DB constraints) guarantees exactly 1 record created and all concurrent duplicates return `23505`.
3. Error classification accurately categorizing `23502`, `23503`, `23505`, and `23514`.
4. Tests pass with `go test ./...` and `go test -race ./...`.
5. Standalone demo in `cmd/demo/main.go` runs with clear output demonstrating safe vs unsafe concurrency and constraint violations.

## Architecture
- `internal/db`: In-memory SQL engine simulation and constraint definitions using pure Go standard library (thread-safe relational engine modeling SQLite/PostgreSQL table and index locking semantics).
- `internal/store`:
  - `UnsafeStore`: Implements read-then-write checks in application memory without DB constraint enforcement.
  - `SafeStore`: Uses DB constraints to enforce integrity at the storage layer.
- `internal/errors`: SQLSTATE class 23 constraint error taxonomy and mappers (`23502`, `23503`, `23505`, `23514`).
- `cmd/demo`: Executable demonstration comparing race conditions and displaying constraint violation scenarios.

## Components
- `internal/engine`: Relational table engine with schema validation, NOT NULL, CHECK expressions, B-tree/Map primary & unique index checking, partial unique indexing, and foreign key referential integrity checks.
- `internal/domain`: Domain entities (e.g., `User`, `Account`, `Order`).
- `internal/service`: Application services demonstrating safe vs unsafe workflows.

## Test Strategy
- Unit tests for each constraint type:
  - `TestNotNullConstraint`
  - `TestCheckConstraint`
  - `TestUniqueConstraint`
  - `TestForeignKeyConstraint`
  - `TestPartialUniqueIndex`
  - `TestErrorMapping`
- Concurrency test:
  - `TestConcurrentRegistration_Unsafe_SuffersRaceCondition`
  - `TestConcurrentRegistration_Safe_EnforcesUniqueness`

## Execution Plan
1. Implement domain models and SQLSTATE error types (`internal/domain`, `internal/errors`).
2. Implement schema constraint engine (`internal/engine`).
3. Implement `UnsafeStore` and `SafeStore` (`internal/store`).
4. Write unit and concurrency tests (`internal/store/...` and `tests/...`).
5. Build runnable CLI demo (`cmd/demo/main.go`).
6. Execute tests and race detector.
7. Record results in `engineering/02-implementation-notes.md` and `engineering/03-execution-result.md`.
8. Write comprehensive `README.md`.

## Implementation Decisions
- Implementation uses a standard library pure Go in-memory ACID-style engine simulating SQL constraints, index lock detection, and SQLSTATE error returns. This eliminates external server dependencies (like Docker PostgreSQL) while faithfully reproducing PostgreSQL SQLSTATE `23xxx` codes and concurrency locking behavior.
