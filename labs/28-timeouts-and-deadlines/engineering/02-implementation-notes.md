# Implementation Notes

## Files Added
- `go.mod`: Module definition for Go 1.22+.
- `README.md`: Lab description and execution instructions.
- `cmd/demo/main.go`: End-to-end runnable demonstration of deadline propagation, backoff jitter, circuit breaker, and idempotency deduplication.
- `internal/deadline/deadline.go`, `deadline_test.go`: Budgeted execution context wrapping.
- `internal/retry/retry.go`, `retry_test.go`: Retrier implementing exponential backoff with full jitter.
- `internal/circuit/circuit.go`, `circuit_test.go`: Thread-safe three-state circuit breaker.
- `internal/idempotency/idempotency.go`, `idempotency_test.go`: Thread-safe deduplication store with TTL.
- `tests/integration_test.go`: Integrated tests combining retrier, circuit breaker, and idempotency store.
- `engineering/01-design.md`: Design document.
- `engineering/02-implementation-notes.md`: Implementation notes.
- `engineering/03-execution-result.md`: Verified execution logs.

## Core Design Decisions
- **Full Jitter Algorithm**: Implemented `sleep = rand.Float64() * min(maxBackoff, baseBackoff * 2^(attempt-1))` as recommended by AWS Architecture Blog and Google SRE to avoid thundering herd.
- **Context Deadline Propagation**: Sub-calls inherit parent context bounds, and context timeout overrides local function timeout budget.
- **Circuit Breaker Integration**: Circuit breaker transitions into Open when reaching consecutive failures threshold, immediately rejecting subsequent retries without hammering downstream dependencies.
- **Idempotency Deduplication**: Memory store with TTL mapping request tokens to responses to ensure retried non-idempotent operations do not trigger duplicate mutations.

## Implementation-Specific Choices
- In-memory synchronization primitives (`sync.Mutex`, `sync.RWMutex`, `context.Context`) chosen over external infrastructure (Redis, PostgreSQL) to keep the lab entirely self-contained.

## Known Limitations
- The in-memory idempotency store does not persist across application crashes. In distributed systems, a persistent store (e.g., transactional DB table or Redis with TTL) is required.

## Trade-offs
- Standard Go library channels and select statements used instead of heavy 3rd-party resilience frameworks (e.g., Resilience4j / Hystrix ports) to minimize attack surface and dependency tree.

## What Is Demonstrated
- Propagation of deadlines across execution layers.
- Randomized exponential backoff preventing lockstep retry storms.
- Fast failure via circuit breaker when downstream service fails repeatedly.
- Idempotency deduplication preventing duplicate execution on retry.

## What Is Not Demonstrated
- Distributed coordination across independent processes or network protocols (e.g. gRPC metadata headers).
