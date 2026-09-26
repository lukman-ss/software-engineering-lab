# Source Map

## Mental Model & Core Concept (Coexistence, Probes, Graceful Shutdown, DB Expand/Contract)

**Research:**
- `research/runs/2026-09-26-zero-downtime-deployment/05-report.md` (Executive Summary, Finding 1, Finding 2, Finding 3)
- `research/runs/2026-09-26-zero-downtime-deployment/01-plan.md` (Research Questions 1-5 on ZDD patterns)

**Implementation:**
- `internal/server/server.go` (`NewServer`, `Shutdown`, `SetReady`) - Probes setup, HTTP connection draining
- `internal/db/db.go` (`InsertLegacy`, `SaveExpand`, `GetUser`) - Expand and Contract schema support
- `internal/worker/worker.go` (`NewWorker`, `Start`, `Stop`, `Enqueue`) - Cooperative termination

**Tests:**
- `tests/server_test.go` (`TestServerProbes`, `TestServerGracefulShutdown`, `TestServerPreStopHook`, `TestServerPreStopContextCancellation`, `TestServerMultiRequestDrain`)
- `tests/worker_test.go` (`TestWorkerGracefulShutdown`, `TestWorkerShutdownTimeout`, `TestWorkerEnqueueAfterStop`)
- `tests/db_test.go` (`TestExpandContractDatabase`, `TestDBLegacyOverwriteWithExpand`)

---

## Health Probes: Liveness, Readiness, Readiness Transition

**Research:**
- `research/runs/2026-09-26-zero-downtime-deployment/03-evidence.md` (Finding 2: Liveness vs Readiness probes)
- `research/runs/2026-09-26-zero-downtime-deployment/05-report.md` (Finding 5: Health Check Depth Requirements)

**Implementation:**
- `internal/server/server.go:28-41` (`/healthz/live`, `/healthz/ready` handlers)

**Tests:**
- `tests/server_test.go` (`TestServerProbes`, `TestServerReadyUnreadyTransition`)

---

## Graceful HTTP Shutdown & PreStop Delay

**Research:**
- `research/runs/2026-09-26-zero-downtime-deployment/03-evidence.md` (Finding 6: Pod Termination Graceful Shutdown)
- `research/runs/2026-09-26-zero-downtime-deployment/06-open-questions.md` (Q6: Kubernetes Endpoint Removal vs Connection Draining Gap)

**Implementation:**
- `internal/server/server.go:85-110` (`Shutdown` method with preStop select-based timer)
- `internal/server/server.go:43-63` (`/work` handler with request tracking via WaitGroup)
- `cmd/demo/main.go:66-72` (10-second server shutdown timeout)

**Tests:**
- `tests/server_test.go` (`TestServerGracefulShutdown`, `TestServerPreStopHook`, `TestServerPreStopContextCancellation`, `TestServerWorkRequestCancellation`, `TestServerMultiRequestDrain`)

---

## Database Expand and Contract Pattern

**Research:**
- `research/runs/2026-09-26-zero-downtime-deployment/05-report.md` (Finding 3, Finding 12: Database Backward Compatibility)
- `research/runs/2026-09-26-zero-downtime-deployment/03-evidence.md` (Finding 4: Renaming/Dropping Columns, Finding 10: Constant Default, Finding 12: Expand-Deploy-Migrate-Contract)
- `research-audit/06-gaps.md` (Gap 2: PostgreSQL DDL locking dynamics, Gap 1: Volatile Default Warning)

**Implementation:**
- `internal/db/db.go:32-51` (`InsertLegacy`, `SaveExpand` methods)
- `internal/db/db.go:53-72` (`GetUser` fallback-read logic)

**Tests:**
- `tests/db_test.go` (`TestExpandContractDatabase`, `TestDBSingleNameLegacy`, `TestDBLegacyOverwriteWithExpand`, `TestDBSaveExpandEmptyFields`)

---

## Background Worker Graceful Termination

**Research:**
- `research/runs/2026-09-26-zero-downtime-deployment/05-report.md` (Finding 2: Graceful Shutdown)
- `research/runs/2026-09-26-zero-downtime-deployment/06-open-questions.md` (Q2: Laravel Application-Level SIGTERM Handling, Q3: Redis Queue Job Persistence)

**Implementation:**
- `internal/worker/worker.go:38-74` (`Start` method with cooperative context checks)
- `internal/worker/worker.go:86-107` (`Stop` method with configurable drain timeout)
- `internal/worker/worker.go:76-84` (`Enqueue` with stopped-state check)
- `internal/worker/worker.go:17-27` (ponytail comments on buffered job abandonment ceiling)

**Tests:**
- `tests/worker_test.go` (`TestWorkerGracefulShutdown`, `TestWorkerShutdownTimeout`, `TestWorkerConcurrency`, `TestWorkerEnqueueAfterStop`, `TestWorkerConcurrentEnqueueStop`)

---

## Demo Orchestrated Lifecycle

**Implementation:**
- `cmd/demo/main.go` (Full sequence: worker start, server with 1s preStop, in-flight request, SIGTERM simulation, graceful shutdown, worker drain)

**Execution Result:**
- `engineering/03-execution-result.md` (Demo output verification)