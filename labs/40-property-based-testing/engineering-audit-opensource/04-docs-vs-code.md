# Docs vs Code Audit: labs/40-property-based-testing

## Comparison Matrix

| Claim / Section in README.md | Code Implementation | Test / Demo Evidence | Alignment Status |
|---|---|---|---|
| Roundtrip Invariant (`Decode(Encode(x)) == x`) | `internal/currency` implements `NaiveCurrency` and `RobustAmount` | `TestPropertyRobustAmountRoundtrip` verifies 1,000 cases; Demo verifies 1,000 cases | MATCH |
| Float precision loss vs exact cents | `NaiveCurrency.Format` uses `%.2f`; `RobustAmount` uses `int64` cents | `TestPropertyNaiveCurrencyFails` and Demo show 1000/1000 float failures | MATCH |
| Idempotence Invariant (`f(f(x)) == f(x)`) | `internal/interval/interval.go:RobustMerge` | `TestPropertyRobustMergeIdempotence` checks 1,000 randomized slices | MATCH |
| Oracle Invariant (`NaiveMerge` vs `RobustMerge`) | `NaiveMerge` ignores sort; `RobustMerge` sorts before merging | `TestPropertyNaiveMergeFails` confirms failure; Demo records 85/100 discrepancies | MATCH |
| Counterexample Shrinking | `internal/shrinker/shrinker.go:FindAndShrink` | `TestFindAndShrink` and Demo show 10-element array shrunk to `[-1]` | MATCH |
| Directory structure in README | `cmd/demo/`, `internal/currency/`, `internal/interval/`, `internal/shrinker/`, `engineering/` | Matches filesystem exactly | MATCH |
| Execution commands documented | `go test -v ./...`, `go test -race ./...`, `go run ./cmd/demo` | All 3 commands run successfully without error | MATCH |

## Findings

- `DOC_CODE_MISMATCH`: None observed.
- `TEST_CLAIM_MISMATCH`: None observed.
- `RESEARCH_IMPLEMENTATION_MISMATCH`: None observed.
