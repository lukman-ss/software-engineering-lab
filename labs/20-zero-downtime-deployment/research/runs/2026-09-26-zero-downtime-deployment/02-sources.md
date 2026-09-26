# 02 — Sources

Research date: 2026-09-26. Accessed: 2026-09-26.

## Source 1
Title: Deployment — Laravel 11.x docs
Publisher: Laravel
URL: https://laravel.com/docs/11.x/deployment
Published: not listed
Accessed: 2026-09-26
Source Tier: 1
Relevance: Official Nginx/PHP-FPM config, optimize caching, `/up` health route + DiagnosingHealth event, Forge/Vapor deployment.

## Source 2
Title: Queues — Laravel 11.x docs
Publisher: Laravel
URL: https://laravel.com/docs/11.x/queues
Published: not listed
Accessed: 2026-09-26
Source Tier: 1
Relevance: queue:work, queue:restart / deployment behavior, Supervisor, after_commit, retry_after, block_for + SIGTERM note.

## Source 3
Title: Laravel Octane — Laravel 11.x docs
Publisher: Laravel
URL: https://laravel.com/docs/11.x/octane
Published: not listed
Accessed: 2026-09-26
Source Tier: 1
Relevance: Octane start/reload/stop behind Nginx, state/memory-leak caveats; relevant because long-lived workers complicate zero-downtime.

## Source 4
Title: Laravel Horizon — Laravel 11.x docs
Publisher: Laravel
URL: https://laravel.com/docs/11.x/horizon
Published: not listed
Accessed: 2026-09-26
Source Tier: 1
Relevance: horizon:terminate graceful deploy step, Supervisor stopwaitsecs guidance, pause/continue/status.

## Source 5
Title: Task Scheduling — Laravel 11.x docs
Publisher: Laravel
URL: https://laravel.com/docs/11.x/scheduling
Published: not listed
Accessed: 2026-09-26
Source Tier: 1
Relevance: schedule:run single-cron, onOneServer, withoutOverlapping, schedule:interrupt for deploys.

## Source 6
Title: Blue Green Deployment — bliki
Publisher: Martin Fowler (martinfowler.com)
URL: https://martinfowler.com/bliki/BlueGreenDeployment.html
Published: 1 March 2010 (page shows 2015-06-05 update note)
Accessed: 2026-09-26
Source Tier: 1
Relevance: Blue-green definition, router switchback rollback, DB-first separation rule.

## Source 7
Title: Parallel Change — bliki
Publisher: Danilo Sato / martinfowler.com
URL: https://martinfowler.com/bliki/ParallelChange.html
Published: 13 May 2014
Accessed: 2026-09-26
Source Tier: 1
Relevance: Expand–migrate–contract pattern; DB refactoring, deployment, API evolution applications.

## Source 8
Title: Pod Lifecycle — Kubernetes docs
Publisher: Kubernetes
URL: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/
Published: not listed
Accessed: 2026-09-26
Source Tier: 1
Relevance: Pod phases/conditions, termination flow (SIGTERM → grace → SIGKILL), endpoint removal ordering.

## Source 9
Title: Liveness, Readiness, and Startup Probes — Kubernetes docs
Publisher: Kubernetes
URL: https://kubernetes.io/docs/concepts/workloads/pods/probes/
Published: not listed
Accessed: 2026-09-26
Source Tier: 1
Relevance: Liveness vs readiness vs startup semantics; readiness gates Service endpoints.

## Source 10
Title: Controlling nginx
Publisher: NGINX
URL: https://nginx.org/en/docs/control.html
Published: not listed
Accessed: 2026-09-26
Source Tier: 1
Relevance: HUP reload (new workers + graceful old-worker drain), QUIT graceful shutdown, USR2/WINCH binary upgrade + rollback path.

## Source 11
Title: Module ngx_http_upstream_module
Publisher: NGINX
URL: https://nginx.org/en/docs/http/ngx_http_upstream_module.html
Published: not listed
Accessed: 2026-09-26
Source Tier: 1
Relevance: upstream server params (max_fails, fail_timeout, backup, down, drain, slow_start); OSS vs commercial health-check limits, sticky/drain notes.

## Source 12
Title: Chapter 5. Data Definition — PostgreSQL 18 docs
Publisher: PostgreSQL Global Development Group
URL: https://www.postgresql.org/docs/current/ddl.html
Published: not listed (current-series doc)
Accessed: 2026-09-26
Source Tier: 1
Relevance: DDL structure/modify-table topic map for backward-compatible migration claims.

## Source 13
Title: ALTER TABLE — PostgreSQL docs
Publisher: PostgreSQL Global Development Group
URL: https://www.postgresql.org/docs/current/sql-altertable.html
Published: not listed (current-series doc)
Accessed: 2026-09-26
Source Tier: 1
Relevance: ADD COLUMN fast path (nullable/no-volatile-default), rewrites for type change/volatile default, ADD CONSTRAINT NOT VALID + VALIDATE CONSTRAINT concurrency pattern, lock levels.

## Source 14 (failed fetch — recorded, not evidence)
Title: Redis upgrade doc (attempted)
Publisher: redis.io
URL: https://redis.io/docs/latest/operate/oss_and_stack/management/upgrading/
Published: unknown
Accessed: 2026-09-26
Source Tier: 1 (unverified — fetch returned 404)
Relevance: NOT USED. Redis rolling-upgrade claim: NOT VERIFIED.
