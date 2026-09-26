# Engineering Audit Plan

Target Lab: `labs/20-zero-downtime-deployment`
Implementation Files:
- `internal/db/db.go`
- `internal/server/server.go`
- `internal/worker/worker.go`

Tests:
- `tests/db_test.go`
- `tests/server_test.go`
- `tests/worker_test.go`

Executable/Demo:
- `cmd/demo/main.go`

Approved Research Inputs:
- `research/runs/2026-09-25-zero-downtime-deployment/05-report.md`
- `research/runs/2026-09-26-zero-downtime-deployment/05-report.md`

Main Claims To Verify:
1. Server executes preStop delay and HTTP connection draining upon SIGTERM/Shutdown request without dropping active requests.
2. Server readiness probe toggles readiness status and detaches from load balancer immediately upon shutdown initialization.
3. Background queue worker stops accepting new jobs upon shutdown signal while draining active jobs to completion.
4. Database Expand-and-Contract schema pattern provides dual-write and backward-compatible read fallback for single/multi-word names.
5. All tests pass with race detector enabled (`go test -race ./...`).

Commands To Run:
- `go test ./...`
- `go test -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Race conditions during concurrent HTTP server shutdown or worker stop operations.
- Improper HTTP request lifecycle tracking causing premature termination before request draining completes.
- Incomplete test coverage for edge cases (e.g., context timeouts during preStop, empty field handling in DB migration logic).
