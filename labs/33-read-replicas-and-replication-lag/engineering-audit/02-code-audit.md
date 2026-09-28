# Code Audit

## Finding 1

Location: `internal/cluster/cluster.go:79-96`
Claimed Behavior: `WaitForLSN` blocks until target LSN is reached or context times out.
Observed Implementation: Uses `sync.Cond` within a separate goroutine. If context is cancelled before LSN is reached, the goroutine remains blocked on `n.cond.Wait()` until the next `Broadcast()` or node update.
Assessment: WARNING
Severity: LOW
Notes: Minor background goroutine leak potential if a replica never receives further WAL updates; acceptable in simulation test lab boundaries.

## Finding 2

Location: `internal/cluster/cluster.go:220-226`
Claimed Behavior: Asynchronous streaming of WAL entries to replica channels.
Observed Implementation: Uses non-blocking channel send `select { case replica.walChannel <- entry: default: }` with buffer size 1024. If buffer overflows, entries are dropped silently without resync.
Assessment: PASS
Severity: LOW
Notes: Safe for expected simulation throughput under test parameters.

## Finding 3

Location: `internal/cluster/cluster.go:137-168`
Claimed Behavior: Replicas consume WAL stream with configurable artificial lag and update applied LSN.
Observed Implementation: `replicaWorker` respects context cancellation during lag sleep and correctly signals waiting goroutines via `replica.cond.Broadcast()`. Mutex protection around state and `appliedLSN` is strictly maintained.
Assessment: PASS
Severity: LOW
Notes: Concurrency safety verified.

## Finding 4

Location: `internal/router/router.go:53-65`
Claimed Behavior: Session write tracking for sticky routing.
Observed Implementation: Thread-safe `sync.Map` stores session timestamp and LSN atomically upon write.
Assessment: PASS
Severity: LOW
Notes: Clean thread-safe implementation.

## Finding 5

Location: `internal/router/router.go:80-90`
Claimed Behavior: Sticky read routes to primary within duration, falls back to lag-aware read afterward.
Observed Implementation: Evaluates `time.Since(state.LastWriteTime) < r.config.StickyDuration`. If valid, routes directly to primary. If expired or unknown session, delegates to `ReadLagAware`.
Assessment: PASS
Severity: LOW
Notes: Accurately mirrors claimed design.

## Finding 6

Location: `internal/router/router.go:117-141`
Claimed Behavior: Dynamic lag-aware replica routing with primary fallback.
Observed Implementation: Compares `primaryLSN - appliedLSN <= MaxLSNDiff`. If all replicas exceed lag, safely routes to primary with `(fallback-lag)` marker.
Assessment: PASS
Severity: LOW
Notes: Clean boundary checking and fallback semantics.
