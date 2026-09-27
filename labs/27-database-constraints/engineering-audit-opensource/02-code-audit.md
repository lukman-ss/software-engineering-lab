# Code Audit

Engine: in-memory, coarse-grained `sync.RWMutex`-locked map store with `atomic.Int64` PK sequences (engine.go). Models PG SQLSTATE codes; no external DB.

## Finding 1
Location: engine.go:46-99 InsertUser
Claimed Behavior: NOT NULL checks email+username; CHECK age>=18/status IN active/suspended/pending; UNIQUE on email, UNIQUE partial active email; atomic.
Observed: All checks evaluated under single `mu.Lock` held for the whole insert; atomic. SQLSTATE 23502/23514/23505 returned with constraint names matching README.
Assessment: PASS
Severity: LOW
Notes: Empty string treated as "not-null violation" (no separate empty-string rejection) — acceptable for in-memory simulator (README claims empty values rejected for required columns).

## Finding 2
Location: engine.go:123-148 InsertOrder
Claimed Behavior: NOT NULL on user_id (23502), CHECK total_cents>0 (23514), FK to users.id (23503).
Observed: Implemented exactly. FK uses `e.users[o.UserID]` map lookup under same lock held for check+insert (consistent snapshot, no orphan possible).
Assessment: PASS
Severity: LOW
Notes: FK does not enforce referenced-table PK/UNIQUE — documented engine simplification; fine for lab.

## Finding 3
Location: engine.go:102-120 SoftDeleteUser + InsertUser usePartialUniqueIndex branch
Claimed Behavior: Partial UNIQUE only active (DeletedAt==nil) rows; soft-delete removes from active index; re-registration after soft-delete succeeds.
Observed: Active email registered; duplicate active rejected; SoftDelete removes entry from activeEmails (and only activeEmails — emailIndex never modified). Second `RegisterUserPartial` (DeletedAt==nil) passes partial check because email no longer in activeEmails; new active created.
Assessment: PASS
Severity: LOW
Notes: `emailIndex` is never touched by the partial path, so a soft-deleted partial email could be re-inserted as soft-deleted again — allowed by spec.

## Finding 4
Location: engine.go:33-43 InsertUserUnsafe
Claimed Behavior: Bypasses all constraints (simulates unconstrained table).
Observed: No field/index checks beyond PK sequence; inserts regardless of email dup/neg total.
Assessment: PASS
Severity: LOW
Notes: Intentionally unsafe; correct per design.

## Finding 5
Location: store.go:24-47 UnsafeStore.RegisterUser
Claimed Behavior: Read-then-write over `GetUser` count; vulnerable under concurrency; sleeps to widen race window.
Observed: Loop `count := GetUsersCount()` then iterates `GetUser(i)` comparing Email; `time.Sleep(1ms)` before insert — race window amplified.
Assessment: PASS
Severity: MEDIUM
Notes: O(n) scan + count race window is the documented vulnerability; demo/test depend on it to materialize duplicates consistently. Flaky-free in practice because sleep guarantees interleaving of the read phase; still inherently probabilistic.

## Finding 6
Location: dberr/errors.go:83-99 MapToDomainError
Claimed Behavior: Maps 23505/23502/23514/23503 to domain-friendly errors.
Observed: Correct switch over Code (typed SQLState); wraps `ConstraintError` only, returns original otherwise.
Assessment: PASS
Severity: MEDIUM
Notes: MAPPING_LOSS — SQLSTATE code not retained in returned `fmt.Errorf` (no `%w`/typed error). Tests only assert `err == nil` vs `err != nil`, so loss undetected. Should return a typed domain error preserving Code for programmatic handling.

## Finding 7
Location: engine.go:73-83 partial-unique vs full-unique branches
Claimed Behavior: Full UNIQUE and PARTIAL UNIQUE both enforced.
Observed: `InsertUser(u, usePartialUniqueIndex bool)` uses one branch or the other — a single row is never checked under BOTH. `RegisterUser` calls InsertUser(false); `RegisterUserPartial` calls InsertUser(true). Tables are effectively independent (full-emailIndex vs activeEmails map).
Assessment: PASS
Severity: LOW
Notes: `RegisterUser` and `RegisterUserPartial` are separate "tables" with separate indexes. Fine for demo; README lists both; no bug.

## Finding 8
Location: demo main.go:97-98
Claimed Behavior: SQLSTATE taxonomy verification: code=23505 isUniqueViolation=true.
Observed: Constructs `NewUniqueViolation("users","users_email_key","duplicate key")` (Code=23505) and verifies `IsConstraintViolation(testErr, SQLStateUniqueViolation)` true.
Assessment: PASS
Severity: LOW
Notes: Pure self-check, output matches execution-result.md.

## Finding 9
Location: store.go:50-65 SafeStore.RegisterUser / RegisterUserPartial / CreateOrder
Claimed Behavior: Delegate InsertUser/InsertOrder to engine constraints; map via dberr.
Observed: Pass-through; no additional app-level duplicate check (correct — relies on engine). Context param unused (demo never cancelled).
Assessment: PASS
Severity: LOW
Notes: `ctx` accepted but ignored (no select on ctx.Done); acceptable in lab but real impl should respect cancellation.

## Finding 10
Location: engine.go:16-21,23-29 struct/init
Claimed Behavior: Thread-safe via RWMutex + atomic counters.
Observed: `userSeq`/`orderSeq` atomic; maps protected by `mu`. All mutation sites acquire `e.mu.Lock()`/`RLock()`. No unprotected access.
Assessment: PASS
Severity: LOW
Notes: None.
