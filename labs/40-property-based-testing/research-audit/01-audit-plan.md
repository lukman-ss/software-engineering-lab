# Audit Plan: Property-Based Testing Research

## Target Lab
`labs/40-property-based-testing`

## Audit Scope
Research documentation audit only (per pipeline override). Implementation and code audit are skipped in this stage.

## Files Reviewed
- `labs/40-property-based-testing/research/01-plan.md`
- `labs/40-property-based-testing/research/02-sources.md`
- `labs/40-property-based-testing/research/03-evidence.md`
- `labs/40-property-based-testing/research/04-contradictions.md`
- `labs/40-property-based-testing/research/05-report.md`
- `labs/40-property-based-testing/research/06-open-questions.md`

## Claims To Verify
1. PBT origin and foundations (Claessen & Hughes ICFP 2000 QuickCheck, Haskell Arbitrary typeclass).
2. Contrast between Example-Based Testing and PBT (human bias vs. universal invariants).
3. Invariant categories (Roundtrip/Encode-Decode, Idempotence, Hard-to-Prove/Easy-to-Verify, Oracle/Equivalence).
4. Shrinking mechanism (counterexample reduction to minimal failing inputs, byte-stream vs value-tree vs type-based).
5. Framework ecosystem behavior (Haskell QuickCheck, Python Hypothesis, Rust proptest, JS fast-check, Go testing/quick & gopter).
6. Intentional distribution bias ("designed for bugs", boundary upweighting, `__proto__` injection).
7. Real-world bug detection evidence (fast-check track record in Jest, Lodash CVEs, React, etc.; Mercurial/Qutebrowser encode-decode bugs; Hypothesis Corpus 2026).
8. Stateful/model-based testing mechanisms (Hypothesis RuleBasedStateMachine, Bundles, rules, invariants).
9. Go-specific tooling status (`testing/quick` frozen status vs `gopter` capabilities).
10. Default test iteration counts across frameworks (100 runs default).

## Code To Execute
*None in this stage (Research-only audit per Pipeline Override).*

## Primary Risks
- Verification of source accuracy and reachability (checking whether cited URLs and materials actually exist and support claims).
- Verification that framework-specific properties (e.g. Hypothesis byte-stream shrinking, fast-check CVE detections) are not generalized inaccurately.
- Assessment of empirical dataset claims (Hypothesis Corpus 2026, fast-check track record).
- Detection of unverified or overgeneralized claims regarding bug-finding superiority without caveats.

## Audit Strategy
1. Audit all 18 cited sources in `02-sources.md` for validity, reachability, tier classification, and alignment with attributed claims.
2. Evaluate all 18 evidence items and 10 report findings against strict quality gates and claim classification.
3. Check internal consistency and verify that reported contradictions or tensions (e.g. PBT vs fuzzing boundary) are accurately analyzed.
4. Record research gaps, weak evidence points, and unverified assumptions.
5. Render a formal evidence-based verdict on the research quality.
