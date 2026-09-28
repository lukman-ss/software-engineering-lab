# 03 — Claim Audit

Target Lab: labs/37-cache-invalidation-strategies

## Claim 1: Cache-Aside Read and Invalidation Sequencing
Claim: Cache-Aside loads data on demand upon a cache miss; updates write to the database first and then invalidate/delete the cache entry.
Location: `05-report.md` Finding 1; `03-evidence.md` Evidence 1.
Evidence Provided: Microsoft Learn Cache-Aside Pattern doc; local `cache_aside.go`.
Source: Source 01 (Microsoft Learn).
Source Actually Supports Claim: YES
Classification: FACT / BEST PRACTICE
Severity: LOW
Notes: Correctly emphasizes ordering: store update before cache deletion to avoid race windows where stale data is re-cached.

---

## Claim 2: Write-Through Immediate Read-After-Write Semantics
Claim: Write-through updates cache and database synchronously during write path, providing immediate read-after-write freshness within the process, but does not provide distributed ACID transaction guarantees across separate DB and Redis instances.
Location: `05-report.md` Finding 2; `04-contradictions.md` C6.
Evidence Provided: Microsoft Learn; Wikipedia Write policies; `labs/04-caching/write_through.go`.
Source: Source 01, Source 06, Source 11.
Source Actually Supports Claim: YES
Classification: FACT / IMPLEMENTATION-SPECIFIC
Severity: LOW
Notes: Properly clarifies the boundary between application-level sequential dual writes and atomic distributed transactions.

---

## Claim 3: Write-Behind Asynchronous Durability Risk
Claim: Write-behind buffers writes in cache and flushes to database asynchronously, yielding high write throughput at the cost of potential data loss if the cache node crashes before flush.
Location: `05-report.md` Finding 3; `03-evidence.md` Evidence 3.
Evidence Provided: Wikipedia Cache (computing) write-back policy; database caching architectural deduction.
Source: Source 06.
Source Actually Supports Claim: YES
Classification: FACT / INTERPRETATION
Severity: LOW
Notes: Properly identifies durability trade-offs in distributed caching tiers.

---

## Claim 4: Cache Stampede Cascading Collapse Mechanics
Claim: Expiration of a hot key under high concurrency causes multiple threads to simultaneously query the backing store, which can exhaust DB connection pools, elevate P99 latency, and cause congestion collapse (hit rate dropping to zero).
Location: `05-report.md` Finding 4; `03-evidence.md` Evidence 5, 6.
Evidence Provided: Wikipedia Cache stampede (citing Vattani et al. 2015).
Source: Source 04, Source 08.
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Synthetic lab numbers (10,000 RPS, 500 goroutines) are correctly labeled as exercise parameters, not empirical benchmarks.

---

## Claim 5: Single-Flight Duplicate Call Suppression
Claim: Go `golang.org/x/sync/singleflight` coalesces concurrent in-flight executions for the same key within a process so only one goroutine executes the recompute function while others wait and share the result.
Location: `05-report.md` Finding 5; `03-evidence.md` Evidence 9.
Evidence Provided: Go pkg.go.dev singleflight documentation.
Source: Source 03.
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Scope accurately delimited to in-process concurrency; does not replace distributed locks across multiple independent node processes.

---

## Claim 6: Correct Formulation of Probabilistic Early Expiration (XFetch)
Claim: XFetch evaluates `-Δ·β·ln(rand()) > TTL_remaining` (where `rand() ∈ (0,1)`), correctly using the negative logarithm to generate a positive time offset for probabilistic early recomputation.
Location: `05-report.md` Finding 6 (Warning block); `04-contradictions.md` C1; `03-evidence.md` Evidence 10.
Evidence Provided: Wikipedia Cache stampede citing Vattani et al. 2015.
Source: Source 04, Source 08.
Source Actually Supports Claim: YES
Classification: FACT / MATHEMATICAL CORRECTION
Severity: LOW (Previously CRITICAL, now resolved by explicit warning and mathematical derivation)
Notes: Research report prominently warns against implementing the unnegated formula `Δ·β·ln(rand()) > TTL_remaining`.

---

## Claim 7: Stale-While-Revalidate Asynchronous Refresh Semantics
Claim: RFC 5861 defines `stale-while-revalidate` permitting caches to serve stale content for delta seconds while revalidating asynchronously; RFC §5 recommends validation be request-triggered to prevent amplification attacks.
Location: `05-report.md` Finding 7; `03-evidence.md` Evidence 12, 13; `04-contradictions.md` C5.
Evidence Provided: RFC 5861 §3, §3.1, §4, §5.
Source: Source 02 (RFC 5861).
Source Actually Supports Claim: YES
Classification: FACT / STANDARD (INFORMATIONAL)
Severity: LOW
Notes: Accurately reflects RFC 5861 text and security considerations.

---

## Claim 8: Jitter as Anti-Synchronization vs Single Hot-Key Stampede
Claim: TTL jitter desynchronizes expiration clocks across multiple keys, but does not cap concurrent rebuilds when a single hot key expires.
Location: `05-report.md` Finding 8; `03-evidence.md` Evidence 16.
Evidence Provided: Wikipedia Thundering herd §Mitigation; single-key contention mechanics.
Source: Source 05, Source 11.
Source Actually Supports Claim: YES
Classification: INTERPRETATION / LOGICAL DEDUCTION
Severity: LOW
Notes: Clearly labeled as an inferential deduction in the revised report.
