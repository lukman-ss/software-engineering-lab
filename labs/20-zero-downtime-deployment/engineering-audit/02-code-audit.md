# Code Audit

Target Lab: labs/20-zero-downtime-deployment

---

## Finding 1

Location: internal/server/server.go:43-63 (/work handler)
Claimed Behavior: Track active in-flight requests; complete them during graceful shutdown.
Observed Implementation:
```go
s.activeCount.Add(1)
s.wg.Add(1)
defer s.wg.Done()
defer s.activeCount.Add(-1)
```
`wg.Add(1)` is called inside the handler body, meaning the WaitGroup count is incremented only after the goroutine is scheduled and the handler begins executing. `http.Server.Shutdown()` stops accepting new connections and waits for active handlers to return via its own internal tracking. The manual `s.wg.Wait()` after `srv.Shutdown()` is therefore redundant — but not harmful. However, between `s.srv.Shutdown(ctx)` returning and `s.wg.Wait()` being reached, there is no logical gap because `http.Server.Shutdown` already guarantees all handlers have returned before it returns. The `wg.Wait()` at line 107 is always a no-op by the time it is reached.

The `activeCount` atomic is used only in tests (via `ActiveRequests()`); it is not used in shutdown logic, which is correct.
Assessment: PASS
Severity: LOW
Notes: The redundant `wg.Wait()` is harmless but misleading. Actual shutdown correctness is provided by `http.Server.Shutdown`, not the manual WaitGroup.

---

## Finding 2

Location: internal/server/server.go:85-109 (Shutdown method)
Claimed Behavior: Mark unready → preStop delay → drain connections → wg.Wait.
Observed Implementation: Sequence is correct. `SetReady(false)` first, then preStop select (context-cancellable), then `srv.Shutdown(ctx)`, then `wg.Wait()`.
Context cancellation during preStop returns `ctx.Err()` immediately without proceeding to listener shutdown — correct fail-fast behavior.
Assessment: PASS
Severity: LOW
Notes: None.

---

## Finding 3

Location: internal/server/server.go:55-62 (/work handler — client disconnect handling)
Claimed Behavior: If client disconnects early, handler exits without writing a response.
Observed Implementation: `r.Context().Done()` is selected; the handler returns without writing. The WaitGroup and activeCount are still decremented via defer. This is correct. `ActiveRequests()` will return to 0 even on client cancellation.
Assessment: PASS
Severity: LOW
Notes: None.

---

## Finding 4

Location: internal/worker/worker.go:35-63 (Start — goroutine loop)
Claimed Behavior: Worker goroutines stop when context is cancelled or channel is closed.
Observed Implementation: Double-select pattern: first checks ctx.Done with a default fallback, then blocks on either ctx.Done or jobChan receive. When `Stop()` closes `jobChan`, the `case job, ok := <-w.jobChan` branch triggers with `ok=false` and the goroutine returns. If context is cancelled (timeout path), `case <-w.ctx.Done()` in the second select fires. Correct in both paths.

However: the first select (lines 41-45) is a non-blocking ctx check before the blocking select. This creates a double-checked pattern that is idiomatic but slightly redundant — the second select already handles ctx.Done. No correctness issue.
Assessment: PASS
Severity: LOW
Notes: None.

---

## Finding 5

Location: internal/worker/worker.go:70-88 (Stop method)
Claimed Behavior: Close channel to stop new job intake; wait for drain; cancel context on timeout.
Observed Implementation: `close(w.jobChan)` stops new enqueues. A goroutine waits on `wg.Wait()` and signals `done`. Main select picks `done` (drain complete) or `time.After(timeout)` (cancel context then wait for done). After `w.cancel()`, goroutines in the second select will pick up `ctx.Done()` and return. `<-done` after cancel ensures all goroutines have exited before `Stop` returns. Correct cooperative shutdown.

Potential issue: after timeout triggers and `w.cancel()` is called, a goroutine currently inside `time.Sleep(job.Duration)` at line 55 will NOT be interrupted — `time.Sleep` is not context-aware. The goroutine will complete its sleep, append the job to completed, then check ctx.Done on next iteration. This means the "timeout" path may still complete the in-flight job before exiting, which is the observed behavior in `TestWorkerShutdownTimeout`.
Assessment: PASS
Severity: LOW
Notes: Behavior is documented in a ponytail comment. Test `TestWorkerShutdownTimeout` expects job-slow (in-flight, 100ms) completes despite 20ms timeout — this is consistent with the implementation since `time.Sleep` is not preemptible.

---

## Finding 6

Location: internal/worker/worker.go:66-68 (Enqueue)
Claimed Behavior: Enqueue a job.
Observed Implementation:
```go
func (w *Worker) Enqueue(job Job) {
    w.jobChan <- job
}
```
No guard against sending to a closed channel. After `Stop()` calls `close(w.jobChan)`, calling `Enqueue` would panic. In demo and tests, Enqueue is always called before Stop — safe in current usage. No runtime protection.
Assessment: WARNING
Severity: MEDIUM
Notes: A caller that calls Enqueue after Stop will panic. No test covers this. In a real system this would need a guard (e.g., recover or a state flag). For a lab demonstrating the pattern, this is acceptable but is a gap.

---

## Finding 7

Location: internal/db/db.go:61-69 (GetUser — fallback logic)
Claimed Behavior: Backward-compatible read for both legacy (Name only) and modern (FirstName+LastName) records.
Observed Implementation:
- If FirstName and LastName are empty but Name is set → split Name on first space to derive FirstName/LastName.
- Else if Name is empty → derive Name from FirstName+LastName.
- Both fields present → return as-is.

Edge case: a user with a single-name (no space) inserted via InsertLegacy would produce LastName="" — which is acceptable behavior. The `strings.SplitN` with n=2 handles the single-word case via `len(parts) > 1` guard.

Edge case: SaveExpand with empty firstName or lastName produces a combined Name via TrimSpace — handled correctly.
Assessment: PASS
Severity: LOW
Notes: No test for single-name legacy user or empty-field expand writes. Minor gap, not a correctness failure.

---

## Finding 8

Location: internal/server/server.go:73-79 (Start method)
Claimed Behavior: Start listening; return nil on graceful close.
Observed Implementation: `http.ErrServerClosed` is explicitly excluded from the error return. This is the standard Go pattern and is correct.
Assessment: PASS
Severity: LOW
Notes: None.

---

## Finding 9

Location: cmd/demo/main.go (overall orchestration)
Claimed Behavior: Simulate startup, in-flight request, SIGTERM, graceful shutdown of server and worker.
Observed Implementation:
- Worker started, job enqueued (DemoJob-1, 2s duration)
- Server started on :8080
- 1s init sleep, then SetReady(true)
- In-flight HTTP request goroutine (2s work)
- 500ms sleep, then simulate SIGTERM via channel
- Server Shutdown with 10s context (includes 1s preStop)
- Worker Stop with 5s timeout

Timeline: SIGTERM at ~1.5s. PreStop completes at ~2.5s. Worker job (DemoJob-1, 2s) finishes at ~2s. HTTP /work request (2s, started at ~1s) finishes at ~3s. All within timeouts.

Actual output confirmed: Client request completed status 200, "Demo finished cleanly."
Assessment: PASS
Severity: LOW
Notes: None.

---

## Finding 10

Location: tests/worker_test.go:10-28 (TestWorkerGracefulShutdown)
Claimed Behavior: Both job-1 and job-2 complete in order.
Observed Implementation: Test asserts `completed[0] == "job-1"` and `completed[1] == "job-2"`. With 1 worker goroutine, this is deterministic. With concurrency=1 and FIFO channel, ordering is guaranteed.
Assessment: PASS
Severity: LOW
Notes: None.

---

## Finding 11

Location: tests/worker_test.go:30-48 (TestWorkerShutdownTimeout)
Claimed Behavior: job-slow (in-flight) completes; job-dropped is abandoned due to timeout.
Observed Implementation: timeout=20ms, job-slow=100ms. Stop is called after 5ms sleep (job-slow is in-flight). After close(jobChan), the worker finishes job-slow (takes ~95ms more), then context is cancelled. job-dropped is never dequeued because the channel was closed before it could be picked up (worker was busy with job-slow). Result: exactly 1 completed job.

This test correctly exercises the timeout/drain behavior. However, it relies on timing: if job-slow finishes within 20ms, the context cancel would not fire and job-dropped might be dequeued. At 100ms job duration and 20ms timeout, there is sufficient margin. Minor timing sensitivity exists on heavily loaded CI machines.
Assessment: PASS
Severity: LOW
Notes: Timing-sensitive but margins are sufficient (5x ratio).
