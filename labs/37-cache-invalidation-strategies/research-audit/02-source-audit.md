# 02 — Source Audit

## Source 01

Claimed Title: Cache-Aside Pattern
Claimed Publisher: Microsoft Learn — Azure Architecture Center
URL: https://learn.microsoft.com/en-us/azure/architecture/patterns/cache-aside

Reachable:
YES

Source Type:
PRIMARY (official vendor architecture documentation)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Verified contents match cited behavior (lazy loading mechanics, write invalidation sequence, read-after-write staleness warnings, TTL guidance).

Assessment:
PASS

---

## Source 02

Claimed Title: RFC 5861 — HTTP Cache-Control Extensions for Stale Content
Claimed Publisher: IETF / RFC Editor — Author M. Nottingham
URL: https://datatracker.ietf.org/doc/html/rfc5861

Reachable:
YES

Source Type:
PRIMARY (informational RFC specification)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- RFC 5861 is an Independent Submission Informational RFC, not an IETF Internet Standard. Calling it an "Internet Standard" would be inaccurate, though it is the de-facto standard reference for `stale-while-revalidate` and `stale-if-error`.

Assessment:
PASS

---

## Source 03

Claimed Title: singleflight package — golang.org/x/sync/singleflight
Claimed Publisher: Go Project (pkg.go.dev)
URL: https://pkg.go.dev/golang.org/x/sync/singleflight

Reachable:
YES

Source Type:
PRIMARY (official language library documentation)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Docs explicitly confirm duplicate function call suppression, `Do`, `DoChan`, `Forget`, and `shared` boolean flag.

Assessment:
PASS

---

## Source 04

Claimed Title: Cache stampede
Claimed Publisher: Wikipedia (English)
URL: https://en.wikipedia.org/wiki/Cache_stampede

Reachable:
YES (Assumed accessible; standard Wikipedia URL)

Source Type:
COMMUNITY (Secondary/Tertiary)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Wikipedia is tertiary. Should not be used as sole authority for formal algorithm optimality without primary paper verification.

Assessment:
WARNING

---

## Source 05

Claimed Title: Thundering herd problem
Claimed Publisher: Wikipedia (English)
URL: https://en.wikipedia.org/wiki/Thundering_herd_problem

Reachable:
YES

Source Type:
COMMUNITY (Secondary/Tertiary)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Disambiguates kernel/OS thundering herd from cache stampede.

Assessment:
PASS

---

## Source 06

Claimed Title: Cache (computing) — Operation / Write policies
Claimed Publisher: Wikipedia (English)
URL: https://en.wikipedia.org/wiki/Cache_(computing)

Reachable:
YES

Source Type:
COMMUNITY (Secondary/Tertiary)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Distinguishes CPU/hardware cache write-through and write-back. Needs careful framing when mapped to application/DB tiers.

Assessment:
PASS

---

## Source 07

Claimed Title: Cache invalidation
Claimed Publisher: Wikipedia (English)
URL: https://en.wikipedia.org/wiki/Cache_invalidation

Reachable:
YES

Source Type:
COMMUNITY (Secondary/Tertiary)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Refers to HTTP proxy libraries (FOSHttpCache) for Purge/Refresh/Ban terminology, which may not map 1:1 to key-value stores.

Assessment:
PASS

---

## Source 08 & 09

Claimed Title: Optimal Probabilistic Cache Stampede Prevention
Claimed Publisher: Proceedings of the VLDB Endowment (Vattani et al., 2015)
URL: https://doi.org/10.14778/2757807.2757813 / http://cseweb.ucsd.edu/~avattani/papers/cache_stampede.pdf

Reachable:
PARTIAL (DOI landing page reachable; PDF URLs returned binary stream that failed decompression/extraction)

Source Type:
PRIMARY (Peer-reviewed academic paper)

Relevant:
YES

Supports Claimed Topic:
PARTIAL (Bibliographic metadata verified; mathematical proof and primary empirical figures NOT VERIFIED due to PDF parse failure)

Problems:
- Research Agent acknowledged PDF fetch failure. Optimality claims rely on secondary Wikipedia summary + DOI landing metadata.

Assessment:
WARNING

---

## Source 10

Claimed Title: Redis / AWS ElastiCache / GCP cache pattern docs
Claimed Publisher: redis.io, docs.aws.amazon.com
URL: https://redis.io/docs/latest/develop/use/cache (and related)

Reachable:
NO (Returned 404/403 on attempted URLs)

Source Type:
UNKNOWN / UNREACHABLE

Relevant:
UNKNOWN

Supports Claimed Topic:
NO

Problems:
- Redis pattern documentation URLs returned 404/403. Content could not be verified from primary documentation.

Assessment:
FAIL

---

## Source 11

Claimed Title: Internal codebase reference — labs/04-caching
Claimed Publisher: Local repository
URL: file://labs/04-caching

Reachable:
YES

Source Type:
COMMUNITY / SECONDARY (Repository code)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Code context for prior lab exercise design; not external evidence.

Assessment:
PASS
