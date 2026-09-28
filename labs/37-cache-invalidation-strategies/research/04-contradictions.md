# 04 — Contradictions / Tensions

## C1 — Lab formula vs Wikipedia / paper formulation (probabilistic early expiration sign)

SOURCE A: Lab specification `Δ · β · ln(rand()) > TTL_remaining`, with text "Semakin dekat waktu expired dan semakin lama durasi query DB, semakin tinggi peluang satu worker me-refresh cache lebih awal."

SOURCE B: Wikipedia Cache stampede (quoting Vattani et al. PVLDB 2015):
`(time() - delta * beta * log(rand(0,1))) ≥ expiry`
Rearranged: `-delta * beta * log(rand(0,1)) ≥ TTL_remaining`, i.e., `-Δ·β·ln(rand) ≥ TTL_remaining`.

ASSESSMENT:
`log(rand(0,1))` is negative (rand in (0,1)), so `-Δ·β·log(rand)` is a positive offset. The lab's `Δ·β·ln(rand()) > TTL_remaining` would be negative > positive for rand∈(0,1), never true. The qualitative description ("higher probability closer to expiry / longer query time") is consistent across both sources, suggesting the lab intends the same exponential-threshold rule but omits the minus sign or implicitly assumes `rand()` returns a value > 1 (unusual). Two interpretations are plausible: (a) formula should read `-Δ·β·ln(rand()) > TTL_remaining`, (b) `ln(rand())` should be `ln(1/rand()) = -ln(rand())`. The paper PDF was fetched but unparseable; optimality proof and exact formula notation cannot be confirmed from primary text. HIGH confidence that as written the lab inequality is directionally wrong for rand∈(0,1). MEDIUM confidence about author intent (qualitative part aligns).

Implication for lab: Engineer Agent should prefer the Wikipedia-formulated condition (or equivalently `-Δ·β·ln(U) > TTL_remaining`) and note the sign convention in implementation comments.

## C2 — Thundering herd vs cache stampede terminology

SOURCE A: Wikipedia Thundering herd: "When the thundering herd problem occurs while attempting to access a cache, it is often referred to as a cache stampede." Distinguishes OS-level "many waiters, one wakes" from cache-level "many recomputes".

SOURCE B: Lab title "Cache Invalidation Strategies --- Memilih Pola Caching dan Mencegah Thundering Herd". Lab section headings use "Cache Stampede (Thundering Herd)" and "Solusi Menghadapi Cache Stampede".

ASSESSMENT: Not a contradiction — a taxonomy overlap. Thundering herd is broader (any resource acquisition contention after a wake event); cache stampede is the cache-miss recomputation case. Using them interchangeably in caching contexts is common (Wikipedia acknowledges this). For lab pedagogy it is acceptable to treat them as synonymous in the cache-access context, but the Engineer should keep the distinction in tests/comments.

## C3 — RFC 5861 status claim

SOURCE A: RFC 5861 header text: "not endorsed by the IETF and has no formal standing in the IETF standards process", "Independent Submission", "Informational".

SOURCE B: Common marketing/casual usage calls RFC 5861 an "Internet standard" for stale-while-revalidate.

ASSESSMENT: The RFC's own status block is unambiguous. Call it an IETF-published informational RFC, not a standards-track RFC. Implementation existence (CDNs, browsers, proxies) is not asserted here because no source was opened to corroborate that.

## C4 — Locking "extra write / doubling writes" vs Go singleflight (no lock key)

SOURCE A: Wikipedia Cache stampede Locking section: "requires an extra write for the locking mechanism... doubling the number of writes".

SOURCE B: Go singleflight package: in-process mutex-based dedup; no separate cache key write for the lock; duplicate callers wait on the same goroutine and receive the same return value. Labs/04-caching/stampede.go `ProtectedStampedeService` uses this in-process.

ASSESSMENT: Not a contradiction — different mechanisms at different scopes. A distributed lock (Redis SET NX PX) adds the extra write; an in-process `singleflight` group does not touch the cache for synchronization. The lab exercise simulates "500 concurrent goroutine" (in-process), so `singleflight` is the correct primitive; the "doubling writes" cost applies only if the mitigation switches to a distributed lock. Explicitly state scope in the report so readers do not conflate the two.

## C5 — Stale-while-revalidate trigger semantics

SOURCE A: RFC 5861 §3: "If a cached response is served stale due to the presence of this extension, the cache SHOULD attempt to revalidate it while still serving stale responses (i.e., without blocking)." §3.1: "asynchronous validation will only happen if a request occurs after the response has become stale, but before the end of the stale-while-revalidate window." Also §5 security note: "suggested that such validation be predicated upon an incoming request, to avoid the possibility of an amplification attack".

SOURCE B: Lab wording "sambil memicu asynchronous job untuk update data baru di latar belakang" reads like an independent background job rather than a request-triggered revalidation.

ASSESSMENT: Apparent tension. Application-level SWR can be implemented either way (request-triggered vs. true background job). The RFC explicitly warns against non-request-triggered refetches (amplification risk). The lab's "asynchronous job" is acceptable at the application tier, but the Engineer should preserve the request-correlation so that a stale response's revalidation is tied to an incoming client request, not a fire-and-forget cron. This keeps the design aligned with RFC §5 intent.

## C6 — Write-through freshness guarantee vs practical ordering

SOURCE A: Microsoft: "write-through caching... updates the data store and the cache in the same write operation so that readers see the new value immediately after a successful write."

SOURCE B: Labs/04-caching/write_through.go comment: "PENTING: Database dan Redis adalah sistem TERPISAH. Tidak ada atomic commit lintas-keduanya." and "Best-effort update cache."

ASSESSMENT: Microsoft's "same write operation" is application-level sequential write (first DB, then cache). It is not a single ACID transaction across both systems. WriteThroughService even treats cache update as best-effort — if cache.Set fails, the DB update still succeeds and TTL acts as safety net. The Engineer should not read "same write operation" as distributed transactional atomicity; the lab's own implementation documents the boundary. Both statements are consistent once "operation" is read as "application transaction" rather than "database-level single transaction."

## No material contradictions (other categories)

No disagreement found between sources on:
- Mechanism of cache-aside (Sources 01, 04, 10 agree).
- Definition of stampede/dog-pile (Sources 04, 05).
- Go singleflight semantics (Source 03 + local code).
- XFetch exponential distribution and beta role (Source 04 quoting Source 08 metadata; Source 09 PDF unverifiable).
- RFC 5861 stale-if-error behavior (Source 02).
- TTL trade-offs (Source 01, 07).
- Jitter purpose (Source 05 + Source 10).
