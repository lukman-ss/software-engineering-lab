# Docs vs Code — labs/27-database-constraints

Sources compared: README.md, engineering/01-design.md, engineering/02-implementation-notes.md, engineering/03-execution-result.md, code, tests, live demo output.

## D1 — README vs code: MATCH
README lists NOT NULL 23502, CHECK 23514, UNIQUE 23505 + race prevention, FK 23503, partial unique with soft-delete re-registration. All five reproduced live via `go run ./cmd/demo` and `go test -v`. Test/demo commands in README (`go test -v ./...`, `go test -race ./...`, `go run ./cmd/demo`) all execute with exit 0.

## D2 — Execution-result vs live run: MATCH
`engineering/03-execution-result.md` demo transcript (50-worker stress, 1 success / 49 rejected, taxonomy line) byte-matches live demo output observed in this audit. Test transcript lists 7 tests; live run has 8 (extra `TestConcurrentRegistration_Unsafe_SuffersRaceCondition`). Result still PASS; doc is stale by one test, not fabricated. Severity LOW.

## D3 — DOC_CODE_MISMATCH (LOW): `tests/...` path does not exist
Design Execution Plan step 4 cites `internal/store/...` and `tests/...`. No `tests/` directory exists; all tests live in `internal/store/store_test.go`. Coverage complete, location differs.

## D4 — DOC_CODE_MISMATCH (LOW): test names vs Test Strategy
Design lists `TestNotNullConstraint`, `TestCheckConstraint`, `TestUniqueConstraint`, `TestForeignKeyConstraint`, `TestErrorMapping`. Actual: `TestNotNullConstraints`, `TestCheckConstraints`, `TestUniqueConstraint`, `TestForeignKeyConstraint`, `TestErrorClassification`, plus two concurrency tests. Logic matches intent; names differ.

## D5 — DOC_CODE_MISMATCH (LOW): CHECK NULL pass-through
Design Expected Behavior says CHECK "passing NULL if nullable". `User.Age` is plain `int`; Age=0 always fails. No nullable check column exists. Untested claim fragment.

## D6 — No fake demo / benchmark
Demo output computed live from engine state (IDs, counts, error strings match code paths). No hardcoded result tables, no benchmark claims. `FAKE_DEMO` / `FAKE_BENCHMARK`: not found.

## D7 — Implementation scope honestly disclosed
`02-implementation-notes.md` states in-memory simulator, no live Postgres, no EXCLUDE constraints, coarse table lock. Matches `go.mod` (stdlib only) and code. No overclaim found.
