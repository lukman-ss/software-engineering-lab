# Engineering Audit Plan

Target Lab: labs/27-database-constraints
Implementation Files: internal/dberr/errors.go, internal/model/model.go, internal/engine/engine.go, internal/store/store.go, cmd/demo/main.go
Tests: internal/store/store_test.go
Executable/Demo: cmd/demo/main.go
Approved Research Inputs: research/01-plan.md, research/05-report.md
Main Claims To Verify:
- Enforces SQL database-level constraints (NOT NULL, CHECK, UNIQUE, FOREIGN KEY, PARTIAL UNIQUE INDEX).
- Maps PostgreSQL SQLSTATE codes (`23502`, `23514`, `23505`, `23503`) to domain error types.
- Demonstrates race condition prevention via storage engine constraints under concurrent stress.
- Demonstrates soft-delete email re-registration via partial unique index (`WHERE deleted_at IS NULL`).
Commands To Run:
- `go test -v ./...`
- `go test -count=1 -race ./...`
- `go run ./cmd/demo`
Primary Risks:
- Simulated in-memory storage engine failing to accurately model thread-safe transaction isolation or lock behavior.
- Soft delete state mutation unindexed or missing validation.
