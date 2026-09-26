# Engineering Design

Target Lab: labs/14-circuit-breaker
Research Status: APPROVED

## Concept To Prove
The Circuit Breaker pattern isolates failing or slow dependencies to prevent cascading resource exhaustion and provides a path to automated recovery.

## Expected Behavior
- **CLOSED**: All requests route to the downstream service. Failures increment failure count. Successes reset failure count. Trips to OPEN when failures reach threshold.
- **OPEN**: All requests immediately fail fast with an error (`ErrCircuitOpen`) without sending downstream network traffic. Enters cooldown period.
- **HALF-OPEN**: Dispatches limited probe requests to test dependency recovery. Success transitions back to CLOSED; failure resets to OPEN.

## Failure Scenario
Downstream service slows down or goes completely down. In the absence of a circuit breaker, caller threads block on timeouts, exhausting thread and connection pools.

## Success Criteria
1. Transitions: CLOSED -> OPEN -> HALF_OPEN -> CLOSED (recovery).
2. Transitions: CLOSED -> OPEN -> HALF_OPEN -> OPEN (repeated failure).
3. Zero downstream requests executed when circuit is OPEN.
4. Concurrency safety verified with Go race detector.

## Architecture
- `Checkout Service`: Application entry point sending payments.
- `Circuit Breaker`: State machine proxy guarding downstream calls.
- `Fake Payment Server`: Controllable HTTP downstream server.

## Components
- `internal/circuitbreaker`: Core state machine implementation.
- `internal/payment`: Fake HTTP server and payment client.
- `internal/checkout`: Checkout business logic consuming payment client.
- `cmd/demo`: Executable demonstration script.

## Test Strategy
- Unit tests covering each state transition.
- Fake server integration tests.
- High-concurrency race condition testing.

## Execution Plan
1. Implement `circuitbreaker.CircuitBreaker`.
2. Implement `payment.FakeServer` and `payment.Client`.
3. Implement `checkout.Service`.
4. Implement runnable demo in `cmd/demo/main.go`.
5. Execute unit, race, and demo tests.

## Implementation Decisions
- Synchronized with `sync.Mutex` rather than complex lock-free data structures for simplicity and readability.
- Injected `now` function for fast, deterministic unit testing without sleeps.
