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
- research/runs/2026-09-25-zero-downtime-deployment/05-report.md
- research/runs/2026-09-26-zero-downtime-deployment/05-report.md
- research-audit/07-verdict.md (APPROVED_WITH_WARNINGS)
- engineering/01-design.md

Main Claims To Verify:
1. HTTP server exposes /healthz/live and /healthz/ready probes with correct semantics
2. Readiness probe returns 503 when unready, 200 when ready
3. Graceful shutdown waits for in-flight requests to complete before terminating
4. Configurable preStop delay executes before HTTP listener close
5. preStop delay aborts on context cancellation
6. Background worker drains queued jobs before stopping
7. Worker stops new job intake on Stop(); finishes active job
8. Worker times out and cancels context if drain exceeds timeout
9. DB Expand/Contract: legacy writes readable via new API (Name split to FirstName/LastName)
10. DB Expand/Contract: new writes readable with combined Name field
11. Demo runs end-to-end with zero dropped requests

Commands To Run:
- go build ./...
- go test -v ./...
- go test -race ./...
- go run ./cmd/demo

Primary Risks:
- Race condition in worker.Stop() between close(jobChan) and in-flight Enqueue()
- TestWorkerShutdownTimeout: job-slow uses time.Sleep (non-preemptible), may not be interrupted by context cancel during active processing
- Server wg.Wait() after srv.Shutdown(ctx) may double-count: both http.Server.Shutdown and wg.Wait() wait for active handlers; redundant but not incorrect
- Port binding conflicts in parallel test execution (fixed ports per test)
- Demo self-sends SIGTERM via channel injection; does not test real OS signal handling path
