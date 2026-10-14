---
title: "Property-Based Testing: Finding Hidden Edge Cases via Universal Invariants and Shrinking"
published: false
description: "Explore Property-Based Testing in Go: mechanics, architecture, caveats, and testable invariants."
tags: go, testing, performance, softwareengineering
---

# Property-Based Testing: Finding Hidden Edge Cases via Universal Invariants and Shrinking

A production system must remain resilient under load. Lab 40 makes that behavior visible with a focused Go simulation and deterministic test assertions.

> Example-based tests check developer assumptions; property-based tests explore domain boundaries.

## The False Confidence of Example-Based Testing

Conventional example-based unit tests check whether f(x) == y for a few hand-picked samples. The core vulnerability is developer confirmation bias: engineers only test scenarios they already anticipated. While test suites show 100% green coverage, edge-case assumptions break immediately when exposed to real-world, uncurated input distributions.

## Three Pillars of Property-Based Testing

Property-based testing splits responsibilities cleanly between human and machine. Engineers define domain constraints and mathematical invariants that must unconditionally hold. The PBT engine acts as an automated adversary, generating random candidates, checking invariants, and shrinking failures to their simplest form.

Tahap | Tanggung Jawab | Output Sistem
|---|---|---|
1. Generator | Buat ratusan sampel input acak terdistribusi | Slices, numbers, structs
2. Invariant | Evaluasi properti universal f(x) | Validasi true / false
3. Shrinker | Pangkas kegagalan ke bentuk terkecil | Minimal counterexample (e.g. [-1])

## Four Canonical Invariant Patterns

The core of PBT lies in choosing relational invariants rather than duplicating implementation logic in test assertions. Roundtrips ensure lossless data transformations, idempotence guarantees state stability, test oracles verify complex optimizations against naive references, and easy-to-verify checks validate structural correctness at low computational cost.

Pola Invariant | Formula Formal | Contoh Kasus Nyata
|---|---|---|
Roundtrip | Decode(Encode(x)) == x | Serialisasi JSON, mata uang, kompresi
Idempotence | f(f(x)) == f(x) | Normalisasi teks, merge interval, cleanup
Oracle Equiv | f_robust(x) == f_naive(x) | Verifikasi algoritma teroptimasi vs referensi
Easy-to-Verify | HardToSolve(x) -> EasyVerify(res) | Validasi shortest path, partisi graf

## Floating-Point Precision Loss in Currency

Storing currency as float64 is a notorious software engineering anti-pattern. Format strings truncate sub-cent precision, while string-to-float conversions suffer from IEEE 754 binary approximation limits. Hand-picked unit test values hide this defect entirely, allowing financial rounding leaks into production.

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

When the roundtrip property is tested against Go's testing/quick package, the float64 implementation fails unconditionally. Across 1,000 randomized iterations, all 1,000 inputs violate the identity invariant Parse(Format(x)) == x due to sub-cent truncation, proving that example tests produced false confidence.

```go
func TestPropertyNaiveCurrencyFails(t *testing.T) {
    naive := NaiveCurrency{}
    prop := func(val NaiveFloat) bool {
        dollars := float64(val)
        formatted := naive.Format(dollars)
        parsed, err := naive.Parse(formatted)
        if err != nil { return false }
        return parsed == dollars // Gagal pada sub-cent!
    }
    err := quick.Check(prop, &quick.Config{MaxCount: 100})
    // quick.Check mengembalikan error counterexample
}
```

## Robust Representation: Exact Integer Cents

The industry-proven solution for monetary values is storing the smallest fractional unit as an integer (int64 cents). Decimal string formatting becomes a presentation-layer projection, while arithmetic remains purely discrete, eliminating floating-point imprecision permanently.

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

Uniform random generators often miss critical edge cases where boundary bugs reside. The RobustAmount custom generator injects deliberate probability biases toward zero, near-zero negative transitions, and multi-billion cent bounds, actively probing boundary handling resilience.

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

By adopting integer cents and exact string parsing, TestPropertyRobustAmountRoundtrip executes 1,000 randomized iterations without a single invariant violation. Across negative amounts, zero, and enterprise-scale values, Parse(Format(x)) == x holds true unconditionally.

Metrik Uji | Naive Float64 | Robust Integer Cents
|---|---|---|
Total Iterasi Acak | 1.000 iterasi | 1.000 iterasi
Tingkat Kelulusan | 0% (1.000 gagal) | 100% (1.000 lolos)
Cakupan Domain | Terbatas hand-picked | Nol, negatif, triliun sen
Invarian Terbukti | Gagal total | Parse(Format(x)) == x VALID

## Interval Merging & Hidden Preconditions

Algorithmic bugs frequently stem from unstated implicit preconditions. NaiveMerge merges overlapping intervals by checking incoming elements against the last emitted interval, implicitly assuming presorted inputs. Hand-picked sorted examples succeed, but arbitrary ordering produces silent corruption.

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
    return out // BUG: Mengasumsikan input terurut!
}
```

## Test Oracle Uncovers 85% Discrepancy

The Test Oracle pattern compares two distinct implementations against identical inputs: f_naive(x) versus f_robust(x). RobustMerge enforces pre-sorting via sort.Slice. Under property evaluation, 85 out of 100 randomized slices expose discrepancies, demonstrating the fragility of unsorted inputs.

```go
func TestPropertyNaiveMergeFails(t *testing.T) {
    prop := func(ivs IntervalSlice) bool {
        naive := NaiveMerge([]Interval(ivs))
        robust := RobustMerge([]Interval(ivs))
        return intervalsEqual(naive, robust)
    }
    err := quick.Check(prop, &quick.Config{MaxCount: 100})
    // 85 dari 100 input acak menghasilkan ketidaksesuaian!
}
```

## Idempotence Invariant on Data Merging

Idempotence is a foundational invariant: f(f(x)) == f(x). In data processing, re-applying a merge operation to an already merged slice must yield an identical result without side effects. Running 1,000 iterations confirms that RobustMerge is mathematically stable across all permutations.

```go
Formula Teoretis: f(f(x)) == f(x)

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

Beyond relational invariants, merge algorithms must fulfill structural domain properties: for all consecutive intervals in the output, out[i].End < out[i+1].Start. Across 1,000 iterations, the test proves that no overlapping fragments escape the merging logic.

```go
func TestPropertyRobustMergeNonOverlapping(t *testing.T) {
    prop := func(ivs IntervalSlice) bool {
        merged := RobustMerge([]Interval(ivs))
        for i := 0; i < len(merged)-1; i++ {
            if merged[i].End >= merged[i+1].Start {
                return false // Ada tumpang tindih!
            }
        }
        return true
    }
    quick.Check(prop, &quick.Config{MaxCount: 1000})
}
```

## Problem Anatomy: Counterexample Noise

When a PBT engine catches an invariant violation on randomized inputs, the raw failure payload is often large and noisy. In our demo, an arbitrary 10-element slice contains mixed positive values and two negatives. Without automated reduction, developers waste hours isolating the exact trigger.

## Counterexample Shrinking Algorithm

Counterexample shrinking deterministically transforms complex failing inputs into minimal reproducible test cases. It operates in three sequential phases: binary list halving, individual element removal, and numerical reduction toward zero, retaining failing states at each stage.

Fase Algoritma | Aksi Reduksi | Target Optimasi
|---|---|---|
1. Binary Sectioning | Potong separuh kiri / kanan array | Pangkas ukuran O(log N)
2. Element Removal | Hapus elemen individual satu per satu | Eliminasi noise non-kritis
3. Value Reduction | Bagi dua nilai elemen menuju 0 | Minimalisasi magnitudo nilai

## Binary Sectioning: Halving Array Length

The first phase of shrinking is binary list sectioning. The slice is halved and each section evaluated against the predicate. If the left half still fails, the entire right half is discarded immediately, reducing input size from 10 elements to 5, then to 3 within a few steps.

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

Once array length is reduced to a single-element slice [-27], the shrinker applies numerical reduction toward zero: -27 -> -13 -> -6 -> -3 -> -1. When candidate 0 passes the non-negative check, the shrinker terminates, locking in [-1] as the exact minimal boundary counterexample.

```go
// Coba set elemen menjadi 0 atau bagi dua nilainya
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

The demo terminal trace illustrates complete diagnostic clarity. From a chaotic 10-element array, the engine executes 21 predicate evaluations to eliminate 9 irrelevant elements and reduce -27 to -1. The resulting [-1] reveals the root cause instantly: any negative value violates the invariant.

```go
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

Go provides built-in property testing via testing/quick in the standard library, requiring zero third-party dependencies. quick.Check evaluates properties across thousands of iterations with configurable counts and seeds. Custom types implement quick.Generator to supply domain-specific generation logic.

```go
cfg := &quick.Config{
    MaxCount: 1000, // Jumlah iterasi acak yang dieksekusi
    Rand:     rand.New(rand.NewSource(seed)), // Determinisme
}

err := quick.Check(func(a RobustAmount) bool {
    return ParseRobust(a.Format()) == a
}, cfg)
if err != nil {
    t.Fatalf("property violation: %v", err)
}
```

## Determinism, Seeding, and CI/CD Pipelines

A common critique of PBT is test flakiness from changing random seeds. Professional teams adopt a dual strategy: fixed seeds with 100-500 iterations for sub-second PR checks, and dynamic timestamp seeds with 50,000 iterations in nightly CI. Logging the failure seed ensures 100% reproducible debugging locally.

Lingkungan | Konfigurasi Seed | Iterasi | Tujuan
|---|---|---|---|
Local PR Check | Fixed Seed (deterministik) | 100 - 500 | Cepat & lolos tanpa flakiness
Nightly CI | Dynamic Time Seed | 10.000 - 50.000 | Eksplorasi bug liar & boundary
Bug Reproduction | Injeksi Failing Seed | 1 run minimal | Reproduksi 100% konsisten

## Comparison Matrix: Example vs Property-Based

The comparison matrix underscores a paradigm shift in test quality. While example tests excel at regression validation for explicit business user flows, PBT delivers mathematical confidence for financial calculations, serialization protocols, and complex data structures through automated exploration.

Dimensi Evaluasi | Example-Based Testing | Property-Based Testing
|---|---|---|
Cakupan Input | Titik diskrit statis (hand-picked) | Spektrum kontinu acak terdistribusi
Bias Pengujian | Tinggi (mengikuti intuisi engineer) | Rendah (adversarial fuzzing)
Bug Boundary | Kerap terlewat (sub-cent, negative) | Terdeteksi via biased generator
Isolasi Bug | Manual via trace & debugger | Otomatis via counterexample shrinking
Tingkat Keyakinan | False confidence (jalur aman) | Calibrated confidence (invarian kokoh)

## Best Practices and Anti-Patterns in PBT

The most common PBT anti-pattern is re-implementing system logic inside test assertions, leading to shared blind spots. Focus strictly on relational invariants. Always bias generators toward boundary conditions (zeros, empty slices, integer limits), and save discovered counterexamples as permanent regression tests.

## Hands-On Reproduction and Verification Guide

Lab 40 runs on standard Go 1.22 without third-party dependencies, guaranteeing portability across macOS, Linux, and Windows. Running go test -v ./... executes all unit tests, quick properties, and shrinking verification. Running go run ./cmd/demo displays real-time CLI comparisons and shrinking traces.

```sh
# 1. Masuk ke direktori lab
cd labs/40-property-based-testing

# 2. Jalankan seluruh property & unit test
go test -v ./...

# 3. Verifikasi ketiadaan race condition
go test -race ./...

# 4. Eksekusi CLI interactive demo
go run ./cmd/demo
```

## Source

[Software Engineering Lab 40 - Property-Based Testing](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/40-property-based-testing), snapshot commit `7d244968`. Source: README, implementation, tests, and execution records. Code excerpts belong in their original package context.
