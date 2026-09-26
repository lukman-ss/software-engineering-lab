# Engineering Audit Plan

Target Lab: labs/26-contract-testing
Implementation Files:
- internal/consumer/client.go (MobileOrderClient, GenerateMobileContract)
- internal/contract/verifier.go (Contract, Interaction, Verifier.Verify, diffValues)
- internal/model/order.go (Order, OrderResponseV1/Breaking/V2)
- internal/provider/server.go (ProviderV1, ProviderBreaking, ProviderDual)
- cmd/demo/main.go (4-stage demo: generate, V1 pass, breaking block, dual pass)
Tests:
- tests/contract_test.go (5 tests: generation, V1 success, breaking fail, dual success, concurrent)
Executable/Demo: cmd/demo/main.go
Approved Research Inputs: PIPELINE OVERRIDE - research/content audit skipped in this stage
Main Claims To Verify:
1. Consumer generates minimal contract (id, status, customer.name, total)
2. Provider V1 passes verification (CI gate ALLOWED)
3. Breaking provider (enum casing, field rename, int->string) fails with >=3 diffs (CI gate BLOCKED)
4. Dual provider preserves V1 compatibility while exposing V2 (CI gate ALLOWED)
5. Mobile client parses V1/dual, rejects breaking schema
6. Concurrent verification safe under race detector
Commands To Run:
- go test ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- diffValues lacks slice/array recursion (known GAP-01, out of scope for single-resource contract)
- Verifier ignores response headers despite contract declaring Content-Type
- Verifier uses shared http.Client concurrently (read-only use, likely safe)
- Client status validation hardcodes enum allowlist (brittle but matches contract)
