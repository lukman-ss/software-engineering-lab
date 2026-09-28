# Engineering Audit Plan

Target Lab: `labs/28-timeouts-and-deadlines`
Implementation Files:
- `internal/deadline/deadline.go`
- `internal/retry/retry.go`
- `internal/circuit/circuit.go`
- `internal/idempotency/idempotency.go`
Tests:
- `internal/deadline/deadline_test.go`
- `internal/retry/retry_test.go`
- `internal/circuit/circuit_test.go`
- `internal/idempotency/idempotency_test.go`
- `tests/integration_test.go`
Executable/Demo:
- `cmd/demo/main.go`
Approved Research Inputs:
- `research/05-report.md`
- `research-audit/07-verdict.md`
- `engineering/01-design.md`
Main Claims To Verify:
1. Context deadline propagation with budgeted execution aborts when context expires.
2. Exponential backoff with full jitter avoids synchronized retries and caps at max backoff.
3. 3-state Circuit breaker (`CLOSED`, `OPEN`, `HALF_OPEN`) trips upon reaching failure thresholds, blocks traffic during cooldown, tests downstream in half-open, and recovers upon reaching success threshold.
4. Thread-safe idempotency deduplication store evicts expired entries and prevents duplicate executions during retries.
5. Integration between retries, circuit breaking, and idempotency works safely without race conditions.
Commands To Run:
- `go test -count=1 -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`
Primary Risks:
- Concurrency races on state transitions in `circuit` and `idempotency`.
- Channel leakage or goroutine blocking on deadline cancellation.
- Mismatch between README documentation and actual implemented API signatures/behaviors.
