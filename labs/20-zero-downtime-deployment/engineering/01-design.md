# Engineering Design

Target Lab: labs/20-zero-downtime-deployment
Research Status: APPROVED_WITH_WARNINGS

## Concept To Prove
Zero-Downtime Deployment (ZDD) requires coordinating:
1. Health Probes (Readiness vs. Liveness) for traffic orchestration.
2. Graceful Shutdown & Connection Draining (handling SIGTERM, allowing in-flight requests to complete, supporting preStop delays for asynchronous routing updates).
3. Database Backward Compatibility (Parallel Change / Expand and Contract pattern).
4. Background Worker Graceful Termination (allowing active jobs to finish before stopping).

## Expected Behavior
- **Traffic Routing**: Traffic is only routed to the instance once the `/ready` probe returns HTTP 200 (indicating readiness). The `/live` probe indicates process vitality.
- **Graceful HTTP Shutdown**: When receiving a shutdown signal (`SIGTERM` or cancellation), the server stops taking new traffic, optionally waits for a preStop delay to simulate network routing detachment, and allows active in-flight requests to complete before terminating.
- **Worker Shutdown**: The background queue worker stops fetching new tasks upon termination signal, but completes any currently running job.
- **Expand/Contract Database Compatibility**: The application can handle a database schema transition where both old fields and new fields coexist, demonstrating the Expand phase and backward-compatible reads/writes.

## Failure Scenario
- **Sudden SIGKILL / Hard Stop**: Dropped active in-flight requests returning network errors or partial writes.
- **Early Shutdown without PreStop**: Traffic routed to a terminated container due to routing table propagation latency.
- **Incompatible DB Schema Migrations**: v1 application fails if a column is abruptly removed, or v2 application fails if a required column is not populated.
- **Abrupt Worker Kill**: Interrupted long-running jobs causing partial state or data corruption.

## Success Criteria
- HTTP server readiness check rejects traffic when unready, and accepts traffic when ready.
- Active in-flight requests successfully complete during graceful shutdown.
- PreStop hook waits for configured delay before closing listeners.
- Worker finishes ongoing job upon shutdown signal and does not drop the task.
- Expand and Contract database reads/writes work seamlessly across v1 and v2 formats.
- All unit and integration tests pass without race conditions (`go test -race ./...`).

## Architecture
- `internal/db`: In-memory simulated database demonstrating the Expand and Contract pattern.
- `internal/server`: HTTP server featuring Liveness, Readiness, in-flight transaction tracking, and graceful shutdown with configurable preStop delay.
- `internal/worker`: Background queue worker that processes tasks and honors cooperative termination signals.
- `cmd/demo`: CLI application executing a mock deployment sequence illustrating readiness, concurrent requests, graceful shutdown, and worker draining.

## Components
1. **Database Repository (`internal/db`)**:
   - `User` model supporting legacy `name` and new `first_name` + `last_name`.
   - Dual-write and fallback-read compatibility logic.
2. **HTTP Server (`internal/server`)**:
   - `/healthz/live` (Liveness)
   - `/healthz/ready` (Readiness)
   - `/work` (Simulates an in-flight business transaction)
   - `Shutdown(ctx, preStopDelay)`: Implements graceful connection draining with a preStop sleep.
3. **Queue Worker (`internal/worker`)**:
   - `Worker`: Consumes tasks from an in-memory queue.
   - `Stop()`: Signals the worker to finish the current job and exit gracefully.

## Test Strategy
- Test Liveness and Readiness state transitions.
- Test in-flight request completion during graceful shutdown.
- Test preStop delay behavior.
- Test worker graceful shutdown ensuring in-flight jobs complete.
- Test database compatibility across legacy and expanded schema states.

## Execution Plan
1. Initialize Go module.
2. Implement `internal/db` with Expand/Contract logic.
3. Implement `internal/server` with probes and graceful termination.
4. Implement `internal/worker` with cooperative job draining.
5. Write comprehensive automated tests in `tests/`.
6. Implement `cmd/demo/main.go` demonstrating the complete lifecycle.
7. Verify all tests with `go test -race ./...` and run the demo.

## Implementation Decisions
- **In-Memory Concurrency Primitives**: Use Go channels, sync primitives, and standard library constructs rather than third-party databases/message brokers to keep the lab portable and self-contained.
- **PreStop Simulation**: Implemented as a configurable duration before triggering the standard `http.Server.Shutdown()` method to simulate K8s preStop hook behavior.
