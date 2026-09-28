# 06 — Open Questions

## Unanswered questions

1. **What is the exact formula notation in Vattani et al. (2015)?** The paper PDF could not be parsed. Wikipedia's pseudocode includes `-delta*beta*log(rand(0,1))`; the lab writes `Δ·β·ln(rand()) > TTL_remaining` without the minus. Which is canonical? Is beta=1 proven optimal or merely recommended? (MEDIUM confidence that beta=1 "works well" comes from Wikipedia; paper-level proof unknown.)
2. **How much early-rebuild traffic does XFetch generate under sustained high load?** The paper claims optimality; the quantitative trade-off (rebuild frequency vs. stampede reduction) was not extracted.
3. **What is the production stampede threshold for a given RDBMS / connection pool size?** No benchmark in opened sources. The lab's 10,000 RPS / 500 goroutine numbers are synthetic.
4. **Does `singleflight` alone suffice in a multi-process (N replicas) deployment?** It prevents in-process duplication but not cross-process. A Redis `SET NX PX` lock or a pub/sub broadcast would be needed; Redis-doc specifics were NOT VERIFIED (fetches 404'd). What is the minimal cross-process pattern?
5. **At what TTL / hit-rate / concurrency does jitter alone fail to prevent stampede?** Conceptually understood (it does not bound concurrency of one hot key), but no numeric envelope found in sources.
6. **What is the real-world amplification risk of request-triggered SWR revalidation on a heavily-repeated key?** RFC 5861 §5 warns against prefetch without request trigger; application-level SWR behavior under high fan-out has not been empirically evaluated here.
7. **Is write-behind safe for financial / payment-domain data?** Mechanism implies data-loss risk; safety depends on WAL, journaling, and crash-recovery — not addressed in opened sources.

## Weak evidence

- **Vattani et al. optimality claim.** Supported by DOI metadata + Wikipedia citation [Sources 04, 08]. Primary text not verified (PDF fetch failed). MEDIUM confidence.
- **Micro-benchmark figures (e.g., P99 improvement, DB connection count reduction).** Not sourced anywhere in opened materials; lab scenario is hypothetical [Evidence 18].
- **Jitter as stampede prevention.** Evidence for jitter desynchronizing clocks exists (Wikipedia thundering herd §Mitigation); direct link to stampede prevention at the cache layer is inferential, not backed by an opened source.

## Claims needing deeper research

1. **Redis / ElastiCache / Memcached-specific stampede APIs.** Attempted redis.io/learn/howtos and docs.aws.amazon.com/AmazonElastiCache/ paths returned 404/403. A direct fetch of the current URL structure (or search of the site map) would be needed. Claims about `SET NX PX` for distributed locks, `SCRIPT` + Lua scripts for conditional revalidation, and native TTL retrieval (`TTL key`) remain NOT VERIFIED in this pass.
2. **Production comparisons of XFetch vs. single-flight vs. SWR on hot keys.** No empirical A/B dataset found. A controlled experiment (matching the lab's 500-goroutine benchmark) is the missing step before engineering can recommend one pattern over the others for the specific workload.
3. **The "Thundering Herd" kernel-level optimization impact.** EPOLLEXCLUSIVE (Linux 4.5) and Windows IOCP mitigate the OS-layer problem; their impact on application-level cache stampede is an open analogy — whether the two are orthogonal deserves a dedicated look.

## Possible next research directions

- Reproduce the lab benchmark (500 concurrent goroutines on an expired key) and measure: (a) DB queries before/after singleflight, (b) P99 latency, (c) connection-pool utilization. Compare singleflight vs. XFetch at beta ∈ {1, 2, 3} vs. TTL-jitter-only baselines.
- Extend to multi-process (N Go servers behind a load balancer) and compare in-process singleflight vs. Redis SET NX PX distributed lock vs. XFetch (which scales horizontally without coordination) vs. SWR background revalidation.
- Evaluate write-behind durability guarantees for the lab's Dashboard model (invoice counts): would a write-ahead log or a commit on close solve the crash-before-flush risk identified in Evidence 3?
- Map RFC 5861's HTTP `Cache-Control` semantics to application-layer cache behavior (e.g., whether the lab's cache should attach a `revalidated_at` timestamp and enforce the `stale-while-revalidate` window explicitly).
- Re-fetch and parse Vattani et al. 2015 PDF from an alternative source (e.g., ACM Digital Library, Semantic Scholar, or university repository) to verify the optimality theorem statement and obtain quantitative bounds.
