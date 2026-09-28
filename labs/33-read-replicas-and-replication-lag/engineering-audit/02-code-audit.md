# Code Audit

Target Lab: `labs/33-read-replicas-and-replication-lag`

## Finding 1

Location: `internal/cluster/cluster.go:79-96` (`Node.WaitForLSN`)
Claimed Behavior: Blocks until replica applies target LSN or context times out.
Observed Implementation: Launches a goroutine that locks `n.mu`, loops checking `n.appliedLSN < targetLSN` via `n.cond.Wait()`, and closes a channel when done. A `select` block waits on context completion or channel closure.
Assessment: WARNING
Severity: LOW
Notes: `n.cond.Wait()` is called inside a separate goroutine launched by `WaitForLSN`. If `ctx` times out early, the launcher goroutine remains blocked on `n.cond.Wait()` until the next broadcast on `n.cond` (or cluster shutdown). In practice for this lab, writes trigger `Broadcast()` frequently so worker goroutines finish quickly, but goroutine cleanup is not immediate on context timeout.

## Finding 2

Location: `internal/cluster/cluster.go:220-225` (`Cluster.Write` async branch)
Claimed Behavior: Streams WAL entries asynchronously to replica channels.
Observed Implementation: Uses non-blocking select send (`select { case replica.walChannel <- entry: default: }`).
Assessment: WARNING
Severity: LOW
Notes: If the buffer of `walChannel` (capacity 1024) fills up, WAL entries are dropped silently without alerting or retrying. Sufficient for in-memory lab bounds, but under extreme backpressure WAL records would drop.

## Finding 3

Location: `internal/cluster/cluster.go:203-218` (`Cluster.Write` sync branch)
Claimed Behavior: Synchronous replication waits for replica WAL application before write returns.
Observed Implementation: Directly applies write data and `appliedLSN` to all replicas under lock with optional delay sleep, simulating `remote_apply`.
Assessment: PASS
Severity: LOW
Notes: Correctly enforces data visibility on all replicas before returning from `Write()`.

## Finding 4

Location: `internal/router/router.go:80-90` (`Router.ReadWithStickySession`)
Claimed Behavior: Routes reads to primary within `StickyDuration` after write, then delegates to lag-aware routing.
Observed Implementation: Loads session state from `sync.Map`; if `time.Since(lastWrite) < StickyDuration`, routes read to primary. Otherwise delegates to `ReadLagAware`.
Assessment: PASS
Severity: LOW
Notes: Correctly implements time-based sticky routing to satisfy read-your-own-writes.

## Finding 5

Location: `internal/router/router.go:92-115` (`Router.ReadWithToken`)
Claimed Behavior: Ensures replica read freshness against `minLSN`, falling back to waiting or primary read.
Observed Implementation: First checks all replicas for `AppliedLSN() >= minLSN`. If none found, blocks on replica 0 using `WaitForLSN` with `WaitTimeout`. On timeout or failure, falls back to primary.
Assessment: PASS
Severity: LOW
Notes: Correctly guarantees read freshness for causal tokens while maintaining primary fallback resiliency.

## Finding 6

Location: `internal/router/router.go:117-140` (`Router.ReadLagAware`)
Claimed Behavior: Excludes replicas whose lag exceeds `MaxLSNDiff` and falls back to primary if no replica is eligible.
Observed Implementation: Calculates `primaryLSN - appliedLSN` for each replica. Filters out replicas exceeding `MaxLSNDiff`. If empty, returns primary with `(fallback-lag)` tag.
Assessment: PASS
Severity: LOW
Notes: Correctly implements dynamic lag filtering and primary failover.
