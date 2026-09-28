# Code Audit

Reviewer: Engineering Auditor (independent)
Lab: labs/27-database-constraints
Module: `github.com/lukman/software-engineering-lab/labs/27-database-constraints`
Go version: 1.22 (per `go.mod`)

Scope: `internal/dberr`, `internal/model`, `internal/engine`, `internal/store`, `cmd/demo`. Static + dynamic (build, `go vet`, `go test -race`, demo).

## Summary of evidence produced during audit

```
go build ./...          -> PASS (exit 0)
go vet ./...            -> PASS (exit 0)
go test -v ./...        -> PASS (8/8 tests)
go test -race ./...     -> PASS (clean, no data races)
go test -race -count=5 -run Concurrent ./internal/store/ -> PASS (6x reproducible)
go run ./cmd/demo       -> PASS (see 06-verdict for actual output)
```

---

## Finding 1 — NOT NULL enforcement is correct and atomic

Location: `internal/engine/engine.go:46-56` (`InsertUser`), `internal/engine/engine.go:128-130` (`InsertOrder`).

Claimed Behavior: NOT NULL constraint rejects missing `email`/`username`/`user_id` and returns SQLSTATE `23502`.

Observed Implementation: `InsertUser` checks `u.Email == ""` and `u.Username == ""` and returns `dberr.NewNotNullViolation(...)`. `InsertOrder` checks `o.UserID == 0` and returns `NewNotNullViolation("orders", "user_id", "orders_user_id_not_null")`. Both checks run under `e.mu.Lock()` before the row is committed. `MapToDomainError` maps `23502` → "invalid input: mandatory field is missing".

Assessment: PASS. Severity: LOW. Notes: Checks are inside the write-lock critical section, so they are consistent with the insert that follows. The only nuance: `OrderID` has no NOT NULL on `total_cents`... it does (CHECK, not NOT NULL). Fine.

---

## Finding 2 — CHECK constraint enforces age boundary and status enum

Location: `internal/engine/engine.go:58-68`.

Claimed Behavior: `CHECK (age >= 18)` and `CHECK (status IN ('active','suspended','pending'))` reject invalid rows with SQLSTATE `23514`.

Observed Implementation: `u.Age < 18` → `NewCheckViolation("users", "users_age_check", ...)`. `switch` on `u.Status` with `default` returning `NewCheckViolation("users", "users_status_check", ...)`. `MapToDomainError` maps `23514` → "validation failed: value outside permissible boundary". Order-side CHECK `total_cents > 0` at `engine.go:133`.

Assessment: PASS. Severity: LOW. Notes: Boundary `age == 18` is accepted (strict `<` check is correct for `>= 18`). Status enum uses a closed switch with explicit default — correct.

---

## Finding 3 — UNIQUE constraint prevents duplicate email (full index path)

Location: `internal/engine/engine.go:70-83,90-96`.

Claimed Behavior: `RegisterUser` (SafeStore) rejects duplicate email with SQLSTATE `23505`; demo concurrency stress test yields exactly 1 success + 49 rejections.

Observed Implementation: `InsertUser` with `usePartialUniqueIndex == false` checks `e.emailIndex[u.Email]` under the write lock; on collision returns `NewUniqueViolation("users", "users_email_key", ...)`. On success, writes `e.emailIndex[u.Email] = u.ID` inside the same locked section. `SafeStore.RegisterUser` → `eng.InsertUser(u, false)`. Because the check + index insertion are both inside one `Lock()`/`defer Unlock()`, two concurrent inserters cannot both observe "not present" — exactly one wins.

Assessment: PASS. Severity: LOW. Notes: This is the core claim of the lab ("database UNIQUE constraint prevents race condition") and it is genuinely implemented with a single mutex critical section. Confirmed live: `go run ./cmd/demo` reported 1 success / 49 rejected / integrity intact; `TestConcurrentRegistration_Safe_EnforcesUniqueness` passed x5 under `-race`.

---

## Finding 4 — FOREIGN KEY constraint rejects orphan orders

Location: `internal/engine/engine.go:137-140`.

Claimed Behavior: Order referencing a non-existent `user_id` is rejected with SQLSTATE `23503`; valid reference succeeds.

Observed Implementation: `e.users[o.UserID]` lookup under the write lock in `InsertOrder`; if absent returns `NewForeignKeyViolation("orders", "fk_orders_user", ...)`. `MapToDomainError` maps `23503` → "reference error: referenced entity does not exist".

Assessment: PASS. Severity: LOW. Notes: Lookup and insert are atomic under the same lock. `TestForeignKeyConstraint` covers both negative and positive cases.

---

## Finding 5 — PARTIAL UNIQUE index allows soft-delete re-registration

Location: `internal/engine/engine.go:71-83,90-96`, `internal/engine/engine.go:101-120` (`SoftDeleteUser`).

Claimed Behavior: `WHERE deleted_at IS NULL` index permits multiple soft-deleted rows with the same email but only one active row; soft-delete frees the active slot.

Observed Implementation:
- `InsertUser` with `usePartialUniqueIndex == true`: duplicates checked only against `e.activeEmails` when `u.DeletedAt == nil`. Committed to `e.activeEmails` only when `u.DeletedAt == nil`.
- `SoftDeleteUser`: removes the email from `e.activeEmails` via `delete`, then writes `DeletedAt` onto the stored row.

Assessment: PASS for the single-path scenario the lab claims. Severity: LOW. Notes: See Finding 7 for the cross-path limitation.

---

## Finding 6 — Error type taxonomy and classification helpers correct

Location: `internal/dberr/errors.go` (lines 8-18 codes, 75-81 `IsConstraintViolation`, 83-100 `MapToDomainError`).

Claimed Behavior: SQLSTATE constants match ANSI/PostgreSQL Class 23 codes (`23001`, `23502`, `23503`, `23505`, `23514`, `23P01`, `40001`); `IsConstraintViolation` type-asserts and compares `Code`; `MapToDomainError` rewrites each code to a domain message.

Observed Implementation: Constants verified against PostgreSQL docs. `IsConstraintViolation` uses `errors.As` to unwrap `*ConstraintError`. `MapToDomainError` handles `23505`/`23502`/`23514`/`23503` with specific messages and a `default` returning the raw detail. The unimplemented codes (`23001`/exclusion, `23P01`/exclusion, `40001`/serialization) fall to the default branch and pass through unchanged — acceptable since they are unused.

Assessment: PASS. Severity: LOW. Notes: `IsConstraintViolation` returns `false` for a plain `error` (no `*ConstraintError`), which is the intended contract.

---

## Finding 7 — UnsafeStore genuinely exhibits the read-then-write race (intended)

Location: `internal/store/store.go:24-47` (`UnsafeStore.RegisterUser`), `internal/engine/engine.go:33-43` (`InsertUserUnsafe`).

Claimed Behavior: App-only uniqueness check races and allows >1 duplicate row; `time.Sleep` widens the window.

Observed Implementation: `UnsafeStore.RegisterUser` calls `GetUsersCount()`/`GetUser()` (read lock, no write lock held across the loop), sleeps 1ms, then calls `InsertUserUnsafe` which takes the write lock only for the insert. The check-then-insert is not atomic — the classic TOCTOU. `InsertUserUnsafe` deliberately performs no uniqueness check, simulating an unconstrained table. Confirmed live: `TestConcurrentRegistration_Unsafe_SuffersRaceCondition` reports `count > 1` (20 rows in observed run).

Assessment: PASS (the bug is intentional and correctly demonstrated). Severity: LOW. Notes: This test is timing-dependent (see Finding 9). It reliably fails the `count <= 1` assertion in practice, but there is an intrinsic flakiness risk: if the scheduler happened to serialize the 20 goroutines such that each inserts before the next checks, `count` could still be >1 anyway. The real risk direction is the opposite test (`Safe`) which is deterministic. The `time.Sleep` is an artificial crutch.

---

## Finding 8 — Data-race safety (concurrency)

Location: `internal/engine/engine.go` mutex usage across all mutating reads.

Claimed Behavior: In-memory engine is thread-safe; no data races.

Observed Implementation: Every access to `e.users`, `e.orders`, `e.emailIndex`, `e.activeEmails` is guarded by `e.mu` (`RLock`/`Lock`). `userSeq`/`orderSeq` use `atomic.Int64`. `GetUsersCount`/`GetUser` use `RLock`. `InsertUserUnsafe`/`InsertUser`/`InsertOrder`/`SoftDeleteUser` use `Lock`. The `UnsafeStore` app-check path releases the read lock before inserting, but that is a *design* race (TOCTOU), not a *data race* — no unsynchronized memory access occurs. `go test -race` is clean across repeated runs.

Assessment: PASS (no data races). Severity: LOW. Notes: `go vet ./...` also clean.

---

## Finding 9 — `time.Sleep(1ms)` in UnsafeStore is a timing-dependent flakiness source

Location: `internal/store/store.go:38`.

Observed Implementation: Hard-coded `time.Sleep(1 * time.Millisecond)` inside the unsafe registration path. This widens the race window so the unsafe test reliably loses.

Assessment: WARNING. Severity: MEDIUM. Notes: Not a correctness bug — the race is intended. But the test's assertion (`count <= 1` must fail) is only statistically guaranteed; on an extremely fast or adversarial-scheduled CI host the artificial delay could still leave all 20 goroutines inserting, which is the desired outcome here, so this test is robust *by design*. The real concern is pedagogical: the "fix" being showcased is the SafeStore (locking + unique constraint), not the sleep. Acceptable but worth flagging to the writer.

---

## Finding 10 — `SoftDeleteUser` API signature is awkward and the second param is misnamed

Location: `internal/engine/engine.go:101-120`.

Observed Implementation: `SoftDeleteUser(id int64, deletedAt model.User) error` takes a whole `model.User` only to read `.DeletedAt`. The parameter name `deletedAt` is misleading (it is a User). If `deletedAt.DeletedAt == nil` it returns an error.

Assessment: WARNING. Severity: LOW. Notes: Works correctly but obscures intent. A cleaner signature would be `SoftDeleteUser(id int64, at time.Time) error`. Not a functional defect.

---

## Finding 11 — Cross-index (full-UNIQUE vs partial-UNIQUE) path is inconsistent

Location: `internal/engine/engine.go` — `emailIndex` (full unique, `InsertUser(...,false)`) vs `activeEmails` (partial, `InsertUser(...,true)`).

Observed Implementation: A user inserted via the partial path is only stored in `activeEmails`, never in `emailIndex`. Conversely a user inserted via the full path is only in `emailIndex`, never in `activeEmails`. The two membership indexes are disjoint. Consequently, registering the same email through `RegisterUserPartial` (active) then `RegisterUser` (full unique) would *not* detect the conflict, because the second path consults only `emailIndex`.

Assessment: WARNING. Severity: MEDIUM. Notes: Within each isolated demo path uniqueness holds; the lab never mixes the two registration methods, so no current test or the demo violates this. But the engine's two public entry points are not mutually aware, so a caller combining them would silently produce duplicate "active" rows under the full-unique lens. This is a latent integrity hole. Recommend either (a) a single unified uniqueness check, or (b) explicit documentation that the two paths are mutually exclusive namespaces. The design doc (Finding in 04) claims UNIQUE enforces single occurrence globally — the implementation only enforces it per-path.

---

## Finding 12 — No rollback / transaction semantics (acceptable for insert, but note)

Location: `internal/engine/engine.go` — no transaction/rollback.

Observed Implementation: Each insert is a single locked critical section; on constraint failure no partial state is committed (the row is written only after all checks pass). There is no multi-statement transaction, but none is claimed.

Assessment: PASS within scope. Severity: LOW. Notes: Fine because every "operation" is atomic. If a future feature did a read-modify-write spanning multiple engine calls, rollback would be needed.

---

## Finding 13 — No recovery / corruption handling (out of scope)

Location: engine does not simulate WAL, checkpoints, or crash recovery.

Observed Implementation: Pure in-memory maps; nothing to recover after restart.

Assessment: NOT_APPLICABLE for this lab's claims. Severity: N/A. Notes: The design doc's "ACID-style" phrasing is qualified by "simulation"; not a defect given the stated in-memory scope.
