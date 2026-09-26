# Engineering Audit Plan

Target Lab: `labs/23-optimistic-vs-pessimistic-locking`
Implementation Files:
- `internal/inventory/model.go`
- `internal/inventory/store.go`
- `internal/inventory/service.go`

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
- `research-audit/07-verdict.md` (Verdict: APPROVED)

Main Claims To Verify:
1. Naive read-modify-write causes lost updates when executed concurrently.
2. Pessimistic locking (`SELECT ... FOR UPDATE` simulation via granular row mutex) serializes access and prevents data corruption / over-allocation.
3. Optimistic locking rejects stale writes based on version mismatch (`ErrOptimisticLock`).
4. Optimistic locking with jittered exponential backoff retries converges successfully under concurrent load.
5. Atomic conditional update (`UPDATE ... WHERE stock >= qty`) prevents negative stock and lost updates locklessly without full transactions.
6. Code compiles cleanly with zero external third-party dependencies.
7. Go race detector (`go test -race ./...`) runs with zero data races.
8. README documentation matches implementation and runnable demo commands.

Commands To Run:
```bash
go test -v ./...
go test -race ./...
go run ./cmd/demo
```

Primary Risks:
- Thread-safety of in-memory datastore internal state during concurrent test execution.
- Flakiness in lost update test if artificial delay is insufficient or scheduling timing varies.
- Live-lock or timeout in optimistic retry test if backoff parameters are misconfigured.
- Discrepancies between database semantic claims and in-memory mock semantics.
