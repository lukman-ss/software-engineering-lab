# Test Audit: labs/40-property-based-testing

## Test Suite Overview

| Test Name | Package | Target Under Test | Invariant / Property Verified |
|---|---|---|---|
| `TestExampleBasedNaiveCurrency` | `currency` | `NaiveCurrency` | Demonstrates passing hand-picked happy cases |
| `TestPropertyNaiveCurrencyFails` | `currency` | `NaiveCurrency` | Proves property failure on arbitrary floats with `testing/quick` |
| `TestPropertyRobustAmountRoundtrip` | `currency` | `RobustAmount` | Proves roundtrip `ParseRobust(amount.Format()) == amount` over 1,000 cases |
| `TestExampleBasedIntervalMerge` | `interval` | `NaiveMerge` | Demonstrates passing hand-picked sorted input |
| `TestPropertyNaiveMergeFails` | `interval` | `NaiveMerge` vs `RobustMerge` | Proves oracle discrepancy between naive and robust implementations |
| `TestPropertyRobustMergeIdempotence` | `interval` | `RobustMerge` | Proves `Merge(Merge(x)) == Merge(x)` over 1,000 cases |
| `TestPropertyRobustMergeNonOverlapping` | `interval` | `RobustMerge` | Proves non-overlapping output invariant over 1,000 cases |
| `TestFindAndShrink` | `shrinker` | `FindAndShrink` | Proves reduction of multi-element array to minimal `[-1]` |

## Execution Results

### 1. `go test -v ./...`
```text
=== RUN   TestExampleBasedNaiveCurrency
--- PASS: TestExampleBasedNaiveCurrency (0.00s)
=== RUN   TestPropertyNaiveCurrencyFails
--- PASS: TestPropertyNaiveCurrencyFails (0.00s)
=== RUN   TestPropertyRobustAmountRoundtrip
--- PASS: TestPropertyRobustAmountRoundtrip (0.00s)
PASS
ok  	labs/40-property-based-testing/internal/currency	0.097s
=== RUN   TestExampleBasedIntervalMerge
--- PASS: TestExampleBasedIntervalMerge (0.00s)
=== RUN   TestPropertyNaiveMergeFails
--- PASS: TestPropertyNaiveMergeFails (0.00s)
=== RUN   TestPropertyRobustMergeIdempotence
--- PASS: TestPropertyRobustMergeIdempotence (0.00s)
=== RUN   TestPropertyRobustMergeNonOverlapping
--- PASS: TestPropertyRobustMergeNonOverlapping (0.00s)
PASS
ok  	labs/40-property-based-testing/internal/interval	0.097s
=== RUN   TestFindAndShrink
--- PASS: TestFindAndShrink (0.00s)
PASS
ok  	labs/40-property-based-testing/internal/shrinker	0.094s
```

### 2. `go test -race ./...`
```text
ok  	labs/40-property-based-testing/internal/currency	(cached)
ok  	labs/40-property-based-testing/internal/interval	(cached)
ok  	labs/40-property-based-testing/internal/shrinker	(cached)
```
Status: PASS (0 race conditions detected).

### 3. Demo Execution (`go run ./cmd/demo`)
Demo executes real computations for 1,000 iterations across currency and interval invariants, reproduces oracle discrepancies (85/100 naive failures), and executes 21 steps of shrinking to isolate `[-1]`. Output matches recorded execution notes.

## Test Depth & Coverage Assessment

- **Generator Edge Cases**: `RobustAmount.Generate` explicitly biases towards `0`, negative ranges, boundary cent offsets, and large bounds up to trillions. `IntervalSlice.Generate` generates negative interval bounds, varying slice lengths (2 to 9), and varying offsets.
- **Negative Testing**: Both `TestPropertyNaiveCurrencyFails` and `TestPropertyNaiveMergeFails` actively assert that naive implementations fail property verification.
- **Conclusion**: Test suite is robust, reproducible, and provides concrete evidence for all stated claims.
