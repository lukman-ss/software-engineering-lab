# 03 — Claim Audit

## Claim 1: Cache-Aside Read/Write Ordering
Claim: Cache-Aside loads on demand; writes update data store first and then invalidate the cache to minimize concurrent stale repopulation.
Location: `05-report.md: Finding 1`, `03-evidence.md: Evidence 1`
Evidence Provided: Microsoft Learn Azure Architecture Center guidance.
Source: Source 01 (Microsoft Learn)
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Well-documented pattern with clear ordering rationale.

---

## Claim 2: Write-Through vs Write-Behind Durability & Latency
Claim: Write-through synchronously updates cache and backing store (higher write latency, fresh read-after-write); Write-behind postpones backing store writes (high write throughput, data loss risk on cache node crash).
Location: `05-report.md: Findings 2 & 3`, `03-evidence.md: Evidence 3 & 4`
Evidence Provided: Wikipedia Cache write policies, Hennessy & Patterson reference, Microsoft Learn.
Source: Source 01, Source 06
Source Actually Supports Claim: YES
Classification: FACT / INTERPRETATION
Severity: LOW
Notes: Crash-before-flush risk for application-level DB sync is a sound logical deduction from deferred writes.

---

## Claim 3: Cache Stampede Failure Dynamics
Claim: Popular key expiry under high concurrency causes massive simultaneous recomputations, degrading hit rate to zero and exhausting DB connections.
Location: `05-report.md: Finding 4`, `03-evidence.md: Evidence 5 & 6`
Evidence Provided: Wikipedia Cache stampede mathematical illustration ($10\text{ req/s} \times 3\text{ s} = 30\text{ concurrent recomputes}$).
Source: Source 04
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Arithmetic behavior verified; research explicitly notes that 10,000 RPS / 500 goroutine scenarios in the lab prompt are synthetic educational parameters.

---

## Claim 4: Single-Flight Request Coalescing
Claim: Go `singleflight.Group.Do` coalesces duplicate concurrent in-flight calls per key, running only one execution while others wait and share the result.
Location: `05-report.md: Finding 5`, `03-evidence.md: Evidence 9`
Evidence Provided: Official Go `x/sync/singleflight` documentation.
Source: Source 03
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Accurately scoped to single-process; multi-process requires distributed locks or Redis equivalents.

---

## Claim 5: Probabilistic Early Expiration (XFetch) Mathematical Formulation
Claim: XFetch evaluates early recomputation when $- \Delta \cdot \beta \cdot \ln(rand()) > \text{TTL\_remaining}$. Lab prompt formulation missing minus sign is mathematically flawed for $rand \in (0,1)$.
Location: `05-report.md: Finding 6`, `04-contradictions.md: C1`
Evidence Provided: Vattani et al. formula via Wikipedia transcription, corroborated by algebraic derivation.
Source: Source 04, Source 08
Source Actually Supports Claim: YES
Classification: FACT
Severity: HIGH (Mitigated via Revision)
Notes: Revision explicitly placed a prominent WARNING block prohibiting the unnegated formula and documenting correct notation.

---

## Claim 6: Stale-While-Revalidate Behavior & Constraints
Claim: RFC 5861 permits serving stale data within an indicated window while revalidating asynchronously in background; validation should be request-triggered to prevent amplification attacks.
Location: `05-report.md: Finding 7`, `03-evidence.md: Evidence 12 & 13`
Evidence Provided: RFC 5861 §3, §3.1, §5.
Source: Source 02
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Primary standards text directly supports the claim.

---

## Claim 7: Jitter Sufficiency Limitations
Claim: Jitter desynchronizes expiration across multiple keys but does not bound concurrent requests to a single expired hot key.
Location: `05-report.md: Finding 8`, `03-evidence.md: Evidence 16`
Evidence Provided: Logical deduction based on hot-key contention mechanics.
Source: Source 05, Source 11
Source Actually Supports Claim: PARTIAL
Classification: INTERPRETATION
Severity: MEDIUM
Notes: Revised research accurately labeled this deduction as inferential rather than directly quoted literature.
