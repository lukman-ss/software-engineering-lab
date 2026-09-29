# Engineering Audit Plan

Target Lab: labs/40-property-based-testing
Implementation Files:
- cmd/demo/main.go
- internal/currency/*.go
- internal/interval/*.go
- internal/shrinker/*.go
Tests:
- internal/**/*_test.go
Executable/Demo: cmd/demo
Approved Research Inputs: research/*
Main Claims To Verify:
- Property-based tests reveal bugs missed by example-based tests.
- NaiveCurrency loses precision; RobustAmount roundtrip holds.
- NaiveMerge fails on unsorted input; RobustMerge idempotent and non-overlapping.
- Shrinker finds minimal counterexample for negative values.
Commands To Run:
- go test ./... 
- go test -race ./... 
- go run ./cmd/demo
Primary Risks:
- Tests may be flaky or insufficient edge coverage.
- Demo output may be fabricated.
