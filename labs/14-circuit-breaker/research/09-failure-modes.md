# Failure Modes

## Anti-Patterns and Failure Modes in Circuit Breaker Design

### 1. Threshold Too Low (Over-sensitive Breaker)

**Claim**: Setting failure threshold too low causes the breaker to trip on momentary blips, dropping valid traffic.

**Evidence**:
- Included in README.md Failure Modes section (spec): "Threshold too low: Breaker trips on momentary blips, dropping valid traffic."
- Corroborated by Resilience4j: `minimumNumberOfCalls` parameter exists specifically to prevent premature tripping — "If only 9 calls have been evaluated the CircuitBreaker will not transition to open even if all 9 calls have failed." — https://resilience4j.readme.io/docs/circuitbreaker

**Mitigation**: Use a minimum call count (minimumNumberOfCalls) before evaluating failure rate. Ensure threshold aligns with expected transient failure rate of the dependency.

**Confidence**: HIGH — Resilience4j explicitly designed its API to prevent this anti-pattern.

---

### 2. Probe Flood in HALF_OPEN

**Claim**: Too many parallel probe calls during HALF_OPEN can overwhelm a recovering service.

**Evidence**:
- Included in README.md Failure Modes section (spec): "Probe flood: Too many parallel probe calls during HALF_OPEN overwhelm recovering service."
- Azure: "The Half-Open state helps prevent a recovering service from suddenly being flooded with requests. As a service recovers, it might be able to support a limited volume of requests until the recovery is complete. But while recovery is in progress, a flood of work can cause the service to time out or fail again." — https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker
- Resilience4j: `permittedNumberOfCallsInHalfOpenState` limits concurrent probes (default: 10).

**Mitigation**: Strictly cap `permittedNumberOfCallsInHalfOpenState`. Do not allow unlimited concurrent probes.

**Confidence**: HIGH — Azure and Resilience4j both explicitly address this.

---

### 3. Ignoring Error Types (4xx vs 5xx)

**Claim**: Counting 4xx client errors (e.g., 400 Bad Request) as circuit-breaking faults, when only 5xx server errors or timeouts should trip the circuit.

**Evidence**:
- Included in README.md Failure Modes section (spec): "Ignoring error types: Counting 4xx (client errors) as circuit-breaking faults rather than restricting to 5xx/timeouts."
- Resilience4j: `recordExceptions` and `ignoreExceptions` allow selective exception classification — "By default all exceptions count as a failure. You can define a list of exceptions which should count as a failure." — https://resilience4j.readme.io/docs/circuitbreaker
- Azure: "The reasons for a request failure can vary in severity... A circuit breaker might be able to examine the types of exceptions that occur and adjust its strategy based on the nature of these exceptions." — https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker

**Mitigation**: Classify exceptions explicitly: count 5xx, timeouts, and connection errors as circuit-breaking failures; ignore 4xx client errors as application-level issues.

**Confidence**: HIGH — Azure and Resilience4j both address error-type classification.

---

### 4. Inappropriate Timeouts on External Services

**Claim**: A circuit breaker configured with too-long timeouts still blocks threads, defeating the circuit breaker's purpose.

**Evidence**:
- Azure: "A circuit breaker might not fully protect applications from failures in external services that have long time-out periods. If the time-out is too long, a thread that runs a circuit breaker might be blocked for an extended period before the circuit breaker indicates that the operation failed. During this time, many other application instances might also try to invoke the service through the circuit breaker and tie up numerous threads before they all fail." — https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker

**Mitigation**: Pair explicit, tightly bounded timeouts with circuit breaker for full protection.

**Confidence**: HIGH — Azure explicitly documents this failure mode.

---

### 5. Retry Storms (Amplification)

**Claim**: Retry without a circuit breaker amplifies load on a failing dependency, accelerating cascade failure.

**Evidence**:
- Google SRE: retry budgets (max 3 attempts), per-client retry budget (10% ratio cap) — https://sre.google/sre-book/handling-overload/
- Azure: "An aggressive retry policy with minimal delay between attempts, and a large number of retries, could further degrade a busy service that's running close to or at capacity." — https://learn.microsoft.com/en-us/azure/architecture/patterns/retry
- Azure Circuit Breaker: "If a request still fails after a significant number of retries, it's better for the application to prevent further requests going to the same resource and report a failure immediately." — https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker

**Mitigation**: Combine bounded retry with circuit breaker. When circuit is OPEN, retry logic must stop.

**Confidence**: HIGH.

---

### 6. Circuit Breaker Fluctuation (Oscillation)

**Claim**: If HALF_OPEN transitions back to OPEN too quickly (and then repeats), the breaker oscillates between OPEN and HALF_OPEN without stabilizing.

**Evidence**:
- Azure: "a circuit breaker can fluctuate and reduce the response times of applications if it switches from the Open state to the Half-Open state too quickly." — https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker

**Mitigation**: Use cooldown that increases progressively, or require consecutive successful probes before CLOSED.

**Confidence**: HIGH — Azure explicitly documents this.

---

## Summary Table

| # | Failure Mode | Severity | Mitigation |
|---|-------------|----------|-----------|
| 1 | Threshold too low | Medium | minimumNumberOfCalls guard |
| 2 | Probe flood | Medium | permittedNumberOfCallsInHalfOpenState cap |
| 3 | Error type confusion | High | Classify 4xx vs 5xx explicitly |
| 4 | Inappropriate timeout | High | Tight timeout paired with circuit breaker |
| 5 | Retry storm | High | Bounded retry + circuit breaker combination |
| 6 | Oscillation | Low | Progressive cooldown / consecutive success probes |

## NOT VERIFIED
- Exact threshold values that trigger oscillation (depends on traffic volume)
- Whether oscillation occurs in all implementations or only specific algorithms