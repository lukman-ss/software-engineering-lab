# Docs vs Code Audit

Sources compared: `README.md`, `engineering/01-design.md`, `engineering/02-implementation-notes.md`, `engineering/03-execution-result.md`.

## Alignment Matrix

### 1. README vs Code: CONSTRAINT SET COVERAGE

README lists 5 implemented constraints.

| README Claim | Code Evidence | Match |
|---|---|---|
| NOT NULL (23502), missing email/username/user_id | engine.go:51-56, 128-130 | MATCH |
| CHECK (23514): age >= 18, status IN (...), total_cents > 0 | engine.go:60-61, 63-68, 133 | MATCH |
| UNIQUE (23505): single occurrence, prevents read-then-write race | engine.go:80-82, 92-93; store.go:61 | MATCH |
| FOREIGN KEY (23503): prevents orphan rows | engine.go:138-140 | MATCH |
| PARTIAL UNIQUE INDEX: WHERE deleted_at IS NULL, soft-delete reuse | engine.go:73-77, 94-95 | MATCH |

Assessment: README coverage matches code. No overstated features.

### 2. README vs Code: HOW TO RUN

README documents `go test -v ./...` and `go test -race ./...` + `go run ./cmd/demo`.

- These commands all execute successfully (verified by auditor).  
- README does not mention `go vet`. Minor omission, not a mismatch.

Assessment: PASS.

### 3. Design doc (01-design.md) vs code: ARCHITECTURE PACKAGES

Design doc claims architecture:

```
internal/db   — in-memory engine + constraint definitions
internal/domain — entities
internal/service — safe vs unsafe workflows
internal/errors — SQLSTATE taxonomy + mappers
```

Actual code:

```
internal/engine  — the engine (what design called internal/db)
internal/model   — entities (what design called internal/domain)
internal/store   — workflows (what design called internal/service)
internal/dberr   — error taxonomy (what design called internal/errors)
```

Assessment: DOC_CODE_MISMATCH. Severity MEDIUM. The *functionality* described matches (`engine` = storage engine with locks + SQLSTATE; `model` = entities; `store` = safe/unsafe; `dberr` = error taxonomy). But the documented package paths do not exist in code — a reader copying the design into code would fail. The implementation notes (02) correct the record by listing the real packages (`internal/dberr`, `internal/model`, `internal/engine`, `internal/store`). So the mismatch is internal between the two engineering docs, not between docs and final code.

### 4. Design doc SUCCESS CRITERIA vs reality

| Design Criterion | Reality | Result |
|---|---|---|
| Full test coverage each constraint | 5 dedicated tests | PASS |
| Concurrency test unsafe loses, safe wins exactly 1 / N-1 23505 | 2 tests, live demo 1/49 | PASS |
| Error classification categorizes 23502/23503/23505/23514 | `IsConstraintViolation` + `MapToDomainError` | PASS (mapper not asserted in tests, see gap) |
| Tests pass with go test ./... and -race | YES | PASS |
| Standalone demo in cmd/demo runs with clear output | YES | PASS |

Assessment: All design success criteria met.

### 5. Execution-result (03-execution-result.md) vs actual

The engineering execution-result lists the test run as including these tests and PASS:
`TestNotNullConstraints`, `TestCheckConstraints`, `TestUniqueConstraint`, `TestForeignKeyConstraint`, `TestPartialUniqueIndex`, `TestConcurrentRegistration_Safe_EnforcesUniqueness`, `TestErrorClassification`.

That is **7 tests**. The actual test file contains **8 tests**: it is missing `TestConcurrentRegistration_Unsafe_SuffersRaceCondition`.

Assessment: DOC_CODE_MISMATCH / TEST_CLAIM_MISMATCH. Severity MEDIUM. The recorded "execution result" omits the unsafe-race test entirely — yet that test is a core deliverable (Success Criterion #2 — "unsafe store suffers race condition duplicates"). The auditor re-ran the suite and this test PASSED. So the result file is under-reported rather than fabricated, but it does not faithfully record the full suite that the design specified. A technical writer cross-referencing would find a test the result doc doesn't mention.

### 6. Implementation notes (02-implementation-notes.md) vs code

Notes accurately enumerate files (`go.mod`, `internal/dberr/errors.go`, `internal/model/model.go`, `internal/engine/engine.go`, `internal/store/store.go`, `internal/store/store_test.go`, `cmd/demo/main.go`). All present and correct. Known limitations section honestly states in-memory (no live PG), EXCLUDE omitted.

Assessment: PASS.

### 7. Demo output vs code vs execution-result

`cmd/demo/main.go` produces:
```
[1] NOT NULL 23502 … missing email
[2] CHECK 23514 … age=15, invalid status
[3] FK 23503 … UserID=9999
[4] Partial unique: create alice ID=1; dup rejected; soft-delete; re-create ID=2
[5] 50 goroutines → 1 success, 49 rejected (23505), integrity intact
SQLSTATE taxonomy: code=23505 isUniqueViolation=true
```

The auditor ran `go run ./cmd/demo` (see 06-verdict for literal captured output) — it **exactly matches** the recorded `03-execution-result.md` for steps [1]-[5] and the taxonomy line. The only deviation is 03-result's test list omitting one test (covered above).

Assessment: Demo output genuine and matches both code and the result doc. PASS.

### 8. README "Implemented Constraints" count vs design "Expected Behavior"

Notably the design doc's "Expected Behavior" mentions PRIMARY KEY ("UNIQUE / PRIMARY KEY constraints reject duplicate keys 23505") but the README's 5-item list does not separately enumerate PRIMARY KEY. Code: primary-key semantics are modeled by `userSeq.Add(1)` (auto PK assignment); there is no explicit `PK` constraint violation test (e.g., inserting a duplicate fixed ID).

Assessment: DOC_CODE_MISMATCH (minor). Severity LOW. PK uniqueness is only implicit via sequence generation; no test asserts PK collision rejection. Not a functional defect, but a slight over-claim in the design doc. The README is fine (it lists 5 concrete implemented constraints and none claim PK).

### 9. Research claims (research/05-report.md headline) vs code

Research headline claim: "Database-level constraints as the foundation for data integrity, race condition prevention, and schema invariants."

Code realizes this: SafeStore delegates to engine constraints; UnsafeStore exhibits the race the research warns about. Demo proves the contrast.

Assessment: RESEARCH_IMPLEMENTATION_ALIGNMENT PASS (implementation-level).

## Summary of doc drift

- Real mismatch: `01-design.md` package paths (`internal/db/domain/service/errors`) ≠ code (`internal/engine/model/store/dberr`). (Severity MEDIUM, but self-corrected in 02-implementation-notes.)
- `03-execution-result.md` test list omits the unsafe concurrency test (Severity MEDIUM).
- Design doc over-stated PRIMARY KEY as a demonstrated constraint; code/README do not. (LOW)
- All behavioral claims (constraint semantics, concurrency outcome, demo output) match code and run faithfully.
