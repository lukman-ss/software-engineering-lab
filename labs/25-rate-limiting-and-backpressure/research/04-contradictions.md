# Contradictions and Disagreements

## No Material Contradictions Discovered

The research sources generally agree on the fundamental concepts of rate limiting and backpressure. However, there are some areas of nuance and implementation differences:

---

## Area 1: Token Bucket vs Leaky Bucket

**Source A (Token Bucket Wikipedia, Redis Documentation):**
- Token bucket adds tokens at a fixed rate, removes tokens per request
- Allows bursts up to bucket capacity
- Favored for API rate limiting (Stripe, AWS API Gateway)

**Source B (Leaky Bucket Wikipedia, NGINX Documentation):**
- Leaky bucket processes at fixed rate, queues/drops excess
- Two variants: as meter (checking) and as queue (enforcing)
- NGINX uses leaky bucket as meter for rate limiting

**Assessment:**
These are mathematically equivalent (mirror images). The difference is primarily in implementation perspective:
- Token bucket: "Do I have tokens?" → accept/reject
- Leaky bucket: "Is there capacity in the queue?" → accept/delay/reject
Both can achieve the same rate limiting behavior with appropriate parameters.

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
- "Figure out what kinds of exceptions to show your users. In practice, you should decide if you want HTTP 429 (Too Many Requests) or HTTP 503 (Service Unavailable) and what is the most accurate depending on the situation."

**Assessment:**
Not a contradiction but a nuance:
- 429 is the standard for client-side rate limiting (you exceeded your quota)
- 503 is more appropriate for server-side overload (I'm too busy right now)
- Stripe correctly notes the choice depends on the specific scenario

---

## Area 4: Exponential Backoff Base Delay Values

**Source A (AWS SDK Documentation):**
- Transient errors: 50 ms base delay
- Throttling errors: 1,000 ms base delay
- Max cap: 20 seconds

**Source B (Wikipedia - Exponential Backoff):**
- SIP protocol: starts at 500ms (T1), doubles to 4s (T2)
- Ethernet: slot time of 51.2μs

**Assessment:**
Different protocols and systems use different base delays based on their specific characteristics:
- AWS SDKs: Optimized for cloud API latency profiles
- SIP: Telephony round-trip time based
- Ethernet: Physical layer collision detection timing
All follow the same exponential backoff principle with protocol-specific parameters.

---

## Area 5: Queue Metrics - What to Monitor

**Source A (Original Topic Specification):**
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