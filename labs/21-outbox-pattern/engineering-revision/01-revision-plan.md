# Engineering Revision Plan

Target Lab: labs/21-outbox-pattern
Previous Verdict: APPROVED / APPROVED_WITH_WARNINGS

## Blocking Issues
None.

## Non-Blocking Issues
1. Concurrency test in `tests/outbox_test.go` used identical order IDs and lacked assertions on published/processed counts.
2. Consumer concurrency not exercised under `-race` test suite.
3. Minor documentation discrepancy in `engineering/01-design.md` referencing SQLite DB schema diagram instead of custom in-memory transactional database.

## Files To Change
- `tests/outbox_test.go`
- `engineering/01-design.md`

## Tests To Add/Modify
- Modify `TestTransactionalOutbox_ConcurrentWrites`: Unique order IDs per worker/iteration, assert total messages published and pending messages queue drained.
- Add `TestTransactionalOutbox_ConcurrentConsumers`: Test concurrent message handling across multiple goroutines with duplicate event IDs under `-race`.

## Validation Commands
```bash
go test -v ./...
go test -race ./...
go run ./cmd/demo
```
