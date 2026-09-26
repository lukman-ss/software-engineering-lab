# Engineering Revision Plan

Target Lab: labs/21-outbox-pattern
Previous Verdict: APPROVED

## Blocking Issues
None.

## Non-Blocking Issues
1. Outbox table cleanup / purge functionality: Noted as omitted in audit warning. Added `PurgeProcessedOutbox` on `DB` to complete outbox lifecycle capabilities cleanly.

## Files To Change
- `internal/outbox/db.go`: Add `PurgeProcessedOutbox` helper to purge processed outbox events.
- `tests/outbox_test.go`: Add test verifying cleanup of processed outbox entries.

## Tests To Add/Modify
- `TestTransactionalOutbox_PurgeProcessed`: Verify `PurgeProcessedOutbox` removes PROCESSED entries while keeping PENDING records.

## Validation Commands
```bash
go test -v ./...
go test -count=1 -race ./...
go run ./cmd/demo
```
