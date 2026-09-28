# Engineering Audit Plan

Target Lab: `labs/32-database-sharding-and-partitioning`
Implementation Files:
- `internal/partitioning/table.go` (single-node range partitioning and pruning)
- `internal/sharding/sharding.go` (hash modulo & consistent hash routing, cluster, GSI, scatter-gather)
- `internal/idgen/idgen.go` (UUIDv7 generator and Vitess-style sequence block allocator)
- `cmd/demo/main.go` (CLI interactive demonstration)

Tests:
- `tests/sharding_test.go` (partition pruning, key relocation, scatter-gather/GSI, ID generators, concurrency)

Executable/Demo:
- `cmd/demo/main.go`

Approved Research Inputs:
- `research/01-plan.md` through `research/05-report.md`
- `research-audit/07-verdict.md` (APPROVED)
- `engineering/01-design.md`, `engineering/02-implementation-notes.md`, `engineering/03-execution-result.md`
- `engineering-revision/01-revision-plan.md`, `engineering-revision/02-changes-made.md`

Main Claims To Verify:
1. Logical table partitioning performs deterministic partition pruning and fast partition drops.
2. Naive hash modulo triggers massive data relocation (~(N-1)/N, ~75-80%) upon scaling shard count.
3. Consistent hashing with virtual nodes limits relocation to ~1/N during scaling.
4. Scatter-gather queries broadcast across all shards while Global Secondary Index (GSI) performs direct point-lookup.
5. Context cancellation is handled cleanly by parallel scatter-gather routines without deadlocks or goroutine leaks.
6. UUIDv7 produces monotonic time-ordered IDs compliant with RFC 9562 format.
7. Sequence block allocator provides sequential integers across batch boundaries without gaps.
8. Concurrent cluster insertions and reads are thread-safe and pass `-race` detector without data races.

Commands To Run:
- `go test -v ./...`
- `go test -race -count=1 ./...`
- `go run ./cmd/demo`

Primary Risks:
- Race conditions during concurrent cluster insertions, shard additions, or scatter-gather aggregation.
- Inaccurate consistent hashing distribution causing false pass criteria in tests.
- Goroutine leaks or channel deadlocks in `ScatterGatherBroadcastWithContext`.
- Discrepancy between demo timings/metrics and actual algorithmic complexity.
