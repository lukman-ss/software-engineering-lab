# Documentation vs Code Verification

Target Lab: `labs/32-database-sharding-and-partitioning`

## Comparison Matrix

| Claim Source | Claimed Behavior | Implementation Status | Verified In Code / Demo / Tests | Discrepancy |
| :--- | :--- | :--- | :--- | :--- |
| `README.md` | Range partitioning with partition pruning and fast partition drop | Implemented in `internal/partitioning/table.go` | Verified in `TestPartitionPruning` & Demo Sec 1 | None |
| `README.md` | `ModuloRouter` vs `ConsistentHashRouter` relocation behavior | Implemented in `internal/sharding/sharding.go` | Verified in `TestRoutingAndConsistentHashRelocation` & Demo Sec 2/3 | None |
| `README.md` | Scatter-gather parallel execution vs GSI point lookups | Implemented in `internal/sharding/sharding.go` | Verified in `TestClusterScatterGatherAndGSI` & Demo Sec 4 | None |
| `README.md` | UUIDv7 time-ordered ID and Vitess sequence block allocation | Implemented in `internal/idgen/idgen.go` | Verified in `TestIDGenerators` & Demo Sec 5 | None |
| `README.md` | Run instructions (`go test -v ./...`, `go test -race ./...`, `go run ./cmd/demo`) | All CLI commands work without failure | Verified by direct execution | None |

## Findings

1. `DOC_CODE_MISMATCH`: None detected.
2. `TEST_CLAIM_MISMATCH`: None detected.
3. `RESEARCH_IMPLEMENTATION_MISMATCH`: None detected.

## Demo Output Verification

Output produced by `go run ./cmd/demo` matches all described metrics in design documentation:
- Range pruning scanned 1/4 partitions (pruned 3).
- Monotonic key sharding resulted in 100% write hotspot on single node, while high-cardinality key distributed writes uniformly across shards.
- Cluster resize from 4 to 5 shards moved 79.84% keys with hash modulo vs 12.00% keys with consistent hashing.
- GSI point lookup avoided scatter-gather broadcast (1 node vs 4 nodes).
- UUIDv7 and sequence allocator generated valid, time-ordered, contiguous IDs.
