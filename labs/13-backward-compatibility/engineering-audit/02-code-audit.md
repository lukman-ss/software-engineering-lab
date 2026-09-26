# Code Audit

Target Lab: labs/13-backward-compatibility

## Finding 1

Location: `internal/compat/store.go:16-86`
Claimed Behavior: Thread-safe in-memory storage supporting legacy, dual-write, and modern mutations.
Observed Implementation: Uses `sync.RWMutex` protecting `users` and `userPhones` maps. `CreateDual` synchronously populates both `users.phone` and `userPhones` entry.
Assessment: PASS
Severity: LOW
Notes: Properly documented with `ponytail: in-memory mock storage; replace with database/sql for persistent store.` Meets lab requirements without introducing external database dependencies.

## Finding 2

Location: `internal/compat/backfill.go:36-85`
Claimed Behavior: Resumable and idempotent batch backfill worker.
Observed Implementation: Checkpoint tracks `LastProcessedID` protected by mutex. Fetches sorted IDs strictly greater than `lastProcessedID`. In `RunBatch`, calls `SavePhoneEntry`, which checks for duplicate numbers before insertion, guaranteeing idempotency.
Assessment: PASS
Severity: LOW
Notes: Correctly handles context cancellation via `ctx.Done()`.

## Finding 3

Location: `internal/compat/service.go:98-148`
Claimed Behavior: Fallback read (dual-read) with lazy backfill when reading unmigrated records.
Observed Implementation: In `ReadFallback` mode, if `len(phones) == 0` and legacy `u.Phone != nil`, it reads `u.Phone` and triggers lazy persistence to `userPhones` via `s.store.SavePhoneEntry(id, phoneVal, true)`.
Assessment: PASS
Severity: LOW
Notes: Prevents data starvation during the migration window before batch backfill runs.

## Finding 4

Location: `internal/compat/service.go:229-245`
Claimed Behavior: Contract phase guarded against premature execution when legacy traffic is still active.
Observed Implementation: `ApplyContract(force bool)` checks `s.obs.LegacyReadHits.Load() > 0` when `force == false`, returning `ErrContractViolation`. When satisfied or forced, switches write mode to `WriteNewOnly`, read mode to `ReadNewOnly`, sets `contractApplied` flag, and invokes `s.store.ApplyContractDropLegacyColumn()`.
Assessment: PASS
Severity: LOW
Notes: Reliably prevents premature column retirement.

## Finding 5

Location: `internal/compat/handler.go:20-41`
Claimed Behavior: HTTP API handler returns RFC 8594 deprecation and sunset headers for legacy endpoints.
Observed Implementation: `GetUserV1` sets `Deprecation: true` and `Sunset: Mon, 31 Dec 2026 23:59:59 GMT`. Returns `410 Gone` once the contract has dropped the legacy endpoint.
Assessment: PASS
Severity: LOW
Notes: Standard-compliant header emission verified.

## Finding 6

Location: `internal/compat/metrics.go:7-29`
Claimed Behavior: Thread-safe observability counters tracking traffic and operational events.
Observed Implementation: Uses `sync/atomic.Int64` for all metric counters (`LegacyReadHits`, `NewReadHits`, `DualWriteCount`, `DualWriteErrors`, `BackfillProcessed`, `DriftDetected`).
Assessment: PASS
Severity: LOW
Notes: No lock contention on metric updates. Concurrency safe.
