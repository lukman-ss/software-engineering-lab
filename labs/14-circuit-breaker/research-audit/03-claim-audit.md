# Claim Audit

## Claim 1

Claim:
The canonical Circuit Breaker state machine has three states (CLOSED, OPEN, HALF_OPEN) with four transitions: CLOSED to OPEN on failure threshold, OPEN to HALF_OPEN after cooldown, HALF_OPEN to CLOSED on success, and HALF_OPEN to OPEN on failure.

Location:
`research/10-final-research.md` § Finding 1: State Machine

Evidence Provided:
Direct alignment across Martin Fowler (2014), Microsoft Azure Architecture Center (2025), and Resilience4j reference documentation.

Source:
Source 1, Source 2, Source 5

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Universally accepted state machine specification across all Tier 1 sources.

---

## Claim 2

Claim:
When OPEN, the breaker rejects calls immediately (fail-fast) without invoking the downstream service.

Location:
`research/10-final-research.md` § Finding 2: Fail-Fast Behavior

Evidence Provided:
Martin Fowler specifies returning an error directly; Azure notes requests fail immediately; Resilience4j specifies throwing `CallNotPermittedException`.

Source:
Source 1, Source 2, Source 5

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Core mechanism of the pattern.

---

## Claim 3

Claim:
Successful calls reset the failure counter to zero in CLOSED state.

Location:
`research/10-final-research.md` § Finding 3: Failure Counter Reset

Evidence Provided:
Martin Fowler and Azure state that consecutive failures increment counter and success resets it back to zero.

Source:
Source 1, Source 2

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Applies to consecutive-failure counting models. (Sliding windows use rolling buffers, which research appropriately notes in areas of disagreement).

---

## Claim 4

Claim:
HALF_OPEN allows only a limited number of concurrent probe calls to prevent overwhelming a recovering service.

Location:
`research/10-final-research.md` § Finding 4: Limited Probe Calls

Evidence Provided:
Azure states a limited number of requests are permitted; Resilience4j implements `permittedNumberOfCallsInHalfOpenState`.

Source:
Source 2, Source 5

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Standard behavior to protect downstream during recovery.

---

## Claim 5

Claim:
Outbound HTTP calls must have explicit timeouts; Circuit Breaker adds fail-fast on top of timeouts and does not replace them.

Location:
`research/10-final-research.md` § Finding 5: Timeout Necessity

Evidence Provided:
Go standard library documentation verifies zero-value timeout hangs indefinitely; Azure documentation confirms circuit breakers without timeouts leave caller threads blocked on slow dependencies.

Source:
Source 2, Source 10

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Critical architectural distinction between bounding single requests and isolating degraded dependencies.

---

## Claim 6

Claim:
Retry should be bounded (retry budget or max attempts) and must respect circuit breaker exceptions by stopping retries when OPEN.

Location:
`research/10-final-research.md` § Finding 6: Retry Must Be Bounded

Evidence Provided:
Google SRE Book defines retry budgets; Azure Retry and Circuit Breaker patterns dictate ceasing retries when the breaker reports downstream is unavailable.

Source:
Source 2, Source 3, Source 6

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Prevents retry storms.

---

## Claim 7

Claim:
Bulkhead isolates resources per dependency, whereas Circuit Breaker gates traffic; they are distinct and complementary patterns.

Location:
`research/10-final-research.md` § Finding 7: Bulkhead Distinctness

Evidence Provided:
Azure Bulkhead pattern and Resilience4j explicitly contrast thread/resource pool isolation against call-gating.

Source:
Source 4, Source 5

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Essential clarification preventing pattern confusion.

---

## Claim 8

Claim:
Every Circuit Breaker state transition should be logged, exposed for monitoring, and trigger alerts.

Location:
`research/10-final-research.md` § Finding 8: Observability

Evidence Provided:
Martin Fowler and Azure explicitly recommend emitting events and logging state transitions for monitoring and alerting.

Source:
Source 1, Source 2

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Well-supported best practice.

---

## Claim 9

Claim:
Circuit breakers should differentiate between client-side errors (e.g., HTTP 4xx) and server/infrastructure failures (e.g., HTTP 5xx, timeouts), tripping only on upstream faults.

Location:
`research/09-failure-modes.md` § Failure Mode 2: Tripping on 4xx Errors

Evidence Provided:
Resilience4j `recordExceptions` / `ignoreExceptions` documentation; Microsoft Azure pattern guidance on filtering errors.

Source:
Source 2, Source 5

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
MEDIUM

Notes:
Tripping on 4xx would allow bad client requests to cause denial-of-service for legitimate callers. Supported by reference implementations.

---

## Claim 10

Claim:
In-memory circuit breakers only track state per process/instance, meaning total probe traffic across N instances scales linearly with instance count.

Location:
`research/05-circuit-states.md` § Concurrency & Multi-Instance Considerations

Evidence Provided:
Architectural reasoning and Azure pattern concurrency section noting local proxy vs shared state tradeoffs.

Source:
Source 2

Source Actually Supports Claim:
YES

Classification:
IMPLEMENTATION-SPECIFIC

Severity:
LOW

Notes:
Properly classified as an implementation characteristic of process-local breakers rather than a universal requirement.
