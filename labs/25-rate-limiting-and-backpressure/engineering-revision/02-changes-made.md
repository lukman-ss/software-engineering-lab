# Engineering Revision Changes Made

## Revision 1

Audit Issue: GAP-01 / Non-blocking stale refill calculation in `TokenBucket.RetryAfterSeconds`
Severity: LOW
Files Changed:
- `internal/ratelimit/bucket.go`
- `internal/ratelimit/bucket_test.go`
Action:
- Updated `TokenBucket.RetryAfterSeconds` to calculate dynamic tokens accounting for elapsed time delta since `lastRefill`.
- Added unit test `TestTokenBucket_RetryAfterSeconds` to test Retry-After calculation pre and post elapsed sleep.
Verification:
- `go test -v -count=1 ./...` PASS
- `go test -race ./...` PASS
Status: RESOLVED
