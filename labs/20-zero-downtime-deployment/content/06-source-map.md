# Source Map

## Mental Model & Core Concept (Coexistence, Probes, Graceful Shutdown, DB Expand/Contract)

**Research:**
- `research/runs/2026-09-25-zero-downtime-deployment/05-report.md` (Executive Summary, Finding 1, Finding 2, Finding 3)

**Implementation:**
- `internal/server/server.go` (Probes setup, HTTP connection draining)
- `internal/db/db.go` (Expand and Contract schema support)
- `internal/worker/worker.go` (Cooperative termination)

**Tests:**
- `tests/server_test.go` (`TestServerProbes`, `TestServerGracefulShutdown`)
- `tests/worker_test.go` (`TestWorkerGracefulShutdown`)
- `tests/db_test.go` (`TestExpandContractDatabase`)

---

## Asynchronous Routing Detachment & PreStop Hook

**Research:**
- `research-audit/06-gaps.md` (Gap 1: Asynchronous endpoint propagation)

**Implementation:**
- `internal/server/server.go` (`Shutdown` method featuring `preStop` delay)
- `cmd/demo/main.go` (`1*time.Second` preStop hook configuration)

**Tests:**
- `tests/server_test.go` (`TestServerPreStopHook`)

---

## Database Expand and Contract & Lock Contention

**Research:**
- `research/runs/2026-09-25-zero-downtime-deployment/05-report.md` (Finding 3: Database and State Backward Compatibility)
- `research-audit/06-gaps.md` (Gap 2: PostgreSQL DDL locking dynamics and `ACCESS EXCLUSIVE`)

**Implementation:**
- `internal/db/db.go` (`InsertLegacy`, `SaveExpand`, `GetUser` logic)
- `engineering/01-design.md` (Database Backward Compatibility concept)

**Tests:**
- `tests/db_test.go` (`TestExpandContractDatabase`)

---

## Worker Graceful Termination & Buffered Queue Issues

**Research:**
- `research/runs/2026-09-25-zero-downtime-deployment/05-report.md` (Finding 2: Background job draining)
- `research-audit/06-gaps.md` (Gap 3: Laravel cache restart vs standard OS POSIX signals)

**Implementation:**
- `internal/worker/worker.go` (`Stop` method, `jobChan` processing loop)
- `engineering-audit/05-gaps.md` (Gap 1: Worker buffered jobs abandonment via immediate `w.cancel()`)

**Tests:**
- `tests/worker_test.go` (`TestWorkerGracefulShutdown`)
- `engineering-audit/05-gaps.md` (Gap 2: Weak assertions masking dropped buffered tasks)
