# 03 — Evidence

Accessed: 2026-09-28.

## Evidence 1

Claim: Cache-Aside (lazy loading) loads data into cache on demand: check cache, on miss fetch from store, then populate cache; writes typically update the store then invalidate the cache entry.

Evidence: Microsoft Azure Architecture Center: "This strategy loads data into the cache on demand." Steps: (1) application attempts to read from cache; (2) on cache miss, retrieve from data store; (3) add item to cache and return. "When an application updates information, it writes the change to the data store and then invalidates the corresponding item in the cache." Order is specified: update store *before* removing cache item, otherwise a concurrent reader can re-cache stale data.

Source: Microsoft Learn — Cache-Aside Pattern
URL: https://learn.microsoft.com/en-us/azure/architecture/patterns/cache-aside
Confidence: HIGH
Corroborated By: Wikipedia Cache stampede typical usage pattern (read-then-recompute-then-write); local labs/04-caching/cache_aside.go implements the same sequence.
Notes: Cache-Aside does not guarantee store/cache consistency; "an external process can change an item in the data store at any time."

## Evidence 2

Claim: Cache-Aside is appropriate when demand is unpredictable and the cache lacks native read-through/write-through; it is a poor fit when hit rate is low, data is sensitive/shared, or the set is small and static.

Evidence: Microsoft: Use when "A cache doesn't provide native read-through and write-through operations" and "Resource demand is unpredictable." Unsuitable when "Most requests don't experience a cache hit" (overhead outweighs benefit), data is security-sensitive, or the set is static and can be primed.

Source: Microsoft Learn — Cache-Aside Pattern
URL: https://learn.microsoft.com/en-us/azure/architecture/patterns/cache-aside
Confidence: HIGH
Corroborated By: None independent at equal authority. Supporting interpretation only.
Notes: Vendor guidance, not a theorem. Transferable beyond Azure Redis.

## Evidence 3

Claim: Write-through writes synchronously to cache and backing store; write-back (write-behind) writes to cache first and postpones store write until eviction (or explicit flush).

Evidence: Wikipedia Cache (computing): "Write-through: Writes are performed synchronously to both the cache and the backing store." "Write-back: Initially, writing is done only to the cache. The write to the backing store is postponed until the modified content is about to be replaced by another cache block." Write-back must track dirty locations; a read miss may require two store accesses (write-back dirty + fetch). Typical pairing: write-back with write-allocate; write-through with no-write-allocate.

Source: Wikipedia — Cache (computing), Write policies
URL: https://en.wikipedia.org/wiki/Cache_(computing)
Confidence: MEDIUM
Corroborated By: Microsoft Cache-Aside page distinguishes write-through as updating store and cache in the same write so readers see the new value immediately; cites write-through caching doc. Hardware textbook pairing cited as Hennessy & Patterson (not independently opened).
Notes: Hardware-cache write-back ≠ application write-behind to a DB (crash-before-sync risk is analogous but not identical). Lab table's "background worker sync" is an application-level write-behind, not CPU write-back.

## Evidence 4

Claim: Cache-Aside does not provide read-after-write freshness the way write-through does; between write+invalidate and next read, a reader may miss or briefly see stale data.

Evidence: Microsoft: "The Cache-Aside pattern invalidates the cached entry on write and repopulates it on the next read. Between the write and the next read, a reader can experience a cache miss or briefly see stale data. This behavior distinguishes the Cache-Aside pattern from write-through caching, which updates the data store and the cache in the same write operation so that readers see the new value immediately after a successful write."

Source: Microsoft Learn — Cache-Aside Pattern
URL: https://learn.microsoft.com/en-us/azure/architecture/patterns/cache-aside
Confidence: HIGH
Corroborated By: Local write_through.go comments: DB and Redis are separate systems; no atomic commit; TTL as safety net if crash between DB commit and cache update.
Notes: "Same write operation" in vendor docs is application-level sequential write, not a distributed transaction.

## Evidence 5

Claim: A cache stampede (also dog-piling) is a cascading failure: under high concurrency, when a popular cached item expires, many threads recompute simultaneously, exhausting shared resources and potentially preventing the item from ever being re-cached.

Evidence: Wikipedia Cache stampede: "A cache stampede is a type of cascading failure that can occur when massively parallel computing systems with caching mechanisms come under a very high load." Concrete example: page takes 3 seconds to render, 10 req/s → 30 processes recompute simultaneously at expiry. Under very heavy load this can cause congestion collapse: "preventing the page from ever being completely re-rendered and re-cached, as every attempt to do so times out. Thus, cache stampede reduces the cache hit rate to zero."

Source: Wikipedia — Cache stampede
URL: https://en.wikipedia.org/wiki/Cache_stampede
Confidence: HIGH (definition widely used; example is illustrative, not a measured production incident)
Corroborated By: Wikipedia Thundering herd: "When the thundering herd problem occurs while attempting to access a cache, it is often referred to as a cache stampede."
Notes: 10 req/s × 3 s = 30 concurrent recomputes is arithmetic, not empirical. Lab's 10,000 RPS flash-sale scenario is hypothetical, not sourced.

## Evidence 6

Claim: Naive cache-aside fetch on miss has no coordination; many processes call recompute_value() at the same expiry.

Evidence: Wikipedia Cache stampede typical pattern:

```
function fetch(key, ttl):
    value ← cache_read(key)
    if (!value):
        value ← recompute_value()
        cache_write(key, value, ttl)
    return value
```

"If the function recompute_value() takes a long time and the key is accessed frequently, many processes will simultaneously call recompute_value() upon expiration of the cache value."

Source: Wikipedia — Cache stampede
URL: https://en.wikipedia.org/wiki/Cache_stampede
Confidence: HIGH
Corroborated By: labs/04-caching/stampede.go BrokenStampedeService.GetData — identical sequence without coalescing.
Notes: This is the lab's "Cache-Aside biasa" baseline.

## Evidence 7

Claim: Thundering herd (OS/network sense) is distinct from but related to cache stampede: many waiters wake for one event; only one proceeds; the rest sleep again, wasting CPU.

Evidence: Wikipedia Thundering herd: "a large number of processes or threads are simultaneously awakened... However, only one process is able to respond to the event or access the new resource, causing most other processes to fail and go back to sleep." Linux serializes waiters on a single fd; EPOLLEXCLUSIVE (kernel 4.5) wakes one epoll set; Windows IOCP can wake one waiter. Jitter on retry intervals desynchronizes clients.

Source: Wikipedia — Thundering herd problem
URL: https://en.wikipedia.org/wiki/Thundering_herd_problem
Confidence: MEDIUM (relies on Jargon File, LKML, blog posts as cited refs)
Corroborated By: Cache stampede page treats stampede as the cache-access instance of the broader phenomenon.
Notes: Lab "Thundering Herd" usage matches cache stampede, not kernel accept-queue thundering herd. Keep terms disambiguated in teaching.

## Evidence 8

Claim: Locking / single-flight: on miss, one process acquires a lock (or is elected) to recompute; others wait, return not-found, or serve stale. Correct locking must handle lock-holder crash, lock TTL, and races.

Evidence: Wikipedia Cache stampede, Locking: "upon a cache miss a process will attempt to acquire the lock for that cache key and recompute it only if it acquires it." Options if lock not acquired: wait; return not-found; keep stale. "If implemented properly, locking can prevent stampedes altogether, but requires an extra write for the locking mechanism... the main drawback is a correct implementation of the locking mechanism which also takes care of edge cases including failure of the process acquiring the lock, tuning of a time-to-live for the lock, race-conditions, and so on."

Source: Wikipedia — Cache stampede
URL: https://en.wikipedia.org/wiki/Cache_stampede
Confidence: HIGH for the mechanism; MEDIUM for "prevent stampedes altogether" (conditional on correct implementation)
Corroborated By: Go singleflight (Evidence 9).
Notes: Distributed Redis SET NX + TTL lock is a multi-process analogue; in-process singleflight is not a distributed lock.

## Evidence 9

Claim: Go `singleflight.Group.Do` executes fn once per key while duplicates wait and share the result (`shared` flag).

Evidence: pkg.go.dev: "Package singleflight provides a duplicate function call suppression mechanism." `Do`: "making sure that only one execution is in-flight for a given key at a time. If a duplicate comes in, the duplicate caller waits for the original to complete and receives the same results. The return value shared indicates whether v was given to multiple callers." Example: two DoChan on same key; only first fn runs; both get "func 1"; Shared: true.

Source: Go — golang.org/x/sync/singleflight
URL: https://pkg.go.dev/golang.org/x/sync/singleflight
Confidence: HIGH
Corroborated By: labs/04-caching/stampede.go ProtectedStampedeService uses `flight.DoChan` around DB+Set; second cache Get inside fn (double-check).
Notes: Process-local. N replicas ⇒ up to N DB queries unless a distributed lock or shared cache lock is added. `Forget(key)` lets a later Do run fn instead of waiting.

## Evidence 10

Claim: Probabilistic early expiration (XFetch) lets each process independently decide to recompute before expiry using an exponential distribution scaled by recompute duration delta and parameter beta.

Evidence: Wikipedia Cache stampede, citing Vattani, Chierichetti, Lowenstein, PVLDB 8(8):886–897, 2015, doi:10.14778/2757807.2757813:

```
function x-fetch(key, ttl, beta=1):
    value, delta, expiry ← cache_read(key)
    if (!value || (time() - delta * beta * log(rand(0,1))) ≥ expiry):
        start ← time()
        value ← recompute_value()
        delta ← time() – start
        cache_write(key, (value, delta), ttl)
    return value
```

"The parameter beta can be set to a value greater than 1 to favor earlier recomputations... the authors show that setting beta=1 works well in practice. The variable delta represents the time to recompute the value." "This approach is simple to implement and effectively reduces cache stampedes by automatically favoring early recomputations when the traffic rate increases." Drawback: extra cache memory for delta (and expiry if TTL cannot be read).

Source: Wikipedia — Cache stampede (algorithm quoted from VLDB 2015 paper)
URL: https://en.wikipedia.org/wiki/Cache_stampede
Primary paper: https://doi.org/10.14778/2757807.2757813
Confidence: MEDIUM — algorithm and beta=1 practice claim are widely cited; optimality proof and experiments not independently extracted from the PDF (parse failure).
Corroborated By: DOI landing confirms title, authors, venue, pages. Wikipedia states the exponential implementation "has been shown to be optimal in terms of its effectiveness in preventing stampedes and how early recomputations can happen" citing the paper.
Notes: Lab formula `Δ · β · ln(rand()) > TTL_remaining` is a rearrangement of the Wikipedia condition, not a verbatim quote. Equivalence: `log(rand(0,1))` is ≤ 0, so `-delta*beta*log(rand)` is a positive offset; early refresh when `now + offset ≥ expiry`, i.e. remaining TTL ≤ offset. Lab uses `ln(rand())` without specifying rand∈(0,1); if rand∈(0,1), ln is negative and the inequality direction matters. Treat lab formula as pedagogical shorthand, not the paper's canonical form.

## Evidence 11

Claim: Vattani, Chierichetti & Lowenstein (2015) published "Optimal Probabilistic Cache Stampede Prevention" in PVLDB 8(8):886–897.

Evidence: DOI.org landing: "Vattani, A., Chierichetti, F., & Lowenstein, K. (2015). Optimal probabilistic cache stampede prevention. Proceedings of the VLDB Endowment, 8(8), 886–897. https://doi.org/10.14778/2757807.2757813"

Source: DOI / Crossref-style landing
URL: https://doi.org/10.14778/2757807.2757813
Confidence: HIGH for bibliographic facts
Corroborated By: Wikipedia Cache stampede reference [3]; PDF URLs http://cseweb.ucsd.edu/~avattani/papers/cache_stampede.pdf and https://www.vldb.org/pvldb/vol8/p886-vattani.pdf both returned PDF bytes (unparsed).
Notes: Paper existence is verified. Theorems, experiment tables, and exact optimality statement: NOT VERIFIED from primary text.

## Evidence 12

Claim: `stale-while-revalidate` allows a cache to serve a response after it becomes stale for up to delta-seconds while revalidating in the background (non-blocking).

Evidence: RFC 5861 §3: "caches MAY serve the response in which it appears after it becomes stale, up to the indicated number of seconds." "If a cached response is served stale due to the presence of this extension, the cache SHOULD attempt to revalidate it while still serving stale responses (i.e., without blocking)." After the window without revalidation, "it SHOULD NOT continue to be served stale." Example: `Cache-Control: max-age=600, stale-while-revalidate=30` — fresh 600s, may serve stale 30s more during async validation. "If the window is too small, or traffic is too sparse, some requests will fall outside of it, and block until the server can validate the cached response."

Source: IETF RFC 5861
URL: https://datatracker.ietf.org/doc/html/rfc5861
Confidence: HIGH
Corroborated By: None needed (standard text). Application-level SWR is an analogue, not specified by this RFC.
Notes: Informational RFC, Independent Submission, "not endorsed by the IETF" / no formal standards standing (status block). Still the de-facto definition of the directive; later HTTP caching practice adopted it (implementation status NOT VERIFIED here). Security: validation SHOULD be request-triggered to avoid amplification/prefetch attacks (§5).

## Evidence 13

Claim: `stale-if-error` allows serving stale content when the origin would return 500/502/503/504, up to a staleness bound — availability under origin failure, not stampede prevention per se.

Evidence: RFC 5861 §4: "when an error is encountered, a cached stale response MAY be used to satisfy the request." Errors: 500, 502, 503, 504. Example: max-age=600, stale-if-error=1200 — after 900s a 500 can be replaced by the stale 200; after age>1800 the error is passed through.

Source: IETF RFC 5861
URL: https://datatracker.ietf.org/doc/html/rfc5861
Confidence: HIGH
Corroborated By: —
Notes: Complementary to SWR. Lab topic lists SWR; stale-if-error is related resilience, not stampede coalescing.

## Evidence 14

Claim: External recomputation (dedicated refresher process) can prevent stampede by moving rebuild off the request path; it fits static keys better than dynamic id-indexed keys.

Evidence: Wikipedia Cache stampede: triggers — approaching expiration; periodically; or on miss. "This approach requires one more moving part... In addition, this solution requires unnatural code separation/duplication and is mostly suited for static cache keys (i.e., not dynamically generated, as in the case of keys indexed by an id)."

Source: Wikipedia — Cache stampede
URL: https://en.wikipedia.org/wiki/Cache_stampede
Confidence: MEDIUM
Corroborated By: RFC 5861 background revalidation is a related (HTTP-cache) form of external/async recomputation.
Notes: Lab does not require this pattern; useful as a third category.

## Evidence 15

Claim: TTL that is too short causes continuous store fetch; too long causes staleness. Caching works best for relatively static or frequently read data.

Evidence: Microsoft Cache-Aside: "Don't make the expiration period too short because premature expiration can cause applications to continually retrieve data from the data store and add it to the cache. Similarly, don't make the expiration period so long that the cached data becomes stale. Caching works best for relatively static data or data that applications read frequently."

Source: Microsoft Learn — Cache-Aside Pattern
URL: https://learn.microsoft.com/en-us/azure/architecture/patterns/cache-aside
Confidence: HIGH
Corroborated By: Wikipedia Cache invalidation alternatives: reducing TTL "can cause issues, as they create high load on the application due to more frequent requests."
Notes: No universal TTL number. Lab 5-minute TTL + 15s jitter is a local convention, not a standard.

## Evidence 16

Claim: Jitter on TTL (or on retry backoff) desynchronizes expiry/retry clocks so many clients do not expire or retry at the same instant.

**NOTE:** The claim that jitter alone fails to cap concurrent rebuilds of a single hot key — i.e., that jitter does not prevent stampede on one key whose expiry still falls within one jitter interval of a burst of miss requests — is an inferential deduction based on single-key contention mechanics, not a direct quote from a primary caching source. No opened source states this exact limitation verbatim. The deduction is logically sound (a single key's expiry window can still receive multiple concurrent requests even with jitter applied across the key's TTL value) but should be labeled as interpretation, not sourced fact.

Evidence: Wikipedia Thundering herd Mitigation: "randomness is added to the wait intervals between retries, so that clients are no longer synchronized." labs/04-caching/stampede.go: `TTLWithJitter` returns `base + random[0, maxJitter)` — "Never reduces TTL below base."

Source: Wikipedia — Thundering herd problem; local stampede.go
URL: https://en.wikipedia.org/wiki/Thundering_herd_problem
Confidence: HIGH for jitter desynchronizes clocks; MEDIUM (inferential) for single-key insufficiency claim.
Corroborated By: Conceptual overlap with XFetch (randomize *when* rebuild happens) but XFetch randomizes per-request near expiry; jitter randomizes the expiry timestamp itself.
Notes: Jitter helps when many keys share the same TTL origin (deploy, cron, cache warm). For a single hot key, jitter alone does not prevent stampede; the key remains vulnerable to concurrent miss requests during its TTL window.

## Evidence 17

Claim: Explicit invalidation methods include purge (immediate delete of all variants), refresh (replace one variant from origin), and ban (blacklist; replace on next request).

Evidence: Wikipedia Cache invalidation (citing FOSHttpCache docs): Purge removes immediately; Refresh fetches even if cached and replaces one variant; Ban adds to blacklist, origin fetch on matching request. Alternatives: short TTL, validate on each request, do not cache volatile content — all increase origin load. Disadvantage: invalidating many objects is complex; proxy invalidation traffic can slow the proxy.

Source: Wikipedia — Cache invalidation
URL: https://en.wikipedia.org/wiki/Cache_invalidation
Confidence: MEDIUM (thin citations; FOSHttpCache is a PHP HTTP cache library, not a general standard)
Corroborated By: Microsoft write-invalidation (delete after store update) is the application analogue of purge.
Notes: HTTP-proxy vocabulary; application Redis DEL is closer to purge.

## Evidence 18

Claim: Quantitative production figures in the lab spec (10,000 RPS flash-sale, 500 goroutines, P99 before/after, connection-pool exhaustion) are not backed by the sources collected here.

Evidence: No opened source reports 10,000 concurrent misses on one key, 500-goroutine lab numbers, or P99 deltas. Wikipedia's only numeric example is 10 req/s × 3s render = 30 concurrent rebuilds. Microsoft and RFC 5861 give no stampede QPS numbers.

Source: Exhaustive check of Sources 01–09, 12–13
URL: n/a
Confidence: HIGH that these numbers are unspecified by the cited literature
Corroborated By: —
Notes: Treat lab numbers as exercise parameters, not empirical claims. Workload-dependent; do not present as industry benchmarks.

## Evidence 19

Claim: Redis official cache-aside / stampede documentation at currently attempted redis.io paths could not be retrieved.

Evidence: HTTP 404/403 on:
- https://redis.io/learn/howtos/solutions/caching-architecture/cache-aside
- https://redis.io/docs/latest/develop/use/cache
- https://redis.io/docs/latest/develop/clients/patterns/cache-invalidation/
- https://redis.io/docs/latest/develop/use/caching/
- several AWS ElastiCache URLs

Source: Direct fetch 2026-09-28
URL: listed above
Confidence: HIGH (fetch failed)
Corroborated By: —
Notes: Redis SET, GET, EXPIRE, SET NX are industry-standard commands; command semantics NOT independently opened in this pass. Do not invent Redis-doc quotations.

## Evidence 20

Claim: Prior lab 04 already implements naive stampede vs singleflight protection and TTL jitter — a local existence proof of the exercise design, not general evidence.

Evidence: `BrokenStampedeService` does Get → miss → repo.GetDashboard → Set. `ProtectedStampedeService` wraps rebuild in `singleflight.DoChan`, re-checks cache inside fn, 30s rebuild timeout via `context.WithoutCancel`. `queryDelay = 5ms`. `TTLWithJitter(base, maxJitter)`.

Source: labs/04-caching/stampede.go
URL: file://labs/04-caching/stampede.go
Confidence: HIGH for this repository
Corroborated By: Microsoft + Go singleflight conceptually.
Notes: Use as implementation hint for Engineer Agent, not as external research finding.
