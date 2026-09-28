# 03 — Claim Audit

## Claim 1

Claim:
Cache-Aside (lazy loading) loads data into the cache on demand: check cache, on miss fetch from store, then populate cache; writes update the store then invalidate the cache entry.

Location:
research/03-evidence.md, Evidence 1 and Evidence 6

Evidence Provided:
Microsoft Azure Architecture Center (Source 01) verbatim. Steps 1–3 documented; write-invalidation order (store first, then delete) stated explicitly.

Source:
Source 01 — Microsoft Learn, Cache-Aside Pattern

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW (well-supported)

Notes:
VERIFIED directly from Microsoft Learn page. The exact write order note ("update store before removing cache item") is present in the source.

---

## Claim 2

Claim:
Write-through writes synchronously to cache and backing store; write-back (write-behind) writes to cache first and postpones store write until eviction (or explicit flush).

Location:
research/03-evidence.md, Evidence 3; research/05-report.md, Finding 2–3

Evidence Provided:
Wikipedia Cache (computing) write policies section (Source 06). Hardware textbook Hennessy & Patterson cited by Wikipedia but not independently opened.

Source:
Source 06 — Wikipedia Cache (computing)

Source Actually Supports Claim:
YES

Classification:
FACT (with note: hardware cache semantics mapped to application-layer)

Severity:
MEDIUM

Notes:
The mapping from hardware write-back to application-layer write-behind is conceptually sound but carries an implied scope extension. The research acknowledges this in Evidence 3 Notes. No application-specific primary source verified.

---

## Claim 3

Claim:
Write-through provides immediate read-after-write freshness because it updates cache and store "in the same write operation."

Location:
research/03-evidence.md, Evidence 4; research/05-report.md, Finding 2; research/04-contradictions.md, C6

Evidence Provided:
Microsoft Learn Cache-Aside Pattern (Source 01) explicitly distinguishes write-through as updating "data store and the cache in the same write operation so that readers see the new value immediately after a successful write."

Source:
Source 01 — Microsoft Learn

Source Actually Supports Claim:
YES

Classification:
FACT (with scope limitation: "same write operation" = application-level sequential; NOT distributed atomic transaction)

Severity:
MEDIUM

Notes:
Research correctly identifies scope limitation in C6. Critical nuance: no distributed ACID guarantee. Claim is accurate if "same operation" is read correctly. Research handles this well.

---

## Claim 4

Claim:
Write-behind (write-back) at application layer risks data loss on cache-node crash before flush to the backing store.

Location:
research/05-report.md, Finding 3

Evidence Provided:
Wikipedia Cache (computing) (Source 06) on hardware write-back semantics. Lab spec text cited as additional support.

Source:
Source 06 — Wikipedia

Source Actually Supports Claim:
PARTIAL

Classification:
INTERPRETATION

Severity:
MEDIUM

Notes:
Data-loss risk is a logical inference from the postponement semantics, not a directly quoted application-layer claim. Research acknowledges: "crash risk is an inference from postponement." Claim is reasonable and industry-accepted, but lacks a primary application-layer source. MEDIUM risk of overgeneralization if presented as definitive.

---

## Claim 5

Claim:
Cache stampede is a cascading failure where expiry of a popular key causes many threads to simultaneously recompute, exhausting shared resources; under very heavy load it can reduce hit rate to zero (congestion collapse).

Location:
research/03-evidence.md, Evidence 5; research/05-report.md, Finding 4

Evidence Provided:
Wikipedia Cache stampede (Source 04). Example: 10 req/s × 3 s recomputation = 30 concurrent rebuilds at expiry.

Source:
Source 04 — Wikipedia Cache stampede

Source Actually Supports Claim:
YES

Classification:
FACT (mechanism) / EXAMPLE (quantitative illustration)

Severity:
LOW

Notes:
Mechanism well-documented. Numeric example (30 concurrent recomputes) is arithmetic/illustrative, not empirical measurement. Research correctly notes this.

---

## Claim 6

Claim:
The XFetch algorithm (Vattani et al. PVLDB 2015) is proven optimal for probabilistic early cache-expiration prevention.

Location:
research/03-evidence.md, Evidence 10–11; research/05-report.md, Finding 6; research/06-open-questions.md, OQ-1

Evidence Provided:
Wikipedia Cache stampede cites the 2015 paper. DOI landing confirms bibliographic existence. Primary PDF could NOT be parsed.

Source:
Source 08 (DOI) + Source 04 (Wikipedia)

Source Actually Supports Claim:
PARTIAL

Classification:
HYPOTHESIS (optimality) / FACT (algorithm description from Wikipedia)

Severity:
HIGH

Notes:
"Proven optimal" is a strong claim. Primary paper proofs cannot be independently verified due to PDF parse failure. Wikipedia states "the authors show that... optimal in terms of effectiveness and how early recomputations can happen" — this is the Wikipedia's characterization of the paper's claim, not independently verified. Treating optimality as established fact is overstated.

---

## Claim 7

Claim:
The XFetch formula is: `value, delta, expiry ← cache_read(key); if (!value || (time() - delta * beta * log(rand(0,1))) ≥ expiry): recompute`.

Location:
research/03-evidence.md, Evidence 10; research/04-contradictions.md, C1

Evidence Provided:
Wikipedia Cache stampede pseudocode. Paper PDF not parsed.

Source:
Source 04 — Wikipedia (citing Source 08/09)

Source Actually Supports Claim:
YES (for the Wikipedia-transcribed formula)

Classification:
FACT (Wikipedia-sourced formula, unverified against primary paper)

Severity:
MEDIUM

Notes:
Formula comes from Wikipedia, not verified from primary paper. Research notes this. Lab formula (`Δ·β·ln(rand()) > TTL_remaining` without the negation) is identified as likely having a sign error in Contradiction C1. The bug identification is correct and critical for implementation.

---

## Claim 8

Claim:
Lab formula `Δ · β · ln(rand()) > TTL_remaining` is a rearrangement of the Wikipedia/paper condition.

Location:
research/04-contradictions.md, C1; research/05-report.md, Finding 6

Evidence Provided:
Mathematical analysis by Research Agent comparing both forms.

Source:
Research Agent analysis

Source Actually Supports Claim:
YES (analysis correct)

Classification:
INTERPRETATION

Severity:
HIGH

Notes:
For rand∈(0,1), ln(rand) is negative; therefore `Δ·β·ln(rand())` is negative, and the inequality `Δ·β·ln(rand()) > TTL_remaining` is always false (never triggers early refresh). This is mathematically correct. Research identified and documented this correctly. Severity HIGH because this is a functional formula error that would break any implementation based on the lab spec's formula verbatim.

---

## Claim 9

Claim:
`stale-while-revalidate` (RFC 5861 §3) allows caches to serve stale responses for up to N seconds while revalidating asynchronously in the background; validation SHOULD be request-triggered.

Location:
research/03-evidence.md, Evidence 12; research/05-report.md, Finding 7

Evidence Provided:
IETF RFC 5861 §3, §3.1, §5 (directly fetched and verified by Auditor).

Source:
Source 02 — RFC 5861

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Fully verified from primary RFC text. "SHOULD attempt to revalidate it while still serving stale responses (i.e., without blocking)" is confirmed. The security note about request-triggered revalidation to avoid amplification is confirmed in §5.

---

## Claim 10

Claim:
TTL should not be too short (continuous refetches) or too long (staleness); optimal TTL depends on workload; no universal number.

Location:
research/03-evidence.md, Evidence 15

Evidence Provided:
Microsoft Cache-Aside Pattern (Source 01).

Source:
Source 01 — Microsoft Learn

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Directly supported. The research note that "Lab 5-minute TTL + 15s jitter is a local convention, not a standard" is appropriate.

---

## Claim 11

Claim:
Jitter desynchronizes expiry timestamps but does NOT prevent stampede on a single hot key.

Location:
research/03-evidence.md, Evidence 16; research/05-report.md, Finding 8

Evidence Provided:
Wikipedia Thundering herd §Mitigation (jitter for retry backoff context). Local code reference.

Source:
Source 05 — Wikipedia Thundering herd

Source Actually Supports Claim:
PARTIAL

Classification:
INTERPRETATION

Severity:
MEDIUM

Notes:
The source describes jitter in the context of retry backoff desynchronization (OS/network level). The extension to cache-stampede prevention is inferential. The specific claim that "jitter does not cap concurrent rebuilds of one hot key" is a correct inference but is not backed by an opened primary source making this exact statement.

---

## Claim 12

Claim:
Lab scenario (10,000 RPS flash-sale, 500 concurrent goroutines, P99 deltas, connection pool exhaustion) are validated production benchmarks.

Location:
research/05-report.md §Executive Summary; research/03-evidence.md, Evidence 18

Evidence Provided:
Evidence 18 explicitly states: "No opened source reports 10,000 concurrent misses on one key, 500-goroutine lab numbers, or P99 deltas."

Source:
Evidence 18 (self-documenting gap)

Source Actually Supports Claim:
YES (claim is that these are NOT supported; this is accurately stated)

Classification:
EXAMPLE (synthetic/hypothetical)

Severity:
LOW

Notes:
Research accurately documents numbers are synthetic lab parameters, not production data. This is correctly handled.

---

## Claim 13

Claim:
Go `singleflight.Group.Do` guarantees only one in-flight execution per key; duplicates wait and receive the same result; `shared` flag is set to true for shared results.

Location:
research/03-evidence.md, Evidence 9; research/05-report.md, Finding 5

Evidence Provided:
pkg.go.dev documentation (Source 03) directly fetched and verified.

Source:
Source 03 — Go pkg.go.dev

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Fully verified from primary source. Exact wording from source matches claim. `Forget(key)` behavior also accurately documented.

---

## Claim 14

Claim:
In-process `singleflight` does not prevent stampede in multi-process (N replicas) deployments without a distributed lock analogue.

Location:
research/03-evidence.md, Evidence 9 Notes; research/05-report.md, Conclusion

Evidence Provided:
Go singleflight docs (process-local). Redis distributed lock behavior NOT VERIFIED from redis.io.

Source:
Source 03 — Go pkg.go.dev

Source Actually Supports Claim:
YES

Classification:
FACT (scope boundary)

Severity:
LOW

Notes:
Correctly inferred from "process-local" nature of singleflight. Accurately stated. Redis SET NX PX as distributed lock alternative is mentioned but NOT VERIFIED.
