# Gap Analysis

Target Lab: `labs/38-mutation-testing`

## Identified Gaps

### GAP-1: Statement Deletion Operator Not Implemented in AST Walker
- Type: `IMPLEMENTATION_OVERCLAIM` (Minor)
- Severity: LOW
- Location: `internal/engine/types.go:12`
- Description: `StatementDelete` enum member is declared in `types.go` and listed in `01-design.md`, but `mutator.go` does not generate statement deletion mutations.
- Mitigation / Status: `engineering/02-implementation-notes.md` explicitly documents this trade-off (no void statements in `discount.go`). `README.md` only advertises the 4 active operators.

### GAP-2: In-Memory Textual Test Oracle Simulation
- Type: `IMPLEMENTATION_OVERCLAIM` (Minor)
- Severity: LOW
- Location: `cmd/demo/main.go:40-43`, `tests/engine_test.go:90-97`
- Description: Mutation engine evaluates `TestFunc` closures via string/AST comparison rather than recompiling and running test binaries via `go test` subprocesses.
- Mitigation / Status: Explicitly documented as an intentional design decision in `engineering/02-implementation-notes.md` (avoids slow subprocess spawns and external dependencies).

## Gaps Not Found (Verified Clean)
- `MISSING_TEST`: None. Comprehensive boundary tests and engine unit tests are present.
- `BROKEN_IMPLEMENTATION`: None. Code parses, mutates, and executes cleanly.
- `DOC_CODE_MISMATCH`: None. README matches codebase.
- `RACE_CONDITION`: None. Race detector passes cleanly with 0 data races.
- `UNHANDLED_ERROR`: None. Parser errors and render errors are handled in runner.
- `FAKE_DEMO`: None. Real terminal output reproduces execution result identically.
- `FAKE_BENCHMARK`: None. No fabricated benchmarks present.
- `UNVERIFIED_RESULT`: None.
