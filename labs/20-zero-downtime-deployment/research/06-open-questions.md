# Open Questions

## Unanswered Questions

1. **Docker Compose Rolling Update Documentation URL** (LOW confidence - missing source)
   - The URL https://docs.docker.com/compose/how-tos/rolling-update/ intended as a source returned 404.
   - Need: Updated/verified Docker Compose reference for zero-downtime compose deployments.
   - Impact: The lab may need to use Docker Swarm or external orchestrator (Kubernetes, Nomad, Ansible) for rolling updates if pure docker-compose is insufficient.

2. **NGINX Open Source Active Health Checks** (MEDIUM confidence)
   - Active health checks, mandatory health checks, and slow_start are NGINX Plus-only features.
   - For NGINX OSS, passive health checks (max_fails/fail_timeout) work at upstream level but cannot block new server from receiving traffic before health passes.
   - Need: Mechanism for NGINX OSS "readiness" — either external orchestrator (Kubernetes), custom script using `nginx -s quit`/`SIGHUP` sequencing, or solution documented in NGINX docs.

3. **Laravel Queue Worker Graceful Shutdown Under SIGKILL** (MEDIUM confidence)
   - Topic spec implies graceful shutdown: "Load Balancer ↓ Stop kirim request baru ↓ Tunggu request aktif selesai ↓ Terminate instance"
   - Laravel `queue:restart` uses cache-based restart signal; workers exit after current job.
   - Need verification: What exactly happens on explicit SIGTERM/SIGKILL to worker process? Does PHP-FPM/Queue handler respect `SIGTERM` for graceful shutdown?
   - Pending deeper sources: Laravel queue worker signal handling docs, Supervisor `STOPWAITSEC`, PHP PCNTL behavior.

## Weak Evidence

4. **Redis Queue Durability Guarantees for In-Flight Jobs**
   - Queue docs show jobs are stored in Redis (or other backend) and re-queued on failure/timeout.
   - However, a job being PROCESSED when the worker dies (crash, OOM) may be lost unless `retry_after` expires and the job is re-queued.
   - Need: Redis documentation on job-at-scale / reliability models; Laravel Horizon's `retry_after` semantics.

5. **PostgreSQL Migration Lock Contention During Rolling Update**
   - ALTER TABLE with ACCESS EXCLUSIVE lock blocks all queries on the table for the lock duration.
   - For large tables, even "metadata-only" operations like ADD COLUMN can take seconds or minutes if constraints trigger validation.
   - Need: Real-world experience reports or official recommendations for ALTER TABLE on large production tables without blocking reads.

6. **Laravel Horizon vs Vanilla Redis Queue Behavior**
   - Topic specifies Redis Queue. Horizon is Laravel's Redis queue dashboard/monitoring.
   - Need clarification: Does the lab use Horizon, vanilla `queue:work`, or both? Behavior during restart may differ.

## Claims Needing Deeper Research

7. **Horizontal Pod Autoscaler (HPA) Interaction with Rolling Updates**
   - Kubernetes Deployments interact with HPA; mid-update scaling may affect rollout speed.
   - Need: Documentation on HPA behavior during rolling update (scale on old pods, scale on new pods, etc.).

8. **Zero-Downtime Schema Changes Without Locks in PostgreSQL 18**
   - PostgreSQL 18 (as of 2026-09-24) introduces new ALTER TABLE features.
   - Need: Verification of which operations are now lock-free or minimal-lock in v18 vs previously documented v16/v17.

9. **Connection Draining Across Nginx → Laravel → PostgreSQL → Redis Chain**
   - The chain is: client → Nginx → Laravel (PHP-FPM/FrankenPHP) → PostgreSQL/Redis.
   - Connection draining at Nginx level (graceful shutdown) must propagate: PHP-FPM must not leave transactions hanging; Redis connections from Laravel should be properly closed to allow queue re-processing.
   - Need: Laravel database transaction rollback on shutdown, Redis client connection cleanup behavior, PHP-FPM `listen.backlog` and `shutdown` docs.

## Possible Next Research Directions

- Fetch official Docker Compose reference for service update/rolling (alternative URL: docs.docker.com/compose/production/).
- Fetch Supervisor configuration best practices for zero-downtime deployment in Laravel context.
- Fetch PostgreSQL 18 ALTER TABLE incremental improvements documentation.
- Fetch Amazon EKS Blue/Green Deployment whitepaper or guide.
- Fetch Laravel Horizon documentation for queue worker lifecycle during deployments.
