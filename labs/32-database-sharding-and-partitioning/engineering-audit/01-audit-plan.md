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
- `research-audit/07-verdict.md`

Main Claims To Verify:
1. In-engine table partitioning enables partition range pruning when querying by partition key (`CreatedAt`).
2. Modulo hashing produces extreme cluster reshuffling on scale-out ($~M/(M+1)$ keys moved), whereas consistent hashing relocates minimal fraction of keys ($~1/(M+1)$).
3. Monotonic shard keys concentrate writes onto a single shard creating write hotspots; high-cardinality keys distribute writes evenly.
4. Non-shard-key queries require full-cluster scatter-gather broadcasts across all shards unless indexed via Global Secondary Index (GSI / Lookup Vindex).
5. Scatter-gather handles timeouts / cancellations cleanly via context.
6. Distributed ID generation produces lexicographically sortable time-ordered UUIDv7 and collision-free central sequence block allocations.

Commands To Run:
```bash
cd labs/32-database-sharding-and-partitioning
go test ./...
go test -race ./...
go run ./cmd/demo
```

Primary Risks:
- Data race conditions in concurrent shard inserts, GSI indexing, or ring traversal.
- Inaccurate metric calculations in relocation or partition pruning.
- Mismatch between README documentation and actual implemented API.
