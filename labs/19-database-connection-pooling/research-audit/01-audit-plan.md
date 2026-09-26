# Audit Plan: Database Connection Pooling Research

## Target Lab
`labs/19-database-connection-pooling`

## Pipeline Scope & Override
- Stage: Research Audit Only.
- Implementation and code testing are deferred / NOT APPLICABLE in this stage per pipeline override.
- Output directory: `labs/19-database-connection-pooling/research-audit/`.

## Files Reviewed
- `labs/19-database-connection-pooling/research/01-plan.md`
- `labs/19-database-connection-pooling/research/02-sources.md`
- `labs/19-database-connection-pooling/research/03-evidence.md`
- `labs/19-database-connection-pooling/research/04-contradictions.md`
- `labs/19-database-connection-pooling/research/05-report.md`
- `labs/19-database-connection-pooling/research/06-open-questions.md`

## Claims To Verify
1. Direct connection establishment imposes severe latency and per-process memory penalties.
2. Pool sizes beyond core/hardware saturation degrade throughput ("the knee" / resource contention).
3. The baseline sizing formula `((core_count * 2) + effective_spindle_count)` approaches `core_count * 2` under modern flash/SSD storage.
4. Independent application worker connection pools multiply linearly, causing backend connection exhaustion unless multiplexed by a proxy pooler (e.g. PgBouncer).
5. Connection pool deadlock formula `pool size = Tn x (Cm - 1) + 1` sets the theoretical minimum pool floor for multi-connection threads.
6. Leaks and starvation are observable via `pg_stat_activity` wait states (`idle in transaction`, `ClientRead`).

## Code To Execute
- Code audit and test suite execution are skipped per pipeline override (Research Audit Only).

## Primary Risks
- Citation accuracy and hallucination in sizing formulas and quotes.
- Conflating classical rotational storage heuristics (`effective_spindle_count`) with modern SSD/NVMe deployment reality.
- Treating application-side pooling and proxy-layer pooling as mutually exclusive rather than complementary architectural layers.

## Audit Strategy
1. Live network verification of all cited source URLs.
2. Exact string and semantic matching between research claims and source text.
3. Analysis of architectural and hardware boundary conditions in sizing claims.
4. Consistency check across research findings, contradiction synthesis, and open questions.
