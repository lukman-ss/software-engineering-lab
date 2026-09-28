# Engineering Design

Target Lab: labs/29-saga-pattern
Research Status: APPROVED

## Concept To Prove
Demonstrate Saga Pattern (Orchestration & Choreography variants) with local transactions, compensating transactions (LIFO rollback), idempotency, and semantic lock countermeasure against data anomalies in Go standard library.

## Expected Behavior
- Orchestrator mode: Central coordinator executes steps sequentially (Order -> Payment -> Inventory -> Delivery).
- Failure handling: Step failure triggers compensating actions in reverse order (LIFO) for completed compensable steps.
- Idempotency: Duplicate executions/retries produce identical state without double-charging or duplicate resource allocation.
- Isolation Countermeasure: Semantic lock prevents concurrent sagas from modifying pending state.

## Failure Scenario
- Payment / Inventory failure during checkout process.
- Intermediate failure causes already committed local transactions to execute business compensations.

## Success Criteria
- Happy path completes all steps and marks saga state SUCCESS.
- Failure path rolls back compensable steps in reverse order and leaves system in consistent state.
- Concurrent executions pass `go test -race ./...`.
- Idempotency verified on retries.

## Architecture
- `pkg/saga`: Core interfaces and orchestrator execution engine.
- `pkg/services`: Mock Order, Payment, Inventory services with local state, compensations, and semantic locking.
- `cmd/demo`: Executable demonstrating happy path, failure rollback, idempotency, and choreography comparison.

## Components
1. `Step`: Encapsulates `Execute` and `Compensate` functions.
2. `Orchestrator`: Manages execution log and compensation stack.
3. `Choreography`: Event-bus based decoupled message passing execution.
4. `Services`: OrderService, PaymentService, InventoryService.

## Test Strategy
- Unit tests: Happy path, failure/rollback path, compensation order verification.
- Concurrency/Race tests: Parallel saga executions with race detection (`go test -race`).
- Idempotency tests: Repeated step execution.

## Execution Plan
1. Create `go.mod` in target directory `labs/29-saga-pattern`.
2. Implement services, orchestrator, and choreography engines.
3. Implement tests in `tests/saga_test.go`.
4. Run `go test ./...`, `go test -race ./...`, `go run ./cmd/demo`.

## Implementation Decisions
- Standard library Go without external heavy framework dependencies.
- In-memory event bus for choreography simulation.
