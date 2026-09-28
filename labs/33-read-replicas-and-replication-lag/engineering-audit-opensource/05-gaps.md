# Gap Analysis

## MISSING_TEST
Location: `router.ReadWithToken` timeout → primary fallback path (router.go:102-114)
Description: The branch where `WaitForLSN` returns error (context deadline) and the function falls back to a primary read is never exercised by any test. `TestReadWithToken_LSN` only tests the success path.
Severity: MEDIUM
Notes: Add a test that sets a very short `WaitTimeout` or uses a replica with huge lag to force timeout, then assert that the read returns from primary (node ID contains "fallback").

## MISSING_TEST
Location: `router.ReadWithToken` fast path (replica already at minLSN)
Description: The loop at router.go:94-100 that returns immediately if any replica has `AppliedLSN() >= minLSN` is untested.
Severity: LOW
Notes: Add a test with zero lag or pre-applied write to hit this branch.

## MISSING_TEST
Location: `cluster.Write` returning `ErrClusterClosed` after `Close()`
Description: No test verifies that writes to a closed cluster return the proper error.
Severity: LOW
Notes: Close cluster, attempt write, assert error equals `ErrClusterClosed`.

## MISSING_TEST
Location: `router.ReadNaive` with zero replicas (primary fallback)
Description: When `numReplicas=0`, `ReadNaive` should read from primary. Not tested.
Severity: LOW
Notes: Create cluster with 0 replicas, call `ReadNaive`, assert nodeID == "primary".

## MISSING_TEST
Location: `router.ReadLagAware` with zero replicas (primary fallback)
Description: When `numReplicas=0`, `ReadLagAware` should read from primary. Not tested.
Severity: LOW
Notes: Similar to above but for lag-aware read.

## MISSING_TEST
Location: `cluster.SetReplicationMode` dynamic switch
Description: No test changes replication mode after cluster creation and verifies behavior changes.
Severity: LOW
Notes: Create async cluster, write some entries, switch to sync, write more, verify new entries appear immediately on replicas.

## MISSING_TEST
Location: Concurrent `SetLag` / `Close` during writes
Description: Stress test where `SetLag` or `Close` is called concurrently with `Write`/`Read` to verify no panic/data corruption.
Severity: MEDIUM
Notes: Extend `TestConcurrentAccess_RaceFree` to also include mutators.

## MISSING_TEST
Location: Correctness under concurrent load (value assertions)
Description: `TestConcurrentAccess_RaceFree` runs many goroutines but never checks that read values match expected latest writes or that LSN ordering holds.
Severity: MEDIUM
Notes: Add assertions like: after each write, a subsequent `ReadWithToken` for same key/lsn must return the written value.

## BROKEN_IMPLEMENTATION
Location: `internal/cluster/cluster.go:220-225` (async Write WAL channel drop)
Description: On channel buffer full, WAL entry is silently dropped. Replica never receives the write, causing permanent inconsistency with no error returned.
Severity: MEDIUM
Notes: This is a silent data-loss bug in the replication path. Fix: block on send (remove `default`) or return error if channel full.

## MISSING_EDGE_CASE
Location: `router.ReadWithToken` with context already cancelled
Description: If input `ctx` is already `Done()`, the function should return quickly, but the `WaitForLSN` goroutine (if spawned) still leaks.
Severity: LOW
Notes: Test with `ctx := context.Background(); cancel(); ReadWithToken(ctx, ...)`.

## RACE_CONDITION
Location: None observed under `-race` (all tests and demo clean).
Description: Go race detector reported zero races across all test runs and demo execution.
Severity: NONE
Notes: No race conditions found in current implementation.

## UNHANDLED_ERROR
Location: `router.ReadWithToken` ignores error from `WaitForLSN` (line 107)
Description: On timeout, `err == context.DeadlineExceeded` is discarded; fallback to primary proceeds as if wait succeeded.
Severity: LOW
Notes: The caller cannot distinguish token-wait success from timeout-fallback. Consider returning the error or a boolean indicating freshness source.

## IMPLEMENTATION_OVERCLAIM
Location: None found.
Description: All claimed behaviors (stale read, sticky routing, token wait, lag-aware fallback, sync freshness) are substantiated by tests or demo.
Severity: NONE
Notes: No overclaim relative to evidence.

## Summary of Gaps:
- 7 `MISSING_TEST` items (coverage gaps)
- 1 `BROKEN_IMPLEMENTATION` (silent WAL drop)
- 1 `MISSING_EDGE_CASE` (pre-cancelled ctx)
- 1 `UNHANDLED_ERROR` (ignored WaitForLSN error)
- 0 `RACE_CONDITION`
- 0 `IMPLEMENTATION_OVERCLAIM`
Total gap count: 10 items.