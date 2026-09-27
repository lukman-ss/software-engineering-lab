# Engineering Revision Plan

Target Lab: labs/21-outbox-pattern
Previous Verdict: APPROVED_WITH_WARNINGS

## Blocking Issues
None.

## Non-Blocking Issues
1. `Relay.Start()` / `Relay.Stop()` not lifecycle-safe: double-close of `stopChan` causes panic; multiple `Start()` leaks goroutines.
2. Missing test exercising relay retry and eventual delivery after transient broker failure.

## Files To Change
- `internal/outbox/relay.go`: Add `sync.Once` guards to `Start()` and `Stop()`.

## Tests To Add/Modify
- `tests/outbox_test.go`: Add `TestTransactionalOutbox_RelayRetryAfterBrokerFailure` verifying transient broker failure retry and `Start()`/`Stop()` idempotency.

## Validation Commands
```bash
cd labs/21-outbox-pattern
go test -v ./...
go test -race ./...
go run ./cmd/demo
```
