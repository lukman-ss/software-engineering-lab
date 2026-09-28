# 02 — Sources

## Source 01

Title: Cache-Aside Pattern
Publisher: Microsoft Learn — Azure Architecture Center
URL: https://learn.microsoft.com/en-us/azure/architecture/patterns/cache-aside
Published: 2025-09-11 (ms.date); updated 2025-12-09
Accessed: 2026-09-28
Source Tier: 1 — official vendor architecture documentation
Relevance: Definitive description of Cache-Aside (lazy loading) mechanism, consistency, TTL, write-invalidation, when to use vs write-through.

## Source 02

Title: RFC 5861 — HTTP Cache-Control Extensions for Stale Content
Publisher: IETF (Independent Submission) / RFC Editor — Author M. Nottingham (Yahoo!)
URL: https://datatracker.ietf.org/doc/html/rfc5861
Published: May 2010
Accessed: 2026-09-28
Source Tier: 1 — standards-track-adjacent informational RFC
Relevance: Normative definition of `stale-while-revalidate` (Sec. 3) and `stale-if-error` (Sec. 4) — basis for Stale-While-Revalidate mitigation.

## Source 03

Title: singleflight package — golang.org/x/sync/singleflight
Publisher: Go Project (pkg.go.dev)
URL: https://pkg.go.dev/golang.org/x/sync/singleflight
Published: v0.23.0 published 2026-08-31 (page viewed); package stable since 2013
Accessed: 2026-09-28
Source Tier: 1 — official language/library documentation
Relevance: Single-Flight / request coalescing primitive used in Go to prevent duplicate in-flight recomputations — maps to distributed-lock / single-flight mitigation.

## Source 04

Title: Cache stampede
Publisher: Wikipedia (English)
URL: https://en.wikipedia.org/wiki/Cache_stampede
Published: Last edited 2025-12-31 (rev 1330515031 viewed)
Accessed: 2026-09-28
Source Tier: 3 — community encyclopedia (cites Tier-1 sources including Vattani et al. VLDB 2015)
Relevance: Core synthesis of cache stampede definition, typical usage pattern pseudocode, three mitigation categories (locking, external recomputation, probabilistic early expiration / x-fetch algorithm), quantitative example.

## Source 05

Title: Thundering herd problem
Publisher: Wikipedia (English)
URL: https://en.wikipedia.org/wiki/Thundering_herd_problem
Published: Last edited 2026-05-27 (rev 1356321785 viewed)
Accessed: 2026-09-28
Source Tier: 3 — community encyclopedia
Relevance: Disambiguates thundering herd (general wake-up) vs cache stampede (cache-miss dog-piling); documents OS mitigations (epoll EPOLLEXCLUSIVE, IOCP, jitter/backoff) relevant to Tier-2 analogy.

## Source 06

Title: Cache (computing) — Operation / Write policies
Publisher: Wikipedia (English)
URL: https://en.wikipedia.org/wiki/Cache_(computing)
Published: Last edited date not recorded at fetch time (content retrieved 2026-09-28)
Accessed: 2026-09-28
Source Tier: 3 — community encyclopedia (section cites Hennessy & Patterson)
Relevance: Defines write-through vs write-back (write-behind), write-allocate vs no-write-allocate pairing; used to cross-check lab's write-through / write-behind table.

## Source 07

Title: Cache invalidation
Publisher: Wikipedia (English)
URL: https://en.wikipedia.org/wiki/Cache_invalidation
Published: Last edited 2023-12-07 (rev 1188797099 viewed)
Accessed: 2026-09-28
Source Tier: 3 — community encyclopedia
Relevance: Defines explicit invalidation (purge / refresh / ban), alternatives (short TTL, per-request validation), and complexity trade-offs.

## Source 08

Title: Optimal Probabilistic Cache Stampede Prevention
Publisher: Proceedings of the VLDB Endowment — Vattani, A.; Chierichetti, F.; Lowenstein, K. (2015), Vol. 8 No. 8, pp. 886–897; DOI 10.14778/2757807.2757813
URL: https://doi.org/10.14778/2757807.2757813
Published: 2015 (VLDB)
Accessed: 2026-09-28
Source Tier: 1 — peer-reviewed academic paper (proceedings)
Relevance: Primary source for XFetch / probabilistic early expiration optimality claim and exponential-distribution algorithm. Full PDF fetched but binary stream not parseable by webfetch; metadata and abstract verified via DOI landing page and secondary citation in Source 04.

## Source 09

Title: Optimal Probabilistic Cache Stampede Prevention — PDF (primary)
Publisher: UCSD CSE / VLDB — http://cseweb.ucsd.edu/~avattani/papers/cache_stampede.pdf and https://www.vldb.org/pvldb/vol8/p886-vattani.pdf
URL: http://cseweb.ucsd.edu/~avattani/papers/cache_stampede.pdf
Published: 2015
Accessed: 2026-09-28
Source Tier: 1 — primary paper PDF
Relevance: Same as Source 08. Fetch returned PDF binary (truncated/decompression failure) — content NOT extracted. Recorded as attempted primary verification; reliance falls back to DOI metadata + Wikipedia's extracted algorithm.

## Source 10

Title: Redis / AWS ElastiCache / GCP cache pattern docs (attempted)
Publisher: redis.io, docs.aws.amazon.com, cloud.google.com
URL: Multiple attempted (e.g., https://redis.io/docs/latest/develop/use/cache, https://redis.io/docs/latest/develop/clients/patterns/cache-invalidation/, https://docs.aws.amazon.com/AmazonElastiCache/latest/mem-protocol/Protocol.html)
Published: —
Accessed: 2026-09-28
Source Tier: 1 (would-be)
Relevance: Intended to corroborate Redis TTL, SETEX, cache-aside/write-through specifics. All returned 404/403 at current site structure — no evidence extracted. Claims requiring Redis-official wording marked NOT VERIFIED.

## Source 11

Title: Internal codebase reference — labs/04-caching (cache_aside.go, write_through.go, stampede.go)
Publisher: Local repository software-engineering-lab
URL: file://labs/04-caching/cache_aside.go, write_through.go, stampede.go
Published: Repository-local (no publication date)
Accessed: 2026-09-28
Source Tier: 2 (reputable but not external authoritative) — used as contextual example, not as evidence for general claims
Relevance: Concrete Go implementation of TTL (5 min) + jitter (15 s), cache-aside with metric classification, broken vs protected (singleflight.DoChan) stampede service; illustrates lab exercise design.

Notes on freshness: Research date 2026-09-28. Source 01 is current (2025-09). RFC 5861 (2010) remains current; not superseded. VLDB 2015 paper is historical primary evidence. Wikipedia pages are periodically edited; revision IDs pinned above.
