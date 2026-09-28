# Engineering Audit Plan

Target Lab: labs/26-contract-testing
Implementation Files:
- internal/consumer/client.go
- internal/contract/verifier.go
- internal/model/order.go
- internal/provider/server.go
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
1. Consumer-Driven Contract (CDC) generation by consumer client.
2. Provider V1 passes contract verification.
3. Breaking provider mutations (enum casing, missing/renamed field, primitive type change) fail contract verification with specific error messages.
4. Dual Provider (V1 + V2) allows safe API evolution by passing V1 contract verification while supporting V2 endpoint.
5. Code compiles cleanly and passes unit tests, race detector (`go test -race ./...`), and demo run (`go run ./cmd/demo`).
6. README accurately reflects implementation, test, and demo instructions.
Commands To Run:
- go test ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Partial or missing contract diff detection (e.g., loose JSON structure validation).
- Race conditions during concurrent verifications.
- Mismatch between README documentation and actual file structure/commands.
