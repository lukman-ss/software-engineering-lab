# Research Report: Rate Limiting & Backpressure

## Research Question

How do rate limiting and backpressure mechanisms work in distributed systems, what algorithms are commonly used, and how should systems handle overload conditions to maintain availability and prevent cascading failures?

---

## Executive Summary

This research investigated rate limiting and backpressure mechanisms in distributed systems. Rate limiting controls the rate of incoming requests at system boundaries, while backpressure propagates pressure signals backward through the system when downstream components are overwhelmed.

Key findings:
- **Token Bucket** and **Leaky Bucket** are the primary rate limiting algorithms, mathematically equivalent but differing in implementation perspective
- **Exponential Backoff with Jitter** prevents thundering herd problems during retries
- **HTTP 429** is the standard status code for rate limiting
- **Little's Law (L = λW)** provides the fundamental relationship for queue analysis
- Multi-layer approaches combining rate limiters with load shedders are most effective for production systems

The topic specification's claim about queue calculations (5,000,000 / 2,000 = 2,500 seconds ≈ 41 minutes 40 seconds) is correct and verifiable through Little's Law.

---

## Findings

### Finding 1: Token Bucket Algorithm

**Claim:** Token bucket allows burst traffic up to bucket capacity while maintaining long-term average rate limits.

**Evidence:** The token bucket algorithm adds tokens at a fixed rate and removes tokens per request. A conforming flow can contain traffic with an average rate up to the token addition rate, with burstiness determined by bucket depth.

**Sources:**
- Token bucket Wikipedia (Tier 2)
- NGINX limit_req_module (Tier 1)
- AWS API Gateway throttling (Tier 1)

**Confidence:** HIGH

**Implementation Notes:**
- Used by Stripe, AWS API Gateway
- Allows "burst" of traffic up to capacity before throttling
- Token addition rate = steady-state limit
- Bucket size = maximum burst size

---

### Finding 2: Backpressure vs Rate Limiting

**Claim:** Rate limiting typically works at system entrance, while backpressure is broader and propagates pressure backward through the system.

**Evidence:** "Rate limiting usually works at the entrance to a system. Backpressure is more general: When downstream is unable to keep up with work from upstream, upstream must slow down."

**Sources:**
- Google SRE Workbook: Handling Overload (Tier 1)
- AWS retry behavior documentation (Tier 1)
- Stripe engineering blog (Tier 2)

**Confidence:** HIGH

**Implementation Notes:**
- Rate limiting: External traffic control (firewall at entrance)
- Backpressure: Internal propagation (signal downstream issues upstream)
- Both work together for system resilience

---

### Finding 3: Exponential Backoff with Jitter

**Claim:** Adding jitter to exponential backoff significantly reduces thundering herd problems and improves system recovery.

**Evidence:** "In the case with 100 contending clients, we've reduced our call count by more than half. We've also significantly improved the time to completion, when compared to un-jittered exponential backoff."

**Sources:**
- AWS Architecture Blog (Tier 2)
- AWS SDK retry behavior (Tier 1)
- Exponential backoff Wikipedia (Tier 2)

**Confidence:** HIGH

**Implementation Notes:**
- AWS SDK formula: `delay = random(0, 1) × min(20000, base_delay × 2^retry)`
- Transient errors: 50ms base delay (AWS SDK v3 default; must be calibrated to downstream service SLA/latency)
- Throttling errors: 1000ms base delay (AWS SDK v3 default; must be calibrated to downstream service SLA/latency)
- Max cap: 20 seconds (AWS SDK v3 default; must be calibrated to downstream service timeout profile)
- AWS SDKs use "Full Jitter" variant (specific to AWS SDK implementation; other ecosystems may differ)

---

### Finding 4: HTTP 429 Status Code

**Claim:** HTTP 429 Too Many Requests is the standard status code for rate limiting.

**Evidence:** "The 429 status code indicates that the user has sent too many requests in a given amount of time ('rate limiting')."

**Sources:**
- RFC 6585: Additional HTTP Status Codes (Tier 1)
- Stripe engineering blog (Tier 2)
- NGINX rate limiting (Tier 1)

**Confidence:** HIGH

**Notes:**
- Should include Retry-After header when possible
- MUST NOT be cached by intermediaries
- 503 may be more appropriate for server overload vs 429 for client quota exceeded

---

### Finding 5: Little's Law for Queue Analysis

**Claim:** Little's Law (L = λW) enables calculation of queue metrics from arrival and processing rates.

**Evidence:** "The long-term average number of customers (L) in a stationary system is equal to the long-term average effective arrival rate (λ) multiplied by the average time that a customer spends in the system (W)."

**Sources:**
- Little's law Wikipedia (Tier 2)
- Google SRE Workbook (Tier 1)
- Topic specification calculations verified

**Confidence:** HIGH

**Application:**
- Topic spec: 5,000,000 items / 2,000 items/sec = 2,500 seconds = 41m 40s ✓
- This calculation is mathematically correct
- Queue age is more useful than raw queue length for detecting issues

---

### Finding 6: Multi-Layer Rate Limiting Approach

**Claim:** Effective production systems use multiple layers of rate limiting and load shedding.

**Evidence:** Stripe uses 4 types of limiters:
1. Request rate limiter - per-user request rate
2. Concurrent requests limiter - max in-flight requests
3. Fleet usage load shedder - reserve capacity for critical traffic
4. Worker utilization load shedder - shed low-priority traffic under pressure

**Sources:**
- Stripe engineering blog (Tier 2)
- AWS SDK retry behavior (Tier 1)
- Google SRE Workbook (Tier 1)

**Confidence:** HIGH

**Implementation Notes:**
- Layered approach allows graceful degradation
- Load shedders help during incidents without affecting core functionality
- Criticality levels allow prioritized request handling

---

### Finding 7: Redis as Rate Limiting Backend

**Claim:** Redis is well-suited for distributed rate limiting due to atomic operations.

**Evidence:** Redis provides atomic INCR/EXPIRE, Lua scripting for transactional operations, and sub-millisecond latency.

**Sources:**
- Redis rate limiter documentation (Tier 1)
- Stripe engineering blog (Tier 2)
- AWS API Gateway (Tier 1)

**Confidence:** HIGH

**Common Patterns:**
- Fixed window: INCR + EXPIRE
- Token bucket: Lua scripts with sorted sets
- Sliding window: Sorted sets with timestamps

---

## Areas of Agreement

1. **Algorithm Equivalence:** Token bucket and leaky bucket are mathematically equivalent (mirror images)
2. **Jitter Benefits:** All sources agree jitter is critical for preventing thundering herd
3. **HTTP 429 Standard:** RFC 6585 and implementations consistently use 429
4. **Queue Theory Foundation:** Little's Law is universally accepted for queue analysis
5. **Multi-Layer Approach:** Production systems benefit from multiple limiting layers

---

## Areas of Disagreement

1. **Token Bucket vs Leaky Bucket Implementation:** Some systems use one algorithm over the other based on implementation preferences, not functional differences
2. **Jitter Algorithm Choice:** AWS explored multiple variants (Full, Equal, Decorrelated) but standardized on Full Jitter
3. **Base Backoff Delays:** Different protocols use different base delays (50ms for AWS cloud APIs vs 500ms for SIP telephony)
4. **HTTP Status Code Selection:** Some systems use 429 for all rate limiting, others use 503 for server overload vs 429 for client quotas

---

## Limitations

1. **Language Barrier:** Many authoritative sources are in English; topic specification is in Bahasa Indonesia
2. **Implementation Details:** Some rate limiter implementations may be proprietary (e.g., Stripe's production code)
3. **Version Variations:** Different systems use different algorithm parameter values (base delays, max caps)
4. **Context-Specific:** What works for one system may not work identically in another due to different traffic patterns

---

## Conclusion

Rate limiting and backpressure are complementary mechanisms for system resilience. Rate limiting controls external traffic at system boundaries using algorithms like token bucket, while backpressure propagates internal pressure signals to prevent cascading failures.

The most effective production systems implement:
1. Token bucket or leaky bucket for request rate limiting
2. Exponential backoff with jitter for retry mechanisms
3. Multi-layered approach combining rate limiters with load shedders
4. Per-user, per-tenant, and per-endpoint limits with cost-based tuning
5. Monitoring of queue age rather than just queue depth

The calculations in the topic specification are mathematically correct and verifiable through Little's Law.

---

**Research Date:** 2026-09-26  
**Total Sources Reviewed:** 13  
**Confidence Level:** HIGH for all major findings