# Open Questions

## Unanswered Questions

### 1. Optimal Sticky Routing Duration
**Question**: What is the optimal time duration (e.g., 2 s vs. 5 s vs. 10 s) for sticky primary routing after a write operation in production workloads?
**Current Evidence**: 5 seconds is widely cited in tutorials and the lab specification as a rule-of-thumb, but database vendor documentation only shows raw replication lag in seconds/bytes without recommending a fixed routing window.
**Next Research Direction**: Benchmark application latency vs. primary write load under different sticky session durations.

### 2. Multi-Region Replication Lag Distribution
**Question**: In cross-region read replica configurations (e.g., US-East to EU-West), what is the typical 99th percentile replication lag under peak write traffic?
**Current Evidence**: Azure and AWS state that cross-region lag is higher due to network latency, but quantitative p95/p99 distributions are not provided in official documentation.
**Next Research Direction**: Review empirical benchmarks and cloud provider SLA monitoring dashboards for cross-region replication.

### 3. Automated LSN-Based Standby Selection Overhead
**Question**: What is the performance overhead of polling `pg_last_wal_receive_lsn` or `SHOW REPLICA STATUS` before every read query to achieve monotonic reads?
**Current Evidence**: MongoDB includes `operationTime` in session metadata with low overhead, but doing this via SQL round-trips in PostgreSQL/MySQL may add significant latency.
**Next Research Direction**: Measure latency overhead of middleware that tracks LSN via connection-level heartbeat vs. per-query check.

### 4. Handling Split-Brain and Stale Reads During Promotion
**Question**: When a replica is promoted to primary during automatic failover, how can the application guarantee no in-flight reads receive stale data from demoted nodes?
**Current Evidence**: Vitess, RDS, and Azure document promotion workflows, but client connection draining and DNS/virtual endpoint propagation have a brief window of inconsistency.
**Next Research Direction**: Investigate consensus-based connection proxies (e.g., ProxySQL with Consul, Vitess VTGate) during failover transitions.

## Weak Evidence Areas

- **Heuristic-based timeouts**: Using fixed time windows (e.g., 5s) instead of deterministic LSN checks is an approximation that can fail under extreme network congestion.
- **ORM-level replica load balancing**: GORM DBResolver implements only Random policy by default; production round-robin or least-lag policies are less documented.

## Claims Needing Deeper Research

- **Aurora Storage-Level Lag**: While AWS claims low replication lag due to shared storage, compute nodes still maintain local buffer caches that must be invalidated via redo log streams, introducing minor read-after-write anomalies.
- **Client Session Scalability**: Storing write timestamps or LSN tokens in stateless JWTs vs. distributed Redis caches.

## Research Date
- **Date**: 2026-09-28
- Topic remains active as distributed SQL and serverless databases evolve read-replica consistency mechanisms.
