# Engineering Audit Plan

Target Lab: labs/26-contract-testing
Implementation Files:
- `internal/consumer/client.go`
- `internal/contract/verifier.go`
- `internal/model/order.go`
- `internal/provider/server.go`

Tests:
- `tests/contract_test.go`

Executable/Demo:
- `cmd/demo/main.go`

Approved Research Inputs:
- `research/05-report.md`
- `research/03-evidence.md`
- `research-audit/07-verdict.md`
- `engineering/01-design.md`
- `engineering/02-implementation-notes.md`

Main Claims To Verify:
1. Consumer generates contract with minimal required schema (tolerant reader).
2. Provider verification validates HTTP endpoints against contract interactions.
3. Breaking changes (casing mutation, field removal/rename, primitive type mutation) cause CI verification gate failure.
4. Backward-compatible dual provider (V1 + V2) satisfies V1 contract without breaking existing consumers.
5. Verifier handles concurrent verification safely without race conditions.
6. Execution outputs match claims in documentation.

Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Race conditions during parallel verification requests.
- False positive passes on schema regressions.
- Discrepancies between README instructions and actual binary/code execution.
- Silent failure or incomplete assertion checks in test suites.
