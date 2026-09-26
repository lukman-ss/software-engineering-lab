# Research Report: Zero-Downtime Deployment

Research Date: 2026-09-26
Topic: Zero-Downtime Deployment — Deploy Versi Baru Tanpa Membuat User Tahu Ada Deployment
Stack: Nginx, Laravel, PostgreSQL, Redis Queue

## Research Question

How can a Laravel application backed by PostgreSQL and Redis Queue be deployed from v1.4 to v1.5 (requiring a new database column and new queue worker version) with zero HTTP downtime, no lost jobs, v1.4/v1.5 coexistence, backward-compatible database state, and fast rollback?

## Executive Summary

Zero-downtime deployment hinges on a single invariant: **start new before stopping old**. All verified patterns — rolling updates, blue-green switching, Nginx graceful reload, and Laravel queue restarts — implement this invariant differently. The hardest constraint is the shared database: migration must support both versions simultaneously via an expand → deploy → migrate → contract sequence. Readiness checks must verify dependencies (DB, Redis, migrations), not just process liveness. In-flight requests and jobs are protected by connection draining and graceful worker shutdown under a process manager (Supervisor). Evidence is HIGH for all core patterns except Docker Compose semantics (NOT VERIFIED) and Redis-as-queue durability edge cases (MEDIUM).

## Findings

### Finding 1 — Core Principle: Start New Before Stopping Old

Claim: Zero-downtime deployment replaces the naive stop → deploy → start sequence with v1 staying live while v2 starts and passes health checks; traffic switches only on PASS, otherwise v2 is discarded and v1 continues.

Evidence: Topic spec sequence diagram and Fowler: "you switch the router so that all incoming requests go to the green... Blue-green deployment also gives you a rapid way to rollback - if anything goes wrong you switch the router back."

Sources:
- https://martinfowler.com/bliki/BlueGreenDeployment.html (2010-03-01)
- https://nginx.org/en/docs/control.html (HUP starts new workers while old workers drain)

Confidence: HIGH

---

### Finding 2 — Rolling Deployment

Claim: With N instances behind a load balancer, instances are updated one-by-one (or in bounded batches). Each new instance must pass a health check before the next is updated, so capacity remains available throughout.

Evidence: Topic spec rolling example A→v2, B→v1, C→v1 → ... → all v2; Kubernetes Deployments implement this as RollingUpdate with maxSurge/maxUnavailable.

Sources:
- https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#strategy
- https://kubernetes.io/docs/tasks/run-application/update-deployment-rolling/

Confidence: HIGH

---

### Finding 3 — Blue-Green Deployment: Fast Rollback at Double Capacity Cost

Claim: Blue (live, v1) and Green (standby, v2) are two near-identical environments. Deploy, migrate, and smoke-test in Green; switch the router Blue→Green only when healthy. Rollback is Green→Blue (router switch). The trade-off is ~2× infrastructure during the switch window.

Evidence: Fowler: "At any time one of them, let's say blue ... is live... Once the software is working in the green environment, you switch the router... There's still the issue of dealing with missed transactions while the green environment was live."

Sources:
- https://martinfowler.com/bliki/BlueGreenDeployment.html (2010-03-01)

Confidence: HIGH

---

### Finding 4 — Database Migrations Must Be Backward-Compatible (Expand → Deploy → Migrate → Contract)

Claim: Because v1 and v2 coexist, migrations cannot break v1. Renaming `name` → `full_name` via `DROP COLUMN name` breaks v1's `SELECT name`. The safe pattern keeps both columns, deploys code that understands both, migrates data, then drops the old column only after all instances are on v2.

Evidence: Fowler: "The trick is to separate the deployment of schema changes from application upgrades. So first apply a database refactoring to change the schema to support both the new and old version... deploy that, check everything is working... then deploy the new version... (And when the upgrade has bedded down remove the database support for the old version.)" PostgreSQL: `ADD COLUMN` with non-volatile DEFAULT is metadata-only and fast; `DROP COLUMN` makes the column invisible to all sessions immediately (no rewrite but instant breakage); most ALTER TABLE forms take ACCESS EXCLUSIVE lock.

Sources:
- https://martinfowler.com/bliki/BlueGreenDeployment.html
- https://www.postgresql.org/docs/current/sql-altertable.html (2026-09-24)

Confidence: HIGH

---

### Finding 5 — Readiness vs Liveness (Health Checks Must Verify Dependencies)

Claim: `GET /health` returning `{"status":"ok"}` is insufficient. Liveness = is the process alive; Readiness = is the instance ready to receive traffic (DB, Redis, storage, migrations). New instances must not enter the load balancer until readiness passes.

Evidence: Kubernetes defines separate liveness/readiness/startup probes; pods failing readiness are removed from service endpoints. Laravel health route at `/up` returns 200 only if the app booted without exception and dispatches `DiagnosingHealth` — a listener can check DB/cache and throw to fail the check. NGINX Plus active `health_check` marks a server unhealthy if response is outside 200–399 and `mandatory` forces new servers to pass before receiving traffic.

Sources:
- https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/
- https://laravel.com/docs/12.x/deployment (health route)
- https://docs.nginx.com/nginx/admin-guide/load-balancer/http-health-check/

Confidence: HIGH

Nuance: Active health checks, `mandatory`, and `slow_start` are NGINX Plus-only. NGINX Open Source has only passive checks (`max_fails`/`fail_timeout`); lab using OSS must implement readiness externally (orchestrator or script).

---

### Finding 6 — Graceful Shutdown / Connection Draining

Claim: Killing a worker mid-request (`docker stop`, SIGTERM) drops in-flight `POST /checkout` and can cause duplicate/missing payments. Correct shutdown: stop sending new requests → wait for active requests to finish → terminate.

Evidence: NGINX: `HUP` starts new workers with new config; old workers "close listen sockets and continue to service old clients. After all clients are serviced, old worker processes are shut down." `QUIT` = graceful shutdown, `TERM`/`INT` = fast shutdown. `USR2` allows zero-downtime binary upgrade with rollback.

Sources:
- https://nginx.org/en/docs/control.html
- https://docs.nginx.com/nginx/admin-guide/load-balancer/http-health-check/ (slow_start/connection draining context)

Confidence: HIGH

---

### Finding 7 — Queue Worker Deployment Without Losing Jobs

Claim: Queue workers are long-lived processes that do not pick up new code until restarted. Jobs are stored durably in Redis/DB outside the worker, so queued jobs survive a worker restart; only the currently-executing job is at risk. Laravel provides `php artisan queue:restart` (graceful: finish current job then exit) and `php artisan reload`; a process manager (Supervisor) must auto-restart workers afterward. `retry_after` must exceed worker `--timeout` to avoid double-processing of frozen jobs.

Evidence: Laravel Deployment docs: "After deploying ... any long-running services such as queue workers ... should be reloaded / restarted... Laravel provides a single reload Artisan command." Queues docs: "This command will instruct all queue workers to gracefully exit after they finish processing their current job so that no existing jobs are lost. Since the queue workers will exit ... you should be running a process manager such as Supervisor to automatically restart the queue workers. The queue uses the cache to store restart signals." And: "The --timeout value should always be at least several seconds shorter than your retry_after."

Sources:
- https://laravel.com/docs/12.x/deployment
- https://laravel.com/docs/12.x/queues (Queue Workers and Deployment; Job Expirations and Timeouts)

Confidence: HIGH for graceful-restart + Supervisor pattern; MEDIUM that Redis queue durability fully prevents job loss (relies on correct Redis persistence config; in-flight job retry semantics depend on driver).

---

## Areas of Agreement

- All authoritative sources agree v1 and v2 must coexist during deployment.
- Expand-contract is the only safe schema change pattern for shared-database zero-downtime.
- Liveness ≠ readiness; all sources distinguish them and require readiness gates before traffic.
- Graceful shutdown requires draining before termination.

## Areas of Disagreement

- No direct contradictions among primary sources on core patterns. Differences are scope/trade-off: blue-green vs rolling (cost vs speed), NGINX Plus vs OSS feature availability. See 04-contradictions.md for detail.
- `DROP COLUMN` disagreement is only apparent: PostgreSQL is fast (no rewrite) but immediately breaks old code (invisibility), which the topic spec correctly flags as "Boom."

## Limitations

- Docker Compose zero-downtime evidence NOT VERIFIED (intended reference URL returned 404).
- Redis Open Source active health checks unavailable in NGINX OSS — lab must document alternative.
- Queue in-flight job behavior under SIGKILL vs graceful stop not fully quoted from PHPRedis/Redis docs in this session.
- Kubernetes rolling-update evidence partially truncated; strategy details summarized from nav/anchors.
- Time-sensitive: Laravel v12.x, PostgreSQL v18 (2026-09-24), Kubernetes v1.37 — behavior may evolve.

## Conclusion

For the lab stack (Nginx, Laravel, PostgreSQL, Redis Queue) the verified zero-downtime sequence is: (1) expand the schema (additive, non-breaking migration), (2) start v1.5 instances with readiness gates that check DB/Redis/migrations, (3) switch traffic incrementally (rolling) or atomically (blue-green) with health checks, (4) drain and terminate v1.4 instances and workers gracefully via Supervisor-managed restarts, (5) backfill/migrate data, (6) contract the schema after all v1.4 instances are gone. Rollback is a traffic switch (blue-green) or re-rolling to the previous image, provided migrations remain backward-compatible.
