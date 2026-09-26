# Engineering Audit Plan

Target Lab: labs/17-architecture-decision-record
Implementation Files:
- internal/adr/models.go
- internal/adr/parser.go
- internal/adr/linter.go
- cmd/demo/main.go
Tests:
- tests/linter_test.go
- tests/parser_test.go
Executable/Demo: cmd/demo/main.go
Approved Research Inputs: (per PIPELINE OVERRIDE, not audited in this stage)
Main Claims To Verify:
1. ADR parser correctly extracts structured data from Markdown.
2. Linter accepts a valid sequence of ADRs (monotonic numbering, correct statuses, valid supersession references).
3. Linter rejects invalid structural states (broken links, unknown status).
4. Concurrency safety verified with Go race detector on the linter.
Commands To Run:
- go test ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Regex parser may be brittle for edge-case Markdown formatting.
- Concurrency in linter may introduce race conditions if shared state not properly protected.
- Demo may not reflect all claimed behaviors if hardcoded.