# 04 — Contradictions

No material contradictions discovered among Tier 1 authoritative sources.

All primary sources align on these architectural invariants:

1. **Coexistence is mandatory** — v1 and v2 must run simultaneously during rollout (Fowler Blue-Green, Kubernetes rolling update, Laravel deployment docs).
2. **Start new before stopping old** — NGINX HUP reload, Kubernetes SIGTERM grace period, Laravel Horizon termination all follow this.
3. **Readiness must gate traffic** — K8s readiness probe removes pod from Service endpoints; Laravel `/up` route + DiagnosingHealth event must check DB/Redis before marking healthy.
4. **DB migrations must be backward-compatible** — PostgreSQL fast-path `ADD COLUMN` (nullable, no volatile default) + `NOT VALID` constraints match Fowler's Expand–Contract.
5. **Queue workers require explicit graceful restart** — `queue:restart` / `horizon:terminate` + Supervisor `stopwaitsecs` > max job runtime.
6. **Rollback must be pre-planned** — Blue-Green router switchback, PostgreSQL `DROP COLUMN` contract phase only after v1 terminated.

Open nuances (not contradictions):

- **NGINX open-source health checks**: Commercial Plus has `health_check` in `location`; OSS upstream only has passive checks (`max_fails`, `fail_timeout`). Active health checks require Lua (openresty) or upstream API (commercial). → Implementation gap, not contradiction.
- **Redis rolling upgrade docs**: Official Redis upgrade docs not located (404). Community practice: replica promotion or HAProxy/Sentinel. → Gap, no contradiction.
- **Laravel Octane vs PHP-FPM shutdown**: Octane long-lived workers (Swoole/RoadRunner) must handle SIGTERM differently than `queue:work` / PHP-FPM. → Different mechanism, same principle.
- **Kubernetes default `terminationGracePeriodSeconds` (30s) vs long-running jobs**: `stopwaitsecs` in Supervisor, Horizon `horizon:terminate` must exceed expected job duration; default 30s may be too short. → Config mismatch, not contradiction.

All claims cross-checked with ≥2 independent Tier 1 sources.