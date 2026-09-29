# Engineering Changes Made

## Revision 1

Audit Issue: Gap 1 (DOC_CODE_MISMATCH) - Design document listed `jitter.go` as a file when `TTLWithJitter` is in `store.go`
Severity: LOW
Files Changed: `engineering/01-design.md`
Action: Updated architecture list in `01-design.md` to indicate TTL jitter calculation is located in `store.go`.
Verification: Verified file structure matches doc.
Status: RESOLVED

---

## Revision 2

Audit Issue: Gap 2 (MISSING_TEST) - Write-Behind overflow write-drop behavior not asserted at DB backing store
Severity: LOW
Files Changed: `tests/cache_test.go`
Action: Updated `TestWriteBehindService_QueueOverflow` to close service and assert `db.WriteCount() < 10` after queue burst.
Verification: `go test -v -run TestWriteBehindService_QueueOverflow ./...` passed.
Status: RESOLVED

---

## Revision 3

Audit Issue: Gap 3 (MISSING_TEST) - Missing test verifying SWR deduplication of background revalidation under concurrency
Severity: LOW
Files Changed: `tests/cache_test.go`
Action: Added `TestSWRService_ConcurrentRevalidationDeduplication` testing 10 concurrent stale GET calls and asserting `svc.RevalidateCount() == 1`.
Verification: `go test -v -race -run TestSWRService_ConcurrentRevalidationDeduplication ./...` passed.
Status: RESOLVED
