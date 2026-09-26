# Engineering Audit Plan

Target Lab: labs/19-database-connection-pooling
Implementation Files: internal/pool/mockdb.go, internal/pool/service.go
Tests: tests/pool_test.go
Executable/Demo: cmd/demo/main.go
Approved Research Inputs: research-audit/07-verdict.md, research/03-evidence.md, research/05-report.md
Main Claims To Verify:
1. Direct connection overhead degrades performance compared to pooling.
2. Oversized pools exhausting backend limits cause failures.
3. External network calls holding DB locks cause pool starvation.
Commands To Run:
- `go test -v ./...`
- `go test -race -v ./...`
- `go run ./cmd/demo`
Primary Risks:
- Simulated mock DB may mask real `database/sql` concurrency bugs.
- Race conditions in mock driver tracking connections.