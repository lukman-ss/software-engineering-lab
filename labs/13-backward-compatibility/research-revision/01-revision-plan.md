# Revision Plan

Target Lab: labs/13-backward-compatibility
Previous Audit Status: NEEDS_REVISION (research-audit/07-verdict.md)

## Blocking Issues

1. **HIGH** — Evidence 7 in `research/03-core-concepts.md` cites Martin Fowler Feature Toggle article for dual-write risks. Source discusses toggle validation complexity, not database dual-write failure modes.

## Non-Blocking Issues

1. **MEDIUM** — Evidence 5 (PostgreSQL `CONCURRENTLY`, additive changes) presented as universal migration strategy; actually PostgreSQL-specific.
2. **MEDIUM** — Evidence 8 (backfill batching/throttling/idempotency) partially supported; PostgreSQL only implies performance issues, doesn't explicitly define patterns.
3. **MEDIUM** — Evidence 6 (dual read fallback) partial support; Fowler doesn't explicitly show code fallback during backfill.
4. **LOW** — Evidence 3 (30-day heuristic) already marked `NOT VERIFIED` in prior revision.

## Files To Modify

- `research/03-core-concepts.md` — Fix Evidence 7 source; narrow Evidence 5, 6, 8 claims
- `research/08-failure-modes.md` — Fix Evidence cite for dual-write failure mode
- `research/02-sources.md` — Add Source 10 (Microservices.io/Transactional Outbox)

## Verification Plan

- source verification: verify microservices.io transactional outbox page describes dual-write risks
- research/03-core-concepts.md: Evidence 7 now cites correct source; Evidence 5 labeled PostgreSQL-specific; Evidence 6 confidence reduced; Evidence 8 narrowed to "implied by PostgreSQL lock behavior"
- research/08-failure-modes.md: Line 29 Evidence now cites microservices.io outbox pattern
- research/02-sources.md: New Source 10 added and verified
- original sources 1-9 remain valid and pass

## Changes Made

- `research/03-core-concepts.md` Evidence 7: replaced Fowler Feature Toggle source with Microservices.io Transactional Outbox
- `research/08-failure-modes.md` Failure Mode 2: replaced evidence with Transactional Outbox pattern
- `research/02-sources.md`: added Source 10 (Microservices.io Transactional Outbox)
- `research/03-core-concepts.md` Evidence 5: added PostgreSQL-specific caveat
- `research/03-core-concepts.md` Evidence 8: narrowed claim and marked LOW confidence
- `research/03-core-concepts.md` Evidence 6: reduced confidence to LOW

## Verification Results

| Check | Result |
|---|---|
| Evidence 7 source matches claim | PASS (verified microservices.io page supports dual-write/2PC/outbox description) |
| Failure Mode 2 evidence | PASS (now cites transactional outbox, not feature toggles) |
| Source 10 reachable | PASS |
| Evidence 5 caveat added | PASS |
| Evidence 6 confidence | PASS (LOW — inferential) |
| Evidence 8 confidence | PASS (LOW — implies only) |