# Source Audit

## Source 1
Claimed Title: Blue Green Deployment
Claimed Publisher: Martin Fowler (ThoughtWorks)
URL: https://martinfowler.com/bliki/BlueGreenDeployment.html

Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: None. Canonical definition and database refactoring guidance verified.
Assessment: PASS

## Source 2
Claimed Title: Pod Lifecycle — Kubernetes Documentation
Claimed Publisher: Kubernetes (CNCF)
URL: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/

Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: None. Probe mechanics and lifecycle status verified.
Assessment: PASS

## Source 3
Claimed Title: HTTP Health Checks — NGINX Documentation
Claimed Publisher: F5 / NGINX
URL: https://docs.nginx.com/nginx/admin-guide/load-balancer/http-health-check/

Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: Features documented (active health checks, mandatory, slow_start) apply exclusively to commercial NGINX Plus, not NGINX Open Source.
Assessment: WARNING

## Source 4
Claimed Title: Deployment — Laravel Documentation (v12.x)
Claimed Publisher: Laravel
URL: https://laravel.com/docs/12.x/deployment

Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: None. Verified built-in `/up` health route and `DiagnosingHealth` event dispatching.
Assessment: PASS

## Source 5
Claimed Title: Queues — Laravel Documentation (v12.x)
Claimed Publisher: Laravel
URL: https://laravel.com/docs/12.x/queues

Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: None. Worker graceful restart mechanics confirmed.
Assessment: PASS

## Source 6
Claimed Title: Controlling NGINX
Claimed Publisher: NGINX
URL: https://nginx.org/en/docs/control.html

Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: None. Verified `HUP`, `QUIT`, `USR2`, and old worker draining mechanics.
Assessment: PASS

## Source 7
Claimed Title: ALTER TABLE — PostgreSQL Documentation (v18)
Claimed Publisher: PostgreSQL
URL: https://www.postgresql.org/docs/current/sql-altertable.html

Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: None. Verified metadata-only addition with non-volatile DEFAULT and ACCESS EXCLUSIVE locking.
Assessment: PASS

## Source 8
Claimed Title: Update a Deployment Without Downtime — Kubernetes
Claimed Publisher: Kubernetes (CNCF)
URL: https://kubernetes.io/docs/tasks/run-application/update-deployment-rolling/

Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: None. Rolling update workflow confirmed.
Assessment: PASS

## Source 9
Claimed Title: Deployments — Kubernetes Strategy
Claimed Publisher: Kubernetes (CNCF)
URL: https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#strategy

Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: None. RollingUpdate vs Recreate strategy semantics confirmed.
Assessment: PASS

## Source 10 (Referenced in Open Questions / Missing Sources)
Claimed Title: Docker Compose Rolling Update Documentation
Claimed Publisher: Docker
URL: https://docs.docker.com/compose/how-tos/rolling-update/

Reachable: NO (HTTP 404)
Source Type: UNKNOWN
Relevant: PARTIAL
Supports Claimed Topic: NO
Problems: Dead URL. Properly flagged by research agent as 404 in `06-open-questions.md`.
Assessment: FAIL
