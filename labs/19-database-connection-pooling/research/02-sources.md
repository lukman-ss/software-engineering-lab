# Sources

## Source 1
Title: 19.3. Connections and Authentication — PostgreSQL 18 Documentation
Publisher: PostgreSQL Global Development Group
URL: https://www.postgresql.org/docs/current/runtime-config-connection.html
Published: 2026 (PostgreSQL 18 documentation)
Accessed: 2026-09-26
Source Tier: Tier 1 — Official documentation
Relevance: Defines `max_connections`, `reserved_connections`, `superuser_reserved_connections`, TCP keepalive, authentication timeout. Primary source for PostgreSQL connection configuration semantics.

## Source 2
Title: Number Of Database Connections — PostgreSQL Wiki
Publisher: PostgreSQL Community (wiki.postgresql.org)
URL: https://wiki.postgresql.org/wiki/Number_Of_Database_Connections
Published: Last edited 14 March 2014
Accessed: 2026-09-26
Source Tier: Tier 1 — Official PostgreSQL community wiki
Relevance: Explains the "knee" in performance vs connections curve, lists reasons for degradation (disk contention, RAM, lock contention, context switches, cache line contention, O(N²) scaling), and provides the pool sizing formula.

## Source 3
Title: About Pool Sizing — HikariCP Wiki
Publisher: HikariCP Project (Brett Wooldridge)
URL: https://github.com/brettwooldridge/HikariCP/wiki/About-Pool-Sizing
Published: Last edited Dec 7, 2021
Accessed: 2026-09-26
Source Tier: Tier 1 — Official project documentation
Relevance: Pool sizing formula, Oracle Real-World Performance 50x improvement demonstration, axiom about saturated pools, pool-locking deadlock avoidance formula.

## Source 4
Title: HikariCP — Configuration (knobs, baby!)
Publisher: HikariCP Project (Brett Wooldridge)
URL: https://github.com/brettwooldridge/HikariCP#configuration-knobs-baby
Published: Active repository (HikariCP 7.1.0)
Accessed: 2026-09-26
Source Tier: Tier 1 — Official project documentation
Relevance: All HikariCP configuration parameters: maximumPoolSize (default 10), connectionTimeout (default 30s), idleTimeout (default 10min), maxLifetime (default 30min), keepaliveTime (default 2min), leakDetectionThreshold (default 0, min 2000ms).

## Source 5
Title: PgBouncer Configuration — pgbouncer.ini
Publisher: PgBouncer Project
URL: https://www.pgbouncer.org/config.html
Published: Active (PgBouncer 1.26.0 as of Sep 2026)
Accessed: 2026-09-26
Source Tier: Tier 1 — Official documentation
Relevance: Pool modes (session/transaction/statement), default_pool_size (20), max_client_conn (100), server_reset_query, timeout settings, connection lifecycle management.

## Source 6
Title: PgBouncer Features
Publisher: PgBouncer Project
URL: https://www.pgbouncer.org/features.html
Published: Active
Accessed: 2026-09-26
Source Tier: Tier 1 — Official documentation
Relevance: Low memory (2kB per connection), three pooling modes with feature compatibility matrix, online reconfiguration support.

## Source 7
Title: Limits in Azure Database for PostgreSQL flexible server
Publisher: Microsoft Azure
URL: https://learn.microsoft.com/en-us/azure/postgresql/flexible-server/concepts-limits
Published: 2026-07-08
Accessed: 2026-09-26
Source Tier: Tier 1 — Official documentation
Relevance: Default max_connections by vCore tier (50-5000), 15 connections reserved for replication/monitoring, recommendation to use built-in PgBouncer transaction mode with 2-5× vCores, error "FATAL: sorry, too many clients already."

## Source 8
Title: Manage database connections — Google Cloud SQL for PostgreSQL
Publisher: Google Cloud
URL: https://cloud.google.com/sql/docs/postgres/manage-connections
Published: Active (2026)
Accessed: 2026-09-26
Source Tier: Tier 1 — Official documentation
Relevance: Best practices for connection pooling with HikariCP (Java), Go, Python, Node.js, Ruby, C#, PHP. Code samples for TCP and Unix socket connection pools.

## Source 9
Title: RDS Connection Limits — Amazon RDS User Guide
Publisher: Amazon Web Services
URL: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/CHAP_Limits.html
Published: Active (2026)
Accessed: 2026-09-26
Source Tier: Tier 1 — Official documentation
Relevance: PostgreSQL `max_connections = LEAST(DBInstanceClassMemory/9531392, 5000)` formula, recommendation to use RDS Proxy for connection pooling.

## Source 10
Title: Dropwizard Metrics — HikariCP Wiki
Publisher: HikariCP Project
URL: https://github.com/brettwooldridge/HikariCP/wiki/Dropwizard-Metrics
Published: Active
Accessed: 2026-09-26
Source Tier: Tier 1 — Official project documentation
Relevance: Pool monitoring metrics: Wait Timer (acquisition time), Usage Histogram (in-use time), TotalConnections, IdleConnections, ActiveConnections, PendingConnections.

## Source 11
Title: PostgreSQL Cumulative Statistics System
Publisher: PostgreSQL Global Development Group
URL: https://www.postgresql.org/docs/current/monitoring-stats.html
Published: 2026 (PostgreSQL 18)
Accessed: 2026-09-26
Source Tier: Tier 1 — Official documentation
Relevance: `pg_stat_activity` for monitoring active connections, `idle_session_timeout` parameter, `pg_stat_database` for database-level stats.

## Source 12
Title: TCP keepalive for a better PostgreSQL experience
Publisher: CYBERTEC PostgreSQL
URL: https://www.cybertec-postgresql.com/en/tcp-keepalive-for-a-better-postgresql-experience/
Published: 2021-10-27
Accessed: 2026-09-26
Source Tier: Tier 2 — Reputable technical publication
Relevance: TCP keepalive configuration for detecting broken connections, OS-level settings for Linux/macOS/Windows, relevance to connection leak prevention.
