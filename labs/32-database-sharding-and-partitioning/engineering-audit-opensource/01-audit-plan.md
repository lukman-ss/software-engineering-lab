# Engineering Audit Plan

Target Lab: labs/32-database-sharding-and-partitioning
Implementation Files: internal/partitioning/*.go, internal/sharding/*.go, internal/idgen/*.go, cmd/demo/main.go
Tests: tests/sharding_test.go
Executable/Demo: cmd/demo/main.go
Approved Research Inputs: (pipeline overridden – not audited)
Main Claims To Verify:
- Partition pruning works
- Modulo vs consistent hash relocation ratios
- Scatter‑gather vs GSI query behavior
- Distributed ID generation uniqueness & ordering
Commands To Run:
```
go test -v ./...
go test -race ./...
go run ./cmd/demo
```
Primary Risks:
- Concurrency safety
- Accurate relocation metrics
- Correct GSI mapping
- UUIDv7 time‑ordering
