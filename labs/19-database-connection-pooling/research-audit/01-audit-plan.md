# Audit Plan

## Target Lab
`labs/19-database-connection-pooling`

## Files Reviewed
- `research/01-plan.md`
- `research/02-sources.md`
- `research/03-evidence.md`
- `research/04-contradictions.md`
- `research/05-report.md`
- `research/06-open-questions.md`

## Claims To Verify
1. PostgreSQL `max_connections` behavior and defaults.
2. Pool size formula: `((core_count * 2) + effective_spindle_count)`.
3. HikariCP Oracle 50x improvement claim.
4. PgBouncer pool modes and features.
5. Cloud provider (Azure, AWS, GCP) connection limits and recommendations.
6. HikariCP leak detection thresholds.
7. HikariCP deadlock prevention formula.

## Code To Execute
*N/A — Pipeline override restricts to research audit only.*

## Primary Risks
- Oracle 50x claim relies on secondary source (HikariCP citing a video).
- AWS RDS connection limit formula cited via "subagent search" without direct verification.
- SSD pool sizing formula unverified.
- "Knee" curve based on 2014 Wiki data.

## Audit Strategy
1. Cross-reference major claims against provided URLs (where valid/accessible).
2. Evaluate if sources directly support the claims.
3. Identify unsupported claims, overgeneralizations, or gaps.
4. Record contradictions and produce final verdict.
