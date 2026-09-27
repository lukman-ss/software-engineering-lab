# Engineering Audit Plan

Target Lab: labs/23-optimistic-vs-pessimistic-locking
Implementation Files:
- internal/inventory/model.go
- internal/inventory/service.go
- internal/inventory/store.go
Tests:
- tests/locking_test.go
Executable/Demo:
- cmd/demo/main.go
Approved Research Inputs:
- engineering/01-design.md
- engineering/02-implementation-notes.md
- engineering/03-execution-result.md
README.md
Main Claims To Verify:
1. Naive read-modify-write causes lost updates under concurrency.
2. Pessimistic locking (per-row mutex) guarantees full consistency.
3. Optimistic locking (version guard) detects conflicts without corruption.
4. Optimistic retry with backoff converges all operations.
5. Atomic single-statement conditional decrement guarantees consistency.
6. Race detector passes with zero warnings.
Commands To Run:
- go build ./...
- go test -v ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Weak test assertions masking incorrect behavior.
- Missing edge-case tests (invalid quantity, missing product).
- Demo output could diverge from recorded engineering notes.
