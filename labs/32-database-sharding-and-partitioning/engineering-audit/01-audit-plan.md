# Engineering Audit Plan

Target Lab: `labs/32-database-sharding-and-partitioning`
Implementation Files:
- `internal/partitioning/table.go`
- `internal/sharding/sharding.go`
- `internal/idgen/idgen.go`
- `cmd/demo/main.go`
Tests:
- `tests/sharding_test.go`
Executable/Demo:
- `cmd/demo/main.go`
Approved Research Inputs:
- `research/05-report.md`
- `research-audit/07-verdict.md`
- `engineering/01-design.md`
- `engineering/02-implementation-notes.md`
- `engineering/03-execution-result.md`
Main Claims To Verify:
1. Logical table partitioning isolates sub-tables and performs range pruning on time ranges.
2. Monotonic sharding keys produce severe write hotspots, while high-cardinality keys distribute evenly across shards.
3. Cluster scale-out under Consistent Hashing relocates roughly $O(K/N)$ keys (~12-25%), compared to Hash Modulo remapping ~75-80% of keys.
4. Queries without sharding keys execute via scatter-gather broadcast across all shards, whereas Global Secondary Index (GSI) lookups perform single-node direct point queries.
5. Distributed ID generators (RFC 9562 UUIDv7 and Vitess-style Sequence Block Allocator) generate unique, monotonic/time-ordered identifiers without cross-shard collisions.
Commands To Run:
- `go test ./...`
- `go test -race ./...`
- `go run ./cmd/demo`
Primary Risks:
- Data race conditions in concurrent scatter-gather, cluster routing, or partition scans.
- Inaccurate modulo or consistent hash rebalancing metrics.
- Discrepancy between README documentation and actual implemented API/demo behavior.
- Incomplete failure/edge case test coverage for empty clusters, unknown shards, or missing keys.
