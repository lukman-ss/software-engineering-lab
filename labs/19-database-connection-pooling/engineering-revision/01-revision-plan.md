# Engineering Revision Plan

Target Lab: labs/19-database-connection-pooling
Previous Verdict: APPROVED_WITH_WARNINGS

## Blocking Issues

None. No CRITICAL or HIGH issues found in audit.

## Non-Blocking Issues

| ID | Severity | Description |
|---|---|---|
| GAP-004 | MEDIUM | TestConnectionStarvationDueToLeak uses 10ms sleep for goroutine sync — fragile under load |
| GAP-001 | LOW | externalCall error path not tested in ProcessOrderSafe or ProcessOrderUnsafeLeak |
| GAP-002 | LOW | Pre-cancelled context path not tested in ProcessOrderSafe |
| GAP-003 | LOW | No test for connection release when ProcessOrderUnsafeLeak externalCall fails (by inspection, defer conn.Close() handles it) |
| GAP-005 | LOW | TestOversizedPoolExhaustsServerConnections asserts errCount > 0 instead of >= 10 |
| GAP-006 | LOW | connectDelay tested only indirectly (TestDirectConnectionOverhead covers this adequately) |

## Files To Change

- `tests/pool_test.go` — fix GAP-004, GAP-005, add tests for GAP-001, GAP-002

## Tests To Add/Modify

1. **Modify** `TestConnectionStarvationDueToLeak`: replace `time.Sleep(10ms)` with channel signal confirming goroutine acquired connection
2. **Modify** `TestOversizedPoolExhaustsServerConnections`: strengthen assertion from `errCount > 0` to `errCount >= 10`
3. **Add** `TestExternalCallErrorPropagation`: covers externalCall returning error in both ProcessOrderSafe and ProcessOrderUnsafeLeak
4. **Add** `TestPreCancelledContextProcessOrderSafe`: covers already-cancelled context passed to ProcessOrderSafe

## Validation Commands

```bash
cd labs/19-database-connection-pooling
go test ./...
go test -race -count=1 ./...
go run ./cmd/demo
```
