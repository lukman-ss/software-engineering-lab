# 04 — Contradictions

## Contradiction 1

Statement A: Lab prompt formula `Δ·β·ln(rand()) > TTL_remaining`

Location: `research/04-contradictions.md: C1`, `research/05-report.md: Finding 6`

Statement B: Wikipedia-transcribed XFetch formulation: `(time() - delta*beta*log(rand(0,1))) ≥ expiry`, which rearranges to `-Δ·β·ln(rand()) > TTL_remaining`

Location: `research/03-evidence.md: Evidence 10`

Type: INTERNAL

Impact: HIGH — as written in the lab prompt, the formula is always false for rand∈(0,1), producing a dead branch that never triggers early refresh. Educational lab would teach incorrect algorithm.

Assessment: RESOLVED — research-revision added prominent WARNING block with correct formula.

---

## Contradiction 2

Statement A: RFC 5861 header: "not endorsed by the IETF and has no formal standing in the IETF standards process" (Informational / Independent Submission).

Location: `research/04-contradictions.md: C3`

Statement B: Common informal usage treating RFC 5861 as an "Internet Standard."

Location: community materials / casual usage

Type: SOURCE_CONFLICT

Impact: LOW — research correctly identifies the RFC's actual informational status.

Assessment: RESOLVED — research explicitly states "Informational RFC" consistently.

---

## Contradiction 3

Statement A: Wikipedia Cache stampede: locking "requires an extra write for the locking mechanism, doubling the number of writes."

Location: `research/04-contradictions.md: C4`

Statement B: Go `singleflight` uses in-process mutex, no external lock write needed.

Location: `research/04-contradictions.md: C4`

Type: SCOPE CLARIFICATION (not a true contradiction)

Impact: MEDIUM — confusion risk if reader conflates in-process with distributed mechanisms.

Assessment: RESOLVED — research explicitly differentiates scope (in-process singleflight vs. Redis SET NX PX distributed lock).

---

## Contradiction 4

Statement A: RFC 5861 §3.1: revalidation triggered by an incoming request; §5 warns against validation without request trigger to avoid amplification.

Location: `research/04-contradictions.md: C5`

Statement B: Lab spec describes "asynchronous background job" for SWR revalidation.

Location: `research/04-contradictions.md: C5`, `05-report.md: Finding 7`

Type: CODE_DOC_MISMATCH (research-level)

Impact: MEDIUM — amplification risk if misimplemented as unconditional background cron.

Assessment: RESOLVED — research notes tension and recommends request-correlated implementation aligned with RFC §5.

---

## No other material contradictions

All other areas — cache-aside ordering, singleflight semantics, XFetch formulation, TTL trade-offs, jitter purpose, stale-if-error — are internally consistent across sources.
