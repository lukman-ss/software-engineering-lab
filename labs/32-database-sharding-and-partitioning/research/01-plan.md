# Research Plan: Database Sharding & Partitioning

## Research Topic
Database Sharding and Table Partitioning in High-Scale Systems.

## Objective
Investigate partitioning versus sharding architecture, sharding key selection tradeoffs for high-write e-commerce systems, consistent hashing routing, cross-shard query strategies, and distributed ID generation to produce an authoritative, evidence-backed research foundation.

## Research Questions
1. What are the architectural differences, mechanical distinctions, and performance implications between Single-Instance Table Partitioning (Logical) and Multi-Instance Sharding (Physical)?
2. How does Sharding Key selection impact scalability, write hot-spotting, query scatter-gather overhead, and data distribution in high-volume transaction environments?
3. How do Consistent Hashing and Hash Modulo algorithms route requests across database shards, and how do they behave during node additions/scaling?
4. What are the core complexities and limitations introduced by sharding (Cross-Shard JOINs, Global Unique ID generation, and Resharding/Rebalancing)?
5. What are the optimal routing and secondary index strategies for lookup queries that do not include the primary Sharding Key?
6. How do Hashed Sharding strategies mitigate monotonic-key hotspots within range-partitioned systems?

## Search Strategy
- Query primary database documentation (PostgreSQL 18 / 19-beta DDL Partitioning, MongoDB 8.0 Sharding & Shard Key Selection & Hashed Sharding, Vitess 24.0 Documentation on Sharding, Vindexes, Sequences, and Distributed Transactions).
- Inspect standards specifications (RFC 9562 for UUIDv7 and time-ordered identifiers).
- Examine computer science literature (Karger et al. 1997 Consistent Hashing paper) plus corroborating reference (Wikipedia Consistent Hashing).

## Expected Primary Sources
- PostgreSQL Documentation (Table Partitioning)
- MongoDB Documentation (Sharding, Shard Key Selection, Hashed Sharding)
- Vitess Documentation (Sharding, Vindexes, Sequences, TwoPC)
- IETF RFC 9562 (UUID Specifications)
- Academic Paper: Karger et al. (1997) "Consistent Hashing and Random Trees"

## Risks / Unknowns
- PostgreSQL `/docs/current/` resolves to PG 18 (released September 2025) in September 2026; PG 19 is in beta. Terminology "partition" is reused across engines (PG partition vs Vitess "shard" vs MongoDB "chunk") and must not be conflated.
- Variations in terminology across database engines (e.g., MongoDB "chunks", Vitess "keyspace IDs", PostgreSQL "partitions").
- Complexity of distributed transaction isolation levels in sharded architectures (TwoPC vs eventual consistency).
- Consistent-hashing remap math: Karger et al. prove adding/removing one node moves 1/n of keys; Wikipedia restates the general resize formula as n/m — both must be cross-checked.
