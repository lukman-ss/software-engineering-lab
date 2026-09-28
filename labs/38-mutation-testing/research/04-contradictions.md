# Contradictions Analysis

## Summary

No material contradictions were discovered between the primary sources consulted. The sources are largely consistent in their description of mutation testing fundamentals, with minor differences in emphasis and scope.

## Potential Tension Points

### Point 1: "Mutation testing requires a test to already exist" vs. "LLMs can generate tests from mutants"

SOURCE A (Meta Engineering Blog):
"Even though mutation testing cannot exist on its own (it requires a test to already exist), it helps engineers and developers identify weak assertions..."

ASSESSMENT:
This is not a contradiction but a clarification of scope. Meta's ACH system uses the *output* of mutation testing (unkilled mutants) to *generate* new tests via LLMs. The system still relies on existing tests as an oracle to detect whether the LLM-generated tests are effective. The direction is: existing tests → mutation → identify gaps → LLM generates new tests to fill gaps. This extends rather than contradicts the traditional workflow.

---

### Point 2: Equivalent mutants as "mathematically undecidable" vs. tools claiming equivalence detection

SOURCE A (Wikipedia, citing [18]):
"The effort needed to check if mutants are equivalent or not can be very high, even for small programs."

SOURCE B (Meta Engineering Blog):
"In our own research and testing with ACH we found that when combined with simple static analysis preprocessing (e.g., stripping comments), this approach achieves high precision (0.79) and recall (0.47) – rising to 0.95 and 0.96 with simple preprocessing – in detecting equivalent mutants."

ASSESSMENT:
These are not contradictory but represent different approaches. The undecidability refers to the general theoretical problem (proving equivalence for arbitrary programs). Practical heuristics (static analysis, LLM-based classification) can achieve high precision/recall for common cases, which is what Meta reports. PIT also avoids mutating certain code patterns (static initializers, enum constructors) as a heuristic approach. The theoretical limit remains, but practical approximations work well enough for industrial use.

---

### Point 3: Computational cost concern vs. tool claims of speed

SOURCE A (PIT FAQ):
"Mutation testing is a computationally expensive process and can take quite some time depending on the size of your codebase and the quality and speed of your test suite."

SOURCE B (Stryker docs):
"Stryker is fast to run and easy to use."

ASSESSMENT:
These are not contradictory but reflect relative improvements over older systems (pre-2010). PIT explicitly states: "PIT is fast compared to other mutation testing systems, but that can still mean that things will take a while." The key insight is that "fast" is relative to earlier tools that took days; modern tools like PIT and Stryker have reduced runtime to minutes for moderate-sized projects. For large codebases, both tools acknowledge that selective targeting (only changed files, incremental mutation) is necessary.

---

## Areas of Agreement (All Sources)

1. Mutation testing fundamentally works by introducing small, controlled changes to source code.
2. The mutation score formula is: (Killed Mutants / Total Mutants) × 100%.
3. Code coverage ≠ test quality; mutation testing exposes the gap.
4. Equivalent mutants remain a practical challenge.
5. The competent programmer hypothesis is foundational.
6. Tools exist for Java (PIT, MuJava) and JS/TS (Stryker), with growing language support.
