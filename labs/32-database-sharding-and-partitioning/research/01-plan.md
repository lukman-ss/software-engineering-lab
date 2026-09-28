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

## Search Strategy
- Query primary database documentation (PostgreSQL 18 DDL Partitioning, MongoDB 8.0 Sharding & Shard Key Selection, Vitess 24.0 Documentation on Sharding, Vindexes, Sequences, and Distributed Transactions).
- Inspect standard specifications (RFC 9562 for UUIDv7 and time-ordered identifiers).
- Examine computer science literature (Karger et al. 1997 Consistent Hashing paper).

## Expected Primary Sources
- PostgreSQL Documentation (Table Partitioning)
- MongoDB Documentation (Sharding & Shard Key Selection)
- Vitess Documentation (Sharding, Vindexes, Sequences, TwoPC)
- IETF RFC 9562 (UUID Specifications)
- Academic Paper: Karger et al. (1997) "Consistent Hashing and Random Trees"

## Risks / Unknowns
- Potential version differences between declarative partitioning (PostgreSQL 10+) vs legacy inheritance.
- Variations in terminology across database engines (e.g., MongoDB "chunks", Vitess "keyspace IDs", PostgreSQL "partitions").
- Complexity of distributed transaction isolation levels in sharded architectures (TwoPC vs eventual consistency).
