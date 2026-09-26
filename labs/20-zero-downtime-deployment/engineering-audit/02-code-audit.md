# Code Audit

Target Lab: labs/20-zero-downtime-deployment

## Finding 1

Location: internal/worker/worker.go:40-45, 62-68
Claimed Behavior: Worker stops accepting new jobs upon shutdown and completes active in-flight jobs, draining existing work.
Observed Implementation:
`w.Stop()` calls `close(w.jobChan)` followed immediately by `w.cancel()`.
Inside `w.Start()`, goroutines listen on a `select` statement between `<-w.ctx.Done()` and `job, ok := <-w.jobChan`.
Because `w.ctx.Done()` is closed, the `select` may prioritize `<-w.ctx.Done()`, exiting immediately before processing remaining buffered jobs inside `w.jobChan`.
Assessment: WARNING
Severity: MEDIUM
Notes: Buffered jobs in the queue might be abandoned if multiple jobs are already enqueued and awaiting pickup when `Stop()` is triggered. For true draining of accepted jobs, `w.cancel()` should only be called after the queue is drained or timeout is reached, or goroutines should drain `w.jobChan` until it is empty.

## Finding 2

Location: internal/server/server.go:84-104
Claimed Behavior: Graceful shutdown with preStop hook delay and connection draining.
Observed Implementation:
Correctly sets readiness probe to false, sleeps for `preStop` delay (simulating iptables/kube-proxy/load-balancer endpoint update window), calls `s.srv.Shutdown(ctx)`, and waits for `s.wg.Wait()` on active handlers.
Assessment: PASS
Severity: LOW
Notes: Implementation correctly matches research-backed Kubernetes graceful termination pattern.

## Finding 3

Location: internal/db/db.go:40-70
Claimed Behavior: Expand and contract pattern backward compatibility.
Observed Implementation:
Supports dual writes in `SaveExpand` and fallback reading in `GetUser` for both legacy records (splitting `Name` into `FirstName`/`LastName`) and expand-phase records. Thread-safe using `sync.RWMutex`.
Assessment: PASS
Severity: LOW
Notes: Correctly mocks expand/contract pattern for schema evolution.
