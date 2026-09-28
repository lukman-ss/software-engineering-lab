# Code Audit

## Finding 1

Location: internal/engine/engine.go:46-99
Claimed Behavior: Atomically enforce NOT NULL, CHECK, UNIQUE, and PARTIAL UNIQUE constraints during record insertion.
Observed Implementation: `Engine.InsertUser` acquires `e.mu.Lock()`, checks empty strings for NOT NULL, validates boundaries for CHECK, checks index maps for UNIQUE / PARTIAL UNIQUE, and updates storage and index maps before releasing mutex.
Assessment: PASS
Severity: LOW
Notes: Implementation correctly simulates database table lock / constraint verification in a single atomic critical section.

## Finding 2

Location: internal/engine/engine.go:123-148
Claimed Behavior: Enforces referential integrity (FOREIGN KEY) and boundary checks on order creation.
Observed Implementation: Checks `o.UserID != 0`, `o.TotalCents > 0`, and confirms existence of `o.UserID` in `e.users` before storing into `e.orders`.
Assessment: PASS
Severity: LOW
Notes: Correctly produces SQLSTATE 23503 on missing foreign key relation.

## Finding 3

Location: internal/dberr/errors.go:61-82
Claimed Behavior: Accurate mapping of SQLSTATE error codes to typed domain errors.
Observed Implementation: `MapToDomainError` inspects `DBError.Code` for SQLSTATE values `23505`, `23503`, `23502`, and `23514`, returning appropriate domain errors (`ErrDuplicate`, `ErrReference`, `ErrInvalidInput`, `ErrValidation`).
Assessment: PASS
Severity: LOW
Notes: Complete coverage of required PostgreSQL integrity constraint violation class (Class 23).

## Finding 4

Location: internal/store/store.go:24-47
Claimed Behavior: UnsafeStore demonstrates vulnerability to read-then-write race condition under concurrent operations.
Observed Implementation: Reads user collection without holding table write lock, sleeps for 1ms to simulate I/O or CPU scheduling delay, then inserts into engine using unconstrained insert.
Assessment: PASS
Severity: LOW
Notes: Successfully reproduces race condition where multiple concurrent registrations with duplicate emails succeed.
