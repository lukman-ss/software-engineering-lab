# Engineering Audit Plan

Target Lab: `labs/40-property-based-testing`
Implementation Files:
- `internal/currency/currency.go`
- `internal/interval/interval.go`
- `internal/shrinker/shrinker.go`

Tests:
- `internal/currency/currency_test.go`
- `internal/interval/interval_test.go`
- `internal/shrinker/shrinker_test.go`

Executable/Demo:
- `cmd/demo/main.go`

Approved Research Inputs:
- `research/01-plan.md`
- `research/02-sources.md`
- `research/03-evidence.md`
- `research/05-report.md`
- `research-audit/07-verdict.md`

Main Claims To Verify:
1. Roundtrip invariant fails on float-based currency and passes on integer cent currency.
2. Idempotence and non-overlapping invariants hold on `RobustMerge` intervals, while `NaiveMerge` fails oracle testing on unsorted inputs.
3. Counterexample shrinking reduces complex slice inputs down to minimal failing counterexamples.
4. Demo executes without mock data and faithfully mirrors the property-based testing behavior.
5. All unit and property tests pass under `go test` and `go test -race`.

Commands To Run:
```bash
go test -v -count=1 ./...
go test -race -count=1 ./...
go run ./cmd/demo
```

Primary Risks:
- Random generator skew or insufficient iteration coverage in `testing/quick`.
- Potential race conditions in random generation or concurrent test execution.
- Discrepancy between README claims and actual package behavior.
