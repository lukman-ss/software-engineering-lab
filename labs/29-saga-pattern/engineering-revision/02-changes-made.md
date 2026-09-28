## Revision 1

Audit Issue: Non-blocking warning: Choreography model lacked failure & compensation unit test.
Severity: LOW
Files Changed: `tests/saga_test.go`
Action: Added `TestChoreography_FailureCompensates` subscribing to `InventoryFailed` and asserting compensation actions (`RefundPayment` and `CancelOrder`) roll back order state to `OrderCancelled` without retaining payment records.
Verification: Ran `go test -v ./...` and `go test -race ./...`. All tests passed.
Status: RESOLVED

## Revision 2

Audit Issue: Compensation error handling ignored in orchestrator.
Severity: MEDIUM
Files Changed: `internal/saga/orchestrator.go`, `tests/saga_test.go`
Action: Updated `compensate` to aggregate non-nil errors from `step.Compensate(ctx)`, mark step status as `COMPENSATE_FAILED`, and surface aggregated error to `Execute` caller. Added `TestOrchestrator_CompensationErrorPropagated`.
Verification: Ran `go test -v ./...` and `go test -race ./...`. All tests passed.
Status: RESOLVED

## Revision 3

Audit Issue: Missing context cancellation check during step execution.
Severity: LOW
Files Changed: `internal/saga/orchestrator.go`, `tests/saga_test.go`
Action: Added `select { case <-ctx.Done(): ... }` check before executing each step. On cancellation, triggers compensation for executed steps and returns cancellation error. Added `TestOrchestrator_ContextCancellation`.
Verification: Ran `go test -v ./...` and `go test -race ./...`. All tests passed.
Status: RESOLVED
