# Research Audit Plan

Target Lab: labs/20-zero-downtime-deployment
Files Reviewed:
- `research/01-plan.md`
- `research/02-sources.md`
- `research/03-evidence.md`
- `research/04-contradictions.md`
- `research/05-report.md`
- `research/06-open-questions.md`

Claims To Verify:
1. Blue-green deployment enables instant rollback by switching load balancer target.
2. Expanding database schema via backward-compatible migrations allows old (v1) and new (v2) app versions to run concurrently.
3. PostgreSQL `ADD COLUMN` with non-volatile DEFAULT does not trigger full table rewrite.
4. Liveness probes are insufficient for zero-downtime readiness checks (must verify DB/Redis/migrations).
5. Nginx graceful reload (`HUP` signal) finishes active client requests before terminating old worker processes.
6. Laravel `queue:restart` allows active background jobs to finish before workers exit.

Code To Execute:
- Pipeline override active: Research audit only. Code execution skipped per instructions.

Primary Risks:
- Invalid or broken source URLs (e.g., Docker Compose 404 URL).
- Misrepresenting Nginx Open Source vs Nginx Plus feature availability (active health checks).
- Over-promising Redis queue durability during process crashes.

Audit Strategy:
- Verify all cited URLs and sources.
- Check cross-consistency between claims, evidence, and report.
- Assess clarity of limitations and research gaps.
