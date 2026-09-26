# Engineering Design

Target Lab: labs/28-timeouts-and-deadlines
Research Status: APPROVED

## Concept To Prove
Demonstrate context deadline propagation, timeout budgets, exponential backoff with full jitter, circuit breaking integrated with retries, and idempotency key deduplication across simulated distributed service calls.

## Expected Behavior
1. Requests propagate context deadlines across downstream handlers; operations abort immediately once context deadline expires.
2. Exponential backoff with jitter prevents synchronized retry storms.
3. Circuit breaker transitions (Closed -> Open -> Half-Open) prevent useless retries during sustained downstream failures.
4. Idempotent request handlers deduplicate retried operations using idempotency keys.

## Failure Scenario
1. Downstream backend delays response beyond client context deadline.
2. Downstream backend experiences high failure rates, causing circuit breaker to trip Open.
3. Retry attempts without idempotency protection risk duplicate execution (mitigated by deduplication store).

## Success Criteria
1. Context cancellation terminates long-running downstream work without leakage.
2. Race detector (`go test -race ./...`) passes cleanly under concurrent operations.
3. All tests pass verifying deadline propagation, backoff jitter, circuit breaker state transitions, and idempotency deduplication.

## Architecture
- `internal/deadline`: Context deadline budgeting and propagation helpers.
- `internal/retry`: Exponential backoff with jitter and retry policy execution.
- `internal/circuit`: 3-state circuit breaker (Closed, Open, Half-Open).
- `internal/idempotency`: In-memory deduplication store with expiration.
- `cmd/demo`: Executable demonstrating cascading timeout prevention, circuit breaking, and idempotent retries.

## Components
- `Deadline`: `context.WithTimeout`, context propagation.
- `Retrier`: Backoff calculation (`base * 2^attempt + jitter`).
- `CircuitBreaker`: State tracking with mutex sync.
- `IdempotencyStore`: Map with RWMutex for deduplicating request IDs.

## Test Strategy
- Unit tests for each internal package (`deadline`, `retry`, `circuit`, `idempotency`).
- Concurrent execution tests under race detector.

## Execution Plan
1. Create design document.
2. Implement packages in `internal/`.
3. Create runnable demo in `cmd/demo/main.go`.
4. Implement unit tests in `internal/*` and `tests/`.
5. Execute `go test ./...`, `go test -race ./...`, `go run ./cmd/demo`.
6. Write implementation notes and execution results.

## Implementation Decisions
- Standard library `context`, `sync`, `time`, `math/rand/v2` used exclusively for minimal external dependencies.
- Deduplication store uses mutex-protected map without external cache/DB to keep lab runnable stand-alone.
