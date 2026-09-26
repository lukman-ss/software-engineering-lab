# Code Audit

Target Lab: labs/27-database-constraints

## Finding 1

Location: `internal/engine/engine.go:33-86`
Claimed Behavior: Atomically checks NOT NULL, CHECK, UNIQUE, and PARTIAL UNIQUE constraints during record insertion and returns standardized SQLSTATE errors.
Observed Implementation: Method acquires `e.mu.Lock()` on table entry, verifies NOT NULL invariants on `Email` and `Username`, evaluates CHECK constraints (`Age >= 18` and `Status` enum), checks index maps (`emailIndex` or `activeEmails`), assigns auto-incrementing ID atomically via `e.userSeq.Add(1)`, and populates storage maps.
Assessment: PASS
Severity: LOW
Notes: Synchronization is coarse-grained table-level locking (`sync.RWMutex`), fully safe for in-memory model invariants.

## Finding 2

Location: `internal/engine/engine.go:109-135`
Claimed Behavior: Enforces referential integrity (FOREIGN KEY) and boundary checks on order insertion.
Observed Implementation: Evaluates NOT NULL on `UserID`, CHECK constraint on `TotalCents > 0`, and FOREIGN KEY referential existence against `e.users[o.UserID]` inside mutex lock before inserting.
Assessment: PASS
Severity: LOW
Notes: Properly returns `SQLStateForeignKeyViolation` (`23503`) when parent user does not exist.

## Finding 3

Location: `internal/dberr/errors.go:8-18, 83-99`
Claimed Behavior: Implements standard ANSI SQL / PostgreSQL SQLSTATE Class 23 error taxonomy and domain mapper.
Observed Implementation: Defines `23001`, `23502`, `23503`, `23505`, `23514`, `23P01`, and `40001`. `MapToDomainError` unwraps `ConstraintError` and translates codes into domain messages.
Assessment: PASS
Severity: LOW
Notes: Implements `errors.As` unwrapping cleanly via `Unwrap()` method on `ConstraintError`.

## Finding 4

Location: `internal/store/store.go:23-43`
Claimed Behavior: `UnsafeStore` exposes race vulnerability in read-then-write check.
Observed Implementation: `UnsafeStore.RegisterUser` checks email presence iterating `e.eng.GetUsersCount()` without transaction locking. However, on line 41, it calls `s.eng.InsertUser(u, true)` (which activates the partial unique index check inside `engine.go:58-64`). While the comment states `// partial unique deactivated for unsafe demo`, calling `InsertUser(u, true)` actually evaluates `activeEmails` index in the engine if `DeletedAt == nil`, unless the caller specifically passes a custom non-indexing mechanism or flag.
Assessment: WARNING
Severity: LOW
Notes: `UnsafeStore` is unused in both `store_test.go` and `cmd/demo/main.go`. It does not affect `SafeStore` or demo execution, but represents dead/incomplete scaffolding.
