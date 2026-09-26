# Code Audit

## Finding 1

Location: internal/db/db.go: GetUser, InsertLegacy, SaveExpand
Claimed Behavior: Expand and Contract pattern with backward-compatible reads/writes. Legacy records store only `Name`. Modern records store `FirstName`/`LastName`. Reading a legacy record derives `FirstName`/`LastName` from `Name`. Reading a modern record derives `Name` from `FirstName`+`LastName`.
Observed Implementation: `InsertLegacy` writes `UserRecord{ID, Name}`. `SaveExpand` writes all four fields with `Name` computed as concatenation. `GetUser` uses a copy of the map value (value semantics) so mutations don't leak into the store. Fallback branches handle: (a) legacy-only record splitting Name into FirstName/LastName; (b) modern-only record computing Name from FirstName+LastName. No-op when both field groups are empty.
Assessment: PASS
Severity: LOW
Notes: `SplitN(" ", 2)` correctly handles single-word names (LastName stays empty). `rec` is a struct copy from the map, so in-place mutation in GetUser does not corrupt the store. Correct value semantics.

## Finding 2

Location: internal/db/db.go
Claimed Behavior: In-memory store must be concurrency-safe.
Observed Implementation: `UserStore` guards `records` map with `sync.RWMutex`. `InsertLegacy` and `SaveExpand` use write lock. `GetUser` uses read lock. No unprotected map access.
Assessment: PASS
Severity: N/A
Notes: RWMutex is appropriate for read-heavy workload. All three mutating/reading methods acquire the lock before touching the map.

## Finding 3

Location: internal/db/db.go
Claimed Behavior: `SaveExpand` writes dual schema (both `Name` and `FirstName`/`LastName`).
Observed Implementation: `SaveExpand` sets `Name` to the trimmed concatenation of `firstName+" "+lastName` and stores `FirstName`/`LastName` individually. `Name` is always a derived value, never independently set on a modern record.
Assessment: PASS
Severity: N/A
Notes: This represents the "Expand" phase where both old and new schemas coexist. The `Name` field is populated for backward compatibility so legacy readers can still read the record.

## Finding 4

Location: internal/server/server.go: Server struct, NewServer, handlers
Claimed Behavior: Liveness probe always 200 OK; Readiness probe returns 200 when ready, 503 otherwise.
Observed Implementation: `/healthz/live` always returns 200 with body "OK". `/healthz/ready` checks `s.ready` (atomic.Bool) — 200/503 with "READY"/"NOT READY". The `/work` endpoint tracks active requests via `atomic.Int32` and `sync.WaitGroup`.
Assessment: PASS
Severity: N/A
Notes: atomic.Bool and atomic.Int32 provide lock-free state access for probes. Work endpoint increments/decrements activeCount and WaitGroup with deferred cleanup.

## Finding 5

Location: internal/server/server.go: Shutdown method
Claimed Behavior: On shutdown, mark unready, execute preStop delay, gracefully close listeners, wait for in-flight requests.
Observed Implementation: Shutdown sets ready=false first (traffic detachment), then sleeps for `preStop` duration if >0 (interruptible by context), then calls `s.srv.Shutdown(ctx)`, then `s.wg.Wait()`. Ordering matches the intended traffic-draining sequence.
Assessment: PASS
Severity: N/A
Notes: preStop sleep uses `select` on `time.After` and `ctx.Done()` — correctly aborts on context cancellation. `s.srv.Shutdown(ctx)` blocks until active HTTP connections finish, then `s.wg.Wait()` ensures handler goroutines complete. This two-layer drain is correct.

## Finding 6

Location: internal/server/server.go: Shutdown, preStop vs. ctx
Claimed Behavior: preStop delay should not outlive the context deadline.
Observed Implementation: The preStop `select` races `time.After(s.preStop)` against `ctx.Done()`. If context expires first, returns `ctx.Err()` immediately.
Assessment: PASS
Severity: N/A
Notes: Verified by `TestServerPreStopContextCancellation` — confirms early abort when context times out during preStop.

## Finding 7

Location: internal/server/server.go: /work handler
Claimed Behavior: In-flight requests complete during graceful shutdown; invalid duration falls back to 50ms default.
Observed Implementation: `d` query param parsed via `time.ParseDuration`; on error defaults to `50 * time.Millisecond`. Handler completes via `time.After(d)` or returns early on `r.Context().Done()` (client disconnect).
Assessment: PASS
Severity: LOW
Notes: `TestServerInvalidDurationFallback` and `TestServerWorkRequestCancellation` verify both behaviors. The handler correctly handles client-side context cancellation, decrementing activeCount via defer.

## Finding 8

Location: internal/server/server.go: Shutdown method
Claimed Behavior: After `s.srv.Shutdown(ctx)` returns all active connections closed, `s.wg.Wait()` ensures handler goroutines fully exit.
Observed Implementation: `s.wg.Wait()` is called after `s.srv.Shutdown`. The WaitGroup is incremented at handler entry and decremented via deferred `Done()`. This is a safety net to ensure handler-side cleanup completes even after the HTTP server reports shutdown complete.
Assessment: PASS
Severity: N/A
Notes: While `http.Server.Shutdown` already waits for handlers, the redundant `wg.Wait()` provides defense-in-depth for the activeCount tracking. No deadlock risk because wg.Done is deferred in every handler path.

## Finding 9

Location: internal/worker/worker.go: Worker struct, Start, Enqueue, Stop
Claimed Behavior: Worker processes queued jobs with concurrency goroutines. Enqueue is safe under concurrent use. Stop drains in-flight jobs and then exits.
Observed Implementation: `Start` spawns `concurrency` goroutines that loop on `select` between `w.ctx.Done()` and `w.jobChan`. Before executing a dequeued job, it re-checks `ctx.Done()` to avoid starting new work during shutdown. Job execution uses `select` on `time.After(job.Duration)` and `w.ctx.Done()`.
Assessment: PASS
Severity: N/A
Notes: The double context check (outer select + inner select before execution) prevents starting jobs when shutdown is in progress. Context cancellation also aborts in-flight long jobs via the second select.

## Finding 10

Location: internal/worker/worker.go: Stop method
Claimed Behavior: Stop closes the job channel, waits for active workers (with timeout), cancels context if timeout exceeded.
Observed Implementation: `Stop` acquires `enqueueMu`, sets `stopped=true`, `close(w.jobChan)`, unlocks. Spawns a goroutine that waits on `w.wg` and signals `done`. Main select races `done` vs `time.After(timeout)`. On timeout, calls `w.cancel()` which cancels `ctx`, then waits on `done`.
Assessment: PASS
Severity: MEDIUM
Notes: The timeout escalation from cooperative drain to forced cancellation is correct. However, calling `Stop` twice will panic due to double `close(w.jobChan)`. No guard prevents re-entry. This is not exercised by tests or demo but is a latent bug. Documented as potential issue; severity downgraded to MEDIUM since not triggered in current usage.

## Finding 11

Location: internal/worker/worker.go: Enqueue method
Claimed Behavior: Concurrent Enqueue calls are safe; Enqueue after Stop is rejected.
Observed Implementation: `Enqueue` holds `enqueueMu` while checking `w.stopped.Load()` and sending to `w.jobChan`. Since `Stop` also acquires `enqueueMu` before `close(w.jobChan)`, there is no TOCTOU window between the stopped check and the channel send. Verified by `TestWorkerConcurrentEnqueueStop`.
Assessment: PASS
Severity: N/A
Notes: The `enqueueMu` mutex serializes the stopped-check and channel-close/stop, eliminating the race. Correct concurrency design.

## Finding 12

Location: internal/worker/worker.go: Start worker loop
Claimed Behavior: Workers exit when either the context is canceled or the channel is closed.
Observed Implementation: Outer `select` has `case <-w.ctx.Done(): return` and `case job, ok := <-w.jobChan:` with `if !ok { return }`. When `Stop` closes the channel, all worker goroutines exit. If `Stop` hits the timeout and cancels the context, workers also exit via `ctx.Done()`.
Assessment: PASS
Severity: N/A
Notes: Both exit paths are covered. Channel closure causes blocked `<-w.jobChan` to return immediately with `ok=false`.

## Finding 13

Location: internal/server/server.go: Start method
Claimed Behavior: Server starts listening and returns nil on graceful close.
Observed Implementation: `s.srv.ListenAndServe()` returns `http.ErrServerClosed` on graceful shutdown, which is treated as success (returns nil). Any other error is returned.
Assessment: PASS
Severity: N/A
Notes: Correct handling of the standard graceful shutdown signal.

## Finding 14

Location: cmd/demo/main.go
Claimed Behavior: Orchestrates worker + server lifecycle; demonstrates readiness, in-flight work, SIGTERM simulation, graceful shutdown, and worker drain.
Observed Implementation: Creates worker (2 goroutines, buffer 100), enqueues DemoJob-1 (2s), starts server with 1s preStop, waits 1s init, sets ready, starts a `/work?d=2s` goroutine, waits 500ms, simulates SIGTERM via self-signal, calls `srv.Shutdown(ctx 10s)`, then `w.Stop(5s)`.
Assessment: PASS
Severity: N/A
Notes: Demo execution output matches the logged sequence. All components coordinate correctly. The 2s worker job and 2s work request are designed to span the shutdown boundary, demonstrating drain behavior.

## Finding 15

Location: internal/server/server.go: Shutdown — preStop sleep before listener close
Claimed Behavior: preStop delay should simulate load balancer detachment latency before the server stops accepting connections.
Observed Implementation: preStop sleep occurs BEFORE `s.srv.Shutdown()`. During preStop, the server still accepts connections. The `ready` flag is set to false before preStop so the readiness probe returns 503, signaling load balancers to detach, but the listener stays open during the sleep.
Assessment: PASS
Severity: N/A
Notes: This matches the Kubernetes preStop hook model: readiness is withdrawn first, then the preStop delay allows routing tables to propagate, then the listener actually closes. Correct sequencing.