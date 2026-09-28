# Engineering Audit Plan

Target Lab:
labs/32-database-sharding-and-partitioning

Implementation Files:
- internal/partitioning/table.go
- internal/sharding/sharding.go
- internal/idgen/idgen.go
- cmd/demo/main.go

Tests:
- tests/sharding_test.go

Executable/Demo:
- cmd/demo/main.go

Approved Research Inputs:
- research/ directory (not audited per pipeline override)

Main Claims To Verify:
1. Logical table partitioning with range pruning
2. Sharding key selection and hotspot effects (monotonic vs high-cardinality)
3. Routing algorithms: hash modulo vs consistent hashing data movement
4. Queries without sharding key: scatter-gather vs Global Secondary Index
5. Distributed ID generation: UUIDv7 and sequence block allocator
6. Concurrent safe access

Commands To Run:
- go test ./...
- go test -race ./...
- go run ./cmd/demo

Primary Risks:
- Race conditions in concurrent access
- Incorrect data movement calculations during resharding
- GSI consistency with writes
- Partition pruning logic errors
- Hash function distribution bias