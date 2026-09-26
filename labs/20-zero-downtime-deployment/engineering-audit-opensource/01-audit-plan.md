# Engineering Audit Plan

Target Lab: labs/20-zero-downtime-deployment
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
Approved Research Inputs: (Not audited per pipeline override - implementation and tests only)
Main Claims To Verify:
1. Database implements expand-contract pattern (parallel change) for zero-downtime schema migration.
2. Server implements liveness/readiness probes, preStop hook for load balancer detachment, and graceful connection draining.
3. Worker stops accepting new jobs on shutdown but completes in-flight jobs.
4. Demo orchestrator shows zero-downtime behavior by simulating traffic, SIGTERM, and observing completion of in-flight work.
Commands To Run:
- go test ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Race conditions in server active request counter or worker job completion tracking.
- Incorrect context cancellation leading to premature job abortion.
- Missing edge cases in database backward/forward compatibility (empty strings, nil).
- Demo may not accurately simulate production signals (uses mocked SIGTERM).