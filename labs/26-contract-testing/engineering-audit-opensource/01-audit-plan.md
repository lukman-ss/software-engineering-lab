# Engineering Audit Plan

Target Lab: labs/26-contract-testing
Implementation Files:
- internal/model/order.go
- internal/consumer/client.go
- internal/contract/verifier.go
- internal/provider/server.go
- cmd/demo/main.go
Tests:
- tests/contract_test.go (5 tests)
Executable/Demo: cmd/demo/main.go
Approved Research Inputs: engineering/01-design.md, 02-implementation-notes.md, 03-execution-result.md (implementation+tests scope only per override; research/ not audited)
Main Claims To Verify:
1. Consumer generates minimal CDC contract (id, status, customer.name, total)
2. Provider V1 passes verification + consumer FetchOrder succeeds (CI gate ALLOWED)
3. Breaking provider (enum casing, field rename, int->string) fails verification with >=3 diffs + client fails (CI gate BLOCKED)
4. Dual provider preserves V1 compatibility while exposing V2 (CI gate ALLOWED)
5. Concurrent verification safe under -race
Commands To Run:
- go build ./...
- go test -v ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Verifier subset logic too lax/strict (false PASS/FAIL)
- Client parsing masking contract violations
- Shared Verifier/http.Client race
- Demo output fabricated vs real
