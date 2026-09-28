# Test Audit

## Finding 1
Location: `tests/replication_test.go` — `TestNaiveReplicationLag_StaleRead`
Claimed Behavior: Stale read occurs on lagging replica immediately after write; succeeds after lag elapses.
Observed Implementation: Asserts `ErrNotFound` on first `ReadNaive` (replica not yet applied), then sleeps 600ms (replica 500ms lag), then asserts successful read.
Assessment: PASS (happy path + failure path both covered)
Severity: LOW
Notes: Covers stale-read anomaly (design success criterion 1). Timing-dependent on 500ms vs 600ms but deterministic given sleep margins.

## Finding 2
Location: `tests/replication_test.go` — `TestStickySessionRouting`
Claimed Behavior: Within TTL, reads route to primary; after TTL, reads route to replica.
Observed Implementation: Sets `StickyDuration=300ms`. Writes, reads within TTL → asserts `nodeID == "primary"`. Reads for other session → asserts error (stale). Sleeps 550ms, re-reads → asserts `nodeID != "primary"` (replica).
Assessment: PASS (happy path + after-TTL transition covered)
Severity: LOW
Notes: Covers sticky routing (design success criterion 2). The "other-session" assertion confirms session isolation.

## Finding 3
Location: `tests/replication_test.go` — `TestReadWithToken_LSN`
Claimed Behavior: Token read waits for replica catch-up and returns correct value.
Observed Implementation: Writes LSN=1, reads with token LSN=1. Asserts value, `readLSN >= lsn`, and `elapsed > 150ms` (i.e., it must have waited ~200ms).
Assessment: WARNING
Severity: MEDIUM
Notes: Only exercises the success path of `ReadWithToken` (replica caught up after wait). The **timeout→fallback-to-primary** branch (router.go:102-114) is untested. Additionally, the `elapsed > 150ms` assertion is fragile: if the replica happens to already have the LSN applied before the check (timing-dependent on 200ms lag), the read returns immediately and `elapsed` ~0 fails the assertion. This is a flaky assertion. Covers causal token happy path; does NOT cover negative/timeout case.

## Finding 4
Location: `tests/replication_test.go` — `TestReplicaLagThreshold_Fallback`
Claimed Behavior: Lag-aware routing falls back to primary when replicas exceed lag threshold.
Observed Implementation: Sets replica lag to 10s, writes 5 keys, asserts `ReadLagAware` returns `nodeID == "primary (fallback-lag)"` and value correct.
Assessment: PASS
Severity: LOW
Notes: Covers lag-aware fallback (design success criterion 3). Correctly exercises the `MaxLSNDiff=2` threshold with replica at LSN 0 vs primary LSN 5.

## Finding 5
Location: `tests/replication_test.go` — `TestSynchronousReplication_Freshness`
Claimed Behavior: Sync replication guarantees immediate replica freshness.
Observed Implementation: Creates sync cluster, writes, `ReadNaive` immediately → asserts value present and `readLSN >= lsn`.
Assessment: PASS
Severity: LOW
Notes: Covers sync replication tradeoff demonstration (design success criterion 4). The immediate read succeeds because `Write` applies to replicas synchronously before returning.

## Finding 6
Location: `tests/replication_test.go` — `TestConcurrentAccess_RaceFree`
Claimed Behavior: No data races under concurrent read/write load.
Observed Implementation: 5 writers × 20 ops, 10 readers × 20 ops, all goroutines call mixed Read/Write/Token/Sticky reads. Uses `wg.Wait()` to join.
Assessment: WARNING
Severity: MEDIUM
Notes: Covers concurrency/race detection (design success criterion 5) — passes under `-race`. However, the test does NOT verify correctness of results — it ignores read values and only checks for crashes/races. No assertions on data consistency. Also does not test concurrent `Close`, `SetLag`, or `SetReplicationMode` during operations. A passing race detector is necessary but not sufficient: silent WAL drops (Finding 1 in Code Audit) could occur under this load but would go undetected since results aren't checked.

## Missing Coverage (no matching test)
- `ReadWithToken` fast path (replica already at LSN, no wait) — MISSING_TEST
- `ReadWithToken` timeout → primary fallback branch — MISSING_TEST
- `ReadWithToken` with pre-ctx already cancelled — MISSING_TEST
- `ReadNaive` / `ReadLagAware` with zero replicas (primary fallback) — MISSING_TEST
- `Cluster.Write` returning `ErrClusterClosed` after `Close()` — MISSING_TEST
- `Cluster.SetReplicationMode` dynamic switch async→sync — MISSING_TEST
- `Node.SetLag` concurrent with reads/writes — MISSING_TEST

## Summary
All 6 tests pass (verified: `go test -v -count=1 ./...` → 6 PASS, exit 0). Race detector clean (`go test -race -count=1 ./...` → ok). Tests cover all 5 design success criteria's happy paths. Gaps: missing negative/timeout/edge-case coverage; `ReadWithToken` timeout branch unexercised; concurrent test does not verify correctness.