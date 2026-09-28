# Engineering Audit Plan

Target Lab: labs/33-read-replicas-and-replication-lag
Implementation Files: internal/cluster/*.go, internal/router/*.go, cmd/demo/main.go, tests/*.go
Tests: tests/replication_test.go
Executable/Demo: cmd/demo/main.go
Approved Research Inputs: (lab scoped, not audited now)
Main Claims To Verify:
- Async replication lag causes stale reads.
- Sticky session routing ensures read‑your‑own‑writes.
- LSN token wait provides freshness.
- Lag‑aware routing falls back to primary.
- Sync replication guarantees immediate replica freshness.
Commands To Run:
- go test -v ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Concurrency bugs, race conditions, incorrect fallback logic.
