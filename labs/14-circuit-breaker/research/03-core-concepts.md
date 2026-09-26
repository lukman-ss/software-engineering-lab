# Core Concepts

## Circuit Breaker

### Definition
A Circuit Breaker is a proxy that wraps calls to a downstream dependency. It monitors failure events (timeouts, HTTP 5xx, connection errors) and, when failures reach a configurable threshold, transitions to an "open" state that rejects all calls immediately without forwarding them to the downstream service. After a cooldown period, the breaker enters a "half-open" state that allows a limited number of probe calls to test whether the dependency has recovered.

### Purpose
1. **Prevent resource exhaustion**: When a dependency is slow, blocked caller threads hold sockets, memory, and connection-pool slots. The circuit breaker fails fast so these resources are released immediately.
2. **Allow dependency recovery**: By cutting off all traffic while OPEN, the downstream service gets idle time to resolve its overload.
3. **Prevent cascade failure**: A failure in one service does not propagate resource starvation to upstream services.

### State Machine
```
CLOSED → (failures ≥ threshold) → OPEN
OPEN   → (cooldown elapsed)   → HALF_OPEN
HALF_OPEN → (probe succeeds) → CLOSED
HALF_OPEN → (probe fails)    → OPEN
```

### Core Configuration Parameters
| Parameter | Description | Source |
|-----------|-------------|--------|
| `failureThreshold` | Number of consecutive failures required to trip the circuit | Martin Fowler §CircuitBreaker; Azure Circuit Breaker; Resilience4j |
| `openTimeout` / `waitDurationInOpenState` | Duration the breaker stays OPEN before transitioning to HALF_OPEN | Same |
| `halfOpenMaxCalls` / `permittedNumberOfCallsInHalfOpenState` | Maximum number of probe calls allowed in HALF_OPEN | Resilience4j |

### Verified Facts
- **CLOSED**: Requests are forwarded to the dependency. Each failure increments a counter. A success resets the counter to zero. When the failure count reaches the threshold, the breaker transitions to OPEN.
  - Martin Fowler: "Should we get a timeout, we increment the failure counter, successful calls reset it back to zero."
  - Azure: "CLOSED: The request from the application is routed to the operation. The proxy maintains a count of the number of recent failures."
- **OPEN**: All calls fail immediately with an exception. No network calls are made to the downstream dependency.
  - Martin Fowler: "raise CircuitBreaker::Open" when `when :open`
  - Azure: "OPEN: The request from the application fails immediately and an exception is returned."
- **HALF_OPEN**: A limited number of probe calls are allowed through. If all permitted calls succeed, the breaker transitions to CLOSED. If any fails, it returns to OPEN.
  - Martin Fowler: "there is now a third state present - half open - meaning the circuit is ready to make a real call as trial to see if the problem is fixed."
  - Azure: "HALF-OPEN: A limited number of requests from the application are allowed to pass through... If these requests are successful, the circuit breaker switches to the Closed state. If any request fails, the circuit breaker reverts to the Open state."

### Monitoring Recommendations (from Martin Fowler)
- Log every state transition
- Expose circuit state for monitoring/diagnostics
- Alert operations staff when the breaker trips
- Allow operations staff to manually trip or reset breakers

### Implementation Note (from Resilience4j)
- CircuitBreaker is thread-safe (AtomicReference for state, synchronized for sliding window)
- CircuitBreaker does NOT synchronize the function call itself — the protected call is not part of the critical section
- CircuitBreaker is distinct from Bulkhead — for concurrency limiting, use Bulkhead separately