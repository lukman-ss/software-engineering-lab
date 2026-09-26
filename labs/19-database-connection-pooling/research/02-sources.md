# Sources

## Source 1
Title: PostgreSQL Documentation: Connections and Authentication
Publisher: The PostgreSQL Global Development Group
URL: https://www.postgresql.org/docs/current/runtime-config-connection.html
Published: N/A (Latest up to PostgreSQL 18/19 Beta)
Accessed: 2026-09-25
Source Tier: Tier 1
Relevance: Official specification for `max_connections`, `reserved_connections`, TCP behavior, TLS overhead, and database process constraints.

## Source 2
Title: About Pool Sizing
Publisher: HikariCP Wiki / Brett Wooldridge (Referencing Oracle Real-World Performance Group)
URL: https://github.com/brettwooldridge/HikariCP/wiki/About-Pool-Sizing
Published: 2021-12-07 (Last edit)
Accessed: 2026-09-25
Source Tier: Tier 1 (Highly authoritative for JVM connection pooling mechanics and universal database sizing math)
Relevance: Key mathematical formulas for pool sizing, concurrency limits, and deadlock avoidance; analysis of throughput collapse when connections exceed core counts.

## Source 3
Title: Number Of Database Connections
Publisher: PostgreSQL Wiki
URL: https://wiki.postgresql.org/wiki/Number_Of_Database_Connections
Published: 2014-03-14 (Historical, foundational architecture reference)
Accessed: 2026-09-25
Source Tier: Tier 1 / Tier 2
Relevance: Primary diagnostic explanation of the "knee" in transaction throughput, context switching, cache line contention, RAM/`work_mem` exhaustion, and disk thrashing under high connections.

## Source 4
Title: PgBouncer Features
Publisher: PgBouncer Authors
URL: https://www.pgbouncer.org/features.html
Published: N/A
Accessed: 2026-09-25
Source Tier: Tier 1
Relevance: Documentation of proxy-layer pooling architectures (Session, Transaction, Statement pooling), low memory footprint (2kB per connection), and transaction state isolation.

## Source 5
Title: PostgreSQL Documentation: The Cumulative Statistics System (pg_stat_activity)
Publisher: The PostgreSQL Global Development Group
URL: https://www.postgresql.org/docs/current/monitoring-stats.html
Published: N/A (Latest up to PostgreSQL 18/19 Beta)
Accessed: 2026-09-25
Source Tier: Tier 1
Relevance: Defines internal tracking for `pg_stat_activity` states (`active`, `idle`, `idle in transaction`), wait events, and transaction boundaries crucial for diagnosing leaks and exhaustion.
