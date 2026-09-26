# Research Plan

## Research Topic
Circuit Breaker Pattern — Cara Mencegah Satu Service Down Menjatuhkan Seluruh Sistem

## Objective
Produce structured research covering:
1. Circuit Breaker core concepts and state machine
2. Cascade failure mechanics
3. State transitions (CLOSED/OPEN/HALF-OPEN)
4. Timeout, retry, backoff, jitter interactions
5. Fallback and Bulkhead patterns
6. Observability metrics
7. Failure modes and anti-patterns
8. Case studies: CMMS WhatsApp and PPOB flows

## Research Questions
1. What are the canonical state transitions for Circuit Breaker?
2. What failure threshold and timeout configurations are recommended?
3. How does Circuit Breaker differ from Retry and Timeout?
4. What are proper fallback strategies (safe vs unsafe)?
5. How does Bulkhead isolate resources differently from Circuit Breaker?
6. What observability metrics are standard?
7. What are common failure modes (threshold too low, probe flood, error type confusion)?
8. How to apply to CMMS (Create Invoice → PDF → WhatsApp)?
9. How to apply to PPOB (User → Order → Provider Pulsa → Payment Gateway → WhatsApp)?

## Search Strategy
- Primary: Martin Fowler bliki (original pattern author)
- Primary: Microsoft Azure Architecture Center patterns
- Primary: Resilience4j reference implementation
- Secondary: Google SRE Book (cascading failures, load shedding)
- Secondary: AWS Builders Library (timeouts, retries, jitter)

## Expected Primary Sources
| Source | URL | Tier |
|--------|-----|------|
| Martin Fowler Circuit Breaker | https://martinfowler.com/bliki/CircuitBreaker.html | 1 |
| Azure Circuit Breaker Pattern | https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker | 1 |
| Azure Retry Pattern | https://learn.microsoft.com/en-us/azure/architecture/patterns/retry | 1 |
| Azure Bulkhead Pattern | https://learn.microsoft.com/en-us/azure/architecture/patterns/bulkhead | 1 |
| Resilience4j CircuitBreaker | https://resilience4j.readme.io/docs/circuitbreaker | 1 |
| Google SRE: Handling Overload | https://sre.google/sre-book/handling-overload/ | 1 |
| Google SRE: Cascading Failures | https://sre.google/sre-book/addressing-cascading-failures/ | 1 |

## Risks / Unknowns
- AWS Builders Library page may be JS-rendered (inaccessible)
- Need to distinguish generic pattern from library-specific behavior
- Must not invent threshold recommendations or benchmarks
- Financial fallback behavior (PPOB) must not make unsafe assumptions