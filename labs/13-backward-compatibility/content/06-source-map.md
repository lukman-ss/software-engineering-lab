# Source Map

## 1. Problem & Core Concepts
- **Research**: `research/03-core-concepts.md`, `research/05-api-compatibility.md`
- **Audit**: `research-audit/03-claim-audit.md`
- **Implementation**: `internal/compat/model.go`

## 2. Expand-Migrate-Contract Pattern
- **Research**: `research/06-expand-migrate-contract.md`
- **Audit**: `research-audit/07-verdict.md`
- **Implementation**: `schema.sql`, `internal/compat/flags.go`
- **Tests**: `tests/migration_test.go`

## 3. Storage, Dual-Write & Fallback Read
- **Engineering Design**: `engineering/01-design.md`
- **Implementation**: `internal/compat/store.go`, `internal/compat/service.go`
- **Tests**: `internal/compat/service_test.go`

## 4. Backfill & Idempotency
- **Research**: `research/04-database-migration.md`
- **Implementation**: `internal/compat/backfill.go`
- **Tests**: `tests/concurrency_test.go`, `internal/compat/service_test.go`

## 5. Rollback, Observability & Contract Phase
- **Research**: `research/07-deployment-and-rollback.md`, `research/08-failure-modes.md`
- **Audit**: `engineering-audit/06-verdict.md`
- **Implementation**: `internal/compat/metrics.go`, `internal/compat/handler.go`
- **Tests**: `tests/migration_test.go`
