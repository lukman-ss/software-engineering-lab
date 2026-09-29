# Source Map: Property-Based Testing

Pemetaan bagian artikel master draft terhadap dokumen riset, implementasi kode sumber, dan pengujian terverifikasi di dalam lab.

---

## 1. Problem & Core Mental Model
- **Research Sources:**
  - `research/01-plan.md`
  - `research/03-evidence.md` (Finding 1)
  - `research/05-report.md` (Section: Research Question & Finding 1)
- **External Citations:**
  - Claessen & Hughes QuickCheck (2000)
  - fast-check documentation ("What is Property-Based Testing?")
- **Implementation & Tests:**
  - `cmd/demo/main.go` (CLI Summary & Side-by-side comparison)

---

## 2. Canonical Invariants (Roundtrip, Idempotence, Oracle)
- **Research Sources:**
  - `research/03-evidence.md` (Finding 2)
  - `research/05-report.md` (Finding 2)
- **Implementation:**
  - `internal/currency/currency.go` (`NaiveCurrency`, `RobustAmount`)
  - `internal/interval/interval.go` (`NaiveMerge`, `RobustMerge`)
- **Tests:**
  - `internal/currency/currency_test.go` (`TestPropertyRobustAmountRoundtrip`, `TestPropertyNaiveCurrencyFails`)
  - `internal/interval/interval_test.go` (`TestPropertyRobustMergeIdempotence`, `TestPropertyRobustMergeNonOverlapping`, `TestPropertyNaiveMergeFails`)

---

## 3. Counterexample Shrinking
- **Research Sources:**
  - `research/03-evidence.md` (Finding 3)
  - `research/05-report.md` (Finding 3)
- **External Citations:**
  - Proptest Book ("Shrinking Basics")
  - Hypothesis ("How Hypothesis Works", David R. MacIver)
- **Implementation:**
  - `internal/shrinker/shrinker.go` (`FindAndShrink`, `ShrinkStep`, `ShrinkResult`, `BuggySortPredicate`)
- **Tests:**
  - `internal/shrinker/shrinker_test.go` (`TestFindAndShrink`)
  - `cmd/demo/main.go` (Demo 3 Shrinking trace logs)

---

## 4. Production Considerations & Common Mistakes
- **Research Sources:**
  - `research/03-evidence.md` (Finding 5: Biased Random Generators)
  - `research/05-report.md` (Finding 5 & Finding 11: Best Practices and Corpus)
  - `research/06-open-questions.md`
- **Implementation:**
  - `internal/currency/currency_test.go` (`RobustAmount.Generate` custom biased distribution)
- **Audit Reports:**
  - `research-audit/07-verdict.md`
  - `engineering-audit/06-verdict.md`
