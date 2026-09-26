# Core Concepts: Circuit Breaker Pattern

## Definition
A circuit breaker is a design pattern that prevents cascading failures in distributed systems by detecting failures and temporarily stopping requests to a failing service, allowing it time to recover.

**Source**: Martin Fowler (Source 1), Microsoft Azure (Source 3)

## Purpose
1. **Fail Fast**: Immediately reject requests when the downstream service is known to be failing
2. **Prevent Cascade**: Stop resource exhaustion (threads, connections, memory) in the caller
3. **Allow Recovery**: Give the failing service time to recover without load
4. **Graceful Degradation**: Enable fallback behavior when the circuit is open

**Source**: Microsoft Azure (Source 3), Google SRE (Source 5)

## The Three States

### CLOSED (Normal Operation)
- Requests flow through to the downstream service
- Failures are counted
- Successes reset or reduce failure count
- If failure threshold reached → transition to OPEN

### OPEN (Failing Fast)
- Requests are rejected immediately without calling downstream
- Returns a circuit-open error (fail fast)
- After a timeout (cooldown period) → transition to HALF_OPEN

### HALF_OPEN (Testing Recovery)
- Allows a limited number of probe requests through
- If probe succeeds → transition to CLOSED
- If probe fails → transition back to OPEN
- Prevents thundering herd on recovery

**Source**: Martin Fowler (Source 1), Netflix Hystrix (Source 2), Microsoft Azure (Source 3)

## Key Configuration Parameters (Illustrative Examples, Not Recommendations)

| Parameter | Purpose | Example Range |
|-----------|---------|---------------|
| FailureThreshold | Number of failures before opening | 5-20 |
| OpenTimeout | Time in OPEN state before HALF_OPEN | 10-60 seconds |
| HalfOpenMaxCalls | Probe requests allowed in HALF_OPEN | 1-10 |
| SuccessThreshold | Successes needed to close from HALF_OPEN | 1-5 |

**Source**: Netflix Hystrix (Source 2), gobreaker (Source 7), cep21/circuit (Source 8)

## Failure Detection
- **Timeout**: Request exceeds configured duration
- **Error Response**: HTTP 5xx, network errors, connection refused
- **Exception**: Panic, unhandled errors in client code

**Source**: Microsoft Azure (Source 3), Netflix Hystrix (Source 2)

## Thread Safety
Circuit breakers must be safe for concurrent access:
- **Mutex**: Protect state transitions and counters (sony/gobreaker uses mutex)
- **Atomic Operations**: For simple counters (cep21 uses atomic for some operations)
- **Channel-based**: State machine via channels (alternative Go pattern)

**Source**: gobreaker source (Source 7), cep21/circuit source (Source 8)

## What Circuit Breaker Does NOT Do
- Does not make a failing service healthy
- Does not retry requests (separate concern)
- Does not replace timeouts (complementary)
- Does not guarantee zero failed requests during transition

**Source**: Martin Fowler (Source 1), Google SRE (Source 5)