# Engineering Audit Plan

Target Lab: labs/17-architecture-decision-record
Implementation Files: `internal/adr/models.go`, `internal/adr/parser.go`, `internal/adr/linter.go`
Tests: `tests/parser_test.go`, `tests/linter_test.go`
Executable/Demo: `cmd/demo/main.go`
Approved Research Inputs: `research/05-report.md`, `engineering/01-design.md`, `engineering/02-implementation-notes.md`
Main Claims To Verify:
1. ADR parser extracts Title, Status, and relationships.
2. Linter enforces monotonic numbering and valid status lifecycle.
3. Linter validates supersession lineage (DAG) and detects broken/orphaned references.
4. Demo executes SaaS ERP modular monolith to microservices scenario.
5. Concurrency safety in linter.
Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`
Primary Risks:
- Thread safety issues in concurrent graph validation.
- Weak parser allowing invalid document structures.
- Circular supersession references bypassing linter checks.