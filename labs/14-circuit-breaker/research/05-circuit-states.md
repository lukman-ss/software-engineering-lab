# 05 Circuit States

```
       ┌───────────┐
       │  CLOSED   │◄─────────────────────────┐
       └─────┬─────┘                          │
             │ failures >= FailureThreshold   │ probe success
             ▼                                │
       ┌───────────┐                          │
       │   OPEN    │                          │
       └─────┬─────┘                          │
             │ cooldown elapsed               │
             ▼                                │
       ┌───────────┐                          │
       │ HALF_OPEN ├──────────────────────────┘
       └─────┬─────┘
             │ probe fail
             ▼
          (to OPEN)
```

## CLOSED — Normal Operation
- Every request dispatched downstream.
- Success → reset failure count (consecutive-count model) or decrement rolling window (Hystrix model).
- Failure → increment count; if >= threshold → OPEN, record openTimestamp.

**Sources**: Martin Fowler (2014-03-06) — "Calling the circuit breaker will call the underlying block if the circuit is closed... I determine the state of the breaker comparing the failure count to the threshold"; Hystrix Wiki — request-volume + error-percentage tripping; Azure Circuit Breaker — "maintains a count of the number of recent failures... If the number of recent failures exceeds a specified threshold within a given time period, the proxy is placed into the Open state".

## OPEN — Fail Fast
- Reject immediately with sentinel error, e.g. `ErrCircuitOpen`.
- Downstream function NOT invoked (verified by downstream call counter in tests).
- After `OpenTimeout` since openTimestamp elapsed → allow transition to HALF_OPEN on next call.
- Optional: manual override to force open/close (operator control).

**Sources**: Azure — "Open: The request from the application fails immediately and an exception is returned"; Fowler — "when :open then raise CircuitBreaker::Open"; Hystrix — "While it is open, it short-circuits all requests made against that circuit-breaker."

## HALF_OPEN — Probe
- Allow limited number of probe calls (`HalfOpenMaxCalls`, typically 1) through.
- Success → CLOSED, reset counters.
- Any failure → OPEN, restart cooldown timer.
- Purpose: prevent flood onto recovering service.

**Sources**: Azure — "Half-Open: A limited number of requests... If these requests are successful, the circuit breaker assumes that the fault... is fixed, and the circuit breaker switches to the Closed state... If any request fails, the circuit breaker assumes that the fault is still present, so it reverts to the Open state"; Fowler — half-open trial call "will either reset the breaker if successful or restart the timeout if not"; Hystrix — single let-through request after sleep window.

## Threshold Models (Disagreement — see 04-contradictions.md)
- Consecutive-count (Fowler basic impl, gobreaker default 5): simple, deterministic, lab choice.
- Rolling-window error % (Hystrix: volume threshold + error %; Azure: failures within time period; gobreaker Interval): robust under fluctuating traffic.
- NOT VERIFIED as superior universally; depends on traffic pattern.

## Concurrency Requirements
- State + counters guarded by mutex (gobreaker uses sync.Mutex; this lab same).
- Azure: "implementation shouldn't block concurrent requests or add excessive overhead".
- Verify with `go test -race ./...`.
