# Contradictions Audit

## Material Contradictions

No material contradictions found.

## Minor Discrepancies and Nuances

### 1. Definitional Boundary Between Fuzzing and Property-Based Testing

- **Statement A:** David R. MacIver (*What is Property-Based Testing?*) defines PBT strictly around user-written invariants/assertions: PBT constructs tests such that when fuzzed, failures reveal system properties not detectable by unguided or assertion-less fuzzing alone.
- **Statement B:** fast-check documentation describes fuzzing more broadly as firing randomly generated values at an algorithm, treating fuzzing as a related concept or loose sub-case of randomized input testing.
- **Type:** SOURCE_CONFLICT (Definitional / Philosophical nuance)
- **Impact:** LOW. Does not impact test execution, generator construction, or invariant design. The research report (`04-contradictions.md`) already noted this distinction accurately.
- **Assessment:** PASS.
