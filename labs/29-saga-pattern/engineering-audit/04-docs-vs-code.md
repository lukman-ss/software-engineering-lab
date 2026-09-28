# Docs vs Code Audit

## Comparisons

1. **README Component Mapping**:
   - `README.md` lists `internal/saga/orchestrator.go`, `internal/saga/choreography.go`, `internal/services/services.go`, `cmd/demo/main.go`, and `tests/saga_test.go`.
   - All files exist and match documented descriptions.

2. **Execution Commands**:
   - `README.md` documents `go test -v ./...`, `go test -race ./...`, `go run ./cmd/demo`.
   - Verified: All commands run cleanly and produce exact expected output.

3. **Engineering Design vs Implementation**:
   - `engineering/01-design.md` specifies Orchestrator struct, EventBus struct, LIFO rollback, semantic locking, and idempotency tracking.
   - Implementation matches the design specification precisely.

4. **Demo Output Verification**:
   - `cmd/demo/main.go` runs two distinct scenarios: happy path forward execution and inventory-exhaustion rollback.
   - Verified output: Real stdout output matches the claims. No fake outputs detected.
