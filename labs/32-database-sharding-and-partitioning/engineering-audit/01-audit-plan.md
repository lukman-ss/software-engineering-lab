# Engineering Audit Plan

Target Lab: labs/32-database-sharding-and-partitioning
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
- `research-audit/07-verdict.md` (APPROVED)
- `engineering/01-design.md`
- `engineering/02-implementation-notes.md`
- `engineering/03-execution-result.md`
- `engineering-revision/01-revision-plan.md`
- `engineering-revision/02-changes-made.md`
- `engineering-revision/03-revision-result.md`

Main Claims To Verify:
1. Range-based table partitioning demonstrates partition pruning on time-range queries and partition dropping.
2. Naive modulo routing vs consistent hashing data relocation on cluster resize: modulo moves ~M/(M+1) keys while consistent hashing moves ~1/(M+1) keys.
3. Monotonic sharding key induces single-shard write hotspot vs high-cardinality shard key distributing writes across nodes.
4. Scatter-gather broadcast queries across shards in parallel with context cancellation/timeout handling.
5. Global Secondary Index (GSI) allows direct point-lookup without broadcast overhead.
6. ID generation via RFC 9562 UUIDv7 provides time-ordered 128-bit identifiers, and sequence block allocation provides sequential chunked IDs.
7. Concurrency safety: data structures support concurrent read/write access without data races under `go test -race`.

Commands To Run:
- `go test ./...`
- `go test -v -count=1 ./tests/...`
- `go test -race -v -count=1 ./tests/...`
- `go run ./cmd/demo`

Primary Risks:
- Race conditions during concurrent cluster insertions, router updates, or partition scans.
- Incorrect mathematical claim on consistent hashing relocation bound in tests/demo.
- Incomplete context propagation or goroutine leaks in scatter-gather queries.
- Discrepancy between README documentation, engineering notes, and actual implementation.
