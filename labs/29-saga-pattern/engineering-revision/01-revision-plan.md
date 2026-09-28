# Engineering Revision Plan

Target Lab: labs/29-saga-pattern
Previous Verdict: APPROVED (with 1 non-blocking test gap warning)

## Blocking Issues
None.

## Non-Blocking Issues
1. `tests/saga_test.go`: Choreography architecture lacked unit test coverage for failure and compensation event handling.

## Files To Change
- `tests/saga_test.go`: Add `TestChoreography_FailureCompensates` to exercise choreography compensation flow.

## Tests To Add/Modify
- Add `TestChoreography_FailureCompensates`: Tests `OrderCreated` -> `PaymentCompleted` -> `InventoryFailed` -> triggers compensation (`RefundPayment` + `CancelOrder`), asserting terminal state `OrderCancelled` and refund status.

## Validation Commands
```bash
go test -v ./...
go test -race ./...
go run ./cmd/demo
```
