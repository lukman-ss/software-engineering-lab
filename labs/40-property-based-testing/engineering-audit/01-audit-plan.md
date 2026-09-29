# Engineering Audit Plan

Target Lab: labs/40-property-based-testing
Implementation Files:
- `internal/currency/currency.go`
- `internal/interval/interval.go`
- `internal/shrinker/shrinker.go`
- `cmd/demo/main.go`
- `go.mod`
Tests:
- `internal/currency/currency_test.go`
- `internal/interval/interval_test.go`
- `internal/shrinker/shrinker_test.go`
Executable/Demo:
- `cmd/demo/main.go`
Approved Research Inputs:
- `research-audit/07-verdict.md` (APPROVED)
- `research/05-report.md`
- `engineering/01-design.md`
- `engineering/02-implementation-notes.md`
Main Claims To Verify:
1. Pure Go standard library PBT implementation using `testing/quick` with zero third-party dependencies.
2. Canonical Invariants verified:
   - Roundtrip: `ParseRobust(Format(amount)) == amount` passes on `RobustAmount` (integer cents) across 1,000 iterations and fails on naive float currency.
   - Idempotence: `Merge(Merge(x)) == Merge(x)` holds for robust interval merging across 1,000 iterations.
   - Equivalence / Oracle: `NaiveMerge` fails compared to `RobustMerge` on arbitrary unsorted inputs.
   - Non-overlapping invariant: `RobustMerge` outputs strictly non-overlapping intervals.
3. Counterexample Shrinking:
   - Automated search discovers a multi-element counterexample violating an invariant.
   - Shrinking mechanics (binary list sectioning + element reduction toward zero) reduce failing array to minimal single-element counterexample `[-1]`.
4. Tests compile cleanly, pass with race detector enabled (`go test -race ./...`), and demo runs deterministically.
Commands To Run:
- `go build ./...`
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`
Primary Risks:
- Randomness non-determinism causing flaky test passes/failures in `testing/quick` or shrinker.
- Unhandled arithmetic overflow in int64 cent conversions.
- Edge case handling with empty slices, single element slices, or negative coordinate intervals.
- Discrepancy between reported demo logs in `engineering/03-execution-result.md` and actual execution output.
