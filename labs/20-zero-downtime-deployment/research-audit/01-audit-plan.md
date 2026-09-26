# Research Audit Plan: Zero-Downtime Deployment

## Target Lab
`labs/20-zero-downtime-deployment`

## Files Reviewed
- `research/runs/2026-09-26-zero-downtime-deployment/01-plan.md`
- `research/runs/2026-09-26-zero-downtime-deployment/02-sources.md`
- `research/runs/2026-09-26-zero-downtime-deployment/03-evidence.md`
- `research/runs/2026-09-26-zero-downtime-deployment/04-contradictions.md`
- `research/runs/2026-09-26-zero-downtime-deployment/05-report.md`
- `research/runs/2026-09-26-zero-downtime-deployment/06-open-questions.md`

## Claims To Verify
1. Rolling deployment mechanics and default `maxUnavailable` / `maxSurge` parameters in Kubernetes.
2. Liveness vs readiness probe separation of duties and traffic routing impact.
3. Blue-green deployment definition, rollback mechanism, and database separation requirements.
4. Database backward compatibility strategies (expand-deploy-migrate-contract / transition views).
5. Health check requirements beyond basic container running state.
6. Pod termination sequence (SIGTERM, grace period, SIGKILL) and endpoint slice condition timing.
7. NGINX passive vs active health check capabilities (NGINX Plus requirement for active checks).
8. Laravel built-in `/up` health endpoint and `DiagnosingHealth` event dispatching.
9. Laravel Horizon graceful termination via `horizon:terminate` and Supervisor `stopwaitsecs` interaction.
10. PostgreSQL `ALTER TABLE ADD COLUMN` constant default metadata-only execution.
11. Docker container stop signal flow and default grace periods.

## Code To Execute
*Pipeline override: Research-only audit stage. Implementation execution not applicable in this phase.*

## Primary Risks
- Overgeneralization of database non-blocking DDL (e.g. constant default vs volatile default or locks).
- Conflating application-level health with container-level liveness.
- Unstated dependency on commercial NGINX Plus features for active health probing.
- Misunderstanding queue worker shutdown during blocking poll (`block_for=0`).

## Audit Strategy
1. Validate external source URLs, publishers, titles, and authority tiers.
2. Cross-check each claimed finding in `05-report.md` and `03-evidence.md` against authoritative documentation.
3. Analyze contradictions and open question completeness.
4. Synthesize audit findings and render an evidence-based verdict.
