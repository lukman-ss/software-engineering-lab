# Docs vs Code Audit

Target Lab: `labs/20-zero-downtime-deployment`

## Documented vs Implemented Comparison

### 1. Database (`internal/db`)
- **README Claim**: Demonstrates the "Expand and Contract" pattern (Parallel Change). Supports writing dual schema versions (legacy `Name` and modern `FirstName`/`LastName`) and transparent fallback logic.
- **Code Reality**: `internal/db/db.go` implements `InsertLegacy`, `SaveExpand`, and `GetUser` with fallback logic.
- **Verdict**: MATCH.

### 2. HTTP Server (`internal/server`)
- **README Claim**: Exposes Liveness and Readiness probes. When a shutdown signal is sent, executes configurable `preStop` delay to simulate load balancer detachment latency, then performs graceful shutdown ensuring in-flight requests complete.
- **Code Reality**: `internal/server/server.go` implements `/healthz/live`, `/healthz/ready`, `preStop` timer check during `Shutdown(ctx)`, and graceful drain via `s.srv.Shutdown(ctx)` + `s.wg.Wait()`.
- **Verdict**: MATCH.

### 3. Background Worker (`internal/worker`)
- **README Claim**: Background daemon pulling jobs from a queue. Upon receiving shutdown signal, stops pulling new jobs but continues processing active jobs until completion.
- **Code Reality**: `internal/worker/worker.go` implements buffered queue, worker pool, atomic state toggle, channel close, and graceful drain timeout.
- **Verdict**: MATCH.

### 4. Demo CLI (`cmd/demo`)
- **README Claim**: Orchestrator wiring components together, simulating startup initialization, executing in-flight workloads, and sending termination signal to demonstrate zero-downtime draining behavior.
- **Code Reality**: `cmd/demo/main.go` runs worker, server, initial sleep for readiness, in-flight work request, mock `SIGTERM`, and orderly component shutdown.
- **Verdict**: MATCH.

### 5. Running Commands
- **README Claim**:
  ```bash
  go run ./cmd/demo
  go test -v ./...
  go test -race ./...
  ```
- **Code Reality**: All specified commands run without error or divergence.
- **Verdict**: MATCH.

## Discrepancies Found
- None. Documentation matches the actual implementation.
