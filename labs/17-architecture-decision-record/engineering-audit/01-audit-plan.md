# Engineering Audit Plan

Target Lab: labs/17-architecture-decision-record
Implementation Files:
- internal/adr/models.go
- internal/adr/parser.go
- internal/adr/linter.go
Tests:
- tests/linter_test.go
- tests/parser_test.go
Executable/Demo:
- cmd/demo/main.go
Approved Research Inputs:
- research/05-report.md
- research-audit/07-verdict.md
Main Claims To Verify:
1. Markdown parser correctly extracts ADR ID, Title, Status, and Superseded/Supersedes relations.
2. Linter enforces monotonic numbering (1, 2, 3...).
3. Linter validates bidirectional integrity of supersession graph (ADR A superseded by ADR B requires ADR B supersedes ADR A).
4. Validation of all lifecycle states (`Proposed`, `Accepted`, `Superseded`, `Deprecated`, `Rejected`).
5. Concurrency safety during validation step.
6. Demo reflects Modular Monolith to Microservices architectural evolution scenario.
Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`
Primary Risks:
- Regex-based parser failing on valid markdown edge cases (whitespace, formatting).
- In-memory mock/string inputs instead of reading real filesystem files.
- Lack of negative tests for duplicate IDs in the test suite.
