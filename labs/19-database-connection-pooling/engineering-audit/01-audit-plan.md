# Engineering Audit Plan

Target Lab: labs/19-database-connection-pooling
Implementation Files: `internal/pool/mockdb.go`, `internal/pool/service.go`
Tests: `tests/pool_test.go`
Executable/Demo: `cmd/demo/main.go`
Approved Research Inputs: `research/05-report.md`, `engineering/01-design.md`
Main Claims To Verify: Connection pooling functionality, concurrency safety, server exhaustion limit enforcement, connection leak starvation.
Commands To Run:
- `go test -v ./...`
- `go test -race -v ./...`
- `go run ./cmd/demo`
Primary Risks: Race conditions in mock driver, invalid simulation of starvation, mismatch between tests and research design.
