# Code Audit Findings

## Finding 1

Location: `internal/server/server.go:85-110`
Claimed Behavior: HTTP server executes preStop delay, switches readiness probe to unready, drains in-flight requests, and shuts down listeners cleanly.
Observed Implementation:
- `s.SetReady(false)` marks readiness atomic flag false.
- `preStop` timer executes with context cancellation awareness (`select { case <-time.After(s.preStop): case <-ctx.Done(): ... }`).
- `s.srv.Shutdown(ctx)` closes listeners and stops accepting new connections.
- `s.wg.Wait()` blocks until all active requests complete.
Assessment: PASS
Severity: LOW
Notes: Concurrency safety ensured via `sync.WaitGroup` and `atomic.Int32`/`atomic.Bool`.

## Finding 2

Location: `internal/worker/worker.go:38-74`, `internal/worker/worker.go:86-107`
Claimed Behavior: Worker consumes queued jobs, stops accepting jobs upon shutdown, drains active in-flight jobs within timeout, and cancels context if timeout is exceeded.
Observed Implementation:
- `w.stopped` atomic flag and `w.enqueueMu` prevent race conditions on channel closure and late enqueue attempts.
- Worker goroutines monitor both `jobChan` and `w.ctx.Done()`.
- `Stop(timeout)` closes the channel, waits on `w.wg`, and enforces the drain timeout via select/cancel fallback.
- `GetCompletedJobs()` uses `w.completedMu` to safely return slice snapshots.
Assessment: PASS
Severity: LOW
Notes: No race conditions detected under `go test -race`.

## Finding 3

Location: `internal/db/db.go:41-71`
Claimed Behavior: Implements expand-and-contract pattern with dual write (`SaveExpand`) and backward compatible read fallback (`GetUser`).
Observed Implementation:
- `SaveExpand` writes both legacy `Name` and split `FirstName`/`LastName` fields.
- `GetUser` falls back to splitting `Name` if `FirstName`/`LastName` are absent, and constructs `Name` if only split fields exist.
- Synchronized with `sync.RWMutex`.
Assessment: PASS
Severity: LOW
Notes: All access synchronized; safe for concurrent reads and writes.
