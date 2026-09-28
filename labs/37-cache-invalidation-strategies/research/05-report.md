# 05 — Research Report

## Research Question
How to select and implement caching patterns (Cache-Aside / Write-Through / Write-Behind) while preventing cache stampede / thundering herd in a high-read backend using Redis, single-flight / probabilistic expiration, and stale-while-revalidate?

## Executive Summary

**Verified facts:**
- Cache stampede (also dog-piling) is a cascading failure where, under high concurrency, expiry of a popular cached item causes many requests to simultaneously recompute, exhausting shared resources (databases, connection pools, CPU) — potentially driving the hit rate to zero until load drops (Wikipedia, citing Vattani et al. PVLDB 2015) [Evidence 5–6].
- Three established mitigation categories exist: (1) locking / single-flight, (2) external recomputation (dedicated worker), (3) probabilistic early expiration (XFetch) using an exponential distribution (Wikipedia summarizing the 2015 paper; paper metadata via DOI) [Evidence 8, 9, 10, 14].
- The Go standard library offers `golang.org/x/sync/singleflight`, which exactly expresses the single-flight mitigation at the in-process level (only one goroutine per key runs a function; duplicates wait and share the result) [Evidence 9].
- Stale-while-revalidate is defined in RFC 5861 (2010, Informational) as a Cache-Control extension: serve stale up to delta seconds, revalidating asynchronously; validation must be request-triggered to avoid amplification [Evidence 12, 13].
- Cache-Aside (lazy) loads on demand; is simple and memory-efficient; but does not guarantee read-after-write freshness — a write+invalidate gap may be observed [Evidence 1, 4].
- Write-through writes synchronously to both cache and store; read-after-write is immediate within the process; "same write operation" in vendor docs means application-level sequential write, NOT a distributed ACID transaction [Evidence 3, 4, 6 (C6)].
- Write-behind (write-back) writes only to cache and defers the backing-store write; crash-before-flush risks data loss; typically paired with write-allocate [Evidence 3].
- TTL choice is a trade-off: too short → continuous fetches (high store load); too long → staleness; works best for read-heavy, relatively static data [Evidence 15].
- Jitter on expiry timestamps desynchronizes simultaneous expires but does not cap concurrent rebuilds of a single hot key — distinct from probabilistic early expiration (Evidence 2, Evidence 16).

**Interpretations (not asserted as verified):**
- "Thundering herd" in the lab title is used synonymously with cache stampede; the underlying OS phenomenon is broader (wake-one contention) — acceptable in caching context but terminology should be kept clear [C2].
- RFC 5861 stale-while-revalidate's request-triggered design is the safer analogue for application-level SWR; an unconditional background job risks amplification (interpretation based on RFC §5 security guidance) [C5].
- Vattani et al.'s XFetch optimality claim exists in title and Wikipedia summary; the paper's proofs and benchmarks were not independently parsed (PDF fetch failed) — the claim is accepted on bibliographic authority alone [Evidence 11].
- Lab's pedagogical scenario (10,000 RPS, 500 concurrent goroutines, P99 comparison) is not cited from production data; it is a synthetic lab design [Evidence 18].

**Confidence profile:** HIGH for Cache-Aside mechanics, singleflight semantics, stampede definition, and RFC 5861 SWR; MEDIUM for XFetch optimality claim, write-behind data-loss risk implication, and jitter's sufficiency; LOW for any Redis-specific behavioral claims due to fetch failures.

## Findings

### Finding 1 — Cache-Aside (lazy loading) is the default pattern for read-heavy workloads

Claim: Cache-Aside loads data into the cache on demand on miss; stores first, then deletes/invalidates the cache on write. It is efficient in memory (only read data is cached) but does not guarantee immediate read-after-write consistency.

Evidence: Microsoft Architecture Center: application checks cache; on miss fetches from store and adds to cache; on update writes to store then invalidates cache — explicitly states the order matters (store before delete) to avoid a window where a reader fetches stale data. Also warns "an external process can change an item in the data store at any time" [Evidence 1, 2, 4]. Local `cache_aside.go` implements the same sequence with a 5-minute TTL and 15-second jitter.

Sources: Microsoft Learn — Cache-Aside Pattern; Wikipedia Cache stampede §Typical; local labs/04-caching/cache_aside.go.
Confidence: HIGH.

### Finding 2 — Write-through provides read-after-write freshness at the cost of synchronous writes to both storage and cache

Claim: Write-through updates both cache and the backing store as part of the same application-level write sequence, so subsequent reads see the new value. It trades higher write latency and more cache usage (even for rarely-read keys) for immediate consistency.

Evidence: Microsoft distinguishes Cache-Aside from write-through: "write-through caching... updates the data store and the cache in the same write operation so that readers see the new value immediately after a successful write" [Evidence 4]. Hardware-cache definitions in Wikipedia specify "writes performed synchronously to both the cache and the backing store" [Evidence 3]. Local `write_through.go` explicitly warns the DB and Redis are separate systems (no atomic cross-system commit) and implements a best-effort cache update after the authoritative DB RETURNING — TTL remains the safety net for crashes between the two writes.

Sources: Microsoft Learn; Wikipedia Cache (computing) §Write policies; local labs/04-caching/write_through.go.
Confidence: HIGH.

### Finding 3 — Write-behind (write-back) decouples writes from the store for throughput, at the risk of data loss on cache failure

Claim: Write-behind writes only to cache; a background process or eviction-driven flush later propagates to the backing store. This yields high write throughput but loses durability guarantees if the cache host crashes before a dirty entry is flushed.

Evidence: Wikipedia: "Initially, writing is done only to the cache. The write to the backing store is postponed until the modified content is about to be replaced by another cache block... a read miss in a write-back cache may require two memory accesses" [Evidence 3]. The lab spec states the same risk ("Risiko data loss jika cache node crash sebelum sync ke DB") for the application-level analogue of write-behind to a relational DB. This is a logical consequence of the postponement semantics, but no opened source quotes the exact crash-scenario phrase for application-tier write-behind.

Sources: Wikipedia Cache (computing); lab spec text.
Confidence: MEDIUM (mechanism HIGH, crash risk is an inference from postponement).

### Finding 4 — Cache stampede is a predictable, quantifiable failure mode

Claim: When a highly-referenced cache key expires, all in-flight requests observe a miss and may recompute in parallel. If each recomputation takes `T` seconds and the request arrival rate during the window is `λ`, roughly `λ·T` workers enter the critical section simultaneously. A concrete illustrative instance in the literature: `T = 3 s`, `λ = 10 req/s` → 30 concurrent recomputes; at higher λ the system can suffer congestion collapse, driving the hit rate to zero as each attempt times out [Evidence 5]. The lab's 10,000 RPS flash-sale scenario and 500-goroutine exercise are synthetic parameters, not documented production measurements [Evidence 18].

Evidence: Wikipedia's explicit example and collapse description [Evidence 5]; local `BrokenStampedeService` mirrors the naive path: multiple goroutines bypass the cache check simultaneously and each hits the repository in `stampede.go` (queryDelay = 5 ms per call) [Evidence 6, 10].

Sources: Wikipedia — Cache stampede; local labs/04-caching/stampede.go.
Confidence: HIGH for the mechanism; MEDIUM for the quantitative example (arithmetic illustration, not empirical).

### Finding 5 — Single-flight / distributed-lock coalesces concurrent misses so only one thread queries the store

Claim: On cache miss, acquire a coordination primitive scoped to the key; the first holder recompute-and-write; all other waiters receive the same result (either by blocking or by reading a stale value). At the in-process level Go's `singleflight.Group.Do` provides this exact behavior [Evidence 9]. Distributed lock variants (Redis SET NX + PX) impose an extra write per lock [Evidence 8, C4].

Evidence: Wikipedia outlines the locking approach and its trade-offs [Evidence 8]. Go singleflight docs specify `Do(key, fn)` runs fn once per key, returns the same result to all waiters, and sets `shared = true` [Evidence 9]. Local `ProtectedStampedeService` calls `flight.DoChan(key, ...)` with a 30-second build timeout and re-checks the cache inside the function (double-checked lock pattern) [Evidence 9 + local code].

Sources: Wikipedia — Cache stampede; Go pkg.go.dev singleflight; local labs/04-caching/stampede.go.
Confidence: HIGH.

### Finding 6 — Probabilistic early expiration (XFetch) spreads rebuilds over time via an exponential draw

Claim: Each request computes a random offset `offset = -Δ·β·ln(U)` (U ~ uniform(0,1)) and re-evaluates the key when `now + offset ≥ expiry`. Longer DB round-trips (`Δ`) and traffic bursts naturally increase the probability of an early rebuild before the official TTL; `beta = 1` is the practical default [Evidence 10]. The approach trades slightly earlier-than-TTL evictions (more frequent, smaller waves of rebuild) for eliminating the synchronized spike at expiry [Evidence 10].

Evidence: Wikipedia algorithm and rationale, attributed to Vattani et al. PVLDB 2015 (DOI confirms paper existence) [Evidence 10, 11]. The authors' optimality claim is accepted on bibliographic authority alone; the paper PDF (both URLs) returned raw compressed streams that could not be parsed, so theorem statements and experimental graphs were not verified [Evidence 11, 09]. The lab's pedagogical formula `Δ·β·ln(rand()) > TTL_remaining` appears to omit the minus sign; as written it is always false for `rand() ∈ (0,1)` [Contradiction C1].

Sources: Wikipedia — Cache stampede (algorithm); DOI metadata for VLDB 2015; labs/04-caching/stampede.go (for context).
Confidence: MEDIUM (algorithm description HIGH; optimality proof and exact formula notation NOT VERIFIED).

### Finding 7 — Stale-while-revalidate hides revalidation latency without reducing cache consistency guarantees at the HTTP level

Claim: A cache MAY continue serving a stale response for up to `stale-while-revalidate` seconds while it revalidates the origin in the background. Revalidation is triggered by a subsequent request (not an independent background job), and if traffic is too sparse some requests will still block [Evidence 12]. The RFC additionally recommends that revalidation be request-triggered to avoid amplification attacks [Evidence 12, C5].

Evidence: RFC 5861 §3 and §3.1: example `max-age=600, stale-while-revalidate=30` — fresh 600 s, may serve stale up to 30 s more; "if validation is inconclusive, or if there is not traffic that triggers it, after 30 seconds the stale-while-revalidate function will cease to operate." §5: "such validation be predicated upon an incoming request, to avoid the possibility of an amplification attack" [Evidence 12]. Complement: `stale-if-error` serves stale on origin 5xx/network errors, extending availability under failure rather than under normal expiry [Evidence 13].

Sources: IETF RFC 5861.
Confidence: HIGH for RFC semantics; MEDIUM for applicability at the application-cache layer (not specified by RFC).

### Finding 8 — Jitter is a complementary anti-synchronization knob, not a standalone stampede solution

Claim: Adding random jitter to TTL values (or to retry backoff intervals) desynchronizes expiry or retry clocks across keys and clients, reducing the likelihood that many entries expire simultaneously [Evidence 16]. However, jitter does not prevent the stampede of a single hot key whose expiry still falls within one jitter interval of a burst of miss requests.

Evidence: Wikipedia cites jitter as a general anti-synchronization technique for retry backoffs [Evidence 7]. The repo's `TTLWithJitter(base, maxJitter)` deterministically returns `base + rand[0, maxJitter)` [Evidence 2, 16].

Sources: Wikipedia — Thundering herd; local stampede.go.
Confidence: HIGH for jitter's purpose; MEDIUM for sufficiency claim.

## Areas of Agreement

- Cache-Aside = lazy load on miss, invalidate on write; order (update store, then delete cache) matters to avoid stale re-population [Sources 01, 04, 06, 10].
- Write-through = synchronous dual write; immediate read-after-write visibility [Sources 01, 03, 04, 06].
- Write-behind = cache-first write, store deferred; durability trade-off [Sources 03, 06].
- Cache stampede occurs at expiry of a popular key under heavy concurrency; can degrade to zero-hit-rate congestion [Sources 04, 05, 06].
- Single-flight / in-process lock is the simplest stampede defense; distributed lock adds complexity and extra writes [Sources 04, 09, C4].
- Probabilistic early expiration shifts rebuild timing forward of official TTL using an exponential draw proportional to measured recompute time [Sources 04, 10, 11].
- Stale-while-revalidate serves stale while asynchronously revalidating; must be request-triggered to avoid amplification [Sources 02, 12, C5].
- Jitter desynchronizes but does not eliminate stampede on a single hot key [Sources 05, 07, 16].

## Areas of Disagreement

See Contradictions file (04-contradictions.md). Key tensions:
- C1 — Lab formula sign vs. Wikipedia / paper formulation for XFetch.
- C2 — Terminology conflation (thundering herd vs. cache stampede).
- C3 — RFC 5861 "standard" status vs. informational classification.
- C4 — Distributed lock extra-write cost vs. in-process singleflight (scope, not contradiction).
- C5 — Request-triggered vs. background-job trigger for SWR revalidation.
- C6 — "Same write operation" in vendor docs vs. separate DB/Redis commit reality.

## Limitations

- Research date: 2026-09-28. Sources span 2010 (RFC 5861), 2015 (VLDB paper), 2025 (Microsoft doc), 2026 (Go singleflight v0.23.0 release). No post-2025 empirical cache-stampede benchmark study was opened.
- The primary XFetch paper (Vattani et al., 2015) was fetched as PDF from two URLs (`http://cseweb.ucsd.edu/~avattani/papers/cache_stampede.pdf` and `https://www.vldb.org/pvldb/vol8/p886-vattani.pdf`) but the binary stream could not be parsed by the fetch tool; the optimality claim and exact formula rely on Wikipedia's summary and DOI metadata (MEDIUM confidence for proof details).
- Redis-specific docs (cache-aside, stampede, SET NX/PX) returned 404/403 at attempted paths; no Redis-official claims were extracted.
- No production incident reports, A/B benchmarks, or P99 numbers were found. The lab's 10,000 RPS / 500 goroutine / P99-before-after numbers are not sourced (marked NOT VERIFIED).
- Wikipedia is a tertiary source; citations were used to identify primary references, not as stand-alone evidence. All important claims were cross-checked against at least one Tier-1 or original-repository source where available.
- The lab's cache-implementation details (TTL=5 min, jitter=15 s, singleflight DoChan with 30 s timeout) are local conventions, not general recommendations.

## Conclusion

Cache invalidation choices trade latency, consistency, and operational complexity. For the read-heavy backend described in the lab, the principal risk is a synchronized expiration spike (cache stampede) that multiplicatively hits the database. The evidence supports three principal mitigations:

1. **Single-flight (in-process)** — simplest and sufficient when all readers share the same process (Go `singleflight.Group`). Verified by vendor docs, language library docs, and existing local code. Does not protect a multi-process deployment without a distributed analogue.
2. **Probabilistic early expiration (XFetch)** — removes synchronization from expiry time by pulling rebuilds forward according to an exponential distribution scaled by recompute duration and a tunable parameter beta. Algorithm details come from a peer-reviewed paper; the paper's proofs are not independently verified in this pass. Care must be taken with formula notation (minus sign on ln(rand)).
3. **Stale-while-revalidate** — serves fresh-to-stale data without blocking while an async revalidation is request-triggered. Defined by RFC 5861; applicable at the HTTP cache layer and conceptually transferable to app caches, provided the request-correlation constraint is preserved.

Jitter and TTL tuning are complementary hygiene (desynchronize keys, avoid overly aggressive or passive TTLs) but are not substitutes for the above three. Write-through is preferable where read-after-write consistency is required and write traffic is modest; write-behind gains throughput but introduces data-loss risk on cache-node failure. Cache-aside remains the pragmatic default for most web-layer reads, combined with single-flight or XFetch at cache-miss time.

No new factual claims are introduced in this conclusion beyond those already documented in the findings.
