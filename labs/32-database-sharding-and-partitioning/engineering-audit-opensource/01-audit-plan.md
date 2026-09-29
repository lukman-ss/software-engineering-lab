# Engineering Audit Plan

Target Lab: labs/32-database-sharding-and-partitioning
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
- `engineering/01-design.md`
- `engineering/02-implementation-notes.md`

Main Claims To Verify:
1. Range partitioning supports partition pruning and partition dropping.
2. Hash modulo sharding causes high data movement during cluster resharding (~N/(N+1) key movement), whereas consistent hashing with virtual nodes minimizes key relocation to ~1/(N+1).
3. Monotonic shard key leads to write hotspots; high-cardinality shard key balances write distribution.
4. Scatter-gather queries broadcast across all shards in parallel with context support, whereas GSI provides point lookup bypassing full broadcast.
5. Distributed ID generation produces lexicographically time-ordered UUIDv7 and sequential sequence block allocations.

Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Race conditions in thread-safe shard or table accesses.
- Inaccurate hash modulo or consistent hashing migration metrics.
- Unhandled context cancellations or thread deadlocks during scatter-gather execution.
- Discrepancy between README documentation claims and actual Go code.
