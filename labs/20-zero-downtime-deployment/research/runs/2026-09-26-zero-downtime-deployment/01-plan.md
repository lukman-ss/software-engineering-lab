# Research Topic
Zero-Downtime Deployment — Deploy Versi Baru Tanpa Membuat User Tahu Ada Deployment

# Objective
Investigate zero-downtime deployment patterns specifically for the Nginx + Laravel + PostgreSQL + Redis Queue stack, collecting authoritative evidence on deployment orchestration, backward-compatible database migrations, health probes, graceful shutdown/connection draining, queue worker deployment strategies, and rollback mechanics.

# Research Questions
1. What are the core zero-downtime deployment patterns (Rolling, Blue-Green, Canary) and their tradeoffs for this stack?
2. How must database schema changes be designed (Expand/Contract, backward-compatible migrations) to allow v1 and v2 coexistence?
3. What constitutes a proper Readiness vs Liveness probe for Laravel (PHP-FPM) and queue workers?
4. How to implement graceful shutdown/connection draining for Nginx upstream, PHP-FPM, and Laravel queue workers?
5. What is the correct deployment sequence for queue workers during rolling/blue-green deployments to avoid job loss?
6. How to achieve fast, safe rollback that restores v1 without data corruption?
7. What monitoring/observability signals confirm zero-downtime success?

# Search Strategy
1. Official docs: Laravel (deployment, queues, Octane), PostgreSQL (DDL concurrency, locking), Redis (rolling upgrade), Nginx (upstream health checks, graceful reload), Docker Compose (healthchecks), Kubernetes (probes, lifecycle).
2. Authoritative patterns: Martin Fowler (Blue-Green, Parallel Change), Kelsey Hightower / Kubernetes best practices.
3. Laravel-specific: Laravel Envoyer/Forge zero-downtime, spatie/laravel-health, spatie/laravel-horizon, queue worker signals.

# Expected Primary Sources (Tier 1)
- laravel.com/docs (deployment, queues, octane)
- postgresql.org/docs (DDL, concurrent index, locking)
- redis.io/docs (clustering, rolling upgrade)
- nginx.org/en/docs (upstream, health_check)
- kubernetes.io/docs (pod lifecycle, probes, termination)
- martinfowler.com (blue-green-deployment, parallel-change)

# Risks / Unknowns
- Laravel Octane (Swoole/RoadRunner) changes shutdown semantics vs PHP-FPM
- Exact SIGTERM handling in php-fpm vs octane workers
- Redis queue job persistence across worker restart (horizon vs raw queue:work)
- Coordinated deployment of web + scheduler + multiple queue workers
- Nginx upstream "drain" support without commercial Plus