# Engineering Revision Plan

Target Lab: labs/20-zero-downtime-deployment
Previous Verdict: APPROVED_WITH_WARNINGS

## Blocking Issues

None.

## Non-Blocking Issues

1. **Timing Flakiness in Worker Timeout Test** — `TestWorkerShutdownTimeout` occasionally executed a second job (`job-dropped`) if worker loop picked up the next job right as the slow job finished before observing `ctx.Done()`.
2. **Server Listener Bind Timing Race** — `TestServerGracefulShutdown` and other server tests used fixed 50ms `time.Sleep` after `go srv.Start()`, intermittently causing connection resets under burst uncached test runs.
3. **Documentation Drift** — `engineering/03-execution-result.md` listed 5 tests originally and then 19/14 inconsistent descriptions instead of the actual stabilized suite.

## Files To Change

- `internal/worker/worker.go` — Check `w.ctx.Done()` before executing job duration delay, allowing prompt worker cancellation upon drain timeout.
- `tests/worker_test.go` — Adjust timing parameters in `TestWorkerShutdownTimeout` to deterministically verify drain timeout abandonment.
- `tests/server_test.go` — Replace fixed 50ms sleep after `srv.Start()` with `waitForServerReady` polling helper in server test cases.
- `engineering/03-execution-result.md` — Update execution result with passing test output and stabilized notes.

## Tests To Add/Modify

| Test | File | Issue |
|------|------|-------|
| TestWorkerShutdownTimeout | tests/worker_test.go | Stabilize timeout drain assertion |
| TestServerProbes | tests/server_test.go | Use `waitForServerReady` probe helper |
| TestServerGracefulShutdown | tests/server_test.go | Use `waitForServerReady` probe helper |
| TestServerPreStopHook | tests/server_test.go | Use `waitForServerReady` probe helper |
| TestServerPreStopContextCancellation | tests/server_test.go | Use `waitForServerReady` probe helper |
| TestServerInvalidDurationFallback | tests/server_test.go | Use `waitForServerReady` probe helper |
| TestServerReadyUnreadyTransition | tests/server_test.go | Use `waitForServerReady` probe helper |
| TestServerMultiRequestDrain | tests/server_test.go | Use `waitForServerReady` probe helper |
| TestServerWorkRequestCancellation | tests/server_test.go | Use `waitForServerReady` probe helper |

## Validation Commands

```bash
cd labs/20-zero-downtime-deployment
go test -v ./...
go test -count=1 ./...
go test -race ./...
go run ./cmd/demo
```
