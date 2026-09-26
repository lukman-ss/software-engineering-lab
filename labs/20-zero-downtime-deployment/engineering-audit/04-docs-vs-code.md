# Docs vs Code Audit

Target Lab: labs/20-zero-downtime-deployment

## Comparison Matrix

| Component / Claim | README & Design Claim | Code & Test Implementation | Status |
|---|---|---|---|
| Database (`internal/db`) | In-memory Expand & Contract pattern with dual schema support and fallback reads | `internal/db/db.go` implements `SaveExpand`, `InsertLegacy`, and `GetUser` with fallback logic; tested in `tests/db_test.go` | PASS |
| Server (`internal/server`) | Readiness/Liveness probes, preStop delay simulation, connection draining | `internal/server/server.go` exposes `/healthz/live`, `/healthz/ready`, `/work`, configurable preStop, and graceful `Shutdown`; tested in `tests/server_test.go` | PASS |
| Worker (`internal/worker`) | Background worker gracefully draining active jobs upon stop signal | `internal/worker/worker.go` implements concurrency, mutex-protected queue stop, and graceful drain with timeout fallback; tested in `tests/worker_test.go` | PASS |
| Demo (`cmd/demo`) | CLI orchestrator demonstrating end-to-end ZDD lifecycle | `cmd/demo/main.go` runs all components and terminates cleanly on SIGTERM | PASS |
| Commands | `go test -v ./...`, `go test -race ./...`, `go run ./cmd/demo` | All documented commands execute successfully | PASS |

## Discrepancies Found
None. Documentation accurately reflects the codebase, test capabilities, and execution results.
