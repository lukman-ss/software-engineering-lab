# Code Audit

## Finding 1

Location: `internal/engine/engine.go:46-99`
Claimed Behavior: Atomic evaluation of NOT NULL, CHECK, UNIQUE, and PARTIAL UNIQUE constraints within storage lock.
Observed Implementation: `Engine.InsertUser` acquires `e.mu.Lock()` and validates mandatory fields, age/status CHECK rules, and email uniqueness prior to inserting record into maps.
Assessment: PASS
Severity: LOW
Notes: Synchronization ensures atomic index/table writes and prevents data corruption.

## Finding 2

Location: `internal/engine/engine.go:102-120`
Claimed Behavior: Soft delete updates `DeletedAt` and removes user from `activeEmails` partial index.
Observed Implementation: `SoftDeleteUser` locks engine mutex, validates non-nil `DeletedAt`, deletes entry from `activeEmails`, and updates user record.
Assessment: PASS
Severity: LOW
Notes: Correctly models PostgreSQL `CREATE UNIQUE INDEX ... WHERE deleted_at IS NULL`.

## Finding 3

Location: `internal/engine/engine.go:123-148`
Claimed Behavior: NOT NULL, CHECK (total_cents > 0), and FOREIGN KEY (user_id exists in users table) validation on order creation.
Observed Implementation: `InsertOrder` checks `UserID != 0`, `TotalCents > 0`, and presence of `UserID` in `e.users` under lock.
Assessment: PASS
Severity: LOW
Notes: Foreign key referential integrity correctly enforced.

## Finding 4

Location: `internal/dberr/errors.go:83-100`
Claimed Behavior: Map storage constraint errors (SQLSTATE 23xxx) to clean domain errors.
Observed Implementation: `MapToDomainError` inspects `ConstraintError.Code` and converts raw SQL state errors into structured domain errors.
Assessment: PASS
Severity: LOW
Notes: Preserves constraint name and classification.

## Finding 5

Location: `internal/store/store.go:24-47`
Claimed Behavior: Application-level check without DB constraints suffers race conditions.
Observed Implementation: `UnsafeStore.RegisterUser` reads current user count without holding engine write lock, sleeps for 1ms to exaggerate context switch window, then calls `InsertUserUnsafe`.
Assessment: PASS
Severity: LOW
Notes: Effectively demonstrates vulnerable TOCTOU pattern.
