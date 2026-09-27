# Engineering Audit Plan

Target Lab: labs/23-optimistic-vs-pessimistic-locking
Implementation Files:
- `internal/inventory/model.go`
- `internal/inventory/store.go`
- `internal/inventory/service.go`
- `cmd/demo/main.go`
Tests:
- `tests/locking_test.go`
Executable/Demo:
- `cmd/demo/main.go`
Approved Research Inputs:
- `research/05-report.md`
- `research-audit/07-verdict.md`
- `engineering/01-design.md`
Main Claims To Verify:
1. Naive read-modify-write causes lost updates under concurrent goroutines.
2. Pessimistic locking (simulated `SELECT ... FOR UPDATE` via row mutex) serializes updates and preserves exact stock.
3. Optimistic locking without retries detects stale version writes, rejects conflicting updates with `ErrOptimisticLock`, and maintains stock invariant without data corruption.
4. Optimistic locking with exponential backoff retries eventually converges all concurrent operations safely.
5. Atomic conditional updates (`UPDATE ... SET stock = stock - N WHERE stock >= N`) serialize updates safely at engine level.
6. Zero data race conditions detected under Go race detector.
Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`
Primary Risks:
- Thread-safety bugs in simulated in-memory store causing false races or undetected races.
- Test flakiness in naive lost update or retry convergence under different CPU schedulers.
- Discrepancies between claims in README/Design and actual code execution.
