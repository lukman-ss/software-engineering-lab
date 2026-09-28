# Contradictions / Disagreements

No material contradictions discovered.

All aligned sources (PostgreSQL, MongoDB, Vitess, RFC 9562, Karger et al.) agree on the fundamental distinctions:

- **Partitioning**: Single-instance, logical; improves query performance and data management; does NOT scale CPU/memory; supports range/list/hash partitions.
- **Sharding**: Multi-instance, physical; scales system limits; introduces complexity for joins, IDs, and rebalancing.

Potential Interpretation Note (Low Contradiction Risk):

- **Cross-Shard JOIN Handling**: 
  - *Interpretation A*: JOIN is not supported natively; requires application-level scatter-gather and merge; or distributed SQL query routing (Vitess Gen4 planner, CockroachDB DTA).
  - *Interpretation B*: Some engines (CockroachDB, Vitess with TwoPC) implement distributed JOIN execution plans with explicit ACID guarantees.
  This is an architectural feature variance, not a contradiction in underlying sharding principles.

- **Partitioning vs Sub-partitioning**:
  - PostgreSQL supports sub-partitioning (partition of partitions).
  - MongoDB uses compound shard keys for similar effect.
  Both sources treat this as a design pattern rather than a contradiction.

## Sources In Agreement

| Claim | Source 1 | Source 2 | Source 3 |
|-------|----------|----------|----------|
| Write hotspot on monotonically increasing keys | MongoDB Shard Key Selection | RFC 9562 §2.1 | Vitess Vindexes |
| Consistent hashing minimizes key remapping | Karger et al. 1997 | Wikipedia | MongoDB Chunk Migration |
| High cardinality shard key recommended | MongoDB Shard Key Selection | Vitess Sequencing | PostgreSQL Best Practices |

## Conclusion

All Tier 1 sources (official documentation and academic standards) are consistent on the primary assertions of this research. No actionable contradictions affecting the research conclusions were identified.