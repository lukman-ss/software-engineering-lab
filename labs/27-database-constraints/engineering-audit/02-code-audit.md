# Code Audit

Target Lab: labs/27-database-constraints

## Finding 1

Location: internal/engine/engine.go:46-99 (`InsertUser`)
Claimed Behavior: Atomic evaluation of NOT NULL, CHECK, UNIQUE, and PARTIAL UNIQUE constraints with SQLSTATE violations.
Observed Implementation: Handled under `e.mu.Lock()` guarding in-memory state. NOT NULL checks empty strings, CHECK enforces `Age >= 18` and valid status enum, UNIQUE and PARTIAL UNIQUE check distinct index maps before insertion.
Assessment: PASS
Severity: LOW
Notes: Correctly models SQL constraint checks before record insertion.

## Finding 2

Location: internal/store/store.go:24-47 (`UnsafeStore.RegisterUser`) vs store.go:59-66 (`SafeStore.RegisterUser`)
Claimed Behavior: SafeStore delegates constraint integrity to engine layer preventing race conditions, whereas UnsafeStore exhibits read-then-write check race conditions.
Observed Implementation: UnsafeStore checks `s.eng.GetUsersCount()` and iterates with an artificial yield (`time.Sleep(1ms)`), allowing concurrent routines to slip through. SafeStore inserts directly through `InsertUser` protected by engine mutex and unique index map, mapping SQLSTATE to domain error.
Assessment: PASS
Severity: LOW
Notes: Demonstrates the structural flaw of application-level validation vs declarative database constraints.

## Finding 3

Location: internal/engine/engine.go:102-120 (`SoftDeleteUser`)
Claimed Behavior: Updates user record `DeletedAt` timestamp and evicts entry from `activeEmails` partial index map.
Observed Implementation: Acquires write lock `e.mu.Lock()`, validates record existence, verifies non-nil `deletedAt`, removes entry from `e.activeEmails`, and updates record in `e.users`. Subsequent inserts with the same email succeed as long as partial indexing is selected.
Assessment: PASS
Severity: LOW
Notes: Properly models PostgreSQL `CREATE UNIQUE INDEX ... WHERE deleted_at IS NULL` semantics.

## Finding 4

Location: internal/engine/engine.go:123-148 (`InsertOrder`)
Claimed Behavior: Enforces foreign key constraint referencing users table and check constraint on total_cents.
Observed Implementation: Rejects `o.UserID == 0` with NOT NULL error (23502), rejects `o.TotalCents <= 0` with CHECK violation (23514), and verifies user exists in `e.users` returning SQLSTATE 23503 on missing key.
Assessment: PASS
Severity: LOW
Notes: Referential integrity correctly guaranteed at engine layer.

## Finding 5

Location: internal/dberr/errors.go:83-100 (`MapToDomainError`)
Claimed Behavior: Maps SQLSTATE constraint errors to idiomatic domain errors.
Observed Implementation: Inspects `ConstraintError.Code` and translates `23505` to conflict error, `23502` to invalid input error, `23514` to validation error, and `23503` to reference error.
Assessment: PASS
Severity: LOW
Notes: Clean domain translation adhering to standard SQLSTATE taxonomy.
