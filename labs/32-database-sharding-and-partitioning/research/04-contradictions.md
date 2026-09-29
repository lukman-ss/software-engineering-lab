# Contradictions / Disagreements

No material contradictions discovered.

All Tier 1 sources (PostgreSQL, MongoDB, Vitess, RFC 9562, Karger et al.) agree on the fundamental distinctions:

- **Partitioning**: Single-instance, logical; improves query performance and data management; does NOT scale CPU/memory; supports range/list/hash partitions.
- **Sharding**: Multi-instance, physical; scales system limits; introduces complexity for joins, IDs, and rebalancing.

## Interpretation Nuances (Not Contradictions)

**Cross-shard JOIN handling**:
- Vitess Gen4 planner and MongoDB distributed transactions *can* execute cross-shard queries, but they impose higher latency and do not provide full ACID isolation across shards.
- The lab topic simplifies this to "JOIN is not supported natively," which is correct at the SQL engine level but is not the full architectural picture.
- Both sources agree: design tables to co-locate related data via the shard key to avoid cross-shard work entirely.

**Rebalancing trigger mechanism**:
- MongoDB's balancer runs automatically in the background.
- Vitess resharding is operator-triggered via tooling (vtctldclient).
- This is a product-design choice, not a disagreement about the underlying principle that resharding redistributes data.

**Consistent Hashing n/m vs 1/n**:
- Karger et al. 1997 proves adding or removing one server moves on average 1/n of keys (where n = current number of servers).
- Wikipedia restates the general resize result as n/m (n = total keys, m = total slots), which is the same quantity expressed differently.
- These are mathematically consistent and not contradictory.

## Sources In Agreement

| Claim | Source 1 | Source 2 | Source 3 |
|-------|----------|----------|----------|
| Write hotspot on monotonically increasing keys | MongoDB Shard Key Selection | MongoDB Hashed Sharding | RFC 9562 §2.1 |
| Consistent hashing minimizes key remapping | Karger et al. 1997 (DOI) | Wikipedia (n/m = 1/n) | MongoDB Chunk Migration |
| High cardinality shard key recommended | MongoDB Shard Key Selection | Vitess Vindexes (Primary Vindex) | PostgreSQL Best Practices |
| Scatter-gather without shard key | MongoDB Sharding (broadcast) | MongoDB Shard Key (scatter-gather) | Vitess Vindexes (absence of Secondary Vindex) |
| TwoPC trades latency for atomicity | Vitess Distributed Transactions | MongoDB Transactions (sharded) | — |

## Conclusion

All Tier 1 sources (official documentation and academic standards) are consistent on the primary assertions of this research. No actionable contradictions affecting the research conclusions were identified. Variations in terminology (chunks vs keyspace IDs vs partitions) and operational behavior (automatic vs operator-triggered resharding) reflect product-specific implementation choices, not fundamental disagreements.
