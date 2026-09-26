# Code Audit

Target Lab: labs/20-zero-downtime-deployment

## Finding 1

Location: `internal/server/server.go:85-110`
Claimed Behavior: HTTP server unreadies first, executes preStop delay, shuts down listener, and waits for in-flight requests.
Observed Implementation: `s.SetReady(false)` is invoked synchronously before preStop sleep. PreStop select respects `ctx.Done()`. `s.srv.Shutdown(ctx)` closes listener while `s.wg.Wait()` ensures all in-flight `/work` handlers return.
Assessment: PASS
Severity: LOW
Notes: Correctly handles preStop timeout interruption and in-flight request tracking.

## Finding 2

Location: `internal/worker/worker.go:76-107`
Claimed Behavior: Background worker stops accepting new jobs upon stop signal while allowing currently executing jobs to finish gracefully or cancel on drain timeout.
Observed Implementation: `Enqueue` is guarded by `w.enqueueMu` and checks `w.stopped.Load()` before pushing to `jobChan`. `Stop` acquires `enqueueMu`, marks `stopped = true`, closes `jobChan`, and waits for `wg.Wait()`. If timeout expires, `w.cancel()` cancels context to abort long jobs.
Assessment: PASS
Severity: LOW
Notes: Synchronization between enqueue and close prevents panics on sending to a closed channel.

## Finding 3

Location: `internal/db/db.go:41-72`
Claimed Behavior: UserStore supports Expand & Contract pattern by dual-writing legacy and expanded fields and reading backward-compatibly with transparent fallbacks.
Observed Implementation: `SaveExpand` writes both `Name` and `FirstName`/`LastName`. `GetUser` uses `RWMutex` read locks and derives `FirstName`/`LastName` if missing from legacy single-field writes, and constructs `Name` if missing.
Assessment: PASS
Severity: LOW
Notes: Thread-safe read/write operations with proper fallback parsing.
