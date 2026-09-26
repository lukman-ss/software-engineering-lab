# Audit Plan: Deadlock Research

## Target Lab
`labs/18-deadlock`

## Files Reviewed
- `labs/18-deadlock/research/01-plan.md`
- `labs/18-deadlock/research/02-sources.md`
- `labs/18-deadlock/research/03-evidence.md`
- `labs/18-deadlock/research/04-contradictions.md`
- `labs/18-deadlock/research/05-report.md`
- `labs/18-deadlock/research/06-open-questions.md`

## Claims To Verify
1. Deadlock requires all 4 Coffman conditions simultaneously (mutual exclusion, hold-and-wait, no preemption, circular wait).
2. Two individually correct transactions can deadlock solely due to inverted lock ordering.
3. RDBMS (PostgreSQL, MySQL/InnoDB, SQL Server) detect deadlocks automatically and abort a victim transaction.
4. Transaction duration directly increases deadlock probability by holding locks longer.
5. Consistent lock ordering is the primary architectural mitigation.
6. Deadlock resolution at the application level requires idempotent retry with exponential backoff and jitter.
7. Deadlock is fundamentally distinct from lock timeout.
8. Engine-specific defaults and behaviors (e.g., PostgreSQL `deadlock_timeout` default 1s, MySQL 1213 error code, SQL Server 1205 error code).

## Code To Execute
None. Pipeline override specifies: "Audit research only. Do not audit implementation/code in this stage."

## Primary Risks
- Use of secondary/tertiary summaries (Wikipedia) for fundamental academic concepts (Coffman conditions, 2PL).
- Use of archived Wayback Machine URLs for MySQL documentation due to direct fetch blocks.
- Potential conflation of implementation-specific behavior (e.g., PostgreSQL `deadlock_timeout`) with generic RDBMS behavior.
- Contextual application to PPOB (Indonesian payment aggregator) being speculative or unbacked by primary domain literature.

## Audit Strategy
1. Audit all 12 listed sources for reachability, tier classification, scope, and evidence fidelity.
2. Cross-examine claims in `05-report.md` and `03-evidence.md` against original text and database standards.
3. Evaluate contradictions and implementation divergences noted in `04-contradictions.md`.
4. Surface research gaps, missing cases, and overgeneralizations into `06-gaps.md`.
5. Issue final verdict in `07-verdict.md`.
