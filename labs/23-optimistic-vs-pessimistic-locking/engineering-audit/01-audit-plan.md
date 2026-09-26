# Engineering Audit Plan

Target Lab: `labs/23-optimistic-vs-pessimistic-locking`
Implementation Files:
- `internal/inventory/model.go`
- `internal/inventory/store.go`
- `internal/inventory/service.go`
- `go.mod`

Tests:
- `tests/locking_test.go`

Executable/Demo:
- `cmd/demo/main.go`

Approved Research Inputs:
- `research/01-plan.md`
- `research/02-sources.md`
- `research/03-evidence.md`
- `research/04-contradictions.md`
- `research/05-report.md`
- `research-audit/07-verdict.md`
- `engineering/01-design.md`

Main Claims To Verify:
1. Unsynchronized read-modify-write causes lost update anomalies under concurrent goroutines.
2. Pessimistic row-level locking (`SELECT ... FOR UPDATE` equivalent) prevents lost updates and serializes access cleanly.
3. Optimistic locking with version checks detects concurrent conflicts, rejects stale writes, and prevents silent corruption.
4. Optimistic locking with jittered exponential backoff retries converges successfully under contention.
5. Atomic conditional updates (`UPDATE ... SET stock = stock - qty WHERE stock >= qty`) guarantee consistency without explicit application-level row lock holding.
6. Zero race conditions occur across the codebase under Go `-race` analysis.
7. Documentation and demo output reflect real executable behavior.

Commands To Run:
```bash
go build ./...
go test -v -count=1 ./...
go test -race -count=1 ./...
go run ./cmd/demo
```

Primary Risks:
- Race conditions in mock storage mutex handling or metrics counters.
- Flaky concurrency tests if artificial sleep / contention timings are poorly calibrated.
- Unhandled negative quantity or insufficient stock edge cases.
- Documentation overclaiming feature capabilities not implemented in code.
