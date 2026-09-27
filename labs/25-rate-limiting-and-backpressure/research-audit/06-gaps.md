# Research Gap Analysis

## Gap 1

Type:
WEAK_SOURCE

Severity:
LOW

Location:
`research/02-sources.md`: Sources 1, 2, 3, 4, 12

Problem:
Multiple algorithmic foundations (Token Bucket, Leaky Bucket, Exponential Backoff, Little's Law) rely on Wikipedia pages rather than original academic papers or network RFCs. While Wikipedia articles here are mathematically accurate and link to primary works (e.g. John Little 1961, Turner 1986), citing the primary literature directly strengthens academic rigor.

Required Revision:
Include direct citations to John D.C. Little (1961) for Little's Law and Jonathan Turner (1986) for Token Bucket in future reference iterations.

Can Be Approved Without Fix:
YES

---

## Gap 2

Type:
MISSING_CASE

Severity:
LOW

Location:
`research/05-report.md`: Finding 7 (Distributed Rate Limiting)

Problem:
Redis-based distributed rate limiting documentation notes Lua scripts and sorted sets, but does not deeply analyze race conditions, Redis cluster multi-key partitioning limitations (hash tags `{user_id}` requirement), or local memory cache fallback when Redis is unreachable (fail-open vs fail-closed strategies).

Required Revision:
Document fail-open vs. fail-closed trade-offs for distributed limiters (noted in Stripe blog: catch exceptions and fail-open to preserve API availability).

Can Be Approved Without Fix:
YES

---

## Gap 3

Type:
OVERGENERALIZATION

Severity:
LOW

Location:
`research/05-report.md`: Finding 3 & 9

Problem:
AWS SDK's full jitter parameters (50 ms base, 1,000 ms throttling base, 20 s cap) and Google SRE's retry budget ratio (< 10%) are specific engineering choices calibrated for their infrastructure scale, not universal constants for every application domain.

Required Revision:
Clearly label AWS SDK parameters and Google SRE retry budgets as reference production configurations rather than universal rules.

Can Be Approved Without Fix:
YES
