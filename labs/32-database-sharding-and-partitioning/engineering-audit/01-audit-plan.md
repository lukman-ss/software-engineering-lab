# Engineering Audit Plan

Target Lab: `labs/32-database-sharding-and-partitioning`

Implementation Files:
- `internal/sharding/sharding.go`
- `internal/partitioning/table.go`
- `internal/idgen/idgen.go`
- `cmd/demo/main.go`

Tests:
- `tests/sharding_test.go`

Executable/Demo:
- `cmd/demo/main.go`

Approved Research Inputs:
- `research/01-plan.md`
- `research/02-sources.md`
- `research/03-evidence.md`
- `research/05-report.md`
- `research-audit/07-verdict.md` (Research verdict: APPROVED)

Main Claims To Verify:
1. In-engine table range partitioning prunes non-matching partitions during range queries.
2. Modulo hashing suffers from catastrophic key relocation (~(M/(M+1))) on scale-out, whereas consistent hashing relocates only ~(1/(M+1)) keys.
3. Monotonic shard keys concentrate writes onto a single shard (write hotspot), whereas high-cardinality keys distribute across shards.
4. Non-shard-key lookups incur full cluster scatter-gather broadcast, whereas Global Secondary Index (GSI) performs point lookups.
5. Distributed ID generation produces time-ordered UUIDv7 and gapless chunked sequence blocks.
6. Concurrency safety across shard, cluster, table partition, and ID generator operations under Go race detector.

Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Weak test assertions (e.g. testing lower bound only on relocation).
- Deadlocks or race conditions on concurrent shard writes or GSI index updates.
- Discrepancy between stated research conclusions and simplified in-memory simulation behavior.
- Failure propagation or timeout omissions in scatter-gather query broadcasting.
