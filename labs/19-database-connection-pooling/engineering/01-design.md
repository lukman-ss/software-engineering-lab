# Engineering Design

Target Lab: labs/19-database-connection-pooling
Research Status: APPROVED

## Concept To Prove
- Direct connection creation penalty (overhead).
- Connection exhaustion and throughput degradation when pools are oversized (beyond hardware limits).
- Connection leaks when connections are held during external I/O.

## Expected Behavior
- Small, saturated pools limit concurrency but maximize throughput and minimize latency.
- Oversized pools lead to connection exhaustion or increased latency.
- External network calls inside transactions hold connections, starving the pool.

## Failure Scenario
- An oversized pool attempting to process more concurrent transactions than hardware limits.
- Code making external I/O without releasing the DB connection, causing timeouts for other borrowers.

## Success Criteria
- Test validates pool sizes and limits.
- Test verifies connection leakage triggers errors.
- Test validates throughput degradation or max connection enforcement.

## Architecture
- Custom `database/sql/driver` to simulate database server limits, connection overhead, and active query counting.
- Service layer handling queries.

## Components
- `internal/pool/mockdb.go`: The mock driver simulating network overhead and backend limits.
- `internal/pool/service.go`: Business logic demonstrating correct and incorrect connection handling.
- `tests/pool_test.go`: Test suite validating concepts.
- `cmd/demo/main.go`: Demo runner.

## Test Strategy
- Happy path: Fixed pool size handling concurrent requests efficiently.
- Failure path: Unbounded/oversized pool exceeding server connection limits.
- Edge case: Connection leak via blocked external IO.

## Execution Plan
1. Implement driver and service.
2. Write tests.
3. Run tests with race detector.
4. Run demo.

## Implementation Decisions
- A mock `database/sql/driver` is used instead of a real PostgreSQL server to eliminate external dependency overhead and deterministically simulate connection latency, `max_connections`, and backend resource saturation.
