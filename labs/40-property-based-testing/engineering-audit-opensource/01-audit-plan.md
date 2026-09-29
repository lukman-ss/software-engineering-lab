# Engineering Audit Plan

Target Lab: labs/40-property-based-testing
Implementation Files:
- cmd/demo/main.go
- internal/currency/currency.go
- internal/interval/interval.go
- internal/shrinker/shrinker.go
Tests:
- internal/currency/currency_test.go
- internal/interval/interval_test.go
- internal/shrinker/shrinker_test.go
Executable/Demo: cmd/demo/main.go
Approved Research Inputs: engineering/01-design.md, engineering/02-implementation-notes.md, engineering/03-execution-result.md
Main Claims To Verify:
- Property-based tests catch float precision bugs and unsorted interval bugs missed by example-based tests.
- ParseRobust(RobustAmount.Format()) == RobustAmount roundtrip holds across 1,000 iterations.
- RobustMerge is idempotent (Merge(Merge(x)) == Merge(x)) and produces non-overlapping output across 1,000 iterations.
- Shrinker reduces complex 10+ element failing arrays to minimal 1-element counterexample (`[-1]`).
Commands To Run:
- go test -v ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Random generators might lack coverage of boundary cases (e.g. 0, negatives, empty slices).
- Demo output might diverge from unit test behavior or hardcode false metrics.
