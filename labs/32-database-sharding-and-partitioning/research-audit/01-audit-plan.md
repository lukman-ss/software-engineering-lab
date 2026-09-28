# Audit Plan: Database Sharding & Partitioning Research

Target Lab: `labs/32-database-sharding-and-partitioning`

## Files Reviewed
- `research/01-plan.md`
- `research/02-sources.md`
- `research/03-evidence.md`
- `research/04-contradictions.md`
- `research/05-report.md`
- `research/06-open-questions.md`

## Claims To Verify
1. Table partitioning splits logically one large table into smaller physical pieces within a single instance, aiding partition pruning and fast dropping, but does not scale CPU/RAM across machines.
2. Sharding distributes data across multiple independent servers/instances to scale CPU, memory, disk I/O, and storage capacity.
3. Monotonically increasing sharding keys create severe write hotspots.
4. High-cardinality, low-frequency keys (`user_id`, `tenant_id`) ensure balanced write distribution and direct point-lookup routing.
5. Consistent Hashing redistributes $O(K/N)$ keys on average during resizing, whereas Hash Modulo ($N \pmod M$) forces nearly all keys to remap.
6. Queries without the sharding key require Scatter-Gather operations across all shards, amplifying tail latency.
7. Secondary Lookup Vindexes provide point-lookup routing for non-sharding-key queries at the cost of additional write overhead and cross-shard consistency maintenance.
8. Native auto_increment cannot generate globally unique IDs across independent database instances; distributed ID generation (UUIDv7, Snowflake, sequences) is required.

## Code To Execute
NOT APPLICABLE per pipeline override (Research audit only).

## Primary Risks
- Inaccurate URL or citation references for PostgreSQL 18, Vitess 24.0, MongoDB, and RFC 9562.
- Over-generalization of engine-specific sharding behaviors (e.g., MongoDB mongos router vs Vitess VTGate vs PostgreSQL declarative table partitioning).
- Arbitrary performance or threshold recommendations without source backing.

## Audit Strategy
1. Inspect all 9 cited sources in `research/02-sources.md` for URL validity, publisher accuracy, and scope.
2. Cross-check all 8 findings in `research/03-evidence.md` and `research/05-report.md` against source content and classification.
3. Analyze internal consistency across plan, evidence, report, contradictions, and open questions.
4. Record gaps, over-generalisations, and issues in structured audit reports.
5. Deliver evidence-based verdict in `audit/07-verdict.md` (mapped to `research-audit/07-verdict.md`).
