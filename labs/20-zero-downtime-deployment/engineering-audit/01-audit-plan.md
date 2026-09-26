# Engineering Audit Plan

Target Lab: labs/20-zero-downtime-deployment
Implementation Files:
- `internal/db/db.go`
- `internal/server/server.go`
- `internal/worker/worker.go`
- `cmd/demo/main.go`
Tests:
- `tests/db_test.go`
- `tests/server_test.go`
- `tests/worker_test.go`
Executable/Demo: `cmd/demo/main.go`
Approved Research Inputs: `engineering/01-design.md`, `engineering/02-implementation-notes.md`
Main Claims To Verify:
1. HTTP Server Readiness (`/healthz/ready`) returns 503 initially, 200 after `SetReady(true)`, and 503 immediately upon shutdown initiation.
2. Graceful Shutdown allows in-flight HTTP requests to complete cleanly while terminating underlying listeners.
3. PreStop delay mechanism simulates orchestrator routing table detachment delay prior to socket teardown, with graceful abortion on context cancellation.
4. Background Worker consumes jobs concurrently and cooperatively drains active in-flight jobs without loss during shutdown.
5. In-memory DB implements Expand & Contract pattern (dual-write and fallback reading across legacy and modern schema representations).
6. Race detector passes cleanly with no data races under concurrent workloads.
Commands To Run:
- `go test -v -count=1 ./...`
- `go test -race -v -count=1 ./...`
- `go run ./cmd/demo`
Primary Risks:
- Race condition on shared worker queue or state maps.
- HTTP server socket leak or hanging in-flight goroutines during shutdown.
- PreStop hook ignoring context cancellation timeouts.
- Worker dropping tasks or hanging when enqueue happens concurrently with stop.
