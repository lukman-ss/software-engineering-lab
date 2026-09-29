# Audit Plan

## Target Lab
`labs/40-property-based-testing`

## Audit Date
2026-09-29

## Mode / Override
PIPELINE OVERRIDE: Research Audit only. Implementation and code files are excluded from this stage.

## Files Reviewed
- `labs/40-property-based-testing/research/01-plan.md`
- `labs/40-property-based-testing/research/02-sources.md`
- `labs/40-property-based-testing/research/03-evidence.md`
- `labs/40-property-based-testing/research/04-contradictions.md`
- `labs/40-property-based-testing/research/05-report.md`
- `labs/40-property-based-testing/research/06-open-questions.md`

## Claims To Verify
1. PBT origin and authorship (QuickCheck, Claessen & Hughes, ICFP 2000).
2. Fundamental difference between Example-Based Testing and PBT (invariants vs hand-picked cases, random input generation).
3. Automatic counterexample shrinking mechanism and minimization of failing inputs.
4. Canonical categories of invariants (Roundtrip, Idempotence, Hard-to-prove/easy-to-verify, Equivalence/Oracle).
5. Framework defaults (QuickCheck, Hypothesis, fast-check default of 100 test runs).
6. Hypothesis architecture (Conjecture byte-stream fuzzer, strategy layer, value-based shrinking vs type-based).
7. Domain vs. distribution separation in modern PBT frameworks (bias toward edge cases, security strings like `__proto__`).
8. Stateful/state-machine testing capabilities (RuleBasedStateMachine, bundles, model vs system verification).
9. Go ecosystem specifics (Go standard library `testing/quick` being frozen/reflection-based; third-party `gopter` feature set).
10. Empirical evidence of PBT efficacy (real-world bug discoveries and CVE detections in Jest, lodash, react, numpy, left-pad, etc.).
11. Complementary nature of PBT and Example-Based Testing (hybrid methodology recommendation).

## Primary Risks
1. Verification of external URLs and academic/primary source authenticity.
2. Distortion or overgeneralization of framework capabilities across languages.
3. Conflating fuzzing with property-based testing without nuance.
4. Unsupported assertions regarding empirical bug detection rates or maintenance costs.

## Audit Strategy
- Verify each cited source URL, publisher, and claimed relevance.
- Map claims in `05-report.md` and `03-evidence.md` against primary documentation.
- Check consistency across research files for contradictions or unacknowledged gaps.
- Classify all claims and record gaps according to severity guidelines.
