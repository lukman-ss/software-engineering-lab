## Revision 1

Audit Issue: DOC_CODE_MISMATCH in design diagram (SQLite reference)
Severity: LOW
Files Changed: `engineering/01-design.md`
Action: Updated architecture diagram and component description to reflect the in-memory transactional database implementation.
Verification: Inspection of `engineering/01-design.md`.
Status: RESOLVED

## Revision 2

Audit Issue: Weak concurrency testing in `TestTransactionalOutbox_ConcurrentWrites`
Severity: MEDIUM
Files Changed: `tests/outbox_test.go`
Action: Refactored test to generate unique order IDs per goroutine iteration and added assertions for total published messages and zero pending outbox state.
Verification: `go test -v ./...` & `go test -race ./...`
Status: RESOLVED

## Revision 3

Audit Issue: Untested concurrent consumer handling
Severity: MEDIUM
Files Changed: `tests/outbox_test.go`
Action: Added `TestTransactionalOutbox_ConcurrentConsumers` to stress consumer deduplication under race conditions.
Verification: `go test -race ./...`
Status: RESOLVED
