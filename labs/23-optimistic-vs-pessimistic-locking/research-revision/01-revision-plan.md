# Revision Plan

Target Lab: labs/23-optimistic-vs-pessimistic-locking

Previous Audit Status: NEEDS_REVISION

## Blocking Issues

1. [HIGH] No research output — directory contains only `research/01-plan.md`; 15 questions unanswered, no evidence, no excerpts. (Gap 1-2)
2. [HIGH] No verifiable sources — 8 expected sources are aspirational only, no URLs, no reachability check, no scope assessment. (Source audit; Gap 1)
3. [HIGH] Isolation/MVCC/locking interactions and version-specific behavior flagged as risks in plan but no mitigation/evidence provided. (Gap 3)
4. [MEDIUM] Universal framing risk (Q2, Q9) without vendor/version/isolation qualification. (Gap 4)
5. [MEDIUM] Academic sources lack edition/year; insufficient alone for modern implementation claims. (Gap 5)
6. [MEDIUM] Performance/decision-threshold guidance (Q4, Q9) has no benchmark methodology or sourced numeric recommendation. (Gap 6)
7. [MEDIUM] Implementation patterns (Q5-Q8: version columns, SELECT FOR UPDATE/SHARE, atomic UPDATE WHERE, ORM support) have no SQL/ORM doc citations. (Gap 7)

## Non-Blocking Issues

1. [LOW] Plan structure sound; search strategy tiers appropriate if actually executed.
2. No README at lab root — not blocking for research-only stage but required before publication.

## Files To Modify

- labs/23-optimistic-vs-pessimistic-locking/research/02-sources.md (NEW)
- labs/23-optimistic-vs-pessimistic-locking/research/03-evidence.md (NEW)
- labs/23-optimistic-vs-pessimistic-locking/research/05-report.md (NEW)
- labs/23-optimistic-vs-pessimistic-locking/research/06-open-questions.md (NEW)

## Verification Plan

- source verification: verify all URLs reachable, excerpts match claims
- classification: each finding classified as FACT / INTERPRETATION / EXAMPLE / HYPOTHESIS / IMPLEMENTATION-SPECIFIC
- scope: all implementation claims scoped by DB/version/isolation level
- numeric claims: any thresholds anchored to vendor recommendation or benchmark with URL
- academic sources: edition/year cited
- no unsupported universal claims