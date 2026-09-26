## Finding 1

Location: `internal/db/db.go`
Claimed Behavior: Expand and Contract pattern — storage supports writing dual schema versions (legacy `Name` and modern `FirstName`/`LastName`) and transparent fallback logic when reading records created by differing versions.
Observed Implementation: `SaveExpand(id, firstName, lastName)` writes both `Name` (combined) and `FirstName`/`LastName`. `InsertLegacy(id, name)` writes only `Name`. `GetUser(id)` has fallback: if `FirstName` == "" && `LastName` == "" && `Name` != "" → derive modern fields from `Name` via `strings.SplitN`; else if `Name` == "" → reconstruct `Name` from `FirstName`/`LastName`. Concurrency: protected by `sync.RWMutex`.
Assessment: PASS
Severity: N/A
Notes: Matches claim exactly. Tests `TestDBLegacyOverwriteWithExpand` and `TestExpandContractDatabase` verify transition and fallback.

## Finding 2

Location: `internal/server/server.go`
Claimed Behavior: Exposes Liveness (`/healthz/live`) and Readiness (`/healthz/ready`) probes. When shutdown signal sent, executes configurable `preStop` delay, then performs graceful shutdown ensuring any in-flight requests complete before termination.
Observed Implementation: `/healthz/live` returns 200 OK; `/healthz/ready` returns 200 if `ready` atomic.Bool true else 503. `Shutdown(ctx)`: sets ready false → executes `preStop` sleep → `s.srv.Shutdown(ctx)` (drains in-flight via http.Server) → `s.wg.Wait()` (redundant safety). `preStop` is constructor parameter.
Assessment: PASS
Severity: N/A
Notes: Tests `TestServerProbes`, `TestServerGracefulShutdown`, `TestServerPreStopHook` verify probes, in-flight completion, preStop delay. Race detector passes.

## Finding 3

Location: `internal/server/server.go`
Claimed Behavior: Graceful shutdown ensures any in-flight requests complete before termination.
Observed Implementation: `Shutdown(ctx)` returns early via `return ctx.Err()` if context expires during `preStop` sleep — *without* calling `s.srv.Shutdown(ctx)` or draining in-flight requests. Listeners remain open; goroutine leak.
Observed Behavior: Test `TestServerPreStopContextCancellation` enforces early return (checks elapsed < preStopDuration) and expects error.
Assessment: WARNING
Severity: LOW
Notes: Design "ensures ... before termination" implies within shutdown context. If context expires, caller signaled unwillingness to wait; early return idiomatic Go. However, leaves HTTP listener + `http.Server.Start()` goroutine leaked until process exit. For lab (process short-lived) acceptable but worth noting. Fix: On preStop context expiry, still attempt `s.srv.Shutdown()` with remaining context or forceful close.

## Finding 4

Location: `internal/worker/worker.go`
Claimed Behavior: Upon receiving termination signal, worker stops pulling new jobs but continues processing the current active job until completion.
Observed Implementation: `Stop(timeout)`: under `enqueueMu`, sets `stopped=true`, closes `jobChan`. Workers on `select { case <-w.ctx.Done(): case job, ok := <-w.jobChan: }` return if `!ok` (closed channel) or `ctx.Done()`. In-flight jobs continue via `select { case <-time.After(job.Duration): ... case <-w.ctx.Done(): return }`. Timeout cancels context to abort in-progress job after `timeout`. Buffered jobs drain (closed chan still yields buffered items until empty). Concurrency: `enqueueMu` guards `Enqueue`/`Stop` TOCTOU; `completedMu` guards `completed` slice append/read.
Assessment: PASS
Severity: N/A
Notes: Tests `TestWorkerGracefulShutdown`, `TestWorkerShutdownTimeout`, `TestWorkerEnqueueAfterStop` verify drain, in-flight completion, timeout abort, post-stop rejection.

## Finding 5

Location: `internal/worker/worker.go`
Claimed Behavior: Worker stops pulling new jobs but continues processing the current active job until completion (implicit safety on repeated Stop).
Observed Implementation: `Stop()` calls `close(w.jobChan)` without guarding against double-close. Calling `Stop` twice panics (double close of closed channel).
Observed Behavior: Not exercised by tests or demo; single Stop only.
Assessment: WARNING
Severity: LOW
Notes: Minor robustness gap. Fix: Add `if !w.stopped.CompareAndSwap(false, true) { return }` or check before close.

## Finding 6

Location: `internal/worker/worker.go`
ClaimedBehavior: Safe concurrent Enqueue and Stop.
Observed Implementation: `Enqueue(job)` holds `enqueueMu` while checking `stopped` and sending to `jobChan`. `Stop()` holds `enqueueMu` while setting `stopped` and closing `jobChan`. Prevents TOCTOU but risks deadlock if Enqueue blocked on full `jobChan` (buffer full, workers draining slowly): Enqueue holds `enqueueMu`; Stop blocks on `enqueueMu`; workers cannot drain to unblock Enqueue → deadlock.
Observed Behavior: Not exercised by tests or demo; requires specific timing/load.
Assessment: WARNING
Severity: LOW
Notes: Theoretical deadlock edge case. Fix: Use non-blocking select on `jobChan` in Enqueue, or decouple mutex from channel send (advanced).

## Finding 7

Location: `internal/server/server.go`
ClaimedBehavior: Active request tracking accurate during shutdown/cancelation.
Observed Implementation: `/work` handler: `s.activeCount.Add(1); s.wg.Add(1); defer s.wg.Done(); defer s.activeCount.Add(-1);`. Test `TestServerWorkRequestCancellation` verifies `ActiveRequests() == 0` after client context cancellation.
ObservedBehavior: Defers run LIFO → `activeCount` decremented before `wg.Done()`. Test sleeps 50ms (enough for handler return) then checks; passes.
Assessment: PASS
Severity: N/A
Notes: Ordering of defers does not affect correctness for this claim. Race detector passes.