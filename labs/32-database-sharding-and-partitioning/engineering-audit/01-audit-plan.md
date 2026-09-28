# Engineering Audit Plan

Target Lab: `labs/32-database-sharding-and-partitioning`
Implementation Files:
- `internal/partitioning/table.go`
- `internal/sharding/sharding.go`
- `internal/idgen/idgen.go`
Tests:
- `tests/sharding_test.go`
Executable/Demo:
- `cmd/demo/main.go`
Approved Research Inputs:
- `research/05-report.md`
Main Claims To Verify:
1. Logical table range partitioning supports partition pruning on range queries.
2. Monotonic sharding keys create write hotspots, whereas high-cardinality keys yield uniform distribution.
3. Scaling cluster from N to N+1 nodes moves ~M/(M+1) keys in Hash Modulo vs ~1/(M+1) keys in Consistent Hashing with virtual nodes.
4. Non-sharded attribute queries use parallel Scatter-Gather broadcast or point-lookup via Global Secondary Index (GSI).
5. Distributed ID generation produces lexicographically time-ordered UUIDv7 and non-overlapping sequence blocks.
6. Code compiles cleanly, tests pass with `-race`, and demo runs accurately.
Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`
Primary Risks:
- Race conditions during concurrent cluster operations.
- Flaky tests due to hash ring variance or timing.
- Implementation bugs in scatter-gather goroutine synchronization or context cancellation.
