# Research Report: Rate Limiting & Backpressure

## Research Question

How do rate limiting and backpressure mechanisms work in distributed systems, what algorithms are commonly used, and how should systems handle overload conditions to maintain availability and prevent cascading failures?

---

## Executive Summary

This research investigated rate limiting and backpressure mechanisms in distributed systems. Rate limiting controls the rate of incoming requests at system boundaries, while backpressure propagates pressure signals backward through the system when downstream components are overwhelmed.

Key findings:
- **Token Bucket** and **Leaky Bucket** are the primary rate limiting algorithms, mathematically equivalent but differing in implementation perspective
- **Exponential Backoff with Full Jitter** prevents thundering herd problems; AWS SDK formula: `delay = random(0,1) × min(20000, base_delay × 2^retry)`
- **HTTP 429 Too Many Requests** is the standard status code (RFC 6585)
- **Little's Law (L = λW)** provides the fundamental relationship for queue analysis
- Multi-layer approaches combining rate limiters with load shedders are most effective for production systems (Stripe uses 4 limiters)

The topic specification's calculation is correct: 5,000,000 items / 2,000 items/sec = 2,500 seconds ≈ 41 minutes 40 seconds.

---

## Findings

### Finding 1: Token Bucket and Leaky Bucket Algorithm Properties

**Claim:** Token bucket allows burst traffic up to bucket capacity while maintaining long-term average rate limits. Leaky bucket has two implementations: as a meter (equivalent to token bucket) and as a queue (special case).

**Evidence:** Wikipedia Token bucket: "A conforming flow can thus contain traffic with an average rate up to the rate at which tokens are added to the bucket, and have a burstiness determined by the depth of the bucket." Wikipedia Leaky bucket: "The leaky bucket as a meter is exactly equivalent to (a mirror image of) the token bucket algorithm."

**Sources:**
- Token bucket Wikipedia (Tier 2) - https://en.wikipedia.org/wiki/Token_bucket
- Leaky bucket Wikipedia (Tier 2) - https://en.wikipedia.org/wiki/Leaky_bucket
- NGINX limit_req_module (Tier 1) - https://nginx.org/en/docs/http/ngx_http_limit_req_module.html

**Confidence:** HIGH

**Implementation Notes:**
- NGINX uses leaky bucket as a meter (checking conformance), not as a queue
- Burst parameter in NGINX allows temporary excess requests
- Used by Stripe, AWS API Gateway, NGINX

---

### Finding 2: Backpressure vs Rate Limiting

**Claim:** Rate limiting typically works at system entrance, while backpressure is broader and propagates pressure backward through the system.

**Evidence:** "Rate limiting usually works at the entrance to a system. Backpressure is more general: When downstream is unable to keep up with work from upstream, upstream must slow down."

**Sources:**
- Google SRE Workbook: Handling Overload (Tier 1) - https://landing.google.com/sre/sre-book/chapters/handling-overload/
- AWS SDK retry behavior (Tier 1) - https://docs.aws.amazon.com/sdkref/latest/guide/feature-retry-behavior.html
- Stripe engineering blog (Tier 2) - https://stripe.com/blog/rate-limiters

**Confidence:** HIGH

**Implementation Notes:**
- Rate limiting: External traffic control at API gateway
- Backpressure: Internal propagation (queue depth, worker saturation)
- Both work together for system resilience

---

### Finding 3: Exponential Backoff with Jitter

**Claim:** Adding jitter to exponential backoff significantly reduces thundering herd problems. AWS SDK uses Full Jitter: `delay = random(0, 1) × min(20,000 ms, base_delay × 2^retry)`.

**Evidence:** AWS Architecture Blog: "In the case with 100 contending clients, we've reduced our call count by more than half." AWS SDK: "The SDK computes each retry delay using this formula: delay = random(0, 1) × min(20,000 ms, base_delay × 2^retry)"

**Sources:**
- AWS Architecture Blog (Tier 2) - https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/
- AWS SDK retry behavior (Tier 1) - https://docs.aws.amazon.com/sdkref/latest/guide/feature-retry-behavior.html
- Exponential backoff Wikipedia (Tier 2) - https://en.wikipedia.org/wiki/Exponential_backoff

**Confidence:** HIGH

**Implementation Notes:**
- Full Jitter: delay = random(0, 1) × min(20000, base_delay × 2^retry)
- Transient errors: 50ms base delay (25ms for DynamoDB)
- Throttling errors: 1000ms base delay
- Max cap: 20 seconds
- Three jitter variants: Full, Equal, Decorrelated; Full recommended

---

### Finding 4: HTTP 429 Status Code

**Claim:** HTTP 429 Too Many Requests is the standard status code for rate limiting, defined in RFC 6585.

**Evidence:** RFC 6585: "The 429 status code indicates that the user has sent too many requests in a given amount of time ('rate limiting')." "Responses with the 429 status code MUST NOT be stored by a cache."

**Sources:**
- RFC 6585 (Tier 1) - https://www.rfc-editor.org/rfc/rfc6585#section-4
- Stripe engineering blog (Tier 2) - https://stripe.com/blog/rate-limiters
- NGINX rate limiting (Tier 1) - https://nginx.org/en/docs/http/ngx_http_limit_req_module.html

**Confidence:** HIGH

**Notes:**
- Should include Retry-After header when possible
- MUST NOT be cached by intermediaries
- 503 may be more appropriate for server overload vs 429 for client quota exceeded
- NGINX defaults to 503 but makes it configurable

---

### Finding 5: Little's Law for Queue Analysis

**Claim:** Little's Law (L = λW) provides the mathematical foundation for understanding queue behavior. Topic specification calculation is correct.

**Evidence:** "the long-term average number of customers (L) in a stationary system is equal to the long-term average effective arrival rate (λ) multiplied by the average time that a customer spends in the system (W)." Calculation: 5,000,000 / 2,000 = 2,500 seconds ≈ 41 minutes 40 seconds.

**Sources:**
- Little's law Wikipedia (Tier 2) - https://en.wikipedia.org/wiki/Little%27s_law
- Google SRE Workbook (Tier 1) - https://landing.google.com/sre/sre-book/chapters/handling-overload/
- Topic specification calculation verified

**Confidence:** HIGH

**Application:**
- Topic spec: 5,000,000 items / 2,000 items/sec = 2,500 seconds = 41m 40s ✓
- Queue age is more useful than raw queue length for detecting issues
- Backlog growth = (arrival_rate - processing_rate) × time

---

### Finding 6: Multi-Layer Rate Limiting Approach (Stripe)

**Claim:** Effective production systems use multiple layers of rate limiting and load shedding.

**Evidence:** Stripe uses 4 types of limiters:
1. Request rate limiter - per-user request rate (most important)
2. Concurrent requests limiter - max in-flight requests
3. Fleet usage load shedder - reserve capacity for critical traffic
4. Worker utilization load shedder - shed low-priority traffic under pressure

**Sources:**
- Stripe engineering blog (Tier 2) - https://stripe.com/blog/rate-limiters
- AWS SDK retry behavior (Tier 1) - https://docs.aws.amazon.com/sdkref/latest/guide/feature-retry-behavior.html
- Google SRE Workbook (Tier 1) - https://landing.google.com/sre/sre-book/chapters/handling-overload/

**Confidence:** HIGH

**Implementation Notes:**
- Layered approach allows graceful degradation
- Load shedders help during incidents without affecting core functionality
- Criticality levels allow prioritized request handling

---

### Finding 7: Redis as Rate Limiting Backend

**Claim:** Redis is well-suited for distributed rate limiting due to atomic operations and data structures.

**Evidence:** "Redis provides the following features that make it a good fit for rate limiting: [INCR and EXPIRE] give you atomic fixed-window counters with automatic time-window cleanup. [Hashes, sorted sets, and strings] cover the data shapes needed for sliding window and token bucket algorithms."

**Sources:**
- Redis rate limiter documentation (Tier 1) - https://redis.io/docs/latest/develop/use-cases/rate-limiter/
- Stripe engineering blog (Tier 2) - https://stripe.com/blog/rate-limiters

**Confidence:** HIGH

**Common Patterns:**
- Fixed window: INCR + EXPIRE
- Token bucket: Lua scripts
- Sliding window: Sorted sets with timestamps
- Sub-millisecond latency for synchronous request path

---

### Finding 8: Google SRE Per-Customer Limits and Criticality

**Claim:** Multi-tenant systems require per-customer rate limiting with criticality levels to prevent noisy neighbors.

**Evidence:** Google SRE implements per-customer CPU quotas (Gmail 4000 CPU seconds/sec, Calendar 4000, Android 3000, Google+ 2000, others 500). Four criticality values: CRITICAL_PLUS, CRITICAL, SHEDDABLE_PLUS, SHEDDABLE.

**Sources:**
- Google SRE Workbook (Tier 1) - https://landing.google.com/sre/sre-book/chapters/handling-overload/
- Stripe Fleet Usage Load Shedder (Tier 2) - https://stripe.com/blog/rate-limiters

**Confidence:** HIGH

**Implementation Notes:**
- System relies on customers not hitting limits simultaneously
- Criticality propagation through RPC system
- Higher criticality = higher threshold before rejection
- Separates batch traffic from interactive traffic

---

### Finding 9: Google SRE Client-Side Throttling and Retry Budget

**Claim:** Client-side throttling and retry budgets prevent retry storms during cascading overload.

**Evidence:** Google SRE: "per-request retry budget of up to three attempts... per-client retry budget... ratio of requests that correspond to retries... below 10%." Client-side throttling: "each client task keeps... requests and accepts... Clients can continue to issue requests until requests is K times as large as accepts."

**Sources:**
- Google SRE Workbook (Tier 1) - https://landing.google.com/sre/sre-book/chapters/handling-overload/
- AWS SDK retry behavior (Tier 1) - https://docs.aws.amazon.com/sdkref/latest/guide/feature-retry-behavior.html

**Confidence:** HIGH

**Implementation Notes:**
- Per-request budget: max 3 attempts
- Per-client budget: max 10% retry ratio
- Threefold growth capped to 1.1x with per-client budget
- AWS SDK uses token bucket retry quota (500 tokens capacity)
- Client-side throttling prevents overwhelming backends with rejected requests

---

## Areas of Agreement

1. **Algorithm Equivalence:** Token bucket and leaky bucket (as meter) are mathematically equivalent (mirror images)
2. **Jitter Benefits:** All sources agree jitter is critical for preventing thundering herd
3. **HTTP 429 Standard:** RFC 6585 and implementations consistently use 429 for rate limiting
4. **Queue Theory Foundation:** Little's Law is universally accepted for queue analysis (L = λW)
5. **Multi-Layer Approach:** Production systems benefit from multiple limiting layers (Stripe's 4 limiters)
6. **Backpressure Definition:** Backpressure is broader than rate limiting; propagates pressure backward

---

## Areas of Disagreement

1. **HTTP Status Code Selection:** 429 (RFC 6585) for rate limiting vs 503 (NGINX default) for overload - no contradiction, different use cases
2. **Base Backoff Delays:** Different protocols use different base delays (50ms AWS SDK vs 500ms SIP vs 51.2μs Ethernet) - protocol-specific, no contradiction
3. **Jitter Algorithm Choice:** AWS explored Full, Equal, Decorrelated; standardized on Full Jitter in SDK - exploration vs standardization, no contradiction

---

## Limitations

1. **Language Barrier:** Many authoritative sources are in English; topic specification is in Bahasa Indonesia
2. **Implementation Details:** Some rate limiter implementations may be proprietary (e.g., Stripe's production code)
3. **Version Variations:** Different systems use different algorithm parameter values (base delays, max caps)
4. **Context-Specific:** What works for one system may not work identically in another due to different traffic patterns
5. **Calculation Scope:** The 2,500-second calculation assumes constant throughput and no additional capacity

---

## Conclusion

Rate limiting and backpressure are complementary mechanisms for system resilience. Rate limiting controls external traffic at system boundaries using algorithms like token bucket, while backpressure propagates internal pressure signals to prevent cascading failures.

The most effective production systems implement:
1. Token bucket or leaky bucket for request rate limiting
2. Exponential backoff with full jitter for retry mechanisms (`delay = random(0,1) × min(20000, base_delay × 2^retry)`)
3. Multi-layered approach combining rate limiters with load shedders (Stripe's 4 limiters pattern)
4. Per-user, per-tenant, and per-endpoint limits with criticality levels (Google SRE 4-level criticality)
5. Monitoring of queue age rather than just queue depth
6. Retry budgets (per-request: 3 attempts, per-client: 10% ratio) to prevent retry storms

The calculations in the topic specification are mathematically correct and verifiable through Little's Law: 5,000,000 items / 2,000 items/sec = 2,500 seconds ≈ 41 minutes 40 seconds.

---

**Research Date:** 2026-09-26  
**Total Sources Reviewed:** 13  
**Confidence Level:** HIGH for all major findings