# Engineering Audit Plan

Target Lab: labs/29-saga-pattern
Implementation Files:
- internal/saga/orchestrator.go
- internal/saga/choreography.go
- internal/services/services.go
- cmd/demo/main.go
Tests:
- tests/saga_test.go
Executable/Demo: cmd/demo/main.go
Approved Research Inputs: research/ (inspected via engineering design notes)
Main Claims To Verify:
1. Orchestrator mode executes steps sequentially and compensates in LIFO order on failure.
2. Idempotency via processedID map prevents duplicate payment processing.
3. Semantic lock prevents concurrent sagas from modifying pending state.
4. Choreography model correctly propagates events and handles failures via compensating actions.
5. Concurrency safety under go test -race.
6. Demo output matches expected happy path and rollback scenarios.
Commands To Run:
- go build ./...
- go test -v ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Concurrency bugs in shared state (services use mutexes; need to verify correct locking).
- Compensation error handling may swallow errors incorrectly.
- Idempotency key collision across different sagas (uses paymentID alone, not combined with orderID).
- Semantic lock released only on approve/cancel; if saga stalls, lock may persist.
- Context cancellation path may not clean up correctly (tested).