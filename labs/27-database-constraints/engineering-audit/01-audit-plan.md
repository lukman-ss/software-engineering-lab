# Engineering Audit Plan

Target Lab: labs/27-database-constraints
Implementation Files:
- `internal/model/model.go`
- `internal/dberr/errors.go`
- `internal/engine/engine.go`
- `internal/store/store.go`
- `cmd/demo/main.go`
Tests:
- `internal/store/store_test.go`
Executable/Demo:
- `cmd/demo/main.go`
Approved Research Inputs:
- `research/05-report.md`
- `research-audit/07-verdict.md` (Verdict: APPROVED)
- `engineering/01-design.md`
Main Claims To Verify:
1. Declarative database constraints (NOT NULL, CHECK, UNIQUE, FOREIGN KEY, and PARTIAL UNIQUE INDEX) enforce invariants at the storage layer.
2. Concurrent read-then-write checks in application space fail with duplicate rows, while storage engine UNIQUE constraints eliminate race conditions and enforce SQLSTATE `23505`.
3. Partial unique index conditionally enforces uniqueness for active records (`WHERE deleted_at IS NULL`), permitting multiple soft-deleted records.
4. Error taxonomy standardizes SQLSTATE Class 23 codes (`23502`, `23503`, `23505`, `23514`) with domain error mapping.
Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`
Primary Risks:
- Race conditions during concurrent inserts or soft-deletes.
- In-memory simulation divergence from ANSI/PostgreSQL SQLSTATE semantics.
- Weak test assertions (e.g. asserting non-nil error without validating constraint classification).
