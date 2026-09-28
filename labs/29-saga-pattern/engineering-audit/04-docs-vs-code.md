# Docs vs Code Comparison

## Document Review
- `README.md`
- `engineering/01-design.md`
- `engineering/02-implementation-notes.md`
- `engineering/03-execution-result.md`

## Comparison Matrix

| Claim in Docs | Location in Code | Status | Details |
|---|---|---|---|
| Orchestrator handles forward steps and LIFO compensation | `internal/saga/orchestrator.go` | MATCH | `Execute` and `compensate` methods implement exact LIFO rollback. |
| Event bus for choreography | `internal/saga/choreography.go` | MATCH | `EventBus` provides `Subscribe` and `Publish`. |
| Domain services with local state, idempotency, semantic locks | `internal/services/services.go` | MATCH | `OrderService` (locks), `PaymentService` (idempotency key), `InventoryService` (stock & reserve). |
| Console demo runnable with `go run ./cmd/demo` | `cmd/demo/main.go` | MATCH | Real demo executes scenarios 1 and 2, producing exact expected output. |
| Running tests instructions (`go test -v ./...`, `go test -race ./...`) | `tests/saga_test.go` | MATCH | Tests pass cleanly under both flags. |

## Mismatch Findings
- DOC_CODE_MISMATCH: None.
- TEST_CLAIM_MISMATCH: None.
- RESEARCH_IMPLEMENTATION_MISMATCH: None.

Assessment: PASS
Severity: LOW
