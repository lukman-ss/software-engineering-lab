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
- `research/01-plan.md`
- `research/02-sources.md`
- `research/03-evidence.md`
- `research/04-contradictions.md`
- `research/05-report.md`
- `research/06-open-questions.md`
- `research-audit/07-verdict.md` (APPROVED)
Main Claims To Verify:
1. Logical table partitioning provides partition pruning on range queries and fast partition drops.
2. Modulo routing exhibits high relocation on scale-out ($~M/(M+1)$) vs Consistent Hashing ring moving minimal fraction ($~1/(M+1)$).
3. Monotonic shard key induces write hotspot; high-cardinality shard key distributes evenly.
4. Non-shard-key lookups execute via Scatter-Gather broadcast or direct point-lookup via Global Secondary Index (GSI).
5. Distributed ID generators (UUIDv7 and Sequence Block Allocator) generate unique, time-ordered, and collision-free identifiers.
Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`
Primary Risks:
- Data race conditions in concurrent shard access, router dynamic modifications, or GSI lookups.
- Scatter-gather deadlocks or unhandled goroutine leaks under cancellation.
- Inaccurate relocation claims or non-deterministic test assertions.
