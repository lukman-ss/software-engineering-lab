# Documentation vs Code Audit

Target Lab: labs/27-database-constraints

## 1. README.md vs Code
- README lists all 5 implemented constraints: NOT NULL (`23502`), CHECK (`23514`), UNIQUE (`23505`), FOREIGN KEY (`23503`), and PARTIAL UNIQUE INDEX.
- Matches code implementations in `internal/engine/engine.go` and `internal/store/store.go`.
- Test commands (`go test -v ./...`, `go test -race ./...`) and demo command (`go run ./cmd/demo`) match repository reality.
- Assessment: PASS

## 2. Engineering Design (`01-design.md`) vs Code
- Planned Architecture: `internal/engine`, `internal/store`, `internal/dberr` (noted as `internal/errors` in early design draft).
- Success Criteria: Planned `TestConcurrentRegistration_Unsafe_SuffersRaceCondition` alongside `TestConcurrentRegistration_Safe_EnforcesUniqueness`. The safe test is fully implemented, but the unsafe race test was omitted from `store_test.go`.
- Assessment: WARNING (Minor test omission from early plan)

## 3. Implementation Notes (`02-implementation-notes.md`) vs Code
- Matches actual files and package names (`internal/dberr`, `internal/model`, `internal/engine`, `internal/store`, `cmd/demo`).
- Accurately details zero external server dependency trade-off and coarse-grained locking approach.
- Assessment: PASS

## 4. Execution Result (`03-execution-result.md`) vs Actual Runtime
- All demo outputs, test output counts, and SQLSTATE codes recorded in `03-execution-result.md` match live terminal execution character-for-character.
- Concurrency demo output (50 goroutines, 1 success, 49 rejected) is genuine and verifiable.
- Assessment: PASS
