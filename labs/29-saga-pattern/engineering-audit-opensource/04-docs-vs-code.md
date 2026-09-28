# Documentation vs Code

## README.md vs Code
- README claims: orchestrator LIFO compensation, event bus choreography, services with idempotency + semantic locks, console demo, test suite.
  Code provides all of the above.
  Assessment: PASS.

- README lists components `internal/saga/orchestrator.go`, `internal/saga/choreography.go`, `internal/services/services.go`, `cmd/demo/main.go`, `tests/saga_test.go`.
  All exist. 
  Assessment: PASS.

- README commands `go test -v ./...` and `go test -race ./...` both execute successfully.
  Assessment: PASS.

## Engineering docs vs Code (internal docs)
- engineering/01-design.md lists `pkg/saga` and `pkg/services` as architecture paths; actual code is `internal/saga` and `internal/services`.
  Type: DOC_CODE_MISMATCH
  Severity: LOW.

- engineering/01-design.md lists orchestrator steps as Order -> Payment -> Inventory -> Delivery; actual steps are CreateOrder, ProcessPayment, ReserveInventory, ApproveOrder.
  Type: DOC_CODE_MISMATCH
  Severity: LOW.

- engineering/02-implementation-notes.md Known Limitations: mentions no compensation retry/backoff, in‑memory persistence simulation – matches code reality.
  Assessment: PASS.

## Execution Results (engineering/03-execution-result.md) vs Actual
- Recorded build: SUCCESS → verified `go build ./...` succeeds.
  Assessment: PASS.

- Recorded test results (9 tests shown) → all pass; actual run shows 9 tests passing (recorded 7; actual adds TestOrchestrator_CompensationErrorPropagated and TestOrchestrator_ContextCancellation). Recorded result is incomplete but not fabricated – it is stale relative to final test set.
  Type: DOC_CODE_MISMATCH (stale test listing)
  Severity: LOW.

- Recorded race detector result `ok` → verified `go test -race ./...` passes.
  Assessment: PASS.

- Recorded demo output → `go run ./cmd/demo` output matches recorded verbatim.
  Assessment: PASS (FAKE_DEMO ruled out).

## Research vs Implementation
- Research/claim (design doc) describes LIFO rollback, idempotency, semantic locking. Code implements all three.
  Assessment: PASS.

## Overall
- README is accurate against code.
- Internal engineering docs drift slightly from code (paths, step names, stale test listing). These are documentation issues, not implementation deception.
  Assessment: WARNING (documentation accuracy).