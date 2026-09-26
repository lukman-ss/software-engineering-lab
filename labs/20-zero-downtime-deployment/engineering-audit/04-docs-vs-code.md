# Documentation vs Code Audit

## README Alignment Analysis

1. **Component Descriptions**:
   - `README.md` describes `internal/db` as demonstrating Expand and Contract pattern with legacy `Name` and modern `FirstName`/`LastName`. Confirmed in `internal/db/db.go`.
   - `README.md` describes `internal/server` exposing Liveness and Readiness probes, preStop hook, and graceful shutdown. Confirmed in `internal/server/server.go`.
   - `README.md` describes `internal/worker` as pulling jobs and completing active jobs upon shutdown. Confirmed in `internal/worker/worker.go`.
   - `README.md` describes `cmd/demo` wiring components together and simulating SIGTERM. Confirmed in `cmd/demo/main.go`.

2. **Execution Commands**:
   - `go run ./cmd/demo`: Verified working.
   - `go test -v ./...`: Verified passing.
   - `go test -race ./...`: Verified passing with zero race warnings.

3. **Mismatches Detected**:
   - None. Documentation matches the actual implementation and test commands.
