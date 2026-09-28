# Docs vs Code Audit

## Consistency Matrix

| Document Claim | Code Implementation | Status |
| --- | --- | --- |
| Logical table partitioning with range pruning | `internal/partitioning/table.go` `QueryRange` | PASS |
| Hash modulo routing vs consistent hashing | `internal/sharding/sharding.go` `ModuloRouter`, `ConsistentHashRouter` | PASS |
| Virtual nodes per physical shard (default 100) | `ConsistentHashRouter.vnodeCount` (100 in demo, configurable in struct) | PASS |
| Global Secondary Index (GSI / Lookup Vindex) point lookups | `internal/sharding/sharding.go` `GlobalSecondaryIndex`, `GetByEmailUsingGSI` | PASS |
| Scatter-gather parallel query broadcast with context timeout support | `internal/sharding/sharding.go` `ScatterGatherBroadcastWithContext` | PASS |
| Distributed unique ID generation (UUIDv7 & Sequence Block Allocator) | `internal/idgen/idgen.go` `NewUUIDv7`, `SequenceBlockAllocator` | PASS |
| Execution instructions in `README.md` (`go test ./...`, `go run ./cmd/demo`) | CLI demo `cmd/demo/main.go` & test suite `tests/sharding_test.go` | PASS |

## Discrepancies Found

None. README, design notes, execution results, code implementation, and demo output are in full alignment.
