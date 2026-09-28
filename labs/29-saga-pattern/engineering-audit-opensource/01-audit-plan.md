# Engineering Audit Plan

Target Lab: labs/29-saga-pattern

Implementation Files:
- internal/saga/orchestrator.go
- internal/saga/choreography.go
- internal/services/services.go
- cmd/demo/main.go

Tests:
- tests/saga_test.go

Executable/Demo:
- go run ./cmd/demo

Approved Research Inputs:
- (audit implementation and tests only; research/content out of scope per pipeline override)

Main Claims To Verify:
1. Centralized saga coordinator (orchestration) manages forward steps + LIFO compensation.
2. Event bus (choreography) coordinates saga steps across decoupled services.
3. Domain services (Order, Payment, Inventory) maintain local state, idempotency keys, semantic locks.
4. Compensating transactions execute in LIFO (reverse) order on failure.
5. Idempotency prevents duplicate step execution.
6. Semantic locking prevents invalid state transitions.
7. Demo shows happy path and failure compensation.
8. Test suite covers happy path, failure path, rollback, idempotency, semantic lock, concurrency.

Commands To Run:
- go test -v ./... (from labs/29-saga-pattern)
- go test -race ./... (from labs/29-saga-pattern)
- go run ./cmd/demo (from labs/29-saga-pattern)
- go vet ./...

Primary Risks:
- Race conditions in concurrent saga execution (state map access without synchronization).
- LIFO compensation ordering not actually enforced.
- Idempotency keys accepted but not enforced on replay.
- Semantic locking absent or insufficient (reentrant ...