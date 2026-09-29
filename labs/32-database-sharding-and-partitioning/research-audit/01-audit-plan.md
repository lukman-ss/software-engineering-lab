# Research Audit Plan: Database Sharding & Partitioning

## Target Lab
`labs/32-database-sharding-and-partitioning`

## Scope & Pipeline Override Notice
Pursuant to PIPELINE OVERRIDE instructions:
- Audit research only (`labs/32-database-sharding-and-partitioning/research/`).
- Do not audit implementation/code in this stage.
- Do not modify research files.
- Write all audit output to `labs/32-database-sharding-and-partitioning/research-audit/`.

## Files Reviewed
- `research/01-plan.md`
- `research/02-sources.md`
- `research/03-evidence.md`
- `research/04-contradictions.md`
- `research/05-report.md`
- `research/06-open-questions.md`

## Claims To Verify
1. Table partitioning splits logically one large table into smaller physical pieces within a single instance, aiding query pruning/bulk drops without scaling CPU/RAM.
2. Multi-instance sharding distributes data horizontally across independent servers to scale CPU, memory, and write capacity.
3. Monotonic keys (timestamps/autoincrement) used in range-based sharding create write hotspots on active maxKey/minKey chunks.
4. Optimal shard keys require high cardinality and low frequency to ensure balanced data distribution.
5. Consistent hashing redistributes ~1/n (or n/m) keys on resize, avoiding mass remapping caused by Hash Modulo (`key % M`).
6. Queries missing the shard key execute scatter-gather (broadcast) across all shards, causing latency amplification.
7. Secondary Lookup Vindexes mitigate scatter-gather at the cost of additional write overhead and cross-shard consistency maintenance.
8. Distributed ID generation (UUIDv7, Vitess Sequences, Snowflake) is required due to central coordination bottlenecks with auto_increment.
9. Cross-shard transactions via TwoPC trade latency for atomicity without providing full cross-shard ACID isolation.
10. Resharding redistributes data with minimal read-only cutover downtime, operated continuously (MongoDB balancer) or manually/tooling-assisted (Vitess).

## Primary Risks Identified During Pre-Audit
- **Source Verification / Accessibility**: Tier 1 sources include web URLs for PostgreSQL 18, MongoDB Manual, Vitess 24.0, RFC 9562, and ACM DOI / Wikipedia. Accessibility and fidelity must be verified.
- **Overgeneralization**: Simplified statements like "JOIN is not supported natively in sharded databases" must be checked against real capabilities (e.g. Vitess Gen4 query planner, MongoDB distributed transactions).
- **Mathematical / Formal Consistency**: Consistent hashing remap formulations (`1/n` vs `n/m`) require strict check against source definitions.
- **Completeness**: Identification of any missing sources or unverified assertions.

## Audit Strategy
1. Inspect source metadata and verify URL reachability and scope appropriateness.
2. Extract all major assertions across research documents and cross-check evidence alignment.
3. Evaluate reported contradictions and nuances.
4. Identify research gaps and missing evidence.
5. Issues verdict based on evidence completeness and rigor.
