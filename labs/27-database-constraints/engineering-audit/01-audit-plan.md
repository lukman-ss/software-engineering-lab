# Engineering Audit Plan

Target Lab: labs/27-database-constraints
Implementation Files:
- internal/dberr/errors.go
- internal/model/model.go
- internal/engine/engine.go
- internal/store/store.go
- cmd/demo/main.go

Tests:
- internal/store/store_test.go

Executable/Demo:
- cmd/demo/main.go

Approved Research Inputs:
- research/01-plan.md
- research/02-sources.md
- research/03-evidence.md
- research/04-contradictions.md
- research/05-report.md
- research/06-open-questions.md
- research-audit/07-verdict.md (APPROVED WITH WARNINGS)

Main Claims To Verify:
1. Implementation of NOT NULL (`23502`), CHECK (`23514`), UNIQUE (`23505`), FOREIGN KEY (`23503`), and PARTIAL UNIQUE INDEX (`WHERE deleted_at IS NULL`) constraints.
2. In-memory relational storage engine concurrency safety and race condition elimination under concurrent registration.
3. Domain error mapping of standard SQLSTATE Class 23 error codes.
4. Accuracy of execution claims in `engineering/01-design.md`, `engineering/02-implementation-notes.md`, `engineering/03-execution-result.md`, and `README.md`.
5. Code completeness and alignment between design claims and code implementation (e.g. UnsafeStore concurrency test claims).

Commands To Run:
- `cd labs/27-database-constraints && go test -v ./...`
- `cd labs/27-database-constraints && go test -race ./...`
- `cd labs/27-database-constraints && go run ./cmd/demo`

Primary Risks:
- Discrepancy between design doc (`engineering/01-design.md`) claims vs actual implementation (e.g., claimed `TestConcurrentRegistration_Unsafe_SuffersRaceCondition` test in design doc vs actual code).
- In-memory relational engine thread-safety locking scope and potential race conditions under full concurrency.
- Incomplete coverage of failure edge cases or unhandled SQLSTATE codes.
