# Engineering Audit Plan

Target Lab: `labs/26-contract-testing`
Implementation Files:
- `internal/contract/verifier.go`
- `internal/model/order.go`
- `internal/provider/server.go`
- `internal/consumer/client.go`

Tests:
- `tests/contract_test.go`

Executable/Demo:
- `cmd/demo/main.go`

Approved Research Inputs:
- `research/01-plan.md`
- `research/05-report.md`
- `research-audit/07-verdict.md`

Main Claims To Verify:
1. Consumer-Driven Contract (CDC) generation defines minimal consumer expectations.
2. Provider V1 passes contract verification and CI gate checks.
3. Breaking changes (enum casing mismatch, missing renamed field, primitive type mutation) fail verification with detailed error reporting.
4. Dual Provider (V1 + V2) maintains backward compatibility while enabling schema evolution.
5. All tests run cleanly under race detector without concurrency issues or memory leaks.
6. Documentation in `README.md` and `engineering/` accurately reflects implementation and demo output.

Commands To Run:
- `cd labs/26-contract-testing && go test -v ./...`
- `cd labs/26-contract-testing && go test -race ./...`
- `cd labs/26-contract-testing && go run ./cmd/demo`

Primary Risks:
- Incomplete JSON diff recursion or loose type assertions leading to false passes.
- Data race or state pollution during concurrent contract verifications.
- Mismatch between documented claims/demo logs and actual runtime behavior.
