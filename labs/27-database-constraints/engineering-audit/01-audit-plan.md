# Engineering Audit Plan

Target Lab: labs/27-database-constraints
Implementation Files:
- internal/model/model.go
- internal/dberr/errors.go
- internal/engine/engine.go
- internal/store/store.go
- cmd/demo/main.go
Tests:
- internal/store/store_test.go
Executable/Demo:
- cmd/demo/main.go
Approved Research Inputs:
- research/05-report.md
- research-audit/07-verdict.md (Approved)
Main Claims To Verify:
1. Declarative constraint validation (NOT NULL, CHECK, UNIQUE, FOREIGN KEY, PARTIAL UNIQUE INDEX).
2. SQLSTATE classification and mapping to domain errors (`23502`, `23514`, `23505`, `23503`).
3. Concurrency race protection via atomic database-level constraint enforcement vs application-level race vulnerability.
4. Soft-delete re-registration behavior using partial unique index (`WHERE deleted_at IS NULL`).
Commands To Run:
- `go test -v ./...`
- `go test -race -count=1 ./...`
- `go run ./cmd/demo`
Primary Risks:
- Thread-safety / race conditions in in-memory simulation engine under concurrent access.
- Mismatch between mapped domain errors and SQLSTATE definitions.
- Mocking gaps in partial index logic or foreign key enforcement.
