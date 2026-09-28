# Code Audit

## Finding 1

Location: internal/services/services.go:42-61 (ApproveOrder/CancelOrder)
Claimed Behavior: Services enforce valid state transitions and reject invalid operations (e.g., approving a non-existent order).
Observed Implementation: 
- ApproveOrder(orderID) only checks if current state == OrderCancelled before setting OrderApproved. It does not verify the order exists (via presence in s.orders map). For a non-existent orderID, s.orders[orderID] returns zero value "", which != OrderCancelled, so it creates a phantom APPROVED entry.
- CancelOrder(orderID) similarly sets OrderCancelled without checking existence, creating a phantom CANCELLED entry.
Assessment: FAIL
Severity: MEDIUM
Notes: While the saga workflow always calls CreateOrder first, exposing the service API with missing validation risks incorrect usage outside the saga. This is an unhandled error case and missing edge case guard. Fix: add existence check (e.g., if _, exists := s.orders[orderID]; !exists) and return error.

## Finding 2

Location: internal/saga/orchestrator.go:61-67
Claimed Behavior: On context cancellation, the orchestrator should halt execution, compensate completed steps, and accurately record failure status.
Observed Implementation: The select statement checks ctx.Done() before each step. If cancelled, it logs the current step as StatusFailed (despite never executing it), then calls compensate() on already-executed steps. 
Assessment: WARNING
Severity: LOW
Notes: Logging a never-executed step as FAILED misrepresents semantics; StepStatus.FAILED implies execution failure. The step was skipped due to cancellation, not failure. A more accurate status (e.g., a new StatusCancelled) would improve clarity, but the compensation behavior is correct. TestOrchestrator_ContextCancellation does not validate logs, so this issue is not caught by tests.

## Finding 3

Location: internal/saga/orchestrator.go:109
Claimed Behavior: When multiple compensations fail, errors should be aggregated and returned in a usable format.
Observed Implementation: `return fmt.Errorf("%v", compErrors)` where compErrors is []error. This produces a string like `[err1 err2]`, which does not permit unwrapping individual errors and is not idiomatic.
Assessment: WARNING
Severity: LOW
Notes: Prefer errors.Join(compErrors...) (Go 1.20+) or a formatted list that preserves error chain. The current format is functional but limits error inspection.

## Finding 4

Location: engineering/01-design.md:26,30,35
Claimed Behavior: Architecture uses pkg/saga and pkg/services; orchestrator steps include Order -> Payment -> Inventory -> Delivery.
Observed Implementation: Code resides in internal/saga and internal/services. Orchestrator steps are CreateOrder, ProcessPayment, ReserveInventory, ApproveOrder (no Delivery step).
Assessment: WARNING
Severity: LOW
Notes: Design doc drifts from actual implementation (pkg vs internal path, missing/incorrect step names). This is a documentation accuracy issue but does not affect code correctness.

## Finding 5

Location: tests/saga_test.go (various)
Claimed Behavior: Test suite verifies idempotency ensures identical state on duplicate executions.
Observed Implementation: TestPayment_Idempotency calls ProcessPayment twice with same ID and asserts no error, but does not verify that the payment amount remains unchanged (i.e., no double-charge).
Assessment: WARNING
Severity: LOW
Notes: The implementation is idempotent (early return via processedID map), but the test lacks an assertion on state equivalence. Strengthen by checking s.payments[id] == expectedAmount after both calls.

## Finding 6

Location: tests/saga_test.go
Claimed Behavior: Test suite covers edge cases: nil Compensate functions, semantic lock interactions, over-release, choreography PaymentFailed path.
Observed Implementation: Missing tests for:
- A middle step with nil Compensate (should be skipped during rollback).
- ApproveOrder after CancelOrder (should fail or be no-op).
- Release reserving more than currently reserved within a saga.
- Choreography path where PaymentFailed triggers compensation (order cancellation only partially tested).
- Context cancellation with timeout (only WithCancel tested).
Assessment: WARNING
Severity: LOW
Notes: While core behavior is verified, these gaps reduce confidence in edge‑case handling. Adding them would strengthen the suite.