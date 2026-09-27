# Engineering Audit Plan

Target Lab: labs/26-contract-testing
Implementation Files:
- internal/model/order.go
- internal/consumer/client.go
- internal/provider/server.go
- internal/contract/verifier.go
- cmd/demo/main.go

Tests:
- tests/contract_test.go

Executable/Demo:
- cmd/demo/main.go

Approved Research Inputs:
- research/05-report.md
- research-audit/07-verdict.md

Main Claims To Verify:
1. Consumer generates contract with minimal required schema subset.
2. Verifier validates running HTTP provider endpoints against contract specification.
3. Breaking provider mutations (enum casing, removed fields, type changes) fail verification and block deployment.
4. Backwards-compatible provider evolution (V1 + V2 dual support) passes verification.
5. All tests pass with `-race` flag enabled without data races.

Commands To Run:
```bash
go test -v ./...
go test -race ./...
go run ./cmd/demo
```

Primary Risks:
- Loose contract validation allowing false negatives on breaking schema changes.
- Concurrency or resource leaks in test HTTP servers / clients.
- Mismatch between README documentation and executable code paths.
