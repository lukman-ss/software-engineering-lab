# Gap Analysis

Target Lab: labs/37-cache-invalidation-strategies

## Gaps Identified

### Gap 1
- **Type**: DOC_CODE_MISMATCH
- **Severity**: LOW
- **Location**: `engineering/01-design.md:46`
- **Description**: Design document architecture section mentions `jitter.go: TTL jitter calculation` as an architecture file. However, `TTLWithJitter` was placed into `internal/cache/store.go` and no `jitter.go` file was created.
- **Impact**: Minor developer confusion if reading the design document. The README correctly lists the actual files in `internal/cache/`.
- **Recommended Action**: Update `engineering/01-design.md` line 46 to reflect that TTL jitter calculation is housed in `store.go`.

### Gap 2
- **Type**: MISSING_TEST
- **Severity**: LOW
- **Location**: `tests/cache_test.go:307-325` (`TestWriteBehindService_QueueOverflow`)
- **Description**: The queue overflow test verifies that `svc.Get` returns the latest cache update after 10 updates to a buffer of size 2, but does not assert the DB write count after worker drains to explicitly prove that excess writes were dropped (`db.WriteCount() < 10`).
- **Impact**: Overflow drop behavior is functionally unasserted at the backing store layer, although cache read consistency is verified.
- **Recommended Action**: Add `time.Sleep(100*time.Millisecond)` and assert `db.WriteCount() <= 3` to explicitly prove write-drop semantics.

### Gap 3
- **Type**: MISSING_TEST
- **Severity**: LOW
- **Location**: `tests/cache_test.go` (`SWRService`)
- **Description**: `SWRService` has a deduplication mechanism (`revalidating` map) to prevent multiple concurrent background goroutines during the stale window. No test asserts that concurrent stale GET requests trigger exactly 1 revalidation.
- **Impact**: Duplicate revalidation suppression under concurrency is unverified by tests.
- **Recommended Action**: Add a concurrent SWR test asserting `svc.RevalidateCount() == 1` when multiple goroutines read a stale key simultaneously.

## Summary of Gap Counts
- CRITICAL: 0
- HIGH: 0
- MEDIUM: 0
- LOW: 3 (1 DOC_CODE_MISMATCH, 2 MISSING_TEST)

## Verification Statement
No CRITICAL or HIGH severity gaps were found. The core functionality (Cache-Aside, Write-Through, Write-Behind, SingleFlight, XFetch formula and execution, SWR, TTL Jitter) is proven by running code and passing tests with the race detector enabled.
