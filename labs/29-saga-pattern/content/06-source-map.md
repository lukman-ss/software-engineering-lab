# Source Map

## Problem & Motivation (2PC vs Saga)
- **Research:** `research/03-evidence.md` (Evidence 1), `research/05-report.md` (Finding 1)
- **Implementation:** `engineering/01-design.md`, `engineering/02-implementation-notes.md`
- **Tests:** `tests/saga_test.go`

## Saga Definition & Local Transactions
- **Research:** `research/03-evidence.md` (Evidence 2, 10), `research/05-report.md` (Finding 2)
- **Implementation:** `internal/saga/orchestrator.go`, `tests/saga_test.go`
- **Tests:** `TestOrchestrator_HappyPath`

## Compensating Transactions & LIFO Rollback
- **Research:** `research/03-evidence.md` (Evidence 3, 9, 9a), `research/05-report.md` (Finding 4, 9)
- **Implementation:** `internal/saga/orchestrator.go` (`compensate` function)
- **Tests:** `tests/saga_test.go` (`TestOrchestrator_FailureCompensatesLIFO`, `TestOrchestrator_CompensationErrorPropagated`)

## Orchestration vs. Choreography
- **Research:** `research/03-evidence.md` (Evidence 4), `research/05-report.md` (Finding 3)
- **Implementation:** `internal/saga/orchestrator.go`, `internal/saga/choreography.go`
- **Tests:** `tests/saga_test.go` (`TestChoreography_Flow`, `TestChoreography_FailureCompensates`)

## Idempotency & Retryable Transactions
- **Research:** `research/03-evidence.md` (Evidence 5, 8), `research/05-report.md` (Finding 5, 8)
- **Implementation:** `internal/services/services.go` (`ProcessPayment` idempotency key)
- **Tests:** `tests/saga_test.go` (`TestPayment_Idempotency`)

## Lack of Isolation & Semantic Locks (Countermeasures)
- **Research:** `research/03-evidence.md` (Evidence 6, 7, 11), `research/05-report.md` (Finding 6, 7)
- **Implementation:** `internal/services/services.go` (`OrderService` semantic locks)
- **Tests:** `tests/saga_test.go` (`TestSemanticLock`)
