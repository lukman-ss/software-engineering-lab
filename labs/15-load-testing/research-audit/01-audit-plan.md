# Audit Plan — Lab 15: Load Testing Research

## Target Lab
`labs/15-load-testing`

## Files Reviewed
- `labs/15-load-testing/research/01-plan.md`
- `labs/15-load-testing/research/02-sources.md`
- `labs/15-load-testing/research/03-evidence.md`
- `labs/15-load-testing/research/04-contradictions.md`
- `labs/15-load-testing/research/05-report.md`
- `labs/15-load-testing/research/06-open-questions.md`

## Scope & Overrides
- **Research Only**: Auditing research methodology, claim validity, and source integrity.
- **Code Audit**: Excluded per pipeline override (`research-only` stage).

## Claims To Verify
1. Standard test types taxonomy: Average-load, Stress, Spike, Soak/Endurance.
2. Metrics foundation: RED method (Rate, Errors, Duration) and Google SRE Four Golden Signals.
3. Tail latency and percentiles (P50, P95, P99) vs average response times.
4. Pass/fail criteria definition through SLO-driven thresholds.
5. Multi-tool comparative tradeoffs: k6 (JS), Locust (Python/gevent), JMeter (Java/XML multi-protocol), Gatling (Scala/async).
6. Bottleneck identification methodology: layer isolation, backend resource correlation.
7. Requirement for realistic test environments and representative data volumes.
8. Third-party and external dependency impact under load.

## Code To Execute
- **NONE** (Pipeline override: Do not audit implementation/code in this stage).

## Primary Risks
- Inaccessible or unverified sources (paywalled standards like ISO/IEC 25010, vendor 403 blocks like Gatling).
- Overgeneralization of tool strengths or architectural claims without direct benchmark data.
- Synthesized investigation procedures presented without single-source authoritative support.
- Misrepresentation of tool support (e.g., Azure Load Testing framework compatibility).

## Audit Strategy
1. **Source Integrity**: Check URLs, publishers, accessibility status, and tier ratings.
2. **Claim Grounding**: Verify if evidence directly supports claims without extrapolation.
3. **Contradiction Analysis**: Ensure internal consistency across all research documents.
4. **Gap Identification**: Highlight open questions, unverified areas, and empirical limits.
5. **Quality Gating**: Deliver verdict based on factual rigor and methodological soundness.
