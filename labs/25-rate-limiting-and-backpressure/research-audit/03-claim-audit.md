# Claim Audit

## Claim 1

Claim:
The Token Bucket algorithm allows burst traffic up to the bucket capacity while maintaining a long-term average rate limit. Leaky Bucket as a meter is mathematically equivalent.

Location:
`research/03-evidence.md`: Evidence 1 & 2; `research/05-report.md`: Finding 1

Evidence Provided:
Wikipedia Token Bucket and Leaky Bucket definitions; NGINX limit_req module leaky bucket implementation.

Source:
Source 1 (Wikipedia Token bucket), Source 5 (NGINX docs), Source 12 (Wikipedia Leaky bucket)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Accurately reflects the dual-viewpoint relationship between token bucket and leaky bucket as a meter.

---

## Claim 2

Claim:
Rate limiting typically operates at the entry point to a system (external traffic control), whereas backpressure is a broader feedback mechanism that propagates pressure signals backward when downstream components are overwhelmed.

Location:
`research/03-evidence.md`: Evidence 9; `research/05-report.md`: Finding 2

Evidence Provided:
Google SRE Chapter 21: "Rate limiting usually works at the entrance to a system. Backpressure is more general: When downstream is unable to keep up with work from upstream, upstream must slow down."

Source:
Source 10 (Google SRE: Handling Overload)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Distinguishes perimeter traffic policing from end-to-end downstream/upstream flow control.

---

## Claim 3

Claim:
Adding full jitter to exponential backoff reduces the thundering herd problem and cuts call contention by more than half for 100 contending clients. AWS SDK standardized the formula: `delay = random(0, 1) * min(20000 ms, base_delay * 2^retry)` with 50 ms base for transient errors (25 ms DynamoDB) and 1,000 ms base for throttling errors.

Location:
`research/03-evidence.md`: Evidence 6 & 7; `research/05-report.md`: Finding 3

Evidence Provided:
AWS Architecture blog simulation results; AWS SDK Retry Behavior specification table and formulas.

Source:
Source 8 (AWS Architecture Blog), Source 9 (AWS SDK Retry Behavior)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Verified directly against AWS SDK documentation and Marc Brooker's architectural research.

---

## Claim 4

Claim:
HTTP 429 Too Many Requests is the standard status code indicating rate limiting per RFC 6585, responses MUST NOT be stored by caches, and SHOULD include a Retry-After header.

Location:
`research/03-evidence.md`: Evidence 4; `research/05-report.md`: Finding 4

Evidence Provided:
RFC 6585 Section 4 text on HTTP 429 semantics and cache restrictions.

Source:
Source 6 (RFC 6585)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Accurate standard citation.

---

## Claim 5

Claim:
Little's Law ($L = \lambda W$) provides the mathematical foundation for queue sizing and backlog growth. A queue backlog of 5,000,000 items processing at 2,000 items/sec requires 2,500 seconds ($\approx 41$ minutes 40 seconds) to clear.

Location:
`research/03-evidence.md`: Evidence 10 & 11; `research/05-report.md`: Finding 5

Evidence Provided:
Little's Law definition ($L = \lambda W$) and arithmetic calculation ($5,000,000 / 2,000 = 2,500\text{ s} = 41\text{m } 40\text{s}$).

Source:
Source 3 (Wikipedia Little's Law)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Calculation is mathematically verified and correctly assumes steady-state departure rate with no net incoming rate during drainage.

---

## Claim 6

Claim:
Stripe operates four distinct limiter/load-shedder layers in production: Request Rate Limiter (per-user RPS), Concurrent Requests Limiter (in-flight CPU limits), Fleet Usage Load Shedder (reserves capacity for critical traffic), and Worker Utilization Load Shedder (sheds traffic in stages: test mode, GETs, POSTs).

Location:
`research/03-evidence.md`: Evidence 5; `research/05-report.md`: Finding 6

Evidence Provided:
Stripe engineering blog post detailing their 4 limiters and production rationale.

Source:
Source 7 (Stripe Engineering Blog)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Accurately cites Stripe's production architecture.

---

## Claim 7

Claim:
Redis provides atomic primitives (`INCR`, `EXPIRE`, sorted sets, Lua scripts) suitable for implementing fixed-window, sliding-window, and token bucket distributed rate limiters with sub-millisecond overhead.

Location:
`research/03-evidence.md`: Evidence 12; `research/05-report.md`: Finding 7

Evidence Provided:
Redis documentation use cases and Stripe architecture blog.

Source:
Source 7 (Stripe Blog), Source 11 (Redis Docs)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Supported by official Redis documentation and industry implementations.

---

## Claim 8

Claim:
Multi-tenant systems require per-customer resource quotas (e.g. CPU seconds/sec) and request criticality classification (`CRITICAL_PLUS`, `CRITICAL`, `SHEDDABLE_PLUS`, `SHEDDABLE`) to prevent noisy neighbors and enable graceful tiered load shedding.

Location:
`research/03-evidence.md`: Evidence 13 & 15; `research/05-report.md`: Finding 8

Evidence Provided:
Google SRE Book Chapter 21 sections on "Per-Customer Limits" and "Criticality".

Source:
Source 10 (Google SRE: Handling Overload)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Faithful representation of Google SRE overload management practices.

---

## Claim 9

Claim:
Client-side adaptive throttling and retry budgets (per-request budget of max 3 attempts, per-client retry budget ratio $< 10\%$, and token bucket retry quotas) are necessary to prevent cascading retry storms.

Location:
`research/03-evidence.md`: Evidence 8, 14, 17, 18; `research/05-report.md`: Finding 9

Evidence Provided:
Google SRE Chapter 21 client-side throttling math ($P_{\text{drop}} = \max(0, \frac{\text{requests} - K \cdot \text{accepts}}{\text{requests} + 1})$), Google retry budget rules, and AWS SDK token bucket retry quota (500 tokens).

Source:
Source 9 (AWS SDK Retry Behavior), Source 10 (Google SRE: Handling Overload)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
All mechanisms and numbers correspond directly to published authoritative documentation.
