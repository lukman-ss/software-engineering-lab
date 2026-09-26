# Engineering Audit Plan

Target Lab: labs/23-optimistic-vs-pessimistic-locking
Implementation Files:
- internal/inventory/model.go
- internal/inventory/store.go
- internal/inventory/service.go
- cmd/demo/main.go
Tests:
- tests/locking_test.go
Executable/Demo:
- cmd/demo/main.go
Approved Research Inputs:
- engineering/01-design.md
- engineering/02-implementation-notes.md
- engineering/03-execution-result.md
Main Claims To Verify:
1. Naive read-modify-write pattern causes lost updates under concurrent execution.
2. Pessimistic locking (simulating SELECT ... FOR UPDATE) guarantees data consistency under concurrent writes.
3. Optimistic locking correctly detects version mismatch conflicts and rejects stale writes without state corruption.
4. Optimistic locking with exponential backoff retry converges all operations under high concurrency.
5. Atomic single-statement updates maintain exact state invariants safely.
6. Clean compilation, zero race conditions via `go test -race ./...`, and truthful executable demo output.
Commands To Run:
- go test -v ./...
- go test -v -race ./...
- go run ./cmd/demo
Primary Risks:
- Race conditions or shared memory hazards during concurrent store operations.
- Flaky tests dependent on non-deterministic micro-sleep timings.
- Unhandled errors or version overflow issues during high retry volume.
