# Engineering Audit Plan

Target Lab: `labs/26-contract-testing`
Implementation Files:
- `internal/contract/verifier.go`
- `internal/consumer/client.go`
- `internal/provider/server.go`
- `internal/model/order.go`
- `go.mod`
Tests:
- `tests/contract_test.go`
Executable/Demo:
- `cmd/demo/main.go`
Approved Research Inputs:
- `research/05-report.md`
- `research-audit/07-verdict.md`
- `engineering/01-design.md`
- `engineering/02-implementation-notes.md`
- `engineering/03-execution-result.md`
Main Claims To Verify:
1. Consumer-Driven Contract (CDC) generation generates interaction specification containing minimal schema requirements.
2. Verification engine checks method, path, headers, status, and recursive JSON body differences with type, value, and field presence checks.
3. Breaking changes (enum casing mismatch, field renaming, primitive type mutation) cause contract verification failure with informative errors.
4. Provider V1 and Dual Provider (V1 + V2) pass contract verification.
5. Verification engine is safe under concurrent execution (`go test -race`).
6. Demo CLI execution matches claimed output and exits cleanly.
Commands To Run:
- `go test -v -count=1 ./...`
- `go test -race ./...`
- `go run ./cmd/demo`
Primary Risks:
- Deep JSON diffing edge cases (missing nested objects vs primitive types).
- Number representation handling in Go `encoding/json` (`json.Number` vs `float64` / `int64`).
- Concurrency race conditions in HTTP client or test doubles.
- Discrepancy between README instructions and codebase.
