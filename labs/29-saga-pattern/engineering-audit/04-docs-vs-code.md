# Docs vs Code Audit

## Comparisons

### 1. README vs Code
- **Claimed Components:**
  - `internal/saga/orchestrator.go`: Present and matches description.
  - `internal/saga/choreography.go`: Present and matches description.
  - `internal/services/services.go`: Present and matches description.
  - `cmd/demo/main.go`: Present and matches description.
  - `tests/saga_test.go`: Present and matches description.
- **Commands:**
  - `go test -v ./...`: Runs cleanly.
  - `go test -race ./...`: Passes cleanly.
  - `go run ./cmd/demo`: Runs and matches documented scenarios.
- **Assessment:** MATCH

### 2. Engineering Notes vs Code
- **LIFO Compensation:** Implemented in `Orchestrator.compensate` iterating `len(executed)-1` down to 0. Matches design.
- **Idempotency Support:** Implemented via `processedID` map in `PaymentService`. Matches design.
- **Semantic Lock Countermeasure:** Implemented via `locks` map in `OrderService`. Matches design.
- **Limitations:** Correctly notes in-memory simulation and standard library scope.
- **Assessment:** MATCH

### 3. Demo Output vs Code Execution
- **Expected Console Output:**
  Matches actual console output character-for-character with valid state transitions (Happy path: APPROVED, Stock 0; Failure path: CANCELLED, HasPayment false, Stock 0).
- **Assessment:** MATCH
