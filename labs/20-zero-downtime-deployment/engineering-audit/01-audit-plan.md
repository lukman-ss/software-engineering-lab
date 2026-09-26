# Engineering Audit Plan

Target Lab: labs/20-zero-downtime-deployment
Implementation Files:
- internal/server/server.go
- internal/worker/worker.go
- internal/db/db.go
- cmd/demo/main.go
Tests:
- tests/server_test.go
- tests/worker_test.go
- tests/db_test.go
Executable/Demo:
- cmd/demo/main.go
Approved Research Inputs:
- research/runs/2026-09-25-zero-downtime-deployment/05-report.md
- research/runs/2026-09-26-zero-downtime-deployment/05-report.md
Main Claims To Verify:
1. Server readiness state changes immediately to 503 on shutdown signal.
2. Configurable preStop delay executes properly to allow routing updates.
3. In-flight HTTP requests complete before server process exits (Graceful Shutdown).
4. Background worker stops taking new jobs on SIGTERM and drains existing active/queued jobs until completion or timeout.
5. In-memory storage correctly implements Expand-Contract pattern for dual schema versions.
6. Race detector passes cleanly with zero data races.
7. README instructions and architectural descriptions match actual Go implementation.
Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`
Primary Risks:
- Race conditions during concurrent HTTP request draining or worker job enqueue/stop operations.
- Improper HTTP context handling causing dropped requests or leaking goroutines.
- Unhandled preStop context cancellation during rapid pod termination.
