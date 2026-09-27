# Source Map

## Lost Update Anomaly (Problem & Core Concept)

Research:
- research/05-report.md — Finding 1 (Lost Update under READ COMMITTED)

Implementation:
- internal/inventory/store.go — `NaiveDeduct()` method (lines 65-89)

Tests:
- tests/locking_test.go — `TestNaiveLostUpdate()` (lines 10-38)

Demo:
- cmd/demo/main.go — Scenario [1] Naive Read-Modify-Write (lines 16-34)

---

## Pessimistic Locking (SELECT ... FOR UPDATE)

Research:
- research/05-report.md — Finding 2 (Pessimistic Locking Prevents Conflicts by Blocking)
- research/05-report.md — Finding 6 (anti-pattern: holding pessimistic locks too long, deadlock risk)

Implementation:
- internal/inventory/store.go — `PessimisticDeduct()` method (lines 93-116)
- internal/inventory/store.go — `GetRowLock()` method (lines 42-51)

Tests:
- tests/locking_test.go — `TestPessimisticLocking()` (lines 40-67)
- tests/locking_test.go — `TestPessimisticLockingInsufficientStock()` (lines 69-83)

Demo:
- cmd/demo/main.go — Scenario [2] Pessimistic Locking (lines 36-53)

---

## Optimistic Locking (Version Guard)

Research:
- research/05-report.md — Finding 4 (Optimistic Locking Detects Conflicts at Commit Time)

Implementation:
- internal/inventory/store.go — `OptimisticDeduct()` method (lines 120-152)
- internal/inventory/service.go — `DeductOptimisticDirect()` (lines 24-26)
- internal/inventory/service.go — `DeductOptimisticWithRetry()` (lines 28-45)
- internal/inventory/model.go — `ErrOptimisticLock` error (line 8)

Tests:
- tests/locking_test.go — `TestOptimisticLockingConflict()` (lines 85-121)
- tests/locking_test.go — `TestOptimisticLockingWithRetry()` (lines 123-145)

Demo:
- cmd/demo/main.go — Scenario [3] Optimistic Direct (lines 55-74)
- cmd/demo/main.go — Scenario [4] Optimistic With Retry (lines 76-97)

---

## Atomic Single-Statement Update

Research:
- research/05-report.md — Finding 5 (Atomic Database Operations Eliminate Race Windows)
- research/05-report.md — Finding 1/6 (Isolation Level Alone Does Not Prevent Lost Update)

Implementation:
- internal/inventory/store.go — `AtomicDeduct()` method (lines 155-173)
- internal/inventory/service.go — `DeductAtomic()` (lines 47-49)

Tests:
- tests/locking_test.go — `TestAtomicConditionalUpdate()` (lines 147-171)

Demo:
- cmd/demo/main.go — Scenario [5] Atomic Operation (lines 99-116)

---

## Isolation Level Limitation

Research:
- research/05-report.md — Finding 3 (Database Isolation Levels Provide Alternative Control)
- research/05-report.md — Finding 6 (Common Mistakes: Assuming Transactions Prevent Lost Updates)

---

## Anti-Patterns

Research:
- research/05-report.md — Finding 6 (Common Anti-Patterns Are Documented)

---

## MVCC Foundation

Research:
- research/05-report.md — Finding 7 (MVCC Is the Foundation Enabling Both Approaches)

---

## Model & Errors

Source File:
- internal/inventory/model.go — Product struct, error definitions

---

## Service Layer

Source File:
- internal/inventory/service.go — Business logic wrappers for all strategies

---

## Store Engine (Central Logic)

Source File:
- internal/inventory/store.go — Simulated storage engine with locking, versioning, atomic operations

---

## Research Status Verification

- research-audit/07-verdict.md — APPROVED (0 unsupported claims, 3 non-blocking gaps: MySQL docs 403, atomic update scope boundary, missing benchmarks)
- research-audit/03-claim-audit.md — All 10 claims verified (8 FACT, 2 PARTIAL)
- research-audit/06-gaps.md — 3 gaps (1 MEDIUM, 2 LOW), all non-blocking

---

## Engineering Status Verification

- engineering-audit/06-verdict.md — APPROVED
  - All 6/6 tests PASS
  - Race detector: PASS (zero warnings)
  - Demo: PASS (verified real output)
  - Research alignment: PASS
  - Documentation accuracy: WARNING (README uses SQL syntax while implementation uses in-memory simulation; documented in design)
- engineering-audit/03-test-audit.md — 6 tests covering naive anomaly, pessimistic success, pessimistic insufficient stock, optimistic conflict, optimistic retry convergence, atomic success
- engineering-audit/05-gaps.md — 4 test-coverage gaps (non-blocking):
  1. MISSING_TEST: `ErrInvalidQuantity` not tested for any strategy
  2. MISSING_TEST: `ErrNotFound` not tested for missing product
  3. MISSING_TEST: Optimistic retry exhaustion (maxRetries reached) not tested
  4. MISSING_TEST: Atomic / pessimistic overdraft under concurrency not tested
- engineering-audit-opensource/06-verdict.md — APPROVED (same findings + 1 doc/code terminology mismatch)

---

## Execution Verification

- engineering/03-execution-result.md — Test results, race detector output, demo output with concrete figures
- engineering/01-design.md — Architecture design, success criteria, test strategy
- engineering/02-implementation-notes.md — Implementation decisions, trade-offs, limitations
- engineering-revision/03-revision-result.md — READY_FOR_ENGINEERING_REAUDIT (no issues resolved; previous APPROVED stands)

---

## Research Findings Reference (from research/05-report.md)

| Finding | Topic |
|---------|-------|
| 1 | Lost Update Not Prevented by Default Isolation Levels |
| 2 | Pessimistic Locking Prevents Conflicts by Blocking |
| 3 | Database Isolation Levels Provide Alternative Control |
| 4 | Optimistic Locking Detects Conflicts at Commit Time |
| 5 | Atomic Database Operations Eliminate Race Windows |
| 6 | Common Mistakes Are Well-Documented |
| 7 | MVCC Is the Foundation Enabling Both Approaches |

Note: The report contains exactly 7 findings (not 9). Earlier drafts incorrectly mapped additional sections as separate findings.