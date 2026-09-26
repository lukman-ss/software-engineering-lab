# Code Audit

Target Lab: labs/20-zero-downtime-deployment

---

## Finding 1

Location: internal/server/server.go:43-63 (`/work` handler)
Claimed Behavior: In-flight requests complete before server terminates; client disconnect is handled cleanly.
Observed Implementation: Handler increments `activeCount` and `wg` atomically at entry. Uses `select` between `time.After(d)` and `r.Context().Done()` for cooperative cancellation. Decrements both on defer. `wg.Wait()` is called after `srv.Shutdown(ctx)`.
Assessment: PASS
Severity: LOW
Notes: `s.wg.Wait()` after `s.srv.Shutdown(ctx)` is technically redundant because `http.Server.Shutdown` already waits for all active handlers to return before returning. The extra `wg.Wait()` provides a belt-and-suspenders guarantee and enables `ActiveRequests()` tracking, but it does not cause incorrectness. Harmless duplication.

---

## Finding 2

Location: internal/server/server.go:85-109 (Shutdown)
Claimed Behavior: Shutdown marks server unready, waits preStop delay, then closes listeners. Context cancellation during preStop aborts early and returns error.
Observed Implementation: `SetReady(false)` called first. Then `select` on `time.After(preStop)` vs `ctx.Done()`. On context cancellation: returns `ctx.Err()` immediately without calling `srv.Shutdown`. On normal path: calls `srv.Shutdown(ctx)`, then `wg.Wait()`.
Assessment: PASS
Severity: LOW
Notes: Context cancellation during preStop short-circuits before `srv.Shutdown` is called, meaning active connections are not drained on hard timeout during preStop phase. This matches the ponytail comment intent and is semantically correct for a hard-abort scenario. Not a defect.

---

## Finding 3

Location: internal/worker/worker.go:68-73 (Enqueue)
Claimed Behavior: Enqueue rejects jobs after Stop() is called.
Observed Implementation: `stopped.Load()` check before channel send. However, `Stop()` calls `close(jobChan)` after setting `stopped=true`. If `Enqueue` passes the `stopped` check and then `Stop()` closes the channel before the send, a panic on send-to-closed-channel occurs.
Assessment: WARNING
Severity: MEDIUM
Notes: There is a TOCTOU race between `Enqueue`'s `stopped.Load()` check and the `close(jobChan)` in `Stop()`. In practice, the demo and tests do not exercise concurrent Enqueue+Stop, so no panic occurs during testing. Under concurrent use (multiple goroutines calling Enqueue while Stop is called concurrently), a send-on-closed-channel panic is possible. The test suite does not cover this scenario.

---

## Finding 4

Location: internal/worker/worker.go:37-66 (Start / goroutine loop)
Claimed Behavior: Worker goroutines process jobs cooperatively and stop on context cancellation or channel close.
Observed Implementation: Double-select pattern: first checks `ctx.Done()` non-blocking, then blocks on either `ctx.Done()` or `jobChan`. After `Stop()` closes the channel, goroutines drain remaining buffered jobs before exiting because channel-close does not cancel the context (unless drain timeout triggers `w.cancel()`). Jobs in the channel buffer are processed, not dropped, on normal Stop.
Assessment: PASS
Severity: LOW
Notes: The drain behavior (process buffered jobs on graceful stop) is correct. Context cancellation only occurs on timeout path, which forces goroutines to exit the inner select via `ctx.Done()` — even mid-`time.Sleep` the goroutine will exit the select but the `time.Sleep` in the job body is NOT preemptible. The goroutine exits the select on ctx.Done() only between jobs, not during a running job's sleep. The currently running job will complete its `time.Sleep` regardless. This means drain timeout cancels future job pickup but cannot interrupt the currently-executing job's sleep. This is documented behavior (ponytail comment) and aligns with the claim "finishes active job".

---

## Finding 5

Location: internal/worker/worker.go:76-95 (Stop)
Claimed Behavior: Stop blocks until drain completes or timeout elapses; on timeout, context is cancelled and Stop still waits for goroutines to exit.
Observed Implementation: `close(jobChan)` → `wg.Wait()` in background goroutine → `done` channel. Select between `done` and `time.After(timeout)`. On timeout: `w.cancel()` then `<-done`. This correctly waits for goroutines to exit even after context cancel.
Assessment: PASS
Severity: LOW
Notes: Correct. Stop never returns while goroutines are still running.

---

## Finding 6

Location: internal/db/db.go:61-69 (GetUser fallback logic)
Claimed Behavior: Legacy records (Name only, no FirstName/LastName) are returned with FirstName/LastName split from Name. Modern records (with FirstName/LastName) return combined Name if Name is empty.
Observed Implementation: If `FirstName==""` AND `LastName==""` AND `Name!=""`: splits Name on first space into FirstName/LastName. If `Name==""`: constructs Name from FirstName+LastName. The check `rec.Name==""` in the else-if is only reached when at least one of FirstName/LastName is non-empty (due to prior if condition).
Assessment: PASS
Severity: LOW
Notes: Edge case: single-name legacy record (e.g., "Madonna") correctly sets FirstName="Madonna", LastName="". Covered by TestDBSingleNameLegacy. Empty firstName with lastName-only also covered. Logic is correct.

---

## Finding 7

Location: internal/db/db.go:41-51 (SaveExpand)
Claimed Behavior: SaveExpand writes all three fields (Name, FirstName, LastName).
Observed Implementation: `Name` is set to `strings.TrimSpace(firstName + " " + lastName)`. When firstName="", Name becomes the trimmed lastName (no leading space). When lastName="", Name becomes trimmed firstName (no trailing space).
Assessment: PASS
Severity: LOW
Notes: TrimSpace correctly handles empty component cases. Verified by TestDBSaveExpandEmptyFields.

---

## Finding 8

Location: cmd/demo/main.go:55-63 (SIGTERM simulation)
Claimed Behavior: Demo simulates orchestrator SIGTERM signal to trigger graceful shutdown.
Observed Implementation: `signal.Notify(sig, ...)` then injects `syscall.SIGTERM` into the channel via goroutine. This exercises the shutdown code path but bypasses actual OS signal delivery. The real OS signal path (kernel → process → Go runtime → channel) is not tested.
Assessment: WARNING
Severity: LOW
Notes: Adequate for demo/illustration purposes. Not a fabricated result — the shutdown logic is genuinely exercised. Acceptable ceiling for a lab context.

---

## Finding 9

Location: tests/server_test.go (all tests using fixed ports 8081-8087)
Claimed Behavior: Tests are isolated and independently runnable.
Observed Implementation: Each test uses a hardcoded port. Go test runs tests within a package sequentially by default (no t.Parallel()), so port conflicts within the package are avoided. However, if tests in the package ran in parallel or if another process occupies these ports, tests would fail with bind errors.
Assessment: WARNING
Severity: LOW
Notes: Not a correctness issue in the current setup. Tests pass cleanly. Not a blocking concern.

---

## Finding 10

Location: tests/worker_test.go:47-65 (TestWorkerShutdownTimeout)
Claimed Behavior: Worker drain timeout causes context cancellation, job-slow completes (it was already running), job-dropped is abandoned.
Observed Implementation: 1 worker, job-slow (100ms), job-dropped (100ms). Stop(20ms timeout). job-slow starts immediately. After 20ms timeout, context cancelled. job-slow's `time.Sleep` is not preemptible → finishes. job-dropped was still in the channel; the goroutine exits because `ctx.Done()` fires between jobs. Test asserts exactly 1 completed job.
Assessment: PASS
Severity: LOW
Notes: Behavior is correctly modeled. The test passes and verifies the claimed drain-with-timeout semantics.
