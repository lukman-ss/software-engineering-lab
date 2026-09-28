# Implementation Notes

## Files Added
- `go.mod`: Go module definition
- `internal/saga/orchestrator.go`: Saga orchestrator and step execution engine
- `internal/saga/choreography.go`: EventBus and choreography constructs
- `internal/services/services.go`: Domain mock services (Order, Payment, Inventory) with semantic locks and idempotency keys
- `tests/saga_test.go`: Unit, compensation, concurrency, and choreography tests
- `cmd/demo/main.go`: Runnable console demonstration

## Core Design Decisions
- LIFO compensation rollback: Failed steps immediately invoke compensating actions of preceding executed steps in reverse order.
- Idempotency support: Payment operations check transaction keys to ensure safe retries.
- Semantic lock countermeasure: Order creation applies a semantic lock flag to protect intermediate state during saga lifetime.

## Implementation-Specific Choices
- Standard library Go concurrency primitives (`sync.Mutex`, `sync.RWMutex`).
- In-memory event dispatching for the choreography model.

## Known Limitations
- Compensation functions are assumed to succeed without automatic retry backoff in this lab implementation.
- Distributed persistence / outbox pattern across separate physical databases is simulated in-memory.

## Trade-offs
- Orchestration provides clear linear failure visibility and easier debugging at the cost of coordinator centralization.
- Choreography decouples services via events at the cost of distributed control flow tracing.

## What Is Demonstrated
- Sequential step execution and state updates.
- LIFO compensation order on step failure.
- Idempotent payment processing.
- Semantic locking on pending records.
- Concurrency safety under `-race`.
- Choreography event flow.

## What Is Not Demonstrated
- Persistent transaction logs backed by disk/database.
- Network partitions or asynchronous distributed retry queues.
