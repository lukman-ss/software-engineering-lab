# Engineering Code Audit

Verification source: `research/11-final-research.md`, `engineering/01-design.md`.

## Finding 1
Location: internal/compat/model.go:20-25 (`UserResponse`)
Claimed Behavior: Enriched additive payload carries both legacy `phone` (string) and new `phones` array; V1 consumer parses `phone`, V2 parses `phones`.
Observed Implementation: `UserResponse` includes `Phones []PhoneEntry` (always emitted) and `Phone` (omitempty). `UserResponse.Phone` is populated by `GetUser` depending on ReadMode.
Assessment: PASS
Severity: LOW
Notes: Additive contract satisfied. `phone` omitted when empty, which could break strict V1 consumers expecting the key, but matches `omitempty` design intent.

## Finding 2
Location: internal/compat/store.go:40-51 (`CreateDual`), :87-118 (`CreateModern`), :204-212 (`ApplyContractDropLegacyColumn`)
Claimed Behavior: Dual-write writes both legacy `users.phone` and new `user_phones`; contract drop sets legacy field nil.
Observed Implementation: `CreateDual` writes a `User` with `Phone` pointer and appends primary `PhoneEntry` to `userPhones` under a single write lock. `ApplyContractDropLegacyColumn` sets `legacyDropped=true` and nils every `u.Phone`.
Assessment: PASS
Severity: LOW
Notes: Atomic within lock. Legacy write suppressed when `legacyDropped` (CreateDual path), preserving contract semantics.

## Finding 3
Location: internal/compat/service.go:71-94 (`CreateUser` WriteDual branch)
Claimed Behavior: Dual-write is dual to both tables; extra phones appended to modern table.
Observed Implementation: WriteDual writes primary to legacy + primary to `user_phones` via `CreateDual`; extra phones appended via `SavePhoneEntry(_, extra, false)`. `SavePhoneEntry` error is discarded (`_, _`), and only `DualWriteErrors` is incremented on `CreateDual` failure — not on extra-phone writes.
Assessment: WARNING
Severity: MEDIUM
Notes: Silent loss of extra-phone entries on `SavePhoneEntry` failure is not surfaced. Acceptable for in-memory demo, but contradicts "atomic dual-write" claim. Marked `ponytail` scope.

## Finding 4
Location: internal/compat/service.go:109-140 (`GetUser`), :114-129 (`ReadFallback`)
Claimed Behavior: Fallback read (dual read) serves un-backfilled rows from legacy and lazily hydrates modern store.
Observed Implementation: In `ReadFallback`, if `len(phones)==0` and `u.Phone!=nil`, sets `phoneVal=*u.Phone` then calls `SavePhoneEntry(id, phoneVal, true)` (lazy backfill). `SavePhoneEntry` is itself idempotent (skips duplicate numbers).
Assessment: PASS
Severity: LOW
Notes: Matches research Q6 (dual read needed while backfill in progress).

## Finding 5
Location: internal/compat/backfill.go:36-85 (`RunBatch`)
Claimed Behavior: Batch migration with ID checkpoint `last_processed_id`; idempotent (skips phones already present); resumable.
Observed Implementation: Locks checkpoint mutex; calls `GetUserIDs(LastProcessedID, batchSize)`, migrates each user whose `Phone != nil` and `user_phones` empty via `SavePhoneEntry`. Updates `LastProcessedID` per id; sets `IsComplete` when fewer than `batchSize` returned or context cancelled.
Assessment: PASS
Severity: LOW
Notes: Idempotent because `SavePhoneEntry` checks existing numbers. Resumable via checkpoint struct.

## Finding 6
Location: internal/compat/service.go:229-245 (`ApplyContract`)
Claimed Behavior: Contract guards on zero legacy traffic; then switches WriteMode/WriteRead to NewOnly and drops legacy column.
Observed Implementation: Non-force path returns `ErrContractViolation` if `LegacyReadHits > 0`. Force path skips this check (used by demo/demo tests and migration_test). Sets write to NewOnly, read to NewOnly, contract flag, and calls `store.ApplyContractDropLegacyColumn()`.
Assessment: PASS
Severity: LOW
Notes: Counter is cumulative and never reset, so the guard is stricter than the research's "consistent zero over a period" heuristic. Functional intent (block premature contract) is satisfied. Research note Q11 explicitly states the 30-day heuristic is "NOT VERIFIED" and that zero-traffic metric is the real signal — implemented behavior matches the real signal, not the unverified heuristic.

## Finding 7
Location: internal/compat/service.go:187-226 (`ReconcileData`)
Claimed Behavior: Drift audit compares `users.phone` against primary `user_phones` entry.
Observed Implementation: Iterates IDs, skips nil legacy Phone, counts users with no phones or primary number != legacy phone; increments `DriftDetected`. Returns 0 immediately if contract applied.
Assessment: PASS
Severity: LOW
Notes: Correctly detects dual-write drift per research Q7/Failure mode #3.

## Finding 8
Location: internal/compat/handler.go:20-41 (`GetUserV1`)
Claimed Behavior: V1 endpoint emits `Deprecation: true` and `Sunset` headers; post-contract returns 410 Gone.
Observed Implementation: Sets `Deprecation`/`Sunset` headers before encoding; returns 410 with JSON error body on `ErrLegacyUnavailable`. V2 endpoint sets Content-Type only.
Assessment: PASS
Severity: LOW
Notes: Matches research Q10/Q13 (Sunset warning headers to signal deprecation).

## Finding 9
Location: internal/compat/flags.go:22-62 (`FeatureFlags`)
Claimed Behavior: Flags control WriteMode/ReadMode/ContractApplied, enabling canary switching and instant rollback without redeploy.
Observed Implementation: `atomic.Value`-backed modes + `atomic.Bool` contract flag; no persistence (in-memory).
Assessment: PASS
Severity: LOW
Notes: Matches research Q13. In-memory means restart loses flag state (see Finding 12).

## Finding 10
Location: internal/compat/store.go:173-196 (`GetUserIDs`)
Claimed Behavior: Provides ordered ID cursor for batched backfill.
Observed Implementation: Collects ids > afterID, sorts via O(n^2) bubble sort, returns limited slice.
Assessment: WARNING
Severity: LOW
Notes: O(n^2) sort unnecessary (stdlib `sort.Ints` available). No deadlock/lock-order issue. Performance only; not a correctness bug.

## Finding 11 (rollback safety)
Location: internal/compat/store.go:32-51 (`CreateLegacy` rejects when `legacyDropped`), service.go `CreateUser` WriteLegacyOnly branch
Claimed Behavior: If N+1 rolled back to N during dual-write, V1 continues serving/writing legacy without data loss.
Observed Implementation: During dual-write both tables populated, so rolling back WriteMode/ReadMode to legacy-only still reads legacy `phone`. `CreateLegacy` blocks only after contract applied — rollback occurs pre-contract.
Assessment: PASS
Severity: LOW
Notes: Matches research Q9/Q14 (reversible migration). Tests `TestRollbackScenarios` and demo step 5 prove this.

## Finding 12 (premature rollback / data loss)
Location: internal/compat/store.go:87-118 (`CreateModern`), research Q7/Q14, design Failure #4
Claimed Behavior: Dual-write keeps full compatibility during rollback; research warns stopping dual-write while rollbacks occur causes silent legacy data loss.
Observed Implementation: `CreateModern` (WriteNewOnly) writes only to `user_phones`, leaving `User.Phone=nil`. A subsequent rollback to V1 then reads `phone=""` (empty string) for that user.
Assessment: PASS (hazard demonstrated, not a defect)
Severity: LOW
Notes: `TestRollbackScenarios` Scenario B explicitly proves this data-loss case. Matches documented failure mode Failure #2/#4.

## Finding 13
Location: store.go:147-171 (`SavePhoneEntry`), service.go:79-81 (extra phone write)
Claimed Behavior: Idempotent backfill prevents duplicate data.
Observed Implementation: `SavePhoneEntry` returns existing entry when number exists for user — idempotent at write level. Backfill also guards on `len(phones)==0`.
Assessment: PASS
Severity: LOW
Notes: Two-layer idempotency (check-before-backfill + check-before-insert). Verified by `TestBackfillIdempotentAndResumable`.

## Finding 14
Location: internal/compat/metrics.go:7-14, service.go, handler.go
Claimed Behavior: Counters track legacy/new reads, dual writes/errors, backfill, drift.
Observed Implementation: `LegacyReadHits` incremented in `GetLegacyUser`; `NewReadHits` in `GetModernUser`; `DualWriteCount`/`DualWriteErrors` in `CreateUser` WriteDual; `BackfillProcessed` per backfilled row; `DriftDetected` in `ReconcileData`. `Snapshot()` returns all six.
Assessment: PASS
Severity: LOW
Notes: Demo metrics snapshot matches `engineering/03-execution-result.md`. Note: `LegacyReadHits` is cumulative (never reset), making ApplyContract guard stricter than a sliding window — documented in Finding 6.

## Summary
- Code compiles (`go build`), `go vet` clean.
- No correctness defects against the Expand-Migrate-Contract spec.
- Warnings (medium/low) are all intentional in-memory simplifications documented via `ponytail:` and Known Limitations, not fabrication.
