# Research Sources: Zero-Downtime Deployment

## Source 1

**Title:** Deployments | Kubernetes  
**Publisher:** The Kubernetes Authors  
**URL:** https://kubernetes.io/docs/concepts/workloads/controllers/deployment/  
**Published:** 2026  
**Accessed:** 2026-09-26  
**Source Tier:** Tier 1 (Official Documentation)  
**Relevance:** Primary source for understanding rolling deployments, update strategies, maxSurge and maxUnavailable parameters, and deployment lifecycle in Kubernetes. Critical for understanding how Kubernetes handles zero-downtime deployments automatically.

## Source 2

**Title:** Pod Lifecycle | Kubernetes  
**Publisher:** The Kubernetes Authors  
**URL:** https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/  
**Published:** 2026  
**Accessed:** 2026-09-26  
**Source Tier:** Tier 1 (Official Documentation)  
**Relevance:** Comprehensive documentation on pod phases, container states, and termination behavior. Contains essential information about how Kubernetes handles pod shutdown and endpoint removal.

## Source 3

**Title:** Liveness, Readiness, and Startup Probes | Kubernetes  
**Publisher:** The Kubernetes Authors  
**URL:** https://kubernetes.io/docs/concepts/workloads/pods/probes/  
**Published:** 2026  
**Accessed:** 2026-09-26  
**Source Tier:** Tier 1 (Official Documentation)  
**Relevance:** Definitive source for understanding the three types of probes. Documents the distinction between liveness (restart unhealthy containers) and readiness (control traffic routing). Includes probe-level terminationGracePeriodSeconds feature (v1.28+).

## Source 4

**Title:** Blue Green Deployment  
**Publisher:** Martin Fowler  
**URL:** https://martinfowler.com/bliki/BlueGreenDeployment.html  
**Published:** 2010  
**Accessed:** 2026-09-26  
**Source Tier:** Tier 1 (Original Concept Paper)  
**Relevance:** Original definition and explanation of blue-green deployment pattern. Contains critical insight about database schema changes needing to be deployed separately from application code. Introduced the concept of having both environments ready for quick rollback.

## Source 5

**Title:** HTTP Health Checks | NGINX Documentation  
**Publisher:** F5 NGINX  
**URL:** https://docs.nginx.com/nginx/admin-guide/load-balancer/http-health-check/  
**Published:** 2026  
**Accessed:** 2026-09-26  
**Source Tier:** Tier 1 (Official Documentation)  
**Relevance:** Documents passive and active health check mechanisms in NGINX. Explains fail_timeout, max_fails, slow_start parameters. Includes information about health check intervals, mandatory health checks, and connection reuse with keepalive.

## Source 6

**Title:** Deployment | Laravel 11.x  
**Publisher:** Laravel  
**URL:** https://laravel.com/docs/11.x/deployment  
**Published:** 2026  
**Accessed:** 2026-09-26  
**Source Tier:** Tier 1 (Official Documentation)  
**Relevance:** Laravel's official deployment guide. Documents health endpoint (/up), optimization commands (optimize:clear, config:cache, route:cache, view:cache), and server requirements. Explains the built-in health route for load balancers and orchestration systems.

## Source 7

**Title:** Queues | Laravel 11.x  
**Publisher:** Laravel  
**URL:** https://laravel.com/docs/11.x/queues  
**Published:** 2026  
**Accessed:** 2026-09-26  
**Source Tier:** Tier 1 (Official Documentation)  
**Relevance:** Comprehensive documentation on Laravel's queue system. Documents retry_after configuration, after_commit option for queue dispatching, job timeouts, and max job attempts. Includes information on Redis queue driver configuration.

## Source 8

**Title:** Laravel Horizon  
**Publisher:** Laravel  
**URL:** https://laravel.com/docs/11.x/horizon  
**Published:** 2026  
**Accessed:** 2026-09-26  
**Source Tier:** Tier 1 (Official Documentation)  
**Relevance:** Documents the recommended approach for deploying queue workers with Horizon. Explains graceful termination using `php artisan horizon:terminate`, the stopwaitsecs supervisor configuration, and balancing strategies (simple, auto, false).

## Source 9

**Title:** Evolutionary Database Design  
**Publisher:** ThoughtWorks  
**URL:** https://martinfowler.com/articles/evodb.html  
**Published:** 2016  
**Accessed:** 2026-09-26  
**Source Tier:** Tier 1 (Original Best Practices Paper)  
**Relevance:** Authoritative source for database migration patterns. Introduces the expand-deploy-migrate-contract pattern and transition phases. Documents how to handle destructive changes that require both old and new versions to coexist.

## Source 10

**Title:** 5.7. Modifying Tables | PostgreSQL Documentation  
**Publisher:** PostgreSQL Global Development Group  
**URL:** https://www.postgresql.org/docs/current/ddl-alter.html  
**Published:** 2026  
**Accessed:** 2026-09-26  
**Source Tier:** Tier 1 (Official Documentation)  
**Relevance:** Documents PostgreSQL's ALTER TABLE capabilities including ADD COLUMN, DROP COLUMN, SET DEFAULT, and data type changes. Critical for understanding which schema changes are non-blocking and which require table rewrites.

## Source 11

**Title:** Update a Deployment Without Downtime | Kubernetes  
**Publisher:** The Kubernetes Authors  
**URL:** https://kubernetes.io/docs/tasks/run-application/update-deployment-rolling/  
**Published:** 2026  
**Accessed:** 2026-09-26  
**Source Tier:** Tier 1 (Official Documentation)  
**Relevance:** Task-based tutorial demonstrating rolling updates in Kubernetes. Shows kubectl rollout status, pause/resume, and rollback commands. Includes configuration examples for maxUnavailable and maxSurge.

## Source 12

**Title:** docker container stop | Docker Documentation  
**Publisher:** Docker  
**URL:** https://docs.docker.com/reference/cli/docker/container/stop/  
**Published:** 2026  
**Accessed:** 2026-09-26  
**Source Tier:** Tier 1 (Official Documentation)  
**Relevance:** Documents the container stop behavior: SIGTERM → grace period → SIGKILL. Explains --timeout flag, --signal flag, STOPSIGNAL Dockerfile instruction, and default timeout values (10s Linux, 30s Windows).

## Source 13

**Title:** Server Configuration | PostgreSQL Documentation  
**Publisher:** PostgreSQL Global Development Group  
**URL:** https://www.postgresql.org/docs/current/runtime-config.html  
**Published:** 2026  
**Accessed:** 2026-09-26  
**Source Tier:** Tier 1 (Official Documentation)  
**Relevance:** Reference for PostgreSQL server configuration parameters. Relevant for understanding connection limits, statement timeouts, and other runtime settings that affect deployment behavior.

## Source 14

**Title:** Pod Termination | Kubernetes  
**Publisher:** The Kubernetes Authors  
**URL:** https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#pod-termination  
**Published:** 2026  
**Accessed:** 2026-09-26  
**Source Tier:** Tier 1 (Official Documentation)  
**Relevance:** Explains pod termination flow, stop signals, forced termination, and how endpoints are updated when pods are deleted. Critical for understanding connection draining behavior.
