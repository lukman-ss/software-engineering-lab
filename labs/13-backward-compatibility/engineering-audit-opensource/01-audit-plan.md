# Engineering Audit Plan

Target Lab:
labs/13-backward-compatibility (module `compat`)

Implementation Files:
- internal/compat/model.go — domain models + DTOs
- internal/compat/store.go — thread-safe MemoryStore (users + user_phones)
- internal/compat/flags.go — feature flags (WriteMode, ReadMode, contract)
- internal/compat/metrics.go — observability counters
- internal/compat/backfill.go — resumable, idempotent batch backfill
- internal/compat/service.go — domain facade (read/write/reconcile/contract)
- internal/compat/handler.go — HTTP v1/v2 endpoints + deprecation headers
- cmd/demo/main.go — runnable lifecycle demo

Tests:
- internal/compat/service_test.go — 5 unit tests
- tests/migration_test.go — 2 lifecycle + rollback integration tests
- tests/concurrency_test.go — concurrent read/write/backfill test

Executable/Demo:
- `go run ./cmd/demo`

Approved Research Inputs:
- engineering/01-design.md (design + success criteria)
- engineering/02-implementation-notes.md (decisions/limitations)
- engineering/03-execution-result.md (recorded build/test/demo output)

Main Claims To Verify:
1. Expand -> Migrate -> Contract (parallel change) phases implemented & executable.
2. Dual-write writes atomically to legacy + modern store.
3. Backfill is resumable (checkpoint last_processed_id) and idempotent (no duplicates on rerun).
4. Fallback read serves un-migrated rows from legacy column + lazy backfill.
5. Safe rollback: V1 still reads data written during DualWrite.
6. Contract phase drops legacy schema only after legacy traffic == 0 (guarded).
7. Deprecation/Sunset headers emitted on legacy endpoint.
8. Concurrency: reads+writes+backfill pass `-race` without data races/deadlocks.
9. Demo output matches engineering/03-execution-result.md exactly.

Commands To Run:
- go version
- go build ./...
- go vet ./...
- go test ./... -v -count=1
- go test -race ./... -count=1
- go run ./cmd/demo

Primary Risks:
- Shared store state + concurrent backfill/read/writer races.
- Checkpoint correctness at batch boundaries / empty final batch.
- Contract guard relying on cumulative counter with no reset (force escape hatch).
- Non-atomic contract application across flags + store.
- Doc vs code drift (recorded demo output vs real output).
