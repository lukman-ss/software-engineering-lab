---
title: "Property-Based Testing: Finding Hidden Edge Cases via Universal Invariants and Shrinking"
published: false
description: "Explore Property-Based Testing in Go: mechanics, architecture, caveats, and testable invariants."
tags: go, testing, performance, softwareengineering
---

# Property-Based Testing: Finding Hidden Edge Cases via Universal Invariants and Shrinking

A production system must remain resilient under uncurated input distributions. Lab 40 demonstrates this vulnerability by contrasting naive assumptions against randomized property evaluation and automated failure minimization.

> Example-based tests check developer assumptions; property-based tests explore domain boundaries.

## The False Confidence of Example-Based Testing

Conventional example-based unit tests assert that f(x) == y for a few hand-picked inputs. Confirmation bias leads developers to test scenarios they already anticipated. While a test suite shows 100% statement coverage, passing static examples does not automatically validate unhandled boundary conditions or unexpected combinations in real workloads.

## Three Pillars of Property-Based Testing

Property-based testing separates domain specification from test case generation. Engineers define relational invariants that must hold true, while the engine generates randomized inputs and shrinks failing cases to their minimal form. This automated search explores boundary spaces, but verifying lab properties does not prove every production concurrency scenario.

Stage | Responsibility | System Output
|---|---|---|
1. Generator | Produce hundreds of distributed random inputs | Slices, numbers, structs
2. Invariant | Evaluate universal property f(x) | Validation true / false
3. Shrinker | Reduce failure to smallest form | Minimal counterexample (e.g. [-1])

## Four Canonical Invariant Patterns

Rather than duplicating application logic inside test assertions, PBT evaluates relational invariants. Roundtrip checks verify serialization symmetry, idempotence verifies repeatable operations without state mutation, test oracles compare optimized paths against reference models, and easy-to-verify checks validate structural results. Passing these properties bounds risk but does not replace domain-specific contract testing.

Invariant Pattern | Formal Definition | Production Use Case
|---|---|---|
Roundtrip | Decode(Encode(x)) == x | JSON serialization, currency, compression
Idempotence | f(f(x)) == f(x) | Text normalization, interval merging, state cleanup
Oracle Equiv | f_robust(x) == f_naive(x) | Optimized algorithm vs reference model
Easy-to-Verify | HardToSolve(x) -> EasyVerify(res) | Shortest path verification, graph partitioning

## Floating-Point Precision Loss in Currency

Storing monetary amounts in float64 creates latent precision bugs due to IEEE 754 binary approximations. Naive format strings truncate sub-cent values without error signaling. Hand-picked tests with whole dollar values pass without issue, yet static test success does not imply decimal integrity across fractional arithmetic.

```go
type NaiveCurrency struct{}

func (n NaiveCurrency) Format(dollars float64) string {
    return fmt.Sprintf("$%.2f", dollars)
}

func (n NaiveCurrency) Parse(s string) (float64, error) {
    s = strings.TrimPrefix(s, "$")
    return strconv.ParseFloat(s, 64)
}
```

## PBT Catches Float Leaks with 100% Failure Rate

When the roundtrip property Parse(Format(x)) == x is evaluated via Go testing/quick, float64 fails on all 1,000 randomized iterations. Truncation destroys fractional cents, proving the roundtrip is not an identity function. However, detecting serialization defects on float64 only applies to this execution path and does not address external database precision types.

```go
func TestPropertyNaiveCurrencyFails(t *testing.T) {
    naive := NaiveCurrency{}
    prop := func(val NaiveFloat) bool {
        dollars := float64(val)
        formatted := naive.Format(dollars)
        parsed, err := naive.Parse(formatted)
        if err != nil { return false }
        return parsed == dollars // Fails on sub-cent truncation
    }
    err := quick.Check(prop, &quick.Config{MaxCount: 100})
    // quick.Check returns counterexample error
}
```

## Robust Representation: Exact Integer Cents

Representing monetary values as discrete int64 cents eliminates floating-point rounding artifacts. String formatting serves strictly as a presentation-layer projection, while internal arithmetic remains discrete. Storing integer cents restricts decimal drift in application memory, but production still requires database schema constraints to mitigate overflow risks.

```go
type RobustAmount struct {
    Cents int64
}

func (r RobustAmount) Format() string {
    absCents := r.Cents; sign := ""
    if absCents < 0 { sign = "-"; absCents = -absCents }
    dollars := absCents / 100; cents := absCents % 100
    return fmt.Sprintf("%s$%d.%02d", sign, dollars, cents)
}
```

## Custom Generator with Boundary Value Bias

Uniform random generators frequently miss critical edge cases where defects concentrate. The RobustAmount generator injects deliberate probability bias toward zero, near-zero negative transitions (-100 to +99 cents), and multi-billion cent bounds. Generating boundary inputs stresses parsing logic, but biased sampling does not eliminate the need for fuzzing arbitrary string inputs.

```go
func (RobustAmount) Generate(r *rand.Rand, size int) reflect.Value {
    switch r.Intn(5) {
    case 0: cents = 0                         // Boundary zero
    case 1: cents = int64(r.Intn(200) - 100)  // Near-zero negatives
    case 2: cents = int64(r.Intn(1000000))    // Normal positive
    case 3: cents = -int64(r.Intn(1000000))   // Normal negative
    default: cents = r.Int63n(1e12) - 5e11    // Multi-billion bounds
    }
    return reflect.ValueOf(RobustAmount{Cents: cents})
}
```

## 1,000 Iterations Pass Integer Roundtrip

With integer cents and exact string parsing, TestPropertyRobustAmountRoundtrip succeeds across 1,000 randomized iterations. Across negative values, zero, and enterprise-scale amounts, ParseRobust(Format(x)) == x holds consistently. This test validates mathematical roundtrip in memory, but this in-memory check does not validate third-party payment gateway serialization.

Metric | Naive Float64 | Robust Integer Cents
|---|---|---|
Random Iterations | 1,000 runs | 1,000 runs
Pass Rate | 0% (1,000 failed) | 100% (1,000 passed)
Domain Scope | Hand-picked examples | Zero, negative, trillion cents
Invariant Verified | Total failure | Parse(Format(x)) == x VALID

## Interval Merging & Hidden Preconditions

Algorithmic flaws often stem from unstated preconditions in function contracts. NaiveMerge merges overlapping intervals by checking adjacent entries against the last emitted interval, implicitly assuming presorted inputs. Presorted test cases succeed, but arbitrary ordering corrupts output. A passing test on sorted slices does not automatically prove correctness for unsorted production payloads.

```go
func NaiveMerge(intervals []Interval) []Interval {
    if len(intervals) == 0 { return nil }
    out := []Interval{intervals[0]}
    for _, iv := range intervals[1:] {
        last := &out[len(out)-1]
        if iv.Start <= last.End {
            if iv.End > last.End { last.End = iv.End }
        } else { out = append(out, iv) }
    }
    return out // BUG: Assumes sorted input
}
```

## Test Oracle Uncovers 85% Discrepancy

The Test Oracle pattern compares an experimental function against a reference implementation: f_naive(x) versus f_robust(x). RobustMerge enforces pre-sorting via sort.Slice before merging. Under property evaluation, 85 out of 100 randomized slices expose discrepancies. This discrepancy rate highlights the vulnerability, but oracle equivalence only holds when the reference model is correct.

```go
func TestPropertyNaiveMergeFails(t *testing.T) {
    prop := func(ivs IntervalSlice) bool {
        naive := NaiveMerge([]Interval(ivs))
        robust := RobustMerge([]Interval(ivs))
        return intervalsEqual(naive, robust)
    }
    err := quick.Check(prop, &quick.Config{MaxCount: 100})
    // 85 of 100 random inputs yield discrepancies
}
```

## Idempotence Invariant on Data Merging

Idempotence requires that repeating an operation on its own output causes no further mutation: f(f(x)) == f(x). Re-applying RobustMerge to an already merged slice produces an identical result across 1,000 iterations. Validating idempotence verifies operational convergence, but scope follows tests and does not prove concurrent thread-safety without external synchronization.

```go
// Mathematical Relation: f(f(x)) == f(x)

func TestPropertyRobustMergeIdempotence(t *testing.T) {
    prop := func(ivs IntervalSlice) bool {
        once := RobustMerge([]Interval(ivs))
        twice := RobustMerge(once)
        return intervalsEqual(once, twice)
    }
    quick.Check(prop, &quick.Config{MaxCount: 1000})
}
```

## Structural Invariant: Non-Overlapping Slices

Beyond relational properties, interval merging must satisfy the structural constraint that consecutive output intervals never overlap (out[i].End < out[i+1].Start). Testing 1,000 random slices verifies that no overlapping segments remain. Satisfying structural checks bounds partition anomalies, but this simulation does not address memory limits when merging millions of entries.

```go
func TestPropertyRobustMergeNonOverlapping(t *testing.T) {
    prop := func(ivs IntervalSlice) bool {
        merged := RobustMerge([]Interval(ivs))
        for i := 0; i < len(merged)-1; i++ {
            if merged[i].End >= merged[i+1].Start {
                return false // Overlapping interval detected
            }
        }
        return true
    }
    quick.Check(prop, &quick.Config{MaxCount: 1000})
}
```

## Problem Anatomy: Counterexample Noise

When randomized generators trigger an invariant failure, raw counterexamples often contain large, noisy data structures. In our demonstration, a failing slice holds 10 elements with arbitrary values and two negatives (-27 and -7). Without automated reduction, developers spend excessive time isolating the minimal failure trigger from irrelevant input clutter.

Input Array: `[137, 18, 100, -27, 95, 107, 126, 78, -7, 79]`
Diagnostic Questions:
- Is the failure triggered by magnitude 137 or slice length 10?
- Do element positions interact to cause the defect?
- What is the minimal trigger required to reproduce the bug?

## Counterexample Shrinking Algorithm

Counterexample shrinking deterministically reduces a complex failing input to its minimal reproducible form. It executes three ordered phases: binary chunk removal to reduce array length, individual element removal to eliminate noise, and numerical reduction toward zero. While shrinking isolates root causes quickly, finding a minimal counterexample does not automatically fix the underlying code defect.

Phase | Action | Optimization Target
|---|---|---|
1. Binary Sectioning | Bisect array into left/right halves | Reduce length O(log N)
2. Element Removal | Drop single elements sequentially | Eliminate non-critical noise
3. Value Reduction | Halve integer values toward zero | Minimize value magnitude

## Binary Sectioning: Halving Array Length

The first phase bisects the input slice and evaluates the failure predicate on each half. If the left half continues to trigger the failure, the entire right half is discarded immediately. This reduces the 10-element array to 5, and then to 3 elements within a few steps. Halving cuts noise rapidly, but binary reduction only succeeds when the failure property is monotonic.

```go
if len(current) > 1 {
    mid := len(current) / 2
    left := current[:mid]
    if !predicate(left) {
        current = append([]int(nil), left...)
        changed = true; continue
    }
    right := current[mid:]
    if !predicate(right) {
        current = append([]int(nil), right...)
        changed = true; continue
    }
}
```

## Element Removal and Value Reduction

Once array length reaches a single element [-27], the shrinker applies numerical reduction toward zero. The value -27 is halved iteratively: -27 -> -13 -> -6 -> -3 -> -1. When candidate 0 passes the non-negative check, the shrinker terminates, locking in [-1]. Value reduction isolates the exact threshold, but this heuristic requires additional validation for non-linear mathematical predicates.

```go
// Try resetting element to 0 or halving its value
candidates := []int{0, current[i] / 2}
for _, cand := range candidates {
    testSlice := append([]int(nil), current...)
    testSlice[i] = cand
    if !predicate(testSlice) {
        current = testSlice
        changed = true; break
    }
}
```

## 21-Step Shrinking Trace to Minimal Bug

The CLI demo traces complete execution progress across 21 predicate evaluations. From a chaotic 10-element array, the engine discards 9 irrelevant elements and shrinks -27 to -1. The resulting [-1] isolates the root cause directly: any negative integer triggers failure. A 21-step trace confirms effective reduction in this lab, but reduction step count varies across different input distributions.

```text
Input Awal (10 elemen): [137 18 100 -27 95 107 126 78 -7 79]
Step  1: remove-right-half         [5 elem] [137 18 100 -27 95] -> FAIL
Step  3: remove-left-half          [3 elem] [100 -27 95]        -> FAIL
Step  5: remove-left-half          [2 elem] [-27 95]            -> FAIL
Step  6: remove-right-half         [1 elem] [-27]               -> FAIL
Step  9: reduce-val to -13         [1 elem] [-13]               -> FAIL
Step 12: reduce-val to -6          [1 elem] [-6]                -> FAIL
Step 15: reduce-val to -3          [1 elem] [-3]                -> FAIL
Step 18: reduce-val to -1          [1 elem] [-1]                -> FAIL
Hasil Akhir (Minimal Counterexample): [-1] (Tepat 21 langkah evaluasi!)
```

## Go Standard Library: testing/quick

Go provides property-based testing in the standard library via testing/quick without third-party dependencies. quick.Check evaluates properties across configurable iteration counts with custom pseudo-random generators. Built-in generators handle primitive types, while quick.Generator supports custom domain types. Using testing/quick avoids external dependencies, but testing/quick lacks integrated shrinking out of the box.

```go
cfg := &quick.Config{
    MaxCount: 1000, // Number of random iterations executed
    Rand:     rand.New(rand.NewSource(seed)), // Deterministic seed
}

err := quick.Check(func(a RobustAmount) bool {
    return ParseRobust(a.Format()) == a
}, cfg)
if err != nil {
    t.Fatalf("property violation: %v", err)
}
```

## Determinism, Seeding, and CI/CD Pipelines

A common critique of PBT is test flakiness caused by changing random seeds. Engineering teams address this by running fixed seeds with 100–500 iterations for fast, reproducible pull request checks, while nightly CI uses dynamic time seeds with 50,000 iterations for broader exploration. Logging failure seeds enables local reproduction, but reproducible seeds do not eliminate environmental timeout discrepancies.

Environment | Seed Strategy | Iterations | Objective
|---|---|---|---|
Local PR Check | Fixed Seed (deterministic) | 100 - 500 | Fast execution without flakiness
Nightly CI | Dynamic Time Seed | 10,000 - 50,000 | Explore rare boundary faults
Bug Reproduction | Inject Failing Seed | 1 run minimal | Consistent failure reproduction

## Comparison Matrix: Example vs Property-Based

Example-based tests validate specific business requirements and regression flows, while property-based testing evaluates relational invariants across broad domain distributions. PBT excels at uncovering numerical truncation and sorting assumptions through biased generators and automated shrinking. While PBT increases confidence across continuous inputs, property tests do not replace end-to-end integration tests.

Evaluation Dimension | Example-Based Testing | Property-Based Testing
|---|---|---|
Input Coverage | Hand-picked discrete values | Continuous randomized domain
Testing Bias | High (developer intuition) | Low (adversarial exploration)
Boundary Bugs | Frequently missed (sub-cents) | Detected via biased generators
Bug Isolation | Manual via logs and debugger | Automated counterexample shrinking
Confidence Level | Verified on specific paths | Calibrated on universal invariants

## Best Practices and Anti-Patterns in PBT

A frequent anti-pattern is re-implementing system logic inside test assertions, creating shared blind spots between code and test. Assertions must focus strictly on relational invariants like roundtrip and idempotence. Generators should explicitly include boundary values (zero, negative boundaries, empty slices), and discovered counterexamples should be retained as permanent regression unit tests.

## Hands-On Reproduction and Verification Guide

Lab 40 runs on standard Go 1.22 without third-party packages. Running `go test -v ./...` executes currency roundtrips, interval oracle checks, and custom shrinking algorithms. The demo command displays interactive CLI output and reduction traces. Passing all unit and property tests verifies the current implementation, but local execution does not replace distributed system verification.

```sh
# 1. Enter lab directory
cd labs/40-property-based-testing

# 2. Run all unit and property tests
go test -v ./...

# 3. Verify absence of memory races
go test -race ./...

# 4. Execute CLI interactive demo
go run ./cmd/demo
```

## Source

[Software Engineering Lab 40 - Property-Based Testing](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/40-property-based-testing), snapshot commit `7d244968`. Source: README, implementation, tests, and execution records. Code excerpts belong in their original package context.
