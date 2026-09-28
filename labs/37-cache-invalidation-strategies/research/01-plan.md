# 01 — Research Plan

## Research Topic
Cache Invalidation Strategies — Memilih Pola Caching dan Mencegah Thundering Herd (Lab 37, lukman-ss)

## Objective
Investigate evidence for: (a) base caching patterns Cache-Aside / Write-Through / Write-Behind — mechanisms, advantages, disadvantages; (b) Cache Stampede / Thundering Herd — definition, trigger conditions, quantified impact; (c) mitigations — Single-Flight / Distributed Lock, Probabilistic Early Expiration / XFetch, Stale-While-Revalidate; (d) supporting concepts — TTL, jitter, Redis — to enable a senior-level practical lab exercise (500 concurrent goroutines on expired key, single-flight protection, P99 / DB-connection comparison).

Research-only. No code, no tests, no publication content.

## Research Questions

| ID | Question | Type |
|---|---|---|
| RQ-01 | What are the mechanisms, benefits, and costs of Cache-Aside (lazy loading), Write-Through, and Write-Behind (write-back)? When is each appropriate? | Descriptive / normative |
| RQ-02 | What is Cache Stampede / Thundering Herd formally? Under what concurrency / TTL-expiry / cost conditions does it occur? What is the failure mode (DB overload, connection exhaustion, congestion collapse)? | Causal / definitional |
| RQ-03 | How does Single-Flight / distributed-lock (request coalescing) prevent stampede? What are implementation trade-offs (extra writes, lock TTL, failure handling, stale fallback)? | Mechanism / trade-off |
| RQ-04 | How does Probabilistic Early Expiration / XFetch work formally (formula, beta, delta, TTL_remaining)? What optimality claim is made and under what model? What are memory/overhead costs? | Algorithmic / formal |
| RQ-05 | How does Stale-While-Revalidate work (RFC 5861 semantics, HTTP and application-cache analogues, background refresh)? What freshness/availability trade-off does it imply? | Mechanism / standard |
| RQ-06 | What role do TTL, jitter, and Redis/in-memory stores play in these patterns? Does jitter alone prevent stampede or only desynchronize? | Supporting |
| RQ-07 | What quantitative evidence exists for stampede impact and mitigation effectiveness (DB hit multiplication, P99 latency, connection load — before/after)? | Empirical |

## Search Strategy

1. Start from Tier-1 seeds: Microsoft Azure Architecture Center (cache-aside pattern), IETF RFC 5861 (stale-while-revalidate), Go pkg.go.dev singleflight, VLDB proceedings (Vattani et al. 2015 via DOI), Wikipedia pages that cite primary sources (cache stampede, thundering herd, cache invalidation, cache computing).
2. For each claim, open the actual source (not snippet). Record URL, publisher, published date, accessed date, tier.
3. Prefer authoritative docs over blogs. Treat community sources (Wikipedia, StackOverflow) as Tier-3 discovery; require Tier-1/2 corroboration for HIGH confidence.
4. Cross-check significant claims with ≥2 independent sources; record disagreements explicitly.
5. For XFetch paper, DOI + Wikipedia citation + attempted PDF fetch; note PDF parsing failure as limitation.
6. Search variants for Redis/official caching docs (redis.io, AWS ElastiCache, GCP) — record NOT VERIFIED where fetch fails rather than inventing.
7. Current research date: 2026-09-28 (UTC). Record source publication dates where available.

## Expected Primary Sources

- Microsoft Learn — Cache-Aside Pattern (Azure Architecture Center)
- IETF RFC 5861 — HTTP Cache-Control Extensions for Stale Content (M. Nottingham, 2010)
- Vattani, Chierichetti, Lowenstein — "Optimal Probabilistic Cache Stampede Prevention", Proc. VLDB Endow. 8(8), 2015
- Go — golang.org/x/sync/singleflight package documentation
- Wikipedia — Cache stampede; Thundering herd problem; Cache (computing) write policies; Cache invalidation (all as Tier-3, citing Tier-1 references therein)
- Redis / AWS / GCP docs (attempted; outcome recorded)

## Risks / Unknowns

- PDF for VLDB 2015 paper could not be parsed (compressed stream); metadata via DOI and Wikipedia citation only — full proof details remain NOT VERIFIED from primary PDF.
- Redis official pattern docs returned 404 at attempted URLs (docs structure changed) — Redis-specific claims will be marked NOT VERIFIED or sourced via secondary Tier-2 if found.
- Lab spec formula `Δ·β·ln(rand()) > TTL_remaining` differs in notation from Wikipedia's `(time() - delta*beta*log(rand)) ≥ expiry`; verify equivalence during evidence extraction.
- Quantitative before/after P99 / connection-pool numbers in lab exercise are hypothetical; real benchmarks are workload-dependent — risk of overgeneralizing.
- Thundering Herd vs Cache Stampede terminology overlap — sources use interchangeably but OS-level thundering herd (accept queue) differs from cache stampede (recomputation) — need to disambiguate.
- Jitter vs probabilistic early expiration conflation risk — both randomize but at different times (TTL assignment vs recomputation decision).
