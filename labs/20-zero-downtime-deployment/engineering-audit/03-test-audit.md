# Test Audit

Target Lab: labs/20-zero-downtime-deployment

---

## Test Inventory

| Test | File | Coverage Area |
|---|---|---|
| TestExpandContractDatabase | db_test.go | DB expand/contract happy path |
| TestServerProbes | server_test.go | Liveness always 200; Readiness 503→200 transition |
| TestServerGracefulShutdown | server_test.go | In-flight request completes during shutdown |
| TestServerPreStopHook | server_test.go | PreStop delay enforced before listener close |
| TestServerPreStopContextCancellation | server_test.go | PreStop aborts on context timeout |
| TestServerWorkRequestCancellation | server_test.go | Client disconnect clears activeCount |
| TestWorkerGracefulShutdown | worker_test.go | Both jobs drain before Stop returns |
| TestWorkerShutdownTimeout | worker_test.go | In-flight job completes; queued job dropped on timeout |

Total: 8 tests across 3 files.

---

## Happy Path Coverage

- Liveness probe always returns 200: COVERED (TestServerProbes)
- Readiness probe: NOT READY → READY transition: COVERED (TestServerProbes)
- In-flight request completes during graceful shutdown: COVERED (TestServerGracefulShutdown)
- PreStop delay enforced: COVERED (TestServerPreStopHook)
- Worker drains all jobs when time allows: COVERED (TestWorkerGracefulShutdown)
- DB expand/contract read/write: COVERED (TestExpandContractDatabase)

## Failure Path Coverage

- PreStop context cancellation (hard deadline): COVERED (TestServerPreStopContextCancellation)
- Worker timeout — in-flight completes, queued dropped: COVERED (TestWorkerShutdownTimeout)
- Client disconnect clears active request counter: COVERED (TestServerWorkRequestCancellation)

## Edge Cases

- Single-name legacy user (no space in Name): NOT COVERED
- Empty firstName or lastName in SaveExpand: NOT COVERED
- DB record not found (ErrNotFound): NOT COVERED
- Enqueue after Stop (panic guard): NOT COVERED
- Zero-duration preStop: implicitly covered (TestServerGracefulShutdown uses preStop=0)
- Ready→Unready→Ready transition: NOT COVERED (only tests NOT READY → READY)
- Concurrent readiness toggle under load: NOT COVERED
- Worker with concurrency > 1: NOT COVERED (all tests use concurrency=1)

## Concurrency / Race

- go test -race -count=1 ./...: PASS (verified by execution)
- Worker completed slice protected by completedMu: correct
- Server activeCount via atomic.Int32: correct
- Server ready via atomic.Bool: correct

## Negative Cases

- /work with invalid duration falls back to 50ms default: NOT COVERED by test (implemented in code but untested)
- Shutdown called with already-expired context: implicitly tested via TestServerPreStopContextCancellation

---

## Test Quality Assessment

Tests are well-structured and test real behavior against a live HTTP server (not mocked). No fake sleeps substituting for assertions — tests verify actual outcomes (response body, status codes, completion counts). Timing margins are adequate. The suite proves the core claimed behaviors. Missing edge cases are minor and do not affect the primary correctness claims.
