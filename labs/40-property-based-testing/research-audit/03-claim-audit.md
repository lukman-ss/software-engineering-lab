# Claim Audit

## Claim 1

Claim:
Property-Based Testing originated with QuickCheck, created by Koen Claessen and John Hughes at Chalmers University of Technology and introduced at ICFP 2000.

Location:
`research/05-report.md: Finding 1`, `research/03-evidence.md: Evidence 1`

Evidence Provided:
Project documentation and paper references at Chalmers website and Hackage documentation.

Source:
Source 1 (Chalmers QuickCheck project page) & Source 2 (Hackage)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Accurately reflects the historical origin of the paradigm.

---

## Claim 2

Claim:
Example-Based Testing relies on hand-crafted input/output pairs that miss unexpected edge cases (e.g., empty slices, duplicate values, extreme bounds), whereas PBT tests universal invariants over thousands of generated inputs.

Location:
`research/05-report.md: Finding 1`, `research/03-evidence.md: Evidence 2 & 3`

Evidence Provided:
Citations from fast-check introduction, Hackage QuickCheck docs, and topic specification.

Source:
Source 13 (fast-check documentation) & Source 2 (Hackage QuickCheck)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Well-established conceptual definition across testing literature.

---

## Claim 3

Claim:
The canonical invariant patterns for PBT are: (1) Roundtrip (Encode/Decode), (2) Idempotence, (3) Hard to Prove / Easy to Verify, and (4) Equivalence / Oracle.

Location:
`research/05-report.md: Finding 2`, `research/03-evidence.md: Evidence 5`

Evidence Provided:
Hypothesis encode/decode article, Hackage QuickCheck manual, and topic specification.

Source:
Source 16 (Hypothesis blog), Source 2 (Hackage QuickCheck), Source 17 (Hypothesis Corpus)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
These four categories represent standard formulations widely referenced by Hughes, MacIver, and Dubien.

---

## Claim 4

Claim:
Test case shrinking automatically reduces large, complex failing inputs to the minimal reproducing counterexample (e.g., from a 500-element array down to a minimal failing slice).

Location:
`research/05-report.md: Finding 3`, `research/03-evidence.md: Evidence 4`

Evidence Provided:
Proptest shrinking tutorial demonstrating ValueTree binary search, fast-check counterexample minimization docs, and Hypothesis Conjecture byte-stream reduction.

Source:
Source 10 (Proptest book), Source 13 (fast-check), Source 6 (Hypothesis)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Core mechanism shared across major PBT libraries.

---

## Claim 5

Claim:
QuickCheck, Hypothesis, and fast-check all use a default of 100 test iterations per property check.

Location:
`research/05-report.md: Areas of Agreement`, `research/03-evidence.md: Evidence 6 & 13`

Evidence Provided:
QuickCheck Args manual (`maxSuccess` default 100), Hypothesis `max_examples` configuration, fast-check default sampling.

Source:
Source 2 (QuickCheck Hackage), Source 13 (fast-check), Source 18 (Hypothesis)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Consistent default value documented across all three frameworks.

---

## Claim 6

Claim:
Hypothesis implements a byte-stream architecture (Conjecture) with value-based shrinking, distinct from Haskell QuickCheck's type-based shrinking (`shrink :: a -> [a]`), avoiding local minima when shrinking complex interrelated values.

Location:
`research/05-report.md: Finding 8`, `research/03-evidence.md: Evidence 7 & 16`

Evidence Provided:
Articles "How Hypothesis Works" and proptest book architectural comparisons.

Source:
Source 6 (hypothesis.works) & Source 9 (Proptest book)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
MEDIUM

Notes:
The technical description of Conjecture's byte stream and lexicographical minimization matches the creator's specification.

---

## Claim 7

Claim:
Modern PBT frameworks deliberately bias input distributions toward boundary values, duplicates, and security vulnerabilities (such as `__proto__`) rather than uniform random sampling.

Location:
`research/05-report.md: Finding 5`, `research/03-evidence.md: Evidence 8 & 11`

Evidence Provided:
Hypothesis "Domain and distribution" documentation and fast-check "Why Property-Based Testing?".

Source:
Source 7 (Hypothesis) & Source 14 (fast-check)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Documents intentional distribution biasing implemented in production PBT libraries.

---

## Claim 8

Claim:
Stateful/state-machine testing generates sequences of operations to discover state-dependent bugs by asserting invariants against a simplified model (e.g. `RuleBasedStateMachine` in Hypothesis).

Location:
`research/05-report.md: Finding 7`, `research/03-evidence.md: Evidence 9`

Evidence Provided:
Hypothesis stateful documentation on rules, bundles, invariants, and preconditions.

Source:
Source 8 (Hypothesis documentation) & Source 12 (gopter README)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Clearly supported by official documentation and gopter feature list.

---

## Claim 9

Claim:
Go's standard library `testing/quick` package is frozen and reflection-based with no built-in shrinker, whereas third-party `gopter` provides tighter generator control, shrinkers, regex generators, and stateful testing.

Location:
`research/05-report.md: Finding 9`, `research/03-evidence.md: Evidence 10 & 18`

Evidence Provided:
Go standard library package note ("frozen and is not accepting new features") and gopter feature comparison.

Source:
Source 11 (pkg.go.dev/testing/quick) & Source 12 (leanovate/gopter)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Accurately captures Go's PBT landscape.

---

## Claim 10

Claim:
Empirical evidence shows PBT has uncovered real vulnerabilities and critical bugs in widely used production open-source libraries, including prototype pollution CVEs in lodash, comparison bugs in Jest (`toStrictEqual`), and unicode handling in left-pad.

Location:
`research/05-report.md: Finding 6`, `research/03-evidence.md: Evidence 12`

Evidence Provided:
fast-check track record with issue URLs (`jestjs/jest#7941`, `left-pad/left-pad#58`, `auth0/node-jsonwebtoken#945`).

Source:
Source 15 (fast-check Track Record)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Verified against public issue tracker references provided in fast-check documentation.

---

## Claim 11

Claim:
PBT is designed to be complementary to example-based unit testing rather than an outright replacement.

Location:
`research/05-report.md: Finding 10`, `research/03-evidence.md: Evidence 14`

Evidence Provided:
Explicit statements in fast-check and Proptest documentation recommending hybrid testing strategies.

Source:
Source 14 (fast-check) & Source 9 (Proptest book)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Industry consensus across major framework maintainers.
