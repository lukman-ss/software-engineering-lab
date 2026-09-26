# Engineering Audit Plan

Target Lab: 17-architecture-decision-record
Implementation Files: internal/adr/models.go, internal/adr/parser.go, internal/adr/linter.go
Tests: tests/parser_test.go, tests/linter_test.go
Executable/Demo: cmd/demo/main.go
Approved Research Inputs: engineering/01-design.md
Main Claims To Verify: Validates monotonic numbering, parses markdown ADRs, verifies statuses, ensures no cyclical supersession, catches broken DAG references.
Commands To Run: go test ./..., go test -race ./..., go run ./cmd/demo
Primary Risks: Race conditions in linter, failure to detect cycle, failure to parse missing sections.
