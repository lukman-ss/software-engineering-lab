# Engineering Audit Plan

Target Lab: Timeouts & Deadlines Laboratory (`labs/28-timeouts-and-deadlines`)
Implementation Files:
- `internal/deadline/deadline.go` — context deadline propagation within an explicit budget
- `internal/retry/retry.go` — exponential backoff with full jitter
- `internal/circuit/circuit.go` — circuit breaker state machine
- `internal/idempotency/idempotency.go` — in-memory deduplication store with TTL
- `cmd/demo/main.go` — wiring demonstration
Tests:
- `internal/deadline/deadline_test.go`
- `internal/retry/retry_test.go`
- `internal/circuit/circuit_test.go`
- `internal/idempotency/idempotency_test.go`
- `tests/integration_test.go`
Executable/Demo: `go run ./cmd/demo`
Approved Research Inputs (out of scope this stage): `research/`, `research-audit/` (not audited per pipeline override)
Main Claims To Verify:
1. Deadline propagation honors both per-call budget and parent context deadline.
2. Retry uses exponential backoff with full jitter, aborts on context cancellation.
3. Circuit breaker transitions `CLOSED → OPEN → HALF_OPEN → CLOSED` / `OPEN` and rejects calls when `OPEN`.
4. Idempotency store deduplicates within TTL and evicts expired keys.
5. Concurrency is safe (race detector clean).
6. Demo output is reproducible and matches described behavior.
Commands To Run:
```bash
cd labs/28-timeouts-and-deadlines
go test -count=1 ./...
go test -race -count=1 ./...
go run ./cmd/demo
```
Primary Risks:
- `deadline.ExecuteWithBudget`: goroutine leak on parent cancellation path (child goroutine `fn` may keep running).
- `circuit.Breaker.State()` mutates state via `Lock()` (not `RLock`) — acceptable but misleading; `RecordSuccess` resets `failures=0` in `CLOSED` which discards useful stats.
- `retry.CalculateBackoff`: `1 << uint(attempt-1)` overflow risk for large `MaxAttempts`; jitter math correct though.
- Idempotency `Get`/`Set` use coarse locking; fine for in-memory.
