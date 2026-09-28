# Engineering Audit Plan

Target Lab: labs/26-contract-testing
Implementation Files:
- internal/contract/verifier.go
- internal/model/order.go
- internal/provider/server.go
- internal/consumer/client.go
Tests: tests/contract_test.go
Executable/Demo: cmd/demo/main.go
Approved Research Inputs: labs/26-contract-testing/research/ (plan, sources, evidence, report)
Main Claims To Verify:
1. Consumer-driven contract generation produces correct JSON schema
2. Provider V1 satisfies consumer contract (verification passes)
3. Breaking provider fails verification with exact field diffs (enum casing, field rename, type change)
4. Dual provider maintains V1 compatibility while exposing V2
5. Concurrent verification passes race detector
6. Demo CLI shows full lifecycle
Commands To Run:
- go build ./...
- go test -v ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Race conditions in shared state (mitigated by httptest isolation)
- Mismatch between documented behavior and actual implementation
- Incomplete test coverage of error branches (partially addressed in revision)