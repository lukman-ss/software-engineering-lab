# Execution Result

## Build
Command:
```
go build ./...
```
Result:
```
(no output — build succeeded)
```

## Tests
Command:
```
go test -v ./...
```
Result:
```
?       labs/40-property-based-testing/cmd/demo        [no test files]
=== RUN   TestExampleBasedNaiveCurrency
--- PASS: TestExampleBasedNaiveCurrency (0.00s)
=== RUN   TestPropertyNaiveCurrencyFails
--- PASS: TestPropertyNaiveCurrencyFails (0.00s)
=== RUN   TestPropertyRobustAmountRoundtrip
--- PASS: TestPropertyRobustAmountRoundtrip (0.00s)
PASS
ok      labs/40-property-based-testing/internal/currency        0.501s
=== RUN   TestExampleBasedIntervalMerge
--- PASS: TestExampleBasedIntervalMerge (0.00s)
=== RUN   TestPropertyNaiveMergeFails
--- PASS: TestPropertyNaiveMergeFails (0.00s)
=== RUN   TestPropertyRobustMergeIdempotence
--- PASS: TestPropertyRobustMergeIdempotence (0.00s)
=== RUN   TestPropertyRobustMergeNonOverlapping
--- PASS: TestPropertyRobustMergeNonOverlapping (0.00s)
PASS
ok      labs/40-property-based-testing/internal/interval        0.373s
=== RUN   TestFindAndShrink
--- PASS: TestFindAndShrink (0.00s)
PASS
ok      labs/40-property-based-testing/internal/shrinker        0.502s
```

## Race Detector
Command:
```
go test -race ./...
```
Result:
```
?       labs/40-property-based-testing/cmd/demo        [no test files]
ok      labs/40-property-based-testing/internal/currency        1.389s
ok      labs/40-property-based-testing/internal/interval        1.405s
ok      labs/40-property-based-testing/internal/shrinker        1.386s
```
No race conditions detected.

## Demo
Command:
```
go run ./cmd/demo
```
Result:
```
========================================================================
  DEMO 1: Currency Formatting – Roundtrip Invariant
========================================================================

-- Example-Based Tests (Naive Float Currency) --
  EXAMPLE PASS: $1.00 -> $1.00 -> $1.00
  EXAMPLE PASS: $5.25 -> $5.25 -> $5.25
  EXAMPLE PASS: $99.99 -> $99.99 -> $99.99
  EXAMPLE PASS: $10.50 -> $10.50 -> $10.50
  EXAMPLE PASS: $0.01 -> $0.01 -> $0.01
  Result: 5/5 examples pass (false confidence!)

-- Property-Based Test: Detect Float Precision Loss --
  Property FAIL: 1000/1000 inputs failed roundtrip
  Counterexample: 3487.3059919921 -> $3487.31 -> 3487.3100000000

-- Property-Based Test: Robust Integer Cents Roundtrip (1000 iterations) --
  All 1000 iterations PASS: ParseRobust(amount.Format()) == amount

========================================================================
  DEMO 2: Interval Merging – Idempotence Invariant
========================================================================

-- Example-Based Test (Sorted Input – Passes Naive) --
  Input:  [{Start:1 End:3} {Start:2 End:6} {Start:8 End:10} {Start:15 End:18}]
  Merged: [{Start:1 End:6} {Start:8 End:10} {Start:15 End:18}]
  Idempotent: true (passes on sorted example – false confidence!)

-- Property-Based Test: NaiveMerge fails oracle check compared to RobustMerge --
  NaiveMerge vs RobustMerge oracle discrepancies: 85/100
  Counterexample input: [{Start:13 End:23} {Start:-8 End:-6} {Start:-7 End:-6}]
  NaiveMerge(input):    [{Start:13 End:23}]
  RobustMerge(input):   [{Start:-8 End:-6} {Start:13 End:23}]

-- Property-Based Test: RobustMerge passes idempotence (1000 iterations) --
  All 1000 iterations PASS: Merge(Merge(x)) == Merge(x)

-- Property-Based Test: Non-overlapping output invariant (1000 iterations) --
  All 1000 iterations PASS: output intervals are non-overlapping

========================================================================
  DEMO 3: Counterexample Shrinking
========================================================================

-- Invariant: all elements must be non-negative (buggy predicate for demo) --
  Initial failing input (10 elements): [137 18 100 -27 95 107 126 78 -7 79]
  Shrinking steps taken: 21

  Shrinking trace (failing steps only):
    Step   1 remove-right-half              [5 elements] [137 18 100 -27 95] -> FAIL (used)
    Step   3 remove-left-half               [3 elements] [100 -27 95] -> FAIL (used)
    Step   5 remove-left-half               [2 elements] [-27 95] -> FAIL (used)
    Step   6 remove-right-half              [1 elements] [-27] -> FAIL (used)
    Step   9 reduce-val-idx-0-to--13        [1 elements] [-13] -> FAIL (used)
    Step  12 reduce-val-idx-0-to--6         [1 elements] [-6] -> FAIL (used)
    Step  15 reduce-val-idx-0-to--3         [1 elements] [-3] -> FAIL (used)
    Step  18 reduce-val-idx-0-to--1         [1 elements] [-1] -> FAIL (used)

  Minimal counterexample (1 elements): [-1]
  Shrinking SUCCESS: reduced 10 elements to 1 element showing exact cause

========================================================================
  SUMMARY: Example-Based vs Property-Based Testing
========================================================================

  Check                                    Example-Based        Property-Based
  --------------------------------------------------------------------------------
  Float precision bugs caught              NO                   YES
  Unsorted interval edge cases caught      NO                   YES
  Negative value handling verified         NO                   YES
  Empty / minimal cases tested             RARELY               ALWAYS
  Automatic bug isolation (shrinking)      NO                   YES
  Confidence on happy path examples        HIGH (false)         MEDIUM (calibrated)
```

## Final Engineering Status
READY_FOR_ENGINEERING_AUDIT
