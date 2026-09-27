# Engineering Audit Plan

Target Lab: labs/23-optimistic-vs-pessimistic-locking
Implementation Files:
- internal/inventory/model.go
- internal/inventory/store.go
- internal/inventory/service.go
- cmd/demo/main.go

Tests:
- tests/locking_test.go

Executable/Demo:
- cmd/demo/main.go

Approved Research Inputs:
- research/05-report.md
- engineering/01-design.md
- engineering/02-implementation-notes.md

Main Claims To Verify:
1. Naive read-modify-write causes lost updates under concurrent access.
2. Pessimistic locking (`SELECT ... FOR UPDATE` simulation via row locks) prevents lost updates and maintains stock invariants.
3. Optimistic locking without retries detects version conflicts and aborts without data corruption.
4. Optimistic locking with retry (exponential backoff with jitter) eventually resolves conflicts under moderate contention.
5. Atomic conditional single-statement updates (`UPDATE ... SET stock = stock - N WHERE stock >= N`) update safely without application-level locks.
6. Code compiles cleanly, tests pass with `-race`, and demo runs producing genuine expected runtime output.

Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Race conditions or data corruption in store lock simulation during high contention.
- Retry loop non-determinism or potential deadlocks in exponential backoff timing.
- Inconsistencies between README documentation claims and actual test/demo implementation.
