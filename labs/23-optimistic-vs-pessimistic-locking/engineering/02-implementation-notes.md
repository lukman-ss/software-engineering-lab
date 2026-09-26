# Implementation Notes

## Files Added
- `go.mod`: Module definition for the lab.
- `internal/inventory/model.go`: Product struct and error definitions (`ErrNotFound`, `ErrInsufficientStock`, `ErrOptimisticLock`, `ErrInvalidQuantity`).
- `internal/inventory/store.go`: Storage engine simulating row-level locking, versioning, and atomic updates.
- `internal/inventory/service.go`: Business service implementing naive, pessimistic, optimistic direct, optimistic with retry, and atomic deduction.
- `tests/locking_test.go`: Concurrency test suite covering naive lost update, pessimistic locking, optimistic conflicts, retry loops, and atomic decrement.
- `cmd/demo/main.go`: Interactive command-line demonstration.
- `engineering/01-design.md`: Engineering design document.
- `engineering/02-implementation-notes.md`: This document.
- `engineering/03-execution-result.md`: Recorded execution output.
- `README.md`: Overview and reproduction steps.

## Core Design Decisions
1. **Simulated Row Locking for Pessimistic Control**: `Store` uses a dedicated `map[int]*sync.Mutex` so each product row can be locked independently without blocking concurrent operations on other product IDs.
2. **Version Checking for Optimistic Control**: Modeled the exact SQL behavior `UPDATE products SET stock=?, version=version+1 WHERE id=? AND version=?` by comparing product versions inside an atomic lock phase.
3. **Exponential Backoff with Jitter**: In `DeductOptimisticWithRetry`, retries apply exponential delay `(1<<attempt)ms` plus randomized jitter (0-5ms) to prevent live-lock under contention.
4. **Statement-Level Conditional Update**: Modeled `UPDATE ... WHERE stock >= qty` via an atomic state modification protected by engine-level lock.

## Implementation-Specific Choices
- **Zero Third-Party Dependencies**: Pure Go standard library (`sync`, `sync/atomic`, `time`, `math/rand`, `testing`) to guarantee maximum portability and instantaneous test execution.
- **Artificial Micro-delays in Naive Method**: Inserted a small `time.Sleep` (100 µs) between read and write in `NaiveDeduct` to reliably replicate the application-level processing window that triggers lost updates in production.

## Known Limitations
- Does not test network latency between an application and a remote RDBMS server.
- Deadlock detection simulation is omitted as lock ordering is single-resource.

## Trade-offs
- **Pessimistic**: High consistency and zero rejected writes, but holds locks during transaction lifespan, reducing throughput under high reader/writer mix.
- **Optimistic**: Non-blocking reads and high throughput under low contention, but high contention incurs CPU/retry overhead and potential `ErrOptimisticLock` exhaustion.
- **Atomic Operations**: Minimal overhead, but restricted to simple in-place arithmetic without complex multi-step validations.

## What Is Demonstrated
- The lost update concurrency anomaly when concurrent processes execute naive read-modify-write.
- Total consistency guarantee with row-level pessimistic locking.
- Detection and safe rejection of concurrent conflicts using optimistic version checking.
- Convergence of optimistic retries using exponential backoff.
- Lockless concurrency protection using single-statement atomic operations.

## What Is Not Demonstrated
- Distributed locking (Redis Redlock, Consul).
- Complex multi-table cross-transaction deadlocks.
