# Evidence

## Evidence 1
Claim: PostgreSQL `max_connections` default is 100 and determines the maximum number of concurrent connections to the database server.
Evidence: PostgreSQL official documentation states: "Determines the maximum number of concurrent connections to the database server. The default is typically 100 connections, but might be less if your kernel settings will not support it."
Source: PostgreSQL 18 Documentation — 19.3. Connections and Authentication
URL: https://www.postgresql.org/docs/current/runtime-config-connection.html
Confidence: HIGH — Official primary documentation
Corroborated By: Azure documentation uses max_connections values in its limits table (50 for B1ms, up to 5000 for large instances)
Notes: The value can only be set at server start. PostgreSQL sizes certain resources (including shared memory) based on max_connections.

## Evidence 2
Claim: PostgreSQL reserves connection slots for superusers via `superuser_reserved_connections` (default 3) and `reserved_connections` (default 0).
Evidence: "Determines the number of connection 'slots' that are reserved for connections by PostgreSQL superusers. At most max_connections connections can ever be active simultaneously." And: "Determines the number of connection 'slots' that are reserved for connections by roles with privileges of the pg_use_reserved_connections role."
Source: PostgreSQL 18 Documentation — 19.3. Connections and Authentication
URL: https://www.postgresql.org/docs/current/runtime-config-connection.html
Confidence: HIGH — Official primary documentation
Corroborated By: Azure documentation states "15 connections reserved for physical replication and monitoring" with formula `max_connections - (reserved_connections + superuser_reserved_connections)`
Notes: Reserved connections ensure administrators can always connect even under connection exhaustion.

## Evidence 3
Claim: Database performance has a "knee" in the connections-vs-throughput curve, after which adding more connections causes performance degradation due to disk contention, RAM pressure, lock contention, context switches, cache line contention, and O(N²) scaling of internal structures.
Evidence: "If you look at any graph of PostgreSQL performance with number of connections on the x axis and tps on the y axis (with nothing else changing), you will see performance climb as connections rise until you hit saturation, and then you have a 'knee' after which performance falls off." Six specific reasons listed: disk contention, RAM usage (work_mem × connections), lock contention (spinlocks, LW locks), context switches, cache line contention, and general scaling (O(N²) or O(N*log(N)) internal structures).
Source: PostgreSQL Wiki — Number Of Database Connections
URL: https://wiki.postgresql.org/wiki/Number_Of_Database_Connections
Confidence: HIGH — Official PostgreSQL community wiki
Corroborated By: HikariCP wiki independently confirms the same principles, citing Computer Science fundamentals about time-slicing and context switching overhead
Notes: The wiki was last edited in 2014 (v9.2 era). Performance has improved since then (v9.2 pushed the knee right), but the fundamental principle remains valid.

## Evidence 4
Claim: Optimal connection pool size formula: `connections = ((core_count * 2) + effective_spindle_count)`. Core count excludes HT threads. Effective spindle count is zero when active data is fully cached.
Evidence: "A formula which has held up pretty well across a lot of benchmarks for years is that for optimal throughput the number of active connections should be somewhere near ((core_count * 2) + effective_spindle_count). Core count should not include HT threads, even if hyperthreading is enabled. Effective spindle count is zero if the active data set is fully cached, and approaches the actual number of spindles as the cache hit rate falls."
Source A: PostgreSQL Wiki — Number Of Database Connections
URL: https://wiki.postgresql.org/wiki/Number_Of_Database_Connections
Source B: HikariCP Wiki — About Pool Sizing
URL: https://github.com/brettwooldridge/HikariCP/wiki/About-Pool-Sizing
Confidence: HIGH — Cited by two independent authoritative sources (PostgreSQL project + HikariCP project)
Corroborated By: Both sources independently present the identical formula
Notes: HikariCP explicitly notes: "There hasn't been any analysis so far regarding how well the formula works with SSDs." The formula is presented as a starting point, not an absolute answer.

## Evidence 5
Claim: Reducing connection pool size from 2048 to 96 improved Oracle database response times by approximately 50x (~100ms → ~2ms).
Evidence: HikariCP wiki references an Oracle Real-World Performance group video demonstration showing this improvement. Quote: "reducing the connection pool size alone, in the absence of any other change, decreased the response times of the application from ~100ms to ~2ms -- over 50x improvement."
Source: HikariCP Wiki — About Pool Sizing
URL: https://github.com/brettwooldridge/HikariCP/wiki/About-Pool-Sizing
Confidence: MEDIUM — Single source citing a video demonstration. The video exists (YouTube link provided) but no independent academic or technical publication has replicated the specific numbers.
Corroborated By: PostgreSQL wiki independently confirms the principle that beyond saturation, more connections = worse performance
Notes: The video URL provided is https://www.youtube.com/watch?v=_C77sBcAtSQ

## Evidence 6
Claim: PostgreSQL performance "knee" occurs at approximately 50 connections in one documented benchmark.
Evidence: HikariCP wiki includes a PostgreSQL benchmark chart showing TPS rates flattening at around 50 connections.
Source: HikariCP Wiki — About Pool Sizing
URL: https://github.com/brettwooldridge/HikariCP/wiki/About-Pool-Sizing
Confidence: MEDIUM — Chart is referenced but the benchmark specifics (hardware, workload, PostgreSQL version) are not documented in detail
Corroborated By: The formula `((4*2)+1) = 9` for a 4-core i7 would predict the knee should be much lower, suggesting the 50-connection knee is for a higher-core machine
Notes: The chart is provided as illustration; exact benchmark conditions are not available

## Evidence 7
Claim: SSDs perform better with FEWER threads/connections than HDDs because faster, no-seek operations mean less blocking and therefore less benefit from additional threads.
Evidence: "Don't be tricked into thinking, 'SSDs are faster and therefore I can have more threads'. That is exactly 180 degrees backwards. Faster, no seeks, no rotational delays means less blocking and therefore fewer threads [closer to core count] will perform better than more threads."
Source: HikariCP Wiki — About Pool Sizing
URL: https://github.com/brettwooldridge/HikariCP/wiki/About-Pool-Sizing
Confidence: MEDIUM — Sound theoretical argument from a credible source, but no specific benchmark data provided for SSDs
Corroborated By: PostgreSQL wiki notes "There hasn't been any analysis so far regarding how well the formula works with SSDs"
Notes: This is a theoretical analysis; empirical SSD-specific benchmarks are noted as lacking

## Evidence 8
Claim: HikariCP `leakDetectionThreshold` defaults to 0 (disabled) and has a minimum enabled value of 2000ms. When enabled, it logs a message with stack trace when a connection is held out of the pool beyond the threshold.
Evidence: "This property controls the amount of time that a connection can be out of the pool before a message is logged indicating a possible connection leak. A value of 0 means leak detection is disabled. Lowest acceptable value for enabling leak detection is 2000 (2 seconds). Default: 0"
Source: HikariCP README — Configuration
URL: https://github.com/brettwooldridge/HikariCP#configuration-knobs-baby
Confidence: HIGH — Official project documentation
Corroborated By: No other connection pool was checked for equivalent features
Notes: HikariCP also supports `logUnclosedConnections` via PostgreSQL JDBC driver for capturing stack traces at connection open time

## Evidence 9
Claim: HikariCP default pool configuration: `maximumPoolSize=10`, `connectionTimeout=30000ms`, `idleTimeout=600000ms`, `maxLifetime=1800000ms`, `keepaliveTime=120000ms`, `validationTimeout=5000ms`.
Evidence: Direct reading from HikariCP README configuration section.
Source: HikariCP README — Configuration
URL: https://github.com/brettwooldridge/HikariCP#configuration-knobs-baby
Confidence: HIGH — Official project documentation, directly observed
Corroborated By: N/A — Primary source
Notes: `maxLifetime` should be "several seconds shorter than any database or infrastructure imposed connection time limit."

## Evidence 10
Claim: PgBouncer supports three pooling modes: session (most polite), transaction (releases after transaction), and statement (most aggressive, disallows multi-statement transactions). Default pool mode is session.
Evidence: "Session: Server is released back to pool after client disconnects. Transaction: Server is released back to pool after transaction finishes. Statement: Server is released back to pool after query finishes. Transactions spanning multiple statements are disallowed in this mode."
Source: PgBouncer Documentation — pgbouncer.ini
URL: https://www.pgbouncer.org/config.html
Confidence: HIGH — Official documentation
Corroborated By: PgBouncer Features page provides feature compatibility matrix for each mode
Notes: Transaction mode breaks session-based features (SET/RESET, LISTEN, WITH HOLD CURSOR, PREPARE/DEALLOCATE, session-level advisory locks, LOAD, temp tables with PRESERVE/DELETE ROWS).

## Evidence 11
Claim: PgBouncer defaults: `default_pool_size=20` (max server connections per user/database pair), `max_client_conn=100` (max client connections), memory usage approximately 2kB per connection.
Evidence: From PgBouncer config: "default_pool_size: The maximum number of server connections to allow per user/database pair. Default: 20" and "max_client_conn: Maximum number of client connections allowed. Default: 100". Features page: "Low memory requirements (2 kB per connection by default)."
Source: PgBouncer Documentation — pgbouncer.ini and Features page
URL: https://www.pgbouncer.org/config.html and https://www.pgbouncer.org/features.html
Confidence: HIGH — Official documentation
Corroborated By: N/A — Primary source
Notes: The theoretical maximum file descriptors used is `max_client_conn + (max pool_size * total databases * total users)`.

## Evidence 12
Claim: Azure Database for PostgreSQL flexible server reserves 15 connections for replication and monitoring. Default max_connections varies by tier: 50 for B1ms (1 vCore, 2 GiB) up to 5000 for 12+ vCore instances.
Evidence: "An Azure Database for PostgreSQL flexible server reserves 15 connections for physical replication and monitoring." Table shows: B1ms=50, B2s=429, B2ms=859, B4ms=1718, B8ms=3437, B12ms+=5000.
Source: Azure Database for PostgreSQL Limits
URL: https://learn.microsoft.com/en-us/azure/postgresql/flexible-server/concepts-limits
Confidence: HIGH — Official Microsoft documentation
Corroborated By: Azure explicitly recommends PgBouncer in transaction mode over increasing max_connections
Notes: Azure warns against increasing max_connections: "Although it's possible to increase the value of max_connections beyond the default setting, we advise against it."

## Evidence 13
Claim: Azure recommends using built-in PgBouncer in transaction mode with pool size of 2-5× vCores rather than increasing max_connections.
Evidence: "If you need more connections, we suggest that you instead use PgBouncer, the built-in Azure solution for connection pool management. Use it in transaction mode. To start, we recommend that you use conservative values by multiplying the vCores within the range of 2 to 5."
Source: Azure Database for PostgreSQL Limits
URL: https://learn.microsoft.com/en-us/azure/postgresql/flexible-server/concepts-limits
Confidence: HIGH — Official Microsoft documentation
Corroborated By: Independent recommendation aligns with PostgreSQL wiki's suggestion to limit active connections
Notes: This is a specific quantitative recommendation from a cloud provider managing thousands of production databases.

## Evidence 14
Claim: AWS RDS for PostgreSQL calculates max_connections as `LEAST(DBInstanceClassMemory/9531392, 5000)`, which is approximately memory_in_MB / 9.
Evidence: Subagent research found this formula documented in AWS RDS limits page.
Source: AWS RDS User Guide — Connection Limits
URL: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/CHAP_Limits.html
Confidence: MEDIUM — Found via subagent search; URL was not directly fetched and verified in full
Corroborated By: Azure uses similar memory-based calculations for its tier limits
Notes: AWS recommends using RDS Proxy for connection pooling. The formula shows that AWS dynamically calculates max_connections based on instance memory.

## Evidence 15
Claim: Google Cloud SQL recommends using connection pooling libraries (HikariCP for Java, SQLAlchemy for Python, Knex for Node.js, etc.) for applications connecting to Cloud SQL PostgreSQL.
Evidence: Google Cloud documentation provides code samples for TCP and Unix socket connection pools using HikariCP, SQLAlchemy, Knex, and other pooling libraries.
Source: Google Cloud SQL — Manage database connections
URL: https://cloud.google.com/sql/docs/postgres/manage-connections
Confidence: HIGH — Official Google Cloud documentation
Corroborated By: All code samples use well-established pooling libraries (HikariCP, SQLAlchemy)
Notes: Google also offers "Managed Connection Pooling" as a Cloud SQL feature, providing an additional pooling layer.

## Evidence 16
Claim: A 4-core i7 server with one hard disk should have a pool size of approximately 9-10 connections and can handle ~3000 front-end users at ~6000 TPS.
Evidence: "Your little 4-Core i7 server with one hard disk should be running a connection pool of: 9 = ((4 * 2) + 1). Call it 10 as a nice round number. Seem low? Give it a try, we'd wager that you could easily handle 3000 front-end users running simple queries at 6000 TPS on such a setup."
Source: HikariCP Wiki — About Pool Sizing
URL: https://github.com/brettwooldridge/HikariCP/wiki/About-Pool-Sizing
Confidence: MEDIUM — Illustrative example from authoritative source; not a benchmark result
Corroborated By: Consistent with the formula but specific numbers are illustrative, not measured
Notes: The claim is hedged ("we'd wager"), suggesting it's based on experience rather than published benchmarks.

## Evidence 17
Claim: Pool-locking deadlock avoidance formula: `pool_size = Tn × (Cm - 1) + 1`, where Tn is max threads and Cm is max simultaneous connections per thread.
Evidence: "The calculation of pool size in order to avoid deadlock is a fairly simple resource allocation formula: pool size = Tn x (Cm - 1) + 1"
Source: HikariCP Wiki — About Pool Sizing
URL: https://github.com/brettwooldridge/HikariCP/wiki/About-Pool-Sizing
Confidence: HIGH — Official project documentation; mathematically provable resource allocation formula
Corroborated By: N/A — Standard resource allocation mathematics
Notes: This is the MINIMUM pool size to avoid deadlock, not the optimal pool size.

## Evidence 18
Claim: HikariCP provides monitoring metrics via Dropwizard/Codahale integration: Wait Timer (connection acquisition wait time), Usage Histogram (connection hold time), TotalConnections, IdleConnections, ActiveConnections, PendingConnections.
Evidence: Metrics available via `metricRegistry` and `registerMbeans=true` for JMX. Specific metrics: `<pool>.pool.Wait`, `<pool>.pool.Usage`, TotalConnections, IdleConnections, ActiveConnections, PendingConnections cached gauges at 1s refresh.
Source: HikariCP Wiki — Dropwizard Metrics
URL: https://github.com/brettwooldridge/HikariCP/wiki/Dropwizard-Metrics
Confidence: HIGH — Official project documentation
Corroborated By: N/A — Primary source
Notes: HikariCP also supports MBean registration for JMX monitoring.

## Evidence 19
Claim: Persistent connections (e.g., Apache mod_php) are NOT connection pooling and still require a connection pool.
Evidence: "Persistent connections (as maintained by Apache's mod_php, for example) are *not* pooling and still require a connection pool."
Source: PostgreSQL Wiki — Number Of Database Connections
URL: https://wiki.postgresql.org/wiki/Number_Of_Database_Connections
Confidence: HIGH — Official PostgreSQL wiki
Corroborated By: Commonly known distinction in database engineering
Notes: This distinction is important for understanding that "keep-alive" is not the same as pooling.

## Evidence 20
Claim: When connections exceed the limit in Azure Database for PostgreSQL, the error is `FATAL: sorry, too many clients already.`
Evidence: "When connections exceed the limit, you might receive the following error: FATAL: sorry, too many clients already."
Source: Azure Database for PostgreSQL Limits
URL: https://learn.microsoft.com/en-us/azure/postgresql/flexible-server/concepts-limits
Confidence: HIGH — Official Microsoft documentation
Corroborated By: This is a standard PostgreSQL error message, documented across all PostgreSQL-based systems
Notes: This is the canonical error for connection exhaustion in PostgreSQL.

## Evidence 21
Claim: Idle and active connections both consume significant database resources. Excessive short-duration connections (< 60 seconds) increase CPU utilization due to connection/disconnection processing overhead.
Evidence: "Each connection in an Azure Database for PostgreSQL flexible server, regardless of whether it's idle or active, consumes a significant amount of resources from your database." And: "significant strain on resources... high CPU utilization, especially when many connections are established simultaneously and when connections have short durations (less than 60 seconds)"
Source: Azure Database for PostgreSQL Limits
URL: https://learn.microsoft.com/en-us/azure/postgresql/flexible-server/concepts-limits
Confidence: HIGH — Official Microsoft documentation
Corroborated By: PostgreSQL wiki confirms RAM usage scales with connection count via work_mem
Notes: This directly contradicts the common assumption that idle connections are "free."

## Evidence 22
Claim: HikariCP does NOT implement PreparedStatement caching at the pool level because it is an anti-pattern; database JDBC drivers handle statement caching more efficiently, sharing execution plans across connections.
Evidence: "At the connection pool layer PreparedStatements can only be cached per connection. If your application has 250 commonly executed queries and a pool of 20 connections you are asking your database to hold on to 5000 query execution plans... Using a statement cache at the pooling layer is an anti-pattern"
Source: HikariCP README — Configuration
URL: https://github.com/brettwooldridge/HikariCP#configuration-knobs-baby
Confidence: HIGH — Official project documentation
Corroborated By: N/A — Design decision documented by HikariCP maintainer
Notes: This is a design philosophy choice, not a widely debated technical claim.
