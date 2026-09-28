# Engineering Audit Plan

Target Lab: labs/29-saga-pattern
Implementation Files:
- internal/saga/orchestrator.go
- internal/saga/choreography.go
- internal/services/services.go
- cmd/demo/main.go
Tests:
- tests/saga_test.go
Executable/Demo:
- cmd/demo/main.go
Approved Research Inputs:
- research/05-report.md
- engineering/01-design.md
Main Claims To Verify:
1. Implementation supports both Orchestration and Choreography Saga models.
2. Failure in orchestrator triggers compensating transactions in reverse order (LIFO).
3. Payment service supports idempotency via processed payment keys.
4. Order service implements semantic locking countermeasure.
5. Code compiles cleanly and passes race detector (`go test -race ./...`).
6. Demo output is real and matches code execution behavior.
Commands To Run:
- go test -v ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Thread-safety / race conditions in orchestrator log state or event bus handler invocation.
- Mismatch between README component paths / descriptions and actual codebase structure.
- Incomplete failure compensation handling (e.g., partial compensation failures or context deadline cancellation).
