# Engineering Audit Plan

Target Lab: `labs/26-contract-testing`
Implementation Files:
- `internal/contract/verifier.go`
- `internal/consumer/client.go`
- `internal/model/order.go`
- `internal/provider/server.go`
- `cmd/demo/main.go`

Tests:
- `tests/contract_test.go`

Executable/Demo:
- `cmd/demo/main.go`

Approved Research Inputs:
- `research/05-report.md`
- `research-audit/07-verdict.md`

Main Claims To Verify:
1. Consumer specifies minimal contract requirements independently.
2. Provider verifier executes interaction tests against provider HTTP endpoints.
3. Breaking changes (enum casing, field rename/omission, primitive type mutation) are reliably detected.
4. Safe API evolution (dual version routing V1/V2) preserves contract compatibility.
5. Verifier and client handle concurrent executions cleanly without race conditions.

Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- False positives in JSON number comparison or loose map assertions.
- Concurrency race conditions in HTTP client or verifier state.
- Documentation vs implementation divergence.
