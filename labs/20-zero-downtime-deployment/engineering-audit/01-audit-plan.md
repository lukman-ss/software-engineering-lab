# Engineering Audit Plan

Target Lab: `labs/20-zero-downtime-deployment`
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

Executable/Demo:
- `cmd/demo/main.go`

Approved Research Inputs:
- `research/05-report.md`
- `research-audit/07-verdict.md`
- `engineering/01-design.md`
- `engineering/02-implementation-notes.md`

Main Claims To Verify:
1. HTTP Server probe routing (`/healthz/live`, `/healthz/ready`) returns correct HTTP status codes matching readiness state.
2. In-flight HTTP requests complete successfully without error or truncation during graceful shutdown.
3. Configurable `preStop` delay executes prior to closing HTTP listeners, aborting promptly on context cancellation.
4. Worker pool processes jobs concurrently and drains in-flight jobs gracefully upon shutdown signal without data loss or race conditions.
5. In-memory database handles Expand and Contract dual-write and fallback-read compatibility without data loss across legacy and modern fields.
6. Race detector passes cleanly on all concurrent operations under high load.
7. Demo executes cleanly end-to-end and mirrors claimed zero-downtime mechanics.

Commands To Run:
- `go test -v ./...`
- `go test -count=1 -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Data race during worker stop and concurrent enqueues.
- Premature listener closure or request cancellation in server graceful shutdown.
- PreStop sleep blocking past shutdown context deadline.
- In-memory database read logic improperly splitting or truncating legacy name fields.
- Mismatch between README documentation and actual source APIs.
