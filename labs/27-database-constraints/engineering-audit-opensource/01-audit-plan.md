# Engineering Audit Plan

Target Lab: labs/27-database-constraints
Implementation Files: 
- internal/engine/engine.go
- internal/model/model.go
- internal/dberr/errors.go
- internal/store/store.go (UnsafeStore, SafeStore)
- cmd/demo/main.go

Tests:
- internal/store/store_test.go

Executable/Demo: go run ./cmd/demo

Approved Research Inputs: 
- engineering/01-design.md (APPROVED)
- engineering/02-implementation-notes.md
- engineering/03-execution-result.md

Main Claims To Verify:
1. Declarative database constraints (NOT NULL, CHECK, UNIQUE, PRIMARY KEY, FOREIGN KEY, and PARTIAL UNIQUE INDEX) enforce relational integrity at storage layer.
2. Under concurrent write workloads, application-only validation fails while DB constraints guarantee consistency with SQLSTATE `23505`.
3. Partial unique index allows conditional uniqueness (multiple soft-deleted, one active).
4. Error mapping from SQLSTATE codes to domain errors is correct.
5. Tests cover all constraint types and concurrency scenarios.
6. Demo output matches claimed behavior.

Commands To Run:
- go build ./...
- go test -v ./...
- go test -race ./...
- go run ./cmd/demo

Primary Risks:
- Tests may pass but lack edge-case coverage.
- Demo output may be hardcoded.
- Implementation may not reflect true concurrency safety.
- Documentation may overclaim capabilities.