# Engineering Audit Plan

Target Lab: labs/29-saga-pattern
Implementation Files:
- internal/saga/orchestrator.go
- internal/saga/choreography.go
- internal/services/services.go
- cmd/demo/main.go
- go.mod

Tests:
- tests/saga_test.go

Executable/Demo:
- cmd/demo/main.go

Approved Research Inputs:
- research/05-report.md (Finding 2: Local transactions sequence; Finding 3: Orchestration vs Choreography; Finding 4 & 9: Compensating transactions & LIFO rollback; Finding 7: Semantic lock countermeasure; Finding 8: Idempotency)
- engineering/01-design.md
- engineering/02-implementation-notes.md

Main Claims To Verify:
1. Orchestrator executes forward steps in sequential order and performs LIFO compensation on failure.
2. Choreography operates via decoupled EventBus pub/sub event dispatching with compensations.
3. Services provide idempotency keys (PaymentService) and semantic locking countermeasure (OrderService).
4. Compensation failure propagates errors and logs `StatusCompensateFailed`.
5. Context cancellation triggers rollback of executed steps and ceases forward execution.
6. Concurrent saga execution does not trigger data races under `-race`.
7. Demo executes cleanly and mirrors claimed output.
8. README accurately reflects components, usage, and behavior.

Commands To Run:
- `go test -v -count=1 ./...`
- `go test -v -count=1 -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Thread-safety of shared service instances across concurrent sagas.
- LIFO ordering correctness when a step fails mid-pipeline.
- Concurrency race conditions in Orchestrator logs or EventBus subscribers.
- Documentation divergence between README/notes and actual code signatures/behaviors.
