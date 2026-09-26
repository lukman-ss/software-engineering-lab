# Source Map

## Lost Update Anomaly (Problem & Core Concept)

Research:
research/05-report.md — Finding 1 (Lost Update under READ COMMITTED)

Implementation:
internal/inventory/store.go — `NaiveDeduct()` method (lines 65-89)

Tests:
tests/locking_test.go — `TestNaiveLostUpdate()` (lines 10-38)

Demo:
cmd/demo/main.go — Scenario [1] Naive Read-Modify-Write (lines 16-34)

## Pessimistic Locking (SELECT ... FOR UPDATE)

Research:
research/05-report.md — Finding 2 (Pessimistic Locking Blocks via FOR UPDATE)
research/05-report.md — Finding 3 (Trade-offs: Deadlock, Hold-Time)

Implementation:
internal/inventory/store.go — `PessimisticDeduct()` method (lines 93-116)
internal/inventory/store.go — `GetRowLock()` method (lines 42-51)

Tests:
tests/locking_test.go — `TestPessimisticLocking()` (lines 40-67)
tests/locking_test.go — `TestPessimisticLockingInsufficientStock()` (lines 69-83)

Demo:
cmd/demo/main.go — Scenario [2] Pessimistic Locking (lines 36-53)

## Optimistic Locking (Version Guard)

Research:
research/05-report.md — Finding 4 (Optimistic Locking via Version Guard)

Implementation:
internal/inventory/store.go — `OptimisticDeduct()` method (lines 120-152)
internal/inventory/service.go — `DeductOptimisticDirect()` (lines 24-26)
internal/inventory/service.go — `DeductOptimisticWithRetry()` (lines 28-45)
internal/inventory/model.go — `ErrOptimisticLock` error (line 8)

Tests:
tests/locking_test.go — `TestOptimisticLockingConflict()` (lines 85-121)
tests/locking_test.go — `TestOptimisticLockingWithRetry()` (lines 123-145)

Demo:
cmd/demo/main.go — Scenario [3] Optimistic Direct (lines 55-74)
cmd/demo/main.go — Scenario [4] Optimistic With Retry (lines 76-97)

## Atomic Single-Statement Update

Research:
research/05-report.md — Finding 6 (Atomic Single-Statement Eliminates RMW Window)
research/05-report.md — Finding 7 (Isolation Level Alone Does Not Prevent Lost Update)

Implementation:
internal/inventory/store.go — `AtomicDeduct()` method (lines 155-173)
internal/inventory/service.go — `DeductAtomic()` (lines 47-49)

Tests:
tests/locking_test.go — `TestAtomicConditionalUpdate()` (lines 147-171)

Demo:
cmd/demo/main.go — Scenario [5] Atomic Operation (lines 99-116)

## Isolation Level Limitation

Research:
research/05-report.md — Finding 7 (Isolation Level Alone Does Not Prevent Lost Update)
research/05-report.md — Finding 8 (Default Isolation Levels Differ Across Databases)

## Anti-Patterns

Research:
research/05-report.md — Finding 9 (Common Anti-Patterns Are Documented)

## Model & Errors

Source File:
internal/inventory/model.go — Product struct, error definitions

## Service Layer

Source File:
internal/inventory/service.go — Business logic wrappers for all strategies

## Store Engine (Central Logic)

Source File:
internal/inventory/store.go — Simulated storage engine with locking, versioning, atomic operations

## Research Status Verification

- research-audit/07-verdict.md — APPROVED
- engineering-audit/06-verdict.md — APPROVED

## Execution Verification

- engineering/03-execution-result.md — Test results, race detector output, demo output
- engineering/01-design.md — Architecture design
- engineering/02-implementation-notes.md — Implementation decisions, trade-offs, limitations
