# Contradictions and Framework Trade-offs

## Contradiction 1: Native Standard Library vs Dedicated Framework Adoption

SOURCE A: Go Standard Library (`testing/quick`)
Claim: Standard library `testing/quick` is sufficient for blackbox quick testing without external dependencies; supports basic generation and equality checks via reflection.

SOURCE B: `gopter` Maintainers & Community
Claim: `testing/quick` is inadequate for non-trivial domain models because it cannot easily compose generators (e.g. valid timestamps, constrained ranges, non-overlapping intervals), lacks stateful model-based testing, and provides primitive shrinking.

ASSESSMENT:
The disagreement stems from design philosophies. `testing/quick` was introduced early in Go's lifecycle as a minimal port of QuickCheck ideas with zero dependencies. For real-world engineering (such as currency parsing with precision constraints or interval merging with ordered bounds), `gopter` or native fuzzing (`testing.F`) provides necessary control over generator bounds and shrink pipelines.

## Contradiction 2: PBT vs Coverage-Guided Fuzzing as Invariant Engine

SOURCE A: Go Fuzzing Docs (`go.dev/doc/security/fuzz`)
Claim: Native fuzzing with `testing.F` should be the default mechanism for finding unexpected edge cases because compiler-instrumented code coverage actively steers inputs into unreached branches.

SOURCE B: Classical PBT (QuickCheck / Hypothesis / Gopter)
Claim: Coverage-guided fuzzing mutates raw bytes or primitive arguments blindly without semantic type preservation, whereas PBT generates high-level domain structures (e.g., ASTs, graphs, composite structs) conforming to strict precondition invariants.

ASSESSMENT:
Both complement each other. Coverage-guided fuzzing is superior for parsing untrusted binary/string formats to prevent crashes (panics, buffer overflows, OOM). PBT is superior for algebraic properties, business invariants, and model equivalence where inputs must satisfy structured domain constraints.
