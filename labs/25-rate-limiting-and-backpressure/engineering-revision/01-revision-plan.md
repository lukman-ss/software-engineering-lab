# Engineering Revision Plan

Target Lab: `labs/25-rate-limiting-and-backpressure`
Previous Verdict: `APPROVED` (with non-blocking lifecycle warning in opensource audit)

## Blocking Issues
None.

## Non-Blocking Issues
1. **MEDIUM — BoundedQueue submission panic after Stop()**: Calling `TrySubmit` after `Stop()` closed `bq.queue`, triggering a panic on closed channel send.

## Files To Change
- `internal/backpressure/queue.go`: Add `stopped` atomic flag, guard `TrySubmit` against closed channel sends by checking `stopped` flag and context cancellation, return `ErrQueueStopped`, make `Stop()` idempotent via CAS.

## Tests To Add/Modify
- `internal/backpressure/queue_test.go`: Add `TestBoundedQueue_SubmitAfterStop` to verify safe handling of `TrySubmit` post-`Stop()` and idempotency of `Stop()`.

## Validation Commands
```bash
go test -v -race ./...
go run ./cmd/demo
```
