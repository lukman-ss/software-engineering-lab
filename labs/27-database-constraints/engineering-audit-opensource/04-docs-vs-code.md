# Docs vs Code Audit — labs/27-database-constraints

## README.md vs Implementation
- Claims: Implements NOT NULL (23502), CHECK (23514), UNIQUE (23505), FOREIGN KEY (23503), PARTIAL UNIQUE INDEX.
- Code: ✅ exactly those via engine; demo shows each.
- Test commands: `go test -v ./...`, `go test -race ./...` match repo.
- No mismatch.

## engineering/01-design.md vs Code
- Design lists constraints; engine enforces NOT NULL, CHECK, UNIQUE, FOREIGN KEY, PARTIAL UNIQUE.
- Concurrency claim: read-then-write race vs engine lock; UnsafeStore unsafe path proven, SafeStore safe via mutex.
- Partial index semantics: engine replicates WHERE deleted_at IS NULL via separate activeEmails map.
- Error mapping: design states mapping 23502/23503/23505/23514 to domain errors; engine returns ConstraintError, store maps.
- Architecture: engine/store split matches internal/engine/internal/store.
- Test strategy: enumerated tests present; concurrency test pair present.
- Success criteria: all satisfied.
- Implementation decisions: in-memory pure-Go simulation matches code.
- No mismatches found.

## engineering/02-implementation-notes.md vs Code
- File list matches.
- Core design decisions: engine mutex/index layer for invariants; matches engine.go lock + maps.
- SQLSTATE standardization: uses exact PostgreSQL codes 23502/23503/23505/23514; dberr defines them.
- Partial Unique Index: engine uses activeEmails map for DeletedAt==nil path.
- Implementation-specific choices: pure Go stdlib (sync.RWMutex, atomic.Int64, maps) — matches.
- Known Limitations: "in-memory simulator" — code shows no external DB; accurate.
- Trade-offs: coarse table lock vs fine B-tree — matches single RWMutex.
- Demonstrated vs Not Demonstrated: matches output.
- No mismatches.

## engineering/03-execution-result.md vs Actual Demo/Tests
- Build: go build ./... → success (inferred from test).
- Tests: output exactly matches captured run (order and PASS lines).
- Race Detector: matches captured output (no races).
- Demo: verbatim match to captured go run ./cmd/demo output.
- Final status READY_FOR_ENGINEERING_AUDIT accurate.

## research/ (not audited per pipeline override)  
- Skipped per instruction.

Result: No DOC_CODE_MISMATCH, no TEST_CLAIM_MISMATCH, no RESEARCH_IMPLEMENTATION_MISMATCH (skipped). Documentation is accurate.