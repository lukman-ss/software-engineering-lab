# Docs vs Code Audit

## Overview

Comparison of `README.md`, `engineering/01-design.md`, `engineering/02-implementation-notes.md` against codebase in `internal/` and `tests/`.

## Comparisons

1. **Components List in README.md vs Actual Files**
   - Listed:
     - `internal/saga/orchestrator.go` -> Exists and matches.
     - `internal/saga/choreography.go` -> Exists and matches.
     - `internal/services/services.go` -> Exists and matches.
     - `cmd/demo/main.go` -> Exists and matches.
     - `tests/saga_test.go` -> Exists and matches.
   - Assessment: PASS

2. **Execution Instructions**
   - Listed commands:
     - `go test -v ./...` -> Runs cleanly.
     - `go test -race ./...` -> Runs cleanly.
     - `go run ./cmd/demo` -> Runs cleanly and prints demo outputs.
   - Assessment: PASS

3. **Design Claims vs Implementation**
   - Claim: LIFO compensation order.
     - Verified: `orchestrator.go` loops backwards through executed steps.
   - Claim: Decoupled Choreography via EventBus.
     - Verified: `EventBus` pub/sub handles async event dispatch and handlers.
   - Claim: Semantic locking on order entities.
     - Verified: `OrderService` locks order keys on create and releases on approve/cancel.
   - Claim: Idempotency keys on payments.
     - Verified: `PaymentService` checks idempotency map on processed IDs.
   - Assessment: PASS
