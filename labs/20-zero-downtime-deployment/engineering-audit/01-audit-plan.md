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

Executable/Demo:
- cmd/demo/main.go (go run ./cmd/demo)

Approved Research Inputs:
- research/runs/2026-09-25-zero-downtime-deployment/
- research/runs/2026-09-26-zero-downtime-deployment/
- engineering/01-design.md
- engineering-revision/03-revision-result.md (previous revision: READY_FOR_ENGINEERING_REAUDIT)

Main Claims To Verify:
1. HTTP readiness probe rejects traffic when unready, accepts when ready
2. In-flight requests complete during graceful HTTP shutdown
3. PreStop hook imposes configurable delay before listener close
4. PreStop is abortable when context is cancelled
5. Worker stops accepting new jobs on shutdown but completes active job
6. Worker drain timeout forces cancel if jobs exceed timeout
7. Expand/Contract DB pattern: legacy reads supply FirstName/LastName via split; expand writes populate all fields; reads normalize Name when missing

Commands To Run:
- go build ./...
- go test -v ./...
- go test -race -count=1 ./...
- go run ./cmd/demo

Primary Risks:
- Race condition in worker completed slice (non-atomic append under concurrent access)
- TestWorkerShutdownTimeout timing sensitivity (job-slow 100ms vs timeout 20ms)
- TestWorkerGracefulShutdown ordering claim (expects exactly job-1 then job-2)
- wg.Add(1) called inside /work handler after activeCount.Add(1); potential race on wg between request start and Shutdown calling wg.Wait()
- Double-wait: srv.Shutdown() calls http.Server.Shutdown (which drains connections) AND then s.wg.Wait() — the wg.Wait may be redundant or racy
