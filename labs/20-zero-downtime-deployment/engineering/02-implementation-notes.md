# Implementation Notes

## Files Added
- `cmd/demo/main.go`: Demo harness showing full lifecycle with readiness, graceful shutdown, preStop hook, and background queue draining.
- `internal/server/server.go`: HTTP server with probe endpoints (`/healthz/live`, `/healthz/ready`), request tracking, and cooperative termination with preStop delay.
- `internal/worker/worker.go`: Background queue worker that processes active jobs cooperatively before exiting on shutdown.
- `internal/db/db.go`: In-memory storage showing the Expand and Contract pattern with fallback read compatibility.
- `tests/server_test.go`: Unit tests for probes, preStop delay, and connection draining.
- `tests/worker_test.go`: Unit tests for graceful worker termination.
- `tests/db_test.go`: Unit tests for legacy and modern database operations.
- `go.mod`: Module definition.

## Core Design Decisions
- **Probes**: Separate `/healthz/live` (status of process) from `/healthz/ready` (traffic availability) based on research Findings.
- **Connection Draining**: Using `net/http` standard library `Server.Shutdown(ctx)` which stops accepting new connections and waits for active ones.
- **PreStop Hook**: Explicitly wait before invoking listener shutdown to mitigate the asynchronous routing table update race condition identified in Gap 1 of the audit.
- **Worker Drain**: Wait for ongoing job to complete via Go channels and sync primitives instead of immediate hard kill.

## Implementation-Specific Choices
- In-memory data structures are used for DB and worker queues to keep the lab completely dependency-free and runnable with standard Go toolchain.
- PreStop hook is parameterized directly on the `Server` struct.

## Known Limitations
- The in-memory database does not replicate real SQL locking behavior (e.g. `ACCESS EXCLUSIVE` table lock queuing in PostgreSQL).
- Does not test physical network routing infrastructure; routing behavior is simulated through the preStop hook.

## Trade-offs
- Used Go channels for queue worker simulation instead of a real Redis/database queue, sacrificing multi-process distributed semantics for simplicity and portability.

## What Is Demonstrated
- Readiness vs Liveness probe distinction.
- PreStop delay handling before graceful server shutdown.
- Non-interrupted completion of active in-flight HTTP requests during a shutdown signal.
- Non-interrupted completion of active worker tasks during a shutdown signal.
- Expand and Contract database compatibility reading old/new record formats.

## What Is Not Demonstrated
- Actual container deployment on Kubernetes clusters.
- Real DDL lock timeout handling against external PostgreSQL/MySQL databases.
