# Docs vs Code Audit

Target Lab: `labs/32-database-sharding-and-partitioning`

## Comparison Matrix

| Component / Claim | README Claim | Code Implementation | Test Verification | Status |
|---|---|---|---|---|
| Logical Partitioning | In-engine range partitioning, pruning, drop partition | `internal/partitioning/table.go` | `TestPartitionPruning` | MATCH |
| Physical Sharding | Independent shard instances, `ModuloRouter`, `ConsistentHashRouter` | `internal/sharding/sharding.go` | `TestRoutingAndConsistentHashRelocation` | MATCH |
| Scatter-Gather & GSI | Parallel scatter-gather broadcast, GSI point lookup | `internal/sharding/sharding.go` | `TestClusterScatterGatherAndGSI` | MATCH |
| Distributed ID Generation | RFC 9562 UUIDv7, Vitess Sequence Block Allocator | `internal/idgen/idgen.go` | `TestIDGenerators` | MATCH |
| Demo Execution | Runnable CLI via `go run ./cmd/demo` | `cmd/demo/main.go` | Executed live and verified | MATCH |

## Discrepancies
None. All components listed in `README.md` are present, tested, and demonstrated in `cmd/demo/main.go`.
