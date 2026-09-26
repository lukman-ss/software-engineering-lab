# 01 Audit Plan

Target Lab: `labs/15-load-testing`
Audit Target: Research Pipeline (`research/` directory)
Audit Date: 2026-09-26

## Files Reviewed
- `labs/15-load-testing/research/01-plan.md`
- `labs/15-load-testing/research/02-sources.md`
- `labs/15-load-testing/research/03-evidence.md`
- `labs/15-load-testing/research/04-contradictions.md`
- `labs/15-load-testing/research/05-report.md`
- `labs/15-load-testing/research/06-open-questions.md`

## Scope & Pipeline Override
Per PIPELINE OVERRIDE instructions:
- Audit research files only.
- Do NOT audit implementation/code files (`internal/`, `cmd/`, `tests/`).
- Do NOT modify research files.
- Deliver audit artifacts to `labs/15-load-testing/research-audit/`.

## Claims To Verify
1. Standardized four test types (Load, Stress, Spike, Soak) and progression order.
2. Metrics classification (RED method, 4 Golden Signals, percentiles vs averages).
3. Tool characteristics (k6, Locust, JMeter, Gatling).
4. External API handling & mocking vs real dependency testing rules.
5. Bottleneck identification & resource correlation methodologies.
6. Realism of test data and production environment mirroring.

## Primary Risks
1. Verification gaps for tool documentation (Gatling 403, JMeter timeouts).
2. Paywalled standard citation (ISO/IEC 25010).
3. Overgeneralized recommendations (e.g. real external calls vs mocks).
4. Unsubstantiated synthesized claims (e.g. Evidence 29).

## Audit Strategy
1. Source Audit: Inspect all 25 cited sources for reachability, tiering, accuracy, and scope.
2. Claim Audit: Extract major findings/claims from `05-report.md` and `03-evidence.md`; verify evidence support.
3. Contradictions Audit: Verify internal consistency across research documents.
4. Gap Analysis: Identify missing sources, weak evidence, overgeneralizations, and unverified claims.
5. Verdict: Assign final verdict and quality gate results.
