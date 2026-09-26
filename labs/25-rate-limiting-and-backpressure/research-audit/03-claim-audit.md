# Claim Audit: Rate Limiting & Backpressure

Target Lab: `labs/25-rate-limiting-and-backpressure`  
Audit Scope: Major Research Claims Verification  
Audit Date: 2026-09-26  

---

## Claim 1

Claim: The token bucket algorithm allows burst traffic up to the bucket capacity while maintaining a long-term average rate limit.

Location: `research/03-evidence.md:3-15` & `research/05-report.md:26-44`

Evidence Provided: Direct quote from Wikipedia Token bucket page and corroborated by NGINX & Redis docs.

Source: Source 1 (Wikipedia: Token bucket), Source 5 (NGINX), Source 11 (Redis)

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Fully supported by all primary and secondary references.

---

## Claim 2

Claim: Adding jitter to exponential backoff significantly reduces the thundering herd problem and improves time to completion during contention.

Location: `research/03-evidence.md:16-28` & `research/05-report.md:67-86`

Evidence Provided: AWS Architecture blog benchmark data showing >50% call count reduction with full jitter.

Source: Source 8 (AWS Architecture Blog), Source 9 (AWS SDK Documentation)

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Supported by AWS benchmarks and queuing theory literature. In `research/05-report.md`, default values (50ms/1000ms/20s) were properly contextualized as AWS SDK defaults rather than universal rules.

---

## Claim 3

Claim: HTTP 429 Too Many Requests is the standard status code for rate limiting and MUST NOT be cached by intermediaries.

Location: `research/03-evidence.md:29-41` & `research/05-report.md:89-106`

Evidence Provided: RFC 6585 Section 4 text on 429 status code definition.

Source: Source 6 (RFC 6585)

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Standard RFC requirement. Report correctly notes 503 as alternative for backend server overload vs 429 for client quota excess.

---

## Claim 4

Claim: Little's Law ($L = \lambda W$) applies to queueing systems, enabling mathematical calculation of queue age and backlog time ($5,000,000 / 2,000 = 2,500\text{s} \approx 41\text{m } 40\text{s}$).

Location: `research/03-evidence.md:42-54` & `research/05-report.md:109-126`

Evidence Provided: Definition of Little's Law formula and verification of backlog processing time.

Source: Source 4 (Wikipedia: Little's Law), Source 10 (Google SRE Workbook)

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Mathematically exact and verified.

---

## Claim 5

Claim: NGINX implements rate limiting using the leaky bucket algorithm via `ngx_http_limit_req_module`.

Location: `research/03-evidence.md:55-67` & `research/05-report.md:33-35`

Evidence Provided: Official NGINX documentation stating limitation is done using leaky bucket method.

Source: Source 5 (NGINX Documentation)

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Directly supported by vendor documentation.

---

## Claim 6

Claim: Stripe operates 4 different types of limiters: Request rate limiter, Concurrent requests limiter, Fleet usage load shedder, and Worker utilization load shedder.

Location: `research/03-evidence.md:68-80` & `research/05-report.md:129-150`

Evidence Provided: Quotation from Stripe Engineering Blog explaining production load-shedding architecture.

Source: Source 7 (Stripe Engineering Blog)

Source Actually Supports Claim: YES

Classification: EXAMPLE / IMPLEMENTATION-SPECIFIC

Severity: LOW

Notes: Classified as implementation example. Accurately reflects Stripe's engineering practice.

---

## Claim 7

Claim: Cost-based rate limiting throttles requests based on estimated resource/computational cost rather than raw request counts.

Location: `research/03-evidence.md:173-184` & `research/06-open-questions.md:7-8`

Evidence Provided: Referenced Stripe's Concurrent Requests Limiter and API Gateway cost-based policies.

Source: Source 7 (Stripe Engineering Blog)

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Cites external primary source (Stripe blog) after revision replaced internal specification citation.

---

## Claim 8

Claim: Multi-tenant systems require per-tenant rate limiting / fair queueing to prevent noisy neighbors from starving other tenants.

Location: `research/03-evidence.md:185-196` & `research/06-open-questions.md:5-6`

Evidence Provided: Stripe load shedders and AWS Well-Architected multi-tenancy isolation guidance.

Source: Source 7 (Stripe Blog), Source 8 (AWS Architecture)

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Supported by cloud reference architectures and multi-tenant design patterns.

---

## Claim 9

Claim: Redis is well-suited for distributed rate limiting due to atomic operations (INCR, EXPIRE, Lua scripts).

Location: `research/03-evidence.md:107-119` & `research/05-report.md:153-170`

Evidence Provided: Redis official rate limiter pattern documentation.

Source: Source 11 (Redis Documentation)

Source Actually Supports Claim: YES

Classification: FACT / EXAMPLE

Severity: LOW

Notes: Supported by official Redis documentation.

---

## Unsupported Claims Analysis

No major claims in the research report or evidence file are unsupported. All 9 key technical claims have direct external citations from standards bodies (IETF), academic/reference works (Wikipedia), vendor docs (NGINX, Redis, AWS, RabbitMQ), or engineering blogs (Stripe, AWS Architecture).

Prior circular citation issues (Evidences 14 & 15 citing topic prompt) were resolved in `research-revision/`.
