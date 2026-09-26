# Research Contradictions: Zero-Downtime Deployment

## Contradiction 1

Statement A: Blue-Green deployment requires maintaining two complete, distinct production environments (blue and green) with traffic routed all at once.

Location: `05-report.md:Finding 3`, `04-contradictions.md:Contradiction 1`

Statement B: Rolling deployment replaces Pods incrementally within a single shared infrastructure environment using `maxSurge` / `maxUnavailable`.

Location: `05-report.md:Finding 1`, `04-contradictions.md:Contradiction 1`

Type: ARCHITECTURAL_CHOICE

Impact: Low. Both achieve zero-downtime under different operational and resource constraints.

Assessment: Analyzed and resolved in research report. Blue-green doubles infrastructure costs for instant rollback capability, whereas rolling deployment conserves resources by updating instances in place.

---

## Contradiction 2

Statement A: Active health probes periodically test endpoints prior to traffic assignment (NGINX Plus `health_check`).

Location: `05-report.md:Finding 7`

Statement B: NGINX Open Source only supports passive health checks (`max_fails` and `fail_timeout`), marking backends down only after failed live client requests.

Location: `05-report.md:Finding 7`, `04-contradictions.md:Contradiction 2`

Type: TOOLING_LIMITATION

Impact: Medium. Relying on passive health checks in NGINX Open Source can expose initial real user requests to a non-ready backend during deployment unless readiness probes are handled at container orchestration layer (Kubernetes/Docker).

Assessment: Properly documented in research report and contradiction analysis.

---

## Contradiction 3

Statement A: PostgreSQL `ALTER TABLE ... ADD COLUMN` with a constant default is fast and metadata-only.

Location: `05-report.md:Finding 10`

Statement B: Evolutionary Database Design recommends nullable columns without default, data backfill via `UPDATE`, followed by adding default / NOT NULL constraint in separate steps.

Location: `04-contradictions.md:Contradiction 4`

Type: METHODOLOGICAL_VARIANCE

Impact: Low. PostgreSQL optimizes constant default column additions, but multi-step evolutionary database refactoring is required for volatile defaults or complex data transformations.

Assessment: Properly resolved in research report.

---

## Summary Assessment
No unresolved or invalidating contradictions found in the research files. All architectural and tooling trade-offs were explicitly identified and analyzed.
