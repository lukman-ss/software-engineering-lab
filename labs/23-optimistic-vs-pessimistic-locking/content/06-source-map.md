# Source Map

## Problem & Mental Model

**Research:**
- `research/05-report.md` (Finding 1: Lost Update Anomaly, Finding 5: Selection Criteria)
- `research/03-evidence.md`

**Engineering:**
- `engineering/01-design.md` (Problem, Concept To Prove)

---

## How It Works — Naive Deduct (Lost Update Path)

**Implementation:**
- `internal/inventory/store.go:65-89` (`NaiveDeduct`)

**Tests:**
- `tests/locking_test.go:10-38` (`TestNaiveLostUpdate`)

**Audit:**
- `engineering-audit/02-code-audit.md` (Finding 1)
- `engineering-audit/03-test-audit.md` (TestNaiveLostUpdate)

---

## How It Works — Pessimistic Locking

**Implementation:**
- `internal/inventory/store.go:93-116` (`PessimisticDeduct`)
- `internal/inventory/store.go:42-51` (`GetRowLock`)

**Tests:**
- `tests/locking_test.go:40-67` (`TestPessimisticLocking`)
- `tests/locking_test.go:69-83` (`TestPessimisticLockingInsufficientStock`)

**Audit:**
- `engineering-audit/02-code-audit.md` (Finding 2)
- `engineering-audit/03-test-audit.md` (TestPessimisticLocking, TestPessimisticLockingInsufficientStock)

---

## How It Works — Optimistic Locking

**Implementation:**
- `internal/inventory/store.go:120-152` (`OptimisticDeduct`)
- `internal/inventory/service.go:28-45` (`DeductOptimisticWithRetry`)

**Tests:**
- `tests/locking_test.go:85-121` (`TestOptimisticLockingConflict`)
- `tests/locking_test.go:123-145` (`TestOptimisticLockingWithRetry`)

**Audit:**
- `engineering-audit/02-code-audit.md` (Finding 3)
- `engineering-audit/03-test-audit.md` (TestOptimisticLockingConflict, TestOptimisticLockingWithRetry)

---

## How It Works — Atomic Single-Statement

**Implementation:**
- `internal/inventory/store.go:154-173` (`AtomicDeduct`)

**Tests:**
- `tests/locking_test.go:147-171` (`TestAtomicConditionalUpdate`)

**Audit:**
- `engineering-audit/02-code-audit.md` (Finding 5)
- `engineering-audit/03-test-audit.md` (TestAtomicConditionalUpdate)

---

## Architecture

**Implementation:**
- `internal/inventory/model.go` (domain model, errors)
- `internal/inventory/store.go` (in-memory engine)
- `internal/inventory/service.go` (business layer)
- `cmd/demo/main.go` (CLI runner)

**Engineering:**
- `engineering/01-design.md` (Architecture, Components)

**Research:**
- `research/05-report.md` (Finding 6: Atomic Single-Statement)

---

## Demo

**Implementation:**
- `cmd/demo/main.go`

**Engineering:**
- `engineering/03-execution-result.md` (Demo output)

---

## Race Detector

**Engineering:**
- `engineering/03-execution-result.md` (Race Detector output)
- `engineering-audit/03-test-audit.md` (Race detector verification)

---

## Selection Criteria

**Research:**
- `research/05-report.md` (Finding 5: Selection Criteria)
- `research/03-evidence.md`

---

## Anti-Patterns & Common Mistakes

**Research:**
- `research/05-report.md` (Finding 9: Common Anti-Patterns)

---

## Isolation Levels

**Research:**
- `research/05-report.md` (Finding 7: Isolation Level Does Not Prevent Lost Update)
- `research/05-report.md` (Finding 8: Default Isolation Levels Differ)
- `research/03-evidence.md`

---

## Recovery / Rollback

**Implementation:**
- `internal/inventory/service.go:28-45` (`DeductOptimisticWithRetry`)

**Tests:**
- `tests/locking_test.go:123-145` (`TestOptimisticLockingWithRetry`)

---

## Input Validation

**Implementation:**
- `internal/inventory/store.go:66` (`NaiveDeduct` validation)
- `internal/inventory/store.go:94` (`PessimisticDeduct` validation)
- `internal/inventory/store.go:121` (`OptimisticDeduct` validation)
- `internal/inventory/store.go:156` (`AtomicDeduct` validation)

**Audit:**
- `engineering-audit/02-code-audit.md` (Finding 6)

---

## Model (Product struct, Errors)

**Implementation:**
- `internal/inventory/model.go`

**Audit:**
- `engineering-audit/05-gaps.md`

---

## Verdict & Approval

**Research:**
- `research-audit/07-verdict.md` (APPROVED)

**Engineering:**
- `engineering-audit/06-verdict.md` (APPROVED)