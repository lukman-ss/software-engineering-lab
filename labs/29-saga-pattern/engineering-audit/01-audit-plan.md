# Engineering Audit Plan

Target Lab: `labs/29-saga-pattern`
Implementation Files:
- `internal/saga/orchestrator.go`
- `internal/saga/choreography.go`
- `internal/services/services.go`
- `cmd/demo/main.go`
- `go.mod`

Tests:
- `tests/saga_test.go` (8 test cases: `TestOrchestrator_HappyPath`, `TestOrchestrator_FailureCompensatesLIFO`, `TestPayment_Idempotency`, `TestSemanticLock`, `TestOrchestrator_Concurrency`, `TestChoreography_Flow`, `TestChoreography_FailureCompensates`, `TestOrchestrator_CompensationErrorPropagated`, `TestOrchestrator_ContextCancellation`)

Executable/Demo:
- `cmd/demo/main.go`

Approved Research Inputs:
- `research/05-report.md`
- `engineering/01-design.md`
- `engineering/02-implementation-notes.md`

Main Claims To Verify:
1. Orchestrator executes steps sequentially and logs states correctly.
2. Step failure triggers LIFO compensating actions for all completed prior steps.
3. Payment service enforces idempotency on duplicate transaction IDs.
4. Order service enforces semantic locks on intermediate pending orders to mitigate lack of isolation.
5. Concurrency safety under high parallel load without race conditions (`go test -race`).
6. Choreography model orchestrates distributed events and compensation via event bus.
7. Context cancellation triggers compensation of already executed steps and aborts remaining steps.
8. Compensation errors are properly tracked and aggregated in return errors.

Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Thread-safety of shared in-memory state during concurrent executions.
- Error masking or silent failures during compensation phase.
- Missing edge case testing on context cancellation or compensation failures.
- Discrepancies between demo claims and actual code execution.
