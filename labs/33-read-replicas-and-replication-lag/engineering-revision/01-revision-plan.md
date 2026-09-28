# Engineering Revision Plan

Target Lab: `labs/33-read-replicas-and-replication-lag`
Previous Verdict: APPROVED (No blocking issues; 1 minor non-blocking warning identified in audit)

## Blocking Issues
None.

## Non-Blocking Issues
- GAP-01: Background helper goroutine waiting on `sync.Cond.Wait()` in `WaitForLSN` stays parked until next WAL broadcast if caller context expires early (`internal/cluster/cluster.go:79-96`).

## Files To Change
- `internal/cluster/cluster.go`
- `tests/replication_test.go`
- `engineering-revision/01-revision-plan.md`
- `engineering-revision/02-changes-made.md`
- `engineering-revision/03-revision-result.md`

## Tests To Add/Modify
- Add `TestWaitForLSN_ContextTimeout` to verify that `WaitForLSN` cleanly unblocks and returns `context.DeadlineExceeded` without hanging or leaking goroutines.

## Validation Commands
```bash
go test -v ./...
go test -race ./...
go run ./cmd/demo
```
