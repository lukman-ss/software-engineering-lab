# Evidence

## Evidence 1
Claim: Property-Based Testing defines invariants that must hold true across a universe of pseudo-randomly generated inputs, fundamentally departing from static Example-Based test cases.
Evidence: QuickCheck introduces executable specifications as properties over arbitrary inputs. Rather than testing individual point assertions, test engines execute properties over generated distributions.
Source: Claessen & Hughes (2000), QuickCheck
URL: https://www.cs.tufts.edu/~nr/cs257/archive/john-hughes/quick.pdf
Confidence: HIGH
Corroborated By: Go `testing/quick` documentation, Gopter documentation
Notes: Foundation of property testing across Haskell, Erlang, ScalaCheck, Hypothesis, and Go.

## Evidence 2
Claim: Shrinking reduces counterexamples to their Minimal Failing Example (MFE) to isolate root causes.
Evidence: Modern PBT engines employ shrinking strategies: binary reduction of slices, integer delta-halving towards 0, string prefix/deletion reduction. Gopter implements explicit `shrink` combinators and integrated shrinking to systematically reduce counterexamples.
Source: Gopter Repository & Documentation
URL: https://github.com/leanovate/gopter
Confidence: HIGH
Corroborated By: Claessen & Hughes (2000)
Notes: Without shrinking, debugging a failure with an arbitrary 500-element slice or 64-bit integer noise requires high cognitive overhead.

## Evidence 3
Claim: Go standard library `testing/quick` has significant operational limitations compared to dedicated tools like `gopter`.
Evidence: `testing/quick` provides `Check(f interface{}, config *Config)` and `CheckEqual(f, g interface{}, config *Config)`, but lacks native recursive shrink combinators, integrated shrinking pipelines, and composable custom generators without implementing the reflection-heavy `quick.Generator` interface on specific types.
Source: Go Standard Library Package testing/quick Documentation
URL: https://pkg.go.dev/testing/quick
Confidence: HIGH
Corroborated By: Gopter Documentation
Notes: `testing/quick` is frozen in API evolution; Go ecosystem projects needing advanced generators and shrinkers rely on `gopter` or native coverage-guided fuzzing (`testing.F`).

## Evidence 4
Claim: Fuzzing (`testing.F`) and Property-Based Testing (`gopter`) serve overlapping but distinct roles in software quality.
Evidence: Go native fuzzing uses coverage-guided mutation engines (instrumented binaries measuring branch coverage) with byte-level corpus mutation, while PBT uses type-aware generator combinators targeting semantic invariants with reproducible random seeds.
Source: Go Fuzzing Documentation & Design Draft
URL: https://go.dev/doc/security/fuzz/
Confidence: HIGH
Corroborated By: Claessen & Hughes (2000)
Notes: Fuzzing excels at crash detection and parser boundary exploitation; PBT excels at verifying domain business invariants (e.g. interval arithmetic, accounting balance).
