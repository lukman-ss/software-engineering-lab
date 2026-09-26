# Research Report: Rate Limiting & Backpressure

## Research Question

What are the evidence-based practices for implementing rate limiting and backpressure in distributed systems to ensure system stability under overload conditions?

## Executive Summary

Rate limiting and backpressure are fundamental techniques for preventing cascading failures in distributed systems. Rate limiting controls the inbound request rate at the API layer using mechanisms like HTTP 429 and algorithms like Token Bucket. Backpressure propagates load signals downstream through queues, workers, and databases using patterns like exponential backoff and jitter. Key findings:

- **Token Bucket** is the dominant rate limiting algorithm, allowing controlled burst while enforcing long-term limits
- **HTTP 429** is the standard rate limiting response with optional Retry-After header
- **Little's Law (L = λW)** provides mathematical foundation for queue capacity planning
- **Exponential backoff with jitter** is essential to prevent retry storms
- **Multi-tenant isolation** requires per-tenant rate limiting and fair queueing
- **System stability** depends on arrival rate vs processing rate comparison

## Findings

### Finding 1: Token Bucket is the Primary Rate Limiting Algorithm

**Claim:** Token Bucket provides controlled burst while enforcing long-term rate limits

**Evidence:** Token Bucket maintains a bucket of tokens added at fixed rate r. When packet of n bytes arrives, if n tokens available, remove n and send; else packet is non-conformant. Bucket capacity b determines maximum burst. Burst time: T_max = b / (M - r) where M is max transmission rate. Token Bucket allows burst traffic while maintaining average rate limits.

**Sources:**
- Wikipedia - Token Bucket: https://en.wikipedia.org/wiki/Token_bucket
- RFC 6585: https://datatracker.ietf.org/doc/html/rfc6585 (mentions rate limiting conceptually)
- IEEE Datacenter Traffic Control: https://www.researchgate.net/publication/321744877

**Confidence:** HIGH

**Notes:** Token Bucket is also used for database IO flow control (ScyllaDB implementation) where tokens = normalized sum of IO request weight and length.

---

### Finding 2: HTTP 429 is the Standard Rate Limiting Response

**Claim:** HTTP 429 "Too Many Requests" is the standard status code for rate limiting responses

**Evidence:** RFC 6585 (April 2012) defines HTTP 429 as indicating rate limiting. Response MAY include Retry-After header indicating wait time. Server may identify users by authentication credentials, stateful cookies, or IP. Servers not required to use 429; may drop connections during attacks.

**Sources:**
- RFC 6585: https://datatracker.ietf.org/doc/html/rfc6585
- Wikipedia - Rate Limiting: https://en.wikipedia.org/wiki/Rate_limiting

**Confidence:** HIGH

**Notes:** Multi-tenant systems should use user_id, tenant_id, API key combinations rather than IP-only for authenticated APIs to avoid fairness issues when multiple users share IP via NAT.

---

### Finding 3: Little's Law Enables Queue Capacity Planning

**Claim:** Little's Law (L = λW) enables mathematical analysis of queue stability

**Evidence:** Little's Law states: "Average number in system (L) = Arrival rate (λ) × Average time in system (W)." Applies to any stationary, ergodic system. If arrival rate exceeds processing rate, system becomes unstable and backlog grows linearly over time.

**Sources:**
- Wikipedia - Little's Law: https://en.wikipedia.org/wiki/Little%27s_law
- MIT Lecture Notes (cited in Wikipedia): http://www.columbia.edu/~ks20/stochastic-I/stochastic-I-LL.pdf

**Confidence:** HIGH

**Notes:** Example calculation: 5M products / 2000 products/second = 2500 seconds ≈ 41 minutes 40 seconds minimum to clear backlog. No queue configuration can change this physical constraint.

---

### Finding 4: Exponential Backoff with Jitter Prevents Retry Storms

**Claim:** Exponential backoff with jitter spreads retry timing to prevent synchronized failures

**Evidence:** Without jitter: many clients retry simultaneously after exponential backoff. Full Jitter = random(0, min(cap, 2^attempt)) spreads retries evenly. AWS SDKs now support this as standard. With 100 contending clients, jitter reduces client work by >50%.

**Sources:**
- AWS Architecture Blog: https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/
- Wikipedia - Rate Limiting (mentions retry storms)

**Confidence:** HIGH

**Notes:** Retry storms cause cascading failures. Jitter variants: Full Jitter (best), Equal Jitter, Decorrelated Jitter. All better than no jitter.

---

### Finding 5: Multi-tenant Systems Require Fair Queueing

**Claim:** Single-tenant monopolization prevented by per-tenant rate limits and fair queueing

**Evidence:** Tenant A import (2M records) can monopolize workers without isolation. Solution: per-tenant rate limit (e.g., max 5 concurrent jobs per tenant), fair queueing, concurrency limits, priority queues.

**Sources:**
- User lab document (original research topic)
- IEEE Datacenter Traffic Control: https://www.researchgate.net/publication/321744877

**Confidence:** MEDIUM

**Notes:** Without fairness, free users, enterprise tenants, and internal services compete unfairly. Resource footprint (memory/CPU) vs precision trade-off in datacenter rate limiters.

---

### Finding 6: Reactive Streams Standardizes Backpressure

**Claim:** Reactive Streams provides standardized backpressure for asynchronous stream processing

**Evidence:** Reactive Streams (2013, JVM 2015, Java 9 JEP 266) defines interfaces for non-blocking back pressure. Goal: prevent downstream from buffering arbitrary amounts of data. Adopted by Akka Streams, Spring Reactor, Netflix RxJava, Vert.x, Apache Kafka, Elasticsearch, Cassandra.

**Sources:**
- Wikipedia - Reactive Streams: https://en.wikipedia.org/wiki/Reactive_Streams
- Reactive Streams Specification: http://www.reactive-streams.org/

**Confidence:** HIGH

**Notes:** Backpressure integral to model to allow bounded queues between threads. specification provides minimal interfaces for interoperability.

---

### Finding 7: System Stability Depends on Arrival Rate vs Processing Rate

**Claim:** Key question for system health: "Does work arrive faster than system can process?"

**Evidence:** Little's Law (L = λW) shows queue depth directly proportional to arrival rate and inversely proportional to processing rate. Monitoring metrics: arrival rate, processing rate, queue depth, oldest job age, retry rate, failure rate. Queue age more informative than raw count.

**Sources:**
- User lab document
- Wikipedia - Little's Law
- AWS Architecture Blog (retry patterns)
- Wikipedia - Rate Limiting (monitoring section)

**Confidence:** HIGH

**Notes:** 10,000 jobs isn't necessarily bad. If oldest job is 45 minutes when normal is 10 seconds, problem exists. Backlog age indicates degradation before capacity exhaustion.

---

### Finding 8: Distributed Rate Limiting Requires Shared State

**Claim:** Rate limiting across distributed instances requires centralized storage

**Evidence:** Wikipedia mentions Redis/Aerospike for in-memory key-value rate limiting. Token bucket in distributed system needs shared state across instances. Options: Redis (as in Lab 25), centralized rate limiter, consistent hashing to sticky instances.

**Sources:**
- Wikipedia - Rate Limiting: https://en.wikipedia.org/wiki/Rate_limiting
- Medium - Alternative Approach: https://medium.com/figma-design/an-alternative-approach-to-rate-limiting-f8a06cf7c94c

**Confidence:** MEDIUM

**Notes:** No standard solution specified. Trade-offs: accuracy vs latency vs availability. Redis most common approach (as in user lab specification).

---

## Areas of Agreement

1. **Rate limiting algorithms**: Token Bucket, Leaky Bucket, Fixed Window, Sliding Window all documented consistently across sources.

2. **HTTP 429**: Standard mechanism universally recognized.

3. **Backpressure**: Consensus that queues are temporary buffers, not infinite capacity.

4. **Exponential backoff + jitter**: Essential for preventing retry storms.

5. **Little's Law**: Universally accepted mathematical relationship.

6. **Multi-tenant fairness**: Per-tenant limits required for fair resource allocation.

---

## Areas of Disagreement

**No significant factual disagreements found.** Minor differences are in:

1. **Leaky Bucket terminology**: Two versions (meter vs queue) cause confusion in literature but same underlying principle.

2. **Implementation choices**: Different approaches (Redis vs sliding window log) are complementary, not contradictory.

3. **Terminology**: Rate limiting vs throttling terminology varies but mechanisms are clear.

---

## Limitations

1. **Source language**: Primary sources in English; limited Indonesian technical resources for rate limiting/backpressure.

2. **Proprietary implementations**: Exact algorithms used by major cloud providers may be opaque.

3. **Performance benchmarks**: Vary significantly by workload; no universal benchmarks.

4. **Emerging patterns**: Some cloud-native patterns (e.g., gRPC flow control, Kubernetes queue proxy) not fully explored in sources reviewed.

5. **Production case studies**: Engineering blogs provide implementations but limited public data on failures.

---

## Conclusion

Rate limiting and backpressure are essential for building resilient distributed systems. The evidence base is strong for:

- **Token Bucket** as primary rate limiting algorithm
- **HTTP 429** as standard response
- **Little's Law** for queue capacity planning
- **Exponential backoff with jitter** to prevent retry storms
- **Multi-tenant isolation** via per-tenant limits

System health depends on monitoring **arrival rate vs processing rate**, not just queue depth. Backpressure is an integral design principle, not an afterthought. The ability to say "cukup" (enough) - to reject, slow down, or queue work - is what makes systems scalable rather than fragile.

**Key metric for system health:** "Apakah pekerjaan masuk lebih cepat daripada kemampuan sistem menyelesaikannya?" If yes, facing capacity or backpressure problem requiring intervention.

**Minimum time to clear 5M backlog at 2000/sec throughput:** ~41 minutes 40 seconds (5,000,000 / 2,000 = 2,500 seconds). No queue configuration changes this physical constraint - only increasing processing capacity or reducing incoming rate.