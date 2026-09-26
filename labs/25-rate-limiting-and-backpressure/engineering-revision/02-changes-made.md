## Revision 1

Audit Issue: BoundedQueue submission after Stop() panic
Severity: MEDIUM
Files Changed:
- `internal/backpressure/queue.go`
- `internal/backpressure/queue_test.go`
Action:
- Added `ErrQueueStopped = errors.New("backpressure: queue stopped")`.
- Added `stopped atomic.Bool` to `BoundedQueue`.
- Guarded `TrySubmit` with `stopped.Load()` and context check before queue channel send.
- Ensured `Stop()` is idempotent using `stopped.CompareAndSwap(false, true)`.
- Added `TestBoundedQueue_SubmitAfterStop` test covering post-stop submission and repeated `Stop()` calls.
Verification: `go test -v -race ./...` PASS
Status: RESOLVED
