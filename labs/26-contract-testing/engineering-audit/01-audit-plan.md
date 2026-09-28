# Engineering Audit Plan

Target Lab: labs/26-contract-testing
Implementation Files:
- internal/contract/verifier.go
- internal/consumer/client.go
- internal/provider/server.go
- internal/model/order.go
- cmd/demo/main.go
Tests:
- tests/contract_test.go
Executable/Demo:
- cmd/demo/main.go
Approved Research Inputs:
- research/05-report.md
- engineering/01-design.md
- engineering/02-implementation-notes.md
Main Claims To Verify:
1. CDC Engine validates JSON schema and payload contracts without running E2E services.
2. Detects breaking API changes (enum casing, field renames, primitive type mutations).
3. CI Gate blocks breaking provider deployment and allows compliant/dual-versioned provider deployment.
4. Concurrency-safe execution during contract verification.
5. README documentation matches actual codebase implementation and run commands.
Commands To Run:
- go test -count=1 ./...
- go test -count=1 -race ./...
- go run ./cmd/demo
Primary Risks:
- Type handling / float vs integer numbers in JSON decoder diff engine.
- Loose or incomplete verification in recursive diff algorithm.
- Unhandled HTTP edge cases or resource leaks in verifier client.
