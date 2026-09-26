# Engineering Revision Plan

Target Lab: labs/20-zero-downtime-deployment
Previous Verdict: APPROVED (with 2 non-blocking warnings and missing test coverage)

## Blocking Issues
None.

## Non-Blocking Issues
1. `internal/server/server.go:91-94`: PreStop delay uses unselectable `time.Sleep`, ignoring context cancellations/timeouts during shutdown.
2. `internal/worker/worker.go:79-85`: Worker graceful drain ceiling documentation. While active jobs intentionally run to completion, the ceiling was not explicitly marked with a `ponytail:` comment.
3. `tests/server_test.go`: Missing test coverage for context cancellation during `preStop` hook and mid-flight client HTTP request cancellations.

## Files To Change
- `internal/server/server.go`
- `internal/worker/worker.go`
- `tests/server_test.go`

## Tests To Add/Modify
- `TestServerPreStopContextCancellation`: Ensure server aborts `preStop` delay immediately if shutdown context is canceled/timed out.
- `TestServerWorkRequestCancellation`: Ensure server cleanly handles early client disconnection without leaving lingering active request count.

## Validation Commands
```bash
go test -v ./...
go test -race ./...
go run ./cmd/demo
```
