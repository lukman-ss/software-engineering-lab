# Code Audit Findings

Target Lab: labs/27-database-constraints

## Finding 1

Location: `internal/engine/engine.go:46-99`
Claimed Behavior: Atomic enforcement of NOT NULL, CHECK, UNIQUE, and PARTIAL UNIQUE constraints during insert operations.
Observed Implementation: Evaluates mandatory fields, check predicates, and map-based index uniqueness under `e.mu.Lock()`.
Assessment: PASS
Severity: LOW
Notes: Concurrency safety maintained via mutex synchronization. Auto-increment sequence generation and index updates happen atomically inside the lock.

## Finding 2

Location: `internal/engine/engine.go:102-120`
Claimed Behavior: Soft deletion updates `DeletedAt` timestamp and clears partial unique index entry.
Observed Implementation: Deletes entry from `activeEmails` map while retaining record in `users` map under `e.mu.Lock()`.
Assessment: PASS
Severity: LOW
Notes: Correctly enables soft-delete re-registration behavior matching partial unique index semantics (`WHERE deleted_at IS NULL`).

## Finding 3

Location: `internal/engine/engine.go:123-148`
Claimed Behavior: Referential integrity verification for foreign key relationships (`fk_orders_user`).
Observed Implementation: Checks presence of `o.UserID` in `e.users` before appending order to `e.orders`.
Assessment: PASS
Severity: LOW
Notes: Correctly surfaces SQLSTATE `23503` when parent key does not exist.

## Finding 4

Location: `internal/store/store.go:24-47` vs `internal/store/store.go:59-66`
Claimed Behavior: Unsafe application store exhibits race conditions under concurrent writes, whereas SafeStore relies on storage engine constraint guarantees.
Observed Implementation: `UnsafeStore.RegisterUser` reads current state, injects sleep yield to expose race window, and writes without index lock, while `SafeStore.RegisterUser` delegates to engine atomic constraint checks.
Assessment: PASS
Severity: LOW
Notes: Accurately isolates storage-level constraint guarantees from application-level check pitfalls.

## Finding 5

Location: `internal/dberr/errors.go:8-99`
Claimed Behavior: ANSI/PostgreSQL SQLSTATE taxonomy for integrity constraint violations and translation to domain errors.
Observed Implementation: Maps SQLSTATE `23502`, `23503`, `23505`, and `23514` into structured `ConstraintError` and translates them into domain messages.
Assessment: PASS
Severity: LOW
Notes: Error typing unwraps and inspects cleanly via `IsConstraintViolation`.
