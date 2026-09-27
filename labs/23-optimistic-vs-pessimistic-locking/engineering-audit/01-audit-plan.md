# Engineering Audit Plan

Target Lab: labs/23-optimistic-vs-pessimistic-locking
Implementation Files:
- `internal/inventory/model.go`
- `internal/inventory/store.go`
- `internal/inventory/service.go`
Tests:
- `tests/locking_test.go`
Executable/Demo:
- `cmd/demo/main.go`
Approved Research Inputs:
- `research/05-report.md`
- `research-audit/07-verdict.md`

Main Claims To Verify:
1. Naive read-modify-write pattern allows concurrent overwrites causing lost updates.
2. Pessimistic row locking (`SELECT ... FOR UPDATE`) guarantees serial execution and zero lost updates.
3. Optimistic version checking detects stale data and rejects conflicting modifications.
4. Optimistic locking with jittered exponential backoff retry converges reliably under contention.
5. Atomic conditional update (`UPDATE ... WHERE stock >= qty`) executes safely without explicit long-lived locks.
6. Code passes `go test -race` with zero data races.

Commands To Run:
```bash
go test -v -count=1 ./...
go test -race -count=1 ./...
go run ./cmd/demo
```

Primary Risks:
- Race conditions or sync bugs in simulated store engine.
- False positive passes where concurrency anomalies are masked by overly coarse mutexes.
- Discrepancy between code semantics and research findings.
- Flaky tests in optimistic retry or race detection under concurrent goroutines.
