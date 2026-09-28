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
- cmd/demo/main.go (run via `go run ./cmd/demo`)
Approved Research Inputs:
- research/05-report.md (approved via content audit)
Main Claims To Verify:
1. Consumer generates CDC JSON contract with minimal schema (id, status, customer.name, total)
2. Provider V1 satisfies contract → verification passes
3. Provider Breaking introduces 3 breaking changes → verification fails with exact diffs (status casing, field rename, type mutation)
4. Provider Dual exposes V2 endpoint while maintaining V1 contract → verification passes
5. CI gate blocks deployment when verification fails (demo shows exit code 1)
6. Demo orchestrates all 4 stages and shows expected outputs
7. Race detector passes (no data races in concurrent verification)
Commands To Run:
- go build ./...
- go test -v ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Verifier diff engine may have false positives/negatives in JSON subset comparison
- Breaking changes may not be detected if verifier logic flawed (mitigated by test coverage)
- Dual provider may accidentally serve V2 on V1 path (routing bug - mitigated by handler prefixes)
- Race conditions in concurrent test (mitigated by sync.WaitGroup)