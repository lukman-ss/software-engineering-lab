# Docs vs Code Audit

## README vs Implementation

| README Claim | Code Reality | Status |
|---|---|---|
| `internal/pool/mockdb.go`: simulates database server constraints (max_connections) and connection establishment latency | Implemented via `maxConnections int32` and `connectDelay time.Duration` in MockDriver | MATCH |
| `internal/pool/service.go`: safe vs unsafe operations (holding DB connection during external I/O) | `ProcessOrderSafe` and `ProcessOrderUnsafeLeak` both implemented correctly | MATCH |
| `tests/pool_test.go`: covers direct overhead, pool exhaustion, connection leaks, concurrent ops | 10 tests covering all stated categories | MATCH |
| `cmd/demo/main.go`: interactive CLI demo demonstrating connection overhead, rejection, and pool starvation | All three demos implemented: overhead, rejection, starvation | MATCH |
| `go test -v ./...` command works | Verified: all PASS | MATCH |
| `go test -race -v ./...` command works | Verified: all PASS, no races | MATCH |
| `go run ./cmd/demo` command works | Verified: runs and produces expected output | MATCH |

## Engineering Design vs Implementation

| Design Claim | Code Reality | Status |
|---|---|---|
| Custom `database/sql/driver` to simulate limits and overhead | `MockDriver` implements `driver.Driver`; `MockConnector` implements `driver.Connector` | MATCH |
| `SetMaxOpenConns`/`SetMaxIdleConns` controls pooling | Used in tests and demo | MATCH |
| Happy path: concurrent requests with fixed pool size | `TestSafeProcessingConcurrently` | MATCH |
| Failure path: oversized pool exceeds server limits | `TestOversizedPoolExhaustsServerConnections` | MATCH |
| Edge case: connection leak via blocked external IO | `TestConnectionStarvationDueToLeak` | MATCH |
| Run with race detector | Executed, clean | MATCH |

## Engineering Execution Results vs Actual Audit Results

| Metric | Claimed (engineering/03-execution-result.md) | Actual (this audit) |
|---|---|---|
| All tests pass | PASS | PASS |
| Race detector clean | PASS | PASS |
| Demo unpooled 5 requests | ~54ms | ~55ms |
| Demo pooled 5 requests | ~6µs | ~13µs |
| Oversized pool: succeeded | 15 | 15 |
| Oversized pool: rejected | 15 | 15 |
| Starvation demo result | "Order 3 Failed: context deadline exceeded" | "Order 3 Failed: context deadline exceeded" |

All demo output values match within expected runtime variance.

## Issues Found

### DOC_CODE_MISMATCH (Minor)

**Location**: go.mod:3  
**Claimed**: `go 1.26.7`  
**Reality**: Go 1.26.7 does not exist (latest stable at audit date is 1.23.x). This is a speculative future version number. Build still succeeds because the installed toolchain handles the directive gracefully.  
**Severity**: LOW — does not affect functionality or reproducibility.

No other mismatches found between README, engineering notes, and implementation.
