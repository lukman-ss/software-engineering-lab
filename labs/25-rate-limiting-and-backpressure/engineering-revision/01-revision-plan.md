# Engineering Revision Plan

Target Lab: labs/25-rate-limiting-and-backpressure
Previous Verdict: APPROVED

## Blocking Issues
None.

## Non-Blocking Issues
1. `TokenBucket.RetryAfterSeconds` evaluated stale `tb.tokens` state without updating refill calculations for elapsed time.

## Files To Change
- `internal/ratelimit/bucket.go`
- `internal/ratelimit/bucket_test.go`

## Tests To Add/Modify
- Added `TestTokenBucket_RetryAfterSeconds` to verify dynamic elapsed token calculations for `RetryAfterSeconds`.

## Validation Commands
```bash
go test -v ./...
go test -race ./...
go run ./cmd/demo
```
