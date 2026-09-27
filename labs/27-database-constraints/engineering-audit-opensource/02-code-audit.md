# Code Audit

Method: inspection + live execution on macOS/darwin, Go 1.22, `go.mod` in-repo.

Execution results:
- `go vet ./...` -> PASS (exit 0)
- `go build ./...` -> PASS
- `go test -v ./...` (9 tests incl. 2 concurrency) -> PASS, 0 failures
- `go test -race -count=1 ./...` -> PASS, no data race reported
- `go run ./cmd/demo` -> PASS, integrity integible=true

## Finding 1
Location: internal/engine/engine.go:102-120 — SoftDeleteUser
Claimed Behavior: Soft-delete removes `activeEmails` entry so a new active user with same email can be inserted afterward, modeling `CREATE UNIQUE INDEX ... WHERE deleted_at IS NULL`.
Observed Implementation: `delete(e.activeEmails, u.Email)` + `u.DeletedAt = deletedAt.DeletedAt` + write back under `e.mu.Lock()`. Index entry removed; row retained with DeletedAt set.
Assessment: PASS — partial-index semantics implemented correctly and atomically under engine mutex.

## Finding 2
Location: internal/engine/engine.go:46-99 — InsertUser
Claimed Behavior: UNIQUE / PARTIAL UNIQUE uniqueness check happens atomically at the storage engine (mutex) layer, defeating read-then-write races.
Observed Implementation: Both `e.emailIndex` and `e.activeEmails` checks live inside the same critical section as the write; no gap between duplicate-check and insert. Use of coarse `e.mu.Lock()` (not `Lock`/`unlock` split) guarantees atomicity.
Assessment: PASS — no TOCTOU window; SafeStore concurrency test asserts exactly 1 survivor and passes.

## Finding 3
Location: internal/store/store.go:24-47 — UnsafeStore.RegisterUser
Claimed Behavior: Application-only count-then-insert loop is racy; 20 concurrent writes should exceed 1 row.
Observed Implementation: Reads `GetUsersCount()` + `GetUser(1..n)` under a *read* lock, then `time.Sleep(1ms)`, then `InsertUserUnsafe` under a write lock. The check and insert use two disjoint critical sections, leaving a deliberate race window. Test asserts `count > 1`.
Assessment: WARNING (intentional). The race is real but timing-dependent (1ms sleep). It reliably produced >1 in tested runs, so verdict is not fabricated; severity kept LOW because the demo/test purpose is illustrative. Not a defect in SafeStore.

## Finding 4
Location: internal/dberr/errors.go:83-100 — MapToDomainError
Claimed Behavior: Bridges SQLSTATE codes (23502/23503/23505/23514) to user-facing domain messages.
Observed Implementation: Uses `fmt.Errorf(...)` (no `%w`), so the returned error loses the underlying `*ConstraintError` wrapping. Callers/tests that need `errors.As` for the code can still recover via `dberr.IsConstraintViolation` on the *original* engine error, but downstream domain consumers cannot.
Assessment: WARNING — correct message mapping, but error chain is severed. No test exercises the returned `error` type beyond nil-vs-non-nil, so behavior is technically unproven at the machine-readable level. Severity MEDIUM given README emphasizes "error mapping."

## Finding 5
Location: store.go / engine.go — `ctx context.Context` parameter
Claimed Behavior: API uses context (signature includes `ctx`).
Observed Implementation: `ctx` is accepted and never read; no timeout/cancel honored.
Assessment: LOW / NOT_A_DEFEATED_CLAIM — harmless for an in-memory demo but misrepresents cancellation support. No test for ctx cancellation.

## Finding 6
Location: engine.go:11-21 — struct field types
Claimed Behavior: Coarse-grained table-level lock.
Observed Implementation: `sync.RWMutex` + `atomic.Int64` sequence counters; `sync.RWMutex` is used consistently (Write lock for inserts/soft-delete, RLock for reads).
Assessment: PASS for concurrency safety under `-race`. Documented trade-off (impl notes) matches reality.

## Finding 7
Location: engine.go:137-147 — InsertOrder FK check
Claimed Behavior: Foreign key on orders.user_id -> users(id).
Observed Implementation: `if _, exists := e.users[o.UserID]; !exists { return ForeignKeyViolation }`. Checked before PK assignment and row write; atomic with insert. Test covers nonexistent (9999) and valid case.
Assessment: PASS.

## Finding 8
Location: engine.go:51-68 — NOT NULL + CHECK in InsertUser
Claimed Behavior: NULL email/username -> 23502; age<18 -> 23514; status not in active/suspended/pending -> 23514.
Observed Implementation: Empty-string checks produce `NewNotNullViolation(...)`. Numeric/ enum checks produce `NewCheckViolation(...)`. Tests assert err != nil only.
Assessment: PASS behaviorally (demo prints correct SQLSTATE taxonomy line); WARNING that tests don't assert specific SQLSTATE on the store/engine call path — they assert generic error presence. Acceptable for this lab's scope.
