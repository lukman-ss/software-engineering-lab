# Audit Plan

Target Lab: labs/18-deadlock (research-only phase)

Files Reviewed:
- research/runs/2026-09-26-deadlock/01-plan.md
- research/runs/2026-09-26-deadlock/02-sources.md
- research/runs/2026-09-26-deadlock/03-evidence.md
- research/runs/2026-09-26-deadlock/04-contradictions.md
- research/runs/2026-09-26-deadlock/05-report.md
- research/runs/2026-09-26-deadlock/06-open-questions.md

Claims To Verify:
1. Deadlock definition and Coffman conditions.
2. PostgreSQL automatic deadlock detection and victim selection.
3. PostgreSQL error codes 40P01 and 40001 requiring retry.
4. Difference between deadlock_timeout and lock_timeout.
5. Lock ordering prevents deadlocks.
6. Transaction duration and contention risk (no external calls in transactions).
7. Deadlock handling: idempotent retry with backoff and jitter in Go.
8. Observability metrics for deadlocks in PostgreSQL.
9. Isolation level behavior related to deadlocks and serialization failures.

Code To Execute:
None (Pipeline override: Audit research only).

Primary Risks:
- Source hallucination (e.g., MySQL references yielding 403).
- Overgeneralization of PostgreSQL behavior to other systems.
- Misinterpreting general application patterns (like PPOB specific external call handling) as database absolute facts.

Audit Strategy:
- Validate PostgreSQL documentation URLs.
- Cross-reference claims against verbatim evidence.
- Ensure research distinguishes between generic engineering practices and database-specific facts.
- Verify contradictions and gaps are reported correctly.
