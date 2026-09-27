# Engineering Audit Plan

## Target Lab
labs/27-database-constraints (`/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/27-database-constraints`)

## Implementation Files
- `go.mod` — module `github.com/lukman/software-engineering-lab/labs/27-database-constraints` (Go 1.22)
- `internal/dberr/errors.go` — SQLSTATE `23xxx` taxonomy, `ConstraintError`, `IsConstraintViolation`, `MapToDomainError`
- `internal/model/model.go` — `User`, `Order` domain entities
- `internal/engine/engine.go` — in-memory relational engine enforcing NOT NULL, CHECK, UNIQUE, FOREIGN KEY, PARTIAL UNIQUE INDEX with `sync.RWMutex`+`atomic.Int64`
- `internal/store/store.go` — `UnsafeStore` (read-then-write, race-prone) and `SafeStore` (delegates to engine constraints)
- `cmd/demo/main.go` — runnable CLI demonstration

## Tests
- `internal/store/store_test.go` (Go `testing`, package `store_test`)
  - `TestNotNullConstraints`, `TestCheckConstraints`, `TestUniqueConstraint`, `TestForeignKeyConstraint`, `TestPartialUniqueIndex`, `TestConcurrentRegistration_Safe_EnforcesUniqueness`, `TestConcurrentRegistration_Unsafe_SuffersRaceCondition`, `TestErrorClassification`

## Executable/Demo
- `go run ./cmd/demo`

## Approved Research Inputs (referenced)
- `research/05-report.md`, `research/06-open-questions.md` (not the audit target this stage)
- `research-audit/` source audit (not the audit target this stage)

## Main Claims To Verify
1. NOT NULL (`23502`) rejects missing `email`, `username`, `orders.user_id`.
2. CHECK (`23514`) rejects `age < 18` and invalid `status`; rejects `total_cents <= 0`.
3. UNIQUE (`23505`) rejects duplicate `email` in `SafeStore`.
4. FOREIGN KEY (`23503`) rejects orders referencing non-existent `user_id`.
5. PARTIAL UNIQUE INDEX allows re-registration after soft delete + blocks duplicate active.
6. Concurrency: `SafeStore` yields exactly 1 success + N `23505`; `UnsafeStore` yields >1 (race condition).
7. `dberr.MapToDomainError` maps each SQLSTATE to a domain message; `IsConstraintViolation` classifies correctly.
8. README matches code; demo output real.

## Commands To Run
```bash
go build ./...
go vet ./...
go test -v -count=1 ./...
go test -race -count=1 ./...
go run ./cmd/demo
```

## Primary Risks
- Concurrency test non-determinism (unsafe path relies on scheduler exposing duplicates; safe path serialized by mutex). Race detector must be clean.
- Missing negative coverage for edge cases (age boundary = 18 accepted; status enum boundary; soft-delete then hard re-insert collision semantics).
- Doc vs code drift (README mentions constraints; does it match code surface?).
