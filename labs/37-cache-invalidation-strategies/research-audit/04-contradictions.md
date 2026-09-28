# 04 — Contradictions / Tensions Audit

Target Lab: labs/37-cache-invalidation-strategies

## Contradiction 1: XFetch Formula Sign in Lab Prompt vs Mathematical Reality

Statement A (Lab Prompt):
`Δ · β · ln(rand()) > TTL_remaining`

Statement B (Mathematical Reality / Wikipedia PVLDB 2015):
`-Δ · β · ln(rand()) > TTL_remaining` (or `(time() - delta * beta * log(rand(0,1))) ≥ expiry`)

Location: `04-contradictions.md` C1; `05-report.md` Finding 6 Warning.
Type: INTERNAL / SPECIFICATION_ERROR
Impact: Critical if uncaught (formula would never trigger early refresh for `rand() ∈ (0,1)`).
Assessment: RESOLVED. The research report contains an explicit, prominent WARNING block explaining why the unnegated formula is mathematically broken and directing the engineering phase to use the negated formula `-Δ·β·ln(rand()) > TTL_remaining`.

---

## Contradiction 2: RFC 5861 "Standard" vs "Informational" Status

Statement A:
Informal engineering references often refer to RFC 5861 as an "IETF standard".

Statement B:
RFC 5861 header states: "Published on the Independent Submission stream. This RFC is not endorsed by the IETF and has no formal standing in the IETF standards process."

Location: `04-contradictions.md` C3; `05-report.md` Finding 7.
Type: SOURCE_CONFLICT
Impact: Minor nuance regarding formal standardization.
Assessment: RESOLVED. Report and evidence correctly designate RFC 5861 as an Informational RFC / Independent Submission.

---

## Contradiction 3: Write-Through "Same Write Operation" vs Distributed Realities

Statement A (Vendor Docs):
Write-through updates data store and cache "in the same write operation".

Statement B (Distributed Systems Fact):
Database and cache (e.g. Postgres and Redis) are independent systems without 2PC/distributed transaction support in standard web stacks.

Location: `04-contradictions.md` C6; `05-report.md` Finding 2.
Type: INTERNAL / SCOPE_CLARIFICATION
Impact: Engineers could falsely assume transactional atomicity across DB + Redis.
Assessment: RESOLVED. Research explicitly documents that "same write operation" means application-level sequential writing with best-effort cache set and TTL safety net.

---

## Contradiction 4: Request-Triggered vs Background Cron SWR

Statement A (RFC 5861 §5):
Revalidation should be predicated upon an incoming request to prevent amplification attacks.

Statement B (Loose SWR Descriptions):
SWR described as an independent periodic background worker refreshing expired keys.

Location: `04-contradictions.md` C5; `05-report.md` Finding 7.
Type: ARCHITECTURAL_DESIGN_TENSION
Impact: Request fan-out vs worker complexity.
Assessment: RESOLVED. Research clarifies that request-triggered revalidation is the canonical RFC 5861 model, while independent background refresh represents the external recomputation pattern.
