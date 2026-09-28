# Research Report

## Research Question

Why does Example-Based Testing (unit test biasa) often miss critical bugs, and how does Property-Based Testing (PBT) test hundreds of random input scenarios automatically through universal invariants, shrinking, and test-case generation?

## Executive Summary

Property-Based Testing (PBT) is a software testing paradigm that replaces hand-picked input/output pairs with universal properties (invariants) that must hold for all valid inputs in a domain. A PBT engine generates thousands of random inputs, checks that each property holds, and when a property is violated, shrinks the failing input to the minimal counterexample that reproduces the bug. PBT was introduced by Claessen and Hughes in 2000 as QuickCheck for Haskell and has since been ported to nearly every major programming language. Empirical evidence from fast-check's track record alone documents real bugs and CVEs found in jest, lodash, react, underscore.js, js-yaml, query-string, left-pad, numpy, and jsonwebtoken. The Hypothesis Corpus (2026) provides the largest known dataset of real-world PBT usage: 28,928 tests across 1,529 repositories. No material contradictions were found among authoritative sources.

## Findings

### Finding 1

**Claim:** PBT tests universal invariants rather than specific examples; the key differentiator from example-based testing is that the framework controls input generation and shrinking, while the developer specifies what must always be true.

**Evidence:**
- QuickCheck project page: "The programmer provides a specification of the program, in the form of properties which functions should satisfy, and QuickCheck then tests that the properties hold in a large number of randomly generated cases."
- fast-check: "A 'property' is an assertion of a relationship between a code's input and output which should hold for all sets of inputs."
- Lab specification: "Alih-alih menentukan input & output spesifik, kita mendefinisikan Sifat / Invariant Universal yang harus selalu benar untuk semua kemungkinan input."

**Sources:**
1. Claessen & Hughes QuickCheck page (2000) — https://www.cse.chalmers.se/~rjmh/QuickCheck/
2. fast-check "What is Property-Based Testing?" — https://fast-check.dev/docs/introduction/what-is-property-based-testing/

**Confidence:** HIGH

---

### Finding 2

**Claim:** The canonical categories of PBT properties are: (1) Roundtrip (Decode(Encode(x)) == x), (2) Idempotence (f(f(x)) == f(x)), (3) Hard-to-Prove-but-Easy-to-Verify (e.g., optimal path existence vs. validity), (4) Equivalence/Oracle (new optimized implementation matches old naive implementation).

**Evidence:**
- Lab specification lists all four categories with concrete formulations.
- Hypothesis encode/decode article demonstrates roundtrip on run-length encoding and finds two bugs: empty-string crash and missing count-reset.
- Hackage QuickCheck manual shows `prop_reverse xs = reverse (reverse xs) == xs` — a classic idempotence property.
- Hypothesis Corpus (2026) AI-classified patterns include "inverse relationship between two functions" and "idempotence of repeated operations."

**Sources:**
1. Lab specification (topic spec)
2. Hypothesis "Encode/Decode Invariant" article — https://hypothesis.works/articles/encode-decode-invariant/
3. Hackage QuickCheck manual — https://hackage.haskell.org/package/QuickCheck-2.14.3/docs/Test-QuickCheck.html
4. Hypothesis Corpus — https://huggingface.co/datasets/HypothesisWorks/Hypothesis-Corpus-2026

**Confidence:** HIGH

---

### Finding 3

**Claim:** Counterexample shrinking automatically reduces complex failing inputs (e.g., 500-element arrays) to minimal counterexamples (e.g., `[]int{0, -1}`), dramatically easing debugging.

**Evidence:**
- Proptest shrinking tutorial: `ValueTree::simplify()` and `complicate()` perform binary search; demonstrates finding boundary condition 501 from range 0..10000.
- fast-check: "When fast-check detects a failed test it will provide the developer with a 'counterexample'... fast-check will try removing elements and check whether the test still fails."
- Hypothesis (MacIver): "In a lot of cases it even works better than heavily customized solutions: Often more structured shrinkers get stuck in local minima... Hypothesis can just spot patterns in the data and speculatively shrink them together."
- Hypothesis encode/decode article: shrinking produces minimal counterexample `'110'` for the missing-count-reset bug (3 characters).

**Sources:**
1. Proptest shrinking basics — https://altsysrq.github.io/proptest-book/proptest/tutorial/shrinking-basics.html
2. fast-check "What is Property-Based Testing?"
3. Hypothesis "How Hypothesis Works" — https://hypothesis.works/articles/how-hypothesis-works/
4. Hypothesis "Encode/Decode Invariant" article

**Confidence:** HIGH

---

### Finding 4

**Claim:** Major PBT frameworks exist for Go (testing/quick, gopter), Rust (proptest, quickcheck), Python (Hypothesis), JavaScript/TypeScript (fast-check), Haskell (QuickCheck), and many other languages.

**Evidence:**
- Go `testing/quick`: Standard library package, frozen, uses reflection-based generation via `Generator` interface; default 100 iterations.
- gopter: Third-party Go library inspired by ScalaCheck/QuickCheck; features tighter generator control, shrinkers, regex generators, stateful test support.
- Rust proptest: Inspired by Hypothesis; per-value generation/shrinking via `Arbitrary` trait and `ValueTree`; supports Rust 2018 macros, forking/timeouts, `no_std`.
- Python Hypothesis: Most widely used PBT library; three-layer architecture (Conjecture fuzzer + strategy library + testing interface); byte-stream shrinking.
- JavaScript/TypeScript fast-check: Three core features (random generation, multiple runs, shrinking); biased generators for bug-finding.
- Haskell QuickCheck: Original; type-based generation via `Arbitrary` class; shrinking via `shrink :: a -> [a]`.

**Sources:**
1. Go testing/quick docs — https://pkg.go.dev/testing/quick
2. gopter README — https://github.com/leanovate/gopter
3. Proptest book — https://altsysrq.github.io/proptest-book/
4. Hypothesis docs — https://hypothesis.readthedocs.io/
5. fast-check docs — https://fast-check.dev/docs/
6. QuickCheck Hackage — https://hackage.haskell.org/package/QuickCheck-2.14.3

**Confidence:** HIGH

---

### Finding 5

**Claim:** PBT frameworks intentionally bias their random generators toward common bug patterns (boundaries, duplicates, empty inputs, security-sensitive strings) rather than uniform distributions, making them "designed for bugs."

**Evidence:**
- fast-check: "It's common for numerical code to have edge cases around the values `0`, `1`, or `-1`... fast-check numerical arbitraries ensure coverage of values close to the edges."
- fast-check: "fast-check even supports detecting some security vulnerabilities, such as by including the potentially dangerous string `__proto__` when generating objects. This has successfully detected CVEs in open source projects."
- Hypothesis: "Hypothesis takes a philosophical stance that while users may be responsible for selecting the domain, the property-based testing library—not the user—should be responsible for selecting the distribution."
- Hypothesis distribution techniques: "integers() upweights range endpoints"; "we collect interesting-looking constants from imported source files as seeds"; "swarm testing adds further randomization in stateful testing."

**Sources:**
1. fast-check "Why Property-Based Testing?" — https://fast-check.dev/docs/introduction/why-property-based/
2. fast-check "Track Record" — https://fast-check.dev/docs/introduction/track-record/
3. Hypothesis "Domain and distribution" — https://hypothesis.readthedocs.io/en/latest/explanation/domain.html

**Confidence:** HIGH

---

### Finding 6

**Claim:** Real-world empirical evidence demonstrates that PBT finds bugs in production open-source projects that example-based testing missed, including CVE-level security issues.

**Evidence:**
fast-check track record documents specific bugs with failing inputs and GitHub issue links:
- `jestjs/jest#7941`: `expect(0).toStrictEqual(5e-324)` passes — floating-point comparison bug.
- `left-pad/left-pad#58`: Inconsistent code-point handling for BMP-outside unicode characters.
- `auth0/node-jsonwebtoken#945`: `jwt.sign({valueOf: 0}, 'some-key')` throws TypeError — prototype pollution via `valueOf` key.
- `trekhleb/javascript-algorithms`: Multiple algorithm bugs (counting sort with negatives, KMP with empty string, Rabin-Karp integer overflow, LCS with unicode).
- `nodeca/js-yaml`: Binary format dumps negative integers as `0b-1010` instead of `-0b1010`.
- "recover most of the CVEs related to prototype poisoning reported on lodash."

Hypothesis Corpus (2026): 28,928 Hypothesis tests across 1,529 repositories with runtime data.

Source: Mercurial bugs 4927 and 5031 found via encode/decode testing. Qutebrowser caught 3 bugs in JavaScript escaping via roundtrip invariant.

**Sources:**
1. fast-check "Track Record" — https://fast-check.dev/docs/introduction/track-record/
2. Hypothesis Corpus — https://huggingface.co/datasets/HypothesisWorks/Hypothesis-Corpus-2026
3. Hypothesis "Encode/Decode Invariant" article — https://hypothesis.works/articles/encode-decode-invariant/

**Confidence:** HIGH

---

### Finding 7

**Claim:** Stateful/state-machine property-based testing allows verification of complex sequential behaviors by having the framework generate and execute sequences of actions against both the real system and an abstract model.

**Evidence:**
- Hypothesis "Stateful tests": "With Hypothesis's stateful testing, Hypothesis instead tries to generate not just data but entire tests. You specify a number of primitive actions that can be combined together, and then Hypothesis will try to find sequences of those actions that result in a failure."
- Bundles: Named collections of generated values reusable across rules. Rules can draw from bundles and remove consumed values.
- Invariants: Checked after every rule execution. Preconditions: filter inapplicable rules before execution.
- Example: Database comparison test — `save` and `delete` operations update both a real `DirectoryBasedExampleDatabase` and an in-memory `defaultdict(set)` model; `values_agree` invariant checks consistency. Commenting out the model update line causes Hypothesis to produce a minimal failing sequence (5 steps).

**Sources:**
1. Hypothesis "Stateful tests" docs — https://hypothesis.readthedocs.io/en/latest/stateful.html
2. gopter README — https://github.com/leanovate/gopter (lists "Support for stateful tests" as feature)

**Confidence:** HIGH

---

### Finding 8

**Claim:** QuickCheck uses type-based generation and shrinking via the `Arbitrary` typeclass (`arbitrary :: Gen a`, `shrink :: a -> [a]`), while modern frameworks like Hypothesis and proptest use value-based generation and shrinking, which is more flexible and composition-friendly.

**Evidence:**
- Hackage QuickCheck docs: `class Arbitrary a where { arbitrary :: Gen a; shrink :: a -> [a] }`. Generation and shrinking are defined per-type.
- Proptest book: "Unlike QuickCheck, generation and shrinking is defined on a per-value basis instead of per-type, which makes it more flexible and simplifies composition."
- Hypothesis (MacIver): "Integrating vs. type based shrinking" article (referenced in blog): "The way shrinking is handled in Haskell QuickCheck is bad and the way it works in Hypothesis... is good."

**Sources:**
1. Hackage QuickCheck manual — https://hackage.haskell.org/package/QuickCheck-2.14.3/docs/Test-QuickCheck.html
2. Proptest book intro — https://altsysrq.github.io/proptest-book/intro.html
3. Hypothesis blog — https://hypothesis.works/

**Confidence:** HIGH

---

### Finding 9

**Claim:** Go's `testing/quick` is officially frozen; for serious PBT in Go, gopter (third-party) is recommended as a more capable alternative with shrinkers, stateful tests, and tighter generator control.

**Evidence:**
- Go pkg.go.dev: "Package quick implements utility functions to help with black box testing. The testing/quick package is frozen and is not accepting new features."
- gopter README: "Gopter tries to bring the goodness of ScalaCheck (and implicitly, the goodness of QuickCheck) to Go. It can also be seen as a more sophisticated version of the testing/quick package." Features listed: tighter generator control, shrinkers, regex match generator, stateful test support.

**Sources:**
1. Go testing/quick docs — https://pkg.go.dev/testing/quick
2. gopter README — https://github.com/leanovate/gopter

**Confidence:** HIGH

---

### Finding 10

**Claim:** PBT and example-based testing are complementary; best practice is to use both (hybrid approach), with example-based tests covering specific known scenarios and PBT discovering unexpected edge cases.

**Evidence:**
- fast-check: "Property-based testing is not only about randomizing inputs to find bugs; it's also about helping users to find and to fix the errors." Under Non-tradeoffs: "Although property-based testing is a powerful technique, it should not be viewed as a substitute for traditional example-based testing. Instead, it should be used in conjunction with example-based testing as a complementary approach."
- Proptest book: "Property testing is best used to compliment traditional unit testing (i.e., using specific inputs chosen by hand). Traditional tests can test specific known edge cases, simple inputs, and inputs that were known in the past to reveal bugs, whereas property tests will search for more complicated inputs that cause problems."

**Sources:**
1. fast-check "Why Property-Based Testing?" — https://fast-check.dev/docs/introduction/why-property-based/
2. Proptest book — https://altsysrq.github.io/proptest-book/intro.html

**Confidence:** HIGH

## Areas of Agreement

1. All frameworks define PBT as testing universal properties against randomly generated inputs with automatic shrinking of failures.
2. Default test count is 100 across QuickCheck, Hypothesis, and fast-check.
3. Roundtrip/encode-decode and idempotence are the two most cited canonical invariant patterns.
4. PBT is complementary to, not a replacement for, example-based testing.
5. All major frameworks support shrinking to minimal counterexamples.
6. Go's `testing/quick` is frozen; third-party libraries offer more features.
7. Distribution optimization (boundary upweighting, duplicate injection) is intentional in modern frameworks.

## Areas of Disagreement

1. **PBT vs. fuzzing boundary**: MacIver (narrower view) argues PBT requires a human-authored property that reveals bugs unfindable by direct fuzzing. fast-check (broader view) treats fuzzing as a sub-case of PBT. This is definitional, not practical — all frameworks overlap in capability.

## Limitations

1. **Lack of large-scale controlled empirical studies**: While fast-check's track record and the Hypothesis Corpus provide evidence, there are no peer-reviewed controlled studies comparing bug-detection rates of PBT vs. example-based testing across large codebases.
2. **Generator difficulty**: Writing good generators and discovering good properties remains a skill bottleneck; the Hypothesis Corpus shows stateful tests are rare (small fraction of total tests).
3. **Non-determinism / flakiness**: Shrinking can produce slightly different minimal inputs across runs; Hypothesis addresses this with flakiness verification but it remains a concern.
4. **Framework-specific behavior**: Shrinking algorithms differ (byte-stream vs. value-tree); results are not portable between frameworks.
5. **Property discovery gap**: How engineers discover good properties for a given system is under-documented and not well-studied empirically.
6. **Limited academic literature in non-English sources**: The lab specification is in Bahasa Indonesia but English-language PBT literature dominates; no significant Bahasa Indonesia PBT resources were found.
7. **Hypothesis Corpus size limitation**: While large (28,928 tests), it covers only Python/Hypothesis ecosystem and may not generalize to other languages/frameworks.

## Conclusion

Property-Based Testing is a well-established, actively-used testing paradigm that shifts the developer's effort from enumerating specific test cases to specifying universal invariants. The canonical property patterns (roundtrip, idempotence, equivalence/oracle, satisfiability) cover a broad range of real-world software. Real-world empirical evidence from fast-check's track record — including CVE recovery in lodash and bugs in jest, react, and numpy — demonstrates that PBT finds genuinely important bugs that example-based testing misses. The shrinking mechanism is a consistently cited differentiator that makes PBT failures actionable. The Go ecosystem benefits from the standard library's `testing/quick` for simple cases, but third-party gopter provides the full QuickCheck-like feature set. For the lab exercise (financial currency formatting or interval merging), the three canonical invariants (roundtrip, idempotence, equivalence) are directly applicable and the 1,000-iteration test runner will surface edge cases like NaN, sub-cent rounding, and negative intervals.