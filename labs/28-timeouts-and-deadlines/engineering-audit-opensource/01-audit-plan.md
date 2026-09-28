# Engineering Audit Plan

Target Lab: labs/28-timeouts-and-deadlines
Implementation Files: 
- internal/deadline/deadline.go
- internal/retry/retry.go
- internal/circuit/circuit.go
- internal/idempotency/idempotency.go
- cmd/demo/main.go
Tests:
- internal/deadline/deadline_test.go
- internal/retry/retry_test.go
- internal/circuit/circuit_test.go
- internal/idempotency/idempotency_test.go
- tests/integration_test.go
Executable/Demo: cmd/demo/main.go
Approved Research Inputs: engineering/01-design.md
Main Claims To Verify:
1. Context deadline propagation and execution within explicit time budgets.
2. Exponential backoff with full jitter to avoid synchronized retry storms.
3. Circuit breaker state machine (CLOSED, OPEN, HALF_OPEN) preventing cascading calls.
4. In-memory deduplication store preventing double execution during retries.
Commands To Run:
- go build ./...
- go test ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Race conditions in shared state (circuit breaker, idempotency store) due to improper locking.
- Incorrect deadline propagation causing early or late timeouts.
- Backoff jitter calculation not adhering to full jitter algorithm.
- Circuit breaker state transitions not matching specification under concurrent access.
- Idempotency store TTL eviction race conditions.