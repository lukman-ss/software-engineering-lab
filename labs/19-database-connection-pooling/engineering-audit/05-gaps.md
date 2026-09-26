# Gap Analysis

## 1. RACE_CONDITION
- **Location:** `internal/pool/mockdb.go:41`
- **Description:** `MockDriver.Open()` reads `d.activeConns` non-atomically while `mockConn.Close()` writes to it atomically. This causes a confirmed data race when `go test -race` is executed.
- **Severity:** HIGH
- **Resolution:** Modify `MockDriver.Open()` to use `atomic.LoadInt32(&d.activeConns)` when evaluating the `if d.activeConns >= d.maxConnections` condition.

## 2. MISSING_EDGE_CASE
- **Location:** `tests/pool_test.go`
- **Description:** The research identified pool-locking deadlock (Finding 11) where a single thread requires multiple connections concurrently. The lab does not test or demonstrate this failure mode.
- **Severity:** LOW
- **Resolution:** Add a demo or test case showing pool exhaustion caused by a transaction requiring two connections from a strictly sized pool.

## 3. IMPLEMENTATION_OVERCLAIM
- **Location:** `engineering/01-design.md` vs `internal/pool/mockdb.go`
- **Description:** The design lists "throughput degradation" as a concept to prove, but the driver only implements a hard limit cutoff (`ErrServerOverloaded`). The performance "knee" is not simulated.
- **Severity:** MEDIUM
- **Resolution:** Either update the engineering docs to reflect only max connection enforcement, or update `MockDriver` to introduce artificial `time.Sleep` latency that scales non-linearly with `d.activeConns` exceeding optimal levels.
