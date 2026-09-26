# Test Audit

Target Lab: labs/17-architecture-decision-record

## Commands Executed

- `go test -v ./...` → PASS (all 7 top-level tests, 16 subtests). No failures.
- `go test -race ./...` → PASS, no data races.
- `go run ./cmd/demo` → PASS, exit 0. Output: parsed 3 ADRs (1 Superseded→2, 2 Accepted→supersedes 1, 3 Rejected), linter passed.
- `go build ./...` → PASS. `go vet ./...` → PASS.

## Coverage vs Required Dimensions

| Dimension | Parser | Linter |
|---|---|---|
| happy path | PASS — TestParse_Valid, TestParse_Superseded, TestParse_Supersedes | PASS — TestLinter_ValidSequence (1 Superseded↔2 Accepted + 3 Rejected) |
| failure path | PASS — missing title/status/context/decision/consequences, invalid status Draft | PASS — superseded→non-existent, supersedes→non-existent, mismatched link, self-supersession, duplicate ID, non-monotonic, cycle |
| edge cases | PASS — empty context/decision/consequences sections rejected | PASS — duplicate ID, self-reference, 3-node cycle, gap (1,3) numbering |
| transitions | PASS — Status + SupersededBy/Supersedes extraction verified | PASS — bidirectional invariant checked both directions in code; one direction explicitly tested (mismatched link) |
| recovery | NOT_APPLICABLE — pure function, no state to recover | NOT_APPLICABLE — stateless validator, no rollback |
| rollback | NOT_APPLICABLE | NOT_APPLICABLE |
| concurrency | NOT_APPLICABLE — single-goroutine parse | PASS — TestLinter_ConcurrencyStress (100 records, 50 pairs) + race detector clean |
| negative cases | PASS — all invalid inputs return error with expected substring | PASS — every invalid graph returns non-empty error list |

## Strengths

- Table-driven negative tests with substring matching on error messages.
- Linter tests cover all major invariants: uniqueness, contiguity, dangling refs, one-sided links, self-links, cycles.
- Stress test proves mutex-protected error aggregation under parallel validation.
- Parser tests enforce non-empty body per required section, not just header presence.

## Weaknesses (non-blocking, LOW)

1. Parser case-insensitivity (`(?i)` regexes) implemented but not explicitly tested.
2. Empty/nil record list behavior (should return no errors) not explicitly tested.
3. Reverse one-sided link (Supersedes declared but target not marked Superseded) exercised only via code path, no dedicated test case name; covered implicitly by symmetric check.
4. Statuses Proposed/Deprecated accepted by `IsValid()` but never appear in linter happy-path test (only Accepted/Superseded/Rejected used).

## Assessment

Test suite is strong for the claimed scope. Passing results are genuine (executed locally, outputs recorded above). No weak-pass concern: failure paths are explicitly asserted, concurrency is stress-tested and race-checked. Minor uncovered branches are LOW severity and do not undermine core behavior proof.
