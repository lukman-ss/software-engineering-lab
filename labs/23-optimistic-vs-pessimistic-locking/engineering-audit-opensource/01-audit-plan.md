# Engineering Audit Plan

Target Lab: labs/23-optimistic-vs-pessimistic-locking
Implementation Files:
- internal/inventory/model.go
- internal/inventory/store.go
- internal/inventory/service.go
- cmd/demo/main.go
- go.mod
Tests:
- tests/locking_test.go (6 tests)
Executable/Demo: cmd/demo/main.go
Approved Research Inputs: (skipped per PIPELINE OVERRIDE — implementation and tests only)
Main Claims To Verify:
1. Naive read-modify-write produces lost update under concurrency
2. Pessimistic row lock preserves exact invariant (100-50=50)
3. Optimistic direct detects version conflict, rejects stale writes, no corruption
4. Optimistic with retry + backoff converges
5. Atomic single-statement conditional update preserves invariant
Commands To Run:
- go build ./...
- go test -v ./...
- go test -race -v ./...
- go run ./cmd/demo
Primary Risks:
- Flaky inverted naive assertion
- Weak retry invariant assertion
- Missing negative/oversell edge tests
- Nondeterministic demo counters misread as fake
