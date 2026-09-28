# Test Audit

## Coverage Analysis

### Currency (`internal/currency/currency_test.go`)
- `TestExampleBasedNaiveCurrency`: Happy path unit test demonstrating false confidence with hand-picked cases.
- `TestPropertyNaiveCurrencyFails`: Uses `testing/quick.Check` with custom `Generate` to prove float precision failure.
- `TestPropertyRobustAmountRoundtrip`: 1,000 iterations over wide range of positive, zero, and negative integer cent values verifying `ParseRobust(amount.Format()) == amount`.
- Coverage: Happy path, negative values, large integer bounds, failure proof.

### Interval (`internal/interval/interval_test.go`)
- `TestExampleBasedIntervalMerge`: Example-based test showing false confidence on pre-sorted arrays.
- `TestPropertyNaiveMergeFails`: Oracle invariant proving `NaiveMerge` differs from `RobustMerge` on arbitrary random intervals.
- `TestPropertyRobustMergeIdempotence`: 1,000 iterations proving `Merge(Merge(x)) == Merge(x)`.
- `TestPropertyRobustMergeNonOverlapping`: 1,000 iterations proving consecutive output intervals never overlap.
- Coverage: Happy path, unsorted inputs, boundary conditions, invariant properties.

### Shrinker (`internal/shrinker/shrinker_test.go`)
- `TestFindAndShrink`: Verifies searching for counterexamples and shrinking array of negative values to a minimal 1-element slice containing a negative number.
- Coverage: Deterministic shrinking, size reduction verification, minimal element validation.

## Execution Results

```text
go test -v -count=1 ./...
=== RUN   TestExampleBasedNaiveCurrency
--- PASS: TestExampleBasedNaiveCurrency (0.00s)
=== RUN   TestPropertyNaiveCurrencyFails
--- PASS: TestPropertyNaiveCurrencyFails (0.00s)
=== RUN   TestPropertyRobustAmountRoundtrip
--- PASS: TestPropertyRobustAmountRoundtrip (0.00s)
PASS
ok  	labs/40-property-based-testing/internal/currency	0.316s
=== RUN   TestExampleBasedIntervalMerge
--- PASS: TestExampleBasedIntervalMerge (0.00s)
=== RUN   TestPropertyNaiveMergeFails
--- PASS: TestPropertyNaiveMergeFails (0.00s)
=== RUN   TestPropertyRobustMergeIdempotence
--- PASS: TestPropertyRobustMergeIdempotence (0.00s)
=== RUN   TestPropertyRobustMergeNonOverlapping
--- PASS: TestPropertyRobustMergeNonOverlapping (0.00s)
PASS
ok  	labs/40-property-based-testing/internal/interval	0.104s
=== RUN   TestFindAndShrink
--- PASS: TestFindAndShrink (0.00s)
PASS
ok  	labs/40-property-based-testing/internal/shrinker	0.315s
```

Race detector:
```text
go test -race -count=1 ./...
PASS in all packages without data races.
```

Demo execution:
```text
go run ./cmd/demo
Executes DEMO 1, DEMO 2, DEMO 3, and SUMMARY cleanly and predictably.
```

Assessment: PASS. Test suite thoroughly exercises both failing and robust implementations.
