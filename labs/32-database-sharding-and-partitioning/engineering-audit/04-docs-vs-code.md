# Documentation vs Code Verification

## Matrix

| Component | README Claim | Code Implementation | Status |
| --- | --- | --- | --- |
| Logical Partitioning | `internal/partitioning`: In-engine range partitioning, partition pruning, fast partition dropping. | `Table.QueryRange` counts scanned vs total partitions; `Table.DropPartition` drops partition slice entry. | PASS |
| Physical Sharding & Routing | `internal/sharding`: ModuloRouter (heavy movement), ConsistentHashRouter (minimal movement ring with vnodes). | `ModuloRouter` uses `h % N`; `ConsistentHashRouter` uses sorted ring of vnodes (`hashKey(shardID#i)`). | PASS |
| Querying & GSI | `Cluster`: Scatter-gather parallel execution and Global Secondary Index (GSI / Lookup Vindex) point-lookups. | `ScatterGatherBroadcastWithContext` executes goroutines per shard; `GetByEmailUsingGSI` maps email to `ShardKey`. | PASS |
| Distributed ID Gen | `internal/idgen`: RFC 9562 UUIDv7 generator, Vitess-style Sequence Block Allocator. | `NewUUIDv7` constructs 48-bit ms timestamp + ver7 + rand; `SequenceBlockAllocator` allocates blocks in chunks. | PASS |
| Quick Start Commands | `go test -v ./...`, `go test -race ./...`, `go run ./cmd/demo` | All commands execute cleanly with zero errors/warnings. | PASS |

## Discrepancies Found

- None. Documentation in `README.md`, `engineering/01-design.md`, and `engineering/02-implementation-notes.md` accurately describes the code and test behavior.
