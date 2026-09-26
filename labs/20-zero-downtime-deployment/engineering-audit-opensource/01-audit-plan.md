# Engineering Audit Plan

Target Lab: Zero-Downtime Deployment Lab
Implementation Files:
- internal/db/db.go
- internal/server/server.go
- internal/worker/worker.go
- cmd/demo/main.go
Tests:
- tests/db_test.go
- tests/server_test.go
- tests/worker_test.go
Executable/Demo: cmd/demo/main.go
Approved Research Inputs: (Not audited per pipeline override)
Main Claims To Verify:
- Database demonstrates "Expand and Contract" pattern with dual schema versions and transparent fallback.
- Server exposes Liveness and Readiness probes, executes preStop delay on shutdown, and gracefully shuts down ensuring in-flight requests complete.
- Worker stops pulling new jobs on shutdown signal but continues processing current active job until completion.
- Demo orchestrator shows zero-downtime draining behavior by simulating startup, in-flight workloads, and termination signal.
Commands To Run:
- go test ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Race conditions in concurrent access to shared state (e.g., UserStore, Worker job queue).
- Improper handling of context cancellation leading to resource leaks or premature termination.
- Inadequate error propagation or logging.
- Demo may not accurately simulate real-world zero-downtime deployment scenarios.