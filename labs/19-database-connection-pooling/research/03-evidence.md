# Evidence

## Evidence 1
Claim: Increasing the connection pool size beyond a certain physical capacity limit degrades database throughput and increases response times.
Evidence: "Once all of the resources are in use, you won't push any more work through by having more connections competing for the resources. In fact, throughput starts to fall off due to the overhead from that contention." "reducing the connection pool size alone, in the absence of any other change, decreased the response times of the application from ~100ms to ~2ms -- over 50x improvement."
Source: PostgreSQL Wiki & HikariCP Wiki
URL: https://wiki.postgresql.org/wiki/Number_Of_Database_Connections & https://github.com/brettwooldridge/HikariCP/wiki/About-Pool-Sizing
Confidence: HIGH
Corroborated By: Source 3 (PostgreSQL Wiki) and Source 2 (HikariCP / Oracle Real-World Performance).
Notes: Both agree that queueing requests externally while maintaining a small active internal pool saturated to resource limits yields higher throughput.

## Evidence 2
Claim: The optimal formula for baseline database connection sizing relies primarily on CPU core count and storage I/O limits.
Evidence: "A formula which has held up pretty well across a lot of benchmarks for years is that for optimal throughput the number of active connections should be somewhere near ((core_count * 2) + effective_spindle_count)."
Source: HikariCP Wiki & PostgreSQL Wiki
URL: https://github.com/brettwooldridge/HikariCP/wiki/About-Pool-Sizing & https://wiki.postgresql.org/wiki/Number_Of_Database_Connections
Confidence: HIGH
Corroborated By: Both sources reference the identical formula derived from PostgreSQL performance testing.
Notes: The formula establishes a starting point for load testing; SSD environments (where spindle count is irrelevant) rely heavily on the `(core_count * 2)` baseline.

## Evidence 3
Claim: High connection counts degrade performance due to disk contention, memory exhaustion (work_mem), lock contention, context switching, and CPU cache line eviction.
Evidence: "Disk contention... RAM usage (work_mem RAM can be allocated for each node of a query on each connection)... Lock contention... Context switches... Cache line contention."
Source: PostgreSQL Wiki
URL: https://wiki.postgresql.org/wiki/Number_Of_Database_Connections
Confidence: HIGH
Corroborated By: Supported by PostgreSQL internals documentation regarding `max_connections` and shared memory allocation scaling.

## Evidence 4
Claim: Distributed deployments without an intermediate pooler easily exceed backend connection limits, resulting in starvation or exhaustion.
Evidence: When applications spawn multiple workers scaling horizontally, independent local pools multiply total active connections. PgBouncer supports "Transaction pooling: A server connection is assigned to a client only during a transaction... Low memory requirements (2 kB per connection by default)."
Source: PgBouncer Features
URL: https://www.pgbouncer.org/features.html
Confidence: HIGH
Corroborated By: Standard PostgreSQL architectural guidance to not handle high-volume distributed client limits directly via the core server.
Notes: In scenarios with 4 instances * 16 workers * 10 connections = 640 connections vs 200 `max_connections`, proxy pooling (transaction mode) collapses the 640 virtual connections into a small real pool.

## Evidence 5
Claim: Connection deadlocks (pool locking) occur when the application thread requires more simultaneous connections than the pool allows.
Evidence: "The calculation of pool size in order to avoid deadlock is a fairly simple resource allocation formula: pool size = Tn x (Cm - 1) + 1. Where Tn is the maximum number of threads, and Cm is the maximum number of simultaneous connections held by a single thread."
Source: HikariCP Wiki
URL: https://github.com/brettwooldridge/HikariCP/wiki/About-Pool-Sizing
Confidence: HIGH
Corroborated By: General distributed systems concurrency theory (resource starvation).
Notes: Underscores that pool size has a mathematical floor governed by application transaction design, and a ceiling governed by database hardware.

## Evidence 6
Claim: Idle connections held open due to application leaks or long-running non-database tasks waste slots, observable via specific state tracking.
Evidence: The `pg_stat_activity` view tracks `state` (active, idle, idle in transaction, idle in transaction (aborted)). "idle in transaction: The backend is in a transaction, but is not currently executing a query." Wait events like `ClientRead` indicate waiting on the application.
Source: PostgreSQL Documentation (pg_stat_activity)
URL: https://www.postgresql.org/docs/current/monitoring-stats.html
Confidence: HIGH
Corroborated By: Standard relational database operational principles.
Notes: "idle in transaction" is the primary indicator of a connection leak or an application holding a connection while awaiting external APIs/file generation.
