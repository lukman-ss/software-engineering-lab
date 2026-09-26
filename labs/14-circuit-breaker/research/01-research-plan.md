# Research Plan

## Research Topic
Circuit Breaker Pattern — Preventing Cascading Failures in Distributed Systems

## Objective
Research the Circuit Breaker pattern comprehensively to build a runnable Go demonstration lab that shows:
1. How cascading failures occur without a circuit breaker
2. How a circuit breaker prevents cascade failures
3. State transitions (CLOSED → OPEN → HALF_OPEN → CLOSED)
4. Recovery and failed recovery scenarios
5. Related patterns: timeout, retry, fallback, bulkhead, load shedding
6. Observability metrics
7. Case studies: CMMS WhatsApp and PPOB

## Research Questions

### Core Circuit Breaker
1. What are the exact state definitions and transitions for CLOSED, OPEN, HALF_OPEN?
2. What are standard failure threshold configurations?
3. How does the HALF_OPEN probe mechanism work?
4. What are the thread-safety considerations?

### Cascade Failure
1. How does a slow/down dependency cause resource exhaustion in callers?
2. What are the typical failure modes (timeout, connection pool exhaustion, thread starvation)?
3. How does fail-fast prevent cascade?

### Timeout vs Retry vs Circuit Breaker
1. What is the relationship between HTTP timeout and circuit breaker?
2. When to use retry with exponential backoff vs circuit breaker?
3. How do they interact (retry amplification)?

### Fallback Patterns
1. What are safe fallback strategies (cached response, degraded response, queued)?
2. What are unsafe fallbacks to avoid?

### Bulkhead vs Circuit Breaker
1. How does bulkhead isolation differ from circuit breaker?
2. When to use each?

### Observability
1. What metrics are essential for circuit breaker monitoring?
2. Standard metric names and meanings?

### Case Studies
1. CMMS: Invoice → PDF → WhatsApp — how to isolate WhatsApp failures
2. PPOB: User → Order → Provider → Payment Gateway → WhatsApp — critical vs async dependencies

## Search Strategy
1. Primary sources: Go standard library patterns, Netflix Hystrix docs, AWS/Azure/GCP resilience docs
2. Standards: RFCs, ISO standards for distributed systems
3. Reputable engineering blogs: Martin Fowler, Netflix Tech Blog, Google SRE books
4. Go-specific: gobreaker, sony/gobreaker, go-resilience libraries

## Expected Primary Sources
- Netflix Hystrix documentation (original circuit breaker implementation)
- Martin Fowler's Circuit Breaker pattern article
- Microsoft Azure Circuit Breaker pattern docs
- AWS Well-Architected Framework reliability pillar
- Google SRE Handbook (circuit breaker, timeout, retry)
- Go resilience libraries source code (sony/gobreaker, cep21/circuit)

## Risks / Unknowns
- Exact HALF_OPEN probe semantics vary across implementations
- Thread-safety mechanisms: mutex vs atomic vs channels
- Production threshold recommendations (not invented)
- Interaction between timeout, retry, and circuit breaker in practice