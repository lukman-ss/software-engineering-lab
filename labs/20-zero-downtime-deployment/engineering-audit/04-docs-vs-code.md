# Docs vs Code Audit

## Documentation Verification

Target: `README.md` compared against `internal/`, `cmd/demo/`, and `tests/`.

### 1. Component Descriptions
- **Database (`internal/db`)**: README states it demonstrates "Expand and Contract" pattern (Parallel Change) with dual schema versions (`Name` and `FirstName`/`LastName`) and transparent fallback logic.
  - Verification: Code in `internal/db/db.go` and tests in `tests/db_test.go` directly mirror this description.
  - Result: MATCH.

- **Server (`internal/server`)**: README states it exposes Liveness and Readiness probes, executes configurable `preStop` delay upon shutdown signal, and drains in-flight requests before termination.
  - Verification: Implemented in `internal/server/server.go:NewServer` and `Shutdown`. Probes at `/healthz/live` and `/healthz/ready`.
  - Result: MATCH.

- **Worker (`internal/worker`)**: README states worker is a background daemon pulling jobs from a queue, stopping new pulls on shutdown, and processing active jobs until completion.
  - Verification: Implemented in `internal/worker/worker.go:Start`, `Enqueue`, and `Stop`.
  - Result: MATCH.

- **Demo (`cmd/demo`)**: README describes CLI orchestrator wiring components, simulating startup, running workloads, and sending termination signal.
  - Verification: Implemented in `cmd/demo/main.go`. Output matches claimed steps.
  - Result: MATCH.

### 2. Execution Commands
- README commands:
  - `go run ./cmd/demo` -> Tested, works as documented.
  - `go test -v ./...` -> Tested, works as documented.
  - `go test -race ./...` -> Tested, works as documented.
  - Result: MATCH.

### 3. Discrepancies
- None identified. No `DOC_CODE_MISMATCH`, `TEST_CLAIM_MISMATCH`, or `RESEARCH_IMPLEMENTATION_MISMATCH`.
