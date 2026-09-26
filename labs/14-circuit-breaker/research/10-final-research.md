# Final Research

## Research Topic
Circuit Breaker Pattern — Cara Mencegah Satu Service Down Menjatuhkan Seluruh Sistem

## Executive Summary

This research synthesizes Tier 1 authoritative sources (Martin Fowler, Microsoft Azure Architecture Center, Resilience4j reference implementation, Google SRE Book) on the Circuit Breaker pattern for distributed systems.

**Core finding**: Circuit breaker does not heal a failing dependency. It limits the blast radius of that failure by failing fast (not calling downstream when OPEN), giving the dependency recovery time, and safely probing recovery via HALF-OPEN state.

**State machine consensus**: All sources agree on three states — CLOSED (normal forwarding), OPEN (fail fast, no downstream calls), HALF-OPEN (limited probe calls) — and four transitions:
1. CLOSED → OPEN (failures ≥ threshold)
2. OPEN → HALF_OPEN (cooldown elapsed)
3. HALF_OPEN → CLOSED (all probes succeed)
4. HALF_OPEN → OPEN (any probe fails)

**Complementary patterns**: Timeout bounds single-request wait; Retry re-issues transient failures with bounded attempts + exponential backoff; Circuit Breaker stops traffic entirely to failing dependency. Bulkhead isolates resources per dependency. Fallback provides degraded responses. All are complementary, not substitutes.

---

## Finding 1: State Machine
**Claim**: The canonical Circuit Breaker state machine has three states (CLOSED/OPEN/HALF_OPEN) with four transitions as specified above.

**Evidence**:
- Martin Fowler (2014): https://martinfowler.com/bliki/CircuitBreaker.html
- Azure (2025): https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker
- Resilience4j: https://resilience4j.readme.io/docs/circuitbreaker

**Sources**: 3 Tier 1 sources converge.

**Confidence**: HIGH.

---

## Finding 2: Fail-Fast Behavior
**Claim**: When OPEN, the breaker rejects calls immediately without invoking the downstream service.

**Evidence**:
- Martin Fowler: "calling the circuit breaker will call the underlying block if the circuit is closed, but return an error if it's open"
- Azure: "OPEN: The request from the application fails immediately"
- Resilience4j: "The CircuitBreaker rejects calls with a CallNotPermittedException when it is OPEN"

**Sources**: 3 Tier 1 sources.

**Confidence**: HIGH.

---

## Finding 3: Failure Counter Reset
**Claim**: Successful calls reset the failure counter to zero in CLOSED state.

**Evidence**:
- Martin Fowler: "successful calls reset it back to zero"
- Azure: "CLOSED: The request from the application is routed to the operation. The proxy maintains a count... If the call to the operation is unsuccessful, the proxy increments this count."
- Resilience4j sliding window replaces old values with new successes

**Sources**: 3 Tier 1 sources.

**Confidence**: HIGH.

---

## Finding 4: Limited Probe Calls
**Claim**: HALF_OPEN allows only a limited number of concurrent probe calls to prevent overwhelming a recovering service.

**Evidence**:
- Azure: "A limited number of requests from the application are allowed to pass through"
- Resilience4j: `permittedNumberOfCallsInHalfOpenState` (default: 10)

**Sources**: 2 Tier 1 sources explicitly; Martin Fowler describes one-at-a-time trial implicitly.

**Confidence**: HIGH.

---

## Finding 5: Timeout Necessity
**Claim**: Outbound HTTP calls must have explicit timeouts; Go `http.Client.Timeout` should be set for every downstream call. Circuit breaker adds fail-fast on top of timeout — they are complementary.

**Evidence**:
- Go docs: https://pkg.go.dev/net/http#Client
- Azure: "A circuit breaker might not fully protect applications from failures in external services that have long time-out periods."

**Sources**: 2 Tier 1 sources.

**Confidence**: HIGH.

---

## Finding 6: Retry Must Be Bounded
**Claim**: Retry should be bounded (max attempts / retry budget) and respect circuit breaker exceptions.

**Evidence**:
- Google SRE: per-request retry budget (max 3 attempts); per-client retry budget (max 10% retries)
- Azure: "the retry logic should be sensitive to any exceptions that the circuit breaker returns and stop retry attempts if the circuit breaker indicates that a fault isn't transient"

**Sources**: 2 Tier 1 sources.

**Confidence**: HIGH.

---

## Finding 7: Bulkhead Distinctness
**Claim**: Bulkhead isolates resources per dependency; Circuit Breaker gates traffic. They are complementary and often combined.

**Evidence**:
- Azure Bulkhead: "consider combining bulkheads with retry, circuit breaker, and throttling patterns"
- Resilience4j: "If you want to restrict the number of concurrent threads, please use a Bulkhead. You can combine a Bulkhead and a CircuitBreaker."

**Sources**: 2 Tier 1 sources.

**Confidence**: HIGH.

---

## Finding 8: Observability
**Claim**: Every Circuit Breaker state transition should be logged, exposed for monitoring, and trigger alerts.

**Evidence**:
- Martin Fowler: "Any change in breaker state should be logged and breakers should reveal details of their state for deeper monitoring."
- Azure: "If the circuit breaker raises an event each time it changes state, this information can help monitor the health of the protected system component"

**Sources**: 2 Tier 1 sources.

**Confidence**: HIGH.

---

## Areas of Agreement
All Tier 1 sources agree on:
1. Three states + four transitions (state machine)
2. Fail-fast when OPEN (no downstream calls)
3. Failure counter reset on success (CLOSED)
4. Limited probe calls in HALF_OPEN
5. Circuit Breaker ≠ Retry ≠ Timeout (complementary patterns)
6. State transitions must be observable
7. Error types should be classified (not all exceptions trip the circuit)

## Areas of Disagreement
1. **Failure evaluation window**: Azure recommends a time-based window; Resilience4j offers both count-based and time-based; Martin Fowler's original example uses simple consecutive-count. No universal algorithm — depends on traffic volume and detection latency.
2. **Automatic vs call-triggered HALF_OPEN transition**: Resilience4j supports both modes (`automaticTransitionFromOpenToHalfOpenEnabled`); Martin Fowler's example is call-triggered; Azure describes a timer. This lab uses call-triggered for simplicity.
3. **Exact numeric recommendations**: No Tier 1 source provides universal numeric values for thresholds, cooldowns, or timeouts. All are "tune to your SLA" guidance.

## Limitations
1. AWS Builders Library page on jitter could not be fully rendered this session — jitter details marked MEDIUM confidence.
2. No quantitative benchmarks exist in authoritative sources for "at what threshold does cascade failure occur" — marked NOT VERIFIED.
3. Financial fallback best practices (PPOB) are derived from the principle "do not make unsafe assumptions" — no authoritative source provides a universal financial fallback recipe.

## Conclusion
The research confirms the Circuit Breaker state machine, transition rules, and complementary patterns (Timeout, Retry, Bulkhead, Fallback) with HIGH confidence across multiple Tier 1 sources. The lab implementation (consecutive-failure counter, call-triggered HALF_OPEN, explicit HTTP timeout, limited probe calls) is consistent with the canonical pattern while simplified for educational clarity.

Per the core thesis: Circuit Breaker does not make a failing dependency healthy. It limits the blast radius of that failure so the caller can fail fast, preserve resources, degrade gracefully, and give the dependency time to recover.

---

## Case Study References

### CMMS WhatsApp (Design Discussion)
Flow: `Create Invoice → Generate PDF → Send WhatsApp`
- Principle: WhatsApp down ≠ Invoice down
- Pattern: Synchronous invoice DB write → async queue (message broker) → worker with Circuit Breaker + exponential backoff retry around WhatsApp API call
- WhatsApp is a non-critical dependency: failures must not roll back invoice creation
- Idempotency key on queued message to prevent duplicate notifications on retry

### PPOB (Criticality Analysis)
Flow: `User → Order → Payment Gateway → Provider Pulsa → WhatsApp`
- **Payment Gateway**: Synchronous critical. If down → fail fast, do not deduct balance, do not complete order. Circuit Breaker OPEN → reject checkout immediately.
- **Provider Pulsa**: Async with strict idempotency key. Retryable, queued for background processing.
- **WhatsApp**: Async non-critical. Failures must not abort fulfillment.

No unsafe financial fallback assumptions made. All financial mutations require explicit upstream confirmation.

---

## Date
Research conducted: 2026-09-26
Sources accessed: 2026-09-26
Freshness: Azure sources published 2024-2026 (current); Martin Fowler 2014 (canonical/unchanged); Google SRE 2017 (canonical/unchanged); Resilience4j (current docs)