# Engineering Audit Plan

Target Lab:
- `labs/26-contract-testing`

Implementation Files:
- `go.mod`
- `internal/model/order.go` (domain models + V1/Breaking/V2 DTOs)
- `internal/contract/verifier.go` (CDC verifier engine)
- `internal/consumer/client.go` (consumer client + contract generator)
- `internal/provider/server.go` (ProviderV1 / ProviderBreaking / ProviderDual HTTP handlers)
- `cmd/demo/main.go` (multi-stage demo executable)

Tests:
- `tests/contract_test.go`

Executable/Demo:
- `cd labs/26-contract-testing && go run ./cmd/demo`

Approved Research Inputs:
- N/A (audit scope limited to implementation + tests + README; research not verified here)

Main Claims To Verify:
1. Consumer-Driven Contract (CDC) engine verifies provider responses against consumer expectations.
2. Provider V1 satisfies the MobileApp contract (verification PASS).
3. Provider "Breaking" is detected and blocked (status casing, field rename, int-vs-string type change).
4. Dual provider (V1 + V2) preserves backwards-compatible V1 contract.
5. Verifier is concurrency-safe under parallel verification.
6. Demo output reflects actual verification results.

Commands To Run:
- `go build ./...`
- `go vet ./...`
- `go test -v -count=1 ./...`
- `go test -race -count=1 ./...`
- `go run ./cmd/demo`

Primary Risks:
- Unimplemented/claimed behavior not exercised by tests.
- Verification logic that silently passes when it should fail (or vice-versa).
- Non-thread-safe HTTP client/verifier state under concurrency.
- Demo output diverging from code behavior.
- README claiming test/demo coverage the suite does not provide.
