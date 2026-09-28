# 04 — Contradiction Audit

## Contradiction 1

Statement A:
Lab specification formula: `Δ · β · ln(rand()) > TTL_remaining`
Location: `research/04-contradictions.md`, C1 (citing lab specification)

Statement B:
Wikipedia / Vattani et al. formula: `(time() - delta * beta * log(rand(0,1))) ≥ expiry`, rearranged as `-Δ·β·ln(rand(0,1)) ≥ TTL_remaining`
Location: `research/04-contradictions.md`, C1 (citing Wikipedia Cache stampede)

Type:
FORMULA_ERROR / SPEC_MISMATCH

Impact:
HIGH — If an engineer implements `Δ · β · ln(rand()) > TTL_remaining` verbatim where `rand() ∈ (0,1)`, `ln(rand())` evaluates to a negative number, resulting in `negative > positive` which is ALWAYS FALSE. Early refresh will NEVER be triggered, silently failing the probabilistic expiration mechanism.

Assessment:
Auditor confirms the Research Agent's contradiction analysis is correct. The lab formula contains a missing negative sign or requires `ln(1/rand())`. The research correctly identified this bug and documented it as a blocker/finding for engineering.

---

## Contradiction 2

Statement A:
Informational header on RFC 5861: "This RFC was published on the Independent Submission stream... not endorsed by the IETF and has no formal standing in the IETF standards process."
Location: `RFC 5861` status block

Statement B:
Common industry material / blog descriptions calling RFC 5861 an "Internet Standard for stale-while-revalidate".
Location: `research/04-contradictions.md`, C3

Type:
STATUS_MISMATCH

Impact:
LOW — Standard informational accuracy issue. RFC 5861 is an informational extension, not an IETF Standards Track RFC.

Assessment:
Auditor verified RFC 5861 header via web fetch. The document category is `Informational` and the stream is `Independent Submission`. Calling it an IETF standard is technically incorrect. Research Agent correctly recorded this distinction.

---

## Contradiction 3

Statement A:
Wikipedia Cache stampede Locking section: "requires an extra write for the locking mechanism... doubling the number of writes".
Location: `research/04-contradictions.md`, C4 (citing Wikipedia)

Statement B:
Go `singleflight` package: in-process mutex-based dedup; no separate cache write for lock management.
Location: `research/04-contradictions.md`, C4 (citing `pkg.go.dev/golang.org/x/sync/singleflight`)

Type:
SCOPE_CONFUSION

Impact:
MEDIUM — Readers conflating distributed locking (Redis SET NX PX) with in-process singleflight might miscalculate operational write costs.

Assessment:
Auditor confirms research resolution: distributed locks add cache write overhead; in-process singleflight adds local memory/mutex cost with ZERO cache write overhead. Explicit scope clarification in the report correctly resolves this tension.

---

## Contradiction 4

Statement A:
RFC 5861 §5 security guidance: validation SHOULD be request-triggered to avoid amplification/prefetch attacks.
Location: `research/04-contradictions.md`, C5 (citing RFC 5861 §5)

Statement B:
Lab description: "sambil memicu asynchronous job untuk update data baru di latar belakang".
Location: `research/04-contradictions.md`, C5 (citing lab specification)

Type:
DESIGN_TENSION

Impact:
MEDIUM — An unconstrained background refresh loop risks amplification attacks if detached from user request traffic.

Assessment:
Auditor verified RFC 5861 §5 text: "suggested that such validation be predicated upon an incoming request, to avoid the possibility of an amplification attack". Research correctly flagged that SWR implementation should tie background jobs to incoming request triggers rather than running autonomous periodic background jobs.

---

## Contradiction 5

Statement A:
Microsoft Learn: "write-through caching... updates the data store and the cache in the same write operation".
Location: `research/04-contradictions.md`, C6 (citing Microsoft Learn)

Statement B:
Application reality: Database and Redis are separate network services without an atomic cross-system commit (two-phase commit/distributed ACID transaction).
Location: `research/04-contradictions.md`, C6 (citing `labs/04-caching/write_through.go`)

Type:
TERMINOLOGY_OVERGENERALIZATION

Impact:
MEDIUM — "Same write operation" in documentation could be misunderstood as ACID atomicity across store and cache.

Assessment:
Research correctly resolves this: vendor documentation means application-level sequential writes, not distributed ACID atomicity. If cache.Set fails post DB commit, staleness/inconsistency can still occur until TTL expires.
