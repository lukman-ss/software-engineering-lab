# Contradictions

## No material contradictions discovered among authoritative sources

All primary sources (QuickCheck documentation, Hypothesis documentation, proptest book, fast-check docs, Go testing/quick docs, gopter README) agree on the core definition and feature set of Property-Based Testing:

1. **Definition**: PBT involves specifying properties (invariants) that must hold for all inputs in a domain, with the framework generating random inputs and checking properties. No disagreement on this core concept.
2. **Three pillars**: Random input generation, multiple runs per test, counterexample shrinking — affirmed by both fast-check and Hypothesis.
3. **Complementary to example-based testing**: Both fast-check and Proptest explicitly recommend hybrid usage. No contradiction.
4. **Default test count**: QuickCheck defaults to 100, Hypothesis defaults to 100 (`max_examples`), fast-check defaults to 100. Consistent across frameworks.

## Minor tension: Definition boundary between PBT and fuzzing

SOURCE A (David MacIver, "What is Property Based Testing?", hypothesis.works):
> "Property based testing is the construction of tests such that, when these tests are fuzzed, failures in the test reveal problems with the system under test that could not have been revealed by direct fuzzing of that system." (with a note: "If you feel strongly that fuzzing should count as property-based testing you can just drop the 'that could not have been etc.' part. I'm on the fence about it myself.")

SOURCE B (fast-check documentation):
> "Fuzzing is the idea of firing lots of randomly generated values onto an algorithm to find bugs. In a way, it's not that far from property-based testing and can even be considered a sub-case of it."

SOURCE C (Hypothesis blog, "What is Property Based Testing?"):
> "So with that in mind, let's provide a definition of fuzzing... Property based testing is the construction of tests such that, when these tests are fuzzed..."

ASSESSMENT:
MacIver is the creator of Hypothesis and the author of the field-defining essay; his position is that PBT requires a human-authored property, while fuzzing can be applied blind. fast-check's documentation draws a softer boundary, treating fuzzing as a sub-case of PBT. This is a definitional/philosophical difference, not a factual contradiction. It reflects a spectrum view: some consider any randomized testing with failure detection as PBT, others require a human-specified invariant.

The lab specification itself (topic spec) supports MacIver's narrower view by emphasizing "property" and "invariant" as the defining distinction from raw fuzzing.

This is the only area of disagreement found among sources.
