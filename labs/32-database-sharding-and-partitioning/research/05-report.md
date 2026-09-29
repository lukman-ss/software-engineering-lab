# Research Report

## Research Question
How to implement database sharding and table partitioning strategies for high-scale systems, addressing sharding key selection, routing algorithms, cross-shard query handling, and global unique ID generation?

## Executive Summary
This research distinguishes logical table partitioning (single instance) from physical sharding (multi-instance), establishes criteria for optimal sharding key selection (high-cardinality, low-frequency, non-monotonic fields), details routing mechanisms (consistent hashing vs hash modulo), exposes sharding complexities (cross-shard JOINs, distributed IDs, resharding), and recommends secondary indexing for non-sharding-key lookups. Evidence is drawn from PostgreSQL, MongoDB, Vitess, RFC 9562, Karger et al. (1997), and corroborated via Wikipedia. Hashed sharding is the documented mitigation for monotonic-key hotspots in range-partitioned systems.

## Findings

### Finding 1: Partitioning vs Sharding Architecture
**Claim**: Partitioning divides a table into smaller physical pieces within a single database instance; sharding distributes data across multiple independent servers to scale beyond single-machine limits.

**Evidence**:
- PostgreSQL 18: "Partitioning refers to splitting what is logically one large table into smaller physical pieces." Benefits: partition pruning, fast `DROP TABLE`/`ALTER TABLE DETACH PARTITION` avoiding VACUUM. Built-in forms: Range, List, Hash.
- MongoDB: "Sharding is a method for distributing data across multiple machines." Horizontal scaling "involves dividing the system dataset and load over multiple servers… Each machine handles a subset of the overall workload." A sharded cluster consists of shards (replica sets), mongos routers, and config servers.
- Vitess: "Sharding is a method of horizontally partitioning a database to store data across two or more database servers." A user keyspace split into two shards stores each user's data in only one shard.

**Sources**: PostgreSQL Documentation (5.12), MongoDB Sharding Overview, Vitess Sharding Overview
**Confidence**: HIGH

### Finding 2: Sharding Key Selection is Critical for Write Distribution
**Claim**: Monotonically increasing range-based shard keys (e.g., `created_at`, autoincrement) concentrate new writes on one active chunk/shard; optimal keys have high cardinality, low frequency, and enable point-lookup routing.

**Evidence**:
- MongoDB: "A shard key on a value that increases or decreases monotonically is more likely to distribute inserts to a single chunk." All new inserts route to the `maxKey` chunk; "the shard containing that chunk becomes the bottleneck for write operations."
- MongoDB: "Where possible, choose a shard key with high cardinality. A shard key with low cardinality reduces the effectiveness of horizontal scaling." Cardinality of 7 (`continent`) caps effective shard count at 7. "The ideal shard key distributes data evenly across the sharded cluster while also facilitating common query patterns."
- Vitess: "The Primary Vindex must be unique: given an input value, it must produce a single keyspace ID. At the time of an insert to the table, the unique mapping produced by the Primary Vindex determines the target shard for the inserted row."

**Sources**: MongoDB Shard Key Selection, MongoDB Hashed Sharding, Vitess Vindexes
**Confidence**: HIGH

### Finding 3: Consistent Hashing Minimizes Data Movement During Resizing
**Claim**: Hash Modulo (key % M) forces most keys to remap when M changes; Consistent Hashing redistributes only ~n/m (≈ 1/n when adding one node) of keys on average.

**Evidence**:
- Karger et al. (1997, STOC '97, DOI 10.1145/258533.258660): "when a hash table is resized, only n/m keys need to be remapped on average where n is the number of keys and m is the number of slots." Traditional hash tables "cause nearly all keys to be remapped because the mapping between the keys and the slots is defined by a modular operation."
- Wikipedia (corroborating restatement): "the addition of n-th server only causes 1/n fraction of the BLOBs to relocate." Virtual nodes (multiple labels per physical node) reduce distribution skew ("skewness … are called virtual nodes").
- MongoDB: "The efficiency of ranged sharding depends on the shard key chosen. Poorly considered shard keys can result in uneven distribution of data, which can negate some benefits of sharding or can cause performance bottlenecks."
- Vitess: Key ranges define responsibility; resharding "reorganizes data to match the new scheme."

**Sources**: ACM STOC '97 (Karger et al.), Wikipedia, MongoDB Sharding, Vitess Resharding
**Confidence**: HIGH

### Finding 4: Core Complexities Introduced by Sharding
**Claim**: Sharding introduces cross-shard JOIN limitations, loss of native auto_increment, and operational resharding/rebalancing requirements; TwoPC trades commit latency for atomicity without full ACID isolation across shards.

**Evidence**:
- MongoDB: "If queries do not include the shard key or the prefix of a compound shard key, mongos performs a broadcast operation, querying all shards… These scatter/gather queries can be long running operations." "Queries that involve multiple shards for each request are less efficient and do not scale linearly when more shards are added."
- RFC 9562 §2.1: "'auto-increment' schemes that are often used by databases do not work well: the effort required to coordinate sequential numeric identifiers across a network can easily become a burden." UUIDv7 addresses this with time-ordered B-tree-locality values.
- Vitess Sequences: "Very high throughput for ID creation, using a configurable in-memory block allocation." Backing table must be in an unsharded keyspace; `cache` value of ~1000 is recommended, with the trade-off that a tablet crash loses up to that many IDs.
- Vitess Distributed Transactions: "Using atomic distributed transactions will impact commit latency." "TwoPC guarantees atomicity but does not provide isolation in the traditional ACID sense. Applications might observe fractured reads in a cross-shard query."
- MongoDB: "when a transaction writes to multiple shards, not all outside read operations need to wait for the result of the committed transaction to be visible across the shards."

**Sources**: MongoDB Sharding Overview, Vitess Sequences, Vitess Distributed Transactions, RFC 9562 §2.1
**Confidence**: HIGH

### Finding 5: Handling Queries Without the Sharding Key
**Claim**: Queries missing the sharding key require scatter-gather (broadcast) unless supplemented by secondary indexes (Lookup Vindexes, Global Secondary Indexes); Lookup Vindexes add per-write overhead to enable targeted point-lookups.

**Evidence**:
- MongoDB: "In a sharded cluster, the mongos routes queries to only the shards that contain the relevant data if the queries contain the shard key. When the queries do not contain the shard key, the queries are broadcast to all shards for evaluation." Exception noted: large aggregations benefit from parallel scatter-gather.
- Vitess: "Secondary Vindexes … offer optimizations for WHERE clauses that do not use the Primary Vindex. In the absence of a Secondary Vindex, VTGate would have to scatter the query to all shards." "Lookup Unique Owned: Used for optimizing high QPS read queries that do not use the Primary Vindex columns in their WHERE clause. There is a price to pay: an extra write to the lookup table for insert and delete operations, and an extra lookup for read operations. The overhead of maintaining the lookup table is amortized as the number of shards grow."
- Vitess `consistent_lookup` type guarantees consistency "without paying the price of 2PC" via careful locking and transaction sequencing.

**Sources**: MongoDB Sharding (scatter-gather), Vitess Vindexes (Lookup Vindex types and costs)
**Confidence**: HIGH

### Finding 6: Hashed Sharding Mitigates Monotonic-Key Hotspots
**Claim**: MongoDB's Hashed Sharding is the documented alternative when the natural shard key is monotonically increasing (timestamps, ObjectIds, autoincrement); it distributes inserts evenly at the cost of range-query efficiency.

**Evidence**:
- MongoDB Hashed Sharding: "Hashed sharding provides a more even data distribution across the sharded cluster at the cost of reducing [Targeted Operations]… Post-hash, documents with 'close' shard key values are unlikely to be on the same chunk or shard - mongos is more likely to perform Broadcast Operations to fulfill a given ranged query." "mongos can target queries with equality matches to a single shard."
- "The field you choose as your hashed shard key should have a good cardinality, or large number of different values. Hashed keys are ideal for shard keys with fields that change monotonically like ObjectId values or timestamps."
- "Given a collection using a monotonically increasing value X as the shard key, using ranged sharding results in [concentrated inserts on the MaxKey chunk]… By using a hashed index on X, the distribution of inserts is [even across the cluster]."

**Sources**: MongoDB Hashed Sharding, MongoDB Choose a Shard Key
**Confidence**: HIGH

### Finding 7: Resharding Reorganizes Data with Brief Downtime
**Claim**: Resharding copies and verifies data on new shards while existing shards serve live traffic; cutover requires only seconds of read-only downtime. Balancing in MongoDB is automatic/continuous; in Vitess it is operator-triggered.

**Evidence**:
- Vitess: "During resharding, Vitess copies, verifies, and keeps data up-to-date on new shards while the existing shards continue to serve live read and write traffic. When you're ready to switch over, the migration occurs with only a few seconds of read-only downtime."
- MongoDB: "a balancer runs in the background to migrate ranges across the shards." "A single shard can only participate in one chunk migration at a time. When MongoDB succeeds in copying a range of data from one shard to another, the range on the donor shard is marked for removal by the range deleter. This process is slow and resource intensive." MongoDB 8.0: `sh.shardAndDistributeCollection()` avoids waiting on the balancer.

**Sources**: Vitess Sharding (Resharding section), MongoDB Sharding Overview
**Confidence**: HIGH

## Areas of Agreement
All Tier 1 sources agree on:
- The architectural distinction between partitioning (single-instance) and sharding (multi-instance).
- The perils of monotonic keys as range-based shard keys; hashed sharding as the documented mitigation.
- The benefit of high-cardinality, low-frequency sharding keys (user_id, tenant_id).
- The superior resizing properties of consistent hashing vs hash modulo.
- The necessity of distributed ID generation (Vitess Sequences, UUIDv7, Snowflake) in sharded environments.
- The need for broadcast/scatter-gather queries or secondary indexes when the shard key is absent.
- The added complexity and latency cost of cross-shard distributed transactions (TwoPC).

## Areas of Disagreement
No material contradictions. Minor implementation-variance differences:
- **Rebalancing trigger**: MongoDB's balancer is automatic; Vitess resharding is operator-triggered. Not a contradiction — both acknowledge resharding is required when the sharding scheme changes.
- **Cross-shard isolation guarantees**: Vitess TwoPC explicitly does not provide full ACID isolation; MongoDB's distributed transactions may also expose partial-commit visibility. Both document this as an accepted trade-off.
- **Scatter-gather acceptability**: MongoDB explicitly notes large aggregations benefit from scatter-gather parallelism; simple reads do not. This is a workload-dependent nuance, not a contradiction.

## Limitations
- Research focused on Tier 1 sources (official documentation, standards, one academic paper); community benchmarks and vendor-managed services (AWS Aurora Limitless, Google Cloud Spanner, Azure Cosmos DB, CockroachDB) were not evaluated in depth.
- No quantitative latency/throughput measurements under load; claims are architectural properties documented by the vendors, not empirical benchmarks.
- PG 19 beta behavior (released after PG 18 current) may introduce partitioning changes not covered.
- Consistent hashing O(K/N + log N) claim in the complexity table is from Wikipedia; the original Karger paper's exact theorem statement could not be directly verified from the ACM HTML landing page.

## Conclusion
For high-scale transaction systems (e.g., e-commerce with millions of writes per day), the evidence supports:
1. Use a high-cardinality, non-monotonic sharding key (user_id or tenant_id) to distribute writes evenly and enable targeted point-lookup routing.
2. When the natural key is monotonic (timestamps, autoincrement), prefer hashed sharding over ranged sharding to prevent write hotspots.
3. Apply consistent hashing (or Vitess-style key ranges) to minimize data movement during cluster scaling.
4. Add secondary Lookup Vindexes or Global Secondary Indexes to avoid scatter-gather for queries lacking the shard key; budget for per-write overhead.
5. Replace native auto_increment with distributed ID generation (UUIDv7 for time-ordered B-tree locality; Vitess Sequences for numeric IDs with block-allocation caching).
6. Plan for cross-shard transaction limitations: co-locate related rows on the same shard key where possible; use TwoPC only where atomicity is required, accepting its latency and isolation trade-offs.
7. Single-instance partitioning remains valuable within each shard for data lifecycle management (time-based archival, partition pruning) but does not replace horizontal scaling via sharding.
