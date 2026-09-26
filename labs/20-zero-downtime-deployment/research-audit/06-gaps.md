# Research Gap Analysis

## Gap 1

Type: MISSING_CASE

Severity: MEDIUM

Location: `05-report.md` (Finding 2: Graceful Shutdown)

Problem:
The research describes Kubernetes Pod termination as sending SIGTERM, waiting for the grace period, and detaching endpoints. However, in Kubernetes, endpoint deregistration (EndpointSlice controller -> kube-proxy / ingress controller iptables update) and kubelet delivering SIGTERM to the container occur concurrently and asynchronously. If the application handles SIGTERM immediately and stops accepting connections before the network proxy routing table updates, in-flight traffic routed to the pod during those milliseconds will fail with connection resets (HTTP 502/504). To achieve true zero downtime, a `preStop` lifecycle hook with a brief sleep (e.g., 5-15 seconds) is usually required before processing SIGTERM.

Required Revision:
Incorporate the asynchronous endpoint propagation race condition and the `preStop` sleep mitigation pattern into the lab architecture specification.

Can Be Approved Without Fix: YES

---

## Gap 2

Type: MISSING_CASE

Severity: MEDIUM

Location: `05-report.md` (Finding 3: Database and State Backward Compatibility)

Problem:
While the Expand/Contract pattern prevents application-level schema incompatibility, it omits relational database DDL locking dynamics. In engines like PostgreSQL, executing `ALTER TABLE ... ADD COLUMN` requires an `ACCESS EXCLUSIVE` lock. If long-running queries are active on that table, the DDL query queues behind them, blocking all subsequent incoming reads and writes, potentially exhausting the database connection pool and creating application downtime.

Required Revision:
Include guidelines for database lock timeouts (`SET lock_timeout = '2s'`) and safe schema migration practices in the lab's migration section.

Can Be Approved Without Fix: YES

---

## Gap 3

Type: OVERGENERALIZATION

Severity: LOW

Location: `03-evidence.md` (Evidence 5), `05-report.md` (Finding 2)

Problem:
The graceful restart pattern for background workers is illustrated solely via Laravel's cache-driven `queue:restart` command. While accurate for Laravel, universal background workers (Go daemons, Celery, Sidekiq) typically rely on OS signal trapping (`SIGTERM` or `SIGQUIT`) to finish active jobs and exit.

Required Revision:
Clarify that while cache timestamp signals are standard in Laravel, direct POSIX signal handling (`SIGTERM`/`SIGQUIT`) is the generic distributed worker standard.

Can Be Approved Without Fix: YES
