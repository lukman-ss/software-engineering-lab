# Engineering Revision Plan

Target Lab: labs/37-cache-invalidation-strategies
Previous Verdict: APPROVED

## Blocking Issues

None.

## Non-Blocking Issues

1. **Failure Path Tests (MISSING_TEST)**: No unit tests exercise DB errors or `ErrNotFound` on `CacheAsideService` and `WriteThroughService`.
2. **XFetch Integration Test (MISSING_TEST)**: `XFetchService.Get` integration with probabilistic early expiration is not covered in `tests/cache_test.go` (only math unit test exists).
3. **Write-Behind Queue Overflow Edge Case (MISSING_EDGE_CASE)**: `WriteBehindService.Update` queue drop behavior on full buffer is unexercised.

## Files To Change

- `tests/cache_test.go`

## Tests To Add/Modify

- `TestCachePatterns_FailurePaths`: Verify DB read/write errors propagate properly in `CacheAsideService` and `WriteThroughService`.
- `TestXFetchService_Get`: Verify end-to-end `XFetchService.Get` hit, miss, and early recompute behaviors using deterministic `SetRandFunc`.
- `TestWriteBehindService_QueueOverflow`: Verify `WriteBehindService` gracefully handles buffer overflow without panicking or blocking.

## Validation Commands

```bash
cd labs/37-cache-invalidation-strategies
go test ./...
go test -race ./...
go run ./cmd/demo
```
