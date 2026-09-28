# 06 — Research Gap Analysis

## Gap 1

Type:
WEAK_SOURCE / UNVERIFIED_CLAIM

Severity:
HIGH

Location:
`research/05-report.md`, Finding 6; `research/03-evidence.md`, Evidence 10–11

Problem:
The optimality claim and mathematical proofs of the XFetch algorithm (Vattani et al., PVLDB 2015) rely entirely on secondary sources (Wikipedia) and DOI landing page metadata. The primary paper PDF fetched from both `cseweb.ucsd.edu` and `vldb.org` returned binary streams that failed parser decompression. Primary proof details, exact variable definitions, and original empirical benchmark figures are NOT VERIFIED from the primary text.

Required Revision:
1. Re-fetch and parse the Vattani et al. 2015 PDF from an accessible plain-text/PDF extraction source.
2. Explicitly mark the XFetch optimality claim as "accepted on bibliographic authority, primary mathematical proof unverified by direct text extraction".

Can Be Approved Without Fix:
YES (with warning disclaimer in research report)

---

## Gap 2

Type:
CONTRADICTION / FORMULA_ERROR

Severity:
HIGH

Location:
`research/04-contradictions.md`, C1; `research/05-report.md`, Finding 6

Problem:
The lab specification formula `Δ · β · ln(rand()) > TTL_remaining` contains a sign error. For uniform random variable `rand() ∈ (0,1)`, `ln(rand())` is negative, making `Δ·β·ln(rand())` negative and the condition `negative > positive` ALWAYS FALSE. The Wikipedia/paper formulation is `(time() - delta * beta * log(rand(0,1))) ≥ expiry`, rearranged as `-Δ·β·ln(rand()) ≥ TTL_remaining`.

Required Revision:
Correct the formula in the research report and lab specification to `-Δ·β·ln(rand()) > TTL_remaining` (or `Δ·β·(-ln(rand())) > TTL_remaining`), explicitly warning engineers against using the unnegated formula.

Can Be Approved Without Fix:
NO (Fix required before engineering implementation)

---

## Gap 3

Type:
UNREACHABLE_SOURCE / MISSING_SOURCE

Severity:
MEDIUM

Location:
`research/02-sources.md`, Source 10; `research/06-open-questions.md`, Sec. 3.1

Problem:
All attempted official Redis / AWS ElastiCache documentation URLs (e.g., `redis.io/docs/latest/develop/use/cache`, `redis.io/docs/latest/develop/clients/patterns/cache-invalidation/`) returned HTTP 404/403 errors. No Redis-official pattern guidance, command specifics (`SET NX PX`), or cache-aside best practices were verified directly from Redis primary documentation.

Required Revision:
Update Redis doc references to current active URLs or explicitly mark Redis-specific pattern claims as relying on general vendor architecture guidelines (Microsoft Azure) rather than Redis-official docs.

Can Be Approved Without Fix:
YES (general cache-aside principles are verified via Microsoft Learn)

---

## Gap 4

Type:
OVERGENERALIZATION / SYNTHETIC_BENCHMARK

Severity:
LOW

Location:
`research/05-report.md`, Finding 4; `research/03-evidence.md`, Evidence 18

Problem:
Quantitative figures in the lab specification (10,000 RPS flash-sale scenario, 500 concurrent goroutines, specific P99 latency deltas, connection pool exhaustion thresholds) are synthetic exercise parameters, not empirical measurements from cited production research.

Required Revision:
Ensure all quantitative benchmark figures are explicitly labeled as "synthetic lab scenario parameters" rather than industry-standard benchmark metrics.

Can Be Approved Without Fix:
YES (Research report already correctly identifies Evidence 18 as synthetic parameters)

---

## Gap 5

Type:
SCOPE_ERROR / MISLEADING_STATUS

Severity:
LOW

Location:
`research/04-contradictions.md`, C3; `research/05-report.md`, Finding 7

Problem:
RFC 5861 is published under the Independent Submission stream as an Informational RFC and header text explicitly notes it is "not endorsed by the IETF and has no formal standing in the IETF standards process". Calling RFC 5861 an "Internet Standard" is inaccurate.

Required Revision:
Refer to RFC 5861 as an "IETF-published Informational RFC" or "informational specification", not an "Internet Standard".

Can Be Approved Without Fix:
YES (Research report already correctly calls it "RFC 5861 (2010, Informational)")

---

## Gap 6

Type:
UNVERIFIED_CLAIM / INFERENTIAL_LEAP

Severity:
MEDIUM

Location:
`research/05-report.md`, Finding 8; `research/03-evidence.md`, Evidence 16

Problem:
The claim that TTL jitter desynchronizes expiration clocks is well-supported (Wikipedia Thundering Herd §Mitigation for retry backoffs), but the claim that "jitter alone fails to prevent stampede on a single hot key" is an inferential leap by the agent rather than a quote from a primary caching source.

Required Revision:
Label the statement about jitter's insufficiency for single hot keys as an "inferential deduction based on single-key contention mechanics", rather than presenting it as an established primary source claim.

Can Be Approved Without Fix:
YES (The deduction is logically sound and mathematically trivial)
