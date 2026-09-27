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
- research-audit/07-verdict.md

Main Claims To Verify:
1. Declarative database constraints (NOT NULL 23502, CHECK 23514, UNIQUE 23505, FOREIGN KEY 23503, PARTIAL UNIQUE INDEX) are correctly implemented and enforced in the engine layer.
2. Concurrency race condition prevention: UnsafeStore (app-level check) fails under concurrent registration with race condition duplicates, while SafeStore (DB constraints) guarantees exactly 1 success and 23505 unique violations for remaining concurrent operations.
3. Partial unique index mechanics permit soft-deleted duplicates while strictly enforcing single active record (`WHERE deleted_at IS NULL`).
4. SQLSTATE error taxonomy (`23502`, `23503`, `23505`, `23514`) correctly maps database errors to domain errors.
5. README and engineering notes match code structure, CLI demo, and test execution.

Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go test -count=1 -v ./...`
- `go test -count=1 -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Thread-safety bugs or data races in memory storage engine under high concurrency.
- Inconsistencies between error codes mapped in `dberr` vs SQLSTATE definitions.
- Discrepancy between README documentation claims and actual code package locations.
