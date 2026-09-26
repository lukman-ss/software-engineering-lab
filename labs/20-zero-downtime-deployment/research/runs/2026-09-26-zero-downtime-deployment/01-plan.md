# Research Plan: Zero-Downtime Deployment

## Research Topic

Zero-Downtime Deployment — Deploy Versi Baru Tanpa Membuat User Tahu Ada Deployment

## Objective

To investigate and document the strategies, patterns, and best practices for implementing zero-downtime deployments in Laravel applications with PostgreSQL and Redis Queue backends. The focus is on ensuring continuous service availability during application upgrades while maintaining database compatibility.

## Research Questions

1. What are the primary zero-downtime deployment patterns (rolling, blue-green)?
2. How do liveness and readiness probes differ and when should each be used?
3. What are the patterns for graceful shutdown and connection draining?
4. How should database migrations be designed for backward compatibility during deployment?
5. What are the strategies for deploying queue workers without job loss?
6. How should rollback mechanisms be designed and tested?
7. What are Laravel/PostgreSQL/Redis-specific considerations for zero-downtime deployments?

## Search Strategy

### Primary Sources (Tier 1)
- Kubernetes official documentation on deployments, probes, and pod lifecycle
- Kubernetes documentation on rolling updates and blue-green deployments
- Docker documentation on container stop signals and graceful shutdown
- Laravel documentation on deployments and queue workers
- PostgreSQL documentation on schema modifications
- ThoughtWorks articles on evolutionary database design

### Secondary Sources (Tier 2)
- Martin Fowler's articles on deployment patterns
- NGINX documentation on health checks
- AWS documentation on deployment strategies

### Tertiary Sources (Tier 3)
- Community articles and blogs (used for discovery only, not primary evidence)

## Expected Primary Sources

1. Kubernetes Deployment documentation: rolling update strategy, maxSurge/maxUnavailable
2. Kubernetes Probe documentation: liveness vs readiness distinction
3. Docker container stop behavior: SIGTERM → grace period → SIGKILL
4. Laravel deployment guide: health endpoints, optimization commands
5. Laravel Horizon deployment: graceful termination of queue workers
6. PostgreSQL ALTER TABLE behavior: additive vs destructive changes
7. Evolutionary Database Design: expand-deploy-migrate-contract pattern

## Risks / Unknowns

1. **Database Migration Backward Compatibility**: Not all database changes can be made backward-compatible. Need to identify which patterns work and which require downtime.

2. **Queue Worker State**: Long-running jobs may be interrupted during deployment. Need to understand how Laravel handles job timeouts and graceful shutdown.

3. **Connection Pooling**: How Nginx and Laravel handle existing connections during deployment.

4. **Health Check Timing**: The gap between container start and application readiness.

5. **Redis Queue Data Persistence**: What happens to queued jobs during deployment of queue workers.
