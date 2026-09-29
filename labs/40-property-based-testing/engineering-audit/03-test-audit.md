# Test Audit

## Test Suite Overview

| Package | Test Name | Invariant / Target | Type | Result |
|---|---|---|---|---|
| `internal/currency` | `TestExampleBasedNaiveCurrency` | Happy-path float parsing | Example-based | PASS |
| `internal/currency` | `TestPropertyNaiveCurrencyFails` | Sub-cent float roundtrip failure | Property-based | PASS |
| `internal/currency` | `TestPropertyRobustAmountRoundtrip` | Roundtrip Invariant (`Parse(Format(x)) == x`) across 1,000 runs | Property-based | PASS |
| `internal/interval` | `TestExampleBasedIntervalMerge` | Hand-picked sorted interval merge | Example-based | PASS |
| `internal/interval` | `TestPropertyNaiveMergeFails` | Oracle failure on unsorted inputs | Property-based | PASS |
| `internal/interval` | `TestPropertyRobustMergeIdempotence` | Idempotence Invariant (`f(f(x)) == f(x)`) across 1,000 runs | Property-based | PASS |
| `internal/interval` | `TestPropertyRobustMergeNonOverlapping` | Non-overlapping output invariant across 1,000 runs | Property-based | PASS |
| `internal/shrinker` | `TestFindAndShrink` | Array size reduction to minimal counterexample `[-1]` | Invariant & Shrinking | PASS |

## Coverage and Test Strength Assessment

1. **Happy Path vs Negative Paths**:
   - Both packages deliberately implement example-based tests showing passes on naive flawed implementations.
   - Property tests explicitly verify negative cases where naive implementations fail oracle properties.
2. **Invariant Verification**:
   - `Roundtrip`: Tested on negative, zero, and extreme integer values up to multi-billion amounts.
   - `Idempotence`: Tested across 1,000 generated slices with negative coordinates and varying interval lengths.
   - `Oracle Check`: Compares naive against robust implementation, confirming discrepancy on arbitrary inputs.
   - `Non-overlapping`: Directly checks structural invariants of resulting slices.
3. **Execution Results**:
   - `go test -v ./...`: All 8 test functions pass cleanly.
   - `go test -race ./...`: Zero data races detected.
   - `go run ./cmd/demo`: Demo output matches execution results reported in `engineering/03-execution-result.md`.
