---
title: "Zero-Downtime Deployment: Probes, PreStop Hook, and DB Parallel Change in Go"
published: false
description: "Coordinating traffic routing, graceful termination, backward-compatible schema changes, and cooperative worker shutdown for zero-downtime releases."
tags: "go, devops, architecture, kubernetes"
---

Application updates frequently drop active HTTP requests, disrupt executing background jobs, and trigger database schema mismatches. A naïve restart sequence terminates running processes with `SIGKILL` or abrupt `SIGTERM` closures, generating client-facing HTTP 502 Bad Gateway errors and leaving records in half-migrated states.

In **Software Engineering Lab 20**, we examine and implement the core engineering patterns required to achieve zero-downtime deployment (ZDD): decoupled health probes, asynchronous routing detachment via `preStop` delays, cooperative HTTP connection draining, background queue worker cooperative termination, and database Expand and Contract (Parallel Change) migrations.

```
Incoming Traffic ---> [ Load Balancer / Ingress ]
                            |               |
                  Ready Pods v              v (Traffic Detached via Probe & PreStop)
                      [ App Instance v2 ]      [ App Instance v1 (Draining) ]
```

---

## The Core Mental Model: Start New Before Stopping Old

Zero-downtime deployment requires version coexistence. Version $N$ (modern) and Version $N-1$ (legacy) must run concurrently on the same infrastructure, share the same database cluster, and handle production load simultaneously during rollout. 

This model requires three strict operational rules:
1. **Traffic Routing Decoupling**: Traffic must be dynamically rerouted to instances only after they pass readiness checks, and detached prior to process shutdown.
2. **Cooperative Connection Draining**: Running processes must finish ongoing work before exiting.
3. **Database Schema Backward Compatibility**: Database schema mutations must support reads and writes from both version $N$ and version $N-1$ at every phase of deployment.

---

## 1. Health Probes: Separating Liveness from Readiness

In modern orchestrators such as Kubernetes, health checks serve two distinct purposes:

- **Liveness Probe (`/healthz/live`)**: Answers "Is the process healthy or stuck in a deadlock/panic?". If liveness fails, the orchestrator restarts the container.
- **Readiness Probe (`/healthz/ready`)**: Answers "Is the instance ready to receive client traffic?". If readiness fails, the orchestrator removes the pod from load balancer endpoints without terminating the process.

Coupling these two states causes cascading deployment outages. An application warming up internal caches or undergoing brief GC pauses must report unready without triggering a restart loop.

In `internal/server/server.go`, readiness is decoupled into an atomic boolean flag:

```go
func (s *Server) SetReady(ready bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ready = ready
}

func (s *Server) liveHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

func (s *Server) readyHandler(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	ready := s.ready
	s.mu.RUnlock()

	if !ready {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("UNREADY"))
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("READY"))
}
```

When initialization begins, readiness returns HTTP 503 Service Unavailable. Once dependencies are verified, `SetReady(true)` flips the probe to HTTP 200 OK.

---

## 2. The Asynchronous Detachment Gap and PreStop Hook

A major misconception in Kubernetes deployments is assuming that receiving `SIGTERM` guarantees the pod has already been detached from the ingress or load balancer. 

In reality, pod termination and service routing updates happen concurrently and asynchronously:

1. **Control Plane Path**: API Server updates Pod state to `Terminating` $\rightarrow$ Kubelet sends `SIGTERM` to container processes.
2. **Data Plane Path**: Endpoints controller updates `EndpointSlice` $\rightarrow$ Kube-proxy / Cloud Controller / Ingress Controller receives update $\rightarrow$ iptables / IPVS / proxy backends reprogram routing rules.

Because the data plane path traverses multiple network hops and asynchronous controllers, the legacy pod often receives `SIGTERM` seconds before the load balancer stops routing traffic to it. If the server stops listening immediately, in-flight traffic routed during this propagation gap encounters refused connections or HTTP 502 errors.

To solve this, our server implementation introduces a configurable `preStop` delay executed before stopping listeners:

```go
func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Server received shutdown request")

	// 1. Immediately signal unready
	s.SetReady(false)
	log.Println("Server marked unready, detached from load balancer")

	// 2. PreStop delay: allow external routing tables to propagate
	if s.preStop > 0 {
		log.Printf("Executing preStop sleep for %v to allow routing table updates...", s.preStop)
		select {
		case <-time.After(s.preStop):
		case <-ctx.Done():
			log.Println("preStop sleep interrupted by context cancellation")
			return ctx.Err()
		}
	}

	// 3. Graceful listener termination
	log.Println("Initiating graceful shutdown of HTTP listeners...")
	err := s.srv.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("shutdown error: %w", err)
	}

	// 4. Wait for in-flight request handlers
	s.wg.Wait()
	log.Println("All in-flight requests completed. Server stopped gracefully.")
	return nil
}
```

Notice that `preStop` uses `select` on `time.After` and `ctx.Done()` rather than a bare `time.Sleep`. If the orchestrator deadline expires or a shutdown timeout triggers, the delay unblocks cooperatively without hanging the process.

---

## 3. Database Parallel Change: Expand and Contract

Relational database schemas represent the hardest constraint in zero-downtime systems. A breaking schema migration—such as renaming a column `name` to `first_name` and `last_name`—cannot be done atomically alongside application rollout. If you rename the column in place, running instances of version $N-1$ will crash with SQL syntax or column-missing errors.

The **Expand and Contract** pattern (Parallel Change) separates the rollout into distinct phases:

1. **Phase 1 (Expand)**: Add the new columns (`first_name`, `last_name`) as nullable fields. The legacy column `name` remains untouched.
2. **Phase 2 (Dual-Write & Fallback Read)**: Deploy application version $N$. When writing, version $N$ writes both formats or populates new fields. When reading, it reads modern fields if present, falling back to legacy fields if absent.
3. **Phase 3 (Backfill)**: Run an asynchronous background migration to backfill legacy rows into modern columns.
4. **Phase 4 (Contract)**: After version $N-1$ is fully decommissioned, remove the legacy column `name` in a subsequent release.

In `internal/db/db.go`, our storage engine models transparent fallback logic:

```go
type UserRecord struct {
	ID        string
	Name      string // Legacy format
	FirstName string // Modern format
	LastName  string // Modern format
}

func (s *UserStore) GetUser(id string) (UserRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	rec, ok := s.records[id]
	if !ok {
		return UserRecord{}, ErrNotFound
	}

	// Fallback: If modern fields are missing, synthesize from legacy Name
	if rec.FirstName == "" && rec.LastName == "" && rec.Name != "" {
		parts := strings.SplitN(rec.Name, " ", 2)
		rec.FirstName = parts[0]
		if len(parts) > 1 {
			rec.LastName = parts[1]
		}
	} else if rec.Name == "" {
		// If legacy Name is missing, synthesize from modern fields
		rec.Name = strings.TrimSpace(rec.FirstName + " " + rec.LastName)
	}

	return rec, nil
}
```

Both old records (`Name = "Alice Smith"`) and new records (`FirstName = "Bob"`, `LastName = "Jones"`) are read seamlessly across versions.

### Production DDL Locking Hazard: PostgreSQL Table Rewrites

Adding columns in production databases requires extreme care. In PostgreSQL:
- Adding a column with a constant default (e.g. `DEFAULT 'active'`) in PostgreSQL 11+ is a metadata-only operation and executes instantly without table rewrites.
- Adding a column with a volatile default (e.g. `DEFAULT clock_timestamp()` or `gen_random_uuid()`) forces PostgreSQL to rewrite the entire physical table on disk, acquiring an `ACCESS EXCLUSIVE` lock.
- Even metadata-only DDL acquires an `ACCESS EXCLUSIVE` lock briefly. If a long-running transaction holds an open query on the table, the `ALTER TABLE` statement queues behind it, blocking all subsequent reads and writes.

Production migrations must always enforce a tight lock timeout before executing DDL:
```sql
SET lock_timeout = '2s';
ALTER TABLE users ADD COLUMN first_name VARCHAR(100);
```

---

## 4. Background Queue Worker Cooperative Termination

Web servers handle request-response lifecycles, but background queue workers execute long-running asynchronous jobs. Terminating a worker abruptly corrupts state or leads to duplicate job processing (at-least-once delivery duplicates).

In `internal/worker/worker.go`, workers listen for shutdown signals cooperatively:

```go
func (w *Worker) Start(concurrency int) {
	for i := 0; i < concurrency; i++ {
		w.wg.Add(1)
		go func(workerID int) {
			defer w.wg.Done()
			for {
				select {
				case <-w.ctx.Done():
					return
				case job, ok := <-w.jobChan:
					if !ok {
						return
					}
					// Double check context after dequeuing
					select {
					case <-w.ctx.Done():
						log.Printf("Worker %d dropping job %s due to shutdown context cancellation", workerID, job.ID)
						return
					default:
					}

					log.Printf("Worker %d starting job %s", workerID, job.ID)
					select {
					case <-time.After(job.Duration):
						w.completedMu.Lock()
						w.completed = append(w.completed, job.ID)
						w.completedMu.Unlock()
						log.Printf("Worker %d finished job %s", workerID, job.ID)
					case <-w.ctx.Done():
						log.Printf("Worker %d aborted job %s due to shutdown timeout", workerID, job.ID)
						return
					}
				}
			}
		}(i)
	}
}
```

When `w.Stop(timeout)` is called:
1. It closes the intake channel, rejecting any new job submissions (`Enqueue` returns error).
2. Active goroutines complete their currently executing job.
3. A drain timer runs. If the worker completes jobs before the timeout, it shuts down cleanly.
4. If the drain timeout expires, the context is canceled, aborting in-flight tasks and dropping uncompleted buffer entries.

---

## Automated Test Verification

The test suite in `labs/20-zero-downtime-deployment` verifies all architectural properties across HTTP, database, and background worker systems:

```bash
$ go test -v ./...
=== RUN   TestDBSaveExpandEmptyFields
--- PASS: TestDBSaveExpandEmptyFields (0.00s)
=== RUN   TestDBLegacyOverwriteWithExpand
--- PASS: TestDBLegacyOverwriteWithExpand (0.00s)
=== RUN   TestExpandContractDatabase
--- PASS: TestExpandContractDatabase (0.00s)
=== RUN   TestServerProbes
--- PASS: TestServerProbes (0.00s)
=== RUN   TestServerGracefulShutdown
--- PASS: TestServerGracefulShutdown (0.16s)
=== RUN   TestServerPreStopHook
--- PASS: TestServerPreStopHook (0.10s)
=== RUN   TestServerPreStopContextCancellation
--- PASS: TestServerPreStopContextCancellation (0.05s)
=== RUN   TestServerInvalidDurationFallback
--- PASS: TestServerInvalidDurationFallback (0.05s)
=== RUN   TestServerReadyUnreadyTransition
--- PASS: TestServerReadyUnreadyTransition (0.00s)
=== RUN   TestServerMultiRequestDrain
--- PASS: TestServerMultiRequestDrain (0.16s)
=== RUN   TestServerWorkRequestCancellation
--- PASS: TestServerWorkRequestCancellation (0.07s)
=== RUN   TestWorkerConcurrency
--- PASS: TestWorkerConcurrency (0.04s)
=== RUN   TestWorkerGracefulShutdown
--- PASS: TestWorkerGracefulShutdown (0.04s)
=== RUN   TestWorkerEnqueueAfterStop
--- PASS: TestWorkerEnqueueAfterStop (0.00s)
=== RUN   TestWorkerConcurrentEnqueueStop
--- PASS: TestWorkerConcurrentEnqueueStop (0.00s)
=== RUN   TestWorkerShutdownTimeout
--- PASS: TestWorkerShutdownTimeout (0.03s)
PASS
ok  	zero-downtime-deployment/tests	1.197s
```

Running under `go test -race ./...` completes with zero data races.

---

## Production Takeaways & Best Practices

1. **Never Sleep Blindly**: Implement `preStop` delays using cancellable context timers (`select` with `ctx.Done()`) to avoid unkillable zombie pods.
2. **Align Timeouts**: Ensure your shutdown sequence timeouts are ordered: `preStop delay < HTTP connection draining timeout < orchestrator terminationGracePeriodSeconds`.
3. **Separate DDL Deployments**: Never bundle breaking database migrations and application code updates into the same deployment step. Always apply Expand, deploy application code, run backfills, and apply Contract in separate releases.
4. **Use Lock Timeouts**: Guard all production DDL with explicit `lock_timeout` statements to prevent table lock pileups during schema expansion.
5. **Differentiate HTTP and Queue Workers**: Web servers drain requests quickly; queue workers may process jobs lasting minutes. Configure queue supervisors (e.g. Supervisor, Horizon) with sufficient `stopwaitsecs` timeouts.

---

## Lab Source & Boundaries

- **Source Code**: [GitHub - Software Engineering Lab 20](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/20-zero-downtime-deployment)
- **Verified Commit**: `04890623a7f333e780a38551686fa7dba2284471`
- **Lab Boundaries**: This lab provides in-memory Go simulations for HTTP draining, queue workers, and database fallback logic. It does not spin up external Kubernetes clusters, cloud load balancers, or PostgreSQL instances. Production deployments must account for external proxy propagation latencies and storage engine table lock queues.
