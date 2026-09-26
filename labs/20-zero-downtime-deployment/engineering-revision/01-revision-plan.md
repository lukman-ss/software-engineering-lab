# Engineering Revision Plan

Target Lab: labs/20-zero-downtime-deployment
Previous Verdict: APPROVED_WITH_WARNINGS

## Blocking Issues

None.

## Non-Blocking Issues

1. **GAP-01 (MEDIUM)** — TOCTOU race in `Enqueue` + `Stop`: `stopped.Load()` check passes, then `Stop()` closes `jobChan` before the send executes → send-on-closed-channel panic. No test exercises this path.
2. **GAP-02 (LOW)** — `engineering/03-execution-result.md` records 5 tests; actual suite has 14. Stale doc.
3. **GAP-03 (LOW)** — No test asserts Enqueue-after-Stop is silently rejected.
4. **GAP-04 (LOW)** — No test overwrites a legacy record with `SaveExpand` and re-reads it.
5. **GAP-05 (LOW)** — No test with multiple concurrent in-flight requests during graceful shutdown.

## Files To Change

- `internal/worker/worker.go` — fix TOCTOU: add mutex guard around stopped check + channel send
- `tests/worker_test.go` — add TestWorkerEnqueueAfterStop, TestWorkerConcurrentEnqueueStop
- `tests/db_test.go` — add TestDBLegacyOverwriteWithExpand
- `tests/server_test.go` — add TestServerMultiRequestDrain
- `engineering/03-execution-result.md` — update test count to actual

## Tests To Add/Modify

| Test | File | GAP |
|------|------|-----|
| TestWorkerEnqueueAfterStop | tests/worker_test.go | GAP-03 |
| TestWorkerConcurrentEnqueueStop | tests/worker_test.go | GAP-01 |
| TestDBLegacyOverwriteWithExpand | tests/db_test.go | GAP-04 |
| TestServerMultiRequestDrain | tests/server_test.go | GAP-05 |

## Validation Commands

```bash
cd labs/20-zero-downtime-deployment
go build ./...
go test -v ./...
go test -race ./...
go run ./cmd/demo
```
