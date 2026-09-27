# Engineering Audit Plan

Target Lab: labs/27-database-constraints
Implementation Files:
- `internal/model/model.go`
- `internal/dberr/errors.go`
- `internal/engine/engine.go`
- `internal/store/store.go`
Tests:
- `internal/store/store_test.go`
Executable/Demo:
- `cmd/demo/main.go`
Approved Research Inputs:
- `research/01-plan.md`
- `research/02-sources.md`
- `research/03-evidence.md`
- `research/04-contradictions.md`
- `research/05-report.md`
- `research-audit/07-verdict.md`
Main Claims To Verify:
1. Storage engine constraint checks enforce NOT NULL (`23502`), CHECK (`23514`), UNIQUE (`23505`), and FOREIGN KEY (`23503`) integrity.
2. Partial unique indexing enforces conditional uniqueness (`WHERE deleted_at IS NULL`) allowing re-registration after soft deletion.
3. Storage engine-level UNIQUE constraints prevent read-then-write concurrency race conditions under simultaneous worker requests (50 goroutines), whereas application-level checks fail.
4. Error taxonomy matches standard SQLSTATE codes and maps to domain error representations.
Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`
Primary Risks:
- Race conditions or false-positive passes under race detector.
- Discrepancies between execution result logs and live demo output.
- Unhandled edge cases in partial index or foreign key verification.
