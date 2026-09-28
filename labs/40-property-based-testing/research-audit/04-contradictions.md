# Contradiction Audit

## Material Contradictions

No material contradictions found across the reviewed research files (`01-plan.md`, `02-sources.md`, `03-evidence.md`, `04-contradictions.md`, `05-report.md`, `06-open-questions.md`).

All primary sources agree on:
1. Core definition and mechanics (invariants, generator, shrinker).
2. Default execution counts (100 iterations default across QuickCheck, Hypothesis, and fast-check).
3. Value-based vs type-based generation tradeoffs.
4. Hybrid/complementary relationship with example-based testing.

---

## Tension 1: Conceptual Boundary between PBT and Fuzzing

### Statement A:
"Property based testing is the construction of tests such that, when these tests are fuzzed, failures in the test reveal problems with the system under test that could not have been revealed by direct fuzzing of that system."
**Location:** David R. MacIver, "What is Property Based Testing?", hypothesis.works (Source 5)

### Statement B:
"Fuzzing is the idea of firing lots of randomly generated values onto an algorithm to find bugs. In a way, it's not that far from property-based testing and can even be considered a sub-case of it."
**Location:** fast-check documentation, "Why Property-Based Testing?" (Source 14)

### Type:
SOURCE_CONFLICT (Definitional nuance)

### Impact:
LOW. Both perspectives agree on the underlying mechanism (automated generation of pseudo-random inputs and automated oracle checking). The difference is purely philosophical: whether PBT is a structured harness around fuzzing or fuzzing is a raw unstructured sub-case of PBT.

### Assessment:
Accurately documented in `research/04-contradictions.md` and appropriately resolved.
