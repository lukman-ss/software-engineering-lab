# 02 — Source Audit

Target Lab: labs/37-cache-invalidation-strategies

## Source 01

Claimed Title: Cache-Aside Pattern
Claimed Publisher: Microsoft Learn — Azure Architecture Center
URL: https://learn.microsoft.com/en-us/azure/architecture/patterns/cache-aside

Reachable: YES
Source Type: PRIMARY (vendor architecture guidance)
Relevant: YES
Supports Claimed Topic: YES
Problems: None. Verified read-on-miss, store-then-delete write ordering, TTL caveats, and read-after-write differences.
Assessment: PASS

---

## Source 02

Claimed Title: RFC 5861 — HTTP Cache-Control Extensions for Stale Content
Claimed Publisher: IETF / RFC Editor (Author: M. Nottingham)
URL: https://datatracker.ietf.org/doc/html/rfc5861

Reachable: YES
Source Type: PRIMARY (RFC specification - Informational)
Relevant: YES
Supports Claimed Topic: YES
Problems: None. Status accurately documented as Informational / Independent Submission stream; §3 stale-while-revalidate and §4 stale-if-error verified.
Assessment: PASS

---

## Source 03

Claimed Title: singleflight package — golang.org/x/sync/singleflight
Claimed Publisher: Go Project (pkg.go.dev)
URL: https://pkg.go.dev/golang.org/x/sync/singleflight

Reachable: YES
Source Type: PRIMARY (official package documentation)
Relevant: YES
Supports Claimed Topic: YES
Problems: None. Verified duplicate call suppression semantics, `Do`/`DoChan`, and in-process scope.
Assessment: PASS

---

## Source 04

Claimed Title: Cache stampede
Claimed Publisher: Wikipedia (English)
URL: https://en.wikipedia.org/wiki/Cache_stampede

Reachable: YES
Source Type: COMMUNITY / TERTIARY
Relevant: YES
Supports Claimed Topic: YES
Problems: None when used for discovery/synthesis. Accurately synthesizes Vattani et al. (2015) XFetch algorithm pseudocode and stampede failure mode.
Assessment: PASS

---

## Source 05

Claimed Title: Thundering herd problem
Claimed Publisher: Wikipedia (English)
URL: https://en.wikipedia.org/wiki/Thundering_herd_problem

Reachable: YES
Source Type: COMMUNITY / TERTIARY
Relevant: YES
Supports Claimed Topic: YES
Problems: None. Correctly distinguishes OS kernel wake-up contention from cache-level miss stampede.
Assessment: PASS

---

## Source 06

Claimed Title: Cache (computing) — Operation / Write policies
Claimed Publisher: Wikipedia (English)
URL: https://en.wikipedia.org/wiki/Cache_(computing)

Reachable: YES
Source Type: COMMUNITY / TERTIARY
Relevant: YES
Supports Claimed Topic: YES
Problems: Defines hardware-level write policies (write-through vs write-back). Research correctly notes this is analogous but distinct from application-tier DB write policies.
Assessment: PASS

---

## Source 07

Claimed Title: Cache invalidation
Claimed Publisher: Wikipedia (English)
URL: https://en.wikipedia.org/wiki/Cache_invalidation

Reachable: YES
Source Type: COMMUNITY / TERTIARY
Relevant: YES
Supports Claimed Topic: YES
Problems: Thin citations on specific proxy invalidations, but adequately used for broad concept coverage.
Assessment: PASS

---

## Source 08

Claimed Title: Optimal Probabilistic Cache Stampede Prevention
Claimed Publisher: Proceedings of the VLDB Endowment (Vattani, Chierichetti, Lowenstein, 2015)
URL: https://doi.org/10.14778/2757807.2757813

Reachable: YES (DOI resolver / landing page)
Source Type: PRIMARY (peer-reviewed academic publication)
Relevant: YES
Supports Claimed Topic: YES
Problems: Landing page and metadata verified. Primary PDF body parsing failed during research; research report properly marks proof/theorems as bibliographic authority and NOT VERIFIED from primary text.
Assessment: PASS

---

## Source 09

Claimed Title: Optimal Probabilistic Cache Stampede Prevention — PDF
Claimed Publisher: UCSD CSE / VLDB
URL: http://cseweb.ucsd.edu/~avattani/papers/cache_stampede.pdf

Reachable: YES (binary stream)
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: Binary decompression unparsed in automated tooling. Handled transparently by research as NOT VERIFIED for internal proof extraction.
Assessment: PASS

---

## Source 10

Claimed Title: Redis / AWS ElastiCache / GCP cache pattern docs
Claimed Publisher: redis.io, docs.aws.amazon.com
URL: Multiple attempted

Reachable: NO (404/403 at attempted legacy paths)
Source Type: PRIMARY (unreachable)
Relevant: YES
Supports Claimed Topic: PARTIAL
Problems: URLs 404'd. Research report and revision explicitly flagged all Redis-specific API claims as NOT VERIFIED and substituted Microsoft Learn as primary vendor guidance.
Assessment: WARNING (Accurately documented and mitigated by research)

---

## Source 11

Claimed Title: Internal codebase reference — labs/04-caching
Claimed Publisher: Local repository
URL: file://labs/04-caching/cache_aside.go

Reachable: YES
Source Type: LOCAL ARTIFACT / SECONDARY
Relevant: YES
Supports Claimed Topic: YES
Problems: Used only as local context for exercise design, not as universal proof.
Assessment: PASS
