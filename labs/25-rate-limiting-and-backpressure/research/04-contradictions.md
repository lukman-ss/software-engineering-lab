# Contradictions and Disagreements

## No Material Contradictions Discovered

The research sources generally agree on the fundamental concepts of rate limiting and backpressure. However, there are some areas of nuance and implementation differences:

---

## Area 1: Token Bucket vs Leaky Bucket Implementation

**Source A (Token Bucket Wikipedia):**
- Token bucket adds tokens at a fixed rate, removes tokens per request
- Allows bursts up to bucket capacity
- Favored for API rate limiting (Stripe, AWS API Gateway)

**Source B (Leaky Bucket Wikipedia, NGINX Documentation):**
- Leaky bucket processes at fixed rate, queues/drops excess
- NGINX uses leaky bucket as meter for rate limiting

**Assessment:**
These are mathematically equivalent (mirror images). The difference is primarily in implementation perspective:
- Token bucket: "Do I have tokens?" → accept/reject
- Leaky bucket as meter: "Will adding water overflow?" → accept/reject
Both can achieve the same rate limiting behavior with appropriate parameters.
NGINX uses leaky bucket as a meter (checking conformance), not as a queue (enforcing).

---

## Area 2: Jitter Algorithm Variants

**Source A (AWS Architecture Blog):**
- Three variants: Full Jitter, Equal Jitter, Decorrelated Jitter
- "Full Jitter" and "Equal Jitter" perform similarly
- Recommends Full Jitter as simple and effective

**Source B (AWS SDK Documentation):**
- Uses "Full Jitter" specifically: `random(0, 1) × min(20000, base_delay × 2^retry)`
- Single standardized implementation across SDKs

**Assessment:**
No contradiction - the SDK documentation specifies the exact implementation (Full Jitter) while the blog post explored alternatives and found Full Jitter and Equal Jitter to be comparable in effectiveness.

---

## Area 3: HTTP Status Code for Rate Limiting

**Source A (RFC 6585):**
- 429 Too Many Requests is the standard code
- "Responses with the 429 status code MUST NOT be stored by a cache"

**Source B (Stripe Blog):**
- Recommends deciding between HTTP 429 and HTTP 503
- "Figure out what kinds of exceptions to show your users. In practice, you should decide if you want HTTP 429 (Too Many Requests) or HTTP 503 (Service Unavailable)"

**Assessment:**
Not a contradiction but a nuance:
- 429 is the standard for client-side rate limiting (you exceeded your quota)
- 503 is more appropriate for server-side overload (I'm too busy right now)
- NGINX defaults to 503 but makes it configurable

---

## Area 4: Base Delay Values for Retry Backoff

**Source A (AWS SDK Documentation):**
- Transient errors: 50 ms base delay (25 ms for DynamoDB)
- Throttling errors: 1000 ms base delay
- Max cap: 20 seconds

**Source B (Wikipedia - Exponential backoff):**
- SIP protocol: starts at 500ms (T1), doubles to 4s (T2)
- Ethernet: slot time of 51.2μs

**Assessment:**
Different protocols and systems use different base delays based on their specific characteristics:
- AWS SDKs: Optimized for cloud API latency profiles (50ms for network timeouts, 1000ms for rate limiting)
- SIP: Telephony round-trip time based (500ms T1 default)
- Ethernet: Physical layer collision detection timing (51.2μs slot time)
All follow the same exponential backoff principle with protocol-specific parameters.

---

## Area 5: Queue Metrics - What to Monitor

**Source A (Topic Specification):**
- Recommends monitoring "queue age" (oldest job wait time) over just queue length
- "10,000 jobs is not always bad. What's more useful is queue age."

**Source B (Google SRE Workbook):**
- Uses "utilization signals" (CPU, memory, executor load average)
- "As utilization approaches configured thresholds, we start rejecting requests based on their criticality"

**Source C (Stripe Blog):**
- Monitors worker utilization and fleet capacity
- "We track the number of workers with available capacity at all times"

**Assessment:**
Different metrics serve different purposes:
- Queue age: Best for detecting stuck/backlogged work
- Utilization: Best for detecting resource saturation
- Worker capacity: Best for load shedding decisions
No contradiction - these are complementary metrics for different aspects of system health.

---

## Area 6: NGINX Rate Limiting Algorithm

**Source A (NGINX Documentation):**
- Uses leaky bucket method for rate limiting

**Source B (Wikipedia Leaky Bucket):**
- NGINX uses leaky bucket as a meter

**Assessment:**
No contradiction - NGINX limit_req_module implements leaky bucket as a meter (checking conformance), not as a queue (enforcing). This aligns with NGINX's use of burst parameter for accepting temporary excess requests.

---

## Area 7: Redis Rate Limiting Support

**Source A (Redis Documentation):**
- Supports INCR/EXPIRE for fixed-window, sorted sets for sliding window, Lua scripts for token bucket

**Source B (Wikipedia Rate Limiting):**
- Lists fixed window counter, sliding window log, sliding window counter as algorithms

**Assessment:**
No contradiction - Redis documentation describes implementation mechanisms, Wikipedia describes algorithm categories. Redis supports all major algorithm variants through its data structures.

---

## Area 8: Retry Behavior During Service Disruption

**Source A (Google SRE Workbook):**
- "If a large subset of backend tasks in the datacenter are overloaded, requests should not be retried and errors should bubble up all the way to the caller."
- Uses "overloaded; don't retry" error for cascading overload situations

**Source B (AWS SDK Documentation):**
- Standard mode retries failed requests using exponential backoff with jitter
- Retry quota prevents unlimited retries during service disruptions

**Assessment:**
No contradiction - different layers of the system have different retry strategies:
- Google SRE: Server-side decision to prevent retry storms (higher layer)
- AWS SDK: Client-side retry quota with token bucket (lower layer)
Both serve to prevent retry storms but at different architectural levels.

---

## Area 9: Rate Limiting vs Backpressure Scope

**Source A (Google SRE Workbook):**
- "Rate limiting usually works at the entrance to a system. Backpressure is more general: When downstream is unable to keep up with work from upstream, upstream must slow down."

**Source B (Topic Specification):**
- Rate limiting at API entrance, backpressure in queue-based systems

**Assessment:**
No contradiction - topic specification's description aligns with Google SRE's distinction. Rate limiting is typically external traffic control, while backpressure is internal propagation of pressure signals.