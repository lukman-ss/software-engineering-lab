# Gap Analysis

## Gap 1

Type: MISSING_EDGE_CASE
Severity: MEDIUM
Location: `internal/adr/linter.go:49-70`
Problem: The linter verifies reciprocal references between `SupersededBy` and `Supersedes` but does not enforce DAG directionality (i.e., that `SupersededBy > ID` and `Supersedes < ID`). A self-referential ADR or backward supersession cycle can pass validation.
Impact: In an adversarial or mistaken repository commit, a decision could technically supersede itself.

---

## Gap 2

Type: MISSING_TEST
Severity: LOW
Location: `tests/linter_test.go`
Problem: Missing automated unit tests for duplicate ADR IDs and `Status: Superseded` records missing the `by <ID>` target clause, even though the linter implementation has error-reporting branches for them.
Impact: Code coverage omits two failure branches.

---

## Gap 3

Type: UNHANDLED_ERROR
Severity: LOW
Location: `internal/adr/parser.go:38`
Problem: Deprecated function `strings.Title` is used for status normalization instead of `cases.Title(language.English)` or a simple map/switch lookup.
Impact: Compiler warning / deprecation notice in modern Go toolchains.
