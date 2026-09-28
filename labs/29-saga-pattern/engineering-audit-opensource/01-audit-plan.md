# Engineering Audit Plan

Target Lab: labs/29-saga-pattern
Audit Scope: implementation + tests only (research/content excluded per pipeline override)

Implementation Files:
- internal/saga/orchestrator.go
- internal/saga/choreography.go
- internal/services/services.go
- cmd/demo/main.go

Tests:
- tests/saga_test.go

Executable/Demo:
- cmd/demo/main.go (go run ./cmd/demo)

Approved Research Inputs:
- engineering/01-design.md (concept, expected behavior, success criteria)
- engineering/02-implementation-notes.md (design decisions, known limitations)

Main Claims To Verify:
1. Orchestrator executes steps sequentially and rolls back compensable steps in LIFO order on failure.
2. Context cancellation mid-saga stops further steps and compensates already-executed steps.
3. Compensation errors are aggregated and propagated.
4. Payment idempotency: duplicate ProcessPayment calls do not double-charge.
5. Semantic lock prevents duplicate order creation during pending state.
6. Choreography event bus drives saga flow and failure compensation.
7. Concurrency safety under `go test -race`.
8. Demo output is real and matches engineering/03-execution-result.md.
9. README matches code.

Commands To Run:
- go build ./...
- go test -v ./...
- go test -race ./...
- go vet ./...
- go run ./cmd/demo

Primary Risks:
- Mock service methods (ApproveOrder/CancelOrder) may lack existence checks, enabling phantom state transitions.
- compensate() error aggregation uses non-unwrapable slice formatting.
- Log semantics: StatusFailed may be recorded for a step that never executed (cancellation path).
- Design doc references pkg/... while code uses internal/... (doc drift).
- Idempotency test asserts no error but not identical state/amount.