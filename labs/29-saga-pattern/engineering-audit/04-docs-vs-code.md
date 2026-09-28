# Docs vs Code Audit

## Consistency Checklist

1. **README.md Components vs File Tree**
   - Listed files:
     - `internal/saga/orchestrator.go` -> Exists & matches description.
     - `internal/saga/choreography.go` -> Exists & matches description.
     - `internal/services/services.go` -> Exists & matches description.
     - `cmd/demo/main.go` -> Exists & matches description.
     - `tests/saga_test.go` -> Exists & matches description.
   - Status: PASS

2. **README.md Commands vs Execution Behavior**
   - `go test -v ./...` -> Executes and passes.
   - `go test -race ./...` -> Executes cleanly with race detector.
   - `go run ./cmd/demo` -> Runs without error.
   - Status: PASS

3. **Design / Implementation Notes vs Actual Code**
   - Claim: LIFO compensation rollback implemented.
   - Reality: Observed in `internal/saga/orchestrator.go:94` and verified in `TestOrchestrator_FailureCompensatesLIFO`.
   - Claim: Idempotency keys in payment processing.
   - Reality: Observed in `internal/services/services.go:86` and verified in `TestPayment_Idempotency`.
   - Claim: Semantic lock countermeasure on OrderService.
   - Reality: Observed in `internal/services/services.go:33` and verified in `TestSemanticLock`.
   - Claim: Choreography model with event bus.
   - Reality: Observed in `internal/saga/choreography.go` and verified in `TestChoreography_Flow`.
   - Status: PASS

4. **Discrepancy Findings**
   - None found.
