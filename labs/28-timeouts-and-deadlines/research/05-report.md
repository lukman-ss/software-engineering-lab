# Research Report: Timeouts — Jangan Biarkan Satu Dependency Lambat Menahan Seluruh Sistem

## Research Question

How should timeouts be implemented in distributed systems to prevent cascading failures while maintaining correctness, considering the interaction between timeouts, retries, circuit breakers, and idempotency requirements?

## Executive Summary

Research from authoritative sources (Google SRE Book, Microsoft Azure Architecture Center, gRPC documentation, PostgreSQL documentation, AWS Architecture Blog, RabbitMQ documentation) establishes that timeouts are a critical reliability mechanism for distributed systems. The core finding is that timeouts must be part of an integrated reliability strategy combining timeout budgeting, exponential backoff with jitter, circuit breakers, and idempotency — no single mechanism is sufficient alone.

Key evidence shows:
1. Single slow dependencies cause cascading failures when resources exhaust (HIGH confidence, 2 sources)
2. Retries amplify load by 4x or more during partial failures (HIGH confidence, 2 sources)
3. Deadline propagation prevents wasted work when users abandon requests (HIGH confidence, 2 sources)
4. Database and queue timeouts are required beyond network API timeouts (HIGH confidence, 3+ sources)
5. Idempotency is required for safe retries of non-idempotent operations (HIGH confidence, 3 sources)

## Findings

### Finding 1: Timeouts Must Be Based on Latency Distributions, Not Arbitrary Values

**Claim:** Timeout values should be determined from production latency distributions (P50, P95, P99) plus safety margins, not arbitrary "safe" values.

**Evidence:** Google SRE Book shows that setting deadlines orders of magnitude larger than mean latency causes thread exhaustion: "With a 100-second deadline, 5% of requests would consume 5,000 threads (50 QPS * 100 seconds), but the frontend doesn't have that many threads available." This leads to cascading failure. All sources agree timeouts should be P99 + margin.

**Sources:** Google SRE Book (Addressing Cascading Failures, Monitoring Distributed Systems)

**Confidence:** HIGH

### Finding 2: Deadline Propagation Prevents Resource Waste

**Claim:** Deadlines should be propagated through service chains as absolute deadlines (not fixed timeouts at each layer) to prevent downstream services from doing work after the original request deadline has passed.

**Evidence:** gRPC documentation shows deadline propagation example: Client → User Server (2s deadline, 0.5s processing) → Billing Server (1.5s timeout). Without propagation, downstream services waste resources processing requests that will arrive too late to benefit the user. Google SRE recommends: "each server in the request tree implements deadline propagation."

**Sources:** gRPC Documentation, Google SRE Book (Addressing Cascading Failures)

**Confidence:** HIGH

### Finding 3: Retries Amplify Load During Partial Failures

**Claim:** Naive retries can amplify load by 4x or more during partial failures, creating retry storms that overwhelm already-struggling dependencies.

**Evidence:** Google SRE example: Backend limit 10,000 QPS, frontend sends 10,100 QPS → 100 rejected → retried every second → total 10,200 QPS → more rejections → exponential growth. With retries at 3 layers (database, backend, frontend), a single user request creates 4^3=64 attempts on the database.

**Source:** Google SRE Book (Addressing Cascading Failures)

**Confidence:** HIGH

**Recommendation:** Implement randomized exponential backoff, limit retries per request (3 attempts recommended), implement per-client retry budgets (<10% retry ratio), and avoid retrying at multiple layers.

### Finding 4: Circuit Breakers Must Integrate with Retries

**Claim:** Circuit breakers should be combined with retry logic where retry policies respect circuit breaker state to prevent retrying when failures are not transient.

**Evidence:** Microsoft Circuit Breaker Pattern shows three states: Closed (normal), Open (fail immediately), Half-Open (trial requests). When circuit is Open, retries should stop. "The Retry pattern should be implemented such that when circuit breaker indicates a fault isn't transient, retry logic stops attempting retries." Resilience4j shows aspect ordering: Retry → CircuitBreaker → RateLimiter → TimeLimiter → Bulkhead.

**Sources:** Microsoft Azure Architecture Center, Resilience4j Documentation

**Confidence:** HIGH

### Finding 5: Idempotency Required for Safe Retries

**Claim:** Retrying non-idempotent operations (like payment processing) without deduplication can cause incorrect results like duplicate charges; idempotency keys are required for safe retry mechanisms.

**Evidence:** Microsoft Idempotent Consumer Pattern: "In at-least-once delivery systems, consumers may receive the same message multiple times due to producer retries (acknowledgment lost), consumer crashes before acknowledgment, or crashes after processing but before acknowledgment. Without idempotency protection, retrying a payment request could cause double charges."

**Source:** Microsoft Azure Architecture Center (Idempotent Consumer Pattern)

**Confidence:** HIGH

**Implementation:** Use idempotency keys with deduplication stores. Commit deduplication marker and business side effects in same transaction to avoid crash windows.

### Finding 6: Database and Queue Timeouts Are Required

**Claim:** Network API timeouts alone are insufficient; database queries, locks, transactions, and queue jobs also require timeout mechanisms to prevent resource exhaustion.

**Evidence:** PostgreSQL provides `statement_timeout` (aborts queries exceeding time), `lock_timeout` (aborts lock waits), `transaction_timeout` (terminates sessions), and `idle_in_transaction_session_timeout`. RabbitMQ documentation emphasizes that at-least-once delivery requires consumer acknowledgements with timeouts to prevent unacknowledged messages from consuming resources indefinitely.

**Sources:** PostgreSQL Documentation, RabbitMQ Documentation

**Confidence:** HIGH

### Finding 7: Client-Side Throttling During Overload

**Claim:** During overload, services should implement client-side throttling and criticality-based request rejection to maintain availability for important traffic while shedding less critical load.

**Evidence:** Google SRE describes per-customer quotas, client-side throttling (clients self-regulate when requests exceed K×accepts), and four criticality levels: CRITICAL_PLUS (most critical), CRITICAL (production), SHEDDABLE_PLUS (batch), SHEDDABLE (frequent unavailability). Services reject lower criticality requests first during overload.

**Source:** Google SRE Book (Handling Overload)

**Confidence:** HIGH

### Finding 8: Monitoring Must Include Timeout-Specific Metrics

**Claim:** Monitoring timeout rates, latency distributions, and circuit breaker state is essential for detecting emerging reliability issues before they cause outages.

**Evidence:** Google SRE Four Golden Signals: Latency, Traffic, Errors, Saturation. "Measuring your 99th percentile response time over some small window can give a very early signal of saturation." Rising P99 latency often precedes increased timeout rates. Health endpoints should expose liveness and readiness separately with dependency-specific timeouts in readiness probes.

**Sources:** Google SRE Book (Monitoring Distributed Systems), Microsoft Azure Architecture Center (Health Endpoint Monitoring)

**Confidence:** HIGH

## Areas of Agreement

All authoritative sources agree on:

1. **Timeouts prevent resource exhaustion** — Timeouts protect threads, connections, and memory from indefinite holding by slow operations
2. **Exponential backoff with jitter is required** — Prevents synchronized retry storms; randomization essential
3. **Idempotency required for safe retries** — Non-idempotent operations (payments, state changes) need deduplication
4. **Circuit breakers complement retries** — Stop retrying when failures are persistent, not transient
5. **Deadline propagation essential** — Use absolute deadlines, not fixed timeouts at each layer
6. **Database/queue timeouts necessary** — Network API timeouts alone insufficient for full reliability
7. **Monitor latency distributions** — P99/P999 matter more than mean; tail latency causes resource exhaustion
8. **Separate liveness from readiness** — Health checks should avoid false alarms from dependency issues

## Areas of Disagreement

No material contradictions discovered between authoritative sources. Minor differences exist in implementation details:

- **Retry budgets:** Some sources suggest 3 retries (Google SRE, Resilience4j default), others show 10 as example of naive implementation. All agree 3 is recommended.
- **Where to implement retries:** Some emphasize retry at highest layer (Azure), others specify retry only at immediate parent layer (Google SRE). Both approaches prevent multiplicative retry amplification.

## Limitations

1. **Research date:** 2026-09-26 — Practices continue evolving, particularly for cloud-native patterns
2. **Proprietary implementations:** Internal implementations at major companies (Google, Amazon) are described at high level; specific configurations not publicly available
3. **Platform-specific variations:** Different languages and frameworks have varying timeout semantics; research focused on general principles applicable across platforms
4. **Synthetic vs real-world scenarios:** Most recommendations based on observed failure patterns; specific configuration requires load testing in each environment

## Conclusion

Timeout configuration in distributed systems requires an integrated approach combining:
1. **Timeout budgets** calculated from production latency distributions (not arbitrary values)
2. **Deadline propagation** through service chains as absolute deadlines
3. **Exponential backoff with jitter** for retries (prevent synchronized retry storms)
4. **Circuit breakers** to stop retrying when failures are persistent
5. **Idempotency keys** for safe retries of non-idempotent operations
6. **Database and queue timeouts** beyond network API timeouts
7. **Criticality-based rejection** during overload to protect important traffic
8. **Timeout-specific monitoring** including P99 latency, timeout rates, and circuit breaker state

No single mechanism is sufficient; they form an interdependent reliability strategy where each component addresses a specific failure mode while relying on others to prevent secondary problems (e.g., retries preventing recovery without circuit breakers, or correct timeouts causing data inconsistency without idempotency).