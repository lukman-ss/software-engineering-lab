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
Approved Research Inputs: research/ (not required per pipeline override)
Main Claims To Verify:
1. ADR parser correctly extracts structured data from Markdown files
2. Linter correctly accepts a valid sequence of ADRs 
3. Linter correctly rejects invalid structural states (broken links, unknown status)
4. Concurrency safety verified with Go race detector
Commands To Run:
- go test ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Parser regex limitations may not handle edge cases in Markdown
- Concurrent linter implementation may have race conditions
- Demo may not accurately reflect research claims