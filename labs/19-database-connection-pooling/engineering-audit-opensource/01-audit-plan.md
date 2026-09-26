# Engineering Audit Plan

Target Lab: labs/19-database-connection-pooling
Implementation Files: 
- internal/pool/mockdb.go
- internal/pool/service.go
- cmd/demo/main.go
Tests: tests/pool_test.go
Executable/Demo: cmd/demo/main.go
Approved Research Inputs: research/05-report.md
Main Claims To Verify:
1. Direct connection creation penalty (overhead)
2. Connection exhaustion when pools are oversized (beyond hardware limits)
3. Connection leaks when connections are held during external I/O
Commands To Run:
- go test ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Mock driver may not accurately simulate real database behavior
- Demo may not reflect production scenarios
- Test coverage might miss edge cases