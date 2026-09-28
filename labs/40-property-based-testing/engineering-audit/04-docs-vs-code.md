# Docs vs Code Audit

## README Claims vs Code

### Claim 1: Roundtrip Invariant — Float fails, integer cents passes 1,000 roundtrips
- Code: `TestPropertyNaiveCurrencyFails` confirms float failure; `TestPropertyRobustAmountRoundtrip` runs 1,000 iterations.
- Demo: Shows "Property FAIL: 1000/1000 inputs failed roundtrip" for float and "All 1000 iterations PASS" for robust.
- Status: MATCH

### Claim 2: Idempotence Invariant — `Merge(Merge(x)) == Merge(x)`
- Code: `TestPropertyRobustMergeIdempotence` runs 1,000 iterations.
- Demo: "All 1000 iterations PASS"
- Status: MATCH

### Claim 3: Oracle Invariant — NaiveMerge fails vs RobustMerge
- Code: `TestPropertyNaiveMergeFails` uses `quick.Check` to confirm NaiveMerge diverges.
- Demo: "NaiveMerge vs RobustMerge oracle discrepancies: 85/100"
- Status: MATCH

### Claim 4: Shrinking reduces array to minimal `[-1]`
- Code: `TestFindAndShrink` asserts `len(res.Minimal) == 1` and `res.Minimal[0] < 0`.
- Demo: "Minimal counterexample (1 elements): [-1]"
- Status: MATCH

### README Directory Layout
- Lists `engineering/`, `go.mod`, `README.md`, all packages and demo.
- Actual directory matches.
- Status: MATCH

## No Discrepancies Found
No `DOC_CODE_MISMATCH`, `TEST_CLAIM_MISMATCH`, or `RESEARCH_IMPLEMENTATION_MISMATCH` detected.
