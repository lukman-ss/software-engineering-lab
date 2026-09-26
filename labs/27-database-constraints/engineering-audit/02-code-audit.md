# Code Audit Report

## Finding 1

Location: `internal/engine/engine.go:46-99` (`InsertUser`)
Claimed Behavior: Evaluates NOT NULL, CHECK, UNIQUE, and PARTIAL UNIQUE constraints atomically under mutex protection.
Observed Implementation: Locks write mutex `e.mu.Lock()` for duration of validations, PK allocation via atomic counter `e.userSeq.Add(1)`, map updates, and index updates. Returns structured `ConstraintError` mapped to PostgreSQL SQLSTATE codes (`23502`, `23514`, `23505`).
Assessment: PASS
Severity: LOW
Notes: Atomic mutex locking guarantees serializable execution of constraints, eliminating concurrency windows inside storage engine.

## Finding 2

Location: `internal/engine/engine.go:102-120` (`SoftDeleteUser`)
Claimed Behavior: Updates `DeletedAt` timestamp and updates partial index (`activeEmails`).
Observed Implementation: Locks mutex `e.mu.Lock()`, verifies user existence, verifies non-nil `deletedAt.DeletedAt`, deletes entry from `e.activeEmails`, updates user record.
Assessment: PASS
Severity: LOW
Notes: Accurately reflects PostgreSQL partial unique index behavior `WHERE deleted_at IS NULL` where soft-deleted rows are removed from index.

## Finding 3

Location: `internal/engine/engine.go:123-148` (`InsertOrder`)
Claimed Behavior: Evaluates NOT NULL (`user_id`), CHECK (`total_cents > 0`), and FOREIGN KEY (`users(id)`) constraints.
Observed Implementation: Verifies non-zero `o.UserID` (SQLSTATE `23502`), positive `o.TotalCents` (SQLSTATE `23514`), and verifies presence of `o.UserID` in `e.users` (SQLSTATE `23503`).
Assessment: PASS
Severity: LOW
Notes: FK checking correctly enforces referential integrity against parent entity store.

## Finding 4

Location: `internal/store/store.go:24-47` (`UnsafeStore.RegisterUser`) vs `internal/store/store.go:59-66` (`SafeStore.RegisterUser`)
Claimed Behavior: Demonstrates vulnerability of application-level check vs safety of database constraint enforcement.
Observed Implementation: `UnsafeStore` iterates through users without mutex isolation across check-and-insert steps, adding artificial delay `time.Sleep(1 * time.Millisecond)` to demonstrate read-then-write race condition window. `SafeStore` delegates atomic insertion directly to `Engine.InsertUser`.
Assessment: PASS
Severity: LOW
Notes: Clearly isolates architectural pattern difference between database-enforced constraints and app-level checks.

## Finding 5

Location: `internal/dberr/errors.go:75-100` (`IsConstraintViolation` and `MapToDomainError`)
Claimed Behavior: Maps SQLSTATE constraint errors to human-readable domain errors and provides type check helpers.
Observed Implementation: `errors.As` cleanly unwraps `ConstraintError` and maps SQLSTATE codes (`23505`, `23502`, `23514`, `23503`) to domain error messages with rule annotations.
Assessment: PASS
Severity: LOW
Notes: Error taxonomy adheres to PostgreSQL error classification standards.
