# Research Audit Plan: Database Connection Pooling

## Target Lab
`labs/19-database-connection-pooling`

## Files Reviewed
- `labs/19-database-connection-pooling/research/01-plan.md`
- `labs/19-database-connection-pooling/research/02-sources.md`
- `labs/19-database-connection-pooling/research/03-evidence.md`
- `labs/19-database-connection-pooling/research/04-contradictions.md`
- `labs/19-database-connection-pooling/research/05-report.md`
- `labs/19-database-connection-pooling/research/06-open-questions.md`

## Claims To Verify
1. Connection establishment overhead and backend process resource scaling in PostgreSQL.
2. Throughput collapse and latency degradation ("the knee") past hardware saturation limits.
3. Universal pool sizing formula `((core_count * 2) + effective_spindle_count)` and SSD behavior.
4. Linear multiplication of client-side pools causing connection exhaustion in distributed topologies.
5. Pool deadlock avoidance formula `pool size = Tn * (Cm - 1) + 1`.
6. Detection of connection leaks via `pg_stat_activity` states (`idle in transaction`) and wait event `ClientRead`.

## Code To Execute
- N/A (Pipeline Override: Audit research only; code audit omitted in this stage).

## Primary Risks
- Overreliance on historical spindle-disk sizing formulas in modern NVMe / cloud deployments.
- Conflation of application-level pool sizing (HikariCP) with proxy-level multiplexing (PgBouncer).
- Unverified mathematical claims regarding deadlock prevention or throughput degradation.
- Misrepresentation of PostgreSQL connection state tracking in operational diagnostics.

## Audit Strategy
1. Independently fetch and verify all cited URLs (PostgreSQL docs, HikariCP wiki, PostgreSQL wiki, PgBouncer docs).
2. Validate claim fidelity against primary source text and engine specifications.
3. Verify consistency between research report, evidence matrix, contradictions, and open questions.
4. Identify research gaps and technical nuances before implementation begins.
