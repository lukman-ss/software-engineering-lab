# Documentation vs Code Audit

Target Lab: `labs/32-database-sharding-and-partitioning`

## Comparison Matrix

| Component / Claim | Research Claim | README Claim | Engineering Implementation | Test / Demo Output | Match Status |
|---|---|---|---|---|---|
| Range Partitioning | Single-node logical range partitioning with partition pruning | `internal/partitioning`: In-engine range partitioning, partition pruning, fast dropping | `Table.QueryRange` filters partitions by time bounds; `DropPartition` deletes range | `TestPartitionPruning` passes; Demo outputs `Partitions Scanned: 1 / 4 (Pruned 3 partitions)` | MATCH |
| Naive Modulo vs Consistent Hash | Modulo moves ~N/(N+1) keys; Consistent hashing moves ~1/(N+1) keys | Modulo routing vs ConsistentHashRouter with virtual nodes minimizing data movement | `ModuloRouter` vs `ConsistentHashRouter` (vnode support) | Demo outputs: Modulo 79.84% moved vs Consistent Hash 12.00% moved on 4->5 resize | MATCH |
| Monotonic vs High-Cardinality Key | Monotonic keys cause single-shard write hotspot | Demonstrated write hotspot vs uniform distribution | `cmd/demo` compares date-string shard key vs `user_id` shard key | Demo outputs: 100% writes to single shard vs balanced distribution across 4 shards | MATCH |
| Scatter-Gather vs GSI | Non-shard key queries require broadcast or secondary index lookup | Scatter-gather parallel execution and GSI point lookups | `ScatterGatherBroadcastWithContext` vs `GetByEmailUsingGSI` | Demo outputs: Scatter-Gather (4 nodes broadcasted) vs GSI (1 node direct point lookup) | MATCH |
| Distributed ID Gen | Time-ordered UUIDv7 and sequence block allocation avoid collision | RFC 9562 UUIDv7 generator and Vitess-style Sequence Block Allocator | `idgen.NewUUIDv7()` and `SequenceBlockAllocator` | Demo outputs valid UUIDv7 and sequential IDs (1 2 3 4 5 6 7 8) | MATCH |

## Audit Discrepancy Checks

1. `DOC_CODE_MISMATCH`: None found. `README.md` accurately describes package layout, execution commands, and exported structures.
2. `TEST_CLAIM_MISMATCH`: None found. Tests verify claimed relocation bounds, pruning count, GSI routing, and context cancellation behavior.
3. `RESEARCH_IMPLEMENTATION_MISMATCH`: None found. Implementation covers all architectural patterns analyzed in `research/05-report.md`.
