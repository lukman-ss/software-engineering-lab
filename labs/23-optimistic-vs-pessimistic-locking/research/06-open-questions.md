# Open Questions & Next Research Directions

## Unanswered Questions
1. **MySQL 8.0 Locking Semantics Verification**:
   - Official MySQL documentation was inaccessible during research (technical difficulties error). MySQL 8.0 `SELECT ... FOR UPDATE`, `LOCK IN SHARE MODE`, and InnoDB gap locking behaviors need direct vendor doc verification before inclusion in future lab implementations.
2. **CockroachDB/Spanner Implementation Details**:
   - Distributed database locking semantics (HLC timestamps, TrueTime, write-write conflict detection) were inferred from academic sources. Direct vendor doc URLs (cockroachdb.com/docs, cloud.google.com/spanner/docs) needed for production use.

## Weak Evidence / Claims Needing Deeper Research
1. **Exact Contention Crossover Point**:
   - While the axiom that optimistic is better at low contention is conceptually verified (Bernstein & Goodman, 1981), empirical crossover thresholds (15-20% from SQL Server benchmarks) are workload-specific. Need controlled benchmark with same hardware, transaction mix, and isolation level across vendors.
2. **PostgreSQL 19+ Locking Changes**:
   - Research based on PostgreSQL 18 docs. PostgreSQL 19 Beta 4 released (2026-09-24). Locking behavior changes in 19 not reviewed.

## Possible Next Research Directions
- Benchmark comparison between optimistic (version column + retry) vs pessimistic (`SELECT FOR UPDATE`) under simulated e-commerce checkout (high contention) and CMS read (low contention) workloads on PostgreSQL 18.
- Investigation of ORM locking fidelity: verify Hibernate 6.6 `@Version` generated SQL matches native locking semantics across PostgreSQL/MySQL/SQL Server dialects.
- Distributed locking patterns: ZooKeeper/etcd-based pessimistic locks vs HLC-based optimistic concurrency in microservice architectures.