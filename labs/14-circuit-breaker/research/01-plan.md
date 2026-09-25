# Research Plan: Circuit Breaker Pattern & Cascading Failure Mitigation

## Research Topic
Circuit Breaker Pattern — Preventing a single degraded dependency from taking down the entire system.

## Objective
Investigate the operational mechanics, failure dynamics, and interplay between Circuit Breaker, Timeout, Retry with Exponential Backoff, Fallback, Bulkhead, and Load Shedding in distributed microservices.

## Research Questions
1. How do unconstrained timeouts lead to resource exhaustion (thread/connection starvation) and cascading failure?
2. What are the formal states and state transition triggers of a Circuit Breaker (CLOSED, OPEN, HALF_OPEN)?
3. How does Circuit Breaker differ fundamentally from Timeout, Retry, Bulkhead, and Load Shedding?
4. What failure modes exist when implementing Circuit Breakers (thundering herd on probe, inappropriate thresholds, metric skew)?
5. What are real-world architectural patterns for resilient integration in asynchronous vs synchronous payment & notification systems (CMMS, PPOB)?

## Search Strategy
1. Primary documentation: Microsoft Azure Architecture Center (Cloud Design Patterns), Martin Fowler's canonical Bliki on Circuit Breaker, AWS Builder's Library on Timeouts/Retries.
2. Verified secondary technical material: Netflix Hystrix architecture documentation, Release It! (Michael Nygard).
3. Cross-reference claims across sources regarding probe mechanics and failure counting.

## Expected Primary Sources
- Martin Fowler, *Circuit Breaker* (https://martinfowler.com/bliki/CircuitBreaker.html)
- Microsoft Azure Architecture Center, *Circuit Breaker Pattern* (https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker)
- AWS Builder's Library, *Timeouts, retries, and backoff with jitter* (https://aws.amazon.com/builders-library/timeouts-retries-and-backoff-with-jitter/)

## Risks / Unknowns
- Differences in state transition triggers: rolling window error-rate percentage vs consecutive failure count.
- Half-open probe concurrency bounds across distributed instances vs single-process in-memory breaker.
