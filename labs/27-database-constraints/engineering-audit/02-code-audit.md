# Code Audit

## Finding 1

Location: `internal/engine/engine.go:46-99`
Claimed Behavior: Enforces NOT NULL, CHECK, UNIQUE, and PARTIAL UNIQUE constraints atomically under concurrent writes.
Observed Implementation: Locks engine mutex (`e.mu.Lock()`), validates fields, checks index maps (`emailIndex`, `activeEmails`), assigns auto-increment ID (`e.userSeq`), and inserts row/indexes atomically.
Assessment: PASS
Severity: LOW
Notes: Synchronization via `RWMutex` correctly guarantees thread-safety and atomic constraint evaluation.

## Finding 2

Location: `internal/store/store.go:23-47`
Claimed Behavior: UnsafeStore demonstrates application-level check vulnerability to read-then-write race conditions.
Observed Implementation: Reads user count, iterates user table to find duplicate email, sleeps 1ms to exaggerate race window, then inserts bypassing engine constraints.
Assessment: PASS
Severity: LOW
Notes: Accurately simulates why application-level validation without database constraints fails under concurrency.

## Finding 3

Location: `internal/dberr/errors.go:83-100`
Claimed Behavior: Maps SQLState errors to domain errors preserving constraint rules.
Observed Implementation: Uses `errors.As` to inspect `ConstraintError` and maps `23505`, `23502`, `23514`, `23503` to structured domain error messages.
Assessment: PASS
Severity: LOW
Notes: Mapping is clean, idiomatic Go, and correctly retains constraint names.

## Finding 4

Location: `internal/engine/engine.go:102-120`
Claimed Behavior: SoftDeleteUser updates `DeletedAt` and updates partial index (`activeEmails`).
Observed Implementation: Verifies user exists, checks `deletedAt.DeletedAt != nil`, deletes email key from `e.activeEmails`, and updates user in map.
Assessment: PASS
Severity: LOW
Notes: Properly maintains partial index state when rows are soft-deleted.
