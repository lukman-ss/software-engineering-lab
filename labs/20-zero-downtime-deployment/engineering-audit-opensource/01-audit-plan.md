# Engineering Audit Plan

## Target Lab
`labs/20-zero-downtime-deployment`

## Implementation Files
- `cmd/demo/main.go` — CLI orchestrator wiring server, worker; simulates deployment lifecycle (readiness, in-flight workload, SIGTERM, graceful drain).
- `internal/db/db.go` — In-memory `UserStore` implementing Expand/Contract (Parallel Change): dual schema (`Name` legacy + `FirstName`/`LastName` modern) with transparent fallback read logic.
- `internal/server/server.go` — `Server` with `/healthz/live`, `/healthz/ready` probes (atomic.Bool), in-flight tracking (`sync.WaitGroup` + `atomic.Int32`), configurable `preStop` delay, and `Shutdown(ctx)` implementing graceful drain.
- `internal/worker/worker.go` — `Worker` consuming an in-memory buffered channel queue with cooperative termination: stops accepting new jobs, drains in-flight + buffered jobs, timeout-abort fallback.

## Tests
- `tests/db_test.go` (5 tests)
- `tests/server_test.go` (8 tests)
- `tests/worker_test.go` (5 tests)
- Package: external `tests` package importing internal packages.

## Executable / Demo
- `go run ./cmd/demo` — demonstrates readiness → in-flight request → SIGTERM → preStop → graceful HTTP drain → worker drain.

## Approved Research Inputs
- `research/01-design.md` (design target): Liveness/Readiness, graceful shutdown + preStop, Expand/Contract DB, cooperative worker termination.
- `research/03-execution-result.md` (expected outputs).

## Main Claims To Verify
1. DB supports dual schema versions with transparent fallback reads/writes (Expand/Contract).
2. Server exposes Liveness and Readiness probes; graceful shutdown drains in-flight requests; configurable preStop delay honored.
3. Worker stops pulling new jobs on shutdown but completes active job.
4. Demo genuinely demonstrates zero-downtime (in-flight request returns 200 during shutdown).
5. All tests pass including under `-race`.

## Commands To Run
1. `go build ./...`
2. `go vet ./...`
3. `go test -v -count=1 ./...`
4. `go test -race ./...`
5. `go run ./cmd/demo`

## Primary Risks
- Concurrency races in server/worker shutdown synchronization (activeCount/wg, completed slice, Enqueue TOCTOU).
- Incomplete failure handling on shutdown-context-expiry paths (preStop cancellation, http.Server.Shutdown timeout).
- Listener/resource leak on early-return shutdown paths.
- Demo output fabricated vs. real.
- Double-close of worker jobChan on repeated Stop calls.
