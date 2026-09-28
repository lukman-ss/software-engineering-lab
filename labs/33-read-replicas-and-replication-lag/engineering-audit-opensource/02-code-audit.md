# Code Audit

## Finding 1
Location: `internal/cluster/cluster.go:220-225` (`Cluster.Write` async branch)
Claimed Behavior: Async replication streams WAL entries to replicas with configurable lag.
Observed Implementation: Uses `select { case replica.walChannel <- entry: default: }` — if the channel buffer (1024) is full, the entry is silently dropped with no error returned to caller. The write returns success LSN even though one or more replicas never receive the entry.
Assessment: WARNING
Severity: MEDIUM
Notes: Silent data loss in replication path. Under burst write load or slow replica processing, replicas permanently miss writes. The lab's core claim of replication is violated when channel fills. Not observed in current tests (low concurrency, short runs) but is a latent correctness bug. Fix: block on send or return error; or monitor channel fullness.

## Finding 2
Location: `internal/cluster/cluster.go:79-96` (`Node.WaitForLSN`)
Claimed Behavior: Block until replica applies target LSN or context cancelled.
Observed Implementation: Spawns a goroutine that locks `n.mu`, waits on `n.cond.Wait()`. On context cancellation, returns error but the goroutine remains blocked in `cond.Wait()` until a future `Broadcast` wakes it. If no further writes occur, the goroutine leaks permanently.
Assessment: WARNING
Severity: MEDIUM
Notes: Goroutine leak on timeout/cancel. In `ReadWithToken`, timeout causes fallback to primary but leaves a goroutine waiting on cond. Repeated timeouts accumulate leaked goroutines. Mitigation: use `sync.WaitGroup` or context-aware cond pattern, or ensure periodic broadcasts eventually clean up (not guaranteed).

## Finding 3
Location: `internal/cluster/cluster.go:203-218` (`Cluster.Write` sync branch)
Claimed Behavior: Sync replication applies write to all replicas before returning.
Observed Implementation: Holds `c.mu.RLock()` while iterating replicas, sleeping `lag` per replica, then applying under `replica.mu.Lock()`. The `time.Sleep` is uninterruptible by context. During sync writes, `Close()` (which needs `c.mu.Lock()`) and `SetReplicationMode` are blocked.
Assessment: PASS
Severity: LOW
Notes: Acceptable for simulation; write latency matches demo (~100ms for 2 replicas × 50ms). Not a race, just a serialization bottleneck. Documented trade-off.

## Finding 4
Location: `internal/router/router.go:92-115` (`ReadWithToken`)
Claimed Behavior: Wait for replica to reach `minLSN`, fallback to primary on timeout.
Observed Implementation: Checks each replica for `AppliedLSN() >= minLSN`. If none, waits on `replicas[0]` only with `WaitForLSN`. On timeout, falls back to primary read. Error from `WaitForLSN` is ignored (line 107: `if err == nil`). The fallback does not distinguish between "replica caught up" and "timeout" — primary read may also be stale if primary crashed (not modeled).
Assessment: WARNING
Severity: LOW
Notes: Error from `WaitForLSN` (context deadline exceeded) is silently discarded. Caller cannot distinguish token-wait success from timeout-fallback. Test `TestReadWithToken_LSN` only exercises the success path; timeout/fallback branch is untested (see Test Audit).

## Finding 5
Location: `internal/router/router.go:117-141` (`ReadLagAware`)
Claimed Behavior: Route reads only to replicas within `MaxLSNDiff` of primary.
Observed Implementation: Computes `diff = primaryLSN - appliedLSN` (uint64). Since `primaryLSN` is always >= `appliedLSN` for valid async replication (primary increments first), `diff` is always defined. The `if primaryLSN >= appliedLSN` guard is always true; else branch dead code. Replicas ahead of primary are impossible in this model.
Assessment: PASS
Severity: LOW
Notes: No bug, but the guard is redundant. Could be removed or replaced with a comment explaining the monotonic LSN invariant.

## Finding 6
Location: `internal/cluster/cluster.go:149-166` (`replicaWorker`)
Claimed Behavior: Replica applies WAL entries after configured lag.
Observed Implementation: Reads `lagDuration` under RLock, then releases lock before `time.After(lag)`. If `SetLag` is called concurrently, the worker may read old or new lag for the same entry — behavior is "last read wins". Acceptable for simulation.
Assessment: PASS
Severity: LOW
Notes: Minor race in lag read timing, but not a correctness issue — lag is an approximation anyway.

## Finding 7
Location: `internal/cluster/cluster.go:235-239` (`SetReplicationMode`)
Claimed Behavior: Change replication mode at runtime.
Observed Implementation: Takes `c.mu.Lock()`. `Write` holds `c.mu.RLock()`. No deadlock cycle. Mode change is atomic and visible to subsequent writes.
Assessment: PASS
Severity: LOW
Notes: Not exercised in tests or demo. Works as designed.

## Finding 8
Location: `internal/router/router.go:80-90` (`ReadWithStickySession`)
Claimed Behavior: Route reads to primary for `StickyDuration` after write, then to replica pool.
Observed Implementation: Checks session `LastWriteTime`. If within TTL, reads primary. Else delegates to `ReadLagAware`. The session state is not cleared after TTL — subsequent reads re-check timestamp and may route to replica. Session state persists indefinitely (map growth).
Assessment: PASS
Severity: LOW
Notes: Memory leak in long-running process (sessions never evicted). Not relevant for short-lived demo/tests. Production would need TTL eviction.

## Finding 9
Location: `internal/cluster/cluster.go:56-65` (`Node.Read`)
Claimed Behavior: Return value and applied LSN for key.
Observed Implementation: Returns `n.appliedLSN` even on `ErrNotFound`. Test `TestNaiveReplicationLag_StaleRead` relies on this to assert `readLSN == 0`. Correct for simulation.
Assessment: PASS
Severity: LOW
Notes: In real DB, `appliedLSN` would be the last applied transaction LSN, not per-key. Here it's per-node global LSN — acceptable model.

## Finding 10
Location: `internal/cluster/cluster.go:182-229` (`Cluster.Write` primary update)
Claimed Behavior: Write commits on primary first, then replicates.
Observed Implementation: Primary data and `appliedLSN` updated under `c.primary.mu.Lock()` while `c.mu.RLock()` held. Then replication proceeds. If sync, replicas updated before return; if async, entries sent to channels. Primary LSN always >= replica LSN.
Assessment: PASS
Severity: LOW
Notes: Consistent with async/sync semantics. No data race detected by `-race`.