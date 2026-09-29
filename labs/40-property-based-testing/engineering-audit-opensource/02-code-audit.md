# Code Audit: labs/40-property-based-testing

## Finding 1: Roundtrip Format and Parse Precision in Currency

Location: `internal/currency/currency.go:9-70`
Claimed Behavior: NaiveCurrency suffers precision loss when formatting float64 values; RobustAmount maintains exact cent values using `int64` and roundtrips cleanly.
Observed Implementation:
- `NaiveCurrency.Format` uses `fmt.Sprintf("$%.2f", dollars)`, discarding fractional cents. `NaiveCurrency.Parse` trims `$` and parses via `strconv.ParseFloat`.
- `RobustAmount` uses integer cents. Negative values format with leading minus sign `-$X.YY`. `ParseRobust` parses sign, verifies `$` prefix, and verifies 2-digit decimal precision before combining `(dollars*100 + cents) * sign`.
Assessment: PASS
Severity: LOW
Notes: Correct mathematical separation of integer cents and decimal representation. Handles negative values, zeroes, and boundary amounts cleanly.

---

## Finding 2: Unsorted Input Handling in Interval Merging

Location: `internal/interval/interval.go:10-56`
Claimed Behavior: NaiveMerge assumes inputs are pre-sorted and fails on unsorted slices. RobustMerge sorts copy of inputs prior to merging.
Observed Implementation:
- `NaiveMerge` performs a single linear scan checking `iv.Start <= last.End`, failing when intervals appear out of ascending start order.
- `RobustMerge` copies slice, sorts by Start ascending then End ascending using `sort.Slice`, and checks `iv.Start <= last.End+1` to merge adjacent intervals.
Assessment: PASS
Severity: LOW
Notes: RobustMerge defensively copies inputs without mutating caller slice. Adjacency rule `+1` handles discrete integer intervals correctly.

---

## Finding 3: Counterexample Shrinking Strategy and Termination

Location: `internal/shrinker/shrinker.go:26-155`
Claimed Behavior: Reduces failing slice counterexamples to minimal failing subset using binary chunk removal, element deletion, and value reduction towards zero.
Observed Implementation:
- Implements 3-phase reduction loop: binary halving (remove right/left half), element deletion (`remove-elem-idx-i`), and magnitude reduction (`cand = 0` or `cand = current[i] / 2`).
- Loop terminates when no strategy yields a failing counterexample (`changed == false`).
- Steps are logged sequentially into `res.Steps`.
Assessment: PASS
Severity: LOW
Notes: Greedy reduction terminates reliably because array size strictly decreases or non-zero values strictly move toward zero.

---

## Finding 4: Error Handling and Propagation

Location: `internal/currency/currency.go:50-67`, `internal/shrinker/shrinker.go:44-46`
Claimed Behavior: Clean error returns on invalid inputs or failed counterexample searches.
Observed Implementation:
- `ParseRobust` returns descriptive formatting errors on missing `$`, bad decimal positions, or invalid integer components.
- `FindAndShrink` returns error if no failing counterexample is found within `maxAttempts`.
Assessment: PASS
Severity: LOW
Notes: No panics or swallowed errors observed.
