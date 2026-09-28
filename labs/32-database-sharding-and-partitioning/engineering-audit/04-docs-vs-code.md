# Docs vs Code Audit

## Comparison Matrix

| Component / Claim | Research / Design Claim | Implementation in Code | Demo / Test Behavior | Status |
|---|---|---|---|---|
| Single-Node Partitioning | Partition pruning on time ranges | `partitioning.Table.QueryRange` checks interval overlap | `QueryRange` scans 1/4 partitions, prunes 3 | MATCH |
| Partition Dropping | Instant partition DROP operation | `Table.DropPartition` removes partition slice entry | Tested in `TestPartitionPruning` | MATCH |
| Modulo Routing Relocation | Naive hash modulo relocates ~M/(M+1) (~80% for 4->5) | `ModuloRouter.GetShard` uses `hashKey(k) % N` | Demo relocates 79.84% (7984/10000) | MATCH |
| Consistent Hashing Relocation | Ring buffer relocates ~1/(M+1) (~20% for 4->5) | `ConsistentHashRouter` uses virtual nodes and binary search ring | Demo relocates 12.00% (1200/10000) | MATCH |
| Monotonic Hotspot | Sequential keys cause write hotspot | Ingesting date keys produces all writes on 1 shard | Demo shows 1000/1000 on shard-2 | MATCH |
| Scatter-Gather Broadcast | Querying without shard key broadcasts to all shards | `ScatterGatherBroadcast` concurrently queries all shards | Demo queries 4/4 shards | MATCH |
| Global Secondary Index | Secondary index enables direct point lookup | `GlobalSecondaryIndex` maps secondary key to ShardKey | Demo performs point lookup without broadcast (1 shard) | MATCH |
| Distributed IDs (UUIDv7 & Sequence) | UUIDv7 time-ordered, chunked sequence allocation | `NewUUIDv7` RFC 9562 & `SequenceBlockAllocator` | Validated in test and demo | MATCH |

## Documentation Accuracy
- `README.md` accurately describes architecture, components, test commands, and demo commands.
- `engineering/01-design.md` matches the implemented package structures and interfaces.
- `engineering/02-implementation-notes.md` accurately identifies limitations and trade-offs (e.g. no 2PC, in-memory transport).
- No misleading or inflated claims detected.
