# Research Plan

## Research Topic
Database Connection Pooling — Connection Pool, Connection Exhaustion, Connection Leak, Pool Size, Connection Wait Time, Database Capacity, Monitoring, Production Troubleshooting

## Objective
Investigate how database connection pooling works, why applications experience connection exhaustion despite database health, the mechanics of connection leaks, optimal pool sizing strategies, and evidence-based monitoring and troubleshooting approaches for production environments.

## Research Questions
1. What is the relationship between `max_connections` and application-side pool configuration, and what happens when they are mismatched?
2. What is the optimal pool size formula, and what evidence supports it?
3. How do connection leaks occur, and what are the observable symptoms versus root causes?
4. What are the mechanisms of external connection poolers (PgBouncer) and how do they mitigate exhaustion?
5. What monitoring metrics are necessary to detect connection pool issues before they cause outages?
6. What are cloud-provider-specific limits and recommendations for connection management?

## Search Strategy
- Query: "database connection pooling production troubleshooting"
- Query: "HikariCP pool sizing formula optimal connections"
- Query: "PostgreSQL max_connections connection exhaustion"
- Query: "PgBouncer connection pooler configuration transaction mode"
- Query: "connection leak detection JDBC pool"
- Query: "database connection pool monitoring metrics active idle wait"
- Query: "AWS RDS max_connections connection pooling limits"
- Query: "Azure PostgreSQL connection limits PgBouncer"
- Query: "Google Cloud SQL managed connection pooling"

Expected Primary Sources:
- PostgreSQL official documentation (runtime-config-connection)
- PostgreSQL wiki (Number Of Database Connections)
- HikariCP official documentation (About Pool Sizing, Configuration)
- PgBouncer official documentation (Configuration, Features)
- Azure Database for PostgreSQL limits documentation
- Google Cloud SQL database connection management docs
- AWS RDS limits documentation

## Risks / Unknowns
- The pool sizing formula `(core_count * 2) + effective_spindle_count` was designed for spinning disks; its applicability to SSDs is explicitly noted as unanalyzed in the source.
- The Oracle Real-World Performance video demonstrating 50x improvement from pool reduction was cited by HikariCP but not independently verified (secondary citation).
- Cloud provider documentation may have changed since publication; Azure docs were last updated 2026-07-08.
- Connection leak detection thresholds are framework-specific (HikariCP); equivalent mechanisms in other pools (e.g., c3p0, DBCP) were not investigated.
- The PostgreSQL wiki article was last edited in 2014 and may not reflect recent PostgreSQL versions' changes.
