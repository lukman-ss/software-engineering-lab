# Gap Analysis

Target Lab: labs/38-mutation-testing

## Gap Findings

### Gap 1: Incomplete Statement Deletion Operator
- Type: `MISSING_EDGE_CASE`
- Severity: LOW
- Location: `internal/engine/types.go:12`, `internal/engine/mutator.go`
- Description: `StatementDelete` exists as an enum constant in `types.go`, but AST node statement deletion logic is not implemented in `mutator.go`.
- Mitigation / Status: Non-blocking. README accurately lists only the 4 working operators, and `02-implementation-notes.md` explicitly lists this limitation.

### Gap 2: Synthetic Strong Test Oracle in Demo Runner
- Type: `IMPLEMENTATION_OVERCLAIM`
- Severity: MEDIUM
- Location: `cmd/demo/main.go:41`, `tests/engine_test.go:92`
- Description: The demo runner's strong test function checks `string(mutatedSrc) != string(src)` rather than running compiled test assertions against a mutated runtime package.
- Mitigation / Status: Non-blocking. `02-implementation-notes.md` clearly discloses that an in-memory AST diff oracle is used to avoid subprocess compilation overhead. The unit tests in `internal/service/discount_strong_test.go` provide real assertion coverage.

## Summary of Allowed Gaps Checked
- `MISSING_TEST`: None.
- `BROKEN_IMPLEMENTATION`: None.
- `DOC_CODE_MISMATCH`: None.
- `RACE_CONDITION`: None.
- `UNHANDLED_ERROR`: None.
- `MISSING_EDGE_CASE`: 1 (StatementDelete omitted, LOW).
- `IMPLEMENTATION_OVERCLAIM`: 1 (Synthetic strong test oracle documented in notes, MEDIUM).
- `RESEARCH_MISMATCH`: None.
- `FAKE_DEMO`: None (demo output is 100% reproducible live).
- `FAKE_BENCHMARK`: None.
- `UNVERIFIED_RESULT`: None.
