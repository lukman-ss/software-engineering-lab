# Code Audit

## Finding 1 — Storage Engine Atomic Constraint Enforcement

Location: `internal/engine/engine.go:46-99`
Claimed Behavior: Atomic evaluation of NOT NULL, CHECK, UNIQUE, and PARTIAL UNIQUE constraints in simulated storage engine.
Observed Implementation: `InsertUser` acquires mutex lock `e.mu.Lock()` and evaluates constraints sequentially before committing row/indexes to state maps.
Assessment: PASS
Severity: LOW
Notes: Thread-safe in-memory simulation accurately mirrors transactional ACID table lock behavior.

## Finding 2 — Foreign Key Reference Verification

Location: `internal/engine/engine.go:123-149`
Claimed Behavior: Prevent inserting orders for non-existent users (SQLSTATE 23503).
Observed Implementation: `InsertOrder` verifies presence of `o.UserID` in `e.users` map within `e.mu.Lock()`.
Assessment: PASS
Severity: LOW
Notes: Enforces referential integrity appropriately.

## Finding 3 — Error Domain Mapping & SQLSTATE Classification

Location: `internal/dberr/errors.go:83-100`
Claimed Behavior: Map database constraint SQLSTATE errors to domain-level errors preserving constraint names.
Observed Implementation: `MapToDomainError` checks error type `*ConstraintError` via `errors.As` and formats domain error messages including constraint rules.
Assessment: PASS
Severity: LOW
Notes: Properly standardizes database errors for upstream consumption.

## Finding 4 — Unsafe vs Safe Store Race Window Simulation

Location: `internal/store/store.go:24-47`
Claimed Behavior: `UnsafeStore` demonstrates application-level check weakness against concurrent writes.
Observed Implementation: `UnsafeStore.RegisterUser` scans `s.eng.GetUser` without table lock, sleeps 1ms to simulate I/O delay, then executes `InsertUserUnsafe`.
Assessment: PASS
Severity: LOW
Notes: Effective demonstration of read-then-write race condition window when database-level constraint is missing.
