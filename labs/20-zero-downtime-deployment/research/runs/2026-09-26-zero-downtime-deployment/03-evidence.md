# 03 — Evidence

## Evidence 1
Claim: NGINX supports graceful configuration reload using `HUP` signal, enabling zero-downtime worker replacement.
Evidence: When `nginx -s reload` (HUP signal) is sent, the master opens new configuration; if successful, new worker processes start and old workers receive graceful shutdown message, close listen sockets, continue serving existing connections until finished, then exit.
Source: NGINX Control.html
URL: https://nginx.org/en/docs/control.html
Confidence: HIGH
Corroborated By: Kubernetes pod lifecycle (termination), Laravel Octane `octane:reload`
Notes: No user-visible downtime; only in-flight requests must complete before old workers exit.

## Evidence 2
Claim: NGINX supports graceful shutdown with `QUIT` signal and in-place binary upgrade with `USR2` signal.
Evidence: `QUIT` triggers graceful shutdown (workers finish, exit). `USR2` upgrades binary: master renames current to `.oldbin`, starts new master/bin; workers of both run, `WINCH` on old master gracefully shuts down old workers; can fallback via HUP on old master (restart workers w/o full reload) or TERM on new master.
Source: NGINX control.html
URL: https://nginx.org/en/docs/control.html
Confidence: HIGH
Corroborated By: nginx.org docs
Notes: Binary upgrade rollback documented; enables seamless version change without losing existing requests.

## Evidence 3
Claim: Kubernetes pod termination sends SIGTERM, waits terminationGracePeriodSeconds (default 30s), then SIGKILL; Readiness probe determines pod inclusion in Service endpoints.
Evidence: Pod Lifecycle page describes: container signals SIGTERM for graceful start, pod level has terminationGracePeriodSeconds; kubelet waits before SIGKILL; Pod Disruption Budgets control max unavailable. Readiness probe controls when pod is removed from Service/LB.
Source: Kubernetes Pod Lifecycle
URL: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/
Confidence: HIGH
Corroborated By: Laravel Horizon `horizon:terminate`
Notes: terminationGracePeriodSeconds defaults to 30s; configurable. Connection draining handled at pod level.

## Evidence 4
Claim: Kubernetes distinguishes Liveness, Readiness, and Startup probes; Readiness determines when pod receives traffic.
Evidence: Probe page defines: Liveness=restart policy; Readiness=governs Service endpoint membership; Startup=allow extended boot time; Application readiness must be true for traffic routing.
Source: Kubernetes Probes
URL: https://kubernetes.io/docs/concepts/workloads/pods/probes/
Confidence: HIGH
Corroborated By: Laravel `/up` route (Liveness equivalent), app-level database health checks
Notes: Default `/up` route insufficient; must check DB/Redis/Migration status.

## Evidence 5
Claim: Laravel `queue:work` can receive SIGTERM for graceful termination; `queue:restart` command triggers graceful reload.
Evidence: Queues docs note signal handling; workers must restart to pick up code changes; `queue:restart` sends message via cache to workers, which finish current job before exiting. Octane docs confirm graceful behavior.
Source: Laravel Queues, Octane
URL: https://laravel.com/docs/11.x/queues, https://laravel.com/docs/11.x/octane
Confidence: HIGH
Corroborated By: Horizon Terminate docs
Notes: `block_for=0` causes indefinite blocking; worker must be able to exit on signal.

## Evidence 6
Claim: Laravel Horizon’s `horizon:terminate` gracefully terminates workers; Supervisor `stopwaitsecs` must exceed longest job duration.
Evidence: Horizon Deploying docs: after deployment, run `php artisan horizon:terminate` so Supervisor restarts process with new code. Supervisor config must set `stopwaitsecs` > max job time or jobs killed.
Source: Horizon Deploying, Supervisor docs
URL: https://laravel.com/docs/11.x/horizon#deploying-horizon
Confidence: HIGH
Corroborated By: Laravel queues SIGTERM handling
Notes: `stopwaitsecs=3600` common for long jobs.

## Evidence 7
Claim: Blue-Green deployment requires two identical environments; router switch enables fast rollback.
Evidence: Martin Fowler describes: Blue (live) vs Green (staging), switch router when Green passes smoke test; rollback = router back; allows disaster-recovery testing on every release.
Source: Martin Fowler — Blue Green Deployment
URL: https://martinfowler.com/bliki/BlueGreenDeployment.html
Confidence: HIGH
Corroborated By: Zero-downtime-deployment theory
Notes: Infrastructure cost doubles; DB requires separate migration strategy.

## Evidence 8
Claim: Parallel Change (Expand–Migrate–Contract) enables backward-compatible DB schema evolution during deployment.
Evidence: Fowler/Sato: 1) Expand: add new column (nullable), both old+new available. 2) Migrate: deploy app handling both schemas. 3) Contract: drop old column after all instances upgraded. Applicable to APIs, DB, deployments.
Source: Martin Fowler — Parallel Change
URL: https://martinfowler.com/bliki/ParallelChange.html
Confidence: HIGH
Corroborated By: PostgreSQL ALTER TABLE docs
Notes: Deploy phase may be long (external clients, migrations).

## Evidence 9
Claim: PostgreSQL `ADD COLUMN` without `NOT NULL` or with volatile default requires full rewrite; without those, it is metadata-only and fast (allows rolling migration).
Evidence: ALTER TABLE docs: adding column with `IF NOT EXISTS` or `DEFAULT` (non-volatile) is instant; changing type/triggering volatile default or adding constraint requires table rewrite + lock.
Source: PostgreSQL ALTER TABLE
URL: https://www.postgresql.org/docs/current/sql-altertable.html
Confidence: HIGH
Corroborated By: Zero-downtime deployment patterns
Notes: ADD COLUMN (nullable, non-volatile default) is safe for zero-downtime; DROP COLUMN fast but not disk reclamation immediate.

## Evidence 10
Claim: PostgreSQL `ADD CONSTRAINT NOT VALID` allows non-blocking constraint addition; `VALIDATE CONSTRAINT` validates pre-existing rows without blocking new writes.
Evidence: ALTER TABLE docs: constraints can be added `NOT ENFORCED`; validation requires scan. `NOT VALID` skips scan, adds constraint immediately; `VALIDATE CONSTRAINT` acquires only `SHARE UPDATE EXCLUSIVE` lock.
Source: PostgreSQL ALTER TABLE
URL: https://www.postgresql.org/docs/current/sql-altertable.html
Confidence: HIGH
Corroborated By: Expand–Contract pattern
Notes: Use `ADD ... NOT VALID` then `VALIDATE CONSTRAINT` in separate steps for zero-downtime.

## Evidence 11
Claim: Laravel Health Rote `/up` returns 200 on successful boot; `DiagnosingHealth` event enables dependency checks.
Evidence: Deployment docs: `/up` route exists, returns 500 on exception during boot; `Illuminate\Foundation\Events\DiagnosingHealth` dispatched; listener can check DB/cache, throw exception on failure.
Source: Laravel Deployment — The Health Route
URL: https://laravel.com/docs/11.x/deployment#the-health-route
Confidence: HIGH
Corroborated By: Kubernetes probes
Notes: Default route is too permissive; must extend for DB/Redis/Migration checks.

## Evidence 12
Claim: NGINX `upstream` supports `backup` servers for failover and `drain` mode (commercial NGINX Plus/open source 1.13.6+ experimental) for connection draining.
Evidence: Upstream module docs: `server ... backup` serves only when primaries down; `drain` marks server as draining, only bound requests proxied (sessionAffinity context). OSS health-check is open-source but `drain` parameter is commercial-only as of 1.13.6.
Source: NGINX upstream_module
URL: https://nginx.org/en/docs/http/ngx_http_upstream_module.html
Confidence: HIGH
Corroborated By: Zero-downtime deployment flow
Notes: Pure OSS requires manual removal from upstream pool via API or configuration reload.

## Evidence 13
Claim: The `queue:restart` command and Horizon’s `horizon:terminate` enable zero-downtime queue worker deployment.
Evidence: Queues docs: `php artisan queue:restart` gracefully restarts workers; after commit, any worker starting a new job will exit after completing current job. Horizon adds `horizon:terminate` + `horizon:pause` / `continue` for supervisor-level control.
Source: Laravel Queues — Queue Workers and Deployment
URL: https://laravel.com/docs/11.x/queues#running-the-queue-worker
Confidence: HIGH
Corroborated By: Supervisor configuration, graceful shutdown principles
Notes: Requires `worker_lifetime` or `timeout` to force periodic checks that new code is loaded.

## Evidence 14
NOT VERIFIED
Claim: Redis supports in-place binary upgrade / rolling upgrade without data loss.
Evidence: Attempted fetch returned HTTP 404; no authoritative source found.
Source: redis.io (attempted)
URL: Not available
Confidence: LOW
Corroborated By: None
Notes: Community sources confirm Redis accepts RDB-AOF; upgrade should be safe but requires replication/cluster coordination.

## Evidence 15
Claim: Laravel Scheduler supports `onOneServer` to ensure single-instance execution across workers.
Evidence: Scheduling docs: `Schedule::command(...)->onOneServer()` uses cache-based atomic lock; first server wins; prevents duplicate execution on multi-server schedulers.
Source: Laravel Task Scheduling
URL: https://laravel.com/docs/11.x/scheduling#running-tasks-on-one-server
Confidence: HIGH
Corroborated By: K8s leader election patterns
Notes: Requires cache driver supporting locks (redis/database/memcached).
