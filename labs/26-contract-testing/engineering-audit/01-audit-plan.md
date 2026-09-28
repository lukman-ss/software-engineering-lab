# Engineering Audit Plan

Target Lab: labs/26-contract-testing
Implementation Files:
- `internal/contract/verifier.go`
- `internal/consumer/client.go`
- `internal/provider/server.go`
- `internal/model/order.go`
Tests:
- `tests/contract_test.go`
Executable/Demo:
- `cmd/demo/main.go`
Approved Research Inputs:
- `research/05-report.md`
- `engineering/01-design.md`
- `engineering/02-implementation-notes.md`
Main Claims To Verify:
1. Consumer generates minimal contract specification (`MobileApp -> OrderService`).
2. Verification engine diffs contract vs running HTTP server response without false positives/negatives.
3. Provider breaking change (type mutation, enum casing, field removal) fails CI gate verification.
4. Dual provider (V1 backwards compatibility + V2 migration route) passes verification for existing consumers.
5. All tests compile and pass under `go test ./...` and `go test -race ./...`.
6. Demo executes cleanly and matches documented expectations.
Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`
Primary Risks:
- Incomplete JSON diff recursion or array/map type assertions causing runtime panics in verifier.
- Flaky concurrent test execution or race conditions in `Verifier` HTTP client usage.
- Mismatch between `README.md` instructions/structure and actual code layout.
