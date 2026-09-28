# Docs vs Code Audit

## Sources Reviewed

- README.md
- internal/saga/orchestrator.go
- internal/saga/choreography.go
- internal/services/services.go
- cmd/demo/main.go
- tests/saga_test.go

# README Claims vs Code

## Component Mappings

### README claim:
- `internal/saga/orchestrator.go`: Centralized saga coordinator managing forward steps and LIFO compensation.
### Code verification:
- Orchestrator struct with AddStep, Execute, compensate (LIFO reverse iteration), Logs. Verified LIFO via lines 94-95 (for i := len(executed) - 1; i >= 0; i--).
### Assessment: PASS

### README claim:
- `internal/saga/choreography.go`: Decoupled event bus for choreographing saga steps across independent services.
### Code verification:
- EventBus struct with Subscribe, Publish, handler slice copy under RWMutex.
### Assessment: PASS

### README claim:
- `internal/services/services.go`: Domain services (Order, Payment, Inventory) with local state, idempotency, and semantic locks.
### Code verification:
- OrderService (semantic lock), PaymentService (idempotency), InventoryService. All present with mutexes.
### Assessment: PASS

### README claim:
- `cmd/demo/main.go`: Console demo showing happy path and failure compensation.
### Code verification:
- Demo runs two scenarios: happy path and rollback on inventory failure.
### Assessment: PASS

### README claim:
- `tests/saga_test.go`: Suite of unit, rollback, idempotency, semantic lock, and concurrency tests.
### Code verification:
- Test file contains: happy path, rollback LIFO, idempotency, semantic lock, concurrency, choreography flow, failure compensation, compensation error propagation, context cancellation.
### Assessment: PASS

## Behavior Claims

### README title claim:
- "both Orchestration and Choreography models, compensating transactions (LIFO rollback), idempotency keys, and semantic locking countermeasures."
### Code verification:
- Orchestration: verified.
- Choreography: verified (EventBus + handlers).
- Compensating transactions / LIFO rollback: verified (orchestrator.compensate reverse loop; demo scenario 2 shows compensation; choreography failure compensation via InventoryFailed handler).
- Idempotency keys: verified (PaymentService.processedID).
- Semantic locking: verified (OrderService.locks).
### Assessment: PASS

## Mismatches Found

No major mismatches found. README accurately describes the implementation.

### Minor observations (not blocking):
- README describes choreography as "across independent services" but the codebase has only a single EventBus with no explicit service isolation. This is a simplification acceptable for a demo.
- README mentions "semantic locking countermeasures" but OrderService's ApproveOrder has a lock-leak risk on failure (see code audit Finding 2).
  - DOC_CODE_MISMATCH: None for README vs code.
  - TEST_CLAIM_MISMATCH: None.
  - RESEARCH_MISMATCH: None (research out of scope).

## Overall Assessment: PASS
