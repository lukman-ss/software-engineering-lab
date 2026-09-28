# Engineering Design

Target Lab: `labs/32-database-sharding-and-partitioning`
Research Status: APPROVED

## Concept To Prove
1. **Partitioning vs Sharding**:
   - In-memory logical table partitioning (single node with sub-partition tables by time/range).
   - Horizontal physical sharding (multiple independent shard nodes with routers).
2. **Sharding Key Selection & Hotspot Effects**:
   - Monotonic keys (e.g. sequential IDs / timestamps) causing write hotspots on single shards.
   - High-cardinality keys (e.g. `user_id` / `tenant_id`) yielding balanced write distribution.
3. **Routing Algorithms**:
   - Hash modulo routing ($N \pmod M$) vs Consistent Hashing with virtual nodes.
   - Data movement measurement during cluster scale-out (from $M$ to $M+1$ nodes): Hash modulo forces $\approx \frac{M}{M+1}$ keys to move, while consistent hashing moves only $\approx \frac{1}{M+1}$ keys.
4. **Queries Without Sharding Key & Global Secondary Indexing**:
   - Scatter-gather broadcast across all shards when querying by non-shard key.
   - Lookup / Global Secondary Index (GSI) routing enabling point-lookup without broadcast.
5. **Distributed ID Generation**:
   - Sequence block allocator / UUIDv7 timestamp-ordered generator avoiding cross-node sequence collision.

## Expected Behavior
- Shard router routes `Insert` and `Query` operations accurately based on sharding key.
- Scale-out demo with Consistent Hashing redistributes only a small fraction of keys ($O(K/N)$), compared to majority redistribution in Hash Modulo.
- Non-shard-key lookups execute via scatter-gather (querying all shards concurrently) or direct point-lookup when backed by a Global Secondary Index (Lookup Index).
- Logical table partitioning on a single node demonstrates range pruning (only accessing relevant sub-partitions).

## Failure Scenario
- Monotonic write hotspot: inserts with sequential timestamps all target the newest shard/partition.
- Query without shard key without GSI incurs full cluster broadcast (scatter-gather latency/overhead).
- Resizing with hash modulo invalidates almost all cache/key locations.

## Success Criteria
- Sharded cluster implementation in pure Go standard library with clean abstractions.
- Automated tests covering:
  - Routing accuracy (Hash Modulo & Consistent Hash with virtual nodes).
  - Minimal data migration during consistent hash rebalancing vs modulo.
  - Logical table partition pruning on range filters.
  - Scatter-gather aggregation & error handling across shards.
  - Global secondary index lookups and synchronization.
  - Concurrent safe access with race detector passing.
- Executable demo printing concrete key distribution and rebalance comparison metrics.

## Architecture
- `internal/partitioning`: In-node table partitioner managing range-based partitions with partition pruning.
- `internal/sharding`:
  - `Shard`: Thread-safe key-value / record storage node simulating an independent database server.
  - `Router`: Pluggable routing interface (`ModuloRouter`, `ConsistentHashRouter`).
  - `ConsistentHashRing`: Ring buffer with configurable virtual nodes per physical shard.
  - `Cluster`: Orchestrator managing multiple `Shard` instances, scatter-gather queries, and resharding migration.
  - `LookupIndex`: Global secondary index mapping non-shard keys (e.g., `email`) to shard keys (`user_id`).
  - `IDGenerator`: UUIDv7 / distributed block allocator.
- `cmd/demo`: CLI demonstrating partitioning, monotonic vs hash distribution, consistent hashing data migration vs modulo, scatter-gather vs secondary lookup, and resharding.

## Components
1. `partitioning.Table`: Range-partitioned table with `Insert(Record)` and `QueryRange(min, max)` with pruning metric.
2. `sharding.ConsistentHash`: Hash ring using `fnv` or `crc32` with virtual nodes (e.g., 50-150 vnodes per node) and binary search ring lookup.
3. `sharding.Cluster`: Sharded cluster with scatter-gather execution with goroutines and `sync.WaitGroup`.
4. `sharding.LookupIndex`: Thread-safe secondary index.

## Test Strategy
- Unit tests for `ConsistentHashRing` correctness, distribution uniformity, and minimal relocation.
- Unit tests for `ModuloRouter`.
- Table partition pruning test verifying only matched partitions are scanned.
- Scatter-gather test with mock failures and partial shard responses.
- Concurrency race detector test (`go test -race ./...`).

## Execution Plan
1. Setup Go module `labs/32-database-sharding-and-partitioning`.
2. Implement partitioning engine (`internal/partitioning`).
3. Implement sharding cluster, consistent hashing, modulo router, scatter-gather, and lookup index (`internal/sharding`).
4. Implement distributed ID generation (`internal/idgen`).
5. Write unit and integration tests (`tests/`).
6. Implement CLI demo (`cmd/demo/main.go`).
7. Run tests with `-race` and execute demo.
8. Document implementation notes and execution results.

## Implementation Decisions
1. **In-Memory Storage**: Use thread-safe in-memory maps per shard to isolate distributed routing/sharding algorithmic mechanics without external DB dependencies.
2. **Hash Function**: Use `hash/fnv` (64-bit) from Go standard library for deterministic hashing without external third-party libraries.
3. **Vnode Count**: Default 100 virtual nodes per physical shard to balance ring size with uniform distribution.
