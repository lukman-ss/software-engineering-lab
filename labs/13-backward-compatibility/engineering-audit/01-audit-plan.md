# Engineering Audit Plan

Target Lab: labs/13-backward-compatibility
Implementation Files:
- `internal/compat/model.go`
- `internal/compat/store.go`
- `internal/compat/flags.go`
- `internal/compat/metrics.go`
- `internal/compat/backfill.go`
- `internal/compat/service.go`
- `internal/compat/handler.go`
- `schema.sql`

Tests:
- `internal/compat/service_test.go`
- `tests/migration_test.go`
- `tests/concurrency_test.go`

Executable/Demo:
- `cmd/demo/main.go`

Approved Research Inputs:
- `research/03-core-concepts.md`
- `research/04-database-migration.md`
- `research/05-api-compatibility.md`
- `research/06-expand-migrate-contract.md`
- `research/07-deployment-and-rollback.md`
- `research/08-failure-modes.md`
- `research/11-final-research.md`
- `research-audit/07-verdict.md`

Main Claims To Verify:
1. Expand-Migrate-Contract lifecycle executes without breaking legacy V1 consumers or starving modern V2 consumers.
2. Dual-write writes atomically to both legacy and modern representations.
3. Batch backfill worker is resumable, idempotent, and non-duplicative.
4. Fallback read (dual-read) hydrates missing records lazily without read failures.
5. Observability metrics accurately count legacy reads, modern reads, dual writes, backfill progress, and drift detection.
6. Safe rollback is verified during dual-write, and unsafe rollback consequences are proven when dual-write is stopped prematurely.
7. Contract phase enforces zero legacy traffic guard before dropping legacy columns.
8. Concurrent execution of reads, writes, and backfills is data-race free under `go test -race`.

Commands To Run:
- `cd labs/13-backward-compatibility && go test -count=1 ./...`
- `cd labs/13-backward-compatibility && go test -count=1 -race ./...`
- `cd labs/13-backward-compatibility && go run ./cmd/demo`

Primary Risks:
- Thread safety in shared in-memory state under high concurrency.
- Incomplete fallback lazy hydration edge cases.
- Discrepancy between SQL schema in `schema.sql` and in-memory store simulation.
- Overclaiming production database capabilities (e.g. distributed locking or ACID outbox) in an in-memory simulation.
