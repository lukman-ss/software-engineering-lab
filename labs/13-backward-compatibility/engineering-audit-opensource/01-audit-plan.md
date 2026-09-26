# Engineering Audit Plan

Target Lab: labs/13-backward-compatibility
Audit Date: 2026-09-25
Output Dir: labs/13-backward-compatibility/engineering-audit-opensource/

Implementation Files:
- internal/compat/model.go (DTOs: User, UserResponse, LegacyConsumerDTO, ModernConsumerDTO)
- internal/compat/store.go (MemoryStore, dual-write, contract drop)
- internal/compat/flags.go (WriteMode, ReadMode, ContractApplied)
- internal/compat/metrics.go (Observability counters)
- internal/compat/backfill.go (BackfillWorker, checkpoint)
- internal/compat/service.go (CompatService, reconcile, ApplyContract)
- internal/compat/handler.go (V1/V2 HTTP, Deprecation/Sunset, 410 Gone)
- cmd/demo/main.go (6-step lifecycle demo)
- schema.sql (V1 baseline, V2 expand, V3 contract)

Tests:
- internal/compat/service_test.go (5 unit tests)
- tests/migration_test.go (lifecycle + rollback scenarios)
- tests/concurrency_test.go (writes/reads/backfill/reconcile under race)

Executable/Demo:
- go run ./cmd/demo

Approved Research Inputs:
- research/11-final-research.md (Expand-Migrate-Contract, dual-write/drift, idempotent+resumable backfill, fallback read, flags, observability, zero-traffic contract guard)
- engineering/01-design.md, engineering/02-implementation-notes.md

Main Claims To Verify:
1. Expand additive payload keeps V1 working (phone string) while V2 gets phones array.
2. Dual-write writes both schemas; rollback to N loses no data while dual-write active.
3. Backfill batch + idempotent + resumable via last_processed_id checkpoint.
4. Fallback read serves un-backfilled rows + lazy hydrates modern store.
5. Reconcile detects drift (legacy vs primary) before cutover.
6. Contract guard blocks while legacy traffic > 0; drops column; V1 -> 410, V2 works.
7. Deprecation/Sunset headers on V1.
8. Thread-safe under concurrent reads/writes/backfill (race clean).

Commands To Run:
- go build ./...
- go vet ./...
- go test -count=1 -v ./...
- go test -race -count=1 ./...
- go run ./cmd/demo

Primary Risks:
- In-memory store cannot prove persistent/restart resumability or real DDL locks.
- ApplyContract guard uses cumulative counter (never reset) so non-force path is stricter than time-window semantics.
- Demo forces contract (ApplyContract(true)) so zero-traffic guard proven only by unit test, not demo.
- Minor unhandled errors (Atoi, extra-phone writes) and O(n^2) ID sort.
