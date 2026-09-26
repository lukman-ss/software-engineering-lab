# Claim Audit

## Claim 1

Claim:
The Token Bucket algorithm allows burst traffic up to the bucket capacity while maintaining a long-term average rate limit.

Location:
`research/03-evidence.md:3-15`, `research/05-report.md:26-44`

Evidence Provided:
Wikipedia: Token Bucket article quote explaining conforming flows and burstiness governed by bucket depth.

Source:
Source 1 (Wikipedia: Token bucket), corroborated by NGINX & Redis documentation.

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Accurate description of Token Bucket mechanics.

---

## Claim 2

Claim:
Adding jitter to exponential backoff significantly reduces the thundering herd problem and improves system recovery during contention.

Location:
`research/03-evidence.md:16-28`, `research/05-report.md:67-86`

Evidence Provided:
AWS Architecture Blog quote showing 100 contending clients call count reduced by >50% and completion time improved.

Source:
Source 8 (AWS Architecture Blog: Exponential Backoff And Jitter).

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Well supported by empirical benchmarks from AWS Architecture.

---

## Claim 3

Claim:
HTTP 429 Too Many Requests is the standard status code for indicating rate limiting has been applied, and responses MUST NOT be cached by default.

Location:
`research/03-evidence.md:29-41`, `research/05-report.md:89-106`

Evidence Provided:
RFC 6585 Section 4 text and definitions.

Source:
Source 6 (RFC 6585).

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Authoritative standard for HTTP rate limiting semantics.

---

## Claim 4

Claim:
Little's Law ($L = \lambda W$) provides the fundamental relationship for understanding queue behavior, and queue backlog increases when $\lambda_{arrival} > \mu_{processing}$.

Location:
`research/03-evidence.md:42-54`, `research/05-report.md:109-127`

Evidence Provided:
Wikipedia Little's Law definition; calculation verifying $5,000,000 / 2,000 = 2,500\text{s} \approx 41\text{m}40\text{s}$.

Source:
Source 4 (Wikipedia: Little's law).

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Calculation is mathematically accurate.

---

## Claim 5

Claim:
NGINX implements rate limiting using the leaky bucket algorithm via `ngx_http_limit_req_module`.

Location:
`research/03-evidence.md:55-67`, `research/05-report.md:34`

Evidence Provided:
NGINX official documentation quote specifying the leaky bucket method.

Source:
Source 5 (NGINX limit_req_module docs).

Source Actually Supports Claim:
YES

Classification:
IMPLEMENTATION-SPECIFIC

Severity:
LOW

Notes:
Accurate documentation excerpt.

---

## Claim 6

Claim:
Stripe operates four distinct layers of limiters in production: Request rate limiter, Concurrent requests limiter, Fleet usage load shedder, and Worker utilization load shedder.

Location:
`research/03-evidence.md:68-80`, `research/05-report.md:129-151`

Evidence Provided:
Stripe Engineering Blog post directly outlining the 4 limiter designs.

Source:
Source 7 (Stripe: Scaling your API with rate limiters).

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Directly supported by Stripe engineering literature.

---

## Claim 7

Claim:
AWS SDKs implement a retry quota using a token bucket to prevent retry storms during service disruptions.

Location:
`research/03-evidence.md:81-93`, `research/05-report.md:81-86`

Evidence Provided:
AWS SDK Developer Guide excerpt detailing the retry quota token bucket mechanism.

Source:
Source 9 (AWS SDK Retry Behavior).

Source Actually Supports Claim:
YES

Classification:
IMPLEMENTATION-SPECIFIC

Severity:
LOW

Notes:
Valid implementation reference from AWS SDK standard mode.

---

## Claim 8

Claim:
Rate limiting typically works at the entry point of a system, whereas backpressure is a broader mechanism propagating downstream strain backward to upstream producers.

Location:
`research/03-evidence.md:120-132`, `research/05-report.md:47-64`

Evidence Provided:
Google SRE Workbook / SRE Book Chapter 22 handling overload principles.

Source:
Source 10 (Google SRE Book).

Source Actually Supports Claim:
YES

Classification:
INTERPRETATION

Severity:
LOW

Notes:
Fundamental architectural distinction clearly formulated.

---

## Claim 9

Claim:
AWS SDKs compute retry delay using `delay = random(0, 1) * min(20000 ms, base_delay * 2^retry)` with 50ms (transient) and 1000ms (throttling) base delays.

Location:
`research/03-evidence.md:133-144`, `research/04-contradictions.md:65-69`, `research/05-report.md:81-85`

Evidence Provided:
AWS SDK documentation formula and parameters.

Source:
Source 9 (AWS SDK Retry Behavior) & Source 8 (AWS Architecture Blog).

Source Actually Supports Claim:
YES

Classification:
IMPLEMENTATION-SPECIFIC

Severity:
MEDIUM

Notes:
Specific to AWS SDK v3/standard retry mode. Must not be generalized as universal constant across all ecosystems.

---

## Claim 10

Claim:
Multi-tenant systems require per-tenant rate limiting / isolation to prevent noisy neighbors from monopolizing worker capacity.

Location:
`research/03-evidence.md:185-196`, `research/05-report.md:209-210`

Evidence Provided:
Topic specification example of Tenant A occupying workers while Tenant B starves.

Source:
Original topic specification (cited as evidence in Evidence 15).

Source Actually Supports Claim:
PARTIAL

Classification:
EXAMPLE

Severity:
MEDIUM

Notes:
Claim is conceptually sound, but citing the prompt/topic specification itself as "Evidence Source" is circular. Should cite literature on multi-tenant fair queueing / noisy neighbor isolation (e.g. Stripe load shedder or cloud tenancy papers).

---

## Claim 11

Claim:
Effective rate limiting should consider computational cost of operations rather than raw request count alone.

Location:
`research/03-evidence.md:172-184`

Evidence Provided:
Topic specification statement; corroborated by Stripe concurrent request limiter and NGINX URI-specific limits.

Source:
Original topic specification (Evidence 14).

Source Actually Supports Claim:
PARTIAL

Classification:
INTERPRETATION

Severity:
LOW

Notes:
Supported in industry practice (e.g. GraphQL query complexity rating, Stripe concurrent limiters for heavy endpoints), but evidence citation directly references the prompt spec.
