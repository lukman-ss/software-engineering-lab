# Source Map

## Content Brief
Sources: `README.md`, `research-audit/07-verdict.md`, `engineering-audit/06-verdict.md`

## Problem & Failure Scenario
Sources:
- Research: `research/05-report.md` (Finding 1: Application-Only Validation Is Vulnerable to Race Conditions)
- Engineering: `engineering/01-design.md` (Failure Scenario section)
- Implementation: `internal/store/store.go:24-47` (UnsafeStore.RegisterUser)
- Test: `internal/store/store_test.go:210-239` (TestConcurrentRegistration_Unsafe_SuffersRaceCondition)

## Why This Matters & Mental Model
Sources:
- Research: `research/05-report.md` (Executive Summary, Finding 1-9)
- Research: `research/02-sources.md` (PostgreSQL 18, SQLite, Stripe)
- Engineering: `engineering/01-design.md` (Concept To Prove)

## Core Concept: NOT NULL, CHECK, UNIQUE, FOREIGN KEY, Partial Index
Sources:
- Research: `research/05-report.md` (Findings 1-11)
- Implementation: `internal/engine/engine.go:45-148`
- Test: `internal/store/store_test.go:16-159`

## Implementation
Sources:
- `internal/model/model.go` (domain models)
- `internal/engine/engine.go` (constraint engine)
- `internal/dberr/errors.go` (error taxonomy and mapping)
- `internal/store/store.go` (UnsafeStore and SafeStore)
- `cmd/demo/main.go` (interactive demo)

## Code Walkthrough (Snippet 5: InsertUser)
Source: `internal/engine/engine.go:45-99`

## What the Tests Prove
Sources:
- `internal/store/store_test.go` (8 tests: NOT NULL, CHECK, UNIQUE, FK, Partial Index, Concurrency Safe, Concurrency Unsafe, Error Classification)
- `engineering/03-execution-result.md` (test output)
- `engineering-audit/03-test-audit.md` (test audit)

## Race Condition Demo (50 goroutines)
Sources:
- Implementation: `cmd/demo/main.go:59-94`
- Output: `engineering/03-execution-result.md` (Section: Demo)
- Test: `internal/store/store_test.go:162-208` (TestConcurrentRegistration_Safe_EnforcesUniqueness)

## Partial Unique Index Demo
Sources:
- Implementation: `internal/engine/engine.go:71-77, 101-120`
- Test: `internal/store/store_test.go:114-160` (TestPartialUniqueIndex)
- Demo: `cmd/demo/main.go:43-56`

## Error Mapping
Sources:
- `internal/dberr/errors.go:83-100` (MapToDomainError)
- Test: `internal/store/store_test.go:241-261` (TestErrorClassification)

## Soft-Delete Flow
Sources:
- Implementation: `internal/engine/engine.go:101-120` (SoftDeleteUser)
- Test: `internal/store/store_test.go:133-157`
- Demo: `cmd/demo/main.go:51-56`

## Recovery / Rollback (NOT VALID + VALIDATE)
Sources:
- Research: `research/05-report.md` (Finding 10)
- Research: `research/02-sources.md` (PostgreSQL ALTER TABLE docs)
- NOT implemented in lab; documented as research-only context

## Production Considerations
Sources:
- Research: `research/05-report.md` (Limitations, Finding 5, Finding 4)
- Research: `research/04-contradictions.md` (MySQL unverified, partial index limitations)
- Research: `research/06-open-questions.md` (open questions)
- Research audit: `research-audit/06-gaps.md` (Gap 1, Gap 2, Gap 3)

## Common Mistakes
Sources:
- Research: `research/05-report.md` (Findings 5, 8)
- Engineering: `engineering/01-design.md` (Concept To Prove)

## Case Study
Sources:
- Implementation: `cmd/demo/main.go` (full demo execution)
- Output: `engineering/03-execution-result.md`
- Test: `internal/store/store_test.go`

## Checklists
Sources:
- `README.md` (5 implemented constraint categories)
- `engineering/01-design.md` (Success Criteria)
- `engineering-audit/06-verdict.md` (Quality Gates)

## Diagrams
Sources:
- Race condition: `internal/store/store.go:24-47`, `internal/store/store_test.go:210-239`
- Atomic enforcement: `internal/engine/engine.go:45-99`
- Constraint evaluation order: `internal/engine/engine.go:45-148`
- Partial unique index: `internal/engine/engine.go:71-96, 101-120`
- Architecture: `internal/store/store.go`, `internal/engine/engine.go`, `internal/dberr/errors.go`
- Error mapping: `internal/dberr/errors.go:83-100`

## Source Map Section
Sources:
- This document
- Research: `research/02-sources.md` (all 10 sources)
- Research audit: `research-audit/07-verdict.md` (APPROVED)
- Engineering audit: `engineering-audit/06-verdict.md` (APPROVED)