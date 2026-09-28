# Engineering Audit Plan

Target Lab: labs/32-database-sharding-and-partitioning
Implementation Files: internal/partitioning/table.go, internal/sharding/sharding.go, internal/idgen/idgen.go, cmd/demo/main.go
Tests: tests/sharding_test.go
Executable/Demo: cmd/demo/main.go
Approved Research Inputs: (pipeline override – not audited)
Main Claims To Verify:
- Partition pruning works
- Hotspot vs uniform distribution
- Resharding data movement percentages
- Scatter‑gather vs GSI query behavior
- UUIDv7 ordering & block allocator monotonic IDs
Commands To Run:
- go test -v ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Concurrency bugs, race conditions, incorrect lock ordering, missing error handling