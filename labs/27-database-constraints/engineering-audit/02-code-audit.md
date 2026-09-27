# Code Audit

## Finding 1

Location: internal/dberr/errors.go:8-18
Claimed Behavior: SQLSTATE Class 23 error codes defined according to ANSI SQL / PostgreSQL standard (23502, 23503, 23505, 23514, 23001, 23P01, 40001).
Observed Implementation: Constants `SQLStateRestrictViolation`, `SQLStateNotNullViolation`, `SQLStateForeignKeyViolation`, `SQLStateUniqueViolation`, `SQLStateCheckViolation`, `SQLStateExclusionViolation`, and `SQLStateSerializationFail` match standard definitions.
Assessment: PASS
Severity: LOW
Notes: Correctly structured with wrapping, error unwrapping, and predicate helpers.

## Finding 2

Location: internal/engine/engine.go:46-99
Claimed Behavior: Declarative database constraints (NOT NULL, CHECK, UNIQUE, PARTIAL UNIQUE) evaluated atomically under engine mutex.
Observed Implementation: Evaluates mandatory field validation (NOT NULL), row-scoped boundary checks (CHECK), index presence checks (UNIQUE & PARTIAL UNIQUE) under `e.mu.Lock()` before assigning sequence IDs and committing to state.
Assessment: PASS
Severity: LOW
Notes: Atomicity and consistency guaranteed across all insert paths.

## Finding 3

Location: internal/engine/engine.go:101-120
Claimed Behavior: Soft-delete removes active keys from the partial index while preserving existing row records.
Observed Implementation: `SoftDeleteUser` locks state, verifies record existence, deletes from `e.activeEmails`, and updates `DeletedAt` timestamp on the stored record.
Assessment: PASS
Severity: LOW
Notes: Accurately replicates PostgreSQL partial unique index semantics (`WHERE deleted_at IS NULL`).

## Finding 4

Location: internal/engine/engine.go:123-148
Claimed Behavior: Foreign key validation verifies parent table presence prior to child insertion.
Observed Implementation: `InsertOrder` validates `o.UserID != 0`, evaluates `TotalCents > 0`, and checks parent `e.users[o.UserID]` under write lock before persisting order.
Assessment: PASS
Severity: LOW
Notes: Emulates referential integrity constraint enforcement.

## Finding 5

Location: internal/store/store.go:24-47
Claimed Behavior: Application-level check-then-act pattern is vulnerable to race conditions without storage constraints.
Observed Implementation: `UnsafeStore.RegisterUser` executes linear search across records, followed by a simulated context yield (`time.Sleep(1ms)`), and calls `InsertUserUnsafe`, reproducing duplicate insertions under concurrency.
Assessment: PASS
Severity: LOW
Notes: Accurately isolates the read-then-write concurrency flaw for educational comparison against safe database-enforced storage.
