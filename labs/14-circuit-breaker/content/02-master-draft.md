# Circuit Breaker Pattern Implementation in Go

## Problem
Remote calls across networks fail or hang. Without protection, callers block waiting for timeouts, holding threads, sockets, and memory until system resources deplete. A single downstream failure can cascade into a total platform outage.

## Why This Matters
Unprotected remote calls exhaust critical system resources (memory, threads, DB connections) leading to cascading failures. Circuit Breaker prevents caller self-destruction by failing fast, preserving system resources, and allowing downstream services time to recover.

## Mental Model
A circuit breaker acts as a state machine proxy that monitors downstream failures. When failures exceed a threshold, it trips OPEN to block all calls immediately (fail-fast). After a cooldown period, it enters HALF_OPEN to allow limited probe calls. Successful probes restore normal operation (CLOSED); failed probes return to OPEN.

## Core Concept
The Circuit Breaker pattern implements three states:
- **CLOSED**: Normal operation — requests flow to downstream
- **OPEN**: Failing fast — requests rejected immediately with zero network calls
- **HALF_OPEN**: Probe state — limited test calls to check recovery

## Implementation
### State Machine
The implementation uses a synchronized struct with mutex protection:

```go
type CircuitBreaker struct {
	mu                   sync.Mutex
	state                State
	failureCount         int
	consecutiveSuccesses int
	halfOpenCalls        int
	lastStateChange      time.Time
	config               Config
	now                  func() time.Time
}
```

### State Transitions
Transitions occur based on failure counts and timeouts:
- CLOSED → OPEN: When `failureCount >= FailureThreshold`
- OPEN → HALF_OPEN: When `now().Sub(lastStateChange) >= OpenTimeout`
- HALF_OPEN → CLOSED: After `HalfOpenMaxCalls` consecutive successes
- HALF_OPEN → OPEN: On any failure during HALF_OPEN

### Execution Flow
```go
func (cb *CircuitBreaker) Execute(fn func() error) error {
	cb.mu.Lock()
	cb.checkStateTransitionLocked()

	switch cb.state {
	case StateOpen:
		cb.mu.Unlock()
		return ErrCircuitOpen

	case StateHalfOpen:
		// Allow limited probe calls
		if cb.halfOpenCalls >= cb.config.HalfOpenMaxCalls {
			cb.mu.Unlock()
			return ErrCircuitOpen
		}
		cb.halfOpenCalls++
		cb.mu.Unlock()
		// Execute function and handle result
		// ...

	default: // StateClosed
		cb.mu.Unlock()
		// Execute function and handle failures
		// ...
	}
}
```

## Code Walkthrough
### Creating a Circuit Breaker
```go
cb := circuitbreaker.New(circuitbreaker.Config{
	FailureThreshold: 3,
	OpenTimeout:      300 * time.Millisecond,
	HalfOpenMaxCalls: 1,
})
```

### Using in Checkout Service
```go
svc := checkout.NewService(paymentClient, cb)
err := svc.Checkout(ctx)
// Returns ErrCircuitOpen when OPEN, otherwise payment errors or success
```

## What the Tests Prove
### Unit Tests
16 unit tests verify:
- Initial state is CLOSED
- Successful calls remain CLOSED
- Failures below threshold stay CLOSED
- Threshold reached triggers OPEN state
- OPEN calls fail fast without downstream execution
- Cooldown moves breaker to HALF_OPEN
- Successful HALF_OPEN probe closes circuit
- Failed HALF_OPEN probe re-opens circuit
- Recovery after dependency heals
- Concurrency safety under race detector
- Panic safety during HALF_OPEN and CLOSED states
- Default configuration values
- HalfOpenMaxCalls > 1 requires consecutive successes

### Integration Tests
Two integration tests verify:
- Downstream failures trip circuit OPEN
- Cooldown and recovery transition HALF_OPEN → CLOSED

### Demo Verification
The executable demo shows four scenarios:
1. **Without CB (Slow Dependency)**: 3 requests block for ~100ms each, all hit downstream
2. **With CB (Fail-Fast)**: After 3 failures, subsequent requests fail in nanoseconds with zero downstream calls
3. **Recovery**: After cooldown, probe succeeds → CLOSED → normal traffic resumes
4. **Failed Recovery**: Probe fails → returns to OPEN → fail-fast resumes

## Recovery / Rollback
Automatic recovery occurs after OpenTimeout elapses:
1. Breaker transitions from OPEN to HALF_OPEN
2. Limited probe calls allowed (HalfOpenMaxCalls)
3. Success: Transitions to CLOSED, resets counters
4. Failure: Returns to OPEN, restarts cooldown timer

## Production Considerations
- Timeouts used in lab (100ms HTTP, 300ms cooldown) are illustrative for testing
- Production values must be tuned to actual service SLA, P99 latency, and recovery profiles
- Implementation uses consecutive failure counting; production may require sliding window or error rate
- Consider externalizing circuit state for microservice fleets (Redis/Consul) vs per-instance
- Classify errors carefully — typically only 5xx and timeouts should trip breaker, not 4xx client errors

## Common Mistakes
- Setting FailureThreshold too low causing premature opening on transient blips
- Allowing too many concurrent probes during HALF_OPEN (thundering herd)
- Counting 4xx errors as circuit-breaking faults instead of restricting to 5xx/timeouts
- Ignoring that circuit breaker does not heal dependencies — it only prevents caller exhaustion

## Case Study
**CMMS Example**: `Create Invoice` → `Generate PDF` → `Send WhatsApp`
- Design rule: WhatsApp down ≠ Invoice creation down
- Pattern:
  1. Save invoice synchronously in primary DB
  2. Push notification job to background queue (RabbitMQ/Kafka)
  3. Worker calls WhatsApp API wrapped in Circuit Breaker + exponential backoff retry
  4. Breaker trips OPEN if WhatsApp API fails; queue retains message; retries pause

## Key Takeaways
- Circuit Breaker does not heal a broken dependency
- Circuit Breaker stops caller self-destruction by failing fast and preserving system threads/memory
- Downstream service gets idle recovery time without request flooding
- Three-state machine prevents cascading failures while enabling automated recovery
- Thread-safe implementation verified under race detector with zero data races
- Fail-fast behavior executes in nanoseconds/microseconds vs millisecond timeouts
- Probe mechanism prevents overwhelming recovering services during HALF_OPEN
- Default configuration provides sensible out-of-the-box behavior
- Implementation demonstrates core pattern without external dependencies

## Sources
- Research: labs/14-circuit-breaker/research/
  - Core concepts: research/03-core-concepts.md
  - Circuit states: research/05-circuit-states.md
  - Timeout/retry/backoff: research/06-timeout-retry-backoff.md
  - Failure modes: research/09-failure-modes.md
  - Final synthesis: research/10-final-research.md
  - Evidence: research/03-evidence.md
  - Sources: research/02-sources.md
- Implementation: labs/14-circuit-breaker/internal/
  - Circuit breaker: internal/circuitbreaker/circuit_breaker.go
  - Tests: internal/circuitbreaker/circuit_breaker_test.go
  - Payment client/server: internal/payment/
  - Checkout service: internal/checkout/service.go
- Demo: cmd/demo/main.go
- Integration tests: tests/integration_test.go