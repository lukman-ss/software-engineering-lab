# Engineering Revision Plan

Target Lab: labs/37-cache-invalidation-strategies
Previous Verdict: APPROVED

## Blocking Issues

None.

## Non-Blocking Issues

1. **Design Doc Architecture Component Mismatch (DOC_CODE_MISMATCH)**: `engineering/01-design.md` listed `jitter.go` as a file, but `TTLWithJitter` is located in `store.go`.
2. **Write-Behind Queue Overflow Assertions (MISSING_TEST)**: `TestWriteBehindService_QueueOverflow` did not assert that write drops actually occurred at the DB layer (`db.WriteCount() < 10`).
3. **SWR Revalidation Deduplication Concurrency Test (MISSING_TEST)**: `SWRService` deduplication via `revalidating` map had no concurrent test asserting `svc.RevalidateCount() == 1` across multiple simultaneous stale reads.

## Files To Change

- `engineering/01-design.md`
- `tests/cache_test.go`

## Tests To Add/Modify

- Modify `TestWriteBehindService_QueueOverflow` to assert write drop semantics (`db.WriteCount() < 10`).
- Add `TestSWRService_ConcurrentRevalidationDeduplication` to test concurrent stale GET requests triggering exactly 1 background revalidation.

## Validation Commands

```bash
cd labs/37-cache-invalidation-strategies
go test ./...
go test -race ./...
go run ./cmd/demo
```
