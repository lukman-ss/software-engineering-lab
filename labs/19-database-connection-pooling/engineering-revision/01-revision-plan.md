# Engineering Revision Plan

Target Lab: labs/19-database-connection-pooling
Previous Verdict: APPROVED (no blocking issues; 5 non-blocking LOW-severity findings)

## Blocking Issues

None.

## Non-Blocking Issues

1. **GAP-01** (DOC_CODE_MISMATCH, LOW): `engineering/03-execution-result.md` lists 6 tests; actual suite has 8.
2. **GAP-02** (MISSING_TEST, LOW): `ErrAcquireTimeout` exported but never returned and never tested — dead code.
3. **GAP-03** (MISSING_TEST, LOW): No test for `ProcessOrderUnsafeLeak` when initial `ExecContext` fails.
4. **GAP-04** (MISSING_EDGE_CASE, LOW): No test for context-cancelled during slow `connectDelay`. Acceptable for mock driver — skip.
5. **GAP-05** (MISSING_TEST, LOW): `TotalCreated()` untested as proof of pool reuse.

## Files To Change

- `internal/pool/service.go` — remove dead `ErrAcquireTimeout`
- `tests/pool_test.go` — add tests for GAP-03 and GAP-05
- `engineering/03-execution-result.md` — update test list from 6 to 8 (+ new tests)

## Tests To Add/Modify

- `TestUnsafeLeakExecContextFailure`: verifies that when MockDriver rejects the ExecContext query (via server overload), `ProcessOrderUnsafeLeak` returns the error and releases the connection.
- `TestTotalCreatedPoolReuse`: verifies pooled reuse results in fewer `TotalCreated` than unpooled.

## Validation Commands

```bash
cd labs/19-database-connection-pooling
go build ./...
go test -v ./...
go test -race ./...
go run ./cmd/demo
```
