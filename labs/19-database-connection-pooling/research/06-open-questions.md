# Open Questions & Next Research Directions

## Unanswered Questions
1. **Dynamic / Elastic Pool Sizing vs Fixed Size Pools**:
   - What are the concrete latency trade-offs between dynamic pooling (min-idle scaling up to max-pool-size on demand) versus fixed-size pools (pre-warmed static connections) under sudden flash traffic spikes?
2. **Prepared Statements with Transaction Pooling**:
   - In modern PgBouncer versions (v1.21+), what are the exact performance and memory overhead trade-offs when enabling `max_prepared_statements` in transaction pooling mode across disparate application schemas?

## Weak Evidence / Claims Needing Deeper Research
1. **Exact SSD Sizing Constants**:
   - While the axiom that faster SSDs mean fewer threads is conceptually verified, empirical benchmark equations replacing `effective_spindle_count` specifically for PCIe Gen4/5 NVMe workloads across PostgreSQL 16+ remain sparsely quantified in academic literature.
2. **Impact of Thread-based PostgreSQL Backends**:
   - Ongoing community experiments exploring a multi-threaded architecture for PostgreSQL (replacing process-forking) may alter per-connection memory overhead in future major releases (PostgreSQL 19+).

## Possible Next Research Directions
- Benchmark comparison between PgBouncer (transaction mode) and native application pooling (HikariCP) under simulated microservice autoscaling (0 to 100 pods).
- Investigation of connection pooling strategies in serverless environments (AWS Lambda with RDS Proxy vs Neon serverless connection multiplexers).
- Formalization of automated connection leak detectors using eBPF or runtime JDBC/pgx interceptors.
