# Contradiction Audit

## Contradiction 1: NGINX Active Health Checks Availability
Statement A: Lab requirements suggest keeping new instances out of load balancer upstream until readiness checks pass via Nginx HTTP health checks.
Location: `research/03-evidence.md`, `research/05-report.md`
Statement B: Active HTTP health checks, `mandatory`, and `slow_start` directives are proprietary NGINX Plus features and unavailable in NGINX Open Source.
Location: `research/04-contradictions.md` (Contradiction 2), `research/05-report.md` (Finding 5 nuance)
Type: SOURCE_CONFLICT
Impact: HIGH for implementation design if using pure Nginx OSS without an orchestrator or external sidecar/script.
Assessment: Explicitly acknowledged and resolved by research agent. Research correctly concludes that OSS requires external readiness coordination (e.g., Kubernetes or script-based upstream reload).

## Contradiction 2: PostgreSQL DROP COLUMN Behavior
Statement A: Dropping a column immediately breaks running v1 instances executing `SELECT column`.
Location: `research/05-report.md` (Finding 4)
Statement B: PostgreSQL `DROP COLUMN` does not rewrite the table or immediately delete physical bytes from disk.
Location: `research/04-contradictions.md` (Contradiction 4)
Type: INTERNAL
Impact: LOW
Assessment: No practical conflict. The breakdown in v1 is due to metadata invisibility and query parse failures in active sessions, not disk rewrite delays.
