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
- `research/06-expand-migrate-contract.md`
- `research/07-deployment-and-rollback.md`
- `research/08-failure-modes.md`
- `research-audit/07-verdict.md` (APPROVED)
- `engineering/01-design.md`
- `engineering/02-implementation-notes.md`

Main Claims To Verify:
1. Expand-Migrate-Contract pattern correctly evolutes 1:1 schema (`users.phone`) to 1:N schema (`user_phones`) without breaking legacy V1 consumers.
2. Dual-write synchronously updates legacy and modern models under feature flag control.
3. Batch backfill worker is resumable via checkpoints and idempotent across repeated invocations.
4. Fallback reading (dual-read) handles reading un-backfilled legacy records in modern mode with lazy backfill.
5. Safe rollback behaves correctly: rolling back during dual-write incurs no data loss for legacy consumers, while stopping dual-write prematurely causes legacy data loss.
6. HTTP handlers return RFC 8594 standard `Deprecation` and `Sunset` headers for legacy endpoints, and return 410 Gone post-contract.
7. Concurrency under `go test -race` demonstrates thread safety across writers, readers, backfill worker, and drift reconciliation.
8. Demo binary executes deterministically with real output matching documented output.

Commands To Run:
```bash
go test -v ./...
go test -race -count=1 ./...
go run ./cmd/demo
```

Primary Risks:
- Race conditions or deadlocks between concurrent writes, fallback read mutations, and batch backfill cursor tracking.
- Desynchronization/drift between legacy `users.phone` and modern `user_phones`.
- Premature drop of legacy column before traffic reaches zero or before consumers migrate.
