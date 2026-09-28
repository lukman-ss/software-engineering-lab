# Engineering Audit Plan

Target Lab: `labs/28-timeouts-and-deadlines`
Implementation Files:
- `internal/deadline/deadline.go`
- `internal/retry/retry.go`
- `internal/circuit/circuit.go`
- `internal/idempotency/idempotency.go`
- `cmd/demo/main.go`

Tests:
- `internal/deadline/deadline_test.go`
- `internal/retry/retry_test.go`
- `internal/circuit/circuit_test.go`
- `internal/idempotency/idempotency_test.go`
- `tests/integration_test.go`

Executable/Demo: `cmd/demo/main.go`
Approved Research Inputs: `research/05-report.md`, `research-audit/07-verdict.md`

Main Claims To Verify:
1. Context deadline propagation & execution budget enforced (`internal/deadline`)
2. Exponential backoff with full jitter preventing retry storms (`internal/retry`)
3. Circuit breaker state machine transitions (`CLOSED` -> `OPEN` -> `HALF_OPEN` -> `CLOSED`) (`internal/circuit`)
4. In-memory idempotency deduplication with TTL (`internal/idempotency`)
5. High concurrency safety across stateful components

Commands To Run:
- `go test ./...`
- `go test -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Worker goroutine leak when `ExecuteWithBudget` hits context timeout (`deadline.go`)
- Expired keys accumulating in idempotency map causing memory leak (`idempotency.go`)
- Non-atomic state transitions or lock contention under concurrent load
