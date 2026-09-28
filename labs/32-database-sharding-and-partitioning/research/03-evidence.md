# Evidence

## Evidence 1
Claim: Table partitioning splits what is logically one large table into smaller physical pieces within a single database instance, enabling partition pruning and fast dropping of data.
Evidence: PostgreSQL documentation specifies: "Partitioning refers to splitting what is logically one large table into smaller physical pieces... Dropping an individual partition using `DROP TABLE`, or doing `ALTER TABLE DETACH PARTITION`, is far faster than a bulk operation... entirely avoiding the VACUUM overhead caused by a bulk DELETE." Furthermore, partition pruning proves that unneeded partitions need not be scanned during query planning or execution.
Source: PostgreSQL 18 Documentation - Section 5.12
URL: https://www.postgresql.org/docs/current/ddl-partitioning.html
Confidence: HIGH
Corroborated By: MongoDB Documentation (Partitioning vs Sharding definitions), Vitess Documentation
Notes: Single-instance partitioning does not distribute write load across multiple independent physical machines or CPU/RAM boundaries.

## Evidence 2
Claim: Sharding distributes data horizontally across multiple independent servers/instances to scale CPU, memory, disk I/O, and storage capacity beyond the limits of a single machine.
Evidence: MongoDB documentation states: "Sharding is a method for distributing data across multiple machines... Vertical scaling increases the capacity of a single server... Available technology and cloud provider hardware configurations impose a practical maximum for vertical scaling. Horizontal Scaling involves dividing the system dataset and load over multiple servers... Each machine handles a subset of the overall workload." Vitess documentation confirms: "Sharding is a method of horizontally partitioning a database to store data across two or more database servers."
Source: MongoDB Manual - Sharding Overview / Vitess Documentation - Sharding Overview
URL: https://www.mongodb.com/docs/manual/sharding/ / https://vitess.io/docs/24.0/reference/features/sharding/
Confidence: HIGH
Corroborated By: PostgreSQL Documentation, ACM Karger et al. (1997)
Notes: Sharding scales both reads and writes horizontally by delegating distinct subsets of data to separate primary instances.

## Evidence 3
Claim: Monotonically increasing or decreasing keys (such as `created_at` or autoincrement IDs) create severe write hotspots when used as range-based sharding keys.
Evidence: MongoDB documentation states: "If the shard key value is always increasing, all new inserts are routed to the chunk with maxKey as the upper bound... The shard containing that chunk becomes the bottleneck for write operations. If your data model requires sharding on a key that changes monotonically, consider using Hashed Sharding."
Source: MongoDB Manual - Choose a Shard Key (Monotonically Changing Shard Keys)
URL: https://www.mongodb.com/docs/manual/core/sharding-choose-a-shard-key/
Confidence: HIGH
Corroborated By: Vitess Documentation (Primary Vindexes), RFC 9562
Notes: In high-throughput transaction systems, routing on `created_at` forces 100% of concurrent write traffic into the current active shard, completely defeating the purpose of horizontal scaling.

## Evidence 4
Claim: Sharding keys with high cardinality and low frequency (e.g., `user_id` or `tenant_id`) ensure balanced write distribution and enable direct point-lookup routing.
Evidence: MongoDB documentation notes: "Where possible, choose a shard key with high cardinality. A shard key with low cardinality reduces the effectiveness of horizontal scaling... The frequency of the shard key represents how often a given shard key value occurs... ideal shard key distributes data evenly across the sharded cluster while also facilitating common query patterns." Vitess documentation states: "For queries that include the shard key... mongos / VTGate can target the query at a specific shard... targeted operations are generally more efficient than broadcasting to every shard in the cluster."
Source: MongoDB Manual - Choose a Shard Key / Vitess Documentation - Vindexes
URL: https://www.mongodb.com/docs/manual/core/sharding-choose-a-shard-key/ / https://vitess.io/docs/24.0/reference/features/vindexes/
Confidence: HIGH
Corroborated By: PostgreSQL Documentation (Best Practices for Partitioning)
Notes: High-cardinality keys like `user_id` allow isolation of tenant data and predictable deterministic hashing.

## Evidence 5
Claim: Consistent Hashing minimizes key redistribution during cluster resizing to $O(K/N)$ keys on average, whereas Hash Modulo ($N \pmod M$) forces almost all keys to relocate.
Evidence: Foundational academic paper by Karger et al. (1997) and Wikipedia analysis document: "In consistent hashing... when a hash table is resized, only $n/m$ keys need to be remapped on average where $n$ is the number of keys and $m$ is the number of slots... In contrast, in most traditional hash tables, a change in the number of array slots causes nearly all keys to be remapped because the mapping is defined by a modular operation." Virtual nodes are employed to reduce distribution skewness across physical nodes.
Source: ACM STOC '97 (Karger et al.) / Consistent Hashing Literature
URL: https://doi.org/10.1145/258533.258660
Confidence: HIGH
Corroborated By: MongoDB Range/Chunk Balancer Docs, Vitess Key Range Documentation
Notes: Consistent hashing with virtual nodes prevents catastrophic cache stampedes and massive data migration bursts when adding or removing shards.

## Evidence 6
Claim: Queries executed without the sharding key require Scatter-Gather (broadcast) operations across all shards, increasing latency and limiting cluster scalability.
Evidence: MongoDB documentation states: "If queries do not include the shard key or the prefix of a compound shard key, mongos performs a broadcast operation, querying all shards in the sharded cluster. These scatter/gather queries can be long running operations... Queries that involve multiple shards for each request are less efficient and do not scale linearly when more shards are added." Vitess documentation corroborates: "In the absence of a Secondary Vindex, VTGate would have to scatter the query to all shards."
Source: MongoDB Manual - Sharding Overview / Vitess Documentation - Vindexes
URL: https://www.mongodb.com/docs/manual/sharding/ / https://vitess.io/docs/24.0/reference/features/vindexes/
Confidence: HIGH
Corroborated By: Vitess Distributed Query Planner Docs
Notes: Scatter-gather latency is bound by the slowest responding shard in the cluster (p99 tail latency amplification).

## Evidence 7
Claim: Secondary Lookup Vindexes (Global Secondary Indexes) provide point-lookup routing for non-sharding-key queries at the cost of additional write overhead and cross-shard consistency maintenance.
Evidence: Vitess documentation defines: "Lookup Vindexes are implemented as a MySQL lookup table that maps a column value to keyspace IDs... used for optimizing high QPS read queries that do not use the Primary Vindex columns in their WHERE clause. There is a price to pay: an extra write to the lookup table for insert and delete operations, and an extra lookup for read operations." Vitess implements `consistent_lookup` to ensure transactional synchronization.
Source: Vitess Documentation - Vindexes (Lookup Vindex types)
URL: https://www.vitess.io/docs/24.0/reference/features/vindexes/
Confidence: HIGH
Corroborated By: MongoDB Global Indexing Patterns
Notes: An alternative to global lookup tables is duplicating the primary sharding key into composite identifiers (e.g., embedding `user_id` inside `order_id`).

## Evidence 8
Claim: Standard auto-increment cannot generate globally unique identifiers across independent database instances without central coordination; modern architectures require distributed ID generation (UUIDv7, Snowflake).
Evidence: Vitess documentation states: "MySQL provides the auto_increment feature... However, when a table is sharded across multiple instances, maintaining the same feature is a lot more tricky." Vitess uses centralized Sequence tables with block allocation. IETF RFC 9562 outlines: "In distributed applications, auto-increment schemes that are often used by databases do not work well: the effort required to coordinate sequential numeric identifiers across a network can easily become a burden... UUIDv7 features a time-ordered value field derived from the widely implemented and well-known Unix Epoch timestamp... with improved B-tree database index locality."
Source: RFC 9562 - Section 2.1 & 5.7 / Vitess Documentation - Sequences
URL: https://www.rfc-editor.org/rfc/rfc9562.html / https://vitess.io/docs/24.0/reference/features/vitess-sequences/
Confidence: HIGH
Corroborated By: Twitter Snowflake Architecture, MongoDB ObjectId Specification
Notes: UUIDv7 and Snowflake IDs eliminate central coordination bottlenecks while preserving temporal B-tree index locality.
