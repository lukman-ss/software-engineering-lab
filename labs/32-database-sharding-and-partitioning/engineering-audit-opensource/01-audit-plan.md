# Engineering Audit Plan

Target Lab: labs/32-database-sharding-and-partitioning
Implementation Files: internal/partitioning/table.go, internal/sharding/sharding.go, internal/idgen/idgen.go, cmd/demo/main.go, tests/sharding_test.go
Tests: tests/sharding_test.go
Executable/Demo: cmd/demo/main.go
Approved Research Inputs: N/A (audit stage skips research)
Main Claims To Verify:
- Partition pruning works
- Modulo router causes high data movement on scale-out
- Consistent hash router minimizes movement
- Scatter‑gather broadcasts all shards; GSI provides point lookup
- UUIDv7 is time‑ordered; sequence allocator yields sequential IDs
- Concurrency safety (race detector)
Commands To Run:
- go test ./... && go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Incorrect migration metrics, race conditions, missing error handling.
