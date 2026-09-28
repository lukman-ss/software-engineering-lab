# Code Audit

## Finding 1

Location: internal/currency/currency.go:13-20 (NaiveCurrency.Format/Parse)
Claimed Behavior: Format and Parse functions are symmetric roundtrips for monetary values.
Observed Implementation: Format uses `fmt.Sprintf("$%.2f", dollars)` which rounds to two decimal places; Parse strips '$' and parses as float64. Roundtrip fails for sub-cent precision due to floating point rounding and formatting truncation.
Assessment: FAIL
Severity: HIGH
Notes: This is intentional demonstration; NaiveCurrency is meant to fail property-based test.

## Finding 2

Location: internal/interval/interval.go:12-28 (NaiveMerge)
Claimed Behavior: Merges overlapping intervals.
Observed Implementation: Assumes input sorted by start; fails when unsorted or when intervals touch (adjacent) because condition `iv.Start <= last.End` does not merge touching intervals (e.g., [1,3] and [4,5] remain separate). However note that RobustMerge merges when `iv.Start <= last.End+1` (touching merges). NaiveMerge also fails on unsorted.
Assessment: FAIL
Severity: MEDIUM
Notes: Demonstrates missing precondition (sorted input) and off-by-one for adjacency.

## Finding 3

Location: internal/shrinker/shrinker.go:24-155 (FindAndShrink)
Claimed Behavior: Performs counterexample shrinking via binary sectioning and element reduction.
Observed Implementation: Implements three strategies: halving, element removal, magnitude reduction toward zero. Works as intended for the demo predicate (BuggySortPredicate). No obvious bugs.
Assessment: PASS
Severity: LOW
Notes: Shrinking algorithm is functional; could be improved but not required.

## Finding 4

Location: internal/currency/currency.go:43-69 (ParseRobust)
Claimed Behavior: Parses currency string with optional sign, dollar sign, and exactly two decimal places.
Observed Implementation: Strict format requires exactly two decimal digits; rejects inputs like "$5" or "$5.0". This matches spec but could be considered restrictive; however used only for roundtrip invariant.
Assessment: PASS
Severity: LOW
Notes: Accepts negative amounts.

## Finding 5

Location: internal/interval/interval.go:30-56 (RobustMerge)
Claimed Behavior: Correctly merges intervals regardless of input order, merges touching intervals.
Observed Implementation: Sorts copy, then merges with condition `iv.Start <= last.End+1`. Works correctly.
Assessment: PASS
Severity: LOW
Notes: None.

## Finding 6

Location: cmd/demo/main.go (demo functions)
Claimed Behavior: Demonstrates PBT vs example-based testing.
Observed Implementation: Runs deterministic examples and property tests using quick.Check. Output matches expectations.
Assessment: PASS
Severity: LOW
Notes: Demo is illustrative and passes.