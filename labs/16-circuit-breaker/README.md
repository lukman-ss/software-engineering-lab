# Lab 16: Circuit Breaker Pattern

A circuit breaker prevents cascading failures by failing fast when a downstream service is struggling.

## Problem

When a downstream service becomes slow or unresponsive:
- Clients wait for timeouts, consuming resources
- Thread pools get exhausted
- Queues fill up
- The entire application becomes unresponsive

## Solution

The Circuit Breaker sits between the caller and the failing service.

### States (Prompt 049)

| State | Behavior |
|-------|----------|
| **CLOSED** | Requests pass through. Failures increment counter. |
| **OPEN** | Requests fail immediately with `ErrCircuitOpen`. |
| **HALF-OPEN** | Allows limited "probe" requests to test recovery. |

### State Transitions

```
CLOSED --(failures >= threshold)--> OPEN
  ^                                     |
  |   (reset timeout elapsed)           |
  +-------- HALF-OPEN --> CLOSED (success)
  |                    \
  +---------------------> OPEN (failure)
```

## Half-Open Probe Limiting (Prompt 051)

When cooldown completes, we **don't allow all traffic through at once**.

A single request could falsely succeed while the service is still unstable, or cause a thundering herd.

**Default**: `MaxHalfOpenProbes = 1` — only one concurrent probe allowed.

Configuration:
```go
cb := breaker.New(breaker.Config{
    FailureThreshold:  5,
    ResetTimeout:      30 * time.Second,
    MaxHalfOpenProbes: 2, // Allow up to 2 concurrent probes
})
```

## Thread Safety (Prompt 050)

The implementation uses `sync.RWMutex` and is safe under concurrent access.

### Concurrent Test

```bash
go test -race -run TestCircuitBreakerConcurrentSafety
```

## Usage

```go
cb := breaker.New(breaker.Config{
    FailureThreshold:  5,
    ResetTimeout:      60 * time.Second,
    MaxHalfOpenProbes: 1,
})

// Use with manual tracking:
err := cb.Execute(func() error {
    return httpClient.Get("http://downstream/api")
})
if errors.Is(err, breaker.ErrCircuitOpen) {
    // Service is unavailable, return cached data or error
}
```

## Trade-offs

| Decision | Trade-off |
|----------|-----------|
| Failure Threshold | Too low: flimsy services cause false opens. Too high: long outages. |
| Reset Timeout | Too short: cause churn. Too long: slow recovery. |
| Max Half-Open Probes | Set higher for faster recovery, lower for safety. |

## What Still Can Fail

1. **Network partitions**: Service appears up but is unreachable
2. **Dependency timeouts**: Circuit opens but root cause is unknown
3. **Partial outages**: Service responds slowly (not failing fast)
4. **Cascading failures**: Multiple circuits open interdependent operations

---

## Navigasi

- **Previous**: [Lab 08 — Database Isolation Level](../08-database-isolation-level/)
- **Next**: [Lab 10 — Rate Limiting](../10-rate-limiting/)