# Engineering Audit Plan

Target Lab: labs/32-database-sharding-and-partitioning
Implementation Files: 
- internal/partitioning/table.go
- internal/sharding/sharding.go
- internal/idgen/idgen.go
- cmd/demo/main.go
Tests:
- tests/sharding_test.go
Executable/Demo: go run ./cmd/demo
Approved Research Inputs: Not required for engineering audit (per pipeline override)
Main Claims To Verify:
1. Logical partition pruning works correctly (scans only relevant partitions)
2. Sharding routers (ModuloRouter and ConsistentHashRouter) distribute keys as expected
3. Resharding key movement percentages match theoretical expectations
4. Scatter-gather broadcast vs Global Secondary Index (GSI) query performance and correctness
5. UUIDv7 generation is time-ordered and SequenceBlockAllocator produces sequential blocks
6. Concurrent access safety (no data races)
7. README.md accurately describes the demonstrated concepts
Commands To Run:
- go test ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Concurrency bugs in sharding/index structures
- Incorrect key movement calculations during resharding
- Demo output may be hardcoded or not reflect actual execution