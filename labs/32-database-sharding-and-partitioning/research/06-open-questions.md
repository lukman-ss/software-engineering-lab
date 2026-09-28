# Open Questions

## Unanswered Questions
- What is the quantitative threshold (row count, data size, QPS) at which partitioning within a single instance becomes insufficient and physical sharding becomes necessary? No authoritative cross-vendor benchmark establishes a universal cutoff.
- How to calculate the optimal number of shards and chunk/partition sizes for a given workload profile (read-heavy vs write-heavy, OLTP vs OLAP) in e-commerce transaction systems?
- Under what conditions does the overhead of maintaining Lookup Vindexes (extra writes, storage) outweigh the benefit of avoiding scatter-gather queries for non-sharding-key lookups?
- What is the precise latency and throughput penalty of Vitess TwoPC distributed transactions vs. eventual consistency or Saga patterns at scale (millions of cross-shard writes per day)?
- How do clock-skew and leap-second handling affect UUIDv7 monotonicity guarantees in distributed ID generation compared to Snowflake-like centralized timestamp+sequence IDs?

## Weak Evidence
- **Hash Modulo vs Consistent Hashing in Production Databases**: Evidence relies on foundational CS literature (Karger et al. 1997) and general hashing theory; direct database-engine documentation on hash modulo implementation details for sharding is limited.
- **UUIDv7 Adoption**: RFC 9562 is recent (May 2024); real-world production adoption metrics and B-tree index performance comparisons with UUIDv4/Snowflake are still emerging.
- **Cross-Shard JOIN Performance**: Vitess Gen4 planner capabilities described qualitatively; no quantitative latency measurements provided for distributed JOIN execution vs. application-level merge.

## Claims Needing Deeper Research
- Virtess Sequences with block allocation (cache=1000) claim "Very high throughput for ID creation, using a configurable in-memory block allocation" — requires empirical QPS validation under concurrent load.
- PostgreSQL claim that "size of the table should exceed the physical memory of the database server" as rule of thumb for partitioning benefit — needs workload-specific validation.
- MongoDB's `analyzeShardKey` (v7.0+) effectiveness in data-driven shard key selection — requires case studies on accuracy vs. manual analysis.

## Possible Next Research Directions
- Empirical comparison of Vitess vs. Citus vs. MongoDB sharded clusters under identical e-commerce workload simulations (millions of writes/day).
- Deep investigation of Google Spanner / CockroachDB true distributed SQL as alternative to manual sharding — external consistency and TrueTime vs. sharding.
- Simulation of shard key cardinality/frequency impact on chunk distribution using synthetic transaction datasets.
- Study of geo-sharding strategies for data residency compliance and latency optimization across regions.
- Analysis of automatic resharding and autoscaling mechanisms in managed cloud databases (Aurora Limitless, Atlas, PlanetScale).
