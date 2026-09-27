# Code Audit — labs/27-database-constraints

Files reviewed: internal/engine/engine.go, internal/store/store.go,
internal/model/model.go, internal/dberr/errors.go, cmd/demo/main.go

## Finding 1

Location: internal/engine/engine.go:46-99 (InsertUser)
Claimed Behavior: NOT NULL (23502), CHECK (23514), UNIQUE (23505), partial unique enforced atomically at storage layer.
Observed Implementation: Single `mu.Lock` covers all checks + row/index commit. Returns before mutation on any violation. Atomic, no partial commit.
Assessment: PASS
Severity: LOW
Notes: Correct ordering NOT NULL -> CHECK -> UNIQUE matches SQL semantics.

## Finding 2

Location: internal/engine/engine.go:71-96 (UNIQUE vs partial branches)
Claimed Behavior: UNIQUE and partial unique index both prevent concurrent duplicates.
Observed Implementation: Two disjoint maps (`emailIndex` vs `activeEmails`); `usePartialUniqueIndex` flag selects one. Mixing `RegisterUser` + `RegisterUserPartial` on same email in one engine bypasses cross-mode uniqueness.
Assessment: WARNING
Severity: LOW
Notes: Tests/demo never mix modes on one engine, so untriggered. Real table would enforce both simultaneously.

## Finding 3

Location: internal/engine/engine.go:102-120 (SoftDeleteUser)
Claimed Behavior: Soft delete frees active email for reuse.
Observed Implementation: Deletes from `activeEmails`, sets `DeletedAt`. Correct. Signature `SoftDeleteUser(id int64, deletedAt model.User)` abuses full User struct to carry one timestamp; nil guard present.
Assessment: WARNING
Severity: LOW
Notes: Awkward API, not a correctness bug.

## Finding 4

Location: internal/engine/engine.go:123-148 (InsertOrder)
Claimed Behavior: NOT NULL user_id (23502), CHECK total_cents>0 (23514), FK existence (23503).
Observed Implementation: All three enforced under one write lock before mutation. FK checks `e.users` existence.
Assessment: PASS
Severity: LOW
Notes: No ON DELETE / cascade handling — acceptable lab scope, disclosed as non-goal.

## Finding 5

Location: internal/store/store.go:24-47 (UnsafeStore.RegisterUser)
Claimed Behavior: Intentionally vulnerable read-then-write (count snapshot + sleep + unconstrained insert).
Observed Implementation: RLocks via GetUsersCount/GetUser, `time.Sleep(1ms)` widens window, `InsertUserUnsafe` skips all validation. No data race (all map access locked); logical race reproducible — race test asserts count>1 and passes.
Assessment: PASS
Severity: LOW
Notes: ID scan `1..count` misses non-contiguous IDs, but that only weakens the app check further — consistent with demo intent.

## Finding 6

Location: internal/store/store.go:59-75 + internal/dberr/errors.go:83-100 (MapToDomainError)
Claimed Behavior: Error mapping bridges SQLSTATE codes to domain errors.
Observed Implementation: Maps to message containing rule name but drops `Code`. Callers of SafeStore cannot use `IsConstraintViolation` on returned errors.
Assessment: WARNING
Severity: MEDIUM
Notes: Tests only classify raw constructor errors, never store-mapped errors. Taxonomy proven at `dberr` layer, not end-to-end.

## Finding 7

Location: internal/engine/engine.go:60-62 (CHECK age>=18 on plain int)
Claimed Behavior: Design mentions CHECK "passing NULL if nullable".
Observed Implementation: No nullable Age field; zero-value Age=0 always fails CHECK. No NULL pass-through path exists.
Assessment: WARNING
Severity: LOW
Notes: Untested claim fragment; harmless since no nullable check column is exercised.

## Finding 8

Location: concurrency (Engine mutex + atomic seq; SafeStore concurrent test)
Claimed Behavior: Exactly 1 success under concurrent same-email inserts.
Observed Implementation: `sync.RWMutex` + `atomic.Int64` correct; verified live: demo 50 workers -> 1 success / 49 rejected; `go test -race` clean.
Assessment: PASS
Severity: LOW
Notes: No shared-state race; error counting in demo/test uses guarded mutex.

## Finding 9

Location: timeout / recovery / rollback / cleanup
Claimed Behavior: None claimed (no timeouts, no multi-step txn).
Observed Implementation: Single-row inserts return before mutation on violation — implicitly atomic. No resources to leak.
Assessment: PASS
Severity: LOW
Notes: NOT_APPLICABLE beyond atomicity, which holds.

## Finding 10

Location: complexity / dependencies (go.mod: stdlib only)
Claimed Behavior: Zero-dependency in-memory engine.
Observed Implementation: `sync`, `sync/atomic`, maps only. Minimal, readable.
Assessment: PASS
Severity: LOW
Notes: Coarse table lock vs page latches is disclosed trade-off.
