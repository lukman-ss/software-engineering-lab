# Engineering Revision Plan

Target Lab: `labs/15-load-testing`
Previous Verdict: APPROVED (with warnings)

## Blocking Issues
None.

## Non-Blocking Issues
1. `internal/server/server.go`: Queue semaphore acquisition and simulated query execution ignore `r.Context().Done()`.
2. `internal/loadtest/runner.go`: Response bodies closed without draining (`io.Copy(io.Discard, resp.Body)`), risking TCP connection reuse failure.
3. `tests/loadtest_test.go`: Missing test coverage for network dial errors and server context cancellation.

## Files To Change
- `labs/15-load-testing/internal/server/server.go`
- `labs/15-load-testing/internal/loadtest/runner.go`
- `labs/15-load-testing/tests/loadtest_test.go`

## Tests To Add/Modify
- Add `TestLoadTest_DialError` in `tests/loadtest_test.go`.
- Add `TestServer_ContextCanceled` in `tests/loadtest_test.go`.

## Validation Commands
```bash
go test -v -count=1 ./...
go test -race ./...
go run ./cmd/demo
```
