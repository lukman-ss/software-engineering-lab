# Test Audit

Target Lab: `labs/33-read-replicas-and-replication-lag`

## Test Coverage Summary

The automated test suite in `tests/replication_test.go` contains 7 test functions testing all required operational modes and edge cases:

1. `TestNaiveReplicationLag_StaleRead`:
   - **Coverage**: Happy/Failure path for naive routing.
   - **Verification**: Asserts that an immediate read on an async lagging replica returns `ErrNotFound`, and after lag duration expires, returns expected data (`"1000"`).

2. `TestStickySessionRouting`:
   - **Coverage**: Time-based sticky session TTL and isolation.
   - **Verification**: Confirms read goes to `primary` within sticky window, other sessions read from replica, and after TTL expires, session reads from replica.

3. `TestReadWithToken_LSN`:
   - **Coverage**: Causal LSN token routing with replica catch-up wait.
   - **Verification**: Measures elapsed time to verify blocking wait (~200ms) until replica applies LSN before returning data.

4. `TestReplicaLagThreshold_Fallback`:
   - **Coverage**: Lag SLA threshold calculation and primary fallback.
   - **Verification**: Sets replica lag to 10s and performs writes exceeding `MaxLSNDiff` (2), verifying router falls back to `"primary (fallback-lag)"`.

5. `TestSynchronousReplication_Freshness`:
   - **Coverage**: Synchronous replication mode (`remote_apply`).
   - **Verification**: Confirms write blocks until replicas apply WAL, ensuring immediate naive read return from replica without stale errors.

6. `TestConcurrentAccess_RaceFree`:
   - **Coverage**: Heavy concurrent read/write operations across multiple goroutines with race detection.
   - **Verification**: Runs 5 concurrent writers and 10 concurrent readers executing 20 ops each across sticky, token, naive, and lag-aware calls. Passes clean under `go test -race`.

7. `TestWaitForLSN_ContextTimeout`:
   - **Coverage**: Negative edge case for LSN wait timeout.
   - **Verification**: Confirms `WaitForLSN` returns `context.DeadlineExceeded` when waiting for an unreached LSN.

## Execution Verification

- `go test -v ./...` -> PASS (2.73s)
- `go test -race ./...` -> PASS (3.79s)
- `go run ./cmd/demo` -> PASS (Output verified clean and matching claims)

All tests are deterministic, robust, and pass without race conditions or memory leaks.
