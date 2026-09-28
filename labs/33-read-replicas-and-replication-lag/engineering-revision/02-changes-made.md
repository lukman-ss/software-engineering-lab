## Revision 1

Audit Issue: Background goroutine parked on `sync.Cond.Wait()` upon context cancellation in `WaitForLSN` (`engineering-audit/02-code-audit.md:7-10`, `engineering-audit/05-gaps.md:5-10`).
Severity: LOW
Files Changed:
- `internal/cluster/cluster.go`
- `tests/replication_test.go`
Action: Added `stop` signaling channel and condition broadcast on context cancellation in `WaitForLSN` to unpark and terminate the background waiter goroutine immediately upon context cancellation. Added test `TestWaitForLSN_ContextTimeout` to verify timeout behavior.
Verification: Ran `go test -count=1 -v ./...`, `go test -count=1 -race ./...`, and `go run ./cmd/demo`.
Status: RESOLVED
