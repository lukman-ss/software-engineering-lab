# Code Audit — labs/27-database-constraints

## Finding 1

Location: `internal/engine/engine.go:46-99` (`InsertUser`)
Claimed Behavior: NOT NULL (23502), CHECK (23514), UNIQUE / PARTIAL UNIQUE (23505) enforced atomically.
Observed Implementation: All checks run under a single `e.mu.Lock()`; indexes (`emailIndex`, `activeEmails`) updated only after checks pass; PK generated via `atomic.Int64` inside the lock.
Assessment: PASS
Severity: LOW
Notes: Check-then-insert is atomic. Coarse table lock is correct for a simulator.

## Finding 2

Location: `internal/engine/engine.go:123-148` (`InsertOrder`)
Claimed Behavior: NOT NULL on `user_id`, CHECK `total_cents > 0`, FOREIGN KEY to `users(id)`.
Observed Implementation: Matches claim exactly; FK checked under same lock as user writes, so no TOCTOU between user insert and order insert.
Assessment: PASS
Severity: LOW
Notes: None.

## Finding 3

Location: `internal/engine/engine.go:102-120` (`SoftDeleteUser`)
Claimed Behavior: Soft delete frees the email in the partial-unique index.
Observed Implementation: Deletes unconditionally from `activeEmails`, but does NOT touch `emailIndex` (full-unique path). A user created via `RegisterUser` (non-partial) then soft-deleted still blocks email reuse via `emailIndex`.
Assessment: WARNING
Severity: LOW
Notes: Demo/tests only soft-delete partial-path users, so behavior is consistent within the demonstrated scope. Cross-path reuse is untested and would surprise.

## Finding 4

Location: `internal/store/store.go:24-47` (`UnsafeStore.RegisterUser`)
Claimed Behavior: Intentionally vulnerable read-then-write check to demonstrate the race.
Observed Implementation: Scans `GetUser(1..count)` with a 1ms sleep widening the race window, then `InsertUserUnsafe` bypasses constraints. Under `-race` with 20 goroutines it reliably produces duplicates (verified).
Assessment: PASS
Severity: LOW
Notes: ID-range scan assumes dense sequential IDs; valid for this simulator but fragile as a general pattern. Acceptable for a deliberate anti-pattern demo.

## Finding 5

Location: `internal/dberr/errors.go:83-99` (`MapToDomainError`)
Claimed Behavior: Maps SQLSTATE codes to domain-friendly errors.
Observed Implementation: Returns plain `fmt.Errorf` strings; SQLSTATE code is dropped, so `IsConstraintViolation(mappedErr, …)` returns false. Callers cannot programmatically classify mapped errors.
Assessment: WARNING
Severity: MEDIUM
Notes: Tests only classify raw constructors, never mapped errors. Either wrap (`%w`) the `*ConstraintError` or document that mapped errors are display-only.

## Finding 6

Location: `internal/engine/engine.go:33-43` (`InsertUserUnsafe`)
Claimed Behavior: Bypasses unique constraints (unconstrained table simulation).
Observed Implementation: Takes the write lock and inserts with no validation at all (also skips NOT NULL/CHECK).
Assessment: PASS
Severity: LOW
Notes: Matches the documented "no constraints" simulation. No data race: full mutex held.

## Finding 7

Location: `internal/engine/engine.go:150-161` (`GetUser`, `GetUsersCount`), `store.go` concurrency paths
Claimed Behavior: Safe under concurrent use.
Observed Implementation: All map access under `RWMutex`; counters atomic; test `mu` guards success/error tallies. `go test -race` clean.
Assessment: PASS
Severity: LOW
Notes: None.

## Finding 8

Location: `internal/model/model.go`
Claimed Behavior: `User` / `Order` domain entities.
Observed Implementation: Plain structs, no validation logic leaking into the model; constraints live in the engine where claimed.
Assessment: PASS
Severity: LOW
Notes: None.
