# Code Audit

## Finding 1

Location: internal/server/server.go:43-47, 85-110
Claimed Behavior: Graceful shutdown drains in-flight `/work` requests after preStop delay.
Observed Implementation: `Shutdown` sets unready, sleeps `preStop` (ctx-aware), calls `http.Server.Shutdown(ctx)`, then `wg.Wait()`. Handler tracks `wg` + `activeCount`.
Assessment: PASS
Severity: LOW
Notes: `wg.Add` inside handler + `Wait` after `srv.Shutdown` redundant but safe; `-race` clean.

## Finding 2

Location: internal/server/server.go:91-99
Claimed Behavior: preStop waits configured delay for routing detachment.
Observed Implementation: `select time.After(preStop)` vs `ctx.Done()`, returns `ctx.Err()` on cancel.
Assessment: PASS
Severity: LOW
Notes: Covered by `TestServerPreStopHook` and `TestServerPreStopContextCancellation`.

## Finding 3

Location: internal/server/server.go:28-42
Claimed Behavior: Liveness always 200; readiness gates traffic.
Observed Implementation: `/healthz/live` static 200; `/healthz/ready` checks `atomic.Bool`, 503 when false.
Assessment: PASS
Severity: LOW
Notes: Matches design; transitions tested.

## Finding 4

Location: internal/worker/worker.go:35-64, 70-88
Claimed Behavior: Stop drains active + buffered jobs, timeout abandons queued jobs.
Observed Implementation: `close(jobChan)`, `wg.Wait()` with timeout ceiling, `cancel()` on timeout. Active job `time.Sleep` uninterruptible.
Assessment: PASS
Severity: LOW
Notes: Timeout does not preempt active sleep, only stops pickup of next job. Documented ceiling, matches tests.

## Finding 5

Location: internal/worker/worker.go:66-68, 70-72
Claimed Behavior: Cooperative termination.
Observed Implementation: `Enqueue` blocking send; `Stop` unconditional `close(jobChan)`.
Assessment: WARNING
Severity: MEDIUM
Notes: `Enqueue` after `Stop` panics (send on closed). Double `Stop` panics (close of closed). No guard/idempotency.

## Finding 6

Location: internal/db/db.go:32-72
Claimed Behavior: Expand/contract dual-write + fallback read.
Observed Implementation: `InsertLegacy` writes `Name` only; `SaveExpand` dual-writes `Name`+`First/Last`; `GetUser` splits legacy or reconstructs modern. `RWMutex` guarded.
Assessment: PASS
Severity: LOW
Notes: Single-word legacy name yields empty LastName, acceptable. Overwrite without error acceptable for lab.

## Finding 7

Location: cmd/demo/main.go:1-78
Claimed Behavior: Demonstrates ready, in-flight request, SIGTERM, preStop, drain.
Observed Implementation: Fixed `127.0.0.1:8080`, 1s init, 1s preStop, self-sent SIGTERM, 10s shutdown ctx, 5s worker stop. Real execution verified.
Assessment: PASS
Severity: LOW
Notes: Fixed port, synthetic signal acceptable for demo scope.
