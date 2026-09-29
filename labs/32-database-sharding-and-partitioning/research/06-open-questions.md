# Open Questions

## Unanswered Questions
- What is the quantitative threshold (row count, data size, QPS) at which partitioning within a single instance becomes insufficient and physical sharding becomes necessary? No authoritative cross-vendor benchmark establishes a universal cutoff.
- How to calculate the optimal number of shards and chunk/partition sizes for a given workload profile (read-heavy vs write-heavy, OLTP vs OLAP) in e-commerce transaction systems?
- Under what conditions does the overhead of maintaining Lookup Vindexes (extra writes, storage) outweigh the benefit of avoiding scatter-gather queries for non-sharding-key lookups?
- What is the precise latency and throughput penalty of Vitess TwoPC distributed transactions vs. eventual consistency or Saga patterns at scale (millions of cross-shard writes per day)?
- How do clock-skew and leap-second handling affect UUIDv7 monotonicity guarantees in distributed ID generation compared to Snowflake-like centralized timestamp+sequence IDs?
- For true consistent hashing across organizations, does the original Karger et al. 1997 result (1/n keys move when one node joins) hold in practice, or do typical implementations (using multiple virtual nodes per physical node) cause variance?
- Given both range and hashed sharding are documented for monotonic keys, does a hybrid (e.g., use range for historical data, hash for new data) provide better write distribution for high-write systems?
- In Vitess, how does the "block allocation" (cache parameter) in Vitess Sequences affect the collision probability and the gap size across a cluster when tablets crash?

## Weak Evidence
- **Hash Modulo vs Consistent Hashing in Production Databases**: Evidence relies on foundational CS literature (Karger et al. 1997) and general hashing theory; direct database-engine documentation on hash modulo implementation details for sharding is limited.
- **UUIDv7 Adoption**: RFC 9562 is recent (May 2024); real-world production adoption metrics and B-tree index performance comparisons with UUIDv4/Snowflake are still emerging.
- **Cross-Shard JOIN Performance**: Vitess Gen4 planner capabilities described qualitatively; no quantitative latency measurements provided for distributed JOIN execution vs. application-level merge.
- **MongoDB Hashed Sharding Consistency**: Hashed sharding feature documented but internal rebalancing and LRU behavior relative to range sharding under load not published.

## Claims Needing Deeper Research
- **Vitess Sequences with block allocation (cache=1000)**: Claims "Very high throughput for ID creation, using a configurable in-memory block allocation" — requires empirical QPS validation under concurrent load.
- **PostgreSQL rule of thumb that "size of the table should exceed the physical memory of the database server"**: Needs workload-specific validation for partitioning benefit threshold.
- **MongoDB's `analyzeShardKey` (v7.0+) effectiveness**: Requires case studies on accuracy vs. manual analysis.
- **Virtual Node Tuning**: Quantitative guidelines for virtual node density (e.g., 100-300 vnodes per physical node) under different deployment scenarios.
- **Real-world Impact of Scatter-Gather**: Performance impact of scatter-gather queries vs. Lookup Vindex overhead across different cluster sizes and node failure rates.

## Possible Next Research Directions
- Empirical comparison of Vitess vs. Citus vs. MongoDB sharded clusters under identical e-commerce workload simulations (millions of writes/day).
- Deep investigation of Google Spanner / CockroachDB true distributed SQL as alternative to manual sharding — external consistency and TrueTime vs. sharding.
- Simulation of shard key cardinality/frequency impact on chunk distribution using synthetic transaction datasets.
- Study of geo-sharding strategies for data residency compliance and latency optimization across regions.
- Analysis of automatic resharding and autoscaling mechanisms in managed cloud databases (Aurora Limitless, Atlas, PlanetScale).
- Field validation of the Karger 1/n remap claim in production clusters (observing key redistribution after node joins/leaves).
- Measurement of UUIDv7 index locality vs. UUIDv4/Snowflake in large-scale production workloads.
- Systematic evaluation of per-write overhead of Lookup Vindexes vs. scatter-gather performance for different query patterns.
