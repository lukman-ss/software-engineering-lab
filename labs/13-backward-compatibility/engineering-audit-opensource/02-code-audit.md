# Code Audit

## Finding 1

Location: internal/compat/store.go:16-23 (MemoryStore), service.go:14-19 (Service)
Claimed Behavior: Thread-safe storage; atomic dual-writes; contract drops legacy column.
Observed Implementation: `sync.RWMutex` on MemoryStore. Read methods use `RLock`, mutate methods `Lock`. `GetUser` returns a pointer-copy to prevent external mutation. `ApplyContractDropLegacyColumn` nils every `User.Phone` under write lock, sets `legacyDropped = true`. `CreateDual` writes to both `users.Phone` and `user_phones` under a single `Lock`.
Assessment: PASS
Severity: LOW
Notes: Per-method lock granularity is correct; dual-write is atomic within the store lock. No torn writes observed.

## Finding 2

Location: internal/compat/flags.go:27-31
Claimed Behavior: Concurrent flag updates (canary switching + instant rollback).
Observed Implementation: `writeMode`/`readMode` via `atomic.Value`; `contractApplied` via `atomic.Bool`. Loads are typed-asserted (`f.writeMode.Load().(WriteMode)`). Defaults to `WriteLegacyOnly`/`ReadLegacyOnly` in factory.
Assessment: PASS
Severity: LOW
Notes: Lock-free flag switches are race-free. Atomic.Value requires all stored values to be the same concrete type — satisfied (WriteMode / ReadMode are `int`-based, stored as the typed int, not the interface pointer). Verified under `-race`: PASS.

## Finding 3

Location: internal/compat/backfill.go:88-101 (RunAll) + store.go:147-171 (SavePhoneEntry)
Claimed Behavior: Idempotent + resumable backfill; rerun produces zero duplicates.
Observed Implementation: `BackfillCheckpoint.LastProcessedID` advances only inside the checkpoint mutex; `RunBatch` short-circuits when `IsComplete`. `SavePhoneEntry` checks for an existing `Number` per user and returns the existing entry (no insert) — idempotency at the storage layer. `RunBatch` marks `IsComplete` when `len(ids) < batchSize`.
Assessment: PASS
Severity: LOW
Notes: Test `TestBackfillIdempotentAndResumable` proves rerun migrates 0; second `RunAll` returns 0. Cancellation via `ctx.Done()` mid-batch leaves `LastProcessedID` at the last fully-processed id; `SavePhoneEntry` is idempotent, so a re-run cannot create duplicates.

## Finding 4

Location: internal/compat/service.go:229-245 (ApplyContract)
Claimed Behavior: Contract applied only when legacy traffic == 0; guarded by metric counters.
Observed Implementation: Guard checks `obs.LegacyReadHits.Load() > 0` (cumulative, never reset). Non-forced apply returns `ErrContractViolation`. Forced apply proceeds. Order: SetWriteMode(NewOnly) -> SetReadMode(NewOnly) -> SetContractApplied(true) -> store.ApplyContractDropLegacyColumn().
Assessment: WARNING
Severity: LOW
Notes: The guard uses a cumulative counter with no reset/sliding-window. In a live system, any single legacy hit ever observed blocks forced migration unless `force=true`. The lab honestly documents this via the `force` parameter ("simulating after 30-day zero legacy traffic window"). `ApplyContract(true)` in demo/test is the documented escape. Not a correctness flaw; a fidelity gap vs "sliding zero-traffic window."

## Finding 5

Location: internal/compat/handler.go:20-22,44-46
Claimed Behavior: HTTP endpoints return deprecation headers + status codes; 410 Gone after contract.
Observed Implementation: Both handlers do `id, _ := strconv.Atoi(idStr)` — error is discarded. Invalid/missing id becomes `id=0` -> handled as `user not found` (404 via store, not a 400 Bad Request). Post-contract V1 returns 410. V2 returns 200 JSON.
Assessment: WARNING
Severity: LOW
Notes: Missing input validation. `strconv.Atoi` error ignored. Acceptable graceful 404 fallback for a lab, but no explicit 400 path for malformed `id`. Not a crash or data-corruption path.

## Finding 6

Location: internal/compat/store.go:173-196 (GetUserIDs)
Claimed Behavior: Returns IDs greater than `afterID`, sorted ascending, limited.
Observed Implementation: Collects matching ids under `RLock`; uses an O(n^2) bubble sort instead of `sort.Ints`.
Assessment: PASS
Severity: LOW
Notes: Correct ordering (asc + limit), which is essential for deterministic checkpoint progression. Bubble sort is needlessly inefficient but functionally correct for the lab's small dataset (50 pre-pop + 200 concurrent writers in the concurrency test). Trivial upgrade path: `sort.Ints(ids)`.

## Finding 7

Location: internal/compat/service.go:98-148 (GetUser)
Claimed Behavior: Fallback read falls back to legacy column when new store empty, then lazy-backfills.
Observed Implementation: `ReadFallback` -> if `phones` non-empty use modern, else if `u.Phone != nil` use legacy value + `SavePhoneEntry` (lazy backfill). Reads `u` (snapshot) and `phones` (snapshot) via separate store calls -> not a single atomic snapshot.
Assessment: PASS
Severity: LOW
Notes: The two separate `RLock` snapshots can see a concurrent mutation between them, but no data race and no corruption — only a benign transitional inconsistency during a live dual-write window (documented behavior in design doc). `SavePhoneEntry` idempotent so lazy backfill is race-safe.

## Finding 8

Location: internal/compat/model.go:19-25
Claimed Behavior: Enriched additive payload supports both V1 (string `phone`) and V2 (`phones` array) consumers simultaneously.
Observed Implementation: `UserResponse` carries both `Phone string` (omitempty) and `Phones []PhoneEntry`. V1 DTO keeps only `phone`; V2 DTO keeps only `phones` -> additive compatibility proven at serialization boundary.
Assessment: PASS
Severity: LOW
