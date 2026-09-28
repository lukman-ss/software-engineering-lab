# Claim Audit

## Claim 1

Claim: Property-Based Testing originated with QuickCheck, created by Koen Claessen and John Hughes at Chalmers University of Technology and first presented at ICFP 2000.

Location:
`research/03-evidence.md` (Evidence 1), `research/05-report.md` (Finding 1)

Evidence Provided:
Direct paper citations and Chalmers project page link.

Source:
Source 1 (QuickCheck Chalmers Page), Source 2 (Hackage docs)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Fully verified historical origin claim.

---

## Claim 2

Claim: Example-Based Testing relies on hand-picked inputs chosen by the developer, missing boundary conditions, empty slices, integer overflow, or unexpected input combinations.

Location:
`research/03-evidence.md` (Evidence 2), `research/05-report.md` (Finding 1)

Evidence Provided:
fast-check docs & topic specification.

Source:
Source 13 (fast-check "What is PBT?"), Source 14 (fast-check "Why PBT?")

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Standard definition of example-based testing limitations.

---

## Claim 3

Claim: PBT replaces specific input/output assertions with universal properties (invariants) that must hold for all inputs, running 100+ (or configured 1,000+) test cases per run.

Location:
`research/03-evidence.md` (Evidence 3 & 6), `research/05-report.md` (Finding 1 & 4)

Evidence Provided:
QuickCheck Hackage docs, fast-check docs, Hypothesis docs.

Source:
Source 2 (Hackage), Source 13 (fast-check), Source 18 (Hypothesis)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Verified default count of 100 across major engines (QuickCheck, Hypothesis, fast-check) and configurable nature.

---

## Claim 4

Claim: Test case shrinking automatically reduces large failing inputs (e.g., 500-element arrays or large numbers) to minimal counterexamples (e.g., `[]int{0, -1}` or single boundary values).

Location:
`research/03-evidence.md` (Evidence 4), `research/05-report.md` (Finding 3)

Evidence Provided:
Proptest tutorial (`ValueTree::simplify()`), fast-check shrinking docs, Hypothesis Conjecture byte-stream shrinking.

Source:
Source 6 (How Hypothesis Works), Source 10 (Proptest shrinking), Source 13 (fast-check docs)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Supported by all major framework documentation reviewed.

---

## Claim 5

Claim: Canonical categories of PBT invariants include Roundtrip (Decode(Encode(x)) == x), Idempotence (f(f(x)) == f(x)), Hard-to-Prove/Easy-to-Verify (oracle verification), and Equivalence (comparing new vs naive implementations).

Location:
`research/03-evidence.md` (Evidence 5 & 15), `research/05-report.md` (Finding 2)

Evidence Provided:
Hypothesis encode/decode essay, Hackage QuickCheck docs, fast-check patterns.

Source:
Source 2 (QuickCheck), Source 14 (fast-check), Source 16 (Hypothesis encode/decode)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Widely acknowledged classification of invariant patterns in PBT literature.

---

## Claim 6

Claim: Modern PBT frameworks (like fast-check and Hypothesis) intentionally upweight boundary values, duplicates, and security-relevant inputs (e.g. `__proto__`) rather than using uniform random distribution.

Location:
`research/03-evidence.md` (Evidence 8 & 11), `research/05-report.md` (Finding 5)

Evidence Provided:
fast-check "Why PBT" and Hypothesis "Domain and distribution" docs.

Source:
Source 7 (Hypothesis distribution), Source 14 (fast-check Why PBT)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Directly backed by maintainer documentation detailing non-uniform distribution strategy.

---

## Claim 7

Claim: PBT has detected production bugs and CVEs in high-profile open source projects (Jest, Lodash prototype pollution, React, Underscore, JS-YAML, left-pad, jsonwebtoken, Mercurial, Qutebrowser).

Location:
`research/03-evidence.md` (Evidence 12 & 15), `research/05-report.md` (Finding 6)

Evidence Provided:
fast-check track record page with explicit GitHub issue/PR links & Hypothesis articles.

Source:
Source 15 (fast-check Track Record), Source 16 (Hypothesis Encode/Decode article)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Concrete empirical track record backed by verifiable issue references.

---

## Claim 8

Claim: Go standard library `testing/quick` is officially frozen and lacks modern features like shrinking or custom generator combinators, while third-party `gopter` provides full ScalaCheck/QuickCheck-style PBT capabilities in Go.

Location:
`research/03-evidence.md` (Evidence 10 & 18), `research/05-report.md` (Finding 4 & 9)

Evidence Provided:
pkg.go.dev testing/quick documentation and leanovate/gopter GitHub repository README.

Source:
Source 11 (testing/quick docs), Source 12 (gopter README)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
`testing/quick` package documentation explicitly states "The testing/quick package is frozen and is not accepting new features."

---

## Claim 9

Claim: Hypothesis uses a unique 3-layer architecture (Conjecture byte-stream fuzzer, strategy layer, test runner) performing byte-stream shrinking lexicographically, whereas QuickCheck uses type-based shrinking.

Location:
`research/03-evidence.md` (Evidence 7 & 16), `research/05-report.md` (Finding 7 & 8)

Evidence Provided:
MacIver's "How Hypothesis Works" and QuickCheck Hackage docs.

Source:
Source 2 (QuickCheck), Source 6 (How Hypothesis Works), Source 9 (Proptest intro)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Accurately distinguishes type-based vs byte-stream/value-based shrinking models.

---

## Claim 10

Claim: PBT is a complete replacement for Example-Based Testing and example unit tests are no longer necessary once properties are defined.

Location:
Evaluated as potential overgeneralization in PBT literature.

Evidence Provided:
None — research explicitly refutes this overgeneralization.

Source:
Source 9 (Proptest book), Source 14 (fast-check docs), `research/05-report.md` (Finding 10)

Source Actually Supports Claim:
NO (Research correctly classifies PBT and example testing as complementary).

Classification:
HYPOTHESIS

Severity:
MEDIUM

Notes:
Research report explicitly notes that PBT is complementary to unit testing, preventing overgeneralization.
