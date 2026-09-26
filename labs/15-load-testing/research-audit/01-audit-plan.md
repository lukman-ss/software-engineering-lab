# Audit Plan

## Target Lab
`labs/15-load-testing`

## Files Reviewed
- `research/01-plan.md`
- `research/02-sources.md`
- `research/03-evidence.md`
- `research/04-contradictions.md` (to be reviewed)
- `research/05-report.md` (to be reviewed)
- `research/06-open-questions.md` (to be reviewed)

## Claims To Verify
- Four primary types of load testing (Load, Stress, Spike, Soak).
- Tools comparison (k6, Locust, JMeter, Gatling).
- Core metrics (P50/P95/P99, RPS, Error Rate).
- RED method and Four Golden Signals.
- The use of synthetic monitoring.
- Handling external APIs during load tests.
- Common load testing mistakes (testing only `/health`, small dummy data).

## Code To Execute
*Pipeline override active: Do not execute code.*

## Primary Risks
- Over-reliance on vendor documentation (k6, Azure, Locust) without considering potential bias.
- Unverified sources (JMeter, ISO/IEC 25010) being used to support major claims.
- Synthesized claims (e.g., Evidence 29) being presented as hard facts.
- Conflation of general software performance testing with specific API load testing.

## Audit Strategy
1. Cross-reference claims in the final report against the provided evidence.
2. Verify that sources marked as `NOT VERIFIED` or `DISCLAIMED` are not used to support critical claims.
3. Check for over-generalization (e.g., assuming all applications need k6).
4. Identify any contradictions between the plan, sources, and final report.
