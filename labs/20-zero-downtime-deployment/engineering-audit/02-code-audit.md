# Code Audit

Target Lab: `labs/20-zero-downtime-deployment`

## Finding 1: Correctness of Server Lifecycle and PreStop Hook

Location: `internal/server/server.go:85-110`
Claimed Behavior: Server marks readiness probe unready (`503`), executes configurable `preStop` delay to simulate routing detachment, then invokes graceful shutdown waiting for active in-flight requests.
Observed Implementation:
- `s.SetReady(false)` is invoked immediately upon `Shutdown()`.
- If `s.preStop > 0`, it waits on `select` between `time.After(s.preStop)` and `ctx.Done()`, avoiding unbounded blocking on cancelled context.
- Calls standard library `s.srv.Shutdown(ctx)`, followed by `s.wg.Wait()` on in-flight requests.
- Active requests are tracked with `atomic.Int32` and `sync.WaitGroup`.
Assessment: PASS
Severity: LOW
Notes: Implementation matches design specifications.

---

## Finding 2: Safe Concurrency and Teardown in Worker Queue

Location: `internal/worker/worker.go:76-107`
Claimed Behavior: Background queue worker accepts jobs, runs workers concurrently, and allows in-flight jobs to complete gracefully during shutdown while dropping jobs enqueued after shutdown.
Observed Implementation:
- `Enqueue()` protects channel send with `enqueueMu` and checks `w.stopped.Load()`.
- `Stop()` acquires `enqueueMu`, flags `stopped = true`, closes `jobChan`, and releases lock before waiting on workers.
- Channel closure guarantees no deadlock or panic on enqueue attempt after stop.
- Worker goroutines check `w.ctx.Done()` both before and after job channel dequeues to prevent stale job dispatch if timeout triggers.
- Protected slice append in `w.completed` with `completedMu`.
Assessment: PASS
Severity: LOW
Notes: Tested under heavy concurrent enqueue/stop loops with Go race detector enabled; no races detected.

---

## Finding 3: Database Expand and Contract Pattern Compatibility

Location: `internal/db/db.go:41-71`
Claimed Behavior: Database store supports legacy single-name format and modern first/last name format, dual-writing during the Expand phase and transparently falling back on reads.
Observed Implementation:
- `InsertLegacy` writes `Name`.
- `SaveExpand` writes `FirstName`, `LastName`, and computes dual-written `Name = strings.TrimSpace(firstName + " " + lastName)`.
- `GetUser` checks if `FirstName == "" && LastName == ""` with non-empty `Name`, splitting via `strings.SplitN(rec.Name, " ", 2)` to populate legacy reads into new format.
- Fallback checks empty `Name` when first/last names exist.
- Thread-safe access using `sync.RWMutex`.
Assessment: PASS
Severity: LOW
Notes: Edge cases (single name without space, empty strings) properly covered and tested.

---

## Finding 4: In-Flight HTTP Request Abort Handling

Location: `internal/server/server.go:43-63`
Claimed Behavior: When a client terminates early during `/work`, server handles context cancellation cleanly without leaking active request counts.
Observed Implementation:
- Uses `s.wg.Add(1)` / `defer s.wg.Done()` and `s.activeCount.Add(1)` / `defer s.activeCount.Add(-1)`.
- Listens on `r.Context().Done()` alongside `time.After(d)`.
- Upon context cancellation, exits handler promptly; `defer` blocks decrement counter and mark WaitGroup done.
Assessment: PASS
Severity: LOW
Notes: Verified by `TestServerWorkRequestCancellation`.
