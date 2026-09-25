# Engineering Audit Plan

Target Lab: labs/17-architecture-decision-record
Implementation Files:
- `internal/adr/models.go`
- `internal/adr/parser.go`
- `internal/adr/linter.go`
Tests:
- `tests/parser_test.go`
- `tests/linter_test.go`
Executable/Demo:
- `cmd/demo/main.go`
Approved Research Inputs:
- `research/runs/2026-09-25-architecture-decision-record/05-report.md`
- `research-audit/07-verdict.md` (APPROVED_WITH_WARNINGS)
Main Claims To Verify:
1. ADR parser extracts title, monotonic ID, status, and bidirectional supersession relationships from markdown.
2. Linter enforces monotonic numbering (starting from 1, contiguous).
3. Linter validates bidirectional integrity for supersession DAG (`Superseded by X` matches `Supersedes Y`).
4. Linter detects broken references to non-existent ADRs.
5. Concurrency safety holds during parallel rule evaluation across records.
6. Execution demo faithfully illustrates the approved case scenario (Modular Monolith ERP to Microservice transition and Event Sourcing rejection).
Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`
Primary Risks:
- Race conditions or data races when collecting validation errors in parallel goroutines.
- Missing edge case validations (duplicate IDs, self-referential supersession, inverted temporal ordering).
- Non-deterministic error ordering across concurrent checks.
