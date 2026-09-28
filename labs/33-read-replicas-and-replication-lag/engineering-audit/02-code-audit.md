# Code Audit

Target Lab: `labs/33-read-replicas-and-replication-lag`

## Finding 1

Location: `internal/cluster/cluster.go:36-46` and `internal/cluster/cluster.go:79-107`
Claimed Behavior: Replicas notify waiters upon catching up to target LSN via condition variable with context timeout support.
Observed Implementation: `newNode` initializes `n.cond = sync.NewCond(&n.mu)`. `WaitForLSN` launches a waiter goroutine waiting on `n.cond.Wait()` holding `n.mu`, while the parent goroutine selects between context cancellation and waiter completion. On context cancellation, `stop` is closed, and `n.cond.Broadcast()` is triggered under lock to wake the waiter.
Assessment: PASS
Severity: LOW
Notes: Synchronization pattern properly avoids lost wakeups and handles context expiration cleanly.

## Finding 2

Location: `internal/cluster/cluster.go:193-240`
Claimed Behavior: Writes commit to primary, increment monotonic LSN, and replicate synchronously or asynchronously.
Observed Implementation: `c.currentLSN.Add(1)` atomically increments LSN. Primary writes are protected by `c.primary.mu`. In synchronous mode, writes sequentially iterate over replicas, apply sleep simulation, update replica map, and broadcast catch-up under replica lock. In asynchronous mode, WAL entries are dispatched to replica `walChannel` buffers non-blockingly (`select ... default`).
Assessment: PASS
Severity: LOW
Notes: Monotonicity guaranteed by `atomic.Uint64`. Channel buffer of 1024 prevents writer starvation.

## Finding 3

Location: `internal/cluster/cluster.go:148-179` and `internal/cluster/cluster.go:252-263`
Claimed Behavior: Background replication workers process WAL entries with configurable simulated lag and terminate gracefully on cluster close.
Observed Implementation: `replicaWorker` listens on `c.ctx.Done()` and `replica.walChannel`. Lag delay is handled via `time.After` checking `c.ctx.Done()`. `Cluster.Close()` triggers `c.cancel()` and waits on `c.wg.Wait()`.
Assessment: PASS
Severity: LOW
Notes: Cluster shutdown cleanly releases all goroutines.

## Finding 4

Location: `internal/router/router.go:53-65`
Claimed Behavior: Session write records last write timestamp and LSN in thread-safe map.
Observed Implementation: `r.sessions.Store(sessionID, SessionState{LastWriteTime: time.Now(), LastWriteLSN: lsn})` uses standard library `sync.Map`.
Assessment: PASS
Severity: LOW
Notes: Concurrency-safe without lock contention across sessions.

## Finding 5

Location: `internal/router/router.go:80-90`
Claimed Behavior: Reads within `StickyDuration` of a write route to primary; after expiration, reads route to replicas via lag-aware routing.
Observed Implementation: Checks `time.Since(state.LastWriteTime) < r.config.StickyDuration`. If true, queries `r.cluster.Primary().Read(key)`; if false or session unknown, calls `r.ReadLagAware(key)`.
Assessment: PASS
Severity: LOW
Notes: Directly mirrors the design and research recommendations.

## Finding 6

Location: `internal/router/router.go:92-115`
Claimed Behavior: `ReadWithToken` guarantees Read-Your-Own-Writes by checking replica LSN, waiting with timeout if not caught up, and falling back to primary.
Observed Implementation: Iterates over replicas; if any replica has `AppliedLSN() >= minLSN`, reads immediately. Otherwise waits on replica 0 up to `WaitTimeout`. If timeout or wait failure, routes to primary with `(fallback)` designation.
Assessment: PASS
Severity: LOW
Notes: Satisfies causal consistency SLA without compromising primary protection.

## Finding 7

Location: `internal/router/router.go:117-141`
Claimed Behavior: `ReadLagAware` filters replicas where `primaryLSN - replicaLSN <= MaxLSNDiff`. If no replica meets SLA, falls back to primary.
Observed Implementation: Reads `primaryLSN = r.cluster.CurrentLSN()`. Computes difference with `replica.AppliedLSN()`. Replicas within `MaxLSNDiff` are collected in `eligible` slice and round-robined via atomic counter. If empty, routes to primary.
Assessment: PASS
Severity: LOW
Notes: Correctly handles dynamic replica exclusion and failover.
