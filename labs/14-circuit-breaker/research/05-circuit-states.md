# Circuit States

## CLOSED State

**Claim**: In CLOSED state, all requests are forwarded to the downstream dependency. Successes reset the failure counter; failures increment it. When failures ≥ threshold, state transitions to OPEN.

**Evidence**:
- Martin Fowler: "Should we get a timeout, we increment the failure counter, successful calls reset it back to zero." — https://martinfowler.com/bliki/CircuitBreaker.html
- Azure: "CLOSED: The request from the application is routed to the operation. The proxy maintains a count of the number of recent failures." — https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker
- Resilience4j: "The state of the CircuitBreaker changes from CLOSED to OPEN when the failure rate is equal or greater than a configurable threshold." — https://resilience4j.readme.io/docs/circuitbreaker

**Confidence**: HIGH — all three Tier 1 sources agree.

**Key detail — failure window/reset behavior**:
- Azure specifies: "The failure counter for the Closed state is time based. It automatically resets at periodic intervals."
- Resilience4j supports both count-based and time-based sliding windows.
- This lab: implements a simple consecutive-failure counter (reset on success). Documented as an educational simplification.

---

## OPEN State

**Claim**: When OPEN, all incoming calls fail immediately without executing the downstream call. A cooldown timer (waitDurationInOpenState) begins. When the timer expires, state transitions to HALF_OPEN.

**Evidence**:
- Martin Fowler: "raise CircuitBreaker::Open" — "calling the circuit breaker will call the underlying block if the circuit is closed, but return an error if it's open" — https://martinfowler.com/bliki/CircuitBreaker.html
- Azure: "OPEN: The request from the application fails immediately and an exception is returned to the application." — https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker
- Resilience4j: "The CircuitBreaker rejects calls with a CallNotPermittedException when it is OPEN. After a wait time duration has elapsed, the CircuitBreaker state changes from OPEN to HALF_OPEN" — https://resilience4j.readme.io/docs/circuitbreaker

**Confidence**: HIGH — unanimous across Tier 1 sources.

**Key detail — what triggers HALF_OPEN transition**:
- Martin Fowler (simpler model): state is computed on each call — if enough time has passed since last failure → HALF_OPEN.
- Azure: a timeout timer starts when entering OPEN; when it expires → HALF_OPEN.
- Resilience4j: `automaticTransitionFromOpenToHalfOpenEnabled` controls whether a background thread triggers the transition, or only the next incoming call triggers it.

---

## HALF-OPEN State

**Claim**: In HALF_OPEN, a limited number of probe calls are permitted. If all probe calls succeed → CLOSED. If any probe call fails → OPEN (restarts cooldown).

**Evidence**:
- Martin Fowler: "there is now a third state present - half open - meaning the circuit is ready to make a real call as trial to see if the problem is fixed." — https://martinfowler.com/bliki/CircuitBreaker.html
- Azure: "HALF-OPEN: A limited number of requests from the application are allowed to pass through and invoke the operation. If these requests are successful, the circuit breaker assumes that the fault that caused the failure is fixed, and the circuit breaker switches to the Closed state. The failure counter is reset. If any request fails, the circuit breaker assumes that the fault is still present, so it reverts to the Open state." — https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker
- Resilience4j: "permits a configurable number of calls to see if the backend is still unavailable or has become available again. Further calls are rejected... If the failure rate... is then equal or greater than the configured threshold, the state changes back to OPEN." — https://resilience4j.readme.io/docs/circuitbreaker

**Confidence**: HIGH — unanimous across Tier 1 sources.

**Key detail — preventing probe flood**:
- Azure: "The Half-Open state helps prevent a recovering service from suddenly being flooded with requests. As a service recovers, it might be able to support a limited volume of requests until the recovery is complete."
- Resilience4j: `permittedNumberOfCallsInHalfOpenState` (default: 10) caps the number of simultaneous probe calls.

---

## State Transition Summary

| From | To | Trigger | Source |
|------|----|---------|--------|
| CLOSED → OPEN | failures ≥ threshold | Martin Fowler, Azure, Resilience4j | HIGH |
| OPEN → HALF_OPEN | cooldown timer elapsed | Martin Fowler, Azure, Resilience4j | HIGH |
| HALF_OPEN → CLOSED | all probe calls succeed | Martin Fowler, Azure, Resilience4j | HIGH |
| HALF_OPEN → OPEN | any probe call fails | Azure, Resilience4j (Martin Fowler implies but less explicit) | HIGH |

## NOT VERIFIED
- Whether there should be a time-based auto-transition from OPEN to HALF_OPEN vs call-triggered only. Resilience4j supports both (`automaticTransitionFromOpenToHalfOpenEnabled`); Martin Fowler's original implementation is call-triggered. This lab uses call-triggered for simplicity.
