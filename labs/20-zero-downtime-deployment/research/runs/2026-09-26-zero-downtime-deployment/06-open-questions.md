# Open Questions: Zero-Downtime Deployment

---

## Unanswered Questions

### 1. PHP-FPM Graceful Shutdown Behavior with SIGTERM
**Question:** Does PHP-FPM handle SIGTERM gracefully by default, or does it require specific configuration to complete in-flight requests before exit?

**Context:** The research confirms Docker sends SIGTERM → grace period → SIGKILL, and Kubernetes follows the same pattern. However, PHP-FPM's signal handling behavior is not explicitly documented in the sources consulted. PHP-FPM may need `process_control_timeout` configuration and specific signal handling in the application code.

**Potential Next Research:** Test PHP-FPM signal handling behavior with and without Laravel Octane/FrankenPHP, which may provide better graceful shutdown support.

### 2. Laravel Application-Level SIGTERM Handling
**Question:** How should Laravel applications handle SIGTERM to complete in-flight HTTP requests gracefully?

**Context:** The research documents infrastructure-level graceful shutdown (Kubernetes, Docker, Supervisor) but does not cover application-level signal handling in PHP/Laravel. PHP applications running under Nginx + PHP-FPM may not have direct access to SIGTERM signals.

**Potential Next Research:** Investigate Laravel Octane, FrankenPHP, or middleware-based approaches for request draining. Check if Nginx's `proxy_next_upstream` and upstream health checks can provide sufficient connection draining without application-level signal handling.

### 3. Redis Queue Job Persistence During Worker Restart
**Question:** What happens to jobs currently being processed by Redis queue workers when workers are restarted during deployment?

**Context:** Laravel Horizon documentation states "Any jobs that are currently being processed will be completed and then Horizon will stop executing." However, the plain `queue:work` documentation warns that `block_for=0` prevents signal handling. The interaction between Redis `retry_after` setting, job timeouts, and graceful termination needs clarification.

**Potential Next Research:** Test job behavior when Horizon terminates vs. when plain workers receive SIGTERM. Document the exact sequence of job completion, requeue, or failure based on configuration.

### 4. Database Migration Lock Behavior During Concurrent Access
**Question:** Do PostgreSQL schema migrations (even safe ones like ADD COLUMN) acquire locks that could block application queries during zero-downtime deployment?

**Context:** PostgreSQL documentation states ADD COLUMN with constant default is metadata-only and fast. However, it doesn't explicitly state whether it acquires any locks that could block concurrent reads/writes. The Evolutionary Database Design article mentions migration ordering but not lock contention.

**Potential Next Research:** Investigate PostgreSQL lock types for each ALTER TABLE operation. Determine if `LOCK_TIMEOUT` settings or `pg_stat_activity` monitoring is needed during migration windows.

### 5. Nginx Upstream Server Removal Timing
**Question:** How long does Nginx take to stop sending new connections to an upstream server marked as unhealthy via passive health checks, and does it wait for existing keep-alive connections?

**Context:** NGINX documentation describes passive health checks with `fail_timeout` and `max_fails`, but doesn't specify the exact timing of when a server is removed from the load balancing pool versus when existing connections are drained. The `proxy_next_upstream` directive behavior during health check failures needs clarification.

**Potential Next Research:** Test Nginx passive health check timing with actual request patterns. Compare with active health check (NGINX Plus) behavior.

### 6. Kubernetes Endpoint Removal vs. Connection Draining Gap
**Question:** What is the actual time window between Kubernetes marking a pod as not-ready (readiness probe failure) and the load balancer completely stopping traffic to that pod?

**Context:** Kubernetes documentation states the endpoint ready condition is set to false immediately upon pod deletion. However, cloud load balancers (AWS ALB/NLB, GCP, Azure) have their own drain timeout configurations. The interaction between Kubernetes endpoint updates and cloud LB connection draining is not fully documented.

**Potential Next Research:** Measure the end-to-end timing in major cloud providers. Document required configuration alignment between `terminationGracePeriodSeconds`, readiness probe intervals, and cloud LB drain timeouts.

### 7. Laravel Octane/FrankenPHP Zero-Downtime Deployment Patterns
**Question:** How does zero-downtime deployment change when using Laravel Octane (Swoole/RoadRunner) or FrankenPHP instead of traditional Nginx + PHP-FPM?

**Context:** The research focuses on traditional Nginx + PHP-FPM deployment. Octane and FrankenPHP run long-lived PHP processes that handle multiple requests, potentially requiring different graceful shutdown approaches and worker management.

**Potential Next Research:** Investigate Octane's `octane:stop` behavior, RoadRunner's graceful shutdown, and FrankenPHP's worker lifecycle management in the context of zero-downtime deployments.

### 8. Multi-Region/Cross-AZ Deployment Coordination
**Question:** How should zero-downtime deployments be coordinated across multiple availability zones or regions to ensure consistent database schema state?

**Context:** The research assumes single-region or single-cluster deployment. Multi-region deployments require coordination of database migrations across regions, potentially with different deployment timing and rollback complexity.

**Potential Next Research:** Investigate blue-green or rolling deployment patterns across regions with shared database. Consider database migration tools that support multi-region coordination (e.g., Flyway with multiple targets, custom migration orchestration).

---

## Weak Evidence Areas

### 1. Actual Measured Downtime for Rolling Deployments
**Gap:** The research documents the theoretical availability (maxUnavailable=25%) but lacks empirical measurements of actual downtime experienced during rolling updates with various configurations.

**Next Steps:** Conduct load testing during rolling deployments to measure actual error rates and latency impact.

### 2. Laravel Health Check Performance Under Load
**Gap:** The `/up` endpoint with database/Redis checks adds latency to every health probe. The impact of frequent readiness probes (default 5s interval) on application performance is not quantified.

**Next Steps:** Benchmark health endpoint performance with and without database/cache checks at various probe intervals.

### 3. PostgreSQL ALTER TABLE Concurrent Performance
**Gap:** While PostgreSQL claims "very fast even on large tables" for ADD COLUMN with constant default, there's no quantitative data on lock duration or impact on concurrent transactions for tables of various sizes.

**Next Steps:** Test ALTER TABLE timing on tables with millions of rows under concurrent read/write load.

---

## Claims Needing Deeper Research

### 1. "Add Column with Constant Default is Always Safe"
**Claim:** PostgreSQL documentation states ADD COLUMN with constant default is metadata-only and fast.
**Uncertainty:** Does this hold for partitioned tables, tables with triggers, or tables with row-level security policies? What about when using `DEFAULT gen_random_uuid()` or other non-constant expressions?

### 2. "Horizon Terminates Gracefully Always"
**Claim:** Horizon's `terminate` command allows currently processing jobs to complete.
**Uncertainty:** What happens if a job exceeds `stopwaitsecs`? What if the job is in a non-interruptible system call? What if Redis connection fails during termination?

### 3. "Blue-Green Instant Rollback is Always Safe"
**Claim:** Blue-green deployment allows instant rollback by switching router back to old environment.
**Uncertainty:** What if database migrations have already been applied and are not backward-compatible? The "instant rollback" only applies to application code, not database state.

### 4. "Readiness Probe Failure Immediately Removes Traffic"
**Claim:** Kubernetes removes pod from endpoints immediately on readiness failure.
**Uncertainty:** EndpointSlice controller reconciliation timing, kube-proxy update latency, and load balancer propagation delays may add seconds to minutes of continued traffic routing.

---

## Possible Next Research Directions

### Direction 1: Empirical Validation Studies
- Conduct controlled experiments measuring actual downtime/error rates during deployments
- Test various maxUnavailable/maxSurge configurations under load
- Benchmark database migration timing with concurrent application traffic

### Direction 2: Advanced Deployment Patterns
- Canary deployment with weighted traffic shifting
- Feature flag integration for gradual functionality rollout
- Progressive delivery with automated rollback on metric anomalies

### Direction 3: Database Migration Tooling
- Compare Flyway, Liquibase, Laravel's native migrations for zero-downtime safety
- Investigate automated backward-compatibility checking for migrations
- Research reversible migration patterns for emergency rollback

### Direction 4: Observability and Validation
- Post-deployment validation strategies beyond basic health checks
- Automated smoke testing in deployment pipeline
- Real-user monitoring (RUM) integration for deployment verification

### Direction 5: Laravel-Specific Patterns
- Laravel Octane/FrankenPHP worker management for zero-downtime
- Vapor (serverless) deployment patterns
- Horizon supervisor configuration for different job profiles

---

## Priority Questions for Lab Implementation

For the specific lab exercise (Nginx + Laravel + PostgreSQL + Redis Queue, v1.4 → v1.5 with new DB column and worker changes), the highest-priority questions are:

1. **What is the exact sequence of commands for a zero-downtime deployment?**
   - Migration first? Workers first? Web servers first?
   - What commands, in what order, with what verification steps?

2. **How to handle the new column and new worker simultaneously?**
   - v1.4 reads old column, v1.5 reads new column
   - Both workers process jobs — how to ensure queue compatibility?

3. **What rollback procedure if v1.5 has issues?**
   - Database column already added — how to rollback?
   - Workers already processing jobs — how to drain?

4. **What health checks to implement in Laravel?**
   - Database connectivity?
   - Redis connectivity?
   - Migration status check?
   - Queue worker health?

5. **Nginx configuration for graceful drain?**
   - `proxy_read_timeout`?
   - `proxy_next_upstream`?
   - Upstream `max_fails`/`fail_timeout`?

These implementation-specific questions should guide the engineering phase following this research.