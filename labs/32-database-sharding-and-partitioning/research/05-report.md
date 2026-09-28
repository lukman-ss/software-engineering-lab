# Research Report

## Research Question
How to implement database sharding and table partitioning strategies for high-scale systems, addressing sharding key selection, routing algorithms, cross-shard query handling, and global unique ID generation?

## Executive Summary
This research distinguishes logical table partitioning (single instance) from physical sharding (multi-instance), establishes criteria for optimal sharding key selection (high-cardinality, non-monotonic fields), details routing mechanisms (consistent hashing vs hash modulo), exposes sharding complexities (cross-shard JOINs, distributed IDs, resharding), and recommends secondary indexing for non-sharding-key lookups. Evidence comes from authoritative sources: PostgreSQL, MongoDB, Vitess, RFC 9562, and the original consistent hashing paper.

## Findings

### Finding 1: Partitioning vs Sharding Architecture
**Claim**: Partitioning divides a table into smaller physical pieces within a single database instance, while sharding distributes data across multiple independent servers/instances to scale beyond single-machine limits.  
**Evidence**: 
- PostgreSQL: "Partitioning refers to splitting what is logically one large table into smaller physical pieces... within one engine database" and notes performance gains via partition pruning and fast bulk drops.  
- MongoDB: "Sharding is a method for distributing data across multiple machines... Horizontal scaling involves dividing the system dataset and load over multiple servers."  
- Vitess: "Sharding is a method of horizontally partitioning a database to store data across two or more database servers."  
**Sources**: PostgreSQL Docs, MongoDB Sharding Overview, Vitess Sharding Overview  
**Confidence**: HIGH

### Finding 2: Sharding Key Selection Critical for Write Distribution
**Claim**: Monotonically increasing sharding keys (e.g., `created_at`, `auto_increment`) cause write hotspots; optimal keys (e.g., `user_id`, `tenant_id`) have high cardinality, low frequency, and enable direct point-lookup routing.  
**Evidence**: 
- MongoDB: "If the shard key value is always increasing, all new inserts are routed to the chunk with maxKey... The shard containing that chunk becomes the bottleneck for write operations."  
- MongoDB: "Where possible, choose a shard key with high cardinality... ideal shard key distributes data evenly across the sharded cluster while also facilitating common query patterns."  
- Vitess: The Primary Vindex determines shard assignment; monotonic keys concentrate inserts onto one shard.  
**Sources**: MongoDB Shard Key Selection, Vitess Vindexes  
**Confidence**: HIGH

### Finding 3: Routing Algorithms and Consistent Hashing Minimize Data Movement
**Claim**: Hash Modulo (N % M) forces most keys to remap on cluster resizing; Consistent Hashing redistributes only O(K/N) keys on average.  
**Evidence**: 
- Karger et al. (1997): "In consistent hashing... when a hash table is resized, only n/m keys need to be remapped on average where n is the number of keys and m is the number of slots... In contrast, in most traditional hash tables, a change in the number of array slots causes nearly all keys to be remapped because the mapping is defined by a modular operation."  
- MongoDB: Describes chunk migration during balancer operation; virtual nodes reduce skewness.  
- Vitess: Key ranges and partitions define responsibility; resharding reorganizes data to match new scheme.  
**Sources**: ACM STOC '97 (Karger et al.), MongoDB Balancer Docs, Vitess Sharding Scheme  
**Confidence**: HIGH

### Finding 4: Core Complexities Introduced by Sharding
**Claim**: Sharding introduces cross-shard JOIN limitations, loss of native auto_increment, and complex resharding/rebalancing requirements.  
**Evidence**: 
- MongoDB: "If queries do not include the shard key or the prefix of a compound shard key, mongos performs a broadcast operation, querying all shards... Queries that involve multiple shards for each request are less efficient and do not scale linearly when more shards are added."  
- MongoDB: "In such cases, 'auto-increment' schemes that are often used by databases do not work well: the effort required to coordinate sequential numeric identifiers across a network can easily become a burden."  
- Vitess: Resharding description includes copying, verifying, and keeping data up-to-date while existing shards serve traffic; requires only seconds of read-only downtime.  
- Vitess Sequences: Explicitly solve auto_increment limitation by using a centralized sequence table with block allocation.  
- RFC 9562: Notes that auto-increment schemes don't work well in distributed systems; UUIDv7 offers time-ordered uniqueness with B-tree index locality.  
**Sources**: MongoDB Sharding Overview, Vitess Resharding, Vitess Sequences, RFC 9562  
**Confidence**: HIGH

### Finding 5: Handling Queries Without the Sharding Key
**Claim**: Queries missing the sharding key require scatter-gather operations unless supplemented by secondary indexes (Lookup Vindexes, Global Secondary Indexes).  
**Evidence**: 
- MongoDB: "If queries do not include the shard key... mongos performs a broadcast operation, querying all shards in the sharded cluster."  
- Vitess: "Secondary Vindexes are additional vindexes... In the absence of a Secondary Vindex, VTGate would have to scatter the query to all shards."  
- Vitess: Lookup Vindexes map a column value to keyspace IDs; extra write overhead but enables point-lookup routing for non-primary key queries.  
**Sources**: MongoDB Sharding Overview, Vitess Vindexes (Secondary Vindexes)  
**Confidence**: HIGH

## Areas of Agreement
All sources agree on:
- The architectural distinction between partitioning (single-instance) and sharding (multi-instance).
- The perils of monotonic keys (e.g., timestamps) as sharding keys.
- The benefit of high-cardinality sharding keys (e.g., user_id, tenant_id).
- The superior resizing properties of consistent hashing vs hash modulo.
- The necessity of distributed ID generation (sequences, UUIDv7) in sharded environments.
- The need for broadcast queries or secondary indexes when sharding key is absent.
- The added complexity of cross-shard transactions and joins.

## Areas of Disagreement
Minor interpretation differences exist on:
- **Cross-shard JOIN implementation**: Some sources imply JOIN is unsupported natively; others (Vitess Gen4 planner, CockroachDB) claim distributed query execution with transactional guarantees.
- **Degree of automatic vs manual rebalancing**: MongoDB’s balancer runs automatically; Vitess provides tooling but requires operator-triggered resharding; PostgreSQL partitioning remains manual DDL.

These reflect feature implementation variations, not contradictions in foundational principles.

## Limitations
- Research focused on Tier 1 sources (official documentation, standards, academic papers); did not include community forums or benchmark comparisons.
- Did not evaluate vendor-specific managed services (AWS Aurora, Google Cloud Spanner, Azure Cosmos DB) in depth.
- Did not measure actual performance metrics (latency, throughput) under load; relied on documented architectural properties.
- Version differences in declarative partitioning (e.g., PostgreSQL 10+ vs earlier inheritance-based) were not exhaustively contrasted.

## Conclusion
For high-scale transaction systems (e.g., e-commerce with millions of writes per day), the optimal sharding strategy uses:
1. A high-cardinality, non-monotonic sharding key (user_id or tenant_id) to evenly distribute writes and enable point-lookup routing.
2. Consistent hashing (or Vitess-style key ranges) to minimize data movement during cluster scaling.
3. Secondary Lookup Vindexes or Global Secondary Indexes to avoid scatter-gather for queries lacking the shard key.
4. Distributed ID generation (Vitess Sequences, UUIDv7/Snowflake) to replace auto_increment.
5. Awareness of inherent sharding tradeoffs: cross-shard JOINs require broadcast or complex planning, and resharding demands operational planning despite tooling assistance.
Partitioning remains valuable within each shard for manageable data lifecycle (time-based archival) but does not replace horizontal scaling via sharding.

This evidence-based approach balances write scalability, query efficiency, and operational simplicity for distributed databases at scale.