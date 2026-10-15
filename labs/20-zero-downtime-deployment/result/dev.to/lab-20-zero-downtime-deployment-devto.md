---
title: "Zero-Downtime Deployment: Graceful Shutdown, Traffic Draining & Expand-Contract DB"
published: false
description: "Implement zero-downtime deployments in Go: health probes, preStop hooks, connection draining, graceful worker termination, and expand-contract database migrations."
tags: go, microservices, architecture, testing
---

# Zero-Downtime Deployment: Graceful Shutdown, Traffic Draining & Expand-Contract DB

## 1. Problem: What Happens Without Draining

Shutting down an instance without coordination drops in-flight HTTP requests, aborts background jobs mid-execution, and breaks schema compatibility when databases mutate. Users receive 502/504 errors; jobs leave partial state; new application versions cannot read old data.

## 2. Mental Model: Start New Before Stopping Old

Zero-downtime deployment relies on coexistence: version two (v2) and version one (v1) run simultaneously, sharing infrastructure and database. Traffic routing, graceful termination, and schema evolution must be coordinated.

## 3. Architecture Components

The implementation has three components:
- `internal/db`: In-memory store implementing expand-contract schema pattern.
- `internal/server`: HTTP server with liveness/readiness probes, configurable preStop delay, and connection draining.
- `internal/worker`: Background queue worker that finishes current jobs before terminating.

## 4. Health Probes: Liveness vs Readiness

Liveness probe (`/healthz/live`) confirms the process is running. Readiness probe (`/healthz/ready`) decides whether the load balancer routes traffic.

Server starts with readiness false, returning HTTP 503. After initialization completes, `SetReady(true)` switches readiness to HTTP 200. Load balancers only send traffic to ready instances.

```go
mux.HandleFunc("/healthz/ready", func(w http.ResponseWriter, r *http.Request) {
    if s.ready.Load() {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("READY"))
    } else {
        w.WriteHeader(http.StatusServiceUnavailable)
        w.Write([]byte("NOT READY"))
    }
})
```

## 5. PreStop Hook: Bridging Asynchronous Routing

When an orchestrator like Kubernetes sends SIGTERM, endpoint slice updates and kube-proxy routing table changes happen asynchronously. If the pod terminates before routing updates propagate, clients still hit the terminating pod and receive 502/504.

The lab implements `preStop` with a configurable delay using `select` on `time.After` or `ctx.Done()`. This allows context cancellation to abort the delay early, avoiding deadlock when the caller times out.

```go
if s.preStop > 0 {
    select {
    case <-time.After(s.preStop):
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

A production `preStop` must exceed external routing latency (kube-proxy, cloud LB).

## 6. Graceful Shutdown: Draining In-Flight Requests

`Server.Shutdown` performs three steps:
1. Sets readiness false, marking the instance as unavailable.
2. Executes the preStop delay.
3. Calls `srv.Shutdown(ctx)` to stop accepting connections, then waits for active requests to finish.

```go
s.SetReady(false)
if s.preStop > 0 {
    select { case <-time.After(s.preStop): case <-ctx.Done(): return ctx.Err() }
}
err := s.srv.Shutdown(ctx)
s.wg.Wait()
```

`wg.Wait()` blocks until `activeCount` drops to zero—every in-flight request completes before the process exits.

## 7. Background Worker: Cooperative Drain

Worker starts goroutines that pull jobs from a channel. Upon `Stop(timeout)`:
1. Sets stopped flag and closes the job channel—no new jobs accepted.
2. Begins draining: waits for active goroutines to finish.
3. If timeout elapses before all workers finish, cancels context to abort remaining jobs.

```go
func (w *Worker) Stop(timeout time.Duration) {
    w.enqueueMu.Lock()
    w.stopped.Store(true)
    close(w.jobChan)
    w.enqueueMu.Unlock()

    done := make(chan struct{})
    go func() { w.wg.Wait(); close(done) }()

    select {
    case <-done:
    case <-time.After(timeout):
        w.cancel()
        <-done
    }
}
```

Workers finish active jobs; queued jobs are dropped when timeout expires. Production systems require job-level preemption or hard deadlines for arbitrary runtimes.

## 8. Expand-Contract Database Pattern

Database migration splits into phases:
- **Expand**: Add new columns (`first_name`, `last_name`) without dropping old ones (`name`). Application supports both read paths.
- **Migrate**: Backfill historical data, dual-write for new records.
- **Contract**: Remove old column (`name`) only after all instances run the new version.

```go
type UserRecord struct {
    ID        string
    Name      string
    FirstName string
    LastName  string
}
```

Read logic auto-converts legacy format: parses `Name` into `FirstName` and `LastName` if the modern fields are empty. Writes store both representations for compatibility.

## 9. Test Evidence

**Server Tests** (8 tests):
- `TestServerProbes`: Liveness returns 200, readiness starts 503, becomes 200 after `SetReady(true)`.
- `TestServerGracefulShutdown`: In-flight request completes during shutdown; no drops.
- `TestServerPreStopHook`: Shutdown delays exactly the configured `preStop` duration.
- `TestServerPreStopContextCancellation`: Context timeout aborts preStop early.
- `TestServerMultiRequestDrain`: Five concurrent requests all complete during shutdown.
- `TestServerWorkRequestCancellation`: Client context cancellation terminates long requests early.
- `TestServerInvalidDurationFallback`: Invalid `?d=INVALID` uses 50ms default.
- `TestServerReadyUnreadyTransition`: Ready→Unready→Ready transitions work correctly.

**Worker Tests** (5 tests):
- `TestWorkerGracefulShutdown`: Two queued jobs finish before worker stops.
- `TestWorkerEnqueueAfterStop`: Jobs rejected after `Stop` returns.
- `TestWorkerConcurrency`: Three workers process six jobs concurrently.
- `TestWorkerConcurrentEnqueueStop`: Race detector confirms no data races.
- `TestWorkerShutdownTimeout`: Short timeout finishes active job, drops queued job.

**DB Tests** (5 tests):
- `TestDBNotFound`: Returns expected error for missing records.
- `TestDBSingleNameLegacy`: Reads legacy `Name` and auto-converts to modern fields.
- `TestDBSaveExpandEmptyFields`: Handles partial empty names correctly.
- `TestDBLegacyOverwriteWithExpand`: Upgrades legacy record to modern format.
- `TestExpandContractDatabase`: Demonstrates full expand-contract cycle.

All 18 tests pass with `go test -v ./...` and zero race conditions under `go test -race ./...`.

## 10. Demo Execution

Running `go run ./cmd/demo` simulates a deployment:
1. Worker starts, processes one 2-second job.
2. Server starts, readiness false.
3. After 1 second, readiness true—load balancer begins routing.
4. HTTP GET request starts, simulating in-flight traffic.
5. SIGTERM received after 0.5 seconds.
6. Server sets readiness false, executes 1-second preStop, drains connections.
7. Client request completes during drain.
8. Worker stops, finishes active job.
9. Demo ends cleanly—zero dropped requests.

```
2026/09/25 14:58:40 Starting Zero-Downtime Deployment Demo
2026/09/25 14:58:40 Background worker started
2026/09/25 14:58:40 Server starting on 127.0.0.1:8080
2026/09/25 14:58:41 Application is ready to receive traffic
2026/09/25 14:58:41 SIGTERM received, initiating graceful shutdown procedures
2026/09/25 14:58:41 Server marked unready, detached from load balancer
2026/09/25 14:58:41 Executing preStop sleep for 1s to allow routing table updates...
2026/09/25 14:58:42 Worker 0 finished job DemoJob-1
2026/09/25 14:58:43 Client request completed with status: 200
2026/09/25 14:58:43 Demo finished cleanly. Zero downtime achieved.
```

## 11. Trade-offs and Production Considerations

**PreStop Duration**: Must exceed external routing latency. Underestimating causes HTTP 502/504.

**DDL Lock Contention**: PostgreSQL `ALTER TABLE ... ADD COLUMN` with constant default is metadata-only and concurrent-safe. Volatile defaults (`clock_timestamp()`, `gen_random_uuid()`) force `ACCESS EXCLUSIVE` lock and full table rewrite. Use `lock_timeout` on database connections during DDL.

**Worker Timeout Behavior**: Draining finishes active jobs while timeout remains. When timeout expires, context cancels, aborting remaining jobs. Production systems need job-level preemption or hard deadline escalation.

**HTTP Worker Patterns**: PHP-FPM does not handle SIGTERM for connection draining by default. Laravel Horizon provides `horizon:terminate` with Supervisor `stopwaitsecs` to finish jobs before shutdown.

## 12. When to Use This Pattern

Choose zero-downtime deployment when:
- SLA requires zero seconds of unavailability.
- Background jobs must complete or fail explicitly, not be killed abruptly.
- Database schema evolves across versions; old and new instances coexist during transition.

Avoid over-engineering for:
- Monoliths with single database: local ACID transactions suffice.
- Non-critical services where brief downtime is acceptable.

## 13. Conclusion

Zero-downtime deployment is not magic—it is orchestration of health probes, preStop hooks, connection draining, and expand-contract migrations. The pattern guarantees no request drops and no job loss, trading brief coexistence complexity for uninterrupted service.

### Laboratory Resources
- Source Code Repository: [software-engineering-lab / labs / 20-zero-downtime-deployment](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/20-zero-downtime-deployment)
- Commit SHA: `f303829f429b7e14668462442f032cf3ba0a4cce`
- Tests: 18 passing (5 DB, 8 server, 5 worker), zero race conditions.
