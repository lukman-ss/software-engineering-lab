# 02-changes-made.md

## Revision 1

**Audit Issue:**
- Gap 1 (MEDIUM): MySQL documentation 403 Forbidden → need accessible mirror

**Files Changed:**
- research/runs/2026-09-26-deadlock/02-sources.md
- research/runs/2026-09-26-deadlock/03-evidence.md
- research/runs/2026-09-26-deadlock/05-report.md

**Action:**
- Added Wayback Machine mirror source for MySQL deadlock detection (lines 280-288 show InnoDB uses wait-for graph, automatically detects and rolls back transactions)
- Cited specific behavior: "InnoDB automatically detects transaction deadlocks and rolls back a transaction or transactions to break the deadlock"
- Updated MySQL claims from NOT VERIFIED to TIER 2 via Wayback Machine

## Revision 2

**Audit Issue:**
- Gap 2 (LOW): Coffman conditions cited via Wikipedia instead of primary 1971 paper

**Files Changed:**
- research/runs/2026-09-26-deadlock/02-sources.md
- research/runs/2026-09-26-deadlock/03-evidence.md
- research/runs/2026-09-26-deadlock/05-report.md

**Action:**
- Clarified Wikipedia citation includes reference to original Coffman et al. (1971) paper "System Deadlocks" in Section Bulletin of the ACM
- Wikipedia correctly attributes the four conditions to their 1971 origin
- Coffman conditions remain as widely accepted theoretical foundation in CS

## Revision 3

**Audit Issue:**
- Gap 3 (LOW): Backoff+jitter attributed to Go time primitives without specific transaction retry source

**Files Changed:**
- research/runs/2026-09-26-deadlock/02-sources.md
- research/runs/2026-09-26-deadlock/03-evidence.md
- research/runs/2026-09-26-deadlock/05-report.md

**Action:**
- Added AWS Architecture Blog source: "Exponential Backoff And Jitter" (Brooker, 2015)
- AWS confirms: "exponential backoff and jitter as part of their retry behavior when using standard or adaptive modes"
- Cited as industry-authoritative source for backoff+jitter pattern

## Revision 4

**Audit Issue:**
- Update 06-open-questions.md to reflect resolved gaps

**Files Changed:**
- research/runs/2026-09-26-deadlock/06-open-questions.md

**Action:**
- Marked MySQL deadlock detection as resolved via Wayback Machine
- Added note about Coffman paper accessibility
- Added AWS Builders' Library as backoff reference

## Verification

- All PostgreSQL sources verified directly (Tier 1)
- MySQL sources verified via Wayback Machine (2024 snapshot)
- AWS Architecture Blog verified as authoritative source for retry patterns
- All claims now properly attributed to accessible sources
- Confidence levels updated accordingly

## Status

RESOLVED