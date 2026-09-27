# Code Audit

## Finding 1

Location:
`internal/engine/engine.go:27-28` (`emailIndex` and `activeEmails` maps)

Claimed Behavior:
Two independent indexes: `emailIndex` for full UNIQUE (email), `activeEmails` for PARTIAL UNIQUE INDEX (email WHERE deleted_at IS NULL).

Observed Implementation:
- `InsertUser` populates `emailIndex` only when `usePartialUniqueIndex = false`.
- `InsertUser` populates `activeEmails` only when `usePartialUniqueIndex = true` AND `u.DeletedAt == nil`.
- `SoftDeleteUser` deletes from `activeEmails` only; does not modify `emailIndex`.

Assessment: WARNING
Severity: MEDIUM
Notes:
Engine supports both indexing modes concurrently on the same instance, leading to inconsistency:
- If a user is inserted via full UNIQUE, `emailIndex` holds mapping but `activeEmails` empty.
- Soft-deleting that user does nothing to `emailIndex`; a subsequent full-unique insert with same email will still see duplicate in `emailIndex` and reject — despite soft delete logically removing the row.
- Conversely, a partial-unique insert after soft delete will see empty `activeEmails` and allow re-use (as expected), but the full `emailIndex` still blocks full-unique re-use.

In practice: tests use fresh engine per test, and each test picks one mode (`false` or `true`), so inconsistency does not surface in test suite. However, documentation implies only one mode is used at a time; mixing modes may produce surprising results.

## Finding 2

Location:
`internal/engine/engine.go:101-120` (`SoftDeleteUser`)

Claimed Behavior:
Soft-delete marks a user as deleted and removes it from the partial unique index, enabling re-registration with same email.

Observed Implementation:
Function signature: `SoftDeleteUser(id int64, deletedAt model.User) error`. The parameter `deletedAt model.User` is odd — it expects a whole `User` struct, but only uses its `DeletedAt` field (`deletedAt.DeletedAt`). Internally:
- Retrieves `u := e.users[id]`.
- Checks `if deletedAt.DeletedAt == nil { return fmt.Errorf("deleted_at cannot be nil for soft delete") }`.
- Deletes from `e.activeEmails` (`delete(e.activeEmails, u.Email)`).
- Sets `u.DeletedAt = deletedAt.DeletedAt`.
- Persists `e.users[id] = u`.

Assessment: WARNING
Severity: LOW
Notes:
API ergonomics: passing a whole `User` just to convey a `*time.Time` is unusual. Could accept `deletedAt *time.Time` directly. However, functionally it works and is exercised in tests and demo. No correctness issue.

## Finding 3

Location:
`internal/engine/engine.go:45-98` (`InsertUser`)

Claimed Behavior:
Atomic evaluation of NOT NULL, CHECK, UNIQUE/PARTIAL UNIQUE under mutex lock; if all pass, commit row + indexes.

Observed Implementation:
Order:
1. NOT NULL checks on `Email`, `Username`.
2. CHECK constraints: `Age >= 18`, `Status in {active,suspended,pending}`.
3. UNIQUE/PARTIAL UNIQUE branching:
   - If `usePartialUniqueIndex`:
       - If `u.DeletedAt == nil`: check `activeEmails`; insert there if free.
       - Else (deleted): skip check, do not insert into `activeEmails`.
   - Else (full UNIQUE):
       - Check `emailIndex`; insert there if free.
4. Primary Key: assign zero ID via `e.userSeq.Add(1)`.
5. Commit:
   - `e.users[u.ID] = u`
   - Conditional index updates as above.

Assessment: PASS
Severity: NONE
Notes:
- All checks happen under `e.mu.Lock()` held for entire function.
- No early returns before commit that could leave partial state.
- Index updates only after successful validation.
- Concurrency safety: RWMutex ensures mutual exclusion with readers via `RLock` in `GetUser`/`GetUsersCount`.

## Finding 4

Location:
`internal/store/store.go:23-47` (`UnsafeStore.RegisterUser`)

Claimed Behavior:
Read-then-write application check for duplicate email (no DB constraint) then insert via engine bypassing unique (`InsertUserUnsafe`). Vulnerable to race.

Observed Implementation:
- Compute `count := s.eng.GetUsersCount()` (RLock).
- Linear scan `i := 1; i <= count` fetching each user via `GetUser` (RLock).
- If email match found → `duplicateFound = true`.
- `time.Sleep(1 * time.Millisecond)` (to widen race window).
- If duplicate → return app error.
- Else → call `s.eng.InsertUserUnsafe(u)` (Lock) which assigns new ID and inserts into `e.users` without touching indexes.

Assessment: PASS
Severity: NONE
Notes:
- Accurately models the vulnerable pattern: check-then-act without transactional guarantee.
- The sleep increases likelihood of interleaving but not required for race.
- Demonstrates duplicate insert under concurrency (see test).

## Finding 5

Location:
`internal/store/store.go:58-66` (`SafeStore.RegisterUser`)

Claimed Behavior:
Delegate uniqueness and validation to storage engine constraints.

Observed Implementation:
- Direct call `s.eng.InsertUser(u, false)` (full UNIQUE) under engine lock.
- On error, map via `dberr.MapToDomainError`.

Assessment: PASS
Severity: NONE
Notes:
- Correctly shifts burden to engine; engine mutex serializes all checks and commit.

## Finding 6

Location:
`internal/dberr/errors.go:75-81` (`IsConstraintViolation`)

Claimed Behavior:
Returns true if error is `*ConstraintError` with matching SQLSTATE.

Observed Implementation:
- Uses `errors.As(err, &cErr)` to unwrap.
- Compares `cErr.Code == code`.

Assessment: PASS
Severity: NONE
Notes:
- Standard pattern; tested in `TestErrorClassification`.

## Finding 7

Location:
`internal/dberr/errors.go:83-99` (`MapToDomainError`)

Claimed Behavior:
Maps SQLSTATE 23502/03/05/14 to domain messages; others to generic.

Observed Implementation:
Switch on `cErr.Code`:
- 23505 → "conflict: resource with this unique attribute already exists (rule: %s)"
- 23502 → "invalid input: mandatory field is missing (rule: %s)"
- 23514 → "validation failed: value outside permissible boundary (rule: %s)"
- 23503 → "reference error: referenced entity does not exist (rule: %s)"
- default → "database integrity violation: %s"

Assessment: PASS
Severity: NONE
Notes:
- Matches demo output; user-facing messages appropriate.

## Finding 8

Location:
`internal/engine/engine.go:14-21` (field declarations)

Claimed Behavior:
`mu` protects `users`, `orders`, `emailIndex`, `activeEmails`; sequences `userSeq`, `orderSeq` atomic.

Observed Implementation:
- `mu sync.RWMutex`
- `userSeq`, `orderSeq` `atomic.Int64`
- Maps: `users`, `orders`, `emailIndex`, `activeEmails`

Assessment: PASS
Severity: NONE
Notes:
- All mutated fields under `mu.Lock()` in writers.
- Readers use `mu.RLock()` (`GetUser`, `GetUsersCount`).
- Sequences only incremented under lock (though atomic would permit unlocked increment; current usage under lock fine).

## Finding 9

Location:
`internal/engine/engine.go:50-56` (NOT NULL email check)

Claimed Behavior:
Reject empty email with `NotNullViolation`.

Observed Implementation:
- `if u.Email == "" { return model.User{}, dberr.NewNotNullViolation("users", "email", "users_email_not_null") }`

Assessment: PASS
Severity: NONE
Notes:
- Correctly maps to SQLSTATE 23502.
- Username checked similarly.

## Finding 10

Location:
`internal/engine/engine.go:59-68` (CHECK age >= 18)

Claimed Behavior:
Reject age < 18 with `CheckViolation`.

Observed Implementation:
- `if u.Age < 18 { return model.User{}, dberr.NewCheckViolation(...)}`

Assessment: PASS
Severity: NONE
Notes:
- Boundary case: age = 18 passes (correct). Test only exercises 16 (<18) and presumably 25 (valid). Edge case 18 not explicitly tested but implied pass.

## Finding 11

Location:
`internal/engine/engine.go:63-68` (CHECK status enum)

Claimed Behavior:
Reject status not in {active,suspended,pending}.

Observed Implementation:
Switch on `u.Status`; only three arms accepted.

Assessment: PASS
Severity: NONE
Notes:
- Case-sensitive; empty string fails (caught by switch default). Valid strings exact match.

## Finding 12

Location:
`internal/engine/engine.go:127-135` (`InsertOrder` NOT NULL user_id)

Claimed Behavior:
Reject zero `UserID`.

Observed Implementation:
- `if o.UserID == 0 { return model.Order{}, dberr.NewNotNullViolation("orders", "user_id", "orders_user_id_not_null") }`

Assessment: PASS
Severity: NONE
Notes:
- Matches expectation.

## Finding 13

Location:
`internal/engine/engine.go:132-135` (`InsertOrder` CHECK total_cents > 0)

Claimed Behavior:
Reject `<= 0`.

Observed Implementation:
- `if o.TotalCents <= 0 { return model.Order{}, dberr.NewCheckViolation(...) }`

Assessment: PASS
Severity: NONE
Notes:
- Boundary: `TotalCents = 0` rejected (correct). Negative also rejected.

## Finding 14

Location:
`internal/engine/engine.go:137-140` (`InsertOrder` FOREIGN KEY)

Claimed Behavior:
Reject if `user_id` not in `e.users`.

Observed Implementation:
- `if _, exists := e.users[o.UserID]; !exists { return ..., dberr.NewForeignKeyViolation(...) }`

Assessment: PASS
Severity: NONE
Notes:
- References `users` map directly; assumes user exists if key present (true because engine never deletes users except via soft-delete which keeps map entry with DeletedAt set).

## Finding 15

Location:
`internal/store/store_test.go:162-208` (`TestConcurrentRegistration_Safe_EnforcesUniqueness`)

Claimed Behavior:
20 goroutines try to insert same email via `SafeStore`; exactly 1 succeeds, 19 get unique violation, final user count = 1.

Observed Implementation:
- Uses `sync.WaitGroup`, mutex-protected counters.
- Calls `s.RegisterUser` (safe path).
- Asserts successCount == 1, errorCount == goroutines-1, eng.GetUsersCount() == 1.

Assessment: PASS
Severity: NONE
Notes:
- Race detector clean; outcome deterministic due to engine mutex.

## Finding 16

Location:
`internal/store/store_test.go:210-239` (`TestConcurrentRegistration_Unsafe_SuffersRaceCondition`)

Claimed Behavior:
20 goroutines via `UnsafeStore` yield >1 user count (race condition duplicates email).

Observed Implementation:
- No error checking; discards results.
- Asserts `eng.GetUsersCount() > 1`.

Assessment: PASS
Severity: NONE
Notes:
- Demonstrates race condition; no data corruption (each gets unique ID via atomic sequence) but logical duplicate email.