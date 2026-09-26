# Docs vs Code Comparison

## README.md Claims

| Claim Location | Claimed Behavior | Code Implementation | Status | Notes |
|----------------|------------------|---------------------|--------|-------|
| Lines 3-4 | DB demonstrates Expand/Contract pattern with legacy `Name` and modern `FirstName`/`LastName`, transparent fallback on read | internal/db/db.go: InsertLegacy/SaveExpand dual writes; GetUser fallback derivation | PASS | Exact match |
| Lines 5-8 | Server exposes Liveness/Readiness probes; shutdown executes preStop delay then graceful shutdown ensuring in-flight requests complete | internal/server/server.go: live/ready probes; Shutdown: SetReady(false), preStop sleep, srv.Shutdown, wg.Wait | PASS | Verified by tests |
| Lines 9-11 | Worker pulls from queue; on shutdown stops pulling new jobs but continues processing current active job until completion | internal/worker/worker.go: Enqueue sends to chan; Stop closes chan, wg.Wait with timeout, active job completes time.Sleep | PASS | Verified by tests |
| Lines 12-13 | Demo orchestrates components, simulates startup, executes in-flight workloads, sends termination signal to demonstrate zero-downtime draining | cmd/demo/main.go: starts worker/server, marks ready, simulates SIGTERM, graceful Shutdown, worker Stop | PASS | Demo output matches claim |
| Lines 14-19 | Demo run via `go run ./cmd/demo` | Manual verification: runs successfully, logs show sequence | PASS | Output matches engineering/03-execution-result.md |
| Lines 21-27 | Tests via `go test -v ./...` and `go test -race ./...` | Manual verification: all tests pass, race detector clean | PASS | - |

## Engineering Notes vs Implementation

### 01-design.md
- Success Criteria lines 26-31 all verified via test suite and demo.
- Architecture lines 34-37: internal/db (Expand/Contract), internal/server (probes+liveness+inflight tracking+preStop+shutdown), internal/worker (cooperative termination), cmd/demo (lifecycle demo) — all match code.

### 02-implementation-notes.md
- Files Added list matches: cmd/demo/main.go, internal/server/server.go, internal/worker/worker.go, internal/db/db.go, tests/*_test.go, go.mod.
- Core Design Decisions: probe distinction, connection draining via Server.Shutdown, preStop simulation, worker drain via chan close+WaitGroup — all implemented.
- Implementation-Specific Choices: in-memory data structures, parameterized preStop — exact.
- Known Limitations: in-memory DB lacks SQL locking, no physical network routing — acknowledged.
- Trade-offs: Go channels vs Redis queue — used for simplicity, documented.
- What Is Demonstrated: probe distinction, preStop delay, in-flight HTTP completion, in-flight worker completion, Expand/Contract DB — all verified by tests/demo.
- What Is Not Demonstrated: actual K8s deployment, real DDL lock handling — correctly out of scope.

### 03-execution-result.md
- Build: `go build ./...` success — verified.
- Tests: `go test -v ./...` output matches — verified.
- Race Detector: `go test -race ./...` success — verified.
- Demo: `go run ./cmd/demo` output matches line-for-line (modulo timestamps) — verified.

## Conclusion
No DOC_CODE_MISMATCH, TEST_CLAIM_MISMATCH, or RESEARCH_IMPLEMENTATION_MISMATCH observed. Documentation accurately reflects implementation.