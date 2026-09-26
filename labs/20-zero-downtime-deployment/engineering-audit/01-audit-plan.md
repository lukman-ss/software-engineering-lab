# Engineering Audit Plan

Target Lab: labs/20-zero-downtime-deployment
Implementation Files:
- `internal/db/db.go`
- `internal/server/server.go`
- `internal/worker/worker.go`
- `cmd/demo/main.go`
- `go.mod`
Tests:
- `tests/db_test.go`
- `tests/server_test.go`
- `tests/worker_test.go`
Executable/Demo: `cmd/demo/main.go`
Approved Research Inputs: `research/runs/2026-09-26-zero-downtime-deployment/05-report.md`, `engineering/01-design.md`, `engineering/02-implementation-notes.md`
Main Claims To Verify:
1. Health Probes: `/healthz/live` returns 200 OK continuously; `/healthz/ready` returns 503 until ready, then 200, and switches back to 503 on shutdown request.
2. Graceful Shutdown & Connection Draining: In-flight HTTP requests complete successfully during server shutdown without dropped connections or premature termination.
3. PreStop Hook Delay: Configurable preStop delay executes before listener shutdown and aborts immediately when context is canceled.
4. Worker Graceful Drain & Concurrency Safety: Background worker drains active jobs upon stop, drops subsequent enqueues safely under high concurrency without panicking on closed channel, and aborts queued jobs cleanly when drain timeout expires.
5. Database Expand/Contract: Backward-compatible reads and writes work seamlessly across legacy (`Name`) and modern (`FirstName`, `LastName`) schemas.
6. Demo Execution: Demo binary runs end-to-end cleanly, proving graceful draining across server and worker under mock SIGTERM.

Commands To Run:
```bash
go test -v -count=1 ./...
go test -race -v -count=1 ./...
go run ./cmd/demo
```

Primary Risks:
- Race conditions during concurrent `Enqueue` and `Stop` on worker channel.
- Deadlocks or premature cancellations during preStop delay or HTTP server graceful shutdown.
- Incomplete coverage of database schema edge cases (empty strings, single names, overwrites).
- Flaky port binding or test race conditions during server test execution.
