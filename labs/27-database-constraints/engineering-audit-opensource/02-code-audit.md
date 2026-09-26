# Engineering Code Audit

## Finding 1

Location: internal/store/store.go:40-42
Claimed Behavior: UnsafeStore "bypasses unique constraints (simulating uncontrained table)" so the app-level read-then-write check is the only uniqueness protection.
Observed Implementation: `s.eng.InsertUser(u, true)` still enforces the partial unique index (`activeEmails`) for any user with `DeletedAt == nil`. There is no engine call path that disables uniqueness; the `else` branch always checks `emailIndex`.
Assessment: WARNING
Severity: MEDIUM
Notes: The "unsafe" store is therefore not actually vulnerable to the read-then-write race the research describes — the engine still serializes and rejects duplicates. The safe-vs-unsafe contrast is not faithfully implemented. No test exercises UnsafeStore; its behavior under concurrency is unverified.

## Finding 2

Location: internal/store/store.go:59, 68, 77; internal/dberr/errors.go:83-100
Claimed Behavior: Research best practice: "Always check SQLSTATE, not error text" and map 23xxx codes to domain errors.
Observed Implementation: `SafeStore` wraps engine errors via `dberr.MapToDomainError`, converting `ConstraintError` into a generic `fmt.Errorf`. The original `SQLSTATE`/`ConstraintError` type is discarded. No store method returns a `ConstraintError`; callers cannot call `IsConstraintViolation` on store return values.
Assessment: WARNING
Severity: MEDIUM
Notes: Programmatic error handling by SQLSTATE (e.g., retry on 40001) is impossible on store return values. Tests do not assert SQLSTATE on store returns.

## Finding 3

Location: internal/engine/engine.go:12-30, 33-86
Claimed Behavior: Engine enforces NOT NULL, CHECK, UNIQUE, and PARTIAL UNIQUE constraints atomically.
Observed Implementation: All constraints evaluated inside `e.mu.Lock()`. Check-then-insert is atomic. ID generated via `atomic.Int64`. Concurrency test passes (20 goroutines → exactly 1 success, 19 errors, count == 1), `go test -race` clean.
Assessment: PASS
Severity: LOW
Notes: Behaviorally correct. The engine uses coarse-grained table locking rather than per-row B-tree page latches — a deliberate, documented trade-off (engineering/02-implementation-notes.md:25). Observable outcome matches research's claim that constraints prevent concurrent duplicate writes.

## Finding 4

Location: internal/dberr/errors.go:11, 16-17
Claimed Behavior: SQLSTATE constants defined for the documented taxonomy.
Observed Implementation: `SQLStateRestrictViolation` (23001) has no constructor and is never returned. `SQLStateExclusionViolation` (23P01) and `SQLStateSerializationFail` (40001) defined but not produced by the engine. Research recommends retry-on-40001 for serialization failures, but no code path generates 40001.
Assessment: WARNING
Severity: LOW
Notes: Dead constants. No fabrication — they are simply unused.

## Finding 5

Location: internal/dberr/errors.go:75-81
Claimed Behavior: `IsConstraintViolation` inspects errors by SQLSTATE via `errors.As`.
Observed Implementation: Correct. Uses `errors.As(err, &cErr)` and compares `cErr.Code`. Works on raw `ConstraintError` returned by the engine (before `MapToDomainError` mapping).
Assessment: PASS

## Finding 6

Location: internal/engine/engine.go:89-107
Claimed Behavior: `SoftDeleteUser` updates `DeletedAt` and removes the user from the active partial index.
Observed Implementation: Removes email from `activeEmails`, sets `DeletedAt` on the stored user. Returns plain `fmt.Errorf` for missing user / nil DeletedAt.
Assessment: PASS
Severity: LOW
Notes: Error returns are plain `fmt.Errorf` (not ConstraintError) — acceptable for a non-constraint error. Missing-user edge case untested.

## Finding 7

Location: internal/store/store.go:22-43
Claimed Behavior: UnsafeStore performs app-level read-then-write uniqueness check.
Observed Implementation: Calls `eng.GetUsersCount()` then `eng.GetUser(i)` for i in 1..count. Correct only when IDs are dense 1..N. Fragile if IDs ever gap (e.g., after deletions with re-used sequence), but IDs are monotonic and never reused in the sim, so it is currently safe-by-luck.
Assessment: WARNING
Severity: LOW
Notes: Not a concurrency flaw but a correctness smell; only the app-check loop is fragile.

## Finding 8

Location: cmd/demo/main.go
Claimed Behavior: Demo exercises every constraint including 50-goroutine concurrency.
Observed Implementation: Builds engine + SafeStore, exercises NOT NULL (missing email), CHECK (age<18, invalid status), FOREIGN KEY (UserID=9999), partial unique (soft delete + reuse), and a 50-goroutine race. Also calls `dberr.NewUniqueViolation` standalone to assert the SQLSTATE taxonomy.
Assessment: PASS
Severity: LOW
Notes: The final "SQLSTATE Taxonomy Verification" line constructs a synthetic error rather than asserting the SQLSTATE on an actual store-level duplicate error — cosmetic rather than end-to-end.

## Finding 9

Location: internal/model/model.go
Observation: `Order` struct has a `Status` field that is never validated, never set, and never returned. The engine does not enforce any constraint on it.
Assessment: WARNING
Severity: LOW
Notes: Unused/incomplete field. Not referenced by tests or demo. No impact on verified behavior.