# Engineering Revision Plan

Target Lab: labs/15-load-testing
Previous Verdict: APPROVED_WITH_WARNINGS

## Blocking Issues
None.

## Non-Blocking Issues
1. **Timeout Claim Mismatch (LOW)**: `01-design.md` claimed timeouts occur for 95th percentile under stress load, whereas test duration (2s) is shorter than client timeout (5s), making timeouts structurally impossible. Clarified design documentation to describe high latency and tail degradation without overclaiming timeouts.
2. **Context Cancellation Handling in Request Creation (LOW)**: In `internal/loadtest/runner.go`, `http.NewRequestWithContext` errors were incrementing error count even when context cancellation was the underlying cause. Added check `if ctx.Err() == nil` before incrementing `errs`.

## Files To Change
- `labs/15-load-testing/engineering/01-design.md`
- `labs/15-load-testing/internal/loadtest/runner.go`

## Tests To Add/Modify
None required; existing unit and integration test suite verifies metrics invariants, error counting, dial errors, and context cancellation.

## Validation Commands
```bash
go test -v ./...
go test -race ./...
go run ./cmd/demo
```
