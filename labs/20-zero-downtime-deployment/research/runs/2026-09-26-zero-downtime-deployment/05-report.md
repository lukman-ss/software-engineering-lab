# Research Report: Zero-Downtime Deployment

## Research Question

What strategies, patterns, and best practices enable zero-downtime deployments for Laravel applications with PostgreSQL and Redis Queue backends, ensuring continuous service availability during upgrades while maintaining data integrity?

## Executive Summary

Zero-downtime deployment requires a coordinated approach across multiple layers: infrastructure orchestration (Kubernetes/Docker), application health monitoring (liveness/readiness probes), database migration patterns (expand-deploy-migrate-contract), and queue worker lifecycle management. The core principle is maintaining backward compatibility between old and new versions during the transition window, with traffic gradually shifting only after new instances pass health verification. Key findings include: (1) Rolling deployments provide incremental updates with configurable availability tradeoffs via maxUnavailable/maxSurge; (2) Liveness probes control container restarts while readiness probes gate traffic; (3) Database schema changes must be backward-compatible, favoring additive changes and transition phases; (4) Graceful shutdown requires coordinated signal handling (SIGTERM → grace period → SIGKILL) at container, application, and queue worker levels; (5) Queue worker deployment needs careful consideration of signal handling during blocking waits; (6) Rollback mechanisms require pre-planned infrastructure capacity and immediate switching capability.

## Findings

### Finding 1: Rolling Deployment Strategy

**Claim:** Rolling deployment gradually replaces old Pods with new ones, keeping the application available throughout the process.

**Evidence:** Kubernetes Deployments support two update strategy types: RollingUpdate (default) and Recreate (which causes downtime). With RollingUpdate, the parameters `maxUnavailable` (default 25%) and `maxSurge` (default 25%) control how many pods can be unavailable/created simultaneously during an update. A rolling update gradually replaces old Pods with new ones, so your application remains available throughout the process.

**Sources:** Kubernetes Documentation - Update a Deployment Without Downtime (Evidence 1), Kubernetes Deployments documentation (Source 1)

**Confidence:** HIGH

### Finding 2: Health Check Separation of Concerns

**Claim:** Liveness probes determine whether to restart a container, while readiness probes determine whether a container should receive traffic.

**Evidence:** "Liveness probes determine when to restart a container." "Readiness probes determine when a container is ready to accept traffic." "If a container fails its readiness probe, the EndpointSlice controller removes the Pod's IP address from the EndpointSlices of all Services that match the Pod. Readiness probes run on the container during its whole lifecycle."

**Sources:** Kubernetes Documentation - Liveness, Readiness, and Startup Probes (Evidence 2, Source 3)

**Confidence:** HIGH

### Finding 3: Blue-Green Deployment Pattern

**Claim:** Blue-green deployment uses two identical production environments, where traffic is switched from old (blue) to new (green) all at once after full deployment and testing.

**Evidence:** "The blue-green deployment approach does this by ensuring you have two production environments, as identical as possible. At any time one of them, let's say blue for the example, is live. As you prepare a new release of your software you do your final stage of testing in the green environment. Once the software is working in the green environment, you switch the router so that all incoming requests go to the green environment - the blue one is now idle."

**Sources:** Martin Fowler - Blue Green Deployment (Evidence 3, Source 4)

**Confidence:** HIGH

### Finding 4: Database Migration Backward Compatibility

**Claim:** Database schema changes that rename or remove columns cause backward compatibility failures when both old and new application versions are running simultaneously.

**Evidence:** "Databases can often be a challenge with this technique [blue-green], particularly when you need to change the schema to support a new version of the software. The trick is to separate the deployment of schema changes from application upgrades." Demonstrates a transition phase pattern: "ALTER TABLE customer RENAME to client; CREATE VIEW customer AS SELECT id, first_name, last_name FROM client;" — a view provides backward compatibility while the table is renamed.

**Sources:** Martin Fowler - Blue Green Deployment (Evidence 4), ThoughtWorks - Evolutionary Database Design (Evidence 4, Source 9)

**Confidence:** HIGH

### Finding 5: Health Check Depth Requirements

**Claim:** A container that is merely "running" does not guarantee the application is ready to serve traffic — health checks must verify database, cache, and dependency availability.

**Evidence:** "container that is merely 'running' doesn't mean the application is ready to accept traffic... The readiness probe checks if your application is actually ready to serve requests." Laravel's health route documentation states: "you may perform additional health checks relevant to your application. Within a listener for [DiagnosingHealth] event, you may check your application's database or cache status."

**Sources:** Kubernetes Documentation - Liveness, Readiness, and Startup Probes (Evidence 5, Source 3), Laravel Documentation - Deployment (Evidence 5, Source 6)

**Confidence:** HIGH

### Finding 6: Pod Termination Graceful Shutdown

**Claim:** Pod termination in Kubernetes follows a graceful shutdown flow: SIGTERM sent → grace period (default 30 seconds) → SIGKILL if process doesn't exit.

**Evidence:** "A pod receives a 'grace period' during which it can perform cleanup operations... The kubelet sends a TERM signal to all of the Pod's containers... Then, the kubelet sends a SIGKILL signal to the container, if it has not gone down for a while." Kubernetes documentation also notes: "when the Pod is deleted, the corresponding endpoint in the EndpointSlice will update its conditions: the endpoint ready condition will be set to false."

**Sources:** Kubernetes Documentation - Pod Lifecycle (Evidence 6, Source 14), Kubernetes Documentation - Pod Termination

**Confidence:** HIGH

### Finding 7: NGINX Health Check Capabilities

**Claim:** NGINX provides both passive and active health checks; active health checks require NGINX Plus (commercial) while passive health checks are available in open source.

**Evidence:** Passive checks use `max_fails` and `fail_timeout` to temporarily remove servers. Active checks use the `health_check` directive to periodically probe servers. The `mandatory` parameter ensures new servers pass health checks before receiving traffic. `slow_start` parameter allows gradual recovery: "A recently recovered server can be easily overwhelmed by connections."

**Sources:** NGINX Documentation - HTTP Health Checks (Evidence 7, Source 5)

**Confidence:** HIGH

### Finding 8: Laravel Health Endpoint

**Claim:** Laravel provides a built-in health endpoint at `/up` (default) that returns 200 if the application booted without exceptions, and allows custom health checks via the `DiagnosingHealth` event.

**Evidence:** "Laravel includes a built-in health check route... By default the health check route is served at /up and will return a 200 HTTP response if the application has booted without exceptions. Otherwise, a 500 HTTP response will be returned." "you may check your application's database or cache status. If you detect a problem with your application, you may simply throw an exception from the listener."

**Sources:** Laravel Documentation - Deployment (Evidence 8, Source 6)

**Confidence:** HIGH

### Finding 9: Horizon Graceful Queue Worker Termination

**Claim:** Laravel Horizon provides graceful termination for queue workers via `php artisan horizon:terminate`, with Supervisor's `stopwaitsecs` ensuring long-running jobs complete before shutdown.

**Evidence:** "You may gracefully terminate the Horizon process using the `horizon:terminate` Artisan command. Any jobs that are currently being processed will be completed and then Horizon will stop executing." "During your application's deployment process, you should instruct the Horizon process to terminate so that it will be restarted by your process monitor and receive your code changes." "you should ensure that the value of `stopwaitsecs` is greater than the number of seconds consumed by your longest running job."

**Sources:** Laravel Documentation - Horizon (Evidence 9, Source 8)

**Confidence:** HIGH

### Finding 10: PostgreSQL Safe Column Addition

**Claim:** PostgreSQL's `ALTER TABLE ... ADD COLUMN` with a constant default value does not rewrite the table — it is a metadata-only operation that is safe for concurrent access.

**Evidence:** "Adding a column with a constant default value does not require each row of the table to be updated when the ALTER TABLE statement is executed. Instead, the default value will be returned the next time the row is accessed, and applied when the table is rewritten, making the ALTER TABLE very fast even on large tables."

**Sources:** PostgreSQL Documentation - Modifying Tables (Evidence 10, Source 10)

**Confidence:** HIGH

### Finding 11: Container Stop Signal Handling

**Claim:** Docker container stop sends SIGTERM and waits for a grace period (10s Linux, 30s Windows) before sending SIGKILL. The grace period is configurable via `--time` flag or Dockerfile's STOPSIGNAL.

**Evidence:** "The main process inside the container will receive SIGTERM, and after a grace period, SIGKILL. The first signal can be changed with the STOPSIGNAL instruction in the container's Dockerfile, or the --stop-signal option to docker run" "The default timeout can be specified using the --stop-timeout option when creating the container." "The default timeout... is 10 seconds for Linux containers, and 30 seconds for Windows containers."

**Sources:** Docker Documentation - docker container stop (Evidence 11, Source 12)

**Confidence:** HIGH

### Finding 12: Expand-Deploy-Migrate-Contract Database Pattern

**Claim:** Database migrations should follow the expand-deploy-migrate-contract pattern: add new columns first (expand), deploy code that uses both (deploy), migrate existing data, then remove old columns (contract).

**Evidence:** From Blue Green Deployment: "The trick is to separate the deployment of schema changes from application upgrades. So first apply a database refactoring to change the schema to support both the new and old version of the application, deploy that, check everything is working fine so you have a rollback point, then deploy the new version of the application. (And when the upgrade has bedded down remove the database support for the old version.)"

**Sources:** Martin Fowler - Blue Green Deployment (Evidence 12, Source 4), ThoughtWorks - Evolutionary Database Design (Evidence 12, Source 9)

**Confidence:** HIGH

### Finding 13: Kubernetes Rolling Update Parameters

**Claim:** The `maxUnavailable` and `maxSurge` parameters in Kubernetes rolling updates control the trade-off between speed and availability during deployment.

**Evidence:** "For the RollingUpdate strategy, these parameters control how Kubernetes performs the update: maxUnavailable (default 25%) — Maximum number of Pods that can be unavailable during the update; maxSurge (default 25%) — Maximum number of extra Pods that can be created during the update." "Kubernetes calculates percentages from the desired replica count, rounding down for maxUnavailable and rounding up for maxSurge."

**Sources:** Kubernetes Documentation - Update a Deployment Without Downtime (Evidence 13, Source 11)

**Confidence:** HIGH

## Areas of Agreement

All sources agree on the fundamental principles of zero-downtime deployment:

1. **Traffic Continuity**: New instances must be healthy and receiving traffic before old instances are removed
2. **Backward Compatibility**: Old and new application versions must coexist during transition
3. **Database Safety**: Schema changes must not break the currently running version
4. **Graceful Shutdown**: In-flight requests and jobs must be allowed to complete
5. **Rollback Preparedness**: Mechanisms must exist to quickly revert to previous version
6. **Health Verification**: Automated checks must confirm instance readiness before traffic routing

There is universal agreement that deployment success is not measured by CI pipeline green status alone, but by production metrics remaining stable post-deployment.

## Areas of Disagreement

Disagreements are primarily contextual rather than conceptual:

1. **Deployment Pattern Choice**: Blue-green vs rolling deployment represents different architectural tradeoffs (infrastructure duplication vs. incremental replacement), not contradictory approaches
2. **Health Check Implementation**: Variances in active vs. passive checking capabilities stem from tool licensing and implementation layers (infrastructure vs. application level)
3. **Termination Timing Nuances**: Differences between theoretical endpoint removal and practical connection draining reflect real-world load balancer behaviors rather than contradictions in the core principles

These represent implementation considerations rather than fundamental conflicts in zero-downtime deployment theory.

## Limitations

1. **Environment Specificity**: Evidence primarily reflects Kubernetes, Docker, Laravel, PostgreSQL, and NGINX ecosystems; patterns may differ in other technology stacks
2. **Version Specificity**: Some features (e.g., probe-level terminationGracePeriodSeconds in Kubernetes v1.28+) require specific version minimums
3. **Scale Considerations**: Very high-traffic systems may require additional considerations like circuit breakers and gradual traffic shifting (canary) not covered in the core research
4. **Data Migration Complexity**: While the expand-deploy-migrate-contract pattern is well-established, complex data transformations (splitting/merging columns, changing data types) may still require careful planning and testing
5. **Observability Gaps**: The research focuses on deployment mechanics but less on post-deployment validation techniques beyond basic health checks

## Conclusion

Zero-downtime deployment is achievable through a combination of well-established patterns: rolling or blue-green deployment strategies for instance replacement, liveness/readiness probes for traffic routing control, expand-deploy-migrate-contract for database schema changes, and coordinated graceful shutdown handling at all levels (container, application, queue worker). 

The most critical insight is that successful zero-downtime deployment requires planning for version coexistence — ensuring that v1 and v2 can run simultaneously, that database migrations are backward-compatible, and that traffic only shifts to new instances after they pass comprehensive health verification. Rollback capability must be designed in advance, not as an afterthought.

For the specific Laravel/PostgreSQL/Redis stack outlined in the lab, the recommendation is to:
1. Use Kubernetes Deployments with RollingUpdate strategy (maxUnavailable=25%, maxSurge=25%)
2. Implement Laravel's `/up` health endpoint enhanced with database/Redis connectivity checks
3. Apply database changes using expand-deploy-migrate-contract with intermediate verification
4. Configure Laravel Horizon with appropriate `stopwaitsecs` for longest job duration
5. Ensure containers handle SIGTERM to complete in-flight requests
6. Pre-validate rollback procedures in staging environments

This approach provides production-grade zero-downtime deployment capability while maintaining operational simplicity and rollback safety.