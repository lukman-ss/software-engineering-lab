# Audit Plan: Research on Optimistic vs Pessimistic Locking

## Target Lab
`labs/23-optimistic-vs-pessimistic-locking`

## Scope
Research files audit only (per pipeline override):
- `research/01-plan.md`
- `research/02-sources.md`
- `research/03-evidence.md`
- `research/04-contradictions.md`
- `research/05-report.md`
- `research/06-open-questions.md`

## Files Reviewed
- `research/01-plan.md` (117 lines)
- `research/02-sources.md` (173 lines)
- `research/03-evidence.md` (291 lines)
- `research/04-contradictions.md` (90 lines)
- `research/05-report.md` (170 lines)
- `research/06-open-questions.md` (89 lines)

## Claims To Verify
1. Definition and reproduction of the "lost update" anomaly under default isolation (READ COMMITTED).
2. Mechanism of pessimistic locking via `SELECT ... FOR UPDATE` (row-level exclusive locks, blocking behavior, release at transaction boundary).
3. Mechanism of optimistic locking (version/timestamp guard in `WHERE` clause, conflict detection via `affected_rows == 0` or exception, requirement for explicit application retry/409 handling).
4. Trade-offs: Pessimistic locking reduces concurrency and risks deadlocks; optimistic locking wastes work upon conflict.
5. Selection criteria: Pessimistic for high contention/correctness-critical; optimistic for low contention/read-heavy/offline workflows.
6. Atomic single-statement update (`UPDATE ... SET stock = stock - N WHERE stock >= N`) as an alternative.
7. Database-specific default isolation levels (PostgreSQL vs MySQL vs Oracle) and isolation differences.
8. Common anti-patterns (assuming transaction alone prevents lost update, holding locks across network calls, ignoring 0 affected rows).

## Code To Execute
Skipped per pipeline override ("Audit research only. Do not audit implementation/code in this stage.").

## Primary Risks
1. Verification of external links (e.g. MySQL 8.0 manual URLs returning 403 or requiring mirrors).
2. Distinction between verified primary source statements and synthesized best practices (e.g., atomic decrement WHERE pattern).
3. Completeness of contradiction and gap analysis.

## Audit Strategy
1. Perform web fetch verification on critical primary sources (PostgreSQL official docs, Fowler EAA patterns, Wikipedia concurrency control).
2. Inspect source metadata (titles, publishers, URLs, relevance, tier classification).
3. Cross-reference all claims in `05-report.md` and `03-evidence.md` against sources.
4. Verify contradiction analysis in `04-contradictions.md`.
5. Verify open questions and gap documentation in `06-open-questions.md`.
6. Issue final verdict in `07-verdict.md`.
