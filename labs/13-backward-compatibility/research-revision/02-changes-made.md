# Changes Made

## Revision 1 — Fix Evidence 7 source mismatch (HIGH — blocking)

Audit Issue:
HIGH — Evidence 7 in `research/03-core-concepts.md` cited Martin Fowler Feature Toggle article for dual-write risks. Source discusses toggle validation complexity, not database dual-write failure modes.

Files Changed:
- `research/03-core-concepts.md` — Evidence 7
- `research/08-failure-modes.md` — Failure Mode 2 Evidence
- `research/02-sources.md` — added Source 10

Action:
- Evidence 7: replaced unsupported source with Microservices.io Transactional Outbox; added direct quotes about 2PC impossibility, transaction commit guarantees, and message relay idempotency
- Failure Mode 2: replaced Martin Fowler Feature Toggle evidence with Transactional Outbox pattern quotes
- Source 10: added Microservices.io Transactional Outbox as Tier 1 source, explicitly linking to Evidence 7 and Failure Mode 2

Verification:
- Source verification: Microservices.io page exists and supports dual-write problem description
- Evidence 7: source now correctly supports claim about dual-write atomicity failure modes
- Evidence 5: PostgreSQL-specific behavior now explicitly caveated
- Evidence 6: confidence reduced to LOW (inferential, no explicit code example)
- Evidence 8: claim narrowed to "umum dilakukan" and labeled LOW confidence

Status:
RESOLVED

---

## Revision 2 — Narrow PostgreSQL-specific claims (MEDIUM)

Audit Issue:
MEDIUM — Evidence 5 overgeneralizes PostgreSQL `CONCURRENTLY` behavior as universal migration rule.

Action:
- Evidence 5: Added explicit caveat: "`CONCURRENTLY`, non-rewriting defaults, and `NOT VALID` are PostgreSQL-specific behaviors, not a universal database rule"
- Noted MySQL tooling examples (`pt-online-schema-change`, `gh-ost`) as illustrative only

Status:
RESOLVED

---

## Revision 3 — Evidence 6 and 8 confidence adjustments (MEDIUM/LOW)

Action:
- Evidence 6: Reduced confidence to LOW (inferential; Fowler doesn't show explicit fallback read code)
- Evidence 8: Claim changed to "umum dilakukan" with LOW confidence; source only implies performance, not explicit patterns

Status:
RESOLVED