# Open Questions

## Unanswered Questions

1. **What is the empirically validated pool sizing formula for SSDs?**
   - Status: Unresearched by authoritative sources. The PostgreSQL wiki and HikariCP both acknowledge: "There hasn't been any analysis so far regarding how well the formula works with SSDs." While HikariCP theorizes that fewer connections may be better for SSDs due to less blocking, no experimental data exists to confirm an optimal formula for NVMe or SAN-based SSD storage.

2. **How does the optimal pool size formula change for different workload types?**
   - Status: Incomplete evidence. Sources provide the base formula but do not break down by workload characteristics:
     - Read-heavy vs write-heavy workloads
     - Transaction-intensive vs analytical workloads
     - Short OLTP transactions vs long-running reporting queries
     - Cache hit rate impact on `effective_spindle_count` term

3. **What are equivalent leak detection mechanisms in connection pools other than HikariCP?**
   - Status: Partially researched. HikariCP's `leakDetectionThreshold` and `logUnclosedConnections` (via PostgreSQL JDBC) are well-documented. Equivalent mechanisms in:
     - Apache DBCP2
     - Tomcat JDBC Pool
     - c3p0
     - Vibur DBCP
     - are not explicitly detailed in the evidence gathered.

4. **What is the performance impact of different PgBouncer pool modes on various PostgreSQL workloads?**
   - Status: Qualitative only. Sources document which features break in each mode, but lack:
     - Benchmark data comparing throughput/latency between session, transaction, and statement modes
     - Workload-specific recommendations (e.g., "use transaction mode for REST API services")
     - Quantitative measurement of CPU/memory savings from more aggressive pooling modes

## Weak Evidence

1. **Oracle Real-World Performance 50x Improvement Claim**
   - Evidence: Single video demonstration cited by HikariCP wiki
   - Weakness: No independent replication, no details on hardware, PostgreSQL version, or exact workload
   - Strength: Principle (fewer connections = less contention) is independently verified by PostgreSQL wiki
   - Next Step: Seek independent benchmark studies or reproduce in controlled environment

2. **AWS RDS max_connections Formula**
   - Evidence: Found via subagent research (`LEAST(DBInstanceClassMemory/9531392, 5000)`)
   - Weakness: URL not directly fetched/verified in this research session; may require checking current AWS documentation
   - Strength: Similar memory-based calculation aligns with Azure's tiered approach
   - Next Step: Direct verification from current AWS RDS documentation

3. **HikariCP Default Performance Claims**
   - Evidence: HikariCP README states: "you could easily handle 3000 front-end users running simple queries at 6000 TPS on such a setup"
   - Weakness: Illustrative example ("we'd wager"), not benchmark result
   - Strength: Consistent with the principle that small pools can handle high concurrency
   - Next Step: Reproduce with benchmark tools like pgbench or HammerDB

## Claims Needing Deeper Research

1. **Connection Pool Sizing for Microservices Architectures**
   - Gap: Research focused on monolithic or traditional multi-instance deployments. Modern microservices architectures with dozens/hundreds of services each needing database access present unique challenges for global pool sizing.

2. **Interaction Between Connection Pooling and Prepared Statement Caching**
   - Gap: HikariCP explicitly states it does not cache PreparedStatements at the pool level (calling it an anti-pattern). Need research on:
     - Performance comparison: pool-level vs driver-level statement caching
     - Memory usage trade-offs
     - Workloads where statement caching provides significant benefit

3. **Effect of Connection Pooling on Database Connection Lifetime Metrics**
   - Gap: Sources mention `maxLifetime` and `server_lifetime` but lack:
     - Empirical data on optimal lifetime settings for different cloud providers
     - Analysis of connection churn vs memory fragmentation trade-offs
     - Impact of TLS session resumption on connection lifetime decisions

4. **Production Thresholds for Connection Pool Monitoring Metrics**
   - Gap: Sources list metrics to monitor (wait time, usage, active/idle) but lack:
     - Evidence-based alert thresholds (e.g., "alert if 95th percentile wait time > 100ms")
     - Correlation between specific metric values and user-visible latency
     - Baselines for different application types (OLTP, batch, reporting)

## Possible Next Research Directions

1. **SSD-Specific Connection Pool Benchmark Study**
   - Design: Systematic benchmark varying pool size on identical hardware with HDD vs SSD storage
   - Metrics: TPS, latency (p50, p95, p99), connection wait time, CPU utilization
   - Variables: Cache hit rate, query complexity, read/write ratio
   - Output: Empirical formula or adjustment factor for SSD storage

2. **Connection Leak Detection Across Pooling Libraries**
   - Design: Comparative analysis of leak detection mechanisms in HikariCP, DBCP2, Tomcat JDBC, c3p0
   - Metrics: Detection latency, false positive/negative rates, overhead, configurability
     - Test scenarios: timed-out connections, forgotten close() in finally blocks, exception paths
   - Output: Feature matrix and recommendations by use case

3. **Cloud Provider Connection Pooling Efficacy Study**
   - Design: Measure actual connection savings from:
     - AWS RDS Proxy vs direct pooling
     - Azure built-in PgBouncer vs application-level pooling
     - Google Cloud SQL Managed Connection Pooling vs self-managed
   - Metrics: Connection count reduction, CPU/memory savings, latency impact, failure rate during scaling events
   - Output: Quantified benefit of managed pooling services

4. **Dynamic Pool Sizing Algorithms for Variable Workloads**
   - Design: Evaluate algorithms that adjust pool size based on:
     - Real-time wait time metrics
     - Time-of-day patterns
     - Queue depth observation
     - Connection acquisition success rate
   - Metrics: Stability, responsiveness to spikes, resource efficiency, oscillation tendency
   - Output: Recommended algorithms and tuning parameters for different workload patterns

5. **PostgreSQL Connection Memory Profiling**
   - Design: Measure actual memory consumption per connection at different:
     - `work_mem` settings
     - `maintenance_work_mem` settings
     - `temp_file_limit` settings
     - With/without prepared statements, cursors, large result sets
   - Output: Accurate memory-per-connection formula for capacity planning

## Research Gaps Identified During Investigation

1. **Temporal Aspects of Pool Exhaustion**
   - Little discussion on how quickly pools exhaust under leak conditions
   - No guidance on safe leak rates (e.g., "X leaks/hour is acceptable for Y pool size")

2. **Interaction Between Pooling and Statement Timeout/Idle Timeout**
   - How do `statement_timeout`, `idle_in_transaction_session_timeout`, and pool timeouts interact?
   - Which should be shorter/longer for optimal resource reclamation?

3. **Prepared Statement Caching in External Poolers**
   - Does PgBouncer's prepared statement caching (`max_prepared_statements`) interact beneficially or adversely with driver-level caching?

4. **Connection Pooling in Serverless/FaaS Environments**
   - How do connection pools work in ephemeral compute environments (AWS Lambda, Azure Functions)?
   - What are the trade-offs between pooling and connection overhead in short-lived contexts?

5. **Observability of Connection State in Distributed Tracing**
   - How to correlate connection pool metrics with distributed traces (OpenTelemetry, Jaeger)?
   - Can connection wait time be attributed to specific business transactions?