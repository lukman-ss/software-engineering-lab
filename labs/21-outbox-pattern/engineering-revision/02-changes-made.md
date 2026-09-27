# Engineering Changes Made

Target Lab: labs/21-outbox-pattern
Previous Verdict: APPROVED_WITH_WARNINGS

## Revision 1

Audit Issue: `Relay.Start()` and `Relay.Stop()` lifecycle safety (prevent multiple start goroutine leaks and double close panic).
Severity: MEDIUM
Files Changed:
- `internal/outbox/relay.go`
Action: Added `startOnce` and `stopOnce` `sync.Once` guards to `Relay.Start()` and `Relay.Stop()`.
Verification:
- `go test -v ./...` passed.
- `go test -race ./...` passed without race or panic.
Status: RESOLVED

## Revision 2

Audit Issue: Lack of test verifying relay retry and delivery after transient broker failure.
Severity: LOW
Files Changed:
- `tests/outbox_test.go`
Action: Added `TestTransactionalOutbox_RelayRetryAfterBrokerFailure` simulating broker failure on initial poll and asserting successful recovery, dispatch, and consumer processing on subsequent poll.
Verification:
- `go test -v -run TestTransactionalOutbox_RelayRetryAfterBrokerFailure ./...` passed.
- `go test -race ./...` passed.
Status: RESOLVED
