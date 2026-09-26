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
- `research/runs/2026-09-26-architecture-decision-record/05-report.md`
- `research-audit/07-verdict.md` (Approved)

Main Claims To Verify:
1. Markdown parser extracts ID, Title, Status, and supersession links (Superseded by N, Supersedes: N).
2. Linter verifies monotonic sequence numbering (1..N).
3. Linter detects broken supersession references and enforces bidirectional supersession consistency.
4. Concurrency safety: Linter executes graph checks concurrently across records without data races.
5. Demo executes real scenario (Modular Monolith -> Microservice extraction + Rejected Event Sourcing) without mock shortcuts or false positives.
6. Documentation in README matches actual implementation, commands, and outputs.

Commands To Run:
```bash
go test -count=1 -v ./...
go test -count=1 -race ./...
go run ./cmd/demo
```

Primary Risks:
- Race conditions in concurrent graph inspection or shared map access.
- Partial supersession declaration without reciprocal link enforcement.
- Divergence between demo text output and documented execution results.
