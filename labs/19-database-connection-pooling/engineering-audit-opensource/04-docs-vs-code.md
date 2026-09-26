# Docs vs Code

Target Lab: labs/19-database-connection-pooling

## README vs Code

### README: "internal/pool/mockdb.go: Implements a Go database/sql/driver to simulate database server constraints (max_connections)
Code: MockDriver struct enforces maxConnections limit in Open(). ✓ MATCH

### README: "internal/pool/service.go: Provides services executing safe operations vs unsafe operations
Code: ProcessOrderSafe (external I/O before DB op) and ProcessOrderUnsafeLeak (I/O during DB lock). ✓ MATCH

### README: "tests/pool_test.go: Automated test suite covering direct overhead, pool exhaustion, connection leaks, and concurrent operations
Code: 4 tests covering exactly these scenarios. ✓ MATCH

### README: "cmd/demo/main.go: Interactive CLI demo demonstrating connection overhead, rejection, and pool starvation
Code: 3 demo functions demonstrating each. ✓ MATCH

### README: "go test -v ./... / go test -race -v ./..."
Verified: Both commands pass. ✓ MATCH

## Engineering Design vs Code

### Design: "Direct connection creation penalty (overhead)"
Code: TestDirectConnectionOverhead + demo. ✓ MATCH

### Design: "Connection exhaustion when pools oversized beyond hardware limits"
Code: TestOversizedPoolExhaustsServerConnections + demo. ✓ MATCH

### Design: "Connection leaks when connections held during external I/O"
Code: ProcessOrderUnsafeLeak + TestConnectionStarvationDueToLeak + demo. ✓ MATCH

### Design Success Criteria: "Test validates pool sizes and limits"
Code: TestOversizedPoolExhaustsServerConnections. ✓ MATCH

### Design Success Criteria: "Test verifies connection leakage triggers errors"
Code: TestConnectionStarvationDueToLeak. ✓ MATCH

### Design Success Criteria: "Test validates throughput degradation or max connection enforcement"
Code: TestDirectConnectionOverhead (throughput degradation). ✓ MATCH

## Engineering Execution-Result vs Actual Execution

| Item | Claimed | Actual | Match |
|------|---------|--------|-------|
| Build | Success, no output | Success, no output | ✓ |
| Tests | 4 PASS | 4 PASS | ✓ |
| Race Detector | 4 PASS | 4 PASS | ✓ |
| Demo: Unpooled | 54.59ms | 55.25ms | ✓ |
| Demo: Pooled | 6.33µs | 2.46µs | ✓ |
| Demo: Oversized | 30, 15, 15 | 30, 15, 15 | ✓ |
| Demo: Leak | Order 3 Failed: context deadline exceeded | Order 3 Failed: context deadline exceeded | ✓ |

## Mismatches Found

None. README, design notes, execution-result, code, tests, and demo all align.

## Notes

Demo timing values differ slightly from engineering/03-execution-result.md (expected — timing varies by run). Core outcomes identical.