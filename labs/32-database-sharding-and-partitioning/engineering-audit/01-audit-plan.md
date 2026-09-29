# Engineering Audit Plan

Target Lab: labs/32-database-sharding-and-partitioning
Implementation Files:
- internal/partitioning/table.go
- internal/sharding/sharding.go
- internal/idgen/idgen.go
- cmd/demo/main.go

Tests:
- tests/sharding_test.go

Executable/Demo:
- cmd/demo/main.go (`go run ./cmd/demo`)

Approved Research Inputs:
- research/05-report.md
- engineering/01-design.md
- engineering/02-implementation-notes.md

Main Claims To Verify:
1. Single-node logical range partitioning correctly prunes scanned partitions based on time range boundaries.
2. Modulo router vs Consistent Hash router redistribution ratio during cluster scale-out (Modulo moves ~M/(M+1) keys, Consistent Hashing moves ~1/(M+1) keys).
3. Monotonic sharding key causes single-shard write hotspot vs high-cardinality key uniform distribution.
4. Scatter-gather queries broadcast to all shards concurrently with context cancellation support.
5. Global Secondary Index (GSI) allows direct single-shard point lookups without broadcasting.
6. RFC 9562 time-ordered UUIDv7 and Vitess-style sequence block allocators generate unique, ordered IDs.

Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Race conditions or deadlocks during concurrent cluster access or scatter-gather context cancellation.
- Inaccurate routing / hash ring sorting edge cases in consistent hashing.
- Mismatch between README documentation/claims and actual implementation behavior.
